// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"os"
	"path/filepath"
	"strings"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"
)

// LoadAnnotations scans all the files of a given service directory and returns a list of
// Annotation from the comments in the format "//netdep: ..." that it discovers.
func LoadAnnotations(servicePath string, serviceName string, annotations map[string]map[callanalyzer.Position]string) error {
	files, err := os.ReadDir(servicePath)
	if err != nil {
		return err
	}

	annotations[serviceName] = make(map[callanalyzer.Position]string)

	for _, file := range files {
		//nolint:goconst
		if filepath.Ext(file.Name()) == ".go" && !strings.HasSuffix(file.Name(), "_test.go") && !strings.HasSuffix(file.Name(), "pb.go") {
			// If the file is a .go file - parse it
			err := parseComments(filepath.Join(servicePath, file.Name()), serviceName, annotations)
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			// If the file is a directory - recursively look for .go files inside it
			err := LoadAnnotations(filepath.Join(servicePath, file.Name()), serviceName, annotations)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
