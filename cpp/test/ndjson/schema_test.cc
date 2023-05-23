// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <iostream>
#include <sstream>

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "../generated/yardl/detail/ndjson/reader_writer.h"

using namespace testing;
using namespace yardl::ndjson;

namespace {

class MyReader : public NDJsonReader {
 public:
  MyReader(std::istream& stream, std::string& schema) : NDJsonReader(stream, schema) {}
};

TEST(NDJsonSchemaTests, ValidSchema) {
  std::string expected_schema = R"({"protocol":{"name":"HelloWorld","sequence":[]},"types":null})";
  std::stringstream ss{R"({"yardl":{"version":1,"schema":{"protocol":{"name":"HelloWorld","sequence":[]},"types":null}}})"};
  MyReader reader(ss, expected_schema);
}

TEST(NDJsonSchemaTests, NotJson) {
  std::string expected_schema = R"({"protocol":{"name":"HelloWorld","sequence":[]},"types":null})";

  EXPECT_THAT(
      [&]() {
        std::stringstream ss{"abc"};
        MyReader reader(ss, expected_schema);
      },
      ThrowsMessage<std::runtime_error>(HasSubstr("Data in the stream is not in the expected Yardl NDJSON format.")));
}

TEST(NDJsonSchemaTests, MissingYardl) {
  std::string expected_schema = R"({"protocol":{"name":"HelloWorld","sequence":[]},"types":null})";
  EXPECT_THAT(
      [&]() {
        std::stringstream ss{R"({"foo": "bar"})"};
        MyReader reader(ss, expected_schema);
      },
      ThrowsMessage<std::runtime_error>(HasSubstr("Data in the stream is not in the expected Yardl NDJSON format.")));
}

TEST(NDJsonSchemaTests, Array) {
  std::string expected_schema = R"([{"protocol":{"name":"HelloWorld","sequence":[]},"types":null}])";
  EXPECT_THAT(
      [&]() {
        std::stringstream ss{R"([1,2,3])"};
        MyReader reader(ss, expected_schema);
      },
      ThrowsMessage<std::runtime_error>(HasSubstr("Data in the stream is not in the expected Yardl NDJSON format.")));
}

TEST(NDJsonSchemaTests, DifferentYardlVersion) {
  std::string expected_schema = R"([{"protocol":{"name":"HelloWorld","sequence":[]},"types":null}])";
  EXPECT_THAT(
      [&]() {
        std::stringstream ss{R"({"yardl":{"version":9999,"schema":{"protocol":{"name":"HelloWorld","sequence":[]},"types":null}}})"};
        MyReader reader(ss, expected_schema);
      },
      ThrowsMessage<std::runtime_error>(HasSubstr("Unsupported Yardl NDJSON format version.")));
}

TEST(NDJsonSchemaTests, WrongProtocol) {
  std::string expected_schema = R"([{"protocol":{"name":"HelloWorld","sequence":[]},"types":null}])";
  EXPECT_THAT(
      [&]() {
        std::stringstream ss{R"({"yardl":{"version":1,"schema":{"protocol":{"name":"WrongProtocol","sequence":[]},"types":null}}})"};
        MyReader reader(ss, expected_schema);
      },
      ThrowsMessage<std::runtime_error>(HasSubstr("The schema of the data to be read is not compatible with the current protocol.")));
}

}  // namespace
