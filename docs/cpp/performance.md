# Performance Tips

## Batched Reads and Writes

Generated protocol reader and writer classes have read and write methods for
each step. When a step is a stream, there will also be overloads that read or
write a batch of values as an `std::vector` in one go. For small data types,
using batched reads and writes can make a dramatic difference in throughput,
especially for HDF5 files.

## Use Fixed Data Types When Possible

For HDF5, using variable-length collections (like `!vector` without a length or
`!array` without fixed dimension sizes) has lower throughput than their fixed-sized
counterparts.

## Avoid NDJSON Encoding

The NDJSON serialization format is great for debugging and interoperability with
other tools (like `jq`) but is it orders of magnitude less efficient than the
binary or HDF5 formats.
