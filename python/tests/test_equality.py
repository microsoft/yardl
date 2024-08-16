import datetime
import numpy as np
import test_model as tm


def test_simple_equality():
    a = tm.SimpleRecord(x=1, y=2, z=3)
    b = tm.SimpleRecord(x=1, y=2, z=3)
    assert a == b

    c = tm.SimpleRecord(x=1, y=2, z=4)
    assert a != c


def test_flags_equality():
    a = tm.DaysOfWeek.MONDAY | tm.DaysOfWeek.TUESDAY
    b = tm.DaysOfWeek.MONDAY | tm.DaysOfWeek.TUESDAY
    assert a == b

    c = tm.DaysOfWeek(0)
    d = tm.DaysOfWeek(0)
    assert c == d
    assert a != c

    e = tm.DaysOfWeek(0xFFFF)
    f = tm.DaysOfWeek(0xFFFF)
    assert e == f


def test_enum_equality():
    a = tm.Fruits.APPLE
    b = tm.Fruits.APPLE
    assert a == b
    assert hash(a) == hash(b)

    c = tm.Fruits(10000)
    d = tm.Fruits(10000)
    assert c == d


def test_record_with_enum_equality():
    a = tm.RecordWithEnums(
        enum=tm.Fruits.APPLE,
        flags=tm.DaysOfWeek.SATURDAY | tm.DaysOfWeek.SUNDAY,
        rec=tm.RecordWithNoDefaultEnum(enum=tm.Fruits.PEAR),
    )
    b = tm.RecordWithEnums(
        enum=tm.Fruits.APPLE,
        flags=tm.DaysOfWeek.SATURDAY | tm.DaysOfWeek.SUNDAY,
        rec=tm.RecordWithNoDefaultEnum(enum=tm.Fruits.PEAR),
    )
    assert a == b

    c = tm.RecordWithEnums(
        enum=tm.Fruits.APPLE,
        flags=tm.DaysOfWeek.SATURDAY,
        rec=tm.RecordWithNoDefaultEnum(enum=tm.Fruits.PEAR),
    )
    assert a != c

    d = tm.RecordWithEnums(
        enum=tm.Fruits.APPLE,
        flags=tm.DaysOfWeek.SATURDAY | tm.DaysOfWeek.SUNDAY,
        rec=tm.RecordWithNoDefaultEnum(enum=tm.Fruits.APPLE),
    )
    assert a != d


def test_date_equality():
    a = tm.RecordWithPrimitives(date_field=datetime.date(2020, 1, 1))
    b = tm.RecordWithPrimitives(date_field=datetime.date(2020, 1, 1))
    assert a == b

    c = tm.RecordWithPrimitives(date_field=datetime.date(2020, 1, 2))
    assert a != c
    assert b != c


def test_string_equality():
    a = tm.RecordWithStrings(a="a", b="b")
    b = tm.RecordWithStrings(a="a", b="b")
    assert a == b

    c = tm.RecordWithStrings(a="a", b="c")
    assert a != c


def test_record_with_primitive_vectors_equality():
    a = tm.RecordWithVectors(
        default_vector=[1, 2],
        default_vector_fixed_length=[1, 2, 3],
        vector_of_vectors=[[1, 2], [3, 4]],
    )

    b = tm.RecordWithVectors(
        default_vector=[1, 2],
        default_vector_fixed_length=[1, 2, 3],
        vector_of_vectors=[[1, 2], [3, 4]],
    )

    assert a == b


def test_optional_int_equality():
    a = tm.RecordWithOptionalFields(optional_int=1)
    b = tm.RecordWithOptionalFields(optional_int=1)
    assert a == b

    c = tm.RecordWithOptionalFields(optional_int=2)
    assert a != c
    d = tm.RecordWithOptionalFields(optional_int=None)
    e = tm.RecordWithOptionalFields(optional_int=None)
    assert d == e
    assert a != d


def test_optional_time_equality():
    a = tm.RecordWithOptionalFields(optional_time=tm.Time.from_components(1, 1, 1, 1))
    b = tm.RecordWithOptionalFields(optional_time=tm.Time.from_components(1, 1, 1, 1))

    assert a == b

    c = tm.RecordWithOptionalFields(optional_time=tm.Time.from_components(1, 1, 1, 2))
    assert a != c
    assert b != c

    d = tm.RecordWithOptionalFields(optional_time=None)
    e = tm.RecordWithOptionalFields(optional_time=None)
    assert d == e
    assert a != d


def test_time_vector_equality():
    a = tm.RecordWithVectorOfTimes(
        times=[tm.Time.from_components(1, 1, 1, 1), tm.Time.from_components(1, 1, 1, 1)]
    )
    b = tm.RecordWithVectorOfTimes(
        times=[tm.Time.from_components(1, 1, 1, 1), tm.Time.from_components(1, 1, 1, 1)]
    )
    assert a == b
    c = tm.RecordWithVectorOfTimes(
        times=[tm.Time.from_components(1, 1, 1, 1), tm.Time.from_components(1, 1, 1, 2)]
    )
    assert a != c


def test_simple_array_equality():
    a = tm.RecordWithArrays(default_array=np.array([1, 2, 3], dtype=np.int32))
    b = tm.RecordWithArrays(default_array=np.array([1, 2, 3], dtype=np.int32))
    assert a == b

    c = tm.RecordWithArrays(default_array=np.array([1, 2, 4], dtype=np.int32))
    assert a != c


def test_array_list_equality():
    a = np.array([1, 2, 3], np.int32)
    b = [1, 2, 3]
    assert tm.structural_equal(a, b)
    assert tm.structural_equal(b, a)


def test_array_structural_equality_with_object():
    a = np.array(
        [(1, 2, "hello"), (3, 4, "world")],
        dtype=[("a", np.int32), ("b", np.int32), ("c", np.object_)],
    )
    b = a
    assert tm.structural_equal(a, b)
    assert tm.structural_equal(b, a)


def test_simple_union_equality():
    a = tm.basic_types.RecordWithUnions(null_or_int_or_string=None)
    b = tm.basic_types.RecordWithUnions(null_or_int_or_string=None)
    assert a == b

    c = tm.basic_types.RecordWithUnions(
        null_or_int_or_string=tm.basic_types.Int32OrString.Int32(1)
    )
    d = tm.basic_types.RecordWithUnions(
        null_or_int_or_string=tm.basic_types.Int32OrString.Int32(1)
    )
    assert c == d
    assert a != c

    e = tm.basic_types.RecordWithUnions(
        null_or_int_or_string=tm.basic_types.Int32OrString.String("hello")
    )
    f = tm.basic_types.RecordWithUnions(
        null_or_int_or_string=tm.basic_types.Int32OrString.String("hello")
    )
    assert e == f
    assert a != e
    assert c != e


def test_time_union_equality():
    a = tm.basic_types.RecordWithUnions(
        date_or_datetime=tm.basic_types.TimeOrDatetime.Time(
            tm.Time.from_components(1, 1, 1, 1)
        )
    )
    b = tm.basic_types.RecordWithUnions(
        date_or_datetime=tm.basic_types.TimeOrDatetime.Time(
            tm.Time.from_components(1, 1, 1, 1)
        )
    )

    assert a == b


def test_generic_equality():
    a = tm.GenericRecord(
        scalar_1=1,
        scalar_2=2.0,
        vector_1=[1, 2, 3],
        image_2=np.array([[1.1, 2.2], [3.3, 4.4]], np.float32),
    )
    b = tm.GenericRecord(
        scalar_1=1,
        scalar_2=2.0,
        vector_1=[1, 2, 3],
        image_2=np.array([[1.1, 2.2], [3.3, 4.4]], np.float32),
    )

    assert a == b

    a = tm.MyTuple(v1=42.0, v2="hello, world")
    b = tm.basic_types.tuples.Tuple(v1=42.0, v2="hello, world")

    assert a == b

    a = tm.MyTuple(v1=42.0, v2="hello, world")
    b = tm.AliasedTuple(v1=42.0, v2="hello, world")
