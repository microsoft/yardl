// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package common

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

// TODO: Populate all Matlab reserved names
var isReservedName = map[string]bool{}

var TypeSyntaxWriter dsl.TypeSyntaxWriter[string] = func(self dsl.TypeSyntaxWriter[string], t dsl.Node, contextNamespace string) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		switch t {
		case dsl.Bool:
			return "bool"
		case dsl.Int8:
			return "int8"
		case dsl.Uint8:
			return "uint8"
		case dsl.Int16:
			return "int16"
		case dsl.Uint16:
			return "uint16"
		case dsl.Int32:
			return "int32"
		case dsl.Uint32:
			return "uint32"
		case dsl.Int64:
			return "int64"
		case dsl.Uint64:
			return "uint64"
		case dsl.Size:
			return "uint64"
		case dsl.Float32:
			return "float32"
		case dsl.Float64:
			return "float64"
		case dsl.ComplexFloat32:
			return "complex"
		case dsl.ComplexFloat64:
			return "complex"
		case dsl.String:
			return "str"
		case dsl.Date:
			return "datetime"
		case dsl.Time:
			return "datetime"
		case dsl.DateTime:
			return "datetime"
		default:
			panic(fmt.Sprintf("primitive '%v' not recognized", t))
		}
	case *dsl.GenericTypeParameter:
		return TypeIdentifierName(t.Name)
	case dsl.TypeDefinition:
		meta := t.GetDefinitionMeta()
		typeName := TypeIdentifierName(meta.Name)
		if t.GetDefinitionMeta().Namespace != contextNamespace {
			typeName = fmt.Sprintf("%s.%s", formatting.ToSnakeCase(meta.Namespace), typeName)
		}

		typeSyntax := typeName

		return typeSyntax

	case nil:
		return "None"
	case *dsl.SimpleType:
		return self.ToSyntax(t.ResolvedDefinition, contextNamespace)
	case *dsl.GeneralizedType:
		scalarString := func() string {
			if t.Cases.IsSingle() {
				return self.ToSyntax(t.Cases[0].Type, contextNamespace)
			}
			if t.Cases.IsOptional() {
				return "yardl.Optional"
			}

			return UnionSyntax(t)
		}()

		switch d := t.Dimensionality.(type) {
		case nil, *dsl.Stream, *dsl.Vector, *dsl.Array:
			return scalarString
		case *dsl.Map:
			// TODO: Update to Matlab `dictionary` class when we require Matlab version >= r2023b
			return "containers.Map"
		default:
			panic(fmt.Sprintf("unexpected type %T", d))
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}

}

func TypeSyntax(typeOrTypeDefinition dsl.Node, contextNamespace string) string {
	return TypeSyntaxWriter.ToSyntax(typeOrTypeDefinition, contextNamespace)
}

func ComputedFieldIdentifierName(name string) string {
	cased := formatting.ToSnakeCase(name)
	if !isReservedName[name] {
		return cased
	}

	return cased + "_"
}

func TypeIdentifierName(name string) string {
	if !isReservedName[name] {
		return name
	}

	return name + "_"
}

func PackageDir(name string) string {
	return fmt.Sprintf("+%s", formatting.ToSnakeCase(name))
}

func NamespaceIdentifierName(namespace string) string {
	return formatting.ToSnakeCase(namespace)
}

func FieldIdentifierName(name string) string {
	snakeCased := formatting.ToSnakeCase(name)
	if !isReservedName[snakeCased] {
		return snakeCased
	}

	return snakeCased + "_"
}

func EnumValueIdentifierName(name string) string {
	cased := formatting.ToUpperSnakeCase(name)
	if !isReservedName[cased] {
		return cased
	}

	return cased + "_"
}

func UnionClassName(gt *dsl.GeneralizedType) (className string, typeParameters string) {
	if !gt.Cases.IsUnion() {
		panic("Not a union")
	}

	cases := make([]string, 0, len(gt.Cases))
	for _, typeCase := range gt.Cases {
		if typeCase.Type == nil {
			continue
		}
		cases = append(cases, formatting.ToPascalCase(typeCase.Tag))
	}

	return strings.Join(cases, "Or"), ""
}

func UnionSyntax(gt *dsl.GeneralizedType) string {
	className, typeParameters := UnionClassName(gt)
	var syntax string
	if len(typeParameters) > 0 {
		// TODO: Remove
		syntax = fmt.Sprintf("%s[%s]", className, typeParameters)
	} else {
		syntax = className
	}

	// TODO: Not needed in Matlab
	// if gt.Cases.HasNullOption() {
	// 	return fmt.Sprintf("typing.Optional[%s]", syntax)
	// }

	return syntax
}

func WriteBlockBody(w *formatting.IndentedWriter, f func()) {
	defer func() {
		w.WriteStringln("end")
	}()
	w.Indented(f)
}

func WriteComment(w *formatting.IndentedWriter, comment string) {
	comment = strings.TrimSpace(comment)
	if comment != "" {
		w = formatting.NewIndentedWriter(w, "% ").Indent()
		w.WriteStringln(comment)
	}
}

func AbstractWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sWriterBase", formatting.ToPascalCase(p.Name))
}

func AbstractReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sReaderBase", formatting.ToPascalCase(p.Name))
}

func ProtocolWriteMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("write_%s", formatting.ToSnakeCase(s.Name))
}

func ProtocolWriteImplMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("write_%s_", formatting.ToSnakeCase(s.Name))
}

func ProtocolReadMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("read_%s", formatting.ToSnakeCase(s.Name))
}

func ProtocolReadImplMethodName(s *dsl.ProtocolStep) string {
	return fmt.Sprintf("read_%s_", formatting.ToSnakeCase(s.Name))
}

func WriteGeneratedFileHeader(w *formatting.IndentedWriter) {
	WriteComment(w, "This file was generated by the \"yardl\" tool. DO NOT EDIT.")
	w.WriteStringln("")
}

type MatlabFileWriter struct {
	PackageDir string
}

func (fw *MatlabFileWriter) WriteFile(name string, writeContents func(w *formatting.IndentedWriter)) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	WriteGeneratedFileHeader(w)

	writeContents(w)

	fname := fmt.Sprintf("%s.m", name)
	definitionsPath := path.Join(fw.PackageDir, fname)
	if err := iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}
