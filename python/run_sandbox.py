#! /usr/bin/env python3

import os
import sandbox
import sandbox._binary
import numpy as np
import typing
from datetime import date, time, datetime
import sys
import ctypes

def main():
    with sandbox.BinaryP1Writer("test.bin") as w:
        w.write_an_int(8)
        w.write_a_stream([1, 2, 3, 4, 5])
        w.write_optional(2)
        w.write_union("oona")

    os.system("hexdump -C test.bin")



if __name__ == "__main__":
    main()
