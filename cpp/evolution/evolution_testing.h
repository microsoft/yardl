// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once
#include <stdexcept>

// We want the evolution test assertions to evaluate regardless of the CMake
// Build mode, so these tests use a very simple custom assert macro.
#define EVO_ASSERT(expr)                                      \
  do {                                                        \
    if (!(expr)) {                                            \
      throw std::runtime_error(__FILE__ ":" +                 \
                               std::to_string(__LINE__) +     \
                               ": Assertion failed: " #expr); \
    }                                                         \
  } while (0)

#define EVO_ASSERT_EQUALISH(a, b) EVO_ASSERT(std::abs(a - b) < 0.0001)

static std::string HelloWorld = "Hello, World!";
