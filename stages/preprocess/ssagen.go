// Package preprocess defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocess

import (
	"fmt"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// LoadPackages takes in project root directory path and the path
// of one service and returns an ssa representation of the service.
func LoadPackages(projectRootDir string, svcPath string) ([]*ssa.Package, error) {
	config := &packages.Config{
		Dir: projectRootDir,
		//nolint // We are using this, as cmd/callgraph is using it.
		Mode:  packages.LoadAllSyntax,
		Tests: false,
	}
	mode := ssa.BuilderMode(0)

	initial, err := packages.Load(config, svcPath)
	if err != nil {
		return nil, err
	}

	if len(initial) == 0 {
		return nil, fmt.Errorf("no packages")
	}

	if packages.PrintErrors(initial) > 0 {
		return nil, fmt.Errorf("packages contain errors")
	}

	prog, pkgs := ssautil.AllPackages(initial, mode)
	// prog has a reference to pkgs internally,
	// and prog.Build() populates pkgs with necessary
	// information
	prog.Build()
	return pkgs, nil
}
