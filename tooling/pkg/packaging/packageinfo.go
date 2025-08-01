// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package packaging

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/validation"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const PackageFileName = "_package.yml"

const MaxImportRecursionDepth = 10

var namespaceNameRegex = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)
var versionLabelRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*(_[a-zA-Z0-9]+)*$`)

type PackageInfo struct {
	FilePath  string `yaml:"-"`
	Namespace string `yaml:"namespace"`

	Versions Versions `yaml:"versions,omitempty"`
	Imports  Imports  `yaml:"imports,omitempty"`

	Json   *JsonCodegenOptions   `yaml:"json,omitempty"`
	Cpp    *CppCodegenOptions    `yaml:"cpp,omitempty"`
	Python *PythonCodegenOptions `yaml:"python,omitempty"`
	Matlab *MatlabCodegenOptions `yaml:"matlab,omitempty"`
}

func (p *PackageInfo) PackageDir() string {
	return filepath.Dir(p.FilePath)
}

func (p *PackageInfo) GetAllReferencedPackages() []*PackageInfo {
	checked := make(map[string]bool)
	var refs []*PackageInfo

	var recurse func(*PackageInfo)
	recurse = func(pInfo *PackageInfo) {
		for _, imp := range pInfo.Imports {
			if !checked[imp.Package.FilePath] {
				recurse(imp.Package)
				checked[imp.Package.FilePath] = true
				refs = append(refs, imp.Package)
			}
		}
	}

	// Get all imported packages
	recurse(p)

	// Get all other versions, and their respective imported packages
	for _, ver := range p.Versions {
		if !checked[ver.Package.FilePath] {
			recurse(ver.Package)
			checked[ver.Package.FilePath] = true
			refs = append(refs, ver.Package)
		}
	}
	return refs
}

func (p *PackageInfo) validate() error {
	errorSink := &validation.ErrorSink{}

	if p.Namespace == "" {
		errorSink.Add(validation.NewValidationError(errors.New("the 'namespace' field is missing"), p.FilePath))
	} else if !namespaceNameRegex.MatchString(p.Namespace) {
		errorSink.Add(validation.NewValidationError(fmt.Errorf("the 'namespace' field must be PascalCased and match the format %s", namespaceNameRegex.String()), p.FilePath))
	}

	for _, ver := range p.Versions {
		if ver.Label == "" {
			errorSink.Add(validation.NewValidationError(errors.New("the version label is missing"), p.FilePath))
		} else if !versionLabelRegex.MatchString(ver.Label) {
			errorSink.Add(validation.NewValidationError(fmt.Errorf("the version label '%s' must match the format %s", ver.Label, versionLabelRegex.String()), p.FilePath))
		}
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

	if p.Matlab != nil {
		p.Matlab.PackageInfo = p
		if p.Matlab.OutputDir == "" {
			errorSink.Add(validation.NewValidationError(errors.New("the 'matlab.outputDir' field must not be empty"), p.FilePath))
		} else {
			p.Matlab.OutputDir = filepath.Join(p.PackageDir(), p.Matlab.OutputDir)
		}
	}

	return errorSink.AsError()
}

type Import struct {
	Url     string
	Package *PackageInfo
}
type Imports []*Import

func (imports *Imports) UnmarshalYAML(value *yaml.Node) error {
	unpacked := []*Import(*imports)

	if value.Tag != "!!seq" {
		return fmt.Errorf("expected import sequence")
	}

	for _, item := range value.Content {
		if item.Tag != "!!str" {
			return fmt.Errorf("expected import url to be a string")
		}

		unpacked = append(unpacked, &Import{Url: item.Value})
	}

	*imports = Imports(unpacked)
	return nil
}

type Version struct {
	Label   string
	Url     string
	Package *PackageInfo
}

type Versions []*Version

func (versions *Versions) UnmarshalYAML(value *yaml.Node) error {
	unpacked := []*Version(*versions)

	if value.Tag != "!!map" {
		return fmt.Errorf("expected versions map")
	}

	for i := 0; i < len(value.Content); i += 2 {
		verKey := value.Content[i]
		verValue := value.Content[i+1]
		if verKey.Tag != "!!str" {
			return fmt.Errorf("expected version label to be a string")
		}
		if verValue.Tag != "!!str" {
			return fmt.Errorf("expected version url to be a string")
		}

		unpacked = append(unpacked, &Version{Label: verKey.Value, Url: verValue.Value})
	}

	*versions = Versions(unpacked)
	return nil
}

type CppCodegenOptions struct {
	PackageInfo         *PackageInfo `yaml:"-"`
	Disabled            bool         `yaml:"disabled"`
	SourcesOutputDir    string       `yaml:"sourcesOutputDir"`
	GenerateHDF5        bool         `yaml:"generateHDF5"`
	GenerateNDJson      bool         `yaml:"generateNDJson"`
	GenerateCMakeLists  bool         `yaml:"generateCMakeLists"`
	OverrideArrayHeader string       `yaml:"overrideArrayHeader"`

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
	o.GenerateHDF5 = true
	o.GenerateNDJson = true
	o.GenerateCMakeLists = true

	type alias CppCodegenOptions
	return value.DecodeWithOptions((*alias)(o), yaml.DecodeOptions{KnownFields: true})
}

type JsonCodegenOptions struct {
	PackageInfo *PackageInfo `yaml:"-"`
	Disabled    bool         `yaml:"disabled"`
	OutputDir   string       `yaml:"outputDir"`
}

type PythonCodegenOptions struct {
	PackageInfo                *PackageInfo `yaml:"-"`
	Disabled                   bool         `yaml:"disabled"`
	OutputDir                  string       `yaml:"outputDir"`
	GenerateNDJson             bool         `yaml:"generateNDJson"`
	InternalSymlinkStaticFiles bool         `yaml:"internalSymlinkStaticFiles"`
}

func (o *PythonCodegenOptions) UnmarshalYAML(value *yaml.Node) error {
	// Set default values
	o.GenerateNDJson = true

	type alias PythonCodegenOptions
	return value.DecodeWithOptions((*alias)(o), yaml.DecodeOptions{KnownFields: true})
}

type MatlabCodegenOptions struct {
	PackageInfo                *PackageInfo `yaml:"-"`
	Disabled                   bool         `yaml:"disabled"`
	OutputDir                  string       `yaml:"outputDir"`
	InternalSymlinkStaticFiles bool         `yaml:"internalSymlinkStaticFiles"`
	InternalGenerateMocks      bool         `yaml:"internalGenerateMocks"`
}

// Parses PackageInfo in dir then loads all package Imports and Predecessors
func LoadPackage(dir string) (*PackageInfo, error) {
	packageInfo, err := loadPackageVersion(dir)
	if err != nil {
		return packageInfo, err
	}

	vdirs, err := collectVersions(packageInfo)
	if err != nil {
		return packageInfo, err
	}

	for i, vdir := range vdirs {
		if vdir == dir {
			// The main package has a version label
			packageInfo.Versions[i].Package = packageInfo
			continue
		}

		versionInfo, err := loadPackageVersion(vdir)
		if err != nil {
			return packageInfo, err
		}

		packageInfo.Versions[i].Package = versionInfo
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
	log.Debug().Msgf("%s- %s from %s", strings.Repeat("  ", indent), p.Namespace, p.PackageDir())
	for _, imp := range p.Imports {
		logImports(imp.Package, indent+1)
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

	log.Info().Msgf("Parsed packageInfo with namespace: %v", packageInfo.Namespace)
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

	log.Info().Msgf("Collecting imports for %v", parentInfo.PackageDir())
	var importUrls []string
	for _, imp := range parentInfo.Imports {
		importUrls = append(importUrls, imp.Url)
	}
	dirs, err := fetchAndCachePackages(parentInfo.PackageDir(), importUrls)
	if err != nil {
		return parentInfo, validation.NewValidationError(err, parentInfo.FilePath)
	}

	for i, dir := range dirs {
		importChain[parentInfo.Namespace] = true
		childInfo, err := collectPackages(dir, alreadyCollected, importChain, depthRemaining-1)
		if err != nil {
			return parentInfo, err
		}
		importChain[parentInfo.Namespace] = false

		// Build the Import tree
		parentInfo.Imports[i].Package = childInfo
	}

	return parentInfo, nil
}

// Fetch and cache each package version in pkgInfo.Versions
func collectVersions(pkgInfo *PackageInfo) ([]string, error) {
	if len(pkgInfo.Versions) <= 0 {
		return nil, nil
	}

	log.Info().Msgf("Collecting referenced versions for %v", pkgInfo.PackageDir())
	var versionUrls []string
	for _, ver := range pkgInfo.Versions {
		versionUrls = append(versionUrls, ver.Url)
	}
	dirs, err := fetchAndCachePackages(pkgInfo.PackageDir(), versionUrls)
	if err != nil {
		err = validation.NewValidationError(err, pkgInfo.FilePath)
	}
	return dirs, err
}
