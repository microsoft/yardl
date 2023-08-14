// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

#pragma once
#include <memory>

#include "format.h"

namespace yardl::testing {

template <typename T>
std::unique_ptr<T> CreateWriter(Format format, std::string const& filename);

template <typename T>
std::unique_ptr<T> CreateReader(Format format, std::string const& filename);

}  // namespace yardl::testing
