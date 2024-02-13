// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package binary

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

func WriteBinary(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	options = options.ChangeOutputDir("binary")
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

#include <cstddef>

#include "../yardl/detail/binary/coded_stream.h"
#include "../yardl/detail/binary/serializers.h"
`)
	writeIsTriviallySerializableSpecializations(w, env)
	writeUnionSerializers(w, env)
	for _, ns := range env.Namespaces {
		fmt.Fprintf(w, "namespace %s::binary {\n", common.NamespaceIdentifierName(ns.Name))
		writeNamespaceDefinitions(w, ns)
		fmt.Fprintf(w, "} // namespace %s::binary\n\n", common.NamespaceIdentifierName(ns.Name))
	}

	filePath := path.Join(options.SourcesOutputDir, "protocols.cc")
	return iocommon.WriteFileIfNeeded(filePath, b.Bytes(), 0644)

}

func writeHeaderFile(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`#pragma once
#include <array>
#include <complex>
#include <memory>
#include <optional>
#include <variant>
#include <vector>

#include "../yardl/detail/binary/reader_writer.h"
#include "../protocols.h"
#include "../types.h"
`)

	for _, ns := range env.Namespaces {
		if !ns.IsTopLevel {
			continue
		}
		fmt.Fprintf(w, "namespace %s::binary {\n", common.NamespaceIdentifierName(ns.Name))
		for _, protocol := range ns.Protocols {
			common.WriteComment(w, fmt.Sprintf("Binary writer for the %s protocol.", protocol.Name))
			common.WriteComment(w, protocol.Comment)
			writerClassName := BinaryWriterClassName(protocol)
			fmt.Fprintf(w, "class %s : public %s, yardl::binary::BinaryWriter {\n", writerClassName, common.QualifiedAbstractWriterName(protocol))
			w.Indented(func() {
				w.WriteStringln("public:")
				fmt.Fprintf(w, "%s(std::ostream& stream, Version version = Version::Current)\n", writerClassName)
				w.Indented(func() {
					w.Indented(func() {
						fmt.Fprintf(w, ": yardl::binary::BinaryWriter(stream, %s::SchemaFromVersion(version)), version_(version) {", common.QualifiedAbstractWriterName(protocol))
					})
				})
				w.WriteStringln("}\n")

				fmt.Fprintf(w, "%s(std::string file_name, Version version = Version::Current)\n", writerClassName)
				w.Indented(func() {
					w.Indented(func() {
						fmt.Fprintf(w, ": yardl::binary::BinaryWriter(file_name, %s::SchemaFromVersion(version)), version_(version) {", common.QualifiedAbstractWriterName(protocol))
					})
				})
				w.WriteStringln("}\n")

				w.WriteString("void Flush() override;\n\n")

				w.WriteStringln("protected:")
				for _, step := range protocol.Sequence {
					endMethodName := common.ProtocolWriteEndImplMethodName(step)
					common.WriteComment(w, step.Comment)

					fmt.Fprintf(w, "void %s(%s const& value) override;\n", common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))

					if step.IsStream() {
						fmt.Fprintf(w, "void %s(std::vector<%s> const& values) override;\n", common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
						fmt.Fprintf(w, "void %s() override;\n", endMethodName)
					}
				}

				w.WriteString("void CloseImpl() override;\n")
				w.WriteStringln("")
				w.WriteStringln("Version version_;")
			})
			fmt.Fprint(w, "};\n\n")

			common.WriteComment(w, fmt.Sprintf("Binary reader for the %s protocol.", protocol.Name))
			common.WriteComment(w, protocol.Comment)
			readerClassName := BinaryReaderClassName(protocol)
			fmt.Fprintf(w, "class %s : public %s, yardl::binary::BinaryReader {\n", readerClassName, common.QualifiedAbstractReaderName(protocol))
			w.Indented(func() {
				fmt.Fprintln(w, "public:")
				fmt.Fprintf(w, "%s(std::istream& stream)\n", readerClassName)
				w.Indented(func() {
					w.Indented(func() {
						fmt.Fprintf(w, ": yardl::binary::BinaryReader(stream), version_(%s::VersionFromSchema(schema_read_)) {", common.QualifiedAbstractReaderName(protocol))
					})
				})
				w.WriteStringln("}\n")

				fmt.Fprintf(w, "%s(std::string file_name)\n", readerClassName)
				w.Indented(func() {
					w.Indented(func() {
						fmt.Fprintf(w, ": yardl::binary::BinaryReader(file_name), version_(%s::VersionFromSchema(schema_read_)) {", common.QualifiedAbstractReaderName(protocol))
					})
				})
				w.WriteStringln("}\n")

				fmt.Fprintf(w, "Version GetVersion() { return version_; }\n\n")

				w.WriteStringln("protected:")
				hasStream := false
				for _, step := range protocol.Sequence {
					if step.IsStream() {
						hasStream = true
					}

					returnType := "void"
					if step.IsStream() {
						returnType = "bool"
					}
					fmt.Fprintf(w, "%s %s(%s& value) override;\n", returnType, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
					if step.IsStream() {
						fmt.Fprintf(w, "bool %s(std::vector<%s>& values) override;\n", common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
					}
				}

				w.WriteString("void CloseImpl() override;\n")
				w.WriteStringln("")
				w.WriteStringln("Version version_;")

				if hasStream {
					w.WriteStringln("\nprivate:")
					w.WriteStringln("size_t current_block_remaining_ = 0;")
				}
			})
			fmt.Fprint(w, "};\n\n")
		}
		fmt.Fprintf(w, "} // namespace %s::binary\n", common.NamespaceIdentifierName(ns.Name))
	}

	filePath := path.Join(options.SourcesOutputDir, "protocols.h")
	return iocommon.WriteFileIfNeeded(filePath, b.Bytes(), 0644)
}

func writeIsTriviallySerializableSpecializations(w *formatting.IndentedWriter, env *dsl.Environment) {
	fmt.Fprintf(w, "namespace yardl::binary {\n")
	w.WriteStringln("#ifndef _MSC_VER")
	common.WriteComment(w, "Values of offsetof() are only used if types are standard-layout.")
	w.WriteStringln("#pragma GCC diagnostic push")
	w.WriteStringln("#pragma GCC diagnostic ignored \"-Winvalid-offsetof\"")
	w.WriteStringln("#endif\n")

	for _, ns := range env.Namespaces {
		for _, td := range ns.TypeDefinitions {
			writeIsTriviallySerializableSpecialization(w, td)
		}
	}
	w.WriteStringln("#ifndef _MSC_VER")
	w.WriteStringln("#pragma GCC diagnostic pop // #pragma GCC diagnostic ignored \"-Winvalid-offsetof\" ")
	w.WriteStringln("#endif")
	fmt.Fprintf(w, "} //namespace yardl::binary \n\n")
}

func writeIsTriviallySerializableSpecialization(w *formatting.IndentedWriter, t dsl.TypeDefinition) {
	switch t.(type) {
	case *dsl.RecordDefinition:
		break
	default:
		return
	}

	meta := t.GetDefinitionMeta()

	w.WriteString("template <")
	formatting.Delimited(w, ", ", meta.TypeParameters, func(w *formatting.IndentedWriter, i int, item *dsl.GenericTypeParameter) {
		fmt.Fprintf(w, "typename %s", common.TypeDefinitionSyntax(item))
	})
	w.WriteStringln(">")
	fmt.Fprintf(w, "struct IsTriviallySerializable<%s> {\n", common.TypeDefinitionSyntax(t))

	w.Indented(func() {
		fmt.Fprintf(w, "using __T__ = %s;\n", common.TypeDefinitionSyntax(t))
		w.WriteString("static constexpr bool value = \n")
		w.Indented(func() {
			w.WriteStringln("std::is_standard_layout_v<__T__> &&")

			switch t := t.(type) {
			case *dsl.RecordDefinition:
				formatting.Delimited(w, " &&\n", t.Fields, func(w *formatting.IndentedWriter, i int, f *dsl.Field) {
					fmt.Fprintf(w, "IsTriviallySerializable<decltype(__T__::%s)>::value", common.FieldIdentifierName(f.Name))
				})

				if len(t.Fields) > 0 {

					w.WriteStringln(" &&")

					fmt.Fprintf(w, "(sizeof(__T__) == (")
					formatting.Delimited(w, " + ", t.Fields, func(w *formatting.IndentedWriter, i int, f *dsl.Field) {
						fmt.Fprintf(w, "sizeof(__T__::%s)", common.FieldIdentifierName(f.Name))
					})

					w.WriteString("))")

					if len(t.Fields) > 1 {
						w.WriteStringln(" &&")
					}

					for i, f := range t.Fields {
						if i > 0 {
							if i > 1 {
								w.WriteString(" && ")
							}
							fmt.Fprintf(w, "offsetof(__T__, %s) < offsetof(__T__, %s)", common.FieldIdentifierName(common.FieldIdentifierName(t.Fields[i-1].Name)), common.FieldIdentifierName(f.Name))
						}
					}
				}
			}
		})

		w.WriteStringln(";")
	})
	fmt.Fprintf(w, "};\n\n")
}

func collectUnionArities(env *dsl.Environment) []int {
	arities := make(map[int]any)

	dsl.Visit(env, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case nil:
			return
		case *dsl.ProtocolDefinition:
			for _, change := range t.Versions {
				if change == nil {
					continue
				}

				for _, tc := range change.StepChanges {
					if tc == nil {
						continue
					}
					self.Visit(tc.OldType())
				}
			}
			self.VisitChildren(node)

		case *dsl.Namespace:
			for _, v := range t.Versions {
				for _, change := range t.DefinitionChanges[v] {
					if change == nil {
						continue
					}

					switch change := change.(type) {
					case *dsl.RecordChange:
						for _, tc := range change.FieldChanges {
							if tc == nil {
								continue
							}
							self.Visit(tc.OldType())
						}
					case *dsl.NamedTypeChange:
						if tc := change.TypeChange; tc != nil {
							self.Visit(tc.OldType())
						}
					}
				}
			}
			self.VisitChildren(node)

		case *dsl.GeneralizedType:
			if t.Cases.IsUnion() {
				arities[len(t.Cases)] = nil
			}

			self.VisitChildren(node)
		default:
			self.VisitChildren(node)
		}

	})

	keys := make([]int, 0, len(arities))
	for k := range arities {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}

func writeUnionSerializers(w *formatting.IndentedWriter, env *dsl.Environment) {
	arites := collectUnionArities(env)
	if len(arites) == 0 {
		return
	}

	w.WriteString("namespace {\n")
	formatting.Delimited(w, "\n", arites, func(w *formatting.IndentedWriter, i int, arity int) {
		elements := make([]int, arity)
		for i := 0; i < arity; i++ {
			elements[i] = i
		}

		writeTemplateParameters := make([]string, 2*arity)
		for i := range elements {
			writeTemplateParameters[2*i] = fmt.Sprintf("typename T%d", i)
			writeTemplateParameters[2*i+1] = fmt.Sprintf("yardl::binary::Writer<T%d> WriteT%d", i, i)
		}

		readTemplateParameters := make([]string, 2*arity)
		for i := range elements {
			readTemplateParameters[2*i] = fmt.Sprintf("typename T%d", i)
			readTemplateParameters[2*i+1] = fmt.Sprintf("yardl::binary::Reader<T%d> ReadT%d", i, i)
		}

		variantTemplateArguments := make([]string, arity)
		for i := range elements {
			variantTemplateArguments[i] = fmt.Sprintf("T%d", i)
		}

		fmt.Fprintf(w, "template<%s>\n", strings.Join(writeTemplateParameters, ", "))
		fmt.Fprintf(w, "void WriteUnion(yardl::binary::CodedOutputStream& stream, std::variant<%s> const& value) {\n", strings.Join(variantTemplateArguments, ", "))
		w.Indented(func() {
			w.WriteStringln("yardl::binary::WriteInteger(stream, value.index());")
			fmt.Fprintf(w, "switch (value.index()) {\n")
			for i := range elements {
				fmt.Fprintf(w, "case %d: {\n", i)
				w.Indented(func() {
					fmt.Fprintf(w, "T%d const& v = std::get<%d>(value);\n", i, i)
					fmt.Fprintf(w, "WriteT%d(stream, v);\n", i)
					w.WriteStringln("break;")
				})
				w.WriteStringln("}")
			}
			fmt.Fprintf(w, "default: throw std::runtime_error(\"Invalid union index.\");\n")
			fmt.Fprintf(w, "}\n")
		})
		fmt.Fprintf(w, "}\n\n")

		fmt.Fprintf(w, "template<%s>\n", strings.Join(readTemplateParameters, ", "))
		fmt.Fprintf(w, "void ReadUnion(yardl::binary::CodedInputStream& stream, std::variant<%s>& value) {\n", strings.Join(variantTemplateArguments, ", "))
		w.Indented(func() {
			w.WriteStringln("size_t index;")
			w.WriteStringln("yardl::binary::ReadInteger(stream, index);")
			w.WriteStringln("switch (index) {")
			w.Indented(func() {
				for i := range elements {
					fmt.Fprintf(w, "case %d: {\n", i)
					w.Indented(func() {
						fmt.Fprintf(w, "T%d v;\n", i)
						fmt.Fprintf(w, "ReadT%d(stream, v);\n", i)
						fmt.Fprintf(w, "value = std::move(v);\n")
						w.WriteStringln("break;")
					})
					w.WriteStringln("}")
				}
				fmt.Fprintf(w, "default: throw std::runtime_error(\"Invalid union index.\");\n")
			})
			w.WriteStringln("}")
		})
		fmt.Fprintf(w, "}\n")

	})

	w.WriteString("} // namespace\n\n")
}

func writeNamespaceDefinitions(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	if len(ns.TypeDefinitions) > 0 {
		w.WriteStringln("namespace {")
		for _, typeDef := range ns.TypeDefinitions {
			writeSerializers(w, typeDef)
		}

		for _, versionLabel := range ns.Versions {
			for _, change := range ns.DefinitionChanges[versionLabel] {
				writeCompatibilitySerializers(w, change, versionLabel)
			}
		}
		w.WriteString("} // namespace\n\n")
	}

	if ns.IsTopLevel {
		for _, protocol := range ns.Protocols {
			writeProtocolMethods(w, protocol)
		}
	}
}

func writeSerializers(w *formatting.IndentedWriter, t dsl.TypeDefinition) {
	switch t.(type) {
	case *dsl.EnumDefinition:
		return
	}

	writeFallbackBody := func(write bool) {
		switch t := t.(type) {
		case *dsl.RecordDefinition:
			for _, field := range t.Fields {
				fmt.Fprintf(w, "%s(stream, value.%s);\n", typeRwFunction(field.Type, write), common.FieldIdentifierName(field.Name))
			}
		case *dsl.NamedType:
			fmt.Fprintf(w, "%s(stream, value);\n", typeRwFunction(t.Type, write))
		default:
			panic(fmt.Sprintf("Unexpected type %T", t))
		}
	}

	writeRwFunctionSignature(t, w, true)
	w.Indented(func() {
		fmt.Fprintf(w, "if constexpr (yardl::binary::IsTriviallySerializable<%s>::value) {\n", common.TypeDefinitionSyntax(t))
		w.Indented(func() {
			fmt.Fprintf(w, "yardl::binary::WriteTriviallySerializable(stream, value);\n")
			w.WriteStringln("return;")
		})
		w.WriteStringln("}\n")
		writeFallbackBody(true)
	})
	w.WriteString("}\n\n")

	writeRwFunctionSignature(t, w, false)
	w.Indented(func() {
		fmt.Fprintf(w, "if constexpr (yardl::binary::IsTriviallySerializable<%s>::value) {\n", common.TypeDefinitionSyntax(t))
		w.Indented(func() {
			fmt.Fprintf(w, "yardl::binary::ReadTriviallySerializable(stream, value);\n")
			w.WriteStringln("return;")
		})
		w.WriteStringln("}\n")
		writeFallbackBody(false)
	})
	w.WriteString("}\n\n")
}

func verb(write bool) string {
	if write {
		return "Write"
	}
	return "Read"
}

func writeRwFunctionSignature(t dsl.TypeDefinition, w *formatting.IndentedWriter, write bool) {
	writeRwFunctionTemplateDeclaration(t, w, write)
	if write {
		fmt.Fprintf(w, "[[maybe_unused]] void Write%s(yardl::binary::CodedOutputStream& stream, %s const& value) {\n", t.GetDefinitionMeta().Name, common.TypeDefinitionSyntax(t))
	} else {
		fmt.Fprintf(w, "[[maybe_unused]] void Read%s(yardl::binary::CodedInputStream& stream, %s& value) {\n", t.GetDefinitionMeta().Name, common.TypeDefinitionSyntax(t))
	}
}

func writeRwFunctionTemplateDeclaration(t dsl.TypeDefinition, w *formatting.IndentedWriter, write bool) {
	meta := t.GetDefinitionMeta()
	if len(meta.TypeParameters) > 0 {
		functionPointerName := "Reader"
		if write {
			functionPointerName = "Writer"
		}
		templateParameters := make([]string, 2*len(meta.TypeParameters))
		for i, p := range meta.TypeParameters {
			templateParameters[2*i] = "typename " + common.TypeDefinitionSyntax(p)
			templateParameters[2*i+1] = fmt.Sprintf("yardl::binary::%s<%s> %s%s", functionPointerName, common.TypeDefinitionSyntax(p), verb(write), common.TypeDefinitionSyntax(p))
		}
		fmt.Fprintf(w, "template<%s>\n", strings.Join(templateParameters, ", "))
	}
}

// Writes code needed to convert sourceName: typeChange.Old() into targetName: typeChange.New() on read, or vice-versa on write
func writeTypeConversion(w *formatting.IndentedWriter, typeChange dsl.TypeChange, sourceName, targetName string, write bool) {
	switch tc := typeChange.(type) {
	case *dsl.TypeChangeNumberToNumber:
		if write {
			writeTypeConversion(w, tc.Inverse(), sourceName, targetName, !write)
		} else {
			oldPrim := tc.OldType().(*dsl.SimpleType).ResolvedDefinition.(dsl.PrimitiveDefinition)
			newPrim := tc.NewType().(*dsl.SimpleType).ResolvedDefinition.(dsl.PrimitiveDefinition)
			rhs := sourceName
			if dsl.GetPrimitiveKind(oldPrim) == dsl.PrimitiveKindFloatingPoint && dsl.GetPrimitiveKind(newPrim) == dsl.PrimitiveKindInteger {
				rhs = fmt.Sprintf("std::round(%s)", rhs)
			}

			fmt.Fprintf(w, "%s = static_cast<%s>(%s);\n", targetName, common.TypeSyntax(tc.NewType()), rhs)
		}

	case *dsl.TypeChangeStringToNumber:
		writeTypeConversion(w, tc.Inverse(), sourceName, targetName, !write)
	case *dsl.TypeChangeNumberToString:
		if write {
			rhs := sourceName
			def := tc.OldType().(*dsl.SimpleType).ResolvedDefinition
			switch def.(dsl.PrimitiveDefinition) {
			case dsl.PrimitiveInt8, dsl.PrimitiveInt16, dsl.PrimitiveInt32:
				rhs = fmt.Sprintf("std::stoi(%s)", sourceName)
			case dsl.PrimitiveInt64:
				rhs = fmt.Sprintf("std::stol(%s)", sourceName)
			case dsl.PrimitiveUint8, dsl.PrimitiveUint16, dsl.PrimitiveUint32, dsl.PrimitiveUint64:
				rhs = fmt.Sprintf("std::stoul(%s)", sourceName)
			case dsl.PrimitiveFloat32:
				rhs = fmt.Sprintf("std::stof(%s)", sourceName)
			case dsl.PrimitiveFloat64:
				rhs = fmt.Sprintf("std::stod(%s)", sourceName)
			default:
				rhs = sourceName
			}

			fmt.Fprintf(w, "try {\n")
			w.Indented(func() {
				fmt.Fprintf(w, "%s = %s;\n", targetName, rhs)
			})
			fmt.Fprintf(w, "} catch (...) {\n")
			w.Indented(func() {
				errMsg := fmt.Sprintf(`Unable to convert string \"" + %s + "\" to number`, sourceName)
				fmt.Fprintf(w, "throw new std::runtime_error(\"%s\");\n", errMsg)
			})
			fmt.Fprintf(w, "}\n")

		} else {
			fmt.Fprintf(w, "%s = std::to_string(%s);\n", targetName, sourceName)
		}

	case *dsl.TypeChangeOptionalTypeChanged:
		fmt.Fprintf(w, "if (%s.has_value()) {\n", sourceName)
		w.Indented(func() {
			writeTypeConversion(w, tc.InnerChange, sourceName+".value()", targetName, write)
		})
		fmt.Fprintf(w, "}\n")

	case *dsl.TypeChangeOptionalToScalar:
		writeTypeConversion(w, tc.Inverse(), sourceName, targetName, !write)
	case *dsl.TypeChangeScalarToOptional:
		if write {
			fmt.Fprintf(w, "if (%s.has_value()) {\n", sourceName)
			w.Indented(func() {
				fmt.Fprintf(w, "%s = %s.value();\n", targetName, sourceName)
			})
			fmt.Fprintf(w, "}\n")
		} else {
			fmt.Fprintf(w, "%s = %s;\n", targetName, sourceName)
		}

	case *dsl.TypeChangeUnionToScalar:
		writeTypeConversion(w, tc.Inverse(), sourceName, targetName, !write)
	case *dsl.TypeChangeScalarToUnion:
		if write {
			// Writing a Union as a Scalar
			fmt.Fprintf(w, "if (%s.index() == %d) {\n", sourceName, tc.TypeIndex)
			w.Indented(func() {
				fmt.Fprintf(w, "%s = std::get<%d>(%s);\n", targetName, tc.TypeIndex, sourceName)
			})
			fmt.Fprintf(w, "}\n")
		} else {
			// Reading a Scalar into a Union
			fmt.Fprintf(w, "%s = %s;\n", targetName, sourceName)
		}

	case *dsl.TypeChangeUnionToOptional:
		writeTypeConversion(w, tc.Inverse(), sourceName, targetName, !write)
	case *dsl.TypeChangeOptionalToUnion:
		if write {
			// Writing a Union as an Optional
			fmt.Fprintf(w, "if (%s.index() == %d) {\n", sourceName, tc.TypeIndex)
			w.Indented(func() {
				fmt.Fprintf(w, "%s = std::get<%d>(%s);\n", targetName, tc.TypeIndex, sourceName)
			})
			fmt.Fprintf(w, "}\n")
		} else {
			// Reading an Optional into a Union
			fmt.Fprintf(w, "if (%s.has_value()) {\n", sourceName)
			w.Indented(func() {
				fmt.Fprintf(w, "%s = %s.value();\n", targetName, sourceName)
			})
			fmt.Fprintf(w, "} else {\n")
			w.Indented(func() {
				fmt.Fprintf(w, "%s = std::monostate{};\n", targetName)
			})
			fmt.Fprintf(w, "}\n")
		}

	case *dsl.TypeChangeUnionTypesetChanged:
		if write {
			writeTypeConversion(w, typeChange.Inverse(), sourceName, targetName, !write)
			return
		}

		fmt.Fprintf(w, "switch (%s.index()) {\n", sourceName)
		w.Indented(func() {
			for i := range tc.OldType().(*dsl.GeneralizedType).Cases {
				fmt.Fprintf(w, "case %d: {\n", i)
				w.Indented(func() {
					if tc.OldMatches[i] {
						fmt.Fprintf(w, "%s = std::get<%d>(%s);\n", targetName, i, sourceName)
					} else {
						fmt.Fprintf(w, "throw new std::runtime_error(\"Union type incompatible with previous version of model\");\n")
					}
					fmt.Fprintf(w, "break;\n")

				})
				fmt.Fprintf(w, "}\n")
			}
			fmt.Fprintf(w, "default: throw new std::runtime_error(\"Invalid union index.\");\n")
		})
		fmt.Fprintf(w, "}\n")

	default:
		panic("Expected a TypeChange")
	}
}

// If a TypeChange is the result of an underlying TypeDefinition change, we don't need to perform
// explicit conversion - we only need to call the corresponding "compatibility" serializer functions
func requiresExplicitConversion(tc dsl.TypeChange) bool {
	switch tc := tc.(type) {
	case *dsl.TypeChangeDefinitionChanged:
		return false
	case *dsl.TypeChangeOptionalTypeChanged:
		return requiresExplicitConversion(tc.InnerChange)
	case *dsl.TypeChangeStreamTypeChanged:
		return requiresExplicitConversion(tc.InnerChange)
	case *dsl.TypeChangeVectorTypeChanged:
		return requiresExplicitConversion(tc.InnerChange)
	}
	return true
}

func writeCompatibilitySerializers(w *formatting.IndentedWriter, change dsl.DefinitionChange, versionLabel string) {
	switch change.LatestDefinition().(type) {
	case *dsl.EnumDefinition:
		return
	}

	writeFallbackBody := func(write bool) {
		switch change := change.(type) {
		case *dsl.RecordChange:
			p := change.PreviousDefinition().(*dsl.RecordDefinition)
			for i, field := range p.Fields {
				tmpVarName := common.FieldIdentifierName(field.Name)
				if change.FieldRemoved[i] {
					// Field was removed: Read it and discard, or Write "default" value
					tmpVarType := common.TypeSyntax(field.Type)
					fmt.Fprintf(w, "%s %s = {};\n", tmpVarType, tmpVarName)
					fmt.Fprintf(w, "%s(stream, %s);\n", typeRwFunction(field.Type, write), tmpVarName)
				} else if tc := change.FieldChanges[i]; tc != nil {
					// Field type change: Handle type conversions
					if requiresExplicitConversion(tc) {
						tmpVarType := common.TypeSyntax(tc.OldType())
						fmt.Fprintf(w, "%s %s = {};\n", tmpVarType, tmpVarName)

						if write {
							writeTypeConversion(w, tc, fmt.Sprintf("value.%s", tmpVarName), tmpVarName, write)
							fmt.Fprintf(w, "%s(stream, %s);\n", typeRwFunction(tc.OldType(), write), tmpVarName)
						} else {
							fmt.Fprintf(w, "%s(stream, %s);\n", typeRwFunction(tc.OldType(), write), tmpVarName)
							writeTypeConversion(w, tc, tmpVarName, fmt.Sprintf("value.%s", tmpVarName), write)
						}
					} else {
						fmt.Fprintf(w, "%s(stream, value.%s);\n", typeRwFunction(tc.OldType(), write), tmpVarName)
					}
				} else {
					fmt.Fprintf(w, "%s(stream, value.%s);\n", typeRwFunction(field.Type, write), tmpVarName)
				}
			}
		case *dsl.NamedTypeChange:
			switch prev := change.PreviousDefinition().(type) {
			case *dsl.NamedType:
				// prev := change.PreviousDefinition().(*dsl.NamedType)
				if tc := change.TypeChange; tc != nil {
					tmpVarName := common.FieldIdentifierName(prev.Name)
					if requiresExplicitConversion(tc) {
						varType := common.TypeSyntax(tc.OldType())
						fmt.Fprintf(w, "%s %s = {};\n", varType, tmpVarName)
						if write {
							writeTypeConversion(w, tc, "value", tmpVarName, write)
							fmt.Fprintf(w, "%s(stream, %s);\n", typeRwFunction(tc.OldType(), write), tmpVarName)
						} else {
							fmt.Fprintf(w, "%s(stream, %s);\n", typeRwFunction(tc.OldType(), write), tmpVarName)
							writeTypeConversion(w, tc, tmpVarName, "value", write)
						}
					} else {
						fmt.Fprintf(w, "%s(stream, value);\n", typeRwFunction(tc.OldType(), write))
					}
				} else {
					fmt.Fprintf(w, "%s(stream, value);\n", typeRwFunction(prev.Type, write))
				}

			case dsl.TypeDefinition:
				fmt.Fprintf(w, "%s(stream, value);\n", typeDefinitionRwFunction(prev, write))
			}

		case *dsl.CompatibilityChange:
			fmt.Fprintf(w, "%s(stream, value);\n", typeDefinitionRwFunction(change.LatestDefinition(), write))

		default:
			panic(fmt.Sprintf("Unexpected type %T", change.PreviousDefinition()))
		}
	}

	writeCompatibilityRwFunctionSignature(change.PreviousDefinition(), change.LatestDefinition(), w, true)
	w.Indented(func() {
		writeFallbackBody(true)
	})
	w.WriteString("}\n\n")

	writeCompatibilityRwFunctionSignature(change.PreviousDefinition(), change.LatestDefinition(), w, false)
	w.Indented(func() {
		writeFallbackBody(false)
	})
	w.WriteString("}\n\n")
}

func writeCompatibilityRwFunctionSignature(old dsl.TypeDefinition, new dsl.TypeDefinition, w *formatting.IndentedWriter, write bool) {
	writeRwFunctionTemplateDeclaration(old, w, write)
	if write {
		fmt.Fprintf(w, "[[maybe_unused]] void Write%s(yardl::binary::CodedOutputStream& stream, %s const& value) {\n", old.GetDefinitionMeta().Name, common.TypeDefinitionSyntax(new))
	} else {
		fmt.Fprintf(w, "[[maybe_unused]] void Read%s(yardl::binary::CodedInputStream& stream, %s& value) {\n", old.GetDefinitionMeta().Name, common.TypeDefinitionSyntax(new))
	}
}

func writeProtocolMethods(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition) {

	extractStepChanges := func(stepIndex int) map[string]dsl.TypeChange {
		stepChanges := make(map[string]dsl.TypeChange)
		for label, protocolChange := range p.Versions {
			if protocolChange == nil {
				continue
			}
			stepChanges[label] = protocolChange.StepChanges[stepIndex]
		}
		return stepChanges
	}

	writerClassName := BinaryWriterClassName(p)
	for i, step := range p.Sequence {
		stepChanges := extractStepChanges(i)

		fmt.Fprintf(w, "void %s::%s(%s const& value) {\n", writerClassName, common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
		w.Indented(func() {
			writeProtocolStep(w, step, stepChanges, false, true)
		})
		w.WriteString("}\n\n")

		if step.IsStream() {
			fmt.Fprintf(w, "void %s::%s(std::vector<%s> const& values) {\n", writerClassName, common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
			w.Indented(func() {
				w.WriteStringln("if (!values.empty()) {")
				w.Indented(func() {
					writeProtocolStep(w, step, stepChanges, true, true)
				})
				w.WriteStringln("}")
			})
			w.WriteString("}\n\n")

			fmt.Fprintf(w, "void %s::%s() {\n", writerClassName, common.ProtocolWriteEndImplMethodName(step))
			w.Indented(func() {
				writeEndStream(w, stepChanges)
			})
			w.WriteString("}\n\n")
		}
	}

	fmt.Fprintf(w, "void %s::Flush() {\n", writerClassName)
	w.Indented(func() {
		w.WriteString("stream_.Flush();\n")
	})
	w.WriteString("}\n\n")

	fmt.Fprintf(w, "void %s::CloseImpl() {\n", writerClassName)
	w.Indented(func() {
		w.WriteString("stream_.Flush();\n")
	})
	w.WriteString("}\n\n")

	readerClassName := BinaryReaderClassName(p)
	for i, step := range p.Sequence {
		stepChanges := extractStepChanges(i)

		returnType := "void"
		if step.IsStream() {
			returnType = "bool"
		}

		fmt.Fprintf(w, "%s %s::%s(%s& value) {\n", returnType, readerClassName, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
		w.Indented(func() {
			writeProtocolStep(w, step, stepChanges, false, false)
		})
		w.WriteString("}\n\n")

		if step.IsStream() {
			fmt.Fprintf(w, "%s %s::%s(std::vector<%s>& values) {\n", returnType, readerClassName, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
			w.Indented(func() {
				writeProtocolStep(w, step, stepChanges, true, false)
				w.WriteStringln("return current_block_remaining_ != 0;")
			})
			w.WriteString("}\n\n")
		}
	}

	fmt.Fprintf(w, "void %s::CloseImpl() {\n", readerClassName)
	w.Indented(func() {
		w.WriteString("stream_.VerifyFinished();\n")
	})
	w.WriteString("}\n\n")
}

// Writes a switch statement with cases for each version of a ProtocolStep
// Only if there are version changes or the ProtocolStep was added
func writeChangeSwitchCase(w *formatting.IndentedWriter, changes map[string]dsl.TypeChange,
	writeDefault, writeAdded func(*formatting.IndentedWriter),
	writeConversion func(*formatting.IndentedWriter, dsl.TypeChange)) {

	allNil := func(vs map[string]dsl.TypeChange) bool {
		for _, v := range vs {
			if v != nil {
				return false
			}
		}
		return true
	}

	if allNil(changes) {
		// Skip switch statement - this ProtocolStep has never changed
		writeDefault(w)
		return
	}

	// Sort the case labels for deterministic ordering
	var versionLabels []string
	for versionLabel := range changes {
		versionLabels = append(versionLabels, versionLabel)
	}
	sort.Strings(versionLabels)

	fmt.Fprintf(w, "switch (version_) {\n")
	for _, versionLabel := range versionLabels {
		change := changes[versionLabel]
		if change == nil {
			continue
		}

		_, stepAdded := change.(*dsl.TypeChangeStepAdded)
		if writeConversion == nil && !stepAdded {
			continue
		}

		fmt.Fprintf(w, "case Version::%s: {\n", versionLabel)
		w.Indented(func() {
			defer func() {
				w.WriteStringln("break;")
			}()

			if stepAdded {
				if writeAdded != nil {
					writeAdded(w)
				}
			} else {
				if writeConversion != nil {
					writeConversion(w, change)
				}
			}
		})
		w.WriteStringln("}")
	}
	fmt.Fprintln(w, "default:")
	w.Indented(func() {
		writeDefault(w)
		fmt.Fprintln(w, "break;")
	})
	fmt.Fprintln(w, "}")
}

func writeStepRw(w *formatting.IndentedWriter, stepType dsl.Type, target string, isStream, isPlural, write bool) {
	if !isStream {
		fmt.Fprintf(w, "%s(stream_, %s);\n", typeRwFunction(stepType, write), target)
		return
	}

	if write {
		if isPlural {
			vectorType := *stepType.(*dsl.GeneralizedType)
			vectorType.Dimensionality = &dsl.Vector{}
			stepType = &vectorType
			fmt.Fprintf(w, "%s(stream_, values);\n", typeRwFunction(stepType, write))
		} else {
			fmt.Fprintf(w, "yardl::binary::WriteBlock<%s, %s>(stream_, value);\n", common.TypeSyntax(stepType), typeRwFunction(stepType, write))
		}
	} else {
		if isPlural {
			stepType = stepType.(*dsl.GeneralizedType).ToScalar()
			fmt.Fprintf(w, "yardl::binary::ReadBlocksIntoVector<%s, %s>(stream_, current_block_remaining_, values);\n", common.TypeSyntax(stepType), typeRwFunction(stepType, write))
		} else {
			fmt.Fprintf(w, "return yardl::binary::ReadBlock<%s, %s>(stream_, current_block_remaining_, value);\n", common.TypeSyntax(stepType), typeRwFunction(stepType, write))
		}
	}
}

func writeProtocolStep(w *formatting.IndentedWriter, step *dsl.ProtocolStep, changes map[string]dsl.TypeChange, isPlural bool, write bool) {

	writeChangeSwitchCase(w, changes,
		func(w *formatting.IndentedWriter) {
			writeStepRw(w, step.Type, "value", step.IsStream(), isPlural, write)
		},
		func(w *formatting.IndentedWriter) {
			// Handle "added" ProtocolSteps
			if !write {
				if isPlural {
					fmt.Fprintln(w, "values.clear();")
				} else {
					tmpVarType := common.TypeSyntax(step.Type)
					tmpVarName := common.FieldIdentifierName(step.Name)
					fmt.Fprintf(w, "%s %s = {};\n", tmpVarType, tmpVarName)
					fmt.Fprintf(w, "value = std::move(%s);\n", tmpVarName)
					if step.IsStream() {
						fmt.Fprintf(w, "return false;\n")
					}
				}
			}
		},
		func(w *formatting.IndentedWriter, change dsl.TypeChange) {
			// Handle conversions for ProtocolStep type changes
			if requiresExplicitConversion(change) {
				tmpVarName := common.FieldIdentifierName(step.Name)
				tmpVarType := common.TypeSyntax(change.OldType())
				fmt.Fprintf(w, "%s %s = {};\n", tmpVarType, tmpVarName)

				if write {
					writeTypeConversion(w, change, "value", tmpVarName, write)
					writeStepRw(w, change.OldType(), tmpVarName, step.IsStream(), isPlural, write)
				} else {
					writeStepRw(w, change.OldType(), tmpVarName, step.IsStream(), isPlural, write)
					writeTypeConversion(w, change, tmpVarName, "value", write)
				}
			} else {
				writeStepRw(w, change.OldType(), "value", step.IsStream(), isPlural, write)
			}
		},
	)
}

func writeEndStream(w *formatting.IndentedWriter, stepChanges map[string]dsl.TypeChange) {

	writeChangeSwitchCase(w, stepChanges,
		func(w *formatting.IndentedWriter) {
			w.WriteString("yardl::binary::WriteInteger(stream_, 0U);\n")
		},
		func(w *formatting.IndentedWriter) {
			// No-op: Don't write stream end for added ProtocolSteps
		},
		nil)
}

func typeDefinitionRwFunction(t dsl.TypeDefinition, write bool) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		suffix := func() string {
			switch t {
			case dsl.Bool, dsl.Int8, dsl.Int16, dsl.Int32, dsl.Int64, dsl.Uint8, dsl.Uint16, dsl.Uint32, dsl.Uint64, dsl.Size:
				return "Integer"
			case dsl.Float32, dsl.Float64, dsl.ComplexFloat32, dsl.PrimitiveComplexFloat64:
				return "FloatingPoint"
			case dsl.String:
				return "String"
			case dsl.Date:
				return "Date"
			case dsl.Time:
				return "Time"
			case dsl.DateTime:
				return "DateTime"
			default:
				panic(fmt.Sprintf("Unknown primitive type %s", t))
			}
		}()
		if write {
			return "yardl::binary::Write" + suffix
		}
		return "yardl::binary::Read" + suffix
	case *dsl.EnumDefinition:
		if t.IsFlags {
			return fmt.Sprintf("yardl::binary::%sFlags<%s>", verb(write), common.TypeDefinitionSyntax(t))
		}
		return fmt.Sprintf("yardl::binary::%sEnum<%s>", verb(write), common.TypeDefinitionSyntax(t))
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("%s%s", verb(write), common.TypeDefinitionSyntax(t))
	default:
		meta := t.GetDefinitionMeta()
		suffix := meta.Name
		verb := verb(write)
		typeArgumentsString := ""
		if len(meta.TypeParameters) > 0 {
			typeArguments := make([]string, 2*len(meta.TypeParameters))
			if len(meta.TypeArguments) > 0 {
				for i, p := range meta.TypeArguments {
					typeArguments[2*i] = common.TypeSyntax(p)
					typeArguments[2*i+1] = typeRwFunction(p, write)
				}
			} else {
				for i, p := range meta.TypeParameters {
					typeArguments[2*i] = common.TypeDefinitionSyntax(p)
					typeArguments[2*i+1] = fmt.Sprintf("%s%s", verb, common.TypeDefinitionSyntax(p))
				}
			}
			typeArgumentsString = fmt.Sprintf("<%s>", strings.Join(typeArguments, ", "))
		}
		return fmt.Sprintf("%s::binary::%s%s%s", common.NamespaceIdentifierName(meta.Namespace), verb, suffix, typeArgumentsString)
	}
}

func typeRwFunction(t dsl.Type, write bool) string {
	switch t := t.(type) {
	case nil:
		return fmt.Sprintf("yardl::binary::%sMonostate", verb(write))
	case *dsl.SimpleType:
		return typeDefinitionRwFunction(t.ResolvedDefinition, write)
	case *dsl.GeneralizedType:
		scalarType := t.ToScalar()
		scalarFunction := func() string {
			if t.Cases.IsSingle() {
				return typeRwFunction(t.Cases[0].Type, write)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("yardl::binary::%sOptional<%s, %s>", verb(write), common.TypeSyntax(t.Cases[1].Type), typeRwFunction(t.Cases[1].Type, write))
			}

			templateArguments := make([]string, 2*len(t.Cases))
			for i, c := range t.Cases {
				templateArguments[2*i] = common.TypeSyntax(c.Type)
				templateArguments[2*i+1] = typeRwFunction(c.Type, write)
			}

			return fmt.Sprintf("%sUnion<%s>", verb(write), strings.Join(templateArguments, ", "))
		}()

		switch td := t.Dimensionality.(type) {
		case nil, *dsl.Stream:
			return scalarFunction
		case *dsl.Vector:
			if td.Length == nil {
				return fmt.Sprintf("yardl::binary::%sVector<%s, %s>", verb(write), common.TypeSyntax(scalarType), scalarFunction)
			}
			return fmt.Sprintf("yardl::binary::%sArray<%s, %s, %d>", verb(write), common.TypeSyntax(scalarType), scalarFunction, *td.Length)
		case *dsl.Array:
			if td.IsFixed() {
				lengths := make([]string, len(*td.Dimensions))
				for i, d := range *td.Dimensions {
					lengths[i] = strconv.FormatUint(*d.Length, 10)
				}
				return fmt.Sprintf("yardl::binary::%sFixedNDArray<%s, %s, %s>", verb(write), common.TypeSyntax(scalarType), scalarFunction, strings.Join(lengths, ", "))
			}
			if td.HasKnownNumberOfDimensions() {
				return fmt.Sprintf("yardl::binary::%sNDArray<%s, %s, %d>", verb(write), common.TypeSyntax(scalarType), scalarFunction, len(*td.Dimensions))
			}

			return fmt.Sprintf("yardl::binary::%sDynamicNDArray<%s, %s>", verb(write), common.TypeSyntax(scalarType), scalarFunction)
		case *dsl.Map:
			return fmt.Sprintf("yardl::binary::%sMap<%s, %s, %s, %s>", verb(write), common.TypeSyntax(td.KeyType), common.TypeSyntax(scalarType), typeRwFunction(td.KeyType, write), scalarFunction)
		default:
			panic(fmt.Sprintf("Unknown dimensionality type %T", td))
		}
	default:
		panic(fmt.Sprintf("Unknown type: %T", t))
	}
}

func BinaryWriterClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sWriter", p.Name)
}

func QualifiedBinaryWriterClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::binary::%s", common.TypeNamespaceIdentifierName(p), BinaryWriterClassName(p))
}

func BinaryReaderClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sReader", p.Name)
}

func QualifiedBinaryReaderClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::binary::%s", common.TypeNamespaceIdentifierName(p), BinaryReaderClassName(p))
}
