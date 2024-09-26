// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package hdf5

import (
	"bytes"
	"fmt"
	"path"

	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

func writeHeaderFile(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`#pragma once
#include <array>
#include <complex>
#include <optional>
#include <variant>
#include <vector>

#include "../protocols.h"
#include "../yardl/detail/hdf5/io.h"
`)

	for _, ns := range env.Namespaces {
		if !ns.IsTopLevel {
			continue
		}
		fmt.Fprintf(w, "namespace %s::hdf5 {\n", common.NamespaceIdentifierName(ns.Name))
		for _, protocol := range ns.Protocols {

			common.WriteComment(w, fmt.Sprintf("HDF5 writer for the %s protocol.", protocol.Name))
			common.WriteComment(w, protocol.Comment)
			writerClassName := Hdf5WriterClassName(protocol)
			fmt.Fprintf(w, "class %s : public %s, public yardl::hdf5::Hdf5Writer {\n", writerClassName, common.QualifiedAbstractWriterName(protocol))
			w.Indented(func() {
				w.WriteStringln("public:")
				fmt.Fprintf(w, "%s(std::string path);\n\n", writerClassName)

				hasUnionStream := false
				w.WriteStringln("protected:")
				for _, step := range protocol.Sequence {
					endMethodName := common.ProtocolWriteEndImplMethodName(step)
					fmt.Fprintf(w, "void %s(%s const& value) override;\n\n", common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))

					if step.IsStream() {
						isUnion := dsl.ToGeneralizedType(dsl.GetUnderlyingType(step.Type.(*dsl.GeneralizedType).ToScalar())).Cases.IsUnion()
						if isUnion {
							hasUnionStream = true
						} else {
							// batch optimization is not supported at the moment for unions because
							// we write each type to separate datasets.
							fmt.Fprintf(w, "void %s(std::vector<%s> const& values) override;\n\n", common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
						}
						fmt.Fprintf(w, "void %s() override;\n\n", endMethodName)
					}
				}

				if hasUnionStream {
					w.WriteStringln("public:")
					w.WriteString("void Flush() override;\n\n")
				}

				w.WriteStringln("private:")
				for _, step := range protocol.Sequence {
					if step.IsStream() {
						fmt.Fprintf(w, "std::unique_ptr<yardl::hdf5::%s> %s;\n", datasetWriterType(step), stepStreamField(step))
					}
				}
			})
			fmt.Fprint(w, "};\n\n")

			common.WriteComment(w, fmt.Sprintf("HDF5 reader for the %s protocol.", protocol.Name))
			common.WriteComment(w, protocol.Comment)
			readerClassName := Hdf5ReaderClassName(protocol)
			fmt.Fprintf(w, "class %s : public %s, public yardl::hdf5::Hdf5Reader {\n", readerClassName, common.QualifiedAbstractReaderName(protocol))
			w.Indented(func() {
				fmt.Fprintln(w, "public:")
				fmt.Fprintf(w, "%s(std::string path);\n\n", readerClassName)

				for _, step := range protocol.Sequence {
					returnType := "void"
					if step.IsStream() {
						returnType = "bool"
					}
					fmt.Fprintf(w, "%s %s(%s& value) override;\n\n", returnType, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))

					if step.IsStream() {
						if !dsl.ToGeneralizedType(dsl.GetUnderlyingType(step.Type.(*dsl.GeneralizedType).ToScalar())).Cases.IsUnion() {
							// batch optimization is not supported at the moment for unions because
							// we write each type to separate datasets.
							fmt.Fprintf(w, "%s %s(std::vector<%s>& values) override;\n\n", returnType, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
						}
					}
				}

				w.WriteStringln("private:")
				for _, step := range protocol.Sequence {
					if step.IsStream() {
						fmt.Fprintf(w, "std::unique_ptr<yardl::hdf5::%s> %s;\n", datasetReaderType(step), stepStreamField(step))
					}
				}
			})
			fmt.Fprint(w, "};\n\n")
		}
		fmt.Fprintf(w, "} // namespace %s\n\n", common.NamespaceIdentifierName(ns.Name))
	}

	filePath := path.Join(options.SourcesOutputDir, "protocols.h")
	return iocommon.WriteFileIfNeeded(filePath, b.Bytes(), 0644)
}

func needsUnionDataset(s *dsl.ProtocolStep) bool {
	if s.IsStream() {
		gt := dsl.ToGeneralizedType(dsl.GetUnderlyingType(s.Type.(*dsl.GeneralizedType).ToScalar()))
		return gt.Cases.IsUnion()
	}
	return false
}

func getUnionDatasetNestedCount(s *dsl.ProtocolStep) int {
	if !s.IsStream() {
		panic("not a stream")
	}

	gt := dsl.ToGeneralizedType(dsl.GetUnderlyingType(s.Type.(*dsl.GeneralizedType).ToScalar()))
	if gt.Cases[0].IsNullType() {
		return len(gt.Cases) - 1
	}

	return len(gt.Cases)
}

func datasetWriterType(s *dsl.ProtocolStep) string {
	if needsUnionDataset(s) {
		return fmt.Sprintf("UnionDatasetWriter<%d>", getUnionDatasetNestedCount(s))
	}

	return "DatasetWriter"
}

func datasetReaderType(s *dsl.ProtocolStep) string {
	if needsUnionDataset(s) {
		return fmt.Sprintf("UnionDatasetReader<%d>", getUnionDatasetNestedCount(s))
	}

	return "DatasetReader"
}

func stepStreamField(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("%s_dataset_state_", s.Name)
}

func Hdf5WriterClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sWriter", p.Name)
}

func QualifiedHdf5WriterClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::hdf5::%s", common.TypeNamespaceIdentifierName(p), Hdf5WriterClassName(p))
}

func Hdf5ReaderClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sReader", p.Name)
}

func QualifiedHdf5ReaderClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::hdf5::%s", common.TypeNamespaceIdentifierName(p), Hdf5ReaderClassName(p))
}
