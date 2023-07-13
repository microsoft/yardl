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
