import io
import test_model as tm


def generate_stream():
    b = io.BytesIO()
    with tm.BinaryStateTestWriter(b) as w:
        w.write_an_int(42)
        w.write_a_stream([1, 2, 3, 4, 5])
        w.write_another_int(153)
    b.seek(0)
    return b


def test_skip_steps():
    bytes = generate_stream()
    with tm.BinaryStateTestReader(bytes, skip_completed_check=True) as r:
        r.read_an_int()

    bytes.seek(0)
    with tm.BinaryStateTestReader(bytes, skip_completed_check=True) as r:
        r.read_an_int()
        for _ in r.read_a_stream():
            pass
        r.close()


def test_skip_stream_items():
    bytes = generate_stream()
    with tm.BinaryStateTestReader(bytes, skip_completed_check=True) as r:
        r.read_an_int()
        stream = r.read_a_stream()
        for i, _ in enumerate(stream):
            if i == 2:
                break
