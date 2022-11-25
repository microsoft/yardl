#include <chrono>

#include "generated/binary/protocols.h"

int main() {
  smoketest::binary::MyProtocolWriter w("smoketest.bin");

  w.WriteHeader({"123"});

  w.WriteSamples({std::chrono::system_clock::now(), {1, 2, 3}});
  w.WriteSamples({std::chrono::system_clock::now(), {4, 5, 6}});
  w.EndSamples();

  return EXIT_SUCCESS;
}
