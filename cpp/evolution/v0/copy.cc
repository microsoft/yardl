#include "lib/binary/protocols.h"

using namespace evo_test;

int main(void) {
  ::binary::ProtocolWithChangesReader r(std::cin);
  ::binary::ProtocolWithChangesWriter w(std::cout, r.GetVersion());
  r.CopyTo(w, 3);
  r.Close();
  w.Close();
  return 0;
}
