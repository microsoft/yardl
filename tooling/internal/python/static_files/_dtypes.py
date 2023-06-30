import functools
from types import GenericAlias, UnionType
from typing import Any, Callable, cast, get_args, get_origin
import numpy as np
import numpy.typing as npt
from . import yardl_types as yardl


def make_get_dtype_func(dtype_map :dict[type | GenericAlias, np.dtype[Any] | Callable[[tuple[type, ...]], np.dtype[Any]]]) -> Callable[[type | GenericAlias], np.dtype[Any]]:
    dtype_map[yardl.Bool] = np.dtype(np.bool_)
    dtype_map[yardl.Int8] = np.dtype(np.int8)
    dtype_map[yardl.UInt8] = np.dtype(np.uint8)
    dtype_map[yardl.Int16] = np.dtype(np.int16)
    dtype_map[yardl.UInt16] = np.dtype(np.uint16)
    dtype_map[yardl.Int32] = np.dtype(np.int32)
    dtype_map[yardl.UInt32] = np.dtype(np.uint32)
    dtype_map[yardl.Int64] = np.dtype(np.int64)
    dtype_map[yardl.UInt64] = np.dtype(np.uint64)
    dtype_map[yardl.Size] = np.dtype(np.uint64)
    dtype_map[yardl.Float32] = np.dtype(np.float32)
    dtype_map[yardl.Float64] = np.dtype(np.float64)
    dtype_map[yardl.ComplexFloat] = np.dtype(np.complex64)
    dtype_map[yardl.ComplexDouble] = np.dtype(np.complex128)
    dtype_map[yardl.Date] = np.dtype(np.datetime64)
    dtype_map[yardl.Time] = np.dtype(np.timedelta64)
    dtype_map[yardl.DateTime] = np.dtype(np.datetime64)

    def get_dtype_impl(dtype_map : dict[type | GenericAlias, np.dtype[Any] | Callable[[tuple[type, ...]], np.dtype[Any]]], t: type | GenericAlias) -> np.dtype[Any]:
        if isinstance(t, type) or isinstance(t, UnionType):
            if (res := dtype_map.get(t, None)) is not None:
                if callable(res):
                    raise RuntimeError(f"Generic type arguments for {t} not provided")
                return res

        origin = get_origin(t)
        if origin == np.ndarray:
            if (res := dtype_map.get(cast(GenericAlias, t), None)) is not None:
                if callable(res):
                    raise RuntimeError(f"Unexpected generic type arguments for {t}")
                return res

        if origin is not None:
            if (res := dtype_map.get(origin, None)) is not None:
                if callable(res):
                    return res(get_args(t))


        raise RuntimeError(f"Cannot find dtype for {t}")


    return lambda t: get_dtype_impl(dtype_map, t)
