# Yardl

Yardl is a simple schema language and command-line tool that generates domain
types and serialization code.

<details>
<summary>Simple example</summary>
Given a Yardl definition like this:

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
    data: int*
```

After running `yardl generate`, you can write code like the following to write
data to standard out in a compact binary format:

```python
import sys
from sandbox import BinaryMyProtocolWriter, Header, Sample, DateTime

def generate_samples():
    yield Sample(timestamp=DateTime.now(), data=[1, 2, 3])
    yield Sample(timestamp=DateTime.now(), data=[4, 5, 6])

with BinaryMyProtocolWriter(sys.stdout.buffer) as w:
    w.write_header(Header(subject="Me"))
    w.write_samples(generate_samples())
```

And then another script can read it in from standard in:

```python
import sys
from sandbox import BinaryMyProtocolReader

with BinaryMyProtocolReader(sys.stdin.buffer) as r:
    print(r.read_header())
    for sample in r.read_samples():
        print(sample)
```

</details>

## Motivation

It is conceptually similar to, and inspired by, [Avro](https://avro.apache.org/),
[Protocol Buffers](https://developers.google.com/protocol-buffers),
[Bond](https://microsoft.github.io/bond/), and others, but it was designed
primarily with raw medical instrument signal data in mind. Some of its features
are:

- Persistence to [HDF5](https://www.hdfgroup.org/solutions/hdf5/) files as well
  as a compact binary format suitable for streaming over a network. There is
  also a much less efficient NDJSON format that is easier to manually inspect or
  use with other tools.
- Built-in support for multidimensional arrays and complex numbers.
- The schema is always embedded in the serialized data
- "Clean" generated code with types that are easy to program against.
- Generics
- Computed fields

Modeling a data domain in Yardl brings a number of benefits over writing the
code by hand:

- Writing correct and efficient serialization code can be tricky
- Schema versioning, compatibility, and conversions are handled for you
- You do not need to worry about consistency across different programming
  languages
- Comments could be used to generate documentation

## Getting Started

Please check out the project [documentation](https://aka.ms/yardl).

## Project Status

We are releasing this project order to get community feedback and contributions.
It is not complete and is **not ready for production use** at this time. We
expect to introduce breaking changes until the project reaches V1.

We currently support C++, Python, and MATLAB codegen. Other planned features include:

- Reading data with a different schema version
- References between packages
- Validating schema evolution is non-breaking
- Constraints
- Improvements to the language and editing experience

## Building the Code in this Repo

We recommend opening repo in a [dev
container](https://code.visualstudio.com/docs/devcontainers/containers) or a
[codespace](https://docs.github.com/en/codespaces/overview). Otherwise, all the
required dependencies are specified in the
[Conda](https://docs.conda.io/en/latest/) [environment.yml](environment.yml)
file in the repo root.

We use the [`just`](https://github.com/casey/just) command runner to build and run tests. To get
started, you should be able to run

```bash
$ just
```

from the repo root.

To enable support for Matlab, you must provide a Matlab license file to the devcontainer.
In your HOST environment, export the environment variable `MATLAB_LICENSE_FILE`,
e.g. in `$HOME/.profile`

```bash
export MATLAB_LICENSE_FILE=/mnt/c/Users/username/Documents/MATLAB/license.lic
```

Then invoke `just matlab=enabled ...`.

## Contributing

This project welcomes contributions and suggestions.  Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution. For
details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., status check,
comment). Simply follow the instructions provided by the bot. You will only need
to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of
Conduct](https://opensource.microsoft.com/codeofconduct/). For more information
see the [Code of Conduct
FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact
[opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional
questions or comments.

## Trademarks

This project may contain trademarks or logos for projects, products, or
services. Authorized use of Microsoft trademarks or logos is subject to and must
follow [Microsoft's Trademark & Brand
Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must
not cause confusion or imply Microsoft sponsorship. Any use of third-party
trademarks or logos are subject to those third-party's policies.
