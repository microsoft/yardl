// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.#include "benchmark.h"

// Some very basic throughput benchmarks.

#include "benchmark.h"

#include <functional>
#include <optional>

#include "factories.h"
#include "format.h"
#include "generated/binary/protocols.h"
#include "generated/hdf5/protocols.h"
#include "generated/ndjson/protocols.h"

using namespace test_model;

using namespace yardl::testing;

namespace {

Result BenchmarkFloat256x256(Format format) {
  auto scale = [format]() -> double {
    switch (format) {
      case Format::kNDJson:
        return 0.01;
      default:
        return 1;
    }
  }();

  yardl::FixedNDArray<float, 256, 256> a;
  int i = 0;
  for (auto& x : a) {
    x = static_cast<float>(++i) - std::numeric_limits<float>::epsilon();
  }

  size_t const repetitions = ScaleRepetitions(10000, scale);
  size_t const total_size = sizeof(a) * repetitions;

  return TimeScenario(
      total_size,
      [&]() {
        auto writer = CreateWriter<BenchmarkFloat256x256WriterBase>(format, kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteFloat256x256(a);
        }
        writer->EndFloat256x256();
      },
      [&]() {
        auto reader = CreateReader<BenchmarkFloat256x256ReaderBase>(format, kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadFloat256x256(a)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

Result BenchmarkFloatVlen(Format format) {
  auto scale = [format]() -> double {
    switch (format) {
      case Format::kHdf5:
        return 0.5;
      case Format::kNDJson:
        return 0.01;
      default:
        return 1;
    }
  }();

  yardl::NDArray<float, 2> a({256, 256});
  int i = 0;
  for (auto& x : a) {
    x = static_cast<float>(++i) - std::numeric_limits<float>::epsilon();
  }

  size_t const repetitions = ScaleRepetitions(10000, scale);
  size_t const total_size = sizeof(float) * a.size() * repetitions;

  return TimeScenario(
      total_size,
      [&]() {
        auto writer = CreateWriter<BenchmarkFloatVlenWriterBase>(format, kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteFloatArray(a);
        }
        writer->EndFloatArray();
      },
      [&]() {
        auto reader = CreateReader<BenchmarkFloatVlenReaderBase>(format, kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadFloatArray(a)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

Result BenchmarkSmallInt256x256(Format format) {
  auto scale = [format]() -> double {
    switch (format) {
      case Format::kNDJson:
        return 0.01;
      default:
        return 1;
    }
  }();

  yardl::FixedNDArray<int, 256, 256> a;
  for (auto& x : a) {
    x = 37;
  }

  size_t const repetitions = ScaleRepetitions(6000, scale);
  size_t const total_size = sizeof(a) * repetitions;

  return TimeScenario(
      total_size,
      [&]() {
        auto writer = CreateWriter<BenchmarkInt256x256WriterBase>(format, kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteInt256x256(a);
        }
        writer->EndInt256x256();
      },
      [&]() {
        auto reader = CreateReader<BenchmarkInt256x256ReaderBase>(format, kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadInt256x256(a)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

Result BenchmarkSmallRecord(Format format) {
  auto scale = [format]() -> double {
    switch (format) {
      case Format::kHdf5:
        return 0.005;
      case Format::kNDJson:
        return 0.01;
      default:
        return 1;
    }
  }();

  SmallBenchmarkRecord record{73278383.23123213, 78323.2820379, -2938923.29882};

  size_t const repetitions = ScaleRepetitions(50000000, scale);
  static_assert(sizeof(record) == 16);
  size_t const total_size = sizeof(record) * repetitions;

  return TimeScenario(
      total_size,
      [&]() {
        auto writer = CreateWriter<BenchmarkSmallRecordWriterBase>(format, kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteSmallRecord(record);
        }
        writer->EndSmallRecord();
      },
      [&]() {
        auto reader = CreateReader<BenchmarkSmallRecordReaderBase>(format, kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadSmallRecord(record)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

Result BenchmarkSmallRecordBatched(Format format) {
  auto scale = [format]() -> double {
    switch (format) {
      case Format::kNDJson:
        return 0.005;
      default:
        return 1;
    }
  }();

  SmallBenchmarkRecord const record{73278383.23123213, 78323.2820379, -2938923.29882};
  std::vector<SmallBenchmarkRecord> batch(8192, record);

  size_t const repetitions = ScaleRepetitions(20000, scale);
  size_t const total_size = batch.size() * sizeof(record) * repetitions;

  return TimeScenario(
      total_size,
      [&]() {
        auto writer = CreateWriter<BenchmarkSmallRecordWriterBase>(format, kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteSmallRecord(batch);
        }
        writer->EndSmallRecord();
      },
      [&]() {
        auto reader = CreateReader<BenchmarkSmallRecordReaderBase>(format, kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadSmallRecord(batch)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

Result BenchmarkSmallOptionalsBatched(Format format) {
  auto scale = [format]() -> double {
    switch (format) {
      case Format::kNDJson:
        return 0.01;
      default:
        return 1;
    }
  }();

  SimpleEncodingCounters const record{26723, 92738, 7899};
  std::vector<SimpleEncodingCounters> batch(8192, record);

  size_t repetitions = ScaleRepetitions(5000, scale);
  size_t total_size = batch.size() * sizeof(record) * repetitions;

  return TimeScenario(
      total_size,
      [&]() {
        auto writer = CreateWriter<BenchmarkSmallRecordWithOptionalsWriterBase>(format, kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteSmallRecord(batch);
        }
        writer->EndSmallRecord();
      },
      [&]() {
        auto reader = CreateReader<BenchmarkSmallRecordWithOptionalsReaderBase>(format, kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadSmallRecord(batch)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

Result BenchmarkSimpleMrd(Format format) {
  auto scale = [format]() -> double {
    switch (format) {
      case Format::kHdf5:
        return 0.5;
      case Format::kNDJson:
        return 0.01;
      default:
        return 1;
    }
  }();

  SimpleAcquisition acq;
  acq.data.resize({32, 256});
  acq.trajectory = yardl::NDArray<float, 2>({32, 2});
  std::variant<SimpleAcquisition, Image<float>> value = acq;

  size_t const repetitions = ScaleRepetitions(30000, scale);
  size_t single_size = sizeof(value) + acq.data.size() * sizeof(std::complex<float>) + acq.trajectory.size() * sizeof(float);
  if (single_size != 66032) {
    throw std::runtime_error("Unexpected size: " + std::to_string(single_size));
  }
  size_t const total_size = single_size * repetitions;

  return TimeScenario(
      total_size,
      [&]() {
        auto writer = CreateWriter<BenchmarkSimpleMrdWriterBase>(format, kOutputFileName);
        for (size_t i = 0; i < repetitions; ++i) {
          writer->WriteData(value);
        }
        writer->EndData();
      },
      [&]() {
        auto reader = CreateReader<BenchmarkSimpleMrdReaderBase>(format, kOutputFileName);
        size_t read_repetitions = 0;
        while (reader->ReadData(value)) {
          read_repetitions++;
        }
        AssertRepetitionsSame(repetitions, read_repetitions);
      });
}

std::map<std::string, std::function<Result(Format)>>
    function_map = {
        {"float256x256", BenchmarkFloat256x256},
        {"smallint256x256", BenchmarkSmallInt256x256},
        {"floatvlen", BenchmarkFloatVlen},
        {"smallrecord", BenchmarkSmallRecord},
        {"smallrecordbatched", BenchmarkSmallRecordBatched},
        {"smalloptionalsbatched", BenchmarkSmallOptionalsBatched},
        {"simplemrd", BenchmarkSimpleMrd},
};

std::optional<std::function<Result(Format)>> GetScenarioFunction(std::string scenario) {
  std::transform(scenario.begin(), scenario.end(), scenario.begin(), [](unsigned char c) { return std::tolower(c); });

  auto func_pair = function_map.find(scenario);
  if (func_pair == function_map.end()) {
    return std::nullopt;
  }
  return func_pair->second;
}

}  // namespace

int main(int argc, char* argv[]) {
  if (argc != 3) {
    std::cerr << "Incorrect number of arguments. Usage: banchmark <scenario> <hdf5 | binary | ndjson>" << std::endl;
    return 1;
  }

  std::string scenario = argv[1];
  Format format = ParseFormat(std::string(argv[2]));
  auto func = GetScenarioFunction(scenario);
  if (!func) {
    std::cerr << "Unknown scenario: " << scenario << std::endl;
    return 0;
  }

  Result result = (*func)(format);

  nlohmann::ordered_json j = result;
  std::cout << j << std::endl;
}
