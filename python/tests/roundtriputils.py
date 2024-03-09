import io
import pathlib
import subprocess
import types
from typing import Callable, TypeVar, Union, cast

import test_model as tm
from .factories import Format
from .factories import get_reader_writer_types

# pyright: basic


_translator_path = (
    pathlib.Path(__file__).parent / "../../cpp/build/translator"
).resolve()


def invoke_translator(
    input: Union[bytes, str], input_format: Format, output_format: Format
) -> Union[bytes, str]:
    if isinstance(input, str):
        input = input.encode("utf-8")

    with subprocess.Popen(
        [_translator_path, str(input_format), str(output_format)],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
    ) as proc:
        cpp_output = proc.communicate(input=bytes(input))[0]
        assert proc.wait() == 0, "translator failed"
        if output_format == Format.NDJSON:
            cpp_output = cpp_output.decode("utf-8")
        return cpp_output


T = TypeVar("T")


def create_validating_writer_class(
    format: Format, base_class: type[T]
) -> Callable[[], T]:
    reader_class, writer_class = get_reader_writer_types(format, base_class)
    in_memory_stream_class = io.BytesIO if format == Format.BINARY else io.StringIO

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
                    if isinstance(args[1], types.GeneratorType) or isinstance(
                        args[1], range
                    ):
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
                    if isinstance(args[1], types.GeneratorType) or isinstance(
                        args[1], range
                    ):
                        arg_list = list(args)
                        arg_list[1] = list(args[1])
                        args = tuple(arg_list)
                    if method_snapshot.__name__ in recorded_args:
                        existing = recorded_args[method_snapshot.__name__]
                        recorded_args[method_snapshot.__name__] = existing + args[1]
                    else:
                        recorded_args[method_snapshot.__name__] = args[1]
                        existing = None

                    try:
                        return method_snapshot(*args, **kwargs)
                    except:
                        if existing is not None:
                            recorded_args[method_snapshot.__name__] = existing
                        else:
                            recorded_args.pop(method_snapshot.__name__)
                        raise

                return wrapper

            attrs[method.__name__] = mk_wrapper()

        def close_wrapper(*args, **kwargs):
            writer_class.close(  # pyright: ignore[reportGeneralTypeIssues]
                *args, **kwargs
            )

            self = args[0]
            this_buffer = self._buffer.getvalue()
            validating_instance = validating_class(in_memory_stream_class())
            validating_instance._recorded_arguments = (  # pyright: ignore[reportGeneralTypeIssues]
                self._recorded_arguments
            )

            validate(this_buffer, validating_instance)

        def validate(this_buffer, validating_instance):
            # 1. Validate that we get the same data when we read the output back in

            reader = reader_class(in_memory_stream_class(this_buffer))
            reader.copy_to(validating_instance)

            # 2. Now run the output though the C++ translator, and read its output back in
            # and validate that we get the same data
            cpp_output = invoke_translator(this_buffer, format, format)

            reader = reader_class(
                in_memory_stream_class(
                    cpp_output  # pyright: ignore[reportGeneralTypeIssues]
                )
            )
            reader.copy_to(validating_instance)

            # 3. Now use the translator to convert the output to the other format,
            # and then read it back in and make sure we get the same data.

            other_format = Format.NDJSON if format == Format.BINARY else Format.BINARY

            cpp_output_other_format = invoke_translator(
                this_buffer, format, other_format
            )

            other_format_reader_class = get_reader_writer_types(
                other_format, base_class
            )[0]
            in_memory_stream_class_other_format = (
                io.BytesIO if other_format == Format.BINARY else io.StringIO
            )

            reader = other_format_reader_class(
                in_memory_stream_class_other_format(
                    cpp_output_other_format  # pyright: ignore[reportGeneralTypeIssues]
                )
            )
            reader.copy_to(validating_instance)

        attrs["close"] = close_wrapper

        def init_wrapper(*args, **kwargs):
            recorded_args = {}
            args[0]._recorded_arguments = recorded_args
            buf = in_memory_stream_class()
            args[0]._buffer = buf
            return writer_class.__init__(
                args[0], buf, **kwargs
            )  # pyright: ignore[reportGeneralTypeIssues]

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
