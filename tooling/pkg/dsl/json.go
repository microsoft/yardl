// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"encoding/json"
	"fmt"
)

// Customize JSON marshaling for DSL types

func (t *SimpleType) MarshalJSON() ([]byte, error) {
	if len(t.TypeArguments) == 0 {
		return json.Marshal(t.Name)
	}
	type expanded struct {
		Name          string `json:"name"`
		TypeArguments []Type `json:"typeArguments,omitempty"`
	}
	return json.Marshal(expanded{Name: t.Name, TypeArguments: t.TypeArguments})
}

func (t *GenericTypeParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Name)
}

func (t *GeneralizedType) MarshalJSON() ([]byte, error) {
	switch d := t.Dimensionality.(type) {
	case nil:
		return json.Marshal(t.Cases)
	case *Vector:
		type vectorView struct {
			Items  TypeCases `json:"items"`
			Length *uint64   `json:"length,omitempty"`
		}
		type vecWrapper struct {
			Vector vectorView `json:"vector"`
		}

		return json.Marshal(vecWrapper{Vector: vectorView{Items: t.Cases, Length: d.Length}})
	case *Array:
		type arrayView struct {
			Items      TypeCases        `json:"items"`
			Dimensions *ArrayDimensions `json:"dimensions,omitempty"`
		}
		type arrayWrapper struct {
			Array arrayView `json:"array"`
		}

		return json.Marshal(arrayWrapper{Array: arrayView{Items: t.Cases, Dimensions: d.Dimensions}})

	case *Stream:
		type streamView struct {
			Items TypeCases `json:"items"`
		}
		type streamWrapper struct {
			Stream streamView `json:"stream"`
		}

		return json.Marshal(streamWrapper{Stream: streamView{Items: t.Cases}})
	case *Map:
		type mapView struct {
			Keys   Type      `json:"keys"`
			Values TypeCases `json:"values"`
		}
		type mapWrapper struct {
			Map mapView `json:"map"`
		}

		return json.Marshal(mapWrapper{Map: mapView{Keys: d.KeyType, Values: t.Cases}})

	default:
		panic(fmt.Sprintf("unexpected type %T", d))
	}
}

func (tcs TypeCases) MarshalJSON() ([]byte, error) {
	if len(tcs) == 1 {
		return json.Marshal(tcs[0])
	}

	type Alias TypeCases
	return json.Marshal(Alias(tcs))
}

func (tc *TypeCase) MarshalJSON() ([]byte, error) {
	if tc.IsNullType() {
		return json.Marshal(nil)
	}
	if tc.Label == "" {
		return json.Marshal(tc.Type)
	}

	type Alias TypeCase
	return json.Marshal((*Alias)(tc))
}

func (dims ArrayDimensions) MarshalJSON() ([]byte, error) {
	for _, dim := range dims {
		if dim.Name != nil || dim.Length != nil || dim.Comment != "" {
			type Alias ArrayDimensions
			return json.Marshal(Alias(dims))
		}
	}

	return json.Marshal(len(dims))
}

func (tc *TypeDefinitions) MarshalJSON() ([]byte, error) {
	type expanded struct {
		Enum   *EnumDefinition   `json:"enum,omitempty"`
		Flags  *EnumDefinition   `json:"flags,omitempty"`
		Record *RecordDefinition `json:"record,omitempty"`
		Alias  *NamedType        `json:"alias,omitempty"`
	}

	expandedTypes := make([]expanded, len(*tc))
	for i, typeDefinition := range *tc {
		e := expanded{}
		switch t := typeDefinition.(type) {
		case *EnumDefinition:
			if t.IsFlags {
				e.Flags = t
			} else {
				e.Enum = t
			}
		case *NamedType:
			e.Alias = t
		case *RecordDefinition:
			e.Record = t
		default:
			panic(fmt.Sprintf("unexpected type %T", t))
		}

		expandedTypes[i] = e
	}

	return json.Marshal(expandedTypes)
}

func (e *IntegerLiteralExpression) MarshalJSON() ([]byte, error) {
	return []byte(e.Value.String()), nil
}

func (e *StringLiteralExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%q", e.Value))
}

func (e *MemberAccessExpression) MarshalJSON() ([]byte, error) {
	type Alias MemberAccessExpression
	return json.Marshal(struct {
		MemberAccess *Alias `json:"memberAccess"`
	}{
		MemberAccess: (*Alias)(e),
	})
}

func (e *IndexExpression) MarshalJSON() ([]byte, error) {
	type Alias IndexExpression
	return json.Marshal(struct {
		Index *Alias `json:"index"`
	}{
		Index: (*Alias)(e),
	})
}

func (e *FunctionCallExpression) MarshalJSON() ([]byte, error) {
	type Alias FunctionCallExpression
	return json.Marshal(struct {
		Call Alias `json:"call"`
	}{
		Call: Alias(*e),
	})
}

func (e *SwitchExpression) MarshalJSON() ([]byte, error) {
	type Alias SwitchExpression
	return json.Marshal(struct {
		Switch *Alias `json:"switch"`
	}{
		Switch: (*Alias)(e),
	})
}

func (e *TypeConversionExpression) MarshalJSON() ([]byte, error) {
	type Alias TypeConversionExpression
	return json.Marshal(struct {
		Convert *Alias `json:"convert"`
	}{
		Convert: (*Alias)(e),
	})
}
