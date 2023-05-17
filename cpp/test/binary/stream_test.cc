// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include "../generated/yardl/detail/stream/stream.h"

#include <fstream>
#include <iostream>
#include <memory>

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "../yardl_testing.h"

using namespace yardl::stream;
using namespace yardl::testing;
using ::testing::StrictMock;

namespace {

struct MockOStream {
  MOCK_METHOD(void, write, (char const* data, std::streamsize size), (const));
  MOCK_METHOD(void, flush, (), ());
  MOCK_METHOD(bool, bad, (), (const));
  MOCK_METHOD(void, Destructed, ());
  ~MockOStream() { Destructed(); }
};

struct MockOStreamNoFlush {
  MOCK_METHOD(void, write, (char const* data, std::streamsize size), (const));
  MOCK_METHOD(bool, bad, (), (const));
};

TEST(WritableStreamTests, RefConstructor) {
  StrictMock<MockOStream> s;
  EXPECT_CALL(s, write("hello", 5)).Times(1);
  EXPECT_CALL(s, flush()).Times(1);
  EXPECT_CALL(s, bad()).Times(1).WillOnce(::testing::Return(true));
  {
    WritableStream stream(s);
    stream.Write("hello", 5);
    ASSERT_TRUE(stream.Bad());
    stream.Flush();
    ASSERT_TRUE(::testing::Mock::VerifyAndClearExpectations(&s));
  }
  EXPECT_CALL(s, Destructed()).Times(1);
}

TEST(WritableStreamTests, UniquePtrConstructor) {
  auto s = std::make_unique<StrictMock<MockOStream>>();
  EXPECT_CALL(*s, write("hello", 5)).Times(1);
  EXPECT_CALL(*s, flush()).Times(1);
  EXPECT_CALL(*s, bad()).Times(1).WillOnce(::testing::Return(true));
  EXPECT_CALL(*s, Destructed()).Times(1);
  {
    WritableStream stream(std::move(s));
    stream.Write("hello", 5);
    stream.Flush();
    ASSERT_TRUE(stream.Bad());
  }
}

TEST(WritableStreamTests, SharedPtrConstructor) {
  auto s = std::make_shared<StrictMock<MockOStream>>();
  EXPECT_CALL(*s, write("hello", 5)).Times(1);
  EXPECT_CALL(*s, bad()).Times(1).WillOnce(::testing::Return(true));
  EXPECT_CALL(*s, flush()).Times(1);
  {
    WritableStream stream(s);
    ASSERT_FALSE(s.unique());
    stream.Write("hello", 5);
    ASSERT_TRUE(stream.Bad());
    stream.Flush();
    ASSERT_TRUE(::testing::Mock::VerifyAndClearExpectations(&*s));
  }
  ASSERT_TRUE(s.unique());
  EXPECT_CALL(*s, Destructed()).Times(1);
}

TEST(WritableStreamTests, NoFlushMethod) {
  StrictMock<MockOStreamNoFlush> s;
  EXPECT_CALL(s, write("hello", 5)).Times(1);
  {
    WritableStream stream(s);
    stream.Write("hello", 5);
    stream.Flush();
  }
}

TEST(WritableStreamTests, FilenameConstructor) {
  auto filename = yardl::testing::TestFilename(Format::kBinary);
  {
    WritableStream stream(filename);
    stream.Write("hello", 5);
  }

  ReadableStream stream(filename);
  std::vector<char> buffer(5);
  stream.Read(buffer.data(), buffer.size());
  ASSERT_EQ("hello", std::string(buffer.data(), buffer.size()));
  ASSERT_EQ(5, stream.GCount());
  ASSERT_FALSE(stream.Eof());
}

struct MockIStream {
  MOCK_METHOD(void, read, (char* data, std::streamsize size), (const));
  MOCK_METHOD(bool, eof, (), (const));
  MOCK_METHOD(std::streamsize, gcount, (), (const));
  MOCK_METHOD(void, Destructed, ());
  ~MockIStream() { Destructed(); }
};

TEST(ReadableStreamTests, RefConstructor) {
  StrictMock<MockIStream> s;
  std::vector<char> buffer(5);
  EXPECT_CALL(s, read(buffer.data(), buffer.size())).Times(1);
  EXPECT_CALL(s, gcount()).Times(1).WillOnce(::testing::Return(5));
  EXPECT_CALL(s, eof()).Times(1).WillOnce(::testing::Return(true));
  {
    ReadableStream stream(s);
    stream.Read(buffer.data(), buffer.size());
    ASSERT_EQ(5, stream.GCount());
    ASSERT_TRUE(stream.Eof());
    ASSERT_TRUE(::testing::Mock::VerifyAndClearExpectations(&s));
  }
  EXPECT_CALL(s, Destructed()).Times(1);
}

TEST(ReadableStreamTests, UniquePtrConstructor) {
  auto s = std::make_unique<StrictMock<MockIStream>>();
  std::vector<char> buffer(5);
  EXPECT_CALL(*s, read(buffer.data(), buffer.size())).Times(1);
  EXPECT_CALL(*s, gcount()).Times(1).WillOnce(::testing::Return(5));
  EXPECT_CALL(*s, eof()).Times(1).WillOnce(::testing::Return(true));
  EXPECT_CALL(*s, Destructed()).Times(1);
  {
    ReadableStream stream(std::move(s));
    stream.Read(buffer.data(), buffer.size());
    ASSERT_EQ(5, stream.GCount());
    ASSERT_TRUE(stream.Eof());
  }
}

TEST(ReadableStreamTests, SharedPtrConstructor) {
  auto s = std::make_shared<StrictMock<MockIStream>>();
  std::vector<char> buffer(5);
  EXPECT_CALL(*s, read(buffer.data(), buffer.size())).Times(1);
  EXPECT_CALL(*s, gcount()).Times(1).WillOnce(::testing::Return(5));
  EXPECT_CALL(*s, eof()).Times(1).WillOnce(::testing::Return(true));
  {
    ReadableStream stream(s);
    ASSERT_FALSE(s.unique());
    stream.Read(buffer.data(), buffer.size());
    ASSERT_EQ(5, stream.GCount());
    ASSERT_TRUE(stream.Eof());
    ASSERT_TRUE(::testing::Mock::VerifyAndClearExpectations(&*s));
  }
  ASSERT_TRUE(s.unique());
  EXPECT_CALL(*s, Destructed()).Times(1);
}
}  // namespace
