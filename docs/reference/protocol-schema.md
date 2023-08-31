# Protocol Schema JSON

A protocol's schema is embedded in a JSON format in the binary, HDF5, and NDJSON
encodings. This JSON format is informally described here.

::: warning
We might make breaking changes to this format before V1.
:::

The JSON schema is meant to be provide enough information for deserializers to
understand the schema at runtime, and therefore does not contain the comments
that may be in the Yardl, nor does it contain computed fields, since those are
not needed for deserialization.

## References to Primitive Types

Primitive type references are represented by name as a JSON string, e.g.
`"int32"` or `"string"`

## References to Top-Level Types

References to top-level types (records, enums, and aliases) are represented by
their namespaced name as a JSON string, e.g. `"MyNamespace.MyRecord"` or
`"MyNamespace.MyEnum"`.

## Unions

Unions are represented as a JSON array:

```JSON
[
  {
    "tag": "int32",
    "type": "int32"
  },
  {
    "tag": "float32",
    "type": "float32"
  }
]
```

The `tag` field is unique name automatically assigned to each union case,
derived from its type name. The labels are used in the Python codegen and in
HDF5 encoding. The tag can be explicitly set using the expanded `!union` syntax
([C++](/cpp/language#unions)).

If `null` is one of the cases, it is represented by `null` in the JSON as well:

```JSON
[
  null,
  {
    "label": "int32",
    "type": "int32"
  },
  {
    "label": "float32",
    "type": "float32"
  }
]
```

For the special case of an optional type, a label for the non-null case is
omitted and the object is simplified to its `type` value, since the label is not
used.

```JSON
[
  null,
  "int32"
]
```

## Vectors

Vectors have the following representation:

```JSON
{
  "vector": {
    "items": "int32",
    "length": 10
  }
}
```

The `length` field is only present if it is given in the Yardl definition.

## Arrays

A fixed array with dimensions `x` and `y` would look like this:

```JSON
{
  "array": {
    "items": "int32",
    "dimensions": [
      {
        "name": "x",
        "length": 3
      },
      {
        "name": "y",
        "length": 4
      }
    ]
  }
}
```

A non-fixed array with dimensions `x` and `y` would look like the above but
without the `length` field:

```JSON
{
  "array": {
    "items": "int32",
    "dimensions": [
      {
        "name": "x"
      },
      {
        "name": "y"
      }
    ]
  }
}
```

A non-fixed array with two unnamed dimensions would look like this:

```JSON
{
  "array": {
    "items": "int32",
    "dimensions": 2
  }
}
```

And finally, an array with an unknown number of dimensions:

```JSON
{
  "array": {
    "items": "int32"
  }
}
```

## Maps

Maps have the following representation:

```JSON
{
  "map": {
    "keys": "string",
    "values": "int32"
  }
}
```

## Streams

Streams in protocols are represented as:

```JSON
{
  "stream": {
    "items": "int32",
  }
}
```

## Enums

Enums are top-level types and cannot be declared inline.

An example enum that looks like this in Yardl:

```yaml
Animals: !enum
  base: uint8
  values: [cat, dog]
```

Is represented in JSON as:

```JSON
{
  "enum": {
    "name": "Animals",
    "base": "uint8",
    "values": [
      {
        "symbol": "cat",
        "value": 0
      },
      {
        "symbol": "dog",
        "value": 1
      }
    ]
  }
}
```

The `base` field is only present in the JSON if it is specified in the Yardl.

## Flags

Like enums, flags are top-level types and cannot be declared inline. They
represented just like enums except for the "flags" field name:

```JSON
{
  "flags": {
    "name": "TextFormat",
    "values": [
      {
        "symbol": "regular",
        "value": 0
      },
      {
        "symbol": "bold",
        "value": 1
      },
      {
        "symbol": "italic",
        "value": 2
      },
      {
        "symbol": "underline",
        "value": 4
      },
      {
        "symbol": "strikethrough",
        "value": 8
      }
    ]
  }
},
```

## Records

Records are top-level types that cannot be declared inline.

An example generic record:

```yaml
MyTuple<T1, T2>: !record
  fields:
    f1: T1
    f2: T2
```

Would look like this:

```JSON
{
  "record": {
    "name": "MyTuple",
    "typeParameters": [
      "T1",
      "T2"
    ],
    "fields": [
      {
        "name": "f1",
        "type": "T1"
      },
      {
        "name": "f2",
        "type": "T2"
      }
    ]
  }
}
```

Computed fields are omitted from the JSON since they are not used during
deserialization.

## Aliases

A simple type alias:

```yaml
MyString: string
```

Is converted to:

```JSON
{
  "alias": {
    "name": "MyString",
    "type": "string"
  }
}
```

If an alias is generic:

```yaml
MyVector<T> : !vector
  items: T
```

its JSON looks like this:

```JSON
{
  "alias": {
    "name": "MyVector",
    "typeParameters": [
      "T"
    ],
    "type": {
      "vector": {
        "items": "T"
      }
    }
  }
}
```

## Protocols

A protocol that looks like this in Yardl:

```yaml
MyProtocol : !protocol
  sequence:
    a: string
    b: !stream
      items: double
```

is represented in JSON like this:

```JSON
{
  "name": "MyProtocol",
  "sequence": [
    {
      "name": "a",
      "type": "string"
    },
    {
      "name": "b",
      "type": {
        "stream": {
          "items": "float64"
        }
      }
    }
  ]
}
```

## Top-Level Schema

The JSON that is embedded in the binary, HDF5, or NDJSON formats contains the
protocol definition (defined above) and the transitive closure of named types
(records, enums, and aliases) used by the protocol.

```JSON
{
  "protocol": <protocol>,
  "types": [ <type>, ... ]
}
```
