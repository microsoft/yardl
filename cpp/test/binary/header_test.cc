// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.


#include "../generated/yardl/detail/binary/header.h"

#include <iostream>
#include <sstream>

#include <gtest/gtest.h>

#include "../generated/binary/protocols.h"
#include "../generated/yardl/detail/binary/serializers.h"

using namespace yardl::binary;
using namespace test_model;
using namespace test_model::binary;

namespace {

TEST(HeaderTests, ValidHeader) {
  std::stringstream ss;
  {
    ScalarsWriter w(ss);
    w.WriteInt32(1);
    w.WriteRecord({1, 2, 3});
  }

  ScalarsReader r(ss);
  int i;
  r.ReadInt32(i);
  RecordWithPrimitives rec;
  r.ReadRecord(rec);
}

TEST(HeaderTests, WrongMagicNumber) {
  std::stringstream ss;
  {
    ScalarsWriter w(ss);
    w.WriteInt32(1);
    w.WriteRecord({1, 2, 3});
  }

  ss.seekp(0, std::ios::beg);
  ss.write("WRONG", 5);

  ASSERT_ANY_THROW(ScalarsReader r(ss));
}

TEST(HeaderTests, WrongVersion) {
  std::stringstream ss;
  {
    ScalarsWriter w(ss);
    w.WriteInt32(1);
    w.WriteRecord({1, 2, 3});
  }

  ss.seekp(MAGIC_BYTES.size(), std::ios::beg);
  ss.put(9);

  ASSERT_ANY_THROW(ScalarsReader r(ss));
}

TEST(HeaderTests, WrongSchema) {
  std::stringstream ss;
  CodedOutputStream w(ss);
  std::string bogus_schema = "{}";
  WriteHeader(w, bogus_schema);

  ASSERT_ANY_THROW(ScalarsReader r(ss));
}

}  // namespace
