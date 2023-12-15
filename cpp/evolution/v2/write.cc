#include "lib/binary/protocols.h"

using namespace evo_test;

int main(void) {
  ::binary::MyProtocolWriter w(std::cout);

  Header h;
  h.subject = "Anonymous Human";
  h.weight = 75.0f;
  h.meta = {{"age", {"42"}}, {"weight", {"75.0f"}}};

  w.WriteHeader({h.subject, h.meta, h.weight});

  w.WriteId("123456789");

  Sample s;
  s.data = {1, 2, 3};
  s.timestamp = std::chrono::system_clock::now();
  w.WriteSamples(s);

  w.WriteSamples({
      {std::chrono::system_clock::now(), {4, 5, 6}},
      {std::chrono::system_clock::now(), {7, 8, 9}},
  });

  w.EndSamples();

  w.WriteFooter(std::nullopt);

  w.Close();

  return 0;
}
