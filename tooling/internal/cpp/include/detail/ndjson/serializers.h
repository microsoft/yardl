// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <complex>
#include <memory>
#include <optional>
#include <unordered_map>
#include <utility>
#include <vector>

#include <nlohmann/json.hpp>

#include "../../yardl.h"

namespace yardl::ndjson {

using ordered_json = nlohmann::ordered_json;

inline void WriteProtocolValue(std::ostream& stream, std::string const& stepName, ordered_json const& value) {
  stream << ordered_json{{stepName, value}} << '\n';
}

template <typename T>
inline bool ReadProtocolValue(std::istream& stream, std::string& line, std::string const& stepName, bool required, std::optional<ordered_json>& unused_step, T& value) {
  ordered_json* json_value;
  if (unused_step) {
    try {
      ordered_json& v = unused_step->at(stepName);
      json_value = &v;
    } catch (ordered_json::out_of_range&) {
      if (required) {
        unused_step.reset();
        throw std::runtime_error("encountered unexpected protocol value");
      }
      return false;
    }

    json_value->get_to(value);
    unused_step.reset();
    return true;
  }

  if (!std::getline(stream, line)) {
    if (!required) {
      return false;
    }
    throw std::runtime_error("missing protocol step " + stepName);
  }

  ordered_json parsed_step = ordered_json::parse(line);
  try {
    ordered_json& v = parsed_step.at(stepName);
    json_value = &v;
  } catch (ordered_json::out_of_range&) {
    if (required) {
      throw std::runtime_error("encountered unexpected protocol value");
    }
    unused_step.emplace(std::move(parsed_step));
    return false;
  }

  json_value->get_to(value);
  return true;
}

}  // namespace yardl::ndjson

NLOHMANN_JSON_NAMESPACE_BEGIN

template <typename T>
struct adl_serializer<std::optional<T>> {
  static void to_json(ordered_json& j, std::optional<T> const& opt) {
    if (opt) {
      j = *opt;
    } else {
      j = nullptr;
    }
  }

  static void from_json(ordered_json const& j, std::optional<T>& opt) {
    if (j.is_null()) {
      opt = std::nullopt;
    } else {
      opt = j.get<T>();
    }
  }
};

NLOHMANN_JSON_NAMESPACE_END
