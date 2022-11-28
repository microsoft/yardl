# Copyright (c) Microsoft Corporation.
# Licensed under the MIT License.

param(
    [string]$VcpkgPath="C:\vcpkg"
)

$ErrorActionPreference = "Stop"

$location = Get-Location

try
{
    $scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

    rm -Recurse -Force $scriptDir\model -ErrorAction SilentlyContinue
    mkdir -Force $scriptDir\model | Out-Null

    # Create model and generate code
    cd $scriptDir\model
    yardl init smoketest
    yardl generate

    cd $scriptDir\cpp
    rm -Recurse -Force $scriptDir\cpp\build
    mkdir -Force $scriptDir\cpp\build | Out-Null
    cd $scriptDir\cpp\build

    # Configure
    $toolchainFile = Join-Path $VcpkgPath scripts\buildsystems\vcpkg.cmake
    $configure_args = "..", "-DCMAKE_TOOLCHAIN_FILE=$toolchainFile"
    cmake $configure_args

    # Build
    cmake --build .

    # Run the binary we just built
    .\Debug\smoketest.exe

    # Verify that it produced the expected file
    if (!(Test-Path -Path "smoketest.h5"))
    {
        throw "The expected output file was not found"
    }
}
finally
{
    cd $location
}
