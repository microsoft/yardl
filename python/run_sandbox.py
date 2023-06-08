#! /usr/bin/env python3

import sandbox
import numpy as np
import typing
from datetime import date, time, datetime

def main():
    e = sandbox.MyFlags.A | sandbox.MyFlags.B
    print(sandbox.MyFlags.A in e)

    r = sandbox.MyRec[int](f2=1, f3=e)
    print(r)


if __name__ == "__main__":
    main()
