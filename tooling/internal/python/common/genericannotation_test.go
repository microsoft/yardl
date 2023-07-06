package common

import (
	"os"
	"path"
	"testing"

	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTypeParameterUseDirect(t *testing.T) {
	src := `
Rec<A, B, C, D, E, F, G, H>: !record
  fields:
    arr1: A[]
    scalar1: B
    scalar2: C
    arr2: C[]
    arr3: R2<D>[]
    arr4: E*
    arr5: !array
      items: [F, G]
    arr6: Image<H>

R2<T>: !record
  fields:
    t: T

Image<T>: T[]
`
	env, err := parseAndValidate(t, src)
	require.NoError(t, err)

	AnnotateGenerics(env)

	rec := env.SymbolTable["test.Rec"].(*dsl.RecordDefinition)

	assert.Equal(t, TypeParameterUseArray, rec.TypeParameters[0].Annotations[TypeParameterUseAnnotationKey])
	assert.Equal(t, TypeParameterUseScalar, rec.TypeParameters[1].Annotations[TypeParameterUseAnnotationKey])
	assert.Equal(t, TypeParameterUseBoth, rec.TypeParameters[2].Annotations[TypeParameterUseAnnotationKey])
	assert.Equal(t, TypeParameterUseNone, rec.TypeParameters[3].Annotations[TypeParameterUseAnnotationKey])
	assert.Equal(t, TypeParameterUseScalar, rec.TypeParameters[4].Annotations[TypeParameterUseAnnotationKey])
	assert.Equal(t, TypeParameterUseNone, rec.TypeParameters[5].Annotations[TypeParameterUseAnnotationKey])
	assert.Equal(t, TypeParameterUseNone, rec.TypeParameters[6].Annotations[TypeParameterUseAnnotationKey])
	assert.Equal(t, TypeParameterUseArray, rec.TypeParameters[7].Annotations[TypeParameterUseAnnotationKey])
}

func parseAndValidate(t *testing.T, src string) (*dsl.Environment, error) {
	d := t.TempDir()
	os.WriteFile(path.Join(d, "t.yaml"), []byte(src), 0644)
	ns, err := dsl.ParseYamlInDir(d, "test")
	if err != nil {
		return nil, err
	}

	return dsl.Validate([]*dsl.Namespace{ns})
}
