from enum import Enum
import inspect
from types import GenericAlias
from typing import Callable, TypeVar, cast


import test_model as tm
from test_model._binary import BinaryProtocolWriter


class Format(Enum):
    BINARY = 0
    NDJSON = 1
    HDF5 = 2

    def __str__(self) -> str:
        return self.name.lower()


_type_map = {
    base: (
        (
            cast(
                type,
                getattr(
                    tm,
                    derived.__name__.removesuffix("Writer") + "Reader",
                ),
            ),
            derived,
        ),
        (
            cast(
                type,
                getattr(
                    tm,
                    "NDJson"
                    + derived.__name__.removeprefix("Binary").removesuffix("Writer")
                    + "Reader",
                ),
            ),
            cast(type, getattr(tm, "NDJson" + derived.__name__.removeprefix("Binary"))),
        ),
    )
    for base, derived in {
        [base for base in inspect.getmro(derived) if base.__name__.endswith("Base")][
            0
        ]: cast(type, derived)
        for _, derived in inspect.getmembers(
            tm,
            lambda x: inspect.isclass(x)
            and not isinstance(x, GenericAlias)
            and issubclass(x, BinaryProtocolWriter),
        )
    }.items()
}


def get_reader_writer_types(
    format: Format, base_writer_class: type
) -> tuple[type, type]:
    return _type_map[base_writer_class][format.value]


T = TypeVar("T")


def get_writer_type(format: Format, base_writer_class: type[T]) -> Callable[[str], T]:
    return _type_map[base_writer_class][format.value][1]


def get_reader_type(format: Format, base_reader_class: type[T]) -> Callable[[str], T]:
    base_writer_class = getattr(
        tm, base_reader_class.__name__.removesuffix("ReaderBase") + "WriterBase"
    )
    return _type_map[base_writer_class][format.value][0]
