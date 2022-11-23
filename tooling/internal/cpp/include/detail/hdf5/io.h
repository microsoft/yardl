// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <array>
#include <stdexcept>
#include <string>
#include <utility>
#include <variant>
#include <vector>

#include <H5Cpp.h>

#include "ddl.h"

namespace yardl::hdf5 {

template <typename TInner, typename TOuter>
static inline void WriteScalarDataset(H5::Group const& group, std::string name,
                                      H5::DataType const& hdf5_type, TOuter const& value) {
  H5::DataSpace dataspace;
  auto dataset = group.createDataSet(name, hdf5_type, dataspace);

  if constexpr (std::is_same_v<TInner, TOuter>) {
    dataset.write(&value, hdf5_type);
  } else {
    TInner inner_value(value);
    dataset.write(&inner_value, hdf5_type);
  }
}

template <typename TInner, typename TOuter>
static inline void ReadScalarDataset(H5::Group const& group, std::string name,
                                     H5::DataType const& hdf5_type, TOuter& value) {
  auto dataset = group.openDataSet(name);
  H5::DataSpace dataspace;
  if constexpr (std::is_same_v<TInner, TOuter>) {
    dataset.read(&value, hdf5_type, dataspace, dataspace);
  } else {
    TInner inner_value;
    dataset.read(&inner_value, hdf5_type, dataspace, dataspace);
    value = static_cast<TOuter>(inner_value);
  }
}

class TransferPropertiesContainer {
 public:
  TransferPropertiesContainer(size_t buffer_size)
      : buffer_(buffer_size), background_buffer_(buffer_size) {
    properties_.setBuffer(buffer_size, buffer_.data(), background_buffer_.data());
  }

  TransferPropertiesContainer(TransferPropertiesContainer const&) = delete;

  H5::DSetMemXferPropList const& properties() const {
    return properties_;
  }

 private:
  H5::DSetMemXferPropList properties_{};
  std::vector<uint8_t> buffer_;
  std::vector<uint8_t> background_buffer_;
};

class DatasetWriter {
 public:
  DatasetWriter(H5::Group const& group, std::string name,
                H5::DataType const& element_type, size_t conversion_buffer_size)
      : dataset_(CreateDataset(group, name, element_type)), element_type_(element_type) {
    hsize_t rows_at_a_time = 1;
    mem_space_ = H5::DataSpace(1, &rows_at_a_time, NULL);
    if (conversion_buffer_size > 0) {
      transfer_properties_.emplace(conversion_buffer_size);
    }
  }

  template <typename TInner, typename TOuter>
  hsize_t Append(TOuter const& value) {
    hsize_t new_size = offset_ + 1;
    dataset_.extend(&new_size);

    H5::DataSpace fileSpace = dataset_.getSpace();
    fileSpace.selectElements(H5S_SELECT_SET, 1, &offset_);

    if constexpr (std::is_same_v<TInner, TOuter>) {
      dataset_.write(&value, element_type_, mem_space_, fileSpace);
    } else {
      TInner inner_value(value);
      dataset_.write(&inner_value, element_type_, mem_space_, fileSpace, transfer_properties());
    }

    return offset_++;
  }

  template <typename TInner, typename TOuter>
  hsize_t AppendBatch(std::vector<TOuter> const& value) {
    hsize_t new_size = offset_ + value.size();
    dataset_.extend(&new_size);
    H5::DataSpace fileSpace = dataset_.getSpace();

    hsize_t new_row_count = value.size();
    fileSpace.selectHyperslab(H5S_SELECT_SET, &new_row_count, &offset_);
    H5::DataSpace mem_space(1, &new_row_count, NULL);

    if constexpr (std::is_same_v<TInner, TOuter>) {
      dataset_.write(value.data(), element_type_, mem_space, fileSpace);
    } else {
      InnerTypeBuffer<TInner, TOuter> inner_value(value);
      dataset_.write(inner_value.data(), element_type_, mem_space, fileSpace, transfer_properties());
    }

    auto initial_offset = offset_;
    offset_ += value.size();
    return initial_offset;
  }

 private:
  static H5::DataSet
  CreateDataset(H5::Group const& group, std::string name,
                H5::DataType const& element_type) {
    hsize_t dims = 0;
    hsize_t maxdims = H5S_UNLIMITED;

    // Because we do not know the dataset size up-front, we need to enable chunking.
    // Based on some testing, 8MB chunk size seems to have good throughput,
    // but we could make this a tunable parameter in the future.
    hsize_t chunk_dims = std::max(static_cast<hsize_t>(1),
                                  static_cast<hsize_t>(8 * 1024 * 1024 / element_type.getSize()));

    H5::DataSpace dataspace(1, &dims, &maxdims);

    H5::DSetCreatPropList prop;
    prop.setChunk(1, &chunk_dims);
    return group.createDataSet(name, element_type, dataspace, prop);
  }

  H5::DSetMemXferPropList const& transfer_properties() const {
    return transfer_properties_.has_value()
               ? transfer_properties_.value().properties()
               : H5::DSetMemXferPropList::DEFAULT;
  }

  H5::DataSet dataset_;
  H5::DataType element_type_;
  hsize_t offset_{};
  H5::DataSpace mem_space_;
  std::optional<TransferPropertiesContainer> transfer_properties_;
};

class DatasetReader {
 public:
  DatasetReader(H5::Group const& group, std::string name,
                H5::DataType const& element_type, size_t conversion_buffer_size = 0)
      : dataset_(group.openDataSet(name)),
        element_type_(element_type),
        filespace_(dataset_.getSpace()) {
    filespace_.getSimpleExtentDims(&total_rows_);
    if (conversion_buffer_size > 0) {
      transfer_properties_.emplace(std::max(static_cast<size_t>(conversion_buffer_size),
                                            dataset_.getDataType().getSize()));
    }
  }

  template <typename TInner, typename TOuter>
  bool Read(TOuter& value) {
    if (offset_ >= total_rows_) {
      return false;
    }

    hsize_t rows_at_a_time = 1;
    H5::DataSpace memspace(1, &rows_at_a_time, NULL);
    filespace_.selectElements(H5S_SELECT_SET, 1, &offset_);

    if constexpr (std::is_same_v<TInner, TOuter>) {
      dataset_.read(&value, element_type_, memspace, filespace_);
    } else {
      TInner inner_value;
      dataset_.read(&inner_value, element_type_, memspace, filespace_, transfer_properties());
      value = static_cast<TOuter>(inner_value);
    }

    offset_++;
    return true;
  }

  template <typename TInner, typename TOuter>
  bool ReadBatch(std::vector<TOuter>& values) {
    if (offset_ >= total_rows_) {
      return false;
    }

    hsize_t rows_to_read = std::min(total_rows_ - offset_, static_cast<hsize_t>(values.capacity()));
    if (rows_to_read != values.size()) {
      values.resize(rows_to_read);
    }

    H5::DataSpace memspace(1, &rows_to_read, NULL);
    filespace_.selectHyperslab(H5S_SELECT_SET, &rows_to_read, &offset_);

    if constexpr (std::is_same_v<TInner, TOuter>) {
      dataset_.read(values.data(), element_type_, memspace, filespace_);
    } else {
      InnerTypeBuffer<TInner, TOuter> inner_value(values.size());
      dataset_.read(inner_value.data(), element_type_, memspace, filespace_, transfer_properties());
      std::transform(inner_value.begin(), inner_value.end(), values.begin(),
                     [](TInner const& v) { return static_cast<TOuter>(v); });
    }

    offset_ += rows_to_read;
    return offset_ < total_rows_;
  }

 private:
  H5::DSetMemXferPropList const& transfer_properties() const {
    return transfer_properties_.has_value()
               ? transfer_properties_.value().properties()
               : H5::DSetMemXferPropList::DEFAULT;
  }

  H5::DataSet dataset_;
  H5::DataType element_type_;
  H5::DataSpace filespace_;
  hsize_t offset_{};
  hsize_t total_rows_;
  std::optional<TransferPropertiesContainer> transfer_properties_;
};

template <size_t N>
class UnionDatasetWriter {
 public:
  template <typename... TypeLabelBufferSizeTuples>
  UnionDatasetWriter(H5::Group const& parent_group, std::string name,
                     bool nullable, TypeLabelBufferSizeTuples const&... type_label_buffer_size_tuples)
      : group_(CreateGroup(parent_group, name)),
        index_writer_(
            group_, "$index",
            UnionIndexDatasetElementTypeDdl(
                UnionTypeEnumDdl(nullable, std::get<1>(type_label_buffer_size_tuples)...)),
            0),
        writers_{DatasetWriter{group_,
                               std::get<1>(type_label_buffer_size_tuples),
                               std::get<0>(type_label_buffer_size_tuples),
                               std::get<2>(type_label_buffer_size_tuples)}...} {
    index_entry_buffer_.reserve(8192);
  }

  ~UnionDatasetWriter() {
    Flush();
  }

  void Flush() {
    if (!index_entry_buffer_.empty()) {
      index_writer_.AppendBatch<IndexEntry, IndexEntry>(index_entry_buffer_);
      index_entry_buffer_.clear();
    }
  }

  template <typename TElementInner, typename TElementOuter>
  void Append(int8_t type, TElementOuter const& value) {
    uint64_t offset;
    if (type == -1) {
      offset = 0;
    } else {
      DatasetWriter& ds = writers_[type];
      offset = static_cast<uint64_t>(ds.Append<TElementInner, TElementOuter>(value));
    }

    index_entry_buffer_.push_back({type, offset});
    if (index_entry_buffer_.size() == index_entry_buffer_.capacity()) {
      Flush();
    }
  }

 private:
  static H5::Group CreateGroup(H5::Group const& parent_group, std::string group_name) {
    if (parent_group.nameExists(group_name)) {
      throw std::runtime_error("Unable to create group '" + group_name +
                               "' for protocol because it already exists.");
    }
    return parent_group.createGroup(group_name);
  }

  H5::Group group_;
  DatasetWriter index_writer_;
  std::vector<IndexEntry> index_entry_buffer_;
  std::array<DatasetWriter, N> writers_;
};

template <size_t N>
class UnionDatasetReader {
 public:
  template <typename... TypeLabelBufferSizeTuples>
  UnionDatasetReader(H5::Group const& parent_group, std::string name,
                     bool nullable, TypeLabelBufferSizeTuples const&... type_label_buffer_size_tuples)
      : group_(parent_group.openGroup(name)),
        index_reader_(group_, "$index",
                      UnionIndexDatasetElementTypeDdl(
                          UnionTypeEnumDdl(nullable, std::get<1>(type_label_buffer_size_tuples)...))),
        readers_{DatasetReader{group_,
                               std::get<1>(type_label_buffer_size_tuples),
                               std::get<0>(type_label_buffer_size_tuples),
                               std::get<2>(type_label_buffer_size_tuples)}...} {
    index_entry_buffer_.reserve(8192);
  }

  std::tuple<bool, int8_t, DatasetReader*> ReadIndex() {
    if (index_entry_buffer_offset_ == index_entry_buffer_.size()) {
      if (!index_reader_has_more_) {
        return {false, 0, nullptr};
      }

      index_reader_has_more_ = index_reader_.ReadBatch<IndexEntry, IndexEntry>(index_entry_buffer_);
      index_entry_buffer_offset_ = 0;

      if (index_entry_buffer_.empty()) {
        return {false, 0, nullptr};
      }
    }

    IndexEntry const& entry = index_entry_buffer_[index_entry_buffer_offset_++];

    if (entry.type_ < 0) {
      return {true, entry.type_, nullptr};
    }

    return {true, entry.type_, &readers_[entry.type_]};
  }

  H5::Group group_;
  DatasetReader index_reader_;
  std::vector<IndexEntry> index_entry_buffer_;
  size_t index_entry_buffer_offset_{};
  bool index_reader_has_more_{true};
  std::array<DatasetReader, N> readers_;
};

class Hdf5Writer {
 public:
  Hdf5Writer(std::string const& path, std::string const& group_name, std::string const& schema)
      : file_(path, H5F_ACC_CREAT | H5F_ACC_RDWR),
        group_(CreateGroup(file_, group_name)) {
    WriteScalarDataset<InnerVlenString, std::string>(group_, "$schema", InnerVlenStringDdl(), schema);
  }

 private:
  static H5::Group CreateGroup(H5::H5File const& file, std::string group_name) {
    if (file.nameExists(group_name)) {
      throw std::runtime_error("Unable to create group '" + group_name +
                               "' for protocol because it already exists.");
    }
    return file.createGroup(group_name);
  }

 protected:
  H5::H5File file_;
  H5::Group group_;
};

class Hdf5Reader {
 public:
  Hdf5Reader(std::string path, std::string group_name, std::string const& expected_schema)
      : file_(path, H5F_ACC_RDONLY), group_(OpenGroup(file_, group_name)) {
    std::string actual_schema;
    ReadScalarDataset<InnerVlenString, std::string>(group_, "$schema", InnerVlenStringDdl(), actual_schema);
    if (actual_schema != expected_schema) {
      throw std::runtime_error("Data to be read is not compatible with the protocol");
    }
  }

 private:
  static H5::Group OpenGroup(H5::H5File const& file, std::string group_name) {
    if (!file.nameExists(group_name)) {
      throw std::runtime_error("Unable to open group '" + group_name +
                               "' for protocol because it does not exist.");
    }

    return file.openGroup(group_name);
  }

 protected:
  H5::H5File file_;
  H5::Group group_;
};

template <class>
inline constexpr bool always_false_v = false;

}  // namespace yardl::hdf5
