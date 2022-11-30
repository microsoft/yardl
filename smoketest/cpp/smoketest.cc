// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <iostream>

#include "generated/hdf5/protocols.h"

int main() {
  std::string filename = "smoketest.h5";
  std::remove(filename.c_str());

  smoketest::hdf5::MyProtocolWriter w(filename);

  w.WriteHeader({"123"});

  w.WriteSamples({std::chrono::system_clock::now(), {1, 2, 3}});
  w.WriteSamples({std::chrono::system_clock::now(), {4, 5, 6}});
  w.EndSamples();

  return EXIT_SUCCESS;
}
