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

class BinaryPWriter(_binary.BinaryProtocolWriter, PWriterBase):
    """Binary writer for the P protocol."""


    def __init__(self, stream: typing.BinaryIO | str) -> None:
        PWriterBase.__init__(self)
        _binary.BinaryProtocolWriter.__init__(self, stream, PWriterBase.schema)

    def _write_a(self, value: Line[yardl.Int32]) -> None:
        _LineSerializer(_binary.int32_serializer).write(self._stream, value)


class BinaryPReader(_binary.BinaryProtocolReader, PReaderBase):
    """Binary writer for the P protocol."""


    def __init__(self, stream: io.BufferedReader | str, read_as_numpy: Types) -> None:
        PReaderBase.__init__(self, read_as_numpy)
        _binary.BinaryProtocolReader.__init__(self, stream, PReaderBase.schema)

    def _read_a(self) -> Line[yardl.Int32]:
        return _LineSerializer(_binary.int32_serializer).read(self._stream, self._read_as_numpy)

class _PTSerializer(typing.Generic[T, T_NP], _binary.RecordSerializer[PT[T]]):
    def __init__(self, t_serializer: _binary.TypeSerializer[T, T_NP]) -> None:
        super().__init__([("x", t_serializer), ("y", t_serializer)])

    def write(self, stream: _binary.CodedOutputStream, value: PT[T]) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.x, value.y)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['x'], value['y'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> PT[T]:
        field_values = self._read(stream, read_as_numpy)
        return PT[T](x=field_values[0], y=field_values[1])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        if not isinstance(value, PT):
            return False
        return (
            self._field_serializers[0][1].is_value_supported(value.x)
            and self._field_serializers[1][1].is_value_supported(value.y)
        )


class _PFloatSerializer(_binary.RecordSerializer[PFloat]):
    def __init__(self) -> None:
        super().__init__([("x", _binary.float32_serializer), ("y", _binary.float32_serializer)])

    def write(self, stream: _binary.CodedOutputStream, value: PFloat) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.x, value.y)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['x'], value['y'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> PFloat:
        field_values = self._read(stream, read_as_numpy)
        return PFloat(x=field_values[0], y=field_values[1])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        return isinstance(value, PFloat)


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


class _RecSerializer(_binary.RecordSerializer[Rec]):
    def __init__(self) -> None:
        super().__init__([("d", _binary.MapSerializer(_binary.string_serializer, _binary.int32_serializer)), ("f", _binary.EnumSerializer(_binary.int32_serializer, F)), ("g", _binary.OptionalSerializer(_PIntSerializer())), ("v", _binary.VectorSerializer(_binary.int32_serializer)), ("pf", _binary.FixedVectorSerializer(_binary.int32_serializer, 3)), ("pt", _binary.FixedVectorSerializer(_PTSerializer(_binary.int32_serializer), 3))])

    def write(self, stream: _binary.CodedOutputStream, value: Rec) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.d, value.f, value.g, value.v, value.pf, value.pt)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['d'], value['f'], value['g'], value['v'], value['pf'], value['pt'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> Rec:
        field_values = self._read(stream, read_as_numpy)
        return Rec(d=field_values[0], f=field_values[1], g=field_values[2], v=field_values[3], pf=field_values[4], pt=field_values[5])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        return isinstance(value, Rec)


class _GenRecSerializer(typing.Generic[T, T_NP], _binary.RecordSerializer[GenRec[T]]):
    def __init__(self, t_serializer: _binary.TypeSerializer[T, T_NP]) -> None:
        super().__init__([("t", t_serializer)])

    def write(self, stream: _binary.CodedOutputStream, value: GenRec[T]) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.t)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['t'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> GenRec[T]:
        field_values = self._read(stream, read_as_numpy)
        return GenRec[T](t=field_values[0])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        if not isinstance(value, GenRec):
            return False
        return (
            self._field_serializers[0][1].is_value_supported(value.t)
        )


class _DualGenRecSerializer(typing.Generic[T, T_NP], _binary.RecordSerializer[DualGenRec[T, T_NP]]):
    def __init__(self, t_serializer: _binary.TypeSerializer[T, T_NP]) -> None:
        super().__init__([("s", t_serializer), ("arr", _binary.DynamicNDArraySerializer(t_serializer))])

    def write(self, stream: _binary.CodedOutputStream, value: DualGenRec[T, T_NP]) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.s, value.arr)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['s'], value['arr'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> DualGenRec[T, T_NP]:
        field_values = self._read(stream, read_as_numpy)
        return DualGenRec[T, T_NP](s=field_values[0], arr=field_values[1])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        if not isinstance(value, DualGenRec):
            return False
        return (
            self._field_serializers[0][1].is_value_supported(value.s)
            and self._field_serializers[1][1].is_value_supported(value.arr)
        )


class _ASerializer(_binary.RecordSerializer[A]):
    def __init__(self) -> None:
        super().__init__([("pi", _binary.int32_serializer)])

    def write(self, stream: _binary.CodedOutputStream, value: A) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.pi)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['pi'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> A:
        field_values = self._read(stream, read_as_numpy)
        return A(pi=field_values[0])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        return isinstance(value, A)


class _LineSerializer(typing.Generic[T, T_NP], _binary.RecordSerializer[Line[T]]):
    def __init__(self, t_serializer: _binary.TypeSerializer[T, T_NP]) -> None:
        super().__init__([("start", _PTSerializer(t_serializer)), ("end", _PTSerializer(t_serializer))])

    def write(self, stream: _binary.CodedOutputStream, value: Line[T]) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.start, value.end)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['start'], value['end'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> Line[T]:
        field_values = self._read(stream, read_as_numpy)
        return Line[T](start=field_values[0], end=field_values[1])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        if not isinstance(value, Line):
            return False
        return (
            self._field_serializers[0][1].is_value_supported(value.start)
            and self._field_serializers[1][1].is_value_supported(value.end)
        )


class _LineIntSerializer(_binary.RecordSerializer[LineInt]):
    def __init__(self) -> None:
        super().__init__([("start", _PIntSerializer()), ("end", _PIntSerializer())])

    def write(self, stream: _binary.CodedOutputStream, value: LineInt) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.start, value.end)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['start'], value['end'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> LineInt:
        field_values = self._read(stream, read_as_numpy)
        return LineInt(start=field_values[0], end=field_values[1])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        return isinstance(value, LineInt)


class _FooSerializer(_binary.RecordSerializer[Foo]):
    def __init__(self) -> None:
        super().__init__([("a", _binary.int32_serializer), ("b", _binary.NDArraySerializer(_binary.float32_serializer, 2))])

    def write(self, stream: _binary.CodedOutputStream, value: Foo) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.a, value.b)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['a'], value['b'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> Foo:
        field_values = self._read(stream, read_as_numpy)
        return Foo(a=field_values[0], b=field_values[1])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        return isinstance(value, Foo)


class _PersonSerializer(_binary.RecordSerializer[Person]):
    def __init__(self) -> None:
        super().__init__([("numbers", _binary.FixedNDArraySerializer(_PTSerializer(_binary.int32_serializer), (2, 2,))), ("d", _binary.datetime_serializer)])

    def write(self, stream: _binary.CodedOutputStream, value: Person) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.numbers, value.d)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['numbers'], value['d'])

    def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> Person:
        field_values = self._read(stream, read_as_numpy)
        return Person(numbers=field_values[0], d=field_values[1])

    def is_value_supported(self, value: Any) -> bool:
        if isinstance(value, np.void) and value.dtype == self.overall_dtype():
            return True

        return isinstance(value, Person)


