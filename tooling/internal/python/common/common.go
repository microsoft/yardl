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

func TypeSyntax(t dsl.Type, contextNamespace string, includeTypeParameters bool) string {
	switch t := t.(type) {
	case nil:
		return "None"
	case *dsl.SimpleType:
		return TypeDefinitionSyntax(t.ResolvedDefinition, contextNamespace, includeTypeParameters)
	case *dsl.GeneralizedType:
		scalarString := func() string {
			if t.Cases.IsSingle() {
				return TypeSyntax(t.Cases[0].Type, contextNamespace, includeTypeParameters)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("%s | None", TypeSyntax(t.Cases[1].Type, contextNamespace, includeTypeParameters))
			}

			typeMap := make(map[string]any)

			caseStrings := make([]string, 0, len(t.Cases))
			for _, typeCase := range t.Cases {
				if typeCase.Type == nil {
					continue
				}

				syntax := TypeSyntax(typeCase.Type, contextNamespace, includeTypeParameters)

				if _, ok := typeMap[syntax]; !ok {
					typeMap[syntax] = nil
					caseStrings = append(caseStrings, syntax)
				}
			}

			if t.Cases.HasNullOption() {
				caseStrings = append(caseStrings, "None")
			}

			return strings.Join(caseStrings, " | ")
		}()

		switch d := t.Dimensionality.(type) {
		case nil, *dsl.Stream:
			return scalarString
		case *dsl.Vector:
			return fmt.Sprintf("list[%s]", scalarString)
		case *dsl.Array:
			return fmt.Sprintf("npt.NDArray[%s]", TypeDTypeTypeArgument(t.ToScalar()))
		case *dsl.Map:
			return fmt.Sprintf("dict[%s, %s]", TypeSyntax(d.KeyType, contextNamespace, includeTypeParameters), scalarString)
		default:
			panic(fmt.Sprintf("unexpected type %T", d))
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func TypeDefinitionSyntax(t dsl.TypeDefinition, contextNamespace string, includeTypeParameters bool) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		return PrimitiveSyntax(t)
	case *dsl.GenericTypeParameter:
		return TypeIdentifierName(t.Name)
	default:
		meta := t.GetDefinitionMeta()
		typeName := TypeIdentifierName(meta.Name)
		if t.GetDefinitionMeta().Namespace != contextNamespace {
			typeName = fmt.Sprintf("%s.%s", formatting.ToSnakeCase(meta.Namespace), typeName)
		}

		if len(meta.TypeParameters) == 0 || !includeTypeParameters && len(meta.TypeArguments) == 0 {
			return typeName
		}

		typeArguments := make([]string, len(meta.TypeParameters))
		if len(meta.TypeArguments) > 0 {
			for i, typeArg := range meta.TypeArguments {
				typeArguments[i] = TypeSyntax(typeArg, contextNamespace, includeTypeParameters)
			}
		} else {
			for i, typeParam := range meta.TypeParameters {
				typeArguments[i] = TypeDefinitionSyntax(typeParam, contextNamespace, includeTypeParameters)
			}
		}

		return fmt.Sprintf("%s[%s]", typeName, strings.Join(typeArguments, ", "))
	}
}

func PrimitiveSyntax(p dsl.PrimitiveDefinition) string {
	switch p {
	case dsl.Bool:
		return "yardl.Bool"
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
		return "yardl.Date"
	case dsl.Time:
		return "yardl.Time"
	case dsl.DateTime:
		return "yardl.DateTime"
	default:
		panic(fmt.Sprintf("primitive '%v' not yet supported", p))
	}
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
		return fmt.Sprintf("%s_serializer.overall_dtype()", formatting.ToSnakeCase(t.Name))
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func TypeDTypeTypeArgument(t dsl.Type) string {
	dTypeSyntax := TypeDTypeSyntax(t)
	// check if this is a record
	// TODO: this is a hack
	if strings.HasSuffix(dTypeSyntax, ")") {
		return "np.void"
	}

	return dTypeSyntax
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
		w.WriteStringln(`"""`)
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

func WriteGeneratedFileHeader(w *formatting.IndentedWriter) {
	WriteComment(w, "This file was generated by the \"yardl\" tool. DO NOT EDIT.")
	w.WriteStringln("")
}

func WriteTypeVars(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	typeVars := make(map[string]any)
	for _, td := range ns.TypeDefinitions {
		for _, tp := range td.GetDefinitionMeta().TypeParameters {
			identifier := TypeIdentifierName(tp.Name)
			if _, ok := typeVars[identifier]; !ok {
				typeVars[identifier] = nil
				fmt.Fprintf(w, "%s = typing.TypeVar('%s')\n", identifier, identifier)
			}
		}
	}
	if len(typeVars) > 0 {
		w.WriteStringln("")
	}
}
