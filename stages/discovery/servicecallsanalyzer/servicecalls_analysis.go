/*
Package servicecallsanalyzer defines servicecalls package specific scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package servicecallsanalyzer

import (
	"os"
	"path/filepath"
	"strings"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"
)

// IntCall stores the name and the number of parameters of a method call,
// currently used to store interesting calls which come from the servicecalls package.
type IntCall struct {
	Name      string
	NumParams int
}

// ParseServiceCallsPackage iterates through all the go files in the servicecalls package (files ending in -service.go !!) and
// scans all the method names defined in the interfaces.
func ParseServiceCallsPackage(serviceCallsDir string) (map[IntCall]string, *[]*callanalyzer.CallTarget, error) {
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
			ParseInterfaces(filepath.Join(serviceCallsDir, file.Name()), serviceName, serviceCalls, &serverTargets)
		}
	}

	return serviceCalls, &serverTargets, nil
}

// LoadServiceCalls scans all the files of a given service directory and returns a list of
// clientTargets based on the method names found in the servicecalls package.
func LoadServiceCalls(servicePath string, serviceName string, internalCalls map[IntCall]string, clientTargets *[]*callanalyzer.CallTarget) error {
	files, err := os.ReadDir(servicePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".go" && !strings.HasSuffix(file.Name(), "_test.go") && !strings.HasSuffix(file.Name(), "pb.go") {
			// If the file is a .go file - parse it
			currClientTargets, err := ParseMethods(filepath.Join(servicePath, file.Name()), internalCalls, serviceName)
			if err != nil {
				return err
			}
			*clientTargets = append(*clientTargets, *currClientTargets...)
		} else if file.IsDir() {
			// If the file is a directory - recursively look for .go files inside it
			err := LoadServiceCalls(filepath.Join(servicePath, file.Name()), serviceName, internalCalls, clientTargets)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
