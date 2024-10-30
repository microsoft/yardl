// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package common

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

var reservedNames = map[string]any{
	"__has_cpp_attribute":      nil,
	"__has_include":            nil,
	"_Pragma":                  nil,
	"alignas":                  nil,
	"alignof":                  nil,
	"and":                      nil,
	"and_eq":                   nil,
	"asm":                      nil,
	"atomic_cancel":            nil,
	"atomic_commit":            nil,
	"atomic_noexcept":          nil,
	"auto":                     nil,
	"bitand":                   nil,
	"bitor":                    nil,
	"bool":                     nil,
	"break":                    nil,
	"case":                     nil,
	"catch":                    nil,
	"char":                     nil,
	"char16_t":                 nil,
	"char32_t":                 nil,
	"char8_t":                  nil,
	"class":                    nil,
	"co_await":                 nil,
	"co_return":                nil,
	"co_yield":                 nil,
	"compl":                    nil,
	"concept":                  nil,
	"const":                    nil,
	"const_cast":               nil,
	"consteval":                nil,
	"constexpr":                nil,
	"constinit":                nil,
	"continue":                 nil,
	"decltype":                 nil,
	"default":                  nil,
	"define":                   nil,
	"defined":                  nil,
	"delete":                   nil,
	"do":                       nil,
	"double":                   nil,
	"dynamic_cast":             nil,
	"elif":                     nil,
	"elifdef":                  nil,
	"elifndef":                 nil,
	"else":                     nil,
	"endif":                    nil,
	"enum":                     nil,
	"error":                    nil,
	"explicit":                 nil,
	"export":                   nil,
	"extern":                   nil,
	"false":                    nil,
	"final":                    nil,
	"float":                    nil,
	"for":                      nil,
	"friend":                   nil,
	"goto":                     nil,
	"if":                       nil,
	"ifdef":                    nil,
	"ifndef":                   nil,
	"import":                   nil,
	"include":                  nil,
	"inline":                   nil,
	"int":                      nil,
	"INT_FAST16_MAX":           nil,
	"int_fast16_t":             nil,
	"INT_FAST32_MAX":           nil,
	"int_fast32_t":             nil,
	"INT_FAST64_MAX":           nil,
	"int_fast64_t":             nil,
	"INT_FAST8_MAX":            nil,
	"int_fast8_t":              nil,
	"INT_FASTN_MAX":            nil,
	"INT_FASTN_MIN":            nil,
	"int_fastN_t":              nil,
	"INT_LEAST16_MAX":          nil,
	"int_least16_t":            nil,
	"INT_LEAST32_MAX":          nil,
	"int_least32_t":            nil,
	"INT_LEAST64_MAX":          nil,
	"int_least64_t":            nil,
	"INT_LEAST8_MAX":           nil,
	"int_least8_t":             nil,
	"INT_LEASTN_MAX":           nil,
	"INT_LEASTN_MIN":           nil,
	"int_leastN_t":             nil,
	"INT16_MAX":                nil,
	"int16_t":                  nil,
	"INT32_MAX":                nil,
	"int32_t":                  nil,
	"INT64_MAX":                nil,
	"int64_t":                  nil,
	"INT8_MAX":                 nil,
	"int8_t":                   nil,
	"INTMAX_C":                 nil,
	"INTMAX_MAX":               nil,
	"INTMAX_MIN":               nil,
	"intmax_t":                 nil,
	"INTN_C":                   nil,
	"INTN_MAX":                 nil,
	"INTN_MIN":                 nil,
	"intN_t":                   nil,
	"INTPTR_MAX":               nil,
	"INTPTR_MIN":               nil,
	"intptr_t":                 nil,
	"line":                     nil,
	"long":                     nil,
	"module":                   nil,
	"mutable":                  nil,
	"namespace":                nil,
	"new":                      nil,
	"noexcept":                 nil,
	"not":                      nil,
	"not_eq":                   nil,
	"nullptr":                  nil,
	"operator":                 nil,
	"or":                       nil,
	"or_eq":                    nil,
	"override":                 nil,
	"pragma":                   nil,
	"private":                  nil,
	"protected":                nil,
	"PTRDIFF_MAX":              nil,
	"PTRDIFF_MIN":              nil,
	"public":                   nil,
	"reflexpr":                 nil,
	"register":                 nil,
	"reinterpret_cast":         nil,
	"requires":                 nil,
	"return":                   nil,
	"short":                    nil,
	"SIG_ATOMIC_MAX":           nil,
	"SIG_ATOMIC_MIN":           nil,
	"signed":                   nil,
	"SIZE_MAX":                 nil,
	"sizeof":                   nil,
	"static":                   nil,
	"static_assert":            nil,
	"static_cast":              nil,
	"std":                      nil,
	"struct":                   nil,
	"switch":                   nil,
	"synchronized":             nil,
	"yardl":                    nil,
	"template":                 nil,
	"this":                     nil,
	"thread_local":             nil,
	"throw":                    nil,
	"transaction_safe":         nil,
	"transaction_safe_dynamic": nil,
	"true":                     nil,
	"try":                      nil,
	"typedef":                  nil,
	"typeid":                   nil,
	"typename":                 nil,
	"UINT_FAST16_MAX":          nil,
	"uint_fast16_t":            nil,
	"UINT_FAST32_MAX":          nil,
	"uint_fast32_t":            nil,
	"UINT_FAST64_MAX":          nil,
	"uint_fast64_t":            nil,
	"UINT_FAST8_MAX":           nil,
	"uint_fast8_t":             nil,
	"UINT_FASTN_MAX":           nil,
	"uint_fastN_t":             nil,
	"UINT_LEAST16_MAX":         nil,
	"uint_least16_t":           nil,
	"UINT_LEAST32_MAX":         nil,
	"uint_least32_t":           nil,
	"UINT_LEAST64_MAX":         nil,
	"uint_least64_t":           nil,
	"UINT_LEAST8_MAX":          nil,
	"uint_least8_t":            nil,
	"UINT_LEASTN_MAX":          nil,
	"uint_leastN_t":            nil,
	"UINT16_MAX":               nil,
	"uint16_t":                 nil,
	"UINT32_MAX":               nil,
	"uint32_t":                 nil,
	"UINT64_MAX":               nil,
	"uint64_t":                 nil,
	"UINT8_MAX":                nil,
	"uint8_t":                  nil,
	"UINTMAX_C":                nil,
	"UINTMAX_MAX":              nil,
	"uintmax_t":                nil,
	"UINTN_C":                  nil,
	"UINTN_MAX":                nil,
	"uintN_t":                  nil,
	"UINTPTR_MAX":              nil,
	"uintptr_t":                nil,
	"undef":                    nil,
	"union":                    nil,
	"unsigned":                 nil,
	"using":                    nil,
	"virtual":                  nil,
	"void":                     nil,
	"volatile":                 nil,
	"warning":                  nil,
	"WCHAR_MAX":                nil,
	"WCHAR_MIN":                nil,
	"wchar_t":                  nil,
	"while":                    nil,
	"WINT_MAX":                 nil,
	"WINT_MIN":                 nil,
	"xor":                      nil,
	"xor_eq":                   nil,
}

func TypeSyntax(t dsl.Type) string {
	switch t := t.(type) {
	case nil:
		return "std::monostate"
	case *dsl.SimpleType:
		return TypeDefinitionSyntax(t.ResolvedDefinition)
	case *dsl.GeneralizedType:
		scalarString := func() string {
			if t.Cases.IsSingle() {
				return TypeSyntax(t.Cases[0].Type)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("std::optional<%s>", TypeSyntax(t.Cases[1].Type))
			}

			caseStrings := make([]string, len(t.Cases))
			for i, typeCase := range t.Cases {
				if typeCase == nil {
					caseStrings[i] = "std::monostate"
				} else {
					caseStrings[i] = TypeSyntax(typeCase.Type)
				}
			}

			return fmt.Sprintf("std::variant<%s>", strings.Join(caseStrings, ", "))
		}()

		switch d := t.Dimensionality.(type) {
		case nil, *dsl.Stream:
			return scalarString
		case *dsl.Vector:
			if d.Length == nil {
				return fmt.Sprintf("std::vector<%s>", scalarString)
			}

			return fmt.Sprintf("std::array<%s, %d>", scalarString, *d.Length)
		case *dsl.Array:
			if !d.HasKnownNumberOfDimensions() {
				return fmt.Sprintf("yardl::DynamicNDArray<%s>", scalarString)
			}

			if !d.IsFixed() {
				return fmt.Sprintf("yardl::NDArray<%s, %d>", scalarString, len(*d.Dimensions))
			}

			dims := make([]string, len(*d.Dimensions))
			for i, dim := range *d.Dimensions {
				dims[i] = fmt.Sprint(*dim.Length)
			}

			return fmt.Sprintf("yardl::FixedNDArray<%s, %s>", scalarString, strings.Join(dims, ", "))
		case *dsl.Map:
			return fmt.Sprintf("std::unordered_map<%s, %s>", TypeSyntax(d.KeyType), scalarString)
		default:
			panic(fmt.Sprintf("unexpected type %T", d))
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func TypeDefinitionSyntax(t dsl.TypeDefinition) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		return PrimitiveSyntax(t)
	case *dsl.GenericTypeParameter:
		return TypeIdentifierName(t.Name)
	default:
		typeWithoutGenericArgs := fmt.Sprintf("%s::%s", TypeNamespaceIdentifierName(t), TypeIdentifierName(t.GetDefinitionMeta().Name))
		meta := t.GetDefinitionMeta()
		if len(meta.TypeParameters) == 0 {
			return typeWithoutGenericArgs
		}
		typeArgs := make([]string, len(meta.TypeParameters))
		if len(meta.TypeArguments) > 0 {
			for i, typeArg := range meta.TypeArguments {
				typeArgs[i] = TypeSyntax(typeArg)
			}
		} else {
			for i, typeParam := range meta.TypeParameters {
				typeArgs[i] = TypeDefinitionSyntax(typeParam)
			}
		}

		return fmt.Sprintf("%s<%s>", typeWithoutGenericArgs, strings.Join(typeArgs, ", "))
	}
}

func PrimitiveSyntax(p dsl.PrimitiveDefinition) string {
	switch p {
	case dsl.Int8, dsl.Uint8, dsl.Int16, dsl.Uint16, dsl.Int32, dsl.Uint32, dsl.Int64, dsl.Uint64:
		return string(p) + "_t"
	case dsl.Size:
		return "yardl::Size"
	case dsl.Float32:
		return "float"
	case dsl.Float64:
		return "double"
	case dsl.ComplexFloat32:
		return "std::complex<float>"
	case dsl.ComplexFloat64:
		return "std::complex<double>"
	case dsl.Bool:
		return "bool"
	case dsl.String:
		return "std::string"
	case dsl.Date:
		return "yardl::Date"
	case dsl.Time:
		return "yardl::Time"
	case dsl.DateTime:
		return "yardl::DateTime"
	default:
		panic(fmt.Sprintf("primitive '%v' not yet supported", p))
	}
}

func WriteComment(w *formatting.IndentedWriter, comment string) {
	comment = strings.TrimSpace(comment)
	if comment != "" {
		w = formatting.NewIndentedWriter(w, "// ").Indent()
		w.WriteStringln(comment)
	}
}

func WriteDefinitionTemplateSpec(w *formatting.IndentedWriter, td dsl.TypeDefinition) {
	meta := td.GetDefinitionMeta()
	if len(meta.TypeParameters) > 0 {
		templateParameters := make([]string, len(meta.TypeParameters))
		for i, p := range meta.TypeParameters {
			templateParameters[i] = "typename " + TypeDefinitionSyntax(p)
		}
		fmt.Fprintf(w, "template <%s>\n", strings.Join(templateParameters, ", "))
	}
}

func NamespaceIdentifierName(namespace string) string {
	return formatting.ToSnakeCase(strings.ReplaceAll(namespace, ".", "::"))
}

func TypeNamespaceIdentifierName(t dsl.TypeDefinition) string {
	return NamespaceIdentifierName(t.GetDefinitionMeta().Namespace)
}

func FieldIdentifierName(name string) string {
	snakeCased := formatting.ToSnakeCase(name)
	if _, reserved := reservedNames[snakeCased]; !reserved {
		return snakeCased
	}

	return fmt.Sprintf("%s_field", snakeCased)
}

func EnumValueIdentifierName(name string) string {
	prefixed := fmt.Sprintf("k%s", formatting.ToPascalCase(name))
	// snakeCased := formatting.ToSnakeCase(name)
	if _, reserved := reservedNames[prefixed]; !reserved {
		return prefixed
	}

	return fmt.Sprintf("%s_value", prefixed)
}

func ComputedFieldIdentifierName(name string) string {
	pascalCased := formatting.ToPascalCase(name)
	if _, reserved := reservedNames[pascalCased]; !reserved {
		return pascalCased
	}

	return fmt.Sprintf("%s_field", pascalCased)
}

func TypeIdentifierName(name string) string {
	if _, reserved := reservedNames[name]; !reserved {
		return name
	}

	return fmt.Sprintf("%s_Type", name)
}

func AbstractWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sWriterBase", p.Name)
}

func QualifiedAbstractWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::%s", TypeNamespaceIdentifierName(p), AbstractWriterName(p))
}

func AbstractReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sReaderBase", p.Name)
}

func QualifiedAbstractReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::%s", TypeNamespaceIdentifierName(p), AbstractReaderName(p))
}

func AbstractIndexedReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sIndexedReaderBase", p.Name)
}

func QualifiedAbstractIndexedReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::%s", TypeNamespaceIdentifierName(p), AbstractIndexedReaderName(p))
}

func ProtocolWriteMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("Write%s", formatting.ToPascalCase(s.Name))
}

func ProtocolWriteImplMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("Write%sImpl", formatting.ToPascalCase(s.Name))
}

func ProtocolWriteEndMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("End%s", formatting.ToPascalCase(s.Name))
}

func ProtocolWriteEndImplMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("End%sImpl", formatting.ToPascalCase(s.Name))
}

func ProtocolReadMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("Read%s", formatting.ToPascalCase(s.Name))
}

func ProtocolReadImplMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("Read%sImpl", formatting.ToPascalCase(s.Name))
}

func ProtocolStreamSizeMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("Count%s", formatting.ToPascalCase(s.Name))
}

func ProtocolStreamSizeImplMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("Count%sImpl", formatting.ToPascalCase(s.Name))
}

func EnumIntegerLiteral(e *dsl.EnumDefinition, v *dsl.EnumValue) string {
	if e.BaseType == nil {
		return fmt.Sprintf("%d", &v.IntegerValue)
	}

	return IntegerLiteral(v.IntegerValue, e.BaseType)
}

func IntegerLiteral(value big.Int, literalType dsl.Type) string {
	primitive, ok := dsl.GetPrimitiveType(literalType)
	if !ok || !dsl.IsIntegralPrimitive(primitive) {
		panic(fmt.Sprintf("type %v is not an integral type", literalType))
	}

	literalSuffix := ""

	switch primitive {
	case dsl.Int64:
		literalSuffix = "LL"
	case dsl.Uint64, dsl.Size:
		literalSuffix = "ULL"
	}

	return fmt.Sprintf("%d%s", &value, literalSuffix)
}

func WriteGeneratedFileHeader(w *formatting.IndentedWriter) {
	WriteComment(w, "This file was generated by the \"yardl\" tool. DO NOT EDIT.")
	w.WriteStringln("")
}
