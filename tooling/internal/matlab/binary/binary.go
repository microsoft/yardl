// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package binary

import (
	"bytes"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/internal/matlab/common"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func WriteBinary(ns *dsl.Namespace, packageDir string) error {

	if ns.IsTopLevel {
		if err := writeProtocols(ns, packageDir); err != nil {
			return err
		}
	}

	return writeRecordSerializers(ns, packageDir)
}

func writeProtocols(ns *dsl.Namespace, packageDir string) error {
	for _, p := range ns.Protocols {

		if err := writeProtocolWriter(p, ns, packageDir); err != nil {
			return err
		}

		if err := writeProtocolReader(p, ns, packageDir); err != nil {
			return err
		}
	}
	return nil
}

func writeProtocolWriter(p *dsl.ProtocolDefinition, ns *dsl.Namespace, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)

	common.WriteComment(w, fmt.Sprintf("Binary writer for the %s protocol", p.Name))
	common.WriteComment(w, p.Comment)
	fmt.Fprintf(w, "classdef %s < yardl.binary.BinaryProtocolWriter & %s\n", BinaryWriterName(p), common.AbstractWriterName(p))
	common.WriteBlockBody(w, func() {

		w.WriteStringln("methods")
		common.WriteBlockBody(w, func() {
			fmt.Fprintf(w, "function obj = %s(filename)\n", BinaryWriterName(p))
			common.WriteBlockBody(w, func() {
				fmt.Fprintf(w, "obj@%s();\n", common.AbstractWriterName(p))
				fmt.Fprintf(w, "obj@yardl.binary.BinaryProtocolWriter(filename, %s.schema);\n", common.AbstractWriterName(p))
			})
		})
		w.WriteStringln("")

		w.WriteStringln("methods (Access=protected)")
		common.WriteBlockBody(w, func() {
			for i, step := range p.Sequence {
				fmt.Fprintf(w, "function %s(obj, value)\n", common.ProtocolWriteImplMethodName(step))
				common.WriteBlockBody(w, func() {
					fmt.Fprintf(w, "w = %s;\n", typeSerializer(step.Type, ns.Name, nil))
					w.WriteStringln("w.write(obj.stream_, value);")
				})
				if i < len(p.Sequence)-1 {
					w.WriteStringln("")
				}
			}
		})
	})

	binaryPath := path.Join(packageDir, fmt.Sprintf("%s.m", BinaryWriterName(p)))
	return iocommon.WriteFileIfNeeded(binaryPath, b.Bytes(), 0644)
}

func writeProtocolReader(p *dsl.ProtocolDefinition, ns *dsl.Namespace, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)

	common.WriteComment(w, fmt.Sprintf("Binary reader for the %s protocol", p.Name))
	common.WriteComment(w, p.Comment)
	fmt.Fprintf(w, "classdef %s < yardl.binary.BinaryProtocolReader & %s\n", BinaryReaderName(p), common.AbstractReaderName(p))
	common.WriteBlockBody(w, func() {

		w.WriteStringln("methods")
		common.WriteBlockBody(w, func() {
			fmt.Fprintf(w, "function obj = %s(filename)\n", BinaryReaderName(p))
			common.WriteBlockBody(w, func() {
				fmt.Fprintf(w, "obj@%s();\n", common.AbstractReaderName(p))
				fmt.Fprintf(w, "obj@yardl.binary.BinaryProtocolReader(filename, %s.schema);\n", common.AbstractReaderName(p))
			})
		})
		w.WriteStringln("")

		w.WriteStringln("methods (Access=protected)")
		common.WriteBlockBody(w, func() {
			for i, step := range p.Sequence {
				fmt.Fprintf(w, "function value = %s(obj, value)\n", common.ProtocolReadImplMethodName(step))
				common.WriteBlockBody(w, func() {
					fmt.Fprintf(w, "r = %s;\n", typeSerializer(step.Type, ns.Name, nil))
					w.WriteStringln("value = r.read(obj.stream_);")
				})

				if i < len(p.Sequence)-1 {
					w.WriteStringln("")
				}
			}
		})
	})

	binaryPath := path.Join(packageDir, fmt.Sprintf("%s.m", BinaryReaderName(p)))
	return iocommon.WriteFileIfNeeded(binaryPath, b.Bytes(), 0644)
}

func writeRecordSerializers(ns *dsl.Namespace, packageDir string) error {
	for _, td := range ns.TypeDefinitions {
		switch td := td.(type) {
		case *dsl.RecordDefinition:
			if err := writeRecordSerializer(td, ns, packageDir); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeRecordSerializer(rec *dsl.RecordDefinition, ns *dsl.Namespace, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)

	typeSyntax := common.TypeSyntax(rec, ns.Name)
	fmt.Fprintf(w, "classdef %s < yardl.binary.RecordSerializer\n", recordSerializerClassName(rec, ns.Name))
	common.WriteBlockBody(w, func() {

		w.WriteStringln("methods")
		common.WriteBlockBody(w, func() {
			fmt.Fprintf(w, "function obj = %s()\n", recordSerializerClassName(rec, ns.Name))
			common.WriteBlockBody(w, func() {
				for i, field := range rec.Fields {
					fmt.Fprintf(w, "field_serializers{%d} = %s;\n", i+1, typeSerializer(field.Type, ns.Name, nil))
				}
				fmt.Fprintf(w, "obj@yardl.binary.RecordSerializer(field_serializers);\n")
			})
			w.WriteStringln("")

			fmt.Fprintf(w, "function write(obj, outstream, value)\n")
			common.WriteBlockBody(w, func() {
				fmt.Fprintf(w, "assert(isa(value, '%s'));\n", typeSyntax)

				fieldAccesses := make([]string, len(rec.Fields))
				for i, field := range rec.Fields {
					fieldAccesses[i] = fmt.Sprintf("value.%s", common.FieldIdentifierName(field.Name))
				}
				fmt.Fprintf(w, "obj.write_(outstream, %s)\n", strings.Join(fieldAccesses, ", "))
			})
			w.WriteStringln("")

			fmt.Fprintf(w, "function value = read(obj, instream)\n")
			common.WriteBlockBody(w, func() {
				w.WriteStringln("field_values = obj.read_(instream);")
				fmt.Fprintf(w, "value = %s(field_values{:});\n", typeSyntax)
			})
		})
	})

	binaryPath := path.Join(packageDir, fmt.Sprintf("%s.m", recordSerializerClassName(rec, ns.Name)))
	return iocommon.WriteFileIfNeeded(binaryPath, b.Bytes(), 0644)
}

func recordSerializerClassName(record *dsl.RecordDefinition, contextNamespace string) string {
	className := fmt.Sprintf("%sSerializer", formatting.ToPascalCase(record.Name))
	if record.Namespace != contextNamespace {
		className = fmt.Sprintf("%s.binary.%s", common.NamespaceIdentifierName(record.Namespace), className)
	}
	return className
}

func typeDefinitionSerializer(t dsl.TypeDefinition, contextNamespace string) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		return fmt.Sprintf("yardl.binary.%sSerializer", formatting.ToPascalCase(string(t)))
	case *dsl.EnumDefinition:
		var baseType dsl.Type
		if t.BaseType != nil {
			baseType = t.BaseType
		} else {
			baseType = dsl.Int32Type
		}

		elementSerializer := typeSerializer(baseType, contextNamespace, nil)
		return fmt.Sprintf("yardl.binary.EnumSerializer(%s, @%s)", elementSerializer, common.TypeSyntax(t, contextNamespace))
	case *dsl.RecordDefinition:
		serializerName := recordSerializerClassName(t, contextNamespace)
		if len(t.TypeParameters) == 0 {
			return fmt.Sprintf("%s()", serializerName)
		}
		if len(t.TypeArguments) == 0 {
			panic("Expected type arguments")
		}

		typeArguments := make([]string, 0, len(t.TypeArguments))
		for _, arg := range t.TypeArguments {
			typeArguments = append(typeArguments, typeSerializer(arg, contextNamespace, nil))
		}

		if len(typeArguments) == 0 {
			return fmt.Sprintf("%s()", serializerName)
		}

		return fmt.Sprintf("%s(%s)", serializerName, strings.Join(typeArguments, ", "))
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("%s_serializer", formatting.ToSnakeCase(t.Name))
	case *dsl.NamedType:
		return typeSerializer(t.Type, contextNamespace, t)
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func typeSerializer(t dsl.Type, contextNamespace string, namedType *dsl.NamedType) string {
	switch t := t.(type) {
	case nil:
		return "yardl.binary.none_serializer"
	case *dsl.SimpleType:
		return typeDefinitionSerializer(t.ResolvedDefinition, contextNamespace)
	case *dsl.GeneralizedType:
		getScalarSerializer := func() string {
			if t.Cases.IsSingle() {
				return typeSerializer(t.Cases[0].Type, contextNamespace, namedType)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("yardl.binary.OptionalSerializer(%s)", typeSerializer(t.Cases[1].Type, contextNamespace, namedType))
			}

			// TODO:!
			return fmt.Sprintf("NOT YET IMPLEMENTED")

			// unionClassName, typeParameters := common.UnionClassName(t)
			// if namedType != nil {
			// 	unionClassName = namedType.Name
			// 	if namedType.Namespace != contextNamespace {
			// 		unionClassName = fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(namedType.Namespace), unionClassName)
			// 	}
			// }

			// var classSyntax string
			// if len(typeParameters) == 0 {
			// 	classSyntax = unionClassName
			// } else {
			// 	classSyntax = fmt.Sprintf("%s[%s]", unionClassName, typeParameters)
			// }
			// options := make([]string, len(t.Cases))
			// for i, c := range t.Cases {
			// 	if c.Type == nil {
			// 		options[i] = "None"
			// 	} else {
			// 		options[i] = fmt.Sprintf("(%s.%s, %s)", classSyntax, formatting.ToPascalCase(c.Tag), typeSerializer(c.Type, contextNamespace, namedType))
			// 	}
			// }

			// return fmt.Sprintf("yardl.binary.UnionSerializer(%s, [%s])", unionClassName, strings.Join(options, ", "))

		}
		switch td := t.Dimensionality.(type) {
		case nil:
			return getScalarSerializer()
		case *dsl.Stream:
			return fmt.Sprintf("yardl.binary.StreamSerializer(%s)", getScalarSerializer())
		case *dsl.Vector:
			if td.Length != nil {
				return fmt.Sprintf("yardl.binary.FixedVectorSerializer(%s, %d)", getScalarSerializer(), *td.Length)
			}

			return fmt.Sprintf("yardl.binary.VectorSerializer(%s)", getScalarSerializer())
		case *dsl.Array:
			if td.IsFixed() {
				dims := make([]string, len(*td.Dimensions))
				for i, d := range *td.Dimensions {
					dims[len(*td.Dimensions)-i-1] = strconv.FormatUint(*d.Length, 10)
				}

				return fmt.Sprintf("yardl.binary.FixedNDArraySerializer(%s, [%s])", getScalarSerializer(), strings.Join(dims, ", "))
			}

			if td.HasKnownNumberOfDimensions() {
				return fmt.Sprintf("yardl.binary.NDArraySerializer(%s, %d)", getScalarSerializer(), len(*td.Dimensions))
			}

			return fmt.Sprintf("yardl.binary.DynamicNDArraySerializer(%s)", getScalarSerializer())

		case *dsl.Map:
			keySerializer := typeSerializer(td.KeyType, contextNamespace, namedType)
			valueSerializer := typeSerializer(t.ToScalar(), contextNamespace, namedType)

			return fmt.Sprintf("yardl.binary.MapSerializer(%s, %s)", keySerializer, valueSerializer)
		default:
			panic(fmt.Sprintf("Not implemented %T", t.Dimensionality))
		}
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func BinaryWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("Binary%sWriter", formatting.ToPascalCase(p.Name))
}

func BinaryReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("Binary%sReader", formatting.ToPascalCase(p.Name))
}
