#include "lib/binary/protocols.h"

using namespace evo_test;

int main(void) {
  ::binary::MyProtocolWriter w(std::cout);

  Header h;
  h.subject = "Anonymous Human";
  h.weight = 75;
  h.meta = {{"age", {"42"}}, {"weight", {"75.0f"}}};
  h.added = "this string is only available in v1";

  w.WriteHeader({h.meta, h.subject, h.weight, h.added});

  w.WriteId(123456789);

  Sample s;
  s.data = {1, 2, 3};
  s.timestamp = std::chrono::system_clock::now();
  w.WriteSamples(s);

  w.WriteSamples({
      {{4, 5, 6}, std::chrono::system_clock::now()},
      {{7, 8, 9}, std::chrono::system_clock::now()},
  });

  w.EndSamples();

  w.WriteMaybe(42);

  w.WriteFooter(std::nullopt);

  w.Close();

  return 0;
}
