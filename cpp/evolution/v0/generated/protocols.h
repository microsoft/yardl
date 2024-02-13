// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include "types.h"

namespace evo_test {
enum class Version {
  Current
};
// Abstract writer for the ProtocolWithChanges protocol.
class ProtocolWithChangesWriterBase {
  public:
  // Ordinal 0.
  void WriteInt8ToInt(int8_t const& value);

  // Ordinal 1.
  void WriteInt8ToLong(int8_t const& value);

  // Ordinal 2.
  void WriteInt8ToUint(int8_t const& value);

  // Ordinal 3.
  void WriteInt8ToUlong(int8_t const& value);

  // Ordinal 4.
  void WriteInt8ToFloat(int8_t const& value);

  // Ordinal 5.
  void WriteInt8ToDouble(int8_t const& value);

  // Ordinal 6.
  void WriteIntToUint(int32_t const& value);

  // Ordinal 7.
  void WriteIntToLong(int32_t const& value);

  // Ordinal 8.
  void WriteIntToFloat(int32_t const& value);

  // Ordinal 9.
  void WriteIntToDouble(int32_t const& value);

  // Ordinal 10.
  void WriteUintToUlong(uint32_t const& value);

  // Ordinal 11.
  void WriteUintToFloat(uint32_t const& value);

  // Ordinal 12.
  void WriteUintToDouble(uint32_t const& value);

  // Ordinal 13.
  void WriteFloatToDouble(float const& value);

  // Ordinal 14.
  void WriteIntToString(int32_t const& value);

  // Ordinal 15.
  void WriteUintToString(uint32_t const& value);

  // Ordinal 16.
  void WriteLongToString(int64_t const& value);

  // Ordinal 17.
  void WriteUlongToString(uint64_t const& value);

  // Ordinal 18.
  void WriteFloatToString(float const& value);

  // Ordinal 19.
  void WriteDoubleToString(double const& value);

  // Ordinal 20.
  void WriteIntToOptional(int32_t const& value);

  // Ordinal 21.
  void WriteFloatToOptional(float const& value);

  // Ordinal 22.
  void WriteStringToOptional(std::string const& value);

  // Ordinal 23.
  void WriteIntToUnion(int32_t const& value);

  // Ordinal 24.
  void WriteFloatToUnion(float const& value);

  // Ordinal 25.
  void WriteStringToUnion(std::string const& value);

  // Ordinal 26.
  void WriteOptionalIntToFloat(std::optional<int32_t> const& value);

  // Ordinal 27.
  void WriteOptionalFloatToString(std::optional<float> const& value);

  // Ordinal 28.
  void WriteAliasedLongToString(evo_test::AliasedLongToString const& value);

  // Ordinal 29.
  void WriteStringToAliasedString(std::string const& value);

  // Ordinal 30.
  void WriteStringToAliasedInt(std::string const& value);

  // Ordinal 31.
  void WriteOptionalIntToUnion(std::optional<int32_t> const& value);

  // Ordinal 32.
  void WriteOptionalRecordToUnion(std::optional<evo_test::RecordWithChanges> const& value);

  // Ordinal 33.
  void WriteRecordWithChanges(evo_test::RecordWithChanges const& value);

  // Ordinal 34.
  void WriteAliasedRecordWithChanges(evo_test::AliasedRecordWithChanges const& value);

  // Ordinal 35.
  void WriteRecordToRenamedRecord(evo_test::RenamedRecord const& value);

  // Ordinal 36.
  void WriteRecordToAliasedRecord(evo_test::RecordWithChanges const& value);

  // Ordinal 37.
  void WriteRecordToAliasedAlias(evo_test::RecordWithChanges const& value);

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
  void WriteRlinkRX(evo_test::RLink const& value);

  // Ordinal 41.
  void WriteRlinkRY(evo_test::RLink const& value);

  // Ordinal 42.
  void WriteRlinkRZ(evo_test::RLink const& value);

  // Ordinal 43.
  void WriteRaRLink(evo_test::RA const& value);

  // Ordinal 44.
  void WriteRaRX(evo_test::RA const& value);

  // Ordinal 45.
  void WriteRaRY(evo_test::RA const& value);

  // Ordinal 46.
  void WriteRaRZ(evo_test::RA const& value);

  // Ordinal 47.
  void WriteRbRLink(evo_test::RB const& value);

  // Ordinal 48.
  void WriteRbRX(evo_test::RB const& value);

  // Ordinal 49.
  void WriteRbRY(evo_test::RB const& value);

  // Ordinal 50.
  void WriteRbRZ(evo_test::RB const& value);

  // Ordinal 51.
  void WriteRcRLink(evo_test::RC const& value);

  // Ordinal 52.
  void WriteRcRX(evo_test::RC const& value);

  // Ordinal 53.
  void WriteRcRY(evo_test::RC const& value);

  // Ordinal 54.
  void WriteRcRZ(evo_test::RC const& value);

  // Ordinal 55.
  void WriteRlinkRNew(evo_test::RLink const& value);

  // Ordinal 56.
  void WriteRaRNew(evo_test::RA const& value);

  // Ordinal 57.
  void WriteRbRNew(evo_test::RB const& value);

  // Ordinal 58.
  void WriteRcRNew(evo_test::RC const& value);

  // Ordinal 59.
  void WriteRlinkRUnion(evo_test::RLink const& value);

  // Ordinal 60.
  void WriteRaRUnion(evo_test::RA const& value);

  // Ordinal 61.
  void WriteRbRUnion(evo_test::RB const& value);

  // Ordinal 62.
  void WriteRcRUnion(evo_test::RC const& value);

  // Ordinal 63.
  void WriteOptionalRecordWithChanges(std::optional<evo_test::RecordWithChanges> const& value);

  // Ordinal 64.
  void WriteAliasedOptionalRecordWithChanges(std::optional<evo_test::AliasedRecordWithChanges> const& value);

  // Ordinal 65.
  void WriteUnionRecordWithChanges(std::variant<evo_test::RecordWithChanges, int32_t> const& value);

  // Ordinal 66.
  void WriteUnionWithSameTypeset(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value);

  // Ordinal 67.
  void WriteUnionWithTypesAdded(std::variant<evo_test::RecordWithChanges, float> const& value);

  // Ordinal 68.
  void WriteUnionWithTypesRemoved(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value);

  // Ordinal 69.
  void WriteRecordToOptional(evo_test::RecordWithChanges const& value);

  // Ordinal 70.
  void WriteRecordToAliasedOptional(evo_test::RecordWithChanges const& value);

  // Ordinal 71.
  void WriteRecordToUnion(evo_test::RecordWithChanges const& value);

  // Ordinal 72.
  void WriteRecordToAliasedUnion(evo_test::RecordWithChanges const& value);

  // Ordinal 73.
  void WriteUnionToAliasedUnion(std::variant<evo_test::RecordWithChanges, int32_t> const& value);

  // Ordinal 74.
  void WriteUnionToAliasedUnionWithChanges(std::variant<evo_test::RecordWithChanges, int32_t> const& value);

  // Ordinal 75.
  void WriteOptionalToAliasedOptional(std::optional<evo_test::RecordWithChanges> const& value);

  // Ordinal 76.
  void WriteOptionalToAliasedOptionalWithChanges(std::optional<int32_t> const& value);

  // Ordinal 77.
  void WriteGenericRecord(evo_test::GenericRecord<int32_t, std::string> const& value);

  // Ordinal 78.
  void WriteGenericRecordToOpenAlias(evo_test::GenericRecord<int32_t, std::string> const& value);

  // Ordinal 79.
  void WriteGenericRecordToClosedAlias(evo_test::GenericRecord<int32_t, std::string> const& value);

  // Ordinal 80.
  void WriteGenericRecordToHalfClosedAlias(evo_test::GenericRecord<int32_t, std::string> const& value);

  // Ordinal 81.
  void WriteAliasedGenericRecordToAlias(evo_test::AliasedHalfClosedGenericRecord<int32_t> const& value);

  // Ordinal 82.
  void WriteClosedGenericRecordToUnion(evo_test::AliasedClosedGenericRecord const& value);

  // Ordinal 83.
  void WriteGenericRecordToAliasedUnion(evo_test::GenericRecord<int32_t, std::string> const& value);

  // Ordinal 84.
  void WriteGenericUnionOfChangedRecord(evo_test::AliasedClosedGenericUnion const& value);

  // Ordinal 85.
  void WriteGenericParentRecord(evo_test::GenericParentRecord<int32_t> const& value);

  // Ordinal 86.
  void WriteGenericNestedRecords(evo_test::GenericRecord<evo_test::UnchangedGeneric<int32_t>, evo_test::ChangedGeneric<std::string, int32_t>> const& value);

  // Ordinal 87.
  // Call this method for each element of the `genericRecordStream` stream, then call `EndGenericRecordStream() when done.`
  void WriteGenericRecordStream(evo_test::GenericRecord<int32_t, std::string> const& value);

  // Ordinal 87.
  // Call this method to write many values to the `genericRecordStream` stream, then call `EndGenericRecordStream()` when done.
  void WriteGenericRecordStream(std::vector<evo_test::GenericRecord<int32_t, std::string>> const& values);

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

  // Optionaly close this writer before destructing. Validates that all steps were completed.
  void Close();

  virtual ~ProtocolWithChangesWriterBase() = default;

  // Flushes all buffered data.
  virtual void Flush() {}

  protected:
  virtual void WriteInt8ToIntImpl(int8_t const& value) = 0;
  virtual void WriteInt8ToLongImpl(int8_t const& value) = 0;
  virtual void WriteInt8ToUintImpl(int8_t const& value) = 0;
  virtual void WriteInt8ToUlongImpl(int8_t const& value) = 0;
  virtual void WriteInt8ToFloatImpl(int8_t const& value) = 0;
  virtual void WriteInt8ToDoubleImpl(int8_t const& value) = 0;
  virtual void WriteIntToUintImpl(int32_t const& value) = 0;
  virtual void WriteIntToLongImpl(int32_t const& value) = 0;
  virtual void WriteIntToFloatImpl(int32_t const& value) = 0;
  virtual void WriteIntToDoubleImpl(int32_t const& value) = 0;
  virtual void WriteUintToUlongImpl(uint32_t const& value) = 0;
  virtual void WriteUintToFloatImpl(uint32_t const& value) = 0;
  virtual void WriteUintToDoubleImpl(uint32_t const& value) = 0;
  virtual void WriteFloatToDoubleImpl(float const& value) = 0;
  virtual void WriteIntToStringImpl(int32_t const& value) = 0;
  virtual void WriteUintToStringImpl(uint32_t const& value) = 0;
  virtual void WriteLongToStringImpl(int64_t const& value) = 0;
  virtual void WriteUlongToStringImpl(uint64_t const& value) = 0;
  virtual void WriteFloatToStringImpl(float const& value) = 0;
  virtual void WriteDoubleToStringImpl(double const& value) = 0;
  virtual void WriteIntToOptionalImpl(int32_t const& value) = 0;
  virtual void WriteFloatToOptionalImpl(float const& value) = 0;
  virtual void WriteStringToOptionalImpl(std::string const& value) = 0;
  virtual void WriteIntToUnionImpl(int32_t const& value) = 0;
  virtual void WriteFloatToUnionImpl(float const& value) = 0;
  virtual void WriteStringToUnionImpl(std::string const& value) = 0;
  virtual void WriteOptionalIntToFloatImpl(std::optional<int32_t> const& value) = 0;
  virtual void WriteOptionalFloatToStringImpl(std::optional<float> const& value) = 0;
  virtual void WriteAliasedLongToStringImpl(evo_test::AliasedLongToString const& value) = 0;
  virtual void WriteStringToAliasedStringImpl(std::string const& value) = 0;
  virtual void WriteStringToAliasedIntImpl(std::string const& value) = 0;
  virtual void WriteOptionalIntToUnionImpl(std::optional<int32_t> const& value) = 0;
  virtual void WriteOptionalRecordToUnionImpl(std::optional<evo_test::RecordWithChanges> const& value) = 0;
  virtual void WriteRecordWithChangesImpl(evo_test::RecordWithChanges const& value) = 0;
  virtual void WriteAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges const& value) = 0;
  virtual void WriteRecordToRenamedRecordImpl(evo_test::RenamedRecord const& value) = 0;
  virtual void WriteRecordToAliasedRecordImpl(evo_test::RecordWithChanges const& value) = 0;
  virtual void WriteRecordToAliasedAliasImpl(evo_test::RecordWithChanges const& value) = 0;
  virtual void WriteStreamOfAliasTypeChangeImpl(evo_test::StreamItem const& value) = 0;
  virtual void WriteStreamOfAliasTypeChangeImpl(std::vector<evo_test::StreamItem> const& value);
  virtual void EndStreamOfAliasTypeChangeImpl() = 0;
  virtual void WriteRlinkImpl(evo_test::RLink const& value) = 0;
  virtual void WriteRlinkRXImpl(evo_test::RLink const& value) = 0;
  virtual void WriteRlinkRYImpl(evo_test::RLink const& value) = 0;
  virtual void WriteRlinkRZImpl(evo_test::RLink const& value) = 0;
  virtual void WriteRaRLinkImpl(evo_test::RA const& value) = 0;
  virtual void WriteRaRXImpl(evo_test::RA const& value) = 0;
  virtual void WriteRaRYImpl(evo_test::RA const& value) = 0;
  virtual void WriteRaRZImpl(evo_test::RA const& value) = 0;
  virtual void WriteRbRLinkImpl(evo_test::RB const& value) = 0;
  virtual void WriteRbRXImpl(evo_test::RB const& value) = 0;
  virtual void WriteRbRYImpl(evo_test::RB const& value) = 0;
  virtual void WriteRbRZImpl(evo_test::RB const& value) = 0;
  virtual void WriteRcRLinkImpl(evo_test::RC const& value) = 0;
  virtual void WriteRcRXImpl(evo_test::RC const& value) = 0;
  virtual void WriteRcRYImpl(evo_test::RC const& value) = 0;
  virtual void WriteRcRZImpl(evo_test::RC const& value) = 0;
  virtual void WriteRlinkRNewImpl(evo_test::RLink const& value) = 0;
  virtual void WriteRaRNewImpl(evo_test::RA const& value) = 0;
  virtual void WriteRbRNewImpl(evo_test::RB const& value) = 0;
  virtual void WriteRcRNewImpl(evo_test::RC const& value) = 0;
  virtual void WriteRlinkRUnionImpl(evo_test::RLink const& value) = 0;
  virtual void WriteRaRUnionImpl(evo_test::RA const& value) = 0;
  virtual void WriteRbRUnionImpl(evo_test::RB const& value) = 0;
  virtual void WriteRcRUnionImpl(evo_test::RC const& value) = 0;
  virtual void WriteOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges> const& value) = 0;
  virtual void WriteAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges> const& value) = 0;
  virtual void WriteUnionRecordWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t> const& value) = 0;
  virtual void WriteUnionWithSameTypesetImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value) = 0;
  virtual void WriteUnionWithTypesAddedImpl(std::variant<evo_test::RecordWithChanges, float> const& value) = 0;
  virtual void WriteUnionWithTypesRemovedImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value) = 0;
  virtual void WriteRecordToOptionalImpl(evo_test::RecordWithChanges const& value) = 0;
  virtual void WriteRecordToAliasedOptionalImpl(evo_test::RecordWithChanges const& value) = 0;
  virtual void WriteRecordToUnionImpl(evo_test::RecordWithChanges const& value) = 0;
  virtual void WriteRecordToAliasedUnionImpl(evo_test::RecordWithChanges const& value) = 0;
  virtual void WriteUnionToAliasedUnionImpl(std::variant<evo_test::RecordWithChanges, int32_t> const& value) = 0;
  virtual void WriteUnionToAliasedUnionWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t> const& value) = 0;
  virtual void WriteOptionalToAliasedOptionalImpl(std::optional<evo_test::RecordWithChanges> const& value) = 0;
  virtual void WriteOptionalToAliasedOptionalWithChangesImpl(std::optional<int32_t> const& value) = 0;
  virtual void WriteGenericRecordImpl(evo_test::GenericRecord<int32_t, std::string> const& value) = 0;
  virtual void WriteGenericRecordToOpenAliasImpl(evo_test::GenericRecord<int32_t, std::string> const& value) = 0;
  virtual void WriteGenericRecordToClosedAliasImpl(evo_test::GenericRecord<int32_t, std::string> const& value) = 0;
  virtual void WriteGenericRecordToHalfClosedAliasImpl(evo_test::GenericRecord<int32_t, std::string> const& value) = 0;
  virtual void WriteAliasedGenericRecordToAliasImpl(evo_test::AliasedHalfClosedGenericRecord<int32_t> const& value) = 0;
  virtual void WriteClosedGenericRecordToUnionImpl(evo_test::AliasedClosedGenericRecord const& value) = 0;
  virtual void WriteGenericRecordToAliasedUnionImpl(evo_test::GenericRecord<int32_t, std::string> const& value) = 0;
  virtual void WriteGenericUnionOfChangedRecordImpl(evo_test::AliasedClosedGenericUnion const& value) = 0;
  virtual void WriteGenericParentRecordImpl(evo_test::GenericParentRecord<int32_t> const& value) = 0;
  virtual void WriteGenericNestedRecordsImpl(evo_test::GenericRecord<evo_test::UnchangedGeneric<int32_t>, evo_test::ChangedGeneric<std::string, int32_t>> const& value) = 0;
  virtual void WriteGenericRecordStreamImpl(evo_test::GenericRecord<int32_t, std::string> const& value) = 0;
  virtual void WriteGenericRecordStreamImpl(std::vector<evo_test::GenericRecord<int32_t, std::string>> const& value);
  virtual void EndGenericRecordStreamImpl() = 0;
  virtual void WriteGenericParentRecordStreamImpl(evo_test::GenericParentRecord<int32_t> const& value) = 0;
  virtual void WriteGenericParentRecordStreamImpl(std::vector<evo_test::GenericParentRecord<int32_t>> const& value);
  virtual void EndGenericParentRecordStreamImpl() = 0;
  virtual void WriteVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges> const& value) = 0;
  virtual void WriteStreamedRecordWithChangesImpl(evo_test::RecordWithChanges const& value) = 0;
  virtual void WriteStreamedRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges> const& value);
  virtual void EndStreamedRecordWithChangesImpl() = 0;
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
  void ReadInt8ToInt(int8_t& value);

  // Ordinal 1.
  void ReadInt8ToLong(int8_t& value);

  // Ordinal 2.
  void ReadInt8ToUint(int8_t& value);

  // Ordinal 3.
  void ReadInt8ToUlong(int8_t& value);

  // Ordinal 4.
  void ReadInt8ToFloat(int8_t& value);

  // Ordinal 5.
  void ReadInt8ToDouble(int8_t& value);

  // Ordinal 6.
  void ReadIntToUint(int32_t& value);

  // Ordinal 7.
  void ReadIntToLong(int32_t& value);

  // Ordinal 8.
  void ReadIntToFloat(int32_t& value);

  // Ordinal 9.
  void ReadIntToDouble(int32_t& value);

  // Ordinal 10.
  void ReadUintToUlong(uint32_t& value);

  // Ordinal 11.
  void ReadUintToFloat(uint32_t& value);

  // Ordinal 12.
  void ReadUintToDouble(uint32_t& value);

  // Ordinal 13.
  void ReadFloatToDouble(float& value);

  // Ordinal 14.
  void ReadIntToString(int32_t& value);

  // Ordinal 15.
  void ReadUintToString(uint32_t& value);

  // Ordinal 16.
  void ReadLongToString(int64_t& value);

  // Ordinal 17.
  void ReadUlongToString(uint64_t& value);

  // Ordinal 18.
  void ReadFloatToString(float& value);

  // Ordinal 19.
  void ReadDoubleToString(double& value);

  // Ordinal 20.
  void ReadIntToOptional(int32_t& value);

  // Ordinal 21.
  void ReadFloatToOptional(float& value);

  // Ordinal 22.
  void ReadStringToOptional(std::string& value);

  // Ordinal 23.
  void ReadIntToUnion(int32_t& value);

  // Ordinal 24.
  void ReadFloatToUnion(float& value);

  // Ordinal 25.
  void ReadStringToUnion(std::string& value);

  // Ordinal 26.
  void ReadOptionalIntToFloat(std::optional<int32_t>& value);

  // Ordinal 27.
  void ReadOptionalFloatToString(std::optional<float>& value);

  // Ordinal 28.
  void ReadAliasedLongToString(evo_test::AliasedLongToString& value);

  // Ordinal 29.
  void ReadStringToAliasedString(std::string& value);

  // Ordinal 30.
  void ReadStringToAliasedInt(std::string& value);

  // Ordinal 31.
  void ReadOptionalIntToUnion(std::optional<int32_t>& value);

  // Ordinal 32.
  void ReadOptionalRecordToUnion(std::optional<evo_test::RecordWithChanges>& value);

  // Ordinal 33.
  void ReadRecordWithChanges(evo_test::RecordWithChanges& value);

  // Ordinal 34.
  void ReadAliasedRecordWithChanges(evo_test::AliasedRecordWithChanges& value);

  // Ordinal 35.
  void ReadRecordToRenamedRecord(evo_test::RenamedRecord& value);

  // Ordinal 36.
  void ReadRecordToAliasedRecord(evo_test::RecordWithChanges& value);

  // Ordinal 37.
  void ReadRecordToAliasedAlias(evo_test::RecordWithChanges& value);

  // Ordinal 38.
  [[nodiscard]] bool ReadStreamOfAliasTypeChange(evo_test::StreamItem& value);

  // Ordinal 38.
  [[nodiscard]] bool ReadStreamOfAliasTypeChange(std::vector<evo_test::StreamItem>& values);

  // Ordinal 39.
  // Comprehensive NamedType changes
  void ReadRlink(evo_test::RLink& value);

  // Ordinal 40.
  void ReadRlinkRX(evo_test::RLink& value);

  // Ordinal 41.
  void ReadRlinkRY(evo_test::RLink& value);

  // Ordinal 42.
  void ReadRlinkRZ(evo_test::RLink& value);

  // Ordinal 43.
  void ReadRaRLink(evo_test::RA& value);

  // Ordinal 44.
  void ReadRaRX(evo_test::RA& value);

  // Ordinal 45.
  void ReadRaRY(evo_test::RA& value);

  // Ordinal 46.
  void ReadRaRZ(evo_test::RA& value);

  // Ordinal 47.
  void ReadRbRLink(evo_test::RB& value);

  // Ordinal 48.
  void ReadRbRX(evo_test::RB& value);

  // Ordinal 49.
  void ReadRbRY(evo_test::RB& value);

  // Ordinal 50.
  void ReadRbRZ(evo_test::RB& value);

  // Ordinal 51.
  void ReadRcRLink(evo_test::RC& value);

  // Ordinal 52.
  void ReadRcRX(evo_test::RC& value);

  // Ordinal 53.
  void ReadRcRY(evo_test::RC& value);

  // Ordinal 54.
  void ReadRcRZ(evo_test::RC& value);

  // Ordinal 55.
  void ReadRlinkRNew(evo_test::RLink& value);

  // Ordinal 56.
  void ReadRaRNew(evo_test::RA& value);

  // Ordinal 57.
  void ReadRbRNew(evo_test::RB& value);

  // Ordinal 58.
  void ReadRcRNew(evo_test::RC& value);

  // Ordinal 59.
  void ReadRlinkRUnion(evo_test::RLink& value);

  // Ordinal 60.
  void ReadRaRUnion(evo_test::RA& value);

  // Ordinal 61.
  void ReadRbRUnion(evo_test::RB& value);

  // Ordinal 62.
  void ReadRcRUnion(evo_test::RC& value);

  // Ordinal 63.
  void ReadOptionalRecordWithChanges(std::optional<evo_test::RecordWithChanges>& value);

  // Ordinal 64.
  void ReadAliasedOptionalRecordWithChanges(std::optional<evo_test::AliasedRecordWithChanges>& value);

  // Ordinal 65.
  void ReadUnionRecordWithChanges(std::variant<evo_test::RecordWithChanges, int32_t>& value);

  // Ordinal 66.
  void ReadUnionWithSameTypeset(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value);

  // Ordinal 67.
  void ReadUnionWithTypesAdded(std::variant<evo_test::RecordWithChanges, float>& value);

  // Ordinal 68.
  void ReadUnionWithTypesRemoved(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value);

  // Ordinal 69.
  void ReadRecordToOptional(evo_test::RecordWithChanges& value);

  // Ordinal 70.
  void ReadRecordToAliasedOptional(evo_test::RecordWithChanges& value);

  // Ordinal 71.
  void ReadRecordToUnion(evo_test::RecordWithChanges& value);

  // Ordinal 72.
  void ReadRecordToAliasedUnion(evo_test::RecordWithChanges& value);

  // Ordinal 73.
  void ReadUnionToAliasedUnion(std::variant<evo_test::RecordWithChanges, int32_t>& value);

  // Ordinal 74.
  void ReadUnionToAliasedUnionWithChanges(std::variant<evo_test::RecordWithChanges, int32_t>& value);

  // Ordinal 75.
  void ReadOptionalToAliasedOptional(std::optional<evo_test::RecordWithChanges>& value);

  // Ordinal 76.
  void ReadOptionalToAliasedOptionalWithChanges(std::optional<int32_t>& value);

  // Ordinal 77.
  void ReadGenericRecord(evo_test::GenericRecord<int32_t, std::string>& value);

  // Ordinal 78.
  void ReadGenericRecordToOpenAlias(evo_test::GenericRecord<int32_t, std::string>& value);

  // Ordinal 79.
  void ReadGenericRecordToClosedAlias(evo_test::GenericRecord<int32_t, std::string>& value);

  // Ordinal 80.
  void ReadGenericRecordToHalfClosedAlias(evo_test::GenericRecord<int32_t, std::string>& value);

  // Ordinal 81.
  void ReadAliasedGenericRecordToAlias(evo_test::AliasedHalfClosedGenericRecord<int32_t>& value);

  // Ordinal 82.
  void ReadClosedGenericRecordToUnion(evo_test::AliasedClosedGenericRecord& value);

  // Ordinal 83.
  void ReadGenericRecordToAliasedUnion(evo_test::GenericRecord<int32_t, std::string>& value);

  // Ordinal 84.
  void ReadGenericUnionOfChangedRecord(evo_test::AliasedClosedGenericUnion& value);

  // Ordinal 85.
  void ReadGenericParentRecord(evo_test::GenericParentRecord<int32_t>& value);

  // Ordinal 86.
  void ReadGenericNestedRecords(evo_test::GenericRecord<evo_test::UnchangedGeneric<int32_t>, evo_test::ChangedGeneric<std::string, int32_t>>& value);

  // Ordinal 87.
  [[nodiscard]] bool ReadGenericRecordStream(evo_test::GenericRecord<int32_t, std::string>& value);

  // Ordinal 87.
  [[nodiscard]] bool ReadGenericRecordStream(std::vector<evo_test::GenericRecord<int32_t, std::string>>& values);

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

  // Optionaly close this writer before destructing. Validates that all steps were completely read.
  void Close();

  void CopyTo(ProtocolWithChangesWriterBase& writer, size_t stream_of_alias_type_change_buffer_size = 1, size_t generic_record_stream_buffer_size = 1, size_t generic_parent_record_stream_buffer_size = 1, size_t streamed_record_with_changes_buffer_size = 1);

  virtual ~ProtocolWithChangesReaderBase() = default;

  protected:
  virtual void ReadInt8ToIntImpl(int8_t& value) = 0;
  virtual void ReadInt8ToLongImpl(int8_t& value) = 0;
  virtual void ReadInt8ToUintImpl(int8_t& value) = 0;
  virtual void ReadInt8ToUlongImpl(int8_t& value) = 0;
  virtual void ReadInt8ToFloatImpl(int8_t& value) = 0;
  virtual void ReadInt8ToDoubleImpl(int8_t& value) = 0;
  virtual void ReadIntToUintImpl(int32_t& value) = 0;
  virtual void ReadIntToLongImpl(int32_t& value) = 0;
  virtual void ReadIntToFloatImpl(int32_t& value) = 0;
  virtual void ReadIntToDoubleImpl(int32_t& value) = 0;
  virtual void ReadUintToUlongImpl(uint32_t& value) = 0;
  virtual void ReadUintToFloatImpl(uint32_t& value) = 0;
  virtual void ReadUintToDoubleImpl(uint32_t& value) = 0;
  virtual void ReadFloatToDoubleImpl(float& value) = 0;
  virtual void ReadIntToStringImpl(int32_t& value) = 0;
  virtual void ReadUintToStringImpl(uint32_t& value) = 0;
  virtual void ReadLongToStringImpl(int64_t& value) = 0;
  virtual void ReadUlongToStringImpl(uint64_t& value) = 0;
  virtual void ReadFloatToStringImpl(float& value) = 0;
  virtual void ReadDoubleToStringImpl(double& value) = 0;
  virtual void ReadIntToOptionalImpl(int32_t& value) = 0;
  virtual void ReadFloatToOptionalImpl(float& value) = 0;
  virtual void ReadStringToOptionalImpl(std::string& value) = 0;
  virtual void ReadIntToUnionImpl(int32_t& value) = 0;
  virtual void ReadFloatToUnionImpl(float& value) = 0;
  virtual void ReadStringToUnionImpl(std::string& value) = 0;
  virtual void ReadOptionalIntToFloatImpl(std::optional<int32_t>& value) = 0;
  virtual void ReadOptionalFloatToStringImpl(std::optional<float>& value) = 0;
  virtual void ReadAliasedLongToStringImpl(evo_test::AliasedLongToString& value) = 0;
  virtual void ReadStringToAliasedStringImpl(std::string& value) = 0;
  virtual void ReadStringToAliasedIntImpl(std::string& value) = 0;
  virtual void ReadOptionalIntToUnionImpl(std::optional<int32_t>& value) = 0;
  virtual void ReadOptionalRecordToUnionImpl(std::optional<evo_test::RecordWithChanges>& value) = 0;
  virtual void ReadRecordWithChangesImpl(evo_test::RecordWithChanges& value) = 0;
  virtual void ReadAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges& value) = 0;
  virtual void ReadRecordToRenamedRecordImpl(evo_test::RenamedRecord& value) = 0;
  virtual void ReadRecordToAliasedRecordImpl(evo_test::RecordWithChanges& value) = 0;
  virtual void ReadRecordToAliasedAliasImpl(evo_test::RecordWithChanges& value) = 0;
  virtual bool ReadStreamOfAliasTypeChangeImpl(evo_test::StreamItem& value) = 0;
  virtual bool ReadStreamOfAliasTypeChangeImpl(std::vector<evo_test::StreamItem>& values);
  virtual void ReadRlinkImpl(evo_test::RLink& value) = 0;
  virtual void ReadRlinkRXImpl(evo_test::RLink& value) = 0;
  virtual void ReadRlinkRYImpl(evo_test::RLink& value) = 0;
  virtual void ReadRlinkRZImpl(evo_test::RLink& value) = 0;
  virtual void ReadRaRLinkImpl(evo_test::RA& value) = 0;
  virtual void ReadRaRXImpl(evo_test::RA& value) = 0;
  virtual void ReadRaRYImpl(evo_test::RA& value) = 0;
  virtual void ReadRaRZImpl(evo_test::RA& value) = 0;
  virtual void ReadRbRLinkImpl(evo_test::RB& value) = 0;
  virtual void ReadRbRXImpl(evo_test::RB& value) = 0;
  virtual void ReadRbRYImpl(evo_test::RB& value) = 0;
  virtual void ReadRbRZImpl(evo_test::RB& value) = 0;
  virtual void ReadRcRLinkImpl(evo_test::RC& value) = 0;
  virtual void ReadRcRXImpl(evo_test::RC& value) = 0;
  virtual void ReadRcRYImpl(evo_test::RC& value) = 0;
  virtual void ReadRcRZImpl(evo_test::RC& value) = 0;
  virtual void ReadRlinkRNewImpl(evo_test::RLink& value) = 0;
  virtual void ReadRaRNewImpl(evo_test::RA& value) = 0;
  virtual void ReadRbRNewImpl(evo_test::RB& value) = 0;
  virtual void ReadRcRNewImpl(evo_test::RC& value) = 0;
  virtual void ReadRlinkRUnionImpl(evo_test::RLink& value) = 0;
  virtual void ReadRaRUnionImpl(evo_test::RA& value) = 0;
  virtual void ReadRbRUnionImpl(evo_test::RB& value) = 0;
  virtual void ReadRcRUnionImpl(evo_test::RC& value) = 0;
  virtual void ReadOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges>& value) = 0;
  virtual void ReadAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges>& value) = 0;
  virtual void ReadUnionRecordWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t>& value) = 0;
  virtual void ReadUnionWithSameTypesetImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value) = 0;
  virtual void ReadUnionWithTypesAddedImpl(std::variant<evo_test::RecordWithChanges, float>& value) = 0;
  virtual void ReadUnionWithTypesRemovedImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value) = 0;
  virtual void ReadRecordToOptionalImpl(evo_test::RecordWithChanges& value) = 0;
  virtual void ReadRecordToAliasedOptionalImpl(evo_test::RecordWithChanges& value) = 0;
  virtual void ReadRecordToUnionImpl(evo_test::RecordWithChanges& value) = 0;
  virtual void ReadRecordToAliasedUnionImpl(evo_test::RecordWithChanges& value) = 0;
  virtual void ReadUnionToAliasedUnionImpl(std::variant<evo_test::RecordWithChanges, int32_t>& value) = 0;
  virtual void ReadUnionToAliasedUnionWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t>& value) = 0;
  virtual void ReadOptionalToAliasedOptionalImpl(std::optional<evo_test::RecordWithChanges>& value) = 0;
  virtual void ReadOptionalToAliasedOptionalWithChangesImpl(std::optional<int32_t>& value) = 0;
  virtual void ReadGenericRecordImpl(evo_test::GenericRecord<int32_t, std::string>& value) = 0;
  virtual void ReadGenericRecordToOpenAliasImpl(evo_test::GenericRecord<int32_t, std::string>& value) = 0;
  virtual void ReadGenericRecordToClosedAliasImpl(evo_test::GenericRecord<int32_t, std::string>& value) = 0;
  virtual void ReadGenericRecordToHalfClosedAliasImpl(evo_test::GenericRecord<int32_t, std::string>& value) = 0;
  virtual void ReadAliasedGenericRecordToAliasImpl(evo_test::AliasedHalfClosedGenericRecord<int32_t>& value) = 0;
  virtual void ReadClosedGenericRecordToUnionImpl(evo_test::AliasedClosedGenericRecord& value) = 0;
  virtual void ReadGenericRecordToAliasedUnionImpl(evo_test::GenericRecord<int32_t, std::string>& value) = 0;
  virtual void ReadGenericUnionOfChangedRecordImpl(evo_test::AliasedClosedGenericUnion& value) = 0;
  virtual void ReadGenericParentRecordImpl(evo_test::GenericParentRecord<int32_t>& value) = 0;
  virtual void ReadGenericNestedRecordsImpl(evo_test::GenericRecord<evo_test::UnchangedGeneric<int32_t>, evo_test::ChangedGeneric<std::string, int32_t>>& value) = 0;
  virtual bool ReadGenericRecordStreamImpl(evo_test::GenericRecord<int32_t, std::string>& value) = 0;
  virtual bool ReadGenericRecordStreamImpl(std::vector<evo_test::GenericRecord<int32_t, std::string>>& values);
  virtual bool ReadGenericParentRecordStreamImpl(evo_test::GenericParentRecord<int32_t>& value) = 0;
  virtual bool ReadGenericParentRecordStreamImpl(std::vector<evo_test::GenericParentRecord<int32_t>>& values);
  virtual void ReadVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges>& value) = 0;
  virtual bool ReadStreamedRecordWithChangesImpl(evo_test::RecordWithChanges& value) = 0;
  virtual bool ReadStreamedRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges>& values);
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
