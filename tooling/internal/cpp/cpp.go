// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cpp

import (
	_ "embed"
	"os"

	"github.com/microsoft/yardl/tooling/internal/cpp/binary"
	"github.com/microsoft/yardl/tooling/internal/cpp/hdf5"
	"github.com/microsoft/yardl/tooling/internal/cpp/include"
	"github.com/microsoft/yardl/tooling/internal/cpp/mocks"
	"github.com/microsoft/yardl/tooling/internal/cpp/ndjson"
	"github.com/microsoft/yardl/tooling/internal/cpp/protocols"
	"github.com/microsoft/yardl/tooling/internal/cpp/translator"
	"github.com/microsoft/yardl/tooling/internal/cpp/types"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
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

	err = hdf5.WriteHdf5(env, options)
	if err != nil {
		return err
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
