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
from . import _ndjson
from . import yardl_types as yardl

class HeaderConverter(_ndjson.JsonConverter[Header, np.void]):
    def __init__(self) -> None:
        self._subject_converter = _ndjson.string_converter
        super().__init__(np.dtype([
            ("subject", self._subject_converter.overall_dtype()),
        ]))

    def to_json(self, value: Header) -> object:
        if not isinstance(value, Header): # pyright: ignore [reportUnnecessaryIsInstance]
            raise TypeError("Expected 'Header' instance")
        json_object = {}

        json_object["subject"] = self._subject_converter.to_json(value.subject)
        return json_object

    def numpy_to_json(self, value: np.void) -> object:
        if not isinstance(value, np.void): # pyright: ignore [reportUnnecessaryIsInstance]
            raise TypeError("Expected 'np.void' instance")
        json_object = {}

        json_object["subject"] = self._subject_converter.numpy_to_json(value["subject"])
        return json_object

    def from_json(self, json_object: object) -> Header:
        if not isinstance(json_object, dict):
            raise TypeError("Expected 'dict' instance")
        return Header(
            subject=self._subject_converter.from_json(json_object["subject"],),
        )

    def from_json_to_numpy(self, json_object: object) -> np.void:
        if not isinstance(json_object, dict):
            raise TypeError("Expected 'dict' instance")
        return (
            self._subject_converter.from_json_to_numpy(json_object["subject"]),
        ) # type:ignore 


class SampleConverter(_ndjson.JsonConverter[Sample, np.void]):
    def __init__(self) -> None:
        self._id_converter = _ndjson.uint32_converter
        self._data_converter = _ndjson.NDArrayConverter(_ndjson.int32_converter, 1)
        super().__init__(np.dtype([
            ("id", self._id_converter.overall_dtype()),
            ("data", self._data_converter.overall_dtype()),
        ]))

    def to_json(self, value: Sample) -> object:
        if not isinstance(value, Sample): # pyright: ignore [reportUnnecessaryIsInstance]
            raise TypeError("Expected 'Sample' instance")
        json_object = {}

        json_object["id"] = self._id_converter.to_json(value.id)
        json_object["data"] = self._data_converter.to_json(value.data)
        return json_object

    def numpy_to_json(self, value: np.void) -> object:
        if not isinstance(value, np.void): # pyright: ignore [reportUnnecessaryIsInstance]
            raise TypeError("Expected 'np.void' instance")
        json_object = {}

        json_object["id"] = self._id_converter.numpy_to_json(value["id"])
        json_object["data"] = self._data_converter.numpy_to_json(value["data"])
        return json_object

    def from_json(self, json_object: object) -> Sample:
        if not isinstance(json_object, dict):
            raise TypeError("Expected 'dict' instance")
        return Sample(
            id=self._id_converter.from_json(json_object["id"],),
            data=self._data_converter.from_json(json_object["data"],),
        )

    def from_json_to_numpy(self, json_object: object) -> np.void:
        if not isinstance(json_object, dict):
            raise TypeError("Expected 'dict' instance")
        return (
            self._id_converter.from_json_to_numpy(json_object["id"]),
            self._data_converter.from_json_to_numpy(json_object["data"]),
        ) # type:ignore 


class NDJsonMyProtocolWriter(_ndjson.NDJsonProtocolWriter, MyProtocolWriterBase):
    """NDJson writer for the MyProtocol protocol."""


    def __init__(self, stream: typing.Union[typing.TextIO, str]) -> None:
        MyProtocolWriterBase.__init__(self)
        _ndjson.NDJsonProtocolWriter.__init__(self, stream, MyProtocolWriterBase.schema)

    def _write_header(self, value: Header) -> None:
        converter = HeaderConverter()
        json_value = converter.to_json(value)
        self._write_json_line({"header": json_value})

    def _write_samples(self, value: collections.abc.Iterable[Sample]) -> None:
        converter = SampleConverter()
        for item in value:
            json_item = converter.to_json(item)
            self._write_json_line({"samples": json_item})


class NDJsonMyProtocolReader(_ndjson.NDJsonProtocolReader, MyProtocolReaderBase):
    """NDJson writer for the MyProtocol protocol."""


    def __init__(self, stream: typing.Union[io.BufferedReader, typing.TextIO, str]) -> None:
        MyProtocolReaderBase.__init__(self)
        _ndjson.NDJsonProtocolReader.__init__(self, stream, MyProtocolReaderBase.schema)

    def _read_header(self) -> Header:
        json_object = self._read_json_line("header", True)
        converter = HeaderConverter()
        return converter.from_json(json_object)

    def _read_samples(self) -> collections.abc.Iterable[Sample]:
        converter = SampleConverter()
        while (json_object := self._read_json_line("samples", False)) is not _ndjson.MISSING_SENTINEL:
            yield converter.from_json(json_object)
