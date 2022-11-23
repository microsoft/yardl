// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cmd

import (
	"fmt"
	"os"

	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
	"github.com/spf13/cobra"
)

func newValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "validate",
		Short:                 "Validate the package in the current directory.",
		Long:                  `Validate the package in the current directory.`,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := validateImpl()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
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

	namespace, err := dsl.ParseYamlInDir(dir, packageInfo.Namespace)
	if err != nil {
		return err
	}

	_, err = dsl.Validate([]*dsl.Namespace{namespace})
	if err != nil {
		return err
	}
	return err
}
