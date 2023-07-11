// This file was generated by the "yardl" tool. DO NOT EDIT.

#include <iostream>

#include "../format.h"
#include "binary/protocols.h"
#include "hdf5/protocols.h"
#include "ndjson/protocols.h"

namespace yardl::testing {
void TranslateStream(std::string const& protocol_name, yardl::testing::Format input_format, std::istream& input, yardl::testing::Format output_format, std::ostream& output) {
  switch (input_format) {
  case yardl::testing::Format::kBinary:
    break;
  case yardl::testing::Format::kNDJson:
    break;
  default:
    throw std::runtime_error("Unsupported input format");
  }

  if (protocol_name == "BenchmarkFloat256x256") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::BenchmarkFloat256x256ReaderBase>(new test_model::binary::BenchmarkFloat256x256Reader(input))
      : std::unique_ptr<test_model::BenchmarkFloat256x256ReaderBase>(new test_model::ndjson::BenchmarkFloat256x256Reader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::BenchmarkFloat256x256WriterBase>(new test_model::binary::BenchmarkFloat256x256Writer(output))
      : std::unique_ptr<test_model::BenchmarkFloat256x256WriterBase>(new test_model::ndjson::BenchmarkFloat256x256Writer(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "BenchmarkFloatVlen") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::BenchmarkFloatVlenReaderBase>(new test_model::binary::BenchmarkFloatVlenReader(input))
      : std::unique_ptr<test_model::BenchmarkFloatVlenReaderBase>(new test_model::ndjson::BenchmarkFloatVlenReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::BenchmarkFloatVlenWriterBase>(new test_model::binary::BenchmarkFloatVlenWriter(output))
      : std::unique_ptr<test_model::BenchmarkFloatVlenWriterBase>(new test_model::ndjson::BenchmarkFloatVlenWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "BenchmarkSmallRecord") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::BenchmarkSmallRecordReaderBase>(new test_model::binary::BenchmarkSmallRecordReader(input))
      : std::unique_ptr<test_model::BenchmarkSmallRecordReaderBase>(new test_model::ndjson::BenchmarkSmallRecordReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::BenchmarkSmallRecordWriterBase>(new test_model::binary::BenchmarkSmallRecordWriter(output))
      : std::unique_ptr<test_model::BenchmarkSmallRecordWriterBase>(new test_model::ndjson::BenchmarkSmallRecordWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "BenchmarkSmallRecordWithOptionals") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::BenchmarkSmallRecordWithOptionalsReaderBase>(new test_model::binary::BenchmarkSmallRecordWithOptionalsReader(input))
      : std::unique_ptr<test_model::BenchmarkSmallRecordWithOptionalsReaderBase>(new test_model::ndjson::BenchmarkSmallRecordWithOptionalsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::BenchmarkSmallRecordWithOptionalsWriterBase>(new test_model::binary::BenchmarkSmallRecordWithOptionalsWriter(output))
      : std::unique_ptr<test_model::BenchmarkSmallRecordWithOptionalsWriterBase>(new test_model::ndjson::BenchmarkSmallRecordWithOptionalsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "BenchmarkSimpleMrd") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::BenchmarkSimpleMrdReaderBase>(new test_model::binary::BenchmarkSimpleMrdReader(input))
      : std::unique_ptr<test_model::BenchmarkSimpleMrdReaderBase>(new test_model::ndjson::BenchmarkSimpleMrdReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::BenchmarkSimpleMrdWriterBase>(new test_model::binary::BenchmarkSimpleMrdWriter(output))
      : std::unique_ptr<test_model::BenchmarkSimpleMrdWriterBase>(new test_model::ndjson::BenchmarkSimpleMrdWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "Scalars") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::ScalarsReaderBase>(new test_model::binary::ScalarsReader(input))
      : std::unique_ptr<test_model::ScalarsReaderBase>(new test_model::ndjson::ScalarsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::ScalarsWriterBase>(new test_model::binary::ScalarsWriter(output))
      : std::unique_ptr<test_model::ScalarsWriterBase>(new test_model::ndjson::ScalarsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "ScalarOptionals") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::ScalarOptionalsReaderBase>(new test_model::binary::ScalarOptionalsReader(input))
      : std::unique_ptr<test_model::ScalarOptionalsReaderBase>(new test_model::ndjson::ScalarOptionalsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::ScalarOptionalsWriterBase>(new test_model::binary::ScalarOptionalsWriter(output))
      : std::unique_ptr<test_model::ScalarOptionalsWriterBase>(new test_model::ndjson::ScalarOptionalsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "NestedRecords") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::NestedRecordsReaderBase>(new test_model::binary::NestedRecordsReader(input))
      : std::unique_ptr<test_model::NestedRecordsReaderBase>(new test_model::ndjson::NestedRecordsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::NestedRecordsWriterBase>(new test_model::binary::NestedRecordsWriter(output))
      : std::unique_ptr<test_model::NestedRecordsWriterBase>(new test_model::ndjson::NestedRecordsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "Vlens") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::VlensReaderBase>(new test_model::binary::VlensReader(input))
      : std::unique_ptr<test_model::VlensReaderBase>(new test_model::ndjson::VlensReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::VlensWriterBase>(new test_model::binary::VlensWriter(output))
      : std::unique_ptr<test_model::VlensWriterBase>(new test_model::ndjson::VlensWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "Strings") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::StringsReaderBase>(new test_model::binary::StringsReader(input))
      : std::unique_ptr<test_model::StringsReaderBase>(new test_model::ndjson::StringsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::StringsWriterBase>(new test_model::binary::StringsWriter(output))
      : std::unique_ptr<test_model::StringsWriterBase>(new test_model::ndjson::StringsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "OptionalVectors") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::OptionalVectorsReaderBase>(new test_model::binary::OptionalVectorsReader(input))
      : std::unique_ptr<test_model::OptionalVectorsReaderBase>(new test_model::ndjson::OptionalVectorsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::OptionalVectorsWriterBase>(new test_model::binary::OptionalVectorsWriter(output))
      : std::unique_ptr<test_model::OptionalVectorsWriterBase>(new test_model::ndjson::OptionalVectorsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "FixedVectors") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::FixedVectorsReaderBase>(new test_model::binary::FixedVectorsReader(input))
      : std::unique_ptr<test_model::FixedVectorsReaderBase>(new test_model::ndjson::FixedVectorsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::FixedVectorsWriterBase>(new test_model::binary::FixedVectorsWriter(output))
      : std::unique_ptr<test_model::FixedVectorsWriterBase>(new test_model::ndjson::FixedVectorsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "Streams") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::StreamsReaderBase>(new test_model::binary::StreamsReader(input))
      : std::unique_ptr<test_model::StreamsReaderBase>(new test_model::ndjson::StreamsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::StreamsWriterBase>(new test_model::binary::StreamsWriter(output))
      : std::unique_ptr<test_model::StreamsWriterBase>(new test_model::ndjson::StreamsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "FixedArrays") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::FixedArraysReaderBase>(new test_model::binary::FixedArraysReader(input))
      : std::unique_ptr<test_model::FixedArraysReaderBase>(new test_model::ndjson::FixedArraysReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::FixedArraysWriterBase>(new test_model::binary::FixedArraysWriter(output))
      : std::unique_ptr<test_model::FixedArraysWriterBase>(new test_model::ndjson::FixedArraysWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "NDArrays") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::NDArraysReaderBase>(new test_model::binary::NDArraysReader(input))
      : std::unique_ptr<test_model::NDArraysReaderBase>(new test_model::ndjson::NDArraysReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::NDArraysWriterBase>(new test_model::binary::NDArraysWriter(output))
      : std::unique_ptr<test_model::NDArraysWriterBase>(new test_model::ndjson::NDArraysWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "NDArraysSingleDimension") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::NDArraysSingleDimensionReaderBase>(new test_model::binary::NDArraysSingleDimensionReader(input))
      : std::unique_ptr<test_model::NDArraysSingleDimensionReaderBase>(new test_model::ndjson::NDArraysSingleDimensionReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::NDArraysSingleDimensionWriterBase>(new test_model::binary::NDArraysSingleDimensionWriter(output))
      : std::unique_ptr<test_model::NDArraysSingleDimensionWriterBase>(new test_model::ndjson::NDArraysSingleDimensionWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "DynamicNDArrays") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::DynamicNDArraysReaderBase>(new test_model::binary::DynamicNDArraysReader(input))
      : std::unique_ptr<test_model::DynamicNDArraysReaderBase>(new test_model::ndjson::DynamicNDArraysReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::DynamicNDArraysWriterBase>(new test_model::binary::DynamicNDArraysWriter(output))
      : std::unique_ptr<test_model::DynamicNDArraysWriterBase>(new test_model::ndjson::DynamicNDArraysWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "Maps") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::MapsReaderBase>(new test_model::binary::MapsReader(input))
      : std::unique_ptr<test_model::MapsReaderBase>(new test_model::ndjson::MapsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::MapsWriterBase>(new test_model::binary::MapsWriter(output))
      : std::unique_ptr<test_model::MapsWriterBase>(new test_model::ndjson::MapsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "Unions") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::UnionsReaderBase>(new test_model::binary::UnionsReader(input))
      : std::unique_ptr<test_model::UnionsReaderBase>(new test_model::ndjson::UnionsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::UnionsWriterBase>(new test_model::binary::UnionsWriter(output))
      : std::unique_ptr<test_model::UnionsWriterBase>(new test_model::ndjson::UnionsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "StreamsOfUnions") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::StreamsOfUnionsReaderBase>(new test_model::binary::StreamsOfUnionsReader(input))
      : std::unique_ptr<test_model::StreamsOfUnionsReaderBase>(new test_model::ndjson::StreamsOfUnionsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::StreamsOfUnionsWriterBase>(new test_model::binary::StreamsOfUnionsWriter(output))
      : std::unique_ptr<test_model::StreamsOfUnionsWriterBase>(new test_model::ndjson::StreamsOfUnionsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "Enums") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::EnumsReaderBase>(new test_model::binary::EnumsReader(input))
      : std::unique_ptr<test_model::EnumsReaderBase>(new test_model::ndjson::EnumsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::EnumsWriterBase>(new test_model::binary::EnumsWriter(output))
      : std::unique_ptr<test_model::EnumsWriterBase>(new test_model::ndjson::EnumsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "Flags") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::FlagsReaderBase>(new test_model::binary::FlagsReader(input))
      : std::unique_ptr<test_model::FlagsReaderBase>(new test_model::ndjson::FlagsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::FlagsWriterBase>(new test_model::binary::FlagsWriter(output))
      : std::unique_ptr<test_model::FlagsWriterBase>(new test_model::ndjson::FlagsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "StateTest") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::StateTestReaderBase>(new test_model::binary::StateTestReader(input))
      : std::unique_ptr<test_model::StateTestReaderBase>(new test_model::ndjson::StateTestReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::StateTestWriterBase>(new test_model::binary::StateTestWriter(output))
      : std::unique_ptr<test_model::StateTestWriterBase>(new test_model::ndjson::StateTestWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "SimpleGenerics") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::SimpleGenericsReaderBase>(new test_model::binary::SimpleGenericsReader(input))
      : std::unique_ptr<test_model::SimpleGenericsReaderBase>(new test_model::ndjson::SimpleGenericsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::SimpleGenericsWriterBase>(new test_model::binary::SimpleGenericsWriter(output))
      : std::unique_ptr<test_model::SimpleGenericsWriterBase>(new test_model::ndjson::SimpleGenericsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "AdvancedGenerics") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::AdvancedGenericsReaderBase>(new test_model::binary::AdvancedGenericsReader(input))
      : std::unique_ptr<test_model::AdvancedGenericsReaderBase>(new test_model::ndjson::AdvancedGenericsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::AdvancedGenericsWriterBase>(new test_model::binary::AdvancedGenericsWriter(output))
      : std::unique_ptr<test_model::AdvancedGenericsWriterBase>(new test_model::ndjson::AdvancedGenericsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "Aliases") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::AliasesReaderBase>(new test_model::binary::AliasesReader(input))
      : std::unique_ptr<test_model::AliasesReaderBase>(new test_model::ndjson::AliasesReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::AliasesWriterBase>(new test_model::binary::AliasesWriter(output))
      : std::unique_ptr<test_model::AliasesWriterBase>(new test_model::ndjson::AliasesWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "StreamsOfAliasedUnions") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::StreamsOfAliasedUnionsReaderBase>(new test_model::binary::StreamsOfAliasedUnionsReader(input))
      : std::unique_ptr<test_model::StreamsOfAliasedUnionsReaderBase>(new test_model::ndjson::StreamsOfAliasedUnionsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::StreamsOfAliasedUnionsWriterBase>(new test_model::binary::StreamsOfAliasedUnionsWriter(output))
      : std::unique_ptr<test_model::StreamsOfAliasedUnionsWriterBase>(new test_model::ndjson::StreamsOfAliasedUnionsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "ProtocolWithComputedFields") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::ProtocolWithComputedFieldsReaderBase>(new test_model::binary::ProtocolWithComputedFieldsReader(input))
      : std::unique_ptr<test_model::ProtocolWithComputedFieldsReaderBase>(new test_model::ndjson::ProtocolWithComputedFieldsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::ProtocolWithComputedFieldsWriterBase>(new test_model::binary::ProtocolWithComputedFieldsWriter(output))
      : std::unique_ptr<test_model::ProtocolWithComputedFieldsWriterBase>(new test_model::ndjson::ProtocolWithComputedFieldsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  if (protocol_name == "ProtocolWithKeywordSteps") {
    auto reader = input_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::ProtocolWithKeywordStepsReaderBase>(new test_model::binary::ProtocolWithKeywordStepsReader(input))
      : std::unique_ptr<test_model::ProtocolWithKeywordStepsReaderBase>(new test_model::ndjson::ProtocolWithKeywordStepsReader(input));

    auto writer = output_format == yardl::testing::Format::kBinary
      ? std::unique_ptr<test_model::ProtocolWithKeywordStepsWriterBase>(new test_model::binary::ProtocolWithKeywordStepsWriter(output))
      : std::unique_ptr<test_model::ProtocolWithKeywordStepsWriterBase>(new test_model::ndjson::ProtocolWithKeywordStepsWriter(output));
    reader->CopyTo(*writer);
    return;
  }
  throw std::runtime_error("Unsupported protocol " + protocol_name);
}
} // namespace yardl::testing
