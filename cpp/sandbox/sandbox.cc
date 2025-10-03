// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <filesystem>
#include <iostream>

#include <xtensor/xio.hpp>

#include "generated/binary/protocols.h"
#include "generated/hdf5/protocols.h"
#include "generated/ndjson/protocols.h"
#include "generated/protocols.h"
#include "generated/types.h"

using namespace sandbox;
using namespace yardl;

template <typename T>
void Write(std::string filename) {
  T w(filename);
  HelloWorldWriterBase& writer = w;

  writer.WriteData({{{892.37889483, -9932.485937837}, {73.383672763878, -33.3394472537}},
                    {{3883.22890980, 373.4933837}, {56985.39384393, -33833.3330128474373}},
                    {{283.383672763878, -33.3394472537}, {3883.22890980, 373.4933837}}});
  writer.EndData();

  writer.Close();  // validates that protocol was completed.
}

template <typename T>
void Read(std::string filename) {
  T r(filename);
  HelloWorldReaderBase& reader = r;
  FixedNDArray<std::complex<double>, 2> a;
  while (reader.ReadData(a)) {
    std::cout << a << std::endl;
  }

  reader.Close();  // validates that protocol was completed.
}

void ExecCommand(std::string command) {
  int ret = std::system(command.c_str());
  if (ret != 0) {
    std::cerr << "Command failed: " << command << std::endl;
    exit(1);
  }
}

int main() {
  std::cout << "=================HDF5=================\n\n";
  std::string filename = "sandbox.h5";
  std::remove(filename.c_str());
  Write<sandbox::hdf5::HelloWorldWriter>(filename);
  Read<sandbox::hdf5::HelloWorldReader>(filename);

  std::cout << "\nh5dump output:\n\n";
  ExecCommand("h5dump " + filename);

  std::cout << "\n================BINARY================\n\n";
  filename = "sandbox.bin";
  std::remove(filename.c_str());
  Write<sandbox::binary::HelloWorldWriter>(filename);
  Read<sandbox::binary::HelloWorldReader>(filename);

  std::cout << "\nhexdump output:\n\n";
  ExecCommand("hexdump -C " + filename);

  std::cout << "\n================NDJSON================\n\n";
  filename = "sandbox.ndjson";
  std::remove(filename.c_str());
  Write<sandbox::ndjson::HelloWorldWriter>(filename);
  Read<sandbox::ndjson::HelloWorldReader>(filename);

  std::cout << "\noutput:\n\n";
  ExecCommand("cat " + filename);

  return 0;
}
