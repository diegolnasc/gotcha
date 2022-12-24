// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.

package main

import (
	root "github.com/diegolnasc/gotcha/cmd"
)

func main() {
	version := root.GotchaVersion{Version: "1.0.0"}
	root.RootCmd.AddCommand(version.Init())
	root.RootCmd.AddCommand(root.Init())
	root.Execute()
}
