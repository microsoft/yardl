// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <chrono>
#include <iomanip>
#include <iostream>
#include <string>

#include "generated/yardl/detail/binary/reader_writer.h"
#include "generated/yardl/detail/hdf5/io.h"

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

inline void WriteBenchmarkTableHeader() {
  std::cout.imbue(std::locale(""));

  std::cout << std::left
            << "| " << std::setw(8) << "Provider"
            << " | " << std::setw(25) << "Scenario"
            << " | " << std::setw(6) << "Action"
            << " | " << std::setw(10) << "MiB/s"
            << " | " << std::setw(6) << "GiB"
            << " |" << std::endl;
  std::cout << std::right << std::setfill('-')
            << "| " << std::setw(8) << ""
            << " | " << std::setw(25) << ""
            << " | " << std::setw(6) << ""
            << " | " << std::setw(10) << ""
            << " | " << std::setw(6) << ""
            << " | " << std::setfill(' ') << std::left << std::endl;
}

inline void WriteBenchmarkTableRow(std::string& provider, std::string& scenario,
                                   std::string& action, double throughput_mi_byte_s,
                                   double total_size_gi_byte) {
  std::string provider_color = provider == "binary"
                                   ? (action == "write" ? BOLDCYAN : CYAN)
                                   : (action == "write" ? BOLDBLUE : BLUE);

  std::cout << std::left
            << "| " << provider_color << std::setw(8) << provider << RESET
            << " | " << provider_color << std::setw(25) << scenario << RESET
            << " | " << provider_color << std::setw(6) << action << RESET
            << " | " << provider_color << std::setw(10) << std::right << std::fixed << std::setprecision(2) << throughput_mi_byte_s << RESET
            << " | " << provider_color << std::setw(6) << std::right << std::fixed << std::setprecision(2) << total_size_gi_byte << RESET
            << " |" << std::endl;
}

template <typename T>
class TimedScope {
 public:
  TimedScope(std::string scenario_name, size_t total_size_bytes)
      : total_size_bytes_(total_size_bytes),
        scenario_name_(std::move(scenario_name)),
        start_(std::chrono::steady_clock::now()) {
  }

  ~TimedScope() {
    auto end = std::chrono::steady_clock::now();
    std::string action;
    std::string provider;

    if constexpr (std::is_base_of_v<yardl::binary::BinaryWriter, T>) {
      action = "write";
      provider = "binary";
    } else if constexpr (std::is_base_of_v<yardl::binary::BinaryReader, T>) {
      action = "read";
      provider = "binary";
    } else if constexpr (std::is_base_of_v<yardl::hdf5::Hdf5Writer, T>) {
      action = "write";
      provider = "hdf5";
    } else if constexpr (std::is_base_of_v<yardl::hdf5::Hdf5Reader, T>) {
      action = "read";
      provider = "hdf5";
    } else {
      throw std::runtime_error("Unknown type");
    }

    if (scenario_name_.find("Benchmark") == 0) {
      scenario_name_ = scenario_name_.substr(9);
    }

    float elapsed_seconds = std::chrono::duration_cast<std::chrono::duration<float>>(end - start_).count();
    float total_size_mi_byte = total_size_bytes_ / 1024.0 / 1024.0;
    float throughput_mi_byte_s = total_size_mi_byte / elapsed_seconds;

    WriteBenchmarkTableRow(provider, scenario_name_, action, throughput_mi_byte_s, total_size_mi_byte / 1024.0);
  }

 private:
  size_t total_size_bytes_;
  std::string scenario_name_;
  std::chrono::steady_clock::time_point start_;
};

}  // namespace yardl::testing
