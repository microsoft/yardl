import io
import subprocess
import test_model as tm
from test_model._binary import BinaryProtocolReader
import test_model.binary as tmb

def test_rt():
    path = "/workspaces/yardl/cpp/build/test_output/binary/RoundTripTests_Scalars_Binary.bin"
    with open(path, "rb") as f:
        expected = f.read()

    i = BinaryProtocolReader(path, None)

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
