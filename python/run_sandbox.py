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
from typing import (
    Generic,
    Iterable,
    Literal,
    NamedTuple,
    cast,
    Any,
    TypeAlias,
    get_args,
    Callable,
)
import typing
import sandbox
import numpy as np
import numpy.typing as npt
import inspect
from pprint import pprint


def print_value(value: Any) -> None:
    print(f"{value} {type(value)} dtype={value.dtype if isinstance(value, np.ndarray) else None}")  # type: ignore


with sandbox.BinaryPWriter("test.bin") as w:
    value = sandbox.WithUnion(f=("string->int32", {"a": 1, "b": 2}))

    print_value(value)
    w.write_value(value)
    pass

# os.system("hexdump -C test.bin")

with sandbox.BinaryPReader("test.bin", sandbox.Types.NONE) as r:
    v = r.read_value()
    print_value(v)
