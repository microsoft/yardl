#include "../evolution_testing.h"
#include "generated/binary/protocols.h"

using namespace evo_test;

void validateRecordWithChanges(RecordWithChanges const& rec) {
  (void)(rec);

  EVO_ASSERT(rec.int_to_long == INT_MIN);
  EVO_ASSERT_EQUALISH(rec.float_to_double, M_PI);
  EVO_ASSERT(rec.optional_long_to_string.has_value());
  EVO_ASSERT(rec.optional_long_to_string.value() == LONG_MIN);

  EVO_ASSERT(rec.unchanged_record.name == "Jane Doe");
  EVO_ASSERT(rec.unchanged_record.age == 42);
  EVO_ASSERT(rec.unchanged_record.meta.at("height") == 161.3);
  EVO_ASSERT(rec.unchanged_record.meta.at("weight") == 75.0f);
}

void validateGenericRecord(GenericRecord<int, std::string> record) {
  if (record.removed.has_value()) {
    EVO_ASSERT(record.removed.value() == true);
  }
  EVO_ASSERT(record.field_1 == 42);
  EVO_ASSERT(record.field_2 == "42");
}

void validateGenericParentRecord(GenericParentRecord<int> parent) {
  validateGenericRecord(parent.record);
  EVO_ASSERT(parent.record_of_union.field_1.index() == 0);
  EVO_ASSERT(std::get<0>(parent.record_of_union.field_1) == 42);
  EVO_ASSERT(parent.record_of_union.field_2 == "Hello, World");
  EVO_ASSERT(parent.union_of_record.index() == 0);
  validateGenericRecord(std::get<0>(parent.union_of_record));
}

int main(void) {
  ::binary::ProtocolWithChangesReader r(std::cin);

  int8_t int8 = 0;
  int32_t int32 = 0;
  int64_t int64 = 0;
  uint32_t uint32;
  uint64_t uint64 = 0;
  float flt = 0.f;
  double dbl = 0.f;
  std::string str;

  r.ReadInt8ToInt(int8);
  EVO_ASSERT(int8 == INT8_MIN);
  r.ReadInt8ToLong(int8);
  EVO_ASSERT(int8 == INT8_MIN);
  r.ReadInt8ToUint(int8);
  EVO_ASSERT(int8 == INT8_MIN);
  r.ReadInt8ToUlong(int8);
  EVO_ASSERT(int8 == INT8_MIN);
  r.ReadInt8ToFloat(int8);
  EVO_ASSERT(int8 == INT8_MIN);
  r.ReadInt8ToDouble(int8);
  EVO_ASSERT(int8 == INT8_MIN);

  r.ReadIntToUint(int32);
  EVO_ASSERT(int32 == INT_MIN);
  r.ReadIntToLong(int32);
  EVO_ASSERT(int32 == INT_MIN);
  r.ReadIntToFloat(int32);
  EVO_ASSERT(int32 == INT8_MIN);
  r.ReadIntToDouble(int32);
  EVO_ASSERT(int32 == INT8_MIN);

  r.ReadUintToUlong(uint32);
  EVO_ASSERT(uint32 == UINT_MAX);
  r.ReadUintToFloat(uint32);
  EVO_ASSERT(uint32 == UINT8_MAX);
  r.ReadUintToDouble(uint32);
  EVO_ASSERT(uint32 == UINT8_MAX);

  r.ReadFloatToDouble(flt);
  EVO_ASSERT_EQUALISH(flt, M_PI);

  r.ReadIntToString(int32);
  EVO_ASSERT(int32 == INT_MIN);
  r.ReadUintToString(uint32);
  EVO_ASSERT(uint32 == UINT_MAX);
  r.ReadLongToString(int64);
  EVO_ASSERT(int64 == LONG_MIN);
  r.ReadUlongToString(uint64);
  EVO_ASSERT(uint64 == ULONG_MAX);
  r.ReadFloatToString(flt);
  EVO_ASSERT_EQUALISH(flt, M_PI);
  r.ReadDoubleToString(dbl);
  EVO_ASSERT_EQUALISH(dbl, M_PI);

  r.ReadIntToOptional(int32);
  EVO_ASSERT(int32 == INT_MIN);
  r.ReadFloatToOptional(flt);
  EVO_ASSERT_EQUALISH(flt, M_PI);

  r.ReadStringToOptional(str);
  EVO_ASSERT(str == HelloWorld);

  r.ReadIntToUnion(int32);
  EVO_ASSERT(int32 == INT_MIN);
  r.ReadFloatToUnion(flt);
  EVO_ASSERT_EQUALISH(flt, M_PI);
  r.ReadStringToUnion(str);
  EVO_ASSERT(str == HelloWorld);

  std::optional<int32_t> maybe_int32;
  r.ReadOptionalIntToFloat(maybe_int32);
  EVO_ASSERT(maybe_int32.has_value());
  EVO_ASSERT(maybe_int32 == INT8_MIN);

  std::optional<float> maybe_flt;
  r.ReadOptionalFloatToString(maybe_flt);
  EVO_ASSERT(maybe_flt.has_value());
  EVO_ASSERT_EQUALISH(maybe_flt.value(), M_PI);

  r.ReadAliasedLongToString(int64);
  EVO_ASSERT(int64 == LONG_MIN);
  r.ReadStringToAliasedString(str);
  EVO_ASSERT(str == HelloWorld);
  r.ReadStringToAliasedInt(str);
  EVO_ASSERT(str == std::to_string(INT_MIN));

  RecordWithChanges rec;
  std::optional<RecordWithChanges> maybe_rec;

  r.ReadOptionalIntToUnion(maybe_int32);
  EVO_ASSERT(maybe_int32.has_value());
  EVO_ASSERT(maybe_int32 == INT_MIN);

  r.ReadOptionalRecordToUnion(maybe_rec);
  EVO_ASSERT(maybe_rec.has_value());
  validateRecordWithChanges(maybe_rec.value());

  r.ReadRecordWithChanges(rec);
  validateRecordWithChanges(rec);

  r.ReadAliasedRecordWithChanges(rec);
  validateRecordWithChanges(rec);

  RenamedRecord renamed;
  r.ReadRecordToRenamedRecord(renamed);
  EVO_ASSERT(renamed.i == INT_MIN);
  EVO_ASSERT(renamed.s == HelloWorld);

  r.ReadRecordToAliasedRecord(rec);
  validateRecordWithChanges(rec);

  r.ReadRecordToAliasedAlias(rec);
  validateRecordWithChanges(rec);

  std::vector<StreamItem> stream_items(10);
  while (r.ReadStreamOfAliasTypeChange(stream_items)) {
    EVO_ASSERT(stream_items.size() == 7);
    for (auto const& item : stream_items) {
      validateRecordWithChanges(item);
    }
  }

  RC record;
  r.ReadRlink(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRlinkRX(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRlinkRY(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRlinkRZ(record);
  EVO_ASSERT(record.subject == "42");

  r.ReadRaRLink(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRaRX(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRaRY(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRaRZ(record);
  EVO_ASSERT(record.subject == "42");

  r.ReadRbRLink(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRbRX(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRbRY(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRbRZ(record);
  EVO_ASSERT(record.subject == "42");

  r.ReadRcRLink(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRcRX(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRcRY(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRcRZ(record);
  EVO_ASSERT(record.subject == "42");

  r.ReadRlinkRNew(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRaRNew(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRbRNew(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRcRNew(record);
  EVO_ASSERT(record.subject == "42");

  r.ReadRlinkRUnion(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRaRUnion(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRbRUnion(record);
  EVO_ASSERT(record.subject == "42");
  r.ReadRcRUnion(record);
  EVO_ASSERT(record.subject == "42");

  r.ReadOptionalRecordWithChanges(maybe_rec);
  EVO_ASSERT(maybe_rec.has_value());
  validateRecordWithChanges(maybe_rec.value());

  r.ReadAliasedOptionalRecordWithChanges(maybe_rec);
  EVO_ASSERT(maybe_rec.has_value());
  validateRecordWithChanges(maybe_rec.value());

  std::variant<RecordWithChanges, int> rec_or_int;
  r.ReadUnionRecordWithChanges(rec_or_int);
  EVO_ASSERT(rec_or_int.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_int));

  // r.ReadAliasedUnionRecordWithChanges(rec_or_int);
  // EVO_ASSERT(rec_or_int.index() == 0);
  // validateRecordWithChanges(std::get<0>(rec_or_int));

  std::variant<RecordWithChanges, int, float, std::string> rec_or_int_or_flt_or_str;
  r.ReadUnionWithSameTypeset(rec_or_int_or_flt_or_str);
  EVO_ASSERT(rec_or_int_or_flt_or_str.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_int_or_flt_or_str));

  std::variant<RecordWithChanges, float> rec_or_flt;
  r.ReadUnionWithTypesAdded(rec_or_flt);
  EVO_ASSERT(rec_or_flt.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_flt));

  r.ReadUnionWithTypesRemoved(rec_or_int_or_flt_or_str);
  EVO_ASSERT(rec_or_int_or_flt_or_str.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_int_or_flt_or_str));

  r.ReadRecordToOptional(rec);
  validateRecordWithChanges(rec);
  r.ReadRecordToAliasedOptional(rec);
  validateRecordWithChanges(rec);
  r.ReadRecordToUnion(rec);
  validateRecordWithChanges(rec);
  r.ReadRecordToAliasedUnion(rec);
  validateRecordWithChanges(rec);

  r.ReadUnionToAliasedUnion(rec_or_int);
  EVO_ASSERT(rec_or_int.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_int));
  r.ReadUnionToAliasedUnionWithChanges(rec_or_int);
  EVO_ASSERT(rec_or_int.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_int));
  r.ReadOptionalToAliasedOptional(maybe_rec);
  EVO_ASSERT(maybe_rec.has_value());
  validateRecordWithChanges(maybe_rec.value());
  std::optional<int> maybe_int;
  r.ReadOptionalToAliasedOptionalWithChanges(maybe_int);
  EVO_ASSERT(maybe_int.has_value());
  EVO_ASSERT(maybe_int == INT_MIN);

  GenericRecord<int, std::string> generic_record;

  r.ReadGenericRecord(generic_record);
  validateGenericRecord(generic_record);
  r.ReadGenericRecordToOpenAlias(generic_record);
  validateGenericRecord(generic_record);
  r.ReadGenericRecordToClosedAlias(generic_record);
  validateGenericRecord(generic_record);
  r.ReadGenericRecordToHalfClosedAlias(generic_record);
  validateGenericRecord(generic_record);
  r.ReadAliasedGenericRecordToAlias(generic_record);
  validateGenericRecord(generic_record);

  r.ReadClosedGenericRecordToUnion(generic_record);
  validateGenericRecord(generic_record);
  r.ReadGenericRecordToAliasedUnion(generic_record);
  validateGenericRecord(generic_record);

  GenericUnion<GenericRecord<int, std::string>, float> generic_union_record;
  r.ReadGenericUnionOfChangedRecord(generic_union_record);
  assert(generic_union_record.index() == 0);
  validateGenericRecord(std::get<0>(generic_union_record));

  GenericParentRecord<int> generic_parent;
  r.ReadGenericParentRecord(generic_parent);

  validateGenericRecord(generic_parent.record);

  assert(generic_parent.record_of_union.field_1.index() == 0);
  assert(std::get<0>(generic_parent.record_of_union.field_1) == 42);
  assert(generic_parent.record_of_union.field_2 == "Hello, World");

  assert(generic_parent.union_of_record.index() == 0);
  validateGenericRecord(std::get<0>(generic_parent.union_of_record));

  GenericRecord<UnchangedGeneric<int>, ChangedGeneric<std::string, int>> generic_nested;
  r.ReadGenericNestedRecords(generic_nested);
  assert(generic_nested.field_1.field == 42);
  assert(generic_nested.field_2.y == "42");
  assert(generic_nested.field_2.z.field == 42);

  std::vector<AliasedClosedGenericRecord> generic_records(10);
  while (r.ReadGenericRecordStream(generic_records)) {
    EVO_ASSERT(generic_records.size() == 7);
    for (auto const& rec : generic_records) {
      validateGenericRecord(rec);
    }
  }

  std::vector<GenericParentRecord<int>> generic_parents(10);
  while (r.ReadGenericParentRecordStream(generic_parents)) {
    EVO_ASSERT(generic_parents.size() == 7);
    for (auto const& parent : generic_parents) {
      validateGenericParentRecord(parent);
    }
  }

  std::vector<RecordWithChanges> vec;
  r.ReadVectorRecordWithChanges(vec);
  EVO_ASSERT(vec.size() == 7);
  for (auto const& rec : vec) {
    validateRecordWithChanges(rec);
  }

  int count = 0;
  while (r.ReadStreamedRecordWithChanges(rec)) {
    validateRecordWithChanges(rec);
    count += 1;
  }

  EVO_ASSERT(count == 7);

  std::vector<std::string> strings;
  r.ReadAddedStringVector(strings);
  switch (r.GetVersion()) {
    case Version::v0:
    case Version::v1:
      EVO_ASSERT(strings.empty());
      break;

    default:
      EVO_ASSERT(strings.size() == 7);
      for (auto const& str : strings) {
        EVO_ASSERT(str == HelloWorld);
      }
  }

  r.ReadAddedOptional(maybe_rec);
  switch (r.GetVersion()) {
    case Version::v0:
      EVO_ASSERT(!maybe_rec.has_value());
      break;

    default:
      EVO_ASSERT(maybe_rec.has_value());
      validateRecordWithChanges(maybe_rec.value());
  }

  std::unordered_map<std::string, std::string> map;
  r.ReadAddedMap(map);
  switch (r.GetVersion()) {
    case Version::v0:
      EVO_ASSERT(map.empty());
      break;

    default:
      EVO_ASSERT(map.size() == 1);
      EVO_ASSERT(map["hello"] == "world");
  }

  std::variant<std::monostate, RecordWithChanges, std::string> nil_or_rec_or_str;
  r.ReadAddedUnion(nil_or_rec_or_str);
  switch (r.GetVersion()) {
    case Version::v0:
    case Version::v1:
      EVO_ASSERT(nil_or_rec_or_str.index() == 0);
      break;

    default:
      EVO_ASSERT(nil_or_rec_or_str.index() == 1);
      validateRecordWithChanges(std::get<1>(nil_or_rec_or_str));
  }

  std::vector<RecordWithChanges> records(10);
  switch (r.GetVersion()) {
    case Version::v0:
      EVO_ASSERT(r.ReadAddedRecordStream(records) == false);
      EVO_ASSERT(records.empty());
      break;

    default:
      EVO_ASSERT(r.ReadAddedRecordStream(records) == true);
      EVO_ASSERT(records.size() == 7);
      for (auto const& rec : records) {
        validateRecordWithChanges(rec);
      }
      EVO_ASSERT(r.ReadAddedRecordStream(records) == false);
  }

  std::variant<RecordWithChanges, RenamedRecord> rec_or_renamed;
  count = 0;
  while (r.ReadAddedUnionStream(rec_or_renamed)) {
    EVO_ASSERT(rec_or_renamed.index() == 0);
    validateRecordWithChanges(std::get<0>(rec_or_renamed));
    count += 1;
  }
  switch (r.GetVersion()) {
    case Version::v0:
    case Version::v1:
      EVO_ASSERT(count == 0);
      break;

    default:
      EVO_ASSERT(count == 7);
  }

  r.Close();

  return 0;
}
