# Copyright (c) Microsoft Corporation.
# Licensed under the MIT License.

"""Tests for the binary CodedOutputStream / CodedInputStream and the
serializers built directly on top of them.

The protocol-roundtrip tests exercise binary serialization end-to-end but only
with the default 65536-byte buffer and inputs that comfortably fit in it. These
tests deliberately use very small buffer_size values to drive flushes at
controlled offsets, mirroring the strategy used by the C++
``cpp/test/binary/coded_stream_test.cc`` suite. The buffer-boundary cases that
exposed bug #289 had no unit-level coverage at all before this file.
"""

import io

import pytest

from test_model._binary import (
    CodedInputStream,
    CodedOutputStream,
    FixedVectorSerializer,
    OptionalSerializer,
    StreamSerializer,
    UnionSerializer,
    bool_serializer,
    int8_serializer,
    int16_serializer,
    int32_serializer,
    int64_serializer,
    size_serializer,
    string_serializer,
    uint8_serializer,
    uint16_serializer,
    uint32_serializer,
    uint64_serializer,
)
from test_model.types import Int32OrSimpleRecord


# ---------------------------------------------------------------------------
# CodedOutputStream — direct byte-level behavior
# ---------------------------------------------------------------------------

class TestCodedOutputStreamBuffering:
    def test_write_byte_when_buffer_exactly_full_does_not_crash(self):
        """Regression for #289. Writing a byte when the buffer is exactly full
        must flush and then write, rather than raising IndexError."""
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=4)
        w.write_bytes(b"abcd")  # fills the buffer exactly; no flush yet
        assert out.getvalue() == b""
        w.write_byte(1)  # pre-fix this raised IndexError
        w.flush()
        assert out.getvalue() == b"abcd\x01"

    def test_write_byte_with_room_left_keeps_buffer(self):
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=8)
        w.write_byte(1)
        w.write_byte(2)
        assert out.getvalue() == b""  # still buffered
        w.flush()
        assert out.getvalue() == b"\x01\x02"

    def test_write_byte_rejects_out_of_range_values(self):
        w = CodedOutputStream(io.BytesIO(), buffer_size=4)
        with pytest.raises(AssertionError):
            w.write_byte(-1)
        with pytest.raises(AssertionError):
            w.write_byte(256)

    def test_write_bytes_at_exact_boundary_flushes_then_buffers(self):
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=4)
        w.write_bytes(b"abcd")  # buffer exactly full, no flush yet
        assert out.getvalue() == b""
        w.write_bytes(b"ef")  # remaining=0, must flush before buffering
        w.flush()
        assert out.getvalue() == b"abcdef"

    def test_write_bytes_larger_than_buffer_bypasses_buffer(self):
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=4)
        w.write_byte(0xFF)
        w.write_bytes(b"0123456789")  # > buffer_size; flushes existing then writes directly
        w.flush()
        assert out.getvalue() == b"\xff0123456789"

    def test_close_flushes_pending_buffer(self):
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=4)
        w.write_byte(5)
        w.close()
        assert out.getvalue() == b"\x05"

    def test_ensure_capacity_flushes_only_when_needed(self):
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=4)
        w.write_byte(1)
        # remaining == 3, ensure_capacity(3) must NOT flush
        w.ensure_capacity(3)
        assert out.getvalue() == b""
        w.write_bytes(b"abc")  # fills buffer exactly
        # remaining == 0, ensure_capacity(1) must flush
        w.ensure_capacity(1)
        assert out.getvalue() == b"\x01abc"


# ---------------------------------------------------------------------------
# Buffer-boundary roundtrip parity with C++ coded_stream_test.cc
# ---------------------------------------------------------------------------

class TestBufferBoundaryRoundtrip:
    def test_read_past_end_raises_eof(self):
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=10)
        uint32_serializer.write(w, 1)
        w.close()
        r = CodedInputStream(io.BytesIO(out.getvalue()), buffer_size=10)
        assert uint32_serializer.read(r) == 1
        with pytest.raises(EOFError):
            uint32_serializer.read(r)

    @pytest.mark.parametrize(
        "writer_buffer_size,reader_buffer_size",
        # buffer_size must be >= MAX_VARINT64_BYTES (10) because
        # write_unsigned_varint does a single ensure_capacity(10) up front and
        # then writes up to 10 bytes unchecked. Same constraint applies in C++.
        [(64, 64), (64, 65), (64, 63), (10, 1024), (1024, 10)],
    )
    def test_string_roundtrip_across_buffer_sizes(
        self, writer_buffer_size, reader_buffer_size
    ):
        s = "a" * 256
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=writer_buffer_size)
        string_serializer.write(w, s)
        w.close()
        r = CodedInputStream(io.BytesIO(out.getvalue()), buffer_size=reader_buffer_size)
        assert string_serializer.read(r) == s

    def test_long_string_roundtrip(self):
        """Analogue of the C++ ``Strings`` test: a string much larger than the
        buffer surrounded by short strings."""
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=10)
        string_serializer.write(w, "hello")
        string_serializer.write(w, "a" * 20000)
        string_serializer.write(w, "world")
        w.close()

        r = CodedInputStream(io.BytesIO(out.getvalue()), buffer_size=10)
        assert string_serializer.read(r) == "hello"
        assert string_serializer.read(r) == "a" * 20000
        assert string_serializer.read(r) == "world"


# ---------------------------------------------------------------------------
# Varint encoding — parity with C++ VarShort / VarUShort / VarIntegers
# ---------------------------------------------------------------------------

# Same value set as the C++ tests.
_VARINT_SAMPLES = [
    0, 1, 5, 33, 0x7E, 0x7F, 0x80, 0x81, 255, 256, 257,
    838, 0x3FFF, 0x4000, 0x4001, 0x7FFF, 0x8000, 0x8001, 0xFFFF,
    283928, 2847772, 0x7FFFFFFF, 0xFFFFFFFF,
]


class TestVarInt:
    def test_unsigned_varint_roundtrip_with_small_buffer(self):
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=10)
        for v in _VARINT_SAMPLES:
            w.write_unsigned_varint(v)
        w.close()
        r = CodedInputStream(io.BytesIO(out.getvalue()), buffer_size=10)
        for v in _VARINT_SAMPLES:
            assert r.read_unsigned_varint() == v

    def test_signed_varint_zigzag_roundtrip(self):
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=10)
        for v in _VARINT_SAMPLES:
            w.write_signed_varint(int(v))
            w.write_signed_varint(-int(v))
        w.close()
        r = CodedInputStream(io.BytesIO(out.getvalue()), buffer_size=10)
        for v in _VARINT_SAMPLES:
            assert r.read_signed_varint() == int(v)
            assert r.read_signed_varint() == -int(v)

    @pytest.mark.parametrize(
        "serializer,values",
        [
            (int8_serializer, [-128, -1, 0, 1, 127]),
            (uint8_serializer, [0, 1, 127, 255]),
            (int16_serializer, [-32768, -1, 0, 1, 32767]),
            (uint16_serializer, [0, 1, 0xFFFF]),
            (int32_serializer, [-(2 ** 31), -1, 0, 1, 2 ** 31 - 1]),
            (uint32_serializer, [0, 1, 0xFFFFFFFF]),
            (int64_serializer, [-(2 ** 63), -1, 0, 1, 2 ** 63 - 1]),
            (uint64_serializer, [0, 1, 2 ** 64 - 1]),
            (size_serializer, [0, 1, 0xFFFF, 2 ** 32, 2 ** 63 - 1]),
            (bool_serializer, [False, True, False]),
        ],
        ids=lambda x: type(x).__name__ if not isinstance(x, list) else "values",
    )
    def test_each_integer_serializer_roundtrip(self, serializer, values):
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=10)
        for v in values:
            serializer.write(w, v)
        w.close()
        r = CodedInputStream(io.BytesIO(out.getvalue()), buffer_size=10)
        for v in values:
            assert serializer.read(r) == v


# ---------------------------------------------------------------------------
# StreamSerializer — regression tests for the per-item path of bug #289
# ---------------------------------------------------------------------------

class TestStreamSerializerRegression:
    def test_per_item_marker_at_exact_buffer_boundary(self):
        """Regression for #289.

        With a 4-byte buffer and a 3-byte element body, each per-item entry is
        marker(1) + body(3) = 4 bytes. After the first element body is written
        the buffer offset lands on exactly buffer_size. Pre-fix, the next
        per-item marker (``write_byte_no_check(1)``) raised IndexError; post-fix
        the marker triggers a flush first.
        """
        element = FixedVectorSerializer(uint8_serializer, 3)
        stream_serializer = StreamSerializer(element)

        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=4)
        items = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
        # iter() bypasses the batch-list fast path inside StreamSerializer.write
        stream_serializer.write(w, iter(items))
        w.write_byte(0)  # block terminator
        w.close()

        r = CodedInputStream(io.BytesIO(out.getvalue()), buffer_size=4)
        assert list(stream_serializer.read(r)) == items

    def test_per_item_marker_with_optional_element(self):
        """Closer to the original #289 report: a stream of small optionals where
        the marker write can land on any offset relative to the buffer."""
        element = OptionalSerializer(int8_serializer)
        stream_serializer = StreamSerializer(element)
        items = [None, 7, None, None, 8, None, 9, 10]

        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=2)
        stream_serializer.write(w, iter(items))
        w.write_byte(0)
        w.close()

        r = CodedInputStream(io.BytesIO(out.getvalue()), buffer_size=2)
        assert list(stream_serializer.read(r)) == items

    def test_large_iterator_at_default_buffer_size(self):
        """End-to-end repro of bug #289 at the default 65536-byte buffer size.

        ``OptionalSerializer(int8_serializer)`` writes 1 byte per None and 2
        bytes per Some(_). Combined with the 1-byte per-item stream marker the
        all-None stream emits 2 bytes per item, so the marker write lands on
        the 65536 boundary at item 32768. Pre-fix that raised IndexError; the
        target here is simply to complete the loop without crashing on a
        realistic buffer size.
        """
        element = OptionalSerializer(int8_serializer)
        stream_serializer = StreamSerializer(element)
        n = 70000
        items = [None] * n

        out = io.BytesIO()
        w = CodedOutputStream(out)  # default buffer_size=65536
        stream_serializer.write(w, iter(items))
        w.write_byte(0)
        w.close()

        r = CodedInputStream(io.BytesIO(out.getvalue()))
        decoded = list(stream_serializer.read(r))
        assert len(decoded) == n
        assert all(v is None for v in decoded)

    def test_batch_and_per_item_paths_decode_to_same_values(self):
        """The batch-list path and the per-item path encode differently (single
        varint count + bodies vs. one marker per body) but must decode to the
        same values."""
        element = OptionalSerializer(int32_serializer)
        stream_serializer = StreamSerializer(element)
        items = [None, 1, None, 2, 3, None]

        out_batch = io.BytesIO()
        w_batch = CodedOutputStream(out_batch, buffer_size=4)
        stream_serializer.write(w_batch, items)  # list -> batch fast path
        w_batch.write_byte(0)
        w_batch.close()

        out_iter = io.BytesIO()
        w_iter = CodedOutputStream(out_iter, buffer_size=4)
        stream_serializer.write(w_iter, iter(items))  # iterator -> per-item path
        w_iter.write_byte(0)
        w_iter.close()

        # The byte streams differ in framing, so the buffers won't match.
        assert out_batch.getvalue() != out_iter.getvalue()

        r_batch = CodedInputStream(io.BytesIO(out_batch.getvalue()), buffer_size=4)
        r_iter = CodedInputStream(io.BytesIO(out_iter.getvalue()), buffer_size=4)
        assert list(stream_serializer.read(r_batch)) == items
        assert list(stream_serializer.read(r_iter)) == items


# ---------------------------------------------------------------------------
# OptionalSerializer
# ---------------------------------------------------------------------------

class TestOptionalSerializer:
    def test_roundtrip_across_buffer_boundary(self):
        serializer = OptionalSerializer(uint8_serializer)
        items = [None, 1, None, 2, 3, None, None, 4]
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=2)
        for v in items:
            serializer.write(w, v)
        w.close()
        r = CodedInputStream(io.BytesIO(out.getvalue()), buffer_size=2)
        for v in items:
            assert serializer.read(r) == v

    def test_optional_of_string_with_small_buffer(self):
        serializer = OptionalSerializer(string_serializer)
        items = [None, "hello", None, "this is a longer payload", ""]
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=3)
        for v in items:
            serializer.write(w, v)
        w.close()
        r = CodedInputStream(io.BytesIO(out.getvalue()), buffer_size=3)
        for v in items:
            assert serializer.read(r) == v


# ---------------------------------------------------------------------------
# UnionSerializer — regression for the latent None-branch bug
# ---------------------------------------------------------------------------

class TestUnionSerializerRegression:
    def _nullable_int_or_record_serializer(self):
        return UnionSerializer(
            Int32OrSimpleRecord,
            [None, (Int32OrSimpleRecord.Int32, int32_serializer)],
        )

    def test_none_branch_when_buffer_exactly_full_does_not_crash(self):
        """Regression for the same class of bug as #289 but on the
        ``UnionSerializer.write`` None branch. Pre-fix that branch wrote a 0
        tag byte with no capacity check; with the buffer exactly full it
        raised IndexError. Post-fix it goes through ``write_byte`` which
        flushes first."""
        serializer = self._nullable_int_or_record_serializer()
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=4)
        w.write_bytes(b"abcd")  # buffer exactly full
        serializer.write(w, None)  # pre-fix: IndexError
        w.close()
        # Tag for the None case is 0 (cases[0] is None).
        assert out.getvalue() == b"abcd\x00"

    def test_tag_branch_when_buffer_exactly_full_does_not_crash(self):
        """The non-None tag branch already had ``ensure_capacity(1)`` before
        the fix, but the refactor routes it through ``write_byte``. Verify the
        boundary behavior is unchanged."""
        serializer = self._nullable_int_or_record_serializer()
        out = io.BytesIO()
        w = CodedOutputStream(out, buffer_size=4)
        w.write_bytes(b"abcd")  # buffer exactly full
        serializer.write(w, Int32OrSimpleRecord.Int32(42))
        w.close()
        # Tag index for Int32 with a leading None case is 1.
        expected = b"abcd" + bytes([1]) + bytes([84])  # 42 zigzag-encoded = 84
        assert out.getvalue() == expected

    def test_union_roundtrip_with_small_buffer(self):
        serializer = self._nullable_int_or_record_serializer()
        values = [None, Int32OrSimpleRecord.Int32(0), None,
                  Int32OrSimpleRecord.Int32(-1), Int32OrSimpleRecord.Int32(123456)]
        out = io.BytesIO()
        # buffer_size must be >= MAX_VARINT64_BYTES (10) since the Int32 case
        # writes a varint of up to 5 bytes after a single ensure_capacity(10).
        w = CodedOutputStream(out, buffer_size=10)
        for v in values:
            serializer.write(w, v)
        w.close()
        r = CodedInputStream(io.BytesIO(out.getvalue()), buffer_size=10)
        for v in values:
            assert serializer.read(r) == v
