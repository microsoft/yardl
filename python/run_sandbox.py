#! /usr/bin/env python3

import dataclasses
import os
import timeit
import sandbox
import sandbox._binary
import numpy as np
import typing
from datetime import date, time, datetime, timedelta
import sys
import ctypes

class Point(ctypes.Structure):
    _fields_ = [("x", ctypes.c_int),("y", ctypes.c_int)]

MyTime = time | np.timedelta64

MyInt = int | np.int32


MyInt2 = int | np.int32


def x(p: int) -> int:
    return p *2



def main():
    with sandbox.BinaryP1Writer("test.bin") as w:
        w.write_an_int(8)
        w.write_a_stream([1, 2, 3, 4, 5])
        w.write_optional(2)
        w.write_union(22)
        w.write_date(date.today())


    r = sandbox.MyRec(f2="hello", f3=sandbox.MyFlags.A | sandbox.MyFlags.B)
    print(r)
    os.system("hexdump -C test.bin")

    c = sandbox.MyRec(f2="hello", f3=sandbox.MyFlags.A | sandbox.MyFlags.B)
    print(c)
    c2 = sandbox.R(f2="hello", f3=sandbox.MyFlags.A | sandbox.MyFlags.B)
    print(c2)



if  __name__ == "__main__"    :
    main()
