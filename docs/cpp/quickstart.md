# Quick Start

::: info Note
Yardl is currently based on YAML. If you are new to YAML, you can get an
overview [here](https://learnxinyminutes.com/docs/yaml/).
:::

To get started, create a new empty directory and `cd` into it. Then run:

``` bash
yardl init playground
```

This creates the initial structure and files for our project:

```txt
$ tree .
.
└── model
    ├── model.yml
    └── _package.yml
```

The Yardl model package is in the `model` directory.

`_package.yml` is the package's manifest.

``` yaml
namespace: Playground

cpp:
  sourcesOutputDir: ../cpp/generated

python:
  outputDir: ../python
```

It specifies the package's namespace along with code generation settings. The
`cpp.sourcesOutputDir` property specifies where the generated C++ code should
go. If you are not interested in generating Python code, you can remove the
`python` property from the file:

``` yaml
namespace: Playground

cpp:
  sourcesOutputDir: ../cpp/generated

python: // [!code --]
  outputDir: ../python // [!code --]
```

All other `.yml` and `.yaml` files in the directory are assumed to be yardl
model files. The contents of `model.yml` look like this:

```yaml
# This is an example protocol, which is defined as a Header value
# followed by a stream of zero or more Sample values
MyProtocol: !protocol
  sequence:

    # A Header value
    header: Header

    # A stream of Samples
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

    # The time the sample was taken
    timestamp: datetime

    # A vector of integers
    data: int*
```

`!protocol`, `!stream` and `!record` are custom YAML tags, which describe the
type of the YAML node that follows.

`MyProtocol` is a protocol, which is a defined sequence of values that can be
written to or read from a file or binary stream (e.g. over a network
connection). This example protocol says that there will be one `Header` value
followed by an unknown number of `Sample`s. `Header` and `Sample` are records.

To generate C++ code for this model, `cd` into the `model` directory and run:

```bash
yardl generate
```

This will generate a number of files in the `sourcesOutputDir` directory:

```txt
$ tree -L 2 --dirsfirst
.
├── binary
│   ├── protocols.cc
│   └── protocols.h
├── hdf5
│   ├── protocols.cc
│   └── protocols.h
├── ndjson
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
records and enums. `protocols.h` declares abstract protocol readers and writers,
which are the base classes for implementations in `binary/protocols.h` and
`hdf5/protocols.h`, and `ndjson/protocols.h`. The `yardl/yardl.h` file defines
core datatypes like arrays and dates, and the header files in `yardl/detail/`
are included in generated files but are not intended to be included by consuming
code.

Ok, let's write some code! In the parent directory of the generated code, `cpp`,
create `playground.cc` that looks like this:

```cpp
#include <iostream>
#include <string>

#include "generated/binary/protocols.h"

int main() {
  std::string filename = "playground.bin";
  std::remove(filename.c_str());

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

  return 0;
}
```

Adjacent to that file, create a `CMakeLists.txt` that looks something like this:

```cmake
cmake_minimum_required(VERSION 3.19)
project(playground)

set(CMAKE_CXX_STANDARD 17)

add_executable(playground playground.cc)
target_link_libraries(playground playground_generated)

add_subdirectory(generated)
```

Now let's compile and run this code. Here are the steps on Linux:

```bash
mkdir build
cd build
cmake ..
cmake --build .
./playground
```

You can inspect the binary file our code produced with:

```bash
hexdump -C playground.bin
```

Note that the binary file contains a JSON representation of the protocol's
schema. This allows code that was not previously aware of this protocol to
deserialize the contents.

In addition to the compact binary format, we can write the protocol out to an
HDF5 file. This requires only a few modifications to our code:

```cpp
#include <iostream>
#include <string>

#include "generated/binary/protocols.h" // [!code --]
#include "generated/hdf5/protocols.h" // [!code ++]

int main() {
  std::string filename = "playground.bin"; // [!code --]
  std::string filename = "playground.h5"; // [!code ++]
  std::remove(filename.c_str());

  {
     playground::binary::MyProtocolWriter writer(filename); // [!code --]
     playground::hdf5::MyProtocolWriter writer(filename); // [!code ++]

    writer.WriteHeader({"123"});

    writer.WriteSamples({std::chrono::system_clock::now(), {1, 2, 3}});
    writer.WriteSamples({std::chrono::system_clock::now(), {4, 5, 6, 7}});
    writer.EndSamples();
  }

  playground::binary::MyProtocolReader reader(filename); // [!code --]
  playground::hdf5::MyProtocolReader reader(filename); // [!code ++]

  playground::Header header;
  reader.ReadHeader(header);

  std::cout << "Header.subject: " << header.subject << std::endl;

  playground::Sample sample;
  while (reader.ReadSamples(sample)) {
    std::cout << "Sample.data.size(): " << sample.data.size() << std::endl;
  }

  return 0;
}
```

You can inspect the HDF5 file with HDFView or by running

```bash
h5dump playground.h5
```
