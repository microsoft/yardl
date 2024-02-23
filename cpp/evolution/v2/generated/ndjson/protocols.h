// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include <array>
#include <complex>
#include <memory>
#include <optional>
#include <variant>
#include <vector>

#include "../yardl/detail/ndjson/reader_writer.h"
#include "../protocols.h"
#include "../types.h"

namespace evo_test::ndjson {
// NDJSON writer for the ProtocolWithChanges protocol.
class ProtocolWithChangesWriter : public evo_test::ProtocolWithChangesWriterBase, yardl::ndjson::NDJsonWriter {
  public:
  ProtocolWithChangesWriter(std::ostream& stream)
      : yardl::ndjson::NDJsonWriter(stream, schema_) {
  }

  ProtocolWithChangesWriter(std::string file_name)
      : yardl::ndjson::NDJsonWriter(file_name, schema_) {
  }

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
  void WriteComplexFloatToComplexDoubleImpl(std::complex<float> const& value) override;
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
  void WriteEnumToAliasedEnumImpl(evo_test::GrowingEnum const& value) override;
  void WriteOptionalIntToUnionImpl(std::optional<int32_t> const& value) override;
  void WriteOptionalRecordToUnionImpl(std::optional<evo_test::RecordWithChanges> const& value) override;
  void WriteRecordWithChangesImpl(evo_test::RecordWithChanges const& value) override;
  void WriteAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges const& value) override;
  void WriteRecordToRenamedRecordImpl(evo_test::RenamedRecord const& value) override;
  void WriteRecordToAliasedRecordImpl(evo_test::RecordWithChanges const& value) override;
  void WriteRecordToAliasedAliasImpl(evo_test::RecordWithChanges const& value) override;
  // Stream and Vector type changes
  void WriteStreamIntToStringToFloatImpl(float const& value) override;
  void EndStreamIntToStringToFloatImpl() override {}
  void WriteVectorIntToStringToFloatImpl(std::vector<float> const& value) override;
  void WriteIntFloatUnionReorderedImpl(std::variant<int32_t, float> const& value) override;
  void WriteVectorUnionReorderedImpl(std::vector<std::variant<int32_t, float>> const& value) override;
  void WriteStreamUnionReorderedImpl(std::variant<int32_t, std::string> const& value) override;
  void EndStreamUnionReorderedImpl() override {}
  void WriteIntToUnionStreamImpl(std::variant<std::string, int32_t> const& value) override;
  void EndIntToUnionStreamImpl() override {}
  void WriteUnionStreamTypeChangeImpl(std::variant<int32_t, float> const& value) override;
  void EndUnionStreamTypeChangeImpl() override {}
  void WriteStreamOfAliasTypeChangeImpl(evo_test::StreamItem const& value) override;
  void EndStreamOfAliasTypeChangeImpl() override {}
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
  void WriteUnionToAliasedUnionImpl(std::variant<evo_test::RecordWithChanges, int32_t> const& value) override;
  void WriteUnionToAliasedUnionWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t> const& value) override;
  void WriteOptionalToAliasedOptionalImpl(std::optional<evo_test::RecordWithChanges> const& value) override;
  void WriteOptionalToAliasedOptionalWithChangesImpl(std::optional<int32_t> const& value) override;
  void WriteGenericRecordImpl(evo_test::GenericRecord<int32_t, std::string> const& value) override;
  void WriteGenericRecordToOpenAliasImpl(evo_test::GenericRecord<int32_t, std::string> const& value) override;
  void WriteGenericRecordToClosedAliasImpl(evo_test::GenericRecord<int32_t, std::string> const& value) override;
  void WriteGenericRecordToHalfClosedAliasImpl(evo_test::GenericRecord<int32_t, std::string> const& value) override;
  void WriteAliasedGenericRecordToAliasImpl(evo_test::AliasedHalfClosedGenericRecord<int32_t> const& value) override;
  void WriteGenericRecordToReversedImpl(evo_test::GenericRecord<int32_t, std::string> const& value) override;
  void WriteClosedGenericRecordToUnionImpl(evo_test::AliasedClosedGenericRecord const& value) override;
  void WriteGenericRecordToAliasedUnionImpl(evo_test::GenericRecord<int32_t, std::string> const& value) override;
  void WriteGenericUnionToReversedImpl(evo_test::GenericUnion<evo_test::GenericRecord<int32_t, std::string>, float> const& value) override;
  void WriteGenericUnionOfChangedRecordImpl(evo_test::GenericUnion<evo_test::GenericRecord<int32_t, std::string>, float> const& value) override;
  void WriteGenericParentRecordImpl(evo_test::GenericParentRecord<int32_t> const& value) override;
  void WriteGenericNestedRecordsImpl(evo_test::GenericRecord<evo_test::UnchangedGeneric<int32_t>, evo_test::ChangedGeneric<std::string, int32_t>> const& value) override;
  void WriteGenericRecordStreamImpl(evo_test::GenericRecord<int32_t, std::string> const& value) override;
  void EndGenericRecordStreamImpl() override {}
  void WriteGenericParentRecordStreamImpl(evo_test::GenericParentRecord<int32_t> const& value) override;
  void EndGenericParentRecordStreamImpl() override {}
  void WriteVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges> const& value) override;
  void WriteStreamedRecordWithChangesImpl(evo_test::RecordWithChanges const& value) override;
  void EndStreamedRecordWithChangesImpl() override {}
  void WriteAddedStringVectorImpl(std::vector<evo_test::AliasedString> const& value) override;
  void WriteAddedOptionalImpl(std::optional<evo_test::RecordWithChanges> const& value) override;
  void WriteAddedMapImpl(std::unordered_map<std::string, std::string> const& value) override;
  void WriteAddedUnionImpl(std::variant<std::monostate, evo_test::RecordWithChanges, std::string> const& value) override;
  void WriteAddedRecordStreamImpl(evo_test::RecordWithChanges const& value) override;
  void EndAddedRecordStreamImpl() override {}
  void WriteAddedUnionStreamImpl(std::variant<evo_test::RecordWithChanges, evo_test::RenamedRecord> const& value) override;
  void EndAddedUnionStreamImpl() override {}
  void CloseImpl() override;
};

// NDJSON reader for the ProtocolWithChanges protocol.
class ProtocolWithChangesReader : public evo_test::ProtocolWithChangesReaderBase, yardl::ndjson::NDJsonReader {
  public:
  ProtocolWithChangesReader(std::istream& stream)
      : yardl::ndjson::NDJsonReader(stream, schema_) {
  }

  ProtocolWithChangesReader(std::string file_name)
      : yardl::ndjson::NDJsonReader(file_name, schema_) {
  }

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
  void ReadComplexFloatToComplexDoubleImpl(std::complex<float>& value) override;
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
  void ReadEnumToAliasedEnumImpl(evo_test::GrowingEnum& value) override;
  void ReadOptionalIntToUnionImpl(std::optional<int32_t>& value) override;
  void ReadOptionalRecordToUnionImpl(std::optional<evo_test::RecordWithChanges>& value) override;
  void ReadRecordWithChangesImpl(evo_test::RecordWithChanges& value) override;
  void ReadAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges& value) override;
  void ReadRecordToRenamedRecordImpl(evo_test::RenamedRecord& value) override;
  void ReadRecordToAliasedRecordImpl(evo_test::RecordWithChanges& value) override;
  void ReadRecordToAliasedAliasImpl(evo_test::RecordWithChanges& value) override;
  bool ReadStreamIntToStringToFloatImpl(float& value) override;
  void ReadVectorIntToStringToFloatImpl(std::vector<float>& value) override;
  void ReadIntFloatUnionReorderedImpl(std::variant<int32_t, float>& value) override;
  void ReadVectorUnionReorderedImpl(std::vector<std::variant<int32_t, float>>& value) override;
  bool ReadStreamUnionReorderedImpl(std::variant<int32_t, std::string>& value) override;
  bool ReadIntToUnionStreamImpl(std::variant<std::string, int32_t>& value) override;
  bool ReadUnionStreamTypeChangeImpl(std::variant<int32_t, float>& value) override;
  bool ReadStreamOfAliasTypeChangeImpl(evo_test::StreamItem& value) override;
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
  void ReadUnionToAliasedUnionImpl(std::variant<evo_test::RecordWithChanges, int32_t>& value) override;
  void ReadUnionToAliasedUnionWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t>& value) override;
  void ReadOptionalToAliasedOptionalImpl(std::optional<evo_test::RecordWithChanges>& value) override;
  void ReadOptionalToAliasedOptionalWithChangesImpl(std::optional<int32_t>& value) override;
  void ReadGenericRecordImpl(evo_test::GenericRecord<int32_t, std::string>& value) override;
  void ReadGenericRecordToOpenAliasImpl(evo_test::GenericRecord<int32_t, std::string>& value) override;
  void ReadGenericRecordToClosedAliasImpl(evo_test::GenericRecord<int32_t, std::string>& value) override;
  void ReadGenericRecordToHalfClosedAliasImpl(evo_test::GenericRecord<int32_t, std::string>& value) override;
  void ReadAliasedGenericRecordToAliasImpl(evo_test::AliasedHalfClosedGenericRecord<int32_t>& value) override;
  void ReadGenericRecordToReversedImpl(evo_test::GenericRecord<int32_t, std::string>& value) override;
  void ReadClosedGenericRecordToUnionImpl(evo_test::AliasedClosedGenericRecord& value) override;
  void ReadGenericRecordToAliasedUnionImpl(evo_test::GenericRecord<int32_t, std::string>& value) override;
  void ReadGenericUnionToReversedImpl(evo_test::GenericUnion<evo_test::GenericRecord<int32_t, std::string>, float>& value) override;
  void ReadGenericUnionOfChangedRecordImpl(evo_test::GenericUnion<evo_test::GenericRecord<int32_t, std::string>, float>& value) override;
  void ReadGenericParentRecordImpl(evo_test::GenericParentRecord<int32_t>& value) override;
  void ReadGenericNestedRecordsImpl(evo_test::GenericRecord<evo_test::UnchangedGeneric<int32_t>, evo_test::ChangedGeneric<std::string, int32_t>>& value) override;
  bool ReadGenericRecordStreamImpl(evo_test::GenericRecord<int32_t, std::string>& value) override;
  bool ReadGenericParentRecordStreamImpl(evo_test::GenericParentRecord<int32_t>& value) override;
  void ReadVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges>& value) override;
  bool ReadStreamedRecordWithChangesImpl(evo_test::RecordWithChanges& value) override;
  void ReadAddedStringVectorImpl(std::vector<evo_test::AliasedString>& value) override;
  void ReadAddedOptionalImpl(std::optional<evo_test::RecordWithChanges>& value) override;
  void ReadAddedMapImpl(std::unordered_map<std::string, std::string>& value) override;
  void ReadAddedUnionImpl(std::variant<std::monostate, evo_test::RecordWithChanges, std::string>& value) override;
  bool ReadAddedRecordStreamImpl(evo_test::RecordWithChanges& value) override;
  bool ReadAddedUnionStreamImpl(std::variant<evo_test::RecordWithChanges, evo_test::RenamedRecord>& value) override;
  void CloseImpl() override;
};

} // namespace evo_test::ndjson
