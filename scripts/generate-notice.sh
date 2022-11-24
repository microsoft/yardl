#! /bin/bash

# Copyright (c) Microsoft Corporation.
# Licensed under the MIT License.

# This script generates the NOTICE file for the repository based on the GO dependencies
# using the go-licenses tool.

set -euo pipefail

go install github.com/google/go-licenses@v1.5.0

repo_root=$(readlink -f "$(dirname "$0")/..")

cd "$repo_root/tooling"

save_dir="/tmp/licenses"

# The tools will write out warnings to stderr for non-go binary dependencies that it cannot follow.
# This pattern will filter the known warnings out of the output.
known_non_go_dependency_patterns="(golang.org/x/sys.*/unix)"

go-licenses save ./... --ignore "github.com/microsoft/yardl/tooling" --save_path=$save_dir --force 2> >(grep -Pv "$known_non_go_dependency_patterns")

# license and notice files will be in directories named after the import path of each library

# get the library names from the directory names
lib_names=$(find $save_dir -type f -print0 | xargs -0 realpath --relative-to $save_dir | xargs dirname | sort | uniq)

# the file we will be writing to
notice_path="$repo_root/NOTICE.txt"

# the header of the notice file
cat > "$notice_path" <<- EOM
NOTICES

This repository incorporates material as listed below or described in the code.

EOM

for lib_name in $lib_names; do
    {
        echo "================================================================================"
        echo -e "\n$lib_name\n"

        notice_pattern="NOTICE*"

        license=$(find "$save_dir/$lib_name" -type f ! -iname "$notice_pattern" -exec cat {} \;)

        if [ -n "$license" ]; then
            echo "$license"
            echo ""
        fi

        notice=$(find "$save_dir/$lib_name" -type f -iname "$notice_pattern" -exec cat {} \;)
        if [ -n "$notice" ]; then
            echo "$notice"
            echo ""
        fi
    } >> "$notice_path"

done
