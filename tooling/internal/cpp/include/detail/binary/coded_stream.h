// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <cassert>
#include <complex>
#include <cstring>
#include <memory>
#include <utility>
#include <vector>

namespace yardl::binary {

static int const MAX_VARINT32_BYTES = 5;
static int const MAX_VARINT64_BYTES = 10;

/**
 * An exception thrown when EOF is reached prematurely.
 */
class EndOfStreamException : public std::exception {
 public:
  char const* what() const noexcept override { return "Unexpected end of stream"; }
};

/**
 * A buffered output stream that provides methods for writing integers in a
 * compact form. Unsigned integers written using Protobuf "varint" encoding,
 * and signed integers are first converted to unsigned using Protobuf's
 * "zig-zag" encoding.
 */
class CodedOutputStream {
 public:
  CodedOutputStream(std::ostream& stream, size_t buffer_size = 65536)
      : stream_(stream),
        buffer_(buffer_size),
        buffer_ptr_(buffer_.data()),
        buffer_end_ptr_(buffer_ptr_ + buffer_.size()),
        pos_(0) {
  }

  ~CodedOutputStream() {
    Flush();
  }

  template <typename T, std::enable_if_t<std::is_integral_v<T> && sizeof(T) == 1, bool> = true>
  void WriteByte(T const& v) {
    if (RemainingBufferSpace() == 0) {
      FlushBuffer();
    }

    *buffer_ptr_++ = v;
    pos_ += 1;
  }

  void WriteVarInt32(uint32_t const& value) {
    if (RemainingBufferSpace() < MAX_VARINT32_BYTES) {
      FlushBuffer();
    }

    WriteVarInt(value);
  }

  void WriteVarInt32(int32_t const& value) {
    WriteVarInt32(ZigZagEncode32(value));
  }

  template <typename T, std::enable_if_t<std::is_integral_v<T> &&
                                             sizeof(T) == 8 &&
                                             std::is_unsigned_v<T>,
                                         bool> = true>
  void WriteVarInt64(T const& value) {
    if (RemainingBufferSpace() < MAX_VARINT64_BYTES) {
      FlushBuffer();
    }

    WriteVarInt(value);
  }

  void WriteVarInt64(int64_t const& value) {
    WriteVarInt64(ZigZagEncode64(value));
  }

  template <typename T, std::enable_if_t<std::is_integral_v<T>, bool> = true>
  void WriteFixedInteger(T const& value) {
    if (RemainingBufferSpace() < sizeof(value)) {
      FlushBuffer();
    }

#if __BYTE_ORDER__ == __ORDER_LITTLE_ENDIAN__
    memcpy(buffer_ptr_, &value, sizeof(value));
#else
    static_assert(false, "Unsupported byte order");
#endif

    buffer_ptr_ += sizeof(value);
    pos_ += sizeof(value);
  }

  void WriteBytes(void const* data, size_t size_in_bytes) {
    while (true) {
      size_t const remaining_buffer_space = RemainingBufferSpace();
      if (remaining_buffer_space >= size_in_bytes) {
        memcpy(buffer_ptr_, data, size_in_bytes);
        buffer_ptr_ += size_in_bytes;
        pos_ += size_in_bytes;
        break;
      }

      if (remaining_buffer_space > 0) {
        memcpy(buffer_ptr_, data, remaining_buffer_space);
        buffer_ptr_ += remaining_buffer_space;
        pos_ += remaining_buffer_space;
        data = static_cast<uint8_t const*>(data) + remaining_buffer_space;
        size_in_bytes -= remaining_buffer_space;
      }

      FlushBuffer();
    }
  }

  void Flush() {
    FlushBuffer();
    stream_.flush();
  }

  size_t Pos() {
    return pos_;
  }

 private:
  size_t RemainingBufferSpace() {
    assert(buffer_ptr_ <= buffer_end_ptr_);
    return buffer_end_ptr_ - buffer_ptr_;
  }

  template <typename T, std::enable_if_t<std::is_integral_v<T> && std::is_unsigned_v<T>, bool> = true>
  void WriteVarInt(T value) {
    auto start = buffer_ptr_;
    while (value > 0x7F) {
      *buffer_ptr_++ = static_cast<uint8_t>(value) | 0x80;
      value >>= 7;
    }

    *buffer_ptr_++ = static_cast<uint8_t>(value);
    pos_ += buffer_ptr_ - start;
  }

  static uint32_t ZigZagEncode32(int32_t v) {
    return (static_cast<uint32_t>(v) << 1) ^ static_cast<uint32_t>(v >> 31);
  }

  static uint64_t ZigZagEncode64(int64_t v) {
    return (static_cast<uint64_t>(v) << 1) ^ static_cast<uint64_t>(v >> 63);
  }

  void FlushBuffer() {
    if (buffer_ptr_ == buffer_.data()) {
      return;
    }

    stream_.write(reinterpret_cast<char*>(const_cast<uint8_t*>(buffer_.data())),
                  buffer_ptr_ - buffer_.data());
    buffer_ptr_ = buffer_.data();
    if (stream_.bad()) {
      throw std::runtime_error("Failed to write to stream");
    }
  }

  std::ostream& stream_;
  std::vector<uint8_t> buffer_;
  uint8_t* buffer_ptr_;
  uint8_t* buffer_end_ptr_;
  size_t pos_;
};

/**
 * A buffered input stream that provides methods reading data written
 * using a CodedOutputStream.
 */
class CodedInputStream {
 public:
  CodedInputStream(std::istream& stream, size_t buffer_size = 65536)
      : stream_(stream),
        buffer_(buffer_size),
        buffer_ptr_(buffer_.data()),
        buffer_end_ptr_(buffer_ptr_),
        last_read_pos_(0) {
  }

 public:
  template <typename T, std::enable_if_t<std::is_integral_v<T> && sizeof(T) == 1, bool> = true>
  void ReadByte(T& v) {
    if (buffer_ptr_ == buffer_end_ptr_) {
      FillBuffer();
    }
    v = *buffer_ptr_++;
  }

  template <typename T, std::enable_if_t<std::is_integral_v<T>, bool> = true>
  void ReadFixedInteger(T& value) {
    if (RemainingBufferSpace() < sizeof(value)) {
      ReadFixedIntegerSlow(value);
    } else {
      ReadFixedIntegerFastFromArray(value, buffer_ptr_);
    }
  }

  void ReadVarInt32(uint32_t& value) {
    if (RemainingBufferSpace() < MAX_VARINT32_BYTES) {
      ReadVarIntegerSlow(value);
    } else {
      ReadVarIntegerFastFromArray(value, buffer_ptr_);
    }
  }

  void ReadVarInt32(int32_t& value) {
    uint32_t v;
    ReadVarInt32(v);
    value = ZigZagDecode32(v);
  }

  template <typename T, std::enable_if_t<std::is_integral_v<T> &&
                                             sizeof(T) == 8 &&
                                             std::is_unsigned_v<T>,
                                         bool> = true>
  void ReadVarInt64(T& value) {
    if (RemainingBufferSpace() < MAX_VARINT64_BYTES) {
      ReadVarIntegerSlow(value);
    } else {
      ReadVarIntegerFastFromArray(value, buffer_ptr_);
    }
  }

  void ReadVarInt64(int64_t& value) {
    uint64_t v;
    ReadVarInt64(v);
    value = ZigZagDecode64(v);
  }

  void ReadBytes(void* data, size_t size_in_bytes) {
    uint8_t* uint8_data = static_cast<uint8_t*>(data);
    while (size_in_bytes > 0) {
      if (buffer_ptr_ == buffer_end_ptr_) {
        FillBuffer();
      }

      size_t bytes_to_copy = std::min(
          size_in_bytes,
          static_cast<size_t>(buffer_end_ptr_ - buffer_ptr_));
      memcpy(uint8_data, buffer_ptr_, bytes_to_copy);
      buffer_ptr_ += bytes_to_copy;
      uint8_data += bytes_to_copy;
      size_in_bytes -= bytes_to_copy;
    }
  }

  void VerifyFinished() {
    if (at_eof_) {
      if (buffer_ptr_ == buffer_end_ptr_) {
        return;
      }
    } else {
      if (buffer_ptr_ == buffer_end_ptr_) {
        FillBuffer();
        if (at_eof_ && buffer_ptr_ == buffer_end_ptr_) {
          return;
        }
      }
    }

    throw std::runtime_error("Stream was not completely read");
  }

  size_t Pos() {
    return last_read_pos_ + (buffer_ptr_ - buffer_.data());
  }

  void Seek(int64_t whence) {
    if (whence < 0) {
      // Seek backward from the end of stream
      stream_.clear();
      auto pos = stream_.tellg();
      stream_.seekg(whence, std::ios::end);
      auto desired = stream_.tellg();
      stream_.seekg(pos);

      if (desired > 0) {
        Seek(desired);
      }
    } else {
      size_t offset = static_cast<size_t>(whence);
      if (offset < last_read_pos_ || offset >= last_read_pos_ + (buffer_end_ptr_ - buffer_.data())) {
        // Seek and re-buffer
        at_eof_ = false;
        stream_.clear();
        stream_.seekg(offset);
        if (stream_.fail()) {
          throw std::runtime_error("Failed to seek in stream");
        }
        FillBuffer();
      } else {
        // Desired offset is already buffered...
        buffer_ptr_ = buffer_.data() + (offset - last_read_pos_);
      }
    }
  }

 private:
  template <typename T, std::enable_if_t<std::is_integral_v<T>, bool> = true>
  static void ReadFixedIntegerFastFromArray(T& value, uint8_t*& local_buffer_ptr) {
#if __BYTE_ORDER__ == __ORDER_LITTLE_ENDIAN__
    memcpy(&value, local_buffer_ptr, sizeof(value));
#else
    static_assert(false, "Unsupported byte order");
#endif

    local_buffer_ptr += sizeof(value);
  }

  template <typename T, std::enable_if_t<std::is_integral_v<T>, bool> = true>
  void ReadFixedIntegerSlow(T& value) {
    if (buffer_ptr_ == buffer_end_ptr_) {
      FillBuffer();
      ReadFixedIntegerFastFromArray(value, buffer_ptr_);
      return;
    }

    uint8_t bytes[sizeof(T)];
    ReadBytes(bytes, sizeof(T));
    uint8_t* bytes_ptr = bytes;
    ReadFixedIntegerFastFromArray(value, bytes_ptr);
  }

  template <typename T, std::enable_if_t<std::is_integral_v<T>, bool> = true>
  static void ReadVarIntegerFastFromArray(T& value, uint8_t*& local_buffer_ptr) {
    value = 0;
    int shift = 0;
    while (true) {
      uint8_t byte = *local_buffer_ptr++;
      value |= static_cast<T>(byte & 0x7F) << shift;
      if ((byte & 0x80) == 0) {
        break;
      }
      shift += 7;
    }
  }

  template <typename T, std::enable_if_t<std::is_integral_v<T>, bool> = true>
  void ReadVarIntegerSlow(T& value) {
    if (buffer_ptr_ == buffer_end_ptr_) {
      FillBuffer();
      ReadVarIntegerFastFromArray(value, buffer_ptr_);
      return;
    }

    value = 0;
    int shift = 0;
    while (true) {
      if (buffer_ptr_ == buffer_end_ptr_) {
        FillBuffer();
      }
      uint8_t byte = *buffer_ptr_++;
      value |= static_cast<T>(byte & 0x7F) << shift;
      if ((byte & 0x80) == 0) {
        break;
      }
      shift += 7;
    }
  }

  static int32_t ZigZagDecode32(uint32_t n) {
    return static_cast<int32_t>((n >> 1) ^ (~(n & 1) + 1));
  }

  static int64_t ZigZagDecode64(uint64_t n) {
    return static_cast<int64_t>((n >> 1) ^ (~(n & 1) + 1));
  }

  size_t FillBuffer() {
    if (at_eof_) {
      throw EndOfStreamException();
    }

    last_read_pos_ = stream_.tellg();
    stream_.read(reinterpret_cast<char*>(buffer_.data()), buffer_.size());
    at_eof_ = stream_.eof();
    auto bytes_read = stream_.gcount();
    buffer_ptr_ = buffer_.data();
    buffer_end_ptr_ = buffer_ptr_ + bytes_read;
    return bytes_read;
  }

  size_t RemainingBufferSpace() {
    assert(buffer_ptr_ <= buffer_end_ptr_);
    return buffer_end_ptr_ - buffer_ptr_;
  }

  std::istream& stream_;
  std::vector<uint8_t> buffer_;
  uint8_t* buffer_ptr_;
  uint8_t* buffer_end_ptr_;
  bool at_eof_ = false;

  size_t last_read_pos_;
};

}  // namespace yardl::binary
