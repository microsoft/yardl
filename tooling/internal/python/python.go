package python

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/internal/python/binary"
	"github.com/microsoft/yardl/tooling/internal/python/common"
	"github.com/microsoft/yardl/tooling/internal/python/protocols"
	"github.com/microsoft/yardl/tooling/internal/python/types"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

//go:embed static_files/*
var staticFiles embed.FS

func Generate(env *dsl.Environment, options packaging.PythonCodegenOptions) error {
	err := os.MkdirAll(options.OutputDir, 0775)
	if err != nil {
		return err
	}

	for _, ns := range env.Namespaces {
		err = writeNamespace(ns, env.SymbolTable, options)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeNamespace(ns *dsl.Namespace, st dsl.SymbolTable, options packaging.PythonCodegenOptions) error {
	packageDir := path.Join(options.OutputDir, formatting.ToSnakeCase(ns.Name))
	if err := os.MkdirAll(packageDir, 0775); err != nil {
		return err
	}

	// Write __init__.py
	if err := writePackageInitFile(packageDir, ns); err != nil {
		return err
	}

	iocommon.CopyEmbeddedStaticFiles(packageDir, false, staticFiles)

	if err := types.WriteTypes(ns, packageDir); err != nil {
		return err
	}

	if err := protocols.WriteProtocols(ns, st, packageDir); err != nil {
		return err
	}

	if err := binary.WriteBinary(ns, packageDir); err != nil {
		return err
	}

	return nil
}

func writePackageInitFile(packageDir string, ns *dsl.Namespace) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)

	typeNames := make([]string, 0)
	for _, t := range ns.TypeDefinitions {
		typeNames = append(typeNames, common.TypeDefinitionSyntax(t, ns.Name, false))
	}

	if len(typeNames) > 0 {
		fmt.Fprintf(w, "from .types import %s\n", strings.Join(typeNames, ", "))
	}

	protocolTypes := make([]string, 0)
	for _, p := range ns.Protocols {
		protocolTypes = append(protocolTypes, common.AbstractWriterName(p), common.AbstractReaderName(p))
	}

	if len(protocolTypes) > 0 {
		fmt.Fprintf(w, "from .protocols import %s\n", strings.Join(protocolTypes, ", "))

		for i, p := range ns.Protocols {
			protocolTypes[i*2] = binary.BinaryWriterName(p)
			protocolTypes[i*2+1] = binary.BinaryReaderName(p)
		}

		fmt.Fprintf(w, "from .binary import %s\n", strings.Join(protocolTypes, ", "))
	}

	return iocommon.WriteFileIfNeeded(path.Join(packageDir, "__init__.py"), b.Bytes(), 0644)
}
