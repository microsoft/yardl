#pragma once

#include <algorithm>
#include <cctype>
#include <iostream>
#include <string>

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

inline Format ParseFormat(std::string format) {
  std::transform(format.begin(), format.end(), format.begin(), [](unsigned char c) { return std::tolower(c); });
  if (format == "binary") {
    return Format::kBinary;
  } else if (format == "hdf5") {
    return Format::kHdf5;
  } else if (format == "ndjson") {
    return Format::kNDJson;
  }
  throw std::runtime_error("Unknown format: " + format);
}

}  // namespace yardl::testing
