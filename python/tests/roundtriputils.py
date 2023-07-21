import inspect
import io
import pathlib
import subprocess
import types
from typing import Callable, TypeVar, cast

import test_model as tm
from test_model._binary import BinaryProtocolWriter

# pyright: basic

_translator_path = (
    pathlib.Path(__file__).parent / "../../cpp/build/translator"
).resolve()


def invoke_translator(py_buf: bytes) -> bytes:
    with subprocess.Popen(
        [_translator_path, "binary", "binary"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
    ) as proc:
        cpp_output = proc.communicate(input=py_buf)[0]
        assert proc.wait() == 0, "translator failed"
        return cpp_output


# base writer type -> (derived writer type, derived reader type)
_type_map = {
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
            lambda x: inspect.isclass(x)
            and not isinstance(x, types.GenericAlias)
            and issubclass(x, BinaryProtocolWriter),
        )
    }.items()
}


T = TypeVar("T")


def create_validating_writer_class(
    base_class: type[T],
) -> Callable[[], T]:
    writer_class, reader_class = _type_map[base_class]

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
                    recorded_args[method_snapshot.__name__] = args[1]
                    return method_snapshot(*args, **kwargs)

                return wrapper

            attrs[method.__name__] = mk_wrapper()

        def exit_wrapper(*args, **kwargs):
            result = writer_class.__exit__(  # pyright: ignore[reportGeneralTypeIssues]
                *args, **kwargs
            )
            if args[1] is not None:
                # There was an exception, don't validate
                return result

            self = args[0]
            this_buffer = self._buffer.getvalue()
            validating_instance = validating_class(io.BytesIO())
            validating_instance._recorded_arguments = (  # pyright: ignore[reportGeneralTypeIssues]
                self._recorded_arguments
            )

            reader = reader_class(io.BytesIO(this_buffer))
            reader.copy_to(validating_instance)

            cpp_output = invoke_translator(this_buffer)
            reader = reader_class(io.BytesIO(cpp_output))
            reader.copy_to(validating_instance)

            return result

        attrs["__exit__"] = exit_wrapper

        def init_wrapper(*args, **kwargs):
            recorded_args = {}
            args[0]._recorded_arguments = recorded_args
            buf = io.BytesIO()
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
