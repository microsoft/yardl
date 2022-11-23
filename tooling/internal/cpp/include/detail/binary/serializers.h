// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <complex>
#include <memory>
#include <optional>
#include <utility>
#include <vector>

#include "../../yardl.h"
#include "coded_stream.h"

#if __cplusplus >= 202002L
#include <bit>
#endif

namespace yardl::binary {

/**
 * bit_cast() is used below to perform bitwise casts between floating point and
 * integer types. This should be used instead of expressions like
 * reinterpret_cast<int&>(float_value) because accessing the result of
 * the latter is technically undefined behavior.
 */
#if __cplusplus >= 202002L
using std::bit_cast;
#else
template <typename Dest,
          typename Source,
          typename std::enable_if_t<
              sizeof(Dest) == sizeof(Source) &&
                  std::is_trivially_copyable_v<Source> &&
                  std::is_trivially_copyable_v<Dest> &&
                  std::is_default_constructible_v<Dest>,
              bool> = true>
inline Dest bit_cast(Source const& source) {
  Dest dest;
  memcpy(static_cast<void*>(std::addressof(dest)),
         static_cast<void const*>(std::addressof(source)), sizeof(dest));
  return dest;
}
#endif

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

template <typename T, size_t... Dims>
struct IsTriviallySerializable<yardl::FixedNDArray<T, Dims...>,
                               typename std::enable_if_t<IsTriviallySerializable<T>::value>>
    : std::true_type {
};

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

inline void WriteFloat(CodedOutputStream& stream, float const& value) {
  stream.WriteFixedInteger(bit_cast<uint32_t>(value));
}

inline void ReadFloat(CodedInputStream& stream, float& value) {
  uint32_t tmp;
  stream.ReadFixedInteger(tmp);
  value = bit_cast<float>(tmp);
}

inline void WriteDouble(CodedOutputStream& stream, double const& value) {
  stream.WriteFixedInteger(bit_cast<uint64_t>(value));
}

inline void ReadDouble(CodedInputStream& stream, double& value) {
  uint64_t tmp;
  stream.ReadFixedInteger(tmp);
  value = bit_cast<double>(tmp);
}

inline void WriteComplexFloat(CodedOutputStream& stream, std::complex<float> const& value) {
  auto arr = bit_cast<std::array<uint32_t, 2>>(value);
  stream.WriteBytes(arr.data(), sizeof(arr));
}

inline void ReadComplexFloat(CodedInputStream& stream, std::complex<float>& value) {
  std::array<uint32_t, 2> arr;
  stream.ReadBytes(arr.data(), sizeof(arr));
  value = bit_cast<std::complex<float>>(arr);
}

inline void WriteComplexDouble(CodedOutputStream& stream, std::complex<double> const& value) {
  auto arr = bit_cast<std::array<uint64_t, 2>>(value);
  stream.WriteBytes(arr.data(), sizeof(arr));
}

inline void ReadComplexDouble(CodedInputStream& stream, std::complex<double>& value) {
  std::array<uint64_t, 2> arr;
  stream.ReadBytes(arr.data(), sizeof(arr));
  value = bit_cast<std::complex<double>>(arr);
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
#if __cplusplus < 202002L
  value = yardl::Date(date::days(days));
#else
  value = yardl::Date{std::chrono::days{days}};
#endif
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
#if __cplusplus < 202002L
  value = yardl::DateTime(date::sys_time<std::chrono::nanoseconds>(std::chrono::nanoseconds(ns)));
#else
  value = yardl::DateTime{
      std::chrono::time_point<std::chrono::system_clock,
                              std::chrono::nanoseconds>(std::chrono::nanoseconds(ns))};
#endif
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
  auto shape = value.shape();
  WriteInteger(stream, shape.size());
  for (auto const& dim : shape) {
    WriteInteger(stream, dim);
  }

  if constexpr (IsTriviallySerializable<T>::value) {
    stream.WriteBytes(value.data(), value.size() * sizeof(T));
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
  value.resize(shape);

  if constexpr (IsTriviallySerializable<T>::value) {
    stream.ReadBytes(value.data(), value.size() * sizeof(T));
    return;
  }

  for (auto& element : value) {
    ReadElement(stream, element);
  }
}

template <typename T, Writer<T> WriteElement, size_t N>
inline void WriteNDArray(CodedOutputStream& stream, yardl::NDArray<T, N> const& value) {
  for (auto const& dim : value.shape()) {
    WriteInteger(stream, dim);
  }

  if constexpr (IsTriviallySerializable<T>::value) {
    stream.WriteBytes(value.data(), value.size() * sizeof(T));
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
  value.resize(shape);

  if constexpr (IsTriviallySerializable<T>::value) {
    stream.ReadBytes(value.data(), value.size() * sizeof(T));
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
    stream.WriteBytes(value.data(), value.size() * sizeof(T));
    return;
  }

  for (auto const& element : value) {
    WriteElement(stream, element);
  }
}

template <typename T, Reader<T> ReadElement, size_t... Dims>
inline void ReadFixedNDArray(CodedInputStream& stream, yardl::FixedNDArray<T, Dims...>& value) {
  if constexpr (IsTriviallySerializable<T>::value) {
    stream.ReadBytes(value.data(), value.size() * sizeof(T));
    return;
  }

  for (auto& element : value) {
    ReadElement(stream, element);
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
