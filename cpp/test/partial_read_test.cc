// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <iostream>
#include <sstream>

#include <gtest/gtest.h>

#include "generated/binary/protocols.h"
#include "generated/hdf5/protocols.h"
#include "generated/ndjson/protocols.h"
#include "yardl_testing.h"

using namespace test_model;

using namespace yardl::testing;

namespace {

std::unique_ptr<test_model::StateTestReaderBase> CreatePartialReader(Format format, std::string const& filename) {
  switch (format) {
    case Format::kHdf5:
      return std::make_unique<test_model::hdf5::StateTestReader>(filename, true);
    case Format::kBinary:
      return std::make_unique<test_model::binary::StateTestReader>(filename, true);
    case Format::kNDJson:
      return std::make_unique<test_model::ndjson::StateTestReader>(filename, true);
    default:
      throw std::runtime_error("Unknown format");
  }
}

class PartialReadTests : public ::testing::TestWithParam<Format> {
 protected:
  void SetUp() override {
    format_ = GetParam();
  }

  std::unique_ptr<StateTestReaderBase> CreateTestingReader() {
    // Use CreateWriter from generated/factories.h
    auto filename = TestFilename(format_);
    auto writer = yardl::testing::CreateWriter<StateTestWriterBase>(format_, filename);
    writer->WriteAnInt(42);
    std::vector<int32_t> stream{1, 2, 3, 4, 5};
    writer->WriteAStream(stream);
    writer->EndAStream();
    writer->WriteAnotherInt(153);
    writer->Close();

    // Use specialized factory that sets `skip_completed_check=true` on the reader
    return CreatePartialReader(format_, filename);
  }

  Format format_;
};

TEST_P(PartialReadTests, SkipSteps) {
  {
    auto reader = CreateTestingReader();
    int value;
    reader->ReadAnInt(value);
    ASSERT_NO_THROW(reader->Close());
  }

  {
    auto reader = CreateTestingReader();
    int value;
    reader->ReadAnInt(value);
    while (reader->ReadAStream(value)) {
      // pass
    }
    ASSERT_NO_THROW(reader->Close());
  }
}

TEST_P(PartialReadTests, SkipStreamItems) {
  auto reader = CreateTestingReader();
  int value;
  reader->ReadAnInt(value);
  GTEST_ASSERT_TRUE(reader->ReadAStream(value));
  GTEST_ASSERT_TRUE(reader->ReadAStream(value));
  ASSERT_NO_THROW(reader->Close());
}

INSTANTIATE_TEST_SUITE_P(,
                         PartialReadTests,
                         ::testing::Values(
                             Format::kBinary,
                             Format::kHdf5,
                             Format::kNDJson),
                         [](::testing::TestParamInfo<Format> const& info) {
  switch (info.param) {
  case Format::kBinary:
    return "Binary";
  case Format::kHdf5:
    return "HDF5";
  case Format::kNDJson:
    return "NDJson";
  default:
    return "Unknown";
  } });

}  // namespace
