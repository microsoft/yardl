// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <gtest/gtest.h>

#include "generated/protocols.h"
#include "generated/types.h"
#include "yardl_testing.h"

using namespace test_model;
using namespace yardl;
using namespace yardl::testing;

#if __cplusplus < 202002L
using year = date::year;
#else
using year = std::chrono::year;
#endif

namespace {

class RoundTripTests : public ::testing::TestWithParam<Format> {
 protected:
  void SetUp() override {
    format_ = GetParam();
  }

  template <typename T>
  std::unique_ptr<T> CreateValidatingWriter() {
    return yardl::testing::CreateValidatingWriter<T>(format_, TestFilename(format_));
  }

  Format format_;
};

TEST_P(RoundTripTests, Scalars) {
  auto tw = CreateValidatingWriter<ScalarsWriterBase>();

  tw->WriteInt32(1);

  RecordWithPrimitives rec;
  rec.bool_field = true;
  rec.int8_field = -33;
  rec.uint8_field = 33;
  rec.uint16_field = -44;
  rec.uint16_field = 44;
  rec.int32_field = -55;
  rec.uint32_field = 55;
  rec.int64_field = -66;
  rec.uint64_field = 66;
  rec.size_field = UINT64_MAX;
  rec.float32_field = 4290.39;
  rec.float64_field = 2234290.39;
  rec.complexfloat32_field = {1.3, 2.2};
  rec.complexfloat64_field = {-2.4, 999.3};
  rec.date_field = Date(year{2022} / 9 / 8);
  rec.time_field = std::chrono::hours(10) + std::chrono::minutes(50) +
                   std::chrono::seconds(25) + std::chrono::milliseconds(777);
  rec.datetime_field = std::chrono::system_clock::now();
  tw->WriteRecord(rec);

  tw->Close();
}

TEST_P(RoundTripTests, ScalarOptionals_NoValue) {
  auto tw = CreateValidatingWriter<ScalarOptionalsWriterBase>();

  std::optional<int> optional_int;
  tw->WriteOptionalInt(optional_int);

  std::optional<SimpleRecord> optional_rec;
  tw->WriteOptionalRecord(optional_rec);

  RecordWithOptionalFields rec_with_optional_field;
  tw->WriteRecordWithOptionalFields(rec_with_optional_field);

  std::optional<RecordWithOptionalFields> optional_rec_with_optional_fields;
  tw->WriteOptionalRecordWithOptionalFields(optional_rec_with_optional_fields);

  tw->Close();
}

TEST_P(RoundTripTests, ScalarOptionals_WithValue) {
  auto tw = CreateValidatingWriter<ScalarOptionalsWriterBase>();

  std::optional<int> optional_int = 55;
  tw->WriteOptionalInt(optional_int);

  std::optional<SimpleRecord> optional_rec = SimpleRecord{8, 9, 10};
  tw->WriteOptionalRecord(optional_rec);

  RecordWithOptionalFields rec_with_optional_field{44};
  tw->WriteRecordWithOptionalFields(rec_with_optional_field);

  std::optional<RecordWithOptionalFields> optional_rec_with_optional_fields{
      RecordWithOptionalFields{66},
  };
  tw->WriteOptionalRecordWithOptionalFields(optional_rec_with_optional_fields);

  tw->Close();
}

TEST_P(RoundTripTests, NestedRecords) {
  auto tw = CreateValidatingWriter<NestedRecordsWriterBase>();

  TupleWithRecords t_written{SimpleRecord{1, 2, 3}, SimpleRecord{4, 5, 6}};
  tw->WriteTupleWithRecords(t_written);

  tw->Close();
}

TEST_P(RoundTripTests, Vlens) {
  auto tw = CreateValidatingWriter<VlensWriterBase>();

  std::vector<int> ints_written = {1, 2, 3};
  tw->WriteIntVector(ints_written);

  std::vector<std::complex<float>> cplx_written = {{1, 2}, {3, 4}, {5, 6}};
  tw->WriteComplexVector(cplx_written);

  RecordWithVlens rec_with_vlens_written{
      {SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}, SimpleRecord{17, 18, 19}},
      11,
      12,
  };
  tw->WriteRecordWithVlens(rec_with_vlens_written);

  std::vector<RecordWithVlens> vlen_of_rec_with_vlens_written =
      {rec_with_vlens_written, rec_with_vlens_written};
  tw->WriteVlenOfRecordWithVlens(vlen_of_rec_with_vlens_written);

  tw->Close();
}

TEST_P(RoundTripTests, Strings) {
  auto tw = CreateValidatingWriter<StringsWriterBase>();

  std::string my_string = "hello world";
  tw->WriteSingleString(my_string);

  RecordWithStrings rec{"Seattle", "Montreal"};
  tw->WriteRecWithString(rec);

  tw->Close();
}

TEST_P(RoundTripTests, OptionalVectors_NoValue) {
  auto tw = CreateValidatingWriter<OptionalVectorsWriterBase>();

  RecordWithOptionalVector rec;
  tw->WriteRecordWithOptionalVector(rec);

  tw->Close();
}

TEST_P(RoundTripTests, OptionalVectors_WithValue) {
  auto tw = CreateValidatingWriter<OptionalVectorsWriterBase>();

  RecordWithOptionalVector rec{{{1, 2, 3, 4}}};
  tw->WriteRecordWithOptionalVector(rec);

  tw->Close();
}

TEST_P(RoundTripTests, FixedVectors) {
  auto tw = CreateValidatingWriter<FixedVectorsWriterBase>();

  std::array<int, 5> int_arr = {1, 2, 3, 4, 5};
  tw->WriteFixedIntVector(int_arr);

  std::array<SimpleRecord, 3> simple_rec_arr =
      {SimpleRecord{1, 2, 3}, SimpleRecord{4, 5, 6}, SimpleRecord{7, 8, 9}};
  tw->WriteFixedSimpleRecordVector(simple_rec_arr);

  std::array<RecordWithVlens, 2> rec_with_vlens_arr = {
      RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
      RecordWithVlens{{{SimpleRecord{17, 28, 39}}}, 133, 142}};
  tw->WriteFixedRecordWithVlensVector(rec_with_vlens_arr);

  RecordWithFixedVectors rec_with_fixed_vectors{
      int_arr,
      simple_rec_arr,
      rec_with_vlens_arr};

  tw->WriteRecordWithFixedVectors(rec_with_fixed_vectors);

  tw->Close();
}

TEST_P(RoundTripTests, FixedArrays) {
  auto tw = CreateValidatingWriter<FixedArraysWriterBase>();

  tw->WriteInts({{1, 2, 3}, {4, 5, 6}});

  tw->WriteFixedSimpleRecordArray({
      {SimpleRecord{1, 2, 3}, SimpleRecord{4, 5, 6}, SimpleRecord{7, 8, 9}},
      {SimpleRecord{11, 12, 13}, SimpleRecord{14, 15, 16}, SimpleRecord{17, 18, 19}},
  });

  tw->WriteFixedRecordWithVlensArray({
      {
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
      },
      {
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
      },
  });

  RecordWithFixedArrays rec{
      {{1, 2, 3}, {4, 5, 6}},
      {
          {SimpleRecord{1, 2, 3}, SimpleRecord{4, 5, 6}, SimpleRecord{7, 8, 9}},
          {SimpleRecord{11, 12, 13}, SimpleRecord{14, 15, 16}, SimpleRecord{17, 18, 19}},
      },
      {
          {
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          },
          {
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          },
      }};

  tw->WriteRecordWithFixedArrays(rec);
  NamedFixedNDArray named = {{1, 2, 3, 4}, {5, 6, 7, 8}};
  tw->WriteNamedArray(named);

  tw->Close();
}

TEST_P(RoundTripTests, NDArrays) {
  auto tw = CreateValidatingWriter<NDArraysWriterBase>();

  NDArray<int, 2> arr = {{1, 2, 3}, {4, 5, 6}};
  ASSERT_EQ(arr.dimension(), 2);
  ASSERT_EQ(arr.shape(0), 2);
  ASSERT_EQ(arr.shape(1), 3);
  tw->WriteInts(arr);

  tw->WriteSimpleRecordArray({
      {SimpleRecord{1, 2, 3}, SimpleRecord{4, 5, 6}, SimpleRecord{7, 8, 9}},
      {SimpleRecord{11, 12, 13}, SimpleRecord{14, 15, 16}, SimpleRecord{17, 18, 19}},
  });

  tw->WriteRecordWithVlensArray({
      {
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
      },
      {
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
      },
  });

  RecordWithNDArrays rec{
      {{1, 2, 3}, {4, 5, 6}},
      {
          {SimpleRecord{1, 2, 3}, SimpleRecord{4, 5, 6}, SimpleRecord{7, 8, 9}},
          {SimpleRecord{11, 12, 13}, SimpleRecord{14, 15, 16}, SimpleRecord{17, 18, 19}},
      },

      {
          {
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          },
          {
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          },
      }};

  tw->WriteRecordWithNDArrays(rec);
  NamedNDArray named = {{1, 2, 3}, {4, 5, 6}};
  tw->WriteNamedArray(named);

  tw->Close();
}

// We optimize storage for NDArrays with a single dimension.
TEST_P(RoundTripTests, NDArraysSingleDimension) {
  auto tw = CreateValidatingWriter<NDArraysSingleDimensionWriterBase>();

  NDArray<int, 1> arr = {{1, 2, 3}};
  ASSERT_EQ(arr.dimension(), 1);
  ASSERT_EQ(arr.shape(0), 3);
  tw->WriteInts(arr);

  tw->WriteSimpleRecordArray({
      {SimpleRecord{1, 2, 3}, SimpleRecord{4, 5, 6}, SimpleRecord{7, 8, 9}},
  });

  tw->WriteRecordWithVlensArray({{
      RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
      RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
  }});

  RecordWithNDArraysSingleDimension rec{
      {{1, 2, 3}},
      {
          {SimpleRecord{1, 2, 3}, SimpleRecord{4, 5, 6}, SimpleRecord{7, 8, 9}},
      },
      {
          {
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          },
      }};

  tw->WriteRecordWithNDArrays(rec);

  tw->Close();
}

TEST_P(RoundTripTests, DynamicNDArrays) {
  auto tw = CreateValidatingWriter<DynamicNDArraysWriterBase>();

  tw->WriteInts({{1, 2, 3}, {4, 5, 6}});

  tw->WriteSimpleRecordArray({
      {SimpleRecord{1, 2, 3}, SimpleRecord{4, 5, 6}, SimpleRecord{7, 8, 9}},
      {SimpleRecord{11, 12, 13}, SimpleRecord{14, 15, 16}, SimpleRecord{17, 18, 19}},
  });

  tw->WriteRecordWithVlensArray({
      {
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
      },
      {
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
      },
  });

  RecordWithDynamicNDArrays rec{
      {{1, 2, 3}, {4, 5, 6}},
      {
          {SimpleRecord{1, 2, 3}, SimpleRecord{4, 5, 6}, SimpleRecord{7, 8, 9}},
          {SimpleRecord{11, 12, 13}, SimpleRecord{14, 15, 16}, SimpleRecord{17, 18, 19}},
      },
      {
          {
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          },
          {
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
              RecordWithVlens{{{SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}}}, 13, 14},
          },
      }};

  tw->WriteRecordWithDynamicNDArrays(rec);

  tw->Close();
}

TEST_P(RoundTripTests, Maps) {
  auto tw = CreateValidatingWriter<MapsWriterBase>();

  tw->WriteStringToInt({{"a", 1}, {"b", 2}, {"c", 3}});
  tw->WriteStringToUnion({{"a", 1}, {"b", "2"}});
  tw->WriteAliasedGeneric({{"a", 1}, {"b", 2}, {"c", 3}});

  tw->Close();
}

TEST_P(RoundTripTests, Unions_FirstOption) {
  auto tw = CreateValidatingWriter<UnionsWriterBase>();

  tw->WriteIntOrSimpleRecord({2});
  tw->WriteIntOrRecordWithVlens({2});
  tw->WriteMonosotateOrIntOrSimpleRecord({});

  tw->Close();
}

TEST_P(RoundTripTests, Unions_SecondOption) {
  auto tw = CreateValidatingWriter<UnionsWriterBase>();

  tw->WriteIntOrSimpleRecord(SimpleRecord{1, 2, 3});
  tw->WriteIntOrRecordWithVlens(
      RecordWithVlens{
          {SimpleRecord{1, 2, 3}, SimpleRecord{7, 8, 9}, SimpleRecord{17, 18, 19}},
          11,
          12,
      });
  tw->WriteMonosotateOrIntOrSimpleRecord({6});

  tw->Close();
}

TEST_P(RoundTripTests, Enums) {
  auto tw = CreateValidatingWriter<EnumsWriterBase>();
  tw->WriteSingle(Fruits::kApple);
  tw->WriteVec({Fruits::kApple, Fruits::kBanana});
  tw->WriteSize(SizeBasedEnum::kC);

  tw->Close();
}

TEST_P(RoundTripTests, SimpleDatasets) {
  auto tw = CreateValidatingWriter<StreamsWriterBase>();

  for (int i = 0; i < 10; i++) {
    tw->WriteIntData(i);
  }

  tw->WriteIntData({10, 11, 12, 13, 14, 15, 16, 17, 18, 19});
  tw->WriteIntData({20, 21, 22, 23, 24, 25});
  tw->EndIntData();

  tw->WriteOptionalIntData({10, 11, {}, 13, 14, 15, {}, 17, 18, 19});
  tw->EndOptionalIntData();

  tw->WriteRecordWithOptionalVectorData({{{1, 2, 3, 4}}});
  tw->WriteRecordWithOptionalVectorData({{{5, 6}}});
  tw->WriteRecordWithOptionalVectorData({{{7, 8, 9, 10, 11}}});

  std::vector<RecordWithOptionalVector> recs;
  for (int i = 0; i < 10; i++) {
    recs.push_back({{{i, i + 1, i + 2}}});
  }
  tw->WriteRecordWithOptionalVectorData(recs);
  tw->EndRecordWithOptionalVectorData();

  tw->WriteFixedVector({1, 2, 3});
  tw->WriteFixedVector({4, 5, 6});
  tw->EndFixedVector();

  tw->Close();
}

TEST_P(RoundTripTests, SimpleDatasets_Empty) {
  auto tw = CreateValidatingWriter<StreamsWriterBase>();

  tw->EndIntData();

  tw->EndOptionalIntData();

  tw->EndRecordWithOptionalVectorData();

  tw->EndFixedVector();

  tw->Close();
}

TEST_P(RoundTripTests, StreamsOfUnions) {
  auto tw = CreateValidatingWriter<StreamsOfUnionsWriterBase>();

  tw->WriteIntOrSimpleRecord(1);
  tw->WriteIntOrSimpleRecord(SimpleRecord{3, 4, 5});
  tw->WriteIntOrSimpleRecord(2);
  tw->EndIntOrSimpleRecord();

  tw->WriteNullableIntOrSimpleRecord(std::monostate{});
  tw->WriteNullableIntOrSimpleRecord(1);
  tw->WriteNullableIntOrSimpleRecord(SimpleRecord{3, 4, 5});
  tw->WriteNullableIntOrSimpleRecord(2);
  tw->EndNullableIntOrSimpleRecord();

  tw->Close();
}

TEST_P(RoundTripTests, Simple) {
  auto tw = CreateValidatingWriter<SimpleGenericsWriterBase>();

  tw->WriteFloatImage({{3.0, 4.0, 5.0}, {6.0, 7.0, 8.0}});
  tw->WriteIntImage({{3, 4, 5}, {6, 7, 8}});
  tw->WriteIntImageAlternateSyntax({{13, 14}, {16, 17}});

  tw->WriteStringImage({{"a", "b", "c"}, {"d", "e", "f"}});

  tw->WriteIntFloatTuple({1, 66.2});
  tw->WriteFloatFloatTuple({99.43, 66.2});
  tw->WriteIntFloatTupleAlternateSyntax({2, 62.2});
  tw->WriteIntStringTuple({1, "bonjour"});

  tw->WriteStreamOfTypeVariants(Image<float>{{3, 4}});
  tw->WriteStreamOfTypeVariants(Image<double>{{3, 4, 5}});
  tw->EndStreamOfTypeVariants();

  tw->Close();
}

TEST_P(RoundTripTests, Advanced) {
  auto tw = CreateValidatingWriter<AdvancedGenericsWriterBase>();

  Image<int> i1 = {{3, 4, 5}, {6, 7, 8}};
  Image<int> i2 = {{30, 40, 50}, {60, 70, 80}};
  Image<int> i3 = {{300, 400, 500}, {600, 700, 800}};
  Image<int> i4 = {{3000, 4000, 5000}, {6000, 7000, 8000}};

  tw->WriteIntImageImage({{i1, i2}, {i3, i4}});

  GenericRecord<int, std::string> r1{
      7,
      "hello",
      {77, 88, 99},
      {{"a", "b", "c"}, {"d", "e", "f"}}};

  tw->WriteGenericRecord1(r1);
  tw->WriteTupleOfOptionals({std::nullopt, "hello"});
  tw->WriteTupleOfOptionalsAlternateSyntax({5, std::nullopt});
  tw->WriteTupleOfVectors({{1, 2, 3}, {33.3, 44.4, 55.5}});

  tw->Close();
}

TEST_P(RoundTripTests, Aliases) {
  auto tw = CreateValidatingWriter<AliasesWriterBase>();

  AliasedString as = "hello";
  tw->WriteAliasedString(as);

  AliasedEnum ae = Fruits::kBanana;
  tw->WriteAliasedEnum(ae);

  AliasedOpenGeneric<AliasedString, AliasedEnum> aog = {as, ae};
  tw->WriteAliasedOpenGeneric(aog);

  AliasedClosedGeneric acg = {as, ae};
  tw->WriteAliasedClosedGeneric(acg);

  AliasedOptional aliased_optional = 42;
  tw->WriteAliasedOptional(aliased_optional);

  AliasedGenericOptional<float> aliased_generic_optional = 94.2;
  tw->WriteAliasedGenericOptional(aliased_generic_optional);

  AliasedGenericUnion2<AliasedString, AliasedEnum> aliased_generic_union2 = as;
  tw->WriteAliasedGenericUnion2(aliased_generic_union2);

  AliasedGenericVector<float> aliased_generic_vector = {1.0, 2.0, 3.0};
  tw->WriteAliasedGenericVector(aliased_generic_vector);

  AliasedGenericFixedVector<float> aliased_generic_fixed_vector = {1.0, 2.0, 3.0};
  tw->WriteAliasedGenericFixedVector(aliased_generic_fixed_vector);

  tw->WriteStreamOfAliasedGenericUnion2(as);
  tw->WriteStreamOfAliasedGenericUnion2(ae);
  tw->EndStreamOfAliasedGenericUnion2();

  tw->Close();
}

TEST_P(RoundTripTests, ReservedNames) {
  auto tw = CreateValidatingWriter<ProtocolWithKeywordStepsWriterBase>();

  RecordWithKeywordFields rec;
  rec.int_field = "some string";
  rec.sizeof_field = {{1, 2, 3}, {4, 5, 6}};
  rec.if_field = EnumWithKeywordSymbols::kCatch;
  ASSERT_EQ(rec.Return(), 6);
  tw->WriteInt(rec);
  tw->EndInt();

  tw->WriteFloat(EnumWithKeywordSymbols::kTry);

  tw->Close();
}

INSTANTIATE_TEST_SUITE_P(,
                         RoundTripTests,
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
