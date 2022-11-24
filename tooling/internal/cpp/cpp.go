// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cpp

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/microsoft/yardl/tooling/internal/cpp/binary"
	"github.com/microsoft/yardl/tooling/internal/cpp/hdf5"
	"github.com/microsoft/yardl/tooling/internal/cpp/mocks"
	"github.com/microsoft/yardl/tooling/internal/cpp/protocols"
	"github.com/microsoft/yardl/tooling/internal/cpp/types"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

//go:embed include/*
var includes embed.FS

func Generate(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	err := os.MkdirAll(options.SourcesOutputDir, 0775)
	if err != nil {
		return err
	}

	err = copyStaticHeaderFiles(options)
	if err != nil {
		return err
	}

	err = types.WriteTypes(env, options)
	if err != nil {
		return err
	}

	err = protocols.WriteProtocols(env, options)
	if err != nil {
		return err
	}

	err = binary.WriteBinary(env, options)
	if err != nil {
		return err
	}

	err = hdf5.WriteHdf5(env, options)
	if err != nil {
		return err
	}

	if options.InternalGenerateMocks {
		err = mocks.WriteMocks(env, options)
		if err != nil {
			return err
		}
	}

	if options.GenerateCMakeLists {
		err = writeCMakeLists(env, options)
	}

	return err
}

func getCompileTimeIncludeDirPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(filepath.Dir(filename), "include")
}

func copyStaticHeaderFiles(options packaging.CppCodegenOptions) error {
	destinationDir := path.Join(options.SourcesOutputDir, "yardl")
	if !options.InternalSymlinkStaticHeaders {
		return copyEmbeddedIncludeDir("include", destinationDir)
	}

	includeDir, _ := filepath.Rel(options.SourcesOutputDir, getCompileTimeIncludeDirPath())

	stat, err := os.Lstat(destinationDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if stat.Mode()&os.ModeSymlink == 0 {
			return fmt.Errorf("static headers destination dir %s exists and is not a symlink", destinationDir)
		}

		currentTarget, err := os.Readlink(destinationDir)
		if err != nil {
			return err
		}
		if currentTarget == includeDir {
			return nil
		}

		err = os.Remove(destinationDir)
		if err != nil {
			return err
		}
	}

	return os.Symlink(includeDir, destinationDir)
}

func copyEmbeddedIncludeDir(sourceDir, destDir string) error {
	entries, err := includes.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return nil
	}

	err = os.MkdirAll(destDir, 0775)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {

			err = copyEmbeddedIncludeDir(path.Join(sourceDir, entry.Name()), path.Join(destDir, entry.Name()))
			if err != nil {
				return err
			}
		} else {
			content, err := includes.ReadFile(path.Join(sourceDir, entry.Name()))
			if err != nil {
				return err
			}
			err = iocommon.WriteFileIfNeeded(path.Join(destDir, entry.Name()), content, 0664)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
