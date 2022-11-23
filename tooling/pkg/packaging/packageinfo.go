// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package packaging

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/microsoft/yardl/tooling/internal/validation"
	"gopkg.in/yaml.v3"
)

const PackageFileName = "_package.yml"

var namespaceNameRegex = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)

type PackageInfo struct {
	FilePath  string `yaml:"-"`
	Namespace string `yaml:"namespace"`

	Json *JsonCodegenOptions `yaml:"json,omitempty"`
	Cpp  *CppCodegenOptions  `yaml:"cpp,omitempty"`
}

type CppCodegenOptions struct {
	PackageInfo        *PackageInfo `yaml:"-"`
	SourcesOutputDir   string       `yaml:"sourcesOutputDir"`
	GenerateCMakeLists bool         `yaml:"generateCMakeLists"`

	InternalSymlinkStaticHeaders bool `yaml:"internalSymlinkStaticHeaders"`
	InternalGenerateMocks        bool `yaml:"internalGenerateMocks"`
}

func (o CppCodegenOptions) ChangeOutputDir(newRelativeDir string) CppCodegenOptions {
	o.SourcesOutputDir = path.Join(o.SourcesOutputDir, newRelativeDir)
	return o
}

func (o *CppCodegenOptions) UnmarshalYAML(value *yaml.Node) error {
	// Set default values
	o.GenerateCMakeLists = true

	type alias CppCodegenOptions
	return value.DecodeWithOptions((*alias)(o), yaml.DecodeOptions{KnownFields: true})
}

type JsonCodegenOptions struct {
	PackageInfo *PackageInfo `yaml:"-"`
	OutputDir   string       `yaml:"outputDir"`
}

func ReadPackageInfo(directory string) (PackageInfo, error) {
	packageFilePath, _ := filepath.Abs(path.Join(directory, PackageFileName))
	packageInfo := PackageInfo{FilePath: packageFilePath}
	f, err := os.Open(packageFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return packageInfo, fmt.Errorf("a '%s' file is missing from the directory '%s'", PackageFileName, directory)
		}
		return packageInfo, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	err = decoder.Decode(&packageInfo)
	if err != nil {
		return packageInfo, validation.NewValidationError(err, packageFilePath)
	}

	return packageInfo, packageInfo.Validate()
}

func (p *PackageInfo) Validate() error {
	errorSink := &validation.ErrorSink{}
	packageDir := path.Dir(p.FilePath)

	if p.Namespace == "" {
		errorSink.Add(validation.NewValidationError(errors.New("the 'namespace' field is missing"), p.FilePath))
	} else if !namespaceNameRegex.MatchString(p.Namespace) {
		errorSink.Add(validation.NewValidationError(fmt.Errorf("the 'namespace' field must be PascalCased and match the format %s", namespaceNameRegex.String()), p.FilePath))
	}

	if p.Json != nil {
		p.Json.PackageInfo = p
		if p.Json.OutputDir == "" {
			errorSink.Add(validation.NewValidationError(errors.New("the 'json.outputDir' field must not be empty"), p.FilePath))
		} else {
			p.Json.OutputDir = path.Join(packageDir, p.Json.OutputDir)
		}
	}

	if p.Cpp != nil {
		p.Cpp.PackageInfo = p
		if p.Cpp.SourcesOutputDir == "" {
			errorSink.Add(validation.NewValidationError(errors.New("the 'cpp.sourcesOutputDir' field must not be empty"), p.FilePath))
		} else {
			p.Cpp.SourcesOutputDir = path.Join(packageDir, p.Cpp.SourcesOutputDir)
		}
	}

	return errorSink.AsError()
}
