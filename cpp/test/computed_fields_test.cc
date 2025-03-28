// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <filesystem>
#include <type_traits>
#include <vector>

#include <gtest/gtest.h>

#include "generated/types.h"

using namespace test_model;

namespace {

TEST(ComputedFieldsTests, AccessFields) {
  RecordWithComputedFields r;

  r.int_field = 42;
  EXPECT_EQ(r.AccessIntField(), r.int_field);

  r.string_field = "hello";
  EXPECT_EQ(r.AccessStringField(), r.string_field);

  r.tuple_field = {2, 3};
  EXPECT_EQ(r.AccessTupleField(), r.tuple_field);
  EXPECT_EQ(r.AccessNestedTupleField(), r.tuple_field.v2);

  r.array_field = {{1, 2, 3}, {4, 5, 6}};
  EXPECT_EQ(r.AccessArrayField(), r.array_field);
  EXPECT_EQ(r.AccessArrayFieldElement(), yardl::at(r.array_field, 0, 1));
  EXPECT_EQ(r.AccessArrayFieldElementByName(), yardl::at(r.array_field, 0, 1));

  EXPECT_EQ(r.AccessOtherComputedField(), r.AccessIntField());

  r.vector_of_vectors_field = {{1, 2, 3}, {4, 5, 6}};
  EXPECT_EQ(r.AccessVectorOfVectorsField(), r.vector_of_vectors_field[1][2]);

  r.map_field = {{"hello", "world"}, {"world", "bye"}};
  EXPECT_EQ(r.AccessMap(), r.map_field);
  EXPECT_EQ(r.AccessMapEntry(), "world");
  EXPECT_EQ(r.AccessMapEntryWithComputedField(), "world");
  EXPECT_EQ(r.AccessMapEntryWithComputedFieldNested(), "bye");
  EXPECT_THROW(r.AccessMissingMapEntry(), std::out_of_range);
}

TEST(ComputedFieldsTest, Literals) {
  RecordWithComputedFields r;

  static_assert(std::is_same_v<decltype(r.IntLiteral()), uint8_t>);
  static_assert(std::is_same_v<decltype(r.LargeNegativeInt64Literal()), int64_t>);
  static_assert(std::is_same_v<decltype(r.LargeUInt64Literal()), uint64_t>);

  ASSERT_EQ(r.StringLiteral(), "hello");
  ASSERT_EQ(r.StringLiteral2(), "hello");
  ASSERT_EQ(r.StringLiteral3(), "hello");
  ASSERT_EQ(r.StringLiteral4(), "hello");
}

TEST(ComputedFieldsTest, DimensionIndex) {
  RecordWithComputedFields r;
  EXPECT_EQ(r.ArrayDimensionXIndex(), 0);
  EXPECT_EQ(r.ArrayDimensionYIndex(), 1);
  r.string_field = "y";
  EXPECT_EQ(r.ArrayDimensionIndexFromStringField(), 1);
  r.string_field = "missing";
  ASSERT_ANY_THROW(r.ArrayDimensionIndexFromStringField());
}

TEST(ComputedFieldsTest, DimensionCount) {
  RecordWithComputedFields r;
  EXPECT_EQ(r.ArrayDimensionCount(), 2);
  r.dynamic_array_field = {{1, 2, 3}, {4, 5, 6}};
  EXPECT_EQ(r.DynamicArrayDimensionCount(), 2);
  r.dynamic_array_field = {1, 2, 3};
  EXPECT_EQ(r.DynamicArrayDimensionCount(), 1);
}

TEST(ComputedFieldsTests, VectorSize) {
  RecordWithComputedFields r;
  EXPECT_EQ(r.VectorSize(), 0);
  r.vector_field = {1, 2, 3};
  EXPECT_EQ(r.VectorSize(), 3);

  EXPECT_EQ(r.FixedVectorSize(), 3);
}

TEST(ComputedFieldsTests, MapSize) {
  RecordWithComputedFields r;
  r.map_field = {{"hello", "bonjour"}, {"world", "monde"}};
  EXPECT_EQ(r.MapSize(), 2);
}

TEST(ComputedFieldsTests, ArraySize) {
  RecordWithComputedFields r;
  r.array_field = {{1, 2, 3}, {4, 5, 6}};
  ASSERT_EQ(r.ArraySize(), 6);
  ASSERT_EQ(r.ArrayXSize(), 2);
  ASSERT_EQ(r.ArrayYSize(), 3);
  ASSERT_EQ(r.Array0Size(), 2);
  ASSERT_EQ(r.Array1Size(), 3);

  ASSERT_EQ(r.ArraySizeFromIntField(), 2);
  r.int_field = 1;
  ASSERT_EQ(r.ArraySizeFromIntField(), 3);

  r.string_field = "x";
  ASSERT_EQ(r.ArraySizeFromStringField(), 2);
  r.string_field = "y";
  ASSERT_EQ(r.ArraySizeFromStringField(), 3);
  r.string_field = "missing";
  ASSERT_ANY_THROW(r.ArraySizeFromStringField());

  r.tuple_field.v1 = 1;
  ASSERT_EQ(r.ArraySizeFromNestedIntField(), 3);

  ASSERT_EQ(r.FixedArraySize(), yardl::size(r.fixed_array_field));
  ASSERT_EQ(r.FixedArrayXSize(), yardl::shape(r.fixed_array_field, 0));
  ASSERT_EQ(r.FixedArray0Size(), yardl::shape(r.fixed_array_field, 0));

  r.array_field_map_dimensions = {{1, 2, 3}, {4, 5, 6}};
  ASSERT_EQ(r.ArrayFieldMapDimensionsXSize(), 2);
}

TEST(ComputedFieldsTest, SwitchExpression) {
  RecordWithComputedFields r;
  r.optional_named_array = {{1, 2, 3}, {4, 5, 6}};
  ASSERT_EQ(r.OptionalNamedArrayLength(), 6);
  ASSERT_EQ(r.OptionalNamedArrayLengthWithDiscard(), 6);
  r.optional_named_array = {};
  ASSERT_EQ(r.OptionalNamedArrayLength(), 0);
  ASSERT_EQ(r.OptionalNamedArrayLengthWithDiscard(), 0);
  static_assert(std::is_same_v<decltype(r.OptionalNamedArrayLength()), yardl::Size>);
  static_assert(std::is_same_v<decltype(r.OptionalNamedArrayLengthWithDiscard()), yardl::Size>);

  r.int_float_union = 42;
  ASSERT_EQ(r.IntFloatUnionAsFloat(), 42.0f);
  static_assert(std::is_same_v<decltype(r.IntFloatUnionAsFloat()), float>);

  r.nullable_int_float_union = {};
  ASSERT_EQ(r.NullableIntFloatUnionString(), "null");
  r.nullable_int_float_union = 42;
  ASSERT_EQ(r.NullableIntFloatUnionString(), "int");
  r.nullable_int_float_union = 42.0f;
  ASSERT_EQ(r.NullableIntFloatUnionString(), "float");

  r.union_with_nested_generic_union = 42;
  ASSERT_EQ(r.NestedSwitch(), -1);
  ASSERT_EQ(r.UseNestedComputedField(), -1);
  r.union_with_nested_generic_union = basic_types::GenericRecordWithComputedFields<std::string, float>{"hi"};
  ASSERT_EQ(r.NestedSwitch(), 10);
  ASSERT_EQ(r.UseNestedComputedField(), 0);
  r.union_with_nested_generic_union = basic_types::GenericRecordWithComputedFields<std::string, float>{42.0f};
  ASSERT_EQ(r.NestedSwitch(), 20);
  ASSERT_EQ(r.UseNestedComputedField(), 1);

  basic_types::GenericRecordWithComputedFields<int, double> gr;
  gr.f1 = 42;
  ASSERT_EQ(gr.TypeIndex(), 0);
  gr.f1 = 42.0;
  ASSERT_EQ(gr.TypeIndex(), 1);
}

TEST(ComputedFieldsTest, Arithmetic) {
  RecordWithComputedFields r;
  ASSERT_EQ(r.Arithmetic1(), 3);
  static_assert(std::is_same_v<decltype(r.Arithmetic1()), int32_t>);
  ASSERT_EQ(r.Arithmetic2(), 11);
  ASSERT_EQ(r.Arithmetic3(), 13);

  r.array_field = {{1, 2, 3}, {4, 5, 6}};
  r.int_field = 1;
  ASSERT_EQ(r.Arithmetic4(), 5);
  ASSERT_EQ(r.Arithmetic5(), 3);

  ASSERT_EQ(r.Arithmetic6(), 3);

  ASSERT_EQ(r.Arithmetic7(), 49.0);
  static_assert(std::is_same_v<decltype(r.Arithmetic7()), double>);

  r.complexfloat32_field = {2.0f, 3.0f};
  ASSERT_EQ(r.Arithmetic8(), std::complex<float>(6.0f, 9.0f));

  ASSERT_EQ(r.Arithmetic9(), 2.2);
  static_assert(std::is_same_v<decltype(r.Arithmetic9()), double>);

  ASSERT_EQ(r.Arithmetic10(), 1e10 + 9e9);

  ASSERT_EQ(r.Arithmetic11(), -(5.3));
}

TEST(ComputedFieldsTest, Casting) {
  RecordWithComputedFields r;
  r.int_field = 42;
  ASSERT_EQ(r.CastIntToFloat(), 42.0f);
  static_assert(std::is_same_v<decltype(r.CastIntToFloat()), float>);

  r.float32_field = 42.9f;
  ASSERT_EQ(r.CastFloatToInt(), 42);
  static_assert(std::is_same_v<decltype(r.CastFloatToInt()), int32_t>);

  ASSERT_EQ(r.CastPower(), 49);
  static_assert(std::is_same_v<decltype(r.CastPower()), int32_t>);

  r.complexfloat32_field = {2.0f, 3.0f};
  r.complexfloat64_field = {2.0, 3.0};
  ASSERT_EQ(r.CastComplex32ToComplex64(), std::complex<double>(2.0, 3.0));
  ASSERT_EQ(r.CastComplex64ToComplex32(), std::complex<float>(2.0f, 3.0f));

  ASSERT_EQ(r.CastFloatToComplex(), std::complex<float>(66.6f, 0.0f));
  static_assert(std::is_same_v<decltype(r.CastFloatToComplex()), std::complex<float>>);
}

}  // namespace
