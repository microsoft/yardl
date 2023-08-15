#! /usr/bin/env python3

import os
from typing import Any
import sandbox
import numpy as np


# def print_value(value: Any) -> None:
#     print(f"{value} {type(value)} dtype={value.dtype if isinstance(value, np.ndarray) else None}")  # type: ignore


# file = "sandbox_py.bin"

# with sandbox.BinaryHelloWorldWriter(file) as w:

#     def data_items_stream():
#         yield np.array(
#             [892.37889483 - 9932.485937837j, 73.383672763878 - 33.3394472537j],
#             dtype=np.complex128,
#         )
#         yield np.array(
#             [3883.22890980 + 373.4933837j, 56985.39384393 - 33833.3330128474373j],
#             dtype=np.complex128,
#         )
#         yield np.array(
#             [283.383672763878 - 33.3394472537j, 3883.22890980 + 373.4933837j],
#             dtype=np.complex128,
#         )

#     w.write_data(data_items_stream())
#     pass

# with sandbox.BinaryHelloWorldReader(file) as r:
#     value = r.read_data()
#     print_value(list(value))

# os.system(f"hexdump -C {file}")


import time

from rich.live import Live
from rich.table import Table

table = Table(title="Throughput in MiB/s (higher is better)")
table.add_column("Scenario")
table.add_column("Format")
table.add_column("C++ Write", justify="right")
table.add_column("C++ Read", justify="right")
table.add_column("Python Write", justify="right")
table.add_column("Python Read", justify="right")

with Live(table, auto_refresh=False) as l:  # update 4 times a second to feel fluid
    for row in range(12):
        time.sleep(0.4)  # arbitrary delay
        # update the renderable internally
        table.add_row(f"{row}", f"description {row}", "[red]ERROR")
        l.refresh()
