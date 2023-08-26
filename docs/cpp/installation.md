# Installation

<!--@include: ../parts/installation-core.md-->

## C++ Dependencies

In order to compile the C++ code that `yardl` generates, you will need to have a
C++17 (or more recent) compiler and the following dependencies installed:

1. HDF5 with the [C++ API](https://support.hdfgroup.org/HDF5/doc/cpplus_RM/).
2. [xtensor](https://xtensor.readthedocs.io/en/latest/)
3. Howard Hinnant's [date](https://howardhinnant.github.io/date/date.html)
   library.
4. [JSON for Modern C++](https://github.com/nlohmann/json).


### Conda

If using the [Conda](https://docs.conda.io/en/latest/) package manager, these
dependencies can be installed with:

``` bash
conda install -c conda-forge hdf5 xtensor howardhinnant_date nlohmann-json
```

Alternatively, you can create a new conda environment with all dependencies and
compilers using an environment.yml like [the one in this
repo](https://github.com/microsoft/yardl/blob/main/environment.yml).

```bash
wget https://raw.githubusercontent.com/microsoft/yardl/main/environment.yml
conda env create -n yardl -f environment.yml
conda activate yardl
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
