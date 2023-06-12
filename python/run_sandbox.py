#! /usr/bin/env python3

import os
import sandbox
import sandbox._binary
import numpy as np
import typing
from datetime import date, time, datetime
import sys
import ctypes

class ImageHeader(ctypes.Structure):
    _pack_ = 2
    _fields_ = [("version", ctypes.c_uint16),
                ("data_type", ctypes.c_uint16),
                ("flags", ctypes.c_uint64),
                ("measurement_uid", ctypes.c_uint32),
                ("attribute_string_len", ctypes.c_uint32), ]


class POINT(ctypes.Structure):
    _fields_ = [("x", ctypes.c_int), ("y", ctypes.c_int)]

def main():
    with sandbox._binary._CodedOutputStream("test.bin") as s:
        w = sandbox._binary.DynamicNDArrayWriter(np.dtype(np.float32), sandbox._binary.write_float32, True)
        w(s, np.array([1, 2, 3, 4, 5, 6, 7, 8, 9], dtype=np.float32))

    os.system("hexdump -C test.bin")





if __name__ == "__main__":
    main()
