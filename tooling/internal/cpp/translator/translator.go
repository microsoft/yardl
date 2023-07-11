// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package translator

import (
	"bytes"
	"fmt"
	"path"

	"github.com/microsoft/yardl/tooling/internal/cpp/binary"
	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/cpp/ndjson"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

func WriteTranslator(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`#include <iostream>

#include "../format.h"
#include "binary/protocols.h"
#include "hdf5/protocols.h"
#include "ndjson/protocols.h"
`)

	w.WriteStringln("namespace yardl::testing {")
	w.WriteStringln("void TranslateStream(std::string const& protocol_name, yardl::testing::Format input_format, std::istream& input, yardl::testing::Format output_format, std::ostream& output) {")
	w.Indented(func() {
		w.WriteStringln("switch (input_format) {")
		w.WriteStringln("case yardl::testing::Format::kBinary:")
		w.Indented(func() {
			w.WriteStringln("break;")
		})
		w.WriteStringln("case yardl::testing::Format::kNDJson:")
		w.Indented(func() {
			w.WriteStringln("break;")
		})
		w.WriteStringln("default:")
		w.Indented(func() {
			w.WriteStringln("throw std::runtime_error(\"Unsupported input format\");")
		})
		w.WriteStringln("}\n")

		for _, ns := range env.Namespaces {
			for _, protocol := range ns.Protocols {
				fmt.Fprintf(w, "if (protocol_name == \"%s\") {\n", protocol.Name)
				w.Indented(func() {
					w.WriteStringln("auto reader = input_format == yardl::testing::Format::kBinary")
					w.Indented(func() {
						fmt.Fprintf(w, "? std::unique_ptr<%s>(new %s(input))\n", common.QualifiedAbstractReaderName(protocol), binary.QualifiedBinaryReaderClassName(protocol))
						fmt.Fprintf(w, ": std::unique_ptr<%s>(new %s(input));\n", common.QualifiedAbstractReaderName(protocol), ndjson.QualifiedNDJsonReaderClassName(protocol))
					})
					w.WriteStringln("")
					w.WriteStringln("auto writer = output_format == yardl::testing::Format::kBinary")
					w.Indented(func() {
						fmt.Fprintf(w, "? std::unique_ptr<%s>(new %s(output))\n", common.QualifiedAbstractWriterName(protocol), binary.QualifiedBinaryWriterClassName(protocol))
						fmt.Fprintf(w, ": std::unique_ptr<%s>(new %s(output));\n", common.QualifiedAbstractWriterName(protocol), ndjson.QualifiedNDJsonWriterClassName(protocol))
					})

					w.WriteStringln("reader->CopyTo(*writer);")
					w.WriteStringln("return;")
				})
				w.WriteStringln("}")
			}
		}

		w.WriteStringln("throw std::runtime_error(\"Unsupported protocol \" + protocol_name);")
	})
	w.WriteStringln("}")
	w.WriteStringln("} // namespace yardl::testing")

	definitionsPath := path.Join(options.SourcesOutputDir, "translator_impl.cc")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}
