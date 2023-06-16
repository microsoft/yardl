import datetime
from types import TracebackType
from typing import BinaryIO, Iterable, TypeVar, Protocol, Generic, Any, Optional, Tuple, cast
from collections.abc import Callable
from abc import ABC, abstractmethod
from functools import partial
import struct
import sys
import numpy as np
from numpy.lib import recfunctions
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
        string_descriptor.write(self._stream, schema)

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

T = TypeVar("T")

class TypeDescriptor(Generic[T], ABC):
    def __init__(self, dtype: npt.DTypeLike) -> None:
        self._dtype: np.dtype[Any] = np.dtype(dtype)

    def overall_dtype(self) -> np.dtype[Any]:
        return self._dtype

    def struct_format_str(self) -> str | None:
        return None

    @abstractmethod
    def write(self, stream: CodedOutputStream, value: T) -> None:
        raise NotImplementedError

    def is_trivially_serializable(self) -> bool:
        return False

class StructDescriptor(TypeDescriptor[T]):
    def __init__(self, dtype: npt.DTypeLike, format_string: str) -> None:
        super().__init__(dtype)
        self.struct = struct.Struct(format_string)

    def write(self, stream: CodedOutputStream, value: T) -> None:
        stream.write(self.struct, value)

    def struct_format_str(self) -> str:
        return self.struct.format

class BoolDescriptor(StructDescriptor[bool]):
    def __init__(self) -> None:
        super().__init__(np.bool_, "<?")

write_descriptor = BoolDescriptor()

class Int8Descriptor(StructDescriptor[Int8]):
    def __init__(self) -> None:
        super().__init__(np.int8, "<b")

    def is_trivially_serializable(self) -> bool:
        return True

int8_descriptor = Int8Descriptor()

class UInt8Descriptor(StructDescriptor[UInt8]):
    def __init__(self) -> None:
        super().__init__(np.uint8, "<B")

    def is_trivially_serializable(self) -> bool:
        return True

uint8_descriptor = UInt8Descriptor()

class Int16Descriptor(TypeDescriptor[Int16]):
    def __init__(self) -> None:
        super().__init__(np.int16)

    def write(self, stream: CodedOutputStream, value: Int16) -> None:
        if isinstance(value, int):
            if value < 0 or value > UINT16_MAX:
                raise ValueError(
                    f"Value {value} is outside the range of an unsigned 16-bit integer"
                )
        elif not isinstance(value, cast(type, np.uint16)):
            raise ValueError(f"Value in not an unsigned 16-bit integer: {value}")

        stream.write_unsigned_varint(value)

int16_descriptor = Int16Descriptor()

class UInt16Descriptor(TypeDescriptor[UInt16]):
    def __init__(self) -> None:
        super().__init__(np.uint16)

    def write(self, stream: CodedOutputStream, value: UInt16) -> None:
        if isinstance(value, int):
            if value < 0 or value > UINT16_MAX:
                raise ValueError(
                    f"Value {value} is outside the range of an unsigned 16-bit integer"
                )
        elif not isinstance(value, cast(type, np.uint16)):
            raise ValueError(f"Value in not an unsigned 16-bit integer: {value}")

        stream.write_unsigned_varint(value)

uint16_descriptor = UInt16Descriptor()

class Int32Descriptor(TypeDescriptor[Int32]):
    def __init__(self) -> None:
        super().__init__(np.int32)

    def write(self, stream: CodedOutputStream, value: Int32) -> None:
        if isinstance(value, int):
            if value < 0 or value > UINT32_MAX:
                raise ValueError(
                    f"Value {value} is outside the range of an unsigned 32-bit integer"
                )
        elif not isinstance(value, cast(type, np.int32)):
            raise ValueError(f"Value in not an unsigned 32-bit integer: {value}")

        stream.write_unsigned_varint(value)

int32_descriptor = Int32Descriptor()

class UInt32Descriptor(TypeDescriptor[UInt32]):
    def __init__(self) -> None:
        super().__init__(np.uint32)

    def write(self, stream: CodedOutputStream, value: UInt32) -> None:
        if isinstance(value, int):
            if value < 0 or value > UINT32_MAX:
                raise ValueError(
                    f"Value {value} is outside the range of an unsigned 32-bit integer"
                )
        elif not isinstance(value, cast(type, np.uint32)):
            raise ValueError(f"Value in not an unsigned 32-bit integer: {value}")

        stream.write_unsigned_varint(value)

uint32_descriptor = UInt32Descriptor()

class Int64Descriptor(TypeDescriptor[Int64]):
    def __init__(self) -> None:
        super().__init__(np.int64)

    def write(self, stream: CodedOutputStream, value: Int64) -> None:
        if isinstance(value, int):
            if value < INT64_MIN or value > INT64_MAX:
                raise ValueError(
                    f"Value {value} is outside the range of a signed 64-bit integer"
                )
        elif not isinstance(value, cast(type, np.int64)):
            raise ValueError(f"Value in not a signed 64-bit integer: {value}")

        stream.write_signed_varint(value)

int64_descriptor = Int64Descriptor()

class UInt64Descriptor(TypeDescriptor[UInt64]):
    def __init__(self) -> None:
        super().__init__(np.uint64)

    def write(self, stream: CodedOutputStream, value: UInt64) -> None:
        if isinstance(value, int):
            if value < 0 or value > UINT64_MAX:
                raise ValueError(
                    f"Value {value} is outside the range of an unsigned 64-bit integer"
                )
        elif not isinstance(value, cast(type, np.uint64)):
            raise ValueError(f"Value in not an unsigned 64-bit integer: {value}")

        stream.write_unsigned_varint(value)

uint64_descriptor = UInt64Descriptor()

size_descriptor = uint64_descriptor


class Float32Descriptor(StructDescriptor[Float32]):
    def __init__(self) -> None:
        super().__init__(np.float32, "<f")

    def is_trivially_serializable(self) -> bool:
        return True

float32_descriptor = Float32Descriptor()

class Float64Descriptor(StructDescriptor[Float64]):
    def __init__(self) -> None:
        super().__init__(np.float64, "<d")

    def is_trivially_serializable(self) -> bool:
        return True

float64_descriptor = Float64Descriptor()

class Complex32Descriptor(StructDescriptor[ComplexFloat]):
    def __init__(self) -> None:
        super().__init__(np.complex64, "<ff")

    def is_trivially_serializable(self) -> bool:
        return True

complex32_descriptor = Complex32Descriptor()

class Complex64Descriptor(StructDescriptor[ComplexDouble]):
    def __init__(self) -> None:
        super().__init__(np.complex128, "<dd")

    def is_trivially_serializable(self) -> bool:
        return True

complex64_descriptor = Complex64Descriptor()

class StringDescriptor(TypeDescriptor[str]):
    def __init__(self) -> None:
        super().__init__(np.object_)

    def write(self, stream: CodedOutputStream, value: str) -> None:
        b = value.encode("utf-8")
        stream.write_unsigned_varint(len(b))
        stream.write_bytes(b)

string_descriptor = StringDescriptor()

EPOCH_ORDINAL_DAYS = datetime.date(1970, 1, 1).toordinal()
DATETIME_DAYS_DTYPE = np.dtype("datetime64[D]")

class DateDescriptor(TypeDescriptor[Date]):
    def __init__(self) -> None:
        super().__init__(np.datetime64)

    def write(self, stream: CodedOutputStream, value: Date) -> None:
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

date_descriptor = DateDescriptor()

TIMEDELTA_NANOSECONDS_DTYPE = np.dtype("timedelta64[ns]")

class TimeDescriptor(TypeDescriptor[Time]):
    def __init__(self) -> None:
        super().__init__(np.timedelta64)

    def write(self, stream: CodedOutputStream, value: Time) -> None:
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

time_descriptor = TimeDescriptor()

DATETIME_NANOSECONDS_DTYPE = np.dtype("datetime64[ns]")
EPOCH_DATETIME = datetime.datetime.utcfromtimestamp(0)

class DateTimeDescriptor(TypeDescriptor[DateTime]):
    def __init__(self) -> None:
        super().__init__(np.datetime64)

    def write(self, stream: CodedOutputStream, value: DateTime) -> None:
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

datetime_descriptor = DateTimeDescriptor()

class NoneDescriptor(TypeDescriptor[None]):
    def __init__(self) -> None:
        super().__init__(np.object_)

    def write(self, stream: CodedOutputStream, value: None) -> None:
        pass

none_descriptor = NoneDescriptor()

def write_none(stream: CodedOutputStream, value: None) -> None:
    pass

class EnumDescriptor(Generic[T], TypeDescriptor[Enum]):
    def __init__(self, integer_descriptor: TypeDescriptor[T]) -> None:
        super().__init__(integer_descriptor.overall_dtype())
        self.integer_descriptor = integer_descriptor

    def write(self, stream: CodedOutputStream, value: Enum) -> None:
        self.integer_descriptor.write(stream, value.value)

    def is_trivially_serializable(self) -> bool:
        return self.integer_descriptor.is_trivially_serializable()


class OptionalDescriptor(TypeDescriptor[T]):
    def __init__(self, element_descriptor: TypeDescriptor[T]) -> None:
        super().__init__(np.dtype([("has_value", np.bool_), ("value", element_descriptor.overall_dtype())]))
        self.element_descriptor = element_descriptor

    def write(self, stream: CodedOutputStream, value: Optional[T]) -> None:
        if value is None:
            stream.write_byte(0)
        else:
            stream.write_byte(1)
            self.element_descriptor.write(stream, value)


class UnionDescriptor(TypeDescriptor[Any]):
    def __init__(self, cases: list[Tuple[type, TypeDescriptor[Any]]]) -> None:
        super().__init__(np.object_)
        self.cases = cases

    def write(self, stream: CodedOutputStream, value: Any) -> None:
        for i, (case_type, case_descriptor) in enumerate(self.cases):
            if isinstance(value, case_type):
                stream.write_byte(i)
                case_descriptor.write(stream, value)
                return

        raise ValueError(f"Incorrect union type {type(value)}")


class StreamDescriptor(TypeDescriptor[Iterable[T]]):
    def __init__(self, element_descriptor: TypeDescriptor[T]) -> None:
        super().__init__(np.object_)
        self.element_descriptor = element_descriptor

    def write(self, stream: CodedOutputStream, value: Iterable[T]) -> None:
        for element in value:
            stream.write_byte(1)
            self.element_descriptor.write(stream, element)


class FixedVectorDescriptor(TypeDescriptor[list[T]]):
    def __init__(self, element_descriptor: TypeDescriptor[T], length: int) -> None:
        super().__init__(np.dtype((element_descriptor.overall_dtype(), length)))
        self.element_descriptor = element_descriptor
        self.length = length

    def write(self, stream: CodedOutputStream, value: list[T]) -> None:
        if len(value) != self.length:
            raise ValueError(
                f"Expected a list of length {self.length}, got {len(value)}"
            )
        for element in value:
            self.element_descriptor.write(stream, element)

    def is_trivially_serializable(self) -> bool:
        return self.element_descriptor.is_trivially_serializable()


class VectorDescriptor(TypeDescriptor[list[T]]):
    def __init__(self, element_descriptor: TypeDescriptor[T]) -> None:
        super().__init__(np.object_)
        self.element_descriptor = element_descriptor

    def write(self, stream: CodedOutputStream, value: list[T]) -> None:
        stream.write_unsigned_varint(len(value))
        for element in value:
            self.element_descriptor.write(stream, element)


TKey = TypeVar("TKey")
TValue = TypeVar("TValue")


class MapDescriptor(TypeDescriptor[dict[TKey, TValue]]):
    def __init__(self, key_descriptor: TypeDescriptor[TKey], value_descriptor: TypeDescriptor[TValue]) -> None:
        super().__init__(np.object_)
        self.key_descriptor = key_descriptor
        self.value_descriptor = value_descriptor

    def write(self, stream: CodedOutputStream, value: dict[TKey, TValue]) -> None:
        stream.write_unsigned_varint(len(value))
        for k, v in value.items():
            self.key_descriptor.write(stream, k)
            self.value_descriptor.write(stream, v)

class NDArrayDescriptorBase(Generic[T], TypeDescriptor[npt.NDArray[Any]]):
    def __init__(
        self,
        overall_dtype: npt.DTypeLike,
        element_descriptor: TypeDescriptor[T],
        dtype: npt.DTypeLike,
    ) -> None:
        super().__init__(overall_dtype)
        self._array_dtype: np.dtype[Any] = dtype if isinstance(dtype, np.dtype) else np.dtype(dtype)
        self._element_descriptor = element_descriptor

    def _write_data(self, stream: CodedOutputStream, value: npt.NDArray[Any]) -> None:
        if value.dtype != self._array_dtype:
            # see if it's the same dtype but packed, not aligned
            packed_dtype = recfunctions.repack_fields(self._array_dtype, align=False, recurse=True) # type: ignore
            if packed_dtype != value.dtype:
                raise ValueError(f"Expected dtype {self._array_dtype} or {packed_dtype}, got {value.dtype}")

        if self._is_current_array_trivially_serializable(value):
            stream.write_bytes_directly(value.data)
        else:
            to_iterate = value if value.dtype.fields is None else cast(npt.NDArray[Any], value.view(np.recarray))
            for element in to_iterate.flat:
                self._element_descriptor.write(stream, element)

    def _is_current_array_trivially_serializable(self, value: npt.NDArray[Any]) -> bool:
        return self._element_descriptor.is_trivially_serializable() and value.flags.c_contiguous \
            and (self._array_dtype.fields is None or all(f != "" for f in self._array_dtype.fields))


class DynamicNDArrayDescriptor(NDArrayDescriptorBase[T]):
    def __init__(
        self,
        element_descriptor: TypeDescriptor[T],
    ) -> None:
        super().__init__(np.object_, element_descriptor, element_descriptor.overall_dtype())

    def write(self, stream: CodedOutputStream, value: npt.NDArray[Any]) -> None:
        stream.write_unsigned_varint(value.ndim)
        for dim in value.shape:
            stream.write_unsigned_varint(dim)

        self._write_data(stream, value)


class NDArrayDescriptor(Generic[T], NDArrayDescriptorBase[T]):
    def __init__(
        self,
        element_descriptor: TypeDescriptor[T],
        potentially_trivially_serializable: bool,
        ndims: int,
    ) -> None:
        super().__init__(np.object_, element_descriptor, element_descriptor.overall_dtype())
        self.ndims = ndims

    def write(self, stream: CodedOutputStream, value: npt.NDArray[Any]) -> None:
        if value.ndim != self.ndims:
            raise ValueError(f"Expected {self.ndims} dimensions, got {value.ndim}")

        for dim in value.shape:
            stream.write_unsigned_varint(dim)

        self._write_data(stream, value)


class FixedNDArrayDescriptor(Generic[T], NDArrayDescriptorBase[T]):
    def __init__(
        self,
        element_descriptor: TypeDescriptor[T],
        shape: tuple[int, ...],
    ) -> None:
        dtype = element_descriptor.overall_dtype()
        super().__init__(np.dtype((dtype, shape)), element_descriptor, dtype)
        self.shape = shape

    def write(self, stream: CodedOutputStream, value: npt.NDArray[Any]) -> None:
        if value.shape != self.shape:
            raise ValueError(f"Expected shape {self.shape}, got {value.shape}")

        self._write_data(stream, value)

    def is_trivially_serializable(self) -> bool:
        return self._element_descriptor.is_trivially_serializable()


class RecordDescriptor(TypeDescriptor[T]):
    def __init__(self, field_descriptors: list[Tuple[str, TypeDescriptor[Any]]]) -> None:
        super().__init__(np.dtype([(name, descriptor.overall_dtype()) for name, descriptor in field_descriptors], align=True))
        combined_format = "<"
        for _, field_descriptor in field_descriptors:
            fmt = field_descriptor.struct_format_str()
            if fmt is None:
                combined_format = None
                break
            else:
                combined_format += fmt[1:] if fmt[0] == '<' else fmt

        if combined_format is None:
            self._struct = None
        else:
            self._struct = struct.Struct(combined_format)

        self._field_descriptors = field_descriptors

    def is_trivially_serializable(self) -> bool:
        return all(descriptor.is_trivially_serializable() for _, descriptor in self._field_descriptors)

    def _write(self, stream: CodedOutputStream, *values: Any) -> None:
        if self._struct:
            stream.write(self._struct, *values)
        else:
            for i, (_,descriptor) in enumerate(self._field_descriptors):
                descriptor.write(stream, values[i])


# Only used in the header
int32_struct = struct.Struct("<i")
assert int32_struct.size == 4


def write_fixed_int32(stream: CodedOutputStream, value: int) -> None:
    if value < INT32_MIN or value > INT32_MAX:
        raise ValueError(
            f"Value {value} is outside the range of a signed 32-bit integer"
        )
    stream.write(int32_struct, value)
