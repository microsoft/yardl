import io
import pytest

import test_model as tm
from .factories import Format, get_reader_writer_types


@pytest.fixture(scope="module", params=[Format.BINARY, Format.NDJSON])
def format(request: pytest.FixtureRequest) -> Format:
    return request.param


def create_testing_reader(
    format: Format, base_class: type
) -> tuple[io.BytesIO | io.StringIO, type]:
    reader_class, writer_class = get_reader_writer_types(format, base_class)
    in_memory_stream_class = io.BytesIO if format == Format.BINARY else io.StringIO

    stream = in_memory_stream_class()
    with writer_class(stream) as w:
        w.write_an_int(42)
        w.write_a_stream([1, 2, 3, 4, 5])
        w.write_another_int(153)
    stream.seek(0)
    return stream, reader_class


def test_skip_steps(format: Format):
    stream, reader_class = create_testing_reader(format, tm.StateTestWriterBase)
    with reader_class(stream, skip_completed_check=True) as r:
        r.read_an_int()

    stream.seek(0)
    with reader_class(stream, skip_completed_check=True) as r:
        r.read_an_int()
        for _ in r.read_a_stream():
            pass
        r.close()


def test_skip_stream_items(format: Format):
    stream, reader_class = create_testing_reader(format, tm.StateTestWriterBase)
    with reader_class(stream, skip_completed_check=True) as r:
        r.read_an_int()
        stream = r.read_a_stream()
        for i, _ in enumerate(stream):
            if i == 2:
                break
