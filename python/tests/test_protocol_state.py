from typing import Iterable

import pytest
import test_model as tm


class _TestStateTestWriter(tm.StateTestWriterBase):
    def _write_an_int(self, value: tm.Int32) -> None:
        pass

    def _write_a_stream(self, value: Iterable[tm.Int32]) -> None:
        pass

    def _write_another_int(self, value: tm.Int32) -> None:
        pass

    def _end_stream(self) -> None:
        pass

    def close(self) -> None:
        return super().close()


def test_proper_sequence_write():
    with _TestStateTestWriter() as w:
        w.write_an_int(1)
        w.write_a_stream([1, 2, 3])
        w.write_another_int(3)


def test_proper_sequence_write_empty_stream():
    with _TestStateTestWriter() as w:
        w.write_an_int(1)
        w.write_a_stream([])
        w.write_another_int(3)


def test_proper_sequence_write_with_multiple_stream_writes():
    with _TestStateTestWriter() as w:
        w.write_an_int(1)
        w.write_a_stream([1, 2, 3])
        w.write_a_stream([4, 5, 6])
        w.write_another_int(3)


def test_sequence_write_missing_first_step():
    with pytest.raises(
        tm.ProtocolError,
        match="Expected to call to 'write_an_int' but received call to 'write_a_stream'.",
    ), _TestStateTestWriter() as w:
        w.write_a_stream([1, 2, 3])


def test_sequence_write_premature_close():
    with pytest.raises(
        tm.ProtocolError,
        match="Protocol writer closed before all steps were called. Expected to call to 'write_another_int'.",
    ), _TestStateTestWriter() as w:
        w.write_an_int(1)
        w.write_a_stream([1, 2, 3])


class _TestStateTestReader(tm.StateTestReaderBase):
    def _read_an_int(self) -> tm.Int32:
        return -2

    def _read_a_stream(self) -> Iterable[tm.Int32]:
        yield -1
        yield -2
        yield -3

    def _read_another_int(self) -> tm.Int32:
        return -4

    def close(self) -> None:
        pass


def test_proper_sequence_read():
    with _TestStateTestReader() as r:
        r.read_an_int()
        for _ in r.read_a_stream():
            pass
        r.read_another_int()


def test_read_without_consuming_stream():
    with pytest.raises(
        tm.ProtocolError,
        match="Received call to 'read_another_int' but the iterable returned by 'read_a_stream' was not fully consumed",
    ), _TestStateTestReader() as r:
        r.read_an_int()
        r.read_a_stream()
        r.read_another_int()


def test_read_without_consuming_stream_and_closing():
    with pytest.raises(
        tm.ProtocolError,
        match="Protocol reader closed before all data was consumed. The iterable returned by 'read_a_stream' was not fully consumed.",
    ), _TestStateTestReader() as r:
        r.read_an_int()
        r.read_a_stream()


def test_read_without_first_step():
    with pytest.raises(
        tm.ProtocolError,
        match="Expected to call to 'read_an_int' but received call to 'read_a_stream'.",
    ), _TestStateTestReader() as r:
        r.read_a_stream()
        r.read_another_int()

    with pytest.raises(
        tm.ProtocolError,
        match="Protocol reader closed before all data was consumed. Expected call to 'read_another_int'.",
    ), _TestStateTestReader() as r:
        r.read_an_int()
        for _ in r.read_a_stream():
            pass
