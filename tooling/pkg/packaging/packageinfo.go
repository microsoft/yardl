// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package packaging

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/validation"
	"gopkg.in/yaml.v3"
)

const PackageFileName = "_package.yml"

var namespaceNameRegex = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)

type PackageInfo struct {
	FilePath  string `yaml:"-"`
	Namespace string `yaml:"namespace"`

	PredecessorUrls []string `yaml:"predecessors,omitempty"`
	ImportUrls      []string `yaml:"imports,omitempty"`

	Predecessors []*PackageInfo `yaml:"-"`
	Imports      []*PackageInfo `yaml:"-"`

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

func ReadPackageInfo(directory string) (*PackageInfo, error) {
	packageFilePath, _ := filepath.Abs(filepath.Join(directory, PackageFileName))
	packageInfo := &PackageInfo{FilePath: packageFilePath}
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
	return packageInfo, packageInfo.validate()
}

// Parses PackageInfo in dir then loads all package Imports and Predecessors
func LoadPackage(dir string) (*PackageInfo, error) {
	packageInfo, err := loadVersion(dir)
	if err != nil {
		return nil, err
	}

	dirs, err := collectPredecessors(packageInfo)
	if err != nil {
		return packageInfo, err
	}

	for _, dir := range dirs {
		predecessorInfo, err := loadVersion(dir)
		if err != nil {
			return packageInfo, err
		}

		packageInfo.Predecessors = append(packageInfo.Predecessors, predecessorInfo)
	}

	return packageInfo, nil
}

// Returns path to package directory
func (p *PackageInfo) PackageDir() string {
	return filepath.Dir(p.FilePath)
}

func (p *PackageInfo) validate() error {
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

func loadVersion(dir string) (*PackageInfo, error) {
	packageInfo, err := ReadPackageInfo(dir)
	if err != nil {
		return packageInfo, err
	}

	err = collectImportsRecursively(packageInfo, MaxImportRecursionDepth)
	if err != nil {
		return packageInfo, err
	}

	log.Printf("Imports collected for %v:", packageInfo.FilePath)
	logImports(packageInfo, 0)

	return packageInfo, nil
}

func logImports(p *PackageInfo, level int) {
	indent := strings.Repeat("  ", level)
	log.Printf("%v- %v: %v", indent, p.PackageDir(), p.Namespace)
	for _, dep := range p.Imports {
		logImports(dep, level+1)
	}
}

func collectImportsRecursively(parent *PackageInfo, depthRemaining int) error {
	if len(parent.ImportUrls) <= 0 {
		return nil
	}

	if depthRemaining <= 0 {
		return validation.NewValidationError(errors.New("reached maximum number of recursive imports"), parent.FilePath)
	}

	log.Printf("Collecting imports for %v", parent.PackageDir())
	dirs, err := fetchAndCachePackages(parent.PackageDir(), parent.ImportUrls)
	if err != nil {
		return validation.NewValidationError(err, parent.FilePath)
	}

	for _, dir := range dirs {
		packageInfo, err := ReadPackageInfo(dir)
		if err != nil {
			return err
		}

		if err := collectImportsRecursively(packageInfo, depthRemaining-1); err != nil {
			return err
		}

		parent.Imports = append(parent.Imports, packageInfo)
	}

	return nil
}

func collectPredecessors(pkgInfo *PackageInfo) ([]string, error) {
	if len(pkgInfo.PredecessorUrls) <= 0 {
		return nil, nil
	}

	log.Printf("Collecting predecessors for %v", pkgInfo.PackageDir())
	dirs, err := fetchAndCachePackages(pkgInfo.PackageDir(), pkgInfo.PredecessorUrls)
	if err != nil {
		err = validation.NewValidationError(err, pkgInfo.FilePath)
	}
	return dirs, err
}
