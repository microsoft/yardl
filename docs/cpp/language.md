# The Yardl Language

Yardl model files use YAML syntax and are required to have either a `.yml` or
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

## Protocols

As explained in the [quick start](quickstart), protocols define a sequence of
values, called "steps", that are required to be transmitted, in order. They are
defined like this:

```yaml
MyProtocol: !protocol
  sequence:
    a: int
    b: !stream
      items: float
    c: !stream
      items: string
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

## Records

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
    recA: !record # NOT SUPPORTED! // [!code error]
      fields: // [!code error]
        a: int // [!code error]
    recB: RecordB # But this is fine.

RecordB: !record
  fields:
    c: int
```

Note that Yardl does not support type inheritance.

## Primitive Types

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
| `complexfloat64` | A complex number where each component is a 64-bit floating-point number |
| `complexdouble`  | Alias of `complexfloat64`                                               |
| `string`         |                                                                         |
| `date`           | A number of days since the epoch                                        |
| `time`           | A number of nanoseconds after midnight                                  |
| `datetime`       | A number of nanoseconds since the epoch                                 |

## Optional Types

If a value is optional, its type has a `?` suffix.

```yaml
Rec: !record
  fields:
    optionalInt: int?
```

They can also be expressed as a YAML array of length two, with `null` in the
first position:

```yaml
Rec: !record
  fields:
    optionalInt: [null, int] # equivalent to the example above
```

In C++, optional types are generated as `std::optional`.

## Unions

Optional types are a special case of unions, which are used when a value can be
one of several types:

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
```

The `null` type in the example above means that no value is also a possibility.

In C++, unions are generated as `std::variant`.

For Python codegen, we generate an identifier "tag" for each case of the union
based on its type. This means that you will get an error when attempting to use
a type in a union that is made up of symbols that are not valid in an identifier.
For example:

```yaml
Rec: !record
  fields:
    floatArrayOrDoubleArray:
      - float[] # Error! // [!code error]
      - double[] # Error! // [!code error]
```

There are two simple solutions to this problem. The first is to give explicit
tag names to each union case using the expanded `!union` syntax:

```yaml
Rec: !record
  fields:
    floatArrayOrDoubleArray: !union
      floatArray: float[]
      doubleArray: double[]
```

The second option is to create [aliases](#type-aliases) for the types:

```yaml

FloatArray: float[]
DoubleArray: double[]

Rec: !record
  fields:
    floatArrayOrDoubleArray:
      - FloatArray
      - DoubleArray
```

The first option is usually preferred, unless the type alias is going to be
used elsewhere.

## Enums

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

Any integer values that are left blank will be:

- 0 if the first value
- 1 greater than the previous value if positive
- 1 less that the previous value if negative.

Enums are generated as C++ enum classes (scoped enumerations).

## Flags

Flags are similar to enums but are meant to represent a bit field, meaning
multiple values can be set at once.

They can be defined with automatic values:

```yaml
Permissions: !flags
  values:
    - read
    - write
    - execute
```

Or with explicit values and an optional base type:

```yaml
Permissions: !flags
  base: unit8
  values:
    read: 1
    write: 2
    execute:
```

Any value without an integer value will have the next power of two bit set that
is greater than the previous value. In the example above, `execute` would have
the value 4.

For C++, we generate a special class with overloaded operators `|`, `&`, `^`,
and `~` and convenience methods `HasFlags()`, `SetFlags()`, `UnsetFlags()`, etc and
static const member variables for each flag value defined in the model.

Example usage:

```cpp
auto permissions = Permissions::kRead | Permissions::kWrite;
// ...
if (permissions.HasFlags(Permissions::kRead)) {
// ...
}
```

## Vectors

Vectors are one-dimensional lists. They can optionally have a fixed length. The
simple syntax for vectors is `<type>*[length]`.

For example:

```yaml
MyRec: !record
  fields:
    vec1: int*
    vec2: int*10
```

In the example above, `vec1` is a vector of integers of unknown length and
`vec2` has length 10. The expanded syntax for vectors is:

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

## Arrays

Arrays are multidimensional. Like vectors, there is a simple
syntax and an expanded syntax for declaring them. Both syntaxes are shown in the
examples below.

There are three kinds of arrays. They can be of a fixed size:

```yaml
MyRec: !record
  fields:
    fixedNdArray: float[3, 4]
    fixedNdArrayExpandedSyntax: !array
      items: float
      dimensions: [3, 4]
```

Or the size might not be fixed but the number of dimensions is known:

```yaml
MyRec: !record
  fields:
    ndArray: float[,]
    ndArrayExpandedSyntax: !array
      items: float
      dimensions: 2
```

Or finally, the number of dimensions may be unknown as well:

```yaml
MyRec: !record
  fields:
    dynamicNdArray: float[]
    dynamicNdArrayExpandedSyntax: !array
      items: float
```

Dimensions can be given names, which can be used in [computed
field](#computed-fields) expressions.

```yaml
MyRec: !record
  fields:
    fixedNdArray: float[x:3, y:4]
    fixedNdArrayExpandedSyntax: !array
      items: float
      dimensions:
        x: 3
        y: 4
    ndArray: !array
      items: float
      dimensions: [x, y]
    ndArrayExpandedSyntax: !array
      items: float
      dimensions: [x, y]
    ndArrayExpandedSyntaxAlternate: !array
      items: float
      dimensions:
        x:
        y:
```

In the simple syntax, `int[]` means an int array with an unknown number of
dimensions, and `int[,]` means an int array with two dimensions. To declare an
array with 1 dimension of unknown length, you can either give the dimension a
name (`int[x]`) or use parentheses to disambiguate from an empty set of
dimensions: `int[()]`.

In C++, we use the xtensor library for arrays.

## Maps

Maps, also known as dictionaries or associative arrays, are an unordered
collection of key-value pairs.

They can be declared like this:

```yaml
MyMap: string->int
```

Or declared with the expanded syntax:

```yaml
MyMap: !map
  keys: string
  values: int
```

Keys are required to be scalar primitive types.

In generated C++ code, these are generated as `std::unordered_map`.

## Type Aliases

Any type can be given one or more aliases:

```yaml
FloatArray: float[]

SignedInteger: [int8, int16, int32, int64]

Id: string
Name: string
```

This simply gives another name to a type, so the `Name` type above is no
different from the `string` type.

## Computed Fields

In addition to fields, records can contain computed fields. These are simple expressions
over the record's other (computed) fields.

```yaml
MyRec: !record
  fields:
    arrayField: int[x,y]
  computedFields:
    accessArray: arrayField
    accessArrayElement: arrayField[0, 1]
    accessArrayElementByName: arrayField[x:0, y:1]
    sizeOfArrayField: size(arrayField)
    sizeOfFirstDimension: size(arrayField, 0)
    sizeOfXDimension: size(arrayField, 'x')
```

To work with union types, you need to use a switch expression with type pattern
matching:

```yaml
NamedArray: int[x, y]

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

## Generics

Yardl supports generic types.

```yaml
Image<T>: T[]

ImageUnion: !union
  float: Image<float>
  double: Image<double>
  complexFloat: Image<complexfloat>
  complexDouble: Image<complexdouble>

RecordWithImages<T, U>: !record
  fields:
    image1: Image<T>
    image2: Image<U>
```

Note that protocols cannot be open generic types, but their steps may be made up of
closed generic types (e.g. `Image<float>`). Enums and Flags cannot be generic either.

Generic types map to C++ template classes.
