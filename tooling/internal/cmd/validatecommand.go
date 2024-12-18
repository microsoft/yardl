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
			configOverrides, err := cmd.Flags().GetStringToString("config")
			if err != nil {
				log.Fatal().Msgf("error getting config: %v", err)
			}

			warnings, err := validateImpl(configOverrides)
			if err != nil {
				log.Error().Msg(err.Error())
				os.Exit(1)
			}
			for _, warning := range warnings {
				log.Warn().Msg(warning)
			}
		},
	}

	return cmd
}

func validateImpl(configArgs map[string]string) ([]string, error) {
	inputDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	packageInfo, err := packaging.LoadPackage(inputDir)
	if err != nil {
		return nil, err
	}

	if err := updatePackageInfoFromArgs(packageInfo, configArgs); err != nil {
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

	var versionEnvs []*dsl.Environment
	var labels []string
	for _, version := range packageInfo.Versions {
		for _, label := range labels {
			if label == version.Label {
				return env, nil, fmt.Errorf("duplicate predecessor label %s", version.Label)
			}
		}
		labels = append(labels, version.Label)

		namespaces, err := parseAndFlattenNamespaces(version.Package)
		if err != nil {
			return nil, nil, err
		}

		oldEnv, err := dsl.Validate(namespaces)
		if err != nil {
			return nil, nil, err
		}

		versionEnvs = append(versionEnvs, oldEnv)
	}

	var warnings []string
	if len(versionEnvs) > 0 {
		env, warnings, err = dsl.ValidateEvolution(env, versionEnvs, labels)
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
