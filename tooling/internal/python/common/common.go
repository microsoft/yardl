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

func TypeSyntax(t dsl.Type, contextNamespace string) string {
	switch t := t.(type) {
	case nil:
		return "None"
	case *dsl.SimpleType:
		return TypeDefinitionSyntax(t.ResolvedDefinition, contextNamespace)
	case *dsl.GeneralizedType:
		scalarString := func() string {
			if t.Cases.IsSingle() {
				return TypeSyntax(t.Cases[0].Type, contextNamespace)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("%s | None", TypeSyntax(t.Cases[1].Type, contextNamespace))
			}

			typeMap := make(map[string]any)

			caseStrings := make([]string, 0, len(t.Cases))
			for _, typeCase := range t.Cases {
				if typeCase.Type == nil {
					continue
				}

				syntax := TypeSyntax(typeCase.Type, contextNamespace)

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
			return "np.array"
		case *dsl.Map:
			return fmt.Sprintf("dict[%s, %s]", TypeSyntax(d.KeyType, contextNamespace), scalarString)
		default:
			panic(fmt.Sprintf("unexpected type %T", d))
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func TypeDefinitionSyntax(t dsl.TypeDefinition, contextNamespace string) string {
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

		if len(meta.TypeArguments) == 0 {
			return typeName
		}

		typeArguments := make([]string, 0, len(meta.TypeArguments))
		for _, typeArgument := range meta.TypeArguments {
			typeArguments = append(typeArguments, TypeSyntax(typeArgument, contextNamespace))
		}

		return fmt.Sprintf("%s[%s]", typeName, strings.Join(typeArguments, ", "))
	}
}

func PrimitiveSyntax(p dsl.PrimitiveDefinition) string {
	switch p {
	case dsl.Int8, dsl.Uint8, dsl.Int16, dsl.Uint16, dsl.Int32, dsl.Uint32, dsl.Int64, dsl.Uint64, dsl.Size:
		return "int"
	case dsl.Float32, dsl.Float64:
		return "float"
	case dsl.ComplexFloat32, dsl.ComplexFloat64:
		return "complex"
	case dsl.Bool:
		return "bool"
	case dsl.String:
		return "str"
	case dsl.Date:
		return "datetime.date"
	case dsl.Time:
		return "datetime.time"
	case dsl.DateTime:
		return "datetime.datetime"
	default:
		panic(fmt.Sprintf("primitive '%v' not yet supported", p))
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
		w.WriteStringln("")
	}
}

func WriteGeneratedFileHeader(w *formatting.IndentedWriter) {
	WriteComment(w, "This file was generated by the \"yardl\" tool. DO NOT EDIT.")
	w.WriteStringln("")
}
