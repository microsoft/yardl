import datetime
from enum import Enum
from io import BufferedReader
from types import TracebackType
from typing import BinaryIO, Iterable, TypeVar, Generic, Any, Optional, Tuple, cast
from abc import ABC, abstractmethod
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
        string_serializer.write(self._stream, schema)

    def close(self) -> None:
        self._stream.close()

class BinaryProtocolReader(ABC):
    def __init__(self, stream: BufferedReader | str, expected_schema: str) -> None:
        self._stream = CodedInputStream(stream)
        magic_bytes = self._stream.read_view(len(MAGIC_BYTES))
        if magic_bytes != MAGIC_BYTES:
            raise RuntimeError("Invalid magic bytes")

        version = read_fixed_int32(self._stream)
        if version != CURRENT_BINARY_FORMAT_VERSION:
            raise RuntimeError("Invalid binary format version")

        self._schema = string_serializer.read(self._stream, Types.NONE)
        if self._schema != expected_schema:
            raise RuntimeError("Invalid schema")

    def close(self) -> None:
        self._stream.close()



class CodedOutputStream:
    def __init__(self, stream: BinaryIO | str, *, buffer_size: int = 65536) -> None:
        if isinstance(stream, str):
            self._stream = cast(BinaryIO, open(stream, "wb"))
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
    def __init__(
        self, stream: BufferedReader | str, *, buffer_size: int = 65536
    ) -> None:
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
            remaining_view = memoryview(self._buffer)[self._offset :]
            self._buffer[:remaining] = remaining_view

        slice = memoryview(self._buffer)[remaining:]
        self._last_read_count = self._stream.readinto(slice)
        self._offset = 0
        if self._last_read_count == 0:
            self._at_end = True
        if min_count > 0 and (self._last_read_count + remaining) < min_count:
            raise EOFError("Unexpected EOF")


T = TypeVar("T")


class TypeSerializer(Generic[T], ABC):
    def __init__(self, dtype: npt.DTypeLike) -> None:
        self._dtype: np.dtype[Any] = np.dtype(dtype)

    def overall_dtype(self) -> np.dtype[Any]:
        return self._dtype

    def struct_format_str(self) -> str | None:
        return None

    @abstractmethod
    def write(self, stream: CodedOutputStream, value: T) -> None:
        raise NotImplementedError

    @abstractmethod
    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> T:
        raise NotImplementedError

    def is_trivially_serializable(self) -> bool:
        return False


class StructSerializer(TypeSerializer[T]):
    def __init__(self, dtype: npt.DTypeLike, format_string: str) -> None:
        super().__init__(dtype)
        self._struct = struct.Struct(format_string)

    def write(self, stream: CodedOutputStream, value: T) -> None:
        stream.write(self._struct, value)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> T:
        return cast(T, stream.read(self._struct)[0])

    def struct_format_str(self) -> str:
        return self._struct.format


class BoolSerializer(StructSerializer[Bool]):
    def __init__(self) -> None:
        super().__init__(np.bool_, "<?")

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Bool:
        python_value = super().read(stream, read_as_numpy)
        return np.bool_(python_value) if Types.BOOL in read_as_numpy else python_value


bool_serializer = BoolSerializer()


class Int8Serializer(StructSerializer[Int8]):
    def __init__(self) -> None:
        super().__init__(np.int8, "<b")

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Int8:
        python_value = super().read(stream, read_as_numpy)
        return np.int8(python_value) if Types.INT8 in read_as_numpy else python_value

    def is_trivially_serializable(self) -> bool:
        return True


int8_serializer = Int8Serializer()


class UInt8Serializer(StructSerializer[UInt8]):
    def __init__(self) -> None:
        super().__init__(np.uint8, "<B")

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> UInt8:
        python_value = super().read(stream, read_as_numpy)
        return np.uint8(python_value) if Types.UINT8 in read_as_numpy else python_value

    def is_trivially_serializable(self) -> bool:
        return True


uint8_serializer = UInt8Serializer()


class Int16Serializer(TypeSerializer[Int16]):
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

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Int16:
        python_value = stream.read_signed_varint()
        return np.int16(python_value) if Types.INT16 in read_as_numpy else python_value


int16_serializer = Int16Serializer()


class UInt16Serializer(TypeSerializer[UInt16]):
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

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> UInt16:
        python_value = stream.read_unsigned_varint()
        return (
            np.uint16(python_value) if Types.UINT16 in read_as_numpy else python_value
        )


uint16_serializer = UInt16Serializer()


class Int32Serializer(TypeSerializer[Int32]):
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

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Int32:
        python_val = stream.read_signed_varint()
        return np.int32(python_val) if Types.INT32 in read_as_numpy else python_val


int32_serializer = Int32Serializer()


class UInt32Serializer(TypeSerializer[UInt32]):
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

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> UInt32:
        python_value = stream.read_unsigned_varint()
        return (
            np.uint32(python_value) if Types.UINT32 in read_as_numpy else python_value
        )


uint32_serializer = UInt32Serializer()


class Int64Serializer(TypeSerializer[Int64]):
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

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Int64:
        python_value = stream.read_signed_varint()
        return np.int64(python_value) if Types.INT64 in read_as_numpy else python_value


int64_serializer = Int64Serializer()


class UInt64Serializer(TypeSerializer[UInt64]):
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

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> UInt64:
        python_value = stream.read_unsigned_varint()
        return (
            np.uint64(python_value) if Types.UINT64 in read_as_numpy else python_value
        )


uint64_serializer = UInt64Serializer()


class SizeSerializer(TypeSerializer[Size]):
    def __init__(self) -> None:
        super().__init__(np.uint64)

    def write(self, stream: CodedOutputStream, value: Size) -> None:
        if isinstance(value, int):
            if value < 0 or value > UINT64_MAX:
                raise ValueError(
                    f"Value {value} is outside the range of an unsigned 64-bit integer"
                )
        elif not isinstance(value, cast(type, np.uint64)):
            raise ValueError(f"Value in not an unsigned 64-bit integer: {value}")

        stream.write_unsigned_varint(value)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Size:
        python_value = stream.read_unsigned_varint()
        return np.uint64(python_value) if Types.SIZE in read_as_numpy else python_value


size_serializer = SizeSerializer()


class Float32Serializer(StructSerializer[Float32]):
    def __init__(self) -> None:
        super().__init__(np.float32, "<f")

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Float32:
        python_value = super().read(stream, read_as_numpy)
        return (
            np.float32(python_value) if Types.FLOAT32 in read_as_numpy else python_value
        )

    def is_trivially_serializable(self) -> bool:
        return True


float32_serializer = Float32Serializer()


class Float64Serializer(StructSerializer[Float64]):
    def __init__(self) -> None:
        super().__init__(np.float64, "<d")

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Float64:
        python_value = super().read(stream, read_as_numpy)
        return (
            np.float64(python_value) if Types.FLOAT64 in read_as_numpy else python_value
        )

    def is_trivially_serializable(self) -> bool:
        return True


float64_serializer = Float64Serializer()


class Complex32Serializer(StructSerializer[ComplexFloat]):
    def __init__(self) -> None:
        super().__init__(np.complex64, "<ff")

    def write(self, stream: CodedOutputStream, value: ComplexFloat) -> None:
        stream.write(self._struct, value.real, value.imag)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> ComplexFloat:
        real, imag = stream.read(self._struct)
        python_value = complex(real, imag)
        return (
            np.complex64(python_value)
            if Types.COMPLEX_FLOAT32 in read_as_numpy
            else python_value
        )

    def is_trivially_serializable(self) -> bool:
        return True


complexfloat32_serializer = Complex32Serializer()


class Complex64Serializer(StructSerializer[ComplexDouble]):
    def __init__(self) -> None:
        super().__init__(np.complex128, "<dd")

    def write(self, stream: CodedOutputStream, value: ComplexDouble) -> None:
        stream.write(self._struct, value.real, value.imag)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> ComplexDouble:
        real, imag = stream.read(self._struct)
        python_value = complex(real, imag)
        return (
            np.complex128(python_value)
            if Types.COMPLEX_FLOAT64 in read_as_numpy
            else python_value
        )

    def is_trivially_serializable(self) -> bool:
        return True


complexfloat64_serializer = Complex64Serializer()


class StringSerializer(TypeSerializer[str]):
    def __init__(self) -> None:
        super().__init__(np.object_)

    def write(self, stream: CodedOutputStream, value: str) -> None:
        b = value.encode("utf-8")
        stream.write_unsigned_varint(len(b))
        stream.write_bytes(b)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> str:
        length = stream.read_unsigned_varint()
        view = stream.read_view(length)
        return str(view, "utf-8")


string_serializer = StringSerializer()

EPOCH_ORDINAL_DAYS = datetime.date(1970, 1, 1).toordinal()
DATETIME_DAYS_DTYPE = np.dtype("datetime64[D]")


class DateSerializer(TypeSerializer[Date]):
    def __init__(self) -> None:
        super().__init__(np.datetime64)

    def write(self, stream: CodedOutputStream, value: Date) -> None:
        if isinstance(value, datetime.date):
            stream.write_signed_varint(value.toordinal() - EPOCH_ORDINAL_DAYS)
        else:
            if not isinstance(value, np.datetime64):
                raise ValueError(
                    f"Expected datetime.date or numpy.datetime64, got {type(value)}"
                )

            if value.dtype == DATETIME_DAYS_DTYPE:
                stream.write_signed_varint(value.astype(np.int32))
            else:
                stream.write_signed_varint(
                    value.astype(DATETIME_DAYS_DTYPE).astype(np.int32)
                )

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Date:
        days_since_epoch = stream.read_signed_varint()
        if Types.DATE in read_as_numpy:
            return np.datetime64(days_since_epoch, "D")

        return datetime.date.fromordinal(days_since_epoch + EPOCH_ORDINAL_DAYS)


date_serializer = DateSerializer()

TIMEDELTA_NANOSECONDS_DTYPE = np.dtype("timedelta64[ns]")


class TimeSerializer(TypeSerializer[Time]):
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
                raise ValueError(
                    f"Expected a datetime.time or np.timedelta64, got {type(value)}"
                )

            if value.dtype == TIMEDELTA_NANOSECONDS_DTYPE:
                stream.write_signed_varint(value.astype(np.int64))
            else:
                stream.write_signed_varint(
                    value.astype(DATETIME_NANOSECONDS_DTYPE).astype(np.int64)
                )

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Time:
        nanoseconds_since_midnight = stream.read_signed_varint()
        if Types.TIME in read_as_numpy:
            return np.timedelta64(nanoseconds_since_midnight, "ns")

        hours, r = divmod(nanoseconds_since_midnight, 3_600_000_000_000)
        minutes, r = divmod(r, 60_000_000_000)
        seconds, r = divmod(r, 1_000_000_000)
        microseconds = r // 1000
        return datetime.time(hours, minutes, seconds, microseconds)


time_serializer = TimeSerializer()

DATETIME_NANOSECONDS_DTYPE = np.dtype("datetime64[ns]")
EPOCH_DATETIME = datetime.datetime.utcfromtimestamp(0)


class DateTimeSerializer(TypeSerializer[DateTime]):
    def __init__(self) -> None:
        super().__init__(np.datetime64)

    def write(self, stream: CodedOutputStream, value: DateTime) -> None:
        if isinstance(value, datetime.datetime):
            delta = value - EPOCH_DATETIME
            nanoseconds_since_epoch = int(delta.total_seconds() * 1e9)
            stream.write_signed_varint(nanoseconds_since_epoch)
        else:
            if not isinstance(value, np.datetime64):
                raise ValueError(
                    f"Expected datetime.datetime or numpy.datetime64, got {type(value)}"
                )

            if value.dtype == DATETIME_NANOSECONDS_DTYPE:
                stream.write_signed_varint(value.astype(np.int64))
            else:
                stream.write_signed_varint(
                    value.astype(DATETIME_NANOSECONDS_DTYPE).astype(np.int64)
                )

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> DateTime:
        nanoseconds_since_epoch = stream.read_signed_varint()
        if Types.DATETIME in read_as_numpy:
            return np.datetime64(nanoseconds_since_epoch, "ns")

        return EPOCH_DATETIME + datetime.timedelta(
            microseconds=nanoseconds_since_epoch / 1000
        )


datetime_serializer = DateTimeSerializer()


class NoneSerializer(TypeSerializer[None]):
    def __init__(self) -> None:
        super().__init__(np.object_)

    def write(self, stream: CodedOutputStream, value: None) -> None:
        pass

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> None:
        return None


none_serializer = NoneSerializer()

TEnum = TypeVar("TEnum", bound=Enum)


class EnumSerializer(Generic[TEnum], TypeSerializer[TEnum]):
    def __init__(
        self, integer_serializer: TypeSerializer[TEnum], enum_type: type
    ) -> None:
        super().__init__(integer_serializer.overall_dtype())
        self._integer_serializer = integer_serializer
        self._enum_type = enum_type

    def write(self, stream: CodedOutputStream, value: TEnum) -> None:
        self._integer_serializer.write(stream, value.value)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> TEnum:
        return self._enum_type(self._integer_serializer.read(stream, read_as_numpy))

    def is_trivially_serializable(self) -> bool:
        return self._integer_serializer.is_trivially_serializable()


class OptionalSerializer(TypeSerializer[Optional[T]]):
    def __init__(self, element_serializer: TypeSerializer[T]) -> None:
        super().__init__(
            np.dtype(
                [("has_value", np.bool_), ("value", element_serializer.overall_dtype())]
            )
        )
        self.element_serializer = element_serializer

    def write(self, stream: CodedOutputStream, value: Optional[T]) -> None:
        if value is None:
            stream.write_byte(0)
        else:
            stream.write_byte(1)
            self.element_serializer.write(stream, value)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Optional[T]:
        has_value = stream.read_byte()
        if has_value == 0:
            return None
        else:
            return self.element_serializer.read(stream, read_as_numpy)


class UnionSerializer(TypeSerializer[Any]):
    def __init__(self, cases: list[Tuple[type, TypeSerializer[Any]]]) -> None:
        super().__init__(np.object_)
        self.cases = cases

    def write(self, stream: CodedOutputStream, value: Any) -> None:
        for i, (case_type, case_serializer) in enumerate(self.cases):
            if isinstance(value, case_type):
                stream.write_byte(i)
                case_serializer.write(stream, value)
                return

        raise ValueError(f"Incorrect union type {type(value)}")

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Any:
        case_index = stream.read_byte()
        _, case_serializer = self.cases[case_index]
        return case_serializer.read(stream, read_as_numpy)


class StreamSerializer(TypeSerializer[Iterable[T]]):
    def __init__(self, element_serializer: TypeSerializer[T]) -> None:
        super().__init__(np.object_)
        self.element_serializer = element_serializer

    def write(self, stream: CodedOutputStream, value: Iterable[T]) -> None:
        for element in value:
            stream.write_byte(1)
            self.element_serializer.write(stream, element)

        stream.write_byte(0)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> Iterable[T]:
        while stream.read_byte():
            yield self.element_serializer.read(stream, read_as_numpy)


class FixedVectorSerializer(TypeSerializer[list[T]]):
    def __init__(self, element_serializer: TypeSerializer[T], length: int) -> None:
        super().__init__(np.dtype((element_serializer.overall_dtype(), length)))
        self.element_serializer = element_serializer
        self.length = length

    def write(self, stream: CodedOutputStream, value: list[T]) -> None:
        if len(value) != self.length:
            raise ValueError(
                f"Expected a list of length {self.length}, got {len(value)}"
            )
        for element in value:
            self.element_serializer.write(stream, element)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> list[T]:
        return [
            self.element_serializer.read(stream, read_as_numpy)
            for _ in range(self.length)
        ]

    def is_trivially_serializable(self) -> bool:
        return self.element_serializer.is_trivially_serializable()


class VectorSerializer(TypeSerializer[list[T]]):
    def __init__(self, element_serializer: TypeSerializer[T]) -> None:
        super().__init__(np.object_)
        self.element_serializer = element_serializer

    def write(self, stream: CodedOutputStream, value: list[T]) -> None:
        stream.write_unsigned_varint(len(value))
        for element in value:
            self.element_serializer.write(stream, element)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> list[T]:
        length = stream.read_unsigned_varint()
        return [
            self.element_serializer.read(stream, read_as_numpy) for _ in range(length)
        ]


TKey = TypeVar("TKey")
TValue = TypeVar("TValue")


class MapSerializer(TypeSerializer[dict[TKey, TValue]]):
    def __init__(
        self,
        key_serializer: TypeSerializer[TKey],
        value_serializer: TypeSerializer[TValue],
    ) -> None:
        super().__init__(np.object_)
        self.key_serializer = key_serializer
        self.value_serializer = value_serializer

    def write(self, stream: CodedOutputStream, value: dict[TKey, TValue]) -> None:
        stream.write_unsigned_varint(len(value))
        for k, v in value.items():
            self.key_serializer.write(stream, k)
            self.value_serializer.write(stream, v)

    def read(
        self, stream: CodedInputStream, read_as_numpy: Types
    ) -> dict[TKey, TValue]:
        length = stream.read_unsigned_varint()
        return {
            self.key_serializer.read(stream, read_as_numpy): self.value_serializer.read(
                stream, read_as_numpy
            )
            for _ in range(length)
        }


class NDArraySerializerBase(Generic[T], TypeSerializer[npt.NDArray[Any]]):
    def __init__(
        self,
        overall_dtype: npt.DTypeLike,
        element_serializer: TypeSerializer[T],
        dtype: npt.DTypeLike,
    ) -> None:
        super().__init__(overall_dtype)
        self._array_dtype: np.dtype[Any] = (
            dtype if isinstance(dtype, np.dtype) else np.dtype(dtype)
        )
        self._element_serializer = element_serializer

    def _write_data(self, stream: CodedOutputStream, value: npt.NDArray[Any]) -> None:
        if value.dtype != self._array_dtype:
            # see if it's the same dtype but packed, not aligned
            packed_dtype = recfunctions.repack_fields(self._array_dtype, align=False, recurse=True)  # type: ignore
            if packed_dtype != value.dtype:
                raise ValueError(
                    f"Expected dtype {self._array_dtype} or {packed_dtype}, got {value.dtype}"
                )

        if self._is_current_array_trivially_serializable(value):
            stream.write_bytes_directly(value.data)
        else:
            for element in value.flat:
                self._element_serializer.write(stream, element)

    def _read_data(
        self, stream: CodedInputStream, shape: tuple[int, ...], read_as_numpy: Types
    ) -> npt.NDArray[Any]:
        flat_length = int(np.prod(shape))  # type: ignore

        if self._element_serializer.is_trivially_serializable():
            flat_byte_length = flat_length * self._array_dtype.itemsize
            byte_array = stream.read_bytearray(flat_byte_length)
            return np.frombuffer(byte_array, dtype=self._array_dtype).reshape(shape)

        result = np.empty((flat_length,), dtype=self._array_dtype)
        for i in range(flat_length):
            result[i] = self._element_serializer.read(stream, read_as_numpy)

        return result.reshape(shape)

    def _is_current_array_trivially_serializable(self, value: npt.NDArray[Any]) -> bool:
        return (
            self._element_serializer.is_trivially_serializable()
            and value.flags.c_contiguous
            and (
                self._array_dtype.fields is None
                or all(f != "" for f in self._array_dtype.fields)
            )
        )


class DynamicNDArraySerializer(NDArraySerializerBase[T]):
    def __init__(
        self,
        element_serializer: TypeSerializer[T],
    ) -> None:
        super().__init__(
            np.object_, element_serializer, element_serializer.overall_dtype()
        )

    def write(self, stream: CodedOutputStream, value: npt.NDArray[Any]) -> None:
        stream.write_unsigned_varint(value.ndim)
        for dim in value.shape:
            stream.write_unsigned_varint(dim)

        self._write_data(stream, value)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> npt.NDArray[Any]:
        ndims = stream.read_unsigned_varint()
        shape = tuple(stream.read_unsigned_varint() for _ in range(ndims))
        return self._read_data(stream, shape, read_as_numpy)


class NDArraySerializer(Generic[T], NDArraySerializerBase[T]):
    def __init__(
        self,
        element_serializer: TypeSerializer[T],
        ndims: int,
    ) -> None:
        super().__init__(
            np.object_, element_serializer, element_serializer.overall_dtype()
        )
        self.ndims = ndims

    def write(self, stream: CodedOutputStream, value: npt.NDArray[Any]) -> None:
        if value.ndim != self.ndims:
            raise ValueError(f"Expected {self.ndims} dimensions, got {value.ndim}")

        for dim in value.shape:
            stream.write_unsigned_varint(dim)

        self._write_data(stream, value)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> npt.NDArray[Any]:
        shape = tuple(stream.read_unsigned_varint() for _ in range(self.ndims))
        return self._read_data(stream, shape, read_as_numpy)


class FixedNDArraySerializer(Generic[T], NDArraySerializerBase[T]):
    def __init__(
        self,
        element_serializer: TypeSerializer[T],
        shape: tuple[int, ...],
    ) -> None:
        dtype = element_serializer.overall_dtype()
        super().__init__(np.dtype((dtype, shape)), element_serializer, dtype)
        self.shape = shape

    def write(self, stream: CodedOutputStream, value: npt.NDArray[Any]) -> None:
        if value.shape != self.shape:
            raise ValueError(f"Expected shape {self.shape}, got {value.shape}")

        self._write_data(stream, value)

    def read(self, stream: CodedInputStream, read_as_numpy: Types) -> npt.NDArray[Any]:
        return self._read_data(stream, self.shape, read_as_numpy)

    def is_trivially_serializable(self) -> bool:
        return self._element_serializer.is_trivially_serializable()


class RecordSerializer(TypeSerializer[T]):
    def __init__(
        self, field_serializers: list[Tuple[str, TypeSerializer[Any]]]
    ) -> None:
        super().__init__(
            np.dtype(
                [
                    (name, serializer.overall_dtype())
                    for name, serializer in field_serializers
                ],
                align=True,
            )
        )
        combined_format: str = "<"
        for _, field_serializer in field_serializers:
            fmt = field_serializer.struct_format_str()
            if fmt is None:
                combined_format = ""
                break
            else:
                combined_format += fmt[1:] if fmt[0] == "<" else fmt

        if combined_format == "":
            self._struct : struct.Struct | None = None
        else:
            self._struct = struct.Struct(combined_format)

        self._field_serializers = field_serializers

    def is_trivially_serializable(self) -> bool:
        return all(
            serializer.is_trivially_serializable()
            for _, serializer in self._field_serializers
        )

    def _write(self, stream: CodedOutputStream, *values: Any) -> None:
        if self._struct:
            stream.write(self._struct, *values)
        else:
            for i, (_, serializer) in enumerate(self._field_serializers):
                serializer.write(stream, values[i])

    def _read(self, stream: CodedInputStream, read_as_numpy: Types) -> tuple[Any, ...]:
        if self._struct:
            return stream.read(self._struct)
        else:
            return tuple(
                serializer.read(stream, read_as_numpy)
                for _, serializer in self._field_serializers
            )


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
