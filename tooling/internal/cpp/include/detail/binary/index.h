// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once

#include "serializers.h"

namespace yardl::binary {

class Index {
 public:
  Index() = default;
  ~Index() = default;
  Index(Index const&) = default;
  Index& operator=(Index const&) = default;

  void set_step_offset(std::string const& step_name, size_t offset) {
    // Don't overwrite the start offset if it already exists.
    if (step_offsets.count(step_name) == 0) {
      step_offsets[step_name] = offset;
    }
  }

  void add_stream_offset(std::string const& step_name, size_t offset) {
    // Save the block start index.
    stream_blocks[step_name].push_back(stream_offsets[step_name].size());

    // Save the item offset.
    auto& dest = stream_offsets[step_name];
    dest.push_back(offset);
  }

  void add_stream_offsets(std::string const& step_name, std::vector<size_t> const& offsets) {
    // Save the block start index.
    stream_blocks[step_name].push_back(stream_offsets[step_name].size());

    auto& dest = stream_offsets[step_name];
    dest.insert(dest.end(), offsets.begin(), offsets.end());
  }

  size_t get_step_offset(std::string const& step_name) const {
    return step_offsets.at(step_name);
  }

  size_t get_stream_size(std::string const& step_name) const {
    return stream_offsets.at(step_name).size();
  }

  // Given the name of a Stream ProtocolStep and an index into the stream, sets the
  //  1. absolute byte offset into the stream, and
  //  2. number of items remaining in corresponding stream block
  //  Returns true on success, false on error, including index out of range
  bool find_stream_item(std::string const& step_name, size_t index, size_t& absolute_offset, size_t& items_remaining) const {
    if (stream_offsets.count(step_name) == 0) {
      if (index == 0) {
        absolute_offset = step_offsets.at(step_name);
        items_remaining = 0;
        return true;
      }
      return false;
    }

    auto& offsets = stream_offsets.at(step_name);
    if (index >= offsets.size()) {
      // Index out-of-bounds
      return false;
    }
    absolute_offset = offsets.at(index);

    auto& blocks = stream_blocks.at(step_name);
    size_t last_block = 0;
    for (auto const& block : blocks) {
      if (block > index) {
        last_block = block;
        break;
      }
    }
    if (last_block == 0) {
      last_block = offsets.size();
    }
    items_remaining = last_block - index;

    return true;
  }

  void dump() const {
    for (auto const& [k, v] : step_offsets) {
      std::cerr << "Step " << k << " at offset " << v << std::endl;

      if (stream_offsets.count(k) > 0) {
        for (auto const& offset : stream_offsets.at(k)) {
          std::cerr << "  Stream item offset " << offset << std::endl;
        }

        for (auto const& block : stream_blocks.at(k)) {
          std::cerr << "  Stream block at " << block << std::endl;
        }

        for (size_t i = 0; i < stream_offsets.at(k).size(); i++) {
          size_t offset, remaining;
          if (!find_stream_item(k, i, offset, remaining)) {
            std::cerr << "Failed to find item " << i << " in stream " << k << std::endl;
            continue;
          }
          std::cerr << "  Stream item " << i << " is at offset " << offset << " with " << remaining << " items remaining in block" << std::endl;
        }
      }
    }
  }

  //  protected:
  std::unordered_map<std::string, size_t> step_offsets;
  std::unordered_map<std::string, std::vector<size_t>> stream_offsets;
  std::unordered_map<std::string, std::vector<size_t>> stream_blocks;
};

static inline std::array<char, 10> INDEX_MAGIC_BYTES = {'y', 'a', 'r', 'd', 'l', 'i', 'n', 'd', 'e', 'x'};
static inline uint32_t kBinaryIndexFormatVersionNumber = 1;

inline Index ReadIndex(CodedInputStream& stream) {
  auto pos = stream.Pos();

  size_t index_offset = 0;
  stream.Seek(-sizeof(index_offset));
  stream.ReadFixedInteger(index_offset);
  try {
    stream.Seek(index_offset);
  } catch (std::exception const& e) {
    throw std::runtime_error("Binary Index not found in stream.");
  }

  std::array<char, sizeof(INDEX_MAGIC_BYTES)> magic_bytes{};
  stream.ReadBytes(magic_bytes.data(), magic_bytes.size());
  if (magic_bytes != INDEX_MAGIC_BYTES) {
    throw std::runtime_error("Binary Index in the stream is not in the expected format.");
  }

  uint32_t version_number;
  stream.ReadFixedInteger(version_number);
  if (version_number != kBinaryIndexFormatVersionNumber) {
    throw std::runtime_error(
        "Binary Index in the stream is not in the expected format. Unsupported version.");
  }

  Index index;
  ReadMap<std::string, size_t, ReadString, ReadInteger>(stream, index.step_offsets);
  ReadMap<std::string, std::vector<size_t>, ReadString, ReadVector<size_t, ReadInteger>>(stream, index.stream_offsets);
  ReadMap<std::string, std::vector<size_t>, ReadString, ReadVector<size_t, ReadInteger>>(stream, index.stream_blocks);

  stream.Seek(pos);

  return index;
}

inline void WriteIndex(CodedOutputStream& stream, Index const& index) {
  size_t pos = stream.Pos();

  stream.WriteBytes(INDEX_MAGIC_BYTES.data(), INDEX_MAGIC_BYTES.size());
  stream.WriteFixedInteger(kBinaryIndexFormatVersionNumber);

  yardl::binary::WriteMap<std::string, size_t, yardl::binary::WriteString, yardl::binary::WriteInteger<size_t>>(stream, index.step_offsets);
  yardl::binary::WriteMap<std::string, std::vector<size_t>, yardl::binary::WriteString, yardl::binary::WriteVector<size_t, yardl::binary::WriteInteger<size_t>>>(stream, index.stream_offsets);
  yardl::binary::WriteMap<std::string, std::vector<size_t>, yardl::binary::WriteString, yardl::binary::WriteVector<size_t, yardl::binary::WriteInteger<size_t>>>(stream, index.stream_blocks);

  stream.WriteFixedInteger(pos);
}

}  // namespace yardl::binary
