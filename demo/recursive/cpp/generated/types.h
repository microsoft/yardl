// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include <array>
#include <complex>
#include <optional>
#include <unordered_map>
#include <variant>
#include <vector>

#include "yardl/yardl.h"

namespace sketch {
struct BinaryTree;

template <typename T>
struct LinkedList;

struct Directory;

struct BinaryTree {
  int32_t value{};
  std::unique_ptr<sketch::BinaryTree> left{};
  std::unique_ptr<sketch::BinaryTree> right{};

  bool operator==(const BinaryTree& other) const {
    return value == other.value &&
      left == other.left &&
      right == other.right;
  }

  bool operator!=(const BinaryTree& other) const {
    return !(*this == other);
  }
};

template <typename T>
struct LinkedList {
  T value{};
  std::unique_ptr<sketch::LinkedList<T>> next{};

  bool operator==(const LinkedList& other) const {
    return value == other.value &&
      next == other.next;
  }

  bool operator!=(const LinkedList& other) const {
    return !(*this == other);
  }
};

struct File {
  std::string name{};
  std::vector<uint8_t> data{};

  bool operator==(const File& other) const {
    return name == other.name &&
      data == other.data;
  }

  bool operator!=(const File& other) const {
    return !(*this == other);
  }
};

using DirectoryEntry = std::variant<sketch::File, std::unique_ptr<sketch::Directory>>;

struct Directory {
  std::string name{};
  std::vector<sketch::DirectoryEntry> entries{};

  bool operator==(const Directory& other) const {
    return name == other.name &&
      entries == other.entries;
  }

  bool operator!=(const Directory& other) const {
    return !(*this == other);
  }
};

} // namespace sketch

