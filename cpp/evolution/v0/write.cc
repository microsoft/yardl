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

  UnchangedRecord unchanged;
  unchanged.name = "Jane Doe";
  unchanged.age = 42;
  unchanged.meta = {{"height", 161.3}, {"weight", 75.0f}};

  RecordWithChanges rec;
  rec.int_to_long = INT_MIN;
  rec.deprecated_vector = {1, 2, 3};
  rec.float_to_double = M_PI;
  rec.deprecated_array = {4, 5, 6};
  rec.optional_long_to_string = LONG_MIN;
  rec.deprecated_map = {{"a", {1, 4, 7}}, {"b", {2, 5, 8}}, {"c", {3, 6, 9}}};
  rec.unchanged_record = unchanged;

  w.WriteRecordWithChanges(rec);
  w.WriteAliasedRecordWithChanges(rec);

  w.WriteOptionalRecordWithChanges(rec);
  w.WriteAliasedOptionalRecordWithChanges(rec);

  // Stream a total of 7 records
  w.WriteStreamedRecordWithChanges(rec);
  w.WriteStreamedRecordWithChanges(rec);
  std::vector<RecordWithChanges> recs(4, rec);
  w.WriteStreamedRecordWithChanges(recs);
  w.WriteStreamedRecordWithChanges(rec);
  w.EndStreamedRecordWithChanges();

  w.Close();

  return 0;
}
