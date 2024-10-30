// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package common

import (
	"fmt"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

var reservedNames = map[string]any{
	"False":    nil,
	"None":     nil,
	"True":     nil,
	"and":      nil,
	"as":       nil,
	"assert":   nil,
	"async":    nil,
	"await":    nil,
	"break":    nil,
	"class":    nil,
	"continue": nil,
	"def":      nil,
	"del":      nil,
	"elif":     nil,
	"else":     nil,
	"except":   nil,
	"finally":  nil,
	"for":      nil,
	"from":     nil,
	"global":   nil,
	"if":       nil,
	"import":   nil,
	"in":       nil,
	"is":       nil,
	"lambda":   nil,
	"nonlocal": nil,
	"not":      nil,
	"or":       nil,
	"pass":     nil,
	"raise":    nil,
	"return":   nil,
	"try":      nil,
	"while":    nil,
	"with":     nil,
	"yield":    nil,
	"case":     nil,
	"match":    nil,
	// builtin types
	"bool":    nil,
	"int":     nil,
	"float":   nil,
	"complex": nil,
	"str":     nil,
}

var TypeSyntaxWriter dsl.TypeSyntaxWriter[string] = func(self dsl.TypeSyntaxWriter[string], t dsl.Node, contextNamespace string) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		switch t {
		case dsl.Bool:
			return "bool"
		case dsl.Int8:
			return "yardl.Int8"
		case dsl.Uint8:
			return "yardl.UInt8"
		case dsl.Int16:
			return "yardl.Int16"
		case dsl.Uint16:
			return "yardl.UInt16"
		case dsl.Int32:
			return "yardl.Int32"
		case dsl.Uint32:
			return "yardl.UInt32"
		case dsl.Int64:
			return "yardl.Int64"
		case dsl.Uint64:
			return "yardl.UInt64"
		case dsl.Size:
			return "yardl.Size"
		case dsl.Float32:
			return "yardl.Float32"
		case dsl.Float64:
			return "yardl.Float64"
		case dsl.ComplexFloat32:
			return "yardl.ComplexFloat"
		case dsl.ComplexFloat64:
			return "yardl.ComplexDouble"
		case dsl.String:
			return "str"
		case dsl.Date:
			return "datetime.date"
		case dsl.Time:
			return "yardl.Time"
		case dsl.DateTime:
			return "yardl.DateTime"
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

		if len(meta.TypeParameters) > 0 {

			typeArguments := make([]string, 0, len(meta.TypeParameters))
			if len(meta.TypeArguments) > 0 {
				for i, typeArg := range meta.TypeArguments {
					typeParameter := meta.TypeParameters[i]
					use := typeParameter.Annotations[TypeParameterUseAnnotationKey].(TypeParameterUse)
					if use&TypeParameterUseScalar != 0 {
						typeArguments = append(typeArguments, self.ToSyntax(typeArg, contextNamespace))
					}
					if use&TypeParameterUseArray != 0 {
						typeArguments = append(typeArguments, TypeArrayTypeArgument(typeArg))
					}
				}
			} else {
				for i, typeParam := range meta.TypeParameters {
					typeParameter := meta.TypeParameters[i]
					use := typeParameter.Annotations[TypeParameterUseAnnotationKey].(TypeParameterUse)
					if use&TypeParameterUseScalar != 0 {
						typeArguments = append(typeArguments, self.ToSyntax(typeParam, contextNamespace))
					}
					if use&TypeParameterUseArray != 0 {
						typeArguments = append(typeArguments, NumpyTypeParameterSyntax(typeParam))
					}
				}
			}

			typeSyntax = fmt.Sprintf("%s[%s]", typeName, strings.Join(typeArguments, ", "))
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

			return UnionSyntax(t)
		}()

		switch d := t.Dimensionality.(type) {
		case nil, *dsl.Stream:
			return scalarString
		case *dsl.Vector:
			return fmt.Sprintf("list[%s]", scalarString)
		case *dsl.Array:
			return fmt.Sprintf("npt.NDArray[%s]", TypeArrayTypeArgument(t.ToScalar()))
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

var typeSyntaxWithoutTypeParametersWriter dsl.TypeSyntaxWriter[string] = func(self dsl.TypeSyntaxWriter[string], t dsl.Node, contextNamespace string) string {
	switch t := t.(type) {
	case dsl.TypeDefinition:
		meta := t.GetDefinitionMeta()
		if len(meta.TypeParameters) > 0 {
			meta := t.GetDefinitionMeta()
			typeName := TypeIdentifierName(meta.Name)
			if t.GetDefinitionMeta().Namespace != contextNamespace {
				typeName = fmt.Sprintf("%s.%s", formatting.ToSnakeCase(meta.Namespace), typeName)
			}

			return typeName
		}
	}

	return TypeSyntaxWriter(self, t, contextNamespace)
}

func TypeSyntaxWithoutTypeParameters(typeOrTypeDefinition dsl.Node, contextNamespace string) string {
	return typeSyntaxWithoutTypeParametersWriter.ToSyntax(typeOrTypeDefinition, contextNamespace)
}

func TypeDTypeSyntax(t dsl.Type) string {
	switch t := t.(type) {
	case *dsl.SimpleType:
		return TypeDefinitionDTypeSyntax(t.ResolvedDefinition)
	case *dsl.GeneralizedType:
		if len(t.Cases) > 1 {
			return "np.object_"
		}
		switch td := t.Dimensionality.(type) {
		case nil:
			return TypeDTypeSyntax(t.Cases[0].Type)
		case *dsl.Vector:
			if td.Length == nil {
				return "np.object_"
			}
			scalarDType := TypeDTypeSyntax(t.ToScalar())
			return fmt.Sprintf("%s, (%d,)", scalarDType, *td.Length)
		case *dsl.Array:
			if !td.IsFixed() {
				return "np.object_"
			}
			scalarDType := TypeDTypeSyntax(t.ToScalar())
			dims := make([]string, len(*td.Dimensions))
			for i, dim := range *td.Dimensions {
				dims[i] = fmt.Sprintf("%d", *dim.Length)
			}
			return fmt.Sprintf("%s, (%s)", scalarDType, strings.Join(dims, ", "))

		default:
			return "np.object_"
		}
	default:
		panic(fmt.Sprintf("Dype for %T not implemented", t))
	}
}

func TypeDefinitionDTypeSyntax(t dsl.TypeDefinition) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		switch t {
		case dsl.Bool:
			return "np.bool_"
		case dsl.Int8, dsl.Uint8, dsl.Int16, dsl.Uint16, dsl.Int32, dsl.Uint32, dsl.Int64, dsl.Uint64, dsl.Float32, dsl.Float64:
			return fmt.Sprintf("np.%s", strings.ToLower(string(t)))
		case dsl.Size:
			return "np.uint64"
		case dsl.ComplexFloat32:
			return "np.complex64"
		case dsl.ComplexFloat64:
			return "np.complex128"
		case dsl.Date:
			return "np.datetime64"
		case dsl.Time:
			return "np.timedelta64"
		case dsl.DateTime:
			return "np.datetime64"
		case dsl.String:
			return "np.object_"
		default:
			panic(fmt.Sprintf("Not implemented %s", t))
		}
	case *dsl.RecordDefinition:
		fields := make([]string, len(t.Fields))
		for i, field := range t.Fields {
			fields[i] = fmt.Sprintf("('%s', %s)", field.Name, TypeDTypeSyntax(field.Type))
		}

		return fmt.Sprintf("np.dtype([%s], align=True)", strings.Join(fields, ", "))
	case *dsl.EnumDefinition:
		if t.BaseType == nil {
			return TypeDefinitionDTypeSyntax(dsl.PrimitiveInt32)
		}

		return TypeDTypeSyntax(t.BaseType)
	case *dsl.NamedType:
		return TypeDTypeSyntax(t.Type)
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("%s_dtype", formatting.ToSnakeCase(t.Name))
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func TypeDefinitionArrayTypeArgument(t dsl.TypeDefinition) string {
	switch t := t.(type) {
	case *dsl.RecordDefinition:
		return "np.void"
	case *dsl.GenericTypeParameter:
		return NumpyTypeParameterSyntax(t)
	default:
		return TypeDefinitionDTypeSyntax(t)
	}
}

func TypeParameterSyntax(p *dsl.GenericTypeParameter, numpy bool) string {
	if numpy {
		return NumpyTypeParameterSyntax(p)
	}

	return TypeIdentifierName(p.Name)
}

func NumpyTypeParameterSyntax(p *dsl.GenericTypeParameter) string {
	return fmt.Sprintf("%s_NP", TypeIdentifierName(p.Name))
}

func TypeArrayTypeArgument(t dsl.Type) string {
	switch t := t.(type) {
	case *dsl.SimpleType:
		return TypeDefinitionArrayTypeArgument(t.ResolvedDefinition)
	case *dsl.GeneralizedType:
		if len(t.Cases) > 1 {
			return "np.object_"
		}
		switch td := t.Dimensionality.(type) {
		case nil:
			return TypeArrayTypeArgument(t.Cases[0].Type)
		case *dsl.Vector:
			if td.Length == nil {
				return "np.object_"
			}
			return TypeArrayTypeArgument(t.ToScalar())
		case *dsl.Array:
			if !td.IsFixed() {
				return "np.object_"
			}
			return TypeArrayTypeArgument(t.ToScalar())
		default:
			return "np.object_"
		}
	default:
		panic(fmt.Sprintf("DType for %T not implemented", t))
	}
}

func NamespaceIdentifierName(namespace string) string {
	return formatting.ToSnakeCase(namespace)
}

func TypeNamespaceIdentifierName(t dsl.TypeDefinition) string {
	return NamespaceIdentifierName(t.GetDefinitionMeta().Namespace)
}

func FieldIdentifierName(name string) string {
	snakeCased := formatting.ToSnakeCase(name)
	if _, reserved := reservedNames[snakeCased]; !reserved {
		return snakeCased
	}

	return snakeCased + "_"
}

func EnumValueIdentifierName(name string) string {
	cased := formatting.ToUpperSnakeCase(name)
	if _, reserved := reservedNames[cased]; !reserved {
		return cased
	}

	return cased + "_"
}

func ComputedFieldIdentifierName(name string) string {
	cased := formatting.ToSnakeCase(name)
	if _, reserved := reservedNames[cased]; !reserved {
		return cased
	}

	return cased + "_"
}

func TypeIdentifierName(name string) string {
	if _, reserved := reservedNames[name]; !reserved {
		return name
	}

	return name + "_"
}

func UnionClassName(gt *dsl.GeneralizedType) (className string, typeParameters string) {
	if !gt.Cases.IsUnion() {
		panic("Not a union")
	}

	cases := make([]string, 0, len(gt.Cases))
	for _, typeCase := range gt.Cases {
		if typeCase.Type == nil {
			continue
		}
		cases = append(cases, formatting.ToPascalCase(typeCase.Tag))
	}

	return strings.Join(cases, "Or"), GetOpenGenericTypeParameters(gt)
}

func UnionSyntax(gt *dsl.GeneralizedType) string {
	className, typeParameters := UnionClassName(gt)
	var syntax string
	if len(typeParameters) > 0 {
		syntax = fmt.Sprintf("%s[%s]", className, typeParameters)
	} else {
		syntax = className
	}

	if gt.Cases.HasNullOption() {
		return fmt.Sprintf("typing.Optional[%s]", syntax)
	}

	return syntax
}

// Returns open type parameters used within the node as a comma-separated string.
// e.g. "T1, T2, T2_NP, T3"
func GetOpenGenericTypeParameters(node dsl.Node) string {
	var res []*dsl.GenericTypeParameter
	var paramStrings []string

	dsl.Visit(node, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case *dsl.GenericTypeParameter:
			for _, existing := range res {
				if t == existing {
					return
				}
			}
			res = append(res, t)
			paramStrings = appendGenericTypeParameterUses(paramStrings, t)
		case *dsl.SimpleType:
			if gtp, ok := t.ResolvedDefinition.(*dsl.GenericTypeParameter); ok {
				found := false
				for _, existing := range res {
					if gtp == existing {
						found = true
					}
				}
				if !found {
					res = append(res, gtp)
					paramStrings = appendGenericTypeParameterUses(paramStrings, gtp)
				}
			}
		}
		self.VisitChildren(node)
	})

	if len(paramStrings) > 0 {
		return strings.Join(paramStrings, ", ")
	}

	return ""
}

func appendGenericTypeParameterUses(slice []string, typeParameter *dsl.GenericTypeParameter) []string {
	if slice == nil {
		slice = make([]string, 0, 1)
	}
	use := typeParameter.Annotations[TypeParameterUseAnnotationKey].(TypeParameterUse)
	if use&TypeParameterUseScalar != 0 {
		slice = append(slice, TypeParameterSyntax(typeParameter, false))
	}
	if use&TypeParameterUseArray != 0 {
		slice = append(slice, TypeParameterSyntax(typeParameter, true))
	}

	return slice
}

func WriteComment(w *formatting.IndentedWriter, comment string) {
	comment = strings.TrimSpace(comment)
	if comment != "" {
		w = formatting.NewIndentedWriter(w, "# ").Indent()
		w.WriteStringln(comment)
	}
}

func WriteDocstring(w *formatting.IndentedWriter, comment string) {
	comment = strings.TrimSpace(comment)
	if comment != "" {
		w.WriteString(`"""`)
		w.WriteString(comment)
		if strings.Contains(comment, "\n") {
			w.WriteStringln("")
		}
		w.WriteStringln("\"\"\"\n")
	}
}

func WriteDocstringWithLeadingLine(w *formatting.IndentedWriter, first, rest string) {
	first = strings.TrimSpace(first)
	if first == "" {
		WriteDocstring(w, rest)
		return
	}

	rest = strings.TrimSpace(rest)
	if rest == "" {
		WriteDocstring(w, first)
		return
	}

	WriteDocstring(w, first+"\n\n"+rest)
}

func AbstractWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sWriterBase", formatting.ToPascalCase(p.Name))
}

func AbstractReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sReaderBase", formatting.ToPascalCase(p.Name))
}

func AbstractIndexedReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sIndexedReaderBase", formatting.ToPascalCase(p.Name))
}

func ProtocolWriteMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("write_%s", formatting.ToSnakeCase(s.Name))
}

func ProtocolWriteImplMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("_write_%s", formatting.ToSnakeCase(s.Name))
}

func ProtocolReadMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("read_%s", formatting.ToSnakeCase(s.Name))
}

func ProtocolReadImplMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("_read_%s", formatting.ToSnakeCase(s.Name))
}

func ProtocolStreamSizeMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("count_%s", formatting.ToSnakeCase(s.Name))
}

func ProtocolStreamSizeImplMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("_count_%s", formatting.ToSnakeCase(s.Name))
}

func WriteGeneratedFileHeader(w *formatting.IndentedWriter) {
	WriteComment(w, "This file was generated by the \"yardl\" tool. DO NOT EDIT.")
	w.WriteStringln("")
}
