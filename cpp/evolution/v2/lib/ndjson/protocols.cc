// This file was generated by the "yardl" tool. DO NOT EDIT.

#include "../yardl/detail/ndjson/serializers.h"
#include "protocols.h"

namespace evo_test {
using ordered_json = nlohmann::ordered_json;

[[maybe_unused]] static void to_json(ordered_json& j, evo_test::UnchangedRecord const& value);
[[maybe_unused]] static void from_json(ordered_json const& j, evo_test::UnchangedRecord& value);

[[maybe_unused]] static void to_json(ordered_json& j, evo_test::RecordWithChanges const& value);
[[maybe_unused]] static void from_json(ordered_json const& j, evo_test::RecordWithChanges& value);

[[maybe_unused]] static void to_json(ordered_json& j, evo_test::RenamedRecord const& value);
[[maybe_unused]] static void from_json(ordered_json const& j, evo_test::RenamedRecord& value);

} // namespace evo_test

NLOHMANN_JSON_NAMESPACE_BEGIN

template <>
struct adl_serializer<std::variant<evo_test::RecordWithChanges, int32_t>> {
  [[maybe_unused]] static void to_json(ordered_json& j, std::variant<evo_test::RecordWithChanges, int32_t> const& value) {
    std::visit([&j](auto const& v) {j = v;}, value);
  }

  [[maybe_unused]] static void from_json(ordered_json const& j, std::variant<evo_test::RecordWithChanges, int32_t>& value) {
    if ((j.is_object())) {
      value = j.get<evo_test::RecordWithChanges>();
      return;
    }
    if ((j.is_number())) {
      value = j.get<int32_t>();
      return;
    }
    throw std::runtime_error("Invalid union value");
  }
};

template <>
struct adl_serializer<std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>> {
  [[maybe_unused]] static void to_json(ordered_json& j, std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value) {
    switch (value.index()) {
      case 0:
        j = ordered_json{ {"RecordWithChanges", std::get<evo_test::RecordWithChanges>(value)} };
        break;
      case 1:
        j = ordered_json{ {"int32", std::get<int32_t>(value)} };
        break;
      case 2:
        j = ordered_json{ {"float32", std::get<float>(value)} };
        break;
      case 3:
        j = ordered_json{ {"string", std::get<std::string>(value)} };
        break;
      default:
        throw std::runtime_error("Invalid union value");
    }
  }

  [[maybe_unused]] static void from_json(ordered_json const& j, std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value) {
    auto it = j.begin();
    std::string tag = it.key();
    if (tag == "RecordWithChanges") {
      value = it.value().get<evo_test::RecordWithChanges>();
      return;
    }
    if (tag == "int32") {
      value = it.value().get<int32_t>();
      return;
    }
    if (tag == "float32") {
      value = it.value().get<float>();
      return;
    }
    if (tag == "string") {
      value = it.value().get<std::string>();
      return;
    }
  }
};

template <>
struct adl_serializer<std::variant<evo_test::RecordWithChanges, float>> {
  [[maybe_unused]] static void to_json(ordered_json& j, std::variant<evo_test::RecordWithChanges, float> const& value) {
    std::visit([&j](auto const& v) {j = v;}, value);
  }

  [[maybe_unused]] static void from_json(ordered_json const& j, std::variant<evo_test::RecordWithChanges, float>& value) {
    if ((j.is_object())) {
      value = j.get<evo_test::RecordWithChanges>();
      return;
    }
    if ((j.is_number())) {
      value = j.get<float>();
      return;
    }
    throw std::runtime_error("Invalid union value");
  }
};

template <>
struct adl_serializer<std::variant<std::monostate, evo_test::RecordWithChanges, std::string>> {
  [[maybe_unused]] static void to_json(ordered_json& j, std::variant<std::monostate, evo_test::RecordWithChanges, std::string> const& value) {
    std::visit([&j](auto const& v) {j = v;}, value);
  }

  [[maybe_unused]] static void from_json(ordered_json const& j, std::variant<std::monostate, evo_test::RecordWithChanges, std::string>& value) {
    if ((j.is_null())) {
      value = j.get<std::monostate>();
      return;
    }
    if ((j.is_object())) {
      value = j.get<evo_test::RecordWithChanges>();
      return;
    }
    if ((j.is_string())) {
      value = j.get<std::string>();
      return;
    }
    throw std::runtime_error("Invalid union value");
  }
};

template <>
struct adl_serializer<std::variant<evo_test::RecordWithChanges, evo_test::RenamedRecord>> {
  [[maybe_unused]] static void to_json(ordered_json& j, std::variant<evo_test::RecordWithChanges, evo_test::RenamedRecord> const& value) {
    switch (value.index()) {
      case 0:
        j = ordered_json{ {"RecordWithChanges", std::get<evo_test::RecordWithChanges>(value)} };
        break;
      case 1:
        j = ordered_json{ {"RenamedRecord", std::get<evo_test::RenamedRecord>(value)} };
        break;
      default:
        throw std::runtime_error("Invalid union value");
    }
  }

  [[maybe_unused]] static void from_json(ordered_json const& j, std::variant<evo_test::RecordWithChanges, evo_test::RenamedRecord>& value) {
    auto it = j.begin();
    std::string tag = it.key();
    if (tag == "RecordWithChanges") {
      value = it.value().get<evo_test::RecordWithChanges>();
      return;
    }
    if (tag == "RenamedRecord") {
      value = it.value().get<evo_test::RenamedRecord>();
      return;
    }
  }
};

NLOHMANN_JSON_NAMESPACE_END

namespace evo_test {
using ordered_json = nlohmann::ordered_json;

[[maybe_unused]] static void to_json(ordered_json& j, evo_test::UnchangedRecord const& value) {
  j = ordered_json::object();
  if (yardl::ndjson::ShouldSerializeFieldValue(value.name)) {
    j.push_back({"name", value.name});
  }
  if (yardl::ndjson::ShouldSerializeFieldValue(value.age)) {
    j.push_back({"age", value.age});
  }
  if (yardl::ndjson::ShouldSerializeFieldValue(value.meta)) {
    j.push_back({"meta", value.meta});
  }
}

[[maybe_unused]] static void from_json(ordered_json const& j, evo_test::UnchangedRecord& value) {
  if (auto it = j.find("name"); it != j.end()) {
    it->get_to(value.name);
  }
  if (auto it = j.find("age"); it != j.end()) {
    it->get_to(value.age);
  }
  if (auto it = j.find("meta"); it != j.end()) {
    it->get_to(value.meta);
  }
}

[[maybe_unused]] static void to_json(ordered_json& j, evo_test::RecordWithChanges const& value) {
  j = ordered_json::object();
  if (yardl::ndjson::ShouldSerializeFieldValue(value.int_to_long)) {
    j.push_back({"intToLong", value.int_to_long});
  }
  if (yardl::ndjson::ShouldSerializeFieldValue(value.deprecated_vector)) {
    j.push_back({"deprecatedVector", value.deprecated_vector});
  }
  if (yardl::ndjson::ShouldSerializeFieldValue(value.float_to_double)) {
    j.push_back({"floatToDouble", value.float_to_double});
  }
  if (yardl::ndjson::ShouldSerializeFieldValue(value.deprecated_array)) {
    j.push_back({"deprecatedArray", value.deprecated_array});
  }
  if (yardl::ndjson::ShouldSerializeFieldValue(value.optional_long_to_string)) {
    j.push_back({"optionalLongToString", value.optional_long_to_string});
  }
  if (yardl::ndjson::ShouldSerializeFieldValue(value.deprecated_map)) {
    j.push_back({"deprecatedMap", value.deprecated_map});
  }
  if (yardl::ndjson::ShouldSerializeFieldValue(value.unchanged_record)) {
    j.push_back({"unchangedRecord", value.unchanged_record});
  }
}

[[maybe_unused]] static void from_json(ordered_json const& j, evo_test::RecordWithChanges& value) {
  if (auto it = j.find("intToLong"); it != j.end()) {
    it->get_to(value.int_to_long);
  }
  if (auto it = j.find("deprecatedVector"); it != j.end()) {
    it->get_to(value.deprecated_vector);
  }
  if (auto it = j.find("floatToDouble"); it != j.end()) {
    it->get_to(value.float_to_double);
  }
  if (auto it = j.find("deprecatedArray"); it != j.end()) {
    it->get_to(value.deprecated_array);
  }
  if (auto it = j.find("optionalLongToString"); it != j.end()) {
    it->get_to(value.optional_long_to_string);
  }
  if (auto it = j.find("deprecatedMap"); it != j.end()) {
    it->get_to(value.deprecated_map);
  }
  if (auto it = j.find("unchangedRecord"); it != j.end()) {
    it->get_to(value.unchanged_record);
  }
}

[[maybe_unused]] static void to_json(ordered_json& j, evo_test::RenamedRecord const& value) {
  j = ordered_json::object();
  if (yardl::ndjson::ShouldSerializeFieldValue(value.i)) {
    j.push_back({"i", value.i});
  }
  if (yardl::ndjson::ShouldSerializeFieldValue(value.s)) {
    j.push_back({"s", value.s});
  }
}

[[maybe_unused]] static void from_json(ordered_json const& j, evo_test::RenamedRecord& value) {
  if (auto it = j.find("i"); it != j.end()) {
    it->get_to(value.i);
  }
  if (auto it = j.find("s"); it != j.end()) {
    it->get_to(value.s);
  }
}

} // namespace evo_test

namespace evo_test::ndjson {
void ProtocolWithChangesWriter::WriteInt8ToIntImpl(int8_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "int8ToInt", json_value);}

void ProtocolWithChangesWriter::WriteInt8ToLongImpl(int8_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "int8ToLong", json_value);}

void ProtocolWithChangesWriter::WriteInt8ToUintImpl(int8_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "int8ToUint", json_value);}

void ProtocolWithChangesWriter::WriteInt8ToUlongImpl(int8_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "int8ToUlong", json_value);}

void ProtocolWithChangesWriter::WriteInt8ToFloatImpl(int8_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "int8ToFloat", json_value);}

void ProtocolWithChangesWriter::WriteInt8ToDoubleImpl(int8_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "int8ToDouble", json_value);}

void ProtocolWithChangesWriter::WriteIntToUintImpl(int32_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "intToUint", json_value);}

void ProtocolWithChangesWriter::WriteIntToLongImpl(int32_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "intToLong", json_value);}

void ProtocolWithChangesWriter::WriteIntToFloatImpl(int32_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "intToFloat", json_value);}

void ProtocolWithChangesWriter::WriteIntToDoubleImpl(int32_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "intToDouble", json_value);}

void ProtocolWithChangesWriter::WriteUintToUlongImpl(uint32_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "uintToUlong", json_value);}

void ProtocolWithChangesWriter::WriteUintToFloatImpl(uint32_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "uintToFloat", json_value);}

void ProtocolWithChangesWriter::WriteUintToDoubleImpl(uint32_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "uintToDouble", json_value);}

void ProtocolWithChangesWriter::WriteFloatToDoubleImpl(float const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "floatToDouble", json_value);}

void ProtocolWithChangesWriter::WriteIntToStringImpl(int32_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "intToString", json_value);}

void ProtocolWithChangesWriter::WriteUintToStringImpl(uint32_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "uintToString", json_value);}

void ProtocolWithChangesWriter::WriteLongToStringImpl(int64_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "longToString", json_value);}

void ProtocolWithChangesWriter::WriteUlongToStringImpl(uint64_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "ulongToString", json_value);}

void ProtocolWithChangesWriter::WriteFloatToStringImpl(float const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "floatToString", json_value);}

void ProtocolWithChangesWriter::WriteDoubleToStringImpl(double const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "doubleToString", json_value);}

void ProtocolWithChangesWriter::WriteIntToOptionalImpl(int32_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "intToOptional", json_value);}

void ProtocolWithChangesWriter::WriteFloatToOptionalImpl(float const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "floatToOptional", json_value);}

void ProtocolWithChangesWriter::WriteStringToOptionalImpl(std::string const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "stringToOptional", json_value);}

void ProtocolWithChangesWriter::WriteIntToUnionImpl(int32_t const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "intToUnion", json_value);}

void ProtocolWithChangesWriter::WriteFloatToUnionImpl(float const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "floatToUnion", json_value);}

void ProtocolWithChangesWriter::WriteStringToUnionImpl(std::string const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "stringToUnion", json_value);}

void ProtocolWithChangesWriter::WriteOptionalIntToFloatImpl(std::optional<int32_t> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "optionalIntToFloat", json_value);}

void ProtocolWithChangesWriter::WriteOptionalFloatToStringImpl(std::optional<float> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "optionalFloatToString", json_value);}

void ProtocolWithChangesWriter::WriteAliasedLongToStringImpl(evo_test::AliasedLongToString const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "aliasedLongToString", json_value);}

void ProtocolWithChangesWriter::WriteStringToAliasedStringImpl(std::string const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "stringToAliasedString", json_value);}

void ProtocolWithChangesWriter::WriteStringToAliasedIntImpl(std::string const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "stringToAliasedInt", json_value);}

void ProtocolWithChangesWriter::WriteOptionalIntToUnionImpl(std::optional<int32_t> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "optionalIntToUnion", json_value);}

void ProtocolWithChangesWriter::WriteOptionalRecordToUnionImpl(std::optional<evo_test::RecordWithChanges> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "optionalRecordToUnion", json_value);}

void ProtocolWithChangesWriter::WriteRecordWithChangesImpl(evo_test::RecordWithChanges const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "recordWithChanges", json_value);}

void ProtocolWithChangesWriter::WriteAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "aliasedRecordWithChanges", json_value);}

void ProtocolWithChangesWriter::WriteRecordToRenamedRecordImpl(evo_test::RenamedRecord const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "recordToRenamedRecord", json_value);}

void ProtocolWithChangesWriter::WriteRecordToAliasedRecordImpl(evo_test::RecordWithChanges const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "recordToAliasedRecord", json_value);}

void ProtocolWithChangesWriter::WriteRecordToAliasedAliasImpl(evo_test::RecordWithChanges const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "recordToAliasedAlias", json_value);}

void ProtocolWithChangesWriter::WriteOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "optionalRecordWithChanges", json_value);}

void ProtocolWithChangesWriter::WriteAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "aliasedOptionalRecordWithChanges", json_value);}

void ProtocolWithChangesWriter::WriteUnionRecordWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "unionRecordWithChanges", json_value);}

void ProtocolWithChangesWriter::WriteUnionWithSameTypesetImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "unionWithSameTypeset", json_value);}

void ProtocolWithChangesWriter::WriteUnionWithTypesAddedImpl(std::variant<evo_test::RecordWithChanges, float> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "unionWithTypesAdded", json_value);}

void ProtocolWithChangesWriter::WriteUnionWithTypesRemovedImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "unionWithTypesRemoved", json_value);}

void ProtocolWithChangesWriter::WriteRecordToOptionalImpl(evo_test::RecordWithChanges const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "recordToOptional", json_value);}

void ProtocolWithChangesWriter::WriteRecordToAliasedOptionalImpl(evo_test::RecordWithChanges const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "recordToAliasedOptional", json_value);}

void ProtocolWithChangesWriter::WriteRecordToUnionImpl(evo_test::RecordWithChanges const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "recordToUnion", json_value);}

void ProtocolWithChangesWriter::WriteRecordToAliasedUnionImpl(evo_test::RecordWithChanges const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "recordToAliasedUnion", json_value);}

void ProtocolWithChangesWriter::WriteVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "vectorRecordWithChanges", json_value);}

void ProtocolWithChangesWriter::WriteStreamedRecordWithChangesImpl(evo_test::RecordWithChanges const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "streamedRecordWithChanges", json_value);}

void ProtocolWithChangesWriter::WriteAddedStringVectorImpl(std::vector<evo_test::AliasedString> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "addedStringVector", json_value);}

void ProtocolWithChangesWriter::WriteAddedOptionalImpl(std::optional<evo_test::RecordWithChanges> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "addedOptional", json_value);}

void ProtocolWithChangesWriter::WriteAddedMapImpl(std::unordered_map<std::string, std::string> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "addedMap", json_value);}

void ProtocolWithChangesWriter::WriteAddedUnionImpl(std::variant<std::monostate, evo_test::RecordWithChanges, std::string> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "addedUnion", json_value);}

void ProtocolWithChangesWriter::WriteAddedRecordStreamImpl(evo_test::RecordWithChanges const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "addedRecordStream", json_value);}

void ProtocolWithChangesWriter::WriteAddedUnionStreamImpl(std::variant<evo_test::RecordWithChanges, evo_test::RenamedRecord> const& value) {
  ordered_json json_value = value;
  yardl::ndjson::WriteProtocolValue(stream_, "addedUnionStream", json_value);}

void ProtocolWithChangesWriter::Flush() {
  stream_.flush();
}

void ProtocolWithChangesWriter::CloseImpl() {
  stream_.flush();
}

void ProtocolWithChangesReader::ReadInt8ToIntImpl(int8_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "int8ToInt", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadInt8ToLongImpl(int8_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "int8ToLong", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadInt8ToUintImpl(int8_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "int8ToUint", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadInt8ToUlongImpl(int8_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "int8ToUlong", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadInt8ToFloatImpl(int8_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "int8ToFloat", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadInt8ToDoubleImpl(int8_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "int8ToDouble", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadIntToUintImpl(int32_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "intToUint", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadIntToLongImpl(int32_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "intToLong", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadIntToFloatImpl(int32_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "intToFloat", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadIntToDoubleImpl(int32_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "intToDouble", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadUintToUlongImpl(uint32_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "uintToUlong", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadUintToFloatImpl(uint32_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "uintToFloat", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadUintToDoubleImpl(uint32_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "uintToDouble", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadFloatToDoubleImpl(float& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "floatToDouble", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadIntToStringImpl(int32_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "intToString", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadUintToStringImpl(uint32_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "uintToString", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadLongToStringImpl(int64_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "longToString", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadUlongToStringImpl(uint64_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "ulongToString", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadFloatToStringImpl(float& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "floatToString", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadDoubleToStringImpl(double& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "doubleToString", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadIntToOptionalImpl(int32_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "intToOptional", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadFloatToOptionalImpl(float& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "floatToOptional", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadStringToOptionalImpl(std::string& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "stringToOptional", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadIntToUnionImpl(int32_t& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "intToUnion", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadFloatToUnionImpl(float& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "floatToUnion", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadStringToUnionImpl(std::string& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "stringToUnion", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadOptionalIntToFloatImpl(std::optional<int32_t>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "optionalIntToFloat", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadOptionalFloatToStringImpl(std::optional<float>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "optionalFloatToString", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadAliasedLongToStringImpl(evo_test::AliasedLongToString& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "aliasedLongToString", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadStringToAliasedStringImpl(std::string& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "stringToAliasedString", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadStringToAliasedIntImpl(std::string& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "stringToAliasedInt", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadOptionalIntToUnionImpl(std::optional<int32_t>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "optionalIntToUnion", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadOptionalRecordToUnionImpl(std::optional<evo_test::RecordWithChanges>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "optionalRecordToUnion", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadRecordWithChangesImpl(evo_test::RecordWithChanges& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "recordWithChanges", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadAliasedRecordWithChangesImpl(evo_test::AliasedRecordWithChanges& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "aliasedRecordWithChanges", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadRecordToRenamedRecordImpl(evo_test::RenamedRecord& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "recordToRenamedRecord", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadRecordToAliasedRecordImpl(evo_test::RecordWithChanges& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "recordToAliasedRecord", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadRecordToAliasedAliasImpl(evo_test::RecordWithChanges& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "recordToAliasedAlias", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadOptionalRecordWithChangesImpl(std::optional<evo_test::RecordWithChanges>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "optionalRecordWithChanges", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadAliasedOptionalRecordWithChangesImpl(std::optional<evo_test::AliasedRecordWithChanges>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "aliasedOptionalRecordWithChanges", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadUnionRecordWithChangesImpl(std::variant<evo_test::RecordWithChanges, int32_t>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "unionRecordWithChanges", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadUnionWithSameTypesetImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "unionWithSameTypeset", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadUnionWithTypesAddedImpl(std::variant<evo_test::RecordWithChanges, float>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "unionWithTypesAdded", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadUnionWithTypesRemovedImpl(std::variant<evo_test::RecordWithChanges, int32_t, float, std::string>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "unionWithTypesRemoved", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadRecordToOptionalImpl(evo_test::RecordWithChanges& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "recordToOptional", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadRecordToAliasedOptionalImpl(evo_test::RecordWithChanges& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "recordToAliasedOptional", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadRecordToUnionImpl(evo_test::RecordWithChanges& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "recordToUnion", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadRecordToAliasedUnionImpl(evo_test::RecordWithChanges& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "recordToAliasedUnion", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadVectorRecordWithChangesImpl(std::vector<evo_test::RecordWithChanges>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "vectorRecordWithChanges", true, unused_step_, value);
}

bool ProtocolWithChangesReader::ReadStreamedRecordWithChangesImpl(evo_test::RecordWithChanges& value) {
  return yardl::ndjson::ReadProtocolValue(stream_, line_, "streamedRecordWithChanges", false, unused_step_, value);
}

void ProtocolWithChangesReader::ReadAddedStringVectorImpl(std::vector<evo_test::AliasedString>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "addedStringVector", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadAddedOptionalImpl(std::optional<evo_test::RecordWithChanges>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "addedOptional", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadAddedMapImpl(std::unordered_map<std::string, std::string>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "addedMap", true, unused_step_, value);
}

void ProtocolWithChangesReader::ReadAddedUnionImpl(std::variant<std::monostate, evo_test::RecordWithChanges, std::string>& value) {
  yardl::ndjson::ReadProtocolValue(stream_, line_, "addedUnion", true, unused_step_, value);
}

bool ProtocolWithChangesReader::ReadAddedRecordStreamImpl(evo_test::RecordWithChanges& value) {
  return yardl::ndjson::ReadProtocolValue(stream_, line_, "addedRecordStream", false, unused_step_, value);
}

bool ProtocolWithChangesReader::ReadAddedUnionStreamImpl(std::variant<evo_test::RecordWithChanges, evo_test::RenamedRecord>& value) {
  return yardl::ndjson::ReadProtocolValue(stream_, line_, "addedUnionStream", false, unused_step_, value);
}

void ProtocolWithChangesReader::CloseImpl() {
  VerifyFinished();
}

} // namespace evo_test::ndjson

