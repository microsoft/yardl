// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

type TypeSyntaxWriter[TContext any] func(self TypeSyntaxWriter[TContext], typeOrTypeDef Node, context TContext) string

func (writer TypeSyntaxWriter[TContext]) ToSyntax(typeOrTypeDef Node, context TContext) string {
	return writer(writer, typeOrTypeDef, context)
}

func MkTypeSyntaxRewriter[TContext any](writer TypeSyntaxWriter[TContext]) func(typeOrTypeDef Node, context TContext) string {
	return func(typeOrTypeDef Node, context TContext) string {
		return writer(writer, typeOrTypeDef, context)
	}
}
