// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/validation"
)

type ComputedFieldScope struct {
	Record          *RecordDefinition
	RewrittenFields map[*ComputedField]*ComputedField
	CurrentFields   []*ComputedField
	Variables       []*DeclarationPattern
}

func resolveComputedFields(env *Environment, errorSink *validation.ErrorSink) *Environment {
	if len(errorSink.Errors) > 0 {
		// computed fields rely on type inference. Don't attempt
		// if there already are errors in the model.
		return env
	}

	return RewriteWithContext(env, &ComputedFieldScope{}, func(node Node, context *ComputedFieldScope, self *RewriterWithContext[*ComputedFieldScope]) Node {
		switch t := node.(type) {
		case *RecordDefinition:
			if len(t.ComputedFields) == 0 {
				return t
			}
			scope := ComputedFieldScope{
				Record:          t,
				RewrittenFields: make(map[*ComputedField]*ComputedField),
			}

			return self.DefaultRewrite(node, &scope)
		case *ComputedField:

			if rewritten, ok := context.RewrittenFields[t]; ok {
				return rewritten
			}

			for _, cf := range context.CurrentFields {
				if cf == t {
					chain := make([]string, len(context.CurrentFields)+1)
					for i, cf := range context.CurrentFields {
						chain[i] = cf.Name
					}
					chain[len(chain)-1] = t.Name

					errorSink.Add(validationError(t, "cycle detected in computed fields: %s", strings.Join(chain, " -> ")))
					return t
				}
			}

			rewritten := self.DefaultRewrite(node, &ComputedFieldScope{context.Record, context.RewrittenFields, append(context.CurrentFields, t), context.Variables})
			context.RewrittenFields[t] = rewritten.(*ComputedField)
			return rewritten
		case *IntegerLiteralExpression:
			if t.Value.Sign() >= 0 {
				if t.Value.Cmp(MaxUint8) <= 0 {
					updated := *t
					updated.ResolvedType = Uint8Type
					return &updated
				}
				if t.Value.Cmp(MaxUint16) <= 0 {
					updated := *t
					updated.ResolvedType = Uint16Type
					return &updated
				}
				if t.Value.Cmp(MaxUint32) <= 0 {
					updated := *t
					updated.ResolvedType = Uint32Type
					return &updated
				}
				if t.Value.Cmp(MaxUint64) <= 0 {
					updated := *t
					updated.ResolvedType = Uint64Type
					return &updated
				}

			} else {
				if t.Value.Cmp(MinInt8) >= 0 {
					updated := *t
					updated.ResolvedType = Int8Type
					return &updated
				}
				if t.Value.Cmp(MinInt16) >= 0 {
					updated := *t
					updated.ResolvedType = Int16Type
					return &updated
				}
				if t.Value.Cmp(MinInt32) >= 0 {
					updated := *t
					updated.ResolvedType = Int32Type
					return &updated
				}
				if t.Value.Cmp(MinInt64) >= 0 {
					updated := *t
					updated.ResolvedType = Int64Type
					return &updated
				}
			}

			errorSink.Add(validationError(t, "integer literal is too large"))
			return t
		case *StringLiteralExpression:
			clone := *t
			clone.ResolvedType = StringType
			return &clone
		case *MemberAccessExpression:
			t = self.DefaultRewrite(t, context).(*MemberAccessExpression)
			target := context.Record
			if t.Target != nil {
				target = nil
				if t.Target.GetResolvedType() == nil {
					return t
				}

				switch t := GetUnderlyingType(t.Target.GetResolvedType()).(type) {
				case *SimpleType:
					if rec, ok := t.ResolvedDefinition.(*RecordDefinition); ok {
						target = rec
					}
				}

				if target == nil {
					errorSink.Add(validationError(t, "member access target must be a !record type"))
					return t
				}
			}

			if t.Target == nil {
				for _, variable := range context.Variables {
					if variable.Identifier == t.Member {
						t.ResolvedType = variable.Type
						return t
					}
				}
			}

			for _, f := range target.Fields {
				if f.Name == t.Member {
					t.ResolvedType = f.Type
					return t
				}
			}

			for _, f := range target.ComputedFields {
				if f.Name == t.Member {
					innerContext := context
					if target != context.Record {
						// we're accessing a computed field on a different record type
						updatedContext := *context
						updatedContext.Record = target
						innerContext = &updatedContext
					}
					rewrittenField := self.Rewrite(f, innerContext).(*ComputedField)
					t.ResolvedType = rewrittenField.Expression.GetResolvedType()
					t.IsComputedField = true
					return t
				}
			}

			errorSink.Add(validationError(t, "there is no variable in scope with the name '%s' nor does the record '%s' does not have a field or computed field named '%s'", t.Member, target.Name, t.Member))
			return t
		case *IndexExpression:
			t = self.DefaultRewrite(t, context).(*IndexExpression)
			if t.Target.GetResolvedType() == nil {
				return t
			}
			targetType := ToGeneralizedType(GetUnderlyingType(t.Target.GetResolvedType()))
			t.ResolvedType = targetType.ToScalar()
			argumentsValidated := false

			switch d := targetType.Dimensionality.(type) {
			case nil:
				errorSink.Add(validationError(t, "index target must be a vector, array, or map"))
				return t
			case *Vector:
				if len(t.Arguments) != 1 {
					errorSink.Add(validationError(t, "vector index must have exactly one argument"))
				}
				if d.Length != nil {
					switch arg := t.Arguments[0].Value.(type) {
					case *IntegerLiteralExpression:
						if arg.Value.Cmp(big.NewInt(int64(*d.Length))) >= 0 {
							errorSink.Add(validationError(t.Arguments[0], "index argument (%s) is too large for the vector of length %d", arg.Value.String(), *d.Length))
						}
					}
				}
			case *Map:
				if len(t.Arguments) != 1 {
					errorSink.Add(validationError(t, "map lookup must have exactly one argument"))
				}

				argType := t.Arguments[0].Value.GetResolvedType()
				if argType == nil {
					return t
				}

				if !TypesEqual(argType, d.KeyType) {
					errorSink.Add(validationError(t.Arguments[0], "incorrect map lookup argument type"))
					return t
				}
				argumentsValidated = true
			case *Array:
				labeledCount := 0
				unlabeledCount := 0
				for _, a := range t.Arguments {
					if a.Label != "" {
						labeledCount++
					} else {
						unlabeledCount++
					}
				}

				if labeledCount > 0 && unlabeledCount > 0 {
					errorSink.Add(validationError(t, "array index cannot mix labeled and unlabeled arguments"))
					return t
				}

				if d.Dimensions != nil {
					if len(t.Arguments) < len(*d.Dimensions) {
						errorSink.Add(validationError(t, "array index must provide arguments for all %d dimensions", len(*d.Dimensions)))
						return t
					}
					if len(t.Arguments) > len(*d.Dimensions) {
						errorSink.Add(validationError(t.Arguments[len(*d.Dimensions)].Value, "array index has more arguments than dimensions"))
						return t
					}

					if labeledCount > 0 {
						orderedArguments := make([]*IndexArgument, len(*d.Dimensions))
						for argIndex, arg := range t.Arguments {
							found := false
							for dimIndex, dim := range *d.Dimensions {
								if *dim.Name == arg.Label {
									found = true
									if orderedArguments[dimIndex] != nil {
										errorSink.Add(validationError(arg.Value, "array index has multiple arguments for dimension '%s'", *dim.Name))
										return t
									}

									if dimIndex != argIndex {
										expectedOrder := make([]string, len(*d.Dimensions))
										for i, dim := range *d.Dimensions {
											expectedOrder[i] = *dim.Name
										}
										errorSink.Add(validationError(arg.Value, "array index has arguments must be specified in order: %s", strings.Join(expectedOrder, ", ")))
										return t
									}

									orderedArguments[dimIndex] = arg
									break
								}
							}

							if !found {
								errorSink.Add(validationError(arg.Value, "the array has no dimension named '%s'", arg.Label))
								return t
							}
						}
						t.Arguments = orderedArguments
					}

					for i, arg := range t.Arguments {
						dimLength := (*d.Dimensions)[i].Length
						if dimLength != nil {
							switch argValue := arg.Value.(type) {
							case *IntegerLiteralExpression:
								if argValue.Value.Cmp(big.NewInt(int64(*dimLength))) >= 0 {
									label := arg.Label
									if label == "" {
										label = strconv.Itoa(i)
									}

									errorSink.Add(validationError(argValue, "index argument (%s) is too large for array dimension '%s' of length %d", argValue.Value.String(), label, *dimLength))
								}
							}
						}
					}
				}
			}

			if !argumentsValidated {
				for _, arg := range t.Arguments {
					argType := arg.Value.GetResolvedType()
					if argType == nil {
						return t
					}
					if !IsIntegralType(argType) {
						errorSink.Add(validationError(arg.Value, "index argument must be an integral type"))
						return t
					}
				}
			}

			t.ResolvedType = targetType.ToScalar()
			return t

		case *FunctionCallExpression:
			switch t.FunctionName {
			case FunctionSize:
				return resolveSizeFunctionCall(t, self, context, errorSink)
			case FunctionDimensionIndex:
				return resolveDimensionIndexFunctionCall(t, self, context, errorSink)
			case FunctionDimensionCount:
				return resolveDimensionCountFunctionCall(t, self, context, errorSink)
			default:
				errorSink.Add(validationError(t, "unknown function '%s'", t.FunctionName))
				return t
			}
		case *SwitchExpression:
			rewrittenTarget := self.Rewrite(t.Target, context).(Expression)

			resolvedTargetType := ToGeneralizedType(GetUnderlyingType(rewrittenTarget.GetResolvedType()))
			if resolvedTargetType == nil {
				return t
			}

			if resolvedTargetType.Dimensionality != nil {
				errorSink.Add(validationError(t.Target, "switch expression cannot be applied to a vector or array"))
				return t
			}

			remainingTypeIndexes := make([]int, len(resolvedTargetType.Cases))
			for i := range resolvedTargetType.Cases {
				remainingTypeIndexes[i] = i
			}

			rewrittenCases := make([]*SwitchCase, len(t.Cases))
			for i, c := range t.Cases {
				rewrittenCase := resolveSwitchCase(c, resolvedTargetType.Cases, self, context, errorSink)
				wasResolved := rewrittenCase != c
				rewrittenCases[i] = rewrittenCase

				if _, isDiscard := rewrittenCase.Pattern.(*DiscardPattern); isDiscard {
					discardedCount := 0
					for ri, rti := range remainingTypeIndexes {
						if rti != -1 {
							discardedCount++
						}
						remainingTypeIndexes[ri] = -1
					}

					if discardedCount == 0 {
						errorSink.Add(validationError(rewrittenCase.Pattern, "switch expression has no remaining cases to discard"))
					}
				} else if wasResolved {
					var patternType Type
					switch p := rewrittenCase.Pattern.(type) {
					case *TypePattern:
						patternType = p.Type
					case *DeclarationPattern:
						patternType = p.Type
					default:
						panic(fmt.Sprintf("unexpected pattern type %T", p))
					}

					notAlreadymatched := false
					for ri, rti := range remainingTypeIndexes {
						if rti == -1 {
							continue
						}

						if TypesEqual(resolvedTargetType.Cases[rti].Type, patternType) {
							remainingTypeIndexes[ri] = -1
							notAlreadymatched = true
							break
						}
					}

					if !notAlreadymatched {
						errorSink.Add(validationError(rewrittenCase.Pattern, "the switch case is not reachable"))
					}
				}
			}

			for _, rti := range remainingTypeIndexes {
				if rti != -1 {
					errorSink.Add(validationError(t, "switch expression is not exhaustive"))
				}
			}

			var commonType Type
			for _, sc := range rewrittenCases {
				if commonType == nil {
					commonType = sc.Expression.GetResolvedType()
				} else {
					resolvedType := sc.Expression.GetResolvedType()
					if resolvedType == nil {
						continue
					}
					ct, err := GetCommonType(commonType, resolvedType)
					if err != nil {
						errorSink.Add(validationError(t, "no best type was found for the switch expression"))
						return t
					}
					commonType = ct
				}
			}

			for i, sc := range rewrittenCases {
				rewrittenCases[i].Expression = insertConversion(sc.Expression, commonType)
			}

			updated := *t
			updated.Target = rewrittenTarget
			updated.Cases = rewrittenCases
			updated.ResolvedType = commonType
			return &updated
		default:
			return self.DefaultRewrite(node, context)
		}
	}).(*Environment)
}

func insertConversion(expression Expression, targetType Type) Expression {
	if TypesEqual(expression.GetResolvedType(), targetType) {
		return expression
	}

	if integerLiteral, ok := expression.(*IntegerLiteralExpression); ok && IsIntegralType(targetType) {
		updated := *integerLiteral
		updated.ResolvedType = targetType
		return &updated
	}

	return &TypeConversionExpression{
		Expression: expression,
		Type:       targetType,
	}
}

func resolveSwitchCase(switchCase *SwitchCase, typeCases TypeCases, self *RewriterWithContext[*ComputedFieldScope], context *ComputedFieldScope, errorSink *validation.ErrorSink) *SwitchCase {
	validateType := func(typePattern *TypePattern) bool {
		isValid := false
		for _, typeCase := range typeCases {
			if TypesEqual(typeCase.Type, typePattern.Type) {
				isValid = true
			}
		}

		if !isValid {
			errorSink.Add(validationError(typePattern, "the type is not a valid case for this switch expression"))
		}

		return isValid
	}

	switch t := switchCase.Pattern.(type) {
	case *DiscardPattern:
		updated := *switchCase
		updated.Expression = self.Rewrite(switchCase.Expression, context).(Expression)
		return &updated
	case *TypePattern:
		isValid := validateType(t)
		if !isValid {
			return switchCase
		}

		updated := *switchCase
		updated.Expression = self.Rewrite(switchCase.Expression, context).(Expression)
		return &updated
	case *DeclarationPattern:
		isValid := validateType(&t.TypePattern)
		if !isValid {
			return switchCase
		}

		if t.Type == nil {
			errorSink.Add(validationError(t, "a declaration pattern cannot be used with the null type"))
		}

		updated := *switchCase
		updated.Expression = self.Rewrite(switchCase.Expression, &ComputedFieldScope{context.Record, context.RewrittenFields, context.CurrentFields, append(context.Variables, t)}).(Expression)
		return &updated
	default:
		panic(fmt.Sprintf("unexpected switch case pattern type %T", switchCase.Pattern))
	}
}

func resolveDimensionCountFunctionCall(functionCall *FunctionCallExpression, visitor *RewriterWithContext[*ComputedFieldScope], context *ComputedFieldScope, errorSink *validation.ErrorSink) Expression {
	functionCall = visitor.DefaultRewrite(functionCall, context).(*FunctionCallExpression)
	functionCall.ResolvedType = SizeType

	if len(functionCall.Arguments) != 1 {
		errorSink.Add(validationError(functionCall, "%s() expects 1 argument, but called with %d", FunctionDimensionCount, len(functionCall.Arguments)))
		return functionCall
	}

	if functionCall.Arguments[0].GetResolvedType() == nil {
		return functionCall
	}

	target := ToGeneralizedType(GetUnderlyingType(functionCall.Arguments[0].GetResolvedType()))
	switch dim := target.Dimensionality.(type) {
	case *Array:
		if dim.Dimensions != nil {
			// count of dimensions as bigint
			dimCount := big.NewInt(int64(len(*dim.Dimensions)))
			// simplification
			return &IntegerLiteralExpression{
				NodeMeta:     functionCall.NodeMeta,
				Value:        *dimCount,
				ResolvedType: SizeType,
			}
		}
	default:
		errorSink.Add(validationError(functionCall, "%s() must be called with an !array argument", FunctionDimensionCount))
	}

	return functionCall
}

func resolveDimensionIndexFunctionCall(functionCall *FunctionCallExpression, visitor *RewriterWithContext[*ComputedFieldScope], context *ComputedFieldScope, errorSink *validation.ErrorSink) Expression {
	functionCall = visitor.DefaultRewrite(functionCall, context).(*FunctionCallExpression)
	functionCall.ResolvedType = SizeType

	if len(functionCall.Arguments) != 2 {
		errorSink.Add(validationError(functionCall, "%s() expects 2 arguments, but called with %d", FunctionDimensionIndex, len(functionCall.Arguments)))
		return functionCall
	}

	if functionCall.Arguments[0].GetResolvedType() == nil {
		return functionCall
	}

	target := ToGeneralizedType(GetUnderlyingType(functionCall.Arguments[0].GetResolvedType()))
	switch dim := target.Dimensionality.(type) {
	case *Array:
		if functionCall.Arguments[1].GetResolvedType() == nil {
			return functionCall
		}

		hasNamedDimension := false
		if dim.Dimensions != nil {
			for _, v := range *dim.Dimensions {
				if v.Name != nil {
					hasNamedDimension = true
					break
				}
			}
		}

		if !hasNamedDimension {
			errorSink.Add(validationError(functionCall, "%s() is only valid for arrays with named dimensions", FunctionDimensionIndex))
			return functionCall
		}

		if primitive, ok := GetPrimitiveType(functionCall.Arguments[1].GetResolvedType()); ok && primitive == String {
			if stringLiteral, ok := functionCall.Arguments[1].(*StringLiteralExpression); ok {
				for i, dim := range *dim.Dimensions {
					if dim.Name != nil && *dim.Name == stringLiteral.Value {
						// simplification
						return &IntegerLiteralExpression{
							NodeMeta:     functionCall.NodeMeta,
							Value:        *big.NewInt(int64(i)),
							ResolvedType: SizeType,
						}
					}
				}

				errorSink.Add(validationError(functionCall, "the array does not have a dimension named '%s'", stringLiteral.Value))
				return functionCall
			}

		} else {
			errorSink.Add(validationError(functionCall.Arguments[1], "the second argument to %s() must be a dimension name string", FunctionDimensionIndex))
			return functionCall
		}

	default:
		errorSink.Add(validationError(functionCall, "%s() must be called with an !array as the first argument", FunctionDimensionIndex))
	}

	return functionCall
}

func resolveSizeFunctionCall(functionCall *FunctionCallExpression, visitor *RewriterWithContext[*ComputedFieldScope], context *ComputedFieldScope, errorSink *validation.ErrorSink) Expression {
	functionCall = visitor.DefaultRewrite(functionCall, context).(*FunctionCallExpression)
	functionCall.ResolvedType = SizeType

	if len(functionCall.Arguments) == 0 || len(functionCall.Arguments) > 2 {
		errorSink.Add(validationError(functionCall, "%s() expects 1 or 2 arguments, but called with %d", FunctionSize, len(functionCall.Arguments)))
		return functionCall
	}

	if functionCall.Arguments[0].GetResolvedType() == nil {
		return functionCall
	}

	target := ToGeneralizedType(GetUnderlyingType(functionCall.Arguments[0].GetResolvedType()))
	switch dim := target.Dimensionality.(type) {
	case *Vector:
		if len(functionCall.Arguments) == 2 {
			errorSink.Add(validationError(functionCall, "%s() does not accept a second argument when called with a !vector", FunctionSize))
			return functionCall
		}

		if dim.Length != nil {
			// simplification
			return &IntegerLiteralExpression{
				NodeMeta:     functionCall.NodeMeta,
				Value:        *big.NewInt(int64(*dim.Length)),
				ResolvedType: SizeType,
			}
		}
	case *Map:
		if len(functionCall.Arguments) == 2 {
			errorSink.Add(validationError(functionCall, "%s() does not accept a second argument when called with a !map", FunctionSize))
			return functionCall
		}

	case *Array:
		if len(functionCall.Arguments) == 1 {
			if dim.IsFixed() {
				// simplification
				size := uint64(0)
				if len(*dim.Dimensions) > 0 {
					size = 1
					for _, v := range *dim.Dimensions {
						size *= *v.Length
					}
				}

				return &IntegerLiteralExpression{
					NodeMeta:     functionCall.NodeMeta,
					Value:        *big.NewInt(int64(size)),
					ResolvedType: SizeType,
				}
			}
			return functionCall
		}
		if functionCall.Arguments[1].GetResolvedType() == nil {
			return functionCall
		}

		arg2Type := functionCall.Arguments[1].GetResolvedType()
		if primitive, ok := GetPrimitiveType(arg2Type); ok {
			if primitive == String {
				if stringLit, ok := functionCall.Arguments[1].(*StringLiteralExpression); ok {
					if dim.Dimensions != nil {
						for i, d := range *dim.Dimensions {
							if d.Name != nil && *d.Name == stringLit.Value {
								if d.Length != nil {
									// simplification
									return &IntegerLiteralExpression{
										NodeMeta:     functionCall.NodeMeta,
										Value:        *big.NewInt(int64(*d.Length)),
										ResolvedType: SizeType,
									}
								}

								intLiteral := &IntegerLiteralExpression{
									Value: *big.NewInt(int64(i)),
								}
								functionCall.Arguments[1] = visitor.Rewrite(intLiteral, context).(Expression)
								return functionCall
							}
						}
					}

					errorSink.Add(validationError(functionCall.Arguments[1], "this array does not have a dimension named '%s'", stringLit.Value))
					return functionCall
				}
				dimensionIndexCall := &FunctionCallExpression{
					NodeMeta:     functionCall.NodeMeta,
					FunctionName: "dimensionIndex",
					Arguments:    []Expression{functionCall.Arguments[0], functionCall.Arguments[1]},
				}
				functionCall.Arguments[1] = visitor.Rewrite(dimensionIndexCall, context).(Expression)
				return functionCall
			}
			if IsIntegralPrimitive(primitive) {
				if intLit, ok := functionCall.Arguments[1].(*IntegerLiteralExpression); ok {
					if intLit.Value.Sign() < 0 {
						errorSink.Add(validationError(functionCall.Arguments[1], "array dimension cannot be negative"))
						return functionCall
					}

					if dim.Dimensions != nil {
						if intLit.Value.Cmp(big.NewInt(int64(len(*dim.Dimensions)))) >= 0 {
							errorSink.Add(validationError(functionCall.Arguments[1], "array dimension index is out of bounds"))
							return functionCall
						}

						if dimension := (*dim.Dimensions)[intLit.Value.Int64()]; dimension.Length != nil {
							// simplification
							return &IntegerLiteralExpression{
								NodeMeta:     functionCall.NodeMeta,
								Value:        *big.NewInt(int64(*(*dim.Dimensions)[intLit.Value.Int64()].Length)),
								ResolvedType: SizeType,
							}
						}
					}
				}
				return functionCall
			}
		}

		errorSink.Add(validationError(functionCall.Arguments[1], "%s() expects a string or integer as its second argument", FunctionSize))
	default:
		errorSink.Add(validationError(functionCall, "%s() must be called with a !vector, !array, or !map as the first argument", FunctionSize))
	}

	return functionCall
}
