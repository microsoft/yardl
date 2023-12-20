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
#include "../types.h"

namespace evo_test::ndjson {
// NDJSON writer for the MyProtocol protocol.
class MyProtocolWriter : public evo_test::MyProtocolWriterBase, yardl::ndjson::NDJsonWriter {
  public:
  MyProtocolWriter(std::ostream& stream)
      : yardl::ndjson::NDJsonWriter(stream, schema_) {
  }

  MyProtocolWriter(std::string file_name)
      : yardl::ndjson::NDJsonWriter(file_name, schema_) {
  }

  void Flush() override;

  protected:
  void WriteHeaderImpl(evo_test::Header const& value) override;
  void WriteIdImpl(std::string const& value) override;
  void WriteSamplesImpl(evo_test::Sample const& value) override;
  void EndSamplesImpl() override {}
  void WriteMaybeImpl(std::optional<std::string> const& value) override;
  void WriteFooterImpl(std::optional<evo_test::Footer> const& value) override;
  void CloseImpl() override;
};

// NDJSON reader for the MyProtocol protocol.
class MyProtocolReader : public evo_test::MyProtocolReaderBase, yardl::ndjson::NDJsonReader {
  public:
  MyProtocolReader(std::istream& stream)
      : yardl::ndjson::NDJsonReader(stream, schema_) {
  }

  MyProtocolReader(std::string file_name)
      : yardl::ndjson::NDJsonReader(file_name, schema_) {
  }

  protected:
  void ReadHeaderImpl(evo_test::Header& value) override;
  void ReadIdImpl(std::string& value) override;
  bool ReadSamplesImpl(evo_test::Sample& value) override;
  void ReadMaybeImpl(std::optional<std::string>& value) override;
  void ReadFooterImpl(std::optional<evo_test::Footer>& value) override;
  void CloseImpl() override;
};

// NDJSON writer for the NewProtocol protocol.
class NewProtocolWriter : public evo_test::NewProtocolWriterBase, yardl::ndjson::NDJsonWriter {
  public:
  NewProtocolWriter(std::ostream& stream)
      : yardl::ndjson::NDJsonWriter(stream, schema_) {
  }

  NewProtocolWriter(std::string file_name)
      : yardl::ndjson::NDJsonWriter(file_name, schema_) {
  }

  void Flush() override;

  protected:
  void WriteCalibrationImpl(std::vector<double> const& value) override;
  void WriteDataImpl(evo_test::NewRecord const& value) override;
  void EndDataImpl() override {}
  void CloseImpl() override;
};

// NDJSON reader for the NewProtocol protocol.
class NewProtocolReader : public evo_test::NewProtocolReaderBase, yardl::ndjson::NDJsonReader {
  public:
  NewProtocolReader(std::istream& stream)
      : yardl::ndjson::NDJsonReader(stream, schema_) {
  }

  NewProtocolReader(std::string file_name)
      : yardl::ndjson::NDJsonReader(file_name, schema_) {
  }

  protected:
  void ReadCalibrationImpl(std::vector<double>& value) override;
  bool ReadDataImpl(evo_test::NewRecord& value) override;
  void CloseImpl() override;
};

} // namespace evo_test::ndjson
