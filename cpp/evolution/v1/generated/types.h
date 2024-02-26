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
using AliasedInt = int32_t;

using AliasedString = std::string;

using AliasedLongToString = std::string;

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
  double float_to_double{};
  evo_test::UnchangedRecord unchanged_record{};
  int64_t int_to_long{};
  std::optional<std::string> optional_long_to_string{};

  bool operator==(const RecordWithChanges& other) const {
    return float_to_double == other.float_to_double &&
      unchanged_record == other.unchanged_record &&
      int_to_long == other.int_to_long &&
      optional_long_to_string == other.optional_long_to_string;
  }

  bool operator!=(const RecordWithChanges& other) const {
    return !(*this == other);
  }
};

using AliasedRecordWithChanges = evo_test::RecordWithChanges;

using AliasOfAliasedRecordWithChanges = evo_test::AliasedRecordWithChanges;

using AliasedOptionalRecord = std::optional<evo_test::RecordWithChanges>;

using AliasedOptionalString = std::optional<std::string>;

using AliasedRecordOrInt = std::variant<evo_test::RecordWithChanges, int32_t>;

using AliasedRecordOrString = std::variant<evo_test::RecordWithChanges, std::string>;

struct DeprecatedRecord {
  std::string s{};
  int32_t i{};

  bool operator==(const DeprecatedRecord& other) const {
    return s == other.s &&
      i == other.i;
  }

  bool operator!=(const DeprecatedRecord& other) const {
    return !(*this == other);
  }
};

using RenamedRecord = evo_test::DeprecatedRecord;

using StreamItem = std::variant<evo_test::RecordWithChanges, evo_test::RenamedRecord>;

struct RZ {
  int32_t subject{};

  bool operator==(const RZ& other) const {
    return subject == other.subject;
  }

  bool operator!=(const RZ& other) const {
    return !(*this == other);
  }
};

using RY = evo_test::RZ;

using RNew = evo_test::RY;

using RLink = evo_test::RNew;

using RX = evo_test::RLink;

using RUnion = std::variant<evo_test::RX, std::string>;

struct UnusedButChangedRecord {
  float age{};
  std::string name{};

  bool operator==(const UnusedButChangedRecord& other) const {
    return age == other.age &&
      name == other.name;
  }

  bool operator!=(const UnusedButChangedRecord& other) const {
    return !(*this == other);
  }
};

template <typename T1, typename T2>
struct GenericRecord {
  T2 field_2{};
  T1 field_1{};
  std::optional<bool> added{};

  bool operator==(const GenericRecord& other) const {
    return field_2 == other.field_2 &&
      field_1 == other.field_1 &&
      added == other.added;
  }

  bool operator!=(const GenericRecord& other) const {
    return !(*this == other);
  }
};

template <typename A, typename B>
using AliasedOpenGenericRecord = evo_test::GenericRecord<A, B>;

template <typename T>
using AliasedHalfClosedGenericRecord = evo_test::AliasedOpenGenericRecord<T, std::string>;

template <typename X, typename Y>
using GenericUnion2 = std::variant<X, Y>;

template <typename T1, typename T2>
using GenericUnion = evo_test::GenericUnion2<T1, T2>;

template <typename A, typename B>
using AliasedOpenGenericUnion = evo_test::GenericUnion<A, B>;

template <typename T>
using AliasedHalfClosedGenericUnion = evo_test::AliasedOpenGenericUnion<T, float>;

using AliasedClosedGenericUnion = evo_test::AliasedHalfClosedGenericUnion<evo_test::GenericRecord<int32_t, std::string>>;

template <typename T>
struct GenericParentRecord {
  evo_test::AliasedHalfClosedGenericRecord<T> record{};
  evo_test::AliasedOpenGenericRecord<evo_test::AliasedOpenGenericUnion<T, float>, std::string> record_of_union{};
  evo_test::AliasedClosedGenericUnion union_of_record{};

  bool operator==(const GenericParentRecord& other) const {
    return record == other.record &&
      record_of_union == other.record_of_union &&
      union_of_record == other.union_of_record;
  }

  bool operator!=(const GenericParentRecord& other) const {
    return !(*this == other);
  }
};

template <typename T, typename U>
using GenericUnionReversed = evo_test::GenericUnion2<U, T>;

using AliasedClosedGenericRecord = evo_test::AliasedHalfClosedGenericRecord<int32_t>;

template <typename X, typename Y>
using GenericRecordReversed = evo_test::GenericRecord<Y, X>;

using AliasedGenericRecordOrString = std::variant<evo_test::GenericRecord<int32_t, std::string>, std::string>;

template <typename T2>
struct OldUnchangedGeneric {
  T2 field{};

  bool operator==(const OldUnchangedGeneric& other) const {
    return field == other.field;
  }

  bool operator!=(const OldUnchangedGeneric& other) const {
    return !(*this == other);
  }
};

template <typename A>
using UnchangedGeneric = evo_test::OldUnchangedGeneric<A>;

using Unchanged = evo_test::UnchangedGeneric<int32_t>;

template <typename Y, typename Z>
struct OldChangedGeneric {
  std::optional<Y> y{};
  std::optional<evo_test::OldUnchangedGeneric<Z>> z{};

  bool operator==(const OldChangedGeneric& other) const {
    return y == other.y &&
      z == other.z;
  }

  bool operator!=(const OldChangedGeneric& other) const {
    return !(*this == other);
  }
};

template <typename I, typename J>
using ChangedGeneric = evo_test::OldChangedGeneric<I, J>;

using Changed = evo_test::ChangedGeneric<std::string, int32_t>;

enum class GrowingEnum : uint16_t {
  kA = 0,
  kB = 1,
  kC = 2,
  kD = 3,
};

using AliasedEnum = evo_test::GrowingEnum;

// Compatibility aliases for version v0.

using AliasedLongToString_v0 = evo_test::AliasedLongToString;

using UnchangedRecord_v0 = evo_test::UnchangedRecord;

using RecordWithChanges_v0 = evo_test::RecordWithChanges;

using AliasedRecordWithChanges_v0 = evo_test::AliasedRecordWithChanges;

using RenamedRecord_v0 = evo_test::DeprecatedRecord;

using StreamItem_v0 = evo_test::StreamItem;

using RC_v0 = evo_test::RZ;

using RB_v0 = evo_test::RZ;

using RA_v0 = evo_test::RZ;

using RLink_v0 = evo_test::RLink;

template <typename T1, typename T2>
using GenericRecord_v0 = evo_test::GenericRecord<T1, T2>;

template <typename T1, typename T2>
using GenericUnion_v0 = evo_test::GenericUnion<T1, T2>;

template <typename T>
using GenericParentRecord_v0 = evo_test::GenericParentRecord<T>;

template <typename T>
using AliasedHalfClosedGenericUnion_v0 = evo_test::GenericUnion<T, float>;

using AliasedClosedGenericUnion_v0 = evo_test::AliasedClosedGenericUnion;

template <typename T>
using AliasedHalfClosedGenericRecord_v0 = evo_test::AliasedHalfClosedGenericRecord<T>;

template <typename T2>
using UnchangedGeneric_v0 = evo_test::OldUnchangedGeneric<T2>;

template <typename Y, typename Z>
using ChangedGeneric_v0 = evo_test::OldChangedGeneric<Y, Z>;

using GrowingEnum_v0 = evo_test::GrowingEnum;

} // namespace evo_test
