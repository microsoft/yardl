// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <fstream>

#include <nlohmann/json.hpp>

#include "header.h"

namespace yardl::ndjson {
using ordered_json = nlohmann::ordered_json;

class NDJsonWriter {
 protected:
  NDJsonWriter(std::ostream& stream, std::string& schema)
      : stream_(stream) {
    WriteHeader(stream_, schema);
  }

  NDJsonWriter(std::string file_name, std::string& schema)
      : owned_file_stream_(open_file(file_name)), stream_(*owned_file_stream_) {
    WriteHeader(stream_, schema);
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
  NDJsonReader(std::istream& stream, std::string& schema)
      : stream_(stream) {
    ReadAndValidateHeader(stream_, schema);
  }

  NDJsonReader(std::string file_name, std::string& schema)
      : owned_file_stream_(open_file(file_name)), stream_(*owned_file_stream_) {
    ReadAndValidateHeader(stream_, schema);
  }

  void VerifyFinished() {
    if (unused_step_ || stream_.peek() != EOF) {
      throw std::runtime_error("The stream was not read to completion.");
    }
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
  std::string line_{};
  std::optional<ordered_json> unused_step_{};
};

}  // namespace yardl::ndjson
