// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package python

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/internal/python/binary"
	"github.com/microsoft/yardl/tooling/internal/python/common"
	"github.com/microsoft/yardl/tooling/internal/python/ndjson"
	"github.com/microsoft/yardl/tooling/internal/python/protocols"
	"github.com/microsoft/yardl/tooling/internal/python/types"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

//go:embed static_files/*
var staticFiles embed.FS

func Generate(env *dsl.Environment, options packaging.PythonCodegenOptions) error {
	common.AnnotateGenerics(env)

	err := os.MkdirAll(options.OutputDir, 0775)
	if err != nil {
		return err
	}

	topNamespace := env.GetTopLevelNamespace()
	topPackageDir := path.Join(options.OutputDir, formatting.ToSnakeCase(topNamespace.Name))
	if err := iocommon.CopyEmbeddedStaticFiles(topPackageDir, options.InternalSymlinkStaticFiles, staticFiles); err != nil {
		return err
	}

	for _, ns := range env.Namespaces {
		packageDir := topPackageDir
		if !ns.IsTopLevel {
			packageDir = path.Join(packageDir, formatting.ToSnakeCase(ns.Name))
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

	// Write __init__.py
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

	if err := ndjson.WriteNDJson(ns, packageDir); err != nil {
		return err
	}

	return nil
}

func writePackageInitFile(ns *dsl.Namespace, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`# pyright: reportUnusedImport=false`)

	w.WriteStringln(`from typing import Tuple as _Tuple
import re as _re
import numpy as _np

_MIN_NUMPY_VERSION = (1, 22, 0)

def _parse_version(version: str) -> _Tuple[int, ...]:
    try:
        return tuple(map(int, version.split(".")))
    except ValueError:
        # ignore any prerelease suffix
        version = _re.sub(r"[^0-9.]", "", version)
        return tuple(map(int, version.split(".")))

if _parse_version(_np.__version__) < _MIN_NUMPY_VERSION:
    raise ImportError(f"Your installed numpy version is {_np.__version__}, but version >= {'.'.join(str(i) for i in _MIN_NUMPY_VERSION)} is required.")
`)

	relativePath := ".."
	if ns.IsTopLevel {
		relativePath = "."
	}
	fmt.Fprintf(w, "from %syardl_types import *\n", relativePath)

	for _, ref := range ns.GetAllChildReferences() {
		fmt.Fprintf(w, "from %s import %s\n", relativePath, common.NamespaceIdentifierName(ref.Name))
	}

	typesMembers := make([]string, 0)
	typesMembers = append(typesMembers, "get_dtype")
	for _, t := range ns.TypeDefinitions {
		typesMembers = append(typesMembers, common.TypeIdentifierName(t.GetDefinitionMeta().Name))
	}

	unions := make(map[string]interface{})
	dsl.Visit(ns, func(self dsl.Visitor, node dsl.Node) {
		switch node := node.(type) {
		case *dsl.NamedType:
			if gt, ok := node.Type.(*dsl.GeneralizedType); ok && gt.Cases.IsUnion() {
				// We use the alias name for the union type, which will be imported
				// below.
				return
			}
		case *dsl.GeneralizedType:
			if node.Cases.IsUnion() {
				unionClassName, _ := common.UnionClassName(node)
				if _, ok := unions[unionClassName]; !ok {
					unions[unionClassName] = nil
					typesMembers = append(typesMembers, unionClassName)
				}
			}
		}
		self.VisitChildren(node)
	})

	sort.Slice(typesMembers, func(i, j int) bool {
		return typesMembers[i] < typesMembers[j]
	})

	fmt.Fprintf(w, "from .types import (\n")
	w.Indented(func() {
		for _, t := range typesMembers {
			fmt.Fprintf(w, "%s,\n", t)
		}
	})
	fmt.Fprintf(w, ")\n")

	var protocolsMembers []string
	for _, p := range ns.Protocols {
		protocolsMembers = append(protocolsMembers, common.AbstractWriterName(p), common.AbstractReaderName(p))
	}
	if ns.IsTopLevel && len(protocolsMembers) > 0 {
		sort.Slice(protocolsMembers, func(i, j int) bool {
			return protocolsMembers[i] < protocolsMembers[j]
		})

		fmt.Fprintf(w, "from .protocols import (\n")
		w.Indented(func() {
			for _, p := range protocolsMembers {
				fmt.Fprintf(w, "%s,\n", p)
			}
		})
		fmt.Fprintf(w, ")\n")

		for i, p := range ns.Protocols {
			protocolsMembers[i*2] = binary.BinaryWriterName(p)
			protocolsMembers[i*2+1] = binary.BinaryReaderName(p)
		}

		sort.Slice(protocolsMembers, func(i, j int) bool {
			return protocolsMembers[i] < protocolsMembers[j]
		})

		fmt.Fprintf(w, "from .binary import (\n")
		w.Indented(func() {
			for _, p := range protocolsMembers {
				fmt.Fprintf(w, "%s,\n", p)
			}
		})
		fmt.Fprintf(w, ")\n")

		for i, p := range ns.Protocols {
			protocolsMembers[i*2] = ndjson.NDJsonWriterName(p)
			protocolsMembers[i*2+1] = ndjson.NDJsonReaderName(p)
		}

		sort.Slice(protocolsMembers, func(i, j int) bool {
			return protocolsMembers[i] < protocolsMembers[j]
		})

		fmt.Fprintf(w, "from .ndjson import (\n")
		w.Indented(func() {
			for _, p := range protocolsMembers {
				fmt.Fprintf(w, "%s,\n", p)
			}
		})
		fmt.Fprintf(w, ")\n")
	} else {
		w.WriteStringln("from . import binary")
		w.WriteStringln("from . import ndjson")
	}

	return iocommon.WriteFileIfNeeded(path.Join(packageDir, "__init__.py"), b.Bytes(), 0644)
}
