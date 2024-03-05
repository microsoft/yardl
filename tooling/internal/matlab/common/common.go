// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package common

import (
	"fmt"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

// TODO: Populate all Matlab reserved names
var isReservedName = map[string]bool{}

var TypeSyntaxWriter dsl.TypeSyntaxWriter[string] = func(self dsl.TypeSyntaxWriter[string], t dsl.Node, contextNamespace string) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		switch t {
		case dsl.Bool:
			return "bool"
		case dsl.Int8:
			return "int8"
		case dsl.Uint8:
			return "uint8"
		case dsl.Int16:
			return "int16"
		case dsl.Uint16:
			return "uint16"
		case dsl.Int32:
			return "int32"
		case dsl.Uint32:
			return "uint32"
		case dsl.Int64:
			return "int64"
		case dsl.Uint64:
			return "uint64"
		case dsl.Size:
			return "uint64"
		case dsl.Float32:
			return "float32"
		case dsl.Float64:
			return "float64"
		case dsl.ComplexFloat32:
			return "complex"
		case dsl.ComplexFloat64:
			return "complex"
		case dsl.String:
			return "str"
		case dsl.Date:
			return "datetime"
		case dsl.Time:
			return "datetime"
		case dsl.DateTime:
			return "datetime"
		default:
			panic(fmt.Sprintf("primitive '%v' not recognized", t))
		}
	case *dsl.GenericTypeParameter:
		return TypeIdentifierName(t.Name)
	case dsl.TypeDefinition:
		meta := t.GetDefinitionMeta()
		typeName := TypeIdentifierName(meta.Name)
		if t.GetDefinitionMeta().Namespace != contextNamespace {
			typeName = fmt.Sprintf("%s.%s", formatting.ToSnakeCase(meta.Namespace), typeName)
		}

		typeSyntax := typeName

		// TODO: Generic TypeDefinitions and TypeParameters
		if len(meta.TypeParameters) > 0 {
			// typeArguments := make([]string, 0, len(meta.TypeParameters))
			// if len(meta.TypeArguments) > 0 {
			// 	for i, typeArg := range meta.TypeArguments {
			// 		typeParameter := meta.TypeParameters[i]
			// 		use := typeParameter.Annotations[TypeParameterUseAnnotationKey].(TypeParameterUse)
			// 		if use&TypeParameterUseScalar != 0 {
			// 			typeArguments = append(typeArguments, self.ToSyntax(typeArg, contextNamespace))
			// 		}
			// 		if use&TypeParameterUseArray != 0 {
			// 			typeArguments = append(typeArguments, TypeArrayTypeArgument(typeArg))
			// 		}
			// 	}
			// } else {
			// 	for i, typeParam := range meta.TypeParameters {
			// 		typeParameter := meta.TypeParameters[i]
			// 		use := typeParameter.Annotations[TypeParameterUseAnnotationKey].(TypeParameterUse)
			// 		if use&TypeParameterUseScalar != 0 {
			// 			typeArguments = append(typeArguments, self.ToSyntax(typeParam, contextNamespace))
			// 		}
			// 		if use&TypeParameterUseArray != 0 {
			// 			typeArguments = append(typeArguments, NumpyTypeParameterSyntax(typeParam))
			// 		}
			// 	}
			// }

			// typeSyntax = fmt.Sprintf("%s[%s]", typeName, strings.Join(typeArguments, ", "))
		}

		if nt, ok := t.(*dsl.NamedType); ok {
			if gt, ok := nt.Type.(*dsl.GeneralizedType); ok && gt.Cases.HasNullOption() && !gt.Cases.IsOptional() {
				typeSyntax = fmt.Sprintf("typing.Optional[%s]", typeSyntax)
			}
		}

		return typeSyntax

	case nil:
		return "None"
	case *dsl.SimpleType:
		return self.ToSyntax(t.ResolvedDefinition, contextNamespace)
	case *dsl.GeneralizedType:
		scalarString := func() string {
			if t.Cases.IsSingle() {
				return self.ToSyntax(t.Cases[0].Type, contextNamespace)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("typing.Optional[%s]", self.ToSyntax(t.Cases[1].Type, contextNamespace))
			}

			// TODO: Union syntax?
			// return UnionSyntax(t)
			panic("Unions not yet supported")
		}()

		switch d := t.Dimensionality.(type) {
		case nil, *dsl.Stream:
			return scalarString
		case *dsl.Vector:
			return fmt.Sprintf("list[%s]", scalarString)
		case *dsl.Array:
			return scalarString
			// return fmt.Sprintf("npt.NDArray[%s]", TypeArrayTypeArgument(t.ToScalar()))
		case *dsl.Map:
			return fmt.Sprintf("dict[%s, %s]", self.ToSyntax(d.KeyType, contextNamespace), scalarString)
		default:
			panic(fmt.Sprintf("unexpected type %T", d))
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}

}

func TypeSyntax(typeOrTypeDefinition dsl.Node, contextNamespace string) string {
	return TypeSyntaxWriter.ToSyntax(typeOrTypeDefinition, contextNamespace)
}

func ComputedFieldIdentifierName(name string) string {
	cased := formatting.ToSnakeCase(name)
	if !isReservedName[name] {
		return cased
	}

	return cased + "_"
}

func TypeIdentifierName(name string) string {
	if !isReservedName[name] {
		return name
	}

	return name + "_"
}

func PackageDir(name string) string {
	return fmt.Sprintf("+%s", formatting.ToSnakeCase(name))
}

func NamespaceIdentifierName(namespace string) string {
	return formatting.ToSnakeCase(namespace)
}

func FieldIdentifierName(name string) string {
	snakeCased := formatting.ToSnakeCase(name)
	if !isReservedName[snakeCased] {
		return snakeCased
	}

	return snakeCased + "_"
}

func EnumValueIdentifierName(name string) string {
	cased := formatting.ToUpperSnakeCase(name)
	if !isReservedName[cased] {
		return cased
	}

	return cased + "_"
}

func WriteBlockBody(w *formatting.IndentedWriter, f func()) {
	defer func() {
		w.WriteStringln("end")
	}()
	w.Indented(f)
}

func WriteComment(w *formatting.IndentedWriter, comment string) {
	comment = strings.TrimSpace(comment)
	if comment != "" {
		w = formatting.NewIndentedWriter(w, "% ").Indent()
		w.WriteStringln(comment)
	}
}

func AbstractWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sWriterBase", formatting.ToPascalCase(p.Name))
}

func AbstractReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sReaderBase", formatting.ToPascalCase(p.Name))
}

func ProtocolWriteMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("write_%s", formatting.ToSnakeCase(s.Name))
}

func ProtocolWriteImplMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("write_%s_", formatting.ToSnakeCase(s.Name))
}

func ProtocolReadMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("read_%s", formatting.ToSnakeCase(s.Name))
}

func ProtocolReadImplMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("read_%s_", formatting.ToSnakeCase(s.Name))
}

func WriteGeneratedFileHeader(w *formatting.IndentedWriter) {
	WriteComment(w, "This file was generated by the \"yardl\" tool. DO NOT EDIT.")
	w.WriteStringln("")
}
