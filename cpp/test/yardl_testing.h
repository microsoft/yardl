// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once
#include <filesystem>
#include <memory>

#include <fmt/core.h>

#include "factories.h"
#include "format.h"

namespace yardl::testing {

template <typename T>
std::unique_ptr<T> CreateValidatingWriter(Format format, std::string const& filename);

static inline std::string TestFilename(Format format, bool ensure_deleted = true) {
  std::filesystem::path test_output_dir = "test_output/";
  std::string extension;
  switch (format) {
    case Format::kBinary:
      test_output_dir /= "binary";
      extension = ".bin";
      break;
    case Format::kHdf5:
      test_output_dir /= "hdf5";
      extension = ".h5";
      break;
    case Format::kNDJson:
      test_output_dir /= "ndjson";
      extension = ".ndjson";
      break;
  }

  std::filesystem::create_directories(test_output_dir);
  auto test_info = ::testing::UnitTest::GetInstance()->current_test_info();
  std::string test_name = test_info->name();
  std::replace(test_name.begin(), test_name.end(), '/', '_');
  auto filename = fmt::format(
      "{}/{}_{}{}",
      test_output_dir.string(), test_info->test_suite_name(), test_name, extension);
  if (ensure_deleted) {
    std::remove(filename.c_str());
  }
  return filename;
}
}  // namespace yardl::testing
