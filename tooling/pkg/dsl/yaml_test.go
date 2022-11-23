// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseVectorOfVectors(t *testing.T) {
	src := `
x: !record
  fields:
    vectorOfVectors: !vector
      items: !vector
        items: int
        length: 2`
	ns, err := parse(t, src)
	require.Nil(t, err)

	f := ns.TypeDefinitions[0].(*RecordDefinition).Fields[0]

	require.Equal(t, "int", RequireType[*SimpleType](t, RequireType[*GeneralizedType](t, RequireType[*GeneralizedType](t, f.Type).Cases[0].Type).Cases[0].Type).Name)

}

func TestBasicErrors(t *testing.T) {
	testCases := []struct {
		src string
		err string
	}{
		{`a:`, "type cannot be empty"},
		{`a: !bogus {}`, "unrecognized type kind '!bogus'"},
		{`a: !array {items: !bogus s}`, "unrecognized type kind '!bogus'"},
		{`a: !array {items: null}`, "type null is only supported in unions"},
		{`a: !array {items: }`, "type null is only supported in unions"},
	}

	for _, tC := range testCases {
		t.Run(tC.src, func(t *testing.T) {
			_, err := parse(t, tC.src)
			require.ErrorContains(t, err, tC.err)
		})
	}
}

func TestCommentsOnRecords(t *testing.T) {
	src := `
# This is a comment on a record
x: !record
  fields:
    # comment on field
    f: int`

	ns, err := parse(t, src)
	require.Nil(t, err)
	require.Equal(t, "This is a comment on a record", ns.TypeDefinitions[0].(*RecordDefinition).Comment)
	require.Equal(t, "comment on field", ns.TypeDefinitions[0].(*RecordDefinition).Fields[0].Comment)
}

func TestCommentsOnEnums(t *testing.T) {
	src := `
# This is a comment on an enum
x: !enum
  values:
    # comment on value
    - A`

	ns, err := parse(t, src)
	require.Nil(t, err)
	require.Equal(t, "This is a comment on an enum", ns.TypeDefinitions[0].(*EnumDefinition).Comment)
	require.Equal(t, "comment on value", ns.TypeDefinitions[0].(*EnumDefinition).Values[0].Comment)
}

func TestCommentsOnNamedTypes(t *testing.T) {
	src := `
# This is a comment on a named type
x: string`

	ns, err := parse(t, src)
	require.Nil(t, err)
	require.Equal(t, "This is a comment on a named type", ns.TypeDefinitions[0].(*NamedType).Comment)
}

func TestCommentsOnProtocols(t *testing.T) {
	src := `
# This is a comment on a protocol
x: !protocol
  sequence:
    # comment on step
    i: int`

	ns, err := parse(t, src)
	require.Nil(t, err)
	require.Equal(t, "This is a comment on a protocol", ns.Protocols[0].Comment)
	require.Equal(t, "comment on step", ns.Protocols[0].Sequence[0].Comment)
}

func RequireType[T any](t *testing.T, value any) T {
	if typed, ok := value.(T); ok {
		return typed
	}
	var x T
	t.Fatalf("value of type %T is not of expected type %T", value, x)
	panic("")
}

func parse(t *testing.T, src string) (*Namespace, error) {
	d := t.TempDir()
	os.WriteFile(path.Join(d, "t.yaml"), []byte(src), 0644)
	return ParseYamlInDir(d, "test")
}
