# Yardl Documentation <!-- omit from toc -->

- [Installation](#installation)
- [C++ Dependencies](#c-dependencies)
- [Quick Start](#quick-start)
- [Packages](#packages)
- [Yardl Syntax](#yardl-syntax)
  - [Protocols](#protocols)
  - [Primitive Types](#primitive-types)
  - [Optional Types](#optional-types)
  - [Unions](#unions)
  - [Records](#records)
  - [Enums](#enums)
  - [Vectors](#vectors)
  - [Arrays](#arrays)
  - [Named Types](#named-types)
  - [Computed Fields](#computed-fields)
- [Command-Line Reference](#command-line-reference)
- [Performance Tips](#performance-tips)

## Installation

Yardl is a single executable file. The installation steps are:

1. Head over to the [latest
   release](https://github.com/microsoft/yardl/releases/latest) page.
2. Download the right archive for your platform.
3. Extract the archive and find the `yardl` executable. Copy it to a directory
   in your `PATH` environment variable.

You should now be able to run `yardl --version`.

## C++ Dependencies

In order to compile the C++ code that `yardl` generates, you will need to have a
C++17 (or more recent) compiler and the following dependencies installed:

1. HDF5 with the [C++ API](https://support.hdfgroup.org/HDF5/doc/cpplus_RM/).
2. [xtensor](https://xtensor.readthedocs.io/en/latest/)
3. If using C++17, Howard Hinnant's
   [date](https://howardhinnant.github.io/date/date.html) library.

If using the [Conda](https://docs.conda.io/en/latest/) package manager, these
can be installed with:

``` bash
conda install -c conda-forge hdf5 xtensor howardhinnant_date
```

If using [vcpkg](https://vcpkg.io/en/index.html), you can use a manifest file
that looks like the one
[here](../smoketest/cpp/vcpkg.json).

The `yardl generate` command emits a `CMakeLists.txt` that defines an
object library and the necessary `find_package()` and `target_link_libraries()`
calls. It has been tested to work with Conda on Linux with Clang and GCC and on
Windows with MSVC with vcpkg. MacOS and homebrew support is coming.

## Quick Start

> **Note**
> Yardl is currently based on YAML. If you are not familiar with its syntax, you
> can get an overview [here](https://learnxinyminutes.com/docs/yaml/).


To get started, create a new empty directory and `cd` into it. This directory
will contain ouy yardl package. To quickly create a package you can run:

``` bash
yardl init playground
```

This creates a package with the the namespace `Playground`, containing the following files:

```text
$ tree .
.
├── model.yml
└── _package.yml
```

_package.yml is the package's manifest.

``` yaml
namespace: Playground

cpp:
  sourcesOutputDir: ../cpp/generated
```

It specifies the package's namespace along with code generation settings. The
`cpp.sourcesOutputDir` property specifies where the generated C++ code should go.

All other `.yml` files in the directory are assumed to be yardl model files. The
contents of `model.yml` look like this:

```yaml
# This is an example protocol, which is defined as a Header value
# followed by a stream of zero or more Sample values
MyProtocol: !protocol
  sequence:
    header: Header
    samples: !stream
      items: Sample

# Header is a record with a single string field
Header: !record
  fields:
    subject: string

# Sample is a record made up of a datetime and
# a vector of integers
Sample: !record
  fields:
    timestamp: datetime
    data: !vector
      items: int
```

`!protocol`, `!stream` and `!record` are custom YAML tags, which describe the
type of the YAML node that follows.

`MyProtocol` is a protocol, which is a defined sequence of values that are to be
transmitted or received. This protocol says that there will be one `Header`
value followed by an unknown number of `Sample`s. `Header` and `Sample` are
records, which are converted to C++ structs.

To generate code for this model, run this from the same directory:

```bash
yardl generate
```

This will generate a number of files in the `sourcesOutputDir` directory:

``` text
$ tree -L 2 --dirsfirst
.
├── binary
│   ├── protocols.cc
│   └── protocols.h
├── hdf5
│   ├── protocols.cc
│   └── protocols.h
├── yardl
│   ├── detail
│   └── yardl.h
├── CMakeLists.txt
├── protocols.cc
├── protocols.h
└── types.h
```

In the root directory, `types.h` contains generated code for named types like
records and enums. `protocols.h` declares abstract protocol readers and writer,
which are the base classes for implementations in `binary/protocols.h` and
`hdf5/protocols.h`. `yardl.yardl.h` defines core datatypes like arrays and
dates, and the header files in `yardl/detail` are included in generated files
but are not intended to be included by consuming code.

Ok, let's write some code! In the parent directory of `cpp/generated/`, `cpp/`
create `playground.cc` that looks like this:

```cpp
#include <iostream>
#include <string>

#include "generated/binary/protocols.h"

int main() {
  std::string filename = "playground.bin";

  {
    playground::binary::MyProtocolWriter writer(filename);

    writer.WriteHeader({"123"});

    writer.WriteSamples({std::chrono::system_clock::now(), {1, 2, 3}});
    writer.WriteSamples({std::chrono::system_clock::now(), {4, 5, 6, 7}});

    // signal the end of the samples stream
    writer.EndSamples();
  }

  playground::binary::MyProtocolReader reader(filename);

  playground::Header header;
  reader.ReadHeader(header);

  std::cout << "Read Header.subject: " << header.subject << std::endl;

  playground::Sample sample;
  while (reader.ReadSamples(sample)) {
    std::cout << "Read Sample.data.size(): " << sample.data.size() << std::endl;
  }

  return EXIT_SUCCESS;
}
```

Adjacent to that file, create a `CMakeLists.txt` that looks something like this (on Linux):

```cmake
cmake_minimum_required(VERSION 3.19)
project(playground)

set(CMAKE_CXX_STANDARD 17)

add_executable(playground playground.cc)
target_link_libraries(playground playground_generated)

add_subdirectory(generated)
```

Now let's compile and run this code.

```bash
mkdir build
cd build
cmake .. -GNinja
ninja
./playground
```

You can inspect the binary file with:

```bash
hexdump -C playground.bin
```

Note that the binary file contains a JSON representation of the protocol's
schema. This allows code that was not previously aware of this protocol to
deserialize this file.

In addition to the compact binary format, we can write the protocol out to an HDF5 file.
This requires only a few modifications to our code:

```diff
  #include <iostream>
  #include <string>

- #include "generated/binary/protocols.h"
+ #include "generated/hdf5/protocols.h"

  int main() {
-   std::string filename = "playground.bin";
+   std::string filename = "playground.h5";

    {
-      playground::binary::MyProtocolWriter writer(filename);
+      playground::hdf5::MyProtocolWriter writer(filename);

      writer.WriteHeader({"123"});

      writer.WriteSamples({std::chrono::system_clock::now(), {1, 2, 3}});
      writer.WriteSamples({std::chrono::system_clock::now(), {4, 5, 6, 7}});
      writer.EndSamples();
    }

-   playground::binary::MyProtocolReader reader(filename);
+   playground::hdf5::MyProtocolReader reader(filename);

    playground::Header header;
    reader.ReadHeader(header);

    std::cout << "Header.subject: " << header.subject << std::endl;

    playground::Sample sample;
    while (reader.ReadSamples(sample)) {
      std::cout << "Sample.data.size(): " << sample.data.size() << std::endl;
    }

    return EXIT_SUCCESS;
  }

```

## Packages

## Yardl Syntax

### Protocols

### Primitive Types

### Optional Types

### Unions

### Records

### Enums

### Vectors

### Arrays

### Named Types

### Computed Fields

## Command-Line Reference

## Performance Tips
