#include "../evolution_testing.h"
#include "lib/binary/protocols.h"

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

  RC record;
  r.ReadRlink(record);
  assert(record.subject == "42");
  r.ReadRlinkRX(record);
  assert(record.subject == "42");
  r.ReadRlinkRY(record);
  assert(record.subject == "42");
  r.ReadRlinkRZ(record);
  assert(record.subject == "42");

  r.ReadRaRLink(record);
  assert(record.subject == "42");
  r.ReadRaRX(record);
  assert(record.subject == "42");
  r.ReadRaRY(record);
  assert(record.subject == "42");
  r.ReadRaRZ(record);
  assert(record.subject == "42");

  r.ReadRbRLink(record);
  assert(record.subject == "42");
  r.ReadRbRX(record);
  assert(record.subject == "42");
  r.ReadRbRY(record);
  assert(record.subject == "42");
  r.ReadRbRZ(record);
  assert(record.subject == "42");

  r.ReadRcRLink(record);
  assert(record.subject == "42");
  r.ReadRcRX(record);
  assert(record.subject == "42");
  r.ReadRcRY(record);
  assert(record.subject == "42");
  r.ReadRcRZ(record);
  assert(record.subject == "42");

  r.ReadRlinkRNew(record);
  assert(record.subject == "42");
  r.ReadRaRNew(record);
  assert(record.subject == "42");
  r.ReadRbRNew(record);
  assert(record.subject == "42");
  r.ReadRcRNew(record);
  assert(record.subject == "42");

  r.ReadRlinkRUnion(record);
  assert(record.subject == "42");
  r.ReadRaRUnion(record);
  assert(record.subject == "42");
  r.ReadRbRUnion(record);
  assert(record.subject == "42");
  r.ReadRcRUnion(record);
  assert(record.subject == "42");

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

  r.Close();

  return 0;
}
