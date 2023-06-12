// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package hdf5

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func collectUnionArities(env *dsl.Environment) []int {
	arities := make(map[int]any)

	dsl.Visit(env, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case *dsl.GeneralizedType:
			if _, isStream := t.Dimensionality.(*dsl.Stream); !isStream {
				if t.Cases[0].IsNullType() {
					if len(t.Cases) > 2 {
						arities[len(t.Cases)-1] = nil
					}
				} else if len(t.Cases) > 1 {
					arities[len(t.Cases)] = nil
				}
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

func writeInnerUnionTypes(w *formatting.IndentedWriter, env *dsl.Environment) {
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

		outerTypeArguments := func() string {
			args := make([]string, len(elements))
			for i := 0; i < arity; i++ {
				args[i] = fmt.Sprintf("TOuter%d", i)
			}
			return strings.Join(args, ", ")
		}()

		w.WriteString("template <")
		formatting.Delimited(w, ", ", elements, func(w *formatting.IndentedWriter, i int, item int) {
			fmt.Fprintf(w, "typename TInner%d, typename TOuter%d", i, i)
		})
		w.WriteString(">\n")
		fmt.Fprintf(w, "class InnerUnion%d {\n", arity)
		w.Indented(func() {
			w.WriteStringln("public:")
			fmt.Fprintf(w, "InnerUnion%d() : type_index_(-1) {} \n", arity)
			fmt.Fprintf(w, "InnerUnion%d(std::variant<%s> const& v) : type_index_(static_cast<int8_t>(v.index())) {\n", arity, outerTypeArguments)
			w.Indented(func() {
				w.WriteStringln("Init(v);")
			})
			w.WriteString("}\n\n")

			fmt.Fprintf(w, "InnerUnion%d(std::variant<std::monostate, %s> const& v) : type_index_(static_cast<int8_t>(v.index()) - 1) {\n", arity, outerTypeArguments)
			w.Indented(func() {
				w.WriteStringln("Init(v);")
			})
			w.WriteString("}\n\n")

			fmt.Fprintf(w, "InnerUnion%d(InnerUnion%d const& v) = delete;\n\n", arity, arity)
			fmt.Fprintf(w, "InnerUnion%d operator=(InnerUnion%d const&) = delete;\n\n", arity, arity)

			fmt.Fprintf(w, "~InnerUnion%d() {\n", arity)
			w.Indented(func() {
				w.WriteString("switch (type_index_) {\n")
				for i := 0; i < arity; i++ {
					fmt.Fprintf(w, "case %d:\n", i)
					w.Indented(func() {
						fmt.Fprintf(w, "value%d_.~TInner%d();\n", i, i)
						w.WriteStringln("break;")
					})
				}
				w.WriteStringln("}")
			})

			w.WriteString("}\n\n")

			fmt.Fprintf(w, "void ToOuter(std::variant<%s>& o) const {\n", outerTypeArguments)
			w.Indented(func() {
				w.WriteStringln("ToOuterImpl(o);")
			})
			w.WriteString("}\n\n")

			fmt.Fprintf(w, "void ToOuter(std::variant<std::monostate, %s>& o) const {\n", outerTypeArguments)
			w.Indented(func() {
				w.WriteStringln("ToOuterImpl(o);")
			})
			w.WriteString("}\n\n")

			w.WriteStringln("int8_t type_index_;")
			for i := 0; i < arity; i++ {
				w.WriteStringln("union {")
				w.Indented(func() {
					fmt.Fprintf(w, "char empty%d_[sizeof(TInner%d)]{};\n", i, i)
					fmt.Fprintf(w, "TInner%d value%d_;\n", i, i)
				})
				w.WriteString("};\n")
			}

			w.WriteString("\nprivate:\n")

			w.WriteStringln("template <typename T>")
			w.WriteStringln("void Init(T const& v) {")
			w.Indented(func() {
				w.WriteStringln("constexpr size_t offset = GetOuterVariantOffset<std::remove_const_t<std::remove_reference_t<decltype(v)>>>();")
				w.WriteString("switch (type_index_) {\n")
				for i := 0; i < arity; i++ {
					fmt.Fprintf(w, "case %d:\n", i)
					w.Indented(func() {
						fmt.Fprintf(w, "new (&value%d_) TInner%d(std::get<%d + offset>(v));\n", i, i, i)
						w.WriteStringln("return;")
					})
				}
				w.WriteStringln("}")
			})
			w.WriteString("}\n\n")

			w.WriteStringln("template <typename TVariant>")
			w.WriteStringln("void ToOuterImpl(TVariant& o) const {")
			w.Indented(func() {
				w.WriteStringln("constexpr size_t offset = GetOuterVariantOffset<TVariant>();")
				w.WriteString("switch (type_index_) {\n")
				w.WriteStringln(`case -1:
  if constexpr (offset == 1) {
    o.template emplace<0>(std::monostate{});
    return;
  }`)
				for i := 0; i < arity; i++ {
					fmt.Fprintf(w, "case %d:\n", i)
					w.Indented(func() {
						fmt.Fprintf(w, "o.template emplace<%d + offset>();\n", i)
						fmt.Fprintf(w, "yardl::hdf5::ToOuter(value%d_, std::get<%d + offset>(o));\n", i, i)
						w.WriteStringln("return;")
					})
				}
				w.WriteStringln("}")
				w.WriteStringln("throw std::runtime_error(\"unrecognized type variant type index \" + std::to_string(type_index_));")
			})

			w.WriteString("}\n\n")

			w.WriteStringln(`template <typename TVariant>
static constexpr size_t GetOuterVariantOffset() {
  constexpr bool has_monostate = std::is_same_v<std::monostate, std::variant_alternative_t<0, TVariant>>;
  if constexpr (has_monostate) {
    return 1;
  }
    return 0;
}`)
		})
		w.WriteString("};\n\n")

		// DDL function
		w.WriteString("template <")
		formatting.Delimited(w, ", ", elements, func(w *formatting.IndentedWriter, i int, item int) {
			fmt.Fprintf(w, "typename TInner%d, typename TOuter%d", i, i)
		})
		w.WriteString(">\n")
		fmt.Fprintf(w, "H5::CompType InnerUnion%dDdl(bool nullable, ", arity)
		formatting.Delimited(w, ", ", elements, func(w *formatting.IndentedWriter, i int, item int) {
			fmt.Fprintf(w, "H5::DataType const& t%d, ", i)
			fmt.Fprintf(w, "std::string const& label%d", i)
		})
		w.WriteString(") {\n")
		w.Indented(func() {
			innerTypeName := func() string {
				args := make([]string, 2*arity)
				for i := 0; i < arity; i++ {
					args[2*i] = fmt.Sprintf("TInner%d", i)
					args[2*i+1] = fmt.Sprintf("TOuter%d", i)
				}
				return fmt.Sprintf("::InnerUnion%d<%s>", arity, strings.Join(args, ", "))
			}()
			fmt.Fprintf(w, "using UnionType = %s;\n", innerTypeName)
			w.WriteStringln("H5::CompType rtn(sizeof(UnionType));")
			labels := make([]string, arity)
			for i := 0; i < arity; i++ {
				labels[i] = fmt.Sprintf("label%d", i)
			}

			fmt.Fprintf(w, "rtn.insertMember(\"$type\", HOFFSET(UnionType, type_index_), yardl::hdf5::UnionTypeEnumDdl(nullable, %s));\n", strings.Join(labels, ", "))

			for i := 0; i < arity; i++ {
				fmt.Fprintf(w, "rtn.insertMember(label%d, HOFFSET(UnionType, value%d_), t%d);\n", i, i, i)
			}

			w.WriteStringln("return rtn;")
		})
		w.WriteString("}\n")
	})

	w.WriteString("}\n\n")
}

func writeInnerType(w *formatting.IndentedWriter, recordDef *dsl.RecordDefinition) {
	if !needsInnerType(recordDef) {
		return
	}

	writeInnerDefinitionTemplateSpec(w, recordDef)
	innerName := innerTypeName(recordDef)
	fmt.Fprintf(w, "struct %s {\n", innerName)

	outerTypeSyntax := common.TypeDefinitionSyntax(recordDef)
	w.Indented(func() {
		fmt.Fprintf(w, "%s() {} \n", innerName)
		fmt.Fprintf(w, "%s(%s const& o) \n", innerName, outerTypeSyntax)
		w.Indented(func() {
			w.Indented(func() {
				w.WriteString(": ")
				formatting.Delimited(w, ",\n", recordDef.Fields, func(w *formatting.IndentedWriter, i int, f *dsl.Field) {
					fmt.Fprintf(w, "%s(o.%s)", common.FieldIdentifierName(f.Name), common.FieldIdentifierName(f.Name))
				})
				w.WriteString(" {\n")
			})
		})
		w.WriteString("}\n\n")

		fmt.Fprintf(w, "void ToOuter (%s& o) const {\n", outerTypeSyntax)
		w.Indented(func() {
			for _, f := range recordDef.Fields {
				fieldName := common.FieldIdentifierName(f.Name)
				fmt.Fprintf(w, "yardl::hdf5::ToOuter(%s, o.%s);\n", fieldName, fieldName)
			}
		})

		w.WriteString("}\n\n")

		for _, f := range recordDef.Fields {
			fmt.Fprintf(w, "%s %s;\n", innerTypeSyntax(f.Type), common.FieldIdentifierName(f.Name))
		}
	})
	fmt.Fprint(w, "};\n\n")
}

func needsInnerType(node dsl.Node) bool {
	result := false
	dsl.Visit(node, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case *dsl.SimpleType:
			self.Visit(t.ResolvedDefinition)
			self.VisitChildren(node)
		case *dsl.GeneralizedType:
			if len(t.Cases) > 1 {
				result = true
				return
			}
			switch d := t.Dimensionality.(type) {
			case *dsl.Vector:
				if !d.IsFixed() {
					result = true
					return
				}
			case *dsl.Array:
				if !d.IsFixed() {
					result = true
					return
				}
			}
			self.VisitChildren(node)
		case dsl.PrimitiveDefinition:
			if t == dsl.String {
				result = true
			}
		case *dsl.GenericTypeParameter:
			result = true
		default:
			self.VisitChildren(node)
		}
	})
	return result
}

func containsVlen(node dsl.Node) bool {
	result := false
	dsl.Visit(node, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case *dsl.SimpleType:
			self.Visit(t.ResolvedDefinition)
			self.VisitChildren(node)
		case *dsl.GeneralizedType:
			switch d := t.Dimensionality.(type) {
			case *dsl.Vector:
				if !d.IsFixed() {
					result = true
					return
				}
			case *dsl.Array:
				if !d.IsFixed() {
					result = true
					return
				}
			}
			self.VisitChildren(node)
		case dsl.PrimitiveDefinition:
			if t == dsl.String {
				result = true
			}
		default:
			self.VisitChildren(node)
		}
	})

	return result
}

func innerTypeSyntax(t dsl.Type) string {
	switch t := t.(type) {
	case nil:
		return common.TypeSyntax(t)
	case *dsl.SimpleType:
		return innerTypeDefinitionSyntax(t.ResolvedDefinition)
	case *dsl.GeneralizedType:
		scalarTOuterSyntax := common.TypeSyntax(t.ToScalar())
		scalarTInnerSyntax := func() string {
			if t.Cases.IsSingle() {
				return innerTypeSyntax(t.Cases[0].Type)
			}

			if t.Cases.IsOptional() {
				innerTypeSyntax := innerTypeSyntax(t.Cases[1].Type)
				return fmt.Sprintf("yardl::hdf5::InnerOptional<%s, %s>", innerTypeSyntax, common.TypeSyntax(t.Cases[1].Type))
			}

			caseStrings := make([]string, 0)
			for _, typeCase := range t.Cases {
				if !typeCase.IsNullType() {
					caseStrings = append(caseStrings, innerTypeSyntax(typeCase.Type), common.TypeSyntax(typeCase.Type))
				}
			}

			return fmt.Sprintf("::InnerUnion%d<%s>", len(caseStrings)/2, strings.Join(caseStrings, ", "))
		}()

		switch d := t.Dimensionality.(type) {
		case nil, *dsl.Stream:
			return scalarTInnerSyntax
		case *dsl.Vector:
			if !d.IsFixed() {
				return fmt.Sprintf("yardl::hdf5::InnerVlen<%s, %s>", scalarTInnerSyntax, scalarTOuterSyntax)
			}
			if needsInnerType(t.ToScalar()) {
				return fmt.Sprintf("yardl::hdf5::InnerFixedVector<%s, %s, %d>", scalarTInnerSyntax, scalarTOuterSyntax, *d.Length)
			}
			return common.TypeSyntax(t)
		case *dsl.Array:
			if !d.HasKnownNumberOfDimensions() {
				return fmt.Sprintf("yardl::hdf5::InnerDynamicNdArray<%s, %s>", scalarTInnerSyntax, scalarTOuterSyntax)
			}

			if d.IsFixed() {
				if needsInnerType(t.ToScalar()) {
					maxes := make([]string, len(*d.Dimensions))
					for i, dim := range *d.Dimensions {
						maxes[i] = fmt.Sprint(*dim.Length)
					}
					return fmt.Sprintf("yardl::hdf5::InnerFixedNdArray<%s, %s, %s>", scalarTInnerSyntax, scalarTOuterSyntax, strings.Join(maxes, ", "))
				}
				return common.TypeSyntax(t)
			}
			if len(*d.Dimensions) == 1 {
				return fmt.Sprintf("yardl::hdf5::InnerVlen<%s, %s>", scalarTInnerSyntax, scalarTOuterSyntax)
			}
			return fmt.Sprintf("yardl::hdf5::InnerNdArray<%s, %s, %d>", scalarTInnerSyntax, scalarTOuterSyntax, len(*d.Dimensions))
		case *dsl.Map:
			return fmt.Sprintf("yardl::hdf5::InnerMap<%s, %s, %s, %s>", innerTypeSyntax(d.KeyType), common.TypeSyntax(d.KeyType), scalarTInnerSyntax, scalarTOuterSyntax)
		default:
			panic(fmt.Sprintf("unexpected type %T", d))
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func innerTypeDefinitionSyntax(t dsl.TypeDefinition) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		switch string(t) {
		case "string":
			return "yardl::hdf5::InnerVlenString"
		default:
			return common.PrimitiveSyntax(t)
		}
	case *dsl.NamedType:
		return innerTypeSyntax(t.Type)
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("_%s_Inner", common.TypeDefinitionSyntax(t))
	default:
		meta := t.GetDefinitionMeta()
		convertedNamespace := common.NamespaceIdentifierName(meta.Namespace)
		typeArgsSpec := ""
		needsInnerType := needsInnerType(t)
		if len(meta.TypeParameters) > 0 {
			typeArgSpecs := make([]string, 0)
			if len(meta.TypeArguments) > 0 {
				for _, typeArg := range meta.TypeArguments {
					if needsInnerType {
						typeArgSpecs = append(typeArgSpecs, innerTypeSyntax(typeArg))
					}
					typeArgSpecs = append(typeArgSpecs, common.TypeSyntax(typeArg))
				}
			} else {
				for _, typeParam := range meta.TypeParameters {
					if needsInnerType {
						typeArgSpecs = append(typeArgSpecs, innerTypeDefinitionSyntax(typeParam))
					}
					typeArgSpecs = append(typeArgSpecs, common.TypeDefinitionSyntax(typeParam))
				}
			}

			typeArgsSpec = fmt.Sprintf("<%s>", strings.Join(typeArgSpecs, ", "))
		}

		if needsInnerType {
			return fmt.Sprintf("%s::hdf5::%s%s", convertedNamespace, innerTypeName(t), typeArgsSpec)
		} else {
			return fmt.Sprintf("%s::%s%s", convertedNamespace, common.TypeIdentifierName(meta.Name), typeArgsSpec)
		}
	}
}

func innerTypeName(t dsl.TypeDefinition) string {
	return fmt.Sprintf("_Inner_%s", common.TypeIdentifierName(t.GetDefinitionMeta().Name))
}

func writeRecordDdlFunction(w *formatting.IndentedWriter, rec *dsl.RecordDefinition) {
	writeInnerDefinitionTemplateSpec(w, rec)

	fmt.Fprintf(w, "[[maybe_unused]] H5::CompType %s(", hdf5DdlFunctionName(rec))
	formatting.Delimited(w, ", ", rec.DefinitionMeta.TypeParameters, func(w *formatting.IndentedWriter, i int, item *dsl.GenericTypeParameter) {
		fmt.Fprintf(w, "H5::DataType const& %s_type", item.Name)
	})

	w.WriteString(") {\n")
	w.Indented(func() {
		cppTypeSyntax := innerTypeDefinitionSyntax(rec)
		fmt.Fprintf(w, "using RecordType = %s;\n", cppTypeSyntax)
		w.WriteStringln("H5::CompType t(sizeof(RecordType));")

		for _, f := range rec.Fields {
			fmt.Fprintf(w, "t.insertMember(\"%s\", HOFFSET(RecordType, %s), %s);\n", f.Name, common.FieldIdentifierName(f.Name), typeDdlExpression(f.Type))
		}

		w.WriteStringln("return t;")
	})
	w.WriteString("}\n\n")
}

func writeEnumDdlFunction(w *formatting.IndentedWriter, e *dsl.EnumDefinition) {
	var baseType dsl.Type
	if e.BaseType != nil {
		baseType = e.BaseType
	} else {
		baseType = dsl.Int32Type
	}

	baseDdlExpression := typeDdlExpression(baseType)
	baseCppDataType := common.TypeSyntax(baseType)

	fmt.Fprintf(w, "[[maybe_unused]] H5::EnumType %s() {\n", hdf5DdlFunctionName(e))
	w.Indented(func() {
		fmt.Fprintf(w, "H5::EnumType t(%s);\n", baseDdlExpression)
		for i, enumValue := range e.Values {
			if i == 0 {
				fmt.Fprintf(w, "%s i = %s;\n", baseCppDataType, common.EnumIntegerLiteral(e, enumValue))
			} else {
				fmt.Fprintf(w, "i = %s;\n", common.EnumIntegerLiteral(e, enumValue))
			}
			fmt.Fprintf(w, "t.insert(\"%s\", &i);\n", enumValue.Symbol)
		}
		w.WriteString("return t;\n")
	})
	w.WriteString("}\n\n")
}

func typeDdlExpression(t dsl.Type) string {
	switch t := t.(type) {
	case *dsl.SimpleType:
		return typeDefinitionDdlExpression(t.ResolvedDefinition)
	case *dsl.GeneralizedType:
		scalarDdl := func() string {
			if t.Cases.IsSingle() {
				return typeDdlExpression(t.Cases[0].Type)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("yardl::hdf5::OptionalTypeDdl<%s, %s>(%s)", innerTypeSyntax(t.Cases[1].Type), common.TypeSyntax(t.Cases[1].Type), typeDdlExpression(t.Cases[1].Type))
			}

			templateArguments := make([]string, 0)
			for _, typeCase := range t.Cases {
				if !typeCase.IsNullType() {
					templateArguments = append(templateArguments, innerTypeSyntax(typeCase.Type), common.TypeSyntax(typeCase.Type))
				}
			}

			arguments := make([]string, 0)
			for _, typeCase := range t.Cases {
				if !typeCase.IsNullType() {
					arguments = append(arguments, typeDdlExpression(typeCase.Type))
					arguments = append(arguments, fmt.Sprintf("\"%s\"", typeCase.Label))
				}
			}

			return fmt.Sprintf(
				"::InnerUnion%dDdl<%s>(%t, %s)",
				len(templateArguments)/2,
				strings.Join(templateArguments, ", "),
				t.Cases[0].IsNullType(),
				strings.Join(arguments, ", "))
		}()

		switch d := t.Dimensionality.(type) {
		case nil, *dsl.Stream:
			return scalarDdl
		case *dsl.Vector:
			if !d.IsFixed() {
				return fmt.Sprintf("yardl::hdf5::InnerVlenDdl(%s)", scalarDdl)
			}

			return fmt.Sprintf("yardl::hdf5::FixedVectorDdl(%s, %d)", scalarDdl, *d.Length)
		case *dsl.Array:
			scalarT := t.ToScalar()
			if !d.HasKnownNumberOfDimensions() {
				return fmt.Sprintf("yardl::hdf5::DynamicNDArrayDdl<%s, %s>(%s)", innerTypeSyntax(scalarT), common.TypeSyntax(scalarT), scalarDdl)
			}

			if d.IsFixed() {
				maxes := make([]string, len(*d.Dimensions))
				for i, dim := range *d.Dimensions {
					maxes[i] = fmt.Sprint(*dim.Length)
				}
				return fmt.Sprintf("yardl::hdf5::FixedNDArrayDdl(%s, {%s})", typeDdlExpression(scalarT), strings.Join(maxes, ", "))
			}

			if len(*d.Dimensions) == 1 {
				return fmt.Sprintf("yardl::hdf5::InnerVlenDdl(%s)", typeDdlExpression(scalarT))
			}

			return fmt.Sprintf("yardl::hdf5::NDArrayDdl<%s, %s, %d>(%s)", innerTypeSyntax(scalarT), common.TypeSyntax(scalarT), len(*d.Dimensions), scalarDdl)
		case *dsl.Map:
			return fmt.Sprintf("yardl::hdf5::InnerMapDdl<%s, %s>(%s, %s)", innerTypeSyntax(d.KeyType), innerTypeSyntax(t.ToScalar()), typeDdlExpression(d.KeyType), scalarDdl)
		default:
			panic(fmt.Sprintf("unexpected type %T", d))
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func typeDefinitionDdlExpression(t dsl.TypeDefinition) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		switch t {
		case dsl.Int8, dsl.Uint8, dsl.Int16, dsl.Uint16, dsl.Int32, dsl.Uint32, dsl.Int64, dsl.Uint64:
			return "H5::PredType::NATIVE_" + strings.ToUpper(string(t))
		case dsl.Size:
			return "yardl::hdf5::SizeTypeDdl()"
		case dsl.Bool:
			return "H5::PredType::NATIVE_HBOOL"
		case dsl.Float32:
			return "H5::PredType::NATIVE_FLOAT"
		case dsl.Float64:
			return "H5::PredType::NATIVE_DOUBLE"
		case dsl.ComplexFloat32:
			return "yardl::hdf5::ComplexTypeDdl<float>()"
		case dsl.ComplexFloat64:
			return "yardl::hdf5::ComplexTypeDdl<double>()"
		case dsl.String:
			return "yardl::hdf5::InnerVlenStringDdl()"
		case dsl.Date:
			return "yardl::hdf5::DateTypeDdl()"
		case dsl.Time:
			return "yardl::hdf5::TimeTypeDdl()"
		case dsl.DateTime:
			return "yardl::hdf5::DateTimeTypeDdl()"
		default:
			log.Panicf("primitive '%v' not yet supported", t)
		}
	case *dsl.NamedType:
		return typeDdlExpression(t.Type)
	case *dsl.EnumDefinition:
		if !t.IsFlags {
			break
		}

		var baseType dsl.Type
		if t.BaseType != nil {
			baseType = t.BaseType
		} else {
			baseType = dsl.Int32Type
		}
		return typeDdlExpression(baseType)
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("%s_type", t.Name)
	}

	typeArgumentsString := ""
	typeDdlArgumentString := ""
	meta := t.GetDefinitionMeta()
	if len(meta.TypeParameters) > 0 {
		typeArguments := make([]string, 2*len(meta.TypeParameters))
		typeDdlExpressions := make([]string, len(meta.TypeParameters))

		if len(meta.TypeArguments) > 0 {
			for i, typeArg := range meta.TypeArguments {
				typeArguments[2*i] = innerTypeSyntax(typeArg)
				typeArguments[2*i+1] = common.TypeSyntax(typeArg)
				typeDdlExpressions[i] = typeDdlExpression(typeArg)
			}
		} else {
			for i, typeArg := range meta.TypeParameters {
				typeArguments[2*i] = innerTypeDefinitionSyntax(typeArg)
				typeArguments[2*i+1] = common.TypeDefinitionSyntax(typeArg)
				typeDdlExpressions[i] = typeDefinitionDdlExpression(typeArg)
			}
		}

		typeArgumentsString = fmt.Sprintf("<%s>", strings.Join(typeArguments, ", "))
		typeDdlArgumentString = strings.Join(typeDdlExpressions, ", ")
	}

	return fmt.Sprintf("%s%s(%s)", qualifiedH5DdlFunctionName(t), typeArgumentsString, typeDdlArgumentString)
}

func hdf5DdlFunctionName(t dsl.TypeDefinition) string {
	return fmt.Sprintf("Get%sHdf5Ddl", t.GetDefinitionMeta().Name)
}

func qualifiedH5DdlFunctionName(t dsl.TypeDefinition) string {
	return fmt.Sprintf("%s::hdf5::%s", common.TypeNamespaceIdentifierName(t), hdf5DdlFunctionName(t))
}

func writeInnerDefinitionTemplateSpec(w *formatting.IndentedWriter, td dsl.TypeDefinition) {
	meta := td.GetDefinitionMeta()
	if len(meta.TypeParameters) > 0 {
		templateParameters := make([]string, 2*len(meta.TypeParameters))
		for i, p := range meta.TypeParameters {
			templateParameters[2*i] = "typename " + fmt.Sprintf("_%s_Inner", common.TypeDefinitionSyntax(p))
			templateParameters[2*i+1] = "typename " + common.TypeDefinitionSyntax(p)
		}
		fmt.Fprintf(w, "template <%s>\n", strings.Join(templateParameters, ", "))
	}
}
