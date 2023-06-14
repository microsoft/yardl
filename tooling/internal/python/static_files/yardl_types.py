from enum import Enum, auto
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
Float = Float32 | Float64
Complex = ComplexFloat | ComplexDouble

class TypePreference(Enum):
    PYTHON = auto()
    NUMPY = auto()

@dataclass(frozen=True, slots=True, kw_only=True)
class ReadTypePreferences:
    bool: TypePreference = TypePreference.PYTHON
    integer: TypePreference = TypePreference.PYTHON
    floatingPoint: TypePreference = TypePreference.PYTHON
    complex: TypePreference = TypePreference.PYTHON
    date: TypePreference = TypePreference.PYTHON
    time: TypePreference = TypePreference.PYTHON
    dateTime: TypePreference = TypePreference.PYTHON
