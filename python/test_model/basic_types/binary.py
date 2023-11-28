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

from .. import _binary
from .. import yardl_types as yardl

class RecordWithUnionsSerializer(_binary.RecordSerializer[RecordWithUnions]):
    def __init__(self) -> None:
        super().__init__([("null_or_int_or_string", _binary.UnionSerializer(Int32OrString, [None, (Int32OrString.Int32, _binary.int32_serializer), (Int32OrString.String, _binary.string_serializer)])), ("date_or_datetime", _binary.UnionSerializer(TimeOrDatetime, [(TimeOrDatetime.Time, _binary.time_serializer), (TimeOrDatetime.Datetime, _binary.datetime_serializer)])), ("null_or_fruits_or_days_of_week", _binary.UnionSerializer(GenericNullableUnion2, [None, (GenericNullableUnion2.T1, _binary.EnumSerializer(_binary.int32_serializer, Fruits)), (GenericNullableUnion2.T2, _binary.EnumSerializer(_binary.int32_serializer, DaysOfWeek))]))])

    def write(self, stream: _binary.CodedOutputStream, value: RecordWithUnions) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.null_or_int_or_string, value.date_or_datetime, value.null_or_fruits_or_days_of_week)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['null_or_int_or_string'], value['date_or_datetime'], value['null_or_fruits_or_days_of_week'])

    def read(self, stream: _binary.CodedInputStream) -> RecordWithUnions:
        field_values = self._read(stream)
        return RecordWithUnions(null_or_int_or_string=field_values[0], date_or_datetime=field_values[1], null_or_fruits_or_days_of_week=field_values[2])


class GenericRecordWithComputedFieldsSerializer(typing.Generic[T0, T0_NP, T1, T1_NP], _binary.RecordSerializer[GenericRecordWithComputedFields[T0, T1]]):
    def __init__(self, t0_serializer: _binary.TypeSerializer[T0, T0_NP], t1_serializer: _binary.TypeSerializer[T1, T1_NP]) -> None:
        super().__init__([("f1", _binary.UnionSerializer(T0OrT1, [(T0OrT1[T0, T1].T0, t0_serializer), (T0OrT1[T0, T1].T1, t1_serializer)]))])

    def write(self, stream: _binary.CodedOutputStream, value: GenericRecordWithComputedFields[T0, T1]) -> None:
        if isinstance(value, np.void):
            self.write_numpy(stream, value)
            return
        self._write(stream, value.f1)

    def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:
        self._write(stream, value['f1'])

    def read(self, stream: _binary.CodedInputStream) -> GenericRecordWithComputedFields[T0, T1]:
        field_values = self._read(stream)
        return GenericRecordWithComputedFields[T0, T1](f1=field_values[0])


