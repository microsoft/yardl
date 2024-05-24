// This file was generated by the "yardl" tool. DO NOT EDIT.

#include <functional>
#include "../factories.h"
#include "binary/protocols.h"
#include "hdf5/protocols.h"
#include "ndjson/protocols.h"

namespace yardl::testing {
template<>
std::unique_ptr<test_model::BenchmarkFloat256x256WriterBase> CreateWriter<test_model::BenchmarkFloat256x256WriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkFloat256x256Writer>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkFloat256x256Writer>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkFloat256x256Writer>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::BenchmarkFloat256x256ReaderBase> CreateReader<test_model::BenchmarkFloat256x256ReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkFloat256x256Reader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkFloat256x256Reader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkFloat256x256Reader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::BenchmarkInt256x256WriterBase> CreateWriter<test_model::BenchmarkInt256x256WriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkInt256x256Writer>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkInt256x256Writer>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkInt256x256Writer>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::BenchmarkInt256x256ReaderBase> CreateReader<test_model::BenchmarkInt256x256ReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkInt256x256Reader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkInt256x256Reader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkInt256x256Reader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::BenchmarkFloatVlenWriterBase> CreateWriter<test_model::BenchmarkFloatVlenWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkFloatVlenWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkFloatVlenWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkFloatVlenWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::BenchmarkFloatVlenReaderBase> CreateReader<test_model::BenchmarkFloatVlenReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkFloatVlenReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkFloatVlenReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkFloatVlenReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::BenchmarkSmallRecordWriterBase> CreateWriter<test_model::BenchmarkSmallRecordWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkSmallRecordWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkSmallRecordWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkSmallRecordWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::BenchmarkSmallRecordReaderBase> CreateReader<test_model::BenchmarkSmallRecordReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkSmallRecordReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkSmallRecordReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkSmallRecordReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::BenchmarkSmallRecordWithOptionalsWriterBase> CreateWriter<test_model::BenchmarkSmallRecordWithOptionalsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkSmallRecordWithOptionalsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkSmallRecordWithOptionalsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkSmallRecordWithOptionalsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::BenchmarkSmallRecordWithOptionalsReaderBase> CreateReader<test_model::BenchmarkSmallRecordWithOptionalsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkSmallRecordWithOptionalsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkSmallRecordWithOptionalsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkSmallRecordWithOptionalsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::BenchmarkSimpleMrdWriterBase> CreateWriter<test_model::BenchmarkSimpleMrdWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkSimpleMrdWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkSimpleMrdWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkSimpleMrdWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::BenchmarkSimpleMrdReaderBase> CreateReader<test_model::BenchmarkSimpleMrdReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::BenchmarkSimpleMrdReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::BenchmarkSimpleMrdReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::BenchmarkSimpleMrdReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::ScalarsWriterBase> CreateWriter<test_model::ScalarsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::ScalarsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::ScalarsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::ScalarsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::ScalarsReaderBase> CreateReader<test_model::ScalarsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::ScalarsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::ScalarsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::ScalarsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::ScalarOptionalsWriterBase> CreateWriter<test_model::ScalarOptionalsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::ScalarOptionalsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::ScalarOptionalsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::ScalarOptionalsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::ScalarOptionalsReaderBase> CreateReader<test_model::ScalarOptionalsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::ScalarOptionalsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::ScalarOptionalsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::ScalarOptionalsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::NestedRecordsWriterBase> CreateWriter<test_model::NestedRecordsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::NestedRecordsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::NestedRecordsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::NestedRecordsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::NestedRecordsReaderBase> CreateReader<test_model::NestedRecordsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::NestedRecordsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::NestedRecordsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::NestedRecordsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::VlensWriterBase> CreateWriter<test_model::VlensWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::VlensWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::VlensWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::VlensWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::VlensReaderBase> CreateReader<test_model::VlensReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::VlensReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::VlensReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::VlensReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::StringsWriterBase> CreateWriter<test_model::StringsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::StringsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::StringsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::StringsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::StringsReaderBase> CreateReader<test_model::StringsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::StringsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::StringsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::StringsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::OptionalVectorsWriterBase> CreateWriter<test_model::OptionalVectorsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::OptionalVectorsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::OptionalVectorsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::OptionalVectorsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::OptionalVectorsReaderBase> CreateReader<test_model::OptionalVectorsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::OptionalVectorsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::OptionalVectorsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::OptionalVectorsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::FixedVectorsWriterBase> CreateWriter<test_model::FixedVectorsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::FixedVectorsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::FixedVectorsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::FixedVectorsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::FixedVectorsReaderBase> CreateReader<test_model::FixedVectorsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::FixedVectorsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::FixedVectorsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::FixedVectorsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::StreamsWriterBase> CreateWriter<test_model::StreamsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::StreamsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::StreamsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::StreamsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::StreamsReaderBase> CreateReader<test_model::StreamsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::StreamsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::StreamsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::StreamsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::FixedArraysWriterBase> CreateWriter<test_model::FixedArraysWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::FixedArraysWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::FixedArraysWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::FixedArraysWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::FixedArraysReaderBase> CreateReader<test_model::FixedArraysReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::FixedArraysReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::FixedArraysReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::FixedArraysReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::SubarraysWriterBase> CreateWriter<test_model::SubarraysWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::SubarraysWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::SubarraysWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::SubarraysWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::SubarraysReaderBase> CreateReader<test_model::SubarraysReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::SubarraysReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::SubarraysReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::SubarraysReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::SubarraysInRecordsWriterBase> CreateWriter<test_model::SubarraysInRecordsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::SubarraysInRecordsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::SubarraysInRecordsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::SubarraysInRecordsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::SubarraysInRecordsReaderBase> CreateReader<test_model::SubarraysInRecordsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::SubarraysInRecordsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::SubarraysInRecordsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::SubarraysInRecordsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::NDArraysWriterBase> CreateWriter<test_model::NDArraysWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::NDArraysWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::NDArraysWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::NDArraysWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::NDArraysReaderBase> CreateReader<test_model::NDArraysReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::NDArraysReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::NDArraysReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::NDArraysReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::NDArraysSingleDimensionWriterBase> CreateWriter<test_model::NDArraysSingleDimensionWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::NDArraysSingleDimensionWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::NDArraysSingleDimensionWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::NDArraysSingleDimensionWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::NDArraysSingleDimensionReaderBase> CreateReader<test_model::NDArraysSingleDimensionReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::NDArraysSingleDimensionReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::NDArraysSingleDimensionReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::NDArraysSingleDimensionReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::DynamicNDArraysWriterBase> CreateWriter<test_model::DynamicNDArraysWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::DynamicNDArraysWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::DynamicNDArraysWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::DynamicNDArraysWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::DynamicNDArraysReaderBase> CreateReader<test_model::DynamicNDArraysReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::DynamicNDArraysReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::DynamicNDArraysReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::DynamicNDArraysReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::MultiDArraysWriterBase> CreateWriter<test_model::MultiDArraysWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::MultiDArraysWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::MultiDArraysWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::MultiDArraysWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::MultiDArraysReaderBase> CreateReader<test_model::MultiDArraysReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::MultiDArraysReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::MultiDArraysReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::MultiDArraysReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::ComplexArraysWriterBase> CreateWriter<test_model::ComplexArraysWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::ComplexArraysWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::ComplexArraysWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::ComplexArraysWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::ComplexArraysReaderBase> CreateReader<test_model::ComplexArraysReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::ComplexArraysReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::ComplexArraysReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::ComplexArraysReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::MapsWriterBase> CreateWriter<test_model::MapsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::MapsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::MapsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::MapsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::MapsReaderBase> CreateReader<test_model::MapsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::MapsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::MapsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::MapsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::UnionsWriterBase> CreateWriter<test_model::UnionsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::UnionsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::UnionsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::UnionsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::UnionsReaderBase> CreateReader<test_model::UnionsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::UnionsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::UnionsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::UnionsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::StreamsOfUnionsWriterBase> CreateWriter<test_model::StreamsOfUnionsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::StreamsOfUnionsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::StreamsOfUnionsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::StreamsOfUnionsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::StreamsOfUnionsReaderBase> CreateReader<test_model::StreamsOfUnionsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::StreamsOfUnionsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::StreamsOfUnionsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::StreamsOfUnionsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::EnumsWriterBase> CreateWriter<test_model::EnumsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::EnumsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::EnumsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::EnumsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::EnumsReaderBase> CreateReader<test_model::EnumsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::EnumsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::EnumsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::EnumsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::FlagsWriterBase> CreateWriter<test_model::FlagsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::FlagsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::FlagsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::FlagsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::FlagsReaderBase> CreateReader<test_model::FlagsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::FlagsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::FlagsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::FlagsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::StateTestWriterBase> CreateWriter<test_model::StateTestWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::StateTestWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::StateTestWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::StateTestWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::StateTestReaderBase> CreateReader<test_model::StateTestReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::StateTestReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::StateTestReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::StateTestReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::SimpleGenericsWriterBase> CreateWriter<test_model::SimpleGenericsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::SimpleGenericsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::SimpleGenericsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::SimpleGenericsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::SimpleGenericsReaderBase> CreateReader<test_model::SimpleGenericsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::SimpleGenericsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::SimpleGenericsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::SimpleGenericsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::AdvancedGenericsWriterBase> CreateWriter<test_model::AdvancedGenericsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::AdvancedGenericsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::AdvancedGenericsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::AdvancedGenericsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::AdvancedGenericsReaderBase> CreateReader<test_model::AdvancedGenericsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::AdvancedGenericsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::AdvancedGenericsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::AdvancedGenericsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::AliasesWriterBase> CreateWriter<test_model::AliasesWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::AliasesWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::AliasesWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::AliasesWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::AliasesReaderBase> CreateReader<test_model::AliasesReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::AliasesReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::AliasesReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::AliasesReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::StreamsOfAliasedUnionsWriterBase> CreateWriter<test_model::StreamsOfAliasedUnionsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::StreamsOfAliasedUnionsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::StreamsOfAliasedUnionsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::StreamsOfAliasedUnionsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::StreamsOfAliasedUnionsReaderBase> CreateReader<test_model::StreamsOfAliasedUnionsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::StreamsOfAliasedUnionsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::StreamsOfAliasedUnionsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::StreamsOfAliasedUnionsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::ProtocolWithComputedFieldsWriterBase> CreateWriter<test_model::ProtocolWithComputedFieldsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::ProtocolWithComputedFieldsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::ProtocolWithComputedFieldsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::ProtocolWithComputedFieldsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::ProtocolWithComputedFieldsReaderBase> CreateReader<test_model::ProtocolWithComputedFieldsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::ProtocolWithComputedFieldsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::ProtocolWithComputedFieldsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::ProtocolWithComputedFieldsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::ProtocolWithKeywordStepsWriterBase> CreateWriter<test_model::ProtocolWithKeywordStepsWriterBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::ProtocolWithKeywordStepsWriter>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::ProtocolWithKeywordStepsWriter>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::ProtocolWithKeywordStepsWriter>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

template<>
std::unique_ptr<test_model::ProtocolWithKeywordStepsReaderBase> CreateReader<test_model::ProtocolWithKeywordStepsReaderBase>(Format format, std::string const& filename) {
  switch (format) {
  case Format::kHdf5:
    return std::make_unique<test_model::hdf5::ProtocolWithKeywordStepsReader>(filename);
  case Format::kBinary:
    return std::make_unique<test_model::binary::ProtocolWithKeywordStepsReader>(filename);
  case Format::kNDJson:
    return std::make_unique<test_model::ndjson::ProtocolWithKeywordStepsReader>(filename);
  default:
    throw std::runtime_error("Unknown format");
  }
}

}
