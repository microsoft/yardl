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
		packageDir := path.Join(options.OutputDir, common.PackageDir(ns.Name))
		if err := os.MkdirAll(packageDir, 0775); err != nil {
			return err
		}

		topLevelFileWriter := &common.MatlabFileWriter{PackageDir: packageDir}
		if err := types.WriteTypes(topLevelFileWriter, ns, env.SymbolTable); err != nil {
			return err
		}

		if ns.IsTopLevel {
			if err := protocols.WriteProtocols(topLevelFileWriter, ns, env.SymbolTable); err != nil {
				return err
			}

			if options.InternalGenerateMocks {
				mocksDir := path.Join(packageDir, common.PackageDir("testing"))
				if err := os.MkdirAll(mocksDir, 0775); err != nil {
					return err
				}
				fw := &common.MatlabFileWriter{PackageDir: mocksDir}
				if err := mocks.WriteMocks(fw, ns); err != nil {
					return err
				}
			}
		}

		binaryDir := path.Join(packageDir, common.PackageDir("binary"))
		if err := os.MkdirAll(binaryDir, 0775); err != nil {
			return err
		}
		binaryFileWriter := &common.MatlabFileWriter{PackageDir: binaryDir}
		if err := binary.WriteBinary(binaryFileWriter, ns); err != nil {
			return err
		}
	}

	return nil
}
