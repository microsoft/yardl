// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package mocks

import (
	"bytes"
	"fmt"
	"path"

	"github.com/microsoft/yardl/tooling/internal/cpp/binary"
	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/cpp/hdf5"
	"github.com/microsoft/yardl/tooling/internal/cpp/ndjson"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

func WriteMocks(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	err := writeFactories(env, options)
	if err != nil {
		return err
	}

	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`#include <functional>
#include <queue>

#include <gtest/gtest.h>

#include "../yardl_testing.h"
#include "binary/protocols.h"
#include "hdf5/protocols.h"
#include "ndjson/protocols.h"
#include "types.h"
`)

	for _, ns := range env.Namespaces {
		if !ns.IsTopLevel {
			continue
		}
		fmt.Fprintf(w, "namespace %s {\n", common.NamespaceIdentifierName(ns.Name))
		w.WriteStringln("namespace {")
		formatting.Delimited(w, "\n", ns.Protocols, func(w *formatting.IndentedWriter, i int, p *dsl.ProtocolDefinition) {
			writeProtocolTestWriter(w, p)
		})
		w.WriteStringln("} // namespace")
		fmt.Fprintf(w, "} // namespace %s\n\n", common.NamespaceIdentifierName(ns.Name))
	}

	w.WriteStringln("namespace yardl::testing {")
	for _, ns := range env.Namespaces {
		if !ns.IsTopLevel {
			continue
		}
		for _, protocol := range ns.Protocols {
			w.WriteStringln("template<>")
			fmt.Fprintf(w, "std::unique_ptr<%s> CreateValidatingWriter<%s>(Format format, std::string const& filename) {\n", common.QualifiedAbstractWriterName(protocol), common.QualifiedAbstractWriterName(protocol))
			w.Indented(func() {
				fmt.Fprintf(w, "return std::make_unique<%s>(\n", qualifiedTestWriterName(protocol))
				w.Indented(func() {
					fmt.Fprintf(w, "CreateWriter<%s>(format, filename),\n", common.QualifiedAbstractWriterName(protocol))
					fmt.Fprintf(w, "[format, filename](){ return CreateReader<%s>(format, filename);}\n", common.QualifiedAbstractReaderName(protocol))
				})
				w.WriteStringln(");")
			})
			w.WriteStringln("}\n")
		}
	}
	w.WriteStringln("}")

	definitionsPath := path.Join(options.SourcesOutputDir, "mocks.cc")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeProtocolTestWriter(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition) {
	writeProtocolMock(w, p)
	fmt.Fprintf(w, "class %s : public %s {\n", testWriterName(p), common.AbstractWriterName(p))
	w.Indented(func() {
		w.WriteStringln("public:")
		fmt.Fprintf(w, "%s(std::unique_ptr<%s> writer, std::function<std::unique_ptr<%s>()> create_reader) : writer_(std::move(writer)), create_reader_(create_reader) {\n}\n\n", testWriterName(p), common.QualifiedAbstractWriterName(p), common.AbstractReaderName(p))

		fmt.Fprintf(w, "~%s() {\n", testWriterName(p))
		w.Indented(func() {
			w.WriteStringln("if (!close_called_ && !std::uncaught_exceptions()) {")
			w.Indented(func() {
				fmt.Fprintf(w, "ADD_FAILURE() << \"Close() needs to be called on '%s' to verify mocks\";\n", testWriterName(p))
			})
			w.WriteStringln("}")
		})
		w.WriteString("}\n\n")

		w.WriteStringln("protected:")
		for _, step := range p.Sequence {
			writeMethodName := common.ProtocolWriteImplMethodName(step)
			fmt.Fprintf(w, "void %s(%s const& value) override {\n", writeMethodName, common.TypeSyntax(step.Type))
			w.Indented(func() {
				fmt.Fprintf(w, "writer_->%s(value);\n", common.ProtocolWriteMethodName(step))
				fmt.Fprintf(w, "mock_writer_.Expect%s(value);\n", writeMethodName)
			})

			w.WriteString("}\n\n")

			if step.IsStream() {

				fmt.Fprintf(w, "void %s(std::vector<%s> const& values) override {\n", writeMethodName, common.TypeSyntax(step.Type))
				w.Indented(func() {
					fmt.Fprintf(w, "writer_->%s(values);\n", common.ProtocolWriteMethodName(step))
					w.WriteStringln("for (auto const& v : values) {")
					w.Indented(func() {
						fmt.Fprintf(w, "mock_writer_.Expect%s(v);\n", writeMethodName)
					})
					w.WriteStringln("}")
				})
				w.WriteString("}\n\n")

				endMethodName := common.ProtocolWriteEndImplMethodName(step)
				fmt.Fprintf(w, "void %s() override {\n", endMethodName)
				w.Indented(func() {
					fmt.Fprintf(w, "writer_->%s();\n", common.ProtocolWriteEndMethodName(step))
					fmt.Fprintf(w, "mock_writer_.Expect%s();\n", endMethodName)
				})
				w.WriteString("}\n\n")
			}
		}

		w.WriteString("void CloseImpl() override {\n")
		w.Indented(func() {
			w.WriteStringln("close_called_ = true;")
			w.WriteStringln("writer_->Close();")
			fmt.Fprintf(w, "std::unique_ptr<%s> reader = create_reader_();\n", common.AbstractReaderName(p))
			w.WriteString("reader->CopyTo(mock_writer_")

			// set a mix of values for buffer sizes
			for i, s := range p.Sequence {
				if s.IsStream() {
					bufferSize := 1
					if i%2 == 0 {
						bufferSize = i + 2
					}

					fmt.Fprintf(w, ", %d", bufferSize)
				}
			}

			w.WriteStringln(");")
			w.WriteStringln("mock_writer_.Verify();")
		})
		w.WriteString("}\n\n")

		w.WriteStringln("private:")
		fmt.Fprintf(w, "std::unique_ptr<%s> writer_;\n", common.QualifiedAbstractWriterName(p))
		fmt.Fprintf(w, "std::function<std::unique_ptr<%s>()> create_reader_;\n", common.QualifiedAbstractReaderName(p))
		fmt.Fprintf(w, "Mock%sWriter mock_writer_;\n", p.Name)
		w.WriteString("bool close_called_ = false;\n")
	})
	w.WriteString("};\n")
}

func writeProtocolMock(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition) {
	fmt.Fprintf(w, "class Mock%sWriter : public %s {\n", p.Name, common.AbstractWriterName(p))
	w.Indented(func() {
		w.WriteStringln("public:")
		for _, step := range p.Sequence {
			fmt.Fprintf(w, "void %s (%s const& value) override {\n", common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
			w.Indented(func() {
				fmt.Fprintf(w, "if (%s_expected_values_.empty()) {\n", common.ProtocolWriteImplMethodName(step))
				w.Indented(func() {
					fmt.Fprintf(w, "throw std::runtime_error(\"Unexpected call to %s\");\n", common.ProtocolWriteImplMethodName(step))
				})
				w.WriteString("}\n")
				fmt.Fprintf(w, "if (%s_expected_values_.front() != value) {\n", common.ProtocolWriteImplMethodName(step))
				w.Indented(func() {
					fmt.Fprintf(w, "throw std::runtime_error(\"Unexpected argument value for call to %s\");\n", common.ProtocolWriteImplMethodName(step))
				})
				w.WriteString("}\n")
				fmt.Fprintf(w, "%s_expected_values_.pop();\n", common.ProtocolWriteImplMethodName(step))
			})
			w.WriteString("}\n\n")

			fmt.Fprintf(w, "std::queue<%s> %s_expected_values_;\n\n", common.TypeSyntax(step.Type), common.ProtocolWriteImplMethodName(step))

			fmt.Fprintf(w, "void Expect%s (%s const& value) {\n", common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
			w.Indented(func() {
				fmt.Fprintf(w, "%s_expected_values_.push(value);\n", common.ProtocolWriteImplMethodName(step))
			})

			w.WriteStringln("}\n")

			if step.IsStream() {
				fmt.Fprintf(w, "void %s () override {\n", common.ProtocolWriteEndImplMethodName(step))
				w.Indented(func() {
					fmt.Fprintf(w, "if (--%s_expected_call_count_ < 0) {\n", common.ProtocolWriteEndImplMethodName(step))
					w.Indented(func() {
						fmt.Fprintf(w, "throw std::runtime_error(\"Unexpected call to %s\");\n", common.ProtocolWriteEndImplMethodName(step))
					})
					w.WriteString("}\n")
				})
				w.WriteStringln("}\n")

				fmt.Fprintf(w, "int %s_expected_call_count_ = 0;\n\n", common.ProtocolWriteEndImplMethodName(step))

				fmt.Fprintf(w, "void Expect%s () {\n", common.ProtocolWriteEndImplMethodName(step))
				w.Indented(func() {
					fmt.Fprintf(w, "%s_expected_call_count_++;\n", common.ProtocolWriteEndImplMethodName(step))
				})
				w.WriteStringln("}\n")
			}
		}

		w.WriteStringln("void Verify() {")
		w.Indented(func() {
			for _, step := range p.Sequence {
				fmt.Fprintf(w, "if (!%s_expected_values_.empty()) {\n", common.ProtocolWriteImplMethodName(step))
				w.Indented(func() {
					fmt.Fprintf(w, "throw std::runtime_error(\"Expected call to %s was not received\");\n", common.ProtocolWriteImplMethodName(step))
				})
				w.WriteString("}\n")

				if step.IsStream() {
					fmt.Fprintf(w, "if (%s_expected_call_count_ > 0) {\n", common.ProtocolWriteEndImplMethodName(step))
					w.Indented(func() {
						fmt.Fprintf(w, "throw std::runtime_error(\"Expected call to %s was not received\");\n", common.ProtocolWriteEndImplMethodName(step))
					})
					w.WriteString("}\n")
				}
			}
		})
		w.WriteStringln("}")

	})
	w.WriteString("};\n\n")
}

func testWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("Test%s", common.AbstractWriterName(p))
}

func qualifiedTestWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::%s", common.TypeNamespaceIdentifierName(p), testWriterName(p))
}

func writeFactories(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`#include <functional>
#include "../factories.h"
#include "binary/protocols.h"
#include "hdf5/protocols.h"
#include "ndjson/protocols.h"
`)

	w.WriteStringln("namespace yardl::testing {")
	for _, ns := range env.Namespaces {
		if !ns.IsTopLevel {
			continue
		}
		for _, protocol := range ns.Protocols {
			w.WriteStringln("template<>")
			fmt.Fprintf(w, "std::unique_ptr<%s> CreateWriter<%s>(Format format, std::string const& filename) {\n", common.QualifiedAbstractWriterName(protocol), common.QualifiedAbstractWriterName(protocol))
			w.Indented(func() {
				w.WriteStringln("switch (format) {")
				w.WriteStringln("case Format::kHdf5:")
				w.Indented(func() {
					fmt.Fprintf(w, "return std::make_unique<%s>(filename);\n", hdf5.QualifiedHdf5WriterClassName(protocol))
				})
				w.WriteStringln("case Format::kBinary:")
				w.Indented(func() {
					fmt.Fprintf(w, "return std::make_unique<%s>(filename);\n", binary.QualifiedBinaryWriterClassName(protocol))
				})
				w.WriteStringln("case Format::kNDJson:")
				w.Indented(func() {
					fmt.Fprintf(w, "return std::make_unique<%s>(filename);\n", ndjson.QualifiedNDJsonWriterClassName(protocol))
				})
				w.WriteStringln("default:")
				w.Indented(func() {
					w.WriteStringln("throw std::runtime_error(\"Unknown format\");")
				})

				w.WriteStringln("}")
			})
			w.WriteStringln("}\n")

			w.WriteStringln("template<>")
			fmt.Fprintf(w, "std::unique_ptr<%s> CreateReader<%s>(Format format, std::string const& filename) {\n", common.QualifiedAbstractReaderName(protocol), common.QualifiedAbstractReaderName(protocol))
			w.Indented(func() {
				w.WriteStringln("switch (format) {")
				w.WriteStringln("case Format::kHdf5:")
				w.Indented(func() {
					fmt.Fprintf(w, "return std::make_unique<%s>(filename);\n", hdf5.QualifiedHdf5ReaderClassName(protocol))
				})
				w.WriteStringln("case Format::kBinary:")
				w.Indented(func() {
					fmt.Fprintf(w, "return std::make_unique<%s>(filename);\n", binary.QualifiedBinaryReaderClassName(protocol))
				})
				w.WriteStringln("case Format::kNDJson:")
				w.Indented(func() {
					fmt.Fprintf(w, "return std::make_unique<%s>(filename);\n", ndjson.QualifiedNDJsonReaderClassName(protocol))
				})
				w.WriteStringln("default:")
				w.Indented(func() {
					w.WriteStringln("throw std::runtime_error(\"Unknown format\");")
				})

				w.WriteStringln("}")
			})
			w.WriteStringln("}\n")
		}
	}
	w.WriteStringln("}")

	definitionsPath := path.Join(options.SourcesOutputDir, "factories.cc")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}
