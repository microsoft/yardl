# This file was generated by the "yardl" tool. DO NOT EDIT.

# pyright: reportUnusedClass=false

import collections.abc
import datetime
import io
import typing
import numpy as np
import numpy.typing as npt

from . import *
from . import _binary
from . import yardl_types as yardl

T = typing.TypeVar('T')
T_NP = typing.TypeVar('T_NP', bound=np.generic)
U = typing.TypeVar('U')
U_NP = typing.TypeVar('U_NP', bound=np.generic)

class BinaryPWriter(_binary.BinaryProtocolWriter, PWriterBase):
    """Binary writer for the P protocol."""


    def __init__(self, stream: typing.BinaryIO | str) -> None:
        PWriterBase.__init__(self)
        _binary.BinaryProtocolWriter.__init__(self, stream, PWriterBase.schema)

    def _write_value(self, value: WithUnion) -> None:
        _WithUnionSerializer().write(self._stream, value)


class BinaryPReader(_binary.BinaryProtocolReader, PReaderBase):
    """Binary writer for the P protocol."""


    def __init__(self, stream: io.BufferedReader | str, read_as_numpy: Types = Types.NONE) -> None:
        PReaderBase.__init__(self, read_as_numpy)
        _binary.BinaryProtocolReader.__init__(self, stream, PReaderBase.schema)

    def _read_value(self) -> WithUnion:
        return _WithUnionSerializer().read(self._stream, self._read_as_numpy)

class _PTSerializer(typing.Generic[T, T_NP], _binary.RecordSerializer[PT[T]]):
    def __init__(self, t_serializer: _binary.TypeSerializer[T, T_NP]) -> None:
        super().__init__([("x", t_serializer), ("y", t_serializer), ("z", _binary.int32_serializer)])

    def write(self, stream: _binary.CodedOutputStream, value: PT[T]) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.x, value.y, value.z)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['x'], value['y'], value['z'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> PT[T]:
        field_values = self._read(stream, read_as_numpy)
        return PT[T](x=field_values[0], y=field_values[1], z=field_values[2])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        if not isinstance(value, PT):
            return False
        return (
            self._field_serializers[0][1].is_value_supported(value.x)
            and self._field_serializers[1][1].is_value_supported(value.y)
        )


class _PIntSerializer(_binary.RecordSerializer[PInt]):
    def __init__(self) -> None:
        super().__init__([("x", _binary.int32_serializer), ("y", _binary.int32_serializer)])

    def write(self, stream: _binary.CodedOutputStream, value: PInt) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.x, value.y)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['x'], value['y'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> PInt:
        field_values = self._read(stream, read_as_numpy)
        return PInt(x=field_values[0], y=field_values[1])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        return isinstance(value, PInt)


class _WithUnionSerializer(_binary.RecordSerializer[WithUnion]):
    def __init__(self) -> None:
        super().__init__([("f", _binary.UnionSerializer([_binary.none_serializer, _binary.int32_serializer, _binary.VectorSerializer(_binary.float32_serializer), _binary.string_serializer, _PIntSerializer(), _binary.MapSerializer(_binary.string_serializer, _binary.int32_serializer)]))])

    def write(self, stream: _binary.CodedOutputStream, value: WithUnion) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.f)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['f'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> WithUnion:
        field_values = self._read(stream, read_as_numpy)
        return WithUnion(f=field_values[0])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        return isinstance(value, WithUnion)


