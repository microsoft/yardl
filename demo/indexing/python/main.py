import sketch
from io import BytesIO
import itertools

HEADER = sketch.Header(subject="Hello World!")


def generate_samples(N, start=0):
    for i in range(N):
        v = i + start
        yield sketch.Sample(id=v, data=[v, v + 1, v + 2])


def main():
    test_missing_index()
    test_empty_stream()
    test_copy_to()
    test_stream_read()
    print("Success!")


def test_missing_index():
    stream = BytesIO()
    with sketch.BinaryMyProtocolWriter(stream) as writer:
        writer.write_header(HEADER)
        writer.write_samples(generate_samples(10))

    stream.seek(0)
    try:
        with sketch.BinaryMyProtocolIndexedReader(stream) as reader:
            pass
    except RuntimeError as e:
        assert "binary index not found" in str(e).lower()


def test_copy_to():
    stream = BytesIO()
    with sketch.BinaryMyProtocolWriter(stream) as writer:
        writer.write_header(HEADER)
        writer.write_samples(generate_samples(10))

    stream.seek(0)
    output = BytesIO()
    with sketch.BinaryMyProtocolReader(stream) as reader:
        with sketch.BinaryMyProtocolIndexedWriter(output) as writer:
            reader.copy_to(writer)

    output.seek(0)
    with sketch.BinaryMyProtocolIndexedReader(output) as reader:
        header_read = reader.read_header()
        samples_read = list(reader.read_samples())

        assert header_read == HEADER
        assert len(samples_read) == 10


def test_empty_stream():
    stream = BytesIO()
    with sketch.BinaryMyProtocolIndexedWriter(stream) as writer:
        writer.write_header(HEADER)
        writer.write_samples(generate_samples(0))

    stream.seek(0)
    with sketch.BinaryMyProtocolIndexedReader(stream) as reader:
        header_read = reader.read_header()
        samples_read = list(reader.read_samples())

        assert header_read == HEADER
        assert len(samples_read) == 0


def test_stream_read():
    stream = BytesIO()
    with sketch.BinaryMyProtocolIndexedWriter(stream) as writer:
        writer.write_header(HEADER)
        writer.write_samples(list(generate_samples(77)))
        writer.write_samples(generate_samples(33, 77))
        writer.write_samples(list(generate_samples(55, 33 + 77)))

    total_samples = 77 + 33 + 55

    stream.seek(0)
    with sketch.BinaryMyProtocolIndexedReader(stream) as reader:
        samples_read = list(reader.read_samples())

        assert reader.read_header() == HEADER
        assert len(samples_read) == total_samples

    stream.seek(0)
    with sketch.BinaryMyProtocolIndexedReader(stream) as reader:
        for start in range(0, total_samples, 15):
            samples_read = list(itertools.islice(reader.read_samples(idx=start), 15))
            for i, s in enumerate(samples_read):
                assert s.id == i + start, f"Expected {i + start}, got {s.id}"

        assert reader.read_header() == HEADER


if __name__ == "__main__":
    main()
