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

	packageInfo, err := packaging.LoadPackage(dir)
	if err != nil {
		return err
	}

	_, err = validatePackage(packageInfo)

	return err
}

func validatePackage(packageInfo packaging.PackageInfo) (*dsl.Environment, error) {
	namespace, err := dsl.ParsePackageContents(packageInfo)
	if err != nil {
		return nil, err
	}

	env, err := dsl.Validate([]*dsl.Namespace{namespace})
	if err != nil {
		return nil, err
	}

	for _, packageInfo := range packageInfo.PreviousVersions {
		namespace, err := dsl.ParsePackageContents(packageInfo)
		if err != nil {
			return env, err
		}

		_, err = dsl.Validate([]*dsl.Namespace{namespace})
		if err != nil {
			return env, err
		}
	}

	return env, nil
}
