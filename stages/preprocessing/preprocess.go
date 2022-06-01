// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"fmt"
	"os"
	"path"

	"golang.org/x/tools/go/ssa"
)

// LoadServices takes a project directory and a service
// directory and for each directory of that service builds
// an SSA representation and a list of Annotation for each service in svcDir.
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
