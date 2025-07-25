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

//go:embed detail/binary/*
var DetailBinaryHeaders embed.FS

//go:embed detail/hdf5/*
var DetailHDF5Headers embed.FS

//go:embed detail/ndarray/*
var DetailArrayHeaders embed.FS

//go:embed detail/ndjson/*
var DetailNdJsonHeaders embed.FS

//go:embed yardl.h.tmpl
var YardlHeaderTmpl string

const DefaultArrayHeader = "detail/ndarray/impl.h"

func GenerateYardlHeaders(options packaging.CppCodegenOptions) error {
	err := os.MkdirAll(path.Join(options.SourcesOutputDir, "yardl"), 0775)
	if err != nil {
		return err
	}

	targetDir := path.Join(options.SourcesOutputDir, "yardl", "detail")

	err = iocommon.CopyEmbeddedStaticFiles(targetDir, false, DetailBinaryHeaders)
	if err != nil {
		return err
	}

	err = iocommon.CopyEmbeddedStaticFiles(targetDir, false, DetailArrayHeaders)
	if err != nil {
		return err
	}

	if options.GenerateHDF5 {
		err = iocommon.CopyEmbeddedStaticFiles(targetDir, false, DetailHDF5Headers)
		if err != nil {
			return err
		}
	}

	if options.GenerateNDJson {
		err = iocommon.CopyEmbeddedStaticFiles(targetDir, false, DetailNdJsonHeaders)
		if err != nil {
			return err
		}
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
