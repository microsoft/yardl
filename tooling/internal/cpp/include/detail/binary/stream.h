// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <fstream>
#include <memory>
#include <variant>

namespace yardl::binary {

template <typename T, typename = void>
struct IsWritable : std::false_type {
};

template <typename T>
struct IsWritable<T, typename std::enable_if_t<
                         std::is_member_function_pointer_v<decltype(&T::write)> &&
                         std::is_member_function_pointer_v<decltype(&T::bad)>>>
    : std::true_type {
};

/**
 * Wraps an std::ostream-like object. Requires the ostream-like object to have
 * methods void write(const char*, size_t) and bool bad(), and can optionally
 * provide void flush().
 */
class WritableStream {
 public:
  /**
   * @brief Construct a new WritableStream instance. The ostream-like object
   * will not owned by this instance.
   *
   * @tparam T the ostream-like type
   * @param stream a reference to an o-stream like object.
   */
  template <class T, std::enable_if_t<IsWritable<T>::value &&
                                          !std::is_same_v<T, WritableStream>,
                                      bool> = true>
  WritableStream(T& stream)
      : impl_(std::make_unique<WritableStreamT<T>>(&stream)) {
  }

  /**
   * @brief Construct a new WritableStream instance.
   * The ostream-like object will be owned by this instance.
   *
   * @tparam T the ostream-like type
   * @param stream a unique pointer to the ostream-like object.
   */
  template <class T, std::enable_if_t<IsWritable<T>::value, bool> = true>
  WritableStream(std::unique_ptr<T> stream)
      : impl_(std::make_unique<WritableStreamT<T>>(std::move(stream))) {
  }

  /**
   * @brief Construct a new WritableStream instance using an existing
   * shared pointer to a stream-like object.
   *
   * @tparam T the ostream-like type
   * @param stream a shared pointer to the ostream-like object.
   */
  template <class T, std::enable_if_t<IsWritable<T>::value, bool> = true>
  WritableStream(std::shared_ptr<T> stream)
      : impl_(std::make_unique<WritableStreamT<T>>(stream)) {
  }

  /**
   * @brief Construct a new WritableStream instance from a file. The file is closed
   * when this instance is destructed.
   *
   * @param filename the path to the file.
   */
  WritableStream(std::string const& filename) {
    auto file_stream = std::make_shared<std::ofstream>(
        filename, std::ios::binary | std::ios::out);
    if (!file_stream->good()) {
      throw std::runtime_error("Failed to open file for writing.");
    }
    impl_ = std::make_unique<WritableStreamT<std::ofstream>>(file_stream);
  }

  WritableStream(WritableStream const&) = delete;
  WritableStream& operator=(WritableStream const&) = delete;

  void Write(char const* buffer, size_t count) {
    impl_->Write(buffer, count);
  }

  void Flush() {
    impl_->Flush();
  }

  bool Bad() const {
    return impl_->Bad();
  }

 private:
  struct WritableStreamImpl {
    virtual ~WritableStreamImpl() = default;
    virtual void Write(char const* buffer, size_t count) = 0;
    virtual void Flush() {}
    virtual bool Bad() const = 0;
  };

  template <typename T, typename = void>
  struct HasFlush : std::false_type {
  };
  template <typename T>
  struct HasFlush<T, typename std::enable_if_t<
                         std::is_member_function_pointer<decltype(&T::flush)>::value>>
      : std::true_type {
  };

  template <class T>
  struct WritableStreamT : public WritableStreamImpl {
    using variant_type = std::variant<std::unique_ptr<T>, std::shared_ptr<T>, T*>;

    WritableStreamT(variant_type writable) : stream_{std::move(writable)} {}

    void Write(char const* buffer, size_t count) override {
      std::visit([buffer, count](auto&& arg) { arg->write(buffer, count); }, stream_);
    }

    void Flush() override {
      if constexpr (HasFlush<T>::value) {
        std::visit([](auto&& arg) { arg->flush(); }, stream_);
      }
    }

    bool Bad() const override {
      return std::visit([](auto&& arg) { return arg->bad(); }, stream_);
    }

    variant_type stream_;
  };

  std::unique_ptr<WritableStreamImpl> impl_;
};

template <typename T, typename = void>
struct IsReadable : std::false_type {
};

template <typename T>
struct IsReadable<T, typename std::enable_if_t<
                         std::is_member_function_pointer_v<decltype(&T::read)> &&
                         std::is_member_function_pointer_v<decltype(&T::eof)> &&
                         std::is_member_function_pointer_v<decltype(&T::gcount)>>>
    : std::true_type {
};

/**
 * Wraps an std::istream-like object. Requires the istream-like object to have
 * methods void read(char*, size_t), bool eof(), and size_t gcount().
 */
class ReadableStream {
 public:
  /**
   * @brief Construct a new ReadableStream instance. The istream-like object
   * will not owned by this instance.
   *
   * @tparam T the istream-like type
   * @param stream a reference to an i-stream like object.
   */
  template <class T, std::enable_if_t<IsReadable<T>::value &&
                                          !std::is_same_v<T, ReadableStream>,
                                      bool> = true>
  ReadableStream(T& stream)
      : impl_(std::make_unique<ReadableStreamT<T>>(&stream)) {
  }

  /**
   * @brief Construct a new ReadableStream instance. The istream-like object
   * will be owned by this instance.
   *
   * @tparam T the istream-like type
   * @param stream a unique pointer to the istream-like object.
   */
  template <class T, std::enable_if_t<IsReadable<T>::value, bool> = true>
  ReadableStream(std::unique_ptr<T> stream)
      : impl_(std::make_unique<ReadableStreamT<T>>(std::move(stream))) {
  }

  /**
   * @brief Construct a new ReadableStream instance using an existing
   * shared pointer to a stream-like object.
   *
   * @tparam T the istream-like type
   * @param stream a shared pointer to the istream-like object.
   */
  template <class T, std::enable_if_t<IsReadable<T>::value, bool> = true>
  ReadableStream(std::shared_ptr<T> stream)
      : impl_(std::make_unique<ReadableStreamT<T>>(stream)) {
  }

  /**
   * @brief Construct a new ReadableStream instance from a file.
   * The file is closed when this instance is destructed.
   *
   * @param filename the path to the file.
   */
  ReadableStream(std::string const& filename) {
    auto file_stream = std::make_shared<std::ifstream>(
        filename, std::ios::binary | std::ios::out);
    if (!file_stream->good()) {
      throw std::runtime_error("Failed to open file for reading.");
    }
    impl_ = std::make_unique<ReadableStreamT<std::ifstream>>(file_stream);
  }

  ReadableStream(ReadableStream const&) = delete;
  ReadableStream& operator=(ReadableStream const&) = delete;

  void Read(char* buffer, size_t count) {
    impl_->Read(buffer, count);
  }

  bool Eof() const {
    return impl_->Eof();
  }

  std::streamsize GCount() const {
    return impl_->GCount();
  }

 private:
  struct ReadableStreamImpl {
    virtual ~ReadableStreamImpl() = default;
    virtual void Read(char* buffer, size_t count) = 0;
    virtual bool Eof() const = 0;
    virtual std::streamsize GCount() const = 0;
  };

  template <class T>
  struct ReadableStreamT : public ReadableStreamImpl {
    using variant_type = std::variant<std::unique_ptr<T>, std::shared_ptr<T>, T*>;

    ReadableStreamT(variant_type readable) : stream_{std::move(readable)} {}

    void Read(char* buffer, size_t count) override {
      std::visit([buffer, count](auto&& arg) { arg->read(buffer, count); }, stream_);
    }

    bool Eof() const override {
      return std::visit([](auto&& arg) { return arg->eof(); }, stream_);
    }

    std::streamsize GCount() const override {
      return std::visit([](auto&& arg) { return arg->gcount(); }, stream_);
    }

    variant_type stream_;
  };

  std::unique_ptr<ReadableStreamImpl> impl_;
};
}  // namespace yardl::binary
