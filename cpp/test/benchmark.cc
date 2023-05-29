// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.#include "benchmark.h"

// Some very basic throughput benchmarks.

#include "benchmark.h"

#include "generated/binary/protocols.h"
#include "generated/hdf5/protocols.h"
#include "generated/ndjson/protocols.h"

using namespace test_model;

using namespace yardl::testing;

namespace {

std::string const kOutputFileName = "/tmp/benchmark_data.dat";

void AssertRepetitionsSame(size_t expected, size_t actual) {
  if (expected != actual) {
    throw std::runtime_error("Expected " + std::to_string(expected) +
                             " repetitions, but got " + std::to_string(actual));
  }
}

template <typename TWriter>
constexpr size_t ScaleRepetitions(size_t repetitions) {
  if constexpr (std::is_base_of_v<yardl::ndjson::NDJsonWriter, TWriter>) {
    return repetitions / 200;
  }

  return repetitions;
}

template <typename TWriter, typename TReader>
void BenchmarkFloat256x256() {
  std::remove(kOutputFileName.c_str());

  yardl::FixedNDArray<float, 256, 256> a;
  int i = 0;
  for (auto& x : a) {
    x = static_cast<float>(++i) - std::numeric_limits<float>::epsilon();
  }

  size_t const repetitions = ScaleRepetitions<TWriter>(10000);
  size_t const total_size = sizeof(a) * repetitions;
  TimeScenario<TWriter>(
      __FUNCTION__,
      total_size,
      [&]() {
        std::unique_ptr<BenchmarkFloat256x256WriterBase> writer = std::make_unique<TWriter>(kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteFloat256x256(a);
        }
        writer->EndFloat256x256();
      },
      [&]() {
        std::unique_ptr<BenchmarkFloat256x256ReaderBase> reader = std::make_unique<TReader>(kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadFloat256x256(a)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

template <typename TWriter, typename TReader>
void BenchmarkFloatVlen() {
  std::remove(kOutputFileName.c_str());

  yardl::NDArray<float, 2> a({256, 256});
  int i = 0;
  for (auto& x : a) {
    x = static_cast<float>(++i) - std::numeric_limits<float>::epsilon();
  }

  size_t const repetitions = ScaleRepetitions<TWriter>(10000);
  size_t const total_size = sizeof(float) * a.size() * repetitions;

  TimeScenario<TWriter>(
      __FUNCTION__,
      total_size,
      [&]() {
        std::unique_ptr<BenchmarkFloatVlenWriterBase> writer = std::make_unique<TWriter>(kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteFloatArray(a);
        }
        writer->EndFloatArray();
      },
      [&]() {
        std::unique_ptr<BenchmarkFloatVlenReaderBase> reader = std::make_unique<TReader>(kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadFloatArray(a)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

template <typename TWriter, typename TReader>
void BenchmarkSmallRecord() {
  std::remove(kOutputFileName.c_str());

  SmallBenchmarkRecord record{73278383.23123213, 78323.2820379, -2938923.29882};

  size_t const repetitions = 1000000;
  size_t const total_size = sizeof(record) * repetitions;

  TimeScenario<TWriter>(
      __FUNCTION__,
      total_size,
      [&]() {
        std::unique_ptr<BenchmarkSmallRecordWriterBase> writer = std::make_unique<TWriter>(kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteSmallRecord(record);
        }
        writer->EndSmallRecord();
      },
      [&]() {
        std::unique_ptr<BenchmarkSmallRecordReaderBase> reader = std::make_unique<TReader>(kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadSmallRecord(record)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

template <typename TWriter, typename TReader>
void BenchmarkSmallRecordBatched() {
  std::remove(kOutputFileName.c_str());

  SmallBenchmarkRecord const record{73278383.23123213, 78323.2820379, -2938923.29882};
  std::vector<SmallBenchmarkRecord> batch(8192, record);

  size_t const repetitions = ScaleRepetitions<TWriter>(50000);
  size_t const total_size = batch.size() * sizeof(record) * repetitions;

  TimeScenario<TWriter>(
      __FUNCTION__,
      total_size,
      [&]() {
        std::unique_ptr<BenchmarkSmallRecordWriterBase> writer = std::make_unique<TWriter>(kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteSmallRecord(batch);
        }
        writer->EndSmallRecord();
      },
      [&]() {
        std::unique_ptr<BenchmarkSmallRecordReaderBase> reader = std::make_unique<TReader>(kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadSmallRecord(batch)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

template <typename TWriter, typename TReader>
void SmallOptionalsBatched() {
  std::remove(kOutputFileName.c_str());

  SimpleEncodingCounters const record{26723, 92738, 7899};
  std::vector<SimpleEncodingCounters> batch(8192, record);

  size_t repetitions = ScaleRepetitions<TWriter>(10000);
  size_t total_size = batch.size() * sizeof(record) * repetitions;

  TimeScenario<TWriter>(
      __FUNCTION__,
      total_size,
      [&]() {
        std::unique_ptr<BenchmarkSmallRecordWithOptionalsWriterBase> writer = std::make_unique<TWriter>(kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteSmallRecord(batch);
        }
        writer->EndSmallRecord();
      },
      [&]() {
        std::unique_ptr<BenchmarkSmallRecordWithOptionalsReaderBase> reader = std::make_unique<TReader>(kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadSmallRecord(batch)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

template <typename TWriter, typename TReader>
void BenchmarkSimpleMrd() {
  std::remove(kOutputFileName.c_str());

  SimpleAcquisition acq;
  acq.data.resize({32, 256});
  acq.trajectory = yardl::NDArray<float, 2>({32, 2});
  std::variant<SimpleAcquisition, Image<float>> value = acq;

  size_t const repetitions = ScaleRepetitions<TWriter>(50000);
  size_t const total_size = (sizeof(value) + acq.data.size() * sizeof(std::complex<float>) + acq.trajectory.size() * sizeof(float)) * repetitions;

  TimeScenario<TWriter>(
      __FUNCTION__,
      total_size,
      [&]() {
        std::unique_ptr<BenchmarkSimpleMrdWriterBase> writer = std::make_unique<TWriter>(kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteData(value);
        }
        writer->EndData();
      },
      [&]() {
        std::unique_ptr<BenchmarkSimpleMrdReaderBase> reader = std::make_unique<TReader>(kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadData(value)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

}  // namespace

int main() {
  WriteBenchmarkTableHeader();

  BenchmarkFloat256x256<binary::BenchmarkFloat256x256Writer, binary::BenchmarkFloat256x256Reader>();
  BenchmarkFloat256x256<hdf5::BenchmarkFloat256x256Writer, hdf5::BenchmarkFloat256x256Reader>();
  BenchmarkFloat256x256<ndjson::BenchmarkFloat256x256Writer, ndjson::BenchmarkFloat256x256Reader>();
  WriteSeparatorRow();

  BenchmarkFloatVlen<binary::BenchmarkFloatVlenWriter, binary::BenchmarkFloatVlenReader>();
  BenchmarkFloatVlen<hdf5::BenchmarkFloatVlenWriter, hdf5::BenchmarkFloatVlenReader>();
  BenchmarkFloatVlen<ndjson::BenchmarkFloatVlenWriter, ndjson::BenchmarkFloatVlenReader>();
  WriteSeparatorRow();

  BenchmarkSmallRecord<binary::BenchmarkSmallRecordWriter, binary::BenchmarkSmallRecordReader>();
  BenchmarkSmallRecord<hdf5::BenchmarkSmallRecordWriter, hdf5::BenchmarkSmallRecordReader>();
  BenchmarkSmallRecord<ndjson::BenchmarkSmallRecordWriter, ndjson::BenchmarkSmallRecordReader>();
  WriteSeparatorRow();

  BenchmarkSmallRecordBatched<binary::BenchmarkSmallRecordWriter, binary::BenchmarkSmallRecordReader>();
  BenchmarkSmallRecordBatched<hdf5::BenchmarkSmallRecordWriter, hdf5::BenchmarkSmallRecordReader>();
  BenchmarkSmallRecordBatched<ndjson::BenchmarkSmallRecordWriter, ndjson::BenchmarkSmallRecordReader>();
  WriteSeparatorRow();

  SmallOptionalsBatched<binary::BenchmarkSmallRecordWithOptionalsWriter, binary::BenchmarkSmallRecordWithOptionalsReader>();
  SmallOptionalsBatched<hdf5::BenchmarkSmallRecordWithOptionalsWriter, hdf5::BenchmarkSmallRecordWithOptionalsReader>();
  SmallOptionalsBatched<ndjson::BenchmarkSmallRecordWithOptionalsWriter, ndjson::BenchmarkSmallRecordWithOptionalsReader>();
  WriteSeparatorRow();

  BenchmarkSimpleMrd<binary::BenchmarkSimpleMrdWriter, binary::BenchmarkSimpleMrdReader>();
  BenchmarkSimpleMrd<hdf5::BenchmarkSimpleMrdWriter, hdf5::BenchmarkSimpleMrdReader>();
  BenchmarkSimpleMrd<ndjson::BenchmarkSimpleMrdWriter, ndjson::BenchmarkSimpleMrdReader>();
}
