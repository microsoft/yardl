// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <chrono>
#include <iostream>
#include <string>

#include <nlohmann/json.hpp>

#include "generated/yardl/detail/binary/reader_writer.h"
#include "generated/yardl/detail/hdf5/io.h"
#include "generated/yardl/detail/ndjson/reader_writer.h"

namespace yardl::testing {
using ordered_json = nlohmann::ordered_json;

std::string const kOutputFileName = "/tmp/benchmark_data.dat";

struct Result {
  double write_mi_bytes_per_second;
  double read_mi_bytes_per_second;
};

void to_json(ordered_json& j, Result const& res) {
  j = ordered_json{{"write_mi_bytes_per_second", res.write_mi_bytes_per_second},
                   {"read_mi_bytes_per_second", res.read_mi_bytes_per_second}};
}

inline void AssertRepetitionsSame(size_t expected, size_t actual) {
  if (expected != actual) {
    throw std::runtime_error("Expected " + std::to_string(expected) +
                             " repetitions, but got " + std::to_string(actual));
  }
}

inline size_t ScaleRepetitions(size_t repetitions, double scale = 1) {
  return static_cast<size_t>(repetitions * scale);
}

template <typename WriteFunc, typename ReadFunc>
Result TimeScenario(size_t total_bytes_size, WriteFunc writeImpl, ReadFunc readImpl) {
  std::remove(kOutputFileName.c_str());

  double total_size_mi_byte = total_bytes_size / 1024.0 / 1024.0;

  auto write_start = std::chrono::high_resolution_clock::now();
  writeImpl();
  auto write_end = std::chrono::high_resolution_clock::now();
  double write_elapsed_seconds = std::chrono::duration_cast<std::chrono::duration<double>>(write_end - write_start).count();
  double write_mi_bytes_per_second = total_size_mi_byte / write_elapsed_seconds;

  auto read_start = std::chrono::high_resolution_clock::now();
  readImpl();
  auto read_end = std::chrono::high_resolution_clock::now();
  double read_elapsed_seconds = std::chrono::duration_cast<std::chrono::duration<double>>(read_end - read_start).count();
  double read_mi_bytes_per_second = total_size_mi_byte / read_elapsed_seconds;

  return Result{write_mi_bytes_per_second, read_mi_bytes_per_second};
}

}  // namespace yardl::testing
