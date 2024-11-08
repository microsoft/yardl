# Packages

A Yardl package is a directory with a `_package.yml` manifest.

## Package Manifest

Here is an example commented `_package.yml` file:

```yaml
# _package.yml

# The namespace of the package. Required.
namespace: MyNamespace

# Import model types from external locations (URLs)
#
# URLs may be one of the following:
#   1. Path to local directory, either
#     - Relative to the `_package.yml` manifest, or
#     - Absolute path
#   2. Remote git repository. Options (provided as query parameters):
#     - ref=<git commit hash>
#     - dir=<relative path to model>
# NOTE: Yardl caches git repositories in `$HOME/.yardl/cache`
imports:
  - ../myCommonTypes
  - /workspaces/yardl/models/more-common-types
  - https://github.com/microsoft/yardl?ref=31a6e29&dir=models/test

# Evolve your schema from previous versions
# See imports above for details on specifying model version locations
versions:
  v0_1: ../models/test/v0.1
  v20240201: https://github.com/microsoft/yardl/models/test/v20240201

# Settings for C++ code generation (optional)
cpp:
  # The directory where generated code will be written.
  # The directory will be created if it does not exist.
  sourcesOutputDir: ../path/relative/to/this/file

  # Whether to generate a CMakeLists.txt file in sourcesOutputDir
  # Default true
  generateCMakeLists: true

  # Include path for custom NDArray implementation header file
  # If provided, the generated C++ code will include this header
  # instead of Yardl's default NDArray implementation
  overrideArrayHeader: path/to/custom/NDArray/header

# Settings for Python code generation (optional)
python:
  # The directory where the generated Python package will be written
  outputDir: ../path/relative/to/this/file

# Settings for MATLAB code generation (optional)
matlab:
  # The directory where the generated MATLAB packages will be written
  outputDir: ../path/relative/to/this/file
```

## Overriding the Package Manifest

Fields in the `_package.yml` manifest can be overriden on the command-line using the `-c/--config` flag, e.g.

```bash
$ yardl generate -c cpp.sourcesOutputDir=/tmp/$(date +%F_%T)/generated
✅ Wrote C++ to /tmp/2024-09-30_16:39:43/generated.
✅ Wrote Python to /workspaces/yardl/python.
✅ Wrote Matlab to /workspaces/yardl/matlab/generated.
```
