# Yardl Documentation <!-- omit from toc -->

- [Installation](#installation)
- [C++ Dependencies](#c-dependencies)
- [Quick Start](#quick-start)
- [Packages](#packages)
- [Yardl Syntax](#yardl-syntax)
  - [Protocols](#protocols)
  - [Records](#records)
  - [Primitive Types](#primitive-types)
  - [Unions](#unions)
  - [Enums](#enums)
  - [Vectors](#vectors)
  - [Arrays](#arrays)
  - [Type Aliases](#type-aliases)
  - [Computed Fields](#computed-fields)
  - [Generics](#generics)
- [Performance Tips](#performance-tips)
  - [Batched Reads and Writes](#batched-reads-and-writes)
  - [Use Fixed Data Types When Possible](#use-fixed-data-types-when-possible)
- [Protocol Schema Reference](#protocol-schema-reference)
  - [References to Primitive Types](#references-to-primitive-types)
  - [References to Top-Level Types](#references-to-top-level-types)
  - [Unions](#unions-1)
  - [Vectors](#vectors-1)
  - [Arrays](#arrays-1)
  - [Streams](#streams)
  - [Enums](#enums-1)
  - [Records](#records-1)
  - [Aliases](#aliases)
  - [Protocols](#protocols-1)
  - [Top-Level Schema](#top-level-schema)

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
[here](../smoketest/cpp/vcpkg.JSON).

The `yardl generate` command emits a `CMakeLists.txt` that defines an
object library and the necessary `find_package()` and `target_link_libraries()`
calls. It has been tested to work with Conda on Linux with Clang and GCC and on
Windows with MSVC with vcpkg. MacOS and homebrew support is coming.

## Quick Start

> **Note**<br>
> Yardl is currently based on YAML. If you are new to YAML, you can get an
> overview [here](https://learnxinyminutes.com/docs/yaml/).

To get started, create a new empty directory and `cd` into it. This directory
will contain the yardl package. To quickly create a package you can run:

``` bash
yardl init playground
```

This creates a package with the the namespace `Playground`, containing the
following files:

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

All other `.yml` and `.yaml` files in the directory are assumed to be yardl
model files. The contents of `model.yml` look like this:

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

`MyProtocol` is a protocol, which is a defined sequence of values that can be
written to or read from a file or binary stream (e.g. over a network
connection). This example protocol says that there will be one `Header` value
followed by an unknown number of `Sample`s. `Header` and `Sample` are records.

To generate C++ code for this model, run this from the same directory:

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
records and enums. `protocols.h` declares abstract protocol readers and writers,
which are the base classes for implementations in `binary/protocols.h` and
`hdf5/protocols.h`. The `yardl/yardl.h` file defines core datatypes like arrays and
dates, and the header files in `yardl/detail/` are included in generated files
but are not intended to be included by consuming code.

Ok, let's write some code! In the parent directory of the generate code, `cpp`,
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
cmake .. -GNinja
ninja
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

    return 0;
  }

```

You can inspect the HDF5 file with HDFView or by running

```bash
h5dump playground.h5
```

## Packages

A package is a directory with a `_package.yml` manifest.

Here is a commented `_package.yml` file:

```yaml
# _package.yml

# The namespace of the package. Required.
namespace: MyNamespace

# settings for C++ code generation
cpp:
  # The directory where generated code will be written.
  # The directory will be created if it does not exist.
  sourcesOutputDir: ../path/relative/to/this/file

  # Whether to generate a CMakeLists.txt file in sourcesOutputDir
  # Default true
  generateCMakeLists: true
```

In the future, this file will be able to reference other packages and specify
options for other languages.

## Yardl Syntax

Yardl model files use YAML syntax and are requires to have either a `.yml` or
`.yaml` file extension.

To efficiently work with yardl, we recommend that you run the following from the
package directory:

```bash
yardl generate --watch
```

This watches the directory for changes and generates code whenever a file is
saved. This allows you to get rapid feedback as you experiment.

Comments placed above top-level types and their fields are captured and added to
the generated code.

`yardl generate` only generates code once the model files in the package have
been validated. It will write out any validation errors to standard error.

### Protocols

As explained in the [quick start](#quick-start), protocols define a sequence of
values, called "steps", that are required to be transmitted, in order. They are
defined like this:

```yaml
MyProtocol: !protocol
sequence:
  a: int
  b: !stream
    items: float
  c: !stream:
    items string
```

In the example, the first step is a single integer named `a`. Following that
will be a stream (named `b`) of zero or more floating-point numbers, and a
stream (named `c`) of strings.

It is an error to attempt to read or write a protocol's steps out of order. In
order to verify that a protocol has been completely written to or read from, you
can call `Close()` on the generated reader or writer instance.

Generated protocol readers have a `CopyTo()` method that allows you to copy the
contents of the protocol to another protocol writer. This makes is easy to, say,
read from an HDF5 file and send the data in the binary format over a network
connection.

### Records

Records have fields and, optionally, [computed fields](#computed-fields). In
generated C++ code, they map to structs.

Fields have a name and can be of any primitive or compound type. For example:

```yaml
MyRecord: !record
  fields:
    myIntField: int
    myStringField: string
```

Records must be declared at the top level and cannot be inlined. For example,
this is not supported:

```yaml
RecordA: !record
  fields:
    recA: !record # NOT SUPPORTED!
      fields:
        a: int
    recB: RecordB # But this is fine.

RecordB: !record
  fields:
    c: int
```

Note that Yardl does not support type inheritance.

### Primitive Types

Yardl has the following primitive types:

| Type             | Comment                                                                 |
| ---------------- | ----------------------------------------------------------------------- |
| `bool`           |                                                                         |
| `int8`           |                                                                         |
| `uint8`          |                                                                         |
| `byte`           | Alias of `uint8`                                                        |
| `int16`          |                                                                         |
| `uint16`         |                                                                         |
| `int32`          |                                                                         |
| `int`            | Alias of `int32`                                                        |
| `uint32`         |                                                                         |
| `uint`           | Alias of `unit32`                                                       |
| `int64`          |                                                                         |
| `long`           | Alias of `int64`                                                        |
| `uint64`         |                                                                         |
| `ulong`          | Alias of `uint64`                                                       |
| `size`           | Equivalent to `uint64`                                                  |
| `float32`        |                                                                         |
| `float`          | Alias of `float32`                                                      |
| `float64`        |                                                                         |
| `double`         | Alias of `float64`                                                      |
| `complexfloat32` | A complex number where each component is a 32-bit floating-point number |
| `complexfloat`   | Alias of `complexfloat32`                                               |
| `complexfloat64` | A complex number where each component is a 63-bit floating-point number |
| `complexdouble`  | Alias of `complexfloat64`                                               |
| `string`         |                                                                         |
| `date`           | A number of days since the epoch                                        |
| `time`           | A number of nanoseconds after midnight                                  |
| `datetime`       | A number of nanoseconds since the epoch                                 |

### Unions

When a value could be one of several types, you can define a union:

```yaml
Rec: !record
  fields:
    intOrFloat: [int, float]
    intOrFloatExpandedForm:
      - int
      - float
    nullableIntOrFloat:
      - null
      - int
      - float
    arrayOfFloatsOrDoubles:
      - !array
        items: float
      - !array
        items: double
```

The `null` type in the example above means that no value is also a possibility.

The `?` suffix can be appended to a type name as a shorthand to define an
*optional type*, a special case of union. For example, `int?` is the same as
`[null, int]`. Note that the expanded form has to be used for complex optional
types:

```yaml
Rec: !record
  fields:
    optionalArray:
      - null
      - !vector
        items: int
```

### Enums

Enums can be defined as a list of values:

```yaml
Fruits: !enum
  values:
    - apple
    - banana
    - pear
```

You can optionally specify the underlying type of the enum and give each symbol
an integer value:

```yaml
UInt64Enum: !enum
  base: uint64
  values:
    a: 0x1
    b: 0x2
    c: 20
```

### Vectors

Vectors are one-dimensional arrays. They can optionally have a fixed length.

```yaml
MyRec: !record
  fields:
    vec1: !vector
      items: int
    vec2: !vector
      items: int
      length: 10
```

In generated C++ code, `vec1` maps to an `std::vector<int>` and `vec2` to an
`std::array<int, 10>`

### Arrays

The `!array` tag denotes a multidimensional array. They can be of a fixed size:

```yaml
MyRec: !record
  fields:
    fixedNdArray: !array
      items: float
      dimensions: [3, 4]
```

Or the size might not be fixed but the number of dimensions is known:

```yaml
MyRec: !record
  fields:
    ndArray: !array
      items: float
      dimensions: 2
```

Or finally, the number of dimensions may be unknown as well:

```yaml
MyRec: !record
  fields:
    dynamicNdArray: !array
      items: float
```

Dimensions can be given names, which can be used in [computed
field](#computed-fields) expressions.

```yaml
MyRec: !record
  fields:
    fixedNdArray: !array
      items: float
      dimensions:
        x: 3
        y: 4
    ndArray: !array
      items: float
      dimensions: [x, y]
    ndArrayAlternate: !array
      items: float
      dimensions:
        x:
        y:
```

### Type Aliases

We've seen records, enums, and protocols defined as top-level, named types, but
any type can be given one or more aliases:

```yaml
FloatArray: !array
  items: float

SignedInteger: [int8, int16, int32, int64]

Id: string
Name: string
```

This simply gives another name to a type, so the `Name` type above is no
different from the `string` type.

### Computed Fields

In addition to fields, records can contain computed fields. These are simple expressions
over the record's other (computed) fields.

```yaml
MyRec: !record
  fields:
    arrayField: !array
        items: int
        dimensions: [x, y]
  computedFields:
    accessArray: arrayField
    accessArrayElement: arrayField[0, 1]
    accessArrayElementByName: arrayField[y:1, x:0]
    sizeOfArrayField: size(arrayField)
    sizeOfFirstDimension: size(arrayField, 0)
    sizeOfXDimension: size(arrayField, 'x')
```

To work with union types, you need to use a switch expression with type pattern
matching:

```yaml
NamedArray: !array
  items: int
  dimensions: [x, y]

MyRec: !record
  fields:
    myUnion: [null, int, NamedArray]
  computedFields:
    myUnionSize:
      !switch myUnion:
        int: 1 # if the union holds an int
        NamedArray arr: size(arr) # if it's a NamedArray. Note the variable declaration.
        _: 0 # all other cases (here it's just null)
```

The following function calls are supported from computed field expressions:

- `size(vector)`: returns the size (length) of the vector
- `size(array)`: returns the total size of the array
- `size(array, integer)`: returns the size of the array's dimension at the given
  index
- `size(array, string)`: returns the size of the array's dimension with the
  given name

- `dimensionIndex(array, string)` returns the index of the dimension with the
  given name

- `dimensionCount(array)` returns the dimension count of the array

### Generics

Yardl supports generic types.

```yaml
Image<T>: !array
  items: T

ImageVariant:
  - Image<float>
  - Image<double>
  - Image<complexfloat>
  - Image<complexdouble>

RecordWithImages<T, U>: !record
  fields:
    image1: Image<T>
    image2: Image<U>
```

Note that protocols cannot be generic types, but its steps may be made up of
closed generic types (e.g. `Image<float>`).

## Performance Tips

### Batched Reads and Writes

Generated protocol reader and writer classes have read and write methods for
each step. When a step is a stream, there will also be overloads that read or
write a batch of values as an `std::vector` in one go. For small data types,
using batched reads and writes can make a dramatic difference in throughput,
especially for HDF5 files.

### Use Fixed Data Types When Possible

For HDF5, using variable-length collections (like `!vector` without a length or
`!array` without fixed dimension sizes) has lower throughput than their fixed-sized
counterparts.

## Protocol Schema Reference

A protocol's schema is embedded in a JSON format in both the HDF5 and binary
formats. This JSON format is informally described here.

> **Warning**<br>
> We might make breaking changes to this format before V1.

The JSON schema is meant to be compact and therefore it does not contain of the
comments that may be in the Yardl, nor does it contain computed fields, since
those are not needed for deserialization. We chose JSON because it can be easily
parsed at runtime during deserialization but also easily read.

### References to Primitive Types

Primitive type references are represented by name as a JSON string, e.g.
`"int32"` or `"string"`

### References to Top-Level Types

References to top-level types (records, enums, and aliases) are represented by
their namespaced name as a JSON string, e.g. `"MyNamespace.MyRecord"` or
`"MyNamespace.MyEnum"`.

### Unions

Unions are represented as a JSON array:

```JSON
[
  {
    "label": "int32",
    "type": "int32"
  },
  {
    "label": "float32",
    "type": "float32"
  }
]
```

The `label` field a unique name given to each union case, derived from its type
name. The labels are used in HDF5 the HDF5 format.

If `null` is one of the cases, then it is represented by `null` in the JSON as well:

```JSON
[
  null,
  {
    "label": "int32",
    "type": "int32"
  },
  {
    "label": "float32",
    "type": "float32"
  }
]
```

For the special case of an optional type, a label for the non-null case is
omitted and the object is simplified to its `type` value:

```JSON
[
  null,
  "int32"
]
```

### Vectors

Vectors have the following representation:

```JSON
{
  "vector": {
    "items": "int32",
    "length": 10
  }
}
```

The `length` field is only present if it is specified in the Yardl definition.

### Arrays

A fixed array with dimensions `x` and `y` would look like this:

```JSON
{
  "array": {
    "items": "int32",
    "dimensions": [
      {
        "name": "x",
        "length": 3
      },
      {
        "name": "y",
        "length": 4
      }
    ]
  }
}
```

A non-fixed array with dimensions `x` and `y` would look like the above but
without the `length` field:

```JSON
{
  "array": {
    "items": "int32",
    "dimensions": [
      {
        "name": "x"
      },
      {
        "name": "y"
      }
    ]
  }
}
```

A non-fixed array with two unnamed dimensions would look like this:

```JSON
{
  "array": {
    "items": "int32",
    "dimensions": 2
  }
}
```

And finally, an array with an unknown number of dimensions:

```JSON
{
  "array": {
    "items": "int32"
  }
}
```

### Streams

Streams in protocols are represented as:

```JSON
{
  "stream": {
    "items": "int32",
  }
}
```

### Enums

Enums are top-level types and cannot be declared inline.

An example enum that looks like this in Yardl:

```yaml
Animals: !enum
  base: uint8
  values: [cat, dog]
```

Is represented in JSON as:

```JSON
{
  "enum": {
    "name": "Animals",
    "base": "uint8",
    "values": [
      {
        "symbol": "cat",
        "value": 0
      },
      {
        "symbol": "dog",
        "value": 1
      }
    ]
  }
}
```

The `base` field is only present in the JSON if it is specified in the Yardl.

### Records

Records are top-level types that cannot be declared inline.

An example generic record:

```yaml
MyTuple<T1, T2>: !record
  fields:
    f1: T1
    f2: T2
```

Would look like this:

```JSON
{
  "record": {
    "name": "MyTuple",
    "typeParameters": [
      "T1",
      "T2"
    ],
    "fields": [
      {
        "name": "f1",
        "type": "T1"
      },
      {
        "name": "f2",
        "type": "T2"
      }
    ]
  }
}
```

Computed fields are omitted from the protocol since they are not used during
deserialization.

### Aliases

A simple type alias:

```yaml
MyString: string
```

Is converted to:

```JSON
{
  "alias": {
    "name": "MyString",
    "type": "string"
  }
}
```

If an alias is generic:

```yaml
MyVector<T> : !vector
  items: T
```

Its JSON looks like this:

```JSON
{
  "alias": {
    "name": "MyVector",
    "typeParameters": [
      "T"
    ],
    "type": {
      "vector": {
        "items": "T"
      }
    }
  }
}
```

### Protocols

A protocol that looks like this in Yardl:

```yaml
MyProtocol : !protocol
  sequence:
    a: string
    b: !stream
      items: double
```

Is represented in JSON like this:

```JSON
{
  "name": "MyProtocol",
  "sequence": [
    {
      "name": "a",
      "type": "string"
    },
    {
      "name": "b",
      "type": {
        "stream": {
          "items": "float64"
        }
      }
    }
  ]
}
```

### Top-Level Schema

The JSON that is embedded in the binary or HDF5 format contains the protocol
definition (defined above) and the transitive closure of named types (records,
enums, and aliases) used by the protocol.

```JSON
{
  "protocol": <protocol>,
  "types": [ <type>, ... ]
}
```
