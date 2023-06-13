from typing import BinaryIO, Iterable, TypeVar, Protocol, Generic, Any, Optional, Tuple
from collections.abc import Callable
from abc import ABC
from functools import partial
import struct
import sys
import numpy as np

MAGIC_BYTES = b"yardl"
CURRENT_BINARY_FORMAT_VERSION = 1

INT32_MIN = np.iinfo(np.int32).min
INT32_MAX = np.iinfo(np.int32).max

UINT32_MAX = np.iinfo(np.uint32).max


class BinaryProtocolWriter(ABC):
    def __init__(self, stream: BinaryIO | str, schema: str) -> None:
        self._stream = CodedOutputStream(stream)
        self._stream.write_bytes(MAGIC_BYTES)
        write_fixed_int32(self._stream, CURRENT_BINARY_FORMAT_VERSION)
        self._stream.write_string(schema)

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_value, traceback) -> None:
        self.close()

    def close(self) -> None:
        try:
            self._close()
        finally:
            self._stream.close()

    def _close(self) -> None:
        pass


class CodedOutputStream:
    def __init__(self, stream: BinaryIO | str,*, buffer_size=65536) -> None:
        if isinstance(stream, str):
            self._stream = open(stream, "wb")
            self._owns_stream = True
        else:
            self._stream  = stream
            self._owns_stream = False

        self._buffer = bytearray(buffer_size)
        self._view = memoryview(self._buffer)

    def __enter__(self) -> 'CodedOutputStream':
        return self

    def __exit__(self, exc_type, exc_value, traceback) -> None:
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

    def write(self, formatter: struct.Struct, *args) -> None:
        if len(self._view) < formatter.size:
            self.flush()

        formatter.pack_into(self._view, 0, *args)
        self._view = self._view[formatter.size:]

    def write_bytes(self, data: bytes | bytearray) -> None:
        if len(data) > len(self._view):
            self.flush()
            self._stream.write(data)
        else:
            self._view[:len(data)] = data
            self._view = self._view[len(data):]

    def write_bytes_directly(self, data: bytes | bytearray | memoryview) -> None:
        self.flush()
        self._stream.write(data)

    def write_byte(self, value: int) -> None:
        self._view[0] = value
        self._view = self._view[1:]

    def write_unsigned_varint(self, value: int) -> None:
        if len(self._view) < 10:
            self.flush()

        while True:
            if value < 0x80:
                self.write_byte(value)
                return

            self.write_byte((value & 0x7F) | 0x80)
            value >>= 7

    def zigzag_encode(self, value: int) -> int:
        return (value << 1) ^ (value >> 63)

    def write_signed_varint(self, value: int) -> None:
        self.write_unsigned_varint(self.zigzag_encode(value))

    def write_string(self, value: str) -> None:
        self.write_bytes(value.encode('utf-8'))


T = TypeVar('T', contravariant=True)
Writer = Callable[[CodedOutputStream, T], None]

TOuter = TypeVar('TOuter', contravariant=True)
TElement = TypeVar('TElement', covariant=True)
class OuterWriter(Protocol, Generic[TElement, TOuter]):
    def __call__(self, write_element: Writer[TElement], stream: CodedOutputStream, value : TOuter) -> None: ...

int32_struct = struct.Struct('<i')
assert int32_struct.size == 4
def write_fixed_int32(stream: CodedOutputStream, value: int) -> None:
    stream.write(int32_struct, value)

float32_struct = struct.Struct('<f')
assert float32_struct.size == 4
def write_float32(stream: CodedOutputStream, value: float) -> None:
    stream.write(float32_struct, value)

float64_struct = struct.Struct('<d')
assert float64_struct.size == 8
def write_float64(stream: CodedOutputStream, value: float) -> None:
    stream.write(float64_struct, value)

def write_int32(stream: CodedOutputStream, value: int) -> None:
    assert INT32_MIN <= value <= INT32_MAX
    stream.write_signed_varint(value)

def write_uint32(stream: CodedOutputStream, value: int) -> None:
    assert 0 <= value <= INT32_MAX
    stream.write_unsigned_varint(value)

complex32_struct = struct.Struct('<ff')
assert complex32_struct.size == 8
def write_complex32(stream: CodedOutputStream, value: complex) -> None:
    stream.write(complex32_struct, value.real, value.imag)

complex64_struct = struct.Struct('<dd')
assert complex64_struct.size == 16
def write_complex64(stream: CodedOutputStream, value: complex) -> None:
    stream.write(complex64_struct, value.real, value.imag)

def write_string(stream: CodedOutputStream, value: str) -> None:
    stream.write_string(value)

def write_none(stream: CodedOutputStream, value: None) -> None:
    pass


class OptionalWriter(Generic[TElement]):
    def __init__(self, write_element: Writer[TElement]) -> None:
        self.write_element = write_element

    def __call__(self, stream: CodedOutputStream, value: Optional[TElement]) -> None:
        if value is None:
            stream.write_byte(0)
        else:
            stream.write_byte(1)
            self.write_element(stream, value)

class UnionWriter:
    def __init__(self, cases: list[Tuple[type | None, Writer]]) -> None:
        self.cases = cases

    def __call__(self, stream: CodedOutputStream, value: Any) -> None:
        for i, (case_type, case_writer) in enumerate(self.cases):
            if case_type is None:
                if value is None:
                    stream.write_byte(i)
                    return
            elif isinstance(value, case_type):
                stream.write_byte(i)
                case_writer(stream, value)
                return

        raise ValueError(f'Incorrect union type {type(value)}')


class StreamWriter(Generic[TElement]):
    def __init__(self, write_element: Writer[TElement]) -> None:
        self.write_element = write_element

    def __call__(self, stream: CodedOutputStream, value: Iterable[TElement]) -> None:
        for element in value:
            stream.write_byte(1)
            self.write_element(stream, element)

class FixedVectorWriter(Generic[TElement]):
    def __init__(self, length: int, write_element: Writer[TElement]) -> None:
        self.length = length
        self.write_element = write_element

    def __call__(self, stream: CodedOutputStream, value: list[TElement]) -> None:
        assert len(value) == self.length
        for element in value:
            self.write_element(stream, element)

class DynamicVectorWriter(Generic[TElement]):
    def __init__(self, write_element: Writer[TElement]) -> None:
        self.write_element = write_element

    def __call__(self, stream: CodedOutputStream, value: list[TElement]) -> None:
        stream.write_unsigned_varint(len(value))
        for element in value:
            self.write_element(stream, element)

TKey = TypeVar('TKey')
TValue = TypeVar('TValue')
class MapWriter(Generic[TKey, TValue]):
    def __init__(self, write_key: Writer[TKey], write_value: Writer[TValue]) -> None:
        self.write_key = write_key
        self.write_value = write_value

    def __call__(self, stream: CodedOutputStream, value: dict[TKey, TValue]) -> None:
        stream.write_unsigned_varint(len(value))
        for k, v in value.items():
            self.write_key(stream, k)
            self.write_value(stream, v)

class DynamicNDArrayWriter(Generic[TElement]):
    def __init__(self, write_element: Writer[TElement], dtype: np.dtype, trivially_serializable: bool) -> None:
        self.dtype = dtype
        self.write_element = write_element
        self.trivially_serializable = trivially_serializable

    def __call__(self, stream: CodedOutputStream, value: np.ndarray) -> None:
        assert value.dtype == self.dtype, "dtype mismatch"
        stream.write_unsigned_varint(value.ndim)
        for dim in value.shape:
            stream.write_unsigned_varint(dim)

        if self.trivially_serializable and value.flags.c_contiguous:
            stream.flush()
            stream._stream.write(value.data)
        else:
            for element in value.flat:
                self.write_element(stream, element)
