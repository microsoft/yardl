#!/usr/bin/env bash

# Copyright (c) Microsoft Corporation.
# Licensed under the MIT License.

# This script tests the functionality of the yardl code generation tool with specific configurations,
# including disabling HDF5 and NDJSON features at generation and build time.

set -euo pipefail

repo_root=$(readlink -f "$(dirname "$0")/..")

# Test with HDF5 and NDJSon disabled at generation time
cd "$repo_root/models/sandbox"
yardl generate \
    -c namespace=OnlyBinary \
    -c cpp.generateHDF5=false \
    -c cpp.generateNDJson=false \
    -c cpp.sourcesOutputDir=../../cpp/onlybinary/generated \
    -c python.disabled=true \
    -c matlab.disabled=true

cd "$repo_root/cpp/onlybinary"
rm -rf ./build && mkdir build && cd build \
    && cmake -G Ninja .. && ninja && ./test_only_binary

# Test with HDF5 and NDJSon disabled at build time
cd "$repo_root/models/sandbox"
yardl generate \
    -c namespace=OnlyBinary \
    -c cpp.sourcesOutputDir=../../cpp/onlybinary/generated \
    -c python.disabled=true \
    -c matlab.disabled=true

cd "$repo_root/cpp/onlybinary"
rm -rf ./build && mkdir build && cd build \
    && cmake -G Ninja -D OnlyBinary_GENERATED_USE_HDF5=Off -D OnlyBinary_GENERATED_USE_NDJSON=Off .. \
    && ninja && ./test_only_binary
