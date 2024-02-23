// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <fstream>
#include <memory>

#include "header.h"

namespace yardl::binary {
class BinaryWriter {
 protected:
  BinaryWriter(std::ostream& stream, std::string const& schema)
      : stream_(stream) {
    WriteHeader(stream_, schema);
  }

  BinaryWriter(std::string file_name, std::string const& schema)
      : owned_file_stream_(open_file(file_name)), stream_(*owned_file_stream_) {
    WriteHeader(stream_, schema);
  }

 private:
  static std::unique_ptr<std::ofstream> open_file(std::string filename) {
    auto file_stream = std::make_unique<std::ofstream>(filename, std::ios::binary | std::ios::out);
    if (!file_stream->good()) {
      throw std::runtime_error("Failed to open file for writing.");
    }

    return file_stream;
  }

 private:
  std::unique_ptr<std::ofstream> owned_file_stream_{};

 protected:
  yardl::binary::CodedOutputStream stream_;
};

class BinaryReader {
 protected:
  BinaryReader(std::istream& stream)
      : stream_(stream) {
    schema_read_ = ReadHeader(stream_);
  }

  BinaryReader(std::string file_name)
      : owned_file_stream_(open_file(file_name)), stream_(*owned_file_stream_) {
    schema_read_ = ReadHeader(stream_);
  }

 private:
  static std::unique_ptr<std::ifstream> open_file(std::string filename) {
    auto file_stream = std::make_unique<std::ifstream>(filename, std::ios::binary | std::ios::in);
    if (!file_stream->good()) {
      throw std::runtime_error("Failed to open file for reading.");
    }

    return file_stream;
  }

 private:
  std::unique_ptr<std::ifstream> owned_file_stream_{};

 protected:
  yardl::binary::CodedInputStream stream_;
  std::string schema_read_;
};

}  // namespace yardl::binary
