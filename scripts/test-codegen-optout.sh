#!/usr/bin/env bash

set -euo pipefail

# Test with HDF5 and NDJSon disabled at generation time
pushd models/sandbox \
    && yardl generate \
        -c namespace=OnlyBinary \
        -c cpp.generateHDF5=false \
        -c cpp.generateNDJson=false \
        -c cpp.sourcesOutputDir=../../cpp/onlybinary/generated \
        -c python.disabled=true \
        -c matlab.disabled=true \
    && popd

pushd cpp/onlybinary \
    && rm -rf build \
    && mkdir build \
    && cd build \
    && cmake .. \
    && ninja \
    && ./test_only_binary \
    && popd

# Test with HDF5 and NDJSon disabled at build time
pushd models/sandbox \
    && yardl generate \
        -c namespace=OnlyBinary \
        -c cpp.sourcesOutputDir=../../cpp/onlybinary/generated \
        -c python.disabled=true \
        -c matlab.disabled=true \
    && popd

pushd cpp/onlybinary \
    && rm -rf build \
    && mkdir build \
    && cd build \
    && cmake -D OnlyBinary_GENERATED_USE_HDF5=Off -D OnlyBinary_GENERATED_USE_NDJSON=Off .. \
    && ninja \
    && ./test_only_binary \
    && popd
