// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

func Compare[T any](newRoot, oldRoot Node, context T, compareFunc CompareFunc[T]) error {
	compareVisitor := CompareVisitor[T]{compareFunc: compareFunc}
	return compareFunc(compareVisitor, newRoot, oldRoot, context)
}

type CompareFunc[T any] func(self CompareVisitor[T], newRoot, oldRoot Node, context T) error

type CompareVisitor[T any] struct {
	compareFunc CompareFunc[T]
}

func (cv CompareVisitor[T]) Compare(newRoot, oldRoot Node, context T) error {
	return cv.compareFunc(cv, newRoot, oldRoot, context)
}
