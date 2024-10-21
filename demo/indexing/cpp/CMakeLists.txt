cmake_minimum_required(VERSION 3.19)

project(sketch)

set(CMAKE_BUILD_TYPE Debug)

SET(CMAKE_CXX_FLAGS_DEBUG "-O0 -ggdb")

message(STATUS "Build type: ${CMAKE_BUILD_TYPE}")

if(NOT CMAKE_CXX_STANDARD)
    set(CMAKE_CXX_STANDARD 17)
endif()

set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR})
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY  ${CMAKE_BINARY_DIR})

add_compile_options(-Wall -Wextra -pedantic -Werror -Wno-psabi)

add_subdirectory(generated)
add_executable(sketch main.cc)
target_link_libraries(sketch sketch_generated)
