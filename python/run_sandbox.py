#! /usr/bin/env python3

import abc
import datetime
import enum
import os
import io
import sys
import typing
import sandbox
import numpy as np


def main():
    with sandbox.BinaryP1Writer("test.bin") as w:

        dt = np.dtype([('x', np.int32), ('y', np.int32)])
        value = np.array([(1, 2), (3, 4)], dtype=dt)
        print(f"{value} {type(value)}")
        w.write_my_value(value)
        pass

    # os.system("hexdump -C test.bin")

    with sandbox.BinaryP1Reader("test.bin", sandbox.Types.INTEGER) as r:
        v = r.read_my_value()
        print(f"{v} {type(v)}")

if __name__ == "__main__":
    main()
