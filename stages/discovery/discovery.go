// Package discovery defines discovery of clients calls and endpoints
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package discovery

import (
	"golang.org/x/tools/go/ssa"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"
)

/*
In the Discovery stages, clients and endpoints are discovered and mapped to their parent service.
Refer to the Project plan, chapter 5.3 for more information.
*/

// CallData stores for each call the full path of the file in which it happens and the exact line in that file
type CallData struct {
	Filepath string `json:"filepath"`
	Line     int    `json:"line"`
}

// ServiceCalls stores for each service its name and the calls that it makes (strings of URLs / method names)
type ServiceCalls struct {
	Service string                `json:"service"` // TODO: is it possible to just specify lowerCamelCase to the JSON marshaller?
	Calls   map[string][]CallData `json:"calls"`   // TODO: see above ^
}

// DiscoveredData is initialised and populated during the discovery stage.
// It stores a list of ServiceCalls for each service and a map of all handled endpoints / methods
// along with the name of the service that handles each one.
type DiscoveredData struct {
	ServCalls []ServiceCalls
	Handled   map[string]string
}

// Discover finds client calls in the specified project directory
func Discover(pkgsToAnalyse []*ssa.Package) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, error) {
	// The current output data structure. TODO: add additional fields
	allClientTargets := make([]*callanalyzer.CallTarget, 0)
	allServerTargets := make([]*callanalyzer.CallTarget, 0)
	config := callanalyzer.DefaultConfigForFindingHTTPCalls()
	for _, pkg := range pkgsToAnalyse {
		// Analyse each package
		clientTargetsOfCurrPkg, serverTargetsOfCurrPkg, err := callanalyzer.AnalysePackageCalls(pkg, &config)
		if err != nil {
			return nil, nil, err
		}

		allClientTargets = append(allClientTargets, clientTargetsOfCurrPkg...)
		allServerTargets = append(allServerTargets, serverTargetsOfCurrPkg...)
	}
	return allClientTargets, allServerTargets, nil
}
