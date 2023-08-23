#! /usr/bin/env python3

from dataclasses import dataclass
import inspect
import json
import os
import pathlib
import subprocess
import sys
import time
from typing import Callable, NamedTuple, Optional

import numpy as np
from rich.live import Live
from rich.table import Table

import test_model as tm
from tests.factories import (
    Format,
    get_reader_type,
    get_writer_type,
)


OUTPUT_FILE = "/tmp/benchmark_data.dat"

# Setting this to True will display the roundtrip duration instead of the throughput
# and should only be done to calibrate the scales of each scenario and format.
DISPLAY_DURATIONS = False


@dataclass
class Result:
    write_mi_bytes_per_second: float
    read_mi_bytes_per_second: float
    roundtrip_duration_seconds: float = 0.0


class MutlitingualResults(NamedTuple):
    cpp: Optional[Result]
    python: Optional[Result]


_cpp_benchmark_path = (
    pathlib.Path(__file__).parent / "../../cpp/build/benchmark"
).resolve()


def scale_repetitions(repetitions: int, scale: float):
    return int(repetitions * scale)


def time_scenario(
    total_size_bytes: int, write_impl: Callable[[], None], read_impl: Callable[[], None]
) -> Result:
    try:
        os.remove(OUTPUT_FILE)
    except FileNotFoundError:
        pass

    total_size_mi_byte = total_size_bytes / 1024.0 / 1024.0

    write_start_time = time.monotonic_ns()
    write_impl()
    write_end_time = time.monotonic_ns()
    write_elapsed_seconds = (write_end_time - write_start_time) / 1e9
    write_mi_bytes_per_second = total_size_mi_byte / write_elapsed_seconds

    read_start_time = time.monotonic_ns()
    read_impl()
    read_end_time = time.monotonic_ns()
    read_elapsed_seconds = (read_end_time - read_start_time) / 1e9
    read_mi_bytes_per_second = total_size_mi_byte / read_elapsed_seconds

    return Result(
        write_mi_bytes_per_second=write_mi_bytes_per_second,
        read_mi_bytes_per_second=read_mi_bytes_per_second,
        roundtrip_duration_seconds=write_elapsed_seconds + read_elapsed_seconds,
    )


def benchmark_float_256x256(format: Format) -> Optional[Result]:
    if format == Format.HDF5:
        return None

    if format == Format.NDJSON:
        scale = 0.002
    else:
        scale = 1

    arr = np.random.random_sample((256, 256)).astype(np.float32)

    repetitions = scale_repetitions(10000, scale)
    total_size_bytes = arr.nbytes * repetitions

    def write():
        with get_writer_type(format, tm.BenchmarkFloat256x256WriterBase)(
            OUTPUT_FILE
        ) as w:
            w.write_float256x256(arr for _ in range(repetitions))

    def read():
        with get_reader_type(format, tm.BenchmarkFloat256x256ReaderBase)(
            OUTPUT_FILE
        ) as r:
            for _ in r.read_float256x256():
                pass

    return time_scenario(total_size_bytes, write, read)


def benchmark_float_vlen(format: Format) -> Optional[Result]:
    if format == Format.HDF5:
        return None

    if format == Format.NDJSON:
        scale = 0.002
    else:
        scale = 1

    arr = np.random.random_sample((256, 256)).astype(np.float32)

    repetitions = scale_repetitions(10000, scale)
    total_size_bytes = arr.nbytes * repetitions

    def write():
        with get_writer_type(format, tm.BenchmarkFloatVlenWriterBase)(OUTPUT_FILE) as w:
            w.write_float_array(arr for _ in range(repetitions))

    def read():
        with get_reader_type(format, tm.BenchmarkFloatVlenReaderBase)(OUTPUT_FILE) as r:
            for _ in r.read_float_array():
                pass

    return time_scenario(total_size_bytes, write, read)


def benchmark_small_int_256x256(format: Format) -> Optional[Result]:
    if format == Format.HDF5:
        return None

    if format == Format.NDJSON:
        scale = 0.03
    else:
        scale = 0.02

    arr = np.full((256, 256), 37, dtype=np.int32)

    repetitions = scale_repetitions(1000, scale)
    total_size_bytes = arr.nbytes * repetitions

    def write():
        with get_writer_type(format, tm.BenchmarkInt256x256WriterBase)(
            OUTPUT_FILE
        ) as w:
            w.write_int256x256(arr for _ in range(repetitions))

    def read():
        with get_reader_type(format, tm.BenchmarkInt256x256ReaderBase)(
            OUTPUT_FILE
        ) as r:
            for _ in r.read_int256x256():
                pass

    return time_scenario(total_size_bytes, write, read)


def benchmark_small_record(format: Format) -> Optional[Result]:
    if format == Format.HDF5:
        return None

    if format == Format.NDJSON:
        scale = 0.002
    else:
        scale = 0.005

    record = tm.SmallBenchmarkRecord(
        a=73278383.23123213, b=78323.2820379, c=-2938923.29882
    )

    repetitions = scale_repetitions(50000000, scale)
    total_size_bytes = 16 * repetitions

    def write():
        with get_writer_type(format, tm.BenchmarkSmallRecordWriterBase)(
            OUTPUT_FILE
        ) as w:
            w.write_small_record(record for _ in range(repetitions))

    def read():
        with get_reader_type(format, tm.BenchmarkSmallRecordReaderBase)(
            OUTPUT_FILE
        ) as r:
            for _ in r.read_small_record():
                pass

    return time_scenario(total_size_bytes, write, read)


def benchmark_small_record_batched(format: Format) -> Optional[Result]:
    # batching has not been implemented in the python version yet
    return None


def benchmark_small_optionals_batched(format: Format) -> Optional[Result]:
    # batching has not been implemented in the python version yet
    return None


def benchmark_simple_mrd(format: Format) -> Optional[Result]:
    if format == Format.HDF5:
        return None

    if format == Format.NDJSON:
        scale = 0.002
    else:
        scale = 0.5

    acq = tm.SimpleAcquisition()
    acq.data.resize((32, 256))
    acq.trajectory.resize((32, 2))

    data = tm.AcquisitionOrImage.Acquisition(acq)

    repetitions = scale_repetitions(30000, scale)
    total_size_bytes = 66032 * repetitions

    def write():
        with get_writer_type(format, tm.BenchmarkSimpleMrdWriterBase)(OUTPUT_FILE) as w:
            w.write_data(data for _ in range(repetitions))

    def read():
        with get_reader_type(format, tm.BenchmarkSimpleMrdReaderBase)(OUTPUT_FILE) as r:
            for _ in r.read_data():
                pass

    return time_scenario(total_size_bytes, write, read)


def invoke_cpp_benchmark(scenario: str, format: Format) -> Optional[Result]:
    start = time.monotonic_ns()

    res = subprocess.run(
        [_cpp_benchmark_path, scenario, str(format)],
        stdout=subprocess.PIPE,
        check=True,
        encoding="utf-8",
    )

    end = time.monotonic_ns()
    elapsed_seconds = (end - start) / 1e9

    if res.stdout == "":
        return None

    return Result(**json.loads(res.stdout), roundtrip_duration_seconds=elapsed_seconds)


def scenario_name(scenario_func: Callable[[Format], Optional[Result]]) -> str:
    return scenario_func.__name__.removeprefix("benchmark").replace("_", "")


def invoke_benchmark(
    scenario_func: Callable[[Format], Optional[Result]], format: Format
) -> MutlitingualResults:
    cpp_res = invoke_cpp_benchmark(scenario_name(scenario_func), format)
    python_res = scenario_func(format)
    return MutlitingualResults(cpp=cpp_res, python=python_res)


def update_table(
    table: Table,
    live: Live,
    scenario_func: Callable[[Format], Optional[Result]],
    format: Format,
    results: MutlitingualResults,
):
    def format_float(thoughput: float) -> str:
        return f"{thoughput:,.2f}"

    def color():
        if format == Format.HDF5:
            return "cyan"
        elif format == Format.BINARY:
            return "blue"
        elif format == Format.NDJSON:
            return "green"
        return None

    if DISPLAY_DURATIONS:
        table.add_row(
            scenario_name(scenario_func),
            str(format),
            format_float(results.cpp.roundtrip_duration_seconds)
            if results.cpp
            else None,
            format_float(results.python.roundtrip_duration_seconds)
            if results.python
            else None,
            style=color(),
        )
    else:
        table.add_row(
            scenario_name(scenario_func),
            str(format),
            format_float(results.cpp.write_mi_bytes_per_second)
            if results.cpp
            else None,
            format_float(results.cpp.read_mi_bytes_per_second) if results.cpp else None,
            format_float(results.python.write_mi_bytes_per_second)
            if results.python
            else None,
            format_float(results.python.read_mi_bytes_per_second)
            if results.python
            else None,
            style=color(),
        )
    live.refresh()


if __name__ == "__main__":
    table = Table()
    if DISPLAY_DURATIONS:
        table.title = "Roundtrip duration"
        table.add_column("Scenario")
        table.add_column("Format")
        table.add_column("C++ Duration", justify="right")
        table.add_column("Python Duration", justify="right")
    else:
        table.title = "Throughput in MiB/s"
        table.add_column("Scenario")
        table.add_column("Format")
        table.add_column("C++ Write", justify="right")
        table.add_column("C++ Read", justify="right")
        table.add_column("Python Write", justify="right")
        table.add_column("Python Read", justify="right")

    with Live(table, auto_refresh=False) as live:
        for _, benchmark_func in inspect.getmembers(
            sys.modules[__name__],
            lambda x: inspect.isfunction(x) and x.__name__.startswith("benchmark_"),
        ):
            table.add_section()
            for format in [Format.HDF5, Format.BINARY, Format.NDJSON]:
                res = invoke_benchmark(benchmark_func, format)
                update_table(table, live, benchmark_func, format, res)
