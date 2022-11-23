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
  EXPECT_EQ(r.AccessArrayFieldElement(), r.array_field.at(0, 1));
  EXPECT_EQ(r.AccessArrayFieldElementByName(), r.array_field.at(0, 1));

  EXPECT_EQ(r.AccessOtherComputedField(), r.AccessIntField());

  r.vector_of_vectors_field = {{1, 2, 3}, {4, 5, 6}};
  EXPECT_EQ(r.AccessVectorOfVectorsField(), r.vector_of_vectors_field[1][2]);
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

  EXPECT_EQ(r.FixedVectorSize(), 3);
}

TEST(ComputedFieldsTests, VectorSize) {
  RecordWithComputedFields r;
  EXPECT_EQ(r.VectorSize(), 0);
  r.vector_field = {1, 2, 3};
  EXPECT_EQ(r.VectorSize(), 3);

  EXPECT_EQ(r.FixedVectorSize(), 3);
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

  ASSERT_EQ(r.FixedArraySize(), r.fixed_array_field.size());
  ASSERT_EQ(r.FixedArrayXSize(), r.fixed_array_field.shape(0));
  ASSERT_EQ(r.FixedArray0Size(), r.fixed_array_field.shape(0));
}

TEST(ComputedFieldsTest, SwitchExpression) {
  RecordWithComputedFields r;
  r.optional_named_array = {{1, 2, 3}, {4, 5, 6}};
  ASSERT_EQ(r.OptionalNamedArrayLength(), 6);
  ASSERT_EQ(r.OptionalNamedArrayLengthWithDiscard(), 6);
  r.optional_named_array = {};
  ASSERT_EQ(r.OptionalNamedArrayLength(), 0);
  ASSERT_EQ(r.OptionalNamedArrayLengthWithDiscard(), 0);
  static_assert(std::is_same_v<decltype(r.OptionalNamedArrayLength()), size_t>);
  static_assert(std::is_same_v<decltype(r.OptionalNamedArrayLengthWithDiscard()), size_t>);

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
  r.union_with_nested_generic_union = GenericRecordWithComputedFields<std::string, float>{"hi"};
  ASSERT_EQ(r.NestedSwitch(), 10);
  ASSERT_EQ(r.UseNestedComputedField(), 0);
  r.union_with_nested_generic_union = GenericRecordWithComputedFields<std::string, float>{42.0f};
  ASSERT_EQ(r.NestedSwitch(), 20);
  ASSERT_EQ(r.UseNestedComputedField(), 1);

  GenericRecordWithComputedFields<int, double> gr;
  gr.f1 = 42;
  ASSERT_EQ(gr.TypeIndex(), 0);
  gr.f1 = 42.0;
  ASSERT_EQ(gr.TypeIndex(), 1);
}

}  // namespace
