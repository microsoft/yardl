// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package include

import (
	"bytes"
	"embed"
	_ "embed"
	"io"
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

	arrayHeader := DefaultArrayHeader
	if options.OverrideArrayHeader != "" {
		arrayHeader = options.OverrideArrayHeader
	}
	data := struct {
		ArrayHeader string
	}{
		ArrayHeader: arrayHeader,
	}

	b := bytes.Buffer{}
	w := io.Writer(&b)
	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return iocommon.WriteFileIfNeeded(path.Join(options.SourcesOutputDir, "yardl", "yardl.h"), b.Bytes(), 0644)
}
