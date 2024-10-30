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

    def _write_tree(self, value: BinaryTree) -> None:
        BinaryTreeSerializer().write(self._stream, value)

    def _write_ptree(self, value: BinaryTree__) -> None:
        _binary.RecursiveSerializer(lambda *args, **kwargs : BinaryTreeSerializer()).write(self._stream, value)

    def _write_list(self, value: typing.Optional[LinkedList[str]]) -> None:
        _binary.OptionalSerializer(LinkedListSerializer(_binary.string_serializer)).write(self._stream, value)

    def _write_cwd(self, value: collections.abc.Iterable[DirectoryEntry]) -> None:
        _binary.StreamSerializer(_binary.UnionSerializer(DirectoryEntry, [(DirectoryEntry.File, FileSerializer()), (DirectoryEntry.Directory, _binary.RecursiveSerializer(lambda *args, **kwargs : DirectorySerializer()))])).write(self._stream, value)


class BinaryMyProtocolReader(_binary.BinaryProtocolReader, MyProtocolReaderBase):
    """Binary writer for the MyProtocol protocol."""


    def __init__(self, stream: typing.Union[io.BufferedReader, io.BytesIO, typing.BinaryIO, str]) -> None:
        MyProtocolReaderBase.__init__(self)
        _binary.BinaryProtocolReader.__init__(self, stream, MyProtocolReaderBase.schema)

    def _read_tree(self) -> BinaryTree:
        return BinaryTreeSerializer().read(self._stream)

    def _read_ptree(self) -> BinaryTree__:
        return _binary.RecursiveSerializer(lambda *args, **kwargs : BinaryTreeSerializer()).read(self._stream)

    def _read_list(self) -> typing.Optional[LinkedList[str]]:
        return _binary.OptionalSerializer(LinkedListSerializer(_binary.string_serializer)).read(self._stream)

    def _read_cwd(self) -> collections.abc.Iterable[DirectoryEntry]:
        return _binary.StreamSerializer(_binary.UnionSerializer(DirectoryEntry, [(DirectoryEntry.File, FileSerializer()), (DirectoryEntry.Directory, _binary.RecursiveSerializer(lambda *args, **kwargs : DirectorySerializer()))])).read(self._stream)


class BinaryMyProtocolIndexedWriter(_binary.BinaryProtocolIndexedWriter, MyProtocolWriterBase):
    """Binary indexed writer for the MyProtocol protocol."""


    def __init__(self, stream: typing.Union[typing.BinaryIO, str]) -> None:
        MyProtocolWriterBase.__init__(self)
        _binary.BinaryProtocolIndexedWriter.__init__(self, stream, MyProtocolWriterBase.schema)

    def _write_tree(self, value: BinaryTree) -> None:
        pos = self._stream.pos()
        self._index.set_step_offset("Tree", pos)
        BinaryTreeSerializer().write(self._stream, value)

    def _write_ptree(self, value: BinaryTree__) -> None:
        pos = self._stream.pos()
        self._index.set_step_offset("Ptree", pos)
        _binary.RecursiveSerializer(lambda *args, **kwargs : BinaryTreeSerializer()).write(self._stream, value)

    def _write_list(self, value: typing.Optional[LinkedList[str]]) -> None:
        pos = self._stream.pos()
        self._index.set_step_offset("List", pos)
        _binary.OptionalSerializer(LinkedListSerializer(_binary.string_serializer)).write(self._stream, value)

    def _write_cwd(self, value: collections.abc.Iterable[DirectoryEntry]) -> None:
        pos = self._stream.pos()
        self._index.set_step_offset("Cwd", pos)
        offsets, num_blocks = _binary.StreamSerializer(_binary.UnionSerializer(DirectoryEntry, [(DirectoryEntry.File, FileSerializer()), (DirectoryEntry.Directory, _binary.RecursiveSerializer(lambda *args, **kwargs : DirectorySerializer()))])).write_and_save_offsets(self._stream, value)
        self._index.add_stream_offsets("Cwd", offsets, num_blocks)


class BinaryMyProtocolIndexedReader(_binary.BinaryProtocolIndexedReader, MyProtocolIndexedReaderBase):
    """Binary indexed writer for the MyProtocol protocol."""


    def __init__(self, stream: typing.Union[io.BufferedReader, io.BytesIO, typing.BinaryIO, str]) -> None:
        MyProtocolIndexedReaderBase.__init__(self)
        _binary.BinaryProtocolIndexedReader.__init__(self, stream, MyProtocolIndexedReaderBase.schema)

    def _read_tree(self) -> BinaryTree:
        pos = self._index.get_step_offset("Tree")
        self._stream.seek(pos)
        return BinaryTreeSerializer().read(self._stream)

    def _read_ptree(self) -> BinaryTree__:
        pos = self._index.get_step_offset("Ptree")
        self._stream.seek(pos)
        return _binary.RecursiveSerializer(lambda *args, **kwargs : BinaryTreeSerializer()).read(self._stream)

    def _read_list(self) -> typing.Optional[LinkedList[str]]:
        pos = self._index.get_step_offset("List")
        self._stream.seek(pos)
        return _binary.OptionalSerializer(LinkedListSerializer(_binary.string_serializer)).read(self._stream)

    def _read_cwd(self, idx: int = 0) -> collections.abc.Iterable[DirectoryEntry]:
        offset, remaining = self._index.find_stream_item("Cwd", idx)
        self._stream.seek(offset)
        return _binary.StreamSerializer(_binary.UnionSerializer(DirectoryEntry, [(DirectoryEntry.File, FileSerializer()), (DirectoryEntry.Directory, _binary.RecursiveSerializer(lambda *args, **kwargs : DirectorySerializer()))])).read_mid_stream(self._stream, remaining)

    def _count_cwd(self) -> int:
        return self._index.get_stream_size("Cwd")

class BinaryTreeSerializer(_binary.RecordSerializer[BinaryTree]):
    def __init__(self) -> None:
        super().__init__([("value", _binary.int32_serializer), ("left", _binary.RecursiveSerializer(lambda *args, **kwargs : BinaryTreeSerializer())), ("right", _binary.RecursiveSerializer(lambda *args, **kwargs : BinaryTreeSerializer()))])

    def write(self, stream: _binary.CodedOutputStream, value: BinaryTree) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.value, value.left, value.right)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['value'], value['left'], value['right'])

    def read(self, stream: _binary.CodedInputStream) -> BinaryTree:
        field_values = self._read(stream)
        return BinaryTree(value=field_values[0], left=field_values[1], right=field_values[2])


class LinkedListSerializer(typing.Generic[T, T_NP], _binary.RecordSerializer[LinkedList[T]]):
    def __init__(self, t_serializer: _binary.TypeSerializer[T, T_NP]) -> None:
        super().__init__([("value", t_serializer), ("next", _binary.RecursiveSerializer(lambda *args, **kwargs : LinkedListSerializer(t_serializer)))])

    def write(self, stream: _binary.CodedOutputStream, value: LinkedList[T]) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.value, value.next)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['value'], value['next'])

    def read(self, stream: _binary.CodedInputStream) -> LinkedList[T]:
        field_values = self._read(stream)
        return LinkedList[T](value=field_values[0], next=field_values[1])


class FileSerializer(_binary.RecordSerializer[File]):
    def __init__(self) -> None:
        super().__init__([("name", _binary.string_serializer), ("data", _binary.VectorSerializer(_binary.uint8_serializer))])

    def write(self, stream: _binary.CodedOutputStream, value: File) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.name, value.data)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['name'], value['data'])

    def read(self, stream: _binary.CodedInputStream) -> File:
        field_values = self._read(stream)
        return File(name=field_values[0], data=field_values[1])


class DirectorySerializer(_binary.RecordSerializer[Directory]):
    def __init__(self) -> None:
        super().__init__([("name", _binary.string_serializer), ("entries", _binary.VectorSerializer(_binary.UnionSerializer(DirectoryEntry, [(DirectoryEntry.File, FileSerializer()), (DirectoryEntry.Directory, _binary.RecursiveSerializer(lambda *args, **kwargs : DirectorySerializer()))])))])

    def write(self, stream: _binary.CodedOutputStream, value: Directory) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.name, value.entries)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['name'], value['entries'])

    def read(self, stream: _binary.CodedInputStream) -> Directory:
        field_values = self._read(stream)
        return Directory(name=field_values[0], entries=field_values[1])


