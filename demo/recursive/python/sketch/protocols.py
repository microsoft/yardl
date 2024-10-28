# This file was generated by the "yardl" tool. DO NOT EDIT.

# pyright: reportUnusedImport=false

import abc
import collections.abc
import datetime
import typing

import numpy as np
import numpy.typing as npt

from .types import *
from .yardl_types import ProtocolError
from . import yardl_types as yardl

class MyProtocolWriterBase(abc.ABC):
    """Abstract writer for the MyProtocol protocol."""


    def __init__(self) -> None:
        self._state = 0

    schema = r"""{"protocol":{"name":"MyProtocol","sequence":[{"name":"tree","type":"Sketch.BinaryTree"},{"name":"ptree","type":"Sketch.BinaryTree"},{"name":"list","type":[null,{"name":"Sketch.LinkedList","typeArguments":["string"]}]},{"name":"cwd","type":{"stream":{"items":"Sketch.DirectoryEntry"}}}]},"types":[{"name":"BinaryTree","fields":[{"name":"value","type":"int32"},{"name":"left","type":"Sketch.BinaryTree"},{"name":"right","type":"Sketch.BinaryTree"}]},{"name":"Directory","fields":[{"name":"name","type":"string"},{"name":"entries","type":{"vector":{"items":"Sketch.DirectoryEntry"}}}]},{"name":"DirectoryEntry","type":[{"tag":"File","type":"Sketch.File"},{"tag":"Directory","type":"Sketch.Directory"}]},{"name":"File","fields":[{"name":"name","type":"string"},{"name":"data","type":{"vector":{"items":"uint8"}}}]},{"name":"LinkedList","typeParameters":["T"],"fields":[{"name":"value","type":"T"},{"name":"next","type":{"name":"Sketch.LinkedList","typeArguments":["T"]}}]}]}"""

    def close(self) -> None:
        if self._state == 7:
            try:
                self._end_stream()
                return
            finally:
                self._close()
        self._close()
        if self._state != 8:
            expected_method = self._state_to_method_name((self._state + 1) & ~1)
            raise ProtocolError(f"Protocol writer closed before all steps were called. Expected to call to '{expected_method}'.")

    def __enter__(self):
        return self

    def __exit__(self, exc_type: typing.Optional[type[BaseException]], exc: typing.Optional[BaseException], traceback: object) -> None:
        try:
            self.close()
        except Exception as e:
            if exc is None:
                raise e

    def write_tree(self, value: BinaryTree) -> None:
        """Ordinal 0"""

        if self._state != 0:
            self._raise_unexpected_state(0)

        self._write_tree(value)
        self._state = 2

    def write_ptree(self, value: BinaryTree__) -> None:
        """Ordinal 1"""

        if self._state != 2:
            self._raise_unexpected_state(2)

        self._write_ptree(value)
        self._state = 4

    def write_list(self, value: typing.Optional[LinkedList[str]]) -> None:
        """Ordinal 2"""

        if self._state != 4:
            self._raise_unexpected_state(4)

        self._write_list(value)
        self._state = 6

    def write_cwd(self, value: collections.abc.Iterable[DirectoryEntry]) -> None:
        """Ordinal 3

        dirs: !stream
          items: Directory
        """

        if self._state & ~1 != 6:
            self._raise_unexpected_state(6)

        self._write_cwd(value)
        self._state = 7

    @abc.abstractmethod
    def _write_tree(self, value: BinaryTree) -> None:
        raise NotImplementedError()

    @abc.abstractmethod
    def _write_ptree(self, value: BinaryTree__) -> None:
        raise NotImplementedError()

    @abc.abstractmethod
    def _write_list(self, value: typing.Optional[LinkedList[str]]) -> None:
        raise NotImplementedError()

    @abc.abstractmethod
    def _write_cwd(self, value: collections.abc.Iterable[DirectoryEntry]) -> None:
        raise NotImplementedError()

    @abc.abstractmethod
    def _close(self) -> None:
        pass

    @abc.abstractmethod
    def _end_stream(self) -> None:
        pass

    def _raise_unexpected_state(self, actual: int) -> None:
        expected_method = self._state_to_method_name(self._state)
        actual_method = self._state_to_method_name(actual)
        raise ProtocolError(f"Expected to call to '{expected_method}' but received call to '{actual_method}'.")

    def _state_to_method_name(self, state: int) -> str:
        if state == 0:
            return 'write_tree'
        if state == 2:
            return 'write_ptree'
        if state == 4:
            return 'write_list'
        if state == 6:
            return 'write_cwd'
        return "<unknown>"

class MyProtocolReaderBase(abc.ABC):
    """Abstract reader for the MyProtocol protocol."""


    def __init__(self) -> None:
        self._state = 0

    def close(self) -> None:
        self._close()
        if self._state != 8:
            if self._state % 2 == 1:
                previous_method = self._state_to_method_name(self._state - 1)
                raise ProtocolError(f"Protocol reader closed before all data was consumed. The iterable returned by '{previous_method}' was not fully consumed.")
            else:
                expected_method = self._state_to_method_name(self._state)
                raise ProtocolError(f"Protocol reader closed before all data was consumed. Expected call to '{expected_method}'.")
            	

    schema = MyProtocolWriterBase.schema

    def __enter__(self):
        return self

    def __exit__(self, exc_type: typing.Optional[type[BaseException]], exc: typing.Optional[BaseException], traceback: object) -> None:
        try:
            self.close()
        except Exception as e:
            if exc is None:
                raise e

    @abc.abstractmethod
    def _close(self) -> None:
        raise NotImplementedError()

    def read_tree(self) -> BinaryTree:
        """Ordinal 0"""

        if self._state != 0:
            self._raise_unexpected_state(0)

        value = self._read_tree()
        self._state = 2
        return value

    def read_ptree(self) -> BinaryTree__:
        """Ordinal 1"""

        if self._state != 2:
            self._raise_unexpected_state(2)

        value = self._read_ptree()
        self._state = 4
        return value

    def read_list(self) -> typing.Optional[LinkedList[str]]:
        """Ordinal 2"""

        if self._state != 4:
            self._raise_unexpected_state(4)

        value = self._read_list()
        self._state = 6
        return value

    def read_cwd(self) -> collections.abc.Iterable[DirectoryEntry]:
        """Ordinal 3

        dirs: !stream
          items: Directory
        """

        if self._state != 6:
            self._raise_unexpected_state(6)

        value = self._read_cwd()
        self._state = 7
        return self._wrap_iterable(value, 8)

    def copy_to(self, writer: MyProtocolWriterBase) -> None:
        writer.write_tree(self.read_tree())
        writer.write_ptree(self.read_ptree())
        writer.write_list(self.read_list())
        writer.write_cwd(self.read_cwd())

    @abc.abstractmethod
    def _read_tree(self) -> BinaryTree:
        raise NotImplementedError()

    @abc.abstractmethod
    def _read_ptree(self) -> BinaryTree__:
        raise NotImplementedError()

    @abc.abstractmethod
    def _read_list(self) -> typing.Optional[LinkedList[str]]:
        raise NotImplementedError()

    @abc.abstractmethod
    def _read_cwd(self) -> collections.abc.Iterable[DirectoryEntry]:
        raise NotImplementedError()

    T = typing.TypeVar('T')
    def _wrap_iterable(self, iterable: collections.abc.Iterable[T], final_state: int) -> collections.abc.Iterable[T]:
        yield from iterable
        self._state = final_state

    def _raise_unexpected_state(self, actual: int) -> None:
        actual_method = self._state_to_method_name(actual)
        if self._state % 2 == 1:
            previous_method = self._state_to_method_name(self._state - 1)
            raise ProtocolError(f"Received call to '{actual_method}' but the iterable returned by '{previous_method}' was not fully consumed.")
        else:
            expected_method = self._state_to_method_name(self._state)
            raise ProtocolError(f"Expected to call to '{expected_method}' but received call to '{actual_method}'.")
        	
    def _state_to_method_name(self, state: int) -> str:
        if state == 0:
            return 'read_tree'
        if state == 2:
            return 'read_ptree'
        if state == 4:
            return 'read_list'
        if state == 6:
            return 'read_cwd'
        return "<unknown>"

