// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include <array>
#include <complex>
#include <memory>
#include <optional>
#include <variant>
#include <vector>

#include "../yardl/detail/binary/reader_writer.h"
#include "../protocols.h"
#include "../types.h"

namespace evo_test::binary {
// Binary writer for the ProtocolWithChanges protocol.
class ProtocolWithChangesWriter : public evo_test::ProtocolWithChangesWriterBase, yardl::binary::BinaryWriter {
  public:
  ProtocolWithChangesWriter(std::ostream& stream, Version version = Version::Latest)
      : yardl::binary::BinaryWriter(stream, evo_test::ProtocolWithChangesWriterBase::SchemaFromVersion(version)), version_(version) {}

  ProtocolWithChangesWriter(std::string file_name, Version version = Version::Latest)
      : yardl::binary::BinaryWriter(file_name, evo_test::ProtocolWithChangesWriterBase::SchemaFromVersion(version)), version_(version) {}

  void Flush() override;

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
  void WriteStringToAliasedStringImpl(std::string const& value) override;
  void WriteStringToAliasedIntImpl(std::string const& value) override;
  void WriteOptionalIntToUnionImpl(std::optional<int32_t> const& value) override;
  void WriteOptionalRecordToUnionImpl(std::optional<evo_test::RecordWithChanges> const& value) override;
  void WriteRecordWithChangesImpl(evo_test::RecordWithChanges const& value) override;
  void WriteAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges const& value) override;
  void WriteRecordToRenamedRecordImpl(evo_test::RenamedRecord const& value) override;
  void WriteRecordToAliasedRecordImpl(evo_test::RecordWithChanges const& value) override;
  void WriteRecordToAliasedAliasImpl(evo_test::RecordWithChanges const& value) override;
  // Comprehensive NamedType changes
  void WriteRlinkImpl(evo_test::RLink const& value) override;
  void WriteRlinkRXImpl(evo_test::RLink const& value) override;
  void WriteRlinkRYImpl(evo_test::RLink const& value) override;
  void WriteRlinkRZImpl(evo_test::RLink const& value) override;
  void WriteRaRLinkImpl(evo_test::RA const& value) override;
  void WriteRaRXImpl(evo_test::RA const& value) override;
  void WriteRaRYImpl(evo_test::RA const& value) override;
  void WriteRaRZImpl(evo_test::RA const& value) override;
  void WriteRbRLinkImpl(evo_test::RB const& value) override;
  void WriteRbRXImpl(evo_test::RB const& value) override;
  void WriteRbRYImpl(evo_test::RB const& value) override;
  void WriteRbRZImpl(evo_test::RB const& value) override;
  void WriteRcRLinkImpl(evo_test::RC const& value) override;
  void WriteRcRXImpl(evo_test::RC const& value) override;
  void WriteRcRYImpl(evo_test::RC const& value) override;
  void WriteRcRZImpl(evo_test::RC const& value) override;
  void WriteRlinkRNewImpl(evo_test::RLink const& value) override;
  void WriteRaRNewImpl(evo_test::RA const& value) override;
  void WriteRbRNewImpl(evo_test::RB const& value) override;
  void WriteRcRNewImpl(evo_test::RC const& value) override;
  void WriteRlinkRUnionImpl(evo_test::RLink const& value) override;
  void WriteRaRUnionImpl(evo_test::RA const& value) override;
  void WriteRbRUnionImpl(evo_test::RB const& value) override;
  void WriteRcRUnionImpl(evo_test::RC const& value) override;
  void WriteOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges> const& value) override;
  void WriteAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges> const& value) override;
  void WriteUnionRecordWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t> const& value) override;
  void WriteUnionWithSameTypesetImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value) override;
  void WriteUnionWithTypesAddedImpl(std::variant<evo_test::RecordWithChanges, float> const& value) override;
  void WriteUnionWithTypesRemovedImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value) override;
  void WriteRecordToOptionalImpl(evo_test::RecordWithChanges const& value) override;
  void WriteRecordToAliasedOptionalImpl(evo_test::RecordWithChanges const& value) override;
  void WriteRecordToUnionImpl(evo_test::RecordWithChanges const& value) override;
  void WriteRecordToAliasedUnionImpl(evo_test::RecordWithChanges const& value) override;
  void WriteVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges> const& value) override;
  void WriteStreamedRecordWithChangesImpl(evo_test::RecordWithChanges const& value) override;
  void WriteStreamedRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges> const& values) override;
  void EndStreamedRecordWithChangesImpl() override;
  void CloseImpl() override;

  Version version_;
};

// Binary reader for the ProtocolWithChanges protocol.
class ProtocolWithChangesReader : public evo_test::ProtocolWithChangesReaderBase, yardl::binary::BinaryReader {
  public:
  ProtocolWithChangesReader(std::istream& stream)
      : yardl::binary::BinaryReader(stream), version_(evo_test::ProtocolWithChangesReaderBase::VersionFromSchema(schema_read_)) {}

  ProtocolWithChangesReader(std::string file_name)
      : yardl::binary::BinaryReader(file_name), version_(evo_test::ProtocolWithChangesReaderBase::VersionFromSchema(schema_read_)) {}

  Version GetVersion() { return version_; }

  protected:
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
  void ReadStringToAliasedStringImpl(std::string& value) override;
  void ReadStringToAliasedIntImpl(std::string& value) override;
  void ReadOptionalIntToUnionImpl(std::optional<int32_t>& value) override;
  void ReadOptionalRecordToUnionImpl(std::optional<evo_test::RecordWithChanges>& value) override;
  void ReadRecordWithChangesImpl(evo_test::RecordWithChanges& value) override;
  void ReadAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges& value) override;
  void ReadRecordToRenamedRecordImpl(evo_test::RenamedRecord& value) override;
  void ReadRecordToAliasedRecordImpl(evo_test::RecordWithChanges& value) override;
  void ReadRecordToAliasedAliasImpl(evo_test::RecordWithChanges& value) override;
  void ReadRlinkImpl(evo_test::RLink& value) override;
  void ReadRlinkRXImpl(evo_test::RLink& value) override;
  void ReadRlinkRYImpl(evo_test::RLink& value) override;
  void ReadRlinkRZImpl(evo_test::RLink& value) override;
  void ReadRaRLinkImpl(evo_test::RA& value) override;
  void ReadRaRXImpl(evo_test::RA& value) override;
  void ReadRaRYImpl(evo_test::RA& value) override;
  void ReadRaRZImpl(evo_test::RA& value) override;
  void ReadRbRLinkImpl(evo_test::RB& value) override;
  void ReadRbRXImpl(evo_test::RB& value) override;
  void ReadRbRYImpl(evo_test::RB& value) override;
  void ReadRbRZImpl(evo_test::RB& value) override;
  void ReadRcRLinkImpl(evo_test::RC& value) override;
  void ReadRcRXImpl(evo_test::RC& value) override;
  void ReadRcRYImpl(evo_test::RC& value) override;
  void ReadRcRZImpl(evo_test::RC& value) override;
  void ReadRlinkRNewImpl(evo_test::RLink& value) override;
  void ReadRaRNewImpl(evo_test::RA& value) override;
  void ReadRbRNewImpl(evo_test::RB& value) override;
  void ReadRcRNewImpl(evo_test::RC& value) override;
  void ReadRlinkRUnionImpl(evo_test::RLink& value) override;
  void ReadRaRUnionImpl(evo_test::RA& value) override;
  void ReadRbRUnionImpl(evo_test::RB& value) override;
  void ReadRcRUnionImpl(evo_test::RC& value) override;
  void ReadOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges>& value) override;
  void ReadAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges>& value) override;
  void ReadUnionRecordWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t>& value) override;
  void ReadUnionWithSameTypesetImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value) override;
  void ReadUnionWithTypesAddedImpl(std::variant<evo_test::RecordWithChanges, float>& value) override;
  void ReadUnionWithTypesRemovedImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value) override;
  void ReadRecordToOptionalImpl(evo_test::RecordWithChanges& value) override;
  void ReadRecordToAliasedOptionalImpl(evo_test::RecordWithChanges& value) override;
  void ReadRecordToUnionImpl(evo_test::RecordWithChanges& value) override;
  void ReadRecordToAliasedUnionImpl(evo_test::RecordWithChanges& value) override;
  void ReadVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges>& value) override;
  bool ReadStreamedRecordWithChangesImpl(evo_test::RecordWithChanges& value) override;
  bool ReadStreamedRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges>& values) override;
  void CloseImpl() override;

  Version version_;

  private:
  size_t current_block_remaining_ = 0;
};

// Binary writer for the UnusedProtocol protocol.
class UnusedProtocolWriter : public evo_test::UnusedProtocolWriterBase, yardl::binary::BinaryWriter {
  public:
  UnusedProtocolWriter(std::ostream& stream, Version version = Version::Latest)
      : yardl::binary::BinaryWriter(stream, evo_test::UnusedProtocolWriterBase::SchemaFromVersion(version)), version_(version) {}

  UnusedProtocolWriter(std::string file_name, Version version = Version::Latest)
      : yardl::binary::BinaryWriter(file_name, evo_test::UnusedProtocolWriterBase::SchemaFromVersion(version)), version_(version) {}

  void Flush() override;

  protected:
  void WriteRecordsImpl(evo_test::UnchangedRecord const& value) override;
  void WriteRecordsImpl(std::vector<evo_test::UnchangedRecord> const& values) override;
  void EndRecordsImpl() override;
  void CloseImpl() override;

  Version version_;
};

// Binary reader for the UnusedProtocol protocol.
class UnusedProtocolReader : public evo_test::UnusedProtocolReaderBase, yardl::binary::BinaryReader {
  public:
  UnusedProtocolReader(std::istream& stream)
      : yardl::binary::BinaryReader(stream), version_(evo_test::UnusedProtocolReaderBase::VersionFromSchema(schema_read_)) {}

  UnusedProtocolReader(std::string file_name)
      : yardl::binary::BinaryReader(file_name), version_(evo_test::UnusedProtocolReaderBase::VersionFromSchema(schema_read_)) {}

  Version GetVersion() { return version_; }

  protected:
  bool ReadRecordsImpl(evo_test::UnchangedRecord& value) override;
  bool ReadRecordsImpl(std::vector<evo_test::UnchangedRecord>& values) override;
  void CloseImpl() override;

  Version version_;

  private:
  size_t current_block_remaining_ = 0;
};

} // namespace evo_test::binary
