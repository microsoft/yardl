import io
import inspect
import pytest
import subprocess
import test_model as tm
import pathlib
from test_model._binary import BinaryProtocolReader

cpp_test_output_dir = (pathlib.Path(__file__).parent / "../../cpp/build/test_output/binary/").resolve()

def cases():
    for path in cpp_test_output_dir.glob("RoundTripTests_*_Binary.bin"):
        yield path.name.removeprefix("RoundTripTests_").removesuffix("_Binary.bin")


def path_from_case_name(name):
    return str(cpp_test_output_dir / f"RoundTripTests_{name}_Binary.bin")


@pytest.fixture(scope="module")
def readers_writers_by_json():
    pairs = {}
    for name, obj in inspect.getmembers(tm.binary):
        if inspect.isclass(obj) and issubclass(obj, BinaryProtocolReader) and obj != BinaryProtocolReader:
            reader = obj
            writer =  getattr(tm.binary, name.removesuffix("Reader") + "Writer")
            schema = getattr(obj, "schema")
            pairs[schema] = (reader, writer)

    return pairs


@pytest.mark.parametrize("read_as_numpy", [tm.Types.ALL, tm.Types.NONE])
@pytest.mark.parametrize("case_name", cases())
def test_cpp_roundtrip(case_name, readers_writers_by_json, read_as_numpy):
    path = path_from_case_name(case_name)

    with open(path, "rb") as f:
        expected = f.read()

    with open(path, "rb") as f:
        i = BinaryProtocolReader(f, None)
        reader_type, writer_type = readers_writers_by_json[i._schema]

    x = io.BytesIO()
    with reader_type(path, read_as_numpy) as r, writer_type(x) as w:
        r.copy_to(w)


    if case_name == "SimpleDatasets":
        pytest.skip("we do not support writing streams with batches yet")

    if case_name == "Scalars" and read_as_numpy == tm.Types.NONE:
        # Times and datetimes only have microsecond precision in Python
        return

    if expected != bytes(x.getbuffer()):
        print("Expected:")
        subprocess.run(["hexdump", "-C", path])
        print("\nActual:")
        with subprocess.Popen(["hexdump", "-C"], stdin=subprocess.PIPE) as p:
            assert p.stdin != None
            p.stdin.write(x.getbuffer())
            p.stdin.close()

        assert expected == bytes(x.getbuffer())
