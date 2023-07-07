# This file was generated by the "yardl" tool. DO NOT EDIT.


import abc
import collections.abc
import datetime
import numpy as np
import numpy.typing as npt
import typing
from . import *
from . import yardl_types as yardl

class HelloWorldWriterBase(abc.ABC):
    """Abstract writer for the HelloWorld protocol."""


    def __init__(self) -> None:
        self._state = 0

    schema = r"""{"protocol":{"name":"HelloWorld","sequence":[{"name":"data","type":{"stream":{"items":{"array":{"items":"complexfloat64","dimensions":[{"length":2}]}}}}}]},"types":null}"""

    def __enter__(self):
        return self

    def __exit__(self, exc_type: type[BaseException] | None, exc: BaseException | None, traceback: typing.Any | None) -> None:
        self.close()
        if exc is None and self._state != 1:
            expected_method = self._state_to_method_name(self._state)
            raise ProtocolException(f"Protocol writer closed before all steps were called. Expected to call to '{expected_method}'.")

    def write_data(self, value: collections.abc.Iterable[npt.NDArray[np.complex128]]) -> None:
        """Ordinal 0"""

        if self._state != 0:
            self._raise_unexpected_state(0)

        self._write_data(value)
        self._state = 1

    @abc.abstractmethod
    def _write_data(self, value: collections.abc.Iterable[npt.NDArray[np.complex128]]) -> None:
        raise NotImplementedError()

    @abc.abstractmethod
    def close(self) -> None:
        raise NotImplementedError()

    def _raise_unexpected_state(self, actual: int) -> None:
        expected_method = self._state_to_method_name(self._state)
        actual_method = self._state_to_method_name(actual)
        raise ProtocolException(f"Expected to call to '{expected_method}' but received call to '{actual_method}'.")

    def _state_to_method_name(self, state: int) -> str:
        if state == 0:
            return 'write_data'
        return "<unknown>"

class HelloWorldReaderBase(abc.ABC):
    """Abstract reader for the HelloWorld protocol."""


    def __init__(self, read_as_numpy: Types = Types.NONE) -> None:
        self._read_as_numpy = read_as_numpy
        self._state = 0

    schema = HelloWorldWriterBase.schema

    def __enter__(self):
        return self

    def __exit__(self, exc_type: type[BaseException] | None, exc: BaseException | None, traceback: typing.Any | None) -> None:
        self.close()
        if exc is None and self._state != 2:
            if self._state % 2 == 1:
                previous_method = self._state_to_method_name(self._state - 1)
                raise ProtocolException(f"Protocol reader closed before all data was consumed. The iterable returned by '{previous_method}' was not fully consumed.")
            else:
                expected_method = self._state_to_method_name(self._state)
                raise ProtocolException(f"Protocol reader closed before all data was consumed. Expected call to '{expected_method}'.")
            	

    @abc.abstractmethod
    def close(self) -> None:
        raise NotImplementedError()

    def read_data(self) -> collections.abc.Iterable[npt.NDArray[np.complex128]]:
        """Ordinal 0"""

        if self._state != 0:
            self._raise_unexpected_state(0)

        value = self._read_data()
        self._state = 1
        return self._wrap_iterable(value, 2)

    def copy_to(self, writer: HelloWorldWriterBase) -> None:
        writer.write_data(self.read_data())

    @abc.abstractmethod
    def _read_data(self) -> collections.abc.Iterable[npt.NDArray[np.complex128]]:
        raise NotImplementedError()

    T = typing.TypeVar('T')
    def _wrap_iterable(self, iterable: collections.abc.Iterable[T], final_state: int) -> collections.abc.Iterable[T]:
        yield from iterable
        self._state = final_state

    def _raise_unexpected_state(self, actual: int) -> None:
        actual_method = self._state_to_method_name(actual)
        if self._state % 2 == 1:
            previous_method = self._state_to_method_name(self._state - 1)
            raise ProtocolException(f"Received call to '{actual_method}' but the iterable returned by '{previous_method}' was not fully consumed.")
        else:
            expected_method = self._state_to_method_name(self._state)
            raise ProtocolException(f"Expected to call to '{expected_method}' but received call to '{actual_method}'.")
        	
    def _state_to_method_name(self, state: int) -> str:
        if state == 0:
            return 'read_data'
        return "<unknown>"

class ProtocolException(Exception):
    """Raised when the contract of a protocol is not respected."""
    pass
