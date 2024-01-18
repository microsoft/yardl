#include "../evolution_testing.h"
#include "lib/binary/protocols.h"

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
  EVO_ASSERT(uint32 == static_cast<uint32_t>(INT8_MIN));
  r.ReadInt8ToUlong(uint64);
  EVO_ASSERT(uint64 == static_cast<uint64_t>(INT8_MIN));
  r.ReadInt8ToFloat(flt);
  EVO_ASSERT_EQUALISH(flt, static_cast<float>(INT8_MIN));
  r.ReadInt8ToDouble(dbl);
  EVO_ASSERT_EQUALISH(dbl, static_cast<double>(INT8_MIN));

  r.ReadIntToUint(uint32);
  EVO_ASSERT(uint32 == static_cast<uint32_t>(INT_MIN));
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
