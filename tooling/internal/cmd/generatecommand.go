// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/inancgumus/screen"
	"github.com/microsoft/yardl/tooling/internal/cpp"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
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
		Run: func(*cobra.Command, []string) {
			if !flags.watch {
				packageInfo, err := generateCore()
				if err != nil {
					// avoiding returning the error here because
					// cobra prefixes the error with "Error: "
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				WriteSuccessfulSummary(packageInfo)

				return
			}

			// Enter watch mode
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.Fatal(err)
			}
			defer watcher.Close()

			completedChannel := make(chan error)
			go dedupLoop(watcher, completedChannel)

			err = watcher.Add(".")
			if err != nil {
				log.Fatal(err)
			}

			err = <-completedChannel
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().BoolVarP(&flags.watch, "watch", "w", false, "Regenerate code whenever a file in the current directory changes.")

	return cmd
}

// dedup fsnotify events
func dedupLoop(w *fsnotify.Watcher, completedChannel chan<- error) {
	generateInWatchMode()

	const waitFor = 5 * time.Millisecond
	timer := time.AfterFunc(math.MaxInt64, generateInWatchMode)
	timer.Stop()

	for {
		select {
		case err, ok := <-w.Errors:
			if !ok {
				// channel was closed
				completedChannel <- err
				return
			}

			log.Printf("ERROR: %s\n", err)
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

func generateInWatchMode() {
	defer func() {
		if err := recover(); err != nil {
			screen.Clear()
			screen.MoveTopLeft()
			fmt.Printf("panic: %v \n%s", err, string(debug.Stack()))
		}
	}()

	packageInfo, err := generateCore()
	screen.Clear()
	screen.MoveTopLeft()

	fmt.Printf("Validated model package '%s' at %v.\n\n", packageInfo.Namespace, time.Now().Format("15:04:05"))

	if err != nil {
		fmt.Println(err)
	} else {
		WriteSuccessfulSummary(packageInfo)
	}
}

func WriteSuccessfulSummary(packageInfo packaging.PackageInfo) {
	if packageInfo.Cpp != nil {
		fmt.Printf("✅ Wrote C++ to %s.\n", packageInfo.Cpp.SourcesOutputDir)
	}
	if packageInfo.Python != nil {
		fmt.Printf("✅ Wrote Python to %s.\n", packageInfo.Python.OutputDir)
	}
	if packageInfo.Json != nil {
		fmt.Printf("✅ Wrote JSON to %s.\n", packageInfo.Json.OutputDir)
	}
}

func generateCore() (packaging.PackageInfo, error) {
	inputDir, _ := os.Getwd()
	packageInfo, err := packaging.ReadPackageInfo(inputDir)
	if err != nil {
		return packageInfo, err
	}

	namespace, err := dsl.ParseYamlInDir(inputDir, packageInfo.Namespace)
	if err != nil {
		return packageInfo, err
	}

	env, err := dsl.Validate([]*dsl.Namespace{namespace})
	if err != nil {
		return packageInfo, err
	}

	if packageInfo.Cpp != nil {
		err = cpp.Generate(env, *packageInfo.Cpp)
		if err != nil {
			return packageInfo, err
		}
	}

	if packageInfo.Python != nil {
		err = python.Generate(env, *packageInfo.Python)
		if err != nil {
			return packageInfo, err
		}
	}

	if packageInfo.Json != nil {
		err = outputJson(env, packageInfo.Json)
		if err != nil {
			return packageInfo, err
		}
	}

	return packageInfo, err
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
