# The Yardl Language

:::info Note

The Python implementation supports the binary and NDJSON formats, but does not
currently support HDF5.

:::

Yardl model files use YAML syntax and are required to have either a `.yml` or
`.yaml` file extension.

To efficiently work with Yardl, we recommend that you run the following from the
yardl package (model) directory:

```bash
yardl generate --watch
```

This watches the directory for changes and generates code whenever a file is
saved. This allows you to get rapid feedback as you experiment.

Comments placed above top-level types and their fields are captured and added to
the generated code as docstrings.

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

You can write to this protocol like this:

```python
with NDJsonMyProtocolWriter(sys.stdout) as w:
    w.write_a(1)

    w.write_b(float(i) for i in range(10))

    w.write_c(["a", "b"])
    w.write_c(["c", "d"])  # Add more to the "c" stream
```

And read the data back like this:

```python
with NDJsonMyProtocolReader(sys.stdout) as r:
    print(r.read_a())
    for b in r.read_b():
        print(b)
    for c in r.read_c():
        print(c)
```

It is an error to attempt to read or write a protocol's steps out of order or to
close a reader or writer without having written or read all steps.

```python
with BinaryMyProtocolWriter(sys.stdout.buffer) as w:
    w.write_b(float(i) for i in range(10)) # Error: Expected to call to 'write_a' but received call to 'write_b' // [!code error]
```

Generated protocol readers have a `copy_to()` method that allows you to copy the
contents of the protocol to another protocol writer. This makes is easy to, say,
read from an NDJSON file and send the data in the binary format over a network
connection.

## Records

Records have fields and, optionally, [computed fields](#computed-fields). In
generated Python code, they map to classes.

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

The generated class constructors take keyword-only arguments for each field and
in most cases they are optional. They are only required when the field type is
[generic](#generics).

## Primitive Types

Yardl has the following primitive types:

| Yardl Type       | Comment                                                                 | Python Type           | Underlying Python Type |
| ---------------- | ----------------------------------------------------------------------- | --------------------- | ---------------------- |
| `bool`           |                                                                         | `bool`                |                        |
| `int8`           |                                                                         | `yardl.Int8`          | `int`                  |
| `uint8`          |                                                                         | `yardl.UInt8`         | `int`                  |
| `byte`           | Alias of `uint8`                                                        |                       |                        |
| `int16`          |                                                                         | `yardl.Int16`         | `int`                  |
| `uint16`         |                                                                         | `yardl.UInt16`        | `int`                  |
| `int32`          |                                                                         | `yardl.Int32`         | `int`                  |
| `int`            | Alias of `int32`                                                        |                       |                        |
| `uint32`         |                                                                         | `yardl.UInt32`        | `int`                  |
| `uint`           | Alias of `unit32`                                                       |                       |                        |
| `int64`          |                                                                         | `yardl.Int64`         | `int`                  |
| `long`           | Alias of `int64`                                                        |                       |                        |
| `uint64`         |                                                                         | `yardl.UInt64`        | `int`                  |
| `ulong`          | Alias of `uint64`                                                       |                       |                        |
| `size`           | Equivalent to `uint64`                                                  | `yardl.Size`          | `int`                  |
| `float32`        |                                                                         | `yardl.Float32`       | `float`                |
| `float`          | Alias of `float32`                                                      |                       |                        |
| `float64`        |                                                                         | `yardl.Float64`       | `float`                |
| `double`         | Alias of `float64`                                                      |                       |                        |
| `complexfloat32` | A complex number where each component is a 32-bit floating-point number | `yardl.ComplexFloat`  | `complex`              |
| `complexfloat`   | Alias of `complexfloat32`                                               |                       |                        |
| `complexfloat64` | A complex number where each component is a 64-bit floating-point number | `yardl.ComplexDouble` | `complex`              |
| `complexdouble`  | Alias of `complexfloat64`                                               |                       |                        |
| `string`         |                                                                         | `str`                 |                        |
| `date`           | A number of days since the epoch                                        | `datetime.date`       |                        |
| `time`           | A number of nanoseconds after midnight                                  | `yardl.Time`          |                        |
| `datetime`       | A number of nanoseconds since the epoch                                 | `yardl.DateTime`      |                        |

`yardl.Int8`, `yardl.UInt8`, `yardl.Int16`, `yardl.UInt16`, `yardl.Int32`,
`yardl.UInt32`, `yardl.Size` are all annotated aliases of `int` for the purposes
of Python [type hinting](https://docs.python.org/3/library/typing.html).
Similarly, `yardl.Float32` and `yardl.Float64` are aliases of `float`, and
`yardl.ComplexFloat` and `yardl.ComplexDouble` are aliases of `complex`.

`yardl.Time` and `yardl.DateTime` are custom time and date-time classes because
Yardl uses nanosecond precision and Python's `datetime.time` and
`datetime.datetime` have only microsecond precision.

:::info Note

Note The `yardl` qualifier used above is probably not how you will reference
these types. They are generated as part of the model package and should be
imported and used like any other type in the package. `mypackage.Int32` is a
more realistic example.

:::

:::info Note

The table above does not tell the full story, because these primitives have
different types within NumPy arrays. See [Arrays](#arrays) below.

:::

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

In Python, these have the type hint `Optional[T]` (equivalent to `T | None` in Python 3.10+) and
have the value `None` when not set.

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

The `null` type in the example above means that no value (`None` in Python) is
also a possibility.

### Generated Union Types

For the Python codegen, we represent unions as a [tagged
union](https://en.wikipedia.org/wiki/Tagged_union) and  we generate a class for
each union in a model. In the example above, the union `[int, float]` would have
a class named `Int32OrFloat32`. Each union case then has a nested type which
inherits from the base union type.

This union class can be used like this:

```python
def process_my_union(u: Int32OrFloat32):
    if isinstance(u, Int32OrFloat32.Int32):
        assert type(u.value) == int
        print(f"{u.value} in an int")
    elif isinstance(u, Int32OrFloat32.Float32):
        assert type(u.value) == float
        print(f"{u.value} is a float")
    else:
        raise ValueError(f"Unrecognized type {u}")


process_my_union(Int32OrFloat32.Int32(2))
process_my_union(Int32OrFloat32.Float32(7.9))
```

Or in Python 3.10+ with the `match` statement:

```python
def process_my_union(u: Int32OrFloat32):
    match u:
        case Int32OrFloat32.Int32():
            assert type(u.value) == int
            print(f"{u.value} in an int")
        case Int32OrFloat32.Float32():
            assert type(u.value) == float
            print(f"{u.value} is a float")
        case _:
            raise ValueError(f"Unrecognized type {u}")


process_my_union(Int32OrFloat32.Int32(2))
process_my_union(Int32OrFloat32.Float32(7.9))
```

The constructors and `value` fields are type hinted to help you use this feature.

::: details Why not use [Python type hinting unions](https://docs.python.org/3/library/stdtypes.html#types-union)?

Python type hinting unions are not suitable for the Yardl type system in all cases.

1. Unions exist as type hints only and there is no union information to query at runtime.
2. Several Yardl types map to the same runtime type. e.g. `int32` and `int64` both map to Python's `int`.
3. Lists and arrays could potentially require fully enumerating their contents
   to determine which union case they represent.

For these reasons, we opted for a mechanism where the union case is clearly indicated at runtime.

 :::

### Union Tag Names

The type case name used in union classes (`Int32` and `Float32` in the example
above) is derived from the name of the type of each case. This means that you
will get an error when attempting to use a type in a union that is made up of
symbols that are not valid in an identifier. For example:

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

The first option is usually preferred, unless the type alias is going to be used
elsewhere. In both cases, the union type will be `FloatArrayOrDoubleArray` and
the cases `FloatArray` and `DoubleArray`.

If you don't like the generated name you can give the union an
[alias](#type-aliases):

```yaml
ArrayUnion: !union
  floatArray: float[]
  doubleArray: double[]

Rec: !record
  fields:
    arrayUnion: ArrayUnion
```

Now the union class is `ArrayUnion` with case types `FloatArray` and `DoubleArray`.

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

Enums are generated as Python `enum.Enum`s, but we customize the behavior to allow integer values that are outside of the defined values. This is support future versioning capabilities.

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

These are generated as Python [`enum.IntFlag`](https://docs.python.org/3/library/enum.html#enum.IntFlag) classes.

Example usage:

```python
permissions = sandbox.Permissions.READ | sandbox.Permissions.EXECUTE
# ...
if sandbox.Permissions.READ in permissions:
    # ...
```

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

These map to Python dictionaries.

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

Both flavors of vectors are generated as Python lists.

## Arrays

Arrays are multidimensional and map to NumPy arrays. Like vectors, there is a simple
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

### NumPy Types

Arrays map to the Numpy `ndarray` type. Whereas in standard Python, types like
`int` and `float` are used to represent numerical values, NumPy introduces its
own set of types, which are generally fixed-size and designed for efficiency
within large arrays. A Yardl `int32[2, 3]` array therefore becomes a
`np.ndarray(shape=(2, 3), dtype=np.int32)`.

The following table summarizes how different Yardl types are represented in
"standard" Python and within NumPy arrays:

| Yardl Type       | Type Outside of NumPy Array | dtype in NumPy array              | Comment                                                                                                                                                                                    |
| ---------------- | --------------------------- | --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `bool`           | `bool`                      | `np.bool_`                        |
| `int8`           | `int`                       | `np.int8`                         |
| `uint8`          | `int`                       | `np.uint8`                        |
| `int16`          | `int`                       | `np.int16`                        |
| `uint16`         | `int`                       | `np.uint16`                       |
| `int32`          | `int`                       | `np.int32`                        |
| `uint32`         | `int`                       | `np.uint32`                       |
| `int64`          | `int`                       | `np.int64`                        |
| `uint64`         | `int`                       | `np.uint64`                       |
| `size`           | `int`                       | `np.uint64`                       |
| `float32`        | `float`                     | `np.float32`                      |
| `float64`        | `float`                     | `np.float64`                      |
| `complexfloat32` | `complex`                   | `np.complex64`                    |
| `complexfloat64` | `complex`                   | `np.complex128`                   |
| `string`         | `str`                       | `np.object_` (`str`)              | Since strings are variable-length, we use normal heap-allocated Python strings.                                                                                                            |
| `date`           | `datetime.date`             | `np.datetime64[D]`                |
| `time`           | `yardl.Time`                | `np.timedelta64[ns]`              | `np.timedelta64[ns]` stores the integer nanoseconds since midnight.                                                                                                                        |
| `datetime`       | `yardl.DateTime`            | `np.datetime64[ns]`               |
| optional         | `Optional[python_type]`     | structured record                 | This becomes a [structured record](https://numpy.org/doc/stable/user/basics.rec.html) with the following fields: `{"has_value": np.bool_, "value": inner_numpy_type }`                        |
| union            | Tagged union class          | `np.object_` (`python_type`)      | Array values the tagged unions of Python types. NumPy types are not used.                                                                                                                  |
| fixed vector     | `list[python_type]`         | subarray                          | An array of fixed vectors becomes a single NumPy array with increased dimensionality. Fixed vectors in records become [subarrays](https://numpy.org/doc/stable/glossary.html#term-subarray).  |
| dynamic vector   | `list[python_type]`         | `np.object` (`list[python_type]`) | Because the size is not fixed, these vectors are stored as lists of the normal Python type.                                                                                                |
| fixed array      | `np.ndarray`                | subarray                          | An array of fixed arrays results in a single NumPy array with increased dimensionality. Fixed arrays in records become [subarrays](https://numpy.org/doc/stable/glossary.html#term-subarray). |
| dynamic array    | `np.ndarray`                | `np.object_` (np.ndarray)         | Because the shape is not fixed, these are stored as nested `np.ndarray`s.                                                                                                                  |
| map              | dictionary                  | `np.object` (dictionary)          | In an array, dictionaries are represented in the same way as they are outside of an array.                                                                                                 |
| record           | class                       | structured record                 | See [note below](#structured-arrays).                                                                                                                                                                            |
| enum/flag        | enum/flag class             | underlying NumPy integer type     |

#### Structured Arrays

Records have a very different representation in a NumPy array compared to
outside an array. Normally represented as a class, a record becomes a fixed-size [structured
record](https://numpy.org/doc/stable/user/basics.rec.html) within a NumPy array.

Suppose we have the following Yardl:

```yaml
Point: !record
  fields:
    x: double
    y: double
```

The normal generated Python class for `Point` has `yardl.Float32` (aliases of
`float`) fields `x` and `y`. But for an array of these records, you would create
a structured array:

```python
import numpy as np

dt = np.dtype([("x", np.float32), ("y", np.float32)])
arr = np.array([(1, 2), (3, 4)], dtype=dt)

# set a field
arr[0]["x"] = 8

# set an array value
arr[0] = (7, 10)

```

#### The `get_dtype()` Function

When passing a NumPy array to a protocol writer, the writer checks that array
has the correct dtype. For structured arrays, the dtype can either be [aligned
or
unaligned](https://numpy.org/doc/stable/user/basics.rec.html#automatic-byte-offsets-and-alignment).
In can be cumbersome to write out the dtype definition by hand, so the generated
Python module contains a `get_dtype()` function to help convert a Python type to
a NumPy dtype.

```python
>>> get_dtype(Int32)
dtype('int32')
>>> get_dtype(Int32)
dtype('int32')
>>> get_dtype(Float32)
dtype('float32')
>>> get_dtype(str)
dtype('O')
>>> get_dtype(Point)
dtype([('x', '<f8'), ('y', '<f8')], align=True)
>>> get_dtype(GenericPoint[Int32])
dtype([('x', '<i4'), ('y', '<i4')], align=True)
```

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

In Python, there are generated as [type aliases](https://docs.python.org/3/library/typing.html#type-aliases).

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
    accessArrayElementAndConvert: arrayField[0, 1] as int
    sizeOfArrayField: size(arrayField)
    sizeOfFirstDimension: size(arrayField, 0)
    sizeOfXDimension: size(arrayField, 'x')
    basicArithmentic:  arrayField[0, 1] * 2
```

The following expression types are supported:
- Numeric literals, such as `1`, `-1`, `0xF`, `3.4`, and `-2e-3`.
- String literals, such as `"abc"` and `'abc'`.
- Simple arithmethic expresions, such as `1 + 2`, `2.0 * 3`, and `2 ** 3` (`**`
  is the power operator and yields a `float64`).
- Type conversions using the `as` operator, such as `1 as float64`.
- Field accesses, such as `myField`. You can access a field on another field
  using the `.` operator, such as `myField.anotherField`.
- Array and vector element access, such as `arrayField[0, 1]` or
  `arrayField[x:0, y:1]` to identify the dimensions by name.
- Function calls:
  - `size(vector)`: returns the size (length) of the vector.
  - `size(array)`: returns the total size of the array.
  - `size(array, integer)`: returns the size of the array's dimension at the given
    index.
  - `size(array, string)`: returns the size of the array's dimension with the
    given name.
  - `dimensionIndex(array, string)` returns the index of the dimension with the
    given name.
  - `dimensionCount(array)` returns the dimension count of the array.

To work with union or optional types, you need to use a switch expression with type pattern
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

Computed fields become parameterless methods on the generated Python class. Here
is an example of invoking the field from the preceding Yardl definition:

```python
>>> rec = MyRec(my_union=Int32OrNamedArray.Int32(4))
>>> rec.my_union_size()
1
```

## Generics

Yardl supports generic types.

```yaml
Point<T>: !record
  fields:
    x: T
    y: T

MyProtocol: !protocol
  sequence:
    p: Point<int>
```

Here `Point` is a generic type with one type parameter `T`, while `MyProtocol`
references `Point` with `int` as its type argument.

Records and type aliases can be generic, but enums, flags, and
protocols cannot.

In Python, generics are supported through [type
hints](https://docs.python.org/3/library/typing.html#generics).

Often type arguments do not have to be specified thanks to type inference:

```python
with NDJsonMyProtocolWriter(sys.stdout) as w:
    w.write_p(Point(x=1, y=2)) # type argument inferred
```

But if necessary or preferred, type arguments can be supplied using subscription:

```python
with NDJsonMyProtocolWriter(sys.stdout) as w:
    w.write_p(Point(x=1, y=2)) # type argument inferred // [!code --]
    w.write_p(Point[Int32](x=1, y=2)) # type argument explicitly provided // [!code ++]
```

### Array Type Arguments

Type arguments that are used in arrays are constrained to be NumPy types.

```yaml
Image<T>: T[]
```

```python
img: Image[Int32] # Error // [!code error]
img2: Image[np.int32] # OK
```

But sometimes a type parameter is used as an array and as a scalar:

```yaml
Rec<T>: !record
  fields:
    scalar: T
    arr: T[]
```

In that case, the Python class is generated with two type parameters, one that
is unconstrained (`T` in this example), the other that is constrained to be a
NumPy type (`T_NP` in this example). There is unfortunately no way to constrain
the two type parameters to be of corresponding types. So for example, this
would ideally raise a typing error but doesn't:

```python
rec = Rec[Float32, np.int32](scalar=1.0, arr=np.array([1, 2, 3], np.int32))
```

However, using this value in a protocol writer will surface a typing error and a
runtime error:

```python
with NDJsonMyProtocolWriter(sys.stdout) as w:
    w.write_rec(rec) # typing error and runtime error // [!code error]
```

## Imported Types

Types can be imported from other packages (see [Packages](packages)) and referenced
through their respective yardl namespace:

```yaml
MyTuple: BasicTypes.Tuple<string, int>
```

Imported types are likewise namespaced in Python submodules:

```python
from my_package import basic_types
myInfo: MyTuple = basic_types.Tuple(v1="John Smith", v2=42)
```

Note that yardl ignores protocols defined in imported packages.
