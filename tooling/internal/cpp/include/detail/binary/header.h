// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include "serializers.h"

namespace yardl::binary {
static inline std::array<char, 5> MAGIC_BYTES = {'y', 'a', 'r', 'd', 'l'};
static inline uint32_t kBinaryFormatVersionNumber = 1;

inline void WriteHeader(CodedOutputStream& w, std::string const& schema) {
  w.WriteBytes(MAGIC_BYTES.data(), MAGIC_BYTES.size());
  w.WriteFixedInteger(kBinaryFormatVersionNumber);
  yardl::binary::WriteString(w, schema);
}

inline std::string ReadHeader(CodedInputStream& r) {
  std::array<char, 5> magic_bytes{};
  r.ReadBytes(magic_bytes.data(), magic_bytes.size());
  if (magic_bytes != MAGIC_BYTES) {
    throw std::runtime_error("Data in the stream is not in the expected format.");
  }

  uint32_t version_number;
  r.ReadFixedInteger(version_number);
  if (version_number != kBinaryFormatVersionNumber) {
    throw std::runtime_error(
        "Data in the stream is not in the expected format. Unsupported version.");
  }

  std::string actual_schema;
  yardl::binary::ReadString(r, actual_schema);
  return actual_schema;
}

inline void MatchSchema(std::string actual_schema, std::string const& current_schema, std::vector<std::string> const& previous_schemas, int& schema_index) {
  if (actual_schema == current_schema) {
    schema_index = -1;
    return;
  }

  for (size_t i = 0; i < previous_schemas.size(); ++i) {
    if (actual_schema == previous_schemas[i]) {
      schema_index = i;
      return;
    }
  }

  throw std::runtime_error(
      "The schema of the data to be read is not compatible with the current protocol.");
}

}  // namespace yardl::binary
