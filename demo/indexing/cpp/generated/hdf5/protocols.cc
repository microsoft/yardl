// This file was generated by the "yardl" tool. DO NOT EDIT.

#include "protocols.h"

#include "../yardl/detail/hdf5/io.h"
#include "../yardl/detail/hdf5/ddl.h"
#include "../yardl/detail/hdf5/inner_types.h"

namespace sketch::hdf5 {
namespace {
struct _Inner_Header {
  _Inner_Header() {} 
  _Inner_Header(sketch::Header const& o) 
      : subject(o.subject) {
  }

  void ToOuter (sketch::Header& o) const {
    yardl::hdf5::ToOuter(subject, o.subject);
  }

  yardl::hdf5::InnerVlenString subject;
};

struct _Inner_Sample {
  _Inner_Sample() {} 
  _Inner_Sample(sketch::Sample const& o) 
      : id(o.id),
      data(o.data) {
  }

  void ToOuter (sketch::Sample& o) const {
    yardl::hdf5::ToOuter(id, o.id);
    yardl::hdf5::ToOuter(data, o.data);
  }

  uint32_t id;
  yardl::hdf5::InnerVlen<int32_t, int32_t> data;
};

[[maybe_unused]] H5::CompType GetHeaderHdf5Ddl() {
  using RecordType = sketch::hdf5::_Inner_Header;
  H5::CompType t(sizeof(RecordType));
  t.insertMember("subject", HOFFSET(RecordType, subject), yardl::hdf5::InnerVlenStringDdl());
  return t;
}

[[maybe_unused]] H5::CompType GetSampleHdf5Ddl() {
  using RecordType = sketch::hdf5::_Inner_Sample;
  H5::CompType t(sizeof(RecordType));
  t.insertMember("id", HOFFSET(RecordType, id), H5::PredType::NATIVE_UINT32);
  t.insertMember("data", HOFFSET(RecordType, data), yardl::hdf5::InnerVlenDdl(H5::PredType::NATIVE_INT32));
  return t;
}

} // namespace 

MyProtocolWriter::MyProtocolWriter(std::string path)
    : yardl::hdf5::Hdf5Writer::Hdf5Writer(path, "MyProtocol", schema_) {
}

void MyProtocolWriter::WriteHeaderImpl(sketch::Header const& value) {
  yardl::hdf5::WriteScalarDataset<sketch::hdf5::_Inner_Header, sketch::Header>(group_, "header", sketch::hdf5::GetHeaderHdf5Ddl(), value);
}

void MyProtocolWriter::WriteSamplesImpl(sketch::Sample const& value) {
  if (!samples_dataset_state_) {
    samples_dataset_state_ = std::make_unique<yardl::hdf5::DatasetWriter>(group_, "samples", sketch::hdf5::GetSampleHdf5Ddl(), std::max(sizeof(sketch::hdf5::_Inner_Sample), sizeof(sketch::Sample)));
  }

  samples_dataset_state_->Append<sketch::hdf5::_Inner_Sample, sketch::Sample>(value);
}

void MyProtocolWriter::WriteSamplesImpl(std::vector<sketch::Sample> const& values) {
  if (!samples_dataset_state_) {
    samples_dataset_state_ = std::make_unique<yardl::hdf5::DatasetWriter>(group_, "samples", sketch::hdf5::GetSampleHdf5Ddl(), std::max(sizeof(sketch::hdf5::_Inner_Sample), sizeof(sketch::Sample)));
  }

  samples_dataset_state_->AppendBatch<sketch::hdf5::_Inner_Sample, sketch::Sample>(values);
}

void MyProtocolWriter::EndSamplesImpl() {
  if (!samples_dataset_state_) {
    samples_dataset_state_ = std::make_unique<yardl::hdf5::DatasetWriter>(group_, "samples", sketch::hdf5::GetSampleHdf5Ddl(), std::max(sizeof(sketch::hdf5::_Inner_Sample), sizeof(sketch::Sample)));
  }

  samples_dataset_state_.reset();
}

MyProtocolReader::MyProtocolReader(std::string path)
    : yardl::hdf5::Hdf5Reader::Hdf5Reader(path, "MyProtocol", schema_) {
}

void MyProtocolReader::ReadHeaderImpl(sketch::Header& value) {
  yardl::hdf5::ReadScalarDataset<sketch::hdf5::_Inner_Header, sketch::Header>(group_, "header", sketch::hdf5::GetHeaderHdf5Ddl(), value);
}

bool MyProtocolReader::ReadSamplesImpl(sketch::Sample& value) {
  if (!samples_dataset_state_) {
    samples_dataset_state_ = std::make_unique<yardl::hdf5::DatasetReader>(group_, "samples", sketch::hdf5::GetSampleHdf5Ddl(), std::max(sizeof(sketch::hdf5::_Inner_Sample), sizeof(sketch::Sample)));
  }

  bool has_value = samples_dataset_state_->Read<sketch::hdf5::_Inner_Sample, sketch::Sample>(value);
  if (!has_value) {
    samples_dataset_state_.reset();
  }

  return has_value;
}

bool MyProtocolReader::ReadSamplesImpl(std::vector<sketch::Sample>& values) {
  if (!samples_dataset_state_) {
    samples_dataset_state_ = std::make_unique<yardl::hdf5::DatasetReader>(group_, "samples", sketch::hdf5::GetSampleHdf5Ddl());
  }

  bool has_more = samples_dataset_state_->ReadBatch<sketch::hdf5::_Inner_Sample, sketch::Sample>(values);
  if (!has_more) {
    samples_dataset_state_.reset();
  }

  return has_more;
}

} // namespace sketch::hdf5
