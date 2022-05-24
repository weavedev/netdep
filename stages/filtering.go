// Package stages defines different stages of analysis
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package stages

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type Annotation struct {
	ServiceName string
	Position    token.Position
	Value       string
}

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

// LoadServices takes a project directory and a service
// directory and for each directory of that service builds
// an SSA representation and a list of annotations for each service in svcDir.
func LoadServices(projectDir string, svcDir string) ([]*ssa.Package, []*Annotation, error) {
	// Collect all files within the services directory
	files, err := os.ReadDir(svcDir)
	if err != nil {
		return nil, nil, err
	}

	packagesToAnalyze := make([]*ssa.Package, 0)
	annotations := make([]*Annotation, 0)

	for _, file := range files {
		if file.IsDir() {
			servicePath := path.Join(svcDir, file.Name())
			fmt.Println(servicePath)

			pkgs, err := LoadPackages(projectDir, servicePath)
			if err != nil {
				return nil, nil, err
			}
			serviceAnnotations, err := LoadAnnotations(servicePath, file.Name())
			if err != nil {
				return nil, nil, err
			}
			annotations = append(annotations, serviceAnnotations...)

			packagesToAnalyze = append(packagesToAnalyze, pkgs...)
		}
	}

	if len(packagesToAnalyze) == 0 {
		return nil, nil, fmt.Errorf("no service to analyse were found")
	}

	return packagesToAnalyze, annotations, nil
}

// LoadAnnotations scans all the files of a given service directory and returns a list of
// Annotation from the comments in the format "//netdep: ..." that it discovers.
func LoadAnnotations(servicePath string, serviceName string) ([]*Annotation, error) {
	files, err := os.ReadDir(servicePath)
	if err != nil {
		return nil, err
	}

	serviceAnnotations := make([]*Annotation, 0)

	for _, file := range files {
		if file.Name()[len(file.Name())-3:] == ".go" {
			// If the file is a .go file - parse it
			fileAnnotations, err := parseComments(path.Join(servicePath, file.Name()), serviceName)
			if err != nil {
				return nil, err
			}
			serviceAnnotations = append(serviceAnnotations, fileAnnotations...)
		} else if file.IsDir() {
			// If the file is a directory - recursively look for .go files inside it
			innerServiceAnnotations, err := LoadAnnotations(path.Join(servicePath, file.Name()), serviceName)
			if err != nil {
				return nil, err
			}
			serviceAnnotations = append(serviceAnnotations, innerServiceAnnotations...)
		}
	}

	return serviceAnnotations, nil
}

// parseComments parses the given file with a parser.ParseComments mode, filters out
// the comments which don't contain a substring "netdep", generates an Annotation for
// every remaining comment and returns a list of them.
func parseComments(path string, serviceName string) ([]*Annotation, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	annotations := make([]*Annotation, 0)

	for _, commentGroup := range f.Comments {
		for _, comment := range commentGroup.List {
			if strings.Contains(comment.Text, "netdep") {
				ann := &Annotation{
					ServiceName: serviceName,
					Position:    fs.Position(comment.Slash),
					Value:       comment.Text,
				}
				annotations = append(annotations, ann)
			}
		}
	}

	return annotations, nil
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
