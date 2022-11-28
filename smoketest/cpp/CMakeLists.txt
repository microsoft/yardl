cmake_minimum_required(VERSION 3.19)
project(smoketest)
set(CMAKE_CXX_STANDARD 17)

add_executable(smoketest smoketest.cpp)
target_link_libraries(smoketest smoketest_generated)

add_subdirectory(generated)
