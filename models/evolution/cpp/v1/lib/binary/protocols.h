// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include <array>
#include <complex>
#include <memory>
#include <optional>
#include <variant>
#include <vector>

#include "../yardl/detail/binary/reader_writer.h"
#include "../protocols.h"
#include "../types.h"

namespace evo_test::binary {
// Binary writer for the MyProtocol protocol.
class MyProtocolWriter : public evo_test::MyProtocolWriterBase, yardl::binary::BinaryWriter {
  public:
  MyProtocolWriter(std::ostream& stream, const std::string& schema=schema_)
      : yardl::binary::BinaryWriter(stream, schema, schema_, previous_schemas_) {
  }

  MyProtocolWriter(std::string file_name, const std::string& schema=schema_)
      : yardl::binary::BinaryWriter(file_name, schema, schema_, previous_schemas_) {
  }

  void Flush() override;

  protected:
  void WriteHeaderImpl(evo_test::Header const& value) override;
  void WriteIdImpl(int64_t const& value) override;
  void WriteSamplesImpl(evo_test::Sample const& value) override;
  void WriteSamplesImpl(std::vector<evo_test::Sample> const& values) override;
  void EndSamplesImpl() override;
  void WriteFooterImpl(std::optional<evo_test::Footer> const& value) override;
  void CloseImpl() override;
};

// Binary reader for the MyProtocol protocol.
class MyProtocolReader : public evo_test::MyProtocolReaderBase, yardl::binary::BinaryReader {
  public:
  MyProtocolReader(std::istream& stream)
      : yardl::binary::BinaryReader(stream, schema_, previous_schemas_) {
  }

  MyProtocolReader(std::string file_name)
      : yardl::binary::BinaryReader(file_name, schema_, previous_schemas_) {
  }

  std::string GetSchema() { if (schema_index_ < 0) { return schema_; } else { return previous_schemas_[schema_index_]; } }

  protected:
  void ReadHeaderImpl(evo_test::Header& value) override;
  void ReadIdImpl(int64_t& value) override;
  bool ReadSamplesImpl(evo_test::Sample& value) override;
  bool ReadSamplesImpl(std::vector<evo_test::Sample>& values) override;
  void ReadFooterImpl(std::optional<evo_test::Footer>& value) override;
  void CloseImpl() override;

  private:
  size_t current_block_remaining_ = 0;
};

// Binary writer for the UnusedProtocol protocol.
class UnusedProtocolWriter : public evo_test::UnusedProtocolWriterBase, yardl::binary::BinaryWriter {
  public:
  UnusedProtocolWriter(std::ostream& stream, const std::string& schema=schema_)
      : yardl::binary::BinaryWriter(stream, schema, schema_, previous_schemas_) {
  }

  UnusedProtocolWriter(std::string file_name, const std::string& schema=schema_)
      : yardl::binary::BinaryWriter(file_name, schema, schema_, previous_schemas_) {
  }

  void Flush() override;

  protected:
  void WriteSamplesImpl(evo_test::Sample const& value) override;
  void WriteSamplesImpl(std::vector<evo_test::Sample> const& values) override;
  void EndSamplesImpl() override;
  void CloseImpl() override;
};

// Binary reader for the UnusedProtocol protocol.
class UnusedProtocolReader : public evo_test::UnusedProtocolReaderBase, yardl::binary::BinaryReader {
  public:
  UnusedProtocolReader(std::istream& stream)
      : yardl::binary::BinaryReader(stream, schema_, previous_schemas_) {
  }

  UnusedProtocolReader(std::string file_name)
      : yardl::binary::BinaryReader(file_name, schema_, previous_schemas_) {
  }

  std::string GetSchema() { if (schema_index_ < 0) { return schema_; } else { return previous_schemas_[schema_index_]; } }

  protected:
  bool ReadSamplesImpl(evo_test::Sample& value) override;
  bool ReadSamplesImpl(std::vector<evo_test::Sample>& values) override;
  void CloseImpl() override;

  private:
  size_t current_block_remaining_ = 0;
};

} // namespace evo_test::binary
