// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package matlab

import (
	"embed"
	"os"
	"path"

	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/internal/matlab/binary"
	"github.com/microsoft/yardl/tooling/internal/matlab/common"
	"github.com/microsoft/yardl/tooling/internal/matlab/mocks"
	"github.com/microsoft/yardl/tooling/internal/matlab/protocols"
	"github.com/microsoft/yardl/tooling/internal/matlab/types"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

//go:embed static_files/*
var staticFiles embed.FS

func Generate(env *dsl.Environment, options packaging.MatlabCodegenOptions) error {
	err := os.MkdirAll(options.OutputDir, 0775)
	if err != nil {
		return err
	}

	staticDir := path.Join(options.OutputDir, common.PackageDir("yardl"))
	if err := iocommon.CopyEmbeddedStaticFiles(staticDir, options.InternalSymlinkStaticFiles, staticFiles); err != nil {
		return err
	}

	for _, ns := range env.Namespaces {
		// Write package Types and Protocol definitions
		packageDir := path.Join(options.OutputDir, common.PackageDir(ns.Name))
		if err := updatePackage(packageDir, func(fw *common.MatlabFileWriter) error {
			if err := types.WriteTypes(fw, ns, env.SymbolTable); err != nil {
				return err
			}
			if ns.IsTopLevel {
				if err := protocols.WriteProtocols(fw, ns, env.SymbolTable); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}

		// Write binary serializers
		binaryDir := path.Join(packageDir, common.PackageDir("binary"))
		if err := updatePackage(binaryDir, func(fw *common.MatlabFileWriter) error {
			return binary.WriteBinary(fw, ns)
		}); err != nil {
			return err
		}

		// Write mocks and test support classes
		if ns.IsTopLevel && options.InternalGenerateMocks {
			mocksDir := path.Join(packageDir, common.PackageDir("testing"))
			if err := updatePackage(mocksDir, func(fw *common.MatlabFileWriter) error {
				return mocks.WriteMocks(fw, ns)
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

// Creates `+package` directory, writes the package implementation, and removes stale files.
func updatePackage(packageDir string, writePackageImpl func(*common.MatlabFileWriter) error) error {
	if err := os.MkdirAll(packageDir, 0775); err != nil {
		return err
	}
	fw := &common.MatlabFileWriter{PackageDir: packageDir}

	writePackageImpl(fw)

	return fw.RemoveStaleFiles()
}
