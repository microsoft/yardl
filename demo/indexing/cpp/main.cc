#include <random>
#include <sstream>

#include "generated/binary/protocols.h"


int main(void)
{
  std::stringstream output;

  sketch::binary::MyProtocolWriter writer(output);
  writer.WriteHeader(sketch::Header{"John Doe"});

  size_t sample_count = 0;
  std::vector<sketch::Sample> samples(77);
  for (auto& sample : samples) {
    sample.id = sample_count++;
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

  {
    bool caught_expected = false;
    std::stringstream input(serialized_without_index);
    try {
      sketch::binary::MyProtocolIndexedReader reader(input);
    } catch (const std::exception& ex) {
      caught_expected = true;
    }

    if (!caught_expected) {
      std::cerr << "Fail: Expected MyProtocolIndexedReader to throw exception!" << std::endl;
    }
  }

  output = std::stringstream{};

  {
    std::stringstream input(serialized_without_index);
    sketch::binary::MyProtocolReader reader(input);

    sketch::binary::MyProtocolIndexedWriter writer(output);

    reader.CopyTo(writer);
    reader.Close();
    writer.Close();
  }

  auto serialized_with_index = output.str();

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

    sketch::Sample sample;
    for (size_t idx : indices) {
      if (!reader.ReadSamples(sample, idx)) {
        std::cerr << "No more samples to read " << idx << std::endl;
        break;
      }
      if (sample.id != idx) {
        std::cerr << "Failed to read sample " << idx << std::endl;
        break;
      }
    }

    reader.Close();
  }

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
    std::cerr << "Batch read all samples: " << (idx == sample_count ? "SUCCESS" : "FAILURE") << std::endl;
    reader.Close();
  }

  return 0;
}
