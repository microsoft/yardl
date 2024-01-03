#include "lib/binary/protocols.h"

using namespace evo_test;

static std::string HelloWorld = "Hello, World!";

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

  UnchangedRecord unchanged;
  unchanged.name = "Jane Doe";
  unchanged.age = 42;
  unchanged.meta = {{"height", 161.3}, {"weight", 75.0f}};

  RecordWithChanges rec;
  rec.int_to_long = static_cast<long>(INT_MIN);
  rec.float_to_double = M_PI;
  rec.optional_long_to_string = std::to_string(LONG_MIN);
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
