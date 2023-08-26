# NDJSON Encoding Reference

::: warning
We might make breaking changes to this format before V1.
:::

The NDJSON format is meant for easy debugging or interoperability in scenarios where performance is less of a concern.

Here is an example protocol illustrates what what the format looks like. We will define the following model:

```yaml
MyRecord: !record
  fields:
    x: int
    y: int
    z: int?

MyEnum: !enum
  values:
    - a
    - b
    - c

MyFlags: !flags
  values:
    - a
    - b
    - c

HelloNDJson: !protocol
  sequence:
    anIntStream: !stream
      items: int
    aBoolean: bool
    aString: string
    aComplex: complexdouble
    aDate: date
    aTime: time
    aDateTime: datetime

    anEnum: MyEnum
    someFlags: MyFlags

    anOptionalIntThatIsNotSet: int?
    anOptionalIntThatIsSet: int?

    aRecordWithOptionalNotSet: MyRecord
    aRecordWithOptionalSet: MyRecord

    aVector: int*
    aDynamicArray: int[]
    aFixedArray: int[2,3]

    aMapWithAStringKey: string->int
    aMapWithAnIntKey: int->int

    aUnionWithSimpleRepresentation: [int, bool]
    aUnionRequiringTag: [string, MyEnum]
  ```

  And then write some data to an NDJSON protocol writer. The output looks like this:

```json
{"yardl":{"version":1,"schema":{"protocol":{"name":"HelloNDJson","sequence":[{"name":"anIntStream","type":{"stream":{"items":"int32"}}},{"name":"aBoolean","type":"bool"},{"name":"aString","type":"string"},{"name":"aComplex","type":"complexfloat64"},{"name":"aDate","type":"date"},{"name":"aTime","type":"time"},{"name":"aDateTime","type":"datetime"},{"name":"anEnum","type":"Sandbox.MyEnum"},{"name":"someFlags","type":"Sandbox.MyFlags"},{"name":"anOptionalIntThatIsNotSet","type":[null,"int32"]},{"name":"anOptionalIntThatIsSet","type":[null,"int32"]},{"name":"aRecordWithOptionalNotSet","type":"Sandbox.MyRecord"},{"name":"aRecordWithOptionalSet","type":"Sandbox.MyRecord"},{"name":"aVector","type":{"vector":{"items":"int32"}}},{"name":"aDynamicArray","type":{"array":{"items":"int32"}}},{"name":"aFixedArray","type":{"array":{"items":"int32","dimensions":[{"length":2},{"length":3}]}}},{"name":"aMapWithAStringKey","type":{"map":{"keys":"string","values":"int32"}}},{"name":"aMapWithAnIntKey","type":{"map":{"keys":"int32","values":"int32"}}},{"name":"aUnionWithSimpleRepresentation","type":[{"label":"int32","type":"int32"},{"label":"bool","type":"bool"}]},{"name":"aUnionRequiringTag","type":[{"label":"string","type":"string"},{"label":"MyEnum","type":"Sandbox.MyEnum"}]}]},"types":[{"name":"MyEnum","values":[{"symbol":"a","value":0},{"symbol":"b","value":1},{"symbol":"c","value":2}]},{"name":"MyFlags","values":[{"symbol":"a","value":1},{"symbol":"b","value":2},{"symbol":"c","value":4}]},{"name":"MyRecord","fields":[{"name":"x","type":"int32"},{"name":"y","type":"int32"},{"name":"z","type":[null,"int32"]}]}]}}}
{"anIntStream":1}
{"anIntStream":2}
{"anIntStream":3}
{"aBoolean":true}
{"aString":"hello"}
{"aComplex":[1.0,2.0]}
{"aDate":"2020-01-17"}
{"aTime":"10:50:25.777888999"}
{"aDateTime":"2023-05-30T18:36:56.708792349Z"}
{"anEnum":"a"}
{"someFlags":["a","b"]}
{"anOptionalIntThatIsNotSet":null}
{"anOptionalIntThatIsSet":42}
{"aRecordWithOptionalNotSet":{"x":1,"y":2}}
{"aRecordWithOptionalSet":{"x":1,"y":2,"z":3}}
{"aVector":[1,2,3]}
{"aDynamicArray":{"shape":[2,3],"data":[1,2,3,4,5,6]}}
{"aFixedArray":[1,2,3,4,5,6]}
{"aMapWithAStringKey":{"b":2,"a":1}}
{"aMapWithAnIntKey":[[2,2],[1,1]]}
{"aUnionWithSimpleRepresentation":22}
{"aUnionRequiringTag":{"string":"a"}}
```

Each line is a JSON document. The first contains the encoding format version
along with the [protocol schema](#protocol-schema-json-reference). Each
subsequent line is a JSON object with a single field. The name of the field is the
protocol step name, and its value is payload value. Note that in the case of
streams, there can be many contiguous lines with same protocol step name.

Datatypes are serialized as follows:

- Booleans are serialized as JSON booleans.
- Integers and floating-point numbers are serialized as JSON numbers.
- Complex numbers are serialized as a JSON array of the real component followed
  by the imaginary component.
- Strings are serialized as JSON strings.
- Enums are serialized as their symbolic string value or as the integer value if
  the value is outside of the defined values.
- Flags are serialized as an array of the symbolic string values that are set.
  If the value is outside of the defined values, the underlying integer value is
  written instead of the array.
- Dates, Times, and DateTimes are formatted as strings. Dates are formatted as
  `YYYY-MM-DD`, times as `HH:mm:SS.FFFFFFFFF`, and datetimes as
  `YYYY-MM-DDTHH:mm:ss:FFFFFFFFFZ`
- Records are serialized as JSON objects, with a JSON field for each record
  fields. Fields are skipped if they are options or unions with `null` as a
  option and are the value is `null`.
- Vectors are serialized as JSON arrays.
- Fixed multidimensional arrays are serialized as a flattened JSON array with
  the values written in row-major order.
- Multidimensional arrays where the dimension sizes are not fixed are serialized
  as a JSON object with two fields: `shape` and `data`. `shape` is an array of
  dimension sizes and `data` is an array of the values in row-major order.
- Maps where the key is a string are written as a JSON object.
- Other maps are written as an array of arrays, with each inner array holding
  the key and value of each entry.
- Optional values are serialized as the inner value if present, and `null` if
  not.
- If each type case of a union serializes to a distict JSON datatype (number,
  string, boolean, array, object), the inner value is serialized directy. For
  example the union `[int, bool]` can be written simply as `29` or `true`. In
  other cases, the value is serialized a JSON object with a single field. The
  field's name is the label of the type case set (see [here](#unions-1) for a
  description of the labels) and the field's value is the JSON serialization of
  the inner value. For example, a value of the union `[float, double]` could be
  written as `{"float32": 29.9}` or `{"float64": 882.2}`.
