// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cmd

import (
	"log"
	"os"

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
				log.Fatalf("Error: %v\n", err)
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

func parseNamespaces(p *packaging.PackageInfo, namespaces *[]*dsl.Namespace) error {
	namespace, err := dsl.ParsePackageContents(p)
	if err != nil {
		return nil
	}

	for _, dep := range p.Imports {
		if err := parseNamespaces(dep, namespaces); err != nil {
			return err
		}
	}

	*namespaces = append(*namespaces, namespace)
	log.Printf("Parsed namespace %v", namespace.Name)

	return nil
}

func validatePackage(packageInfo *packaging.PackageInfo) (*dsl.Environment, error) {
	var namespaces []*dsl.Namespace

	if err := parseNamespaces(packageInfo, &namespaces); err != nil {
		return nil, err
	}

	env, err := dsl.Validate(namespaces)
	if err != nil {
		return nil, err
	}

	for _, packageInfo := range packageInfo.Predecessors {
		var namespaces []*dsl.Namespace
		if err := parseNamespaces(packageInfo, &namespaces); err != nil {
			return nil, err
		}

		_, err = dsl.Validate(namespaces)
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}
