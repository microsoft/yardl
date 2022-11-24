// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package main

import (
	"github.com/microsoft/yardl/tooling/internal/cmd"
)

var (
	version = ""
	commit  = ""
)

func main() {
	cmd.Execute(version, commit)
}
