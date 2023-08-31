# Packages

A Yardl package is a directory with a `_package.yml` manifest.

Here is a commented `_package.yml` file:

```yaml
# _package.yml

# The namespace of the package. Required.
namespace: MyNamespace

# settings for C++ code generation (optional)
cpp:
  # The directory where generated code will be written.
  # The directory will be created if it does not exist.
  sourcesOutputDir: ../path/relative/to/this/file

  # Whether to generate a CMakeLists.txt file in sourcesOutputDir
  # Default true
  generateCMakeLists: true

# settings for Python code generation (optional)
python:
  # The directory where the generated Python package will be written
  outputDir: ../path/relative/to/this/file
```

In the future, this file will be able to reference other packages and specify
options for other languages.
