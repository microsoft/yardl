// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include "serializers.h"

namespace yardl::binary {
static inline std::array<char, 5> MAGIC_BYTES = {'y', 'a', 'r', 'd', 'l'};
static inline uint32_t kBinaryFormatVersionNumber = 1;

inline void WriteHeader(CodedOutputStream& w, std::string& schema) {
  w.WriteBytes(MAGIC_BYTES.data(), MAGIC_BYTES.size());
  w.WriteFixedInteger(kBinaryFormatVersionNumber);
  yardl::binary::WriteString(w, schema);
}

inline void ReadHeader(CodedInputStream& r, std::string& schema) {
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
  if (actual_schema != schema) {
    throw std::runtime_error(
        "The schema of the data to be read is not compatible with the current protocol.");
  }
}

}  // namespace yardl::binary
