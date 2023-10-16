// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package packaging

import (
	"errors"
	"fmt"
	"os"
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

	Predecessors     []string      `yaml:"predecessors,omitempty"`
	Imports          []string      `yaml:"imports,omitempty"`
	PreviousVersions []PackageInfo `yaml:"-"`

	Json   *JsonCodegenOptions   `yaml:"json,omitempty"`
	Cpp    *CppCodegenOptions    `yaml:"cpp,omitempty"`
	Python *PythonCodegenOptions `yaml:"python,omitempty"`
}

type CppCodegenOptions struct {
	PackageInfo        *PackageInfo `yaml:"-"`
	SourcesOutputDir   string       `yaml:"sourcesOutputDir"`
	GenerateCMakeLists bool         `yaml:"generateCMakeLists"`

	InternalSymlinkStaticHeaders bool `yaml:"internalSymlinkStaticHeaders"`
	InternalGenerateMocks        bool `yaml:"internalGenerateMocks"`
	InternalGenerateTranslator   bool `yaml:"internalGenerateTranslator"`
}

func (o CppCodegenOptions) ChangeOutputDir(newRelativeDir string) CppCodegenOptions {
	o.SourcesOutputDir = filepath.Join(o.SourcesOutputDir, newRelativeDir)
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

type PythonCodegenOptions struct {
	PackageInfo                *PackageInfo `yaml:"-"`
	OutputDir                  string       `yaml:"outputDir"`
	InternalSymlinkStaticFiles bool         `yaml:"internalSymlinkStaticFiles"`
}

func ReadPackageInfo(directory string) (PackageInfo, error) {
	packageFilePath, _ := filepath.Abs(filepath.Join(directory, PackageFileName))
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

	//log.Printf("Parsed packageInfo with namespace: %v", packageInfo.Namespace)

	return packageInfo, packageInfo.Validate()
}

// Parses PackageInfo in dir then loads all package Imports and Predecessors
func LoadPackage(dir string) (PackageInfo, error) {
	packageInfo, err := ReadPackageInfo(dir)
	if err != nil {
		return packageInfo, err
	}

	err = CollectImports(packageInfo)
	if err != nil {
		return packageInfo, err
	}

	dirs, err := CollectPredecessors(packageInfo)
	if err != nil {
		return packageInfo, err
	}

	for _, dir := range dirs {
		predecessorInfo, err := ReadPackageInfo(dir)
		if err != nil {
			return packageInfo, err
		}

		err = CollectImports(predecessorInfo)
		if err != nil {
			return packageInfo, err
		}

		packageInfo.PreviousVersions = append(packageInfo.PreviousVersions, predecessorInfo)
	}

	return packageInfo, nil
}

func (p *PackageInfo) PackageDir() string {
	return filepath.Dir(p.FilePath)
}

func (p *PackageInfo) Validate() error {
	errorSink := &validation.ErrorSink{}

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
			p.Json.OutputDir = filepath.Join(p.PackageDir(), p.Json.OutputDir)
		}
	}

	if p.Cpp != nil {
		p.Cpp.PackageInfo = p
		if p.Cpp.SourcesOutputDir == "" {
			errorSink.Add(validation.NewValidationError(errors.New("the 'cpp.sourcesOutputDir' field must not be empty"), p.FilePath))
		} else {
			p.Cpp.SourcesOutputDir = filepath.Join(p.PackageDir(), p.Cpp.SourcesOutputDir)
		}
	}

	if p.Python != nil {
		p.Python.PackageInfo = p
		if p.Python.OutputDir == "" {
			errorSink.Add(validation.NewValidationError(errors.New("the 'python.outputDir' field must not be empty"), p.FilePath))
		} else {
			p.Python.OutputDir = filepath.Join(p.PackageDir(), p.Python.OutputDir)
		}
	}

	return errorSink.AsError()
}
