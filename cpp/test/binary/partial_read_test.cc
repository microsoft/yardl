// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <iostream>
#include <sstream>

#include <gtest/gtest.h>

#include "../generated/binary/protocols.h"
#include "../generated/yardl/detail/binary/header.h"
#include "../generated/yardl/detail/binary/serializers.h"

using namespace yardl::binary;
using namespace test_model;
using namespace test_model::binary;

namespace {

std::stringstream GenerateStream() {
  std::stringstream ss;
  {
    StateTestWriter w(ss);
    w.WriteAnInt(42);
    std::vector<int> stream{1, 2, 3, 4, 5};
    w.WriteAStream(stream);
    w.EndAStream();
    w.WriteAnotherInt(153);
  }
  ss.seekg(0);
  return ss;
}

TEST(PartialReadTests, SkipSteps) {
  {
    std::stringstream ss = GenerateStream();
    StateTestReader r(ss, true);
    int value;
    r.ReadAnInt(value);

    ASSERT_NO_THROW(r.Close());
  }

  {
    std::stringstream ss = GenerateStream();
    StateTestReader r(ss, true);
    int value;
    r.ReadAnInt(value);
    while (r.ReadAStream(value)) {
      // pass
    }

    ASSERT_NO_THROW(r.Close());
  }
}

TEST(PartialReadTests, SkipStreamItems) {
  std::stringstream ss = GenerateStream();
  {
    StateTestReader r(ss, true);
    int value;
    r.ReadAnInt(value);
    GTEST_ASSERT_TRUE(r.ReadAStream(value));
    GTEST_ASSERT_TRUE(r.ReadAStream(value));

    ASSERT_NO_THROW(r.Close());
  }
}

}  // namespace
