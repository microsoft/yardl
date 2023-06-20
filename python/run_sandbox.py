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
import numpy.typing as npt


class Z(enum.Flag):
    A = 1
    B = 2
    C = 4

    @staticmethod
    def dtype() -> npt.DTypeLike:
        return np.int32


def main():
    with sandbox.BinaryP1Writer("test.bin") as w:

        value = np.array([((1, 2), (3, 4))], dtype=sandbox.Line.dtype(np.int32))
        print(f"{value} {type(value)}")
        w.write_my_value(value)
        pass

    # os.system("hexdump -C test.bin")

    with sandbox.BinaryP1Reader("test.bin", sandbox.Types.INTEGER) as r:
        v = r.read_my_value()
        print(f"{v} {type(v)}")

if __name__ == "__main__":
    main()
