#!/usr/bin/env bash

# Copyright (c) Microsoft Corporation.
# Licensed under the MIT License.

# This script generates the NOTICE file for the repository based on the GO dependencies
# using the go-licenses tool.

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
mkdir -p build1 && cd build1 \
    && cmake -G Ninja .. && ninja && ./test_only_binary

# Test with HDF5 and NDJSon disabled at build time
cd "$repo_root/models/sandbox"
yardl generate \
    -c namespace=OnlyBinary \
    -c cpp.sourcesOutputDir=../../cpp/onlybinary/generated \
    -c python.disabled=true \
    -c matlab.disabled=true

cd "$repo_root/cpp/onlybinary"
mkdir -p build2 && cd build2 \
    && cmake -G Ninja -D OnlyBinary_GENERATED_USE_HDF5=Off -D OnlyBinary_GENERATED_USE_NDJSON=Off .. \
    && ninja && ./test_only_binary
