// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include <array>
#include <complex>
#include <memory>
#include <optional>
#include <variant>
#include <vector>

#include "../yardl/detail/ndjson/reader_writer.h"
#include "../protocols.h"

namespace sketch::ndjson {
// NDJSON writer for the MyProtocol protocol.
class MyProtocolWriter : public sketch::MyProtocolWriterBase, yardl::ndjson::NDJsonWriter {
  public:
  MyProtocolWriter(std::ostream& stream)
      : yardl::ndjson::NDJsonWriter(stream, schema_) {
  }

  MyProtocolWriter(std::string file_name)
      : yardl::ndjson::NDJsonWriter(file_name, schema_) {
  }

  void Flush() override;

  protected:
  void WriteTreeImpl(sketch::BinaryTree const& value) override;
  void WritePtreeImpl(std::unique_ptr<sketch::BinaryTree> const& value) override;
  void WriteListImpl(std::optional<sketch::LinkedList<std::string>> const& value) override;
  // dirs: !stream
  //   items: Directory
  void WriteCwdImpl(sketch::DirectoryEntry const& value) override;
  void EndCwdImpl() override {}
  void CloseImpl() override;
};

// NDJSON reader for the MyProtocol protocol.
class MyProtocolReader : public sketch::MyProtocolReaderBase, yardl::ndjson::NDJsonReader {
  public:
  MyProtocolReader(std::istream& stream)
      : yardl::ndjson::NDJsonReader(stream, schema_) {
  }

  MyProtocolReader(std::string file_name)
      : yardl::ndjson::NDJsonReader(file_name, schema_) {
  }

  protected:
  void ReadTreeImpl(sketch::BinaryTree& value) override;
  void ReadPtreeImpl(std::unique_ptr<sketch::BinaryTree>& value) override;
  void ReadListImpl(std::optional<sketch::LinkedList<std::string>>& value) override;
  bool ReadCwdImpl(sketch::DirectoryEntry& value) override;
  void CloseImpl() override;
};

} // namespace sketch::ndjson