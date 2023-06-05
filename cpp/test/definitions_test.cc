// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <filesystem>
#include <type_traits>
#include <vector>

#include <gtest/gtest.h>

#include "generated/types.h"

using namespace test_model;

namespace {

TEST(DefinitionsTest, PrimitiveAliasesHaveExpectedType) {
  static_assert(std::is_same_v<decltype(RecordWithPrimitiveAliases::byte_field), uint8_t>);
  static_assert(std::is_same_v<decltype(RecordWithPrimitiveAliases::int_field), int32_t>);
  static_assert(std::is_same_v<decltype(RecordWithPrimitiveAliases::uint_field), uint32_t>);
  static_assert(std::is_same_v<decltype(RecordWithPrimitiveAliases::long_field), int64_t>);
  static_assert(std::is_same_v<decltype(RecordWithPrimitiveAliases::ulong_field), uint64_t>);
  static_assert(std::is_same_v<decltype(RecordWithPrimitiveAliases::float_field), float>);
  static_assert(std::is_same_v<decltype(RecordWithPrimitiveAliases::double_field), double>);
  static_assert(std::is_same_v<decltype(RecordWithPrimitiveAliases::complexfloat_field), std::complex<float>>);
  static_assert(std::is_same_v<decltype(RecordWithPrimitiveAliases::complexdouble_field), std::complex<double>>);
}

TEST(DefinitionsTests, VectorFieldsHaveExpectedType) {
  static_assert(std::is_same_v<decltype(RecordWithVectors::default_vector), std::vector<int>>);
  static_assert(std::is_same_v<decltype(RecordWithVectors::default_vector_fixed_length), std::array<int, 3>>);
}

TEST(DefinitionsTests, ArrayFieldsHaveExpectedType) {
  static_assert(std::is_same_v<decltype(RecordWithArrays::default_array),
                               yardl::DynamicNDArray<int32_t>>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::default_array_with_empty_dimension),
                               yardl::DynamicNDArray<int32_t>>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::rank_1_array),
                               yardl::NDArray<int32_t, 1>>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::rank_2_array),
                               yardl::NDArray<int32_t, 2>>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::rank_2_array_with_named_dimensions),
                               yardl::NDArray<int32_t, 2>>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::rank_2_fixed_array),
                               yardl::FixedNDArray<int32_t, 3, 4>>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::rank_2_fixed_array_with_named_dimensions),
                               yardl::FixedNDArray<int32_t, 3, 4>>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::dynamic_array),
                               yardl::DynamicNDArray<int32_t>>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::array_of_vectors),
                               yardl::FixedNDArray<std::array<int, 4>, 5>>);
}

TEST(DefinitionsTests, ArrayFieldsWithSimpleSyntaxHaveExpectedType) {
  static_assert(std::is_same_v<decltype(RecordWithArrays::default_array),
                               decltype(RecordWithArraysSimpleSyntax::default_array)>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::default_array_with_empty_dimension),
                               decltype(RecordWithArraysSimpleSyntax::default_array_with_empty_dimension)>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::rank_1_array),
                               decltype(RecordWithArraysSimpleSyntax::rank_1_array)>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::rank_2_array),
                               decltype(RecordWithArraysSimpleSyntax::rank_2_array)>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::rank_2_array_with_named_dimensions),
                               decltype(RecordWithArraysSimpleSyntax::rank_2_array_with_named_dimensions)>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::rank_2_fixed_array),
                               decltype(RecordWithArraysSimpleSyntax::rank_2_fixed_array)>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::rank_2_fixed_array_with_named_dimensions),
                               decltype(RecordWithArraysSimpleSyntax::rank_2_fixed_array_with_named_dimensions)>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::dynamic_array),
                               decltype(RecordWithArraysSimpleSyntax::dynamic_array)>);
  static_assert(std::is_same_v<decltype(RecordWithArrays::array_of_vectors),
                               decltype(RecordWithArraysSimpleSyntax::array_of_vectors)>);
}

TEST(DefinitionsTests, OptionalsHaveExpectedType) {
  static_assert(std::is_same_v<decltype(RecordWithOptionalFields::optional_int),
                               std::optional<int32_t>>);
  static_assert(std::is_same_v<decltype(RecordWithOptionalFields::optional_int_alternate_syntax),
                               std::optional<int32_t>>);
}

TEST(DefinitionsTests, EnumsHaveExpectedUnderlyingType) {
  static_assert(std::is_same_v<std::underlying_type_t<UInt64Enum>, uint64_t>);
  static_assert(std::is_same_v<std::underlying_type_t<Int64Enum>, int64_t>);
}

}  // namespace
