// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <gtest/gtest.h>

#include "../generated/hdf5/protocols.h"
#include "../generated/types.h"
#include "../yardl_testing.h"

using namespace test_model;
using namespace test_model::hdf5;
using namespace yardl;
using namespace yardl::testing;

namespace {

TEST(Hdf5Tests, MultipleProtocolsInFile) {
  std::string filename = TestFilename(Format::kHdf5);

  {
    ScalarsWriter w(filename);
    w.WriteInt32(1);
    w.WriteRecord({true});
  }

  {
    EnumsWriter w(filename);
    w.WriteSingle(Fruits::kBanana);
    w.WriteVec({Fruits::kApple, Fruits::kBanana});
  }

  ScalarsReader r1(filename);
  int i;
  r1.ReadInt32(i);
  ASSERT_EQ(i, 1);
  RecordWithPrimitives rec;
  r1.ReadRecord(rec);
  ASSERT_TRUE(rec.bool_field);

  EnumsReader r2(filename);

  Fruits fruit;
  r2.ReadSingle(fruit);
  ASSERT_EQ(fruit, Fruits::kBanana);
  std::vector<Fruits> fruit_vec;
  r2.ReadVec(fruit_vec);
  ASSERT_EQ(fruit_vec.size(), 2);
  ASSERT_EQ(fruit_vec[0], Fruits::kApple);
  ASSERT_EQ(fruit_vec[1], Fruits::kBanana);
}

TEST(Hdf5Tests, ProtocolAlreadyExists) {
  std::string filename = TestFilename(Format::kHdf5);

  {
    ScalarsWriter w(filename);
    w.WriteInt32(1);
    w.WriteRecord({true});
  }

  ASSERT_ANY_THROW(static_cast<void>(ScalarsWriter(filename)));
}

TEST(Hdf5Tests, WrongProtocol) {
  std::string filename = TestFilename(Format::kHdf5);

  {
    ScalarsWriter w(filename);
    w.WriteInt32(1);
    w.WriteRecord({true});
  }

  ASSERT_ANY_THROW(static_cast<void>(ScalarOptionalsReader(filename)));
}

TEST(Hdf5Tests, ProtocolSchemaMismatch) {
  std::string filename = TestFilename(Format::kHdf5);

  {
    ScalarsWriter w(filename);
    w.WriteInt32(1);
    w.WriteRecord({true});
  }

  ASSERT_ANY_THROW(static_cast<void>(yardl::hdf5::Hdf5Reader(filename, "Scalars", "{}")));
}

TEST(Hdf5Tests, StreamsOfUnionsInSeparateDatasets) {
  std::string filename = TestFilename(Format::kHdf5);
  {
    StreamsOfUnionsWriter w(filename);
    w.WriteIntOrSimpleRecord(1);
    w.WriteIntOrSimpleRecord(SimpleRecord{1, 2, 3});
    w.EndIntOrSimpleRecord();

    w.WriteNullableIntOrSimpleRecord(std::monostate{});
    w.WriteNullableIntOrSimpleRecord(2);
    w.WriteNullableIntOrSimpleRecord(SimpleRecord{1, 2, 3});
    w.EndNullableIntOrSimpleRecord();

    w.Close();
  }

  H5::H5File file(filename, H5F_ACC_RDONLY);
  auto protocolGroup = file.openGroup("StreamsOfUnions");
  auto stepGroup1 = protocolGroup.openGroup("intOrSimpleRecord");
  stepGroup1.openDataSet("$index");
  stepGroup1.openDataSet("int32");
  stepGroup1.openDataSet("SimpleRecord");
  auto stepGroup2 = protocolGroup.openGroup("nullableIntOrSimpleRecord");
  stepGroup2.openDataSet("$index");
  stepGroup2.openDataSet("int32");
  stepGroup2.openDataSet("SimpleRecord");
}

TEST(Hdf5Tests, StreamsOfAliasedUnionsInSeparateDatasets) {
  std::string filename = TestFilename(Format::kHdf5);
  {
    StreamsOfAliasedUnionsWriter w(filename);
    w.WriteIntOrSimpleRecord(1);
    w.WriteIntOrSimpleRecord(SimpleRecord{1, 2, 3});
    w.EndIntOrSimpleRecord();

    w.WriteNullableIntOrSimpleRecord(std::monostate{});
    w.WriteNullableIntOrSimpleRecord(2);
    w.WriteNullableIntOrSimpleRecord(SimpleRecord{1, 2, 3});
    w.EndNullableIntOrSimpleRecord();

    w.Close();
  }

  H5::H5File file(filename, H5F_ACC_RDONLY);
  auto protocolGroup = file.openGroup("StreamsOfAliasedUnions");
  auto stepGroup1 = protocolGroup.openGroup("intOrSimpleRecord");
  stepGroup1.openDataSet("$index");
  stepGroup1.openDataSet("int32");
  stepGroup1.openDataSet("SimpleRecord");
  auto stepGroup2 = protocolGroup.openGroup("nullableIntOrSimpleRecord");
  stepGroup2.openDataSet("$index");
  stepGroup2.openDataSet("int32");
  stepGroup2.openDataSet("SimpleRecord");
}

}  // namespace
