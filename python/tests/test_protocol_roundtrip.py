import datetime
import re
from typing import TypeVar

import numpy as np
import numpy.typing as npt
import pytest

import test_model as tm
from .roundtriputils import create_validating_writer_class, Format


T = TypeVar("T")


@pytest.fixture(scope="module", params=[Format.BINARY, Format.NDJSON])
def format(request: pytest.FixtureRequest):
    return request.param


def test_scalar_primitives(format: Format):
    with create_validating_writer_class(format, tm.ScalarsWriterBase)() as w:
        w.write_int32(42)
        rec = tm.RecordWithPrimitives(
            bool_field=True,
            int8_field=-88,
            uint8_field=88,
            int16_field=-1616,
            uint16_field=1616,
            int32_field=-3232,
            uint32_field=3232,
            int64_field=-64646464,
            uint64_field=64646464,
            size_field=64646464,
            float32_field=32.0,
            float64_field=64.64,
            complexfloat32_field=complex(32.0, 64.0),
            complexfloat64_field=64.64 + 32.32j,
            date_field=datetime.date(2024, 4, 2),
            time_field=tm.Time.from_components(12, 34, 56),
            datetime_field=tm.DateTime.from_components(
                2024, 4, 2, 12, 34, 56, 111222333
            ),
        )
        w.write_record(rec)


def test_scalar_optionals(format: Format):
    c = create_validating_writer_class(format, tm.ScalarOptionalsWriterBase)

    with c() as w:
        w.write_optional_int(None)
        w.write_optional_record(None)
        w.write_record_with_optional_fields(tm.RecordWithOptionalFields())
        w.write_optional_record_with_optional_fields(None)

    with c() as w:
        w.write_optional_int(55)
        w.write_optional_record(tm.SimpleRecord(x=8, y=9, z=10))
        w.write_record_with_optional_fields(
            tm.RecordWithOptionalFields(
                optional_int=44,
                optional_time=tm.Time.from_components(12, 34, 56),
            )
        )
        w.write_optional_record_with_optional_fields(
            tm.RecordWithOptionalFields(
                optional_int=12,
                optional_time=tm.Time.from_components(11, 32, 26),
            )
        )


def test_nested_records(format: Format):
    with create_validating_writer_class(format, tm.NestedRecordsWriterBase)() as w:
        w.write_tuple_with_records(
            tm.TupleWithRecords(
                a=tm.SimpleRecord(x=1, y=2, z=3), b=tm.SimpleRecord(x=4, y=5, z=6)
            )
        )


def test_variable_length_vectors(format: Format):
    with create_validating_writer_class(format, tm.VlensWriterBase)() as w:
        w.write_int_vector([1, 2, 3])
        w.write_complex_vector([complex(1, 2), complex(3, 4)])
        rec_with_vlens = tm.RecordWithVlens(
            a=[tm.SimpleRecord(x=1, y=2, z=3), tm.SimpleRecord(x=4, y=5, z=6)],
            b=4,
            c=2,
        )

        w.write_record_with_vlens(rec_with_vlens)
        w.write_vlen_of_record_with_vlens([rec_with_vlens, rec_with_vlens])


def test_strings(format: Format):
    with create_validating_writer_class(format, tm.StringsWriterBase)() as w:
        w.write_single_string("hello")
        w.write_rec_with_string(tm.RecordWithStrings(a="Montréal", b="臺北市"))


def test_optional_vectors(format: Format):
    c = create_validating_writer_class(format, tm.OptionalVectorsWriterBase)
    with c() as w:
        w.write_record_with_optional_vector(tm.RecordWithOptionalVector())

    with c() as w:
        w.write_record_with_optional_vector(
            tm.RecordWithOptionalVector(optional_vector=[1, 2, 3])
        )


def test_fixed_vectors(format: Format):
    with create_validating_writer_class(format, tm.FixedVectorsWriterBase)() as w:
        int_list: list[tm.Int32] = [1, 2, 3, 4, 5]
        w.write_fixed_int_vector(int_list)
        simple_rec_list = [
            tm.SimpleRecord(x=1, y=2, z=3),
            tm.SimpleRecord(x=4, y=5, z=6),
            tm.SimpleRecord(x=7, y=8, z=9),
        ]
        w.write_fixed_simple_record_vector(simple_rec_list)
        rec_with_vlens_list = [
            tm.RecordWithVlens(
                a=[tm.SimpleRecord(x=1, y=2, z=3), tm.SimpleRecord(x=4, y=5, z=6)],
                b=4,
                c=2,
            ),
            tm.RecordWithVlens(
                a=[tm.SimpleRecord(x=7, y=8, z=9), tm.SimpleRecord(x=10, y=11, z=12)],
                b=5,
                c=3,
            ),
        ]
        w.write_fixed_record_with_vlens_vector(rec_with_vlens_list)
        rec_with_fixed_list = tm.RecordWithFixedVectors(
            fixed_int_vector=int_list,
            fixed_simple_record_vector=simple_rec_list,
            fixed_record_with_vlens_vector=rec_with_vlens_list,
        )
        w.write_record_with_fixed_vectors(rec_with_fixed_list)


def test_fixed_arrays(format: Format):
    with create_validating_writer_class(format, tm.FixedArraysWriterBase)() as w:
        w.write_ints(np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32))
        simple_record_array = np.array(
            [
                [(1, 2, 3), (4, 5, 6)],
                [(11, 12, 13), (14, 15, 16)],
                [(21, 22, 23), (24, 25, 26)],
            ],
            dtype=tm.get_dtype(tm.SimpleRecord),
        )

        w.write_fixed_simple_record_array(simple_record_array)

        fixed_record_with_vlen_arrays = np.array(
            [
                [
                    (
                        [
                            tm.SimpleRecord(x=1, y=2, z=3),
                            tm.SimpleRecord(x=7, y=8, z=9),
                        ],
                        13,
                        14,
                    ),
                    (
                        [
                            tm.SimpleRecord(x=21, y=22, z=23),
                        ],
                        113,
                        114,
                    ),
                ],
                [
                    (
                        [
                            tm.SimpleRecord(x=31, y=32, z=33),
                            tm.SimpleRecord(x=34, y=35, z=36),
                            tm.SimpleRecord(x=37, y=38, z=39),
                        ],
                        213,
                        214,
                    ),
                    (
                        [
                            tm.SimpleRecord(x=41, y=42, z=43),
                        ],
                        313,
                        314,
                    ),
                ],
            ],
            dtype=tm.get_dtype(tm.RecordWithVlens),
        )

        w.write_fixed_record_with_vlens_array(fixed_record_with_vlen_arrays)

        w.write_record_with_fixed_arrays(
            tm.RecordWithFixedArrays(
                ints=np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32),
                fixed_simple_record_array=simple_record_array,
                fixed_record_with_vlens_array=fixed_record_with_vlen_arrays,
            )
        )

        # TODO: named fixed arrays are kind of broken since it
        # doesn't seem to be possible to specify the shape of the array in the type
        named_fixed_array = np.array([[1, 2, 3, 4], [5, 6, 7, 8]], dtype=np.int32)
        w.write_named_array(named_fixed_array)


def test_complex_arrays(format: Format):
    with create_validating_writer_class(format, tm.ComplexArraysWriterBase)() as w:
        fs = np.zeros((2, 16), dtype=np.complex64)
        w.write_floats(fs)
        ds = np.zeros((2, 16), dtype=np.complex128)
        w.write_doubles(ds)

    # Again but with arrays in Fortran order
    with create_validating_writer_class(format, tm.ComplexArraysWriterBase)() as w:
        fs = np.zeros((2, 16), dtype=np.complex64, order="F")
        w.write_floats(fs)
        ds = np.zeros((2, 16), dtype=np.complex128, order="F")
        w.write_doubles(ds)


def test_subarrays(format: Format):
    with create_validating_writer_class(format, tm.SubarraysWriterBase)() as w:
        with pytest.raises(
            ValueError, match=re.escape("The array is required to have shape (..., 3)")
        ):
            w.write_dynamic_with_fixed_int_subarray(
                np.array([[1, 2, 3, 4], [11, 12, 13, 14]], dtype=np.int32)
            )

        with pytest.raises(
            ValueError, match=re.escape("The array is required to have shape (..., 3)")
        ):
            w.write_dynamic_with_fixed_int_subarray(np.ndarray((), dtype=np.int32))

        w.write_dynamic_with_fixed_int_subarray(
            np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32)
        )
        w.write_dynamic_with_fixed_float_subarray(
            np.array([[1, 2, 3], [4, 5, 6]], dtype=np.float32)
        )

        with pytest.raises(
            ValueError, match=re.escape("The array is required to have shape (..., 3)")
        ):
            w.write_known_dim_count_with_fixed_int_subarray(
                np.array([[1, 2, 3, 4], [11, 12, 13, 14]], dtype=np.int32)
            )

        w.write_known_dim_count_with_fixed_int_subarray(
            np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32)
        )
        w.write_known_dim_count_with_fixed_float_subarray(
            np.array([[1, 2, 3], [4, 5, 6]], dtype=np.float32)
        )

        with pytest.raises(
            ValueError, match=re.escape("Expected shape (2, 3), got (2, 4)")
        ):
            w.write_fixed_with_fixed_int_subarray(
                np.array([[1, 2, 3, 4], [11, 12, 13, 14]], dtype=np.int32)
            )

        w.write_fixed_with_fixed_int_subarray(
            np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32)
        )
        w.write_fixed_with_fixed_float_subarray(
            np.array([[1, 2, 3], [4, 5, 6]], dtype=np.float32)
        )

        with pytest.raises(
            ValueError,
            match=re.escape("The array is required to have shape (..., 2, 3)"),
        ):
            w.write_nested_subarray(
                np.array(
                    [
                        [[1, 2, 3, 9], [4, 5, 6, 9]],
                        [[10, 20, 30, 90], [40, 50, 60, 90]],
                        [[100, 200, 300, 900], [400, 500, 600, 900]],
                    ],
                    dtype=np.int32,
                )
            )

        with pytest.raises(
            ValueError,
            match=re.escape("The array is required to have shape (..., 2, 3)"),
        ):
            w.write_nested_subarray(
                np.array(
                    [
                        [[1, 2, 3], [4, 5, 6], [7, 8, 9]],
                        [[10, 20, 30], [40, 50, 60], [70, 80, 90]],
                        [[100, 200, 300], [400, 500, 600], [700, 800, 900]],
                    ],
                    dtype=np.int32,
                )
            )

        w.write_nested_subarray(
            np.array(
                [
                    [[1, 2, 3], [4, 5, 6]],
                    [[10, 20, 30], [40, 50, 60]],
                    [[100, 200, 300], [400, 500, 600]],
                ],
                dtype=np.int32,
            )
        )

        w.write_dynamic_with_fixed_vector_subarray(
            np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32)
        )

        w.write_generic_subarray(
            np.array(
                [
                    [[1, 2, 3], [4, 5, 6]],
                    [[10, 12, 13], [14, 15, 16]],
                ],
                dtype=np.int32,
            )
        )


def test_subarrays_in_records(format: Format):
    fixed_dtype = tm.get_dtype(tm.RecordWithFixedCollections)
    assert fixed_dtype == np.dtype(
        [("fixed_vector", "<i4", (3,)), ("fixed_array", "<i4", (2, 3))], align=True
    )

    vlen_dtype = tm.get_dtype(tm.RecordWithVlenCollections)
    assert vlen_dtype == np.dtype([("vector", "O"), ("array", "O")], align=True)

    with create_validating_writer_class(format, tm.SubarraysInRecordsWriterBase)() as w:
        w.write_with_fixed_subarrays(
            np.array(
                [
                    ([1, 2, 3], [[11, 12, 13], [14, 15, 16]]),
                    ([101, 102, 103], [[1011, 1012, 1013], [1014, 15, 16]]),
                ],
                dtype=fixed_dtype,
            )
        )

        w.write_with_vlen_subarrays(
            np.array(
                [
                    (
                        np.array([1, 2, 3], dtype=np.int32),
                        np.array([[11, 12, 13], [14, 15, 16]], dtype=np.int32),
                    )
                ],
                dtype=vlen_dtype,
            )
        )


def test_arrays_with_known_dimension_count(format: Format):
    with create_validating_writer_class(format, tm.NDArraysWriterBase)() as w:
        w.write_ints(np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32))
        w.write_simple_record_array(
            np.array(
                [
                    [(1, 2, 3), (4, 5, 6)],
                    [(11, 12, 13), (14, 15, 16)],
                    [(21, 22, 23), (24, 25, 26)],
                ],
                dtype=tm.get_dtype(tm.SimpleRecord),
            )
        )

        w.write_record_with_vlens_array(
            np.array(
                [
                    [
                        (
                            [
                                tm.SimpleRecord(x=1, y=2, z=3),
                                tm.SimpleRecord(x=4, y=5, z=6),
                            ],
                            -7,
                            22,
                        ),
                        (
                            [
                                tm.SimpleRecord(x=1, y=2, z=3),
                                tm.SimpleRecord(x=4, y=5, z=6),
                            ],
                            -7,
                            22,
                        ),
                    ],
                    [
                        (
                            [
                                tm.SimpleRecord(x=1, y=2, z=3),
                                tm.SimpleRecord(x=4, y=5, z=6),
                            ],
                            -7,
                            22,
                        ),
                        (
                            [
                                tm.SimpleRecord(x=1, y=2, z=3),
                                tm.SimpleRecord(x=4, y=5, z=6),
                            ],
                            -7,
                            22,
                        ),
                    ],
                ],
                tm.get_dtype(tm.RecordWithVlens),
            )
        )

        w.write_record_with_nd_arrays(
            tm.RecordWithNDArrays(
                ints=np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32),
                fixed_simple_record_array=np.array(
                    [[(1, 2, 3)], [(4, 5, 6)]], dtype=tm.get_dtype(tm.SimpleRecord)
                ),
                fixed_record_with_vlens_array=np.array(
                    [
                        [
                            (
                                [
                                    tm.SimpleRecord(x=1, y=2, z=3),
                                    tm.SimpleRecord(x=4, y=5, z=6),
                                ],
                                -33,
                                44,
                            )
                        ],
                        [
                            (
                                [
                                    tm.SimpleRecord(x=1, y=2, z=3),
                                    tm.SimpleRecord(x=432, y=235, z=342),
                                ],
                                233,
                                347,
                            )
                        ],
                    ],
                    tm.get_dtype(tm.RecordWithVlens),
                ),
            )
        )

        w.write_named_array(np.array([[1, 2, 3], [4, 5, 6]], np.int32))


def test_dynamic_ndarrays(format: Format):
    with create_validating_writer_class(format, tm.DynamicNDArraysWriterBase)() as w:
        w.write_ints(np.ndarray((4, 3), dtype=np.int32))
        w.write_simple_record_array(
            np.ndarray((2, 3), dtype=tm.get_dtype(tm.SimpleRecord))
        )
        w.write_record_with_vlens_array(
            np.array(
                [
                    [
                        (
                            [tm.SimpleRecord(x=1, y=2, z=3)],
                            -33,
                            44,
                        )
                    ],
                    [
                        (
                            [
                                tm.SimpleRecord(x=8, y=2, z=9),
                                tm.SimpleRecord(x=28, y=3, z=34),
                            ],
                            233,
                            347,
                        )
                    ],
                ],
                tm.get_dtype(tm.RecordWithVlens),
            ),
        )

        w.write_record_with_dynamic_nd_arrays(
            tm.RecordWithDynamicNDArrays(
                ints=np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32),
                simple_record_array=np.array(
                    [[(1, 2, 3)], [(4, 5, 6)]], dtype=tm.get_dtype(tm.SimpleRecord)
                ),
                record_with_vlens_array=np.array(
                    [
                        [
                            (
                                [tm.SimpleRecord(x=1, y=2, z=3)],
                                -33,
                                44,
                            )
                        ],
                        [
                            (
                                [
                                    tm.SimpleRecord(x=8, y=2, z=9),
                                    tm.SimpleRecord(x=28, y=3, z=34),
                                ],
                                233,
                                347,
                            )
                        ],
                    ],
                    tm.get_dtype(tm.RecordWithVlens),
                ),
            )
        )


def test_maps(format: Format):
    with create_validating_writer_class(format, tm.MapsWriterBase)() as w:
        d = {"a": 1, "b": 2, "c": 3}
        w.write_string_to_int(d)
        w.write_int_to_string({1: "a", 2: "b", 3: "c"})
        w.write_string_to_union(
            {"a": tm.StringOrInt32.Int32(1), "b": tm.StringOrInt32.String("2")}
        )
        w.write_aliased_generic({"a": 1, "b": 2, "c": 3})
        w.write_records(
            [
                tm.RecordWithMaps(set_1={1: 1, 2: 2}, set_2={-1: True, 3: False}),
                tm.RecordWithMaps(set_1={1: 2, 2: 1}, set_2={-1: False, 3: True}),
            ]
        )


def test_unions(format: Format):
    c = create_validating_writer_class(format, tm.UnionsWriterBase)

    # first option
    with c() as w:
        w.write_int_or_simple_record(tm.Int32OrSimpleRecord.Int32(1))
        w.write_int_or_record_with_vlens(tm.Int32OrRecordWithVlens.Int32(2))
        w.write_monosotate_or_int_or_simple_record(None)
        w.write_record_with_unions(tm.basic_types.RecordWithUnions())

    # second option
    with c() as w:
        w.write_int_or_simple_record(
            tm.Int32OrSimpleRecord.SimpleRecord(tm.SimpleRecord(x=1, y=2, z=3))
        )
        w.write_int_or_record_with_vlens(
            tm.Int32OrRecordWithVlens.RecordWithVlens(
                tm.RecordWithVlens(a=[tm.SimpleRecord(x=1, y=2, z=3)], b=12, c=13)
            )
        )
        w.write_monosotate_or_int_or_simple_record(tm.Int32OrSimpleRecord.Int32(6))
        w.write_record_with_unions(
            tm.basic_types.RecordWithUnions(
                null_or_int_or_string=tm.basic_types.Int32OrString.Int32(7),
                date_or_datetime=tm.basic_types.TimeOrDatetime.Datetime(
                    tm.DateTime.from_components(2025, 3, 4),
                ),
                null_or_fruits_or_days_of_week=tm.basic_types.GenericNullableUnion2[
                    tm.basic_types.Fruits, tm.basic_types.DaysOfWeek
                ].T1(tm.basic_types.Fruits.APPLE),
            )
        )


def test_enums(format: Format):
    with create_validating_writer_class(format, tm.EnumsWriterBase)() as w:
        w.write_single(tm.Fruits.APPLE)
        w.write_vec([tm.Fruits.APPLE, tm.Fruits.BANANA, tm.Fruits(233983)])
        w.write_size(tm.SizeBasedEnum.C)
        w.write_rec(
            tm.RecordWithEnums(
                enum=tm.Fruits.PEAR,
                rec=tm.RecordWithNoDefaultEnum(enum=tm.Fruits.BANANA),
            )
        )


def test_flags(format: Format):
    def days():
        yield tm.DaysOfWeek.SUNDAY
        yield tm.DaysOfWeek.MONDAY | tm.DaysOfWeek.WEDNESDAY | tm.DaysOfWeek.FRIDAY
        yield tm.DaysOfWeek(0)
        yield tm.DaysOfWeek(282839)
        yield tm.DaysOfWeek(234532)

    with create_validating_writer_class(format, tm.FlagsWriterBase)() as w:
        w.write_days(days())
        w.write_formats(
            [
                tm.TextFormat.BOLD,
                tm.TextFormat.BOLD | tm.TextFormat.ITALIC,
                tm.TextFormat.REGULAR,
                tm.TextFormat(232932),
            ]
        )


def test_simple_streams(format: Format):
    c = create_validating_writer_class(format, tm.StreamsWriterBase)

    # non-empty streams
    with c() as w:
        w.write_int_data(range(10))
        w.write_int_data(range(20))

        w.write_optional_int_data([1, 2, None, 4, 5, None, 7, 8, 9, 10])
        w.write_record_with_optional_vector_data(
            [
                tm.RecordWithOptionalVector(),
                tm.RecordWithOptionalVector(optional_vector=[1, 2, 3]),
                tm.RecordWithOptionalVector(
                    optional_vector=[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
                ),
            ]
        )
        w.write_fixed_vector(([1, 2, 3] for _ in range(4)))

    # mixed empty and non-empty streams
    with c() as w:
        w.write_int_data(range(0))
        w.write_optional_int_data([1, 2, None, 4, 5, None, 7, 8, 9, 10])
        w.write_record_with_optional_vector_data([])
        w.write_fixed_vector(([1, 2, 3] for _ in range(4)))

    # empty streams
    with c() as w:
        w.write_int_data(range(0))
        w.write_optional_int_data([])
        w.write_record_with_optional_vector_data([])
        w.write_fixed_vector([])


def test_streams_of_unions(format: Format):
    with create_validating_writer_class(format, tm.StreamsOfUnionsWriterBase)() as w:
        w.write_int_or_simple_record(
            [
                tm.Int32OrSimpleRecord.Int32(1),
                tm.Int32OrSimpleRecord.SimpleRecord(tm.SimpleRecord(x=1, y=2, z=3)),
                tm.Int32OrSimpleRecord.Int32(2),
            ]
        )
        w.write_nullable_int_or_simple_record(
            [
                None,
                tm.Int32OrSimpleRecord.Int32(1),
                tm.Int32OrSimpleRecord.SimpleRecord(tm.SimpleRecord(x=1, y=2, z=3)),
                None,
                tm.Int32OrSimpleRecord.Int32(2),
                None,
            ]
        )
        w.write_many_cases(
            [
                tm.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.Int32(3),
                tm.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.Float32(7.0),
                tm.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.String(
                    "hello"
                ),
                tm.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.SimpleRecord(
                    tm.SimpleRecord(x=1, y=2, z=3)
                ),
                tm.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.NamedFixedNDArray(
                    np.array([[1, 2, 3, 4], [5, 6, 7, 8]], dtype=np.int32)
                ),
            ]
        )


def test_streams_of_aliased_unions(format: Format):
    with create_validating_writer_class(
        format, tm.StreamsOfAliasedUnionsWriterBase
    )() as w:
        w.write_int_or_simple_record(
            [
                tm.AliasedIntOrSimpleRecord.Int32(1),
                tm.AliasedIntOrSimpleRecord.SimpleRecord(
                    tm.SimpleRecord(x=1, y=2, z=3)
                ),
                tm.AliasedIntOrSimpleRecord.Int32(2),
            ]
        )
        w.write_nullable_int_or_simple_record(
            [
                None,
                tm.AliasedNullableIntSimpleRecord.Int32(1),
                tm.AliasedNullableIntSimpleRecord.SimpleRecord(
                    tm.SimpleRecord(x=1, y=2, z=3)
                ),
                None,
                tm.AliasedNullableIntSimpleRecord.Int32(2),
                None,
            ]
        )


def test_simple_generics(format: Format):
    with create_validating_writer_class(format, tm.SimpleGenericsWriterBase)() as w:
        w.write_float_image(np.array([[1.0, 2.0], [3.0, 4.0]], dtype=np.float32))
        w.write_int_image(np.array([[1, 2], [3, 4]], dtype=np.int32))
        w.write_int_image_alternate_syntax(np.array([[1, 2], [3, 4]], dtype=np.int32))
        w.write_string_image(np.array([["a", "b"], ["c", "d"]], dtype=np.object_))

        w.write_int_float_tuple(tm.MyTuple(v1=1, v2=2.0))
        w.write_float_float_tuple(tm.MyTuple(v1=1.0, v2=2.0))

        t = tm.MyTuple(v1=1, v2=2.0)
        w.write_int_float_tuple_alternate_syntax(t)

        w.write_int_string_tuple(tm.MyTuple(v1=1, v2="2"))

        w.write_stream_of_type_variants(
            [
                tm.ImageFloatOrImageDouble.ImageFloat(
                    np.array([[1.0, 2.0], [3.0, 4.0]], dtype=np.float32),
                ),
                tm.ImageFloatOrImageDouble.ImageDouble(
                    np.array([[1.0, 2.0], [3.0, 4.0]], dtype=np.float64),
                ),
            ]
        )


def test_advanced_generics(format: Format):
    with create_validating_writer_class(format, tm.AdvancedGenericsWriterBase)() as w:
        i1: tm.Image[np.float32] = np.array([[3, 4, 5], [6, 7, 8]], dtype=np.float32)
        i2: tm.Image[np.float32] = np.array(
            [[30, 40, 50], [60, 70, 80]], dtype=np.float32
        )
        i3: tm.Image[np.float32] = np.array(
            [[300, 400, 500], [600, 700, 800]], dtype=np.float32
        )
        i4: tm.Image[np.float32] = np.array(
            [[3000, 4000, 5000], [6000, 7000, 8000]], dtype=np.float32
        )

        img_img_array: npt.NDArray[np.object_] = np.ndarray((2, 2), dtype=np.object_)
        img_img_array[:] = [[i1, i2], [i3, i4]]

        w.write_float_image_image(img_img_array)

        w.write_generic_record_1(
            tm.GenericRecord(
                scalar_1=1,
                scalar_2="hello",
                vector_1=[1, 2, 3],
                image_2=np.array([["abc", "def"], ["a", "b"]], tm.get_dtype(str)),
            )
        )

        w.write_tuple_of_optionals(tm.MyTuple(v1=None, v2="hello"))
        w.write_tuple_of_optionals_alternate_syntax(tm.MyTuple(v1=34, v2=None))
        w.write_tuple_of_vectors(tm.MyTuple(v1=[1, 2, 3], v2=[4.0, 5.0, 6.0]))


def test_aliases(format: Format):
    with create_validating_writer_class(format, tm.AliasesWriterBase)() as w:
        w.write_aliased_string("hello")
        w.write_aliased_enum(tm.Fruits.APPLE)
        w.write_aliased_open_generic(
            tm.AliasedClosedGeneric(v1="hello", v2=tm.Fruits.BANANA)
        )
        w.write_aliased_closed_generic(
            tm.AliasedClosedGeneric(v1="hello", v2=tm.Fruits.PEAR)
        )
        w.write_aliased_optional(23)
        w.write_aliased_generic_optional(44.0)

        w.write_aliased_generic_union_2(
            tm.AliasedGenericUnion2[tm.AliasedString, tm.AliasedEnum].T1("hello")
        )

        w.write_aliased_generic_vector([1.0, 33.0, 44.0])
        w.write_aliased_generic_fixed_vector([1.0, 33.0, 44.0])

        w.write_stream_of_aliased_generic_union_2(
            [
                tm.AliasedGenericUnion2[tm.AliasedString, tm.AliasedEnum].T1("hello"),
                tm.AliasedGenericUnion2[tm.AliasedString, tm.AliasedEnum].T2(
                    tm.Fruits.APPLE
                ),
            ]
        )


def test_streams_of_unions_manual_close(format: Format):
    w = create_validating_writer_class(format, tm.StreamsOfUnionsWriterBase)()
    w.write_int_or_simple_record(
        [
            tm.Int32OrSimpleRecord.Int32(1),
            tm.Int32OrSimpleRecord.SimpleRecord(tm.SimpleRecord(x=1, y=2, z=3)),
            tm.Int32OrSimpleRecord.Int32(2),
        ]
    )
    w.write_nullable_int_or_simple_record(
        [
            None,
            tm.Int32OrSimpleRecord.Int32(1),
            tm.Int32OrSimpleRecord.SimpleRecord(tm.SimpleRecord(x=1, y=2, z=3)),
            None,
            tm.Int32OrSimpleRecord.Int32(2),
            None,
        ]
    )
    w.write_many_cases(
        [
            tm.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.Int32(3),
            tm.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.Float32(7.0),
            tm.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.String("hello"),
            tm.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.SimpleRecord(
                tm.SimpleRecord(x=1, y=2, z=3)
            ),
            tm.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.NamedFixedNDArray(
                np.array([[1, 2, 3, 4], [5, 6, 7, 8]], dtype=np.int32)
            ),
        ]
    )
    w.close()
