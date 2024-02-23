# Copyright (c) Microsoft Corporation.
# Licensed under the MIT License.

import sys
import subprocess

V0 = "v0"
V1 = "v1"
V2 = "v2"


def main():
    # Test for Regressions
    write_copy_validate(V0, V0, V0)
    write_copy_validate(V1, V1, V1)
    write_copy_validate(V2, V2, V2)

    # Test Backward Compatibility
    write_copy_validate(V0, V0, V1)
    write_copy_validate(V0, V0, V2)
    write_copy_validate(V1, V1, V2)

    # Test Round-Trip Compatibility
    write_copy_validate(V0, V1, V0)
    write_copy_validate(V0, V2, V0)
    write_copy_validate(V1, V2, V1)

    print("Evolution tests passed.", flush=True)


def write_copy_validate(w: str, c: str, v: str) -> None:
    try:
        pipetest(w, c, v)
    except subprocess.CalledProcessError as e:
        log("Error occurred in write {} -> copy {} -> validate {}".format(w, c, v))
        log("Command: {}".format(e.cmd))
        if e.stdout:
            log(e.stdout.decode("utf-8"))
        if e.stderr:
            log(e.stderr.decode("utf-8"))
        sys.exit(1)


def log(msg: str) -> None:
    print("Evolution Test: {}".format(msg), flush=True)


def pipetest(w: str, c: str, v: str) -> None:
    write = subprocess.run(write_command(w), capture_output=True, check=True)
    copy = subprocess.run(
        copy_command(c),
        capture_output=True,
        check=True,
        input=write.stdout,
    )
    validate = subprocess.run(
        validate_command(v),
        capture_output=True,
        check=True,
        input=copy.stdout,
    )
    validate.check_returncode()


def write_command(v: str) -> list[str]:
    return ["./{}_write".format(v)]


def copy_command(v: str) -> list[str]:
    return ["./{}_copy".format(v)]


def validate_command(v: str) -> list[str]:
    return ["./{}_validate".format(v)]


if __name__ == "__main__":
    main()
