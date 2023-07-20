import datetime

import numpy as np
import pytest

import test_model.yardl_types as yt


def test_datetime_from_valid_datetime():
    dt = yt.DateTime.from_datetime(datetime.datetime(2020, 2, 29, 12, 22, 44, 111222))
    assert str(dt) == "2020-02-29T12:22:44.111222000"


def test_datetime_from_invalid_datetime():
    with pytest.raises(ValueError, match="day is out of range for month"):
        yt.DateTime.from_datetime(datetime.datetime(2021, 2, 29, 12, 22, 44, 111222))


def test_datetime_from_int():
    dt = yt.DateTime(1577967764111222333)
    assert str(dt) == "2020-01-02T12:22:44.111222333"


def test_datetime_from_datetime64():
    dt = yt.DateTime(np.datetime64(1577967764111222, "us"))
    assert str(dt) == "2020-01-02T12:22:44.111222000"


def test_now():
    dt = yt.DateTime.now()
    assert isinstance(dt, yt.DateTime)


def test_time_from_time():
    t = yt.Time.from_time(datetime.time(12, 22, 44, 111222))
    assert str(t) == "12:22:44.111222"


def test_parse():
    assert str(yt.Time.parse("12:22:44.111222")) == "12:22:44.111222"
    assert str(yt.Time.parse("12:22:44.111222333")) == "12:22:44.111222333"
    assert str(yt.Time.parse("12:22:44.1")) == "12:22:44.1"
    assert str(yt.Time.parse("12:22:44")) == "12:22:44"
    assert str(yt.Time.parse("12:22")) == "12:22:00"
    assert str(yt.Time.parse("12:22:44.000000001")) == "12:22:44.000000001"
    assert str(yt.Time.parse("12:22:44.0000000001")) == "12:22:44"
    assert str(yt.Time.parse("12:22:44.0000000000000000001")) == "12:22:44"


def test_parse_invalid():
    with pytest.raises(ValueError):
        yt.Time.parse("a")
    with pytest.raises(ValueError):
        yt.Time.parse("25:22:44.111222333444")
