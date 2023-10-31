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

const MaxImportRecursionDepth = 10

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

func (p *PackageInfo) PackageDir() string {
	return filepath.Dir(p.FilePath)
}

func (p *PackageInfo) GetAllImportedPackages() []*PackageInfo {
	checked := make(map[string]bool)
	var imports []*PackageInfo
	var recurse func(*PackageInfo)
	recurse = func(pInfo *PackageInfo) {
		for _, ref := range pInfo.Imports {
			if !checked[ref.FilePath] {
				recurse(ref)
				checked[ref.FilePath] = true
				imports = append(imports, ref)
			}
		}
	}
	recurse(p)
	return imports
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

// Parses PackageInfo in dir then loads all package Imports and Predecessors
func LoadPackage(dir string) (*PackageInfo, error) {
	packageInfo, err := loadPackageVersion(dir)
	if err != nil {
		return nil, err
	}

	dirs, err := collectPredecessors(packageInfo)
	if err != nil {
		return packageInfo, err
	}

	for _, dir := range dirs {
		predecessorInfo, err := loadPackageVersion(dir)
		if err != nil {
			return packageInfo, err
		}

		packageInfo.Predecessors = append(packageInfo.Predecessors, predecessorInfo)
	}

	return packageInfo, nil
}

func loadPackageVersion(dir string) (*PackageInfo, error) {
	pkgsCollected := make(map[string]*PackageInfo)
	importChain := make(map[string]bool)
	packageInfo, err := collectPackages(dir, pkgsCollected, importChain, MaxImportRecursionDepth)
	if err != nil {
		return packageInfo, err
	}

	logImports(packageInfo, 0)

	return packageInfo, nil
}

func logImports(p *PackageInfo, indent int) {
	log.Printf("%s- %s from %s (%p)", strings.Repeat("  ", indent), p.Namespace, p.PackageDir(), p)
	for _, imp := range p.Imports {
		logImports(imp, indent+1)
	}
}

func readPackageInfo(directory string) (*PackageInfo, error) {
	packageDir, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(packageDir)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("package directory '%s' not found", packageDir)
	}

	packageFilePath := filepath.Join(packageDir, PackageFileName)
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

	log.Printf("Parsed packageInfo with namespace: %v", packageInfo.Namespace)
	return packageInfo, packageInfo.validate()
}

// Recursively collects all packages starting with parentDir, building an Import tree of *PackageInfo
// alreadyCollected is used to check for namespace conflicts (e.g. same namespace but different package directory)
// importChain is used to check for import cycles
// depthRemaining is used to limit the depth of the import tree
func collectPackages(parentDir string, alreadyCollected map[string]*PackageInfo, importChain map[string]bool, depthRemaining int) (*PackageInfo, error) {
	parentInfo, err := readPackageInfo(parentDir)
	if err != nil {
		return nil, err
	}

	if importChain[parentInfo.Namespace] {
		return parentInfo, validation.NewValidationError(fmt.Errorf("import cycle detected"), parentInfo.FilePath)
	}

	if collected, found := alreadyCollected[parentInfo.Namespace]; found {
		if collected.FilePath != parentInfo.FilePath {
			return collected, validation.NewValidationError(fmt.Errorf("namespace '%s' conflicts with '%s'", parentInfo.Namespace, collected.FilePath), parentInfo.FilePath)
		} else {
			return collected, nil
		}
	}

	alreadyCollected[parentInfo.Namespace] = parentInfo

	if depthRemaining <= 0 {
		return parentInfo, validation.NewValidationError(errors.New("reached maximum number of recursive imports"), parentInfo.FilePath)
	}

	log.Printf("Collecting imports for %v", parentInfo.PackageDir())
	dirs, err := fetchAndCachePackages(parentInfo.PackageDir(), parentInfo.ImportUrls)
	if err != nil {
		return parentInfo, validation.NewValidationError(err, parentInfo.FilePath)
	}

	for _, dir := range dirs {
		importChain[parentInfo.Namespace] = true
		childInfo, err := collectPackages(dir, alreadyCollected, importChain, depthRemaining-1)
		if err != nil {
			return parentInfo, err
		}
		importChain[parentInfo.Namespace] = false

		// Build the Import tree
		parentInfo.Imports = append(parentInfo.Imports, childInfo)
	}

	return parentInfo, nil
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
