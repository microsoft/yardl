#! /usr/bin/env python3

import os
import io
import sys
import typing
import sandbox
import numpy as np


def main():
    with sandbox.BinaryP1Writer("test.bin") as w:
        w.write_arr(np.array([1, 2, 3], dtype=np.uint32))

        pass


    os.system("hexdump -C test.bin")




if __name__ == "__main__":
    main()
