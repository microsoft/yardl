// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <fstream>

#include <nlohmann/json.hpp>

namespace yardl::ndjson {
using ordered_json = nlohmann::ordered_json;

static inline uint32_t kNDJsonFormatVersionNumber = 1;

class NDJsonWriter {
 protected:
  NDJsonWriter(std::ostream& stream, std::string& schema)
      : stream_(stream) {
    WriteHeader(schema);
  }

  NDJsonWriter(std::string file_name, std::string& schema)
      : owned_file_stream_(open_file(file_name)), stream_(*owned_file_stream_) {
    WriteHeader(schema);
  }

 private:
  static std::unique_ptr<std::ofstream> open_file(std::string filename) {
    auto file_stream = std::make_unique<std::ofstream>(filename, std::ios::out);
    if (!file_stream->good()) {
      throw std::runtime_error("Failed to open file for writing.");
    }

    return file_stream;
  }

  void WriteHeader(std::string& schema) {
    auto parsed_schema = ordered_json::parse(schema);

    ordered_json metadata = {{"yardl", {{"version", kNDJsonFormatVersionNumber}, {"schema", parsed_schema}}}};
    stream_ << metadata << "\n";
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
    ReadHeader(schema);
  }

  NDJsonReader(std::string file_name, std::string& schema)
      : owned_file_stream_(open_file(file_name)), stream_(*owned_file_stream_) {
    ReadHeader(schema);
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

  void ReadHeader(std::string& expected_schema) {
    ordered_json expected_schema_json = ordered_json::parse(expected_schema);
    std::string line;
    std::getline(stream_, line);
    try {
      ordered_json actual_header_json = ordered_json::parse(line);
      actual_header_json = actual_header_json.at("yardl");
      if (actual_header_json["version"] != kNDJsonFormatVersionNumber) {
        throw std::runtime_error(
            "Unsupported Yardl NDJSON format version.");
      }
      if (expected_schema_json != actual_header_json.at("schema")) {
        throw std::runtime_error(
            "The schema of the data to be read is not compatible with the current protocol.");
      }
    } catch (ordered_json::exception const&) {
      throw std::runtime_error(
          "Data in the stream is not in the expected Yardl NDJSON format.");
    }
  }

 private:
  std::unique_ptr<std::ifstream> owned_file_stream_{};

 protected:
  std::istream& stream_;
  std::string line_{};
  std::optional<ordered_json> unused_step_{};
};

}  // namespace yardl::ndjson
