# The Yardl Language

:::info Note

The Matlab implementation supports the binary format, but does not
currently support HDF5 or NDJSON.

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

```matlab
w = sandbox.binary.MyProtocolWriter("sandbox.bin");

w.write_a(1);           % Write single protocol step

w.write_b(1:10);        % Write stream items
w.end_b();              % Call stream end method to signal that it is complete

w.write_c(["a", "b"]);  % Write stream items
w.write_c(["c", "d"]);  % Write more stream items
w.end_c();              % Signal that the 'c' stream is complete

w.close();              % Must close the Writer
```

And read the data back like this:

```matlab
r = sandbox.binary.MyProtocolReader("sandbox.bin");

disp(r.read_a());       % Read single protocol step

while r.has_b()         % Check whether stream has ended
    disp(r.read_b());   % Read stream items in a loop
end

while r.has_c()         % Check whether stream has ended
    disp(r.read_c());   % Read stream items in a loop
end

r.close();              % Must close the reader
```

It is an error to attempt to read or write a protocol's steps out of order or to
close a reader or writer without having written or read all steps.

```matlab
w = sandbox.binary.MyProtocolWriter("sandbox.bin");
w.write_b(1:10);  % Error: Expected to call to 'write_a' but received call to 'write_b' // [!code error]
```

Generated protocol readers have a `copy_to()` method that allows you to copy the
contents of the protocol to another protocol writer. This makes is easy to, say,
read from an NDJSON file and send the data in the binary format over a network
connection.

## Records

Records have fields and, optionally, [computed fields](#computed-fields). In
generated Matlab code, they map to class definitions.

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

The generated class constructors take arguments for each field and
in most cases they are optional. They are only required when the field type is
[generic](#generics).

## Primitive Types

Yardl has the following primitive types:

| Yardl Type       | Comment                                                                 | Matlab Type        |
| ---------------- | ----------------------------------------------------------------------- | -------------------|
| `bool`           |                                                                         | `logical`          |
| `int8`           |                                                                         | `int8`             |
| `uint8`          |                                                                         | `uint8`            |
| `byte`           | Alias of `uint8`                                                        |                    |
| `int16`          |                                                                         | `int16`            |
| `uint16`         |                                                                         | `uint16`           |
| `int32`          |                                                                         | `int32`            |
| `int`            | Alias of `int32`                                                        |                    |
| `uint32`         |                                                                         | `uint32`           |
| `uint`           | Alias of `unit32`                                                       |                    |
| `int64`          |                                                                         | `int64`            |
| `long`           | Alias of `int64`                                                        |                    |
| `uint64`         |                                                                         | `uint64`           |
| `ulong`          | Alias of `uint64`                                                       |                    |
| `size`           | Equivalent to `uint64`                                                  |                    |
| `float32`        |                                                                         | `single`           |
| `float`          | Alias of `float32`                                                      |                    |
| `float64`        |                                                                         | `double`           |
| `double`         | Alias of `float64`                                                      |                    |
| `complexfloat32` | A complex number where each component is a 32-bit floating-point number | `complex(single)`  |
| `complexfloat`   | Alias of `complexfloat32`                                               |                    |
| `complexfloat64` | A complex number where each component is a 64-bit floating-point number | `complex(double)`  |
| `complexdouble`  | Alias of `complexfloat64`                                               |                    |
| `string`         |                                                                         | `string`           |
| `date`           | A number of days since the epoch                                        | `yardl.Date`       |
| `time`           | A number of nanoseconds after midnight                                  | `yardl.Time`       |
| `datetime`       | A number of nanoseconds since the epoch                                 | `yardl.DateTime`   |

`yardl.Date`, `yardl.Time`, and `yardl.DateTime` are custom classes because
Yardl uses nanosecond precision and Matlab's `datetime` has only microsecond precision.
Each of them can be easily converted to/from a Matlab `datetime` by calling
the corresponding `to_datetime()` or `from_datetime()` method.

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

In Matlab, optional values can be instantiated using `yardl.Optional(value)`, and have the value `yardl.None` when not set.

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

The `null` type in the example above means that no value (`yardl.None` in Matlab) is
also a possibility.

### Generated Union Types

For the Matlab codegen, we represent unions as a [tagged
union](https://en.wikipedia.org/wiki/Tagged_union) and we generate a class for
each union in a model. In the example above, the union `[int, float]` would have
a class named `Int32OrFloat32`. Each union case then has a corresponding static
constructor method in the union class.

This union class can be used like this:

```matlab
function process_my_union(u)
    if u.isInt32()
        fprintf("%d is an int32\n", u.value);
    elseif u.isFloat32()
        fprintf("%f is a float\n", u.value);
    else
        error("Unrecognized union type");
    end
end

process_my_union(sandbox.Int32OrFloat32.Int32(2))
process_my_union(sandbox.Int32OrFloat32.Float32(7.9))
```

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

Enums are generated as custom Matlab class definitions, *not* using Matlab's `enumeration` support, which doesn't allow integer values that are not explicitly defined in the enum definition.

In Matlab, use the enum constructor to create new values, or the generated static methods for predefined values:

```matlab
fruit1 = sandbox.Fruits.APPLE;
fruit2 = sandbox.Fruits.BANANA;
newFruit = sandbox.Fruits(42);
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
  base: uint8
  values:
    read: 1
    write: 2
    execute:
```

Any value without an integer value will have the next power of two bit set that
is greater than the previous value. In the example above, `execute` would have
the value 4.

Like Enums, Flags are generated as custom Matlab class definitions.

Example usage:

```matlab
permissions = bitor(sandbox.Permissions.READ, sandbox.Permissions.EXECUTE);
% ...
if bitand(sandbox.Permissions.READ, permissions)
    % ...
end
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

These map to the Matlab `dictionary` type.

:::info Note

Matlab's `dictionary` type was introduced in Matlab r2022b, effectively replacing
the `containers.Map` type. The `containers.Map` does not provide sufficient
support for yardl types (including primitive strings) as keys or values.

:::


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

Both flavors of vectors are generated as Matlab arrays.
The generated protocol readers/writers also support cell arrays as vectors.
When working with vectors of vectors, the last dimension represents the outer vector:

```matlab
% Create a vector of integers
vs = [1, 2, 3, 4];

% Create a vector of 3 vectors, each of length 4
vs = [[1; 2; 3; 4], [5; 6; 7; 8], [9; 10; 11; 12]];

% Create a vector of 3 vectors, each of varying length
vs = { [1, 2, 3], [4, 5], [7] };
```

## Arrays

Arrays are multidimensional and map to Matlab arrays. Like vectors, there is a simple
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

### Matlab Arrays

In Matlab, arrays are always created with dimensions reversed with respect to the model definition.
This means that an array defined as `Image: double[x, y, z]` has the shape `[z, y, x]` in Matlab.

Yardl currently supports serializing multi-dimensional arrays only in
C-continguous order, where the last dimension increases most rapidly.
Matlab, however, uses Fortran-order to store and serialize multi-dimensional
arrays, where the first dimension increases most rapidly.

By reversing Matlab array dimensions, yardl maintains consistency with Matlab's
support for multi-dimensional array indexing, and provides optimal serialization performance.

As a side effect, if you define a *matrix* in yardl as `matrix: double[row, col]`,
you will need to transpose the array in Matlab.

Example:

```yaml
MyProtocol: !protocol
  sequence:
    fixedArray: double[x:2, y:4]
```

```matlab
>> r = sandbox.binary.MyProtocolReader(filename);
>> image = r.read_image();
>> size(image)

ans =

     4     2

```

To create an array with more than two dimensions, use Matlab pages:

```yaml
ndarray: double[2, 3, 4]
```

```matlab
>> ndarray(:, :, 1) = [[ 1;  2;  3;  4], [ 5;  6;  7;  8], [ 9; 10; 11; 12]];
>> ndarray(:, :, 2) = [[13; 14; 15; 16], [17; 18; 19; 20], [21; 22; 23; 24]];
>> size(ndarray)

ans =

     4     3     2

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

In Matlab, aliases are generated in one of three forms:
1. Union class definition for union types
2. Function wrapper for optionals/vectors/arrays
3. Subclass definitions for all other types

In all cases, you can use the generated syntax to construct the aliased type.

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

```matlab
>> rec = sandbox.MyRec(sandbox.Int32OrNamedArray.Int32(4));
>> rec.my_union_size()

ans =

    1

>> arr = sandbox.NamedArray(int32(ones(7)));
>> rec = sandbox.MyRec(sandbox.Int32OrNamedArray.NamedArray(arr));
>> rec.my_union_size()

ans =

    49

>> rec = sandbox.MyRec(yardl.None);
>> rec.my_union_size()

ans =

     0

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
    point: Point<int>
```

Here `Point` is a generic type with one type parameter `T`, while `MyProtocol`
references `Point` with `int` as its type argument.

Records and type aliases can be generic, but enums, flags, and protocols cannot.

In Matlab, generics are treated as open types.
Type validation occurs when values are written using a ProtocolWriter.

```matlab
p = sandbox.Point(1, 2);

w = sandbox.binary.MyProtocolWriter("sandbox.bin");
w.write_point(p);
w.close();

p = sandbox.Point("a", "b");
w = sandbox.binary.MyProtocolWriter("sandbox.bin");
w.write_point(p); % Error: ...Value must be of type int32 or be convertible to int32. // [!code error]
w.close();
```

## Imported Types

Types can be imported from other packages (see [Packages](packages)) and referenced
through their respective yardl namespace:

```yaml
MyTuple: BasicTypes.Tuple<string, int>
```

Imported types are likewise namespaced in Matlab packages:

```matlab
t1 = basic_types.Tuple("John Smith", 42);
t2 = sandbox.MyTuple("John Smith", 42);
assert(t1 == t2);
```

Note that yardl ignores protocols defined in imported packages.
