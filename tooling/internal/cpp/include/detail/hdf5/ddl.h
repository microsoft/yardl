// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <array>
#include <complex>
#include <cstring>
#include <memory>
#include <optional>
#include <string>
#include <tuple>
#include <type_traits>
#include <utility>
#include <vector>

#include <H5Cpp.h>

#include "../../yardl.h"
#include "inner_types.h"

// Helper functions for defining HDF5 datatypes.

namespace yardl::hdf5 {

/**
 * @brief Returns the HDF5 type for yardl::Size.
 */
static inline H5::PredType const& SizeTypeDdl() {
  static_assert(sizeof(hsize_t) == sizeof(size_t));
  static_assert(std::is_signed_v<hsize_t> == std::is_signed_v<size_t>);
  return H5::PredType::NATIVE_HSIZE;
}

/**
 * @brief Creates an HDF5 type for dates. These are stored as an integer number
 * of days since the epoch.
 */
static inline H5::DataType DateTypeDdl() {
#if __cplusplus < 202002L
  static_assert(sizeof(yardl::Date) == sizeof(int32_t));
  static_assert(std::is_same_v<yardl::Date::rep, int32_t>);
  return H5::PredType::NATIVE_INT32;
#else
  static_assert(sizeof(yardl::Date) == sizeof(int64_t));
  static_assert(std::is_same_v<yardl::Date::rep, int64_t>);
  return H5::PredType::NATIVE_INT64;
#endif
}

/**
 * @brief Creates an HDF5 type for times. These are stored as an int32 number
 * of nanoseconds since midnight.
 */
static inline H5::DataType TimeTypeDdl() {
  static_assert(sizeof(yardl::Time) == sizeof(int64_t));
  static_assert(std::is_same_v<yardl::Time::rep, int64_t>);
  return H5::PredType::NATIVE_INT64;
}

/**
 * @brief Creates an HDF5 type for datetimes. These are stored as an int64
 * number of nanoseconds since the epoch, ignoring leap seconds
 */
static inline H5::DataType DateTimeTypeDdl() {
  static_assert(sizeof(yardl::DateTime) == sizeof(int64_t));
  static_assert(std::is_same_v<yardl::DateTime::rep, int64_t>);
  return H5::PredType::NATIVE_INT64;
}

/**
 * @brief Creates an HDF5 optional type.
 */
template <typename TInner, typename TOuter>
H5::CompType OptionalTypeDdl(H5::DataType const& value_type) {
  using InnerOptionalType = InnerOptional<TInner, TOuter>;
  H5::CompType type(sizeof(InnerOptionalType));
  type.insertMember("has_value", HOFFSET(InnerOptionalType, has_value), H5::PredType::NATIVE_HBOOL);
  type.insertMember("value", HOFFSET(InnerOptionalType, value), value_type);
  return type;
}

/**
 * @brief Creates an HDF5 v-len data type.
 */
static inline H5::VarLenType InnerVlenDdl(H5::DataType const& element_type) {
  return H5::VarLenType(element_type);
}

/**
 * @brief Creates a compund datatype for std::complex
 */
template <typename T>
H5::CompType ComplexTypeDdl() {
  H5::DataType inner_type;
  if constexpr (std::is_same_v<T, float>) {
    inner_type = H5::PredType::NATIVE_FLOAT;
  } else {
    static_assert(std::is_same_v<T, double>, "Unsupported type parameter");
    inner_type = H5::PredType::NATIVE_DOUBLE;
  }

  H5::CompType type(sizeof(std::complex<T>));
  type.insertMember("real", 0, inner_type);
  type.insertMember("imaginary", sizeof(T), inner_type);
  return type;
}

/**
 * @brief Creates a variable-length string datatype.
 */
static inline H5::StrType InnerVlenStringDdl() {
  static_assert(sizeof(InnerVlenString) == sizeof(char*));
  H5::StrType type(0, H5T_VARIABLE);
  type.setCset(H5T_CSET_UTF8);
  return type;
}

/**
 * @brief Creates a datatype for a fixed-length vector.
 */
static inline H5::ArrayType FixedVectorDdl(H5::DataType const& element_type, hsize_t length) {
  hsize_t size = length;
  return H5::ArrayType(element_type, 1, &size);
}

/**
 * @brief Creates a datatype for a fixed-size NDArray
 */
static inline H5::ArrayType FixedNDArrayDdl(H5::DataType const& element_type,
                                            std::initializer_list<hsize_t> dimensions) {
  return H5::ArrayType(element_type, static_cast<int>(dimensions.size()), std::data(dimensions));
}

/**
 * @brief Creates a datatype for an NDArray with a known number of dimensions.
 */
template <typename TInner, typename TOuter, size_t N>
H5::CompType NDArrayDdl(H5::DataType const& element_type) {
  using ArrayType = InnerNdArray<TInner, TOuter, N>;
  H5::CompType compType(sizeof(ArrayType));
  hsize_t dims = N;
  compType.insertMember("dimensions", HOFFSET(ArrayType, dimensions_),
                        H5::ArrayType(H5::PredType::NATIVE_UINT64, 1, &dims));
  compType.insertMember("data", HOFFSET(ArrayType, data_), H5::VarLenType(element_type));

  assert(H5::PredType::NATIVE_UINT64.getSize() == sizeof(size_t));
  return compType;
}

/**
 * @brief Creates a datatype for DynamicNDArray (unknown number of dimensions)
 */
template <typename TInner, typename TOuter>
H5::CompType DynamicNDArrayDdl(H5::DataType const& element_type) {
  using ArrayType = InnerDynamicNdArray<TInner, TOuter>;
  H5::CompType compType(sizeof(ArrayType));
  compType.insertMember("dimensions", HOFFSET(ArrayType, dimensions_),
                        H5::VarLenType(H5::PredType::NATIVE_UINT64));
  compType.insertMember("data", HOFFSET(ArrayType, data_), H5::VarLenType(element_type));
  return compType;
}

template <typename... Labels>
H5::EnumType UnionTypeEnumDdl(bool nullable, Labels const&... labels) {
  H5::EnumType type_enum(H5::PredType::NATIVE_INT8);
  int8_t type_value = -1;
  if (nullable) {
    type_enum.insert("null", &type_value);
  }

  ((type_value++, type_enum.insert(labels, &type_value)), ...);

  return type_enum;
}

struct IndexEntry {
  int8_t type_;
  uint64_t offset_;
};

static inline H5::CompType UnionIndexDatasetElementTypeDdl(H5::EnumType type_enum) {
  H5::CompType element_type(sizeof(IndexEntry));
  element_type.insertMember("type", HOFFSET(IndexEntry, type_), type_enum);
  element_type.insertMember("offset", HOFFSET(IndexEntry, offset_), H5::PredType::NATIVE_UINT64);
  return element_type;
}

}  // namespace yardl::hdf5
