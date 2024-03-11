// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package binary

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/matlab/common"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func WriteBinary(fw *common.MatlabFileWriter, ns *dsl.Namespace) error {

	if ns.IsTopLevel {
		if err := writeProtocols(fw, ns); err != nil {
			return err
		}
	}

	return writeRecordSerializers(fw, ns)
}

func writeProtocols(fw *common.MatlabFileWriter, ns *dsl.Namespace) error {
	for _, p := range ns.Protocols {
		if err := writeProtocolWriter(fw, p, ns); err != nil {
			return err
		}

		if err := writeProtocolReader(fw, p, ns); err != nil {
			return err
		}
	}
	return nil
}

func writeProtocolWriter(fw *common.MatlabFileWriter, p *dsl.ProtocolDefinition, ns *dsl.Namespace) error {
	return fw.WriteFile(BinaryWriterName(p), func(w *formatting.IndentedWriter) {
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
	})
}

func writeProtocolReader(fw *common.MatlabFileWriter, p *dsl.ProtocolDefinition, ns *dsl.Namespace) error {
	return fw.WriteFile(BinaryReaderName(p), func(w *formatting.IndentedWriter) {
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
	})
}

func writeRecordSerializers(fw *common.MatlabFileWriter, ns *dsl.Namespace) error {
	for _, td := range ns.TypeDefinitions {
		switch td := td.(type) {
		case *dsl.RecordDefinition:
			if err := writeRecordSerializer(fw, td, ns); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeRecordSerializer(fw *common.MatlabFileWriter, rec *dsl.RecordDefinition, ns *dsl.Namespace) error {
	return fw.WriteFile(recordSerializerClassName(rec, ns.Name), func(w *formatting.IndentedWriter) {

		typeSyntax := common.TypeSyntax(rec, ns.Name)
		fmt.Fprintf(w, "classdef %s < yardl.binary.RecordSerializer\n", recordSerializerClassName(rec, ns.Name))
		common.WriteBlockBody(w, func() {

			w.WriteStringln("methods")
			common.WriteBlockBody(w, func() {
				if len(rec.TypeParameters) > 0 {
					typeParamSerializers := make([]string, 0, len(rec.TypeParameters))
					for _, tp := range rec.TypeParameters {
						typeParamSerializers = append(
							typeParamSerializers,
							typeDefinitionSerializer(tp, ns.Name))
					}

					fmt.Fprintf(w, "function obj = %s(%s)\n", recordSerializerClassName(rec, ns.Name), strings.Join(typeParamSerializers, ", "))
				} else {
					fmt.Fprintf(w, "function obj = %s()\n", recordSerializerClassName(rec, ns.Name))
				}

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
	})
}

func recordSerializerClassName(record *dsl.RecordDefinition, contextNamespace string) string {
	className := fmt.Sprintf("%sSerializer", formatting.ToPascalCase(record.Name))
	if record.Namespace != contextNamespace {
		className = fmt.Sprintf("%s.binary.%s", common.NamespaceIdentifierName(record.Namespace), className)
	}
	return className
}

func typeDefinitionSerializer(td dsl.TypeDefinition, contextNamespace string) string {
	switch td := td.(type) {
	case dsl.PrimitiveDefinition:
		return fmt.Sprintf("yardl.binary.%sSerializer", formatting.ToPascalCase(string(td)))
	case *dsl.EnumDefinition:
		var baseType dsl.Type
		if td.BaseType != nil {
			baseType = td.BaseType
		} else {
			baseType = dsl.Int32Type
		}

		elementSerializer := typeSerializer(baseType, contextNamespace, nil)
		return fmt.Sprintf("yardl.binary.EnumSerializer(%s, @%s)", elementSerializer, common.TypeSyntax(td, contextNamespace))
	case *dsl.RecordDefinition:
		serializerName := recordSerializerClassName(td, contextNamespace)
		if len(td.TypeParameters) == 0 {
			return fmt.Sprintf("%s()", serializerName)
		}
		if len(td.TypeArguments) == 0 {
			panic("Expected type arguments")
		}

		typeArguments := make([]string, 0, len(td.TypeArguments))
		for _, arg := range td.TypeArguments {
			typeArguments = append(typeArguments, typeSerializer(arg, contextNamespace, nil))
		}

		if len(typeArguments) == 0 {
			return fmt.Sprintf("%s()", serializerName)
		}

		return fmt.Sprintf("%s(%s)", serializerName, strings.Join(typeArguments, ", "))
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("%s_serializer", formatting.ToSnakeCase(td.Name))
	case *dsl.NamedType:
		return typeSerializer(td.Type, contextNamespace, td)
	default:
		panic(fmt.Sprintf("Not implemented %T", td))
	}
}

func typeSerializer(t dsl.Type, contextNamespace string, namedType *dsl.NamedType) string {
	switch t := t.(type) {
	case nil:
		return "yardl.binary.NoneSerializer"
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

			unionClassName := common.UnionClassName(t)
			if namedType != nil {
				unionClassName = namedType.Name
				if namedType.Namespace != contextNamespace {
					unionClassName = fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(namedType.Namespace), unionClassName)
				}
			}

			serializers := make([]string, len(t.Cases))
			factories := make([]string, len(t.Cases))
			for i, c := range t.Cases {
				if c.Type == nil {
					serializers[i] = "yardl.None"
					factories[i] = "yardl.None"
				} else {
					serializers[i] = typeSerializer(c.Type, contextNamespace, namedType)
					factories[i] = fmt.Sprintf("@%s.%s", unionClassName, formatting.ToPascalCase(c.Tag))
				}
			}

			return fmt.Sprintf("yardl.binary.UnionSerializer('%s', {%s}, {%s})", unionClassName, strings.Join(serializers, ", "), strings.Join(factories, ", "))
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
