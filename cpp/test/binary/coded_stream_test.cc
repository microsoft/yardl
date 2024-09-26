// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <fstream>
#include <iostream>
#include <numeric>

#include <gtest/gtest.h>

#include "../generated/yardl/detail/binary/serializers.h"

using namespace yardl::binary;

namespace {

[[maybe_unused]] void WriteStringStreamToDebugFile(std::string const str) {
  std::ofstream file("debug.bin", std::ios::binary);
  file.write(str.data(), str.size());
  file.close();
}

TEST(CodedStreamTests, ReadTooFar) {
  std::stringstream ss;
  {
    CodedOutputStream w(ss, 10);
    WriteInteger(w, 1);
  }

  CodedInputStream r(ss, 10);
  int i;
  ReadInteger(r, i);
  ASSERT_THROW(ReadInteger(r, i), EndOfStreamException);
}

TEST(CodedStreamTests, IncompleteRead) {
  std::stringstream ss;
  {
    CodedOutputStream w(ss, 10);
    WriteInteger(w, 1);
    WriteInteger(w, 1);
  }

  CodedInputStream r(ss, 10);
  int i;
  ReadInteger(r, i);
  ASSERT_ANY_THROW(r.VerifyFinished());
}

TEST(CodedStreamTests, ReadExactBufferLength) {
  std::stringstream ss;
  std::string data(256, 'a');
  {
    CodedOutputStream w(ss);
    WriteString(w, data);
  }

  ss.seekg(0, std::ios::end);
  auto length = ss.tellg();
  ss.seekg(0, std::ios::beg);

  CodedInputStream r(ss, length);
  std::string data2;
  ReadString(r, data2);
  r.VerifyFinished();
}

TEST(CodedStreamTests, ReadOneMoreThanBufferLength) {
  std::stringstream ss;
  std::string data(256, 'a');
  {
    CodedOutputStream w(ss);
    WriteString(w, data);
  }

  ss.seekg(0, std::ios::end);
  auto length = ss.tellg();
  ss.seekg(0, std::ios::beg);

  CodedInputStream r(ss, static_cast<size_t>(length) - 1);
  std::string data2;
  ReadString(r, data2);
  r.VerifyFinished();
}

TEST(CodedStreamTests, ReadOneLessThanBufferLength) {
  std::stringstream ss;
  std::string data(256, 'a');
  {
    CodedOutputStream w(ss);
    WriteString(w, data);
  }

  ss.seekg(0, std::ios::end);
  auto length = ss.tellg();
  ss.seekg(0, std::ios::beg);

  CodedInputStream r(ss, static_cast<size_t>(length) + 1);
  std::string data2;
  ReadString(r, data2);
  r.VerifyFinished();
}

TEST(CodedStreamTests, BadStreamThrows) {
  std::stringstream ss;
  CodedOutputStream w(ss);
  WriteInteger(w, 1);
  ss.setstate(std::ios::badbit);
  EXPECT_ANY_THROW(w.Flush());
}

TEST(CodedStreamTests, ScalarByte) {
  std::stringstream ss;
  {
    CodedOutputStream w(ss, 10);
    for (uint8_t i = 0; i < 15; i++) {
      WriteInteger(w, i);
      WriteInteger(w, static_cast<int8_t>(i));
    }
  }

  CodedInputStream r(ss, 10);
  for (uint8_t i = 0; i < 15; i++) {
    uint8_t val1;
    ReadInteger(r, val1);
    ASSERT_EQ(i, val1);
    int8_t val2;
    ReadInteger(r, val2);
    ASSERT_EQ(i, val2);
  }

  r.VerifyFinished();
}

TEST(CodedStreamTests, FixedInteger) {
  std::stringstream ss;
  {
    CodedOutputStream w(ss, 10);
    for (int i = 0; i < 15; i++) {
      w.WriteFixedInteger(i);
      w.WriteFixedInteger(static_cast<uint32_t>(i));
    }
  }

  CodedInputStream r(ss, 10);
  for (int i = 0; i < 15; i++) {
    int val1;
    r.ReadFixedInteger(val1);
    ASSERT_EQ(i, val1);

    uint32_t val2;
    r.ReadFixedInteger(val2);
    ASSERT_EQ(i, val2);
  }

  r.VerifyFinished();
}

TEST(CodedStreamTests, VarShort) {
  std::vector<int16_t> entries = {
      0, 1, 5, 33, 0x7E, 0x7F, 0x80, 0x81, 255,
      256, 257, 838, 0x3FFF, 0x4000, 0x4001, 0x7FFF};
  std::stringstream ss;
  {
    CodedOutputStream w(ss, 10);
    for (auto e : entries) {
      WriteInteger(w, e);
      WriteInteger(w, -e);
    }
  }

  CodedInputStream r(ss, 10);
  for (auto e : entries) {
    int16_t val1;
    ReadInteger(r, val1);
    ASSERT_EQ(e, val1);

    int16_t val2;
    ReadInteger(r, val2);
    ASSERT_EQ(-e, val2);
  }
}

TEST(CodedStreamTests, VarUShort) {
  std::vector<uint16_t> entries = {
      0, 1, 5, 33, 0x7E, 0x7F, 0x80, 0x81, 255, 256, 257,
      838, 0x3FFF, 0x4000, 0x4001, 0x7FFF, 0x8000, 0x8001, 0xFFFF};
  std::stringstream ss;
  {
    CodedOutputStream w(ss, 10);
    for (auto e : entries) {
      WriteInteger(w, e);
    }
  }

  CodedInputStream r(ss, 10);
  for (auto e : entries) {
    uint16_t val1;
    ReadInteger(r, val1);
    ASSERT_EQ(e, val1);
  }
}

TEST(CodedStreamTests, VarIntegers) {
  std::vector<uint32_t> entries = {
      0, 1, 5, 33, 0x7E, 0x7F, 0x80, 0x81, 255, 256, 257,
      838, 283928, 2847772, 3443, 0x7FFFFFFF, 0xFFFFFFFF};

  std::stringstream ss;
  {
    CodedOutputStream w(ss);
    for (auto e : entries) {
      WriteInteger(w, e);
      WriteInteger(w, static_cast<int>(e));
      WriteInteger(w, -static_cast<int>(e));

      WriteInteger(w, static_cast<uint64_t>(e));
      WriteInteger(w, static_cast<uint64_t>(e) | 0x800000000UL);

      WriteInteger(w, static_cast<int64_t>(e));
      WriteInteger(w, -static_cast<int64_t>(e));

      WriteInteger(w, static_cast<int64_t>(e) | 0x400000000L);
      WriteInteger(w, -(static_cast<int64_t>(e) | 0x400000000L));
    }
  }

  CodedInputStream r(ss);
  for (auto e : entries) {
    uint32_t val1;
    ReadInteger(r, val1);
    ASSERT_EQ(e, val1);

    int32_t val2;
    ReadInteger(r, val2);
    ASSERT_EQ(static_cast<int>(e), val2);

    int32_t val3;
    ReadInteger(r, val3);
    ASSERT_EQ(-static_cast<int>(e), val3);

    uint64_t val4;
    ReadInteger(r, val4);
    ASSERT_EQ(static_cast<uint64_t>(e), val4);

    uint64_t val5;
    ReadInteger(r, val5);
    ASSERT_EQ(static_cast<uint64_t>(e) | 0x800000000UL, val5);

    int64_t val6;
    ReadInteger(r, val6);
    ASSERT_EQ(static_cast<int64_t>(e), val6);

    int64_t val7;
    ReadInteger(r, val7);
    ASSERT_EQ(-static_cast<int64_t>(e), val7);

    int64_t val8;
    ReadInteger(r, val8);
    ASSERT_EQ(static_cast<int64_t>(e) | 0x400000000L, val8);

    int64_t val9;
    ReadInteger(r, val9);
    ASSERT_EQ(-(static_cast<int64_t>(e) | 0x400000000L), val9);
  }

  r.VerifyFinished();
}

TEST(CodedStreamTests, FloatingPoint) {
  std::stringstream ss;
  {
    CodedOutputStream w(ss, 10);
    for (int i = 0; i < 15; i++) {
      WriteFloatingPoint(w, static_cast<float>(i) - std::numeric_limits<float>::epsilon());
      WriteFloatingPoint(w, static_cast<double>(i) - std::numeric_limits<double>::epsilon());
    }
  }

  CodedInputStream r(ss, 10);
  for (int i = 0; i < 15; i++) {
    float val1;
    ReadFloatingPoint(r, val1);
    ASSERT_EQ(static_cast<float>(i) - std::numeric_limits<float>::epsilon(), val1);

    double val2;
    ReadFloatingPoint(r, val2);
    ASSERT_EQ(static_cast<double>(i) - std::numeric_limits<double>::epsilon(), val2);
  }

  r.VerifyFinished();
}

TEST(CodedStreamTests, Strings) {
  std::stringstream ss;
  std::string long_string(20000, 'a');
  {
    CodedOutputStream w(ss, 10);
    WriteString(w, "hello");
    WriteString(w, long_string);
    WriteString(w, "world");
  }

  CodedInputStream r(ss, 10);
  std::string value;
  ReadString(r, value);
  ASSERT_EQ("hello", value);

  ReadString(r, value);
  ASSERT_EQ(long_string, value);

  ReadString(r, value);
  ASSERT_EQ("world", value);

  r.VerifyFinished();
}

TEST(CodedStreamTests, Vector) {
  std::stringstream ss;
  std::vector<int> v1 = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
  std::vector<std::string> v2 = {"hello", "world", "this", "is", "a", "test"};
  std::vector<float> v3 = {1, 2, 3, 4, 5, 6, 7, 8, 9};
  {
    CodedOutputStream w(ss, 10);
    WriteVector<int, &WriteInteger>(w, v1);
    WriteVector<std::string, &WriteString>(w, v2);
    WriteVector<float, &WriteFloatingPoint>(w, v3);
  }

  CodedInputStream r(ss, 10);
  std::vector<int> v1_read;
  ReadVector<int, &ReadInteger>(r, v1_read);
  ASSERT_EQ(v1, v1_read);

  std::vector<std::string> v2_read;
  ReadVector<std::string, &ReadString>(r, v2_read);
  ASSERT_EQ(v2, v2_read);

  std::vector<float> v3_read;
  ReadVector<float, &ReadFloatingPoint>(r, v3_read);
  ASSERT_EQ(v3, v3_read);

  r.VerifyFinished();
}

TEST(CodedStreamTests, Array) {
  std::stringstream ss;
  std::array<int, 10> a1 = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
  std::array<std::string, 6> a2 = {"hello", "world", "this", "is", "a", "test"};
  std::array<float, 10> a3 = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
  {
    CodedOutputStream w(ss, 10);
    WriteArray<int, &WriteInteger, 10>(w, a1);
    WriteArray<std::string, &WriteString, 6>(w, a2);
    WriteArray<float, &WriteFloatingPoint, 10>(w, a3);
  }

  CodedInputStream r(ss, 10);
  std::array<int, 10> v1_read;
  ReadArray<int, &ReadInteger, 10>(r, v1_read);
  ASSERT_EQ(a1, v1_read);

  std::array<std::string, 6> v2_read;
  ReadArray<std::string, &ReadString, 6>(r, v2_read);
  ASSERT_EQ(a2, v2_read);

  std::array<float, 10> v3_read;
  ReadArray<float, &ReadFloatingPoint, 10>(r, v3_read);
  ASSERT_EQ(a3, v3_read);

  r.VerifyFinished();
}

TEST(CodedStreamTests, DynamicNDArray) {
  std::stringstream ss;
  yardl::DynamicNDArray<int> a1({{1, 2, 3}, {4, 5, 6}});
  yardl::DynamicNDArray<float> a2({{1, 2, 3}, {4, 5, 6}});
  {
    CodedOutputStream w(ss, 10);
    WriteDynamicNDArray<int, &WriteInteger>(w, a1);
    WriteDynamicNDArray<float, &WriteFloatingPoint>(w, a2);
  }

  CodedInputStream r(ss, 10);
  yardl::DynamicNDArray<int> a1_read;
  ReadDynamicNDArray<int, &ReadInteger>(r, a1_read);
  ASSERT_EQ(a1, a1_read);

  yardl::DynamicNDArray<float> a2_read;
  ReadDynamicNDArray<float, &ReadFloatingPoint>(r, a2_read);
  ASSERT_EQ(a2, a2_read);

  r.VerifyFinished();
}

TEST(CodedStreamTests, NDArray) {
  std::stringstream ss;
  yardl::NDArray<int, 2> a1({{1, 2, 3}, {4, 5, 6}});
  yardl::NDArray<float, 2> a2({{1, 2, 3}, {4, 5, 6}});
  {
    CodedOutputStream w(ss, 10);
    WriteNDArray<int, &WriteInteger, 2>(w, a1);
    WriteNDArray<float, &WriteFloatingPoint, 2>(w, a2);
  }

  CodedInputStream r(ss, 10);
  yardl::NDArray<int, 2> a1_read;
  ReadNDArray<int, &ReadInteger, 2>(r, a1_read);
  ASSERT_EQ(a1, a1_read);

  yardl::NDArray<float, 2> a2_read;
  ReadNDArray<float, &ReadFloatingPoint, 2>(r, a2_read);
  ASSERT_EQ(a2, a2_read);

  r.VerifyFinished();
}

TEST(CodedStreamTests, FixedNDArray) {
  std::stringstream ss;
  yardl::FixedNDArray<int, 2, 3> a1({{1, 2, 3}, {4, 5, 6}});
  yardl::FixedNDArray<float, 2, 3> a2({{1, 2, 3}, {4, 5, 6}});
  {
    CodedOutputStream w(ss, 10);
    WriteFixedNDArray<int, &WriteInteger, 2, 3>(w, a1);
    WriteFixedNDArray<float, &WriteFloatingPoint, 2, 3>(w, a2);
  }

  CodedInputStream r(ss, 10);
  yardl::FixedNDArray<int, 2, 3> a1_read;
  ReadFixedNDArray<int, &ReadInteger, 2, 3>(r, a1_read);
  ASSERT_EQ(a1, a1_read);

  yardl::FixedNDArray<float, 2, 3> a2_read;
  ReadFixedNDArray<float, &ReadFloatingPoint, 2, 3>(r, a2_read);
  ASSERT_EQ(a2, a2_read);

  r.VerifyFinished();
}

TEST(CodedStreamTests, BatchedReads_TriviallySerializable_SmallBatches) {
  std::vector<float> expected(128);
  std::iota(expected.begin(), expected.end(), 0);

  std::stringstream ss;
  {
    CodedOutputStream w(ss);
    WriteVector<float, &WriteFloatingPoint>(w, expected);
    WriteInteger(w, 0);
  }

  CodedInputStream r(ss);
  size_t current_block_remaining = 0;
  std::vector<float> inputBuf(7);
  std::vector<float> actual;
  do {
    ReadBlocksIntoVector<float, &ReadFloatingPoint>(r, current_block_remaining, inputBuf);
    actual.insert(actual.end(), inputBuf.begin(), inputBuf.end());

  } while (current_block_remaining > 0);

  ASSERT_EQ(expected, actual);
  r.VerifyFinished();
}

TEST(CodedStreamTests, BatchedReads_NotTriviallySerializable_SmallBatches) {
  std::vector<int> expected(128);
  std::iota(expected.begin(), expected.end(), 0);

  std::stringstream ss;
  {
    CodedOutputStream w(ss);
    WriteVector<int, &WriteInteger>(w, expected);
    WriteInteger(w, 0);
  }

  CodedInputStream r(ss);
  size_t current_block_remaining = 0;
  std::vector<int> inputBuf(7);
  std::vector<int> actual;
  do {
    ReadBlocksIntoVector<int, &ReadInteger>(r, current_block_remaining, inputBuf);
    actual.insert(actual.end(), inputBuf.begin(), inputBuf.end());

  } while (current_block_remaining > 0);

  ASSERT_EQ(expected, actual);
  r.VerifyFinished();
}

TEST(CodedStreamTests, BatchedReads_TriviallySerializable_SingleBatchExactSize) {
  std::vector<float> expected(128);
  std::iota(expected.begin(), expected.end(), 0);

  std::stringstream ss;
  {
    CodedOutputStream w(ss);
    WriteVector<float, &WriteFloatingPoint>(w, expected);
    WriteInteger(w, 0);
  }

  CodedInputStream r(ss);
  size_t current_block_remaining = 0;
  std::vector<float> actual(128);
  ReadBlocksIntoVector<float, &ReadFloatingPoint>(r, current_block_remaining, actual);
  ASSERT_EQ(0, current_block_remaining);
  ASSERT_EQ(expected, actual);
  r.VerifyFinished();
}

TEST(CodedStreamTests, BatchedReads_TriviallySerializable_SingleBatchLargerThanNecessary) {
  std::vector<float> expected(128);
  std::iota(expected.begin(), expected.end(), 0);

  std::stringstream ss;
  {
    CodedOutputStream w(ss);
    WriteVector<float, &WriteFloatingPoint>(w, expected);
    WriteInteger(w, 0);
  }

  CodedInputStream r(ss);
  size_t current_block_remaining = 0;
  std::vector<float> actual(129);
  ReadBlocksIntoVector<float, &ReadFloatingPoint>(r, current_block_remaining, actual);
  ASSERT_EQ(0, current_block_remaining);
  ASSERT_EQ(expected, actual);
  r.VerifyFinished();
}

TEST(CodedStreamTests, TypesThatAreTriviallySerializable) {
  static_assert(IsTriviallySerializable<char>::value);
  static_assert(IsTriviallySerializable<bool>::value);
  static_assert(IsTriviallySerializable<float>::value);
  static_assert(IsTriviallySerializable<double>::value);
  static_assert(IsTriviallySerializable<std::complex<float>>::value);
  static_assert(IsTriviallySerializable<std::complex<double>>::value);
  static_assert(IsTriviallySerializable<std::array<float, 2>>::value);
}

TEST(CodedStreamTests, TypesThatAreNotTriviallySerializable) {
  static_assert(!IsTriviallySerializable<int>::value);
  static_assert(!IsTriviallySerializable<long>::value);
  static_assert(!IsTriviallySerializable<std::vector<float>>::value);
  static_assert(!IsTriviallySerializable<yardl::DynamicNDArray<float>>::value);
  static_assert(!IsTriviallySerializable<yardl::NDArray<float, 2>>::value);
  static_assert(!IsTriviallySerializable<yardl::FixedNDArray<float, 2, 3>>::value);
}

}  // namespace
