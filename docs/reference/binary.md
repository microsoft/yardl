# Compact Binary Encoding Reference

The details of the binary encoding format are provided here for reference. You
do not need to be familiar with this in order to use Yardl, as it takes care of
serialization and deserialization for you. To learn about the semantics of the
data types and how to use them, refer to the [Python](../python/language.md) or
[C++](../cpp/language.md) language guides.

::: warning
We might make breaking changes to this format before V1.
:::

The binary format starts with five magic bytes: `0x79 0x61 0x72 0x64 0x6c`
(ASCII 'y' 'a' 'r' 'd' 'l') followed by four bytes containing a little-endian
32-bit integer representing the encoding version number (currently 1). Then the
[protocol schema](protocol-schema) in JSON format written as a string
in the format described below.

After the schema, the protocol step values are written in order. The sections
below describe how each data type is encoded.

## Booleans

Booleans are encoded as a byte with the value 0 or 1. However, vectors or arrays
of booleans could use a single bit per value in order to save space. This issue
is being tracked [here](https://github.com/microsoft/yardl/issues/19).

## Unsigned Integers

Unsigned integers (`uint8`, `uint16`, `uint32`, `uint64`, and `size`) are
written as variable-width integers, or *varints* (in same way as Protocol
Buffers). The high-order bit of each byte serves as a continuation and indicates
whether more bytes remain. The lower seven bits of each byte are appended as
increasingly significant bits in the resulting value.

This allows smaller values to be encoded in fewer bytes.

Some examples:

| Integer value | First encoded byte | Second encoded byte |
| ------------- | ------------------ | ------------------- |
| 0             | 00000000           |                     |
| 1             | 00000001           |                     |
| 127           | 01111111           |                     |
| 128           | 10000000           | 00000001            |
| 129           | 10000001           | 00000001            |

## Signed integers

Signed integers (`int8`, `int16`, `int32`, and `int64`) are first converted to
unsigned integers using *zig-zag* encoding, and then encoded as unsigned
integers as above.

Because two's complement sets the highest bit of negative numbers, a negative
64-bit integer would always require 10 bytes to be encoded as a varint, making
it a poor choice. Instead of two's complement, zig-zag encoding stores positive
numbers as `2 * n` and negative numbers as `2 * abs(n) + 1`. This way, small
negative values can still be represented with a smaller number of bytes when
encoded as a varint.

Some examples:

| Original value | Encoded value |
| -------------- | ------------- |
| 0              | 0             |
| -1             | 1             |
| 1              | 2             |
| -2             | 3             |
| 2              | 4             |

## Floating-Point Numbers

`float32` and `float64` are written as a little-endian IEEE 754 of byte length
4 and 8, respectively. `complexfloat32` and `complexfloat64` write out first the
real part followed by the imaginary part.

## Strings

For strings, the length of the UTF8-encoded bytes is first written out as an
unsigned varint, followed by the UTF8-encoded bytes.

For example, the string "hello" is encoded as `Ox05 Ox68 Ox65 Ox6c Ox6c Ox6f`.

## Dates, Times, and DateTimes

Dates are written as a signed varint number of days since the epoch.

Times are written an a signed varint number of nanoseconds since midnight.

DateTimes are written as a signed varint number of nanoseconds since the epoch.

## Unions

Unions are written as the 0-based index of the type followed by the
value.

The index is written as an unsigned varint and the value if skipped if the type
is `null`.

Example values for the union `[null, uint, float]`:

| Value           | Encoded Bytes              |
| --------------- | -------------------------- |
| `<null>`        | `0x00`                     |
| 6 (`uint`)      | `0x01 0x06`                |
| 95.72 (`float`) | `0x02 0xa4 0x70 0xbf 0x42` |

## Vectors

If the vector length is not specified in Yardl:

```yaml
!vector
items: int
```

then the format is:

1. The length of the array as an unsigned varint
2. The array values written one after another.

If the length is given, then (1) is omitted.

## Arrays

If the number of array dimensions is not given in the Yardl schema, then the
format is:

1. The number of dimensions as a unsigned varint
2. Each dimension length as an unsigned varint
3. The values of the array in row-major order

If the number of dimensions is given in the schema, (1) is omitted. If the
length of each dimension is specified, (1) and (2) are omitted.

A future version of the binary format may support column-major layout. See
discussion [here](https://github.com/microsoft/yardl/issues/23).

## Maps

The format is:

1. The length of map as an unsigned varint
2. For each entry, the key followed by the value.

## Enums and Flags

Enums and flags are written as a varint encoding of the underlying integer
value. Note that the value is signed if the underlying type is signed, which is
the default case if the `base` property is not specified.

## Records

Records are encoded as the concatenation of the value of its fields, in the
order they appear in the schema. Note that there is no padding between fields.

## Streams

Streams are written as one or more blocks. Each block starts with a length as an
unsigned varint followed by that number of values. The last block will have
length 0 and will simply be `0x0`, which signals that the stream is complete.
Only the last block can have length 0.

## Example

Let's work through an example. Here is a sample model:

```yaml
MyProtocol: !protocol
  sequence:
    floatArray: float[2,2]
    points: !stream
      items: Point

Point: !record
  fields:
    x: uint64
    y: int32
```

We will write the values `{1.2, 3.4}, {5.6, 7.8}` as the `floatArray` step, and
5 points with coordinates `{1, 2}`, `{3, 4}`, `{5, 6}`, `{700, 800}`, and
`{800000, -900000}`. The points will be written in two blocks, the first of
length 3, the second of length 2. The C++ to write these values looks like this:

```cpp
writer.WriteFloatArray({{1.2, 3.4}, {5.6, 7.8}});

writer.WritePoints({{1, 2}, {3, 4}, {5, 6}});
writer.WritePoints({{700, 800}, {800000, -900000}});
writer.EndPoints();
```

Now let's look at the binary file. The first section of the file is the header
and schema. It begins with the magic bytes, the binary version, then the schema
as a JSON string. The string is written as its length (304) encoded as an an
unsigned varint followed by 304 chars.

```txt
ASCII:  y  a  r  d  l  .  .  .  .  .  .  {  "  p  r  o  t  o  c  o  l  "  :  {  "  n  a  m  e  "  :  "  M  y  P  r  o  t  o  c  o  l  "  ,  "  s  e  q  u  e  n  c  e  "  :  [  {  "  n  a  m  e  "  :  "  f  l  o  a  t  A  r  r  a  y  "  ,  "  t  y  p  e  "  :  {  "  a  r  r  a  y  "  :  {  "  i  t  e  m  s  "  :  "  f  l  o  a  t  3  2  "  ,  "  d  i  m  e  n  s  i  o  n  s  "  :  [  {  "  l  e  n  g  t  h  "  :  2  }  ,  {  "  l  e  n  g  t  h  "  :  2  }  ]  }  }  }  ,  {  "  n  a  m  e  "  :  "  p  o  i  n  t  s  "  ,  "  t  y  p  e  "  :  {  "  s  t  r  e  a  m  "  :  {  "  i  t  e  m  s  "  :  "  S  a  n  d  b  o  x  .  P  o  i  n  t  "  }  }  }  ]  }  ,  "  t  y  p  e  s  "  :  [  {  "  n  a  m  e  "  :  "  P  o  i  n  t  "  ,  "  f  i  e  l  d  s  "  :  [  {  "  n  a  m  e  "  :  "  x  "  ,  "  t  y  p  e  "  :  "  u  i  n  t  6  4  "  }  ,  {  "  n  a  m  e  "  :  "  y  "  ,  "  t  y  p  e  "  :  "  i  n  t  3  2  "  }  ]  }  ]  }
HEX:    79 61 72 64 6c 01 00 00 00 b0 02 7b 22 70 72 6f 74 6f 63 6f 6c 22 3a 7b 22 6e 61 6d 65 22 3a 22 4d 79 50 72 6f 74 6f 63 6f 6c 22 2c 22 73 65 71 75 65 6e 63 65 22 3a 5b 7b 22 6e 61 6d 65 22 3a 22 66 6c 6f 61 74 41 72 72 61 79 22 2c 22 74 79 70 65 22 3a 7b 22 61 72 72 61 79 22 3a 7b 22 69 74 65 6d 73 22 3a 22 66 6c 6f 61 74 33 32 22 2c 22 64 69 6d 65 6e 73 69 6f 6e 73 22 3a 5b 7b 22 6c 65 6e 67 74 68 22 3a 32 7d 2c 7b 22 6c 65 6e 67 74 68 22 3a 32 7d 5d 7d 7d 7d 2c 7b 22 6e 61 6d 65 22 3a 22 70 6f 69 6e 74 73 22 2c 22 74 79 70 65 22 3a 7b 22 73 74 72 65 61 6d 22 3a 7b 22 69 74 65 6d 73 22 3a 22 53 61 6e 64 62 6f 78 2e 50 6f 69 6e 74 22 7d 7d 7d 5d 7d 2c 22 74 79 70 65 73 22 3a 5b 7b 22 6e 61 6d 65 22 3a 22 50 6f 69 6e 74 22 2c 22 66 69 65 6c 64 73 22 3a 5b 7b 22 6e 61 6d 65 22 3a 22 78 22 2c 22 74 79 70 65 22 3a 22 75 69 6e 74 36 34 22 7d 2c 7b 22 6e 61 6d 65 22 3a 22 79 22 2c 22 74 79 70 65 22 3a 22 69 6e 74 33 32 22 7d 5d 7d 5d 7d
        mmmmmmmmmmmmmm vvvvvvvvvvv sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss
                                   uuuuu ccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc
        m = magic bytes
        v = version (unsigned int)
        s = string
        u = unsigned varint
        c = char
```

Following the schema string is the `floatArray` value, followed by the `points`
stream, shown below. `floatArray` is made up of four consecutive 32-bit
floating-point values. The `points` stream is made up of blocks with lengths 3, 2,
and 0. The 0-length block indicates the end of the stream, and in this case the
end of the file as well, since there are no more steps in the protocol. Each nonempty
block has `Point`s, each of which is an unsigned varint followed by a signed varint.

```txt
ASCII:  .  .  .  ?  .  .  Y  @  3  3  .  @  .  .  .  @  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0  .  .  m  .
HEX:    9a 99 99 3f 9a 99 59 40 33 33 b3 40 9a 99 f9 40 03 01 04 03 08 05 0c 02 bc 05 c0 0c 80 ea 30 bf ee 6d 00
        aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa ssssssssssssssssssssssssssssssssssssssssssssssssssssssss
        fffffffffff fffffffffff fffffffffff fffffffffff bbbbbbbbbbbbbbbbbbbb bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb bb
                                                        uu ppppp ppppp ppppp uu ppppppppppp ppppppppppppppppp uu
                                                           uu ii uu ii uu ii    uuuuu iiiii uuuuuuuu iiiiiiii

        a = array                                       s = stream
        f = float                                       b = block
                                                        u = unsigned varint
                                                        i = signed varint
                                                        p = point
```
