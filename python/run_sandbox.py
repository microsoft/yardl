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
from typing import cast, Any
import sandbox
import numpy as np
import numpy.typing as npt
import inspect


def print_value(value: Any) -> None:
    print(f"{value} {type(value)} dtype={value.dtype if isinstance(value, np.ndarray) else None}") # type: ignore

with sandbox.BinaryP1Writer("test.bin") as w:
    # dt = sandbox.PT.dtype(np.float32)
    # value = np.array([(1, 2), (3, 4)], dtype=dt)

    dt=np.dtype([('x', np.complex64), ('y', np.complex64)])
    p = np.array([(1+2j, 3+4j), (5+6j, 7+8j)], dtype=dt)[0]

    value = p

    print_value(value)
    w.write_my_value(value)
    pass

# os.system("hexdump -C test.bin")

with sandbox.BinaryP1Reader("test.bin", sandbox.Types.NONE) as r:
    v = r.read_my_value()
    print_value(v)
