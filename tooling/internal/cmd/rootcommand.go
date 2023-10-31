// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func newRootCommand(version, commit string) *cobra.Command {
	if version == "" {
		version = "unknown"
	}
	if commit != "" {
		version = fmt.Sprintf("%s commit %s", version, commit)
	}

	verbose := false
	log.SetOutput(io.Discard)

	cmd := &cobra.Command{
		Use: "yardl",
		Long: `yardl generates domain types and serialization code from a simple schema language.

Read more at https://github.com/microsoft/yardl`,
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				log.SetOutput(os.Stderr)
			}
		},
	}

	// hide --help as a flag in the usage output
	cmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	cmd.PersistentFlags().Lookup("help").Hidden = true

	cobra.EnableCommandSorting = false

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "verbose output")
	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newGenerateCommand())
	cmd.AddCommand(newValidateCommand())

	return cmd
}

func Execute(version, commit string) {
	err := newRootCommand(version, commit).Execute()
	if err != nil {
		os.Exit(1)
	}
}
