import datetime
from io import BufferedIOBase, BufferedReader, RawIOBase
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

if sys.byteorder != "little":
    raise RuntimeError("Only little-endian systems are currently supported")

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

class BinaryProtocolReader(ABC):
    def __init__(self, stream: BufferedReader | str, expected_schema : str) -> None:
        self._stream = CodedInputStream(stream)
        magic_bytes = self._stream.read_view(len(MAGIC_BYTES))
        if magic_bytes != MAGIC_BYTES:
            raise RuntimeError("Invalid magic bytes")

        version = read_fixed_int32(self._stream)
        if version != CURRENT_BINARY_FORMAT_VERSION:
            raise RuntimeError("Invalid binary format version")

        self._schema = string_descriptor.read(self._stream)
        if self._schema != expected_schema:
            raise RuntimeError("Invalid schema")

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

class CodedInputStream:
    def __init__(self, stream: BufferedReader | str, *, buffer_size: int = 65536) -> None:
        if isinstance(stream, str):
            self._stream = open(stream, "rb")
            self._owns_stream = True
        else:
            if not isinstance(stream, BufferedReader):
                self._stream = BufferedReader(stream)
            else:
                self._stream = stream
            self._owns_stream = False

        self._last_read_count = 0
        self._buffer = bytearray(buffer_size)
        self._view = memoryview(self._buffer)
        self._offset = 0
        self._at_end = False

    def close(self) -> None:
        if self._owns_stream:
            self._stream.close()

    def read(self, formatter: struct.Struct) -> tuple[Any, ...]:
        if self._last_read_count - self._offset < formatter.size:
            self._fill_buffer(formatter.size)

        result = formatter.unpack_from(self._buffer, self._offset)
        self._offset += formatter.size
        return result

    def read_byte(self) -> int:
        if self._last_read_count - self._offset < 1:
            self._fill_buffer(1)

        result = self._buffer[self._offset]
        self._offset += 1
        return result

    def read_unsigned_varint(self) -> int:
        result = 0
        shift = 0
        while True:
            if self._last_read_count - self._offset < 1:
                self._fill_buffer(1)

            byte = self._buffer[self._offset]
            self._offset += 1
            result |= (byte & 0x7F) << shift
            if byte < 0x80:
                return result
            shift += 7

    def zigzag_decode(self, value: int) -> int:
        return (value >> 1) ^ -(value & 1)

    def read_signed_varint(self) -> int:
        return self.zigzag_decode(self.read_unsigned_varint())

    def read_view(self, count: int) -> memoryview:
        if count <= (self._last_read_count - self._offset):
            res = self._view[self._offset : self._offset + count]
            self._offset += count
            return res

        if count > len(self._buffer):
            local_buf = bytearray(count)
            local_view = memoryview(local_buf)
            remaining = self._last_read_count - self._offset
            local_view[:remaining] = self._view[self._offset : self._last_read_count]
            self._offset = self._last_read_count
            if self._stream.readinto(local_view[remaining:]) < count - remaining:
                raise EOFError("Unexpected EOF")
            return local_view

        self._fill_buffer(count)
        result = self._view[self._offset : self._offset + count]
        self._offset += count
        return result

    def read_bytearray(self, count: int) -> bytearray:
        if count <= (self._last_read_count - self._offset):
            res = bytearray(self._view[self._offset : self._offset + count])
            self._offset += count
            return res

        if count > len(self._buffer):
            local_buf = bytearray(count)
            local_view = memoryview(local_buf)
            remaining = self._last_read_count - self._offset
            local_view[:remaining] = self._view[self._offset : self._last_read_count]
            self._offset = self._last_read_count
            if self._stream.readinto(local_view[remaining:]) < count - remaining:
                raise EOFError("Unexpected EOF")
            return local_buf

        self._fill_buffer(count)
        result = self._view[self._offset : self._offset + count]
        self._offset += count
        return bytearray(result)

    def _fill_buffer(self, min_count: int = 0) -> None:
        remaining = self._last_read_count - self._offset
        if remaining > 0:
            remaining_view = memoryview(self._buffer)[self._offset:]
            self._buffer[:remaining] = remaining_view

        slice = memoryview(self._buffer)[remaining:]
        self._last_read_count = self._stream.readinto(slice)
        self._offset = 0
        if self._last_read_count == 0:
            self._at_end = True
        if min_count > 0 and (self._last_read_count + remaining) < min_count:
            raise EOFError("Unexpected EOF")



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

    #@abstractmethod
    def read(self, stream: CodedInputStream) -> T:
        raise NotImplementedError

    def is_trivially_serializable(self) -> bool:
        return False

class StructDescriptor(TypeDescriptor[T]):
    def __init__(self, dtype: npt.DTypeLike, format_string: str) -> None:
        super().__init__(dtype)
        self.struct = struct.Struct(format_string)

    def write(self, stream: CodedOutputStream, value: T) -> None:
        stream.write(self.struct, value)

    def read(self, stream: CodedInputStream) -> T:
        return cast(T, stream.read(self.struct)[0])

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

        stream.write_signed_varint(value)

    def read(self, stream: CodedInputStream) -> Int16:
        return cast(Int16, stream.read_signed_varint())

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

    def read(self, stream: CodedInputStream) -> UInt16:
        return cast(UInt16, stream.read_unsigned_varint())

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

        stream.write_signed_varint(value)

    def read(self, stream: CodedInputStream) -> Int32:
        return cast(Int32, stream.read_signed_varint())

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

    def read(self, stream: CodedInputStream) -> UInt32:
        return cast(UInt32, stream.read_unsigned_varint())


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

    def read(self, stream: CodedInputStream) -> Int64:
        return cast(Int64, stream.read_signed_varint())

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

    def read(self, stream: CodedInputStream) -> UInt64:
        return cast(UInt64, stream.read_unsigned_varint())

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

    def read(self, stream: CodedInputStream) -> str:
        length = stream.read_unsigned_varint()
        view = stream.read_view(length)
        return str(view, "utf-8")

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

    def read(self, stream: CodedInputStream) -> Date:
        days_since_epoch = stream.read_signed_varint()
        return datetime.datetime.fromordinal(days_since_epoch + EPOCH_ORDINAL_DAYS)


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

    def read(self, stream: CodedInputStream) -> Time:
        nanoseconds_since_midnight = stream.read_signed_varint()
        hours, r = divmod(nanoseconds_since_midnight, 3_600_000_000_000)
        minutes, r = divmod(r, 60_000_000_000)
        seconds, r = divmod(r, 1_000_000_000)
        microseconds = r // 1000
        return datetime.time(hours, minutes, seconds, microseconds)

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

    def read(self, stream: CodedInputStream) -> DateTime:
        nanoseconds_since_epoch = stream.read_signed_varint()
        return EPOCH_DATETIME + datetime.timedelta(microseconds=nanoseconds_since_epoch / 1000)


datetime_descriptor = DateTimeDescriptor()

class NoneDescriptor(TypeDescriptor[None]):
    def __init__(self) -> None:
        super().__init__(np.object_)

    def write(self, stream: CodedOutputStream, value: None) -> None:
        pass

    def read(self, stream: CodedInputStream) -> None:
        return None

none_descriptor = NoneDescriptor()

TEnum = TypeVar("TEnum", bound=Enum)
class EnumDescriptor(Generic[TEnum], TypeDescriptor[TEnum]):
    def __init__(self, integer_descriptor: TypeDescriptor[TEnum], enum_type: type) -> None:
        super().__init__(integer_descriptor.overall_dtype())
        self._integer_descriptor = integer_descriptor
        self._enum_type = enum_type

    def write(self, stream: CodedOutputStream, value: TEnum) -> None:
        self._integer_descriptor.write(stream, value.value)

    def read(self, stream: CodedInputStream) -> TEnum:
        return self._enum_type(self._integer_descriptor.read(stream))

    def is_trivially_serializable(self) -> bool:
        return self._integer_descriptor.is_trivially_serializable()


class OptionalDescriptor(TypeDescriptor[Optional[T]]):
    def __init__(self, element_descriptor: TypeDescriptor[T]) -> None:
        super().__init__(np.dtype([("has_value", np.bool_), ("value", element_descriptor.overall_dtype())]))
        self.element_descriptor = element_descriptor

    def write(self, stream: CodedOutputStream, value: Optional[T]) -> None:
        if value is None:
            stream.write_byte(0)
        else:
            stream.write_byte(1)
            self.element_descriptor.write(stream, value)

    def read(self, stream: CodedInputStream) -> Optional[T]:
        has_value = stream.read_byte()
        if has_value == 0:
            return None
        else:
            return self.element_descriptor.read(stream)


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

    def read(self, stream: CodedInputStream) -> Any:
        case_index = stream.read_byte()
        _, case_descriptor = self.cases[case_index]
        return case_descriptor.read(stream)


class StreamDescriptor(TypeDescriptor[Iterable[T]]):
    def __init__(self, element_descriptor: TypeDescriptor[T]) -> None:
        super().__init__(np.object_)
        self.element_descriptor = element_descriptor

    def write(self, stream: CodedOutputStream, value: Iterable[T]) -> None:
        for element in value:
            stream.write_byte(1)
            self.element_descriptor.write(stream, element)

        stream.write_byte(0)

    def read(self, stream: CodedInputStream) -> Iterable[T]:
        while stream.read_byte():
            yield self.element_descriptor.read(stream)


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

    def read(self, stream: CodedInputStream) -> list[T]:
        return [
            self.element_descriptor.read(stream)
            for _ in range(self.length)
        ]

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

    def read(self, stream: CodedInputStream) -> list[T]:
        length = stream.read_unsigned_varint()
        return [
            self.element_descriptor.read(stream)
            for _ in range(length)
        ]


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

    def read(self, stream: CodedInputStream) -> dict[TKey, TValue]:
        length = stream.read_unsigned_varint()
        return {
            self.key_descriptor.read(stream): self.value_descriptor.read(stream)
            for _ in range(length)
        }

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


def read_fixed_int32(stream: CodedInputStream) -> int:
    return stream.read(int32_struct)[0]
