// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#include <filesystem>
#include <iostream>

#include "generated/binary/protocols.h"
#include "generated/protocols.h"
#include "generated/types.h"

using namespace yardl;

int main() {
  std::cout << "\n================BINARY================\n\n";
  std::string filename("onlybinary.bin");
  std::remove(filename.c_str());

  only_binary::binary::HelloWorldWriter writer(filename);
  writer.WriteData({{{892.37889483, -9932.485937837}, {73.383672763878, -33.3394472537}},
                    {{3883.22890980, 373.4933837}, {56985.39384393, -33833.3330128474373}},
                    {{283.383672763878, -33.3394472537}, {3883.22890980, 373.4933837}}});
  writer.EndData();
  writer.Close();  // validates that protocol was completed.

  only_binary::binary::HelloWorldReader reader(filename);
  FixedNDArray<std::complex<double>, 2> a;
  while (reader.ReadData(a)) {
    std::cout << a << std::endl;
  }
  reader.Close();  // validates that protocol was completed.

  return 0;
}
