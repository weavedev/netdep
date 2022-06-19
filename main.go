// Package main
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package main

import (
	"os"
	"runtime/debug"

	"lab.weave.nl/internships/tud-2022/netDep/cmd"
)

// main is the entry point to the program
func main() {
	debug.SetMaxStack(1000000000)
	// execute the main logic
	runRoot()
}

// runRoot registers any subcommands and runs the root command of netDep
func runRoot() {
	rootCmd := cmd.RootCmd()
	// add the subcommand for generating a manpage
	rootCmd.AddCommand(cmd.GenManpageCmd(rootCmd))
	err := rootCmd.Execute()
	if err != nil {
		// report an unsuccessful run
		os.Exit(1)
	}
}
