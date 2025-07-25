// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/fsnotify/fsnotify"
	"github.com/inancgumus/screen"
	"github.com/microsoft/yardl/tooling/internal/cpp"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/internal/matlab"
	"github.com/microsoft/yardl/tooling/internal/python"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
	"github.com/spf13/cobra"
)

func newGenerateCommand() *cobra.Command {
	var flags struct {
		watch bool
	}

	cmd := &cobra.Command{
		Use:                   "generate [--watch]",
		Aliases:               []string{"gen"},
		Short:                 "generate code for the package in the current directory",
		Long:                  `generate code for the package in the current directory`,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			configOverrides, err := cmd.Flags().GetStringToString("config")
			if err != nil {
				log.Fatal().Msgf("error getting config: %v", err)
			}

			if !flags.watch {
				packageInfo, warnings, err := generateImpl(configOverrides)
				if err != nil {
					// avoiding returning the error here because
					// cobra prefixes the error with "Error: "
					log.Error().Msg(err.Error())
					os.Exit(1)
				}

				for _, warning := range warnings {
					log.Warn().Msg(warning)
				}
				WriteSuccessfulSummary(packageInfo)

				return
			}

			// Enter watch mode
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.Fatal().Err(err).Msg("")
			}
			defer watcher.Close()

			completedChannel := make(chan error)
			go dedupLoop(configOverrides, watcher, completedChannel)

			err = watcher.Add(".")
			if err != nil {
				log.Fatal().Err(err).Msg("")
			}

			err = <-completedChannel
			if err != nil {
				log.Fatal().Msgf("%s", err)
			}
		},
	}

	cmd.Flags().BoolVarP(&flags.watch, "watch", "w", false, "Regenerate code whenever a file in the current directory changes.")

	return cmd
}

// dedup fsnotify events
func dedupLoop(configArgs map[string]string, w *fsnotify.Watcher, completedChannel chan<- error) {
	regenerate := func() {
		dirsToWatch := generateInWatchMode(configArgs)
		if dirsToWatch != nil && len(dirsToWatch) > len(w.WatchList()) {
			for _, dir := range dirsToWatch {
				if err := w.Add(dir); err != nil {
					completedChannel <- err
					return
				}
			}
		}

	}

	regenerate()

	const waitFor = 5 * time.Millisecond
	timer := time.AfterFunc(math.MaxInt64, regenerate)
	timer.Stop()

	for {
		select {
		case err, ok := <-w.Errors:
			if !ok {
				// channel was closed
				completedChannel <- err
				return
			}

			log.Error().Err(err).Msg("")
		case _, ok := <-w.Events:
			if !ok {
				// channel was closed
				completedChannel <- nil
				return
			}

			timer.Reset(waitFor)
		}
	}
}

// Returns the directories to watch after parsing all package imports, or nil on error
func generateInWatchMode(configArgs map[string]string) []string {
	defer func() {
		if err := recover(); err != nil {
			screen.Clear()
			screen.MoveTopLeft()
			fmt.Printf("panic: %v \n%s", err, string(debug.Stack()))
		}
	}()

	packageInfo, warnings, err := generateImpl(configArgs)
	screen.Clear()
	screen.MoveTopLeft()

	if err != nil {
		log.Error().Msg(err.Error())
	} else {
		fmt.Printf("Validated model package '%s' at %v.\n\n", packageInfo.Namespace, time.Now().Format("15:04:05"))
		for _, warning := range warnings {
			log.Warn().Msg(warning)
		}
		WriteSuccessfulSummary(packageInfo)

		var dirsToWatch []string
		for _, ref := range packageInfo.GetAllReferencedPackages() {
			dirsToWatch = append(dirsToWatch, ref.PackageDir())
		}
		return dirsToWatch
	}
	return nil
}

func WriteSuccessfulSummary(packageInfo *packaging.PackageInfo) {
	if packageInfo.Cpp != nil {
		fmt.Printf("✅ Wrote C++ to %s.\n", packageInfo.Cpp.SourcesOutputDir)
	}
	if packageInfo.Python != nil {
		fmt.Printf("✅ Wrote Python to %s.\n", packageInfo.Python.OutputDir)
	}
	if packageInfo.Json != nil {
		fmt.Printf("✅ Wrote JSON to %s.\n", packageInfo.Json.OutputDir)
	}
	if packageInfo.Matlab != nil {
		fmt.Printf("✅ Wrote Matlab to %s.\n", packageInfo.Matlab.OutputDir)
	}
}

func generateImpl(configArgs map[string]string) (*packaging.PackageInfo, []string, error) {
	inputDir, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	packageInfo, err := packaging.LoadPackage(inputDir)
	if err != nil {
		return packageInfo, nil, err
	}

	if err := updatePackageInfoFromArgs(packageInfo, configArgs); err != nil {
		return packageInfo, nil, err
	}

	env, warnings, err := validatePackage(packageInfo)
	if err != nil {
		return packageInfo, warnings, err
	}

	if packageInfo.Cpp != nil && !packageInfo.Cpp.Disabled {
		err = cpp.Generate(env, *packageInfo.Cpp)
		if err != nil {
			return packageInfo, warnings, err
		}
	}

	if packageInfo.Python != nil && !packageInfo.Python.Disabled {
		err = python.Generate(env, *packageInfo.Python)
		if err != nil {
			return packageInfo, warnings, err
		}
	}

	if packageInfo.Json != nil && !packageInfo.Json.Disabled {
		err = outputJson(env, packageInfo.Json)
		if err != nil {
			return packageInfo, warnings, err
		}
	}

	if packageInfo.Matlab != nil && !packageInfo.Matlab.Disabled {
		err = matlab.Generate(env, *packageInfo.Matlab)
		if err != nil {
			return packageInfo, warnings, err
		}
	}

	return packageInfo, warnings, err
}

func outputJson(env *dsl.Environment, options *packaging.JsonCodegenOptions) error {
	if err := os.MkdirAll(options.OutputDir, 0775); err != nil {
		return err
	}

	b, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return err
	}

	jsonPath := path.Join(options.OutputDir, "model.json")
	return iocommon.WriteFileIfNeeded(jsonPath, b, 0644)
}
