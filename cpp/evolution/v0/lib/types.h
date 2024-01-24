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
using AliasedLongToString = int64_t;

struct UnchangedRecord {
  std::string name{};
  int32_t age{};
  std::unordered_map<std::string, double> meta{};

  bool operator==(const UnchangedRecord& other) const {
    return name == other.name &&
      age == other.age &&
      meta == other.meta;
  }

  bool operator!=(const UnchangedRecord& other) const {
    return !(*this == other);
  }
};

struct RecordWithChanges {
  int32_t int_to_long{};
  std::vector<int32_t> deprecated_vector{};
  float float_to_double{};
  yardl::FixedNDArray<uint8_t, 7> deprecated_array{};
  std::optional<int64_t> optional_long_to_string{};
  std::unordered_map<std::string, std::vector<int32_t>> deprecated_map{};
  evo_test::UnchangedRecord unchanged_record{};

  bool operator==(const RecordWithChanges& other) const {
    return int_to_long == other.int_to_long &&
      deprecated_vector == other.deprecated_vector &&
      float_to_double == other.float_to_double &&
      deprecated_array == other.deprecated_array &&
      optional_long_to_string == other.optional_long_to_string &&
      deprecated_map == other.deprecated_map &&
      unchanged_record == other.unchanged_record;
  }

  bool operator!=(const RecordWithChanges& other) const {
    return !(*this == other);
  }
};

using AliasedRecordWithChanges = evo_test::RecordWithChanges;

struct RenamedRecord {
  int32_t i{};
  std::string s{};

  bool operator==(const RenamedRecord& other) const {
    return i == other.i &&
      s == other.s;
  }

  bool operator!=(const RenamedRecord& other) const {
    return !(*this == other);
  }
};

struct RC {
  std::string subject{};

  bool operator==(const RC& other) const {
    return subject == other.subject;
  }

  bool operator!=(const RC& other) const {
    return !(*this == other);
  }
};

using RB = evo_test::RC;

using RA = evo_test::RB;

using RLink = evo_test::RA;

struct UnusedButChangedRecord {
  std::string name{};
  int32_t age{};

  bool operator==(const UnusedButChangedRecord& other) const {
    return name == other.name &&
      age == other.age;
  }

  bool operator!=(const UnusedButChangedRecord& other) const {
    return !(*this == other);
  }
};

} // namespace evo_test

