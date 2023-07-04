import glob
import io
import inspect
import pytest
import subprocess
import json
import test_model as tm
from test_model._binary import BinaryProtocolReader
import test_model.binary as tmb

def test_rt():
    path = "/workspaces/yardl/cpp/build/test_output/binary/RoundTripTests_Scalars_Binary.bin"
    with open(path, "rb") as f:
        expected = f.read()

    x = io.BytesIO()
    with tmb.BinaryScalarsReader(path, tm.Types.ALL) as r, tmb.BinaryScalarsWriter(x) as w:
        r.copy_to(w)

    if expected != bytes(x.getbuffer()):
        print("Expected:")
        subprocess.run(["hexdump", "-C", path])
        print("\nActual:")
        with subprocess.Popen(["hexdump", "-C"], stdin=subprocess.PIPE) as p:
            assert p.stdin != None
            p.stdin.write(x.getbuffer())
            p.stdin.close()

        assert expected == bytes(x.getbuffer())


def files():
    # get files in this directory
    return glob.glob("/workspaces/yardl/cpp/build/test_output/binary/*.bin")


@pytest.fixture(scope="module")
def readers_writers_by_json():
    pairs = {}
    for name, obj in inspect.getmembers(tm.binary):
        if inspect.isclass(obj) and issubclass(obj, BinaryProtocolReader) and obj != BinaryProtocolReader:
            reader = obj
            writer =  getattr(tm.binary, name.removesuffix("Reader") + "Writer")
            schema = normalize_schema(getattr(obj, "schema"))
            pairs[schema] = (reader, writer)

    return pairs


def normalize_schema(schema):
    return json.dumps(json.loads(schema))

@pytest.mark.parametrize("file", files())
def test_cpp_roundtrip(file, readers_writers_by_json):
    with open(file, "rb") as f:
        expected = f.read()

    with open(file, "rb") as f:
        i = BinaryProtocolReader(f, None)
        schema = normalize_schema(i._schema)
        reader_type, writer_type = readers_writers_by_json[schema]

    x = io.BytesIO()
    with reader_type(file, tm.Types.ALL) as r, writer_type(x) as w:
        r.copy_to(w)


    if "RoundTripTests_SimpleDatasets_Binary.bin" in file:
        pytest.skip("Do not yet support writing streams with batches")

    if expected != bytes(x.getbuffer()):
        print("Expected:")
        subprocess.run(["hexdump", "-C", file])
        print("\nActual:")
        with subprocess.Popen(["hexdump", "-C"], stdin=subprocess.PIPE) as p:
            assert p.stdin != None
            p.stdin.write(x.getbuffer())
            p.stdin.close()

        assert expected == bytes(x.getbuffer())
