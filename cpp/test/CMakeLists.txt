enable_testing()

find_package(fmt CONFIG REQUIRED)
find_package(GTest REQUIRED)

# This is only needed by `yardl.h` in the custom NDArray test (i.e. `just cpp-test-ndarray`)
include_directories(${CMAKE_CURRENT_SOURCE_DIR})

add_executable(tests
  binary/coded_stream_test.cc
  binary/header_test.cc
  computed_fields_test.cc
  definitions_test.cc
  hdf5/hdf5_test.cc
  ndjson/schema_test.cc
  protocol_state_test.cc
  roundtrip_test.cc
)

# This additional test must be compiled with C++20 to test compatibility with
# the C++17 version of the generated library
add_library(abi_test OBJECT
  abi_test.cc
)
set_property(TARGET abi_test PROPERTY CXX_STANDARD 20)

target_link_libraries(
  tests
  abi_test
  fmt
  gtest
  gtest_main
  pthread
  test_model_generated
  test_model_generated_mocks
)

add_executable(benchmark
  benchmark.cc
)

target_link_libraries(benchmark
  test_model_generated
)

add_subdirectory(generated)

gtest_discover_tests(tests)

add_executable(translator
  translator.cc
)

target_link_libraries(translator
  fmt
  test_model_generated
)
