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

	if err := iocommon.CopyEmbeddedStaticFiles(options.OutputDir, options.InternalSymlinkStaticFiles, staticFiles); err != nil {
		return err
	}

	for _, ns := range env.Namespaces {
		packageDir := path.Join(options.OutputDir, common.PackageDir(ns.Name))

		if err := os.MkdirAll(packageDir, 0775); err != nil {
			return err
		}

		fw := &common.MatlabFileWriter{PackageDir: packageDir}
		err = writeNamespace(fw, ns, env.SymbolTable)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeNamespace(fw *common.MatlabFileWriter, ns *dsl.Namespace, st dsl.SymbolTable) error {
	if err := types.WriteTypes(fw, ns, st); err != nil {
		return err
	}

	if ns.IsTopLevel {
		if err := protocols.WriteProtocols(fw, ns, st); err != nil {
			return err
		}
	}

	if err := binary.WriteBinary(fw, ns); err != nil {
		return err
	}

	// if err := ndjson.WriteNDJson(ns, packageDir); err != nil {
	// 	return err
	// }

	return nil
}
