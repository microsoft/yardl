# Schema Evolution

Yardl supports schema evolution over time, including backward compatibility with previous schema versions.

See [Packages](packages) for details on specifying previous schema versions.


## Automated Change Detection

For each previous version of your schema, yardl will automatically determine:

1. How the schema changed, and
2. Whether the changes are backward compatible.

Automated change detection is performed by:

1. Matching named types across the schema for each previous version
2. Comparing semantically equivalent "base" type definitions
3. Generating compatibility serializers for all type definitions that changed between versions
4. Comparing protocol sequences step by step to validate each change

In the future, yardl will allow you to explicitly define how your schema is meant to evolve.


## Backward Compatibility

Yardl recursively detects changes to named types, fields, and protocol steps, classifying each granular change into one of three categories:

1. Backward compatible
2. Partially-compatible
3. Incompatible


### Backward Compatible Changes

Backward compatible changes are fully supported by yardl. Examples include:

1. Adding streams, vectors, and/or optional steps to a Protocol sequence
1. Adding or removing optional fields
2. Adding or removing aliases to types
3. Reordering Record fields


### Partially-Compatible Changes

Partially-compatible changes are valid, but may result in errors at runtime, depending on the data you serialize for older versions of your Protocol.

1. Changing between primitive types, including integers, floating point values, and strings
2. Making a field optional
3. Changing an optional field to a union, and vice versa
3. Adding or removing non-optional fields to/from a Record
4. Adding or removing types to/from a Union

yardl will emit a warning for each of these types of changes.


#### Default conversion

Some partially-compatible changes use default type conversions, for example:

- Converting between numbers and strings in C++ relies on the standard library numeric parsing utilities
- Converting floating point numbers to integers may `round` to the nearest whole number.

In the future, yardl will parse user-defined type conversions for each schema version.


#### Default values

In instances where partially-compatible change may result in invalid values at runtime, yardl defaults to the "zero" value for a type, e.g. `0` for numbers, `""` for string, empty vectors, null Optional/Union, etc.

For example, say you change field `description` from a `string` to a `string?` in the latest version of your schema. Now, when `description` is empty, its value is just `null`, because its type is Optional string. To maintain compatibility with software using the older version of the schema, however, when `description` is empty, yardl will write the empty string `""`.


### Incompatible Changes

Changes for which yardl (currently) cannot generate valid code include:

1. Removing or reordering any steps in a Protocol sequence
2. Changing the inner type of a stream, vector, array, or map
3. Changing generic type parameters
4. Changing enum/flag definitions

Detecting these types of changes will cause yardl to emit one or more errors and stop.

## Examples

### Adding a new type to a Protocol stream

Tip: Use an alias (e.g. `StreamItem` in the example below) to easily add new types in the future

Previous Version:
```yaml
ImageFloat: Image<float>
StreamItem: ImageFloat

MyProtocol: !protocol
  sequence:
    data: StreamItem
```

Latest Version, with new types added to my data stream:
```yaml
Acquisition: !record
  ...
WaveformUint32: Waveform<uint32>
ImageInt16: Image<int16>
ImageFloat: Image<float>
ImageComplexDouble: Image<complexdouble>

StreamItem: [Acquisition, WaveformUint32, ImageInt16, ImageFloat, ImageComplexDouble]

MyProtocol: !protocol
  sequence:
    data: StreamItem
```
