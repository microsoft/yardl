// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include <complex>
#include <memory>
#include <optional>
#include <unordered_map>
#include <utility>
#include <variant>
#include <vector>
#pragma GCC diagnostic ignored "-Wstringop-overflow"
#include <date/date.h>
#pragma GCC diagnostic pop
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

template <typename T>
constexpr bool ShouldSerializeFieldValue(T const&) {
  return true;
}

template <typename T>
constexpr bool ShouldSerializeFieldValue(std::optional<T> const& value) {
  return value.has_value();
}

template <typename... Ts>
constexpr bool ShouldSerializeFieldValue(std::variant<Ts...> const& var) {
  if constexpr (std::is_same_v<std::monostate, std::variant_alternative_t<0, std::variant<Ts...>>>) {
    return var.index() != 0;
  }

  return true;
}

}  // namespace yardl::ndjson

NLOHMANN_JSON_NAMESPACE_BEGIN

template <typename T>
struct adl_serializer<std::optional<T>> {
  static void to_json(ordered_json& j, std::optional<T> const& value) {
    if (value) {
      j = *value;
    } else {
      j = nullptr;
    }
  }

  static void from_json(ordered_json const& j, std::optional<T>& value) {
    if (j.is_null()) {
      value = std::nullopt;
    } else {
      value = j.get<T>();
    }
  }
};

template <>
struct adl_serializer<std::monostate> {
  static void to_json(ordered_json& j, [[maybe_unused]] std::monostate const& value) {
    j = nullptr;
  }

  static void from_json(ordered_json const& j, [[maybe_unused]] std::monostate& value) {
    if (!j.is_null()) {
      throw std::runtime_error("expected null");
    }
  }
};

template <>
struct adl_serializer<yardl::Date> {
  static void to_json(ordered_json& j, yardl::Date const& value) {
#if __cplusplus < 202002L
    j = date::format("%F", value);
#else
    date::local_days ld(value.time_since_epoch());
    j = date::format("%F", ld);
#endif
  }

  static void from_json(ordered_json const& j, yardl::Date& value) {
    std::stringstream ss{j.get<std::string>()};

#if __cplusplus < 202002L
    ss >> date::parse("%F", value);
#else
    date::local_days ld;
    ss >> date::parse("%F", ld);
    value = yardl::Date(ld.time_since_epoch());
#endif

    if (ss.fail()) {
      throw std::runtime_error("invalid date format");
    }
  }
};

template <>
struct adl_serializer<yardl::Time> {
  static void to_json(ordered_json& j, yardl::Time const& value) {
    j = date::format("%T", value);
  }

  static void from_json(ordered_json const& j, yardl::Time& value) {
    std::stringstream ss{j.get<std::string>()};
    ss >> date::parse("%T", value);
    if (ss.fail()) {
      throw std::runtime_error("invalid time format");
    }
  }
};

template <>
struct adl_serializer<yardl::DateTime> {
  static void to_json(ordered_json& j, yardl::DateTime const& value) {
    j = date::format("%FT%T", value);
  }

  static void from_json(ordered_json const& j, yardl::DateTime& value) {
    std::stringstream ss{j.get<std::string>()};
    ss >> date::parse("%FT%T", value);
    if (ss.fail()) {
      throw std::runtime_error("invalid datetime format");
    }
  }
};

template <typename T>
struct adl_serializer<std::complex<T>> {
  static void to_json(ordered_json& j, std::complex<T> const& value) {
    j = ordered_json::array({value.real(), value.imag()});
  }

  static void from_json(ordered_json const& j, std::complex<T>& value) {
    value = std::complex<T>{j.at(0).get<T>(), j.at(1).get<T>()};
  }
};

template <typename T>
struct adl_serializer<yardl::DynamicNDArray<T>> {
  static void to_json(ordered_json& j, yardl::DynamicNDArray<T> const& value) {
    auto shape = value.shape();
    auto data_array = ordered_json::array();
    for (auto const& v : value) {
      data_array.push_back(v);
    }
    j = ordered_json{{"shape", shape}, {"data", data_array}};
  }

  static void from_json(ordered_json const& j, yardl::DynamicNDArray<T>& value) {
    value.resize(j.at("shape").get<std::vector<size_t>>());
    auto data_array = j.at("data").get<std::vector<T>>();
    size_t i = 0;
    for (auto& element : value) {
      element = data_array[i++];
    }
  }
};

template <typename T, size_t N>
struct adl_serializer<yardl::NDArray<T, N>> {
  static void to_json(ordered_json& j, yardl::NDArray<T, N> const& value) {
    auto shape = value.shape();
    auto data_array = ordered_json::array();
    for (auto const& v : value) {
      data_array.push_back(v);
    }
    j = ordered_json{{"shape", shape}, {"data", data_array}};
  }

  static void from_json(ordered_json const& j, yardl::NDArray<T, N>& value) {
    value.resize(j.at("shape").get<std::vector<size_t>>());
    auto data_array = j.at("data").get<std::vector<T>>();
    size_t i = 0;
    for (auto& element : value) {
      element = data_array[i++];
    }
  }
};

template <typename T, size_t... Dims>
struct adl_serializer<yardl::FixedNDArray<T, Dims...>> {
  static void to_json(ordered_json& j, yardl::FixedNDArray<T, Dims...> const& value) {
    auto data_array = ordered_json::array();
    for (auto const& v : value) {
      data_array.push_back(v);
    }
    j = data_array;
  }

  static void from_json(ordered_json const& j, yardl::FixedNDArray<T, Dims...>& value) {
    auto data_array = j.get<std::vector<T>>();
    size_t i = 0;
    for (auto& element : value) {
      element = data_array[i++];
    }
  }
};

NLOHMANN_JSON_NAMESPACE_END
