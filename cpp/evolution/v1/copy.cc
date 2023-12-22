#include "lib/binary/protocols.h"

using namespace evo_test;

int main(void) {
  ::binary::MyProtocolReader r(std::cin);
  ::binary::MyProtocolWriter w(std::cout, r.GetVersion());
  r.CopyTo(w);
  r.Close();
  w.Close();
  return 0;
}
