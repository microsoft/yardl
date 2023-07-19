#! /usr/bin/env python3

import datetime
import os
import time
from typing import (
    Any,
)
import sandbox
import numpy as np
import pandas as pd


def print_value(value: Any) -> None:
    print(f"{value} {type(value)} dtype={value.dtype if isinstance(value, np.ndarray) else None}")  # type: ignore


file = "sandbox_py.bin"

with sandbox.BinaryHelloWorldWriter(file) as w:

    def data_items_stream():
        yield np.array(
            [892.37889483 - 9932.485937837j, 73.383672763878 - 33.3394472537j],
            dtype=np.complex128,
        )
        yield np.array(
            [3883.22890980 + 373.4933837j, 56985.39384393 - 33833.3330128474373j],
            dtype=np.complex128,
        )
        yield np.array(
            [283.383672763878 - 33.3394472537j, 3883.22890980 + 373.4933837j],
            dtype=np.complex128,
        )

    w.write_data(data_items_stream())
    pass

with sandbox.BinaryHelloWorldReader(file, sandbox.Types.NONE) as r:
    value = r.read_data()
    print_value(list(value))

os.system(f"hexdump -C {file}")


print(np.datetime64("2023-07-18T19:19:17.732594999", "D").dtype)
print(time.gmtime(0))


fa = np.array([1, 2, 3], dtype=np.int32)
dt = np.dtype((fa.dtype, fa.shape))
print(f"dt={dt.shape}")
da = np.ndarray((2,), dtype=dt)
print(f"da.dtype={da.dtype}")
da[0:] = fa

print(da.dtype)
print(da.shape)
print(da)

with sandbox.BinaryHello2Writer(file) as w:
    w.write_data(da)
