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
  w.WriteStringToAliasedString(HelloWorld);
  w.WriteStringToAliasedInt(INT_MIN);

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
