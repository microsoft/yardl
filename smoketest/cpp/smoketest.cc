// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include "generated/hdf5/protocols.h"

int main() {
  smoketest::hdf5::MyProtocolWriter w("smoketest.h5");

  w.WriteHeader({"123"});

  w.WriteSamples({std::chrono::system_clock::now(), {1, 2, 3}});
  w.WriteSamples({std::chrono::system_clock::now(), {4, 5, 6}});
  w.EndSamples();

  return EXIT_SUCCESS;
}
