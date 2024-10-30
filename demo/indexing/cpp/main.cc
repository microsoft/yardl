#include <random>
#include <sstream>

#include "generated/binary/protocols.h"

#define ASSERT(cond, msg)                                  \
  if (!(cond)) {                                           \
    std::cerr << "Assertion failed: " << msg << std::endl; \
    exit(1);                                               \
  }

int main(void) {
  std::stringstream output;

  sketch::binary::MyProtocolWriter writer(output);
  writer.WriteHeader(sketch::Header{"John Doe"});

  size_t sample_count = 0;
  std::vector<sketch::Sample> samples(77);
  for (auto& sample : samples) {
    sample.id = sample_count++;
    sample.data = xt::arange<int32_t>(sample_count, sample_count + 1000);
  }

  writer.WriteSamples(samples);

  samples.resize(33);
  for (auto& sample : samples) {
    sample.id = sample_count++;
    writer.WriteSamples(sample);
  }

  samples.resize(55);
  for (auto& sample : samples) {
    sample.id = sample_count++;
  }
  writer.WriteSamples(samples);

  writer.EndSamples();
  writer.Close();

  auto serialized_without_index = output.str();

  // Try to load IndexedReader without an index. Should throw an exception.
  {
    bool caught_expected = false;
    std::stringstream input(serialized_without_index);
    try {
      sketch::binary::MyProtocolIndexedReader reader(input);
    } catch (std::exception const& ex) {
      caught_expected = true;
    }

    ASSERT(caught_expected, "Expected MyProtocolIndexedReader to throw exception!");
  }

  output = std::stringstream{};

  // Copy the protocol stream to a new stream with indexing
  {
    std::stringstream input(serialized_without_index);
    sketch::binary::MyProtocolReader reader(input);

    sketch::binary::MyProtocolIndexedWriter writer(output);

    reader.CopyTo(writer);
    reader.Close();
    writer.Close();
  }

  auto serialized_with_index = output.str();

  // Test reading stream element-by-element from index
  {
    std::stringstream input(serialized_with_index);
    sketch::binary::MyProtocolIndexedReader reader(input);

    std::random_device rd;
    std::mt19937 g(rd());

    std::vector<size_t> indices(sample_count);
    for (size_t i = 0; i < sample_count; i++) {
      indices[i] = i;
    }
    std::shuffle(indices.begin(), indices.end(), g);

    ASSERT(reader.CountSamples() == sample_count, "CountSamples() failed");

    sketch::Sample sample;
    for (size_t idx : indices) {
      ASSERT(reader.ReadSamples(sample, idx), "Failed to read sample");
      ASSERT(sample.id == idx, "Failed to read correct sample");
    }

    reader.Close();
  }

  // Test batch reading stream from index
  {
    std::stringstream input(serialized_with_index);
    sketch::binary::MyProtocolIndexedReader reader(input);

    std::vector<sketch::Sample> samples;
    samples.reserve(3);
    size_t idx = 0;
    while (reader.ReadSamples(samples, idx)) {
      // Do something with samples
      idx += samples.size();
    }
    ASSERT(idx == sample_count, "Batch read all samples failed");
    reader.Close();
  }

  // Test indexing with an empty stream
  {
    std::stringstream output;

    sketch::binary::MyProtocolIndexedWriter writer(output);
    writer.WriteHeader(sketch::Header{"John Doe"});

    writer.EndSamples();
    writer.Close();

    auto serialized = output.str();

    std::stringstream input(serialized);
    sketch::binary::MyProtocolIndexedReader reader(input);

    ASSERT(reader.CountSamples() == 0, "CountSamples() failed");

    sketch::Sample sample;
    size_t idx = 0;
    while (reader.ReadSamples(sample, idx)) {
      // Do something with samples
      idx += 1;
    }
    ASSERT(idx == 0, "Read empty samples failed");
    reader.Close();
  }

  return 0;
}
