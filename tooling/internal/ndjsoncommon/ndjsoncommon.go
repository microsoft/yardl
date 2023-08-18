package ndjsoncommon

import (
	"fmt"

	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

type JsonDataType int

const (
	JsonNull JsonDataType = 1 << iota
	JsonBoolean
	JsonNumber
	JsonString
	JsonArray
	JsonObject
)

func GetJsonDataType(t dsl.Type) JsonDataType {
	if t == nil {
		return JsonNull
	}
	gt := dsl.ToGeneralizedType(t)
	switch d := gt.Dimensionality.(type) {
	case *dsl.Vector:
		return JsonArray
	case *dsl.Array:
		if d.IsFixed() {
			return JsonArray
		}
		return JsonObject
	case *dsl.Map:
		if p, ok := dsl.GetPrimitiveType(d.KeyType); ok && p == dsl.String {
			return JsonObject
		}
		return JsonArray
	}

	if len(gt.Cases) > 1 {
		panic("unexpected union type")
	}

	scalarType := gt.Cases[0].Type.(*dsl.SimpleType)
	switch td := scalarType.ResolvedDefinition.(type) {
	case dsl.PrimitiveDefinition:
		switch td {
		case dsl.String:
			return JsonString
		case dsl.Int8, dsl.Int16, dsl.Int32, dsl.Int64, dsl.Uint8, dsl.Uint16, dsl.Uint32, dsl.Uint64, dsl.Size, dsl.Float32, dsl.Float64:
			return JsonNumber
		case dsl.Bool:
			return JsonBoolean
		case dsl.ComplexFloat32, dsl.ComplexFloat64:
			return JsonArray
		case dsl.Date, dsl.Time, dsl.DateTime:
			return JsonNumber
		default:
			panic(fmt.Sprintf("unexpected primitive type %s", td))
		}
	case *dsl.EnumDefinition:
		if td.IsFlags {
			return JsonArray
		}
		return JsonString | JsonNumber
	case *dsl.RecordDefinition:
		return JsonObject
	case *dsl.GenericTypeParameter:
		return JsonObject
	case *dsl.NamedType:
		return GetJsonDataType(td.Type)
	default:
		panic(fmt.Sprintf("unexpected type %T", td))
	}
}
