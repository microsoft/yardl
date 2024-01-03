#include "lib/binary/protocols.h"

using namespace evo_test;

static std::string HelloWorld = "Hello, World!";

#define assertCloseEnough(a, b) assert(std::abs(a - b) < 0.0001)

void validateRecordWithChanges(RecordWithChanges const& rec) {
  assert(rec.int_to_long == static_cast<long>(INT_MIN));
  assertCloseEnough(rec.float_to_double, M_PI);
  assert(rec.optional_long_to_string.has_value());
  assert(std::stol(rec.optional_long_to_string.value()) == LONG_MIN);

  assert(rec.unchanged_record.name == "Jane Doe");
  assert(rec.unchanged_record.age == 42);
  assert(rec.unchanged_record.meta.at("height") == 161.3);
  assert(rec.unchanged_record.meta.at("weight") == 75.0f);
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
  assert(int32 == INT8_MIN);
  r.ReadInt8ToLong(int64);
  assert(int64 == INT8_MIN);
  r.ReadInt8ToUint(uint32);
  assert(uint32 == static_cast<uint32_t>(INT8_MIN));
  r.ReadInt8ToUlong(uint64);
  assert(uint64 == static_cast<uint64_t>(INT8_MIN));
  r.ReadInt8ToFloat(flt);
  assertCloseEnough(flt, static_cast<float>(INT8_MIN));
  r.ReadInt8ToDouble(dbl);
  assertCloseEnough(dbl, static_cast<double>(INT8_MIN));

  r.ReadIntToUint(uint32);
  assert(uint32 == static_cast<uint32_t>(INT_MIN));
  r.ReadIntToLong(int64);
  assert(int64 == static_cast<int64_t>(INT_MIN));
  r.ReadIntToFloat(flt);
  assertCloseEnough(flt, static_cast<float>(INT8_MIN));
  r.ReadIntToDouble(dbl);
  assertCloseEnough(dbl, static_cast<double>(INT8_MIN));

  r.ReadUintToUlong(uint64);
  assert(uint64 == static_cast<uint32_t>(UINT_MAX));
  r.ReadUintToFloat(flt);
  assertCloseEnough(flt, static_cast<float>(UINT8_MAX));
  r.ReadUintToDouble(dbl);
  assertCloseEnough(dbl, static_cast<double>(UINT8_MAX));

  r.ReadFloatToDouble(dbl);
  assertCloseEnough(dbl, M_PI);

  r.ReadIntToString(str);
  assert(std::stoi(str) == INT_MIN);
  r.ReadUintToString(str);
  assert(std::stoul(str) == UINT_MAX);
  r.ReadLongToString(str);
  assert(std::stol(str) == LONG_MIN);
  r.ReadUlongToString(str);
  assert(std::stoul(str) == ULONG_MAX);
  r.ReadFloatToString(str);
  assertCloseEnough(std::stof(str), M_PI);
  r.ReadDoubleToString(str);
  assertCloseEnough(std::stod(str), M_PI);

  std::optional<int> maybe_int;
  r.ReadIntToOptional(maybe_int);
  if (maybe_int.has_value()) {
    assert(maybe_int.value() == INT_MIN);
  }

  std::optional<float> maybe_float;
  r.ReadFloatToOptional(maybe_float);
  if (maybe_float.has_value()) {
    assertCloseEnough(maybe_float.value(), M_PI);
  }

  std::optional<std::string> maybe_string;
  r.ReadStringToOptional(maybe_string);
  if (maybe_string.has_value()) {
    assert(maybe_string.value() == HelloWorld);
  }

  std::variant<int, bool> intOrBool;
  r.ReadIntToUnion(intOrBool);
  assert(intOrBool.index() == 0);
  assert(std::get<0>(intOrBool) == INT_MIN);

  std::variant<float, bool> floatOrBool;
  r.ReadFloatToUnion(floatOrBool);
  assert(floatOrBool.index() == 0);
  assertCloseEnough(std::get<0>(floatOrBool), M_PI);

  std::variant<std::string, bool> stringOrBool;
  r.ReadStringToUnion(stringOrBool);
  assert(stringOrBool.index() == 0);
  assert(std::get<0>(stringOrBool) == HelloWorld);

  r.ReadOptionalIntToFloat(maybe_float);
  assert(maybe_float.has_value());
  assertCloseEnough(maybe_float.value(), static_cast<float>(INT8_MIN));

  r.ReadOptionalFloatToString(maybe_string);
  assert(maybe_string.has_value());
  assertCloseEnough(std::stof(maybe_string.value()), M_PI);

  r.ReadAliasedLongToString(str);
  assert(std::stol(str) == LONG_MIN);

  RecordWithChanges rec;
  r.ReadRecordWithChanges(rec);
  validateRecordWithChanges(rec);

  r.ReadAliasedRecordWithChanges(rec);
  validateRecordWithChanges(rec);

  std::optional<RecordWithChanges> maybe_rec;
  r.ReadOptionalRecordWithChanges(maybe_rec);
  assert(maybe_rec.has_value());
  validateRecordWithChanges(maybe_rec.value());

  r.ReadAliasedOptionalRecordWithChanges(maybe_rec);
  assert(maybe_rec.has_value());
  validateRecordWithChanges(maybe_rec.value());

  int count = 0;
  while (r.ReadStreamedRecordWithChanges(rec)) {
    validateRecordWithChanges(rec);
    count += 1;
  }

  assert(count == 7);

  r.Close();

  return 0;
}
