package types

import (
	"bytes"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/internal/python/common"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func WriteTypes(ns *dsl.Namespace, st dsl.SymbolTable, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)
	w.WriteStringln(`
import dataclasses
import datetime
import enum
import types
import typing
import numpy as np
import numpy.typing as npt
from . import yardl_types as yardl
from . import _dtypes
`)

	writeTypes(w, st, ns)

	writeGetDTypeFunc(w, ns)

	definitionsPath := path.Join(packageDir, "types.py")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeTypes(w *formatting.IndentedWriter, st dsl.SymbolTable, ns *dsl.Namespace) {
	common.WriteTypeVars(w, ns)

	for _, td := range ns.TypeDefinitions {
		switch td := td.(type) {
		case *dsl.EnumDefinition:
			writeEnum(w, td)
		case *dsl.RecordDefinition:
			writeRecord(w, td, st)
		case *dsl.NamedType:
			writeNamedType(w, td)
		default:
			panic(fmt.Sprintf("unsupported type definition: %T", td))
		}
	}
}

func writeNamedType(w *formatting.IndentedWriter, td *dsl.NamedType) {
	fmt.Fprintf(w, "%s = %s\n", common.TypeSyntaxWithoutTypeParameters(td, td.Namespace), common.TypeSyntax(td.Type, td.Namespace))
	common.WriteDocstring(w, td.Comment)
	w.Indent().WriteStringln("")
}

func writeRecord(w *formatting.IndentedWriter, rec *dsl.RecordDefinition, st dsl.SymbolTable) {
	w.WriteStringln("@dataclasses.dataclass(slots=True, kw_only=True)")
	fmt.Fprintf(w, "class %s%s:\n", common.TypeSyntaxWithoutTypeParameters(rec, rec.Namespace), GetGenericBase(rec))
	w.Indented(func() {
		common.WriteDocstring(w, rec.Comment)
		for _, field := range rec.Fields {
			fmt.Fprintf(w, "%s: %s", common.FieldIdentifierName(field.Name), common.TypeSyntax(field.Type, rec.Namespace))

			if dsl.ContainsGenericTypeParameter(field.Type) {
				// cannot default generic type parameters
				// because they don't really exist at runtime
				w.WriteStringln("")
				continue
			}

			defaultExpr, defaultKind := typeDefault(field.Type, rec.Namespace, st)
			if defaultKind == defaultValueKindNone || defaultExpr == "" {
				w.WriteStringln("")
			} else if defaultKind == defaultValueKindValue {
				fmt.Fprintf(w, " = %s\n", defaultExpr)
			} else if defaultKind == defaultValueKindFactory {
				fmt.Fprintf(w, " = dataclasses.field(default_factory=%s)\n", defaultExpr)
			} else if defaultKind == defaultValueKindLambda {
				fmt.Fprintf(w, " = dataclasses.field(default_factory=lambda: %s)\n", defaultExpr)
			}

			common.WriteDocstring(w, field.Comment)
			w.WriteStringln("")
		}

		for _, computedField := range rec.ComputedFields {
			expressionTypeSyntax := common.TypeSyntax(computedField.Expression.GetResolvedType(), rec.Namespace)
			fieldName := common.ComputedFieldIdentifierName(computedField.Name)
			fmt.Fprintf(w, "def %s(self) -> %s:\n", fieldName, expressionTypeSyntax)
			w.Indented(func() {
				common.WriteDocstring(w, computedField.Comment)
				writeComputedFieldExpression(w, computedField.Expression, rec.Namespace)
				w.WriteStringln("")
			})
		}

		if len(rec.Fields)+len(rec.ComputedFields) == 0 {
			w.WriteStringln("pass")
		}
	})
	w.WriteStringln("")
}

type tailHandler func(next func())

type tailWrapper struct {
	outer   *tailWrapper
	handler tailHandler
}

func (t tailWrapper) Append(handler tailHandler) tailWrapper {
	return tailWrapper{
		outer:   &t,
		handler: handler,
	}
}

func (t tailWrapper) Run(body func()) {
	t.composeFunc(body)()
}

func (t tailWrapper) composeFunc(next func()) func() {
	if t.handler == nil {
		return next
	}
	this := func() { t.handler(next) }
	if t.outer == nil {
		return this
	}

	return t.outer.composeFunc(this)
}

func writeComputedFieldExpression(w *formatting.IndentedWriter, expression dsl.Expression, contextNamespace string) {
	helperFunctionLookup := make(map[any]string)
	dsl.Visit(expression, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case *dsl.FunctionCallExpression:
			if t.FunctionName == dsl.FunctionDimensionIndex {
				arrType := (t.Arguments[0].GetResolvedType())
				if _, ok := helperFunctionLookup[arrType]; !ok {
					funcName := fmt.Sprintf("_helper_%d", len(helperFunctionLookup))
					helperFunctionLookup[arrType] = funcName
					fmt.Fprintf(w, "def %s(dim_name: str) -> int:\n", funcName)
					w.Indented(func() {
						dims := dsl.ToGeneralizedType(arrType).Dimensionality.(*dsl.Array).Dimensions
						for i, d := range *dims {
							fmt.Fprintf(w, "if dim_name == \"%s\":\n", *d.Name)
							w.Indented(func() {
								fmt.Fprintf(w, "return %d\n", i)
							})
						}
						fmt.Fprintf(w, "raise KeyError(f\"Unknown dimension name: '{dim_name}'\")\n")
						w.WriteStringln("")
					})
				}
			}
		}
		self.VisitChildren(node)
	})

	varCounter := 0
	newVarName := func() string {
		varName := fmt.Sprintf("_var%d", varCounter)
		varCounter++
		return varName
	}

	tail := tailWrapper{}.Append(func(next func()) {
		w.WriteString("return ")
		next()
		w.WriteStringln("")
	})

	dsl.VisitWithContext(expression, tail, func(self dsl.VisitorWithContext[tailWrapper], node dsl.Node, tail tailWrapper) {
		switch t := node.(type) {
		case *dsl.IntegerLiteralExpression:
			tail.Run(func() {
				fmt.Fprintf(w, "%d", &t.Value)
			})
		case *dsl.StringLiteralExpression:
			tail.Run(func() {
				fmt.Fprintf(w, "%q", t.Value)
			})
		case *dsl.MemberAccessExpression:
			tail.Run(func() {
				if t.Target == nil {
					if t.Kind == dsl.MemberAccessVariable {
						w.WriteString(common.FieldIdentifierName(t.Member))
						return
					}

					w.WriteString("self")
				} else {
					self.Visit(t.Target, tailWrapper{})
				}
				w.WriteString(".")
				if t.Kind == dsl.MemberAccessComputedField {
					fmt.Fprintf(w, "%s()", common.ComputedFieldIdentifierName(t.Member))
				} else {
					w.WriteString(common.FieldIdentifierName(t.Member))
				}
			})
		case *dsl.IndexExpression:
			tail.Run(func() {
				isTargetArray := false
				if t.Target != nil {
					if gt, ok := t.Target.GetResolvedType().(*dsl.GeneralizedType); ok {
						if _, ok := gt.Dimensionality.(*dsl.Array); ok {
							isTargetArray = true
						}
					}
				}
				if isTargetArray {
					// a cast is needed for numpy subscripting
					fmt.Fprintf(w, "typing.cast(%s, ", common.TypeSyntax(t.GetResolvedType(), contextNamespace))
				}

				self.Visit(t.Target, tailWrapper{})
				w.WriteString("[")
				formatting.Delimited(w, ", ", t.Arguments, func(w *formatting.IndentedWriter, i int, a *dsl.IndexArgument) {
					self.Visit(a.Value, tailWrapper{})
				})
				w.WriteString("]")

				if isTargetArray {
					w.WriteString(")")
				}
			})
		case *dsl.FunctionCallExpression:
			tail.Run(func() {
				switch t.FunctionName {
				case dsl.FunctionSize:
					switch dsl.ToGeneralizedType(dsl.GetUnderlyingType(t.Arguments[0].GetResolvedType())).Dimensionality.(type) {
					case *dsl.Vector, *dsl.Map:
						fmt.Fprintf(w, "len(")
						self.Visit(t.Arguments[0], tailWrapper{})
						fmt.Fprintf(w, ")")
					case *dsl.Array:
						self.Visit(t.Arguments[0], tailWrapper{})

						if len(t.Arguments) == 1 {
							fmt.Fprintf(w, ".size")
						} else {
							fmt.Fprintf(w, ".shape[")
							remainingArgs := t.Arguments[1:]
							formatting.Delimited(w, ", ", remainingArgs, func(w *formatting.IndentedWriter, i int, arg dsl.Expression) {
								self.Visit(arg, tailWrapper{})
							})
							fmt.Fprintf(w, "]")
						}
					}
				case dsl.FunctionDimensionIndex:
					helperFuncName := helperFunctionLookup[t.Arguments[0].GetResolvedType()]
					fmt.Fprintf(w, "%s(", helperFuncName)
					self.Visit(t.Arguments[1], tailWrapper{})
					w.WriteString(")")

				case dsl.FunctionDimensionCount:
					self.Visit(t.Arguments[0], tailWrapper{})
					fmt.Fprintf(w, ".ndim")
				default:
					panic(fmt.Sprintf("Unknown function '%s'", t.FunctionName))
				}
			})
		case *dsl.SwitchExpression:
			targetType := dsl.ToGeneralizedType(dsl.GetUnderlyingType(t.Target.GetResolvedType()))

			unionVariableName := newVarName()
			fmt.Fprintf(w, "%s = ", unionVariableName)
			self.Visit(t.Target, tailWrapper{})
			w.WriteStringln("")

			if targetType.Cases.IsOptional() {
				for i, switchCase := range t.Cases {
					writeSwitchCaseOverOptional(w, switchCase, unionVariableName, i == len(targetType.Cases)-1, self, tail)
				}
				return
			}

			if targetType.Cases.IsUnion() {
				for _, switchCase := range t.Cases {
					writeSwitchCaseOverUnion(w, targetType, switchCase, unionVariableName, self, tail)
				}

				fmt.Fprintf(w, "raise RuntimeError(\"Unexpected union case\")\n")
				return
			}

			// this is over a single type
			if len(t.Cases) != 1 {
				panic("switch expression over a single type expected to have exactly one case")
			}
			switchCase := t.Cases[0]
			switch pattern := switchCase.Pattern.(type) {
			case *dsl.DeclarationPattern:
				fmt.Fprintf(w, "%s = %s\n", common.FieldIdentifierName(pattern.Identifier), unionVariableName)
				self.Visit(switchCase.Expression, tail)
			case *dsl.TypePattern, *dsl.DiscardPattern:
				self.Visit(switchCase.Expression, tail)
			default:
				panic(fmt.Sprintf("Unexpected pattern type %T", t.Cases[0].Pattern))
			}

		case *dsl.TypeConversionExpression:
			tail = tail.Append(func(next func()) {
				fmt.Fprintf(w, "%s(", typeConversionCallable(t.Type))
				next()
				w.WriteString(")")
			})

			self.Visit(t.Expression, tail)
		}
	})
}

func writeSwitchCaseOverOptional(w *formatting.IndentedWriter, switchCase *dsl.SwitchCase, variableName string, isLastCase bool, visitor dsl.VisitorWithContext[tailWrapper], tail tailWrapper) {
	writeCore := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		if declarationIdentifier != "" {
			fmt.Fprintf(w, "%s = %s\n", declarationIdentifier, variableName)
		}
		visitor.Visit(switchCase.Expression, tail)
	}

	writeTypeCase := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		if isLastCase {
			writeCore(typePattern, declarationIdentifier)
			return
		}

		if typePattern.Type == nil {
			fmt.Fprintf(w, "if %s is None:\n", variableName)
		} else {
			fmt.Fprintf(w, "if %s is not None:\n", variableName)
		}

		w.Indented(func() {
			writeCore(typePattern, declarationIdentifier)
		})
	}

	switch t := switchCase.Pattern.(type) {
	case *dsl.TypePattern:
		writeTypeCase(t, "")
	case *dsl.DeclarationPattern:
		writeTypeCase(&t.TypePattern, common.FieldIdentifierName(t.Identifier))
	case *dsl.DiscardPattern:
		writeCore(nil, "")
	default:
		panic(fmt.Sprintf("Unknown pattern type '%T'", switchCase.Pattern))
	}
}

func writeSwitchCaseOverUnion(w *formatting.IndentedWriter, unionType *dsl.GeneralizedType, switchCase *dsl.SwitchCase, variableName string, visitor dsl.VisitorWithContext[tailWrapper], tail tailWrapper) {
	writeTypeCase := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		for _, typeCase := range unionType.Cases {
			if dsl.TypesEqual(typePattern.Type, typeCase.Type) {
				if typePattern.Type == nil {
					fmt.Fprintf(w, "if %s is None:\n", variableName)
					w.Indented(func() {
						visitor.Visit(switchCase.Expression, tail)
					})
				} else {
					fmt.Fprintf(w, "if %s[0] == \"%s\":\n", variableName, typeCase.Tag)
					w.Indented(func() {
						if declarationIdentifier != "" {
							fmt.Fprintf(w, "%s = %s[1]\n", declarationIdentifier, variableName)
						}
						visitor.Visit(switchCase.Expression, tail)
					})
				}
				return
			}
		}
		panic(fmt.Sprintf("Did not find pattern type  '%s'", dsl.TypeToShortSyntax(typePattern.Type, false)))
	}

	switch t := switchCase.Pattern.(type) {
	case *dsl.TypePattern:
		writeTypeCase(t, "")
	case *dsl.DeclarationPattern:
		writeTypeCase(&t.TypePattern, common.FieldIdentifierName(t.Identifier))
	case *dsl.DiscardPattern:
		visitor.Visit(switchCase.Expression, tail)
	default:
		panic(fmt.Sprintf("Unknown pattern type '%T'", switchCase.Pattern))
	}
}

func typeConversionCallable(t dsl.Type) string {
	switch t := t.(type) {
	case *dsl.SimpleType:
		switch t := t.ResolvedDefinition.(type) {
		case dsl.PrimitiveDefinition:
			switch t {
			case dsl.Bool:
				return "bool"
			case dsl.Int8, dsl.Uint8, dsl.Int16, dsl.Uint16, dsl.Int32, dsl.Uint32, dsl.Int64, dsl.Uint64, dsl.Size:
				return "int"
			case dsl.Float32, dsl.Float64:
				return "float"
			case dsl.ComplexFloat32, dsl.ComplexFloat64:
				return "complex"
			case dsl.String:
				return "str"
			case dsl.Date:
				return "datetime.date"
			case dsl.Time:
				return "datetime.time"
			case dsl.DateTime:
				return "datetime.datetime"
			}
		}
	}
	panic(fmt.Sprintf("Unsupported type '%s'", t))
}

func GetGenericBase(t dsl.TypeDefinition) string {
	meta := t.GetDefinitionMeta()
	if len(meta.TypeParameters) == 0 {
		return ""
	}

	var typeParams []string
	for _, tp := range meta.TypeParameters {
		use := tp.Annotations[common.TypeParameterUseAnnotationKey].(common.TypeParameterUse)
		if use&common.TypeParameterUseScalar != 0 {
			typeParams = append(typeParams, common.TypeIdentifierName(tp.Name))
		}
		if use&common.TypeParameterUseArray != 0 {
			typeParams = append(typeParams, common.NumpyTypeParameterSyntax(tp))
		}
	}

	if len(typeParams) == 0 {
		return ""
	}

	return fmt.Sprintf("(typing.Generic[%s])", strings.Join(typeParams, ", "))
}

func writeEnum(w *formatting.IndentedWriter, enum *dsl.EnumDefinition) {
	var base string
	if enum.IsFlags {
		base = "enum.Flag, boundary=enum.KEEP"
	} else {
		base = "enum.Enum"
	}
	fmt.Fprintf(w, "class %s(%s):\n", common.TypeSyntaxWithoutTypeParameters(enum, enum.Namespace), base)

	w.Indented(func() {
		common.WriteDocstring(w, enum.Comment)
		for _, value := range enum.Values {
			fmt.Fprintf(w, "%s = %d\n", common.EnumValueIdentifierName(value.Symbol), &value.IntegerValue)
			common.WriteDocstring(w, value.Comment)
		}
	})
	w.WriteStringln("")
}

type defaultValueKind int

const (
	defaultValueKindNone defaultValueKind = iota
	defaultValueKindValue
	defaultValueKindFactory
	defaultValueKindLambda
)

func typeDefault(t dsl.Type, contextNamespace string, st dsl.SymbolTable) (string, defaultValueKind) {
	switch t := t.(type) {
	case nil:
		return "None", defaultValueKindValue
	case *dsl.SimpleType:
		return typeDefinitionDefault(t.ResolvedDefinition, contextNamespace, st)
	case *dsl.GeneralizedType:
		switch td := t.Dimensionality.(type) {
		case nil:
			defaultExpression, defaultKind := typeDefault(t.Cases[0].Type, contextNamespace, st)
			if t.Cases.IsSingle() || t.Cases.HasNullOption() {
				return defaultExpression, defaultKind
			}

			return fmt.Sprintf(`("%s", %s)`, t.Cases[0].Tag, defaultExpression), defaultValueKindValue
		case *dsl.Vector:
			if td.Length == nil {
				return "list", defaultValueKindFactory
			}

			scalarDefault, scalarDefaultKind := typeDefault(t.Cases[0].Type, contextNamespace, st)

			switch scalarDefaultKind {
			case defaultValueKindNone:
				return "", defaultValueKindNone
			case defaultValueKindValue:
				return fmt.Sprintf("[%s] * %d", scalarDefault, *td.Length), defaultValueKindLambda
			case defaultValueKindFactory:
				return fmt.Sprintf("[%s() for _ in range(%d)]", scalarDefault, *td.Length), defaultValueKindLambda
			case defaultValueKindLambda:
				return fmt.Sprintf("[%s for _ in range(%d)]", scalarDefault, *td.Length), defaultValueKindLambda
			}

		case *dsl.Array:
			context := dTypeExpressionContext{
				namespace: contextNamespace,
				root:      false,
			}
			dtype := typeDTypeExpression(t.ToScalar(), context)

			if td.IsFixed() {
				dims := make([]string, len(*td.Dimensions))
				for i, d := range *td.Dimensions {
					dims[i] = strconv.FormatUint(*d.Length, 10)
				}

				return fmt.Sprintf("np.zeros((%s,), dtype=%s)", strings.Join(dims, ", "), dtype), defaultValueKindLambda
			}

			if td.HasKnownNumberOfDimensions() {
				shape := fmt.Sprintf("(%s)", strings.Repeat("0,", len(*td.Dimensions))[0:len(*td.Dimensions)*2-1])
				return fmt.Sprintf("np.zeros(%s, dtype=%s)", shape, dtype), defaultValueKindLambda
			}

			return fmt.Sprintf("np.zeros((), dtype=%s)", dtype), defaultValueKindLambda

		case *dsl.Map:
			return "dict", defaultValueKindFactory
		}
	}

	return "", defaultValueKindNone
}

func typeDefinitionDefault(t dsl.TypeDefinition, contextNamespace string, st dsl.SymbolTable) (string, defaultValueKind) {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		switch t {
		case dsl.Bool:
			return "False", defaultValueKindValue
		case dsl.Int8, dsl.Uint8, dsl.Int16, dsl.Uint16, dsl.Int32, dsl.Uint32, dsl.Int64, dsl.Uint64, dsl.Size:
			return "0", defaultValueKindValue
		case dsl.Float32, dsl.Float64:
			return "0.0", defaultValueKindValue
		case dsl.ComplexFloat32, dsl.ComplexFloat64:
			return "0j", defaultValueKindValue
		case dsl.String:
			return `""`, defaultValueKindValue
		case dsl.Date:
			return "datetime.date(1970, 1, 1)", defaultValueKindLambda
		case dsl.Time:
			return "datetime.time(0, 0, 0)", defaultValueKindLambda
		case dsl.DateTime:
			return "datetime.datetime(1970, 1, 1, 0, 0, 0)", defaultValueKindLambda
		}
	case *dsl.EnumDefinition:
		zeroValue := t.GetZeroValue()
		if t.IsFlags {
			if zeroValue == nil {
				return fmt.Sprintf("%s(0)", common.TypeSyntax(t, contextNamespace)), defaultValueKindValue
			} else {
				return fmt.Sprintf("%s.%s", common.TypeSyntax(t, contextNamespace), common.EnumValueIdentifierName(zeroValue.Symbol)), defaultValueKindValue
			}
		}

		if zeroValue == nil {
			return "", defaultValueKindNone
		}

		return fmt.Sprintf("%s.%s", common.TypeSyntax(t, contextNamespace), common.EnumValueIdentifierName(zeroValue.Symbol)), defaultValueKindValue
	case *dsl.NamedType:
		return typeDefault(t.Type, contextNamespace, st)

	case *dsl.RecordDefinition:
		if len(t.TypeArguments) == 0 {
			if len(t.TypeParameters) > 0 {
				return "", defaultValueKindNone
			}

			for _, f := range t.Fields {
				_, fieldDefaultKind := typeDefault(f.Type, contextNamespace, st)
				if fieldDefaultKind == defaultValueKindNone {
					return "", defaultValueKindNone
				}
			}

			return common.TypeSyntaxWithoutTypeParameters(t, contextNamespace), defaultValueKindFactory
		}

		// generic record with type arguments

		genericDef := st[t.GetQualifiedName()].(*dsl.RecordDefinition)

		args := make([]string, 0)

		for i, f := range t.Fields {
			fieldDefaultExpr, fieldDefaultKind := typeDefault(f.Type, contextNamespace, st)
			if fieldDefaultKind == defaultValueKindNone {
				return "", defaultValueKindNone
			}

			_, genDefaultKind := typeDefault(genericDef.Fields[i].Type, contextNamespace, st)
			if genDefaultKind == defaultValueKindNone {
				switch fieldDefaultKind {
				case defaultValueKindValue:
					args = append(args, fmt.Sprintf("%s=%s", common.FieldIdentifierName(f.Name), fieldDefaultExpr))
				case defaultValueKindFactory:
					args = append(args, fmt.Sprintf("%s=%s()", common.FieldIdentifierName(f.Name), fieldDefaultExpr))
				case defaultValueKindLambda:
					args = append(args, fmt.Sprintf("%s=%s", common.FieldIdentifierName(f.Name), fieldDefaultExpr))
				}
			}
		}

		return fmt.Sprintf("%s(%s)", common.TypeSyntaxWithoutTypeParameters(t, contextNamespace), strings.Join(args, ", ")), defaultValueKindLambda
	}

	return "", defaultValueKindNone
}

func writeGetDTypeFunc(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	w.WriteStringln("def _mk_get_dtype():")
	w.Indented(func() {
		w.WriteStringln("dtype_map: dict[type | types.GenericAlias, np.dtype[typing.Any] | typing.Callable[[tuple[type, ...]], np.dtype[typing.Any]]] = {}")
		w.WriteStringln("get_dtype = _dtypes.make_get_dtype_func(dtype_map)\n")

		context := dTypeExpressionContext{
			namespace: ns.Name,
			root:      true,
		}

		for _, t := range ns.TypeDefinitions {
			fmt.Fprintf(w, "dtype_map[%s] = %s\n", common.TypeSyntaxWithoutTypeParameters(t, ns.Name), typeDefinitionDTypeExpression(t, context))
		}

		w.WriteStringln("\nreturn get_dtype")
	})
	w.WriteStringln("")

	w.WriteStringln("get_dtype = _mk_get_dtype()\n")
}

type dTypeExpressionContext struct {
	namespace            string
	root                 bool
	typeParameterIndexes map[*dsl.GenericTypeParameter]int
}

func typeDefinitionDTypeExpression(t dsl.TypeDefinition, context dTypeExpressionContext) string {
	if !context.root {
		switch t := t.(type) {
		case dsl.PrimitiveDefinition:
			switch t {
			case dsl.Bool:
				return "np.dtype(np.bool_)"
			case dsl.Int8, dsl.Uint8, dsl.Int16, dsl.Uint16, dsl.Int32, dsl.Uint32, dsl.Int64, dsl.Uint64, dsl.Float32, dsl.Float64:
				return fmt.Sprintf("np.dtype(np.%s)", strings.ToLower(string(t)))
			case dsl.Size:
				return "np.dtype(np.uint64)"
			case dsl.ComplexFloat32:
				return "np.dtype(np.complex64)"
			case dsl.ComplexFloat64:
				return "np.dtype(np.complex128)"
			case dsl.Date:
				return "np.dtype(np.datetime64)"
			case dsl.Time:
				return "np.dtype(np.timedelta64)"
			case dsl.DateTime:
				return "np.dtype(np.datetime64)"
			case dsl.String:
				return "np.dtype(np.object_)"
			default:
				panic(fmt.Sprintf("Not implemented %s", t))
			}
		case *dsl.GenericTypeParameter:
			index, ok := context.typeParameterIndexes[t]
			if !ok {
				panic("type parameter not found")
			}
			return fmt.Sprintf("get_dtype(type_args[%d])", index)
		}

		if len(t.GetDefinitionMeta().TypeParameters) > 0 {
			typeArgs := make([]string, 0)
			for _, ta := range t.GetDefinitionMeta().TypeArguments {
				typeArgs = append(typeArgs, getTypeSyntaxWithGenricArgsReadFromTupleArgs(ta, context))
			}

			return fmt.Sprintf("get_dtype(types.GenericAlias(%s, (%s,)))", common.TypeSyntaxWithoutTypeParameters(t, context.namespace), strings.Join(typeArgs, ", "))
		}

		return fmt.Sprintf("get_dtype(%s)", common.TypeSyntaxWithoutTypeParameters(t, context.namespace))
	}

	meta := t.GetDefinitionMeta()
	lambdaDeclaration := ""
	if len(meta.TypeParameters) > 0 {
		context.typeParameterIndexes = make(map[*dsl.GenericTypeParameter]int)
		for i, p := range meta.TypeParameters {
			context.typeParameterIndexes[p] = i
		}

		lambdaDeclaration = "lambda type_args: "
	}

	switch t := t.(type) {
	case *dsl.NamedType:
		return lambdaDeclaration + typeDTypeExpression(t.Type, context)
	case *dsl.EnumDefinition:
		base := t.BaseType
		if base == nil {
			base = dsl.Int32Type
		}

		return typeDTypeExpression(base, context)

	case *dsl.RecordDefinition:
		fields := make([]string, len(t.Fields))
		for i, f := range t.Fields {
			subarrayShape := ""
			underyingType := dsl.GetUnderlyingType(f.Type)
			if gt, ok := underyingType.(*dsl.GeneralizedType); ok {
				if vec, ok := gt.Dimensionality.(*dsl.Vector); ok && vec.Length != nil {
					subarrayShape = fmt.Sprintf("(%d,)", *vec.Length)
				} else if arr, ok := gt.Dimensionality.(*dsl.Array); ok && arr.IsFixed() {
					dims := make([]string, len(*arr.Dimensions))
					for i, d := range *arr.Dimensions {
						dims[i] = strconv.FormatUint(*d.Length, 10)
					}
					subarrayShape = fmt.Sprintf("(%s,)", strings.Join(dims, ", "))
				}
			}

			if subarrayShape != "" {
				subarrayShape = fmt.Sprintf(", %s", subarrayShape)
			}

			fields[i] = fmt.Sprintf("('%s', %s%s)", f.Name, typeDTypeExpression(f.Type, context), subarrayShape)
		}

		return fmt.Sprintf("%snp.dtype([%s], align=True)", lambdaDeclaration, strings.Join(fields, ", "))
	}

	return "np.dtype(np.object_)"
}

func typeDTypeExpression(t dsl.Type, context dTypeExpressionContext) string {
	switch t := t.(type) {
	case *dsl.SimpleType:
		context.root = false
		return typeDefinitionDTypeExpression(t.ResolvedDefinition, context)

	case *dsl.GeneralizedType:
		switch td := t.Dimensionality.(type) {
		case nil:
			if t.Cases.IsOptional() {
				return fmt.Sprintf("np.dtype([('has_value', np.dtype(np.bool_)), ('value', %s)], align=True)", typeDTypeExpression(t.Cases[1].Type, context))
			}
		case *dsl.Vector:
			if td.Length != nil {
				return typeDTypeExpression(t.ToScalar(), context)
			}

		case *dsl.Array:
			if td.IsFixed() {
				return typeDTypeExpression(t.ToScalar(), context)
			}
		}
	}

	return "np.dtype(np.object_)"
}

func getTypeSyntaxWithGenricArgsReadFromTupleArgs(t dsl.Type, context dTypeExpressionContext) string {
	var f dsl.TypeSyntaxWriter[string] = func(self dsl.TypeSyntaxWriter[string], typeOrTypeDef dsl.Node, _ string) string {
		switch t := typeOrTypeDef.(type) {
		case *dsl.GenericTypeParameter:
			return fmt.Sprintf("type_args[%d]", context.typeParameterIndexes[t])
		}

		return common.TypeSyntaxWriter(self, typeOrTypeDef, context.namespace)
	}

	return f.ToSyntax(t, context.namespace)
}
