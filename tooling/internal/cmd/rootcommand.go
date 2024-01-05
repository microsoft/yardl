// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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
	quiet := false
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	cmd := &cobra.Command{
		Use: "yardl",
		Long: `yardl generates domain types and serialization code from a simple schema language.

Read more at https://github.com/microsoft/yardl`,
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}
			if quiet {
				zerolog.SetGlobalLevel(zerolog.ErrorLevel)
			}
		},
	}

	// hide --help as a flag in the usage output
	cmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	cmd.PersistentFlags().Lookup("help").Hidden = true

	cobra.EnableCommandSorting = false

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "show debug output")
	cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "", false, "hide warnings")
	cmd.MarkFlagsMutuallyExclusive("verbose", "quiet")

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
