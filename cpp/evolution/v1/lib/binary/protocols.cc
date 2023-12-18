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
    IsTriviallySerializable<decltype(__T__::meta)>::value &&
    IsTriviallySerializable<decltype(__T__::subject)>::value &&
    IsTriviallySerializable<decltype(__T__::weight)>::value &&
    IsTriviallySerializable<decltype(__T__::added)>::value &&
    (sizeof(__T__) == (sizeof(__T__::meta) + sizeof(__T__::subject) + sizeof(__T__::weight) + sizeof(__T__::added))) &&
    offsetof(__T__, meta) < offsetof(__T__, subject) && offsetof(__T__, subject) < offsetof(__T__, weight) && offsetof(__T__, weight) < offsetof(__T__, added);
};

template <>
struct IsTriviallySerializable<evo_test::Sample> {
  using __T__ = evo_test::Sample;
  static constexpr bool value = 
    std::is_standard_layout_v<__T__> &&
    IsTriviallySerializable<decltype(__T__::data)>::value &&
    IsTriviallySerializable<decltype(__T__::timestamp)>::value &&
    (sizeof(__T__) == (sizeof(__T__::data) + sizeof(__T__::timestamp))) &&
    offsetof(__T__, data) < offsetof(__T__, timestamp);
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

#ifndef _MSC_VER
#pragma GCC diagnostic pop // #pragma GCC diagnostic ignored "-Winvalid-offsetof" 
#endif
} //namespace yardl::binary 

namespace evo_test::binary {
namespace {
[[maybe_unused]] static void WriteHeader(yardl::binary::CodedOutputStream& stream, evo_test::Header const& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Header>::value) {
    yardl::binary::WriteTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::WriteMap<std::string, std::vector<std::string>, yardl::binary::WriteString, yardl::binary::WriteVector<std::string, yardl::binary::WriteString>>(stream, value.meta);
  yardl::binary::WriteString(stream, value.subject);
  yardl::binary::WriteInteger(stream, value.weight);
  yardl::binary::WriteOptional<std::string, yardl::binary::WriteString>(stream, value.added);
}

[[maybe_unused]] static void ReadHeader(yardl::binary::CodedInputStream& stream, evo_test::Header& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Header>::value) {
    yardl::binary::ReadTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::ReadMap<std::string, std::vector<std::string>, yardl::binary::ReadString, yardl::binary::ReadVector<std::string, yardl::binary::ReadString>>(stream, value.meta);
  yardl::binary::ReadString(stream, value.subject);
  yardl::binary::ReadInteger(stream, value.weight);
  yardl::binary::ReadOptional<std::string, yardl::binary::ReadString>(stream, value.added);
}

[[maybe_unused]] static void WriteHeader_v0(yardl::binary::CodedOutputStream& stream, evo_test::Header const& value) {
  yardl::binary::WriteString(stream, value.subject);
  yardl::binary::WriteInteger(stream, value.weight);
  yardl::binary::WriteMap<std::string, std::vector<std::string>, yardl::binary::WriteString, yardl::binary::WriteVector<std::string, yardl::binary::WriteString>>(stream, value.meta);
}

[[maybe_unused]] static void ReadHeader_v0(yardl::binary::CodedInputStream& stream, evo_test::Header& value) {
  yardl::binary::ReadString(stream, value.subject);
  yardl::binary::ReadInteger(stream, value.weight);
  yardl::binary::ReadMap<std::string, std::vector<std::string>, yardl::binary::ReadString, yardl::binary::ReadVector<std::string, yardl::binary::ReadString>>(stream, value.meta);
}

[[maybe_unused]] static void WriteSample(yardl::binary::CodedOutputStream& stream, evo_test::Sample const& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Sample>::value) {
    yardl::binary::WriteTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::WriteVector<int32_t, yardl::binary::WriteInteger>(stream, value.data);
  yardl::binary::WriteDateTime(stream, value.timestamp);
}

[[maybe_unused]] static void ReadSample(yardl::binary::CodedInputStream& stream, evo_test::Sample& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Sample>::value) {
    yardl::binary::ReadTriviallySerializable(stream, value);
    return;
  }

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
  yardl::binary::WriteInteger(stream, value.number);
}

[[maybe_unused]] static void ReadSignature(yardl::binary::CodedInputStream& stream, evo_test::Signature& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::Signature>::value) {
    yardl::binary::ReadTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::ReadString(stream, value.name);
  yardl::binary::ReadString(stream, value.email);
  yardl::binary::ReadInteger(stream, value.number);
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

[[maybe_unused]] static void WriteAliasedPrimitive(yardl::binary::CodedOutputStream& stream, evo_test::AliasedPrimitive const& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::AliasedPrimitive>::value) {
    yardl::binary::WriteTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::WriteString(stream, value);
}

[[maybe_unused]] static void ReadAliasedPrimitive(yardl::binary::CodedInputStream& stream, evo_test::AliasedPrimitive& value) {
  if constexpr (yardl::binary::IsTriviallySerializable<evo_test::AliasedPrimitive>::value) {
    yardl::binary::ReadTriviallySerializable(stream, value);
    return;
  }

  yardl::binary::ReadString(stream, value);
}

} // namespace

void MyProtocolWriter::WriteHeaderImpl(evo_test::Header const& value) {
  switch (schema_index_) {
  case 0:
    evo_test::binary::WriteHeader_v0(stream_, value);
    break;
  default:
    evo_test::binary::WriteHeader(stream_, value);
    break;
  }
}

void MyProtocolWriter::WriteIdImpl(int64_t const& value) {
  switch (schema_index_) {
  default:
    yardl::binary::WriteInteger(stream_, value);
    break;
  }
}

void MyProtocolWriter::WriteSamplesImpl(evo_test::Sample const& value) {
  yardl::binary::WriteInteger(stream_, 1U);
  switch (schema_index_) {
  default:
    evo_test::binary::WriteSample(stream_, value);
    break;
  }
}

void MyProtocolWriter::WriteSamplesImpl(std::vector<evo_test::Sample> const& values) {
  if (!values.empty()) {
    switch (schema_index_) {
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
  default:
    evo_test::binary::ReadHeader(stream_, value);
    break;
  }
}

void MyProtocolReader::ReadIdImpl(int64_t& value) {
  switch (schema_index_) {
  default:
    yardl::binary::ReadInteger(stream_, value);
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
  default:
    evo_test::binary::ReadSample(stream_, value);
    break;
  }
  current_block_remaining_--;
  return true;
}

bool MyProtocolReader::ReadSamplesImpl(std::vector<evo_test::Sample>& values) {
  switch (schema_index_) {
  default:
    yardl::binary::ReadBlocksIntoVector<evo_test::Sample, evo_test::binary::ReadSample>(stream_, current_block_remaining_, values);
    break;
  }
  return current_block_remaining_ != 0;
}

void MyProtocolReader::ReadFooterImpl(std::optional<evo_test::Footer>& value) {
  switch (schema_index_) {
  default:
    yardl::binary::ReadOptional<evo_test::Footer, evo_test::binary::ReadFooter>(stream_, value);
    break;
  }
}

void MyProtocolReader::CloseImpl() {
  stream_.VerifyFinished();
}

void UnusedProtocolWriter::WriteSamplesImpl(evo_test::Sample const& value) {
  yardl::binary::WriteInteger(stream_, 1U);
  evo_test::binary::WriteSample(stream_, value);
}

void UnusedProtocolWriter::WriteSamplesImpl(std::vector<evo_test::Sample> const& values) {
  if (!values.empty()) {
    yardl::binary::WriteVector<evo_test::Sample, evo_test::binary::WriteSample>(stream_, values);
  }
}

void UnusedProtocolWriter::EndSamplesImpl() {
  yardl::binary::WriteInteger(stream_, 0U);
}

void UnusedProtocolWriter::Flush() {
  stream_.Flush();
}

void UnusedProtocolWriter::CloseImpl() {
  stream_.Flush();
}

bool UnusedProtocolReader::ReadSamplesImpl(evo_test::Sample& value) {
  if (current_block_remaining_ == 0) {
    yardl::binary::ReadInteger(stream_, current_block_remaining_);
    if (current_block_remaining_ == 0) {
      return false;
    }
  }
  evo_test::binary::ReadSample(stream_, value);
  current_block_remaining_--;
  return true;
}

bool UnusedProtocolReader::ReadSamplesImpl(std::vector<evo_test::Sample>& values) {
  yardl::binary::ReadBlocksIntoVector<evo_test::Sample, evo_test::binary::ReadSample>(stream_, current_block_remaining_, values);
  return current_block_remaining_ != 0;
}

void UnusedProtocolReader::CloseImpl() {
  stream_.VerifyFinished();
}

} // namespace evo_test::binary

