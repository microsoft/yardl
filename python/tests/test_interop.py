import datetime
import inspect
import io
import pathlib
import subprocess
import types
from typing import Callable, TypeVar, cast

import numpy as np
import test_model as tm
import test_model.binary as tmb
from test_model._binary import BinaryProtocolWriter


translator_path = (
    pathlib.Path(__file__).parent / "../../cpp/build/translator"
).resolve()


def invoke_translator(py_buf):
    with subprocess.Popen(
        [translator_path, "binary", "binary"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
    ) as proc:
        cpp_output = proc.communicate(input=py_buf)[0]
        assert proc.wait() == 0, "translator failed"

        assert cpp_output == py_buf


# base writer type -> (derived writer type, derived reader type)
type_map = {
    base: (
        derived,
        cast(
            type,
            getattr(
                inspect.getmodule(derived),
                derived.__name__.removesuffix("Writer") + "Reader",
            ),
        ),
    )
    for base, derived in {
        [base for base in inspect.getmro(derived) if base.__name__.endswith("Base")][
            0
        ]: cast(type, derived)
        for _, derived in inspect.getmembers(
            tm,
            lambda x: inspect.isclass(x) and issubclass(x, BinaryProtocolWriter),
        )
    }.items()
}


T = TypeVar("T")


def create_validating_writer_class(
    base_class: type[T],
) -> Callable[[], T]:
    writer_class, reader_class = type_map[base_class]

    write_methods = [
        cast(types.FunctionType, attr)
        for attr in [getattr(writer_class, name) for name in dir(writer_class)]
        if callable(attr) and attr.__name__.startswith("write")
    ]

    def create_validating_class() -> type[T]:
        attrs = {}
        for method in write_methods:

            def mk_wrapper(method_snapshot=method):
                def wrapper(*args, **kwargs):
                    recorded_args = args[0]._recorded_arguments
                    if isinstance(args[1], types.GeneratorType):
                        arg_list = list(args)
                        arg_list[1] = list(args[1])
                        args = tuple(arg_list)
                    assert tm.structural_equal(
                        recorded_args[method_snapshot.__name__], args[1]
                    )
                    return None

                return wrapper

            attrs[method.__name__] = mk_wrapper()

        return cast(
            type[T],
            types.new_class(
                "Validating" + writer_class.__name__,
                (writer_class,),
                {},
                lambda ns: ns.update(attrs),
            ),
        )

    validating_class = create_validating_class()

    def create_recording_class() -> type[T]:
        attrs = {}
        for method in write_methods:

            def mk_wrapper(method_snapshot=method):
                def wrapper(*args, **kwargs):
                    recorded_args = args[0]._recorded_arguments
                    recorded_args[method_snapshot.__name__] = args[1]
                    return method_snapshot(*args, **kwargs)

                return wrapper

            attrs[method.__name__] = mk_wrapper()

        def exit_wrapper(*args, **kwargs):
            result = writer_class.__exit__(*args, **kwargs)
            if args[1] is not None:
                # There was an exception, don't validate
                return result

            self = args[0]
            this_buffer = self._buffer.getvalue()
            validating_instance = validating_class(io.BytesIO())
            validating_instance._recorded_arguments = self._recorded_arguments

            # read as python types
            reader = reader_class(io.BytesIO(this_buffer), tm.Types.NONE)
            reader.copy_to(validating_instance)

            # now read as numpy types
            reader = reader_class(io.BytesIO(this_buffer), tm.Types.ALL)
            reader.copy_to(validating_instance)

            invoke_translator(this_buffer)

            return result

        attrs["__exit__"] = exit_wrapper

        def init_wrapper(*args, **kwargs):
            recorded_args = {}
            args[0]._recorded_arguments = recorded_args
            buf = io.BytesIO()
            args[0]._buffer = buf
            return writer_class.__init__(args[0], buf, **kwargs)

        attrs["__init__"] = init_wrapper

        return cast(
            type[T],
            types.new_class(
                "Recording" + writer_class.__name__,
                (writer_class,),
                {},
                lambda ns: ns.update(attrs),
            ),
        )

    return create_recording_class()


def test_scalar_primitives():
    with create_validating_writer_class(tm.ScalarsWriterBase)() as w:
        w.write_int_32(42)
        rec = tm.RecordWithPrimitives(
            bool_field=True,
            int_8_field=-88,
            uint_8_field=88,
            int_16_field=-1616,
            uint_16_field=1616,
            int_32_field=-3232,
            uint_32_field=3232,
            int_64_field=-64646464,
            uint_64_field=64646464,
            size_field=64646464,
            float_32_field=32.0,
            float_64_field=64.64,
            complexfloat_32_field=complex(32.0, 64.0),
            complexfloat_64_field=64.64 + 32.32j,
            date_field=datetime.date(2024, 4, 2),
            time_field=datetime.time(12, 34, 56),
            datetime_field=datetime.datetime(2024, 4, 2, 12, 34, 56, 111222),
        )
        w.write_record(rec)


def test_scalar_optionals():
    c = create_validating_writer_class(tm.ScalarOptionalsWriterBase)

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
                optional_int=44, optional_time=datetime.time(12, 34, 56)
            )
        )
        w.write_optional_record_with_optional_fields(
            tm.RecordWithOptionalFields(
                optional_int=12, optional_time=datetime.time(11, 32, 26)
            )
        )


def test_nested_records():
    with create_validating_writer_class(tm.NestedRecordsWriterBase)() as w:
        w.write_tuple_with_records(
            tm.TupleWithRecords(
                a=tm.SimpleRecord(x=1, y=2, z=3), b=tm.SimpleRecord(x=4, y=5, z=6)
            )
        )


def test_variable_length_vectors():
    with create_validating_writer_class(tm.VlensWriterBase)() as w:
        w.write_int_vector([1, 2, 3])
        w.write_complex_vector([complex(1, 2), complex(3, 4)])
        rec_with_vlens = tm.RecordWithVlens(
            a=[tm.SimpleRecord(x=1, y=2, z=3), tm.SimpleRecord(x=4, y=5, z=6)],
            b=4,
            c=2,
        )

        w.write_record_with_vlens(rec_with_vlens)
        w.write_vlen_of_record_with_vlens([rec_with_vlens, rec_with_vlens])


def test_strings():
    with create_validating_writer_class(tm.StringsWriterBase)() as w:
        w.write_single_string("hello")
        w.write_rec_with_string(tm.RecordWithStrings(a="Montréal", b="臺北市"))


def test_optional_vectors():
    c = create_validating_writer_class(tm.OptionalVectorsWriterBase)
    with c() as w:
        w.write_record_with_optional_vector(tm.RecordWithOptionalVector())

    with c() as w:
        w.write_record_with_optional_vector(
            tm.RecordWithOptionalVector(optional_vector=[1, 2, 3])
        )


def test_fixed_vectors():
    with create_validating_writer_class(tm.FixedVectorsWriterBase)() as w:
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


def test_fixed_arrays():
    with create_validating_writer_class(tm.FixedArraysWriterBase)() as w:
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

        # TODO: Note the inner lists of the record classes, not the tuples!
        # If the inner vector were fixed, it would be treated as a subarray.
        # Not sure that's best in this cases.
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


def test_subarrays():
    fixed_dtype = tm.get_dtype(tm.RecordWithFixedCollections)
    assert fixed_dtype == np.dtype(
        [("fixed_vector", "<i4", (3,)), ("fixed_array", "<i4", (2, 3))], align=True
    )

    vlen_dtype = tm.get_dtype(tm.RecordWithVlenCollections)
    assert vlen_dtype == np.dtype(
        [("fixed_vector", "O"), ("fixed_array", "O")], align=True
    )

    with create_validating_writer_class(tm.SubarraysWriterBase)() as w:
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


def test_arrays_with_known_dimension_count():
    with create_validating_writer_class(tm.NDArraysWriterBase)() as w:
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

        simple_record_type = tm.get_dtype(tm.SimpleRecord)
        w.write_record_with_vlens_array(
            np.array(
                [
                    [
                        (
                            np.array([(1, 2, 3), (4, 5, 6)], simple_record_type),
                            -7,
                            22,
                        ),
                        (
                            np.array([(1, 2, 3), (4, 5, 6)], simple_record_type),
                            -7,
                            22,
                        ),
                    ],
                    [
                        (
                            np.array([(1, 2, 3), (4, 5, 6)], simple_record_type),
                            -7,
                            22,
                        ),
                        (
                            np.array([(1, 2, 3), (4, 5, 6)], simple_record_type),
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
                                np.array([(1, 2, 3)], tm.get_dtype(tm.SimpleRecord)),
                                -33,
                                44,
                            )
                        ],
                        [
                            (
                                np.array(
                                    [(8, 2, 9), (28, 43, 9)],
                                    tm.get_dtype(tm.SimpleRecord),
                                ),
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


def test_dynamic_ndarrays():
    with create_validating_writer_class(tm.DynamicNDArraysWriterBase)() as w:
        w.write_ints(np.ndarray((4, 3), dtype=np.int32))
        w.write_simple_record_array(
            np.ndarray((2, 3), dtype=tm.get_dtype(tm.SimpleRecord))
        )
        w.write_record_with_vlens_array(
            np.array(
                [
                    [
                        (
                            np.array([(1, 2, 3)], tm.get_dtype(tm.SimpleRecord)),
                            -33,
                            44,
                        )
                    ],
                    [
                        (
                            np.array(
                                [(8, 2, 9), (28, 43, 9)],
                                tm.get_dtype(tm.SimpleRecord),
                            ),
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
                                np.array([(1, 2, 3)], tm.get_dtype(tm.SimpleRecord)),
                                -33,
                                44,
                            )
                        ],
                        [
                            (
                                np.array(
                                    [(8, 2, 9), (28, 43, 9)],
                                    tm.get_dtype(tm.SimpleRecord),
                                ),
                                233,
                                347,
                            )
                        ],
                    ],
                    tm.get_dtype(tm.RecordWithVlens),
                ),
            )
        )
