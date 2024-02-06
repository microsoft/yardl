#include "lib/binary/protocols.h"

using namespace evo_test;

static std::string HelloWorld = "Hello, World!";

int main(void) {
  ::binary::ProtocolWithChangesWriter w(std::cout);

  w.WriteInt8ToInt(INT8_MIN);
  w.WriteInt8ToLong(INT8_MIN);
  w.WriteInt8ToUint(INT8_MIN);
  w.WriteInt8ToUlong(INT8_MIN);
  w.WriteInt8ToFloat(INT8_MIN);
  w.WriteInt8ToDouble(INT8_MIN);

  w.WriteIntToUint(INT_MIN);
  w.WriteIntToLong(INT_MIN);
  w.WriteIntToFloat(INT8_MIN);
  w.WriteIntToDouble(INT8_MIN);
  w.WriteUintToUlong(UINT_MAX);
  w.WriteUintToFloat(UINT8_MAX);
  w.WriteUintToDouble(UINT8_MAX);

  w.WriteFloatToDouble(static_cast<float>(M_PI));

  w.WriteIntToString(INT_MIN);
  w.WriteUintToString(UINT_MAX);
  w.WriteLongToString(LONG_MIN);
  w.WriteUlongToString(ULONG_MAX);
  w.WriteFloatToString(M_PI);
  w.WriteDoubleToString(M_PI);

  w.WriteIntToOptional(INT_MIN);
  w.WriteFloatToOptional(M_PI);
  w.WriteStringToOptional(HelloWorld);

  w.WriteIntToUnion(INT_MIN);
  w.WriteFloatToUnion(M_PI);
  w.WriteStringToUnion(HelloWorld);

  w.WriteOptionalIntToFloat(INT8_MIN);
  w.WriteOptionalFloatToString(M_PI);

  w.WriteAliasedLongToString(LONG_MIN);
  w.WriteStringToAliasedString(HelloWorld);
  w.WriteStringToAliasedInt(std::to_string(INT_MIN));

  UnchangedRecord unchanged;
  unchanged.name = "Jane Doe";
  unchanged.age = 42;
  unchanged.meta = {{"height", 161.3}, {"weight", 75.0f}};

  RenamedRecord renamed;
  renamed.i = INT_MIN;
  renamed.s = HelloWorld;

  RecordWithChanges rec;
  rec.int_to_long = INT_MIN;
  rec.deprecated_vector = {1, 2, 3};
  rec.float_to_double = M_PI;
  rec.deprecated_array = {4, 5, 6};
  rec.optional_long_to_string = LONG_MIN;
  rec.deprecated_map = {{"a", {1, 4, 7}}, {"b", {2, 5, 8}}, {"c", {3, 6, 9}}};
  rec.unchanged_record = unchanged;

  w.WriteOptionalIntToUnion(INT_MIN);
  w.WriteOptionalRecordToUnion(rec);

  w.WriteRecordWithChanges(rec);
  w.WriteAliasedRecordWithChanges(rec);
  w.WriteRecordToRenamedRecord(renamed);
  w.WriteRecordToAliasedRecord(rec);
  w.WriteRecordToAliasedAlias(rec);

  w.WriteStreamOfAliasTypeChange(std::vector<StreamItem>(7, rec));
  w.EndStreamOfAliasTypeChange();

  RC record;
  record.subject = "42";

  w.WriteRlink(record);
  w.WriteRlinkRX(record);
  w.WriteRlinkRY(record);
  w.WriteRlinkRZ(record);

  w.WriteRaRLink(record);
  w.WriteRaRX(record);
  w.WriteRaRY(record);
  w.WriteRaRZ(record);

  w.WriteRbRLink(record);
  w.WriteRbRX(record);
  w.WriteRbRY(record);
  w.WriteRbRZ(record);

  w.WriteRcRLink(record);
  w.WriteRcRX(record);
  w.WriteRcRY(record);
  w.WriteRcRZ(record);

  w.WriteRlinkRNew(record);
  w.WriteRaRNew(record);
  w.WriteRbRNew(record);
  w.WriteRcRNew(record);

  w.WriteRlinkRUnion(record);
  w.WriteRaRUnion(record);
  w.WriteRbRUnion(record);
  w.WriteRcRUnion(record);

  w.WriteOptionalRecordWithChanges(rec);
  w.WriteAliasedOptionalRecordWithChanges(rec);

  w.WriteUnionRecordWithChanges(rec);
  // w.WriteAliasedUnionRecordWithChanges(rec);

  w.WriteUnionWithSameTypeset(rec);
  w.WriteUnionWithTypesAdded(rec);
  w.WriteUnionWithTypesRemoved(rec);

  w.WriteRecordToOptional(rec);
  w.WriteRecordToAliasedOptional(rec);
  w.WriteRecordToUnion(rec);
  w.WriteRecordToAliasedUnion(rec);

  w.WriteUnionToAliasedUnion(rec);
  w.WriteUnionToAliasedUnionWithChanges(rec);
  w.WriteOptionalToAliasedOptional(rec);
  w.WriteOptionalToAliasedOptionalWithChanges(INT_MIN);

  GenericRecord<int, std::string> generic_record;
  generic_record.removed = true;
  generic_record.field_1 = 42;
  generic_record.field_2 = "42";

  w.WriteGenericRecord(generic_record);
  w.WriteGenericRecordToOpenAlias(generic_record);
  w.WriteGenericRecordToClosedAlias(generic_record);
  w.WriteGenericRecordToHalfClosedAlias(generic_record);
  w.WriteAliasedGenericRecordToAlias(generic_record);

  w.WriteClosedGenericRecordToUnion(generic_record);
  w.WriteGenericRecordToAliasedUnion(generic_record);

  w.WriteGenericUnionOfChangedRecord(generic_record);

  GenericParentRecord<int> generic_parent;
  generic_parent.record = generic_record;

  GenericRecord<GenericUnion<int, float>, std::string> generic_record_of_union;
  generic_record_of_union.field_1 = 42;
  generic_record_of_union.field_2 = "Hello, World";
  generic_parent.record_of_union = generic_record_of_union;

  generic_parent.union_of_record = generic_record;

  w.WriteGenericParentRecord(generic_parent);

  GenericRecord<UnchangedGeneric<int>, ChangedGeneric<std::string, int>> generic_nested;
  generic_nested.field_1.field = 42;
  generic_nested.field_2.y = "42";
  generic_nested.field_2.z.field = 42;
  w.WriteGenericNestedRecords(generic_nested);

  w.WriteGenericRecordStream(std::vector<AliasedClosedGenericRecord>(7, generic_record));
  w.EndGenericRecordStream();
  w.WriteGenericParentRecordStream(std::vector<GenericParentRecord<int>>(7, generic_parent));
  w.EndGenericParentRecordStream();

  // Write a vector of size 7 records
  std::vector<RecordWithChanges> recs(7, rec);
  w.WriteVectorRecordWithChanges(recs);

  // Stream a total of 7 records
  w.WriteStreamedRecordWithChanges(rec);
  w.WriteStreamedRecordWithChanges(rec);
  std::vector<RecordWithChanges> more_recs(4, rec);
  w.WriteStreamedRecordWithChanges(more_recs);
  w.WriteStreamedRecordWithChanges(rec);
  w.EndStreamedRecordWithChanges();

  w.WriteAddedStringVector(std::vector<std::string>(7, HelloWorld));
  w.WriteAddedOptional(rec);
  w.WriteAddedMap({{"hello", "world"}});
  w.WriteAddedUnion(rec);
  w.WriteAddedRecordStream(recs);
  w.EndAddedRecordStream();

  for (int i = 0; i < 7; ++i) {
    w.WriteAddedUnionStream(rec);
  }
  w.EndAddedUnionStream();

  w.Close();

  return 0;
}
