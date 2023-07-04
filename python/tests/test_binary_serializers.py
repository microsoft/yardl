import io
import test_model as tm
import test_model._binary as tmb
import pytest


# @pytest.mark.parametrize("value", range(200))
# def test_unsigned_varint(value: int):
#     buf = io.BytesIO()
#     w = tmb.CodedOutputStream(buf)
#     w.write_unsigned_varint(value)
#     w.flush()

#     buf.seek(0)
#     r = tmb.CodedInputStream(buf)
#     assert r.read_unsigned_varint() == value


# @pytest.mark.parametrize("value", range(-300, 300))
# def test_signed_varint(value: int):
#     buf = io.BytesIO()
#     w = tmb.CodedOutputStream(buf)
#     w.write_signed_varint(value)
#     w.flush()

#     buf.seek(0)
#     r = tmb.CodedInputStream(buf)
#     assert r.read_signed_varint() == value
