cmake_minimum_required(VERSION 3.19)

project(sketch)

set(CMAKE_BUILD_TYPE Debug)

SET(CMAKE_CXX_FLAGS_DEBUG "-O0 -ggdb")

message(STATUS "Build type: ${CMAKE_BUILD_TYPE}")

if(NOT CMAKE_CXX_STANDARD)
    set(CMAKE_CXX_STANDARD 17)
endif()

# Enable ccache if not already enabled by symlink masquerading and if no other
# CMake compiler launchers are already defined
find_program(CCACHE_EXECUTABLE ccache)
mark_as_advanced(CCACHE_EXECUTABLE)
if(CCACHE_EXECUTABLE)
  message(STATUS "ccache found")
  foreach(LANG C CXX)
    if(NOT DEFINED CMAKE_${LANG}_COMPILER_LAUNCHER AND NOT CMAKE_${LANG}_COMPILER MATCHES ".*/ccache")
      message(STATUS "Enabling ccache for ${LANG}")
      set(CMAKE_${LANG}_COMPILER_LAUNCHER ${CCACHE_EXECUTABLE} CACHE STRING "")
    endif()
  endforeach()
else()
  message(STATUS "ccache not found")
endif()

set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR})
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY  ${CMAKE_BINARY_DIR})

add_compile_options(-Wall -Wextra -pedantic -Werror -Wno-psabi)

add_subdirectory(generated)
add_executable(sketch main.cc)
target_link_libraries(sketch sketch_generated)
