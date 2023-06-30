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
from typing import Iterable, cast, Any, TypeAlias, get_args, Callable
import typing
import sandbox
import mrd
import numpy as np
import numpy.typing as npt
import inspect
from pprint import pprint


def print_value(value: Any) -> None:
    print(f"{value} {type(value)} dtype={value.dtype if isinstance(value, np.ndarray) else None}")  # type: ignore


# d = sandbox.Rec()
# pprint(d)

# X = sandbox.PT[int]

# print(get_args(X))


# pprint(type(sandbox.PT[int]))

# T = typing.TypeVar('T')
# @dataclasses.dataclass(slots=True, kw_only=True)
# class PT(typing.Generic[T]):
#     x: T
#     y: T

#     @staticmethod
#     def dtype(t_dtype: npt.DTypeLike) -> npt.DTypeLike:
#         return np.dtype([('x', t_dtype), ('y', t_dtype)], align=True)


# def my_dtype(t: type | types.GenericAlias) -> np.dtype[Any]:
#     if isinstance(t, type) or isinstance(t, types.UnionType):
#         if (res := _dtypes.get(t, None)) is not None:
#             if callable(res):
#                 raise RuntimeError(f"Generic type arguments for {t} not provided")
#             return res

#     origin = typing.get_origin(t)
#     if origin is not None:
#         if (res := _dtypes.get(origin, None)) is not None:
#             if callable(res):
#                 return res(typing.get_args(t))


#     raise RuntimeError(f"Cannot find dtype for {t}")


# _dtypes: dict[type, np.dtype[Any] | Callable[[tuple[type, ...]], np.dtype[Any]] ] = {
#     sandbox.PInt: np.dtype([('x', np.int32), ('y', np.int32)], align=True),
#     sandbox.Int32: np.dtype(np.int32),
#     sandbox.PT: lambda args: np.dtype([('x', my_dtype(args[0])), ('y', my_dtype(args[0]))], align=True),
# }

# # _dtypes[sandbox.Line] = lambda args: np.dtype([('start', my_dtype(types.GenericAlias(sandbox.PT, (args[0],)))), ('end', my_dtype(types.GenericAlias(sandbox.PT, (args[0],))))], align=True)

# a = types.GenericAlias(sandbox.PT, (sandbox.Int32,))
# print(a)

# print(typing.get_overloads(sandbox.PT))
# print(my_dtype(types.GenericAlias(sandbox.PT, (sandbox.Int32,))))
# print(my_dtype(sandbox.PT[sandbox.Int32]))
# print(my_dtype(sandbox.Line[sandbox.Int32]))
# print(my_dtype(sandbox.Int32))

print(sandbox.get_dtype(sandbox.PT[sandbox.Int32]))
print(sandbox.get_dtype(sandbox.Line[sandbox.Float32]))
print(sandbox.get_dtype(sandbox.Line[sandbox.A]))

print(sandbox.Person())

# with sandbox.BinaryP1Writer("test.bin") as w:
#     # dt = sandbox.PT.dtype(np.float32)
#     # value = np.array([(1, 2), (3, 4)], dtype=dt)

#     dt=np.dtype([('x', np.complex64), ('y', np.complex64)])
#     p = np.array([(1+2j, 3+4j), (5+6j, 7+8j)], dtype=dt)[0]

#     value = p

#     print_value(value)
#     w.write_my_value(value)
#     pass

# # os.system("hexdump -C test.bin")

# with sandbox.BinaryP1Reader("test.bin", sandbox.Types.NONE) as r:
#     v = r.read_my_value()
#     print_value(v)


# def produce_data() -> Iterable[mrd.StreamItem]:
#     acq = mrd.Acquisition(flags = mrd.AcquisitionFlags.FIRST_IN_SLICE,
#                           idx = mrd.EncodingCounters(),
#                           measurement_uid=2,
#                           )

# with mrd.BinaryMrdWriter("test.bin") as w:
#     h = mrd.Header(experimental_conditions=mrd.ExperimentalConditionsType(h1resonance_frequency_hz=440))
#     w.write_header(h)
#     w.write_data(produce_data())
