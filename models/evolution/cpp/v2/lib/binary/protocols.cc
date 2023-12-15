// This file was generated by the "yardl" tool. DO NOT EDIT.

#include "protocols.h"

#include <cstddef>

#include "../yardl/detail/binary/coded_stream.h"
#include "../yardl/detail/binary/serializers.h"

namespace yardl::binary {
#ifndef _MSC_VER
// Values of offsetof() are only used if types are standard-layout.
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Winvalid-offsetof"
#endif

template <>
struct IsTriviallySerializable<evo_test::Header> {
  using __T__ = evo_test::Header;
  static constexpr bool value = 
    std::is_standard_layout_v<__T__> &&
    IsTriviallySerializable<decltype(__T__::subject)>::value &&
    IsTriviallySerializable<decltype(__T__::meta)>::value &&
    IsTriviallySerializable<decltype(__T__::weight)>::value &&
    (sizeof(__T__) == (sizeof(__T__::subject) + sizeof(__T__::meta) + sizeof(__T__::weight))) &&
    offsetof(__T__, subject) < offsetof(__T__, meta) && offsetof(__T__, meta) < offsetof(__T__, weight);
};

template <>
struct IsTriviallySerializable<evo_test::Sample> {
  using __T__ = evo_test::Sample;
  static constexpr bool value = 
    std::is_standard_layout_v<__T__> &&
    IsTriviallySerializable<decltype(__T__::timestamp)>::value &&
    IsTriviallySerializable<decltype(__T__::data)>::value &&
    (sizeof(__T__) == (sizeof(__T__::timestamp) + sizeof(__T__::data))) &&
    offsetof(__T__, timestamp) < offsetof(__T__, data);
};

template <>
struct IsTriviallySerializable<evo_test::Signature> {
  using __T__ = evo_test::Signature;
  static constexpr bool value = 
    std::is_standard_layout_v<__T__> &&
    IsTriviallySerializable<decltype(__T__::name)>::value &&
    IsTriviallySerializable<decltype(__T__::email)>::value &&
    IsTriviallySerializable<decltype(__T__::number)>::value &&
    (sizeof(__T__) == (sizeof(__T__::name) + sizeof(__T__::email) + sizeof(__T__::number))) &&
    offsetof(__T__, name) < offsetof(__T__, email) && offsetof(__T__, email) < offsetof(__T__, number);
};

template <>
struct IsTriviallySerializable<evo_test::Footer> {
  using __T__ = evo_test::Footer;
  static constexpr bool value = 
    std::is_standard_layout_v<__T__> &&
    IsTriviallySerializable<decltype(__T__::signature)>::value &&
    (sizeof(__T__) == (sizeof(__T__::signature)));
};

template <>
struct IsTriviallySerializable<evo_test::NewRecord> {
  using __T__ = evo_test::NewRecord;
  static constexpr bool value = 
    std::is_standard_layout_v<__T__> &&
    IsTriviallySerializable<decltype(__T__::stuff)>::value &&
    (sizeof(__T__) == (sizeof(__T__::stuff)));
};

#ifndef _MSC_VER
#pragma GCC diagnostic pop // #pragma GCC diagnostic ignored "-Winvalid-offsetof" 
#endif
} //namespace yardl::binary 

namespace {
template<typename T0, yardl::binary::Writer<T0> WriteT0, typename T1, yardl::binary::Writer<T1> WriteT1>
void WriteUnion(yardl::binary::CodedOutputStream& stream, std::variant<T0, T1> const& value) {
  yardl::binary::WriteInteger(stream, value.index());
  switch (value.index()) {
  case 0: {
    T0 const& v = std::get<0>(value);
    WriteT0(stream, v);
    break;
  }
  case 1: {
    T1 const& v = std::get<1>(value);
    WriteT1(stream, v);
    break;
  }
  default: throw std::runtime_error("Invalid union index.");
  }
}

template<typename T0, yardl::binary::Reader<T0> ReadT0, typename T1, yardl::binary::Reader<T1> ReadT1>
void ReadUnion(yardl::binary::CodedInputStream& stream, std::variant<T0, T1>& value) {
  size_t index;
  yardl::binary::ReadInteger(stream, index);
  switch (index) {
    case 0: {
      T0 v;
      ReadT0(stream, v);
      value = std::move(v);
      break;
    }
    case 1: {
      T1 v;
      ReadT1(stream, v);
      value = std::move(v);
      break;
    }
    default: throw std::runtime_error("Invalid union index.");
  }
}
} // namespace

namespace evo_test::binary {
namespace {
[[maybe_unused]] static void WriteHeader(yardl::binary::CodedOutputStream& stream, evo_test::Header const& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Header>::value) {
    yardl::binary::WriteTriviallySerializable(stream, value);
    return;
  }

  WriteUnion<std::string, yardl::binary::WriteString, int64_t, yardl::binary::WriteInteger>(stream, value.subject);
  yardl::binary::WriteMap<std::string, std::vector<std::string>, yardl::binary::WriteString, yardl::binary::WriteVector<std::string, yardl::binary::WriteString>>(stream, value.meta);
  yardl::binary::WriteFloatingPoint(stream, value.weight);
}

[[maybe_unused]] static void ReadHeader(yardl::binary::CodedInputStream& stream, evo_test::Header& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Header>::value) {
    yardl::binary::ReadTriviallySerializable(stream, value);
    return;
  }

  ReadUnion<std::string, yardl::binary::ReadString, int64_t, yardl::binary::ReadInteger>(stream, value.subject);
  yardl::binary::ReadMap<std::string, std::vector<std::string>, yardl::binary::ReadString, yardl::binary::ReadVector<std::string, yardl::binary::ReadString>>(stream, value.meta);
  yardl::binary::ReadFloatingPoint(stream, value.weight);
}

[[maybe_unused]] static void WriteHeader_v0(yardl::binary::CodedOutputStream& stream, evo_test::Header const& value) {
  std::string subject;
  subject = std::get<0>(value.subject);
  yardl::binary::WriteString(stream, subject);
  int64_t weight;
  weight = (double)(value.weight);
  yardl::binary::WriteInteger(stream, weight);
  yardl::binary::WriteMap<std::string, std::vector<std::string>, yardl::binary::WriteString, yardl::binary::WriteVector<std::string, yardl::binary::WriteString>>(stream, value.meta);
}

[[maybe_unused]] static void ReadHeader_v0(yardl::binary::CodedInputStream& stream, evo_test::Header& value) {
  std::string subject;
  yardl::binary::ReadString(stream, subject);
  value.subject = subject;
  int64_t weight;
  yardl::binary::ReadInteger(stream, weight);
  value.weight = (double)(weight);
  yardl::binary::ReadMap<std::string, std::vector<std::string>, yardl::binary::ReadString, yardl::binary::ReadVector<std::string, yardl::binary::ReadString>>(stream, value.meta);
}

[[maybe_unused]] static void WriteHeader_v1(yardl::binary::CodedOutputStream& stream, evo_test::Header const& value) {
  yardl::binary::WriteMap<std::string, std::vector<std::string>, yardl::binary::WriteString, yardl::binary::WriteVector<std::string, yardl::binary::WriteString>>(stream, value.meta);
  std::string subject;
  subject = std::get<0>(value.subject);
  yardl::binary::WriteString(stream, subject);
  int64_t weight;
  weight = (double)(value.weight);
  yardl::binary::WriteInteger(stream, weight);
  std::optional<std::string> added;
  yardl::binary::WriteOptional<std::string, yardl::binary::WriteString>(stream, added);
}

[[maybe_unused]] static void ReadHeader_v1(yardl::binary::CodedInputStream& stream, evo_test::Header& value) {
  yardl::binary::ReadMap<std::string, std::vector<std::string>, yardl::binary::ReadString, yardl::binary::ReadVector<std::string, yardl::binary::ReadString>>(stream, value.meta);
  std::string subject;
  yardl::binary::ReadString(stream, subject);
  value.subject = subject;
  int64_t weight;
  yardl::binary::ReadInteger(stream, weight);
  value.weight = (double)(weight);
  std::optional<std::string> added;
  yardl::binary::ReadOptional<std::string, yardl::binary::ReadString>(stream, added);
}

[[maybe_unused]] static void WriteSample(yardl::binary::CodedOutputStream& stream, evo_test::Sample const& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Sample>::value) {
    yardl::binary::WriteTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::WriteDateTime(stream, value.timestamp);
  yardl::binary::WriteVector<int32_t, yardl::binary::WriteInteger>(stream, value.data);
}

[[maybe_unused]] static void ReadSample(yardl::binary::CodedInputStream& stream, evo_test::Sample& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Sample>::value) {
    yardl::binary::ReadTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::ReadDateTime(stream, value.timestamp);
  yardl::binary::ReadVector<int32_t, yardl::binary::ReadInteger>(stream, value.data);
}

[[maybe_unused]] static void WriteSample_v0(yardl::binary::CodedOutputStream& stream, evo_test::Sample const& value) {
  yardl::binary::WriteVector<int32_t, yardl::binary::WriteInteger>(stream, value.data);
  yardl::binary::WriteDateTime(stream, value.timestamp);
}

[[maybe_unused]] static void ReadSample_v0(yardl::binary::CodedInputStream& stream, evo_test::Sample& value) {
  yardl::binary::ReadVector<int32_t, yardl::binary::ReadInteger>(stream, value.data);
  yardl::binary::ReadDateTime(stream, value.timestamp);
}

[[maybe_unused]] static void WriteSample_v1(yardl::binary::CodedOutputStream& stream, evo_test::Sample const& value) {
  yardl::binary::WriteVector<int32_t, yardl::binary::WriteInteger>(stream, value.data);
  yardl::binary::WriteDateTime(stream, value.timestamp);
}

[[maybe_unused]] static void ReadSample_v1(yardl::binary::CodedInputStream& stream, evo_test::Sample& value) {
  yardl::binary::ReadVector<int32_t, yardl::binary::ReadInteger>(stream, value.data);
  yardl::binary::ReadDateTime(stream, value.timestamp);
}

[[maybe_unused]] static void WriteSignature(yardl::binary::CodedOutputStream& stream, evo_test::Signature const& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Signature>::value) {
    yardl::binary::WriteTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::WriteString(stream, value.name);
  yardl::binary::WriteString(stream, value.email);
  yardl::binary::WriteString(stream, value.number);
}

[[maybe_unused]] static void ReadSignature(yardl::binary::CodedInputStream& stream, evo_test::Signature& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Signature>::value) {
    yardl::binary::ReadTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::ReadString(stream, value.name);
  yardl::binary::ReadString(stream, value.email);
  yardl::binary::ReadString(stream, value.number);
}

[[maybe_unused]] static void WriteSignature_v0(yardl::binary::CodedOutputStream& stream, evo_test::Signature const& value) {
  yardl::binary::WriteString(stream, value.name);
  yardl::binary::WriteString(stream, value.email);
  int64_t number;
  number = std::stol(value.number);
  yardl::binary::WriteInteger(stream, number);
}

[[maybe_unused]] static void ReadSignature_v0(yardl::binary::CodedInputStream& stream, evo_test::Signature& value) {
  yardl::binary::ReadString(stream, value.name);
  yardl::binary::ReadString(stream, value.email);
  int64_t number;
  yardl::binary::ReadInteger(stream, number);
  value.number = std::to_string(number);
}

[[maybe_unused]] static void WriteSignature_v1(yardl::binary::CodedOutputStream& stream, evo_test::Signature const& value) {
  yardl::binary::WriteString(stream, value.name);
  yardl::binary::WriteString(stream, value.email);
  int64_t number;
  number = std::stol(value.number);
  yardl::binary::WriteInteger(stream, number);
}

[[maybe_unused]] static void ReadSignature_v1(yardl::binary::CodedInputStream& stream, evo_test::Signature& value) {
  yardl::binary::ReadString(stream, value.name);
  yardl::binary::ReadString(stream, value.email);
  int64_t number;
  yardl::binary::ReadInteger(stream, number);
  value.number = std::to_string(number);
}

[[maybe_unused]] static void WriteFooter(yardl::binary::CodedOutputStream& stream, evo_test::Footer const& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Footer>::value) {
    yardl::binary::WriteTriviallySerializable(stream, value);
    return;
  }

  evo_test::binary::WriteSignature(stream, value.signature);
}

[[maybe_unused]] static void ReadFooter(yardl::binary::CodedInputStream& stream, evo_test::Footer& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Footer>::value) {
    yardl::binary::ReadTriviallySerializable(stream, value);
    return;
  }

  evo_test::binary::ReadSignature(stream, value.signature);
}

[[maybe_unused]] static void WriteFooter_v0(yardl::binary::CodedOutputStream& stream, evo_test::Footer const& value) {
  evo_test::binary::WriteSignature_v0(stream, value.signature);
}

[[maybe_unused]] static void ReadFooter_v0(yardl::binary::CodedInputStream& stream, evo_test::Footer& value) {
  evo_test::binary::ReadSignature_v0(stream, value.signature);
}

[[maybe_unused]] static void WriteFooter_v1(yardl::binary::CodedOutputStream& stream, evo_test::Footer const& value) {
  evo_test::binary::WriteSignature_v1(stream, value.signature);
}

[[maybe_unused]] static void ReadFooter_v1(yardl::binary::CodedInputStream& stream, evo_test::Footer& value) {
  evo_test::binary::ReadSignature_v1(stream, value.signature);
}

[[maybe_unused]] static void WriteNewRecord(yardl::binary::CodedOutputStream& stream, evo_test::NewRecord const& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::NewRecord>::value) {
    yardl::binary::WriteTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::WriteDynamicNDArray<double, yardl::binary::WriteFloatingPoint>(stream, value.stuff);
}

[[maybe_unused]] static void ReadNewRecord(yardl::binary::CodedInputStream& stream, evo_test::NewRecord& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::NewRecord>::value) {
    yardl::binary::ReadTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::ReadDynamicNDArray<double, yardl::binary::ReadFloatingPoint>(stream, value.stuff);
}

} // namespace

void MyProtocolWriter::WriteHeaderImpl(evo_test::Header const& value) {
  switch (schema_index_) {
  case 0:
    evo_test::binary::WriteHeader_v0(stream_, value);
    break;
  case 1:
    evo_test::binary::WriteHeader_v1(stream_, value);
    break;
  default:
    evo_test::binary::WriteHeader(stream_, value);
    break;
  }
}

void MyProtocolWriter::WriteIdImpl(std::string const& value) {
  switch (schema_index_) {
  case 0:
    int64_t id_0;
    id_0 = std::stol(value);
    yardl::binary::WriteInteger(stream_, id_0);
    break;
  case 1:
    int64_t id_1;
    id_1 = std::stol(value);
    yardl::binary::WriteInteger(stream_, id_1);
    break;
  default:
    yardl::binary::WriteString(stream_, value);
    break;
  }
}

void MyProtocolWriter::WriteSamplesImpl(evo_test::Sample const& value) {
  yardl::binary::WriteInteger(stream_, 1U);
  switch (schema_index_) {
  case 0:
    evo_test::binary::WriteSample_v0(stream_, value);
    break;
  case 1:
    evo_test::binary::WriteSample_v1(stream_, value);
    break;
  default:
    evo_test::binary::WriteSample(stream_, value);
    break;
  }
}

void MyProtocolWriter::WriteSamplesImpl(std::vector<evo_test::Sample> const& values) {
  if (!values.empty()) {
    switch (schema_index_) {
    case 0:
      yardl::binary::WriteVector<evo_test::Sample, evo_test::binary::WriteSample_v0>(stream_, values);
      break;
    case 1:
      yardl::binary::WriteVector<evo_test::Sample, evo_test::binary::WriteSample_v1>(stream_, values);
      break;
    default:
      yardl::binary::WriteVector<evo_test::Sample, evo_test::binary::WriteSample>(stream_, values);
      break;
    }
  }
}

void MyProtocolWriter::EndSamplesImpl() {
  yardl::binary::WriteInteger(stream_, 0U);
}

void MyProtocolWriter::WriteFooterImpl(std::optional<evo_test::Footer> const& value) {
  switch (schema_index_) {
  case 0:
    yardl::binary::WriteOptional<evo_test::Footer, evo_test::binary::WriteFooter_v0>(stream_, value);
    break;
  case 1:
    yardl::binary::WriteOptional<evo_test::Footer, evo_test::binary::WriteFooter_v1>(stream_, value);
    break;
  default:
    yardl::binary::WriteOptional<evo_test::Footer, evo_test::binary::WriteFooter>(stream_, value);
    break;
  }
}

void MyProtocolWriter::Flush() {
  stream_.Flush();
}

void MyProtocolWriter::CloseImpl() {
  stream_.Flush();
}

void MyProtocolReader::ReadHeaderImpl(evo_test::Header& value) {
  switch (schema_index_) {
  case 0:
    evo_test::binary::ReadHeader_v0(stream_, value);
    break;
  case 1:
    evo_test::binary::ReadHeader_v1(stream_, value);
    break;
  default:
    evo_test::binary::ReadHeader(stream_, value);
    break;
  }
}

void MyProtocolReader::ReadIdImpl(std::string& value) {
  switch (schema_index_) {
  case 0:
    int64_t id_0;
    yardl::binary::ReadInteger(stream_, id_0);
    value = std::to_string(id_0);
    break;
  case 1:
    int64_t id_1;
    yardl::binary::ReadInteger(stream_, id_1);
    value = std::to_string(id_1);
    break;
  default:
    yardl::binary::ReadString(stream_, value);
    break;
  }
}

bool MyProtocolReader::ReadSamplesImpl(evo_test::Sample& value) {
  if (current_block_remaining_ == 0) {
    yardl::binary::ReadInteger(stream_, current_block_remaining_);
    if (current_block_remaining_ == 0) {
      return false;
    }
  }
  switch (schema_index_) {
  case 0:
    evo_test::binary::ReadSample_v0(stream_, value);
    break;
  case 1:
    evo_test::binary::ReadSample_v1(stream_, value);
    break;
  default:
    evo_test::binary::ReadSample(stream_, value);
    break;
  }
  current_block_remaining_--;
  return true;
}

bool MyProtocolReader::ReadSamplesImpl(std::vector<evo_test::Sample>& values) {
  switch (schema_index_) {
  case 0:
    yardl::binary::ReadBlocksIntoVector<evo_test::Sample, evo_test::binary::ReadSample_v0>(stream_, current_block_remaining_, values);
    break;
  case 1:
    yardl::binary::ReadBlocksIntoVector<evo_test::Sample, evo_test::binary::ReadSample_v1>(stream_, current_block_remaining_, values);
    break;
  default:
    yardl::binary::ReadBlocksIntoVector<evo_test::Sample, evo_test::binary::ReadSample>(stream_, current_block_remaining_, values);
    break;
  }
  return current_block_remaining_ != 0;
}

void MyProtocolReader::ReadFooterImpl(std::optional<evo_test::Footer>& value) {
  switch (schema_index_) {
  case 0:
    yardl::binary::ReadOptional<evo_test::Footer, evo_test::binary::ReadFooter_v0>(stream_, value);
    break;
  case 1:
    yardl::binary::ReadOptional<evo_test::Footer, evo_test::binary::ReadFooter_v1>(stream_, value);
    break;
  default:
    yardl::binary::ReadOptional<evo_test::Footer, evo_test::binary::ReadFooter>(stream_, value);
    break;
  }
}

void MyProtocolReader::CloseImpl() {
  stream_.VerifyFinished();
}

void NewProtocolWriter::WriteCalibrationImpl(std::vector<double> const& value) {
  yardl::binary::WriteVector<double, yardl::binary::WriteFloatingPoint>(stream_, value);
}

void NewProtocolWriter::WriteDataImpl(evo_test::NewRecord const& value) {
  yardl::binary::WriteInteger(stream_, 1U);
  evo_test::binary::WriteNewRecord(stream_, value);
}

void NewProtocolWriter::WriteDataImpl(std::vector<evo_test::NewRecord> const& values) {
  if (!values.empty()) {
    yardl::binary::WriteVector<evo_test::NewRecord, evo_test::binary::WriteNewRecord>(stream_, values);
  }
}

void NewProtocolWriter::EndDataImpl() {
  yardl::binary::WriteInteger(stream_, 0U);
}

void NewProtocolWriter::Flush() {
  stream_.Flush();
}

void NewProtocolWriter::CloseImpl() {
  stream_.Flush();
}

void NewProtocolReader::ReadCalibrationImpl(std::vector<double>& value) {
  yardl::binary::ReadVector<double, yardl::binary::ReadFloatingPoint>(stream_, value);
}

bool NewProtocolReader::ReadDataImpl(evo_test::NewRecord& value) {
  if (current_block_remaining_ == 0) {
    yardl::binary::ReadInteger(stream_, current_block_remaining_);
    if (current_block_remaining_ == 0) {
      return false;
    }
  }
  evo_test::binary::ReadNewRecord(stream_, value);
  current_block_remaining_--;
  return true;
}

bool NewProtocolReader::ReadDataImpl(std::vector<evo_test::NewRecord>& values) {
  yardl::binary::ReadBlocksIntoVector<evo_test::NewRecord, evo_test::binary::ReadNewRecord>(stream_, current_block_remaining_, values);
  return current_block_remaining_ != 0;
}

void NewProtocolReader::CloseImpl() {
  stream_.VerifyFinished();
}

} // namespace evo_test::binary

