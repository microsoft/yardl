// This file was generated by the "yardl" tool. DO NOT EDIT.

#include "protocols.h"

#ifdef _MSC_VER
#define unlikely(x) x
#else
#define unlikely(x) __builtin_expect((x), 0)
#endif

namespace evo_test {
namespace {
void MyProtocolWriterBaseInvalidState(uint8_t attempted, [[maybe_unused]] bool end, uint8_t current) {
  std::string expected_method;
  switch (current) {
  case 0: expected_method = "WriteHeader()"; break;
  case 1: expected_method = "WriteId()"; break;
  case 2: expected_method = "WriteSamples() or EndSamples()"; break;
  case 3: expected_method = "WriteFooter()"; break;
  }
  std::string attempted_method;
  switch (attempted) {
  case 0: attempted_method = "WriteHeader()"; break;
  case 1: attempted_method = "WriteId()"; break;
  case 2: attempted_method = end ? "EndSamples()" : "WriteSamples()"; break;
  case 3: attempted_method = "WriteFooter()"; break;
  case 4: attempted_method = "Close()"; break;
  }
  throw std::runtime_error("Expected call to " + expected_method + " but received call to " + attempted_method + " instead.");
}

void MyProtocolReaderBaseInvalidState(uint8_t attempted, uint8_t current) {
  auto f = [](uint8_t i) -> std::string {
    switch (i/2) {
    case 0: return "ReadHeader()";
    case 1: return "ReadId()";
    case 2: return "ReadSamples()";
    case 3: return "ReadFooter()";
    case 4: return "Close()";
    default: return "<unknown>";
    }
  };
  throw std::runtime_error("Expected call to " + f(current) + " but received call to " + f(attempted) + " instead.");
}

} // namespace 

std::string MyProtocolWriterBase::schema_ = R"({"protocol":{"name":"MyProtocol","sequence":[{"name":"header","type":"EvoTest.Header"},{"name":"id","type":"int64"},{"name":"samples","type":{"stream":{"items":"EvoTest.Sample"}}},{"name":"footer","type":[null,"EvoTest.Footer"]}]},"types":[{"name":"Footer","fields":[{"name":"signature","type":"EvoTest.Signature"}]},{"name":"Header","fields":[{"name":"meta","type":{"map":{"keys":"string","values":{"vector":{"items":"string"}}}}},{"name":"subject","type":"string"},{"name":"weight","type":"int64"},{"name":"added","type":[null,"string"]}]},{"name":"Sample","fields":[{"name":"data","type":{"vector":{"items":"int32"}}},{"name":"timestamp","type":"datetime"}]},{"name":"Signature","fields":[{"name":"name","type":"string"},{"name":"email","type":"string"},{"name":"number","type":"int64"}]}]})";

std::vector<std::string> MyProtocolWriterBase::previous_schemas_ = {
  R"({"protocol":{"name":"MyProtocol","sequence":[{"name":"header","type":"EvoTest.Header"},{"name":"id","type":"int64"},{"name":"samples","type":{"stream":{"items":"EvoTest.Sample"}}},{"name":"footer","type":[null,"EvoTest.Footer"]}]},"types":[{"name":"Footer","fields":[{"name":"signature","type":"EvoTest.Signature"}]},{"name":"Header","fields":[{"name":"subject","type":"string"},{"name":"weight","type":"int64"},{"name":"meta","type":{"map":{"keys":"string","values":{"vector":{"items":"string"}}}}}]},{"name":"Sample","fields":[{"name":"data","type":{"vector":{"items":"int32"}}},{"name":"timestamp","type":"datetime"}]},{"name":"Signature","fields":[{"name":"name","type":"string"},{"name":"email","type":"string"},{"name":"number","type":"int64"}]}]})",
};

void MyProtocolWriterBase::WriteHeader(evo_test::Header const& value) {
  if (unlikely(state_ != 0)) {
    MyProtocolWriterBaseInvalidState(0, false, state_);
  }

  WriteHeaderImpl(value);
  state_ = 1;
}

void MyProtocolWriterBase::WriteId(int64_t const& value) {
  if (unlikely(state_ != 1)) {
    MyProtocolWriterBaseInvalidState(1, false, state_);
  }

  WriteIdImpl(value);
  state_ = 2;
}

void MyProtocolWriterBase::WriteSamples(evo_test::Sample const& value) {
  if (unlikely(state_ != 2)) {
    MyProtocolWriterBaseInvalidState(2, false, state_);
  }

  WriteSamplesImpl(value);
}

void MyProtocolWriterBase::WriteSamples(std::vector<evo_test::Sample> const& values) {
  if (unlikely(state_ != 2)) {
    MyProtocolWriterBaseInvalidState(2, false, state_);
  }

  WriteSamplesImpl(values);
}

void MyProtocolWriterBase::EndSamples() {
  if (unlikely(state_ != 2)) {
    MyProtocolWriterBaseInvalidState(2, true, state_);
  }

  EndSamplesImpl();
  state_ = 3;
}

// fallback implementation
void MyProtocolWriterBase::WriteSamplesImpl(std::vector<evo_test::Sample> const& values) {
  for (auto const& v : values) {
    WriteSamplesImpl(v);
  }
}

void MyProtocolWriterBase::WriteFooter(std::optional<evo_test::Footer> const& value) {
  if (unlikely(state_ != 3)) {
    MyProtocolWriterBaseInvalidState(3, false, state_);
  }

  WriteFooterImpl(value);
  state_ = 4;
}

void MyProtocolWriterBase::Close() {
  if (unlikely(state_ != 4)) {
    MyProtocolWriterBaseInvalidState(4, false, state_);
  }

  CloseImpl();
}

std::string MyProtocolReaderBase::schema_ = MyProtocolWriterBase::schema_;

std::vector<std::string> MyProtocolReaderBase::previous_schemas_ = MyProtocolWriterBase::previous_schemas_;

void MyProtocolReaderBase::ReadHeader(evo_test::Header& value) {
  if (unlikely(state_ != 0)) {
    MyProtocolReaderBaseInvalidState(0, state_);
  }

  ReadHeaderImpl(value);
  state_ = 2;
}

void MyProtocolReaderBase::ReadId(int64_t& value) {
  if (unlikely(state_ != 2)) {
    MyProtocolReaderBaseInvalidState(2, state_);
  }

  ReadIdImpl(value);
  state_ = 4;
}

bool MyProtocolReaderBase::ReadSamples(evo_test::Sample& value) {
  if (unlikely(state_ != 4)) {
    if (state_ == 5) {
      state_ = 6;
      return false;
    }
    MyProtocolReaderBaseInvalidState(4, state_);
  }

  bool result = ReadSamplesImpl(value);
  if (!result) {
    state_ = 6;
  }
  return result;
}

bool MyProtocolReaderBase::ReadSamples(std::vector<evo_test::Sample>& values) {
  if (values.capacity() == 0) {
    throw std::runtime_error("vector must have a nonzero capacity.");
  }
  if (unlikely(state_ != 4)) {
    if (state_ == 5) {
      state_ = 6;
      values.clear();
      return false;
    }
    MyProtocolReaderBaseInvalidState(4, state_);
  }

  if (!ReadSamplesImpl(values)) {
    state_ = 5;
    return values.size() > 0;
  }
  return true;
}

// fallback implementation
bool MyProtocolReaderBase::ReadSamplesImpl(std::vector<evo_test::Sample>& values) {
  size_t i = 0;
  while (true) {
    if (i == values.size()) {
      values.resize(i + 1);
    }
    if (!ReadSamplesImpl(values[i])) {
      values.resize(i);
      return false;
    }
    i++;
    if (i == values.capacity()) {
      return true;
    }
  }
}

void MyProtocolReaderBase::ReadFooter(std::optional<evo_test::Footer>& value) {
  if (unlikely(state_ != 6)) {
    if (state_ == 5) {
      state_ = 6;
    } else {
      MyProtocolReaderBaseInvalidState(6, state_);
    }
  }

  ReadFooterImpl(value);
  state_ = 8;
}

void MyProtocolReaderBase::Close() {
  if (unlikely(state_ != 8)) {
    MyProtocolReaderBaseInvalidState(8, state_);
  }

  CloseImpl();
}
void MyProtocolReaderBase::CopyTo(MyProtocolWriterBase& writer, size_t samples_buffer_size) {
  {
    evo_test::Header value;
    ReadHeader(value);
    writer.WriteHeader(value);
  }
  {
    int64_t value;
    ReadId(value);
    writer.WriteId(value);
  }
  if (samples_buffer_size > 1) {
    std::vector<evo_test::Sample> values;
    values.reserve(samples_buffer_size);
    while(ReadSamples(values)) {
      writer.WriteSamples(values);
    }
    writer.EndSamples();
  } else {
    evo_test::Sample value;
    while(ReadSamples(value)) {
      writer.WriteSamples(value);
    }
    writer.EndSamples();
  }
  {
    std::optional<evo_test::Footer> value;
    ReadFooter(value);
    writer.WriteFooter(value);
  }
}

namespace {
void UnusedProtocolWriterBaseInvalidState(uint8_t attempted, [[maybe_unused]] bool end, uint8_t current) {
  std::string expected_method;
  switch (current) {
  case 0: expected_method = "WriteSamples() or EndSamples()"; break;
  }
  std::string attempted_method;
  switch (attempted) {
  case 0: attempted_method = end ? "EndSamples()" : "WriteSamples()"; break;
  case 1: attempted_method = "Close()"; break;
  }
  throw std::runtime_error("Expected call to " + expected_method + " but received call to " + attempted_method + " instead.");
}

void UnusedProtocolReaderBaseInvalidState(uint8_t attempted, uint8_t current) {
  auto f = [](uint8_t i) -> std::string {
    switch (i/2) {
    case 0: return "ReadSamples()";
    case 1: return "Close()";
    default: return "<unknown>";
    }
  };
  throw std::runtime_error("Expected call to " + f(current) + " but received call to " + f(attempted) + " instead.");
}

} // namespace 

std::string UnusedProtocolWriterBase::schema_ = R"({"protocol":{"name":"UnusedProtocol","sequence":[{"name":"samples","type":{"stream":{"items":"EvoTest.Sample"}}}]},"types":[{"name":"Sample","fields":[{"name":"data","type":{"vector":{"items":"int32"}}},{"name":"timestamp","type":"datetime"}]}]})";

std::vector<std::string> UnusedProtocolWriterBase::previous_schemas_ = {
  R"()",
};

void UnusedProtocolWriterBase::WriteSamples(evo_test::Sample const& value) {
  if (unlikely(state_ != 0)) {
    UnusedProtocolWriterBaseInvalidState(0, false, state_);
  }

  WriteSamplesImpl(value);
}

void UnusedProtocolWriterBase::WriteSamples(std::vector<evo_test::Sample> const& values) {
  if (unlikely(state_ != 0)) {
    UnusedProtocolWriterBaseInvalidState(0, false, state_);
  }

  WriteSamplesImpl(values);
}

void UnusedProtocolWriterBase::EndSamples() {
  if (unlikely(state_ != 0)) {
    UnusedProtocolWriterBaseInvalidState(0, true, state_);
  }

  EndSamplesImpl();
  state_ = 1;
}

// fallback implementation
void UnusedProtocolWriterBase::WriteSamplesImpl(std::vector<evo_test::Sample> const& values) {
  for (auto const& v : values) {
    WriteSamplesImpl(v);
  }
}

void UnusedProtocolWriterBase::Close() {
  if (unlikely(state_ != 1)) {
    UnusedProtocolWriterBaseInvalidState(1, false, state_);
  }

  CloseImpl();
}

std::string UnusedProtocolReaderBase::schema_ = UnusedProtocolWriterBase::schema_;

std::vector<std::string> UnusedProtocolReaderBase::previous_schemas_ = UnusedProtocolWriterBase::previous_schemas_;

bool UnusedProtocolReaderBase::ReadSamples(evo_test::Sample& value) {
  if (unlikely(state_ != 0)) {
    if (state_ == 1) {
      state_ = 2;
      return false;
    }
    UnusedProtocolReaderBaseInvalidState(0, state_);
  }

  bool result = ReadSamplesImpl(value);
  if (!result) {
    state_ = 2;
  }
  return result;
}

bool UnusedProtocolReaderBase::ReadSamples(std::vector<evo_test::Sample>& values) {
  if (values.capacity() == 0) {
    throw std::runtime_error("vector must have a nonzero capacity.");
  }
  if (unlikely(state_ != 0)) {
    if (state_ == 1) {
      state_ = 2;
      values.clear();
      return false;
    }
    UnusedProtocolReaderBaseInvalidState(0, state_);
  }

  if (!ReadSamplesImpl(values)) {
    state_ = 1;
    return values.size() > 0;
  }
  return true;
}

// fallback implementation
bool UnusedProtocolReaderBase::ReadSamplesImpl(std::vector<evo_test::Sample>& values) {
  size_t i = 0;
  while (true) {
    if (i == values.size()) {
      values.resize(i + 1);
    }
    if (!ReadSamplesImpl(values[i])) {
      values.resize(i);
      return false;
    }
    i++;
    if (i == values.capacity()) {
      return true;
    }
  }
}

void UnusedProtocolReaderBase::Close() {
  if (unlikely(state_ != 2)) {
    UnusedProtocolReaderBaseInvalidState(2, state_);
  }

  CloseImpl();
}
void UnusedProtocolReaderBase::CopyTo(UnusedProtocolWriterBase& writer, size_t samples_buffer_size) {
  if (samples_buffer_size > 1) {
    std::vector<evo_test::Sample> values;
    values.reserve(samples_buffer_size);
    while(ReadSamples(values)) {
      writer.WriteSamples(values);
    }
    writer.EndSamples();
  } else {
    evo_test::Sample value;
    while(ReadSamples(value)) {
      writer.WriteSamples(value);
    }
    writer.EndSamples();
  }
}
} // namespace evo_test
