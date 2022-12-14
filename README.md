# Yardl

Yardl is a simple schema language and command-line tool that generates domain
types and serialization code.

![A DSL on the left is translated to C++ code on the
right](docs/images/overview.png)

It is conceptually similar to, and inspired by, [Avro](https://avro.apache.org/),
[Protocol Buffers](https://developers.google.com/protocol-buffers),
[Bond](https://microsoft.github.io/bond/), and others, but it was designed
primarily with raw medical instrument signal data in mind. Some of its features
are:

- Persistence to [HDF5](https://www.hdfgroup.org/solutions/hdf5/) files as well
  as a compact binary format suitable for streaming over a network
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

Please check out the project [documentation](docs/docs.md).

## Project Status

We are releasing this project order to get community feedback and contributions.
It is not complete and is **not ready for production use** at this time. We
expect to introduce breaking changes until the project reaches V1.

We currently support C++ codegen and work will begin on Python soon. Other
planned features include:

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
