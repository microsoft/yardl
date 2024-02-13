// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
	"github.com/spf13/cobra"
)

func newValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "validate",
		Short:                 "Validate the package in the current directory",
		Long:                  `Validate the package in the current directory`,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			warnings, err := validateImpl()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			for _, warning := range warnings {
				log.Warn().Msg(warning)
			}
		},
	}

	return cmd
}

func validateImpl() ([]string, error) {
	inputDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	packageInfo, err := packaging.LoadPackage(inputDir)
	if err != nil {
		return nil, err
	}

	_, warnings, err := validatePackage(packageInfo)

	return warnings, err
}

func validatePackage(packageInfo *packaging.PackageInfo) (*dsl.Environment, []string, error) {
	namespaces, err := parseAndFlattenNamespaces(packageInfo)
	if err != nil {
		return nil, nil, err
	}

	env, err := dsl.Validate(namespaces)
	if err != nil {
		return nil, nil, err
	}

	var predecessors []*dsl.Environment
	var labels []string
	for _, predecessor := range packageInfo.Versions {
		for _, label := range labels {
			if label == predecessor.Label {
				return env, nil, fmt.Errorf("duplicate predecessor label %s", predecessor.Label)
			}
		}
		labels = append(labels, predecessor.Label)

		namespaces, err := parseAndFlattenNamespaces(predecessor.Package)
		if err != nil {
			return nil, nil, err
		}

		oldEnv, err := dsl.Validate(namespaces)
		if err != nil {
			return nil, nil, err
		}

		predecessors = append(predecessors, oldEnv)
	}

	var warnings []string
	if len(predecessors) > 0 {
		env, warnings, err = dsl.ValidateEvolution(env, predecessors, labels)
		if err != nil {
			return nil, warnings, err
		}
	}

	return env, warnings, nil
}

func parseAndFlattenNamespaces(p *packaging.PackageInfo) ([]*dsl.Namespace, error) {
	alreadyParsed := make(map[string]*dsl.Namespace)
	namespace, err := parsePackageNamespaces(p, alreadyParsed)
	if err != nil {
		return nil, err
	}

	namespace.IsTopLevel = true

	deduplicator := make(map[*dsl.Namespace]bool)
	return flattenNamespaces(namespace, deduplicator), nil
}

func parsePackageNamespaces(p *packaging.PackageInfo, alreadyParsed map[string]*dsl.Namespace) (*dsl.Namespace, error) {
	if existing, found := alreadyParsed[p.Namespace]; found {
		log.Debug().Msgf("Already parsed namespace %s", existing.Name)
		return existing, nil
	}

	namespace, err := dsl.ParsePackageContents(p)
	if err != nil {
		return nil, err
	}

	alreadyParsed[p.Namespace] = namespace
	log.Debug().Msgf("Parsed namespace %s", namespace.Name)

	for _, imp := range p.Imports {
		ns, err := parsePackageNamespaces(imp.Package, alreadyParsed)
		if err != nil {
			return namespace, nil
		}
		namespace.References = append(namespace.References, ns)
	}

	return namespace, nil
}

func flattenNamespaces(ns *dsl.Namespace, duplicate map[*dsl.Namespace]bool) (flat []*dsl.Namespace) {
	if duplicate[ns] {
		return flat
	}
	duplicate[ns] = true

	for _, child := range ns.References {
		flat = append(flat, flattenNamespaces(child, duplicate)...)
	}

	return append(flat, ns)
}
