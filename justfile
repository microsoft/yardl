set shell := ['bash', '-ceuo', 'pipefail']

cpp_version := "17"

@default: validate
  echo $'\n\e[1;34mNote: you can run \'just test\' to a run an inner-loop subset of the complete validation'

@ensure-build-dir:
    mkdir -p cpp/build

@configure: ensure-build-dir
    cd cpp/build; \
    cmake -GNinja -D CMAKE_CXX_STANDARD={{ cpp_version }} ..

@install:
    cd tooling/cmd/yardl; \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-s -w" .

@generate: install
    cd models/test && yardl generate

@generate-sandbox: install
    cd models/sandbox && yardl generate

@build-sandbox: generate-sandbox ensure-build-dir
    cd cpp/build && ninja sandbox_exec

@run-sandbox: build-sandbox
    cd cpp/build && ./sandbox_exec

@run-sandbox-quiet: build-sandbox
    cd cpp/build && ./sandbox_exec > /dev/null

@build-all: generate generate-sandbox configure
    cd cpp/build && ninja

@tooling-test:
    cd tooling; \
    go test ./... | { grep -v "\\[[no test files\\]" || true; }

@watch-tooling-test:
    cd tooling; \
    watchexec -r -c -w . -- 'go test ./... | { grep -v "\\[[no test files\\]" || true; }'

@python-test: generate
    cd python; \
    python3 -m pytest tests

@cpp-test: tooling-test generate ensure-build-dir
    cd cpp/build; \
    ninja tests; \
    ./tests --gtest_brief=1

@test: tooling-test cpp-test python-test

@benchmark: generate ensure-build-dir
    cd cpp/build; \
    ninja benchmark; \
    ./benchmark

@watch-generate-test: install
    watchexec -r -c -w tooling/ -- "just install && cd models/test && yardl generate --watch"

@watch-generate-sandbox:
    watchexec -r -c -w tooling/ -- "just install && cd models/sandbox && yardl generate --watch"

@watch-exec-sandbox: configure
    watchexec -c -w models/sandbox/ -w cpp/ -i **/cpp/build/** -i **/cpp/test/** -i **/cpp/sandbox/generated/** -w tooling/ -- "just generate-sandbox && cd cpp/build && printf 'Building... ' && ninja --quiet sandbox_exec && printf 'done.\n\n' && ./sandbox_exec"

@watch-exec-sandbox-python:
    watchexec -c -w models/sandbox/ -w python/run_sandbox.py -w tooling/ -- "just generate-sandbox && echo "" && ./python/run_sandbox.py"

@validate: build-all test run-sandbox-quiet benchmark

validate-with-no-changes: validate
    #!/usr/bin/env bash
    set -euo pipefail

    if [[ `git status --porcelain` ]]; then
      echo "there are changed files"
      exit 1
    fi

@format:
    find . \( -name generated -prune \) -o \( -name "*.h" -o -name "*.cc" \) -exec clang-format -i {} \;


@watch-python-test:
    watchexec -c -w models/test/ -w python/ -i **/__pycache__/** -w tooling/ --on-busy-update do-nothing -- "just python-test"
