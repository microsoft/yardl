#include "lib/binary/protocols.h"

using namespace evo_test;

int main(void) {
  ::binary::MyProtocolReader r(std::cin);

  Header h;
  r.ReadHeader(h);

  assert(h.subject == "Anonymous Human");
  assert(h.weight == 75);
  assert(h.meta["age"][0] == "42");

  long id;
  r.ReadId(id);
  assert(id == 123456789);

  Sample s;
  int i = 1;
  while (r.ReadSamples(s)) {
    for (auto& d : s.data) {
      assert(d == i++);
      (void)i;
      (void)d;
    }

    // std::cout << s.timestamp.time_since_epoch << std::endl;
  };

  std::optional<Footer> footer;
  r.ReadFooter(footer);

  r.Close();

  return 0;
}
