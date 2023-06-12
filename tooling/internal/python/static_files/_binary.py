from typing import BinaryIO, TypeVar, Protocol, Generic, Any, Optional
from collections.abc import Callable
from functools import partial
import struct
import numpy as np

class _CodedOutputStream:
    def __init__(self, stream: BinaryIO | str,*, buffer_size=65536) -> None:
        if isinstance(stream, str):
            self._stream = open(stream, "wb")
            self._owns_stream = True
        else:
            self._stream  = stream
            self._owns_stream = False

        self._buffer = bytearray(buffer_size)
        self._view = memoryview(self._buffer)

    def __enter__(self) -> '_CodedOutputStream':
        return self

    def __exit__(self, exc_type, exc_value, traceback) -> None:
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
class Writer(Protocol, Generic[T]):
    def __call__(self, stream: _CodedOutputStream, value : T) -> None: ...

TOuter = TypeVar('TOuter', contravariant=True)
TElement = TypeVar('TElement', covariant=True)
class OuterWriter(Protocol, Generic[TElement, TOuter]):
    def __call__(self, write_element: Writer[TElement], stream: _CodedOutputStream, value : TOuter) -> None: ...

float32_struct = struct.Struct('<f')
assert float32_struct.size == 4
def write_float32(stream: _CodedOutputStream, value: float) -> None:
    stream.write(float32_struct, value)

float64_struct = struct.Struct('<d')
assert float64_struct.size == 8
def write_float64(stream: _CodedOutputStream, value: float) -> None:
    stream.write(float64_struct, value)

def write_unsigned_varint(stream: _CodedOutputStream, value: int) -> None:
    stream.write_unsigned_varint(value)


TElement = TypeVar('TElement')
def write_optional(write_element: Writer[TElement], stream: _CodedOutputStream, value: TElement | None) -> None:
    if value is None:
        stream.write_byte(0)
    else:
        stream.write_byte(1)
        write_element(stream, value)

def write_dynamic_vector(write_element: Writer[TElement], stream: _CodedOutputStream, value: list[TElement]) -> None:
    stream.write_unsigned_varint(len(value))
    for element in value:
        write_element(stream, element)

def write_fixed_vector(length: int) -> OuterWriter[TElement, list[TElement]]:
    def write_fixed_vector_inner(write_element: Writer[TElement], stream: _CodedOutputStream, value: list[TElement]) -> None:
        assert len(value) == length
        for element in value:
            write_element(stream, element)
    return write_fixed_vector_inner

class FixedVectorWriter(Generic[TElement]):
    def __init__(self, length: int) -> None:
        self.length = length

    def __call__(self, write_element: Writer[TElement], stream: _CodedOutputStream, value: list[TElement]) -> None:
        assert len(value) == self.length
        for element in value:
            write_element(stream, element)

complex32_struct = struct.Struct('<ff')
assert complex32_struct.size == 8
def write_complex32(stream: _CodedOutputStream, value: complex) -> None:
    stream.write(complex32_struct, value.real, value.imag)

complex64_struct = struct.Struct('<dd')
assert complex64_struct.size == 16
def write_complex64(stream: _CodedOutputStream, value: complex) -> None:
    stream.write(complex64_struct, value.real, value.imag)

def write_string(stream: _CodedOutputStream, value: str) -> None:
    stream.write_string(value)

TKey = TypeVar('TKey')
TValue = TypeVar('TValue')
class MapWriter(Generic[TKey, TValue]):
    def __init__(self, write_key: Writer[TKey], write_value: Writer[TValue]) -> None:
        self.write_key = write_key
        self.write_value = write_value

    def __call__(self, stream: _CodedOutputStream, value: dict[TKey, TValue]) -> None:
        stream.write_unsigned_varint(len(value))
        for k, v in value.items():
            self.write_key(stream, k)
            self.write_value(stream, v)

class DynamicNDArrayWriter(Generic[TElement]):
    def __init__(self, dtype: np.dtype,  write_element: Writer[TElement], trivially_serializable: bool) -> None:
        self.dtype = dtype
        self.write_element = write_element
        self.trivially_serializable = trivially_serializable

    def __call__(self, stream: _CodedOutputStream, value: np.ndarray) -> None:
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


def combine_writers(inner: Writer[TElement], outer: OuterWriter[TElement, TOuter], ) -> Writer[TOuter]:
    return partial(outer, inner)

x = write_fixed_vector(2)

w2 = combine_writers(write_float32, FixedVectorWriter(2))

w3 = MapWriter(write_string, write_float32)


stream = _CodedOutputStream(open('test.bin', 'wb'))


w2(stream, [1.0, 2.0])

# optional [ record ]

# write_optional(stream, value, write_record)
