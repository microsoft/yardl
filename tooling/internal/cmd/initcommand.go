// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
	"github.com/spf13/cobra"
)

//go:embed initcontent/package.tpl
var packageFileTemplate string

//go:embed initcontent/model.yml
var modelFileContents string

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "init PACKAGE_NAME",
		Short:                 "Create a package in the current directory.",
		Long:                  `Create a package in the current directory.`,
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := initImpl(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}

func initImpl(namespace string) error {
	template := template.Must(template.New("package").Parse(packageFileTemplate))
	packageFile, err := os.OpenFile(packaging.PackageFileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("%s already exists", packaging.PackageFileName)
		}
		return err
	}
	defer packageFile.Close()

	data := struct{ Namespace string }{
		Namespace: formatting.ToPascalCase(namespace),
	}

	err = template.Execute(packageFile, data)
	if err != nil {
		return err
	}

	modelFileName := "model.yml"

	modelFile, err := os.OpenFile(modelFileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("%s already exists", modelFileName)
		}
		return err
	}
	defer packageFile.Close()

	_, err = modelFile.WriteString(modelFileContents)
	return err
}
