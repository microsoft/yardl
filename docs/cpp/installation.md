# Installation

<!--@include: ../parts/installation-core.md-->

## C++ Dependencies

In order to compile the C++ code that `yardl` generates, you will need to have a
C++17 (or more recent) compiler and the following dependencies installed:

1. HDF5 with the [C++ API](https://support.hdfgroup.org/HDF5/doc/cpplus_RM/),
   version 1.10.5 or later.
2. [xtensor](https://xtensor.readthedocs.io/en/latest/), version 0.21.10 or
   later.
3. Howard Hinnant's [date](https://howardhinnant.github.io/date/date.html)
   library, version 3.0.0 or later.
4. [JSON for Modern C++](https://github.com/nlohmann/json), version: 3.11.1 or
   later.

### Pixi

An easy way to get all dependencies and a complete C++ toolchain is to
create an isolated environment with [pixi](https://pixi.sh/). To set up a new
project from scratch:

```bash
pixi init my-project
cd my-project
pixi add cmake cxx-compiler hdf5 xtensor howardhinnant_date nlohmann_json
pixi shell
```

This installs a complete C++ toolchain alongside the libraries, fully isolated
from your system. The `cxx-compiler` meta-package pulls in the appropriate
compiler for your platform.

### Conda

Alternatively, if using the [Conda](https://docs.conda.io/en/latest/) package
manager, these dependencies can be installed with:

``` bash
conda install -c conda-forge cmake cxx-compiler hdf5 xtensor howardhinnant_date nlohmann-json
```

### vcpkg

If using [vcpkg](https://vcpkg.io/en/index.html), you can use a manifest file
that looks like the one
[here](https://github.com/microsoft/yardl/blob/main/smoketest/cpp/vcpkg.json).

### Homebrew

On macOS, you can use [Homebrew](https://brew.sh/) to install the dependencies:

```bash
brew install hdf5 xtensor howard-hinnant-date
```

## CMake

The `yardl generate` command emits a `CMakeLists.txt` that defines an object
library and the necessary `find_package()` and `target_link_libraries()` calls.
It has been tested to work on Linux with Clang and GCC, on macOS with
Clang, and on Windows with MSVC with vcpkg.
