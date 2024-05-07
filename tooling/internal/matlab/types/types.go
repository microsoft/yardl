// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package types

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/matlab/common"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func WriteTypes(fw *common.MatlabFileWriter, ns *dsl.Namespace, st dsl.SymbolTable) error {
	unionGenerated := make(map[string]bool)

	for _, td := range ns.TypeDefinitions {
		if err := writeUnionClasses(fw, td, unionGenerated); err != nil {
			return err
		}
	}

	for _, td := range ns.TypeDefinitions {
		var err error
		switch td := td.(type) {
		case *dsl.NamedType:
			if !unionGenerated[td.Name] {
				err = writeNamedType(fw, td)
			}
		case *dsl.EnumDefinition:
			err = writeEnum(fw, td, nil)
		case *dsl.RecordDefinition:
			err = writeRecord(fw, td, st)
		default:
			panic(fmt.Sprintf("unsupported type definition: %T", td))
		}
		if err != nil {
			return err
		}
	}

	for _, p := range ns.Protocols {
		if err := writeUnionClasses(fw, p, unionGenerated); err != nil {
			return err
		}
	}

	return nil
}

func writeUnionClasses(fw *common.MatlabFileWriter, td dsl.TypeDefinition, unionGenerated map[string]bool) error {
	var writeError error
	dsl.Visit(td, func(self dsl.Visitor, node dsl.Node) {
		switch node := node.(type) {
		case *dsl.GeneralizedType:
			if node.Cases.IsUnion() {
				unionClassName := common.UnionClassName(node)
				if !unionGenerated[unionClassName] {
					if _, isNamedType := td.(*dsl.NamedType); isNamedType {
						// This is a named type defining a union, so we will use the named type's name instead
						unionClassName = td.GetDefinitionMeta().Name
					}
					writeError = fw.WriteFile(unionClassName, func(w *formatting.IndentedWriter) {
						writeUnionClass(w, unionClassName, node, td.GetDefinitionMeta().Namespace)
					})
					if writeError != nil {
						return
					}
					unionGenerated[unionClassName] = true
				}
			}
		}
		self.VisitChildren(node)
	})

	return writeError
}

func writeUnionClass(w *formatting.IndentedWriter, className string, generalizedType *dsl.GeneralizedType, contextNamespace string) error {
	qualifiedClassName := fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(contextNamespace), className)
	fmt.Fprintf(w, "classdef %s < yardl.Union\n", className)
	common.WriteBlockBody(w, func() {
		w.WriteStringln("methods (Static)")
		common.WriteBlockBody(w, func() {
			index := 1
			for _, tc := range generalizedType.Cases {
				if tc.Type == nil {
					continue
				}

				fmt.Fprintf(w, "function res = %s(value)\n", formatting.ToPascalCase(tc.Tag))
				common.WriteBlockBody(w, func() {
					fmt.Fprintf(w, "res = %s(%d, value);\n", qualifiedClassName, index)
				})
				index += 1
				w.WriteStringln("")
			}

			writeZerosStaticMethod(w, qualifiedClassName, []string{"0", "yardl.None"})
		})
		w.WriteStringln("")

		w.WriteStringln("methods")
		common.WriteBlockBody(w, func() {
			index := 1
			for _, tc := range generalizedType.Cases {
				if tc.Type == nil {
					continue
				}
				fmt.Fprintf(w, "function res = is%s(self)\n", formatting.ToPascalCase(tc.Tag))
				common.WriteBlockBody(w, func() {
					fmt.Fprintf(w, "res = self.index == %d;\n", index)
				})
				index += 1
				w.WriteStringln("")
			}

			w.WriteStringln("function eq = eq(self, other)")
			common.WriteBlockBody(w, func() {
				fmt.Fprintf(w, "eq = isa(other, \"%s\") && other.index == self.index && other.value == self.value;\n", qualifiedClassName)
			})
			w.WriteStringln("")
			w.WriteStringln("function ne = ne(self, other)")
			common.WriteBlockBody(w, func() {
				w.WriteStringln("ne = ~self.eq(other);")
			})
		})

	})

	return nil
}

func writeNamedType(fw *common.MatlabFileWriter, td *dsl.NamedType) error {
	return fw.WriteFile(common.TypeIdentifierName(td.Name), func(w *formatting.IndentedWriter) {
		ut := dsl.GetUnderlyingType(td.Type)
		// If the underlying type is a RecordDefinition or Optional, we will generate a "function" alias
		if st, ok := ut.(*dsl.SimpleType); ok {
			if _, ok := st.ResolvedDefinition.(*dsl.RecordDefinition); ok {
				fmt.Fprintf(w, "function c = %s(varargin) \n", common.TypeIdentifierName(td.Name))
				common.WriteBlockBody(w, func() {
					common.WriteComment(w, td.Comment)
					fmt.Fprintf(w, "c = %s(varargin{:});\n", common.TypeSyntax(td.Type, td.Namespace))
				})
				return
			}
		} else if gt, ok := ut.(*dsl.GeneralizedType); ok {
			if gt.Cases.IsOptional() {
				innerType := gt.Cases[1].Type
				fmt.Fprintf(w, "function o = %s(value) \n", common.TypeIdentifierName(td.Name))
				common.WriteBlockBody(w, func() {
					common.WriteComment(w, td.Comment)
					if !dsl.TypeContainsGenericTypeParameter(innerType) {
						w.WriteStringln("arguments")
						common.WriteBlockBody(w, func() {
							fmt.Fprintf(w, "value %s\n", common.TypeSyntax(innerType, td.Namespace))
						})
					}
					fmt.Fprintf(w, "o = %s(value);\n", common.TypeSyntax(td.Type, td.Namespace))
				})
				return
			}

			switch gt.Dimensionality.(type) {
			case *dsl.Vector, *dsl.Array:
				scalar := gt.ToScalar()
				fmt.Fprintf(w, "function a = %s(array) \n", common.TypeIdentifierName(td.Name))
				common.WriteBlockBody(w, func() {
					common.WriteComment(w, td.Comment)
					if !dsl.TypeContainsGenericTypeParameter(scalar) {
						w.WriteStringln("arguments")
						common.WriteBlockBody(w, func() {
							fmt.Fprintf(w, "array %s\n", common.TypeSyntax(scalar, td.Namespace))
						})
					}
					w.WriteStringln("a = array;")
				})
				return
			}
		}

		// Otherwise, it's a subclass of the underlying type
		fmt.Fprintf(w, "classdef %s < %s\n", common.TypeIdentifierName(td.Name), common.TypeSyntax(td.Type, td.Namespace))
		common.WriteBlockBody(w, func() {
			common.WriteComment(w, td.Comment)
		})
	})
}

func writeEnum(fw *common.MatlabFileWriter, enum *dsl.EnumDefinition, namedType *dsl.NamedType) error {
	enumName := common.TypeIdentifierName(enum.Name)
	if namedType != nil {
		enumName = common.TypeIdentifierName(namedType.Name)
	}
	return fw.WriteFile(enumName, func(w *formatting.IndentedWriter) {
		var base string
		if enum.BaseType == nil {
			base = "uint64"
		} else {
			base = common.TypeSyntax(enum.BaseType, enum.Namespace)
		}

		fmt.Fprintf(w, "classdef %s < %s\n", enumName, base)
		common.WriteBlockBody(w, func() {
			common.WriteComment(w, enum.Comment)
			w.WriteStringln("methods (Static)")
			common.WriteBlockBody(w, func() {
				for _, value := range enum.Values {
					common.WriteComment(w, value.Comment)
					fmt.Fprintf(w, "function v = %s\n", common.EnumValueIdentifierName(value.Symbol))
					common.WriteBlockBody(w, func() {
						fmt.Fprintf(w, "v = %s(%d);\n", common.TypeSyntax(enum, enum.Namespace), &value.IntegerValue)
					})
				}
				w.WriteStringln("")
				writeZerosStaticMethod(w, common.TypeSyntax(enum, enum.Namespace), []string{"0"})
			})

			if enum.IsFlags {
				// Additional methods for flag checks
				w.WriteStringln("")
				w.WriteStringln("methods")
				common.WriteBlockBody(w, func() {
					fmt.Fprintf(w, "function self = %s(varargin)\n", enumName)
					common.WriteBlockBody(w, func() {
						w.WriteStringln("if nargin == 0")
						w.Indented(func() {
							w.WriteStringln("value = 0;")
						})
						w.WriteStringln("elseif nargin == 1")
						w.Indented(func() {
							w.WriteStringln("value = varargin{1};")
						})
						w.WriteStringln("else")
						common.WriteBlockBody(w, func() {
							w.WriteStringln("value = 0;")
							w.WriteStringln("for i = 1:nargin")
							common.WriteBlockBody(w, func() {
								w.WriteStringln("value = bitor(value, varargin{i});")
							})
						})
						fmt.Fprintf(w, "self@%s(value);\n", base)
					})
					w.WriteStringln("")
					w.WriteStringln("function res = has_flags(self, flag)")
					common.WriteBlockBody(w, func() {
						w.WriteStringln("res = bitand(self, flag) == flag;")
					})
					w.WriteStringln("")
					w.WriteStringln("function res = with_flags(self, flag)")
					common.WriteBlockBody(w, func() {
						fmt.Fprintf(w, "res = %s(bitor(self, flag));\n", common.TypeSyntax(enum, enum.Namespace))
					})
					w.WriteStringln("")
					w.WriteStringln("function res = without_flags(self, flag)")
					common.WriteBlockBody(w, func() {
						fmt.Fprintf(w, "res = %s(bitand(self, bitcmp(flag)));\n", common.TypeSyntax(enum, enum.Namespace))
					})
				})
			}
		})
	})
}

func writeRecord(fw *common.MatlabFileWriter, rec *dsl.RecordDefinition, st dsl.SymbolTable) error {
	recordName := common.TypeIdentifierName(rec.Name)
	return fw.WriteFile(recordName, func(w *formatting.IndentedWriter) {
		fmt.Fprintf(w, "classdef %s < handle\n", recordName)
		common.WriteBlockBody(w, func() {
			common.WriteComment(w, rec.Comment)

			w.WriteStringln("properties")
			var fieldNames []string
			requiresArgsToConstruct := false
			common.WriteBlockBody(w, func() {
				for _, field := range rec.Fields {
					common.WriteComment(w, field.Comment)
					fieldName := common.FieldIdentifierName(field.Name)
					fieldNames = append(fieldNames, fieldName)
					w.WriteStringln(fieldName)
					_, defaultExpressionKind := typeDefault(field.Type, rec.Namespace, "", st)
					switch defaultExpressionKind {
					case defaultValueKindNone:
						requiresArgsToConstruct = true
					}
				}
			})
			w.WriteStringln("")

			w.WriteStringln("methods")
			common.WriteBlockBody(w, func() {

				fmt.Fprintf(w, "function self = %s(kwargs)\n", recordName)
				common.WriteBlockBody(w, func() {
					w.WriteStringln("arguments")
					common.WriteBlockBody(w, func() {
						for _, field := range rec.Fields {
							fieldName := common.FieldIdentifierName(field.Name)
							defaultExpression, defaultExpressionKind := typeDefault(field.Type, rec.Namespace, "", st)
							switch defaultExpressionKind {
							case defaultValueKindNone:
								fmt.Fprintf(w, "kwargs.%s;\n", fieldName)
							case defaultValueKindImmutable, defaultValueKindMutable:
								fmt.Fprintf(w, "kwargs.%s = %s;\n", fieldName, defaultExpression)
							}
						}
					})
					for _, field := range rec.Fields {
						fieldName := common.FieldIdentifierName(field.Name)
						_, defaultExpressionKind := typeDefault(field.Type, rec.Namespace, "", st)
						switch defaultExpressionKind {
						case defaultValueKindNone:
							fmt.Fprintf(w, "if ~isfield(kwargs, \"%s\")\n", fieldName)
							common.WriteBlockBody(w, func() {
								fmt.Fprintf(w, "throw(yardl.TypeError(\"Missing required keyword argument '%s'\"))\n", fieldName)
							})
						}
						fmt.Fprintf(w, "self.%s = kwargs.%s;\n", fieldName, fieldName)
					}
				})
				w.WriteStringln("")

				// Computed Fields
				if len(rec.ComputedFields) > 0 {
					for _, computedField := range rec.ComputedFields {
						fieldName := common.ComputedFieldIdentifierName(computedField.Name)

						fmt.Fprintf(w, "function res = %s(self)\n", fieldName)
						common.WriteBlockBody(w, func() {
							common.WriteComment(w, computedField.Comment)
							writeComputedFieldExpression(w, computedField.Expression, rec.Namespace)
						})
						w.WriteStringln("")
					}
					w.WriteStringln("")
				}

				// eq method
				w.WriteStringln("function res = eq(self, other)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("res = ...")
					w.Indented(func() {
						fmt.Fprintf(w, "isa(other, \"%s\")", common.TypeSyntax(rec, rec.Namespace))
						for _, field := range rec.Fields {
							w.WriteStringln(" && ...")
							fieldIdentifier := common.FieldIdentifierName(field.Name)
							w.WriteString(typeEqualityExpression(field.Type, "self."+fieldIdentifier, "other."+fieldIdentifier))
						}
						w.WriteStringln(";")
					})
				})
				w.WriteStringln("")

				// neq method
				w.WriteStringln("function res = ne(self, other)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("res = ~self.eq(other);")
				})
			})
			w.WriteStringln("")

			if !requiresArgsToConstruct {
				w.WriteStringln("methods (Static)")
				common.WriteBlockBody(w, func() {
					writeZerosStaticMethod(w, common.TypeSyntax(rec, rec.Namespace), []string{})
				})
			}
		})
	})
}

func writeZerosStaticMethod(w *formatting.IndentedWriter, typeSyntax string, defaultArgs []string) {
	// zeros method, only if can be constructed without arguments
	w.WriteStringln("function z = zeros(varargin)")
	common.WriteBlockBody(w, func() {
		fmt.Fprintf(w, "elem = %s(%s);\n", typeSyntax, strings.Join(defaultArgs, ", "))
		w.WriteStringln("if nargin == 0")
		common.WriteBlockBody(w, func() {
			w.WriteStringln("z = elem;")
			w.WriteStringln("return;")
		})

		w.WriteStringln("sz = [varargin{:}];")
		w.WriteStringln("if isscalar(sz)")
		common.WriteBlockBody(w, func() {
			w.WriteStringln("sz = [sz, sz];")
		})
		w.WriteStringln("z = reshape(repelem(elem, prod(sz)), sz);")
	})
}

func typeEqualityExpression(t dsl.Type, a, b string) string {
	if hasSimpleEquality(t) {
		return fmt.Sprintf("all([%s] == [%s])", a, b)
	}

	return fmt.Sprintf("isequal(%s, %s)", a, b)
}

func hasSimpleEquality(t dsl.Node) bool {
	res := true
	dsl.Visit(t, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case *dsl.SimpleType:
			self.Visit(t.ResolvedDefinition)
		case *dsl.Array, *dsl.GenericTypeParameter:
			res = false
			return
		}

		self.VisitChildren(node)
	})
	return res
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
					funcName := fmt.Sprintf("helper_%d_", len(helperFunctionLookup))
					helperFunctionLookup[arrType] = funcName
					fmt.Fprintf(w, "function dim = %s(dim_name)\n", funcName)
					common.WriteBlockBody(w, func() {
						dims := dsl.ToGeneralizedType(arrType).Dimensionality.(*dsl.Array).Dimensions
						for i, d := range *dims {
							fmt.Fprintf(w, "if dim_name == \"%s\"\n", *d.Name)
							w.Indented(func() {
								fmt.Fprintf(w, "dim = %d;\n", i)
							})
							w.WriteString("else")
						}
						w.WriteStringln("")
						common.WriteBlockBody(w, func() {
							w.WriteStringln(`throw(yardl.ValueError("Unknown dimension name: '%s'", dim_name));`)
						})
						w.WriteStringln("")
					})
				}
			}
		}
		self.VisitChildren(node)
	})

	varCounter := 0
	newVarName := func() string {
		varCounter++
		return fmt.Sprintf("var%d", varCounter)
	}

	tail := tailWrapper{}.Append(func(next func()) {
		w.WriteString("res = ")
		next()
		w.WriteStringln(";")
		w.WriteStringln("return")
	})

	dsl.VisitWithContext(expression, tail, func(self dsl.VisitorWithContext[tailWrapper], node dsl.Node, tail tailWrapper) {
		switch t := node.(type) {
		case *dsl.UnaryExpression:
			tail.Run(func() {
				if t.Operator != dsl.UnaryOpNegate {
					panic(fmt.Sprintf("unexpected unary operator %d", t.Operator))
				}
				w.WriteString("-(")
				self.Visit(t.Expression, tailWrapper{})
				w.WriteString(")")
			})

		case *dsl.BinaryExpression:
			tail.Run(func() {
				requiresParentheses := false
				if l, ok := t.Left.(*dsl.BinaryExpression); ok && l.Operator.Precedence() < t.Operator.Precedence() {
					requiresParentheses = true
				}

				if requiresParentheses {
					w.WriteString("(")
				}
				self.Visit(t.Left, tailWrapper{})
				if requiresParentheses {
					w.WriteString(")")
				}

				w.WriteString(" ")

				switch t.Operator {
				case dsl.BinaryOpAdd:
					w.WriteString("+")
				case dsl.BinaryOpSub:
					w.WriteString("-")
				case dsl.BinaryOpMul:
					w.WriteString(".*")
				case dsl.BinaryOpDiv:
					w.WriteString("./")
				case dsl.BinaryOpPow:
					w.WriteString("^")
				default:
					panic(fmt.Sprintf("unexpected binary operator %d", t.Operator))
				}

				w.WriteString(" ")

				requiresParentheses = false
				if r, ok := t.Right.(*dsl.BinaryExpression); ok && r.Operator.Precedence() < t.Operator.Precedence() {
					requiresParentheses = true
				}

				if requiresParentheses {
					w.WriteString("(")
				}
				self.Visit(t.Right, tailWrapper{})
				if requiresParentheses {
					w.WriteString(")")
				}
			})

		case *dsl.IntegerLiteralExpression:
			tail.Run(func() {
				fmt.Fprintf(w, "%d", &t.Value)
			})
		case *dsl.FloatingPointLiteralExpression:
			tail.Run(func() {
				w.WriteString(t.Value)
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

		case *dsl.SubscriptExpression:
			// Collapse adjacent subscript expressions (you can't do `array[x][y]` in Matlab, only `array(x,y)`)
			target := t.Target
			arguments := t.Arguments
			dsl.Visit(t.Target, func(self dsl.Visitor, node dsl.Node) {
				if t, ok := node.(*dsl.SubscriptExpression); ok {
					target = t.Target
					arguments = append(t.Arguments, arguments...)
					self.VisitChildren(t.Target)
				}
			})

			tail.Run(func() {
				self.Visit(target, tailWrapper{})

				startSubscript := "("
				delimeter := ", "
				dsl.Visit(target.GetResolvedType(), func(self dsl.Visitor, node dsl.Node) {
					switch t := node.(type) {
					case *dsl.GeneralizedType:
						switch t.Dimensionality.(type) {
						case *dsl.Vector, *dsl.Array:
							startSubscript = "(1+"
							delimeter = ", 1+"
							return
						}
					}
					self.VisitChildren(node)
				})

				w.WriteString(startSubscript)
				slices.Reverse(arguments)
				formatting.Delimited(w, delimeter, arguments, func(w *formatting.IndentedWriter, i int, a *dsl.SubscriptArgument) {
					self.Visit(a.Value, tailWrapper{})
				})
				w.WriteString(")")
			})

		case *dsl.FunctionCallExpression:
			tail.Run(func() {
				switch t.FunctionName {
				case dsl.FunctionSize:
					switch dsl.ToGeneralizedType(dsl.GetUnderlyingType(t.Arguments[0].GetResolvedType())).Dimensionality.(type) {
					case *dsl.Map:
						fmt.Fprintf(w, "numEntries(")
						self.Visit(t.Arguments[0], tailWrapper{})
						fmt.Fprintf(w, ")")

					case *dsl.Vector:
						fmt.Fprintf(w, "length(")
						self.Visit(t.Arguments[0], tailWrapper{})
						fmt.Fprintf(w, ")")

					case *dsl.Array:
						if len(t.Arguments) > 1 {
							w.WriteString("size(")
							self.Visit(t.Arguments[0], tailWrapper{})
							for _, arg := range t.Arguments[1:] {
								w.WriteString(", ")
								w.WriteString("ndims(")
								self.Visit(t.Arguments[0], tailWrapper{})
								w.WriteString(")-(")
								self.Visit(arg, tailWrapper{})
								w.WriteString(")")
							}
						} else {
							w.WriteString("numel(")
							self.Visit(t.Arguments[0], tailWrapper{})
						}
						w.WriteString(")")
					}

				case dsl.FunctionDimensionIndex:
					helperFuncName := helperFunctionLookup[t.Arguments[0].GetResolvedType()]
					fmt.Fprintf(w, "%s(", helperFuncName)
					self.Visit(t.Arguments[1], tailWrapper{})
					w.WriteString(")")

				case dsl.FunctionDimensionCount:
					w.WriteString("yardl.dimension_count(")
					self.Visit(t.Arguments[0], tailWrapper{})
					w.WriteString(")")

				default:
					panic(fmt.Sprintf("Unknown function '%s'", t.FunctionName))
				}
			})

		case *dsl.SwitchExpression:
			targetType := dsl.ToGeneralizedType(dsl.GetUnderlyingType(t.Target.GetResolvedType()))

			unionVariableName := newVarName()
			fmt.Fprintf(w, "%s = ", unionVariableName)
			self.Visit(t.Target, tailWrapper{})
			w.WriteStringln(";")

			if targetType.Cases.IsOptional() {
				for i, switchCase := range t.Cases {
					writeSwitchCaseOverOptional(w, switchCase, unionVariableName, i == len(targetType.Cases)-1, self, tail)
				}
				return
			}

			if targetType.Cases.IsUnion() {
				// Special handling for SwitchExpression over a Union from an imported namespace
				targetTypeNamespace := ""
				dsl.Visit(t.Target, func(self dsl.Visitor, node dsl.Node) {
					switch node := node.(type) {
					case *dsl.SimpleType:
						self.Visit(node.ResolvedDefinition)
					case *dsl.RecordDefinition:
						for _, field := range node.Fields {
							u := dsl.GetUnderlyingType(field.Type)
							if u == t.Target.GetResolvedType() {
								if targetTypeNamespace == "" {
									meta := node.GetDefinitionMeta()
									targetTypeNamespace = meta.Namespace
								}
								return
							}
						}
					case dsl.Expression:
						t := node.GetResolvedType()
						if t != nil {
							self.Visit(t)
						}
					}
					self.VisitChildren(node)
				})
				if targetTypeNamespace == "" {
					targetTypeNamespace = contextNamespace
				}

				unionClassName := fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(targetTypeNamespace), common.UnionClassName(targetType))
				for _, switchCase := range t.Cases {
					writeSwitchCaseOverUnion(w, targetType, unionClassName, switchCase, unionVariableName, self, tail)
				}

				if _, ok := t.Cases[len(t.Cases)-1].Pattern.(*dsl.DiscardPattern); !ok {
					w.WriteStringln(`throw(yardl.RuntimeError("Unexpected union case"))`)
				}
				return
			}

			// this is over a single type
			if len(t.Cases) != 1 {
				panic("switch expression over a single type expected to have exactly one case")
			}
			switchCase := t.Cases[0]
			switch pattern := switchCase.Pattern.(type) {
			case *dsl.DeclarationPattern:
				fmt.Fprintf(w, "%s = %s;\n", common.FieldIdentifierName(pattern.Identifier), unionVariableName)
				self.Visit(switchCase.Expression, tail)
			case *dsl.TypePattern, *dsl.DiscardPattern:
				self.Visit(switchCase.Expression, tail)
			default:
				panic(fmt.Sprintf("Unexpected pattern type %T", t.Cases[0].Pattern))
			}

		case *dsl.TypeConversionExpression:
			tail = tail.Append(func(next func()) {
				writeTypeConversion(w, t.Type, next)
			})

			self.Visit(t.Expression, tail)

		default:
			panic(fmt.Sprintf("Unknown expression type '%T'", t))
		}
	})
}

func writeSwitchCaseOverOptional(w *formatting.IndentedWriter, switchCase *dsl.SwitchCase, variableName string, isLastCase bool, visitor dsl.VisitorWithContext[tailWrapper], tail tailWrapper) {
	writeCore := func(declarationIdentifier string) {
		if declarationIdentifier != "" {
			fmt.Fprintf(w, "%s = %s;\n", declarationIdentifier, variableName)
		}
		visitor.Visit(switchCase.Expression, tail)
	}

	writeTypeCase := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		if isLastCase {
			writeCore(declarationIdentifier)
			return
		}

		if typePattern.Type == nil {
			fmt.Fprintf(w, "if %s == yardl.None\n", variableName)
		} else {
			fmt.Fprintf(w, "if %s ~= yardl.None\n", variableName)
		}

		common.WriteBlockBody(w, func() {
			writeCore(declarationIdentifier)
		})
	}

	switch t := switchCase.Pattern.(type) {
	case *dsl.TypePattern:
		writeTypeCase(t, "")
	case *dsl.DeclarationPattern:
		writeTypeCase(&t.TypePattern, common.FieldIdentifierName(t.Identifier))
	case *dsl.DiscardPattern:
		writeCore("")
	default:
		panic(fmt.Sprintf("Unknown pattern type '%T'", switchCase.Pattern))
	}
}

func writeSwitchCaseOverUnion(w *formatting.IndentedWriter, unionType *dsl.GeneralizedType, unionClassName string, switchCase *dsl.SwitchCase, variableName string, visitor dsl.VisitorWithContext[tailWrapper], tail tailWrapper) {
	caseIndexOffset := 1
	writeTypeCase := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		for i, typeCase := range unionType.Cases {
			if typeCase.Type == nil {
				caseIndexOffset = 0
			}

			if dsl.TypesEqual(typePattern.Type, typeCase.Type) {
				if typePattern.Type == nil {
					fmt.Fprintf(w, "if %s == yardl.None\n", variableName)
					common.WriteBlockBody(w, func() {
						visitor.Visit(switchCase.Expression, tail)
					})
				} else {
					fmt.Fprintf(w, "if isa(%s, \"%s\") && %s.index == %d\n", variableName, unionClassName, variableName, i+caseIndexOffset)
					common.WriteBlockBody(w, func() {
						if declarationIdentifier != "" {
							fmt.Fprintf(w, "%s = %s.value;\n", declarationIdentifier, variableName)
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

func writeTypeConversion(w *formatting.IndentedWriter, t dsl.Type, next func()) {
	getWrapper := func(t dsl.Type) (string, string) {
		switch t := t.(type) {
		case *dsl.SimpleType:
			switch t := t.ResolvedDefinition.(type) {
			case dsl.PrimitiveDefinition:
				switch t {
				case dsl.Bool:
					return "logical(", ")"
				case dsl.Int8:
					return "int8(", ")"
				case dsl.Uint8:
					return "uint8(", ")"
				case dsl.Int16:
					return "int32(", ")"
				case dsl.Uint16:
					return "uint16(", ")"
				case dsl.Int32:
					return "int32(", ")"
				case dsl.Uint32:
					return "uint32(", ")"
				case dsl.Int64:
					return "int64(", ")"
				case dsl.Uint64, dsl.Size:
					return "uint64(", ")"
				case dsl.Float32:
					return "single(", ")"
				case dsl.Float64:
					return "double(", ")"
				case dsl.ComplexFloat32:
					return "complex(single(", "))"
				case dsl.ComplexFloat64:
					return "complex(double(", "))"
				case dsl.String:
					return "string(", ")"
				case dsl.Date:
					return "yardl.Date(", ")"
				case dsl.Time:
					return "yardl.Time(", ")"
				case dsl.DateTime:
					return "yardl.DateTime(", ")"
				}
			}
		}
		panic(fmt.Sprintf("Unsupported type '%s'", t))
	}

	open, close := getWrapper(t)
	w.WriteString(open)
	next()
	w.WriteString(close)
}

type defaultValueKind int

const (
	defaultValueKindNone defaultValueKind = iota
	defaultValueKindImmutable
	defaultValueKindMutable
)

func typeDefault(t dsl.Type, contextNamespace string, namedType string, st dsl.SymbolTable) (string, defaultValueKind) {
	switch t := t.(type) {
	case nil:
		return "yardl.None", defaultValueKindImmutable
	case *dsl.SimpleType:
		return typeDefinitionDefault(t.ResolvedDefinition, contextNamespace, st)
	case *dsl.GeneralizedType:
		switch td := t.Dimensionality.(type) {
		case nil:
			defaultExpression, defaultKind := typeDefault(t.Cases[0].Type, contextNamespace, "", st)
			if t.Cases.IsSingle() || t.Cases.HasNullOption() {
				return defaultExpression, defaultKind
			}

			var unionClassName string
			if namedType != "" {
				unionClassName = namedType
			} else {
				unionClassName = fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(contextNamespace), common.UnionClassName(t))
			}

			unionCaseConstructor := fmt.Sprintf("%s.%s", unionClassName, formatting.ToPascalCase(t.Cases[0].Tag))

			switch defaultKind {
			case defaultValueKindNone:
				return "", defaultKind
			case defaultValueKindImmutable:
				return fmt.Sprintf(`%s(%s)`, unionCaseConstructor, defaultExpression), defaultKind
			case defaultValueKindMutable:
				if t, ok := dsl.GetUnderlyingType(t.Cases[0].Type).(*dsl.SimpleType); ok {
					if _, ok := t.ResolvedDefinition.(*dsl.RecordDefinition); ok {
						return fmt.Sprintf(`%s(%s())`, unionCaseConstructor, defaultExpression), defaultValueKindMutable
					}
				}
				return fmt.Sprintf(`%s(%s)`, unionCaseConstructor, defaultExpression), defaultValueKindMutable
			}

			return fmt.Sprintf(`("%s", %s)`, t.Cases[0].Tag, defaultExpression), defaultValueKindImmutable

		case *dsl.Vector:
			scalar := t.ToScalar()
			if dsl.TypeContainsGenericTypeParameter(scalar) {
				return "", defaultValueKindNone
			}

			dtype := common.TypeSyntax(scalar, contextNamespace)
			if td.Length == nil {
				return fmt.Sprintf("%s.empty()", dtype), defaultValueKindMutable
			}

			scalarDefault, scalarDefaultKind := typeDefault(t.Cases[0].Type, contextNamespace, "", st)
			switch scalarDefaultKind {
			case defaultValueKindNone:
				return "", defaultValueKindNone
			case defaultValueKindImmutable, defaultValueKindMutable:
				if gt, ok := t.Cases[0].Type.(*dsl.GeneralizedType); ok {
					switch gt.Dimensionality.(type) {
					case *dsl.Vector, *dsl.Array:
						return fmt.Sprintf("repelem({%s}, %d)", scalarDefault, *td.Length), scalarDefaultKind
					}
				}
				return fmt.Sprintf("repelem(%s, %d)", scalarDefault, *td.Length), defaultValueKindMutable
			}

		case *dsl.Array:
			scalar := t.ToScalar()
			if dsl.TypeContainsGenericTypeParameter(scalar) {
				return "", defaultValueKindNone
			}

			if td.IsFixed() {
				dims := make([]string, len(*td.Dimensions))
				for i, d := range *td.Dimensions {
					dims[len(*td.Dimensions)-i-1] = strconv.FormatUint(*d.Length, 10)
				}
				if len(dims) == 1 {
					dims = append(dims, "1")
				}

				scalarDefault, _ := typeDefault(scalar, contextNamespace, "", st)

				return fmt.Sprintf("repelem(%s, %s)", scalarDefault, strings.Join(dims, ", ")), defaultValueKindMutable
			}

			dtype := common.TypeSyntax(scalar, contextNamespace)
			if td.HasKnownNumberOfDimensions() {
				shape := strings.Repeat("0, ", len(*td.Dimensions))[0 : len(*td.Dimensions)*3-2]
				return fmt.Sprintf("%s.empty(%s)", dtype, shape), defaultValueKindMutable
			}
			return fmt.Sprintf("%s.empty()", dtype), defaultValueKindMutable

		case *dsl.Map:
			return "dictionary", defaultValueKindMutable
		}
	}

	return "", defaultValueKindNone
}

func typeDefinitionDefault(t dsl.TypeDefinition, contextNamespace string, st dsl.SymbolTable) (string, defaultValueKind) {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		switch t {
		case dsl.Bool:
			return "false", defaultValueKindImmutable
		case dsl.Int8:
			return "int8(0)", defaultValueKindImmutable
		case dsl.Uint8:
			return "uint8(0)", defaultValueKindImmutable
		case dsl.Int16:
			return "int16(0)", defaultValueKindImmutable
		case dsl.Uint16:
			return "uint16(0)", defaultValueKindImmutable
		case dsl.Int32:
			return "int32(0)", defaultValueKindImmutable
		case dsl.Uint32:
			return "uint32(0)", defaultValueKindImmutable
		case dsl.Int64:
			return "int64(0)", defaultValueKindImmutable
		case dsl.Uint64, dsl.Size:
			return "uint64(0)", defaultValueKindImmutable
		case dsl.Float32:
			return "single(0)", defaultValueKindImmutable
		case dsl.Float64:
			return "double(0)", defaultValueKindImmutable
		case dsl.ComplexFloat32:
			return "complex(single(0))", defaultValueKindImmutable
		case dsl.ComplexFloat64:
			return "complex(0)", defaultValueKindImmutable
		case dsl.String:
			return `""`, defaultValueKindImmutable
		case dsl.Date:
			return "yardl.Date()", defaultValueKindImmutable
		case dsl.Time:
			return "yardl.Time()", defaultValueKindImmutable
		case dsl.DateTime:
			return "yardl.DateTime()", defaultValueKindImmutable
		}

	case *dsl.EnumDefinition:
		zeroValue := t.GetZeroValue()
		if t.IsFlags {
			if zeroValue == nil {
				return fmt.Sprintf("%s(0)", common.TypeSyntax(t, contextNamespace)), defaultValueKindImmutable
			} else {
				return fmt.Sprintf("%s.%s", common.TypeSyntax(t, contextNamespace), common.EnumValueIdentifierName(zeroValue.Symbol)), defaultValueKindImmutable
			}
		}

		if zeroValue == nil {
			return "", defaultValueKindNone
		}

		return fmt.Sprintf("%s.%s", common.TypeSyntax(t, contextNamespace), common.EnumValueIdentifierName(zeroValue.Symbol)), defaultValueKindImmutable

	case *dsl.NamedType:
		return typeDefault(t.Type, contextNamespace, common.TypeSyntax(t, contextNamespace), st)

	case *dsl.RecordDefinition:
		if len(t.TypeArguments) == 0 {
			if len(t.TypeParameters) > 0 {
				// *Open* Generic Record type
				// Should never get here - typeDefault is only called on Fields, which must be closed if generic
				panic(fmt.Sprintf("No typeDefault for open generic record %s", t.Name))
			}

			for _, f := range t.Fields {
				_, fieldDefaultKind := typeDefault(f.Type, contextNamespace, "", st)
				if fieldDefaultKind == defaultValueKindNone {
					// Basic, closed record type
					// Should never get here - a Field in a closed record should always have a default type
					panic(fmt.Sprintf("No typeDefault for record field %s.%s", t.Name, f.Name))
				}
			}

			// Basic record type
			return fmt.Sprintf("%s()", common.TypeSyntax(t, contextNamespace)), defaultValueKindMutable
		}

		// t is a *closed* generic record type
		// genericDef is its original generic type definition
		genericDef := st[t.GetQualifiedName()].(*dsl.RecordDefinition)
		args := make([]string, 0)
		for i, f := range t.Fields {
			fieldDefaultExpr, fieldDefaultKind := typeDefault(f.Type, contextNamespace, "", st)
			if fieldDefaultKind == defaultValueKindNone {
				return "", defaultValueKindNone
			}

			// Only write a constructor argument if it is needed, e.g. the record definition's field is generic and doesn't have a default value
			_, genDefaultKind := typeDefault(genericDef.Fields[i].Type, contextNamespace, "", st)
			if genDefaultKind == defaultValueKindNone {
				args = append(args, fmt.Sprintf("%s=%s", common.FieldIdentifierName(f.Name), fieldDefaultExpr))
			}
		}

		return fmt.Sprintf("%s(%s)", common.TypeSyntax(t, contextNamespace), strings.Join(args, ", ")), defaultValueKindMutable
	}

	return "", defaultValueKindNone
}
