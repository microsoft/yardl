// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package hdf5

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

func WriteHdf5(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	options = options.ChangeOutputDir("hdf5")
	if err := os.MkdirAll(options.SourcesOutputDir, 0775); err != nil {
		return err
	}

	err := writeHeaderFile(env, options)
	if err != nil {
		return err
	}

	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`#include "protocols.h"

#include "../yardl/detail/hdf5/io.h"
#include "../yardl/detail/hdf5/ddl.h"
#include "../yardl/detail/hdf5/inner_types.h"
`)

	writeInnerUnionTypes(w, env)

	for _, ns := range env.Namespaces {
		fmt.Fprintf(w, "namespace %s::hdf5 {\n", common.NamespaceIdentifierName(ns.Name))
		writeNamespaceDefinitions(w, ns)
		fmt.Fprintf(w, "} // namespace %s::hdf5", common.NamespaceIdentifierName(ns.Name))
	}

	filePath := path.Join(options.SourcesOutputDir, "protocols.cc")
	return iocommon.EnsureFileContents(filePath, b.Bytes(), 0644)
}

func writeNamespaceDefinitions(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	w.WriteString("namespace {\n")

	for _, t := range ns.TypeDefinitions {
		switch t := t.(type) {
		case *dsl.EnumDefinition:
			writeEnumDdlFunction(w, t)
		}
	}

	for _, t := range ns.TypeDefinitions {
		switch t := t.(type) {
		case *dsl.RecordDefinition:
			writeInnerType(w, t)
		}
	}

	for _, t := range ns.TypeDefinitions {
		switch t := t.(type) {
		case *dsl.RecordDefinition:
			writeRecordDdlFunction(w, t)
		}
	}

	w.WriteString("} // namespace \n\n")

	for _, p := range ns.Protocols {
		writeProtocolMethods(w, p)
	}
}

func getConversionBufferSizeExpression(t dsl.Type) string {
	if containsVlen(t) {
		return fmt.Sprintf("std::max(sizeof(%s), sizeof(%s))", innerTypeSyntax(t), common.TypeSyntax(t))
	}

	return "0"
}

func writeProtocolMethods(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition) {
	writerClassName := Hdf5WriterClassName(p)
	fmt.Fprintf(w, "%s::%s(std::string path)\n", writerClassName, writerClassName)
	w.Indented(func() {
		w.Indented(func() {
			fmt.Fprintf(w, ": yardl::hdf5::Hdf5Writer::Hdf5Writer(path, \"%s\", schema_) {\n", p.Name)
		})
	})

	w.WriteString("}\n\n")

	unionStreams := make([]*dsl.ProtocolStep, 0)
	for _, step := range p.Sequence {
		writeDatasetStateEnsureCreated := func() {}
		if step.IsStream() {
			writeDatasetStateEnsureCreated = func() {
				underlyingStepType := dsl.ToGeneralizedType(dsl.GetUnderlyingType(step.Type.(*dsl.GeneralizedType).ToScalar()))
				isUnion := underlyingStepType.Cases.IsUnion()
				dsWriterName := stepStreamField(step)
				fmt.Fprintf(w, "if (!%s) {\n", dsWriterName)
				w.Indented(func() {
					if isUnion {
						typeLabelPairs := make([]string, 0)
						for _, typeCase := range underlyingStepType.Cases {
							if !typeCase.IsNullType() {
								typeLabelPairs = append(typeLabelPairs, fmt.Sprintf("std::make_tuple(%s, \"%s\", static_cast<size_t>(%s))", typeDdlExpression(typeCase.Type), typeCase.Label, getConversionBufferSizeExpression(underlyingStepType)))
							}
						}
						fmt.Fprintf(w, "%s = std::make_unique<yardl::hdf5::UnionDatasetWriter<%d>>(group_, \"%s\", %t, %s);\n", dsWriterName, getUnionDatasetNestedCount(step), step.Name, underlyingStepType.Cases.HasNullOption(), strings.Join(typeLabelPairs, ", "))
					} else {
						fmt.Fprintf(w, "%s = std::make_unique<yardl::hdf5::DatasetWriter>(group_, \"%s\", %s, %s);\n", dsWriterName, step.Name, typeDdlExpression(underlyingStepType), getConversionBufferSizeExpression(underlyingStepType))
					}
				})
				w.WriteString("}\n\n")
			}
		}

		fmt.Fprintf(w, "void %s::%s(%s const& value) {\n", writerClassName, common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
		w.Indented(func() {
			if step.IsStream() {
				writeDatasetStateEnsureCreated()

				underlyingStepType := dsl.ToGeneralizedType(dsl.GetUnderlyingType(step.Type.(*dsl.GeneralizedType).ToScalar()))
				isUnion := underlyingStepType.Cases.IsUnion()
				dsWriterName := stepStreamField(step)

				if isUnion {
					unionStreams = append(unionStreams, step)
					indexAdjustment := ""
					if underlyingStepType.Cases.HasNullOption() {
						indexAdjustment = " -1"
					}

					w.WriteStringln("std::visit(")
					w.Indented(func() {
						w.WriteStringln("[&](auto const& arg) {")
						w.Indented(func() {
							w.WriteStringln("using T = std::decay_t<decltype(arg)>;")
							for _, typeCase := range underlyingStepType.Cases {
								innerType := innerTypeSyntax(typeCase.Type)
								outerType := common.TypeSyntax(typeCase.Type)
								fmt.Fprintf(w, "if constexpr (std::is_same_v<T, %s>) {\n", outerType)
								w.Indented(func() {
									fmt.Fprintf(w, "%s->Append<%s, %s>(value.index()%s, arg);\n", dsWriterName, innerType, outerType, indexAdjustment)
								})
								w.WriteString("} else ")
							}

							w.WriteStringln("{")
							w.Indented(func() {
								w.WriteStringln("static_assert(yardl::hdf5::always_false_v<T>, \"non-exhaustive visitor!\");")
							})
							w.WriteStringln("}")
						})
						w.WriteStringln("},")
						w.WriteStringln("value);")
					})

				} else {
					fmt.Fprintf(w, "%s->Append<%s, %s>(value);\n", dsWriterName, innerTypeSyntax(step.Type), common.TypeSyntax(step.Type))
				}
			} else {
				fmt.Fprintf(w, "yardl::hdf5::WriteScalarDataset<%s, %s>(group_, \"%s\", %s, value);\n", innerTypeSyntax(step.Type), common.TypeSyntax(step.Type), step.Name, typeDdlExpression(step.Type))
			}
		})
		w.WriteString("}\n\n")

		if step.IsStream() {
			if !dsl.ToGeneralizedType(dsl.GetUnderlyingType(step.Type.(*dsl.GeneralizedType).ToScalar())).Cases.IsUnion() {
				dsWriterName := stepStreamField(step)
				fmt.Fprintf(w, "void %s::%s(std::vector<%s> const& values) {\n", writerClassName, common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
				w.Indented(func() {
					writeDatasetStateEnsureCreated()
					fmt.Fprintf(w, "%s->AppendBatch<%s, %s>(values);\n", dsWriterName, innerTypeSyntax(step.Type), common.TypeSyntax(step.Type))
				})
				w.WriteString("}\n\n")
			}

			fmt.Fprintf(w, "void %s::%s() {\n", writerClassName, common.ProtocolWriteEndImplMethodName(step))
			w.Indented(func() {
				writeDatasetStateEnsureCreated()
				fmt.Fprintf(w, "%s.reset();\n", stepStreamField(step))
			})
			w.WriteString("}\n\n")
		}
	}

	if len(unionStreams) > 0 {
		fmt.Fprintf(w, "void %s::Flush() {\n", writerClassName)
		w.Indented(func() {
			for _, step := range unionStreams {
				field := stepStreamField(step)
				fmt.Fprintf(w, "if (%s) {\n", field)
				w.Indented(func() {
					fmt.Fprintf(w, "%s->Flush();\n", field)
				})
				w.WriteString("}\n")
			}
		})
		w.WriteString("}\n\n")
	}

	readerClassName := Hdf5ReaderClassName(p)
	fmt.Fprintf(w, "%s::%s(std::string path)\n", readerClassName, readerClassName)
	w.Indented(func() {
		w.Indented(func() {
			fmt.Fprintf(w, ": yardl::hdf5::Hdf5Reader::Hdf5Reader(path, \"%s\", schema_) {\n", p.Name)
		})
	})

	w.WriteString("}\n\n")

	for _, step := range p.Sequence {
		returnType := "void"
		if step.IsStream() {
			returnType = "bool"
		}

		fmt.Fprintf(w, "%s %s::%s(%s& value) {\n", returnType, readerClassName, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
		w.Indented(func() {
			if step.IsStream() {
				underlyingStepType := dsl.ToGeneralizedType(dsl.GetUnderlyingType(step.Type.(*dsl.GeneralizedType).ToScalar()))
				isUnion := underlyingStepType.Cases.IsUnion()
				readerFieldName := stepStreamField(step)
				fmt.Fprintf(w, "if (!%s) {\n", readerFieldName)
				w.Indented(func() {
					if isUnion {
						typeLabelPairs := make([]string, 0)
						for _, typeCase := range underlyingStepType.Cases {
							if !typeCase.IsNullType() {
								typeLabelPairs = append(typeLabelPairs, fmt.Sprintf("std::make_tuple(%s, \"%s\", static_cast<size_t>(%s))", typeDdlExpression(typeCase.Type), typeCase.Label, getConversionBufferSizeExpression(underlyingStepType)))
							}
						}
						fmt.Fprintf(w, "%s = std::make_unique<yardl::hdf5::UnionDatasetReader<%d>>(group_, \"%s\", %t, %s);\n", readerFieldName, getUnionDatasetNestedCount(step), step.Name, underlyingStepType.Cases.HasNullOption(), strings.Join(typeLabelPairs, ", "))
					} else {
						fmt.Fprintf(w, "%s = std::make_unique<yardl::hdf5::DatasetReader>(group_, \"%s\", %s, %s);\n", readerFieldName, step.Name, typeDdlExpression(step.Type), getConversionBufferSizeExpression(underlyingStepType))
					}
				})
				w.WriteString("}\n\n")

				if isUnion {
					fmt.Fprintf(w, "auto [has_result, type_index, reader] = %s->ReadIndex();\n", readerFieldName)
					w.WriteStringln("if (!has_result) {")
					w.Indented(func() {
						fmt.Fprintf(w, "%s.reset();\n", readerFieldName)
						w.WriteStringln("return false;")
					})
					w.WriteString("}\n\n")

					w.WriteStringln("switch (type_index) {")
					for i, typeCase := range underlyingStepType.Cases {
						if typeCase.IsNullType() {
							w.WriteStringln("case -1:")
							w.Indented(func() {
								w.WriteStringln("value.emplace<0>();")
								w.WriteStringln("break;")
							})
							continue
						}
						caseValue := i
						if underlyingStepType.Cases.HasNullOption() {
							caseValue--
						}
						fmt.Fprintf(w, "case %d: {\n", caseValue)
						w.Indented(func() {
							innerType := innerTypeSyntax(typeCase.Type)
							outerType := common.TypeSyntax(typeCase.Type)
							fmt.Fprintf(w, "%s& ref = value.emplace<%d>();\n", outerType, i)
							fmt.Fprintf(w, "reader->Read<%s, %s>(ref);\n", innerType, outerType)
							w.WriteStringln("break;")
						})
						w.WriteStringln("}")
					}
					w.WriteString("}\n\n")
					w.WriteStringln("return true;")

				} else {
					fmt.Fprintf(w, "bool has_value = %s->Read<%s, %s>(value);\n", readerFieldName, innerTypeSyntax(step.Type), common.TypeSyntax(step.Type))
					w.WriteString("if (!has_value) {\n")
					w.Indented(func() {
						fmt.Fprintf(w, "%s.reset();\n", readerFieldName)
					})
					w.WriteString("}\n\n")
					w.WriteString("return has_value;\n")
				}
			} else {
				fmt.Fprintf(w, "yardl::hdf5::ReadScalarDataset<%s, %s>(group_, \"%s\", %s, value);\n", innerTypeSyntax(step.Type), common.TypeSyntax(step.Type), step.Name, typeDdlExpression(step.Type))
			}
		})
		w.WriteString("}\n\n")

		if step.IsStream() && !dsl.ToGeneralizedType(dsl.GetUnderlyingType(step.Type.(*dsl.GeneralizedType).ToScalar())).Cases.IsUnion() {
			fmt.Fprintf(w, "bool %s::%s(std::vector<%s>& values) {\n", readerClassName, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
			w.Indented(func() {
				readerFieldName := stepStreamField(step)
				fmt.Fprintf(w, "if (!%s) {\n", readerFieldName)
				w.Indented(func() {
					fmt.Fprintf(w, "%s = std::make_unique<yardl::hdf5::DatasetReader>(group_, \"%s\", %s);\n", readerFieldName, step.Name, typeDdlExpression(step.Type))
				})
				w.WriteString("}\n\n")
				fmt.Fprintf(w, "bool has_more = %s->ReadBatch<%s, %s>(values);\n", readerFieldName, innerTypeSyntax(step.Type), common.TypeSyntax(step.Type))
				w.WriteString("if (!has_more) {\n")
				w.Indented(func() {
					fmt.Fprintf(w, "%s.reset();\n", readerFieldName)
				})
				w.WriteString("}\n\n")
				w.WriteString("return has_more;\n")
			})
			w.WriteString("}\n\n")
		}
	}
}
