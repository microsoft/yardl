// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include "types.h"

namespace evo_test {
// Abstract writer for the MyProtocol protocol.
class MyProtocolWriterBase {
  public:
  // Ordinal 0.
  void WriteHeader(evo_test::Header const& value);

  // Ordinal 1.
  void WriteId(int64_t const& value);

  // Ordinal 2.
  // Call this method for each element of the `samples` stream, then call `EndSamples() when done.`
  void WriteSamples(evo_test::Sample const& value);

  // Ordinal 2.
  // Call this method to write many values to the `samples` stream, then call `EndSamples()` when done.
  void WriteSamples(std::vector<evo_test::Sample> const& values);

  // Marks the end of the `samples` stream.
  void EndSamples();

  // Ordinal 3.
  void WriteFooter(std::optional<evo_test::Footer> const& value);

  // Optionaly close this writer before destructing. Validates that all steps were completed.
  void Close();

  virtual ~MyProtocolWriterBase() = default;

  // Flushes all buffered data.
  virtual void Flush() {}

  protected:
  virtual void WriteHeaderImpl(evo_test::Header const& value) = 0;
  virtual void WriteIdImpl(int64_t const& value) = 0;
  virtual void WriteSamplesImpl(evo_test::Sample const& value) = 0;
  virtual void WriteSamplesImpl(std::vector<evo_test::Sample> const& value);
  virtual void EndSamplesImpl() = 0;
  virtual void WriteFooterImpl(std::optional<evo_test::Footer> const& value) = 0;
  virtual void CloseImpl() {}

  static std::string schema_;

  static std::vector<std::string> previous_schemas_;

  private:
  uint8_t state_ = 0;

  friend class MyProtocolReaderBase;
};

// Abstract reader for the MyProtocol protocol.
class MyProtocolReaderBase {
  public:
  // Ordinal 0.
  void ReadHeader(evo_test::Header& value);

  // Ordinal 1.
  void ReadId(int64_t& value);

  // Ordinal 2.
  [[nodiscard]] bool ReadSamples(evo_test::Sample& value);

  // Ordinal 2.
  [[nodiscard]] bool ReadSamples(std::vector<evo_test::Sample>& values);

  // Ordinal 3.
  void ReadFooter(std::optional<evo_test::Footer>& value);

  // Optionaly close this writer before destructing. Validates that all steps were completely read.
  void Close();

  void CopyTo(MyProtocolWriterBase& writer, size_t samples_buffer_size = 1);

  virtual ~MyProtocolReaderBase() = default;

  protected:
  virtual void ReadHeaderImpl(evo_test::Header& value) = 0;
  virtual void ReadIdImpl(int64_t& value) = 0;
  virtual bool ReadSamplesImpl(evo_test::Sample& value) = 0;
  virtual bool ReadSamplesImpl(std::vector<evo_test::Sample>& values);
  virtual void ReadFooterImpl(std::optional<evo_test::Footer>& value) = 0;
  virtual void CloseImpl() {}
  static std::string schema_;

  static std::vector<std::string> previous_schemas_;

  private:
  uint8_t state_ = 0;
};
} // namespace evo_test
