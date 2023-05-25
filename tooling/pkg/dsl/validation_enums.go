// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"math/big"

	"github.com/microsoft/yardl/tooling/internal/validation"
)

func validateEnums(env *Environment, errorSink *validation.ErrorSink) *Environment {
	Visit(env, func(self Visitor, node Node) {
		enum, ok := node.(*EnumDefinition)
		if !ok {
			self.VisitChildren(node)
			return
		}

		var enumKind string
		if enum.IsFlags {
			enumKind = "flags"
		} else {
			enumKind = "enum"
		}

		// verify that the enum symbols and integer values are unique
		symbols := make(map[string]any)
		symbolsByVal := make(map[string][]string)
		for _, enumValue := range enum.Values {
			if !memberNameRegex.MatchString(enumValue.Symbol) {
				errorSink.Add(validationError(enumValue, "in %s '%s', the symbol name '%s' must be camelCased matching the format %s", enumKind, enum.Name, enumValue.Symbol, memberNameRegex.String()))
			}

			symbolsByVal[enumValue.IntegerValue.String()] = append(symbolsByVal[enumValue.IntegerValue.String()], enumValue.Symbol)
			if _, found := symbols[enumValue.Symbol]; found {
				errorSink.Add(validationError(enum, "in %s '%s', the symbol '%s' is defined more than once", enumKind, enum.Name, enumValue.Symbol))
			} else {
				symbols[enumValue.Symbol] = nil
			}
		}

		for v, syms := range symbolsByVal {
			if len(syms) > 1 {
				errorSink.Add(validationError(enum, "in %s '%s', the symbols %v have the same value of %s", enumKind, enum.Name, syms, v))
			}
		}

		var baseType PrimitiveDefinition
		if enum.BaseType == nil {
			baseType = PrimitiveInt32
		} else {
			underlyingType := GetUnderlyingType(enum.BaseType)
			switch bt := underlyingType.(type) {
			case *SimpleType:
				switch bt := bt.ResolvedDefinition.(type) {
				case nil:
					// already a type resolution error for this
					return
				case PrimitiveDefinition:
					baseType = bt
				}
			}
		}
		var minValue *big.Int
		var maxValue *big.Int
		switch baseType {
		case Int8:
			minValue = MinInt8
			maxValue = MaxInt8
		case Uint8:
			minValue = Zero
			maxValue = MaxUint8
		case Int16:
			minValue = MinInt16
			maxValue = MaxInt16
		case Uint16:
			minValue = Zero
			maxValue = MaxUint16
		case Int32:
			minValue = MinInt32
			maxValue = MaxInt32
		case Uint32:
			minValue = Zero
			maxValue = MaxUint32
		case Int64:
			minValue = MinInt64
			maxValue = MaxInt64
		case Uint64, Size:
			minValue = Zero
			maxValue = MaxUint64
		default:
			errorSink.Add(validationError(enum, "in %s '%s', the base type must be an integer type", enumKind, enum.Name))
			return
		}

		for _, enumValue := range enum.Values {
			if enumValue.IntegerValue.Cmp(minValue) < 0 || enumValue.IntegerValue.Cmp(maxValue) > 0 {
				errorSink.Add(validationError(enumValue, "in %s '%s', the value '%s' for symbol '%s' is out of range for the base type '%s'", enumKind, enum.Name, enumValue.IntegerValue.String(), enumValue.Symbol, baseType))
			}
		}
	})

	return env
}
