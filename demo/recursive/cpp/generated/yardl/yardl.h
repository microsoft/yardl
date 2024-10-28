// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <chrono>

#include <date/date.h>

#include "detail/ndarray/impl.h"

namespace yardl {

/**
 * @brief Represents a date as a number of days since the epoch.
 */
using Date = date::local_days;

/**
 * @brief Represents a time of day as the number of nanoseconds since midnight.
 */
using Time = std::chrono::duration<int64_t, std::nano>;

/**
 * @brief Represents a datetime as the number of nanoseconds since the epoch.
 */
using DateTime = std::chrono::time_point<std::chrono::system_clock, std::chrono::nanoseconds>;

/**
 * @brief The same as size_t when it is 64 bits, otherwise uint64_t.
 */
using Size = std::conditional_t<sizeof(size_t) == sizeof(uint64_t), size_t, uint64_t>;

/**
 * @brief A base template for generated flags classes

 * @tparam TValue the underlying integral type
 * @tparam TDerived the derived flags class
 */
template <typename TValue, typename TDerived>
struct BaseFlags {
  static_assert(std::is_integral_v<TValue>, "TValue must be an integral type");

 public:
  BaseFlags() = default;
  BaseFlags(TValue value) : value_(value) {}
  BaseFlags(TDerived const& other) : value_(other.value_) {};

  using value_type = TValue;

  [[nodiscard]] TDerived operator|(TDerived rhs) const {
    return TDerived(value_ | rhs.value_);
  }

  [[nodiscard]] TDerived operator&(TDerived rhs) const {
    return TDerived(value_ & rhs.value_);
  }

  [[nodiscard]] TDerived operator^(TDerived rhs) const {
    return TDerived(value_ ^ rhs.value_);
  }

  TDerived operator~() const {
    return TDerived(~value_);
  }

  TDerived& operator=(TDerived const& rhs) {
    value_ = rhs.value_;
    return static_cast<TDerived&>(*this);
  }

  TDerived& operator=(TValue rhs) {
    value_ = rhs;
    return static_cast<TDerived&>(*this);
  }

  TDerived& operator|=(TDerived rhs) {
    value_ |= rhs.value_;
    return static_cast<TDerived&>(*this);
  }

  TDerived& operator&=(TDerived rhs) {
    value_ &= rhs.value_;
    return static_cast<TDerived&>(*this);
  }

  TDerived& operator^=(TDerived rhs) {
    value_ ^= rhs.value_;
    return static_cast<TDerived&>(*this);
  }

  bool operator==(TDerived rhs) const {
    return value_ == rhs.value_;
  }

  bool operator==(TValue rhs) const {
    return value_ == rhs;
  }

  bool operator!=(TDerived rhs) const {
    return value_ != rhs.value_;
  }

  bool operator!=(TValue rhs) const {
    return value_ != rhs;
  }

  [[nodiscard]] bool HasFlags(TDerived flag) const {
    return (value_ & flag.value_) == flag.value_;
  }

  void SetFlags(TDerived flag) {
    value_ |= flag.value_;
  }

  void UnsetFlags(TDerived flag) {
    value_ &= ~flag.value_;
  }

  void Clear() {
    value_ = 0;
  }

  [[nodiscard]] TDerived WithFlags(TDerived flag) const {
    return *this | flag;
  }

  [[nodiscard]] TDerived WithoutFlags(TDerived flag) const {
    return *this & ~flag;
  }

  [[nodiscard]] explicit operator TValue() const {
    return value_;
  }

  [[nodiscard]] TValue Value() const {
    return value_;
  }

 private:
  TValue value_{};
};

}  // namespace yardl
