// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"math/big"
	"strings"
)

// ----------------------------------------------------------------------------
// Node
type Node interface {
	GetNodeMeta() *NodeMeta
}

type NodeMeta struct {
	File   string `json:"-"`
	Line   int    `json:"-"`
	Column int    `json:"-"`
}

func (n *NodeMeta) String() string {
	return fmt.Sprintf("%s:%d:%d", n.File, n.Line, n.Column)
}

func (n *NodeMeta) Equals(other *NodeMeta) bool {
	return n == other || (n != nil && other != nil &&
		n.File == other.File &&
		n.Line == other.Line &&
		n.Column == other.Column)
}

func (n *NodeMeta) GetNodeMeta() *NodeMeta {
	return n
}

// ----------------------------------------------------------------------------
// Environment + Namespace
type Environment struct {
	Namespaces  []*Namespace `json:"namespaces"`
	SymbolTable SymbolTable  `json:"-"`
}

func (n *Environment) GetNodeMeta() *NodeMeta {
	return &NodeMeta{}
}

type SymbolTable map[string]TypeDefinition

func (st SymbolTable) GetGenericTypeDefinition(possiblyGenericType TypeDefinition) TypeDefinition {
	meta := possiblyGenericType.GetDefinitionMeta()
	if len(meta.TypeParameters) == 0 {
		return possiblyGenericType
	}

	genericTypeDefinition, ok := st[meta.GetQualifiedName()]
	if !ok {
		panic("unable to find generic type definition in symbol table")
	}

	return genericTypeDefinition
}

func (st SymbolTable) Clone() SymbolTable {
	newSymbolTable := make(SymbolTable, len(st))
	for k, v := range st {
		newSymbolTable[k] = v
	}

	return newSymbolTable
}

type Namespace struct {
	Name            string                `json:"name"`
	TypeDefinitions TypeDefinitions       `json:"types,omitempty"`
	Protocols       []*ProtocolDefinition `json:"protocols,omitempty"`
}

func (n *Namespace) GetNodeMeta() *NodeMeta {
	return &NodeMeta{}
}

// ----------------------------------------------------------------------------
// TypeDefinition
type TypeDefinition interface {
	Node
	GetDefinitionMeta() *DefinitionMeta
}
type TypeDefinitions []TypeDefinition

// ----------------------------------------------------------------------------
// Primitive types
type PrimitiveDefinition string

const (
	Bool           = "bool"
	Int8           = "int8"
	Uint8          = "uint8"
	Int16          = "int16"
	Uint16         = "uint16"
	Int32          = "int32"
	Uint32         = "uint32"
	Int64          = "int64"
	Uint64         = "uint64"
	Size           = "size"
	Float32        = "float32"
	Float64        = "float64"
	ComplexFloat32 = "complexfloat32"
	ComplexFloat64 = "complexfloat64"
	String         = "string"
	Date           = "date"
	Time           = "time"
	DateTime       = "datetime"
)

var (
	PrimitiveBool           = PrimitiveDefinition(Bool)
	PrimitiveInt8           = PrimitiveDefinition(Int8)
	PrimitiveUint8          = PrimitiveDefinition(Uint8)
	PrimitiveInt16          = PrimitiveDefinition(Int16)
	PrimitiveUint16         = PrimitiveDefinition(Uint16)
	PrimitiveInt32          = PrimitiveDefinition(Int32)
	PrimitiveUint32         = PrimitiveDefinition(Uint32)
	PrimitiveInt64          = PrimitiveDefinition(Int64)
	PrimitiveUint64         = PrimitiveDefinition(Uint64)
	PrimitiveSize           = PrimitiveDefinition(Size)
	PrimitiveFloat32        = PrimitiveDefinition(Float32)
	PrimitiveFloat64        = PrimitiveDefinition(Float64)
	PrimitiveComplexFloat32 = PrimitiveDefinition(ComplexFloat32)
	PrimitiveComplexFloat64 = PrimitiveDefinition(ComplexFloat64)
	PrimitiveString         = PrimitiveDefinition(String)
	PrimitiveDate           = PrimitiveDefinition(Date)
	PrimitiveTime           = PrimitiveDefinition(Time)
	PrimitiveDateTime       = PrimitiveDefinition(DateTime)
)

var primitiveTypes = func(names ...string) map[string]PrimitiveDefinition {
	m := make(map[string]PrimitiveDefinition)
	for _, v := range names {
		toks := strings.Split(v, "=")
		if len(toks) == 1 {
			m[v] = PrimitiveDefinition(v)
		} else {
			m[toks[0]] = PrimitiveDefinition(toks[1])
		}
	}
	return m
}(
	Bool,
	Int8,
	Uint8,
	Int16,
	Uint16,
	Int32,
	Uint32,
	Int64,
	Uint64,
	Size,
	Float32,
	Float64,
	ComplexFloat32,
	ComplexFloat64,
	String,
	Date,
	Time,
	DateTime,

	// aliases
	"byte="+Uint8,
	"int="+Int32,
	"uint="+Uint32,
	"long="+Int64,
	"ulong="+Uint64,
	"float="+Float32,
	"double="+Float64,
	"complexfloat="+ComplexFloat32,
	"complexdouble="+ComplexFloat64,
)

func (n PrimitiveDefinition) GetNodeMeta() *NodeMeta {
	return &NodeMeta{}
}

func (t PrimitiveDefinition) GetDefinitionMeta() *DefinitionMeta {
	return &DefinitionMeta{Name: string(t)}
}

// ----------------------------------------------------------------------------
// DefinitionMeta
type DefinitionMeta struct {
	NodeMeta
	Name           string                  `json:"name"`
	Namespace      string                  `json:"-"`
	TypeParameters []*GenericTypeParameter `json:"typeParameters,omitempty"`
	TypeArguments  []Type                  `json:"typeArguments,omitempty"`
	Comment        string                  `json:"comment,omitempty"`
}

func (meta *DefinitionMeta) GetQualifiedName() string {
	if meta.Namespace == "" {
		return meta.Name
	}

	return fmt.Sprintf("%s.%s", meta.Namespace, meta.Name)
}

// ----------------------------------------------------------------------------
// Named types (arrays etc.)
type NamedType struct {
	*DefinitionMeta
	Type `json:"type"`
}

func (n *NamedType) GetNodeMeta() *NodeMeta {
	return &n.NodeMeta
}

func (nt *NamedType) GetDefinitionMeta() *DefinitionMeta {
	return nt.DefinitionMeta
}

// ----------------------------------------------------------------------------
// Records
type RecordDefinition struct {
	*DefinitionMeta
	Fields         Fields         `json:"fields"`
	ComputedFields ComputedFields `json:"computedFields,omitempty"`
}

func (r *RecordDefinition) GetDefinitionMeta() *DefinitionMeta {
	return r.DefinitionMeta
}

// ----------------------------------------------------------------------------
// Generics
type GenericTypeParameter struct {
	NodeMeta
	Name string
}

func (p *GenericTypeParameter) GetDefinitionMeta() *DefinitionMeta {
	return &DefinitionMeta{Name: p.Name}
}

// ----------------------------------------------------------------------------
// Fields
type Fields []*Field

type Field struct {
	NodeMeta
	Name    string `json:"name"`
	Comment string `json:"comment,omitempty"`
	Type    Type   `json:"type"`
}

// ----------------------------------------------------------------------------
// Types

type ArrayDimensions []*ArrayDimension

type ArrayDimension struct {
	NodeMeta
	Comment string  `json:"comment,omitempty"`
	Name    *string `json:"name,omitempty"`
	Length  *uint64 `json:"length,omitempty"`
}

type Dimensionality interface {
	Node
	dimensionality()
}

type Vector struct {
	NodeMeta
	Length *uint64 `json:"length,omitempty"`
}

func (v *Vector) dimensionality() {}

func (v *Vector) IsFixed() bool {
	return v.Length != nil
}

type Array struct {
	NodeMeta
	Dimensions *ArrayDimensions `json:"dimensions,omitempty"`
}

func (a *Array) dimensionality() {}

func (a *Array) HasKnownNumberOfDimensions() bool {
	return a.Dimensions != nil
}

func (a *Array) IsFixed() bool {
	if a.Dimensions == nil {
		return false
	}

	for _, d := range *a.Dimensions {
		if d.Length == nil {
			return false
		}
	}

	return true
}

// Map implements Dimensionality on a GeneralizedType.
// The Value type of the map is the Cases of the GeneralizedType.
type Map struct {
	NodeMeta
	KeyType Type `json:"keyType"`
}

func (m *Map) dimensionality() {}

type Stream struct {
	NodeMeta
}

func (s *Stream) dimensionality() {}

// ----------------------------------------------------------------------------
// Type
type Type interface {
	Node
	_type()
}

type SimpleType struct {
	NodeMeta
	Name               string
	TypeArguments      []Type
	ResolvedDefinition TypeDefinition
}

type GeneralizedType struct {
	NodeMeta
	Cases          TypeCases
	Dimensionality Dimensionality
}

func (t *GeneralizedType) _type() {}

func (t *GeneralizedType) ToScalar() Type {
	if t.Cases.IsSingle() {
		if simpleType, ok := t.Cases[0].Type.(*SimpleType); ok {
			return simpleType
		}
	}

	return &GeneralizedType{NodeMeta: t.NodeMeta, Cases: t.Cases}
}

func (t *SimpleType) _type() {}

type TypeCases []*TypeCase

func (tcs *TypeCases) IsSingle() bool {
	return len(*tcs) == 1
}

func (tcs *TypeCases) IsOptional() bool {
	return len(*tcs) == 2 && (*tcs)[0].IsNullType()
}

func (tcs *TypeCases) IsUnion() bool {
	return len(*tcs) > 2 || (len(*tcs) == 2 && !(*tcs)[0].IsNullType())
}

func (tcs *TypeCases) HasNullOption() bool {
	return len(*tcs) > 1 && (*tcs)[0].IsNullType()
}

type TypeCase struct {
	NodeMeta
	Label string `json:"label,omitempty"`
	Type  Type   `json:"type"`
}

func (tc *TypeCase) IsNullType() bool {
	return tc.Type == nil
}

// ----------------------------------------------------------------------------
// Enums
type EnumDefinition struct {
	*DefinitionMeta `yaml:"-,inline"`
	BaseType        Type       `json:"base,omitempty"`
	Values          EnumValues `json:"values"`
}

func (e *EnumDefinition) GetDefinitionMeta() *DefinitionMeta {
	return e.DefinitionMeta
}

type EnumValues []*EnumValue

type EnumValue struct {
	NodeMeta
	Symbol       string  `json:"symbol"`
	Comment      string  `json:"comment,omitempty"`
	IntegerValue big.Int `json:"value"`
}

// ----------------------------------------------------------------------------
// Protocols
type ProtocolDefinition struct {
	*DefinitionMeta
	Sequence ProtocolSteps `json:"sequence"`
}

func (p *ProtocolDefinition) GetDefinitionMeta() *DefinitionMeta {
	return p.DefinitionMeta
}

type ProtocolSteps []*ProtocolStep

type ProtocolStep Field

func (s *ProtocolStep) IsStream() bool {
	if gt, ok := s.Type.(*GeneralizedType); ok {
		_, isStream := gt.Dimensionality.(*Stream)
		return isStream
	}

	return false
}

// ----------------------------------------------------------------------------
// Computed fields

type ComputedFields []*ComputedField

type ComputedField struct {
	NodeMeta
	Name       string     `json:"name"`
	Comment    string     `json:"comment,omitempty"`
	Expression Expression `json:"expression"`
}

type Expression interface {
	Node
	GetResolvedType() Type
	IsReference() bool
	_expression()
}

type IntegerLiteralExpression struct {
	NodeMeta
	Value        big.Int
	ResolvedType Type
}

func (e *IntegerLiteralExpression) _expression() {}
func (e *IntegerLiteralExpression) GetResolvedType() Type {
	return e.ResolvedType
}
func (e *IntegerLiteralExpression) IsReference() bool {
	return false
}

type StringLiteralExpression struct {
	NodeMeta
	Value        string
	ResolvedType Type
}

func (e *StringLiteralExpression) _expression() {}
func (e *StringLiteralExpression) GetResolvedType() Type {
	return e.ResolvedType
}
func (e *StringLiteralExpression) IsReference() bool {
	return false
}

type MemberAccessExpression struct {
	NodeMeta
	Target          Expression `json:"target,omitempty"`
	Member          string     `json:"member,omitempty"`
	IsComputedField bool       `json:"isComputed,omitempty"`
	ResolvedType    Type       `json:"-"`
}

func (e *MemberAccessExpression) _expression() {}
func (e *MemberAccessExpression) GetResolvedType() Type {
	return e.ResolvedType
}
func (e *MemberAccessExpression) IsReference() bool {
	return true
}

type IndexExpression struct {
	NodeMeta
	Target       Expression       `json:"target,omitempty"`
	Arguments    []*IndexArgument `json:"arguments"`
	ResolvedType Type             `json:"-"`
}

func (e *IndexExpression) _expression() {}
func (e *IndexExpression) GetResolvedType() Type {
	return e.ResolvedType
}
func (e *IndexExpression) IsReference() bool {
	return true
}

type IndexArgument struct {
	NodeMeta
	Label string     `json:"label,omitempty"`
	Value Expression `json:"expression"`
}

type FunctionCallExpression struct {
	NodeMeta
	FunctionName string       `json:"function"`
	Arguments    []Expression `json:"arguments"`
	ResolvedType Type         `json:"-"`
}

func (e *FunctionCallExpression) _expression() {}
func (e *FunctionCallExpression) GetResolvedType() Type {
	return e.ResolvedType
}
func (e *FunctionCallExpression) IsReference() bool {
	return false
}

const (
	FunctionSize           = "size"
	FunctionDimensionIndex = "dimensionIndex"
	FunctionDimensionCount = "dimensionCount"
)

type TypeConversionExpression struct {
	NodeMeta
	Expression Expression `json:"expression"`
	Type       Type       `json:"type"`
}

func (e *TypeConversionExpression) _expression() {}
func (e *TypeConversionExpression) GetResolvedType() Type {
	return e.Type
}
func (e *TypeConversionExpression) IsReference() bool {
	return false
}

type SwitchExpression struct {
	NodeMeta
	Target       Expression    `json:"target"`
	Cases        []*SwitchCase `json:"cases"`
	ResolvedType Type          `json:"-"`
}

func (e *SwitchExpression) _expression() {}
func (e *SwitchExpression) GetResolvedType() Type {
	return e.ResolvedType
}
func (e *SwitchExpression) IsReference() bool {
	return false
}

type SwitchCase struct {
	NodeMeta
	Pattern    Pattern    `json:"pattern"`
	Expression Expression `json:"expression"`
}

type Pattern interface {
	Node
	_pattern()
}

type TypePattern struct {
	NodeMeta
	Type Type `json:"type"`
}

func (p *TypePattern) _pattern() {}

type DeclarationPattern struct {
	TypePattern
	Identifier string `json:"identifier"`
}

type DiscardPattern struct {
	NodeMeta
}

func (p *DiscardPattern) _pattern() {}

// interface implementation checks
var (
	_ Node = (*Environment)(nil)
	_ Node = (*Namespace)(nil)
	_ Node = (*DefinitionMeta)(nil)
	_ Node = TypeDefinition(nil)
	_ Node = PrimitiveDefinition("")
	_ Node = (*NamedType)(nil)
	_ Node = (*RecordDefinition)(nil)
	_ Node = (*Field)(nil)
	_ Node = (*ArrayDimension)(nil)
	_ Node = Dimensionality(nil)
	_ Node = (*Vector)(nil)
	_ Node = (*Array)(nil)
	_ Node = (*Stream)(nil)
	_ Node = Type(nil)
	_ Node = (*TypeCase)(nil)
	_ Node = (*GeneralizedType)(nil)
	_ Node = (*EnumDefinition)(nil)
	_ Node = (*ProtocolDefinition)(nil)
	_ Node = (*ProtocolStep)(nil)
	_ Node = (*ComputedField)(nil)

	_ TypeDefinition = (*RecordDefinition)(nil)
	_ TypeDefinition = (*EnumDefinition)(nil)
	_ TypeDefinition = (PrimitiveDefinition)("")
	_ TypeDefinition = (*NamedType)(nil)
	_ TypeDefinition = (*ProtocolDefinition)(nil)
	_ TypeDefinition = (*GenericTypeParameter)(nil)

	_ Dimensionality = (*Vector)(nil)
	_ Dimensionality = (*Array)(nil)
	_ Dimensionality = (*Map)(nil)
	_ Dimensionality = (*Stream)(nil)

	_ Type = (*SimpleType)(nil)
	_ Type = (*GeneralizedType)(nil)

	_ Expression = (*IntegerLiteralExpression)(nil)
	_ Expression = (*StringLiteralExpression)(nil)
	_ Expression = (*MemberAccessExpression)(nil)
	_ Expression = (*IndexExpression)(nil)
	_ Expression = (*FunctionCallExpression)(nil)
	_ Expression = (*TypeConversionExpression)(nil)

	_ Pattern = (*TypePattern)(nil)
	_ Pattern = (*DeclarationPattern)(nil)
	_ Pattern = (*DiscardPattern)(nil)
)
