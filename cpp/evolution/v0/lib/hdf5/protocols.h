// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include <array>
#include <complex>
#include <optional>
#include <variant>
#include <vector>

#include "../protocols.h"
#include "../types.h"
#include "../yardl/detail/hdf5/io.h"

namespace evo_test::hdf5 {
// HDF5 writer for the ProtocolWithChanges protocol.
class ProtocolWithChangesWriter : public evo_test::ProtocolWithChangesWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  ProtocolWithChangesWriter(std::string path);

  protected:
  void WriteInt8ToIntImpl(int8_t const& value) override;

  void WriteInt8ToLongImpl(int8_t const& value) override;

  void WriteInt8ToUintImpl(int8_t const& value) override;

  void WriteInt8ToUlongImpl(int8_t const& value) override;

  void WriteInt8ToFloatImpl(int8_t const& value) override;

  void WriteInt8ToDoubleImpl(int8_t const& value) override;

  void WriteIntToUintImpl(int32_t const& value) override;

  void WriteIntToLongImpl(int32_t const& value) override;

  void WriteIntToFloatImpl(int32_t const& value) override;

  void WriteIntToDoubleImpl(int32_t const& value) override;

  void WriteUintToUlongImpl(uint32_t const& value) override;

  void WriteUintToFloatImpl(uint32_t const& value) override;

  void WriteUintToDoubleImpl(uint32_t const& value) override;

  void WriteFloatToDoubleImpl(float const& value) override;

  void WriteIntToStringImpl(int32_t const& value) override;

  void WriteUintToStringImpl(uint32_t const& value) override;

  void WriteLongToStringImpl(int64_t const& value) override;

  void WriteUlongToStringImpl(uint64_t const& value) override;

  void WriteFloatToStringImpl(float const& value) override;

  void WriteDoubleToStringImpl(double const& value) override;

  void WriteIntToOptionalImpl(int32_t const& value) override;

  void WriteFloatToOptionalImpl(float const& value) override;

  void WriteStringToOptionalImpl(std::string const& value) override;

  void WriteIntToUnionImpl(int32_t const& value) override;

  void WriteFloatToUnionImpl(float const& value) override;

  void WriteStringToUnionImpl(std::string const& value) override;

  void WriteOptionalIntToFloatImpl(std::optional<int32_t> const& value) override;

  void WriteOptionalFloatToStringImpl(std::optional<float> const& value) override;

  void WriteAliasedLongToStringImpl(evo_test::AliasedLongToString const& value) override;

  void WriteRecordWithChangesImpl(evo_test::RecordWithChanges const& value) override;

  void WriteAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges const& value) override;

  void WriteOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges> const& value) override;

  void WriteAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges> const& value) override;

  void WriteStreamedRecordWithChangesImpl(evo_test::RecordWithChanges const& value) override;

  void WriteStreamedRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges> const& values) override;

  void EndStreamedRecordWithChangesImpl() override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetWriter> streamedRecordWithChanges_dataset_state_;
};

// HDF5 reader for the ProtocolWithChanges protocol.
class ProtocolWithChangesReader : public evo_test::ProtocolWithChangesReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  ProtocolWithChangesReader(std::string path);

  void ReadInt8ToIntImpl(int8_t& value) override;

  void ReadInt8ToLongImpl(int8_t& value) override;

  void ReadInt8ToUintImpl(int8_t& value) override;

  void ReadInt8ToUlongImpl(int8_t& value) override;

  void ReadInt8ToFloatImpl(int8_t& value) override;

  void ReadInt8ToDoubleImpl(int8_t& value) override;

  void ReadIntToUintImpl(int32_t& value) override;

  void ReadIntToLongImpl(int32_t& value) override;

  void ReadIntToFloatImpl(int32_t& value) override;

  void ReadIntToDoubleImpl(int32_t& value) override;

  void ReadUintToUlongImpl(uint32_t& value) override;

  void ReadUintToFloatImpl(uint32_t& value) override;

  void ReadUintToDoubleImpl(uint32_t& value) override;

  void ReadFloatToDoubleImpl(float& value) override;

  void ReadIntToStringImpl(int32_t& value) override;

  void ReadUintToStringImpl(uint32_t& value) override;

  void ReadLongToStringImpl(int64_t& value) override;

  void ReadUlongToStringImpl(uint64_t& value) override;

  void ReadFloatToStringImpl(float& value) override;

  void ReadDoubleToStringImpl(double& value) override;

  void ReadIntToOptionalImpl(int32_t& value) override;

  void ReadFloatToOptionalImpl(float& value) override;

  void ReadStringToOptionalImpl(std::string& value) override;

  void ReadIntToUnionImpl(int32_t& value) override;

  void ReadFloatToUnionImpl(float& value) override;

  void ReadStringToUnionImpl(std::string& value) override;

  void ReadOptionalIntToFloatImpl(std::optional<int32_t>& value) override;

  void ReadOptionalFloatToStringImpl(std::optional<float>& value) override;

  void ReadAliasedLongToStringImpl(evo_test::AliasedLongToString& value) override;

  void ReadRecordWithChangesImpl(evo_test::RecordWithChanges& value) override;

  void ReadAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges& value) override;

  void ReadOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges>& value) override;

  void ReadAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges>& value) override;

  bool ReadStreamedRecordWithChangesImpl(evo_test::RecordWithChanges& value) override;

  bool ReadStreamedRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges>& values) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetReader> streamedRecordWithChanges_dataset_state_;
};

// HDF5 writer for the UnusedProtocol protocol.
class UnusedProtocolWriter : public evo_test::UnusedProtocolWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  UnusedProtocolWriter(std::string path);

  protected:
  void WriteSamplesImpl(evo_test::UnchangedRecord const& value) override;

  void WriteSamplesImpl(std::vector<evo_test::UnchangedRecord> const& values) override;

  void EndSamplesImpl() override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetWriter> samples_dataset_state_;
};

// HDF5 reader for the UnusedProtocol protocol.
class UnusedProtocolReader : public evo_test::UnusedProtocolReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  UnusedProtocolReader(std::string path);

  bool ReadSamplesImpl(evo_test::UnchangedRecord& value) override;

  bool ReadSamplesImpl(std::vector<evo_test::UnchangedRecord>& values) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetReader> samples_dataset_state_;
};

} // namespace evo_test

