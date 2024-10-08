// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package include

import (
	"embed"
	_ "embed"
	"os"
	"path"
	"text/template"

	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

//go:embed detail/*
var DetailHeaders embed.FS

//go:embed yardl.h.tmpl
var YardlHeaderTmpl string

const DefaultArrayHeader = "detail/ndarray/impl.h"

func GenerateYardlHeaders(options packaging.CppCodegenOptions) error {
	err := os.MkdirAll(path.Join(options.SourcesOutputDir, "yardl"), 0775)
	if err != nil {
		return err
	}

	err = iocommon.CopyEmbeddedStaticFiles(path.Join(options.SourcesOutputDir, "yardl", "detail"), false, DetailHeaders)
	if err != nil {
		return err
	}

	tmpl, err := template.New("yardl_h").Parse(YardlHeaderTmpl)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(path.Join(options.SourcesOutputDir, "yardl", "yardl.h"))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	arrayHeader := DefaultArrayHeader
	if options.OverrideArrayHeader != "" {
		arrayHeader = options.OverrideArrayHeader
	}
	data := struct {
		ArrayHeader string
	}{
		ArrayHeader: arrayHeader,
	}
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return err
	}

	return nil
}
