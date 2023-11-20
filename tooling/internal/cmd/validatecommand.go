// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cmd

import (
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
			err := validateImpl()
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}
		},
	}

	return cmd
}

func validateImpl() error {
	inputDir, err := os.Getwd()
	if err != nil {
		return err
	}

	packageInfo, err := packaging.LoadPackage(inputDir)
	if err != nil {
		return err
	}

	_, err = validatePackage(packageInfo)

	return err
}

func validatePackage(packageInfo *packaging.PackageInfo) (*dsl.Environment, error) {
	namespaces, err := parseAndFlattenNamespaces(packageInfo)
	if err != nil {
		return nil, err
	}

	env, err := dsl.Validate(namespaces)
	if err != nil {
		return nil, err
	}

	for _, packageInfo := range packageInfo.Predecessors {
		namespaces, err := parseAndFlattenNamespaces(packageInfo)
		if err != nil {
			return nil, err
		}

		_, err = dsl.Validate(namespaces)
		if err != nil {
			return nil, err
		}
	}

	return env, nil
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
		log.Debug().Msgf("Already parsed namespace %s (%p)", existing.Name, existing)
		return existing, nil
	}

	namespace, err := dsl.ParsePackageContents(p)
	if err != nil {
		return nil, err
	}

	alreadyParsed[p.Namespace] = namespace
	log.Debug().Msgf("Parsed namespace %s (%p)", namespace.Name, namespace)

	for _, dep := range p.Imports {
		ns, err := parsePackageNamespaces(dep, alreadyParsed)
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
