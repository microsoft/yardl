enable_testing()

find_package(fmt CONFIG REQUIRED)
find_package(GTest REQUIRED)

add_executable(tests
  binary/coded_stream_test.cc
  binary/header_test.cc
  computed_fields_test.cc
  definitions_test.cc
  hdf5/hdf5_test.cc
  protocol_state_test.cc
  roundtrip_test.cc
)

target_link_libraries(
  tests
  fmt
  gmock
  gtest
  gtest_main
  pthread
  test_model_generated
)

add_executable(benchmark
  benchmark.cc
)

target_link_libraries(benchmark
  test_model_generated
)

add_subdirectory(generated)

gtest_discover_tests(tests)
