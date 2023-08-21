#! /usr/bin/env python3

from dataclasses import dataclass
import enum
import os
from typing import Any
import sandbox
import numpy as np
import numpy.typing as npt


def print_value(value: Any) -> None:
    print(f"{value} {type(value)} dtype={value.dtype if isinstance(value, np.ndarray) else None}")  # type: ignore


# file = "/tmp/sandbox_py.bin"

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


file = "/tmp/sandbox.ndjson"

with sandbox.NDJsonHelloNDJsonWriter(file) as w:
    w.write_a_boolean(True)
    w.write_a_boolean_stream([True, False, True])
    w.write_a_boolean_stream([True, False, True])
    w.write_an_enum(sandbox.MyEnum.B)
    w.write_some_flags(sandbox.MyFlags.NONE)
    w.write_an_optional_int_that_is_not_set(None)
    w.write_an_optional_int_that_is_set(123)
    w.write_a_union_with_simple_representation(sandbox.Int32OrBool.Bool(True))
    w.write_a_union_requiring_tag(None)
    w.write_a_record_with_optional_not_set(sandbox.MyRecord(x=123, y=456))
    w.write_a_record_with_optional_set(sandbox.MyRecord(x=123, y=456, z=789))

with open(file, "r") as f:
    print(f.read())

with sandbox.NDJsonHelloNDJsonReader(file) as r:
    print(r.read_a_boolean())
    for b in r.read_a_boolean_stream():
        print(f"stream: {b}")
    print(r.read_an_enum())
    print(r.read_some_flags())
    print(r.read_an_optional_int_that_is_not_set())
    print(r.read_an_optional_int_that_is_set())
    print(r.read_a_union_with_simple_representation())
    print(r.read_a_union_requiring_tag())
    print(r.read_a_record_with_optional_not_set())
    print(r.read_a_record_with_optional_set())


@dataclass
class MyClass:
    a: npt.NDArray[np.int32]
    b: int
    c: int | None


c = MyClass(np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32), 123, None)
