# Quick Start

## Installation

<!--@include: ../parts/installation-core.md-->

### Dependencies

The generated Python code requires Python 3.9 or newer and you need to have
[NumPy](https://numpy.org/install/) version 1.22.0 or later installed.

## Getting our Feet Wet


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
`python.outputDir` property specifies where the generated python package should
go. If you are not interested in generating C++ code, you can remove the `cpp`
property from the file:

``` yaml
namespace: Playground

cpp: // [!code --]
  sourcesOutputDir: ../cpp/generated // [!code --]

python:
  outputDir: ../python
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

To generate Python code for this model, `cd` into the `model` directory and run:

```bash
yardl generate
```

This will generate a Python package in the `outputDir` directory:

```txt
$ tree .
.
└── playground
    ├── _binary.py
    ├── binary.py
    ├── _dtypes.py
    ├── __init__.py
    ├── _ndjson.py
    ├── ndjson.py
    ├── protocols.py
    ├── types.py
    └── yardl_types.py
```

`yardl_types.py` contains definitions of primitive data types. `types.py`
contains the definitions of the non-protocol types defined in our model (in this
case, `Header` and `Sample`). `protocols.py` contains abstract protocol reader
and writer classes, from which concrete implementations inherit from in
`binary.py` and `ndjson.py`.

Ok, let's write some code! in our `python` directory (containing the generated
`playground` directory), create `run_playground` that looks like this:

```python
from playground import (
    BinaryMyProtocolWriter,
    BinaryMyProtocolReader,
    Header,
    Sample,
    DateTime,
)


def generate_samples():
    yield Sample(timestamp=DateTime.now(), data=[1, 2, 3])
    yield Sample(timestamp=DateTime.now(), data=[4, 5, 6])


path = "playground.bin"

with BinaryMyProtocolWriter(path) as w:
    w.write_header(Header(subject="Me"))
    w.write_samples(generate_samples())

with BinaryMyProtocolReader(path) as r:
    print(r.read_header())
    for sample in r.read_samples():
        print(sample)

```

You can inspect the binary file our code produced with:

```bash
hexdump -C playground.bin
```

Note that the binary file contains a JSON representation of the protocol's
schema. This allows code that was not previously aware of this protocol to
deserialize the contents.

In addition to the compact binary format, we can write the protocol out to an
NDJSON file. This requires only a few modifications to our code:

```python
from playground import (
    BinaryMyProtocolWriter, // [!code --]
    NDJsonMyProtocolWriter, // [!code ++]
    BinaryMyProtocolReader, // [!code --]
    NDJsonMyProtocolReader, // [!code ++]
    Header,
    Sample,
    DateTime,
)


def generate_samples():
    yield Sample(timestamp=DateTime.now(), data=[1, 2, 3])
    yield Sample(timestamp=DateTime.now(), data=[4, 5, 6])


path = "playground.bin" // [!code --]
path = "playground.ndjson" // [!code ++]

with BinaryMyProtocolWriter(path) as w: // [!code --]
with NDJsonMyProtocolWriter(path) as w: // [!code ++]
    w.write_header(Header(subject="Me"))
    w.write_samples(generate_samples())

with BinaryMyProtocolReader(path) as r: // [!code --]
with NDJsonMyProtocolReader(path) as r: // [!code ++]
    print(r.read_header())
    for sample in r.read_samples():
        print(sample)

```
