// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include "header.h"

namespace yardl::binary {
class BinaryWriter {
 protected:
  // The stream_arg parameter can either be a std::string filename
  // or a reference, std::unique_ptr, or std::shared_ptr to a stream-like object, such as std::ostream.
  template <typename TStreamArg>
  BinaryWriter(TStreamArg&& stream_arg, std::string& schema)
      : stream_(std::forward<TStreamArg>(stream_arg)) {
    WriteHeader(stream_, schema);
  }

  yardl::binary::CodedOutputStream stream_;
};

class BinaryReader {
 protected:
  // The stream_arg parameter can either be a std::string filename
  // or a reference, std::unique_ptr, or std::shared_ptr to a stream-like object, such as std::istream.
  template <typename TStreamArg>
  BinaryReader(TStreamArg&& stream_arg, std::string& schema)
      : stream_(std::forward<TStreamArg>(stream_arg)) {
    ReadHeader(stream_, schema);
  }

 protected:
  yardl::binary::CodedInputStream stream_;
};

}  // namespace yardl::binary
