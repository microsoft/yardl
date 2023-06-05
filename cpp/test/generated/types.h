// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include <array>
#include <complex>
#include <optional>
#include <unordered_map>
#include <variant>
#include <vector>

#include "yardl/yardl.h"

namespace test_model {
struct SmallBenchmarkRecord {
  double a{};
  float b{};
  float c{};

  bool operator==(const SmallBenchmarkRecord& other) const {
    return a == other.a &&
      b == other.b &&
      c == other.c;
  }

  bool operator!=(const SmallBenchmarkRecord& other) const {
    return !(*this == other);
  }
};

struct SimpleEncodingCounters {
  std::optional<uint32_t> e1{};
  std::optional<uint32_t> e2{};
  std::optional<uint32_t> slice{};
  std::optional<uint32_t> repetition{};

  bool operator==(const SimpleEncodingCounters& other) const {
    return e1 == other.e1 &&
      e2 == other.e2 &&
      slice == other.slice &&
      repetition == other.repetition;
  }

  bool operator!=(const SimpleEncodingCounters& other) const {
    return !(*this == other);
  }
};

struct SimpleAcquisition {
  uint64_t flags{};
  test_model::SimpleEncodingCounters idx{};
  yardl::NDArray<std::complex<float>, 2> data{};
  yardl::NDArray<float, 2> trajectory{};

  bool operator==(const SimpleAcquisition& other) const {
    return flags == other.flags &&
      idx == other.idx &&
      data == other.data &&
      trajectory == other.trajectory;
  }

  bool operator!=(const SimpleAcquisition& other) const {
    return !(*this == other);
  }
};

struct SimpleRecord {
  int32_t x{};
  int32_t y{};
  int32_t z{};

  bool operator==(const SimpleRecord& other) const {
    return x == other.x &&
      y == other.y &&
      z == other.z;
  }

  bool operator!=(const SimpleRecord& other) const {
    return !(*this == other);
  }
};

struct RecordWithPrimitives {
  bool bool_field{};
  int8_t int_8_field{};
  uint8_t uint_8_field{};
  int16_t int_16_field{};
  uint16_t uint_16_field{};
  int32_t int_32_field{};
  uint32_t uint_32_field{};
  int64_t int_64_field{};
  uint64_t uint_64_field{};
  yardl::Size size_field{};
  float float_32_field{};
  double float_64_field{};
  std::complex<float> complexfloat_32_field{};
  std::complex<double> complexfloat_64_field{};
  yardl::Date date_field{};
  yardl::Time time_field{};
  yardl::DateTime datetime_field{};

  bool operator==(const RecordWithPrimitives& other) const {
    return bool_field == other.bool_field &&
      int_8_field == other.int_8_field &&
      uint_8_field == other.uint_8_field &&
      int_16_field == other.int_16_field &&
      uint_16_field == other.uint_16_field &&
      int_32_field == other.int_32_field &&
      uint_32_field == other.uint_32_field &&
      int_64_field == other.int_64_field &&
      uint_64_field == other.uint_64_field &&
      size_field == other.size_field &&
      float_32_field == other.float_32_field &&
      float_64_field == other.float_64_field &&
      complexfloat_32_field == other.complexfloat_32_field &&
      complexfloat_64_field == other.complexfloat_64_field &&
      date_field == other.date_field &&
      time_field == other.time_field &&
      datetime_field == other.datetime_field;
  }

  bool operator!=(const RecordWithPrimitives& other) const {
    return !(*this == other);
  }
};

struct RecordWithPrimitiveAliases {
  uint8_t byte_field{};
  int32_t int_field{};
  uint32_t uint_field{};
  int64_t long_field{};
  uint64_t ulong_field{};
  float float_field{};
  double double_field{};
  std::complex<float> complexfloat_field{};
  std::complex<double> complexdouble_field{};

  bool operator==(const RecordWithPrimitiveAliases& other) const {
    return byte_field == other.byte_field &&
      int_field == other.int_field &&
      uint_field == other.uint_field &&
      long_field == other.long_field &&
      ulong_field == other.ulong_field &&
      float_field == other.float_field &&
      double_field == other.double_field &&
      complexfloat_field == other.complexfloat_field &&
      complexdouble_field == other.complexdouble_field;
  }

  bool operator!=(const RecordWithPrimitiveAliases& other) const {
    return !(*this == other);
  }
};

struct TupleWithRecords {
  test_model::SimpleRecord a{};
  test_model::SimpleRecord b{};

  bool operator==(const TupleWithRecords& other) const {
    return a == other.a &&
      b == other.b;
  }

  bool operator!=(const TupleWithRecords& other) const {
    return !(*this == other);
  }
};

struct RecordWithVectors {
  std::vector<int32_t> default_vector{};
  std::array<int32_t, 3> default_vector_fixed_length{};
  std::vector<std::array<int32_t, 2>> vector_of_vectors{};

  bool operator==(const RecordWithVectors& other) const {
    return default_vector == other.default_vector &&
      default_vector_fixed_length == other.default_vector_fixed_length &&
      vector_of_vectors == other.vector_of_vectors;
  }

  bool operator!=(const RecordWithVectors& other) const {
    return !(*this == other);
  }
};

struct RecordWithArrays {
  yardl::DynamicNDArray<int32_t> default_array{};
  yardl::DynamicNDArray<int32_t> default_array_with_empty_dimension{};
  yardl::NDArray<int32_t, 1> rank_1_array{};
  yardl::NDArray<int32_t, 2> rank_2_array{};
  yardl::NDArray<int32_t, 2> rank_2_array_with_named_dimensions{};
  yardl::FixedNDArray<int32_t, 3, 4> rank_2_fixed_array{};
  yardl::FixedNDArray<int32_t, 3, 4> rank_2_fixed_array_with_named_dimensions{};
  yardl::DynamicNDArray<int32_t> dynamic_array{};
  yardl::FixedNDArray<std::array<int32_t, 4>, 5> array_of_vectors{};

  bool operator==(const RecordWithArrays& other) const {
    return default_array == other.default_array &&
      default_array_with_empty_dimension == other.default_array_with_empty_dimension &&
      rank_1_array == other.rank_1_array &&
      rank_2_array == other.rank_2_array &&
      rank_2_array_with_named_dimensions == other.rank_2_array_with_named_dimensions &&
      rank_2_fixed_array == other.rank_2_fixed_array &&
      rank_2_fixed_array_with_named_dimensions == other.rank_2_fixed_array_with_named_dimensions &&
      dynamic_array == other.dynamic_array &&
      array_of_vectors == other.array_of_vectors;
  }

  bool operator!=(const RecordWithArrays& other) const {
    return !(*this == other);
  }
};

struct RecordWithArraysSimpleSyntax {
  yardl::DynamicNDArray<int32_t> default_array{};
  yardl::DynamicNDArray<int32_t> default_array_with_empty_dimension{};
  yardl::NDArray<int32_t, 1> rank_1_array{};
  yardl::NDArray<int32_t, 2> rank_2_array{};
  yardl::NDArray<int32_t, 2> rank_2_array_with_named_dimensions{};
  yardl::FixedNDArray<int32_t, 3, 4> rank_2_fixed_array{};
  yardl::FixedNDArray<int32_t, 3, 4> rank_2_fixed_array_with_named_dimensions{};
  yardl::DynamicNDArray<int32_t> dynamic_array{};
  yardl::FixedNDArray<std::array<int32_t, 4>, 5> array_of_vectors{};

  bool operator==(const RecordWithArraysSimpleSyntax& other) const {
    return default_array == other.default_array &&
      default_array_with_empty_dimension == other.default_array_with_empty_dimension &&
      rank_1_array == other.rank_1_array &&
      rank_2_array == other.rank_2_array &&
      rank_2_array_with_named_dimensions == other.rank_2_array_with_named_dimensions &&
      rank_2_fixed_array == other.rank_2_fixed_array &&
      rank_2_fixed_array_with_named_dimensions == other.rank_2_fixed_array_with_named_dimensions &&
      dynamic_array == other.dynamic_array &&
      array_of_vectors == other.array_of_vectors;
  }

  bool operator!=(const RecordWithArraysSimpleSyntax& other) const {
    return !(*this == other);
  }
};

struct RecordWithOptionalFields {
  std::optional<int32_t> optional_int{};
  std::optional<int32_t> optional_int_alternate_syntax{};

  bool operator==(const RecordWithOptionalFields& other) const {
    return optional_int == other.optional_int &&
      optional_int_alternate_syntax == other.optional_int_alternate_syntax;
  }

  bool operator!=(const RecordWithOptionalFields& other) const {
    return !(*this == other);
  }
};

struct RecordWithVlens {
  std::vector<test_model::SimpleRecord> a{};
  int32_t b{};
  int32_t c{};

  bool operator==(const RecordWithVlens& other) const {
    return a == other.a &&
      b == other.b &&
      c == other.c;
  }

  bool operator!=(const RecordWithVlens& other) const {
    return !(*this == other);
  }
};

struct RecordWithStrings {
  std::string a{};
  std::string b{};

  bool operator==(const RecordWithStrings& other) const {
    return a == other.a &&
      b == other.b;
  }

  bool operator!=(const RecordWithStrings& other) const {
    return !(*this == other);
  }
};

struct RecordWithOptionalVector {
  std::optional<std::vector<int32_t>> optional_vector{};

  bool operator==(const RecordWithOptionalVector& other) const {
    return optional_vector == other.optional_vector;
  }

  bool operator!=(const RecordWithOptionalVector& other) const {
    return !(*this == other);
  }
};

struct RecordWithFixedVectors {
  std::array<int32_t, 5> fixed_int_vector{};
  std::array<test_model::SimpleRecord, 3> fixed_simple_record_vector{};
  std::array<test_model::RecordWithVlens, 2> fixed_record_with_vlens_vector{};

  bool operator==(const RecordWithFixedVectors& other) const {
    return fixed_int_vector == other.fixed_int_vector &&
      fixed_simple_record_vector == other.fixed_simple_record_vector &&
      fixed_record_with_vlens_vector == other.fixed_record_with_vlens_vector;
  }

  bool operator!=(const RecordWithFixedVectors& other) const {
    return !(*this == other);
  }
};

struct RecordWithFixedArrays {
  yardl::FixedNDArray<int32_t, 2, 3> ints{};
  yardl::FixedNDArray<test_model::SimpleRecord, 3, 2> fixed_simple_record_array{};
  yardl::FixedNDArray<test_model::RecordWithVlens, 2, 2> fixed_record_with_vlens_array{};

  bool operator==(const RecordWithFixedArrays& other) const {
    return ints == other.ints &&
      fixed_simple_record_array == other.fixed_simple_record_array &&
      fixed_record_with_vlens_array == other.fixed_record_with_vlens_array;
  }

  bool operator!=(const RecordWithFixedArrays& other) const {
    return !(*this == other);
  }
};

struct RecordWithNDArrays {
  yardl::NDArray<int32_t, 2> ints{};
  yardl::NDArray<test_model::SimpleRecord, 2> fixed_simple_record_array{};
  yardl::NDArray<test_model::RecordWithVlens, 2> fixed_record_with_vlens_array{};

  bool operator==(const RecordWithNDArrays& other) const {
    return ints == other.ints &&
      fixed_simple_record_array == other.fixed_simple_record_array &&
      fixed_record_with_vlens_array == other.fixed_record_with_vlens_array;
  }

  bool operator!=(const RecordWithNDArrays& other) const {
    return !(*this == other);
  }
};

struct RecordWithNDArraysSingleDimension {
  yardl::NDArray<int32_t, 1> ints{};
  yardl::NDArray<test_model::SimpleRecord, 1> fixed_simple_record_array{};
  yardl::NDArray<test_model::RecordWithVlens, 1> fixed_record_with_vlens_array{};

  bool operator==(const RecordWithNDArraysSingleDimension& other) const {
    return ints == other.ints &&
      fixed_simple_record_array == other.fixed_simple_record_array &&
      fixed_record_with_vlens_array == other.fixed_record_with_vlens_array;
  }

  bool operator!=(const RecordWithNDArraysSingleDimension& other) const {
    return !(*this == other);
  }
};

struct RecordWithDynamicNDArrays {
  yardl::DynamicNDArray<int32_t> ints{};
  yardl::DynamicNDArray<test_model::SimpleRecord> fixed_simple_record_array{};
  yardl::DynamicNDArray<test_model::RecordWithVlens> fixed_record_with_vlens_array{};

  bool operator==(const RecordWithDynamicNDArrays& other) const {
    return ints == other.ints &&
      fixed_simple_record_array == other.fixed_simple_record_array &&
      fixed_record_with_vlens_array == other.fixed_record_with_vlens_array;
  }

  bool operator!=(const RecordWithDynamicNDArrays& other) const {
    return !(*this == other);
  }
};

using NamedFixedNDArray = yardl::FixedNDArray<int32_t, 2, 4>;

using NamedNDArray = yardl::NDArray<int32_t, 2>;

template <typename K, typename V>
using AliasedMap = std::unordered_map<K, V>;

struct RecordWithUnions {
  std::variant<std::monostate, int32_t, std::string> null_or_int_or_string{};

  bool operator==(const RecordWithUnions& other) const {
    return null_or_int_or_string == other.null_or_int_or_string;
  }

  bool operator!=(const RecordWithUnions& other) const {
    return !(*this == other);
  }
};

enum class Fruits {
  kApple = 0,
  kBanana = 1,
  kPear = 2,
};

enum class UInt64Enum : uint64_t {
  kA = 9223372036854775808ULL,
};

enum class Int64Enum : int64_t {
  kB = -4611686018427387904LL,
};

enum class SizeBasedEnum : yardl::Size {
  kA = 0ULL,
  kB = 1ULL,
  kC = 2ULL,
};

struct DaysOfWeek : yardl::BaseFlags<int32_t, DaysOfWeek> {
  using BaseFlags::BaseFlags;
  static const DaysOfWeek kMonday;
  static const DaysOfWeek kTuesday;
  static const DaysOfWeek kWednesday;
  static const DaysOfWeek kThursday;
  static const DaysOfWeek kFriday;
  static const DaysOfWeek kSaturday;
  static const DaysOfWeek kSunday;
};

struct TextFormat : yardl::BaseFlags<uint64_t, TextFormat> {
  using BaseFlags::BaseFlags;
  static const TextFormat kRegular;
  static const TextFormat kBold;
  static const TextFormat kItalic;
  static const TextFormat kUnderline;
  static const TextFormat kStrikethrough;
};

template <typename T>
using Image = yardl::NDArray<T, 2>;

template <typename T1, typename T2>
struct GenericRecord {
  T1 scalar_1{};
  T2 scalar_2{};
  std::vector<T1> vector_1{};
  test_model::Image<T2> image_2{};

  bool operator==(const GenericRecord& other) const {
    return scalar_1 == other.scalar_1 &&
      scalar_2 == other.scalar_2 &&
      vector_1 == other.vector_1 &&
      image_2 == other.image_2;
  }

  bool operator!=(const GenericRecord& other) const {
    return !(*this == other);
  }
};

template <typename T1, typename T2>
struct MyTuple {
  T1 v1{};
  T2 v2{};

  bool operator==(const MyTuple& other) const {
    return v1 == other.v1 &&
      v2 == other.v2;
  }

  bool operator!=(const MyTuple& other) const {
    return !(*this == other);
  }
};

using AliasedString = std::string;

using AliasedEnum = test_model::Fruits;

template <typename T1, typename T2>
using AliasedOpenGeneric = test_model::MyTuple<T1, T2>;

using AliasedClosedGeneric = test_model::MyTuple<test_model::AliasedString, test_model::AliasedEnum>;

using AliasedOptional = std::optional<int32_t>;

template <typename T>
using AliasedGenericOptional = std::optional<T>;

template <typename T1, typename T2>
using AliasedGenericUnion2 = std::variant<T1, T2>;

template <typename T>
using AliasedGenericVector = std::vector<T>;

template <typename T>
using AliasedGenericFixedVector = std::array<T, 3>;

using AliasedIntOrSimpleRecord = std::variant<int32_t, test_model::SimpleRecord>;

using AliasedNullableIntSimpleRecord = std::variant<std::monostate, int32_t, test_model::SimpleRecord>;

template <typename T0, typename T1>
struct GenericRecordWithComputedFields {
  std::variant<T0, T1> f1{};

  uint8_t TypeIndex() const {
    return std::visit(
      [&](auto&& __case_arg__) -> uint8_t {
        if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, T0>) {
          return 0;
        }
        if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, T1>) {
          return 1;
        }
      },
      f1);
  }

  bool operator==(const GenericRecordWithComputedFields& other) const {
    return f1 == other.f1;
  }

  bool operator!=(const GenericRecordWithComputedFields& other) const {
    return !(*this == other);
  }
};

struct RecordWithComputedFields {
  yardl::NDArray<int32_t, 2> array_field{};
  yardl::NDArray<int32_t, 2> array_field_map_dimensions{};
  yardl::DynamicNDArray<int32_t> dynamic_array_field{};
  yardl::FixedNDArray<int32_t, 3, 4> fixed_array_field{};
  int32_t int_field{};
  std::string string_field{};
  test_model::MyTuple<int32_t, int32_t> tuple_field{};
  std::vector<int32_t> vector_field{};
  std::vector<std::vector<int32_t>> vector_of_vectors_field{};
  std::array<int32_t, 3> fixed_vector_field{};
  std::optional<test_model::NamedNDArray> optional_named_array{};
  std::variant<int32_t, float> int_float_union{};
  std::variant<std::monostate, int32_t, float> nullable_int_float_union{};
  std::variant<int32_t, test_model::GenericRecordWithComputedFields<std::string, float>> union_with_nested_generic_union{};
  std::unordered_map<std::string, std::string> map_field{};

  uint8_t IntLiteral() const {
    return 42;
  }

  int64_t LargeNegativeInt64Literal() const {
    return -4611686018427387904LL;
  }

  uint64_t LargeUInt64Literal() const {
    return 9223372036854775808ULL;
  }

  std::string StringLiteral() const {
    return "hello";
  }

  std::string StringLiteral2() const {
    return "hello";
  }

  std::string StringLiteral3() const {
    return "hello";
  }

  std::string StringLiteral4() const {
    return "hello";
  }

  int32_t const& AccessOtherComputedField() const {
    return int_field;
  }

  int32_t& AccessOtherComputedField() {
    return const_cast<int32_t&>(std::as_const(*this).AccessOtherComputedField());
  }

  int32_t const& AccessIntField() const {
    return int_field;
  }

  int32_t& AccessIntField() {
    return const_cast<int32_t&>(std::as_const(*this).AccessIntField());
  }

  std::string const& AccessStringField() const {
    return string_field;
  }

  std::string& AccessStringField() {
    return const_cast<std::string&>(std::as_const(*this).AccessStringField());
  }

  test_model::MyTuple<int32_t, int32_t> const& AccessTupleField() const {
    return tuple_field;
  }

  test_model::MyTuple<int32_t, int32_t>& AccessTupleField() {
    return const_cast<test_model::MyTuple<int32_t, int32_t>&>(std::as_const(*this).AccessTupleField());
  }

  int32_t const& AccessNestedTupleField() const {
    return tuple_field.v2;
  }

  int32_t& AccessNestedTupleField() {
    return const_cast<int32_t&>(std::as_const(*this).AccessNestedTupleField());
  }

  yardl::NDArray<int32_t, 2> const& AccessArrayField() const {
    return array_field;
  }

  yardl::NDArray<int32_t, 2>& AccessArrayField() {
    return const_cast<yardl::NDArray<int32_t, 2>&>(std::as_const(*this).AccessArrayField());
  }

  int32_t const& AccessArrayFieldElement() const {
    return array_field.at(0, 1);
  }

  int32_t& AccessArrayFieldElement() {
    return const_cast<int32_t&>(std::as_const(*this).AccessArrayFieldElement());
  }

  int32_t const& AccessArrayFieldElementByName() const {
    return array_field.at(0, 1);
  }

  int32_t& AccessArrayFieldElementByName() {
    return const_cast<int32_t&>(std::as_const(*this).AccessArrayFieldElementByName());
  }

  std::vector<int32_t> const& AccessVectorField() const {
    return vector_field;
  }

  std::vector<int32_t>& AccessVectorField() {
    return const_cast<std::vector<int32_t>&>(std::as_const(*this).AccessVectorField());
  }

  int32_t const& AccessVectorFieldElement() const {
    return vector_field.at(1);
  }

  int32_t& AccessVectorFieldElement() {
    return const_cast<int32_t&>(std::as_const(*this).AccessVectorFieldElement());
  }

  int32_t const& AccessVectorOfVectorsField() const {
    return vector_of_vectors_field.at(1).at(2);
  }

  int32_t& AccessVectorOfVectorsField() {
    return const_cast<int32_t&>(std::as_const(*this).AccessVectorOfVectorsField());
  }

  yardl::Size ArraySize() const {
    return array_field.size();
  }

  yardl::Size ArrayXSize() const {
    return array_field.shape(0);
  }

  yardl::Size ArrayYSize() const {
    return array_field.shape(1);
  }

  yardl::Size Array0Size() const {
    return array_field.shape(0);
  }

  yardl::Size Array1Size() const {
    return array_field.shape(1);
  }

  yardl::Size ArraySizeFromIntField() const {
    return array_field.shape(int_field);
  }

  yardl::Size ArraySizeFromStringField() const {
    return array_field.shape(([](std::string dim_name) {
      if (dim_name == "x") return 0;
      if (dim_name == "y") return 1;
      throw std::invalid_argument("Unknown dimension name: " + dim_name);
    })(string_field));
  }

  yardl::Size ArraySizeFromNestedIntField() const {
    return array_field.shape(tuple_field.v1);
  }

  yardl::Size ArrayFieldMapDimensionsXSize() const {
    return array_field_map_dimensions.shape(0);
  }

  yardl::Size FixedArraySize() const {
    return 12ULL;
  }

  yardl::Size FixedArrayXSize() const {
    return 3ULL;
  }

  yardl::Size FixedArray0Size() const {
    return 3ULL;
  }

  yardl::Size VectorSize() const {
    return vector_field.size();
  }

  yardl::Size FixedVectorSize() const {
    return 3ULL;
  }

  yardl::Size ArrayDimensionXIndex() const {
    return 0ULL;
  }

  yardl::Size ArrayDimensionYIndex() const {
    return 1ULL;
  }

  yardl::Size ArrayDimensionIndexFromStringField() const {
    return ([](std::string dim_name) {
      if (dim_name == "x") return 0;
      if (dim_name == "y") return 1;
      throw std::invalid_argument("Unknown dimension name: " + dim_name);
    })(string_field);
  }

  yardl::Size ArrayDimensionCount() const {
    return 2ULL;
  }

  yardl::Size DynamicArrayDimensionCount() const {
    return dynamic_array_field.dimension();
  }

  std::unordered_map<std::string, std::string> const& AccessMap() const {
    return map_field;
  }

  std::unordered_map<std::string, std::string>& AccessMap() {
    return const_cast<std::unordered_map<std::string, std::string>&>(std::as_const(*this).AccessMap());
  }

  yardl::Size MapSize() const {
    return map_field.size();
  }

  std::string const& AccessMapEntry() const {
    return map_field.at("hello");
  }

  std::string& AccessMapEntry() {
    return const_cast<std::string&>(std::as_const(*this).AccessMapEntry());
  }

  std::string StringComputedField() const {
    return "hello";
  }

  std::string const& AccessMapEntryWithComputedField() const {
    return map_field.at(StringComputedField());
  }

  std::string& AccessMapEntryWithComputedField() {
    return const_cast<std::string&>(std::as_const(*this).AccessMapEntryWithComputedField());
  }

  std::string const& AccessMapEntryWithComputedFieldNested() const {
    return map_field.at(map_field.at(StringComputedField()));
  }

  std::string& AccessMapEntryWithComputedFieldNested() {
    return const_cast<std::string&>(std::as_const(*this).AccessMapEntryWithComputedFieldNested());
  }

  std::string const& AccessMissingMapEntry() const {
    return map_field.at("missing");
  }

  std::string& AccessMissingMapEntry() {
    return const_cast<std::string&>(std::as_const(*this).AccessMissingMapEntry());
  }

  yardl::Size OptionalNamedArrayLength() const {
    return [](auto&& __case_arg__) -> yardl::Size {
      if (__case_arg__.has_value()) {
        test_model::NamedNDArray const& arr = __case_arg__.value();
        return arr.size();
      }
      return 0ULL;
    }(optional_named_array);
  }

  yardl::Size OptionalNamedArrayLengthWithDiscard() const {
    return [](auto&& __case_arg__) -> yardl::Size {
      if (__case_arg__.has_value()) {
        test_model::NamedNDArray const& arr = __case_arg__.value();
        return arr.size();
      }
      return 0ULL;
    }(optional_named_array);
  }

  float IntFloatUnionAsFloat() const {
    return std::visit(
      [&](auto&& __case_arg__) -> float {
        if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, int32_t>) {
          int32_t const& i = __case_arg__;
          return static_cast<float>(i);
        }
        if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, float>) {
          float const& f = __case_arg__;
          return f;
        }
      },
      int_float_union);
  }

  std::string NullableIntFloatUnionString() const {
    return std::visit(
      [&](auto&& __case_arg__) -> std::string {
        if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, std::monostate>) {
          return "null";
        }
        if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, int32_t>) {
          return "int";
        }
        return "float";
      },
      nullable_int_float_union);
  }

  int16_t NestedSwitch() const {
    return std::visit(
      [&](auto&& __case_arg__) -> int16_t {
        if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, int32_t>) {
          return -1;
        }
        if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, test_model::GenericRecordWithComputedFields<std::string, float>>) {
          test_model::GenericRecordWithComputedFields<std::string, float> const& rec = __case_arg__;
          return static_cast<int16_t>(std::visit(
            [&](auto&& __case_arg__) -> uint8_t {
              if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, float>) {
                return 20;
              }
              if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, std::string>) {
                return 10;
              }
            },
            rec.f1));
        }
      },
      union_with_nested_generic_union);
  }

  int16_t UseNestedComputedField() const {
    return std::visit(
      [&](auto&& __case_arg__) -> int16_t {
        if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, int32_t>) {
          return -1;
        }
        if constexpr (std::is_same_v<std::decay_t<decltype(__case_arg__)>, test_model::GenericRecordWithComputedFields<std::string, float>>) {
          test_model::GenericRecordWithComputedFields<std::string, float> const& rec = __case_arg__;
          return static_cast<int16_t>(rec.TypeIndex());
        }
      },
      union_with_nested_generic_union);
  }

  bool operator==(const RecordWithComputedFields& other) const {
    return array_field == other.array_field &&
      array_field_map_dimensions == other.array_field_map_dimensions &&
      dynamic_array_field == other.dynamic_array_field &&
      fixed_array_field == other.fixed_array_field &&
      int_field == other.int_field &&
      string_field == other.string_field &&
      tuple_field == other.tuple_field &&
      vector_field == other.vector_field &&
      vector_of_vectors_field == other.vector_of_vectors_field &&
      fixed_vector_field == other.fixed_vector_field &&
      optional_named_array == other.optional_named_array &&
      int_float_union == other.int_float_union &&
      nullable_int_float_union == other.nullable_int_float_union &&
      union_with_nested_generic_union == other.union_with_nested_generic_union &&
      map_field == other.map_field;
  }

  bool operator!=(const RecordWithComputedFields& other) const {
    return !(*this == other);
  }
};

template <typename INT16_MAX_Type>
using ArrayWithKeywordDimensionNames = yardl::NDArray<INT16_MAX_Type, 2>;

enum class EnumWithKeywordSymbols {
  kTry = 2,
  kCatch = 1,
};

// BEGIN delibrately using C++ keywords and macros as identitiers
struct RecordWithKeywordFields {
  std::string int_field{};
  test_model::ArrayWithKeywordDimensionNames<int32_t> sizeof_field{};
  test_model::EnumWithKeywordSymbols if_field{};

  std::string const& Float() const {
    return int_field;
  }

  std::string& Float() {
    return const_cast<std::string&>(std::as_const(*this).Float());
  }

  std::string const& Double() const {
    return Float();
  }

  std::string& Double() {
    return const_cast<std::string&>(std::as_const(*this).Double());
  }

  int32_t const& Return() const {
    return sizeof_field.at(1, 2);
  }

  int32_t& Return() {
    return const_cast<int32_t&>(std::as_const(*this).Return());
  }

  bool operator==(const RecordWithKeywordFields& other) const {
    return int_field == other.int_field &&
      sizeof_field == other.sizeof_field &&
      if_field == other.if_field;
  }

  bool operator!=(const RecordWithKeywordFields& other) const {
    return !(*this == other);
  }
};

} // namespace test_model
