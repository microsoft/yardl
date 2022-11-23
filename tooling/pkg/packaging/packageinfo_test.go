// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package packaging

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMissingPackageFile(t *testing.T) {
	d := t.TempDir()
	_, err := ReadPackageInfo(d)
	require.ErrorContains(t, err, "a '_package.yml' file is missing from the directory")
}

func TestBasicPackageFile(t *testing.T) {
	packageFileContents := `
namespace: Foo
cpp:
  sourcesOutputDir: cpp
`
	packageInfo, err := writeAndReadPackageFile(t, packageFileContents)
	require.Nil(t, err)
	require.Equal(t, "Foo", packageInfo.Namespace)
	require.Equal(t, path.Join(path.Dir(packageInfo.FilePath), "cpp"), packageInfo.Cpp.SourcesOutputDir)
}

func TestPackageFileWithUnknownField(t *testing.T) {
	packageFileContents := `
namespace: Foo
whatisthis: field?
`
	_, err := writeAndReadPackageFile(t, packageFileContents)
	// This error message comes from the yaml decoder and is a bit crude
	require.ErrorContains(t, err, "field whatisthis not found in type")
}

func TestPackageFileWithUnknownFieldOnCpp(t *testing.T) {
	packageFileContents := `
namespace: Foo
cpp:
  sourcesOutputDir: cpp
  whatisthis: field?
`
	_, err := writeAndReadPackageFile(t, packageFileContents)
	// This error message comes from the yaml decoder and is a bit crude
	require.ErrorContains(t, err, "field whatisthis not found in type")
}

func TestPackageFileWithMissingRequiredField(t *testing.T) {
	packageFileContents := `
cpp:
  sourcesOutputDir: cpp
`
	_, err := writeAndReadPackageFile(t, packageFileContents)
	require.ErrorContains(t, err, "the 'namespace' field is missing")
}

func TestPackageFileWithInvalidNamespace(t *testing.T) {
	packageFileContents := `
namespace: 123
`
	_, err := writeAndReadPackageFile(t, packageFileContents)
	require.ErrorContains(t, err, "the 'namespace' field must be PascalCased and match the format")
}

func writeAndReadPackageFile(t *testing.T, packageFileContents string) (PackageInfo, error) {
	d := t.TempDir()
	os.WriteFile(path.Join(d, PackageFileName), []byte(packageFileContents), 0644)
	packageInfo, err := ReadPackageInfo(d)
	return packageInfo, err
}
