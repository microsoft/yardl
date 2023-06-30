import datetime

import numpy as np
import pytest
import test_model as tm

def test_defaulting():
    p = tm.RecordWithPrimitives()
    assert p.bool_field == False
    assert p.int_32_field == 0
    assert p.date_field == datetime.date(1970, 1, 1)
    assert p.time_field == datetime.time(0, 0, 0)
    assert p.datetime_field == datetime.datetime(1970, 1, 1, 0, 0, 0)

    v = tm.RecordWithVectors()
    assert v.default_vector == []
    assert v.default_vector_fixed_length == [0, 0, 0]
    assert v.vector_of_vectors == []

    a = tm.RecordWithArrays()
    assert a.default_array.shape == () and a.default_array.dtype == np.int32
    assert a.rank_1_array.shape == (0,)
    assert a.rank_2_array.shape == (0,0)
    assert a.rank_2_array_with_named_dimensions.shape == (0,0)
    assert a.rank_2_fixed_array.shape == (3,4)
    assert a.rank_2_fixed_array_with_named_dimensions.shape == (3,4)
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