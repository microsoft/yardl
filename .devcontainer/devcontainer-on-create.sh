#! /usr/bin/env bash

set -euo pipefail

host_license_file="/tmp/matlab.host.lic"

if [ -s "${host_license_file}" ]; then
    # We don't bind mount directly to $MATLAB_LICENSE_FILE so that we can write
    # to it later without affecting the host.
    cp "${host_license_file}" "${MATLAB_LICENSE_FILE}"
fi
