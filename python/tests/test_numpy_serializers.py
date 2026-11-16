import io

import numpy as np
from numpy.lib.recfunctions import repack_fields
import pytest

import test_model as tm
from test_model import _binary, _ndjson
from test_model.binary import SimpleRecordSerializer, RecordWithOptionalFieldsSerializer
from test_model.ndjson import RecordWithOptionalFieldsConverter


@pytest.mark.parametrize("as_array", [False, True])
@pytest.mark.parametrize("empty", [False, True])
def test_numpy_vector_can_be_read_and_written_again(as_array: bool, empty: bool):
    expected = [] if empty else [tm.SimpleRecord(x=1, y=2, z=3)]
    value = expected
    if as_array:
        value = np.array([] if empty else [(1, 2, 3)], dtype=tm.get_dtype(tm.SimpleRecord))

    serializer = _binary.VectorSerializer(SimpleRecordSerializer())
    first = io.BytesIO()
    stream = _binary.CodedOutputStream(first)
    serializer.write_numpy(stream, value)
    stream.close()

    decoded = serializer.read_numpy(_binary.CodedInputStream(io.BytesIO(first.getvalue())))
    assert isinstance(decoded, list)
    assert decoded == expected

    second = io.BytesIO()
    stream = _binary.CodedOutputStream(second)
    serializer.write_numpy(stream, decoded)
    stream.close()
    assert second.getvalue() == first.getvalue()


@pytest.mark.parametrize(
    "value", [np.array([[1, 2]], dtype=np.int32), np.array(1), (1, 2), 1]
)
def test_numpy_vector_rejects_non_vector_values(value):
    serializer = _binary.VectorSerializer(_binary.int32_serializer)
    output = io.BytesIO()
    stream = _binary.CodedOutputStream(output)
    with pytest.raises(ValueError):
        serializer.write_numpy(stream, value)
    stream.close()
    assert output.getvalue() == b""


@pytest.mark.parametrize("format", ["binary", "ndjson"])
def test_numpy_optional_time_preserves_nanoseconds(format: str):
    dtype = repack_fields(tm.get_dtype(tm.RecordWithOptionalFields), align=False, recurse=True)
    assert dtype["optional_time"]["value"] == np.dtype("timedelta64[ns]")
    expected = np.array(
        [
            ((True, 5), (False, 0), (True, 3_723_000_000_004)),
            ((False, 0), (True, 9), (False, 0)),
        ],
        dtype=dtype,
    )
    if format == "binary":
        serializer = _binary.DynamicNDArraySerializer(RecordWithOptionalFieldsSerializer())
        output = io.BytesIO()
        stream = _binary.CodedOutputStream(output)
        serializer.write(stream, expected)
        stream.close()
        actual = serializer.read(_binary.CodedInputStream(io.BytesIO(output.getvalue())))
    else:
        converter = _ndjson.DynamicNDArrayConverter(RecordWithOptionalFieldsConverter())
        actual = converter.from_json(converter.to_json(expected))

    assert actual.dtype["optional_time"]["value"] == np.dtype("timedelta64[ns]")
    assert actual[0]["optional_time"]["has_value"]
    assert actual[0]["optional_time"]["value"] == np.timedelta64(3_723_000_000_004, "ns")
    assert not actual[1]["optional_time"]["has_value"]
    assert actual[0]["optional_int"]["value"] == 5
    assert actual[1]["optional_int_alternate_syntax"]["value"] == 9
