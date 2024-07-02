#include <sstream>

#include <gtest/gtest.h>

#include "generated/binary/protocols.h"

TEST(ABITest, OptionalDate) {
  using namespace test_model;

  std::ostringstream oss;

  binary::ProtocolWithOptionalDateWriter w(oss);
  RecordWithOptionalDate rec1;
  w.WriteRecord(rec1);
  w.Close();

  std::string data = oss.str();
  std::istringstream iss(data);

  binary::ProtocolWithOptionalDateReader r(iss);
  std::optional<RecordWithOptionalDate> rec2_opt;
  r.ReadRecord(rec2_opt);
  r.Close();

  EXPECT_TRUE(rec2_opt.has_value());
}
