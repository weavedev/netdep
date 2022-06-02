// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"fmt"
	"os"
	"path"
)

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
