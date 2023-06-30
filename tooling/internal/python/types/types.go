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
		}

		if len(rec.Fields) == 0 {
			w.WriteStringln("pass")
		}
	})
	w.WriteStringln("")
}

func GetGenericBase(t dsl.TypeDefinition) string {
	meta := t.GetDefinitionMeta()
	if len(meta.TypeParameters) == 0 {
		return ""
	}

	var typeParams []string
	for _, tp := range meta.TypeParameters {
		use := tp.Tags[common.TypeParameterUseTagKey].(common.TypeParameterUse)
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
	var baseType string
	if enum.IsFlags {
		baseType = "enum.Flag"
	} else {
		baseType = "enum.Enum"
	}
	fmt.Fprintf(w, "class %s(%s):\n", common.TypeSyntaxWithoutTypeParameters(enum, enum.Namespace), baseType)

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
			return typeDefault(t.Cases[0].Type, contextNamespace, st)
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

			return fmt.Sprintf("np.zeros(0, dtype=%s)", dtype), defaultValueKindLambda

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

		for _, t2 := range t.TypeArguments {
			if containsGenericTypeParameter(t2) {
				return "", defaultValueKindNone
			}
		}

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

func containsGenericTypeParameter(node dsl.Node) bool {
	contains := false
	dsl.Visit(node, func(self dsl.Visitor, node dsl.Node) {
		switch node := node.(type) {
		case *dsl.GenericTypeParameter:
			contains = true
			return
		case *dsl.SimpleType:
			self.Visit(node.ResolvedDefinition)
		}
		self.VisitChildren(node)
	})
	return contains
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
