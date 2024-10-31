# This file was generated by the "yardl" tool. DO NOT EDIT.

# pyright: reportUnusedClass=false
# pyright: reportUnusedImport=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false
# pyright: reportUnknownVariableType=false

import collections.abc
import io
import typing

import numpy as np
import numpy.typing as npt

from .types import *

from .protocols import *
from . import _binary
from . import yardl_types as yardl

class BinaryMyProtocolWriter(_binary.BinaryProtocolWriter, MyProtocolWriterBase):
    """Binary writer for the MyProtocol protocol."""

    def __init__(self, stream: typing.Union[typing.BinaryIO, str]) -> None:
        MyProtocolWriterBase.__init__(self)
        _binary.BinaryProtocolWriter.__init__(self, stream, MyProtocolWriterBase.schema)

    def _write_header(self, value: Header) -> None:
        HeaderSerializer().write(self._stream, value)

    def _write_samples(self, value: collections.abc.Iterable[Sample]) -> None:
        _binary.StreamSerializer(SampleSerializer()).write(self._stream, value)


class BinaryMyProtocolReader(_binary.BinaryProtocolReader, MyProtocolReaderBase):
    """Binary writer for the MyProtocol protocol."""

    def __init__(self, stream: typing.Union[io.BufferedReader, io.BytesIO, typing.BinaryIO, str]) -> None:
        MyProtocolReaderBase.__init__(self)
        _binary.BinaryProtocolReader.__init__(self, stream, MyProtocolReaderBase.schema)

    def _read_header(self) -> Header:
        return HeaderSerializer().read(self._stream)

    def _read_samples(self) -> collections.abc.Iterable[Sample]:
        return _binary.StreamSerializer(SampleSerializer()).read(self._stream)


class BinaryMyProtocolIndexedWriter(_binary.BinaryProtocolIndexedWriter, MyProtocolWriterBase):
    """Binary indexed writer for the MyProtocol protocol."""

    def __init__(self, stream: typing.Union[typing.BinaryIO, str]) -> None:
        MyProtocolWriterBase.__init__(self)
        _binary.BinaryProtocolIndexedWriter.__init__(self, stream, MyProtocolWriterBase.schema)

    def _write_header(self, value: Header) -> None:
        pos = self._stream.pos()
        self._index.set_step_offset("Header", pos)
        HeaderSerializer().write(self._stream, value)

    def _write_samples(self, value: collections.abc.Iterable[Sample]) -> None:
        pos = self._stream.pos()
        self._index.set_step_offset("Samples", pos)
        offsets, num_blocks = _binary.StreamSerializer(SampleSerializer()).write_and_save_offsets(self._stream, value)
        self._index.add_stream_offsets("Samples", offsets, num_blocks)


class BinaryMyProtocolIndexedReader(_binary.BinaryProtocolIndexedReader, MyProtocolIndexedReaderBase):
    """Binary indexed writer for the MyProtocol protocol."""

    def __init__(self, stream: typing.Union[io.BufferedReader, io.BytesIO, typing.BinaryIO, str]) -> None:
        MyProtocolIndexedReaderBase.__init__(self)
        _binary.BinaryProtocolIndexedReader.__init__(self, stream, MyProtocolIndexedReaderBase.schema)

    def _read_header(self) -> Header:
        pos = self._index.get_step_offset("Header")
        self._stream.seek(pos)
        return HeaderSerializer().read(self._stream)

    def _read_samples(self, idx: int) -> collections.abc.Iterable[Sample]: # pyright: ignore [reportIncompatibleMethodOverride]
        offset, remaining = self._index.find_stream_item("Samples", idx)
        self._stream.seek(offset)
        return _binary.StreamSerializer(SampleSerializer()).read_mid_stream(self._stream, remaining)

    def _count_samples(self) -> int:
        return self._index.get_stream_size("Samples")

class HeaderSerializer(_binary.RecordSerializer[Header]):
    def __init__(self) -> None:
        super().__init__([("subject", _binary.string_serializer)])

    def write(self, stream: _binary.CodedOutputStream, value: Header) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.subject)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['subject'])

    def read(self, stream: _binary.CodedInputStream) -> Header:
        field_values = self._read(stream)
        return Header(subject=field_values[0])


class SampleSerializer(_binary.RecordSerializer[Sample]):
    def __init__(self) -> None:
        super().__init__([("id", _binary.uint32_serializer), ("data", _binary.NDArraySerializer(_binary.int32_serializer, 1))])

    def write(self, stream: _binary.CodedOutputStream, value: Sample) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.id, value.data)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['id'], value['data'])

    def read(self, stream: _binary.CodedInputStream) -> Sample:
        field_values = self._read(stream)
        return Sample(id=field_values[0], data=field_values[1])

