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
using NDArray = xt::xtensor<T, N>;

/**
 * @brief  A multidimensional array where the number of dimensions
 * is not known at compile-time
 *
 * @tparam T the element type
 */
template <typename T>
using DynamicNDArray = xt::xarray<T>;

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

}  // namespace yardl
