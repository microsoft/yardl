// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <fstream>

#include <nlohmann/json.hpp>

#include "serializers.h"

namespace yardl::ndjson {
using json = nlohmann::json;

class NDJsonWriter {
 protected:
  NDJsonWriter(std::ostream& stream, [[maybe_unused]]std::string& schema)
      : stream_(stream) {
    //WriteHeader(stream_, schema);
  }

  NDJsonWriter(std::string file_name, [[maybe_unused]]std::string& schema) : owned_file_stream_(open_file(file_name)), stream_(*owned_file_stream_) {
    //WriteHeader(stream_, schema);
  }

 private:
  static std::unique_ptr<std::ofstream> open_file(std::string filename) {
    auto file_stream = std::make_unique<std::ofstream>(filename, std::ios::out);
    if (!file_stream->good()) {
      throw std::runtime_error("Failed to open file for writing.");
    }

    return file_stream;
  }

 private:
  std::unique_ptr<std::ofstream> owned_file_stream_{};
 protected:
  std::ostream& stream_;
};

class NDJsonReader {
 protected:
  NDJsonReader(std::istream& stream, [[maybe_unused]]std::string& schema)
      : stream_(stream) {
    //ReadHeader(stream_, schema);
  }

  NDJsonReader(std::string file_name, [[maybe_unused]]std::string& schema) : owned_file_stream_(open_file(file_name)), stream_(*owned_file_stream_) {
    //ReadHeader(stream_, schema);
  }

 private:
  static std::unique_ptr<std::ifstream> open_file(std::string filename) {
    auto file_stream = std::make_unique<std::ifstream>(filename, std::ios::in);
    if (!file_stream->good()) {
      throw std::runtime_error("Failed to open file for reading.");
    }

    return file_stream;
  }

 private:
  std::unique_ptr<std::ifstream> owned_file_stream_{};
 protected:
  std::istream& stream_;
};

}  // namespace yardl::ndjson
