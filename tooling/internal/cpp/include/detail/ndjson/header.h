#pragma once

#include <iostream>
#include <string>

#include <nlohmann/json.hpp>

namespace yardl::ndjson {
using ordered_json = nlohmann::ordered_json;

static inline uint32_t kNDJsonFormatVersionNumber = 1;

inline void WriteHeader(std::ostream& stream, std::string& schema) {
  auto parsed_schema = ordered_json::parse(schema);

  ordered_json metadata = {{"yardl", {{"version", kNDJsonFormatVersionNumber}, {"schema", parsed_schema}}}};
  stream << metadata << "\n";
}

inline ordered_json ReadHeader(std::istream& stream) {
  std::string line;
  std::getline(stream, line);
  try {
    ordered_json actual_header_json = ordered_json::parse(line);
    actual_header_json = actual_header_json.at("yardl");
    if (actual_header_json["version"] != kNDJsonFormatVersionNumber) {
      throw std::runtime_error(
          "Unsupported Yardl NDJSON format version.");
    }

    return actual_header_json.at("schema");
  } catch (ordered_json::exception const&) {
    throw std::runtime_error(
        "Data in the stream is not in the expected Yardl NDJSON format.");
  }
}

inline void ReadAndValidateHeader(std::istream& stream, std::string& expected_schema) {
  ordered_json expected_schema_json = ordered_json::parse(expected_schema);
  ordered_json actual_schema_json = ReadHeader(stream);
  if (expected_schema_json != actual_schema_json) {
    throw std::runtime_error(
        "The schema of the data to be read is not compatible with the current protocol.");
  }
}
}  // namespace yardl::ndjson
