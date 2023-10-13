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
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	packageInfo, err := packaging.ReadPackageInfo(dir)
	if err != nil {
		return err
	}

	err = packaging.CollectImports(dir, packageInfo.Imports)
	if err != nil {
		return err
	}

	namespace, err := dsl.ParseYamlInDir(dir, packageInfo.Namespace)
	if err != nil {
		return err
	}

	_, err = dsl.Validate([]*dsl.Namespace{namespace})
	if err != nil {
		return err
	}

	// Now, load all previous versions

	dirs, err := packaging.CollectPredecessors(dir, packageInfo.Predecessors)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		packageInfo, err := packaging.ReadPackageInfo(dir)
		if err != nil {
			return err
		}

		err = packaging.CollectImports(dir, packageInfo.Imports)
		if err != nil {
			return err
		}

		namespace, err := dsl.ParseYamlInDir(dir, packageInfo.Namespace)
		if err != nil {
			return err
		}

		_, err = dsl.Validate([]*dsl.Namespace{namespace})
		if err != nil {
			return err
		}
	}

	return nil
}
