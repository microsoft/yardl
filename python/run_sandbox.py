#! /usr/bin/env python3

import collections
import dataclasses
import os
import timeit
import sandbox
import sandbox._binary
import numpy as np
import numpy.typing as npt
import typing
from datetime import date, time, datetime, timedelta
import sys
import ctypes
import sandbox.yardl_types as yardl


class Point(ctypes.Structure):
    _fields_ = [("x", ctypes.c_int), ("y", ctypes.c_int)]


def main():
    with sandbox.BinaryP1Writer("test.bin") as w:
        # w.write_an_int(8)
        # w.write_a_stream([1, 2, 3, 4, 5])
        # w.write_optional(2)
        # w.write_union(None)
        # w.write_flag(sandbox.MyFlags.A | sandbox.MyFlags.B)
        # w.write_vec([1,2,3])
        new_var = np.array([[1,2,3],[4,5,6]], dtype=np.uint32)
        w.write_arr(new_var)
        # w.write_map({"a": 1, "b": 2})


        pass

    r = sandbox.MyRec(f2="hello", f3=sandbox.MyFlags.A | sandbox.MyFlags.B)
    print(r)
    os.system("hexdump -C test.bin")

    c = sandbox.MyRec(f2="hello", f3=sandbox.MyFlags.A | sandbox.MyFlags.B)
    print(c)
    c2 = sandbox.R(f2="hello", f3=sandbox.MyFlags.A | sandbox.MyFlags.B)
    print(c2)

    print(np.dtype('object'))

    def f(arr : npt.NDArray[np.int32]):
        print(arr)

    new_var = np.array([[1,2,3],[4,5,6]], dtype=np.int32)
    f(new_var)




if __name__ == "__main__":
    main()
