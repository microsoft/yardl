#! /usr/bin/env python3

import abc
import dataclasses
import datetime
import enum
import os
import io
import sys
import typing
import sandbox
import numpy as np
import numpy.typing as npt
import inspect


T_NP = typing.TypeVar("T_NP", bound=np.generic)
Img : typing.TypeAlias = npt.NDArray[T_NP]

def f(arr: Img[np.int32]) -> None:
    pass

print(inspect.getmro(np.int32))

T = typing.TypeVar('T')
@dataclasses.dataclass(slots=True, kw_only=True)
class MyClass(typing.Generic[T, T_NP]):
    tScalar: T
    tArr: npt.NDArray[T_NP]
del(T)

myc = MyClass(tScalar=1, tArr=np.array([1,2,3], dtype=np.int32))




arr = np.array([1,2,3], dtype=np.int32)

# with sandbox.BinaryP1Writer("test.bin") as w:

#     value = np.array([((1, 2), (3, 4))], dtype=sandbox.Line.dtype(np.int32))
#     print(f"{value} {type(value)}")
#     w.write_my_value(value)
#     pass

# # os.system("hexdump -C test.bin")

# with sandbox.BinaryP1Reader("test.bin", sandbox.Types.INTEGER) as r:
#     v = r.read_my_value()
#     print(f"{v} {type(v)}")
