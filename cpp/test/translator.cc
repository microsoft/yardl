#include <algorithm>
#include <iostream>
#include <sstream>

#include <nlohmann/json.hpp>

#include "format.h"
#include "generated/binary/protocols.h"
#include "generated/ndjson/protocols.h"
#include "generated/yardl/detail/binary/header.h"
#include "generated/yardl/detail/ndjson/header.h"

using ordered_json = nlohmann::ordered_json;
using yardl::testing::Format;

namespace yardl::testing {
void TranslateStream(std::string const& protocol_name, Format input_format, std::istream& input, Format output_format, std::ostream& output);
}

namespace {
Format parseFormat(std::string arg) {
  std::transform(arg.begin(), arg.end(), arg.begin(), [](unsigned char c) { return std::toupper(c); });

  if (arg == "BINARY") {
    return Format::kBinary;
  } else if (arg == "NDJSON") {
    return Format::kNDJson;
  }

  std::cerr << "Invalid format '" << arg << "'. Expected 'binary' or 'ndjson'" << std::endl;
  exit(1);
}

std::string GetProtocolName(std::istream& input, Format format) {
  ordered_json parsed_schema;
  switch (format) {
    case Format::kBinary: {
      yardl::binary::CodedInputStream input_stream(input);
      std::string schema = yardl::binary::ReadHeader(input_stream);
      parsed_schema = ordered_json::parse(schema);
      break;
    }
    case Format::kNDJson:
      parsed_schema = yardl::ndjson::ReadHeader(input);
      break;
    default:
      throw std::runtime_error("Format not supported");
  }

  return parsed_schema["protocol"]["name"].get<std::string>();
}

}  // namespace

int main(int argc, char* argv[]) {
  if (argc != 3) {  // Check if argument count is correct
    std::cerr << "Incorrect number of arguments. Usage: translator <binary | ndjson> <binary | ndjson>" << std::endl;
    return 1;
  }

  Format inputFormat = parseFormat(argv[1]);
  Format outputFormat = parseFormat(argv[2]);

  std::stringstream buffered_input;
  buffered_input << std::cin.rdbuf();
  if (buffered_input.fail()) {
    std::cerr << "Failed to buffer input" << std::endl;
    return 1;
  }

  std::string protocol_name = GetProtocolName(buffered_input, inputFormat);

  buffered_input.clear();  // clear possible failbit from last read where we may have reached EOF
  buffered_input.seekg(0);

  yardl::testing::TranslateStream(protocol_name, inputFormat, buffered_input, outputFormat, std::cout);

  return 0;
}
