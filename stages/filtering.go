// Package stages defines different stages of analysis
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package stages

import (
	"fmt"
	"go/ast"
	"os"
	"path"

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

// FindServices takes a directory which contains services
// and returns a list of service directories.
func FindServices(servicesDir string) ([]string, error) {
	// Collect all files within the services directory
	files, err := os.ReadDir(servicesDir)
	if err != nil {
		return nil, err
	}

	packagesToAnalyze := make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			servicePath := path.Join(servicesDir, file.Name())
			packagesToAnalyze = append(packagesToAnalyze, servicePath)
		}
	}

	if len(packagesToAnalyze) == 0 {
		return nil, fmt.Errorf("no service to analyse were found")
	}

	return packagesToAnalyze, err
}

/*
In the Filtering stages, irrelevant files and directories are removed from the target project.
Refer to the Project plan, chapter 5.1 for more information.
*/

// ScanAndFilter returns a map with:
// - Key: service name
// - Value: array of the services' ASTs per file.
func ScanAndFilter(svcDir string) map[string][]*ast.File {
	// TODO: perhaps, for each service, filter its contents?
	servicesList := findAllServices(svcDir)
	for i := 0; i < len(servicesList); i++ {
		_ = filter(servicesList[i], nil)
		// TODO: add to map the resulting AST array
	}
	filter("test", nil)

	return nil
}

// FindAllServices
// is a method for finding all services, which takes the path of the svc directory as an argument
// Returns a list of all services.
//
// TODO: Remove the following line when implementing this method
// goland:noinspection GoUnusedParameter,GoUnusedFunction
func findAllServices(svcDir string) []string {
	// TODO extract a list of the paths of each service
	return nil
}

// Filter
// is currently a placeholder method for filtering the directory of a specified service.
//
// Return a list of ASTs (of each of the files).
//
// TODO: Remove the following line when implementing this method
// goland:noinspection GoUnusedParameter,GoUnusedFunction
func filter(serviceLoc string, filterList []string) []*ast.File {
	// TODO: This is a placeholder; the signature of this method might need to be changed.
	// TODO: Loop over all subdirectories/files of this service, looking for relevant files
	// Return empty slice for now.
	return make([]*ast.File, 0)
}
