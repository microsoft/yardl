#pragma once

#include <utility>

#ifndef XTENSOR_VERSION_MAJOR
// need to include this for version info
#  include <xtensor/xtensor.hpp>
#endif

#if XTENSOR_VERSION_MAJOR == 0 && XTENSOR_VERSION_MINOR < 26
#include <xtensor/xarray.hpp>
#include <xtensor/xfixed.hpp>
#include <xtensor/xio.hpp>
#include <xtensor/xtensor.hpp>
#else
#include <xtensor/containers/xarray.hpp>
#include <xtensor/containers/xfixed.hpp>
#include <xtensor/io/xio.hpp>
#include <xtensor/containers/xtensor.hpp>
#endif

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

/**** FixedNDArray Implementation ****/
template <typename T, size_t... Dims>
constexpr size_t size(FixedNDArray<T, Dims...> const& arr) { return arr.size(); }

template <typename T, size_t... Dims>
constexpr size_t dimension(FixedNDArray<T, Dims...> const& arr) { return arr.dimension(); }

template <typename T, size_t... Dims>
constexpr std::array<size_t, sizeof...(Dims)> shape(FixedNDArray<T, Dims...> const& arr) { return arr.shape(); }

template <typename T, size_t... Dims>
constexpr size_t shape(FixedNDArray<T, Dims...> const& arr, size_t dim) { return arr.shape(dim); }

template <typename T, size_t... Dims>
constexpr T* dataptr(FixedNDArray<T, Dims...>& arr) { return arr.data(); }
template <typename T, size_t... Dims>
constexpr T const* dataptr(FixedNDArray<T, Dims...> const& arr) { return arr.data(); }

template <typename T, size_t... Dims, class... Args>
constexpr T const& at(FixedNDArray<T, Dims...> const& arr, Args... idx) { return arr.at(idx...); }

/**** NDArray Implementation ****/
template <typename T, size_t N>
size_t size(NDArray<T, N> const& arr) { return arr.size(); }

template <typename T, size_t N>
size_t dimension(NDArray<T, N> const& arr) { return arr.dimension(); }

template <typename T, size_t N>
std::array<size_t, N> shape(NDArray<T, N> const& arr) { return arr.shape(); }

template <typename T, size_t N>
size_t shape(NDArray<T, N> const& arr, size_t dim) { return arr.shape(dim); }

template <typename T, size_t N>
void resize(NDArray<T, N>& arr, std::array<size_t, N> const& shape) { arr.resize(shape, true); }

template <typename T, size_t N>
T* dataptr(NDArray<T, N>& arr) { return arr.data(); }
template <typename T, size_t N>
T const* dataptr(NDArray<T, N> const& arr) { return arr.data(); }

template <typename T, size_t N, class... Args>
T const& at(NDArray<T, N> const& arr, Args... idx) { return arr.at(idx...); }

/**** DynamicNDArray Implementation ****/
template <typename T>
size_t size(DynamicNDArray<T> const& arr) { return arr.size(); }

template <typename T>
size_t dimension(DynamicNDArray<T> const& arr) { return arr.dimension(); }

template <typename T>
std::vector<size_t> shape(DynamicNDArray<T> const& arr) {
  // Xtensor xarray.shape() is an xt::svector, not a std::vector
  auto shape = arr.shape();
  std::vector<size_t> vshape;
  std::copy(shape.begin(), shape.end(), std::back_inserter(vshape));
  return vshape;
}

template <typename T>
size_t shape(DynamicNDArray<T> const& arr, size_t dim) { return arr.shape(dim); }

template <typename T>
void resize(DynamicNDArray<T>& arr, std::vector<size_t> const& shape) { arr.resize(shape, true); }

template <typename T>
T* dataptr(DynamicNDArray<T>& arr) { return arr.data(); }
template <typename T>
T const* dataptr(DynamicNDArray<T> const& arr) { return arr.data(); }

template <typename T, class... Args>
T const& at(DynamicNDArray<T> const& arr, Args... idx) { return arr.at(idx...); }

}  // namespace yardl
