#include "../evolution_testing.h"
#include "generated/binary/protocols.h"

using namespace evo_test;

int main(void) {
  ::binary::ProtocolWithChangesWriter w(std::cout);

  w.WriteInt8ToInt(static_cast<int>(INT8_MIN));
  w.WriteInt8ToLong(static_cast<long>(INT8_MIN));
  w.WriteInt8ToUint(static_cast<uint32_t>(INT8_MIN));
  w.WriteInt8ToUlong(static_cast<uint64_t>(INT8_MIN));
  w.WriteInt8ToFloat(static_cast<float>(INT8_MIN));
  w.WriteInt8ToDouble(static_cast<double>(INT8_MIN));

  w.WriteIntToUint(static_cast<uint32_t>(INT_MIN));
  w.WriteIntToLong(static_cast<long>(INT_MIN));
  w.WriteIntToFloat(static_cast<float>(INT8_MIN));
  w.WriteIntToDouble(static_cast<double>(INT8_MIN));
  w.WriteUintToUlong(static_cast<uint64_t>(UINT_MAX));
  w.WriteUintToFloat(static_cast<float>(UINT8_MAX));
  w.WriteUintToDouble(static_cast<double>(UINT8_MAX));

  w.WriteFloatToDouble(M_PI);

  w.WriteIntToString(std::to_string(INT_MIN));
  w.WriteUintToString(std::to_string(UINT_MAX));
  w.WriteLongToString(std::to_string(LONG_MIN));
  w.WriteUlongToString(std::to_string(ULONG_MAX));
  w.WriteFloatToString(std::to_string(M_PI));
  w.WriteDoubleToString(std::to_string(M_PI));

  w.WriteIntToOptional(INT_MIN);
  w.WriteFloatToOptional(M_PI);
  w.WriteStringToOptional(HelloWorld);

  w.WriteIntToUnion(std::variant<int, bool>(INT_MIN));
  w.WriteFloatToUnion(std::variant<float, bool>(static_cast<float>(M_PI)));
  w.WriteStringToUnion(HelloWorld);

  w.WriteOptionalIntToFloat(static_cast<float>(INT8_MIN));
  w.WriteOptionalFloatToString(std::to_string(M_PI));

  w.WriteAliasedLongToString(std::to_string(LONG_MIN));
  w.WriteStringToAliasedString(HelloWorld);
  w.WriteStringToAliasedInt(INT_MIN);

  w.WriteEnumToAliasedEnum(GrowingEnum::kC);

  UnchangedRecord unchanged;
  unchanged.name = "Jane Doe";
  unchanged.age = 42;
  unchanged.meta = {{"height", 161.3}, {"weight", 75.0f}};

  RenamedRecord renamed;
  renamed.s = HelloWorld;
  renamed.i = INT_MIN;

  RecordWithChanges rec;
  rec.int_to_long = static_cast<long>(INT_MIN);
  rec.float_to_double = M_PI;
  rec.optional_long_to_string = std::to_string(LONG_MIN);
  rec.unchanged_record = unchanged;

  w.WriteOptionalIntToUnion(std::variant<std::monostate, int, std::string>(INT_MIN));
  w.WriteOptionalRecordToUnion(std::variant<std::monostate, RecordWithChanges, std::string>(rec));

  w.WriteRecordWithChanges(rec);
  w.WriteAliasedRecordWithChanges(rec);
  w.WriteRecordToRenamedRecord(renamed);
  w.WriteRecordToAliasedRecord(rec);
  w.WriteRecordToAliasedAlias(rec);

  for (int i = 0; i < 42; i++) {
    w.WriteStreamIntToStringToFloat(std::to_string(42));
  }
  w.EndStreamIntToStringToFloat();

  std::vector<std::string> str_vec(42, std::to_string(42));
  w.WriteVectorIntToStringToFloat(str_vec);

  w.WriteIntFloatUnionReordered(static_cast<float>(M_PI));

  std::vector<std::variant<float, int>> vec(42, static_cast<float>(M_PI));
  w.WriteVectorUnionReordered(vec);

  for (int i = 0; i < 42; i++) {
    w.WriteStreamUnionReordered(HelloWorld);
  }
  w.EndStreamUnionReordered();

  std::vector<int> int_vec(42, 42);
  w.WriteIntToUnionStream(int_vec);
  w.EndIntToUnionStream();

  std::vector<std::variant<int, bool>> int_bool_vec(42, 42);
  w.WriteUnionStreamTypeChange(int_bool_vec);
  w.EndUnionStreamTypeChange();

  w.WriteStreamOfAliasTypeChange(std::vector<StreamItem>(7, rec));
  w.EndStreamOfAliasTypeChange();

  RZ record;
  record.subject = 42;

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

  RNew rnew = record;
  w.WriteRlinkRNew(rnew);
  w.WriteRaRNew(rnew);
  w.WriteRbRNew(rnew);
  w.WriteRcRNew(rnew);

  RUnion runion = record;
  w.WriteRlinkRUnion(runion);
  w.WriteRaRUnion(runion);
  w.WriteRbRUnion(runion);
  w.WriteRcRUnion(runion);

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
  w.WriteOptionalToAliasedOptionalWithChanges(std::to_string(INT_MIN));

  GenericRecord<int, std::string> generic_record;
  generic_record.field_2 = HelloWorld;
  generic_record.field_1 = 42;
  generic_record.added = true;

  w.WriteGenericRecord(generic_record);
  w.WriteGenericRecordToOpenAlias(generic_record);
  w.WriteGenericRecordToClosedAlias(generic_record);
  w.WriteGenericRecordToHalfClosedAlias(generic_record);
  w.WriteAliasedGenericRecordToAlias(generic_record);
  w.WriteGenericRecordToReversed(generic_record);

  w.WriteClosedGenericRecordToUnion(generic_record);
  w.WriteGenericRecordToAliasedUnion(generic_record);

  w.WriteGenericUnionToReversed(generic_record);
  w.WriteGenericUnionOfChangedRecord(generic_record);

  GenericParentRecord<int> generic_parent;
  generic_parent.record = generic_record;

  GenericRecord<GenericUnion<int, float>, std::string> generic_record_of_union;
  generic_record_of_union.field_1 = 42;
  generic_record_of_union.field_2 = HelloWorld;
  generic_parent.record_of_union = generic_record_of_union;

  generic_parent.union_of_record = generic_record;

  w.WriteGenericParentRecord(generic_parent);

  GenericRecord<Unchanged, Changed> generic_nested;
  generic_nested.field_1.field = 42;
  generic_nested.field_2.y = "42";
  generic_nested.field_2.z = Unchanged{42};
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

  w.WriteAddedOptional(rec);
  w.WriteAddedMap({{"hello", "world"}});
  w.WriteAddedRecordStream(recs);
  w.EndAddedRecordStream();

  w.Close();

  return 0;
}
