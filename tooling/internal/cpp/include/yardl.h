// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <chrono>

#if __cplusplus < 202002L
// This functionality is part of the standard library as of C++20
#include <date/date.h>
#endif

#include <xtensor/xarray.hpp>
#include <xtensor/xfixed.hpp>
#include <xtensor/xtensor.hpp>

namespace yardl {

/**
 * @brief A multidimensional array where all dimension sizes
 * are known at compile-time.
 *
 * @tparam T the element type
 * @tparam Dims the array dimensions
 */
template <typename T, size_t... Dims>
using FixedNDArray = xt::xtensor_fixed<T, xt::xshape<Dims...>,
                                       xt::layout_type::row_major, false>;

/**
 * @brief  A multidimensional array where the number of dimensions
 * is known at compile-time
 *
 * @tparam T the element type
 * @tparam N the number of dimensions
 */
template <typename T, size_t N>
using NDArray = xt::xtensor<T, N, xt::layout_type::row_major>;

/**
 * @brief  A multidimensional array where the number of dimensions
 * is not known at compile-time
 *
 * @tparam T the element type
 */
template <typename T>
using DynamicNDArray = xt::xarray<T, xt::layout_type::row_major>;

#if __cplusplus < 202002L
/**
 * @brief Represents a date as a number of days since the epoch.
 */
using Date = date::local_days;
#else
/**
 * @brief Represents a date as a number of days since the epoch.
 */
using Date = std::chrono::local_days;
#endif

/**
 * @brief Represents a time of day as the number of nanoseconds since midnight.
 */
using Time = std::chrono::duration<int64_t, std::nano>;

/**
 * @brief Represents a datetime as the number of nanoseconds since the epoch.
 */
using DateTime = std::chrono::time_point<std::chrono::system_clock,
                                         std::chrono::duration<int64_t, std::nano>>;

/**
 * @brief The same as size_t when it is 64 bits, otherwise uint64_t.
 */
using Size = std::conditional_t<sizeof(size_t) == sizeof(uint64_t), size_t, uint64_t>;

template <class E, class Enabler = void>
struct is_flags_enum_t
    : public std::false_type {};

template <typename _Tp>
inline constexpr bool is_flags_enum_v = is_flags_enum_t<_Tp>::value;

template <typename TEnum, std::enable_if_t<std::is_enum_v<TEnum> && is_flags_enum_v<TEnum>, bool> = true>
constexpr bool HasFlag(TEnum value, TEnum flag) {
  using integer_type = std::underlying_type_t<TEnum>;
  return (static_cast<integer_type>(value) & static_cast<integer_type>(flag)) == static_cast<integer_type>(flag);
}

}  // namespace yardl

// operators that only apply to flags enums
template <typename TEnum, std::enable_if_t<std::is_enum_v<TEnum> && yardl::is_flags_enum_v<TEnum>, bool> = true>
constexpr TEnum operator|(TEnum lhs, TEnum rhs) {
  using integer_type = std::underlying_type_t<TEnum>;
  return static_cast<TEnum>(static_cast<integer_type>(lhs) | static_cast<integer_type>(rhs));
}

template <typename TEnum, std::enable_if_t<std::is_enum_v<TEnum> && yardl::is_flags_enum_v<TEnum>, bool> = true>
constexpr TEnum operator&(TEnum lhs, TEnum rhs) {
  using integer_type = std::underlying_type_t<TEnum>;
  return static_cast<TEnum>(static_cast<integer_type>(lhs) & static_cast<integer_type>(rhs));
}

template <typename TEnum, std::enable_if_t<std::is_enum_v<TEnum> && yardl::is_flags_enum_v<TEnum>, bool> = true>
constexpr TEnum operator^(TEnum lhs, TEnum rhs) {
  using integer_type = std::underlying_type_t<TEnum>;
  return static_cast<TEnum>(static_cast<integer_type>(lhs) ^ static_cast<integer_type>(rhs));
}

template <typename TEnum, std::enable_if_t<std::is_enum_v<TEnum> && yardl::is_flags_enum_v<TEnum>, bool> = true>
constexpr void operator|=(TEnum& lhs, TEnum rhs) {
  lhs = lhs | rhs;
}
template <typename TEnum, std::enable_if_t<std::is_enum_v<TEnum> && yardl::is_flags_enum_v<TEnum>, bool> = true>
constexpr void operator&=(TEnum& lhs, TEnum rhs) {
  lhs = lhs & rhs;
}
template <typename TEnum, std::enable_if_t<std::is_enum_v<TEnum> && yardl::is_flags_enum_v<TEnum>, bool> = true>
constexpr void operator^=(TEnum& lhs, TEnum rhs) {
  lhs = lhs ^ rhs;
}

template <typename TEnum, std::enable_if_t<std::is_enum_v<TEnum> && yardl::is_flags_enum_v<TEnum>, bool> = true>
constexpr TEnum operator~(TEnum value) {
  using integer_type = std::underlying_type_t<TEnum>;
  return static_cast<TEnum>(~static_cast<integer_type>(value));
}
