// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path"
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
		Short:                 "Generate scaffolding for a new package",
		Long:                  `Creates a new package directory named 'model' under the current directory with an example model.`,
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
	modelDir := "model"
	if err := os.MkdirAll(modelDir, 0775); err != nil {
		return err
	}

	packageFilePath := path.Join(modelDir, packaging.PackageFileName)

	template := template.Must(template.New("package").Parse(packageFileTemplate))
	packageFile, err := os.OpenFile(packageFilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("%s already exists", packageFilePath)
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

	modelFilePath := path.Join(modelDir, "model.yml")

	modelFile, err := os.OpenFile(modelFilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("%s already exists", modelFilePath)
		}
		return err
	}
	defer packageFile.Close()

	_, err = modelFile.WriteString(modelFileContents)
	if err != nil {
		return err
	}

	fmt.Println("Initialized new package in the 'model' directory.")
	fmt.Println("To generate code for it, run the following commands:")
	fmt.Println("  cd model")
	fmt.Println("  yardl generate")

	return nil
}
