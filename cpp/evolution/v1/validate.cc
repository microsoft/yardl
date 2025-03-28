#include "../evolution_testing.h"
#include "generated/binary/protocols.h"

using namespace evo_test;

void validateRecordWithChanges(RecordWithChanges const& rec) {
  (void)(rec);

  EVO_ASSERT(rec.int_to_long == static_cast<long>(INT_MIN));
  EVO_ASSERT_EQUALISH(rec.float_to_double, M_PI);
  EVO_ASSERT(rec.optional_long_to_string.has_value());
  EVO_ASSERT(std::stol(rec.optional_long_to_string.value()) == LONG_MIN);

  EVO_ASSERT(rec.unchanged_record.name == "Jane Doe");
  EVO_ASSERT(rec.unchanged_record.age == 42);
  EVO_ASSERT(rec.unchanged_record.meta.at("height") == 161.3);
  EVO_ASSERT(rec.unchanged_record.meta.at("weight") == 75.0f);
}

#define validateGenericRecord(record)           \
  do {                                          \
    EVO_ASSERT(record.field_2 == HelloWorld);   \
    EVO_ASSERT(record.field_1 == 42);           \
    if (record.added.has_value()) {             \
      EVO_ASSERT(record.added.value() == true); \
    }                                           \
  } while (0)

#define validateGenericParentRecord(parent)                        \
  do {                                                             \
    validateGenericRecord(parent.record);                          \
    EVO_ASSERT(parent.record_of_union.field_1.index() == 0);       \
    EVO_ASSERT(std::get<0>(parent.record_of_union.field_1) == 42); \
    EVO_ASSERT(parent.record_of_union.field_2 == HelloWorld);      \
    EVO_ASSERT(parent.union_of_record.index() == 0);               \
    validateGenericRecord(std::get<0>(parent.union_of_record));    \
  } while (0)

int main(void) {
  ::binary::ProtocolWithChangesReader r(std::cin);

  int32_t int32 = 0;
  int64_t int64 = 0;
  uint32_t uint32 = 0;
  uint64_t uint64 = 0;
  float flt = 0.f;
  double dbl = 0.0;
  std::string str;

  r.ReadInt8ToInt(int32);
  EVO_ASSERT(int32 == INT8_MIN);
  r.ReadInt8ToLong(int64);
  EVO_ASSERT(int64 == INT8_MIN);
  r.ReadInt8ToUint(uint32);
  EVO_ASSERT(uint32 == static_cast<uint32_t>(INT8_MAX));
  r.ReadInt8ToUlong(uint64);
  EVO_ASSERT(uint64 == static_cast<uint64_t>(INT8_MAX));
  r.ReadInt8ToFloat(flt);
  EVO_ASSERT_EQUALISH(flt, static_cast<float>(INT8_MIN));
  r.ReadInt8ToDouble(dbl);
  EVO_ASSERT_EQUALISH(dbl, static_cast<double>(INT8_MIN));

  r.ReadIntToUint(uint32);
  EVO_ASSERT(uint32 == static_cast<uint32_t>(INT_MAX));
  r.ReadIntToLong(int64);
  EVO_ASSERT(int64 == static_cast<int64_t>(INT_MIN));
  r.ReadIntToFloat(flt);
  EVO_ASSERT_EQUALISH(flt, static_cast<float>(INT8_MIN));
  r.ReadIntToDouble(dbl);
  EVO_ASSERT_EQUALISH(dbl, static_cast<double>(INT8_MIN));

  r.ReadUintToUlong(uint64);
  EVO_ASSERT(uint64 == static_cast<uint32_t>(UINT_MAX));
  r.ReadUintToFloat(flt);
  EVO_ASSERT_EQUALISH(flt, static_cast<float>(UINT8_MAX));
  r.ReadUintToDouble(dbl);
  EVO_ASSERT_EQUALISH(dbl, static_cast<double>(UINT8_MAX));

  r.ReadFloatToDouble(dbl);
  EVO_ASSERT_EQUALISH(dbl, M_PI);

  std::complex<double> cxdbl;
  r.ReadComplexFloatToComplexDouble(cxdbl);
  EVO_ASSERT_EQUALISH(cxdbl.real(), M_PI);
  EVO_ASSERT_EQUALISH(cxdbl.imag(), -M_PI);

  r.ReadIntToString(str);
  EVO_ASSERT(std::stoi(str) == INT_MIN);
  r.ReadUintToString(str);
  EVO_ASSERT(std::stoul(str) == UINT_MAX);
  r.ReadLongToString(str);
  EVO_ASSERT(std::stol(str) == LONG_MIN);
  r.ReadUlongToString(str);
  EVO_ASSERT(std::stoul(str) == ULONG_MAX);
  r.ReadFloatToString(str);
  EVO_ASSERT_EQUALISH(std::stof(str), M_PI);
  r.ReadDoubleToString(str);
  EVO_ASSERT_EQUALISH(std::stod(str), M_PI);

  std::optional<int> maybe_int;
  r.ReadIntToOptional(maybe_int);
  if (maybe_int.has_value()) {
    EVO_ASSERT(maybe_int.value() == INT_MIN);
  }

  std::optional<float> maybe_float;
  r.ReadFloatToOptional(maybe_float);
  if (maybe_float.has_value()) {
    EVO_ASSERT_EQUALISH(maybe_float.value(), M_PI);
  }

  std::optional<std::string> maybe_string;
  r.ReadStringToOptional(maybe_string);
  if (maybe_string.has_value()) {
    EVO_ASSERT(maybe_string.value() == HelloWorld);
  }

  std::variant<int, bool> intOrBool;
  r.ReadIntToUnion(intOrBool);
  EVO_ASSERT(intOrBool.index() == 0);
  EVO_ASSERT(std::get<0>(intOrBool) == INT_MIN);

  std::variant<float, bool> floatOrBool;
  r.ReadFloatToUnion(floatOrBool);
  EVO_ASSERT(floatOrBool.index() == 0);
  EVO_ASSERT_EQUALISH(std::get<0>(floatOrBool), M_PI);

  std::variant<std::string, bool> stringOrBool;
  r.ReadStringToUnion(stringOrBool);
  EVO_ASSERT(stringOrBool.index() == 0);
  EVO_ASSERT(std::get<0>(stringOrBool) == HelloWorld);

  r.ReadOptionalIntToFloat(maybe_float);
  EVO_ASSERT(maybe_float.has_value());
  EVO_ASSERT_EQUALISH(maybe_float.value(), static_cast<float>(INT8_MIN));

  r.ReadOptionalFloatToString(maybe_string);
  EVO_ASSERT(maybe_string.has_value());
  EVO_ASSERT_EQUALISH(std::stof(maybe_string.value()), M_PI);

  r.ReadAliasedLongToString(str);
  EVO_ASSERT(std::stol(str) == LONG_MIN);
  r.ReadStringToAliasedString(str);
  EVO_ASSERT(str == HelloWorld);
  r.ReadStringToAliasedInt(int32);
  EVO_ASSERT(int32 == INT_MIN);

  AliasedEnum e;
  r.ReadEnumToAliasedEnum(e);
  EVO_ASSERT(e == GrowingEnum::kC);

  RecordWithChanges rec;

  std::variant<std::monostate, int, std::string> nullOrIntOrString;
  std::variant<std::monostate, RecordWithChanges, std::string> nullOrRecOrString;

  r.ReadOptionalIntToUnion(nullOrIntOrString);
  EVO_ASSERT(nullOrIntOrString.index() == 1);
  EVO_ASSERT(std::get<1>(nullOrIntOrString) == INT_MIN);

  r.ReadOptionalRecordToUnion(nullOrRecOrString);
  EVO_ASSERT(nullOrRecOrString.index() == 1);
  validateRecordWithChanges(std::get<1>(nullOrRecOrString));

  r.ReadRecordWithChanges(rec);
  validateRecordWithChanges(rec);

  r.ReadAliasedRecordWithChanges(rec);
  validateRecordWithChanges(rec);

  RenamedRecord renamed;
  r.ReadRecordToRenamedRecord(renamed);
  EVO_ASSERT(renamed.s == HelloWorld);
  EVO_ASSERT(renamed.i == INT_MIN);

  r.ReadRecordToAliasedRecord(rec);
  validateRecordWithChanges(rec);

  r.ReadRecordToAliasedAlias(rec);
  validateRecordWithChanges(rec);

  std::vector<std::string> str_vec(10);
  while (r.ReadStreamIntToStringToFloat(str_vec)) {
    for (auto& s : str_vec) {
      EVO_ASSERT(std::stoi(s) == 42);
    }
  }

  str_vec.clear();
  r.ReadVectorIntToStringToFloat(str_vec);
  for (auto& s : str_vec) {
    EVO_ASSERT(std::stoi(s) == 42);
  }

  std::variant<float, int> int_float;
  r.ReadIntFloatUnionReordered(int_float);
  EVO_ASSERT(int_float.index() == 0);
  EVO_ASSERT(std::get<0>(int_float) == static_cast<float>(M_PI));

  std::vector<std::variant<float, int>> int_float_vec;
  r.ReadVectorUnionReordered(int_float_vec);
  for (auto& v : int_float_vec) {
    EVO_ASSERT(v.index() == 0);
    EVO_ASSERT(std::get<0>(v) == static_cast<float>(M_PI));
  }

  std::vector<std::variant<std::string, int>> int_string_vec(10);
  while (r.ReadStreamUnionReordered(int_string_vec)) {
    for (auto& v : int_string_vec) {
      EVO_ASSERT(v.index() == 0);
      EVO_ASSERT(std::get<0>(v) == HelloWorld);
    }
  }

  std::vector<int> int_vec(10);
  while (r.ReadIntToUnionStream(int_vec)) {
    for (auto& i : int_vec) {
      EVO_ASSERT(i == 42);
    }
  }

  std::variant<int, bool> int_bool;
  while (r.ReadUnionStreamTypeChange(int_bool)) {
    EVO_ASSERT(int_bool.index() == 0);
    EVO_ASSERT(std::get<0>(int_bool) == 42);
  }

  std::vector<StreamItem> stream_items(10);
  while (r.ReadStreamOfAliasTypeChange(stream_items)) {
    EVO_ASSERT(stream_items.size() == 7);
    for (auto const& item : stream_items) {
      EVO_ASSERT(item.index() == 0);
      validateRecordWithChanges(std::get<0>(item));
    }
  }

  RZ record;
  r.ReadRlink(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRlinkRX(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRlinkRY(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRlinkRZ(record);
  EVO_ASSERT(record.subject == 42);

  r.ReadRaRLink(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRaRX(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRaRY(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRaRZ(record);
  EVO_ASSERT(record.subject == 42);

  r.ReadRbRLink(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRbRX(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRbRY(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRbRZ(record);
  EVO_ASSERT(record.subject == 42);

  r.ReadRcRLink(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRcRX(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRcRY(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRcRZ(record);
  EVO_ASSERT(record.subject == 42);

  r.ReadRlinkRNew(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRaRNew(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRbRNew(record);
  EVO_ASSERT(record.subject == 42);
  r.ReadRcRNew(record);
  EVO_ASSERT(record.subject == 42);

  RUnion runion;
  r.ReadRlinkRUnion(runion);
  EVO_ASSERT(runion.index() == 0);
  EVO_ASSERT(std::get<0>(runion).subject == 42);
  r.ReadRaRUnion(runion);
  EVO_ASSERT(runion.index() == 0);
  EVO_ASSERT(std::get<0>(runion).subject == 42);
  r.ReadRbRUnion(runion);
  EVO_ASSERT(runion.index() == 0);
  EVO_ASSERT(std::get<0>(runion).subject == 42);
  r.ReadRcRUnion(runion);
  EVO_ASSERT(runion.index() == 0);
  EVO_ASSERT(std::get<0>(runion).subject == 42);

  std::optional<RecordWithChanges> maybe_rec;
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

  std::variant<float, RecordWithChanges, std::string, int> flt_or_rec_or_str_or_int;
  r.ReadUnionWithSameTypeset(flt_or_rec_or_str_or_int);
  EVO_ASSERT(flt_or_rec_or_str_or_int.index() == 1);
  validateRecordWithChanges(std::get<1>(flt_or_rec_or_str_or_int));

  std::variant<RecordWithChanges, int, float, std::string> rec_or_int_or_flt_or_str;
  r.ReadUnionWithTypesAdded(rec_or_int_or_flt_or_str);
  EVO_ASSERT(rec_or_int_or_flt_or_str.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_int_or_flt_or_str));

  std::variant<RecordWithChanges, std::string> rec_or_str;
  r.ReadUnionWithTypesRemoved(rec_or_str);
  EVO_ASSERT(rec_or_str.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_str));

  AliasedOptionalRecord aliased_rec;
  r.ReadRecordToOptional(aliased_rec);
  EVO_ASSERT(aliased_rec.has_value());
  validateRecordWithChanges(aliased_rec.value());

  r.ReadRecordToAliasedOptional(aliased_rec);
  EVO_ASSERT(aliased_rec.has_value());
  validateRecordWithChanges(aliased_rec.value());

  r.ReadRecordToUnion(rec_or_str);
  EVO_ASSERT(rec_or_str.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_str));

  AliasedRecordOrString aliased_rec_or_str;
  r.ReadRecordToAliasedUnion(aliased_rec_or_str);
  EVO_ASSERT(aliased_rec_or_str.index() == 0);
  validateRecordWithChanges(std::get<0>(aliased_rec_or_str));

  r.ReadUnionToAliasedUnion(rec_or_int);
  EVO_ASSERT(rec_or_int.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_int));
  r.ReadUnionToAliasedUnionWithChanges(rec_or_str);
  EVO_ASSERT(rec_or_str.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_str));
  r.ReadOptionalToAliasedOptional(maybe_rec);
  EVO_ASSERT(maybe_rec.has_value());
  validateRecordWithChanges(maybe_rec.value());
  std::optional<std::string> maybe_str;
  r.ReadOptionalToAliasedOptionalWithChanges(maybe_str);
  EVO_ASSERT(maybe_str.has_value());
  EVO_ASSERT(maybe_str == std::to_string(INT_MIN));

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
  r.ReadGenericRecordToReversed(generic_record);
  validateGenericRecord(generic_record);

  std::variant<GenericRecord<int, std::string>, std::string> generic_record_or_string;
  r.ReadClosedGenericRecordToUnion(generic_record_or_string);
  EVO_ASSERT(generic_record_or_string.index() == 0);
  validateGenericRecord(std::get<0>(generic_record_or_string));
  r.ReadGenericRecordToAliasedUnion(generic_record_or_string);
  EVO_ASSERT(generic_record_or_string.index() == 0);
  validateGenericRecord(std::get<0>(generic_record_or_string));

  GenericUnion<GenericRecord<int, std::string>, float> generic_union_record;
  r.ReadGenericUnionToReversed(generic_union_record);
  EVO_ASSERT(generic_union_record.index() == 0);
  validateGenericRecord(std::get<0>(generic_union_record));
  r.ReadGenericUnionOfChangedRecord(generic_union_record);
  EVO_ASSERT(generic_union_record.index() == 0);
  validateGenericRecord(std::get<0>(generic_union_record));

  GenericParentRecord<int> generic_parent;
  r.ReadGenericParentRecord(generic_parent);

  validateGenericRecord(generic_parent.record);

  EVO_ASSERT(generic_parent.record_of_union.field_1.index() == 0);
  EVO_ASSERT(std::get<0>(generic_parent.record_of_union.field_1) == 42);
  EVO_ASSERT(generic_parent.record_of_union.field_2 == HelloWorld);

  EVO_ASSERT(generic_parent.union_of_record.index() == 0);
  validateGenericRecord(std::get<0>(generic_parent.union_of_record));

  GenericRecord<Unchanged, Changed> generic_nested;
  r.ReadGenericNestedRecords(generic_nested);
  EVO_ASSERT(generic_nested.field_1.field == 42);
  EVO_ASSERT(generic_nested.field_2.y.has_value());
  EVO_ASSERT(generic_nested.field_2.y.value() == "42");
  EVO_ASSERT(generic_nested.field_2.z.has_value());
  EVO_ASSERT(generic_nested.field_2.z.value().field == 42);

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

  r.Close();

  return 0;
}
