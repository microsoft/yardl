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

    schema = r"""{"protocol":{"name":"MyProtocol","sequence":[{"name":"header","type":"Sketch.Header"},{"name":"samples","type":{"stream":{"items":"Sketch.Sample"}}}]},"types":[{"name":"Header","fields":[{"name":"subject","type":"string"}]},{"name":"Sample","fields":[{"name":"id","type":"uint32"},{"name":"data","type":{"array":{"items":"int32","dimensions":1}}}]}]}"""

    def close(self) -> None:
        if self._state == 3:
            try:
                self._end_stream()
                return
            finally:
                self._close()
        self._close()
        if self._state != 4:
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

    def write_header(self, value: Header) -> None:
        """Ordinal 0"""

        if self._state != 0:
            self._raise_unexpected_state(0)

        self._write_header(value)
        self._state = 2

    def write_samples(self, value: collections.abc.Iterable[Sample]) -> None:
        """Ordinal 1"""

        if self._state & ~1 != 2:
            self._raise_unexpected_state(2)

        self._write_samples(value)
        self._state = 3

    @abc.abstractmethod
    def _write_header(self, value: Header) -> None:
        raise NotImplementedError()

    @abc.abstractmethod
    def _write_samples(self, value: collections.abc.Iterable[Sample]) -> None:
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
            return 'write_header'
        if state == 2:
            return 'write_samples'
        return "<unknown>"

class MyProtocolReaderBase(abc.ABC):
    """Abstract indexed reader for the MyProtocol protocol."""

    def __init__(self) -> None:
        self._state = 0

    def close(self) -> None:
        self._close()
        if self._state != 4:
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

    def read_header(self) -> Header:
        """Ordinal 0"""

        if self._state != 0:
            self._raise_unexpected_state(0)

        value = self._read_header()
        self._state = 2
        return value

    def read_samples(self) -> collections.abc.Iterable[Sample]:
        """Ordinal 1"""

        if self._state != 2:
            self._raise_unexpected_state(2)

        value = self._read_samples()
        self._state = 3
        return self._wrap_iterable(value, 4)

    def copy_to(self, writer: MyProtocolWriterBase) -> None:
        writer.write_header(self.read_header())
        writer.write_samples(self.read_samples())

    @abc.abstractmethod
    def _read_header(self) -> Header:
        raise NotImplementedError()

    @abc.abstractmethod
    def _read_samples(self) -> collections.abc.Iterable[Sample]:
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
            return 'read_header'
        if state == 2:
            return 'read_samples'
        return "<unknown>"

class MyProtocolIndexedReaderBase(MyProtocolReaderBase):
    """Abstract reader for the MyProtocol protocol."""

    def close(self) -> None:
        self._close()

    def read_samples(self, idx: int = 0) -> collections.abc.Iterable[Sample]:
        value = self._read_samples(idx)
        return self._wrap_iterable(value, 4)

    def count_samples(self) -> int:
        return self._count_samples()

    @abc.abstractmethod
    def _read_samples(self, idx: int = 0) -> collections.abc.Iterable[Sample]:
        raise NotImplementedError()

    @abc.abstractmethod
    def _count_samples(self) -> int:
        raise NotImplementedError()

    def _raise_unexpected_state(self, actual: int) -> None:
        pass
