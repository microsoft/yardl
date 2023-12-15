// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include <array>
#include <complex>
#include <optional>
#include <unordered_map>
#include <variant>
#include <vector>

#include "yardl/yardl.h"

namespace evo_test {
struct Header {
  std::string subject{};
  int64_t weight{};
  std::unordered_map<std::string, std::vector<std::string>> meta{};

  bool operator==(const Header& other) const {
    return subject == other.subject &&
      weight == other.weight &&
      meta == other.meta;
  }

  bool operator!=(const Header& other) const {
    return !(*this == other);
  }
};

struct Sample {
  std::vector<int32_t> data{};
  yardl::DateTime timestamp{};

  bool operator==(const Sample& other) const {
    return data == other.data &&
      timestamp == other.timestamp;
  }

  bool operator!=(const Sample& other) const {
    return !(*this == other);
  }
};

struct Signature {
  std::string name{};
  std::string email{};
  int64_t number{};

  bool operator==(const Signature& other) const {
    return name == other.name &&
      email == other.email &&
      number == other.number;
  }

  bool operator!=(const Signature& other) const {
    return !(*this == other);
  }
};

struct Footer {
  evo_test::Signature signature{};

  bool operator==(const Footer& other) const {
    return signature == other.signature;
  }

  bool operator!=(const Footer& other) const {
    return !(*this == other);
  }
};

struct UnusedRecord {
  std::string subject{};

  bool operator==(const UnusedRecord& other) const {
    return subject == other.subject;
  }

  bool operator!=(const UnusedRecord& other) const {
    return !(*this == other);
  }
};

} // namespace evo_test

