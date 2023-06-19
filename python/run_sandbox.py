#! /usr/bin/env python3

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
        value = datetime.datetime.now()
        print(f"{value} {type(value)}")

        w.write_my_value(value)
        pass

    # os.system("hexdump -C test.bin")

    with sandbox.BinaryP1Reader("test.bin", sandbox.Types.DATETIME) as r:
        v = r.read_my_value()
        print(f"{v} {type(v)}")






if __name__ == "__main__":
    main()
