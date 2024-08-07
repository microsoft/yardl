import datetime
import typing
import sys

import numpy as np
import pytest
import test_model as tm

# pyright: basic


def test_defaulting():
    p = tm.RecordWithPrimitives()
    assert p.bool_field == False
    assert p.int32_field == 0
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
    # array_of_vectors is a fixed vector within a record, so it becomes a subarray of its inner dtype
    assert a.array_of_vectors.shape == (5,) and a.array_of_vectors.dtype == np.int32

    o = tm.RecordWithOptionalFields()
    assert o.optional_int == None

    c = tm.RecordWithUnionsOfContainers()
    assert c.map_or_scalar == tm.MapOrScalar.Map({})
    assert c.vector_or_scalar == tm.VectorOrScalar.Vector([])
    assert c.array_or_scalar == tm.ArrayOrScalar.Array(np.zeros((), dtype=np.int32))

    ag = tm.RecordWithAliasedGenerics()
    assert ag.my_strings.v1 == ""
    assert ag.my_strings.v2 == ""
    assert ag.aliased_strings.v1 == ""
    assert ag.aliased_strings.v2 == ""

    ## Test defaults for doubly nested generic records
    g1 = tm.RecordWithOptionalGenericField()
    assert g1.v == None
    g1a = tm.RecordWithAliasedOptionalGenericField()
    assert g1a.v == g1.v
    g2 = tm.RecordWithOptionalGenericUnionField()
    assert g2.v == None
    g2a = tm.RecordWithAliasedOptionalGenericUnionField()
    assert g2a.v == g2.v

    g4 = tm.RecordWithGenericVectors()
    assert g4.v == []
    assert g4.av == g4.v
    with pytest.raises(TypeError, match="missing 2 required keyword-only arguments"):
        g5 = tm.RecordWithGenericFixedVectors()  # type: ignore

    with pytest.raises(TypeError, match="missing 6 required keyword-only arguments"):
        g6 = tm.RecordWithGenericArrays()  # type: ignore

    g7 = tm.RecordWithGenericMaps()
    assert g7.m == {}
    assert g7.am == g7.m

    c = tm.RecordContainingNestedGenericRecords()
    assert c.f1 == g1
    assert c.f1a == g1a
    assert c.f2 == g2
    assert c.f2a == g2a
    assert c.nested.g1 == g1
    assert c.nested.g1a == g1a
    assert c.nested.g2 == g2
    assert c.nested.g2a == g2a
    assert c.nested.g3.v1 == ""
    assert c.nested.g3.v2 == 0
    assert c.nested.g3a.v1 == ""
    assert c.nested.g3a.v2 == 0
    assert c.nested.g4.v == []
    assert c.nested.g4.av == []
    assert type(c.nested.g5.fv) == list
    assert len(c.nested.g5.fv) == 3
    assert type(c.nested.g5.afv) == list
    assert len(c.nested.g5.afv) == 3

    assert np.array_equal(c.nested.g6.nd, np.zeros((0, 0)))
    assert c.nested.g6.nd.dtype == np.int32
    assert np.array_equal(c.nested.g6.fixed_nd, np.zeros((16, 8)))
    assert c.nested.g6.fixed_nd.dtype == np.int32
    assert np.array_equal(c.nested.g6.dynamic_nd, np.zeros(()))
    assert c.nested.g6.dynamic_nd.dtype == np.int32
    assert np.array_equal(c.nested.g6.nd, c.nested.g6.aliased_nd)
    assert np.array_equal(c.nested.g6.fixed_nd, c.nested.g6.aliased_fixed_nd)
    assert np.array_equal(c.nested.g6.dynamic_nd, c.nested.g6.aliased_dynamic_nd)

    assert c.nested.g7 == g7

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
    with pytest.raises(RuntimeError, match="Generic type arguments not provided"):
        tm.get_dtype(tm.AliasedTuple)
    with pytest.raises(RuntimeError, match="Generic type arguments not provided"):
        tm.get_dtype(tm.AliasedOpenGeneric)
    with pytest.raises(RuntimeError, match="Generic type arguments not provided"):
        tm.get_dtype(tm.AliasedGenericUnion2)
    with pytest.raises(RuntimeError, match="Generic type arguments not provided"):
        tm.get_dtype(tm.tuples.Tuple)

    assert tm.get_dtype(tm.MyTuple[int, str]) == tm.get_dtype(
        tm.basic_types.tuples.Tuple[int, str]
    )

    assert tm.get_dtype(tm.AliasedClosedGeneric) == np.dtype(
        [("v1", tm.get_dtype(tm.AliasedString)), ("v2", tm.get_dtype(tm.AliasedEnum))],
        align=True,
    )

    assert tm.get_dtype(tm.RecordWithAliasedGenerics) == np.dtype(
        [
            ("my_strings", tm.get_dtype(tm.MyTuple[str, str])),
            ("aliased_strings", tm.get_dtype(tm.AliasedTuple[str, str])),
        ],
        align=True,
    )

    assert tm.get_dtype(tm.RecordWithFixedArrays) == np.dtype(
        [
            ("ints", tm.get_dtype(tm.Int32), (2, 3)),
            ("fixed_simple_record_array", tm.get_dtype(tm.SimpleRecord), (3, 2)),
            ("fixed_record_with_vlens_array", tm.get_dtype(tm.RecordWithVlens), (2, 2)),
        ],
        align=True,
    )

    assert tm.get_dtype(tm.IntFixedArray) == np.object_
    assert tm.get_dtype(tm.SimpleRecordFixedArray) == np.object_
    assert tm.get_dtype(tm.RecordWithVlensFixedArray) == np.object_
    assert tm.get_dtype(tm.RecordWithNamedFixedArrays) == tm.get_dtype(
        tm.RecordWithFixedArrays
    )

    assert tm.get_dtype(tm.MyTuple[tm.Int32, tm.Float32]) == np.dtype(
        [("v1", "<i4"), ("v2", "<f4")], align=True
    )
    assert tm.get_dtype(tm.MyTuple[tm.SimpleRecord, tm.Float32]) == np.dtype(
        [("v1", tm.get_dtype(tm.SimpleRecord)), ("v2", "<f4")], align=True
    )

    if sys.version_info >= (3, 10):
        assert tm.get_dtype(tm.Int32 | None) == np.dtype(
            [("has_value", "?"), ("value", "<i4")], align=True
        )
        assert tm.get_dtype(tm.SimpleRecord | None) == np.dtype(
            [("has_value", "?"), ("value", tm.get_dtype(tm.SimpleRecord))], align=True
        )
        assert tm.get_dtype(tm.Int32 | tm.Float32) == np.object_

    assert tm.get_dtype(tm.AliasedGenericUnion2[tm.SimpleRecord, bool]) == np.object_

    assert tm.get_dtype(typing.Optional[tm.SimpleRecord]) == np.dtype(
        [("has_value", "?"), ("value", tm.get_dtype(tm.SimpleRecord))], align=True
    )

    assert tm.get_dtype(typing.Union[tm.Int32, tm.Float32]) == np.object_

    assert tm.get_dtype(tm.basic_types.Int32OrString) == np.object_
    assert tm.get_dtype(tm.basic_types.Int32OrString.Int32) == np.int32
    assert tm.get_dtype(tm.basic_types.Int32OrString.String) == np.object_

    assert tm.get_dtype(tm.basic_types.TimeOrDatetime) == np.object_
    assert tm.get_dtype(tm.basic_types.TimeOrDatetime.Time) == np.timedelta64
    assert tm.get_dtype(tm.basic_types.TimeOrDatetime.Datetime) == np.datetime64

    assert tm.get_dtype(tm.Int32OrSimpleRecord) == np.object_
    assert tm.get_dtype(tm.Int32OrSimpleRecord.Int32) == np.int32
    assert tm.get_dtype(tm.Int32OrSimpleRecord.SimpleRecord) == tm.get_dtype(
        tm.SimpleRecord
    )

    assert tm.get_dtype(tm.AliasedOptional) == np.dtype(
        [("has_value", "?"), ("value", np.int32)], align=True
    )

    assert tm.get_dtype(tm.AliasedMultiGenericOptional[str, int]) == np.object_
    assert tm.get_dtype(
        tm.RecordWithAliasedOptionalGenericUnionField[str, int]
    ) == np.dtype(
        [
            (
                "v",
                np.object_,
            )
        ],
        align=True,
    )

    assert tm.get_dtype(tm.AliasedNullableIntSimpleRecord) == np.object_
    assert tm.get_dtype(tm.AliasedNullableIntSimpleRecord.Int32) == np.int32
    assert tm.get_dtype(tm.AliasedNullableIntSimpleRecord.SimpleRecord) == tm.get_dtype(
        tm.SimpleRecord
    )
    assert (
        tm.get_dtype(typing.Optional[tm.AliasedNullableIntSimpleRecord]) == np.object_
    )

    # Standalone vectors, arrays, and maps all have dtype np.object_
    assert tm.get_dtype(tm.AliasedGenericVector[int]) == np.object_
    assert tm.get_dtype(tm.AliasedGenericFixedVector[int]) == np.object_
    assert tm.get_dtype(tm.AliasedGenericFixedArray[np.int32]) == np.object_
    assert tm.get_dtype(tm.AliasedGenericRank2Array[np.int32]) == np.object_
    assert tm.get_dtype(tm.AliasedGenericDynamicArray[np.int32]) == np.object_
    assert tm.get_dtype(tm.AliasedMap[str, int]) == np.object_
    assert (
        tm.get_dtype(tm.AliasedVectorOfGenericRecords[int, float, np.float32])
        == np.object_
    )
