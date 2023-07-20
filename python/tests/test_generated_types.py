import datetime
import typing

import numpy as np
import pytest
import test_model as tm


def test_defaulting():
    p = tm.RecordWithPrimitives()
    assert p.bool_field == False
    assert p.int_32_field == 0
    assert p.date_field == datetime.date(1970, 1, 1)
    assert p.time_field == tm.Time(0)
    assert p.datetime_field == tm.DateTime(0)

    v = tm.RecordWithVectors()
    assert v.default_vector == []
    assert v.default_vector_fixed_length == [0, 0, 0]
    assert v.vector_of_vectors == []

    a = tm.RecordWithArrays()
    assert a.default_array.shape == () and a.default_array.dtype == np.int32
    assert a.rank_1_array.shape == (0,)
    assert a.rank_2_array.shape == (0, 0)
    assert a.rank_2_array_with_named_dimensions.shape == (0, 0)
    assert a.rank_2_fixed_array.shape == (3, 4)
    assert a.rank_2_fixed_array_with_named_dimensions.shape == (3, 4)
    assert a.array_of_vectors.shape == (5,) and a.array_of_vectors.dtype == np.object_

    o = tm.RecordWithOptionalFields()
    assert o.optional_int == None

    ## Need to provide default values for generic fields
    with pytest.raises(TypeError):
        eval("tm.MyTuple()")

    with pytest.raises(TypeError):
        eval("tm.MyTuple[tm.Int32, tm.Float32]()")

    # The error goes away when you provide values for the fields.
    tm.MyTuple(v1=1, v2=2.0)


def test_get_dtype():
    assert tm.get_dtype(tm.Int32) == np.int32
    assert tm.get_dtype(bool) == np.bool_
    assert tm.get_dtype(int) == np.int64
    assert tm.get_dtype(str) == np.object_
    assert tm.get_dtype(tm.SimpleRecord) == np.dtype(
        [("x", "<i4"), ("y", "<i4"), ("z", "<i4")], align=True
    )
    assert tm.get_dtype(tm.TupleWithRecords) == np.dtype(
        [("a", tm.get_dtype(tm.SimpleRecord)), ("b", tm.get_dtype(tm.SimpleRecord))],
        align=True,
    )

    with pytest.raises(RuntimeError, match="Generic type arguments not provided"):
        tm.get_dtype(tm.MyTuple)
    assert tm.get_dtype(tm.MyTuple[tm.Int32, tm.Float32]) == np.dtype(
        [("v1", "<i4"), ("v2", "<f4")], align=True
    )
    assert tm.get_dtype(tm.MyTuple[tm.SimpleRecord, tm.Float32]) == np.dtype(
        [("v1", tm.get_dtype(tm.SimpleRecord)), ("v2", "<f4")], align=True
    )

    assert tm.get_dtype(tm.Int32 | None) == np.dtype(
        [("has_value", "?"), ("value", "<i4")], align=True
    )
    assert tm.get_dtype(tm.SimpleRecord | None) == np.dtype(
        [("has_value", "?"), ("value", tm.get_dtype(tm.SimpleRecord))], align=True
    )
    assert tm.get_dtype(typing.Optional[tm.SimpleRecord]) == np.dtype(
        [("has_value", "?"), ("value", tm.get_dtype(tm.SimpleRecord))], align=True
    )
    assert tm.get_dtype(tm.Int32 | tm.Float32) == np.object_
    assert tm.get_dtype(typing.Union[tm.Int32, tm.Float32]) == np.object_
