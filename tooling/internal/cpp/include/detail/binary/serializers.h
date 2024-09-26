// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <complex>
#include <memory>
#include <optional>
#include <unordered_map>
#include <utility>
#include <variant>
#include <vector>

#include "../../yardl.h"
#include "coded_stream.h"

namespace yardl::binary {

/**
 * @brief Function pointer type for writing a single value to a stream.
 *
 * @tparam T the type of the value to write
 */
template <typename T>
using Writer = void (*)(CodedOutputStream& stream, T const& value);

/**
 * @brief Function pointer type for reading a single value from a stream.
 *
 * @tparam T the type of the value to read
 */
template <typename T>
using Reader = void (*)(CodedInputStream& stream, T& value);

/**
 * If T can be serialized using a simple memcpy, provides the member
 * constant value equal to true. Otherwise value is false.
 */
template <typename T, typename = void>
struct IsTriviallySerializable
    : std::false_type {
};

template <typename T>
struct IsTriviallySerializable<T, typename std::enable_if_t<std::is_integral_v<T> &&
                                                            sizeof(T) == 1>>
    : std::true_type {
};

#if __BYTE_ORDER__ == __ORDER_LITTLE_ENDIAN__
template <typename T>
struct IsTriviallySerializable<T, typename std::enable_if_t<std::is_floating_point_v<T>>>
    : std::true_type {
};

template <typename T>
struct IsTriviallySerializable<T, typename std::enable_if_t<
                                      std::is_same_v<T, std::complex<float>> ||
                                      std::is_same_v<T, std::complex<double>>>>
    : std::true_type {
};
#endif

template <typename T, size_t N>
struct IsTriviallySerializable<std::array<T, N>,
                               typename std::enable_if_t<IsTriviallySerializable<T>::value>>
    : std::true_type {
};

// TODO Joe: We can no longer assume FixedNDArray is trivially serializable e.g. if implemented by user
// template <typename T, size_t... Dims>
// struct IsTriviallySerializable<yardl::FixedNDArray<T, Dims...>,
//                                typename std::enable_if_t<IsTriviallySerializable<T>::value>>
//     : std::true_type {
// };

template <typename T>
inline void WriteTriviallySerializable(CodedOutputStream& stream, T const& value) {
  static_assert(IsTriviallySerializable<T>::value, "T must be trivially serializable");
  stream.WriteBytes(reinterpret_cast<char const*>(std::addressof(value)), sizeof(value));
}

template <typename T>
inline void ReadTriviallySerializable(CodedInputStream& stream, T& value) {
  static_assert(IsTriviallySerializable<T>::value, "T must be trivially serializable");
  stream.ReadBytes(reinterpret_cast<char*>(std::addressof(value)), sizeof(value));
}

template <typename T, std::enable_if_t<std::is_integral_v<T> && sizeof(T) == 1, bool> = true>
inline void WriteInteger(CodedOutputStream& stream, T const& value) {
  stream.WriteByte(value);
}

template <typename T, std::enable_if_t<std::is_integral_v<T> && sizeof(T) == 1, bool> = true>
inline void ReadInteger(CodedInputStream& stream, T& value) {
  stream.ReadByte(value);
}

template <typename T, std::enable_if_t<std::is_integral_v<T> && sizeof(T) == 2, bool> = true>
inline void WriteInteger(CodedOutputStream& stream, T const& value) {
  if constexpr (std::is_signed_v<T>) {
    stream.WriteVarInt32(static_cast<int32_t>(value));
  } else {
    stream.WriteVarInt32(static_cast<uint32_t>(value));
  }
}

template <typename T, std::enable_if_t<std::is_integral_v<T> && sizeof(T) == 2, bool> = true>
inline void ReadInteger(CodedInputStream& stream, T& value) {
  if constexpr (std::is_signed_v<T>) {
    int32_t val;
    stream.ReadVarInt32(val);
    value = val;
  } else {
    uint32_t val;
    stream.ReadVarInt32(val);
    value = val;
  }
}

template <typename T, std::enable_if_t<std::is_integral_v<T> && sizeof(T) == 4, bool> = true>
inline void WriteInteger(CodedOutputStream& stream, T const& value) {
  stream.WriteVarInt32(value);
}

template <typename T, std::enable_if_t<std::is_integral_v<T> && sizeof(T) == 4, bool> = true>
inline void ReadInteger(CodedInputStream& stream, T& value) {
  stream.ReadVarInt32(value);
}

template <typename T, std::enable_if_t<std::is_integral_v<T> && sizeof(T) == 8, bool> = true>
inline void WriteInteger(CodedOutputStream& stream, T const& value) {
  stream.WriteVarInt64(value);
}

template <typename T, std::enable_if_t<std::is_integral_v<T> && sizeof(T) == 8, bool> = true>
inline void ReadInteger(CodedInputStream& stream, T& value) {
  stream.ReadVarInt64(value);
}

template <typename T, std::enable_if_t<std::is_floating_point_v<T> ||
                                           std::is_same_v<T, std::complex<float>> ||
                                           std::is_same_v<T, std::complex<double>>,
                                       bool> = true>
inline void WriteFloatingPoint(CodedOutputStream& stream, T const& value) {
  WriteTriviallySerializable(stream, value);
}

template <typename T, std::enable_if_t<std::is_floating_point_v<T> ||
                                           std::is_same_v<T, std::complex<float>> ||
                                           std::is_same_v<T, std::complex<double>>,
                                       bool> = true>
inline void ReadFloatingPoint(CodedInputStream& stream, T& value) {
  ReadTriviallySerializable(stream, value);
}

inline void WriteString(CodedOutputStream& stream, std::string const& value) {
  WriteInteger(stream, value.size());
  stream.WriteBytes(value.data(), value.size());
}

inline void ReadString(CodedInputStream& stream, std::string& value) {
  size_t size;
  ReadInteger(stream, size);
  value.resize(size);
  stream.ReadBytes(value.data(), size);
}

inline void WriteDate(CodedOutputStream& stream, yardl::Date const& value) {
  auto days = value.time_since_epoch().count();
  WriteInteger(stream, days);
}

inline void ReadDate(CodedInputStream& stream, yardl::Date& value) {
  int64_t days;
  ReadInteger(stream, days);
  value = yardl::Date(date::days(days));
}

inline void WriteTime(CodedOutputStream& stream, yardl::Time const& value) {
  WriteInteger(stream, value.count());
}

inline void ReadTime(CodedInputStream& stream, yardl::Time& value) {
  int64_t count;
  ReadInteger(stream, count);
  value = yardl::Time(count);
}

inline void WriteDateTime(CodedOutputStream& stream, yardl::DateTime const& value) {
  auto ns = value.time_since_epoch().count();
  WriteInteger(stream, ns);
}

inline void ReadDateTime(CodedInputStream& stream, yardl::DateTime& value) {
  int64_t ns;
  ReadInteger(stream, ns);
  value = yardl::DateTime{
      std::chrono::time_point<std::chrono::system_clock,
                              std::chrono::nanoseconds>(std::chrono::nanoseconds(ns))};
}

template <typename T, Writer<T> WriteElement>
inline void WriteOptional(CodedOutputStream& stream, std::optional<T> const& value) {
  stream.WriteByte(value.has_value());
  if (value.has_value()) {
    WriteElement(stream, value.value());
  }
}

template <typename T, Reader<T> ReadElement>
inline void ReadOptional(CodedInputStream& stream, std::optional<T>& value) {
  bool has_value;
  stream.ReadByte(has_value);
  if (has_value) {
    T tmp;
    ReadElement(stream, tmp);
    value = std::move(tmp);
  } else {
    value = std::nullopt;
  }
}

template <typename T, Writer<T> WriteElement>
inline void WriteVector(CodedOutputStream& stream, std::vector<T> const& value) {
  WriteInteger(stream, value.size());

  if constexpr (IsTriviallySerializable<T>::value) {
    stream.WriteBytes(value.data(), value.size() * sizeof(T));
    return;
  }

  for (auto const& element : value) {
    WriteElement(stream, element);
  }
}

template <typename T, Reader<T> ReadElement>
inline void ReadVector(CodedInputStream& stream, std::vector<T>& value) {
  uint64_t size;
  ReadInteger(stream, size);
  value.resize(size);

  if constexpr (IsTriviallySerializable<T>::value) {
    stream.ReadBytes(value.data(), value.size() * sizeof(T));
    return;
  }

  for (size_t i = 0; i < size; i++) {
    ReadElement(stream, value[i]);
  }
}

template <typename T, Writer<T> WriteElement, size_t N>
inline void WriteArray(CodedOutputStream& stream, std::array<T, N> const& value) {
  if constexpr (IsTriviallySerializable<T>::value) {
    stream.WriteBytes(value.data(), value.size() * sizeof(T));
    return;
  }

  for (size_t i = 0; i < N; i++) {
    WriteElement(stream, value[i]);
  }
}

template <typename T, Reader<T> ReadElement, size_t N>
inline void ReadArray(CodedInputStream& stream, std::array<T, N>& value) {
  if constexpr (IsTriviallySerializable<T>::value) {
    stream.ReadBytes(value.data(), value.size() * sizeof(T));
    return;
  }

  for (size_t i = 0; i < N; i++) {
    ReadElement(stream, value[i]);
  }
}

template <typename T, Writer<T> WriteElement>
inline void WriteDynamicNDArray(CodedOutputStream& stream, yardl::DynamicNDArray<T> const& value) {
  auto shape = yardl::shape(value);
  WriteInteger(stream, shape.size());
  for (auto const& dim : shape) {
    WriteInteger(stream, dim);
  }

  if constexpr (IsTriviallySerializable<T>::value) {
    stream.WriteBytes(yardl::dataptr(value), yardl::size(value) * sizeof(T));
    return;
  }

  for (auto const& element : value) {
    WriteElement(stream, element);
  }
}

template <typename T, Reader<T> ReadElement>
inline void ReadDynamicNDArray(CodedInputStream& stream, yardl::DynamicNDArray<T>& value) {
  std::vector<size_t> shape;
  ReadVector<size_t, &ReadInteger>(stream, shape);
  yardl::resize(value, shape);

  if constexpr (IsTriviallySerializable<T>::value) {
    stream.ReadBytes(yardl::dataptr(value), yardl::size(value) * sizeof(T));
    return;
  }

  for (auto& element : value) {
    ReadElement(stream, element);
  }
}

template <typename T, Writer<T> WriteElement, size_t N>
inline void WriteNDArray(CodedOutputStream& stream, yardl::NDArray<T, N> const& value) {
  for (auto const& dim : yardl::shape(value)) {
    WriteInteger(stream, dim);
  }

  if constexpr (IsTriviallySerializable<T>::value) {
    stream.WriteBytes(yardl::dataptr(value), yardl::size(value) * sizeof(T));
    return;
  }

  for (auto const& element : value) {
    WriteElement(stream, element);
  }
}

template <typename T, Reader<T> ReadElement, size_t N>
inline void ReadNDArray(CodedInputStream& stream, yardl::NDArray<T, N>& value) {
  std::array<size_t, N> shape;
  ReadArray<size_t, &ReadInteger, N>(stream, shape);
  yardl::resize(value, shape);

  if constexpr (IsTriviallySerializable<T>::value) {
    stream.ReadBytes(yardl::dataptr(value), yardl::size(value) * sizeof(T));
    return;
  }

  for (auto& element : value) {
    ReadElement(stream, element);
  }
}

template <typename T, Writer<T> WriteElement, size_t... Dims>
inline void WriteFixedNDArray(CodedOutputStream& stream,
                              yardl::FixedNDArray<T, Dims...> const& value) {
  if constexpr (IsTriviallySerializable<T>::value) {
    stream.WriteBytes(yardl::dataptr(value), yardl::size(value) * sizeof(T));
    return;
  }

  for (auto const& element : value) {
    WriteElement(stream, element);
  }
}

template <typename T, Reader<T> ReadElement, size_t... Dims>
inline void ReadFixedNDArray(CodedInputStream& stream, yardl::FixedNDArray<T, Dims...>& value) {
  if constexpr (IsTriviallySerializable<T>::value) {
    stream.ReadBytes(yardl::dataptr(value), yardl::size(value) * sizeof(T));
    return;
  }

  for (auto& element : value) {
    ReadElement(stream, element);
  }
}

template <typename TKey, typename TValue, Writer<TKey> WriteKey, Writer<TValue> WriteValue>
inline void WriteMap(CodedOutputStream& stream, std::unordered_map<TKey, TValue> const& value) {
  WriteInteger(stream, value.size());
  for (auto const& [key, value] : value) {
    WriteKey(stream, key);
    WriteValue(stream, value);
  }
}

template <typename TKey, typename TValue, Reader<TKey> ReadKey, Reader<TValue> ReadValue>
inline void ReadMap(CodedInputStream& stream, std::unordered_map<TKey, TValue>& value) {
  uint64_t size;
  ReadInteger(stream, size);

  for (size_t i = 0; i < size; i++) {
    TKey k;
    ReadKey(stream, k);
    TValue v;
    ReadValue(stream, v);
    value.emplace(std::move(k), std::move(v));
  }
}

inline void WriteMonostate([[maybe_unused]] CodedOutputStream& stream,
                           [[maybe_unused]] std::monostate const& value) {
}

inline void ReadMonostate([[maybe_unused]] CodedInputStream& stream,
                          [[maybe_unused]] std::monostate& value) {
}

template <typename T>
inline void WriteEnum(CodedOutputStream& stream, T const& value) {
  using underlying_type = std::underlying_type_t<T>;
  underlying_type underlying_value = static_cast<underlying_type>(value);
  yardl::binary::WriteInteger(stream, underlying_value);
}

template <typename T>
inline void ReadEnum(CodedInputStream& stream, T& value) {
  using underlying_type = std::underlying_type_t<T>;
  underlying_type underlying_value;
  yardl::binary::ReadInteger(stream, underlying_value);
  value = static_cast<T>(underlying_value);
}

template <typename T>
inline void WriteFlags(CodedOutputStream& stream, T const& value) {
  yardl::binary::WriteInteger(stream, value.Value());
}

template <typename T>
inline void ReadFlags(CodedInputStream& stream, T& value) {
  using underlying_type = typename T::value_type;
  underlying_type underlying_value;
  yardl::binary::ReadInteger(stream, underlying_value);
  value = underlying_value;
}

template <typename T, Writer<T> WriteElement>
inline void WriteBlock(CodedOutputStream& stream, T const& source) {
  WriteInteger(stream, 1U);
  WriteElement(stream, source);
}

template <typename T, Reader<T> ReadElement>
inline bool ReadBlock(CodedInputStream& stream, size_t& current_block_remaining, T& destination) {
  if (current_block_remaining == 0) {
    ReadInteger(stream, current_block_remaining);
    if (current_block_remaining == 0) {
      return false;
    }
  }

  ReadElement(stream, destination);
  current_block_remaining--;
  return true;
}

template <typename T, Reader<T> ReadElement>
inline void ReadBlocksIntoVector(CodedInputStream& stream, size_t& current_block_remaining, std::vector<T>& destination) {
  if (current_block_remaining == 0) {
    ReadInteger(stream, current_block_remaining);
  }

  size_t offset = 0;
  size_t remaining_capacity = destination.capacity();
  while (current_block_remaining > 0) {
    size_t read_count = std::min(current_block_remaining, remaining_capacity);
    if (read_count + offset > destination.size()) {
      destination.resize(offset + read_count);
    }

    if constexpr (IsTriviallySerializable<T>::value) {
      stream.ReadBytes(destination.data() + offset, read_count * sizeof(T));
    } else {
      for (size_t i = 0; i < read_count; i++) {
        ReadElement(stream, destination[offset + i]);
      }
    }

    current_block_remaining -= read_count;
    offset += read_count;
    remaining_capacity -= read_count;

    if (current_block_remaining == 0) {
      ReadInteger(stream, current_block_remaining);
    }

    if (remaining_capacity == 0) {
      return;
    }
  }

  if (remaining_capacity > 0) {
    destination.resize(offset);
  }
}
}  // namespace yardl::binary
