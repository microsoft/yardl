#include "lib/binary/protocols.h"

using namespace evo_test;

static std::string HelloWorld = "Hello, World!";

#define assertCloseEnough(a, b) assert(std::abs(a - b) < 0.0001)

void validateRecordWithChanges(RecordWithChanges const& rec) {
  (void)(rec);

  assert(rec.int_to_long == INT_MIN);
  assertCloseEnough(rec.float_to_double, M_PI);
  assert(rec.optional_long_to_string.has_value());
  assert(rec.optional_long_to_string.value() == LONG_MIN);

  assert(rec.unchanged_record.name == "Jane Doe");
  assert(rec.unchanged_record.age == 42);
  assert(rec.unchanged_record.meta.at("height") == 161.3);
  assert(rec.unchanged_record.meta.at("weight") == 75.0f);
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
  assert(int8 == INT8_MIN);
  r.ReadInt8ToLong(int8);
  assert(int8 == INT8_MIN);
  r.ReadInt8ToUint(int8);
  assert(int8 == INT8_MIN);
  r.ReadInt8ToUlong(int8);
  assert(int8 == INT8_MIN);
  r.ReadInt8ToFloat(int8);
  assert(int8 == INT8_MIN);
  r.ReadInt8ToDouble(int8);
  assert(int8 == INT8_MIN);

  r.ReadIntToUint(int32);
  assert(int32 == INT_MIN);
  r.ReadIntToLong(int32);
  assert(int32 == INT_MIN);
  r.ReadIntToFloat(int32);
  assert(int32 == INT8_MIN);
  r.ReadIntToDouble(int32);
  assert(int32 == INT8_MIN);

  r.ReadUintToUlong(uint32);
  assert(uint32 == UINT_MAX);
  r.ReadUintToFloat(uint32);
  assert(uint32 == UINT8_MAX);
  r.ReadUintToDouble(uint32);
  assert(uint32 == UINT8_MAX);

  r.ReadFloatToDouble(flt);
  assertCloseEnough(flt, M_PI);

  r.ReadIntToString(int32);
  assert(int32 == INT_MIN);
  r.ReadUintToString(uint32);
  assert(uint32 == UINT_MAX);
  r.ReadLongToString(int64);
  assert(int64 == LONG_MIN);
  r.ReadUlongToString(uint64);
  assert(uint64 == ULONG_MAX);
  r.ReadFloatToString(flt);
  assertCloseEnough(flt, M_PI);
  r.ReadDoubleToString(dbl);
  assertCloseEnough(dbl, M_PI);

  r.ReadIntToOptional(int32);
  assert(int32 == INT_MIN);
  r.ReadFloatToOptional(flt);
  assertCloseEnough(flt, M_PI);

  r.ReadStringToOptional(str);
  assert(str == HelloWorld);

  r.ReadIntToUnion(int32);
  assert(int32 == INT_MIN);
  r.ReadFloatToUnion(flt);
  assertCloseEnough(flt, M_PI);
  r.ReadStringToUnion(str);
  assert(str == HelloWorld);

  std::optional<int32_t> maybe_int32;
  r.ReadOptionalIntToFloat(maybe_int32);
  assert(maybe_int32.has_value());
  assert(maybe_int32 == INT8_MIN);

  std::optional<float> maybe_flt;
  r.ReadOptionalFloatToString(maybe_flt);
  assert(maybe_flt.has_value());
  assertCloseEnough(maybe_flt.value(), M_PI);

  r.ReadAliasedLongToString(int64);
  assert(int64 == LONG_MIN);

  RecordWithChanges rec;
  std::optional<RecordWithChanges> maybe_rec;

  r.ReadOptionalIntToUnion(maybe_int32);
  assert(maybe_int32.has_value());
  assert(maybe_int32 == INT_MIN);

  r.ReadOptionalRecordToUnion(maybe_rec);
  assert(maybe_rec.has_value());
  validateRecordWithChanges(maybe_rec.value());

  r.ReadRecordWithChanges(rec);
  validateRecordWithChanges(rec);

  r.ReadAliasedRecordWithChanges(rec);
  validateRecordWithChanges(rec);

  r.ReadOptionalRecordWithChanges(maybe_rec);
  assert(maybe_rec.has_value());
  validateRecordWithChanges(maybe_rec.value());

  r.ReadAliasedOptionalRecordWithChanges(maybe_rec);
  assert(maybe_rec.has_value());
  validateRecordWithChanges(maybe_rec.value());

  std::variant<RecordWithChanges, int> rec_or_int;
  r.ReadUnionRecordWithChanges(rec_or_int);
  assert(rec_or_int.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_int));

  // r.ReadAliasedUnionRecordWithChanges(rec_or_int);
  // assert(rec_or_int.index() == 0);
  // validateRecordWithChanges(std::get<0>(rec_or_int));

  std::variant<RecordWithChanges, int, float, std::string> rec_or_int_or_flt_or_str;
  r.ReadUnionWithSameTypeset(rec_or_int_or_flt_or_str);
  assert(rec_or_int_or_flt_or_str.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_int_or_flt_or_str));

  std::variant<RecordWithChanges, float> rec_or_flt;
  r.ReadUnionWithTypesAdded(rec_or_flt);
  assert(rec_or_flt.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_flt));

  r.ReadUnionWithTypesRemoved(rec_or_int_or_flt_or_str);
  assert(rec_or_int_or_flt_or_str.index() == 0);
  validateRecordWithChanges(std::get<0>(rec_or_int_or_flt_or_str));

  std::vector<RecordWithChanges> vec;
  r.ReadVectorRecordWithChanges(vec);
  assert(vec.size() == 7);
  for (auto const& rec : vec) {
    validateRecordWithChanges(rec);
  }

  int count = 0;
  while (r.ReadStreamedRecordWithChanges(rec)) {
    validateRecordWithChanges(rec);
    count += 1;
  }

  assert(count == 7);

  r.Close();

  return 0;
}
