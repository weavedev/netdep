// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"
)

type IntCall struct {
	Name      string
	NumParams int
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
			servicePath := filepath.Join(servicesDir, file.Name())
			packagesToAnalyze = append(packagesToAnalyze, servicePath)
		}
	}

	if len(packagesToAnalyze) == 0 {
		return nil, fmt.Errorf("no service to analyse were found")
	}

	return packagesToAnalyze, err
}

// FindServiceCalls iterates through all the go files in the servicecalls package (files ending in -service.go !!) and
// scans all the method names defined in the interfaces.
func FindServiceCalls(serviceCallsDir string) (map[IntCall]string, *[]*callanalyzer.CallTarget, error) {
	serviceCalls := make(map[IntCall]string)
	serverTargets := make([]*callanalyzer.CallTarget, 0)

	if serviceCallsDir == "" {
		return serviceCalls, &serverTargets, nil
	}

	// Collect all files within the packages directory
	files, err := os.ReadDir(serviceCallsDir)
	if err != nil {
		return serviceCalls, &serverTargets, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".go" && strings.HasSuffix(file.Name(), "-service.go") {
			serviceName := file.Name()[:len(file.Name())-11]
			parseInterfaces(filepath.Join(serviceCallsDir, file.Name()), serviceName, serviceCalls, &serverTargets)
		}
	}

	return serviceCalls, &serverTargets, nil
}
