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

	nonErroredPackages, count := filterOutErroredPackages(loadedPackages)

	if count < 1 {
		return nil, fmt.Errorf("no usable packages found")
	}

	program, processedPackages := ssautil.AllPackages(nonErroredPackages, builderMode)

	// Why we can build program but only return processedPackages:
	// *ssa+Program has a reference to processedPackages. *ssa+Program.Build() populates processedPackages, too.
	program.Build()

	return processedPackages, nil
}

// filterOutErroredPackages removes errored packages from the list of analyzable packages
func filterOutErroredPackages(loadedPackages []*packages.Package) ([]*packages.Package, int) {
	nonErroredPackages := make([]*packages.Package, 0)
	validPackageCount := 0
	for _, loadedPackage := range loadedPackages {
		if len(loadedPackage.Errors) == 0 {
			nonErroredPackages = append(nonErroredPackages, loadedPackage)
			validPackageCount++
		}
	}
	return nonErroredPackages, validPackageCount
}
