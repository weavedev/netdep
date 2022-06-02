// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"fmt"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// LoadAndBuildPackages takes in project root directory path and the path
// of one service and returns the SSA representation of the service.
func LoadAndBuildPackages(projectRootDir string, svcPath string) ([]*ssa.Package, error) {
	// setup build buildConfig
	buildConfig := &packages.Config{
		Dir: projectRootDir,
		//nolint // We are using this, as cmd/callgraph is using it.
		Mode:  packages.LoadAllSyntax,
		Tests: false,
	}

	builderMode := ssa.BuilderMode(0)

	// load all packages in the service directory
	loadedPackages, err := packages.Load(buildConfig, svcPath)
	if err != nil {
		return nil, err
	}

	if len(loadedPackages) == 0 {
		return nil, fmt.Errorf("no packages")
	}

	if packages.PrintErrors(loadedPackages) > 0 {
		return nil, fmt.Errorf("packages contain errors")
	}

	prog, pkgs := ssautil.AllPackages(loadedPackages, builderMode)
	// prog has a reference to pkgs internally,
	// and prog.Build() populates pkgs with necessary
	// information
	prog.Build()

	return pkgs, nil
}
