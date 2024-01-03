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
  void WriteInt8ToIntImpl(int32_t const& value) override;

  void WriteInt8ToLongImpl(int64_t const& value) override;

  void WriteInt8ToUintImpl(uint32_t const& value) override;

  void WriteInt8ToUlongImpl(uint64_t const& value) override;

  void WriteInt8ToFloatImpl(float const& value) override;

  void WriteInt8ToDoubleImpl(double const& value) override;

  void WriteIntToUintImpl(uint32_t const& value) override;

  void WriteIntToLongImpl(int64_t const& value) override;

  void WriteIntToFloatImpl(float const& value) override;

  void WriteIntToDoubleImpl(double const& value) override;

  void WriteUintToUlongImpl(uint64_t const& value) override;

  void WriteUintToFloatImpl(float const& value) override;

  void WriteUintToDoubleImpl(double const& value) override;

  void WriteFloatToDoubleImpl(double const& value) override;

  void WriteIntToStringImpl(std::string const& value) override;

  void WriteUintToStringImpl(std::string const& value) override;

  void WriteLongToStringImpl(std::string const& value) override;

  void WriteUlongToStringImpl(std::string const& value) override;

  void WriteFloatToStringImpl(std::string const& value) override;

  void WriteDoubleToStringImpl(std::string const& value) override;

  void WriteIntToOptionalImpl(std::optional<int32_t> const& value) override;

  void WriteFloatToOptionalImpl(std::optional<float> const& value) override;

  void WriteStringToOptionalImpl(std::optional<std::string> const& value) override;

  void WriteIntToUnionImpl(std::variant<int32_t, bool> const& value) override;

  void WriteFloatToUnionImpl(std::variant<float, bool> const& value) override;

  void WriteStringToUnionImpl(std::variant<std::string, bool> const& value) override;

  void WriteOptionalIntToFloatImpl(std::optional<float> const& value) override;

  void WriteOptionalFloatToStringImpl(std::optional<std::string> const& value) override;

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

  void ReadInt8ToIntImpl(int32_t& value) override;

  void ReadInt8ToLongImpl(int64_t& value) override;

  void ReadInt8ToUintImpl(uint32_t& value) override;

  void ReadInt8ToUlongImpl(uint64_t& value) override;

  void ReadInt8ToFloatImpl(float& value) override;

  void ReadInt8ToDoubleImpl(double& value) override;

  void ReadIntToUintImpl(uint32_t& value) override;

  void ReadIntToLongImpl(int64_t& value) override;

  void ReadIntToFloatImpl(float& value) override;

  void ReadIntToDoubleImpl(double& value) override;

  void ReadUintToUlongImpl(uint64_t& value) override;

  void ReadUintToFloatImpl(float& value) override;

  void ReadUintToDoubleImpl(double& value) override;

  void ReadFloatToDoubleImpl(double& value) override;

  void ReadIntToStringImpl(std::string& value) override;

  void ReadUintToStringImpl(std::string& value) override;

  void ReadLongToStringImpl(std::string& value) override;

  void ReadUlongToStringImpl(std::string& value) override;

  void ReadFloatToStringImpl(std::string& value) override;

  void ReadDoubleToStringImpl(std::string& value) override;

  void ReadIntToOptionalImpl(std::optional<int32_t>& value) override;

  void ReadFloatToOptionalImpl(std::optional<float>& value) override;

  void ReadStringToOptionalImpl(std::optional<std::string>& value) override;

  void ReadIntToUnionImpl(std::variant<int32_t, bool>& value) override;

  void ReadFloatToUnionImpl(std::variant<float, bool>& value) override;

  void ReadStringToUnionImpl(std::variant<std::string, bool>& value) override;

  void ReadOptionalIntToFloatImpl(std::optional<float>& value) override;

  void ReadOptionalFloatToStringImpl(std::optional<std::string>& value) override;

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

