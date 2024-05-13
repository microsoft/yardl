## Motivation

Yardl is conceptually similar to, and inspired by, [Avro](https://avro.apache.org/),
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

## Project Status

We are releasing this project order to get community feedback and contributions.
It is not complete and is **not ready for production use** at this time. We
expect to introduce breaking changes until the project reaches V1.

We currently support C++ and Python codegen. Other planned features include:

- Python and MATLAB support for reading data with a different schema version
- Constraints
- Improvements to the language and editing experience
