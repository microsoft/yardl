# This file was generated by the "yardl" tool. DO NOT EDIT.


import dataclasses
import datetime
import enum
import types
import typing
import numpy as np
import numpy.typing as npt
from . import yardl_types as yardl
from . import _dtypes

K = typing.TypeVar('K')
K_NP = typing.TypeVar('K_NP', bound=np.generic)
V = typing.TypeVar('V')
V_NP = typing.TypeVar('V_NP', bound=np.generic)
T = typing.TypeVar('T')
T_NP = typing.TypeVar('T_NP', bound=np.generic)
T1 = typing.TypeVar('T1')
T1_NP = typing.TypeVar('T1_NP', bound=np.generic)
T2 = typing.TypeVar('T2')
T2_NP = typing.TypeVar('T2_NP', bound=np.generic)
T0 = typing.TypeVar('T0')
T0_NP = typing.TypeVar('T0_NP', bound=np.generic)

@dataclasses.dataclass(slots=True, kw_only=True)
class SmallBenchmarkRecord:
    a: yardl.Float64 = 0.0

    b: yardl.Float32 = 0.0

    c: yardl.Float32 = 0.0


@dataclasses.dataclass(slots=True, kw_only=True)
class SimpleEncodingCounters:
    e1: yardl.UInt32 | None = None

    e2: yardl.UInt32 | None = None

    slice: yardl.UInt32 | None = None

    repetition: yardl.UInt32 | None = None


@dataclasses.dataclass(slots=True, kw_only=True)
class SimpleAcquisition:
    flags: yardl.UInt64 = 0

    idx: SimpleEncodingCounters = dataclasses.field(default_factory=SimpleEncodingCounters)

    data: npt.NDArray[np.complex64] = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=np.dtype(np.complex64)))

    trajectory: npt.NDArray[np.float32] = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=np.dtype(np.float32)))


@dataclasses.dataclass(slots=True, kw_only=True)
class SimpleRecord:
    x: yardl.Int32 = 0

    y: yardl.Int32 = 0

    z: yardl.Int32 = 0


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithPrimitives:
    bool_field: yardl.Bool = False

    int_8_field: yardl.Int8 = 0

    uint_8_field: yardl.UInt8 = 0

    int_16_field: yardl.Int16 = 0

    uint_16_field: yardl.UInt16 = 0

    int_32_field: yardl.Int32 = 0

    uint_32_field: yardl.UInt32 = 0

    int_64_field: yardl.Int64 = 0

    uint_64_field: yardl.UInt64 = 0

    size_field: yardl.Size = 0

    float_32_field: yardl.Float32 = 0.0

    float_64_field: yardl.Float64 = 0.0

    complexfloat_32_field: yardl.ComplexFloat = 0j

    complexfloat_64_field: yardl.ComplexDouble = 0j

    date_field: yardl.Date = dataclasses.field(default_factory=lambda: datetime.date(1970, 1, 1))

    time_field: yardl.Time = dataclasses.field(default_factory=lambda: datetime.time(0, 0, 0))

    datetime_field: yardl.DateTime = dataclasses.field(default_factory=lambda: datetime.datetime(1970, 1, 1, 0, 0, 0))


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithPrimitiveAliases:
    byte_field: yardl.UInt8 = 0

    int_field: yardl.Int32 = 0

    uint_field: yardl.UInt32 = 0

    long_field: yardl.Int64 = 0

    ulong_field: yardl.UInt64 = 0

    float_field: yardl.Float32 = 0.0

    double_field: yardl.Float64 = 0.0

    complexfloat_field: yardl.ComplexFloat = 0j

    complexdouble_field: yardl.ComplexDouble = 0j


@dataclasses.dataclass(slots=True, kw_only=True)
class TupleWithRecords:
    a: SimpleRecord = dataclasses.field(default_factory=SimpleRecord)

    b: SimpleRecord = dataclasses.field(default_factory=SimpleRecord)


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithVectors:
    default_vector: list[yardl.Int32] = dataclasses.field(default_factory=list)

    default_vector_fixed_length: list[yardl.Int32] = dataclasses.field(default_factory=lambda: [0] * 3)

    vector_of_vectors: list[list[yardl.Int32]] = dataclasses.field(default_factory=list)


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithArrays:
    default_array: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((), dtype=np.dtype(np.int32)))

    default_array_with_empty_dimension: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((), dtype=np.dtype(np.int32)))

    rank_1_array: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((0), dtype=np.dtype(np.int32)))

    rank_2_array: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=np.dtype(np.int32)))

    rank_2_array_with_named_dimensions: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=np.dtype(np.int32)))

    rank_2_fixed_array: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((3, 4,), dtype=np.dtype(np.int32)))

    rank_2_fixed_array_with_named_dimensions: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((3, 4,), dtype=np.dtype(np.int32)))

    dynamic_array: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((), dtype=np.dtype(np.int32)))

    array_of_vectors: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((5,), dtype=np.dtype(np.object_)))


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithArraysSimpleSyntax:
    default_array: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((), dtype=np.dtype(np.int32)))

    default_array_with_empty_dimension: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((), dtype=np.dtype(np.int32)))

    rank_1_array: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((0), dtype=np.dtype(np.int32)))

    rank_2_array: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=np.dtype(np.int32)))

    rank_2_array_with_named_dimensions: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=np.dtype(np.int32)))

    rank_2_fixed_array: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((3, 4,), dtype=np.dtype(np.int32)))

    rank_2_fixed_array_with_named_dimensions: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((3, 4,), dtype=np.dtype(np.int32)))

    dynamic_array: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((), dtype=np.dtype(np.int32)))

    array_of_vectors: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((5,), dtype=np.dtype(np.object_)))


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithOptionalFields:
    optional_int: yardl.Int32 | None = None

    optional_int_alternate_syntax: yardl.Int32 | None = None


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithVlens:
    a: list[SimpleRecord] = dataclasses.field(default_factory=list)

    b: yardl.Int32 = 0

    c: yardl.Int32 = 0


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithStrings:
    a: str = ""

    b: str = ""


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithOptionalVector:
    optional_vector: list[yardl.Int32] | None = None


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithFixedVectors:
    fixed_int_vector: list[yardl.Int32] = dataclasses.field(default_factory=lambda: [0] * 5)

    fixed_simple_record_vector: list[SimpleRecord] = dataclasses.field(default_factory=lambda: [SimpleRecord() for _ in range(3)])

    fixed_record_with_vlens_vector: list[RecordWithVlens] = dataclasses.field(default_factory=lambda: [RecordWithVlens() for _ in range(2)])


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithFixedArrays:
    ints: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((2, 3,), dtype=np.dtype(np.int32)))

    fixed_simple_record_array: npt.NDArray[np.void] = dataclasses.field(default_factory=lambda: np.zeros((3, 2,), dtype=get_dtype(SimpleRecord)))

    fixed_record_with_vlens_array: npt.NDArray[np.void] = dataclasses.field(default_factory=lambda: np.zeros((2, 2,), dtype=get_dtype(RecordWithVlens)))


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithNDArrays:
    ints: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=np.dtype(np.int32)))

    fixed_simple_record_array: npt.NDArray[np.void] = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=get_dtype(SimpleRecord)))

    fixed_record_with_vlens_array: npt.NDArray[np.void] = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=get_dtype(RecordWithVlens)))


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithNDArraysSingleDimension:
    ints: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((0), dtype=np.dtype(np.int32)))

    fixed_simple_record_array: npt.NDArray[np.void] = dataclasses.field(default_factory=lambda: np.zeros((0), dtype=get_dtype(SimpleRecord)))

    fixed_record_with_vlens_array: npt.NDArray[np.void] = dataclasses.field(default_factory=lambda: np.zeros((0), dtype=get_dtype(RecordWithVlens)))


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithDynamicNDArrays:
    ints: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((), dtype=np.dtype(np.int32)))

    fixed_simple_record_array: npt.NDArray[np.void] = dataclasses.field(default_factory=lambda: np.zeros((), dtype=get_dtype(SimpleRecord)))

    fixed_record_with_vlens_array: npt.NDArray[np.void] = dataclasses.field(default_factory=lambda: np.zeros((), dtype=get_dtype(RecordWithVlens)))


NamedFixedNDArray = npt.NDArray[np.int32]

NamedNDArray = npt.NDArray[np.int32]

AliasedMap = dict[K, V]

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithUnions:
    null_or_int_or_string: (
        tuple[typing.Literal["int32"], yardl.Int32]
        | tuple[typing.Literal["string"], str]
        | None
    ) = None


class Fruits(enum.Enum):
    APPLE = 0
    BANANA = 1
    PEAR = 2

class UInt64Enum(enum.Enum):
    A = 9223372036854775808

class Int64Enum(enum.Enum):
    B = -4611686018427387904

class SizeBasedEnum(enum.Enum):
    A = 0
    B = 1
    C = 2

class DaysOfWeek(enum.Flag, boundary=enum.KEEP):
    MONDAY = 1
    TUESDAY = 2
    WEDNESDAY = 4
    THURSDAY = 8
    FRIDAY = 16
    SATURDAY = 32
    SUNDAY = 64

class TextFormat(enum.Flag, boundary=enum.KEEP):
    REGULAR = 0
    BOLD = 1
    ITALIC = 2
    UNDERLINE = 4
    STRIKETHROUGH = 8

Image = npt.NDArray[T_NP]

@dataclasses.dataclass(slots=True, kw_only=True)
class GenericRecord(typing.Generic[T1, T2, T2_NP]):
    scalar_1: T1
    scalar_2: T2
    vector_1: list[T1]
    image_2: Image[T2_NP]

@dataclasses.dataclass(slots=True, kw_only=True)
class MyTuple(typing.Generic[T1, T2]):
    v1: T1
    v2: T2

AliasedString = str

AliasedEnum = Fruits

AliasedOpenGeneric = MyTuple[T1, T2]

AliasedClosedGeneric = MyTuple[AliasedString, AliasedEnum]

AliasedOptional = yardl.Int32 | None

AliasedGenericOptional = T | None

AliasedGenericUnion2 = (
    tuple[typing.Literal["T1"], T1]
    | tuple[typing.Literal["T2"], T2]
)

AliasedGenericVector = list[T]

AliasedGenericFixedVector = list[T]

AliasedIntOrSimpleRecord = (
    tuple[typing.Literal["int32"], yardl.Int32]
    | tuple[typing.Literal["SimpleRecord"], SimpleRecord]
)

AliasedNullableIntSimpleRecord = (
    tuple[typing.Literal["int32"], yardl.Int32]
    | tuple[typing.Literal["SimpleRecord"], SimpleRecord]
    | None
)

@dataclasses.dataclass(slots=True, kw_only=True)
class GenericRecordWithComputedFields(typing.Generic[T0, T1]):
    f1: (
        tuple[typing.Literal["T0"], T0]
        | tuple[typing.Literal["T1"], T1]
    )
    def type_index(self) -> yardl.UInt8:
        _var0 = self.f1
        if _var0[0] == "T0":
            return 0
        if _var0[0] == "T1":
            return 1
        raise RuntimeError("Unexpected union case")


@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithComputedFields:
    array_field: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=np.dtype(np.int32)))

    array_field_map_dimensions: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=np.dtype(np.int32)))

    dynamic_array_field: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((), dtype=np.dtype(np.int32)))

    fixed_array_field: npt.NDArray[np.int32] = dataclasses.field(default_factory=lambda: np.zeros((3, 4,), dtype=np.dtype(np.int32)))

    int_field: yardl.Int32 = 0

    string_field: str = ""

    tuple_field: MyTuple[yardl.Int32, yardl.Int32] = dataclasses.field(default_factory=lambda: MyTuple(v1=0, v2=0))

    vector_field: list[yardl.Int32] = dataclasses.field(default_factory=list)

    vector_of_vectors_field: list[list[yardl.Int32]] = dataclasses.field(default_factory=list)

    fixed_vector_field: list[yardl.Int32] = dataclasses.field(default_factory=lambda: [0] * 3)

    optional_named_array: NamedNDArray | None = None

    int_float_union: (
        tuple[typing.Literal["int32"], yardl.Int32]
        | tuple[typing.Literal["float32"], yardl.Float32]
    ) = ("int32", 0)

    nullable_int_float_union: (
        tuple[typing.Literal["int32"], yardl.Int32]
        | tuple[typing.Literal["float32"], yardl.Float32]
        | None
    ) = None

    union_with_nested_generic_union: (
        tuple[typing.Literal["int32"], yardl.Int32]
        | tuple[typing.Literal["GenericRecordWithComputedFields<string, float32>"], GenericRecordWithComputedFields[str, yardl.Float32]]
    ) = ("int32", 0)

    map_field: dict[str, str] = dataclasses.field(default_factory=dict)

    def int_literal(self) -> yardl.UInt8:
        return 42

    def large_negative_int_64_literal(self) -> yardl.Int64:
        return -4611686018427387904

    def large_u_int_64_literal(self) -> yardl.UInt64:
        return 9223372036854775808

    def string_literal(self) -> str:
        return "hello"

    def string_literal_2(self) -> str:
        return "hello"

    def string_literal_3(self) -> str:
        return "hello"

    def string_literal_4(self) -> str:
        return "hello"

    def access_other_computed_field(self) -> yardl.Int32:
        return self.int_field

    def access_int_field(self) -> yardl.Int32:
        return self.int_field

    def access_string_field(self) -> str:
        return self.string_field

    def access_tuple_field(self) -> MyTuple[yardl.Int32, yardl.Int32]:
        return self.tuple_field

    def access_nested_tuple_field(self) -> yardl.Int32:
        return self.tuple_field.v2

    def access_array_field(self) -> npt.NDArray[np.int32]:
        return self.array_field

    def access_array_field_element(self) -> yardl.Int32:
        return typing.cast(yardl.Int32, self.array_field[0, 1])

    def access_array_field_element_by_name(self) -> yardl.Int32:
        return typing.cast(yardl.Int32, self.array_field[0, 1])

    def access_vector_field(self) -> list[yardl.Int32]:
        return self.vector_field

    def access_vector_field_element(self) -> yardl.Int32:
        return self.vector_field[1]

    def access_vector_of_vectors_field(self) -> yardl.Int32:
        return self.vector_of_vectors_field[1][2]

    def array_size(self) -> yardl.Size:
        return self.array_field.size

    def array_x_size(self) -> yardl.Size:
        return self.array_field.shape[0]

    def array_y_size(self) -> yardl.Size:
        return self.array_field.shape[1]

    def array_0_size(self) -> yardl.Size:
        return self.array_field.shape[0]

    def array_1_size(self) -> yardl.Size:
        return self.array_field.shape[1]

    def array_size_from_int_field(self) -> yardl.Size:
        return self.array_field.shape[self.int_field]

    def array_size_from_string_field(self) -> yardl.Size:
        def _helper_0(dim_name: str) -> int:
            if dim_name == "x":
                return 0
            if dim_name == "y":
                return 1
            raise KeyError(f"Unknown dimension name: '{dim_name}'")

        return self.array_field.shape[_helper_0(self.string_field)]

    def array_size_from_nested_int_field(self) -> yardl.Size:
        return self.array_field.shape[self.tuple_field.v1]

    def array_field_map_dimensions_x_size(self) -> yardl.Size:
        return self.array_field_map_dimensions.shape[0]

    def fixed_array_size(self) -> yardl.Size:
        return 12

    def fixed_array_x_size(self) -> yardl.Size:
        return 3

    def fixed_array_0_size(self) -> yardl.Size:
        return 3

    def vector_size(self) -> yardl.Size:
        return len(self.vector_field)

    def fixed_vector_size(self) -> yardl.Size:
        return 3

    def array_dimension_x_index(self) -> yardl.Size:
        return 0

    def array_dimension_y_index(self) -> yardl.Size:
        return 1

    def array_dimension_index_from_string_field(self) -> yardl.Size:
        def _helper_0(dim_name: str) -> int:
            if dim_name == "x":
                return 0
            if dim_name == "y":
                return 1
            raise KeyError(f"Unknown dimension name: '{dim_name}'")

        return _helper_0(self.string_field)

    def array_dimension_count(self) -> yardl.Size:
        return 2

    def dynamic_array_dimension_count(self) -> yardl.Size:
        return self.dynamic_array_field.ndim

    def access_map(self) -> dict[str, str]:
        return self.map_field

    def map_size(self) -> yardl.Size:
        return len(self.map_field)

    def access_map_entry(self) -> str:
        return self.map_field["hello"]

    def string_computed_field(self) -> str:
        return "hello"

    def access_map_entry_with_computed_field(self) -> str:
        return self.map_field[self.string_computed_field()]

    def access_map_entry_with_computed_field_nested(self) -> str:
        return self.map_field[self.map_field[self.string_computed_field()]]

    def access_missing_map_entry(self) -> str:
        return self.map_field["missing"]

    def optional_named_array_length(self) -> yardl.Size:
        _var0 = self.optional_named_array
        if _var0 is not None:
            arr = _var0
            return arr.size
        return 0

    def optional_named_array_length_with_discard(self) -> yardl.Size:
        _var0 = self.optional_named_array
        if _var0 is not None:
            arr = _var0
            return arr.size
        return 0

    def int_float_union_as_float(self) -> yardl.Float32:
        _var0 = self.int_float_union
        if _var0[0] == "int32":
            i_foo = _var0[1]
            return float(i_foo)
        if _var0[0] == "float32":
            f = _var0[1]
            return f
        raise RuntimeError("Unexpected union case")

    def nullable_int_float_union_string(self) -> str:
        _var0 = self.nullable_int_float_union
        if _var0 is None:
            return "null"
        if _var0[0] == "int32":
            return "int"
        return "float"
        raise RuntimeError("Unexpected union case")

    def nested_switch(self) -> yardl.Int16:
        _var0 = self.union_with_nested_generic_union
        if _var0[0] == "int32":
            return -1
        if _var0[0] == "GenericRecordWithComputedFields<string, float32>":
            rec = _var0[1]
            _var1 = rec.f1
            if _var1[0] == "T1":
                return int(20)
            if _var1[0] == "T0":
                return int(10)
            raise RuntimeError("Unexpected union case")
        raise RuntimeError("Unexpected union case")

    def use_nested_computed_field(self) -> yardl.Int16:
        _var0 = self.union_with_nested_generic_union
        if _var0[0] == "int32":
            return -1
        if _var0[0] == "GenericRecordWithComputedFields<string, float32>":
            rec = _var0[1]
            return int(rec.type_index())
        raise RuntimeError("Unexpected union case")

    def switch_over_sigle_value(self) -> yardl.Int32:
        _var0 = self.int_field
        ii = _var0
        return ii


ArrayWithKeywordDimensionNames = npt.NDArray[np.int32]

class EnumWithKeywordSymbols(enum.Enum):
    TRY = 2
    CATCH = 1

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithKeywordFields:
    int_: str = ""

    sizeof: ArrayWithKeywordDimensionNames = dataclasses.field(default_factory=lambda: np.zeros((0,0), dtype=np.dtype(np.int32)))

    if_: EnumWithKeywordSymbols

    def float_(self) -> str:
        return self.int_

    def double(self) -> str:
        return self.float_()

    def return_(self) -> yardl.Int32:
        return self.sizeof[1, 2]


def _mk_get_dtype():
    dtype_map: dict[type | types.GenericAlias, np.dtype[typing.Any] | typing.Callable[[tuple[type, ...]], np.dtype[typing.Any]]] = {}
    get_dtype = _dtypes.make_get_dtype_func(dtype_map)

    dtype_map[SmallBenchmarkRecord] = np.dtype([('a', np.dtype(np.float64)), ('b', np.dtype(np.float32)), ('c', np.dtype(np.float32))], align=True)
    dtype_map[SimpleEncodingCounters] = np.dtype([('e1', np.dtype([('has_value', np.dtype(np.bool_)), ('value', np.dtype(np.uint32))], align=True)), ('e2', np.dtype([('has_value', np.dtype(np.bool_)), ('value', np.dtype(np.uint32))], align=True)), ('slice', np.dtype([('has_value', np.dtype(np.bool_)), ('value', np.dtype(np.uint32))], align=True)), ('repetition', np.dtype([('has_value', np.dtype(np.bool_)), ('value', np.dtype(np.uint32))], align=True))], align=True)
    dtype_map[SimpleAcquisition] = np.dtype([('flags', np.dtype(np.uint64)), ('idx', get_dtype(SimpleEncodingCounters)), ('data', np.dtype(np.object_)), ('trajectory', np.dtype(np.object_))], align=True)
    dtype_map[SimpleRecord] = np.dtype([('x', np.dtype(np.int32)), ('y', np.dtype(np.int32)), ('z', np.dtype(np.int32))], align=True)
    dtype_map[RecordWithPrimitives] = np.dtype([('boolField', np.dtype(np.bool_)), ('int8Field', np.dtype(np.int8)), ('uint8Field', np.dtype(np.uint8)), ('int16Field', np.dtype(np.int16)), ('uint16Field', np.dtype(np.uint16)), ('int32Field', np.dtype(np.int32)), ('uint32Field', np.dtype(np.uint32)), ('int64Field', np.dtype(np.int64)), ('uint64Field', np.dtype(np.uint64)), ('sizeField', np.dtype(np.uint64)), ('float32Field', np.dtype(np.float32)), ('float64Field', np.dtype(np.float64)), ('complexfloat32Field', np.dtype(np.complex64)), ('complexfloat64Field', np.dtype(np.complex128)), ('dateField', np.dtype(np.datetime64)), ('timeField', np.dtype(np.timedelta64)), ('datetimeField', np.dtype(np.datetime64))], align=True)
    dtype_map[RecordWithPrimitiveAliases] = np.dtype([('byteField', np.dtype(np.uint8)), ('intField', np.dtype(np.int32)), ('uintField', np.dtype(np.uint32)), ('longField', np.dtype(np.int64)), ('ulongField', np.dtype(np.uint64)), ('floatField', np.dtype(np.float32)), ('doubleField', np.dtype(np.float64)), ('complexfloatField', np.dtype(np.complex64)), ('complexdoubleField', np.dtype(np.complex128))], align=True)
    dtype_map[TupleWithRecords] = np.dtype([('a', get_dtype(SimpleRecord)), ('b', get_dtype(SimpleRecord))], align=True)
    dtype_map[RecordWithVectors] = np.dtype([('defaultVector', np.dtype(np.object_)), ('defaultVectorFixedLength', np.dtype(np.int32), (3,)), ('vectorOfVectors', np.dtype(np.object_))], align=True)
    dtype_map[RecordWithArrays] = np.dtype([('defaultArray', np.dtype(np.object_)), ('defaultArrayWithEmptyDimension', np.dtype(np.object_)), ('rank1Array', np.dtype(np.object_)), ('rank2Array', np.dtype(np.object_)), ('rank2ArrayWithNamedDimensions', np.dtype(np.object_)), ('rank2FixedArray', np.dtype(np.int32), (3, 4,)), ('rank2FixedArrayWithNamedDimensions', np.dtype(np.int32), (3, 4,)), ('dynamicArray', np.dtype(np.object_)), ('arrayOfVectors', np.dtype(np.object_), (5,))], align=True)
    dtype_map[RecordWithArraysSimpleSyntax] = np.dtype([('defaultArray', np.dtype(np.object_)), ('defaultArrayWithEmptyDimension', np.dtype(np.object_)), ('rank1Array', np.dtype(np.object_)), ('rank2Array', np.dtype(np.object_)), ('rank2ArrayWithNamedDimensions', np.dtype(np.object_)), ('rank2FixedArray', np.dtype(np.int32), (3, 4,)), ('rank2FixedArrayWithNamedDimensions', np.dtype(np.int32), (3, 4,)), ('dynamicArray', np.dtype(np.object_)), ('arrayOfVectors', np.dtype(np.object_), (5,))], align=True)
    dtype_map[RecordWithOptionalFields] = np.dtype([('optionalInt', np.dtype([('has_value', np.dtype(np.bool_)), ('value', np.dtype(np.int32))], align=True)), ('optionalIntAlternateSyntax', np.dtype([('has_value', np.dtype(np.bool_)), ('value', np.dtype(np.int32))], align=True))], align=True)
    dtype_map[RecordWithVlens] = np.dtype([('a', np.dtype(np.object_)), ('b', np.dtype(np.int32)), ('c', np.dtype(np.int32))], align=True)
    dtype_map[RecordWithStrings] = np.dtype([('a', np.dtype(np.object_)), ('b', np.dtype(np.object_))], align=True)
    dtype_map[RecordWithOptionalVector] = np.dtype([('optionalVector', np.dtype([('has_value', np.dtype(np.bool_)), ('value', np.dtype(np.object_))], align=True))], align=True)
    dtype_map[RecordWithFixedVectors] = np.dtype([('fixedIntVector', np.dtype(np.int32), (5,)), ('fixedSimpleRecordVector', get_dtype(SimpleRecord), (3,)), ('fixedRecordWithVlensVector', get_dtype(RecordWithVlens), (2,))], align=True)
    dtype_map[RecordWithFixedArrays] = np.dtype([('ints', np.dtype(np.int32), (2, 3,)), ('fixedSimpleRecordArray', get_dtype(SimpleRecord), (3, 2,)), ('fixedRecordWithVlensArray', get_dtype(RecordWithVlens), (2, 2,))], align=True)
    dtype_map[RecordWithNDArrays] = np.dtype([('ints', np.dtype(np.object_)), ('fixedSimpleRecordArray', np.dtype(np.object_)), ('fixedRecordWithVlensArray', np.dtype(np.object_))], align=True)
    dtype_map[RecordWithNDArraysSingleDimension] = np.dtype([('ints', np.dtype(np.object_)), ('fixedSimpleRecordArray', np.dtype(np.object_)), ('fixedRecordWithVlensArray', np.dtype(np.object_))], align=True)
    dtype_map[RecordWithDynamicNDArrays] = np.dtype([('ints', np.dtype(np.object_)), ('fixedSimpleRecordArray', np.dtype(np.object_)), ('fixedRecordWithVlensArray', np.dtype(np.object_))], align=True)
    dtype_map[NamedFixedNDArray] = np.dtype(np.int32)
    dtype_map[NamedNDArray] = np.dtype(np.object_)
    dtype_map[AliasedMap] = lambda type_args: np.dtype(np.object_)
    dtype_map[RecordWithUnions] = np.dtype([('nullOrIntOrString', np.dtype(np.object_))], align=True)
    dtype_map[Fruits] = np.dtype(np.int32)
    dtype_map[UInt64Enum] = np.dtype(np.uint64)
    dtype_map[Int64Enum] = np.dtype(np.int64)
    dtype_map[SizeBasedEnum] = np.dtype(np.uint64)
    dtype_map[DaysOfWeek] = np.dtype(np.int32)
    dtype_map[TextFormat] = np.dtype(np.uint64)
    dtype_map[Image] = lambda type_args: np.dtype(np.object_)
    dtype_map[GenericRecord] = lambda type_args: np.dtype([('scalar1', get_dtype(type_args[0])), ('scalar2', get_dtype(type_args[1])), ('vector1', np.dtype(np.object_)), ('image2', get_dtype(types.GenericAlias(Image, (type_args[1],))))], align=True)
    dtype_map[MyTuple] = lambda type_args: np.dtype([('v1', get_dtype(type_args[0])), ('v2', get_dtype(type_args[1]))], align=True)
    dtype_map[AliasedString] = np.dtype(np.object_)
    dtype_map[AliasedEnum] = get_dtype(Fruits)
    dtype_map[AliasedOpenGeneric] = lambda type_args: get_dtype(types.GenericAlias(MyTuple, (type_args[0], type_args[1],)))
    dtype_map[AliasedClosedGeneric] = get_dtype(types.GenericAlias(MyTuple, (AliasedString, AliasedEnum,)))
    dtype_map[AliasedOptional] = np.dtype([('has_value', np.dtype(np.bool_)), ('value', np.dtype(np.int32))], align=True)
    dtype_map[AliasedGenericOptional] = lambda type_args: np.dtype([('has_value', np.dtype(np.bool_)), ('value', get_dtype(type_args[0]))], align=True)
    dtype_map[AliasedGenericUnion2] = lambda type_args: np.dtype(np.object_)
    dtype_map[AliasedGenericVector] = lambda type_args: np.dtype(np.object_)
    dtype_map[AliasedGenericFixedVector] = lambda type_args: get_dtype(type_args[0])
    dtype_map[AliasedIntOrSimpleRecord] = np.dtype(np.object_)
    dtype_map[AliasedNullableIntSimpleRecord] = np.dtype(np.object_)
    dtype_map[GenericRecordWithComputedFields] = lambda type_args: np.dtype([('f1', np.dtype(np.object_))], align=True)
    dtype_map[RecordWithComputedFields] = np.dtype([('arrayField', np.dtype(np.object_)), ('arrayFieldMapDimensions', np.dtype(np.object_)), ('dynamicArrayField', np.dtype(np.object_)), ('fixedArrayField', np.dtype(np.int32), (3, 4,)), ('intField', np.dtype(np.int32)), ('stringField', np.dtype(np.object_)), ('tupleField', get_dtype(types.GenericAlias(MyTuple, (yardl.Int32, yardl.Int32,)))), ('vectorField', np.dtype(np.object_)), ('vectorOfVectorsField', np.dtype(np.object_)), ('fixedVectorField', np.dtype(np.int32), (3,)), ('optionalNamedArray', np.dtype([('has_value', np.dtype(np.bool_)), ('value', get_dtype(NamedNDArray))], align=True)), ('intFloatUnion', np.dtype(np.object_)), ('nullableIntFloatUnion', np.dtype(np.object_)), ('unionWithNestedGenericUnion', np.dtype(np.object_)), ('mapField', np.dtype(np.object_))], align=True)
    dtype_map[ArrayWithKeywordDimensionNames] = np.dtype(np.object_)
    dtype_map[EnumWithKeywordSymbols] = np.dtype(np.int32)
    dtype_map[RecordWithKeywordFields] = np.dtype([('int', np.dtype(np.object_)), ('sizeof', get_dtype(ArrayWithKeywordDimensionNames)), ('if', get_dtype(EnumWithKeywordSymbols))], align=True)

    return get_dtype

get_dtype = _mk_get_dtype()

