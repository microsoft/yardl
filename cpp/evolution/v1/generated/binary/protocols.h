// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include <array>
#include <complex>
#include <memory>
#include <optional>
#include <variant>
#include <vector>

#include "../protocols.h"
#include "../yardl/detail/binary/reader_writer.h"

namespace evo_test::binary {
// Binary writer for the ProtocolWithChanges protocol.
class ProtocolWithChangesWriter : public evo_test::ProtocolWithChangesWriterBase, yardl::binary::BinaryWriter {
  public:
  ProtocolWithChangesWriter(std::ostream& stream, Version version = Version::Current)
      : yardl::binary::BinaryWriter(stream, evo_test::ProtocolWithChangesWriterBase::SchemaFromVersion(version)), version_(version) {}

  ProtocolWithChangesWriter(std::string file_name, Version version = Version::Current)
      : yardl::binary::BinaryWriter(file_name, evo_test::ProtocolWithChangesWriterBase::SchemaFromVersion(version)), version_(version) {}

  void Flush() override;

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
  void WriteComplexFloatToComplexDoubleImpl(std::complex<double> const& value) override;
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
  void WriteStringToAliasedStringImpl(evo_test::AliasedString const& value) override;
  void WriteStringToAliasedIntImpl(evo_test::AliasedInt const& value) override;
  void WriteEnumToAliasedEnumImpl(evo_test::AliasedEnum const& value) override;
  void WriteOptionalIntToUnionImpl(std::variant<std::monostate, int32_t, std::string> const& value) override;
  void WriteOptionalRecordToUnionImpl(std::variant<std::monostate, evo_test::RecordWithChanges, std::string> const& value) override;
  void WriteRecordWithChangesImpl(evo_test::RecordWithChanges const& value) override;
  void WriteAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges const& value) override;
  void WriteRecordToRenamedRecordImpl(evo_test::RenamedRecord const& value) override;
  void WriteRecordToAliasedRecordImpl(evo_test::AliasedRecordWithChanges const& value) override;
  void WriteRecordToAliasedAliasImpl(evo_test::AliasOfAliasedRecordWithChanges const& value) override;
  // Stream and Vector type changes
  void WriteStreamIntToStringToFloatImpl(std::string const& value) override;
  void WriteStreamIntToStringToFloatImpl(std::vector<std::string> const& values) override;
  void EndStreamIntToStringToFloatImpl() override;
  void WriteVectorIntToStringToFloatImpl(std::vector<std::string> const& value) override;
  void WriteIntFloatUnionReorderedImpl(std::variant<float, int32_t> const& value) override;
  void WriteVectorUnionReorderedImpl(std::vector<std::variant<float, int32_t>> const& value) override;
  void WriteStreamUnionReorderedImpl(std::variant<std::string, int32_t> const& value) override;
  void WriteStreamUnionReorderedImpl(std::vector<std::variant<std::string, int32_t>> const& values) override;
  void EndStreamUnionReorderedImpl() override;
  void WriteIntToUnionStreamImpl(int32_t const& value) override;
  void WriteIntToUnionStreamImpl(std::vector<int32_t> const& values) override;
  void EndIntToUnionStreamImpl() override;
  void WriteUnionStreamTypeChangeImpl(std::variant<int32_t, bool> const& value) override;
  void WriteUnionStreamTypeChangeImpl(std::vector<std::variant<int32_t, bool>> const& values) override;
  void EndUnionStreamTypeChangeImpl() override;
  void WriteStreamOfAliasTypeChangeImpl(evo_test::StreamItem const& value) override;
  void WriteStreamOfAliasTypeChangeImpl(std::vector<evo_test::StreamItem> const& values) override;
  void EndStreamOfAliasTypeChangeImpl() override;
  // Comprehensive NamedType changes
  void WriteRlinkImpl(evo_test::RLink const& value) override;
  void WriteRlinkRXImpl(evo_test::RX const& value) override;
  void WriteRlinkRYImpl(evo_test::RY const& value) override;
  void WriteRlinkRZImpl(evo_test::RZ const& value) override;
  void WriteRaRLinkImpl(evo_test::RLink const& value) override;
  void WriteRaRXImpl(evo_test::RX const& value) override;
  void WriteRaRYImpl(evo_test::RY const& value) override;
  void WriteRaRZImpl(evo_test::RZ const& value) override;
  void WriteRbRLinkImpl(evo_test::RLink const& value) override;
  void WriteRbRXImpl(evo_test::RX const& value) override;
  void WriteRbRYImpl(evo_test::RY const& value) override;
  void WriteRbRZImpl(evo_test::RZ const& value) override;
  void WriteRcRLinkImpl(evo_test::RLink const& value) override;
  void WriteRcRXImpl(evo_test::RX const& value) override;
  void WriteRcRYImpl(evo_test::RY const& value) override;
  void WriteRcRZImpl(evo_test::RZ const& value) override;
  void WriteRlinkRNewImpl(evo_test::RNew const& value) override;
  void WriteRaRNewImpl(evo_test::RNew const& value) override;
  void WriteRbRNewImpl(evo_test::RNew const& value) override;
  void WriteRcRNewImpl(evo_test::RNew const& value) override;
  void WriteRlinkRUnionImpl(evo_test::RUnion const& value) override;
  void WriteRaRUnionImpl(evo_test::RUnion const& value) override;
  void WriteRbRUnionImpl(evo_test::RUnion const& value) override;
  void WriteRcRUnionImpl(evo_test::RUnion const& value) override;
  void WriteOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges> const& value) override;
  void WriteAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges> const& value) override;
  void WriteUnionRecordWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t> const& value) override;
  void WriteUnionWithSameTypesetImpl(std::variant<float, evo_test::RecordWithChanges, std::string, int32_t> const& value) override;
  void WriteUnionWithTypesAddedImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value) override;
  void WriteUnionWithTypesRemovedImpl(std::variant<evo_test::RecordWithChanges, std::string> const& value) override;
  void WriteRecordToOptionalImpl(std::optional<evo_test::RecordWithChanges> const& value) override;
  void WriteRecordToAliasedOptionalImpl(evo_test::AliasedOptionalRecord const& value) override;
  void WriteRecordToUnionImpl(std::variant<evo_test::RecordWithChanges, std::string> const& value) override;
  void WriteRecordToAliasedUnionImpl(evo_test::AliasedRecordOrString const& value) override;
  void WriteUnionToAliasedUnionImpl(evo_test::AliasedRecordOrInt const& value) override;
  void WriteUnionToAliasedUnionWithChangesImpl(evo_test::AliasedRecordOrString const& value) override;
  void WriteOptionalToAliasedOptionalImpl(evo_test::AliasedOptionalRecord const& value) override;
  void WriteOptionalToAliasedOptionalWithChangesImpl(evo_test::AliasedOptionalString const& value) override;
  void WriteGenericRecordImpl(evo_test::GenericRecord<int32_t, std::string> const& value) override;
  void WriteGenericRecordToOpenAliasImpl(evo_test::AliasedOpenGenericRecord<int32_t, std::string> const& value) override;
  void WriteGenericRecordToClosedAliasImpl(evo_test::AliasedClosedGenericRecord const& value) override;
  void WriteGenericRecordToHalfClosedAliasImpl(evo_test::AliasedHalfClosedGenericRecord<int32_t> const& value) override;
  void WriteAliasedGenericRecordToAliasImpl(evo_test::AliasedOpenGenericRecord<int32_t, std::string> const& value) override;
  void WriteGenericRecordToReversedImpl(evo_test::GenericRecordReversed<std::string, int32_t> const& value) override;
  void WriteClosedGenericRecordToUnionImpl(std::variant<evo_test::GenericRecord<int32_t, std::string>, std::string> const& value) override;
  void WriteGenericRecordToAliasedUnionImpl(evo_test::AliasedGenericRecordOrString const& value) override;
  void WriteGenericUnionToReversedImpl(evo_test::GenericUnionReversed<float, evo_test::GenericRecord<int32_t, std::string>> const& value) override;
  void WriteGenericUnionOfChangedRecordImpl(evo_test::GenericUnion<evo_test::GenericRecord<int32_t, std::string>, float> const& value) override;
  void WriteGenericParentRecordImpl(evo_test::GenericParentRecord<int32_t> const& value) override;
  void WriteGenericNestedRecordsImpl(evo_test::GenericRecord<evo_test::Unchanged, evo_test::Changed> const& value) override;
  void WriteGenericRecordStreamImpl(evo_test::AliasedClosedGenericRecord const& value) override;
  void WriteGenericRecordStreamImpl(std::vector<evo_test::AliasedClosedGenericRecord> const& values) override;
  void EndGenericRecordStreamImpl() override;
  void WriteGenericParentRecordStreamImpl(evo_test::GenericParentRecord<int32_t> const& value) override;
  void WriteGenericParentRecordStreamImpl(std::vector<evo_test::GenericParentRecord<int32_t>> const& values) override;
  void EndGenericParentRecordStreamImpl() override;
  void WriteVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges> const& value) override;
  void WriteStreamedRecordWithChangesImpl(evo_test::RecordWithChanges const& value) override;
  void WriteStreamedRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges> const& values) override;
  void EndStreamedRecordWithChangesImpl() override;
  void WriteAddedOptionalImpl(std::optional<evo_test::RecordWithChanges> const& value) override;
  void WriteAddedMapImpl(std::unordered_map<std::string, std::string> const& value) override;
  void WriteAddedRecordStreamImpl(evo_test::RecordWithChanges const& value) override;
  void WriteAddedRecordStreamImpl(std::vector<evo_test::RecordWithChanges> const& values) override;
  void EndAddedRecordStreamImpl() override;
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
  void ReadComplexFloatToComplexDoubleImpl(std::complex<double>& value) override;
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
  void ReadStringToAliasedStringImpl(evo_test::AliasedString& value) override;
  void ReadStringToAliasedIntImpl(evo_test::AliasedInt& value) override;
  void ReadEnumToAliasedEnumImpl(evo_test::AliasedEnum& value) override;
  void ReadOptionalIntToUnionImpl(std::variant<std::monostate, int32_t, std::string>& value) override;
  void ReadOptionalRecordToUnionImpl(std::variant<std::monostate, evo_test::RecordWithChanges, std::string>& value) override;
  void ReadRecordWithChangesImpl(evo_test::RecordWithChanges& value) override;
  void ReadAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges& value) override;
  void ReadRecordToRenamedRecordImpl(evo_test::RenamedRecord& value) override;
  void ReadRecordToAliasedRecordImpl(evo_test::AliasedRecordWithChanges& value) override;
  void ReadRecordToAliasedAliasImpl(evo_test::AliasOfAliasedRecordWithChanges& value) override;
  bool ReadStreamIntToStringToFloatImpl(std::string& value) override;
  bool ReadStreamIntToStringToFloatImpl(std::vector<std::string>& values) override;
  void ReadVectorIntToStringToFloatImpl(std::vector<std::string>& value) override;
  void ReadIntFloatUnionReorderedImpl(std::variant<float, int32_t>& value) override;
  void ReadVectorUnionReorderedImpl(std::vector<std::variant<float, int32_t>>& value) override;
  bool ReadStreamUnionReorderedImpl(std::variant<std::string, int32_t>& value) override;
  bool ReadStreamUnionReorderedImpl(std::vector<std::variant<std::string, int32_t>>& values) override;
  bool ReadIntToUnionStreamImpl(int32_t& value) override;
  bool ReadIntToUnionStreamImpl(std::vector<int32_t>& values) override;
  bool ReadUnionStreamTypeChangeImpl(std::variant<int32_t, bool>& value) override;
  bool ReadUnionStreamTypeChangeImpl(std::vector<std::variant<int32_t, bool>>& values) override;
  bool ReadStreamOfAliasTypeChangeImpl(evo_test::StreamItem& value) override;
  bool ReadStreamOfAliasTypeChangeImpl(std::vector<evo_test::StreamItem>& values) override;
  void ReadRlinkImpl(evo_test::RLink& value) override;
  void ReadRlinkRXImpl(evo_test::RX& value) override;
  void ReadRlinkRYImpl(evo_test::RY& value) override;
  void ReadRlinkRZImpl(evo_test::RZ& value) override;
  void ReadRaRLinkImpl(evo_test::RLink& value) override;
  void ReadRaRXImpl(evo_test::RX& value) override;
  void ReadRaRYImpl(evo_test::RY& value) override;
  void ReadRaRZImpl(evo_test::RZ& value) override;
  void ReadRbRLinkImpl(evo_test::RLink& value) override;
  void ReadRbRXImpl(evo_test::RX& value) override;
  void ReadRbRYImpl(evo_test::RY& value) override;
  void ReadRbRZImpl(evo_test::RZ& value) override;
  void ReadRcRLinkImpl(evo_test::RLink& value) override;
  void ReadRcRXImpl(evo_test::RX& value) override;
  void ReadRcRYImpl(evo_test::RY& value) override;
  void ReadRcRZImpl(evo_test::RZ& value) override;
  void ReadRlinkRNewImpl(evo_test::RNew& value) override;
  void ReadRaRNewImpl(evo_test::RNew& value) override;
  void ReadRbRNewImpl(evo_test::RNew& value) override;
  void ReadRcRNewImpl(evo_test::RNew& value) override;
  void ReadRlinkRUnionImpl(evo_test::RUnion& value) override;
  void ReadRaRUnionImpl(evo_test::RUnion& value) override;
  void ReadRbRUnionImpl(evo_test::RUnion& value) override;
  void ReadRcRUnionImpl(evo_test::RUnion& value) override;
  void ReadOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges>& value) override;
  void ReadAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges>& value) override;
  void ReadUnionRecordWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t>& value) override;
  void ReadUnionWithSameTypesetImpl(std::variant<float, evo_test::RecordWithChanges, std::string, int32_t>& value) override;
  void ReadUnionWithTypesAddedImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value) override;
  void ReadUnionWithTypesRemovedImpl(std::variant<evo_test::RecordWithChanges, std::string>& value) override;
  void ReadRecordToOptionalImpl(std::optional<evo_test::RecordWithChanges>& value) override;
  void ReadRecordToAliasedOptionalImpl(evo_test::AliasedOptionalRecord& value) override;
  void ReadRecordToUnionImpl(std::variant<evo_test::RecordWithChanges, std::string>& value) override;
  void ReadRecordToAliasedUnionImpl(evo_test::AliasedRecordOrString& value) override;
  void ReadUnionToAliasedUnionImpl(evo_test::AliasedRecordOrInt& value) override;
  void ReadUnionToAliasedUnionWithChangesImpl(evo_test::AliasedRecordOrString& value) override;
  void ReadOptionalToAliasedOptionalImpl(evo_test::AliasedOptionalRecord& value) override;
  void ReadOptionalToAliasedOptionalWithChangesImpl(evo_test::AliasedOptionalString& value) override;
  void ReadGenericRecordImpl(evo_test::GenericRecord<int32_t, std::string>& value) override;
  void ReadGenericRecordToOpenAliasImpl(evo_test::AliasedOpenGenericRecord<int32_t, std::string>& value) override;
  void ReadGenericRecordToClosedAliasImpl(evo_test::AliasedClosedGenericRecord& value) override;
  void ReadGenericRecordToHalfClosedAliasImpl(evo_test::AliasedHalfClosedGenericRecord<int32_t>& value) override;
  void ReadAliasedGenericRecordToAliasImpl(evo_test::AliasedOpenGenericRecord<int32_t, std::string>& value) override;
  void ReadGenericRecordToReversedImpl(evo_test::GenericRecordReversed<std::string, int32_t>& value) override;
  void ReadClosedGenericRecordToUnionImpl(std::variant<evo_test::GenericRecord<int32_t, std::string>, std::string>& value) override;
  void ReadGenericRecordToAliasedUnionImpl(evo_test::AliasedGenericRecordOrString& value) override;
  void ReadGenericUnionToReversedImpl(evo_test::GenericUnionReversed<float, evo_test::GenericRecord<int32_t, std::string>>& value) override;
  void ReadGenericUnionOfChangedRecordImpl(evo_test::GenericUnion<evo_test::GenericRecord<int32_t, std::string>, float>& value) override;
  void ReadGenericParentRecordImpl(evo_test::GenericParentRecord<int32_t>& value) override;
  void ReadGenericNestedRecordsImpl(evo_test::GenericRecord<evo_test::Unchanged, evo_test::Changed>& value) override;
  bool ReadGenericRecordStreamImpl(evo_test::AliasedClosedGenericRecord& value) override;
  bool ReadGenericRecordStreamImpl(std::vector<evo_test::AliasedClosedGenericRecord>& values) override;
  bool ReadGenericParentRecordStreamImpl(evo_test::GenericParentRecord<int32_t>& value) override;
  bool ReadGenericParentRecordStreamImpl(std::vector<evo_test::GenericParentRecord<int32_t>>& values) override;
  void ReadVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges>& value) override;
  bool ReadStreamedRecordWithChangesImpl(evo_test::RecordWithChanges& value) override;
  bool ReadStreamedRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges>& values) override;
  void ReadAddedOptionalImpl(std::optional<evo_test::RecordWithChanges>& value) override;
  void ReadAddedMapImpl(std::unordered_map<std::string, std::string>& value) override;
  bool ReadAddedRecordStreamImpl(evo_test::RecordWithChanges& value) override;
  bool ReadAddedRecordStreamImpl(std::vector<evo_test::RecordWithChanges>& values) override;
  void CloseImpl() override;

  Version version_;

  private:
  size_t current_block_remaining_ = 0;
};

// Binary writer for the UnusedProtocol protocol.
class UnusedProtocolWriter : public evo_test::UnusedProtocolWriterBase, yardl::binary::BinaryWriter {
  public:
  UnusedProtocolWriter(std::ostream& stream, Version version = Version::Current)
      : yardl::binary::BinaryWriter(stream, evo_test::UnusedProtocolWriterBase::SchemaFromVersion(version)), version_(version) {}

  UnusedProtocolWriter(std::string file_name, Version version = Version::Current)
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
