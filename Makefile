.ONESHELL:
SHELL = /bin/bash
.SHELLFLAGS = -ecuo pipefail

CPP_VERSION=17

.SILENT:  testc tooling-test watch-generate-cpp-test

configure:
	cd cpp
	mkdir -p build
	cd build
	cmake -GNinja -D CMAKE_CXX_STANDARD=${CPP_VERSION} ..

install:
	cd tooling/cmd/yardl
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-s -w" .

generate: install
	cd models/test
	yardl generate

generate-sandbox: install
	cd models/sandbox
	yardl generate

build-sandbox: generate-sandbox
	cd cpp/build
	ninja sandbox_exec

run-sandbox: build-sandbox
	cd cpp/build
	./sandbox_exec

run-sandbox-quiet: build-sandbox
	cd cpp/build
	./sandbox_exec > /dev/null

build-all: generate generate-sandbox
	cd cpp/build
	ninja

tooling-test:
	cd tooling
	go test ./... | { grep -v "\\[[no test files\\]" || true; }

watch-tooling-test:
	cd tooling
	watchexec -r -c -w . -- 'go test ./... | { grep -v "\\[[no test files\\]" || true; }'

test: tooling-test generate
	cd cpp/build
	ninja tests
	./tests --gtest_brief=1

benchmark: generate
	cd cpp/build
	ninja benchmark
	./benchmark

watch-generate-test: install
	watchexec -r -c -w tooling/ -- "make install && cd models/test && yardl generate --watch"

watch-generate-sandbox:
	watchexec -r -c -w tooling/ -- "make install && cd models/sandbox && yardl generate --watch"

watch-exec-sandbox: configure
	watchexec -c -w models/sandbox/ -w cpp/ -i **/cpp/build/** -i **/cpp/test/** -i **/cpp/sandbox/generated/** -w tooling/ -- "make generate-sandbox && cd cpp/build && printf 'Building... ' && ninja --quiet sandbox_exec && printf 'done.\n\n' && ./sandbox_exec"

validate: generate generate-sandbox configure build-all test run-sandbox-quiet benchmark

validate-with-no-changes: validate
	if [[ `git status --porcelain` ]]; then
		echo "there are changed files"
		exit 1
	fi
