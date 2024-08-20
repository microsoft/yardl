cmake_minimum_required(VERSION 3.19)
project(smoketest)

set(CMAKE_CXX_STANDARD 17)

if(WIN32)
  add_compile_options(/W3 /WX)
else()
  add_compile_options(-Wall -Wextra -pedantic -Werror)
endif()

add_executable(smoketest smoketest.cc)
target_link_libraries(smoketest smoketest_generated)

if (MSVC)
  target_compile_options(test_model_generated PRIVATE /bigobj)
endif()

add_subdirectory(generated)
