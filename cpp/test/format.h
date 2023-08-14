#pragma once

#include <iostream>

namespace yardl::testing {

enum class Format {
  kBinary,
  kHdf5,
  kNDJson,
};

inline std::ostream& operator<<(std::ostream& os, Format const& format) {
  switch (format) {
    case Format::kBinary:
      return os << "Binary";
    case Format::kHdf5:
      return os << "HDF5";
    case Format::kNDJson:
      return os << "NDJson";
    default:
      return os << std::to_string(static_cast<int>(format));
  }
}
}  // namespace yardl::testing
