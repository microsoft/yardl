#! /usr/bin/env python3

import abc
import array
import dataclasses
import datetime
import enum
import os
import io
import struct
import sys
import timeit
import types
from typing import Generic, Iterable, Literal, NamedTuple, cast, Any, TypeAlias, get_args, Callable
import typing
import sandbox
import numpy as np
import numpy.typing as npt
import inspect
from pprint import pprint


def print_value(value: Any) -> None:
    print(f"{value} {type(value)} dtype={value.dtype if isinstance(value, np.ndarray) else None}")  # type: ignore


with sandbox.BinaryPWriter("test.bin") as w:

    value = sandbox.WithUnion(f=42.0)

    print_value(value)
    w.write_value(value)
    pass

# os.system("hexdump -C test.bin")

with sandbox.BinaryPReader("test.bin", sandbox.Types.NONE) as r:
    v = r.read_value()
    print_value(v)


# print(sandbox.types._get__int32_float32_MyString_PInt_int32Vector_union_index(np.array([1, 2, 3], dtype=np.int32)))

pt = sandbox.PT[int](x=1, y=2, z=3)

print(typing.get_args(type(pt)))

# T_NP = typing.TypeVar('T_NP', bound=np.generic)
# class MyUnion(typing.Generic[T_NP]):
#     class Tag(enum.Enum):
#         NONE =  0,
#         INT32 = 1,
#         FLOATS = 2,
#         SOMETHING = 3,
#         SOMETHING_ELSE = 4,

#     def __init__(self,
#                  *,
#                  _tag: Tag | None = None,
#                  _value: sandbox.Int32 | sandbox.Float32 | npt.NDArray[T_NP] | None = None,

#                  int32 : sandbox.Int32 | None = None,
#                  floats : sandbox.Float32 | None = None,
#                  something: npt.NDArray[T_NP] | None = None,
#                  something_else : npt.NDArray[T_NP] | None = None) -> None:

#         if _tag is not None:
#             self._tag = _tag
#             self._value = _value
#             return

#         if int32 is not None:
#             self._tag = MyUnion.Tag.INT32
#             self._value = int32
#             return

#         if floats is not None:
#             self._tag = MyUnion.Tag.FLOATS
#             self._value = floats
#             return

#         if something is not None:
#             self._tag = MyUnion.Tag.SOMETHING
#             self._value = something
#             return

#         if something_else is not None:
#             self._tag = MyUnion.Tag.SOMETHING_ELSE
#             self._value = something_else
#             return

#     def tag(self) -> Tag:
#         return self._tag

#     def value(self) -> sandbox.Int32 | sandbox.Float32 | npt.NDArray[T_NP] | None:
#         return self._value

#     def __repr__(self) -> str:
#         return f"MyUnion(tag={self._tag}, value={self._value})"

# # use it like this:
# u = MyUnion(int32=1)

# # or like this:
# u = MyUnion(_tag=MyUnion.Tag.INT32, _value=1)

# # MyTag : TypeAlias =

TTag = typing.TypeVar('TTag')
TValue = typing.TypeVar('TValue')
class AUnion(Generic[TTag, TValue]):
    def __init__(self, tag: TTag, value: TValue) -> None:
        self.tag = tag
        self.value = value

# class TaggedUnion(NamedTuple, Generic[TTag, TValue]):
#     tag: TTag
#     value: TValue

TaggedUnion : TypeAlias = tuple[TTag, TValue]

MyUnion = TaggedUnion[
    Literal['int32', 'float', 'none'],
    int | float | None]

def do_something(u: MyUnion) -> None:
    t, v = u
    if t == "int32":
        print(v)

do_something(("int32", 1))
