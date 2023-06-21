# This file was generated by the "yardl" tool. DO NOT EDIT.


import dataclasses
import datetime
import enum
import typing
import numpy as np
import numpy.typing as npt
from . import yardl_types as yardl

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
INT16_MAX = typing.TypeVar('INT16_MAX')
INT16_MAX_NP = typing.TypeVar('INT16_MAX_NP', bound=np.generic)

@dataclasses.dataclass(slots=True, kw_only=True)
class SmallBenchmarkRecord:
    a: yardl.Float64
    b: yardl.Float32
    c: yardl.Float32

@dataclasses.dataclass(slots=True, kw_only=True)
class SimpleEncodingCounters:
    e1: yardl.UInt32 | None = None
    e2: yardl.UInt32 | None = None
    slice: yardl.UInt32 | None = None
    repetition: yardl.UInt32 | None = None

@dataclasses.dataclass(slots=True, kw_only=True)
class SimpleAcquisition:
    flags: yardl.UInt64
    idx: SimpleEncodingCounters
    data: npt.NDArray[np.complex64]
    trajectory: npt.NDArray[np.float32]

@dataclasses.dataclass(slots=True, kw_only=True)
class SimpleRecord:
    x: yardl.Int32
    y: yardl.Int32
    z: yardl.Int32

    @staticmethod
    def dtype() -> npt.DTypeLike:
        return np.dtype([('x', np.int32), ('y', np.int32), ('z', np.int32)], align=True)

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithPrimitives:
    bool_field: yardl.Bool
    int_8_field: yardl.Int8
    uint_8_field: yardl.UInt8
    int_16_field: yardl.Int16
    uint_16_field: yardl.UInt16
    int_32_field: yardl.Int32
    uint_32_field: yardl.UInt32
    int_64_field: yardl.Int64
    uint_64_field: yardl.UInt64
    size_field: yardl.Size
    float_32_field: yardl.Float32
    float_64_field: yardl.Float64
    complexfloat_32_field: yardl.ComplexFloat
    complexfloat_64_field: yardl.ComplexDouble
    date_field: yardl.Date
    time_field: yardl.Time
    datetime_field: yardl.DateTime

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithPrimitiveAliases:
    byte_field: yardl.UInt8
    int_field: yardl.Int32
    uint_field: yardl.UInt32
    long_field: yardl.Int64
    ulong_field: yardl.UInt64
    float_field: yardl.Float32
    double_field: yardl.Float64
    complexfloat_field: yardl.ComplexFloat
    complexdouble_field: yardl.ComplexDouble

@dataclasses.dataclass(slots=True, kw_only=True)
class TupleWithRecords:
    a: SimpleRecord
    b: SimpleRecord

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithVectors:
    default_vector: list[yardl.Int32]
    default_vector_fixed_length: list[yardl.Int32]
    vector_of_vectors: list[list[yardl.Int32]]

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithArrays:
    default_array: npt.NDArray[np.int32]
    default_array_with_empty_dimension: npt.NDArray[np.int32]
    rank_1_array: npt.NDArray[np.int32]
    rank_2_array: npt.NDArray[np.int32]
    rank_2_array_with_named_dimensions: npt.NDArray[np.int32]
    rank_2_fixed_array: npt.NDArray[np.int32]
    rank_2_fixed_array_with_named_dimensions: npt.NDArray[np.int32]
    dynamic_array: npt.NDArray[np.int32]
    array_of_vectors: npt.NDArray[np.int32, (4,)]

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithArraysSimpleSyntax:
    default_array: npt.NDArray[np.int32]
    default_array_with_empty_dimension: npt.NDArray[np.int32]
    rank_1_array: npt.NDArray[np.int32]
    rank_2_array: npt.NDArray[np.int32]
    rank_2_array_with_named_dimensions: npt.NDArray[np.int32]
    rank_2_fixed_array: npt.NDArray[np.int32]
    rank_2_fixed_array_with_named_dimensions: npt.NDArray[np.int32]
    dynamic_array: npt.NDArray[np.int32]
    array_of_vectors: npt.NDArray[np.int32, (4,)]

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithOptionalFields:
    optional_int: yardl.Int32 | None = None
    optional_int_alternate_syntax: yardl.Int32 | None = None

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithVlens:
    a: list[SimpleRecord]
    b: yardl.Int32
    c: yardl.Int32

    @staticmethod
    def dtype() -> npt.DTypeLike:
        return np.dtype([('a', np.object_), ('b', np.int32), ('c', np.int32)], align=True)

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithStrings:
    a: str
    b: str

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithOptionalVector:
    optional_vector: list[yardl.Int32] | None = None

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithFixedVectors:
    fixed_int_vector: list[yardl.Int32]
    fixed_simple_record_vector: list[SimpleRecord]
    fixed_record_with_vlens_vector: list[RecordWithVlens]

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithFixedArrays:
    ints: npt.NDArray[np.int32]
    fixed_simple_record_array: npt.NDArray[np.void]
    fixed_record_with_vlens_array: npt.NDArray[np.void]

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithNDArrays:
    ints: npt.NDArray[np.int32]
    fixed_simple_record_array: npt.NDArray[np.void]
    fixed_record_with_vlens_array: npt.NDArray[np.void]

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithNDArraysSingleDimension:
    ints: npt.NDArray[np.int32]
    fixed_simple_record_array: npt.NDArray[np.void]
    fixed_record_with_vlens_array: npt.NDArray[np.void]

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithDynamicNDArrays:
    ints: npt.NDArray[np.int32]
    fixed_simple_record_array: npt.NDArray[np.void]
    fixed_record_with_vlens_array: npt.NDArray[np.void]

NamedFixedNDArray = npt.NDArray[np.int32]

NamedNDArray = npt.NDArray[np.int32]

AliasedMap = dict[K, V]

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithUnions:
    null_or_int_or_string: yardl.Int32 | str | None = None

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

class DaysOfWeek(enum.Flag):
    MONDAY = 1
    TUESDAY = 2
    WEDNESDAY = 4
    THURSDAY = 8
    FRIDAY = 16
    SATURDAY = 32
    SUNDAY = 64

class TextFormat(enum.Flag):
    REGULAR = 0
    BOLD = 1
    ITALIC = 2
    UNDERLINE = 4
    STRIKETHROUGH = 8

Image = npt.NDArray[T_NP]

@dataclasses.dataclass(slots=True, kw_only=True)
class GenericRecord(typing.Generic[T1, T2]):
    scalar_1: T1
    scalar_2: T2
    vector_1: list[T1]
    image_2: Image[T2]

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

AliasedGenericUnion2 = T1 | T2

AliasedGenericVector = list[T]

AliasedGenericFixedVector = list[T]

AliasedIntOrSimpleRecord = yardl.Int32 | SimpleRecord

AliasedNullableIntSimpleRecord = yardl.Int32 | SimpleRecord | None

@dataclasses.dataclass(slots=True, kw_only=True)
class GenericRecordWithComputedFields(typing.Generic[T0, T1]):
    f1: T0 | T1

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithComputedFields:
    array_field: npt.NDArray[np.int32]
    array_field_map_dimensions: npt.NDArray[np.int32]
    dynamic_array_field: npt.NDArray[np.int32]
    fixed_array_field: npt.NDArray[np.int32]
    int_field: yardl.Int32
    string_field: str
    tuple_field: MyTuple[yardl.Int32, yardl.Int32]
    vector_field: list[yardl.Int32]
    vector_of_vectors_field: list[list[yardl.Int32]]
    fixed_vector_field: list[yardl.Int32]
    optional_named_array: NamedNDArray | None = None
    int_float_union: yardl.Int32 | yardl.Float32
    nullable_int_float_union: yardl.Int32 | yardl.Float32 | None = None
    union_with_nested_generic_union: yardl.Int32 | GenericRecordWithComputedFields[str, yardl.Float32]
    map_field: dict[str, str]

ArrayWithKeywordDimensionNames = npt.NDArray[INT16_MAX_NP]

class EnumWithKeywordSymbols(enum.Enum):
    TRY = 2
    CATCH = 1

@dataclasses.dataclass(slots=True, kw_only=True)
class RecordWithKeywordFields:
    """BEGIN delibrately using C++ keywords and macros as identitiers"""
    int_: str
    sizeof: ArrayWithKeywordDimensionNames[yardl.Int32]
    if_: EnumWithKeywordSymbols

