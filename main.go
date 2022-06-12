// Package main
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package main

import (
	"os"

	"lab.weave.nl/internships/tud-2022/netDep/cmd"
)

// main is the entry point to the program
func main() {
	// execute the main logic
	err := runRoot()
	if err != nil {
		// report an unsuccessful run
		os.Exit(1)
		return
	}
}

// runRoot registers any subcommands and runs the root command of netDep
func runRoot() error {
	rootCmd := cmd.RootCmd()
	// add the subcommand for generating a manpage
	rootCmd.AddCommand(cmd.GenManpageCmd(rootCmd))
	return rootCmd.Execute()
}
