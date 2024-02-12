// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package matlab

import (
	"bytes"
	"embed"
	"os"
	"path"

	"github.com/microsoft/yardl/tooling/internal/formatting"
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
	// common.AnnotateGenerics(env)

	err := os.MkdirAll(options.OutputDir, 0775)
	if err != nil {
		return err
	}

	topNamespace := env.GetTopLevelNamespace()
	// topPackageDir := path.Join(options.OutputDir, common.PackageDir(topNamespace.Name))
	topPackageDir := path.Join(options.OutputDir, formatting.ToSnakeCase(topNamespace.Name))
	if err := iocommon.CopyEmbeddedStaticFiles(topPackageDir, options.InternalSymlinkStaticFiles, staticFiles); err != nil {
		return err
	}

	for _, ns := range env.Namespaces {
		packageDir := topPackageDir
		if !ns.IsTopLevel {
			packageDir = path.Join(packageDir, common.PackageDir(ns.Name))
		}
		err = writeNamespace(ns, env.SymbolTable, packageDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeNamespace(ns *dsl.Namespace, st dsl.SymbolTable, packageDir string) error {
	if err := os.MkdirAll(packageDir, 0775); err != nil {
		return err
	}

	if err := writePackageInitFile(ns, packageDir); err != nil {
		return err
	}

	if err := types.WriteTypes(ns, st, packageDir); err != nil {
		return err
	}

	if ns.IsTopLevel {
		if err := protocols.WriteProtocols(ns, st, packageDir); err != nil {
			return err
		}
	}

	if err := binary.WriteBinary(ns, packageDir); err != nil {
		return err
	}

	// if err := ndjson.WriteNDJson(ns, packageDir); err != nil {
	// 	return err
	// }

	return nil
}

func writePackageInitFile(ns *dsl.Namespace, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	return iocommon.WriteFileIfNeeded(path.Join(packageDir, "init.m"), b.Bytes(), 0644)
}
