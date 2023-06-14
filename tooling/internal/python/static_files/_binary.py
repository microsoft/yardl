import datetime
from types import TracebackType
from typing import BinaryIO, Iterable, TypeVar, Protocol, Generic, Any, Optional, Tuple
from collections.abc import Callable
from abc import ABC
from functools import partial
import struct
import sys
import numpy as np
import numpy.typing as npt
from .yardl_types import *

MAGIC_BYTES: bytes = b"yardl"
CURRENT_BINARY_FORMAT_VERSION: int = 1

INT8_MIN: int = np.iinfo(np.int8).min
INT8_MAX: int = np.iinfo(np.int8).max

UINT8_MAX: int = np.iinfo(np.uint8).max

INT16_MIN: int = np.iinfo(np.int16).min
INT16_MAX: int = np.iinfo(np.int16).max

UINT16_MAX: int = np.iinfo(np.uint16).max

INT32_MIN: int = np.iinfo(np.int32).min
INT32_MAX: int = np.iinfo(np.int32).max

UINT32_MAX: int = np.iinfo(np.uint32).max

INT64_MIN: int = np.iinfo(np.int64).min
INT64_MAX: int = np.iinfo(np.int64).max

UINT64_MAX: int = np.iinfo(np.uint64).max


class BinaryProtocolWriter(ABC):
    def __init__(self, stream: BinaryIO | str, schema: str) -> None:
        self._stream = CodedOutputStream(stream)
        self._stream.write_bytes(MAGIC_BYTES)
        write_fixed_int32(self._stream, CURRENT_BINARY_FORMAT_VERSION)
        self._stream.write_string(schema)

    def __enter__(self):
        return self

    def __exit__(
        self,
        exc_type: Optional[type[BaseException]],
        exc: Optional[BaseException],
        traceback: Optional[TracebackType],
    ) -> None:
        self.close()

    def close(self) -> None:
        try:
            self._close()
        finally:
            self._stream.close()

    def _close(self) -> None:
        pass


class CodedOutputStream:
    def __init__(self, stream: BinaryIO | str, *, buffer_size: int = 65536) -> None:
        if isinstance(stream, str):
            self._stream = open(stream, "wb")
            self._owns_stream = True
        else:
            self._stream = stream
            self._owns_stream = False

        self._buffer = bytearray(buffer_size)
        self._view = memoryview(self._buffer)

    def __enter__(self) -> "CodedOutputStream":
        return self

    def __exit__(
        self,
        exc_type: Optional[type[BaseException]],
        exc: Optional[BaseException],
        traceback: Optional[TracebackType],
    ) -> None:
        self.close()

    def close(self) -> None:
        self.flush()
        if self._owns_stream:
            self._stream.close()

    def flush(self) -> None:
        buffer_filled_count = len(self._buffer) - len(self._view)
        if buffer_filled_count > 0:
            self._stream.write(self._buffer[:buffer_filled_count])
            self._stream.flush()
            self._view = memoryview(self._buffer)

    def write(self, formatter: struct.Struct, *args: Any) -> None:
        if len(self._view) < formatter.size:
            self.flush()

        formatter.pack_into(self._view, 0, *args)
        self._view = self._view[formatter.size :]

    def write_bytes(self, data: bytes | bytearray) -> None:
        if len(data) > len(self._view):
            self.flush()
            self._stream.write(data)
        else:
            self._view[: len(data)] = data
            self._view = self._view[len(data) :]

    def write_bytes_directly(self, data: bytes | bytearray | memoryview) -> None:
        self.flush()
        self._stream.write(data)

    def write_byte(self, value: Integer) -> None:
        assert 0 <= value <= UINT8_MAX
        self._view[0] = value
        self._view = self._view[1:]

    def write_unsigned_varint(self, value: Integer) -> None:
        if len(self._view) < 10:
            self.flush()

        while True:
            if value < 0x80:
                self.write_byte(value)
                return

            self.write_byte((value & 0x7F) | 0x80)
            value >>= 7

    def zigzag_encode(self, value: Integer) -> Integer:
        return (value << 1) ^ (value >> 63)

    def write_signed_varint(self, value: Integer) -> None:
        self.write_unsigned_varint(self.zigzag_encode(value))

    def write_string(self, value: str) -> None:
        self.write_bytes(value.encode("utf-8"))


bool_struct = struct.Struct("<?")
assert bool_struct.size == 1


def write_bool(stream: CodedOutputStream, value: Bool) -> None:
    stream.write(bool_struct, value)


int8_struct = struct.Struct("<b")
assert int8_struct.size == 1


def write_int8(stream: CodedOutputStream, value: Int8) -> None:
    stream.write(int8_struct, value)


uint8_struct = struct.Struct("<B")
assert uint8_struct.size == 1


def write_uint8(stream: CodedOutputStream, value: UInt8) -> None:
    stream.write(uint8_struct, value)


def write_int16(stream: CodedOutputStream, value: Int16) -> None:
    if value < INT16_MIN or value > INT16_MAX:
        raise ValueError(
            f"Value {value} is outside the range of a signed 16-bit integer"
        )
    stream.write_signed_varint(value)


def write_uint16(stream: CodedOutputStream, value: UInt16) -> None:
    if value < 0 or value > UINT16_MAX:
        raise ValueError(
            f"Value {value} is outside the range of an unsigned 16-bit integer"
        )
    stream.write_unsigned_varint(value)


def write_int32(stream: CodedOutputStream, value: Int32) -> None:
    if value < INT32_MIN or value > INT32_MAX:
        raise ValueError(
            f"Value {value} is outside the range of a signed 32-bit integer"
        )
    stream.write_signed_varint(value)


def write_uint32(stream: CodedOutputStream, value: UInt32) -> None:
    if value < 0 or value > UINT32_MAX:
        raise ValueError(
            f"Value {value} is outside the range of an unsigned 32-bit integer"
        )
    stream.write_unsigned_varint(value)


def write_int64(stream: CodedOutputStream, value: Int64) -> None:
    if value < INT64_MIN or value > INT64_MAX:
        raise ValueError(
            f"Value {value} is outside the range of a signed 64-bit integer"
        )
    stream.write_signed_varint(value)


def write_uint64(stream: CodedOutputStream, value: UInt64) -> None:
    if value < 0 or value > UINT64_MAX:
        raise ValueError(
            f"Value {value} is outside the range of an unsigned 64-bit integer"
        )
    stream.write_unsigned_varint(value)


def write_size(stream: CodedOutputStream, value: Size) -> None:
    write_uint64(stream, value)


float32_struct = struct.Struct("<f")
assert float32_struct.size == 4


def write_float32(stream: CodedOutputStream, value: float) -> None:
    stream.write(float32_struct, value)


float64_struct = struct.Struct("<d")
assert float64_struct.size == 8


def write_float64(stream: CodedOutputStream, value: float) -> None:
    stream.write(float64_struct, value)


complex32_struct = struct.Struct("<ff")
assert complex32_struct.size == 8


def write_complex32(stream: CodedOutputStream, value: complex) -> None:
    stream.write(complex32_struct, value.real, value.imag)


complex64_struct = struct.Struct("<dd")
assert complex64_struct.size == 16


def write_complex64(stream: CodedOutputStream, value: complex) -> None:
    stream.write(complex64_struct, value.real, value.imag)


def write_string(stream: CodedOutputStream, value: str) -> None:
    stream.write_string(value)


EPOCH_ORDINAL_DAYS = datetime.date(1970, 1, 1).toordinal()
DATETIME_DAYS_DTYPE = np.dtype("datetime64[D]")


def write_date(stream: CodedOutputStream, value: Date) -> None:
    if isinstance(value, datetime.date):
        stream.write_signed_varint(value.toordinal() - EPOCH_ORDINAL_DAYS)
    else:
        if not isinstance(value, np.datetime64):
            raise TypeError(
                f"Expected datetime.date or numpy.datetime64, got {type(value)}"
            )

        if value.dtype == DATETIME_DAYS_DTYPE:
            stream.write_signed_varint(value.astype(np.int32))
        else:
            stream.write_signed_varint(
                value.astype(DATETIME_DAYS_DTYPE).astype(np.int32)
            )


TIMEDELTA_NANOSECONDS_DTYPE = np.dtype("timedelta64[ns]")


def write_time(stream: CodedOutputStream, value: Time) -> None:
    if isinstance(value, datetime.time):
        nanoseconds_since_midnight = (
            value.hour * 3_600_000_000_000
            + value.minute * 60_000_000_000
            + value.second * 1_000_000_000
            + value.microsecond * 1_000
        )
        stream.write_signed_varint(nanoseconds_since_midnight)
    else:
        if not isinstance(value, np.timedelta64):
            raise TypeError(
                f"Expected a datetime.time or np.timedelta64, got {type(value)}"
            )

        if value.dtype == TIMEDELTA_NANOSECONDS_DTYPE:
            stream.write_signed_varint(value.astype(np.int64))
        else:
            stream.write_signed_varint(
                value.astype(DATETIME_NANOSECONDS_DTYPE).astype(np.int64)
            )


DATETIME_NANOSECONDS_DTYPE = np.dtype("datetime64[ns]")
EPOCH_DATETIME = datetime.datetime.utcfromtimestamp(0)


def write_datetime(stream: CodedOutputStream, value: DateTime) -> None:
    if isinstance(value, datetime.datetime):
        delta = value - EPOCH_DATETIME
        nanoseconds_since_epoch = int(delta.total_seconds() * 1e9)
        stream.write_signed_varint(nanoseconds_since_epoch)
    else:
        if not isinstance(value, np.datetime64):
            raise TypeError(
                f"Expected datetime.datetime or numpy.datetime64, got {type(value)}"
            )

        if value.dtype == DATETIME_NANOSECONDS_DTYPE:
            stream.write_signed_varint(value.astype(np.int64))
        else:
            stream.write_signed_varint(
                value.astype(DATETIME_NANOSECONDS_DTYPE).astype(np.int64)
            )


def write_none(stream: CodedOutputStream, value: None) -> None:
    pass


def write_enum(stream: CodedOutputStream, value: Enum) -> None:
    stream.write_signed_varint(value.value)


T = TypeVar("T")
Writer = Callable[[CodedOutputStream, T], None]

class EnumWriter(Generic[T]):
    def __init__(self, write_integer: Writer[T]) -> None:
        self.write_integer = write_integer

    def __call__(self, stream: CodedOutputStream, value: Enum) -> None:
        self.write_integer(stream, value.value)


class OptionalWriter(Generic[T]):
    def __init__(self, write_element: Writer[T]) -> None:
        self.write_element = write_element

    def __call__(self, stream: CodedOutputStream, value: Optional[T]) -> None:
        if value is None:
            stream.write_byte(0)
        else:
            stream.write_byte(1)
            self.write_element(stream, value)


class UnionWriter:
    def __init__(self, cases: list[Tuple[type, Writer[Any]]]) -> None:
        self.cases = cases

    def __call__(self, stream: CodedOutputStream, value: Any) -> None:
        for i, (case_type, case_writer) in enumerate(self.cases):
            if isinstance(value, case_type):
                stream.write_byte(i)
                case_writer(stream, value)
                return

        raise ValueError(f"Incorrect union type {type(value)}")


class StreamWriter(Generic[T]):
    def __init__(self, write_element: Writer[T]) -> None:
        self.write_element = write_element

    def __call__(self, stream: CodedOutputStream, value: Iterable[T]) -> None:
        for element in value:
            stream.write_byte(1)
            self.write_element(stream, element)


class FixedVectorWriter(Generic[T]):
    def __init__(self, length: int, write_element: Writer[T]) -> None:
        self.length = length
        self.write_element = write_element

    def __call__(self, stream: CodedOutputStream, value: list[T]) -> None:
        assert len(value) == self.length
        for element in value:
            self.write_element(stream, element)


class DynamicVectorWriter(Generic[T]):
    def __init__(self, write_element: Writer[T]) -> None:
        self.write_element = write_element

    def __call__(self, stream: CodedOutputStream, value: list[T]) -> None:
        stream.write_unsigned_varint(len(value))
        for element in value:
            self.write_element(stream, element)


TKey = TypeVar("TKey")
TValue = TypeVar("TValue")


class MapWriter(Generic[TKey, TValue]):
    def __init__(self, write_key: Writer[TKey], write_value: Writer[TValue]) -> None:
        self.write_key = write_key
        self.write_value = write_value

    def __call__(self, stream: CodedOutputStream, value: dict[TKey, TValue]) -> None:
        stream.write_unsigned_varint(len(value))
        for k, v in value.items():
            self.write_key(stream, k)
            self.write_value(stream, v)


class DynamicNDArrayWriter(Generic[T]):
    def __init__(
        self,
        write_element: Writer[T],
        dtype: np.dtype[Any],
        trivially_serializable: bool,
    ) -> None:
        self.dtype = dtype
        self.write_element = write_element
        self.trivially_serializable = trivially_serializable

    def __call__(self, stream: CodedOutputStream, value: npt.NDArray[Any]) -> None:
        assert value.dtype == self.dtype, "dtype mismatch"
        stream.write_unsigned_varint(value.ndim)
        for dim in value.shape:
            stream.write_unsigned_varint(dim)

        if self.trivially_serializable and value.flags.c_contiguous:
            stream.write_bytes_directly(value.data)
        else:
            for element in value.flat:
                self.write_element(stream, element)


class NDArrayWriter(Generic[T]):
    def __init__(
        self,
        write_element: Writer[T],
        dtype: np.dtype[Any],
        ndims: int,
        trivially_serializable: bool,
    ) -> None:
        self.write_element = write_element
        self.dtype = dtype
        self.ndims = ndims
        self.trivially_serializable = trivially_serializable

    def __call__(self, stream: CodedOutputStream, value: npt.NDArray[Any]) -> None:
        assert value.dtype == self.dtype, "dtype mismatch"
        assert value.ndim == self.ndims, "dimension count mismatch"

        for dim in value.shape:
            stream.write_unsigned_varint(dim)

        if self.trivially_serializable and value.flags.c_contiguous:
            stream.write_bytes_directly(value.data)
        else:
            for element in value.flat:
                self.write_element(stream, element)


class FixedNDArrayWriter(Generic[T]):
    def __init__(
        self,
        write_element: Writer[T],
        dtype: np.dtype[Any],
        shape: tuple[int, ...],
        trivially_serializable: bool,
    ) -> None:
        self.write_element = write_element
        self.dtype = dtype
        self.shape = shape
        self.trivially_serializable = trivially_serializable

    def __call__(self, stream: CodedOutputStream, value: npt.NDArray[Any]) -> None:
        assert value.dtype == self.dtype, "dtype mismatch"
        assert value.shape == self.shape, "shape mismatch"

        if self.trivially_serializable and value.flags.c_contiguous:
            stream.write_bytes_directly(value.data)
        else:
            for element in value.flat:
                self.write_element(stream, element)


# Only used in the header
int32_struct = struct.Struct("<i")
assert int32_struct.size == 4


def write_fixed_int32(stream: CodedOutputStream, value: int) -> None:
    if value < INT32_MIN or value > INT32_MAX:
        raise ValueError(
            f"Value {value} is outside the range of a signed 32-bit integer"
        )
    stream.write(int32_struct, value)
