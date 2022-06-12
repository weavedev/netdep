// Package main
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package main

import (
	"os"

	"lab.weave.nl/internships/tud-2022/netDep/cmd"
)

// main registers any subcommands and runs the root command of netDep
func main() {
	rootCmd := cmd.RootCmd()
	// add the subcommand for generating a manpage
	rootCmd.AddCommand(cmd.GenManpageCmd(rootCmd))
	// execute the root command
	err := rootCmd.Execute()
	if err != nil {
		// report an unsuccessful run
		os.Exit(1)
		return
	}
}
