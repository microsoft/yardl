import datetime
import numpy as np
import test_model as tm


def test_simple_equality():
    a = tm.SimpleRecord(x=1, y=2, z=3)
    b = tm.SimpleRecord(x=1, y=2, z=3)
    assert a == b

    c = tm.SimpleRecord(x=1, y=2, z=4)
    assert a != c

    b = tm.SimpleRecord(x=np.int32(1), y=np.int32(2), z=np.int32(3))
    assert a == b


def test_date_equality():
    a = tm.RecordWithPrimitives(date_field=datetime.date(2020, 1, 1))
    b = tm.RecordWithPrimitives(
        date_field=np.datetime64(datetime.date(2020, 1, 1), "ns")
    )
    assert a == b

    c = tm.RecordWithPrimitives(date_field=np.datetime64(datetime.date(2020, 1, 2)))
    assert a != c
    assert b != c


def test_datetime_equality():
    a = tm.RecordWithPrimitives(
        datetime_field=datetime.datetime(2020, 1, 1, 1, 1, 1, 1)
    )
    b = tm.RecordWithPrimitives(
        datetime_field=np.datetime64(datetime.datetime(2020, 1, 1, 1, 1, 1, 1), "ns")
    )
    assert a == b

    c = tm.RecordWithPrimitives(
        datetime_field=np.datetime64(datetime.datetime(2020, 1, 1, 1, 1, 1, 2))
    )
    assert a != c
    assert b != c


def test_time_equality():
    a = tm.RecordWithPrimitives(time_field=datetime.time(1, 1, 1, 1))
    b = tm.RecordWithPrimitives(
        time_field=np.timedelta64(
            datetime.timedelta(hours=1, minutes=1, seconds=1, microseconds=1)
        )
    )
    assert a == b
    assert b == a

    c = tm.RecordWithPrimitives(
        time_field=np.timedelta64(
            datetime.timedelta(hours=1, minutes=1, seconds=1, microseconds=2)
        )
    )
    assert a != c
    assert c != a
    assert b != c
    assert c != b


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
    a = tm.RecordWithOptionalFields(optional_time=datetime.time(1, 1, 1, 1))
    b = tm.RecordWithOptionalFields(
        optional_time=np.timedelta64(
            datetime.timedelta(hours=1, minutes=1, seconds=1, microseconds=1)
        )
    )
    assert a == b

    c = tm.RecordWithOptionalFields(
        optional_time=np.timedelta64(
            datetime.timedelta(hours=1, minutes=1, seconds=1, microseconds=2)
        )
    )
    assert a != c
    assert b != c

    d = tm.RecordWithOptionalFields(optional_time=None)
    e = tm.RecordWithOptionalFields(optional_time=None)
    assert d == e
    assert a != d


def test_time_vector_equality():
    a = tm.RecordWithVectorOfTimes(
        times=[datetime.time(1, 1, 1, 1), datetime.time(1, 1, 1, 1)]
    )
    b = tm.RecordWithVectorOfTimes(
        times=[
            np.timedelta64(
                datetime.timedelta(hours=1, minutes=1, seconds=1, microseconds=1)
            ),
            np.timedelta64(
                datetime.timedelta(hours=1, minutes=1, seconds=1, microseconds=1)
            ),
        ]
    )
    assert a == b
    c = tm.RecordWithVectorOfTimes(
        times=[datetime.time(1, 1, 1, 1), datetime.time(1, 1, 1, 2)]
    )
    assert a != c


def test_simple_array_equality():
    a = tm.RecordWithArrays(default_array=np.array([1, 2, 3], dtype=np.int32))
    b = tm.RecordWithArrays(default_array=np.array([1, 2, 3], dtype=np.int32))
    assert a == b

    c = tm.RecordWithArrays(default_array=np.array([1, 2, 4], dtype=np.int32))
    assert a != c


def test_simple_union_equality():
    a = tm.RecordWithUnions(null_or_int_or_string=None)
    b = tm.RecordWithUnions(null_or_int_or_string=None)
    assert a == b

    c = tm.RecordWithUnions(null_or_int_or_string=("int32", 1))
    d = tm.RecordWithUnions(null_or_int_or_string=("int32", 1))
    assert c == d
    assert a != c

    e = tm.RecordWithUnions(null_or_int_or_string=("string", "hello"))
    f = tm.RecordWithUnions(null_or_int_or_string=("string", "hello"))
    assert e == f
    assert a != e
    assert c != e


def test_time_union_equality():
    a = tm.RecordWithUnions(date_or_datetime=("time", datetime.time(1, 1, 1, 1)))
    b = tm.RecordWithUnions(
        date_or_datetime=(
            "time",
            np.timedelta64(
                datetime.timedelta(hours=1, minutes=1, seconds=1, microseconds=1)
            ),
        )
    )

    assert a == b


def test_enum_equality():
    a = tm.RecordWithEnums(
        enum=tm.Fruits.APPLE, flags=tm.DaysOfWeek.SATURDAY | tm.DaysOfWeek.SUNDAY
    )
    b = tm.RecordWithEnums(
        enum=tm.Fruits.APPLE, flags=tm.DaysOfWeek.SATURDAY | tm.DaysOfWeek.SUNDAY
    )
    assert a == b

    c = tm.RecordWithEnums(enum=tm.Fruits.APPLE, flags=tm.DaysOfWeek.SATURDAY)
    assert a != c


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
