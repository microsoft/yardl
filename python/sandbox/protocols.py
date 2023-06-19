# This file was generated by the "yardl" tool. DO NOT EDIT.


import abc
import collections.abc
import datetime
import numpy as np
import numpy.typing as npt
import typing
from . import *
from . import yardl_types as yardl

class P1WriterBase(abc.ABC):
    """Abstract writer for the P1 protocol."""

    schema = """{"protocol":{"name":"P1","sequence":[{"name":"myValue","type":"datetime"}]},"types":null}"""

    def write_my_value(self, value: yardl.DateTime) -> None:
        """Ordinal 0"""
        self._write_my_value(value)

    @abc.abstractmethod
    def _write_my_value(self, value: yardl.DateTime) -> None:
        raise NotImplementedError()

class P1ReaderBase(abc.ABC):
    """Abstract reader for the P1 protocol."""

    def __init__(self, read_as_numpy: Types = Types.NONE) -> None:
        self._read_as_numpy = read_as_numpy

    schema = P1WriterBase.schema

    def read_my_value(self) -> yardl.DateTime:
        """Ordinal 0"""
        return self._read_my_value()

    @abc.abstractmethod
    def _read_my_value(self) -> yardl.DateTime:
        raise NotImplementedError()

