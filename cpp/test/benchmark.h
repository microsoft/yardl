// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <chrono>
#include <iomanip>
#include <iostream>
#include <string>

#include "generated/yardl/detail/binary/reader_writer.h"
#include "generated/yardl/detail/hdf5/io.h"
#include "generated/yardl/detail/ndjson/reader_writer.h"

namespace yardl::testing {

#define RESET "\033[0m"
#define BLACK "\033[30m"
#define RED "\033[31m"
#define GREEN "\033[32m"
#define YELLOW "\033[33m"
#define BLUE "\033[34m"
#define MAGENTA "\033[35m"
#define CYAN "\033[36m"
#define WHITE "\033[37m"
#define BOLDBLACK "\033[1m\033[30m"
#define BOLDRED "\033[1m\033[31m"
#define BOLDGREEN "\033[1m\033[32m"
#define BOLDYELLOW "\033[1m\033[33m"
#define BOLDBLUE "\033[1m\033[34m"
#define BOLDMAGENTA "\033[1m\033[35m"
#define BOLDCYAN "\033[1m\033[36m"
#define BOLDWHITE "\033[1m\033[37m"

inline void AssertRepetitionsSame(size_t expected, size_t actual) {
  if (expected != actual) {
    throw std::runtime_error("Expected " + std::to_string(expected) +
                             " repetitions, but got " + std::to_string(actual));
  }
}

inline size_t ScaleRepetitions(size_t repetitions, double scale = 1) {
  return static_cast<size_t>(repetitions * scale);
}

inline void WriteSeparatorRow() {
  std::cout << std::right << std::setfill('-')
            << "| " << std::setw(25) << ""
            << " | " << std::setw(8) << ""
            << " | " << std::setw(12) << ""
            << " | " << std::setw(12) << ""
            << " | " << std::setfill(' ') << std::left << std::endl;
}

inline void WriteBenchmarkTableHeader() {
  std::cout.imbue(std::locale(""));

  std::cout << std::left
            << "| " << std::setw(25) << "Scenario"
            << " | " << std::setw(8) << "Provider"
            << " | " << std::setw(12) << "Write MiB/s"
            << " | " << std::setw(12) << "Read MiB/s"
            << " |" << std::endl;

  WriteSeparatorRow();
}

inline void WriteBenchmarkTableRow(std::string scenario, std::string provider, double write_throughput_mi_byte_s, double read_throughput_mi_byte_s) {
  std::string provider_color;
  if (provider == "binary") {
    provider_color = CYAN;
  } else if (provider == "hdf5") {
    provider_color = BLUE;
  } else if (provider == "ndjson") {
    provider_color = GREEN;
  } else {
    provider_color = WHITE;
  }

  std::cout << std::left
            << "| " << provider_color << std::setw(25) << scenario << RESET
            << " | " << provider_color << std::setw(8) << provider << RESET
            << " | " << provider_color << std::setw(12) << std::right << std::fixed << std::setprecision(2) << write_throughput_mi_byte_s << RESET
            << " | " << provider_color << std::setw(12) << std::right << std::fixed << std::setprecision(2) << read_throughput_mi_byte_s << RESET
            << " |" << std::endl;
}

template <typename TWriter, typename WriteFunc, typename ReadFunc>
void TimeScenario(std::string scenario_name, size_t total_bytes_size, WriteFunc writeImpl, ReadFunc readImpl) {
  double total_size_mi_byte = total_bytes_size / 1024.0 / 1024.0;

  auto write_start = std::chrono::high_resolution_clock::now();
  writeImpl();
  auto write_end = std::chrono::high_resolution_clock::now();
  double write_elapsed_seconds = std::chrono::duration_cast<std::chrono::duration<double>>(write_end - write_start).count();
  double write_throughput_mi_byte_s = total_size_mi_byte / write_elapsed_seconds;

  auto read_start = std::chrono::high_resolution_clock::now();
  readImpl();
  auto read_end = std::chrono::high_resolution_clock::now();
  double read_elapsed_seconds = std::chrono::duration_cast<std::chrono::duration<double>>(read_end - read_start).count();
  double read_throughput_mi_byte_s = total_size_mi_byte / read_elapsed_seconds;

  std::string provider;
  if constexpr (std::is_base_of_v<yardl::binary::BinaryWriter, TWriter>) {
    provider = "binary";
  } else if constexpr (std::is_base_of_v<yardl::hdf5::Hdf5Writer, TWriter>) {
    provider = "hdf5";
  } else if constexpr (std::is_base_of_v<yardl::ndjson::NDJsonWriter, TWriter>) {
    provider = "ndjson";
  } else {
    throw std::runtime_error("Unknown writer type");
  }

  if (scenario_name.find("Benchmark") == 0) {
    scenario_name = scenario_name.substr(9);
  }

  WriteBenchmarkTableRow(scenario_name, provider, write_throughput_mi_byte_s, read_throughput_mi_byte_s);
}

}  // namespace yardl::testing
