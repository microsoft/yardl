# This file was generated by the "yardl" tool. DO NOT EDIT.

# pyright: reportUnusedImport=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false
# pyright: reportUnknownVariableType=false

import datetime
import enum
import types
import typing

import numpy as np
import numpy.typing as npt

from . import yardl_types as yardl
from . import _dtypes


T = typing.TypeVar("T")
T_NP = typing.TypeVar("T_NP", bound=np.generic)


BinaryTree__: typing.TypeAlias = "BinaryTree"
LinkedList__: typing.TypeAlias = "LinkedList[T]"
Directory__: typing.TypeAlias = "Directory"

class BinaryTree:
    value: yardl.Int32
    left: BinaryTree__
    right: BinaryTree__

    def __init__(self, *,
        value: yardl.Int32 = 0,
        left: BinaryTree__ = None,
        right: BinaryTree__ = None,
    ):
        self.value = value
        self.left = left
        self.right = right

    def __eq__(self, other: object) -> bool:
        return (
            isinstance(other, BinaryTree)
            and self.value == other.value
            and self.left == other.left
            and self.right == other.right
        )

    def __str__(self) -> str:
        return f"BinaryTree(value={self.value}, left={self.left}, right={self.right})"

    def __repr__(self) -> str:
        return f"BinaryTree(value={repr(self.value)}, left={repr(self.left)}, right={repr(self.right)})"


class LinkedList(typing.Generic[T]):
    value: T
    next: LinkedList__

    def __init__(self, *,
        value: T,
        next: LinkedList__ = None,
    ):
        self.value = value
        self.next = next

    def __eq__(self, other: object) -> bool:
        return (
            isinstance(other, LinkedList)
            and yardl.structural_equal(self.value, other.value)
            and self.next == other.next
        )

    def __str__(self) -> str:
        return f"LinkedList(value={self.value}, next={self.next})"

    def __repr__(self) -> str:
        return f"LinkedList(value={repr(self.value)}, next={repr(self.next)})"


class File:
    name: str
    data: list[yardl.UInt8]

    def __init__(self, *,
        name: str = "",
        data: typing.Optional[list[yardl.UInt8]] = None,
    ):
        self.name = name
        self.data = data if data is not None else []

    def __eq__(self, other: object) -> bool:
        return (
            isinstance(other, File)
            and self.name == other.name
            and self.data == other.data
        )

    def __str__(self) -> str:
        return f"File(name={self.name}, data={self.data})"

    def __repr__(self) -> str:
        return f"File(name={repr(self.name)}, data={repr(self.data)})"


_T = typing.TypeVar('_T')

class DirectoryEntry:
    File: typing.ClassVar[type["DirectoryEntryUnionCase[File]"]]
    Directory: typing.ClassVar[type["DirectoryEntryUnionCase[Directory__]"]]

class DirectoryEntryUnionCase(DirectoryEntry, yardl.UnionCase[_T]):
    pass

DirectoryEntry.File = type("DirectoryEntry.File", (DirectoryEntryUnionCase,), {"index": 0, "tag": "File"})
DirectoryEntry.Directory = type("DirectoryEntry.Directory", (DirectoryEntryUnionCase,), {"index": 1, "tag": "Directory"})
del DirectoryEntryUnionCase

class Directory:
    name: str
    entries: list[DirectoryEntry]

    def __init__(self, *,
        name: str = "",
        entries: typing.Optional[list[DirectoryEntry]] = None,
    ):
        self.name = name
        self.entries = entries if entries is not None else []

    def __eq__(self, other: object) -> bool:
        return (
            isinstance(other, Directory)
            and self.name == other.name
            and self.entries == other.entries
        )

    def __str__(self) -> str:
        return f"Directory(name={self.name}, entries={self.entries})"

    def __repr__(self) -> str:
        return f"Directory(name={repr(self.name)}, entries={repr(self.entries)})"


def _mk_get_dtype():
    dtype_map: dict[typing.Union[type, types.GenericAlias], typing.Union[np.dtype[typing.Any], typing.Callable[[tuple[type, ...]], np.dtype[typing.Any]]]] = {}
    get_dtype = _dtypes.make_get_dtype_func(dtype_map)

    dtype_map.setdefault(BinaryTree, np.dtype([('value', np.dtype(np.int32)), ('left', np.dtype(np.object_)), ('right', np.dtype(np.object_))], align=True))
    dtype_map.setdefault(LinkedList, lambda type_args: np.dtype([('value', get_dtype(type_args[0])), ('next', np.dtype(np.object_))], align=True))
    dtype_map.setdefault(File, np.dtype([('name', np.dtype(np.object_)), ('data', np.dtype(np.object_))], align=True))
    dtype_map.setdefault(DirectoryEntry, np.dtype(np.object_))
    dtype_map.setdefault(Directory, np.dtype([('name', np.dtype(np.object_)), ('entries', np.dtype(np.object_))], align=True))

    return get_dtype

get_dtype = _mk_get_dtype()

