package binary

import (
	"bytes"
	"path"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/internal/python/common"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func WriteBinary(ns *dsl.Namespace, st dsl.SymbolTable, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)
	w.WriteStringln(`
import abc
import collections.abc
import datetime
import numpy as np
`)

	writeProtocols(w, ns, st)

	definitionsPath := path.Join(packageDir, "binary.py")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}
