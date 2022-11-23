// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <gtest/gtest.h>

#include "generated/protocols.h"
#include "generated/types.h"

using namespace test_model;
using namespace yardl;

namespace {

class TestStateTestWriter : public StateTestWriterBase {
  void WriteAnIntImpl([[maybe_unused]] int32_t const& value) override {}
  void WriteAStreamImpl([[maybe_unused]] int32_t const& value) override {}
  void EndAStreamImpl() override {}
  void WriteAnotherIntImpl([[maybe_unused]] int32_t const& value) override {}
};

TEST(WriterStateTest, ProperSequenceWrite) {
  TestStateTestWriter w;
  w.WriteAnInt(1);

  w.WriteAStream(1);
  w.WriteAStream(1);
  w.EndAStream();

  w.WriteAnotherInt(1);

  w.Close();
}

TEST(WriterStateTest, ProperSequenceWriteEmptyStream) {
  TestStateTestWriter w;
  w.WriteAnInt(1);

  w.EndAStream();

  w.WriteAnotherInt(1);

  w.Close();
}

TEST(WriterStateTest, MissingFirstStep) {
  TestStateTestWriter w;
  ASSERT_ANY_THROW(w.WriteAStream(1));
}

TEST(WriterStateTest, MissingEndStream) {
  TestStateTestWriter w;
  w.WriteAnInt(1);

  w.WriteAStream(1);

  ASSERT_ANY_THROW(w.WriteAnotherInt(1));
}

TEST(WriterStateTest, PrematureClose) {
  TestStateTestWriter w;
  w.WriteAnInt(1);
  ASSERT_ANY_THROW(w.Close());
}

class TestStateTestReader : public StateTestReaderBase {
  void ReadAnIntImpl([[maybe_unused]] int32_t& value){};
  bool ReadAStreamImpl([[maybe_unused]] int32_t& value) { return stream_count_++ < 2; };
  void ReadAnotherIntImpl([[maybe_unused]] int32_t& value){};

 private:
  int stream_count_ = 0;
};

TEST(ReaderStateTest, ProperSequenceRead) {
  TestStateTestReader r;
  int i;
  r.ReadAnInt(i);

  while (r.ReadAStream(i)) {
  }

  r.ReadAnotherInt(i);
  r.Close();
}

TEST(ReaderStateTest, ProperSequenceReadWithBatches) {
  TestStateTestReader r;
  int i;
  r.ReadAnInt(i);

  std::vector<int> batch(10);
  ASSERT_TRUE(r.ReadAStream(batch));
  ASSERT_EQ(batch.size(), 2);
  ASSERT_FALSE(r.ReadAStream(batch));
  ASSERT_EQ(batch.size(), 0);

  r.ReadAnotherInt(i);
  r.Close();
}

TEST(ReaderStateTest, NotObservingEndOfStreamAfterBatchCall) {
  TestStateTestReader r;
  int i;
  r.ReadAnInt(i);

  std::vector<int> batch(10);
  ASSERT_TRUE(r.ReadAStream(batch));
  ASSERT_EQ(batch.size(), 2);

  r.ReadAnotherInt(i);
  r.Close();
}

TEST(ReaderStateTest, ProperSequenceReadWithBatches_ReadLastWithoutBatch) {
  TestStateTestReader r;
  int i;
  r.ReadAnInt(i);

  std::vector<int> batch(10);
  ASSERT_TRUE(r.ReadAStream(batch));
  ASSERT_EQ(batch.size(), 2);
  ASSERT_FALSE(r.ReadAStream(batch[0]));

  r.ReadAnotherInt(i);
  r.Close();
}

TEST(ReaderStateTest, MissingFirstStep) {
  TestStateTestReader r;
  int i;
  ASSERT_ANY_THROW(static_cast<void>(r.ReadAStream(i)));
}

TEST(ReaderStateTest, ReadStreamPastEnd) {
  TestStateTestReader r;
  int i;
  r.ReadAnInt(i);
  while (r.ReadAStream(i)) {
  }
  ASSERT_ANY_THROW(static_cast<void>(r.ReadAStream(i)));
}

TEST(ReaderStateTest, PrematureClose) {
  TestStateTestReader r;
  int i;
  r.ReadAnInt(i);
  ASSERT_ANY_THROW(r.Close());
}

}  // namespace
