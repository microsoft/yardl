import importlib.util
import io
import sys
import types
from pathlib import Path


def load_binary_runtime():
    static_dir = (
        Path(__file__).resolve().parents[2]
        / "tooling"
        / "internal"
        / "python"
        / "static_files"
    )
    package_name = "_yardl_static_files_for_test"

    if package_name not in sys.modules:
        package = types.ModuleType(package_name)
        package.__path__ = [str(static_dir)]
        sys.modules[package_name] = package

    spec = importlib.util.spec_from_file_location(
        f"{package_name}._binary", static_dir / "_binary.py"
    )
    assert spec is not None
    assert spec.loader is not None

    module = importlib.util.module_from_spec(spec)
    sys.modules[spec.name] = module
    spec.loader.exec_module(module)
    return module


def test_stream_serializer_flushes_before_item_marker_at_buffer_boundary():
    binary = load_binary_runtime()

    output = io.BytesIO()
    stream = binary.CodedOutputStream(output, buffer_size=2)
    serializer = binary.StreamSerializer(
        binary.OptionalSerializer(binary.uint32_serializer)
    )

    serializer.write(stream, (None for _ in range(2)))
    stream.close()

    assert output.getvalue() == b"\x01\x00\x01\x00"
