// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/microsoft/yardl/tooling/internal/validation"
)

func IsGeneric(t TypeDefinition) bool {
	return len(t.GetDefinitionMeta().TypeParameters) > 0
}

func MakeGenericType(genericTypeDefinition TypeDefinition, typeArguments []Type, shallow bool) (TypeDefinition, error) {
	meta := genericTypeDefinition.GetDefinitionMeta()
	if len(meta.TypeParameters) != len(typeArguments) {
		return nil, errors.New("incorrect number of type arguments given")
	}

	errorSink := validation.ErrorSink{}
	rewritten := Rewrite(genericTypeDefinition, func(self Rewriter, node Node) Node {
		switch t := node.(type) {
		case *DefinitionMeta:
			newMeta := *t
			newMeta.TypeArguments = typeArguments
			return &newMeta
		case *SimpleType:
			if shallow {
				return t
			}
			if targetParam, ok := t.ResolvedDefinition.(*GenericTypeParameter); ok {
				for i, genericTypeParam := range meta.TypeParameters {
					if genericTypeParam == targetParam {
						return typeArguments[i]
					}
				}
				errorSink.Add(validationError(t, "internal error: unable to substitute generic type parameter"))
			}

			rewrittenResolvedType := self.Rewrite(t.ResolvedDefinition)
			if rewrittenResolvedType == t.ResolvedDefinition {
				return t
			}
			newType := *t
			newType.ResolvedDefinition = rewrittenResolvedType.(TypeDefinition)
			return &newType
		case *GenericTypeParameter:
			return t
		default:
			return self.DefaultRewrite(node)
		}
	}).(TypeDefinition)

	return rewritten, errorSink.AsError()
}

func TypeDefinitionsEqual(a, b TypeDefinition) bool {
	if a == b {
		return true
	}
	if a == nil {
		return b == nil
	}
	if b == nil {
		return a == nil
	}

	aMeta := a.GetDefinitionMeta()
	bMeta := b.GetDefinitionMeta()
	if aMeta.Namespace != bMeta.Namespace || aMeta.Name != bMeta.Name {
		return false
	}

	if len(aMeta.TypeParameters) != len(bMeta.TypeParameters) {
		return false
	}

	for i, aTypeParam := range aMeta.TypeParameters {
		if !TypeDefinitionsEqual(aTypeParam, bMeta.TypeParameters[i]) {
			return false
		}
	}

	if len(aMeta.TypeArguments) != len(bMeta.TypeArguments) {
		return false
	}

	for i, aTypeArg := range aMeta.TypeArguments {
		if !TypesEqual(aTypeArg, bMeta.TypeArguments[i]) {
			return false
		}
	}

	switch ta := a.(type) {
	case *PrimitiveDefinition, *GenericTypeParameter:
		return true
	case *RecordDefinition:
		tb, ok := b.(*RecordDefinition)
		if !ok {
			return false
		}

		if len(ta.Fields) != len(tb.Fields) {
			return false
		}

		if len(ta.ComputedFields) != len(tb.ComputedFields) {
			return false
		}

		for i, fa := range ta.Fields {
			fb := tb.Fields[i]
			if fa.Name != fb.Name || !TypesEqual(fa.Type, fb.Type) {
				return false
			}
		}

		for i, fa := range ta.ComputedFields {
			fb := tb.ComputedFields[i]
			if fa.Name != fb.Name || !ExpressionsEqual(fa.Expression, fb.Expression) {
				return false
			}
		}

		return true
	case *ProtocolDefinition:
		tb, ok := b.(*ProtocolDefinition)
		if !ok {
			return false
		}

		if len(ta.Sequence) != len(tb.Sequence) {
			return false
		}

		for i, sa := range ta.Sequence {
			sb := tb.Sequence[i]
			if sa.Name != sb.Name || !TypesEqual(sa.Type, sb.Type) {
				return false
			}
		}

		return true
	case *EnumDefinition:
		tb, ok := b.(*EnumDefinition)
		if !ok {
			return false
		}

		if len(ta.Values) != len(tb.Values) {
			return false
		}

		for i, va := range ta.Values {
			vb := tb.Values[i]
			if va.Symbol != vb.Symbol || va.IntegerValue.Cmp(&vb.IntegerValue) != 0 {
				return false
			}
		}

		return true
	case *NamedType:
		tb, ok := b.(*NamedType)
		if !ok {
			return false
		}
		return TypesEqual(ta.Type, tb.Type)
	default:
		panic(fmt.Sprintf("unexpected type %T", ta))
	}
}

func TypesEqual(a, b Type) bool {
	if a == b {
		return true
	}
	if a == nil {
		return b == nil
	}
	if b == nil {
		return a == nil
	}

	a = GetUnderlyingType(a)
	b = GetUnderlyingType(b)

	switch ta := a.(type) {
	case *SimpleType:
		tb, ok := b.(*SimpleType)
		if !ok {
			return false
		}

		if (ta.ResolvedDefinition == nil) != (tb.ResolvedDefinition == nil) {
			return false
		}

		if ta.ResolvedDefinition != nil {
			return TypeDefinitionsEqual(ta.ResolvedDefinition, tb.ResolvedDefinition)
		}

		return ta.Name == tb.Name
	case *GeneralizedType:
		tb, ok := b.(*GeneralizedType)
		if !ok {
			return false
		}
		if len(ta.Cases) != len(tb.Cases) {
			return false
		}

		for i := 0; i < len(ta.Cases); i++ {
			if !TypesEqual(ta.Cases[i].Type, tb.Cases[i].Type) {
				return false
			}
		}

		if ta.Dimensionality == nil {
			return tb.Dimensionality == nil
		}
		if tb.Dimensionality == nil {
			return ta.Dimensionality == nil
		}

		switch da := ta.Dimensionality.(type) {
		case *Vector:
			db, ok := tb.Dimensionality.(*Vector)
			if !ok {
				return false
			}

			if da.Length == nil {
				return db.Length == nil
			}

			if db.Length == nil {
				return da.Length == nil
			}

			return *da.Length == *db.Length
		case *Array:
			db, ok := tb.Dimensionality.(*Array)
			if !ok {
				return false
			}

			if da.Dimensions == nil {
				return db.Dimensions == nil
			}
			if db.Dimensions == nil {
				return da.Dimensions == nil
			}

			if len(*da.Dimensions) != len(*db.Dimensions) {
				return false
			}

			for i := 0; i < len(*da.Dimensions); i++ {
				da := (*da.Dimensions)[i]
				db := (*db.Dimensions)[i]

				// not taking labels into account

				if da.Length == nil {
					return db.Length == nil
				}
				if db.Length == nil {
					return da.Length == nil
				}

				if *da.Length != *db.Length {
					return false
				}
			}

			return true
		case *Stream:
			_, ok := tb.Dimensionality.(*Stream)
			return ok

		default:
			panic(fmt.Sprintf("unexpected type %T", da))
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", ta))
	}
}

func ExpressionsEqual(a, b Expression) bool {
	if a == b {
		return true
	}
	if a == nil {
		return b == nil
	}
	if b == nil {
		return a == nil
	}

	switch ta := a.(type) {
	case *IntegerLiteralExpression:
		tb, ok := b.(*IntegerLiteralExpression)
		if !ok {
			return false
		}
		return TypesEqual(ta.ResolvedType, tb.ResolvedType) && ta.Value.Cmp(&tb.Value) == 0
	case *StringLiteralExpression:
		tb, ok := b.(*StringLiteralExpression)
		if !ok {
			return false
		}
		return TypesEqual(ta.ResolvedType, tb.ResolvedType) && ta.Value == tb.Value
	case *MemberAccessExpression:
		tb, ok := b.(*MemberAccessExpression)
		if !ok {
			return false
		}
		return ExpressionsEqual(ta.Target, tb.Target) && ta.Member == tb.Member && TypesEqual(ta.ResolvedType, tb.ResolvedType)
	case *FunctionCallExpression:
		tb, ok := b.(*FunctionCallExpression)
		if !ok {
			return false
		}
		if ta.FunctionName != tb.FunctionName || !TypesEqual(ta.ResolvedType, tb.ResolvedType) || len(ta.Arguments) != len(tb.Arguments) {
			return false
		}

		for i := 0; i < len(ta.Arguments); i++ {
			if !ExpressionsEqual(ta.Arguments[i], tb.Arguments[i]) {
				return false
			}
		}

		return true

	case *SwitchExpression:
		tb, ok := b.(*SwitchExpression)
		if !ok {
			return false
		}

		if !ExpressionsEqual(ta.Target, tb.Target) || len(ta.Cases) != len(tb.Cases) {
			return false
		}
		for i, taCase := range ta.Cases {
			tbCase := tb.Cases[i]
			if !PatternsEqual(taCase.Pattern, tbCase.Pattern) || !ExpressionsEqual(taCase.Expression, tbCase.Expression) {
				return false
			}
		}

		return true

	case *IndexExpression:
		tb, ok := b.(*IndexExpression)
		if !ok {
			return false
		}

		if !TypesEqual(ta.ResolvedType, tb.ResolvedType) || !ExpressionsEqual(ta.Target, tb.Target) || len(ta.Arguments) != len(tb.Arguments) {
			return false
		}

		for i := 0; i < len(ta.Arguments); i++ {
			iaa := ta.Arguments[i]
			iab := tb.Arguments[i]
			if iaa.Label != iab.Label || !ExpressionsEqual(iaa.Value, iab.Value) {
				return false
			}
		}

		return true
	default:
		panic(fmt.Sprintf("unexpected type %T", ta))
	}
}

func PatternsEqual(a, b Pattern) bool {
	if a == b {
		return true
	}
	if a == nil {
		return b == nil
	}
	if b == nil {
		return a == nil
	}

	switch ta := a.(type) {
	case *DiscardPattern:
		_, ok := b.(*DiscardPattern)
		return ok
	case *TypePattern:
		tb, ok := b.(*TypePattern)
		if !ok {
			return false
		}
		return TypesEqual(ta.Type, tb.Type)
	case *DeclarationPattern:
		tb, ok := b.(*DeclarationPattern)
		if !ok {
			return false
		}

		return PatternsEqual(&ta.TypePattern, &tb.TypePattern) && ta.Identifier == tb.Identifier

	default:
		panic(fmt.Sprintf("unexpected type %T", ta))
	}
}

func ToGeneralizedType(t Type) *GeneralizedType {
	switch t := t.(type) {
	case nil:
		return nil
	case *GeneralizedType:
		return t
	case *SimpleType:
		return &GeneralizedType{NodeMeta: t.NodeMeta, Cases: TypeCases{&TypeCase{Type: t}}}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func GetUnderlyingType(t Type) Type {
	underlyingTypeFromTypeDefinition := func(t TypeDefinition) Type {
		switch t := t.(type) {
		case *NamedType:
			return GetUnderlyingType(t.Type)
		default:
			return nil
		}
	}

	switch t := t.(type) {
	case nil:
		return nil
	case *SimpleType:
		if underlyingType := underlyingTypeFromTypeDefinition(t.ResolvedDefinition); underlyingType != nil {
			return underlyingType
		}
	case *GeneralizedType:
		switch t.Dimensionality.(type) {
		case nil:
			if t.Cases.IsSingle() {
				return GetUnderlyingType(t.Cases[0].Type)
			}
		}
	}

	return t
}

func GetPrimitiveType(t Type) (primitive PrimitiveDefinition, ok bool) {
	switch t := GetUnderlyingType(t).(type) {
	case *SimpleType:
		primitive, ok := t.ResolvedDefinition.(PrimitiveDefinition)
		return primitive, ok
	case *GeneralizedType:
		switch t.Dimensionality.(type) {
		case nil:
			if t.Cases.IsSingle() {
				return GetPrimitiveType(t.Cases[0].Type)
			}
		}
	}

	return "", false
}

func IsIntegralPrimitive(prim PrimitiveDefinition) bool {
	switch prim {
	case Int8, Int16, Int32, Int64, Uint8, Uint16, Uint32, Uint64, Size:
		return true
	default:
		return false
	}
}

func IsIntegralType(t Type) bool {
	primitive, ok := GetPrimitiveType(t)
	return ok && IsIntegralPrimitive(primitive)
}

var ErrNoCommonType = errors.New("no common type")

type primitivePair struct {
	A, B TypeDefinition
}

var (
	Zero = big.NewInt(0)

	MaxInt8 = big.NewInt(math.MaxInt8)
	MinInt8 = big.NewInt(math.MinInt8)

	MaxInt16 = big.NewInt(math.MaxInt16)
	MinInt16 = big.NewInt(math.MinInt16)

	MaxInt32 = big.NewInt(math.MaxInt32)
	MinInt32 = big.NewInt(math.MinInt32)

	MaxInt64 = big.NewInt(math.MaxInt64)
	MinInt64 = big.NewInt(math.MinInt64)

	MaxUint8 = big.NewInt(math.MaxUint8)

	MaxUint16 = big.NewInt(math.MaxUint16)

	MaxUint32 = big.NewInt(math.MaxUint32)

	MaxUint64 = func() *big.Int {
		var i big.Int
		i.SetUint64(math.MaxUint64)
		return &i
	}()

	MaxSize = MaxUint64
)

var (
	Int8Type           = &SimpleType{Name: Int8, ResolvedDefinition: PrimitiveInt8}
	Int16Type          = &SimpleType{Name: Int16, ResolvedDefinition: PrimitiveInt16}
	Int32Type          = &SimpleType{Name: Int32, ResolvedDefinition: PrimitiveInt32}
	Int64Type          = &SimpleType{Name: Int64, ResolvedDefinition: PrimitiveInt64}
	Uint8Type          = &SimpleType{Name: Uint8, ResolvedDefinition: PrimitiveUint8}
	Uint16Type         = &SimpleType{Name: Uint16, ResolvedDefinition: PrimitiveUint16}
	Uint32Type         = &SimpleType{Name: Uint32, ResolvedDefinition: PrimitiveUint32}
	Uint64Type         = &SimpleType{Name: Uint64, ResolvedDefinition: PrimitiveUint64}
	SizeType           = &SimpleType{Name: Size, ResolvedDefinition: PrimitiveSize}
	BoolType           = &SimpleType{Name: Bool, ResolvedDefinition: PrimitiveBool}
	Float32Type        = &SimpleType{Name: Float32, ResolvedDefinition: PrimitiveFloat32}
	Float64Type        = &SimpleType{Name: Float64, ResolvedDefinition: PrimitiveFloat64}
	ComplexFloat32Type = &SimpleType{Name: ComplexFloat32, ResolvedDefinition: PrimitiveComplexFloat32}
	ComplexFloat64Type = &SimpleType{Name: ComplexFloat64, ResolvedDefinition: PrimitiveComplexFloat64}
	StringType         = &SimpleType{Name: String, ResolvedDefinition: PrimitiveString}
)

var commonTypeMap = func() map[primitivePair]PrimitiveDefinition {
	var commonTypes = []struct {
		a      PrimitiveDefinition
		b      PrimitiveDefinition
		common PrimitiveDefinition
	}{
		{Int8, Int16, Int16},
		{Int8, Int32, Int32},
		{Int8, Int64, Int64},
		{Int8, Uint8, Int16},
		{Int8, Uint16, Int32},
		{Int8, Uint32, Int64},
		{Int8, Uint64, Int64},
		{Int8, Float32, Float32},
		{Int8, Float64, Float64},

		{Int16, Int32, Int32},
		{Int16, Int64, Int64},
		{Int16, Uint8, Int16},
		{Int16, Uint16, Int32},
		{Int16, Uint32, Int64},
		{Int16, Float32, Float32},
		{Int16, Float64, Float64},

		{Int32, Int64, Int64},
		{Int32, Uint8, Int32},
		{Int32, Uint16, Int32},
		{Int32, Uint32, Int64},
		{Int32, Float32, Float32},
		{Int32, Float64, Float64},

		{Uint8, Uint16, Uint16},
		{Uint8, Uint32, Uint32},
		{Uint8, Uint64, Uint64},
		{Uint8, Size, Size},
		{Uint8, Float32, Float32},
		{Uint8, Float64, Float64},

		{Uint16, Uint32, Uint32},
		{Uint16, Uint64, Uint64},
		{Uint16, Size, Size},
		{Uint16, Float32, Float32},
		{Uint16, Float64, Float64},

		{Uint32, Uint64, Uint64},
		{Uint32, Size, Size},
		{Uint32, Float32, Float32},
		{Uint32, Float64, Float64},

		{Uint64, Size, Size},
		{Uint64, Float32, Float32},
		{Uint64, Float64, Float64},

		{Float32, Float64, Float64},
		{ComplexFloat32, ComplexFloat64, ComplexFloat64},
	}

	m := make(map[primitivePair]PrimitiveDefinition)
	for _, ct := range commonTypes {
		m[primitivePair{ct.a, ct.b}] = ct.common
		m[primitivePair{ct.b, ct.a}] = ct.common
	}

	return m
}()

func GetCommonType(a, b Type) (Type, error) {
	if a == b {
		return a, nil
	}

	underlyingA := GetUnderlyingType(a)
	underlyingB := GetUnderlyingType(b)
	if underlyingA == underlyingB {
		return underlyingA, nil
	}

	if primA, ok := GetPrimitiveType(underlyingA); ok {
		if primB, ok := GetPrimitiveType(underlyingB); ok {
			if primA != primB {
				if common, ok := commonTypeMap[primitivePair{primA, primB}]; ok {
					return &SimpleType{Name: string(common), ResolvedDefinition: common}, nil
				}

				return nil, ErrNoCommonType
			}

			return underlyingA, nil
		}
	}

	return nil, ErrNoCommonType
}
