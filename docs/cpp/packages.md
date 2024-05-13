# Packages

A Yardl package is a directory with a `_package.yml` manifest.

Here is a commented `_package.yml` file:

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

# Settings for Python code generation (optional)
python:
  # The directory where the generated Python package will be written
  outputDir: ../path/relative/to/this/file

# Settings for MATLAB code generation (optional)
matlab:
  # The directory where the generated MATLAB packages will be written
  outputDir: ../path/relative/to/this/file
```

In the future, this file will be able to reference previous versions of
your packages and specify options for other languages.
