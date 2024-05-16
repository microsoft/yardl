# Quick Start

## Installation

<!--@include: ../parts/installation-core.md-->

### Dependencies

The generated MATLAB code requires MATLAB R2022b or newer.

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

matlab:
  outputDir: ../matlab
```

It specifies the package's namespace along with code generation settings. The
`matlab.outputDir` property specifies where the generated MATLAB package should
go. If you are not interested in generating Python or C++ code, you can remove
the corresponding property from the file:

``` yaml
namespace: Playground

cpp: // [!code --]
  sourcesOutputDir: ../cpp/generated // [!code --]

python: // [!code --]
  outputDir: ../python // [!code --]

matlab:
  outputDir: ../matlab
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

To generate MATLAB code for this model, `cd` into the `model` directory and run:

```bash
yardl generate
```

This will generate a MATLAB package in the `outputDir` directory:

```txt
$ tree .
.
├── +playground
│   ├── +binary
│   │   ├── HeaderSerializer.m
│   │   ├── MyProtocolReader.m
│   │   ├── MyProtocolWriter.m
│   │   └── SampleSerializer.m
│   ├── Header.m
│   ├── MyProtocolReaderBase.m
│   ├── MyProtocolWriterBase.m
│   └── Sample.m
└── +yardl
    ├── +binary
    │   └── ...
    └── ...
```

The top-level package, e.g. `+playground`, contains the class definitions for (1) the non-protocol types defined in our model (in this case, `Header.m` and `Sample.m`), and (2) the abstract protocol reader and writer classes, from which concrete implementations inherit from (e.g. in the `+binary` subpackage).

The adjacent `+yardl` package contains definitions for primitive types, error handling, and serializers.

To use these packages from outside of the `matlab` directory, use MATLAB's `addpath` function, e.g. `addpath("../path/to/parent/directory");`.

Ok, let's write some code! In our `matlab` directory (containing the generated
`+playground` package), create `run_playground.m` that looks like this:

```matlab
Sample = @playground.Sample;
samples = [Sample(yardl.DateTime.now(), [1, 2, 3]), Sample(yardl.DateTime.now(), [4, 5, 6])];

path = "playground.bin";

w = playground.binary.MyProtocolWriter(path);
w.write_header(playground.Header("Me"));
w.write_samples(samples);
w.end_samples();
w.close();

r = playground.binary.MyProtocolReader(path);
disp(r.read_header());
while r.has_samples()
    sample = r.read_samples();
    disp(sample.timestamp.to_datetime());
    disp(sample.data);
end
r.close();
```

Run this directly in MATLAB, e.g. `run_playground`, or on the command-line with `matlab -batch run_playground`.

You can inspect the binary file our code produced with:

```bash
hexdump -C playground.bin
```

Note that the binary file contains a JSON representation of the protocol's
schema. This allows code that was not previously aware of this protocol to
deserialize the contents.
