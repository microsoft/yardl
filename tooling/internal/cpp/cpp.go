// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cpp

import (
	"bytes"
	_ "embed"
	"os"
	"path"

	"github.com/microsoft/yardl/tooling/internal/cpp/binary"
	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/cpp/hdf5"
	"github.com/microsoft/yardl/tooling/internal/cpp/include"
	"github.com/microsoft/yardl/tooling/internal/cpp/mocks"
	"github.com/microsoft/yardl/tooling/internal/cpp/ndjson"
	"github.com/microsoft/yardl/tooling/internal/cpp/protocols"
	"github.com/microsoft/yardl/tooling/internal/cpp/translator"
	"github.com/microsoft/yardl/tooling/internal/cpp/types"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
	"github.com/rs/zerolog/log"
)

func Generate(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	err := os.MkdirAll(options.SourcesOutputDir, 0775)
	if err != nil {
		return err
	}

	err = include.GenerateYardlHeaders(options)
	if err != nil {
		return err
	}

	err = types.WriteTypes(env, options)
	if err != nil {
		return err
	}

	err = protocols.WriteProtocols(env, options)
	if err != nil {
		return err
	}

	err = ndjson.WriteNdJson(env, options)
	if err != nil {
		return err
	}

	err = binary.WriteBinary(env, options)
	if err != nil {
		return err
	}

	if modelHasRecursiveTypes(env) {
		log.Warn().Msg("Model has recursive types, skipping HDF5 code generation.")

		// This is a temporary workaround to avoid generating HDF5 code for models with recursive types.
		// We write an empty protocols.cc file to avoid breaking the main Yardl test build.
		options = options.ChangeOutputDir("hdf5")
		if err := os.MkdirAll(options.SourcesOutputDir, 0775); err != nil {
			return err
		}
		b := bytes.Buffer{}
		w := formatting.NewIndentedWriter(&b, "  ")
		common.WriteGeneratedFileHeader(w)
		filePath := path.Join(options.SourcesOutputDir, "protocols.cc")
		return iocommon.WriteFileIfNeeded(filePath, b.Bytes(), 0644)
	} else {
		err = hdf5.WriteHdf5(env, options)
		if err != nil {
			return err
		}
	}

	if options.InternalGenerateMocks {
		err = mocks.WriteMocks(env, options)
		if err != nil {
			return err
		}
	}

	if options.InternalGenerateTranslator {
		err = translator.WriteTranslator(env, options)
		if err != nil {
			return err
		}
	}

	if options.GenerateCMakeLists {
		err = writeCMakeLists(env, options)
	}

	return err
}

func modelHasRecursiveTypes(env *dsl.Environment) bool {
	hasRecursiveTypes := false
	dsl.Visit(env, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case *dsl.SimpleType:
			if t.IsRecursive {
				hasRecursiveTypes = true
				return
			}
		}
		if !hasRecursiveTypes {
			self.VisitChildren(node)
		}
	})
	return hasRecursiveTypes
}
