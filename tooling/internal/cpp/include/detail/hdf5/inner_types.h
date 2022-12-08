// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <array>
#include <complex>
#include <cstring>
#include <memory>
#include <optional>
#include <string>
#include <type_traits>
#include <utility>
#include <vector>

#include <H5Cpp.h>

#include "../../yardl.h"

// HDF5 cannot handle c++ containers like std::vector, so we make use of the
// templates in this file to represent data in a form that is compatible with HDF5.
// If the normal "outer" type is trivially copyable and does not contain
// unsupported constructs, is does not need an inner type.
// The types below generally have a constructor that takes a corresponding
// "outer" type, and override the conversion opertator back to that type.
namespace yardl::hdf5 {

static inline void* MallocOrThrow(size_t size) {
  auto res = std::malloc(size);
  if (!res) {
    throw std::bad_alloc{};
  }

  return res;
}

/**
 * @brief An HDF5-compatible representation of std::optional. T
 * std::optional does not document the inner layout so we we can't use it directly.
 */
template <typename TInner, typename TOuter>
struct InnerOptional {
  InnerOptional() {
  }

  InnerOptional(std::optional<TOuter> const& o) {
    if (o.has_value()) {
      has_value = true;
      if constexpr (std::is_same_v<TInner, TOuter>) {
        value = o.value();
      } else {
        new (&value) TInner(o.value());
      }
    }
  }

  ~InnerOptional() {
    if (has_value) {
      value.~TInner();
    }
  }

  explicit operator std::optional<TOuter>() const {
    if (has_value) {
      return static_cast<TOuter>(value);
    }

    return std::nullopt;
  }

  union {
    char empty[sizeof(TInner)]{};
    TInner value;
  };

  bool has_value = false;
};

/**
 * @brief An HDF5-compatible representation of variable-length
 * std::vector and NDArray<TOuter, 1>
 */
template <typename TInner, typename TOuter>
struct InnerVlen : public hvl_t {
  InnerVlen() : hvl_t{0, nullptr} {
  }

  InnerVlen(std::vector<TOuter> const& v)
      : hvl_t{v.size(), MallocOrThrow(v.size() * sizeof(TInner))} {
    if constexpr (std::is_same_v<TInner, TOuter>) {
      // TODO: we could avoid this copy by having separate read/write types
      std::memcpy(p, const_cast<TInner*>(v.data()), len * sizeof(TInner));
    } else {
      for (size_t i = 0; i < len; i++) {
        auto dest = (TInner*)p + i;
        new (dest) TInner(v[i]);
      }
    }
  }

  InnerVlen(NDArray<TOuter, 1> const& o)
      : hvl_t{o.size(), MallocOrThrow(o.size() * sizeof(TInner))} {
    if constexpr (std::is_same_v<TInner, TOuter>) {
      static_assert(std::is_trivially_copyable_v<TInner>);
      std::memcpy(p, o.data(), len * sizeof(TInner));
    } else {
      auto o_iter = o.begin();
      auto p_ptr = static_cast<TInner*>(p);
      for (size_t i = 0; i < len; i++) {
        new (p_ptr++) TInner(*o_iter++);
      }
    }
  }

  InnerVlen(InnerVlen const&) = delete;

  ~InnerVlen() {
    if (p != nullptr) {
      if constexpr (!std::is_trivially_destructible_v<TInner>) {
        for (size_t i = 0; i < len; i++) {
          auto inner_object = (TInner*)p + i;
          inner_object->~TInner();
        }
      }

      free(p);
      p = nullptr;
      len = 0;
    }
  }

  InnerVlen& operator=(InnerVlen const&) = delete;

  explicit operator std::vector<TOuter>() const {
    std::vector<TOuter> vec;
    if (len > 0) {
      vec.reserve(len);
      TInner* inner_objects = static_cast<TInner*>(p);

      for (size_t i = 0; i < len; i++) {
        TInner& inner_object = inner_objects[i];
        vec.push_back(static_cast<TOuter>(inner_object));
      }
    }

    return vec;
  }

  explicit operator NDArray<TOuter, 1>() const {
    std::array<size_t, 1> shape = {len};
    NDArray<TOuter, 1> rtn(shape);
    if (len > 0) {
      if constexpr (std::is_same_v<TInner, TOuter>) {
        static_assert(std::is_trivially_copyable_v<TInner>);
        std::memcpy(rtn.data(), static_cast<TInner*>(p), len * sizeof(TInner));
      } else {
        TInner* inner_objects = static_cast<TInner*>(p);
        auto rtn_iter = rtn.begin();
        for (size_t i = 0; i < len; i++) {
          TInner& inner_object = inner_objects[i];
          *rtn_iter++ = static_cast<TOuter>(inner_object);
        }
      }
    }

    return rtn;
  }
};

/**
 * @brief An HDF5-compatible representation of variable-length std::string
 */
struct InnerVlenString {
  InnerVlenString() : c_str(nullptr) {}
  InnerVlenString(std::string s) {
    c_str = static_cast<char*>(MallocOrThrow((s.size() + 1) * sizeof(std::string::value_type)));
    std::memcpy(c_str, s.c_str(), (s.size() + 1) * sizeof(std::string::value_type));
  }

  ~InnerVlenString() {
    free(c_str);
    c_str = nullptr;
  }

  InnerVlenString(InnerVlenString const&) = delete;

  InnerVlenString& operator=(InnerVlenString const& other) = delete;

  explicit operator std::string() const {
    return std::string(c_str);
  }

  char* c_str;
};

/**
 * @brief An HDF5-compatible representation of a fixed-size vector (std::array).
 */
template <typename TInner, typename TOuter, size_t N>
class InnerFixedVector : public std::array<TInner, N> {
 public:
  InnerFixedVector() {}
  InnerFixedVector(std::array<TOuter, N> const& o) {
    for (size_t i = 0; i < N; i++) {
      new (&(*this)[i]) TInner(o[i]);
    }
  }

  explicit operator std::array<TOuter, N>() const {
    std::array<TOuter, N> rtn;
    for (size_t i = 0; i < N; i++) {
      rtn[i] = static_cast<TOuter>((*this)[i]);
    }

    return rtn;
  }
};

/**
 * @brief An HDF5-compatible representation a fixed-size multidimensional array.
 */
template <typename TInner, typename TOuter, size_t... Dims>
class InnerFixedNdArray : public yardl::FixedNDArray<TInner, Dims...> {
 public:
  InnerFixedNdArray() {}
  InnerFixedNdArray(yardl::FixedNDArray<TOuter, Dims...> const& o) {
    auto o_iter = o.begin();
    auto i_iter = this->begin();
    for (size_t i = 0; i < length; i++) {
      new (i_iter++) TInner(*o_iter++);
    }
  }

  template <typename TFixedNDArray,
            std::enable_if_t<std::is_base_of_v<
                                 yardl::FixedNDArray<TOuter, Dims...>,
                                 TFixedNDArray>,
                             bool> = true>
  explicit operator TFixedNDArray() const {
    TFixedNDArray o{};
    auto o_iter = o.begin();
    auto i_iter = this->begin();
    constexpr size_t length = (Dims * ...);
    for (size_t i = 0; i < length; i++) {
      *o_iter++ = static_cast<TOuter>(*i_iter++);
    }

    return o;
  }

 private:
  static constexpr size_t length = (Dims * ...);
};

/**
 * @brief An HDF5-compatible representation of an NDArray with a
 * known number of dimensions
 */
template <typename TInner, typename TOuter, size_t N>
struct InnerNdArray {
  InnerNdArray() : dimensions_{}, data_{0, nullptr} {}

  InnerNdArray(NDArray<TOuter, N> const& o)
      : dimensions_(o.shape()), data_{o.size(), malloc(o.size() * sizeof(TInner))} {
    if constexpr (std::is_same_v<TInner, TOuter>) {
      std::memcpy(data_.p, o.data(), data_.len * sizeof(TInner));
    } else {
      auto o_iter = o.begin();
      auto p = static_cast<TInner*>(data_.p);
      for (size_t i = 0; i < data_.len; i++) {
        new (p++) TInner(*o_iter++);
      }
    }
  }

  InnerNdArray(InnerNdArray<TInner, TOuter, N> const&) = delete;

  ~InnerNdArray() {
    if (data_.p != nullptr) {
      if constexpr (!std::is_trivially_destructible_v<TInner>) {
        for (size_t i = 0; i < data_.len; i++) {
          auto inner_object = static_cast<TInner*>(data_.p) + i;
          inner_object->~TInner();
        }
      }

      free(data_.p);
      data_.p = nullptr;
      data_.len = 0;
    }
  }

  InnerNdArray<TInner, TOuter, N>& operator=(InnerNdArray<TInner, TOuter, N> const&) = delete;

  template <typename TNDArray, std::enable_if_t<std::is_base_of_v<NDArray<TOuter, N>, TNDArray>, bool> = true>
  explicit operator TNDArray() const {
    TNDArray rtn(dimensions_);
    if (data_.len > 0) {
      if constexpr (std::is_same_v<TInner, TOuter>) {
        static_assert(std::is_trivially_copyable_v<TInner>);
        std::memcpy(rtn.data(), static_cast<TInner*>(data_.p), data_.len * sizeof(TInner));
      } else {
        TInner* inner_objects = static_cast<TInner*>(data_.p);
        auto rtn_iter = rtn.begin();
        for (size_t i = 0; i < data_.len; i++) {
          TInner& inner_object = inner_objects[i];
          *rtn_iter++ = static_cast<TOuter>(inner_object);
        }
      }
    }

    return rtn;
  }

  std::array<size_t, N> dimensions_;
  hvl_t data_;
};

/**
 * @brief An HDF5-compatible representation of a DynamicNDArray (unknown number of dimensions).
 */
template <typename TInner, typename TOuter>
struct InnerDynamicNdArray {
  InnerDynamicNdArray() : dimensions_{0, nullptr}, data_{0, nullptr} {}

  InnerDynamicNdArray(DynamicNDArray<TOuter> const& o)
      : dimensions_{o.dimension(), MallocOrThrow(o.dimension() * sizeof(size_t))},
        data_{o.size(), MallocOrThrow(o.size() * sizeof(TInner))} {
    memcpy(dimensions_.p, o.shape().data(), o.dimension() * sizeof(size_t));
    if constexpr (std::is_same_v<TInner, TOuter>) {
      std::memcpy(data_.p, o.data(), data_.len * sizeof(TInner));
    } else {
      auto o_iter = o.begin();
      auto p = static_cast<TInner*>(data_.p);
      for (size_t i = 0; i < data_.len; i++) {
        new (p++) TInner(*o_iter++);
      }
    }
  }

  InnerDynamicNdArray(InnerDynamicNdArray<TInner, TOuter> const&) = delete;

  ~InnerDynamicNdArray() {
    if (dimensions_.p != nullptr) {
      free(dimensions_.p);
      dimensions_ = hvl_t{};
    }

    if (data_.p != nullptr) {
      if constexpr (!std::is_trivially_destructible_v<TInner>) {
        for (size_t i = 0; i < data_.len; i++) {
          auto inner_object = static_cast<TInner*>(data_.p) + i;
          inner_object->~TInner();
        }
      }

      free(data_.p);
      data_.p = nullptr;
      data_.len = 0;
    }
  }

  InnerDynamicNdArray<TInner, TOuter>& operator=(InnerDynamicNdArray<TInner, TOuter> const&) = delete;

  template <typename TDynamicNDArray,
            std::enable_if_t<std::is_base_of_v<DynamicNDArray<TOuter>, TDynamicNDArray>, bool> = true>
  explicit operator TDynamicNDArray() const {
    std::vector<size_t> dims(static_cast<size_t*>(dimensions_.p),
                             static_cast<size_t*>(dimensions_.p) + dimensions_.len);
    TDynamicNDArray rtn(dims);
    if (data_.len > 0) {
      if constexpr (std::is_same_v<TInner, TOuter>) {
        static_assert(std::is_trivially_copyable_v<TInner>);
        std::memcpy(rtn.data(), static_cast<TInner*>(data_.p), data_.len * sizeof(TInner));
      } else {
        TInner* inner_objects = static_cast<TInner*>(data_.p);
        auto rtn_iter = rtn.begin();
        for (size_t i = 0; i < data_.len; i++) {
          TInner& inner_object = inner_objects[i];
          *rtn_iter++ = static_cast<TOuter>(inner_object);
        }
      }
    }

    return rtn;
  }

  hvl_t dimensions_;
  hvl_t data_;
};

template <typename TInner, typename TOuter>
class InnerTypeBuffer {
 public:
  InnerTypeBuffer(size_t size) : data_(size * sizeof(TInner)) {
  }

  InnerTypeBuffer(std::vector<TOuter> const& o) : data_(o.size() * sizeof(TInner)) {
    static_assert(!std::is_same_v<TInner, TOuter>, "InnerTypeBuffer should only be used for type conversion");
    auto p = reinterpret_cast<TInner*>(data_.data());
    for (size_t i = 0; i < o.size(); i++) {
      new (p + i) TInner(o[i]);
    }
  }

  ~InnerTypeBuffer() {
    if constexpr (!std::is_trivially_destructible_v<TInner>) {
      auto p = reinterpret_cast<TInner*>(data_.data());
      size_t count = data_.size() / sizeof(TInner);
      for (size_t i = 0; i < count; i++) {
        (p + i)->~TInner();
      }
    }
  }

  TInner* data() {
    return reinterpret_cast<TInner*>(data_.data());
  }

  TInner const* begin() const {
    return reinterpret_cast<TInner const*>(data_.data());
  }

  TInner const* end() const {
    return begin() + data_.size() / sizeof(TInner);
  }

 private:
  std::vector<uint8_t> data_;
};

}  // namespace yardl::hdf5
