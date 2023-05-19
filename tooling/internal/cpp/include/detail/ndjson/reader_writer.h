// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <nlohmann/json.hpp>

#include "../stream/stream.h"
#include "serializers.h"

namespace yardl::ndjson {
using json = nlohmann::json;

// static inline uint32_t kJsonFormatVersionNumber = 1;
class NDJsonWriter {
 protected:
  // The stream_arg parameter can either be a std::string filename
  // or a reference, std::unique_ptr, or std::shared_ptr to a stream-like object, such as std::ostream.
  template <typename TStreamArg>
  NDJsonWriter(TStreamArg&& stream_arg, std::string& schema)
      : stream_(std::forward<TStreamArg>(stream_arg)) {
    [[maybe_unused]] auto j = json::parse(schema);

    // WriteStartObject(stream_);
    // WriteFieldName(stream_, "yardl");
    // WriteStartObject(stream_);
    // WriteFieldName(stream_, "version");
    // WriteInteger(stream_, kJsonFormatVersionNumber);
    // WriteComma(stream_);
    // WriteFieldName(stream_, "schema");
    // stream_.Write(schema.data(), schema.size());
    // WriteEndObject(stream_);
    // WriteEndObject(stream_);
  }

  yardl::stream::WritableStream stream_;
};

class NDJsonReader {
 protected:
  // The stream_arg parameter can either be a std::string filename
  // or a reference, std::unique_ptr, or std::shared_ptr to a stream-like object, such as std::istream.
  template <typename TStreamArg>
  NDJsonReader(TStreamArg&& stream_arg, std::string& schema)
      : stream_(std::forward<TStreamArg>(stream_arg)) {
    // ReadHeader(stream_, schema);
  }

 protected:
  yardl::stream::ReadableStream stream_;
};

}  // namespace yardl::ndjson
