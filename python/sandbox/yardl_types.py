from enum import Flag, auto
from typing import Any
import numpy as np
import datetime
from dataclasses import dataclass

Bool = bool | np.bool_
Int8 = int | np.int8
UInt8 = int | np.uint8
Int16 = int | np.int16
UInt16 = int | np.uint16
Int32 = int | np.int32
UInt32 = int | np.uint32
Int64 = int | np.int64
UInt64 = int | np.uint64
Size = int | np.uint64
Float32 = float | np.float32
Float64 = float | np.float64
ComplexFloat = complex | np.complex64
ComplexDouble = complex | np.complex128

Date = datetime.date | np.datetime64
Time = datetime.time | np.timedelta64
DateTime = datetime.datetime | np.datetime64


Integer = Int8 | UInt8 | Int16 | UInt16 | Int32 | UInt32 | Int64 | UInt64 | Size
Floating = Float32 | Float64
Complex = ComplexFloat | ComplexDouble


class Types(Flag):
    NONE = 0
    BOOL = auto()
    INT8 = auto()
    UINT8 = auto()
    INT16 = auto()
    UINT16 = auto()
    INT32 = auto()
    UINT32 = auto()
    INT64 = auto()
    UINT64 = auto()
    SIZE = auto()
    INTEGER = INT8 | UINT8 | INT16 | UINT16 | INT32 | UINT32 | INT64 | UINT64 | SIZE
    FLOAT32 = auto()
    FLOAT64 = auto()
    FLOATS = FLOAT32 | FLOAT64
    COMPLEX_FLOAT32 = auto()
    COMPLEX_FLOAT64 = auto()
    COMPLEX = COMPLEX_FLOAT32 | COMPLEX_FLOAT64
    STRING = auto()
    DATE = auto()
    TIME = auto()
    DATETIME = auto()
    VECTOR = auto()
    ARRAY = auto()

class ProtocolError(Exception):
    pass
