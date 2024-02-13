// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include "types.h"

namespace evo_test {
enum class Version {
  v0,
  Current
};
// Abstract writer for the ProtocolWithChanges protocol.
class ProtocolWithChangesWriterBase {
  public:
  // Ordinal 0.
  void WriteInt8ToInt(int32_t const& value);

  // Ordinal 1.
  void WriteInt8ToLong(int64_t const& value);

  // Ordinal 2.
  void WriteInt8ToUint(uint32_t const& value);

  // Ordinal 3.
  void WriteInt8ToUlong(uint64_t const& value);

  // Ordinal 4.
  void WriteInt8ToFloat(float const& value);

  // Ordinal 5.
  void WriteInt8ToDouble(double const& value);

  // Ordinal 6.
  void WriteIntToUint(uint32_t const& value);

  // Ordinal 7.
  void WriteIntToLong(int64_t const& value);

  // Ordinal 8.
  void WriteIntToFloat(float const& value);

  // Ordinal 9.
  void WriteIntToDouble(double const& value);

  // Ordinal 10.
  void WriteUintToUlong(uint64_t const& value);

  // Ordinal 11.
  void WriteUintToFloat(float const& value);

  // Ordinal 12.
  void WriteUintToDouble(double const& value);

  // Ordinal 13.
  void WriteFloatToDouble(double const& value);

  // Ordinal 14.
  void WriteIntToString(std::string const& value);

  // Ordinal 15.
  void WriteUintToString(std::string const& value);

  // Ordinal 16.
  void WriteLongToString(std::string const& value);

  // Ordinal 17.
  void WriteUlongToString(std::string const& value);

  // Ordinal 18.
  void WriteFloatToString(std::string const& value);

  // Ordinal 19.
  void WriteDoubleToString(std::string const& value);

  // Ordinal 20.
  void WriteIntToOptional(std::optional<int32_t> const& value);

  // Ordinal 21.
  void WriteFloatToOptional(std::optional<float> const& value);

  // Ordinal 22.
  void WriteStringToOptional(std::optional<std::string> const& value);

  // Ordinal 23.
  void WriteIntToUnion(std::variant<int32_t, bool> const& value);

  // Ordinal 24.
  void WriteFloatToUnion(std::variant<float, bool> const& value);

  // Ordinal 25.
  void WriteStringToUnion(std::variant<std::string, bool> const& value);

  // Ordinal 26.
  void WriteOptionalIntToFloat(std::optional<float> const& value);

  // Ordinal 27.
  void WriteOptionalFloatToString(std::optional<std::string> const& value);

  // Ordinal 28.
  void WriteAliasedLongToString(evo_test::AliasedLongToString const& value);

  // Ordinal 29.
  void WriteStringToAliasedString(evo_test::AliasedString const& value);

  // Ordinal 30.
  void WriteStringToAliasedInt(evo_test::AliasedInt const& value);

  // Ordinal 31.
  void WriteOptionalIntToUnion(std::variant<std::monostate, int32_t, std::string> const& value);

  // Ordinal 32.
  void WriteOptionalRecordToUnion(std::variant<std::monostate, evo_test::RecordWithChanges, std::string> const& value);

  // Ordinal 33.
  void WriteRecordWithChanges(evo_test::RecordWithChanges const& value);

  // Ordinal 34.
  void WriteAliasedRecordWithChanges(evo_test::AliasedRecordWithChanges const& value);

  // Ordinal 35.
  void WriteRecordToRenamedRecord(evo_test::RenamedRecord const& value);

  // Ordinal 36.
  void WriteRecordToAliasedRecord(evo_test::AliasedRecordWithChanges const& value);

  // Ordinal 37.
  void WriteRecordToAliasedAlias(evo_test::AliasOfAliasedRecordWithChanges const& value);

  // Ordinal 38.
  // Call this method for each element of the `streamOfAliasTypeChange` stream, then call `EndStreamOfAliasTypeChange() when done.`
  void WriteStreamOfAliasTypeChange(evo_test::StreamItem const& value);

  // Ordinal 38.
  // Call this method to write many values to the `streamOfAliasTypeChange` stream, then call `EndStreamOfAliasTypeChange()` when done.
  void WriteStreamOfAliasTypeChange(std::vector<evo_test::StreamItem> const& values);

  // Marks the end of the `streamOfAliasTypeChange` stream.
  void EndStreamOfAliasTypeChange();

  // Ordinal 39.
  // Comprehensive NamedType changes
  void WriteRlink(evo_test::RLink const& value);

  // Ordinal 40.
  void WriteRlinkRX(evo_test::RX const& value);

  // Ordinal 41.
  void WriteRlinkRY(evo_test::RY const& value);

  // Ordinal 42.
  void WriteRlinkRZ(evo_test::RZ const& value);

  // Ordinal 43.
  void WriteRaRLink(evo_test::RLink const& value);

  // Ordinal 44.
  void WriteRaRX(evo_test::RX const& value);

  // Ordinal 45.
  void WriteRaRY(evo_test::RY const& value);

  // Ordinal 46.
  void WriteRaRZ(evo_test::RZ const& value);

  // Ordinal 47.
  void WriteRbRLink(evo_test::RLink const& value);

  // Ordinal 48.
  void WriteRbRX(evo_test::RX const& value);

  // Ordinal 49.
  void WriteRbRY(evo_test::RY const& value);

  // Ordinal 50.
  void WriteRbRZ(evo_test::RZ const& value);

  // Ordinal 51.
  void WriteRcRLink(evo_test::RLink const& value);

  // Ordinal 52.
  void WriteRcRX(evo_test::RX const& value);

  // Ordinal 53.
  void WriteRcRY(evo_test::RY const& value);

  // Ordinal 54.
  void WriteRcRZ(evo_test::RZ const& value);

  // Ordinal 55.
  void WriteRlinkRNew(evo_test::RNew const& value);

  // Ordinal 56.
  void WriteRaRNew(evo_test::RNew const& value);

  // Ordinal 57.
  void WriteRbRNew(evo_test::RNew const& value);

  // Ordinal 58.
  void WriteRcRNew(evo_test::RNew const& value);

  // Ordinal 59.
  void WriteRlinkRUnion(evo_test::RUnion const& value);

  // Ordinal 60.
  void WriteRaRUnion(evo_test::RUnion const& value);

  // Ordinal 61.
  void WriteRbRUnion(evo_test::RUnion const& value);

  // Ordinal 62.
  void WriteRcRUnion(evo_test::RUnion const& value);

  // Ordinal 63.
  void WriteOptionalRecordWithChanges(std::optional<evo_test::RecordWithChanges> const& value);

  // Ordinal 64.
  void WriteAliasedOptionalRecordWithChanges(std::optional<evo_test::AliasedRecordWithChanges> const& value);

  // Ordinal 65.
  void WriteUnionRecordWithChanges(std::variant<evo_test::RecordWithChanges, int32_t> const& value);

  // Ordinal 66.
  void WriteUnionWithSameTypeset(std::variant<float, evo_test::RecordWithChanges, std::string, int32_t> const& value);

  // Ordinal 67.
  void WriteUnionWithTypesAdded(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value);

  // Ordinal 68.
  void WriteUnionWithTypesRemoved(std::variant<evo_test::RecordWithChanges, std::string> const& value);

  // Ordinal 69.
  void WriteRecordToOptional(std::optional<evo_test::RecordWithChanges> const& value);

  // Ordinal 70.
  void WriteRecordToAliasedOptional(evo_test::AliasedOptionalRecord const& value);

  // Ordinal 71.
  void WriteRecordToUnion(std::variant<evo_test::RecordWithChanges, std::string> const& value);

  // Ordinal 72.
  void WriteRecordToAliasedUnion(evo_test::AliasedRecordOrString const& value);

  // Ordinal 73.
  void WriteUnionToAliasedUnion(evo_test::AliasedRecordOrInt const& value);

  // Ordinal 74.
  void WriteUnionToAliasedUnionWithChanges(evo_test::AliasedRecordOrString const& value);

  // Ordinal 75.
  void WriteOptionalToAliasedOptional(evo_test::AliasedOptionalRecord const& value);

  // Ordinal 76.
  void WriteOptionalToAliasedOptionalWithChanges(evo_test::AliasedOptionalString const& value);

  // Ordinal 77.
  void WriteGenericRecord(evo_test::GenericRecord<int32_t, std::string> const& value);

  // Ordinal 78.
  void WriteGenericRecordToOpenAlias(evo_test::AliasedOpenGenericRecord<int32_t, std::string> const& value);

  // Ordinal 79.
  void WriteGenericRecordToClosedAlias(evo_test::AliasedClosedGenericRecord const& value);

  // Ordinal 80.
  void WriteGenericRecordToHalfClosedAlias(evo_test::AliasedHalfClosedGenericRecord<int32_t> const& value);

  // Ordinal 81.
  void WriteAliasedGenericRecordToAlias(evo_test::AliasedOpenGenericRecord<int32_t, std::string> const& value);

  // Ordinal 82.
  void WriteClosedGenericRecordToUnion(std::variant<evo_test::GenericRecord<int32_t, std::string>, std::string> const& value);

  // Ordinal 83.
  void WriteGenericRecordToAliasedUnion(evo_test::AliasedGenericRecordOrString const& value);

  // Ordinal 84.
  void WriteGenericUnionOfChangedRecord(evo_test::GenericUnion<evo_test::GenericRecord<int32_t, std::string>, float> const& value);

  // Ordinal 85.
  void WriteGenericParentRecord(evo_test::GenericParentRecord<int32_t> const& value);

  // Ordinal 86.
  void WriteGenericNestedRecords(evo_test::GenericRecord<evo_test::Unchanged, evo_test::Changed> const& value);

  // Ordinal 87.
  // Call this method for each element of the `genericRecordStream` stream, then call `EndGenericRecordStream() when done.`
  void WriteGenericRecordStream(evo_test::AliasedClosedGenericRecord const& value);

  // Ordinal 87.
  // Call this method to write many values to the `genericRecordStream` stream, then call `EndGenericRecordStream()` when done.
  void WriteGenericRecordStream(std::vector<evo_test::AliasedClosedGenericRecord> const& values);

  // Marks the end of the `genericRecordStream` stream.
  void EndGenericRecordStream();

  // Ordinal 88.
  // Call this method for each element of the `genericParentRecordStream` stream, then call `EndGenericParentRecordStream() when done.`
  void WriteGenericParentRecordStream(evo_test::GenericParentRecord<int32_t> const& value);

  // Ordinal 88.
  // Call this method to write many values to the `genericParentRecordStream` stream, then call `EndGenericParentRecordStream()` when done.
  void WriteGenericParentRecordStream(std::vector<evo_test::GenericParentRecord<int32_t>> const& values);

  // Marks the end of the `genericParentRecordStream` stream.
  void EndGenericParentRecordStream();

  // Ordinal 89.
  void WriteVectorRecordWithChanges(std::vector<evo_test::RecordWithChanges> const& value);

  // Ordinal 90.
  // Call this method for each element of the `streamedRecordWithChanges` stream, then call `EndStreamedRecordWithChanges() when done.`
  void WriteStreamedRecordWithChanges(evo_test::RecordWithChanges const& value);

  // Ordinal 90.
  // Call this method to write many values to the `streamedRecordWithChanges` stream, then call `EndStreamedRecordWithChanges()` when done.
  void WriteStreamedRecordWithChanges(std::vector<evo_test::RecordWithChanges> const& values);

  // Marks the end of the `streamedRecordWithChanges` stream.
  void EndStreamedRecordWithChanges();

  // Ordinal 91.
  void WriteAddedOptional(std::optional<evo_test::RecordWithChanges> const& value);

  // Ordinal 92.
  void WriteAddedMap(std::unordered_map<std::string, std::string> const& value);

  // Ordinal 93.
  // Call this method for each element of the `addedRecordStream` stream, then call `EndAddedRecordStream() when done.`
  void WriteAddedRecordStream(evo_test::RecordWithChanges const& value);

  // Ordinal 93.
  // Call this method to write many values to the `addedRecordStream` stream, then call `EndAddedRecordStream()` when done.
  void WriteAddedRecordStream(std::vector<evo_test::RecordWithChanges> const& values);

  // Marks the end of the `addedRecordStream` stream.
  void EndAddedRecordStream();

  // Optionaly close this writer before destructing. Validates that all steps were completed.
  void Close();

  virtual ~ProtocolWithChangesWriterBase() = default;

  // Flushes all buffered data.
  virtual void Flush() {}

  protected:
  virtual void WriteInt8ToIntImpl(int32_t const& value) = 0;
  virtual void WriteInt8ToLongImpl(int64_t const& value) = 0;
  virtual void WriteInt8ToUintImpl(uint32_t const& value) = 0;
  virtual void WriteInt8ToUlongImpl(uint64_t const& value) = 0;
  virtual void WriteInt8ToFloatImpl(float const& value) = 0;
  virtual void WriteInt8ToDoubleImpl(double const& value) = 0;
  virtual void WriteIntToUintImpl(uint32_t const& value) = 0;
  virtual void WriteIntToLongImpl(int64_t const& value) = 0;
  virtual void WriteIntToFloatImpl(float const& value) = 0;
  virtual void WriteIntToDoubleImpl(double const& value) = 0;
  virtual void WriteUintToUlongImpl(uint64_t const& value) = 0;
  virtual void WriteUintToFloatImpl(float const& value) = 0;
  virtual void WriteUintToDoubleImpl(double const& value) = 0;
  virtual void WriteFloatToDoubleImpl(double const& value) = 0;
  virtual void WriteIntToStringImpl(std::string const& value) = 0;
  virtual void WriteUintToStringImpl(std::string const& value) = 0;
  virtual void WriteLongToStringImpl(std::string const& value) = 0;
  virtual void WriteUlongToStringImpl(std::string const& value) = 0;
  virtual void WriteFloatToStringImpl(std::string const& value) = 0;
  virtual void WriteDoubleToStringImpl(std::string const& value) = 0;
  virtual void WriteIntToOptionalImpl(std::optional<int32_t> const& value) = 0;
  virtual void WriteFloatToOptionalImpl(std::optional<float> const& value) = 0;
  virtual void WriteStringToOptionalImpl(std::optional<std::string> const& value) = 0;
  virtual void WriteIntToUnionImpl(std::variant<int32_t, bool> const& value) = 0;
  virtual void WriteFloatToUnionImpl(std::variant<float, bool> const& value) = 0;
  virtual void WriteStringToUnionImpl(std::variant<std::string, bool> const& value) = 0;
  virtual void WriteOptionalIntToFloatImpl(std::optional<float> const& value) = 0;
  virtual void WriteOptionalFloatToStringImpl(std::optional<std::string> const& value) = 0;
  virtual void WriteAliasedLongToStringImpl(evo_test::AliasedLongToString const& value) = 0;
  virtual void WriteStringToAliasedStringImpl(evo_test::AliasedString const& value) = 0;
  virtual void WriteStringToAliasedIntImpl(evo_test::AliasedInt const& value) = 0;
  virtual void WriteOptionalIntToUnionImpl(std::variant<std::monostate, int32_t, std::string> const& value) = 0;
  virtual void WriteOptionalRecordToUnionImpl(std::variant<std::monostate, evo_test::RecordWithChanges, std::string> const& value) = 0;
  virtual void WriteRecordWithChangesImpl(evo_test::RecordWithChanges const& value) = 0;
  virtual void WriteAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges const& value) = 0;
  virtual void WriteRecordToRenamedRecordImpl(evo_test::RenamedRecord const& value) = 0;
  virtual void WriteRecordToAliasedRecordImpl(evo_test::AliasedRecordWithChanges const& value) = 0;
  virtual void WriteRecordToAliasedAliasImpl(evo_test::AliasOfAliasedRecordWithChanges const& value) = 0;
  virtual void WriteStreamOfAliasTypeChangeImpl(evo_test::StreamItem const& value) = 0;
  virtual void WriteStreamOfAliasTypeChangeImpl(std::vector<evo_test::StreamItem> const& value);
  virtual void EndStreamOfAliasTypeChangeImpl() = 0;
  virtual void WriteRlinkImpl(evo_test::RLink const& value) = 0;
  virtual void WriteRlinkRXImpl(evo_test::RX const& value) = 0;
  virtual void WriteRlinkRYImpl(evo_test::RY const& value) = 0;
  virtual void WriteRlinkRZImpl(evo_test::RZ const& value) = 0;
  virtual void WriteRaRLinkImpl(evo_test::RLink const& value) = 0;
  virtual void WriteRaRXImpl(evo_test::RX const& value) = 0;
  virtual void WriteRaRYImpl(evo_test::RY const& value) = 0;
  virtual void WriteRaRZImpl(evo_test::RZ const& value) = 0;
  virtual void WriteRbRLinkImpl(evo_test::RLink const& value) = 0;
  virtual void WriteRbRXImpl(evo_test::RX const& value) = 0;
  virtual void WriteRbRYImpl(evo_test::RY const& value) = 0;
  virtual void WriteRbRZImpl(evo_test::RZ const& value) = 0;
  virtual void WriteRcRLinkImpl(evo_test::RLink const& value) = 0;
  virtual void WriteRcRXImpl(evo_test::RX const& value) = 0;
  virtual void WriteRcRYImpl(evo_test::RY const& value) = 0;
  virtual void WriteRcRZImpl(evo_test::RZ const& value) = 0;
  virtual void WriteRlinkRNewImpl(evo_test::RNew const& value) = 0;
  virtual void WriteRaRNewImpl(evo_test::RNew const& value) = 0;
  virtual void WriteRbRNewImpl(evo_test::RNew const& value) = 0;
  virtual void WriteRcRNewImpl(evo_test::RNew const& value) = 0;
  virtual void WriteRlinkRUnionImpl(evo_test::RUnion const& value) = 0;
  virtual void WriteRaRUnionImpl(evo_test::RUnion const& value) = 0;
  virtual void WriteRbRUnionImpl(evo_test::RUnion const& value) = 0;
  virtual void WriteRcRUnionImpl(evo_test::RUnion const& value) = 0;
  virtual void WriteOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges> const& value) = 0;
  virtual void WriteAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges> const& value) = 0;
  virtual void WriteUnionRecordWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t> const& value) = 0;
  virtual void WriteUnionWithSameTypesetImpl(std::variant<float, evo_test::RecordWithChanges, std::string, int32_t> const& value) = 0;
  virtual void WriteUnionWithTypesAddedImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value) = 0;
  virtual void WriteUnionWithTypesRemovedImpl(std::variant<evo_test::RecordWithChanges, std::string> const& value) = 0;
  virtual void WriteRecordToOptionalImpl(std::optional<evo_test::RecordWithChanges> const& value) = 0;
  virtual void WriteRecordToAliasedOptionalImpl(evo_test::AliasedOptionalRecord const& value) = 0;
  virtual void WriteRecordToUnionImpl(std::variant<evo_test::RecordWithChanges, std::string> const& value) = 0;
  virtual void WriteRecordToAliasedUnionImpl(evo_test::AliasedRecordOrString const& value) = 0;
  virtual void WriteUnionToAliasedUnionImpl(evo_test::AliasedRecordOrInt const& value) = 0;
  virtual void WriteUnionToAliasedUnionWithChangesImpl(evo_test::AliasedRecordOrString const& value) = 0;
  virtual void WriteOptionalToAliasedOptionalImpl(evo_test::AliasedOptionalRecord const& value) = 0;
  virtual void WriteOptionalToAliasedOptionalWithChangesImpl(evo_test::AliasedOptionalString const& value) = 0;
  virtual void WriteGenericRecordImpl(evo_test::GenericRecord<int32_t, std::string> const& value) = 0;
  virtual void WriteGenericRecordToOpenAliasImpl(evo_test::AliasedOpenGenericRecord<int32_t, std::string> const& value) = 0;
  virtual void WriteGenericRecordToClosedAliasImpl(evo_test::AliasedClosedGenericRecord const& value) = 0;
  virtual void WriteGenericRecordToHalfClosedAliasImpl(evo_test::AliasedHalfClosedGenericRecord<int32_t> const& value) = 0;
  virtual void WriteAliasedGenericRecordToAliasImpl(evo_test::AliasedOpenGenericRecord<int32_t, std::string> const& value) = 0;
  virtual void WriteClosedGenericRecordToUnionImpl(std::variant<evo_test::GenericRecord<int32_t, std::string>, std::string> const& value) = 0;
  virtual void WriteGenericRecordToAliasedUnionImpl(evo_test::AliasedGenericRecordOrString const& value) = 0;
  virtual void WriteGenericUnionOfChangedRecordImpl(evo_test::GenericUnion<evo_test::GenericRecord<int32_t, std::string>, float> const& value) = 0;
  virtual void WriteGenericParentRecordImpl(evo_test::GenericParentRecord<int32_t> const& value) = 0;
  virtual void WriteGenericNestedRecordsImpl(evo_test::GenericRecord<evo_test::Unchanged, evo_test::Changed> const& value) = 0;
  virtual void WriteGenericRecordStreamImpl(evo_test::AliasedClosedGenericRecord const& value) = 0;
  virtual void WriteGenericRecordStreamImpl(std::vector<evo_test::AliasedClosedGenericRecord> const& value);
  virtual void EndGenericRecordStreamImpl() = 0;
  virtual void WriteGenericParentRecordStreamImpl(evo_test::GenericParentRecord<int32_t> const& value) = 0;
  virtual void WriteGenericParentRecordStreamImpl(std::vector<evo_test::GenericParentRecord<int32_t>> const& value);
  virtual void EndGenericParentRecordStreamImpl() = 0;
  virtual void WriteVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges> const& value) = 0;
  virtual void WriteStreamedRecordWithChangesImpl(evo_test::RecordWithChanges const& value) = 0;
  virtual void WriteStreamedRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges> const& value);
  virtual void EndStreamedRecordWithChangesImpl() = 0;
  virtual void WriteAddedOptionalImpl(std::optional<evo_test::RecordWithChanges> const& value) = 0;
  virtual void WriteAddedMapImpl(std::unordered_map<std::string, std::string> const& value) = 0;
  virtual void WriteAddedRecordStreamImpl(evo_test::RecordWithChanges const& value) = 0;
  virtual void WriteAddedRecordStreamImpl(std::vector<evo_test::RecordWithChanges> const& value);
  virtual void EndAddedRecordStreamImpl() = 0;
  virtual void CloseImpl() {}

  static std::string schema_;

  static std::vector<std::string> previous_schemas_;

  static std::string SchemaFromVersion(Version version);

  private:
  uint8_t state_ = 0;

  friend class ProtocolWithChangesReaderBase;
};

// Abstract reader for the ProtocolWithChanges protocol.
class ProtocolWithChangesReaderBase {
  public:
  // Ordinal 0.
  void ReadInt8ToInt(int32_t& value);

  // Ordinal 1.
  void ReadInt8ToLong(int64_t& value);

  // Ordinal 2.
  void ReadInt8ToUint(uint32_t& value);

  // Ordinal 3.
  void ReadInt8ToUlong(uint64_t& value);

  // Ordinal 4.
  void ReadInt8ToFloat(float& value);

  // Ordinal 5.
  void ReadInt8ToDouble(double& value);

  // Ordinal 6.
  void ReadIntToUint(uint32_t& value);

  // Ordinal 7.
  void ReadIntToLong(int64_t& value);

  // Ordinal 8.
  void ReadIntToFloat(float& value);

  // Ordinal 9.
  void ReadIntToDouble(double& value);

  // Ordinal 10.
  void ReadUintToUlong(uint64_t& value);

  // Ordinal 11.
  void ReadUintToFloat(float& value);

  // Ordinal 12.
  void ReadUintToDouble(double& value);

  // Ordinal 13.
  void ReadFloatToDouble(double& value);

  // Ordinal 14.
  void ReadIntToString(std::string& value);

  // Ordinal 15.
  void ReadUintToString(std::string& value);

  // Ordinal 16.
  void ReadLongToString(std::string& value);

  // Ordinal 17.
  void ReadUlongToString(std::string& value);

  // Ordinal 18.
  void ReadFloatToString(std::string& value);

  // Ordinal 19.
  void ReadDoubleToString(std::string& value);

  // Ordinal 20.
  void ReadIntToOptional(std::optional<int32_t>& value);

  // Ordinal 21.
  void ReadFloatToOptional(std::optional<float>& value);

  // Ordinal 22.
  void ReadStringToOptional(std::optional<std::string>& value);

  // Ordinal 23.
  void ReadIntToUnion(std::variant<int32_t, bool>& value);

  // Ordinal 24.
  void ReadFloatToUnion(std::variant<float, bool>& value);

  // Ordinal 25.
  void ReadStringToUnion(std::variant<std::string, bool>& value);

  // Ordinal 26.
  void ReadOptionalIntToFloat(std::optional<float>& value);

  // Ordinal 27.
  void ReadOptionalFloatToString(std::optional<std::string>& value);

  // Ordinal 28.
  void ReadAliasedLongToString(evo_test::AliasedLongToString& value);

  // Ordinal 29.
  void ReadStringToAliasedString(evo_test::AliasedString& value);

  // Ordinal 30.
  void ReadStringToAliasedInt(evo_test::AliasedInt& value);

  // Ordinal 31.
  void ReadOptionalIntToUnion(std::variant<std::monostate, int32_t, std::string>& value);

  // Ordinal 32.
  void ReadOptionalRecordToUnion(std::variant<std::monostate, evo_test::RecordWithChanges, std::string>& value);

  // Ordinal 33.
  void ReadRecordWithChanges(evo_test::RecordWithChanges& value);

  // Ordinal 34.
  void ReadAliasedRecordWithChanges(evo_test::AliasedRecordWithChanges& value);

  // Ordinal 35.
  void ReadRecordToRenamedRecord(evo_test::RenamedRecord& value);

  // Ordinal 36.
  void ReadRecordToAliasedRecord(evo_test::AliasedRecordWithChanges& value);

  // Ordinal 37.
  void ReadRecordToAliasedAlias(evo_test::AliasOfAliasedRecordWithChanges& value);

  // Ordinal 38.
  [[nodiscard]] bool ReadStreamOfAliasTypeChange(evo_test::StreamItem& value);

  // Ordinal 38.
  [[nodiscard]] bool ReadStreamOfAliasTypeChange(std::vector<evo_test::StreamItem>& values);

  // Ordinal 39.
  // Comprehensive NamedType changes
  void ReadRlink(evo_test::RLink& value);

  // Ordinal 40.
  void ReadRlinkRX(evo_test::RX& value);

  // Ordinal 41.
  void ReadRlinkRY(evo_test::RY& value);

  // Ordinal 42.
  void ReadRlinkRZ(evo_test::RZ& value);

  // Ordinal 43.
  void ReadRaRLink(evo_test::RLink& value);

  // Ordinal 44.
  void ReadRaRX(evo_test::RX& value);

  // Ordinal 45.
  void ReadRaRY(evo_test::RY& value);

  // Ordinal 46.
  void ReadRaRZ(evo_test::RZ& value);

  // Ordinal 47.
  void ReadRbRLink(evo_test::RLink& value);

  // Ordinal 48.
  void ReadRbRX(evo_test::RX& value);

  // Ordinal 49.
  void ReadRbRY(evo_test::RY& value);

  // Ordinal 50.
  void ReadRbRZ(evo_test::RZ& value);

  // Ordinal 51.
  void ReadRcRLink(evo_test::RLink& value);

  // Ordinal 52.
  void ReadRcRX(evo_test::RX& value);

  // Ordinal 53.
  void ReadRcRY(evo_test::RY& value);

  // Ordinal 54.
  void ReadRcRZ(evo_test::RZ& value);

  // Ordinal 55.
  void ReadRlinkRNew(evo_test::RNew& value);

  // Ordinal 56.
  void ReadRaRNew(evo_test::RNew& value);

  // Ordinal 57.
  void ReadRbRNew(evo_test::RNew& value);

  // Ordinal 58.
  void ReadRcRNew(evo_test::RNew& value);

  // Ordinal 59.
  void ReadRlinkRUnion(evo_test::RUnion& value);

  // Ordinal 60.
  void ReadRaRUnion(evo_test::RUnion& value);

  // Ordinal 61.
  void ReadRbRUnion(evo_test::RUnion& value);

  // Ordinal 62.
  void ReadRcRUnion(evo_test::RUnion& value);

  // Ordinal 63.
  void ReadOptionalRecordWithChanges(std::optional<evo_test::RecordWithChanges>& value);

  // Ordinal 64.
  void ReadAliasedOptionalRecordWithChanges(std::optional<evo_test::AliasedRecordWithChanges>& value);

  // Ordinal 65.
  void ReadUnionRecordWithChanges(std::variant<evo_test::RecordWithChanges, int32_t>& value);

  // Ordinal 66.
  void ReadUnionWithSameTypeset(std::variant<float, evo_test::RecordWithChanges, std::string, int32_t>& value);

  // Ordinal 67.
  void ReadUnionWithTypesAdded(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value);

  // Ordinal 68.
  void ReadUnionWithTypesRemoved(std::variant<evo_test::RecordWithChanges, std::string>& value);

  // Ordinal 69.
  void ReadRecordToOptional(std::optional<evo_test::RecordWithChanges>& value);

  // Ordinal 70.
  void ReadRecordToAliasedOptional(evo_test::AliasedOptionalRecord& value);

  // Ordinal 71.
  void ReadRecordToUnion(std::variant<evo_test::RecordWithChanges, std::string>& value);

  // Ordinal 72.
  void ReadRecordToAliasedUnion(evo_test::AliasedRecordOrString& value);

  // Ordinal 73.
  void ReadUnionToAliasedUnion(evo_test::AliasedRecordOrInt& value);

  // Ordinal 74.
  void ReadUnionToAliasedUnionWithChanges(evo_test::AliasedRecordOrString& value);

  // Ordinal 75.
  void ReadOptionalToAliasedOptional(evo_test::AliasedOptionalRecord& value);

  // Ordinal 76.
  void ReadOptionalToAliasedOptionalWithChanges(evo_test::AliasedOptionalString& value);

  // Ordinal 77.
  void ReadGenericRecord(evo_test::GenericRecord<int32_t, std::string>& value);

  // Ordinal 78.
  void ReadGenericRecordToOpenAlias(evo_test::AliasedOpenGenericRecord<int32_t, std::string>& value);

  // Ordinal 79.
  void ReadGenericRecordToClosedAlias(evo_test::AliasedClosedGenericRecord& value);

  // Ordinal 80.
  void ReadGenericRecordToHalfClosedAlias(evo_test::AliasedHalfClosedGenericRecord<int32_t>& value);

  // Ordinal 81.
  void ReadAliasedGenericRecordToAlias(evo_test::AliasedOpenGenericRecord<int32_t, std::string>& value);

  // Ordinal 82.
  void ReadClosedGenericRecordToUnion(std::variant<evo_test::GenericRecord<int32_t, std::string>, std::string>& value);

  // Ordinal 83.
  void ReadGenericRecordToAliasedUnion(evo_test::AliasedGenericRecordOrString& value);

  // Ordinal 84.
  void ReadGenericUnionOfChangedRecord(evo_test::GenericUnion<evo_test::GenericRecord<int32_t, std::string>, float>& value);

  // Ordinal 85.
  void ReadGenericParentRecord(evo_test::GenericParentRecord<int32_t>& value);

  // Ordinal 86.
  void ReadGenericNestedRecords(evo_test::GenericRecord<evo_test::Unchanged, evo_test::Changed>& value);

  // Ordinal 87.
  [[nodiscard]] bool ReadGenericRecordStream(evo_test::AliasedClosedGenericRecord& value);

  // Ordinal 87.
  [[nodiscard]] bool ReadGenericRecordStream(std::vector<evo_test::AliasedClosedGenericRecord>& values);

  // Ordinal 88.
  [[nodiscard]] bool ReadGenericParentRecordStream(evo_test::GenericParentRecord<int32_t>& value);

  // Ordinal 88.
  [[nodiscard]] bool ReadGenericParentRecordStream(std::vector<evo_test::GenericParentRecord<int32_t>>& values);

  // Ordinal 89.
  void ReadVectorRecordWithChanges(std::vector<evo_test::RecordWithChanges>& value);

  // Ordinal 90.
  [[nodiscard]] bool ReadStreamedRecordWithChanges(evo_test::RecordWithChanges& value);

  // Ordinal 90.
  [[nodiscard]] bool ReadStreamedRecordWithChanges(std::vector<evo_test::RecordWithChanges>& values);

  // Ordinal 91.
  void ReadAddedOptional(std::optional<evo_test::RecordWithChanges>& value);

  // Ordinal 92.
  void ReadAddedMap(std::unordered_map<std::string, std::string>& value);

  // Ordinal 93.
  [[nodiscard]] bool ReadAddedRecordStream(evo_test::RecordWithChanges& value);

  // Ordinal 93.
  [[nodiscard]] bool ReadAddedRecordStream(std::vector<evo_test::RecordWithChanges>& values);

  // Optionaly close this writer before destructing. Validates that all steps were completely read.
  void Close();

  void CopyTo(ProtocolWithChangesWriterBase& writer, size_t stream_of_alias_type_change_buffer_size = 1, size_t generic_record_stream_buffer_size = 1, size_t generic_parent_record_stream_buffer_size = 1, size_t streamed_record_with_changes_buffer_size = 1, size_t added_record_stream_buffer_size = 1);

  virtual ~ProtocolWithChangesReaderBase() = default;

  protected:
  virtual void ReadInt8ToIntImpl(int32_t& value) = 0;
  virtual void ReadInt8ToLongImpl(int64_t& value) = 0;
  virtual void ReadInt8ToUintImpl(uint32_t& value) = 0;
  virtual void ReadInt8ToUlongImpl(uint64_t& value) = 0;
  virtual void ReadInt8ToFloatImpl(float& value) = 0;
  virtual void ReadInt8ToDoubleImpl(double& value) = 0;
  virtual void ReadIntToUintImpl(uint32_t& value) = 0;
  virtual void ReadIntToLongImpl(int64_t& value) = 0;
  virtual void ReadIntToFloatImpl(float& value) = 0;
  virtual void ReadIntToDoubleImpl(double& value) = 0;
  virtual void ReadUintToUlongImpl(uint64_t& value) = 0;
  virtual void ReadUintToFloatImpl(float& value) = 0;
  virtual void ReadUintToDoubleImpl(double& value) = 0;
  virtual void ReadFloatToDoubleImpl(double& value) = 0;
  virtual void ReadIntToStringImpl(std::string& value) = 0;
  virtual void ReadUintToStringImpl(std::string& value) = 0;
  virtual void ReadLongToStringImpl(std::string& value) = 0;
  virtual void ReadUlongToStringImpl(std::string& value) = 0;
  virtual void ReadFloatToStringImpl(std::string& value) = 0;
  virtual void ReadDoubleToStringImpl(std::string& value) = 0;
  virtual void ReadIntToOptionalImpl(std::optional<int32_t>& value) = 0;
  virtual void ReadFloatToOptionalImpl(std::optional<float>& value) = 0;
  virtual void ReadStringToOptionalImpl(std::optional<std::string>& value) = 0;
  virtual void ReadIntToUnionImpl(std::variant<int32_t, bool>& value) = 0;
  virtual void ReadFloatToUnionImpl(std::variant<float, bool>& value) = 0;
  virtual void ReadStringToUnionImpl(std::variant<std::string, bool>& value) = 0;
  virtual void ReadOptionalIntToFloatImpl(std::optional<float>& value) = 0;
  virtual void ReadOptionalFloatToStringImpl(std::optional<std::string>& value) = 0;
  virtual void ReadAliasedLongToStringImpl(evo_test::AliasedLongToString& value) = 0;
  virtual void ReadStringToAliasedStringImpl(evo_test::AliasedString& value) = 0;
  virtual void ReadStringToAliasedIntImpl(evo_test::AliasedInt& value) = 0;
  virtual void ReadOptionalIntToUnionImpl(std::variant<std::monostate, int32_t, std::string>& value) = 0;
  virtual void ReadOptionalRecordToUnionImpl(std::variant<std::monostate, evo_test::RecordWithChanges, std::string>& value) = 0;
  virtual void ReadRecordWithChangesImpl(evo_test::RecordWithChanges& value) = 0;
  virtual void ReadAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges& value) = 0;
  virtual void ReadRecordToRenamedRecordImpl(evo_test::RenamedRecord& value) = 0;
  virtual void ReadRecordToAliasedRecordImpl(evo_test::AliasedRecordWithChanges& value) = 0;
  virtual void ReadRecordToAliasedAliasImpl(evo_test::AliasOfAliasedRecordWithChanges& value) = 0;
  virtual bool ReadStreamOfAliasTypeChangeImpl(evo_test::StreamItem& value) = 0;
  virtual bool ReadStreamOfAliasTypeChangeImpl(std::vector<evo_test::StreamItem>& values);
  virtual void ReadRlinkImpl(evo_test::RLink& value) = 0;
  virtual void ReadRlinkRXImpl(evo_test::RX& value) = 0;
  virtual void ReadRlinkRYImpl(evo_test::RY& value) = 0;
  virtual void ReadRlinkRZImpl(evo_test::RZ& value) = 0;
  virtual void ReadRaRLinkImpl(evo_test::RLink& value) = 0;
  virtual void ReadRaRXImpl(evo_test::RX& value) = 0;
  virtual void ReadRaRYImpl(evo_test::RY& value) = 0;
  virtual void ReadRaRZImpl(evo_test::RZ& value) = 0;
  virtual void ReadRbRLinkImpl(evo_test::RLink& value) = 0;
  virtual void ReadRbRXImpl(evo_test::RX& value) = 0;
  virtual void ReadRbRYImpl(evo_test::RY& value) = 0;
  virtual void ReadRbRZImpl(evo_test::RZ& value) = 0;
  virtual void ReadRcRLinkImpl(evo_test::RLink& value) = 0;
  virtual void ReadRcRXImpl(evo_test::RX& value) = 0;
  virtual void ReadRcRYImpl(evo_test::RY& value) = 0;
  virtual void ReadRcRZImpl(evo_test::RZ& value) = 0;
  virtual void ReadRlinkRNewImpl(evo_test::RNew& value) = 0;
  virtual void ReadRaRNewImpl(evo_test::RNew& value) = 0;
  virtual void ReadRbRNewImpl(evo_test::RNew& value) = 0;
  virtual void ReadRcRNewImpl(evo_test::RNew& value) = 0;
  virtual void ReadRlinkRUnionImpl(evo_test::RUnion& value) = 0;
  virtual void ReadRaRUnionImpl(evo_test::RUnion& value) = 0;
  virtual void ReadRbRUnionImpl(evo_test::RUnion& value) = 0;
  virtual void ReadRcRUnionImpl(evo_test::RUnion& value) = 0;
  virtual void ReadOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges>& value) = 0;
  virtual void ReadAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges>& value) = 0;
  virtual void ReadUnionRecordWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t>& value) = 0;
  virtual void ReadUnionWithSameTypesetImpl(std::variant<float, evo_test::RecordWithChanges, std::string, int32_t>& value) = 0;
  virtual void ReadUnionWithTypesAddedImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value) = 0;
  virtual void ReadUnionWithTypesRemovedImpl(std::variant<evo_test::RecordWithChanges, std::string>& value) = 0;
  virtual void ReadRecordToOptionalImpl(std::optional<evo_test::RecordWithChanges>& value) = 0;
  virtual void ReadRecordToAliasedOptionalImpl(evo_test::AliasedOptionalRecord& value) = 0;
  virtual void ReadRecordToUnionImpl(std::variant<evo_test::RecordWithChanges, std::string>& value) = 0;
  virtual void ReadRecordToAliasedUnionImpl(evo_test::AliasedRecordOrString& value) = 0;
  virtual void ReadUnionToAliasedUnionImpl(evo_test::AliasedRecordOrInt& value) = 0;
  virtual void ReadUnionToAliasedUnionWithChangesImpl(evo_test::AliasedRecordOrString& value) = 0;
  virtual void ReadOptionalToAliasedOptionalImpl(evo_test::AliasedOptionalRecord& value) = 0;
  virtual void ReadOptionalToAliasedOptionalWithChangesImpl(evo_test::AliasedOptionalString& value) = 0;
  virtual void ReadGenericRecordImpl(evo_test::GenericRecord<int32_t, std::string>& value) = 0;
  virtual void ReadGenericRecordToOpenAliasImpl(evo_test::AliasedOpenGenericRecord<int32_t, std::string>& value) = 0;
  virtual void ReadGenericRecordToClosedAliasImpl(evo_test::AliasedClosedGenericRecord& value) = 0;
  virtual void ReadGenericRecordToHalfClosedAliasImpl(evo_test::AliasedHalfClosedGenericRecord<int32_t>& value) = 0;
  virtual void ReadAliasedGenericRecordToAliasImpl(evo_test::AliasedOpenGenericRecord<int32_t, std::string>& value) = 0;
  virtual void ReadClosedGenericRecordToUnionImpl(std::variant<evo_test::GenericRecord<int32_t, std::string>, std::string>& value) = 0;
  virtual void ReadGenericRecordToAliasedUnionImpl(evo_test::AliasedGenericRecordOrString& value) = 0;
  virtual void ReadGenericUnionOfChangedRecordImpl(evo_test::GenericUnion<evo_test::GenericRecord<int32_t, std::string>, float>& value) = 0;
  virtual void ReadGenericParentRecordImpl(evo_test::GenericParentRecord<int32_t>& value) = 0;
  virtual void ReadGenericNestedRecordsImpl(evo_test::GenericRecord<evo_test::Unchanged, evo_test::Changed>& value) = 0;
  virtual bool ReadGenericRecordStreamImpl(evo_test::AliasedClosedGenericRecord& value) = 0;
  virtual bool ReadGenericRecordStreamImpl(std::vector<evo_test::AliasedClosedGenericRecord>& values);
  virtual bool ReadGenericParentRecordStreamImpl(evo_test::GenericParentRecord<int32_t>& value) = 0;
  virtual bool ReadGenericParentRecordStreamImpl(std::vector<evo_test::GenericParentRecord<int32_t>>& values);
  virtual void ReadVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges>& value) = 0;
  virtual bool ReadStreamedRecordWithChangesImpl(evo_test::RecordWithChanges& value) = 0;
  virtual bool ReadStreamedRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges>& values);
  virtual void ReadAddedOptionalImpl(std::optional<evo_test::RecordWithChanges>& value) = 0;
  virtual void ReadAddedMapImpl(std::unordered_map<std::string, std::string>& value) = 0;
  virtual bool ReadAddedRecordStreamImpl(evo_test::RecordWithChanges& value) = 0;
  virtual bool ReadAddedRecordStreamImpl(std::vector<evo_test::RecordWithChanges>& values);
  virtual void CloseImpl() {}
  static std::string schema_;

  static std::vector<std::string> previous_schemas_;

  static Version VersionFromSchema(const std::string& schema);

  private:
  uint8_t state_ = 0;
};

// Abstract writer for the UnusedProtocol protocol.
class UnusedProtocolWriterBase {
  public:
  // Ordinal 0.
  // Call this method for each element of the `records` stream, then call `EndRecords() when done.`
  void WriteRecords(evo_test::UnchangedRecord const& value);

  // Ordinal 0.
  // Call this method to write many values to the `records` stream, then call `EndRecords()` when done.
  void WriteRecords(std::vector<evo_test::UnchangedRecord> const& values);

  // Marks the end of the `records` stream.
  void EndRecords();

  // Optionaly close this writer before destructing. Validates that all steps were completed.
  void Close();

  virtual ~UnusedProtocolWriterBase() = default;

  // Flushes all buffered data.
  virtual void Flush() {}

  protected:
  virtual void WriteRecordsImpl(evo_test::UnchangedRecord const& value) = 0;
  virtual void WriteRecordsImpl(std::vector<evo_test::UnchangedRecord> const& value);
  virtual void EndRecordsImpl() = 0;
  virtual void CloseImpl() {}

  static std::string schema_;

  static std::vector<std::string> previous_schemas_;

  static std::string SchemaFromVersion(Version version);

  private:
  uint8_t state_ = 0;

  friend class UnusedProtocolReaderBase;
};

// Abstract reader for the UnusedProtocol protocol.
class UnusedProtocolReaderBase {
  public:
  // Ordinal 0.
  [[nodiscard]] bool ReadRecords(evo_test::UnchangedRecord& value);

  // Ordinal 0.
  [[nodiscard]] bool ReadRecords(std::vector<evo_test::UnchangedRecord>& values);

  // Optionaly close this writer before destructing. Validates that all steps were completely read.
  void Close();

  void CopyTo(UnusedProtocolWriterBase& writer, size_t records_buffer_size = 1);

  virtual ~UnusedProtocolReaderBase() = default;

  protected:
  virtual bool ReadRecordsImpl(evo_test::UnchangedRecord& value) = 0;
  virtual bool ReadRecordsImpl(std::vector<evo_test::UnchangedRecord>& values);
  virtual void CloseImpl() {}
  static std::string schema_;

  static std::vector<std::string> previous_schemas_;

  static Version VersionFromSchema(const std::string& schema);

  private:
  uint8_t state_ = 0;
};
} // namespace evo_test
