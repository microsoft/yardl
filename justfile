set shell := ['bash', '-ceuo', 'pipefail']

cpp_version := "17"

matlab := "disabled"
matlab-test-cmd := if matlab != "disabled" { "run-matlab-command run_tests" } else { "echo Skipping Matlab tests..." }
matlab-sandbox-cmd := if matlab != "disabled" { "run-matlab-command run_sandbox" } else { "echo Skipping Matlab sandbox..." }
benchmark-cmd := if matlab != "disabled" { "python python/benchmark.py --include-matlab" } else { "python python/benchmark.py" }


@default: validate
  echo $'\n\e[1;34mNote: you can run \'just test\' to a run an inner-loop subset of the complete validation'

@ensure-build-dir:
    mkdir -p cpp/build

@configure: ensure-build-dir
    cd cpp/build; \
    cmake -GNinja -D CMAKE_CXX_STANDARD={{ cpp_version }} ..

@install:
    cd tooling/cmd/yardl; \
    go install -ldflags="-s -w" .

@generate: install
    cd models/test && yardl generate

@generate-sandbox: install
    cd models/sandbox && yardl generate

@generate-remote-import: install
    cd models/remote-import && yardl generate

@generate-evolution: install
    cd models/evolution/model_v0 && yardl generate --quiet 2>&1
    cd models/evolution/model_v1 && yardl generate --quiet 2>&1
    cd models/evolution/model_v2 && yardl generate --quiet 2>&1

@generate-ndarray: install
    cd models/ndarray; \
    ln -sf ../test/unittests.yml ../test/benchmarking.yml ./; \
    yardl generate

@build-sandbox: generate-sandbox ensure-build-dir
    cd cpp/build && ninja sandbox_exec

@run-sandbox: build-sandbox
    cd cpp/build && ./sandbox_exec

@run-sandbox-quiet: build-sandbox
    cd cpp/build && ./sandbox_exec > /dev/null

@run-sandbox-python: generate-sandbox
    python python/run_sandbox.py

@run-sandbox-python-quiet: generate-sandbox
    python python/run_sandbox.py > /dev/null

@run-sandbox-matlab: generate-sandbox
    cd matlab; \
    {{ matlab-sandbox-cmd }}

@run-sandbox-matlab-quiet: generate-sandbox
    cd matlab; \
    {{ matlab-sandbox-cmd }} > /dev/null

@build-all: generate generate-sandbox generate-remote-import generate-evolution configure
    cd cpp/build && ninja

@tooling-test:
    cd tooling; \
    go test ./... | { grep -v "\\[no test files\\]" || true; }

@watch-tooling-test:
    cd tooling; \
    watchexec -r -c -w . -- 'go test ./... | { grep -v "\\[[no test files\\]" || true; }'

@build-translator:
    cd cpp/build; \
    ninja translator; \

@python-test: generate build-translator
    cd python; \
    python3 -m pytest tests

@cpp-test: generate ensure-build-dir
    cd cpp/build; \
    ninja tests; \
    ./tests --gtest_brief=1

@cpp-test-ndarray: generate-ndarray ensure-build-dir
    cd cpp/build; \
    ninja tests; \
    ./tests --gtest_brief=1

@matlab-test: generate build-translator
    cd matlab/test; \
    {{ matlab-test-cmd }}

@evolution-test: generate-evolution ensure-build-dir
    cd cpp/build; \
    ninja evolution/all; \
    python ../evolution/test-evolution.py

@test: tooling-test cpp-test python-test matlab-test evolution-test cpp-test-ndarray

@benchmark: generate ensure-build-dir
    cd cpp/build; \
    ninja benchmark; \
    cd ../..; \
    {{ benchmark-cmd }}

@watch-generate-test: install
    watchexec -r -c -w tooling/ -- "just install && cd models/test && yardl generate --watch"

@watch-generate-sandbox:
    watchexec -r -c -w tooling/ -- "just install && cd models/sandbox && yardl generate --watch"

@watch-exec-sandbox: configure
    watchexec -c -w models/sandbox/ -w cpp/ -i **/cpp/build/** -i **/cpp/test/** -i **/cpp/sandbox/generated/** -w tooling/ -- "just generate-sandbox && cd cpp/build && printf 'Building... ' && ninja --quiet sandbox_exec && printf 'done.\n\n' && ./sandbox_exec"

@watch-exec-sandbox-python:
    watchexec -c -w models/sandbox/ -w python/run_sandbox.py -w tooling/ -- "just generate-sandbox && echo "" && python python/run_sandbox.py"

type-check: generate generate-sandbox
    #! /usr/bin/env bash
    set -euo pipefail
    echo "Running Pyright..."
    cd python
    pyright .

@validate: build-all test type-check run-sandbox-quiet run-sandbox-python-quiet run-sandbox-matlab-quiet benchmark

validate-with-no-changes: validate
    #!/usr/bin/env bash
    set -euo pipefail

    if [[ `git status --porcelain` ]]; then
      echo "there are changed files:"
      git status --porcelain
      exit 1
    fi

@format:
    find . \( -name generated -prune \) -o \( -name "*.h" -o -name "*.cc" \) -exec clang-format -i {} \;


@watch-python-test:
    watchexec -c -w models/test/ -w python/ -i **/__pycache__/** -w tooling/ --on-busy-update do-nothing -- "just python-test"

@start-docs-website:
    cd docs && npm install && npm run docs:dev
