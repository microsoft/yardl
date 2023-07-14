from enum import Flag, auto
import numpy as np
import datetime

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

String = str

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

    ALL = (
        BOOL
        | INTEGER
        | FLOATS
        | COMPLEX
        | STRING
        | DATE
        | TIME
        | DATETIME
        | VECTOR
        | ARRAY
    )


class ProtocolError(Exception):
    pass


_EPOCH_ORDINAL_DAYS = datetime.date(1970, 1, 1).toordinal()


def dates_equal(a: Date, b: Date) -> bool:
    if type(a) == type(b):
        return a == b

    if isinstance(a, datetime.date):
        if isinstance(b, np.datetime64):
            b, a = a, b
        else:
            return False
    else:
        if not isinstance(a, np.datetime64) or not isinstance(b, datetime.date):
            return False

    # a is now a datetime64 and b is a datetime.date
    b_days_since_epoch = b.toordinal() - _EPOCH_ORDINAL_DAYS

    return a == np.datetime64(b_days_since_epoch, "D")


def times_equal(a: Time, b: Time) -> bool:
    if type(a) == type(b):
        return a == b

    if isinstance(a, datetime.time):
        if isinstance(b, np.timedelta64):
            b, a = a, b
        else:
            return False
    else:
        if not isinstance(a, np.timedelta64) or not isinstance(b, datetime.time):
            return False

    # a is now a timedelta64 and b is a datetime.time
    b_nanoseconds_since_midnight = (
        b.hour * 3_600_000_000_000
        + b.minute * 60_000_000_000
        + b.second * 1_000_000_000
        + b.microsecond * 1_000
    )

    return a == np.timedelta64(b_nanoseconds_since_midnight, "ns")


_EPOCH_DATETIME = datetime.datetime.utcfromtimestamp(0)


def datetimes_equal(a: DateTime, b: DateTime) -> bool:
    if type(a) == type(b):
        return a == b

    if isinstance(a, datetime.datetime):
        if isinstance(b, np.datetime64):
            b, a = a, b
        else:
            return False
    else:
        if not isinstance(a, np.datetime64) or not isinstance(b, datetime.datetime):
            return False

    # a is now a datetime64 and b is a datetime.datetime
    b_delta = b - _EPOCH_DATETIME
    b_nanoseconds_since_epoch = int(b_delta.total_seconds() * 1e6) * 1000

    return a == np.datetime64(b_nanoseconds_since_epoch, "ns")


def structural_equal(a: object, b: object) -> bool:
    if a is None:
        return b is None

    if isinstance(a, list):
        if not isinstance(b, list):
            if isinstance(b, np.ndarray):
                return b.shape == (len(a),) and all(
                    structural_equal(x, y) for x, y in zip(a, b)
                )
            return False
        return len(a) == len(b) and all(structural_equal(x, y) for x, y in zip(a, b))

    if isinstance(a, np.ndarray):
        if not isinstance(b, np.ndarray):
            if isinstance(b, list):
                return a.shape == (len(b),) and all(
                    structural_equal(x, y) for x, y in zip(a, b)
                )
            return False
        if a.dtype.hasobject:
            return (
                a.dtype == b.dtype
                and a.shape == b.shape
                and all(structural_equal(x, y) for x, y in zip(a, b))
            )
        return np.array_equal(a, b)

    if isinstance(a, np.void):
        if not isinstance(b, np.void):
            return b == a
        return a.dtype == b.dtype and all(structural_equal(x, y) for x, y in zip(a, b))

    if isinstance(b, np.void):
        return a == b

    if isinstance(a, tuple):
        return (
            isinstance(b, tuple)
            and len(a) == len(b)
            and all(structural_equal(x, y) for x, y in zip(a, b))
        )

    if isinstance(a, datetime.time) or isinstance(a, np.timedelta64):
        return times_equal(a, b)

    if isinstance(a, datetime.datetime) or isinstance(a, np.datetime64):
        return datetimes_equal(a, b)

    return a == b
