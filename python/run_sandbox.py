#! /usr/bin/env python3

import datetime
import os
import io
import sys
import typing
import sandbox
import numpy as np


def main():
    with sandbox.BinaryP1Writer("test.bin") as w:
        w.write_my_value({"a": 1, "b": 2})
        pass

    os.system("hexdump -C test.bin")

    with sandbox.BinaryP1Reader("test.bin") as r:
        v = r.read_my_value()
        print(v)

if __name__ == "__main__":
    main()
