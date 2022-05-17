// Package discovery defines discovery of clients calls and endpoints
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package discovery

import (
	"fmt"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
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
	Service string                `json:"service"`
	Calls   map[string][]CallData `json:"calls"`
}

// DiscoveredData is initialised and populated during the discovery stage.
// It stores a list of ServiceCalls for each service and a map of all handled endpoints / methods
// along with the name of the service that handles each one.
type DiscoveredData struct {
	ServCalls []ServiceCalls
	Handled   map[string]string
}

// FindCallersForEndpoint is a sample method for locating
// the callers of a specific endpoint, which is specified
// by the name of its parent service, its path in the target
// project, and its URI.
func FindCallersForEndpoint(parentService, endpointPath, endpointURI string) []interface{} {
	// This is a placeholder; the signature of this method might need to be changed.
	// Return empty slice for now.
	return nil
}

// Discover finds client calls in the specified project directory
func Discover(projDir, svcDir string) ([]*callanalyzer.Caller, []*callanalyzer.Caller, error) {
	conf := callanalyzer.SSAConfig{
		Mode:    ssa.BuilderMode(0),
		SvcDir:  svcDir,
		ProjDir: projDir,
	}

	_, pkg, err := callanalyzer.CreateSSA(conf)
	if err != nil {
		return nil, nil, err
	}

	mains := ssautil.MainPackages(pkg)

	allTargetsClient := make([]*callanalyzer.Caller, 0)
	allTargetsServer := make([]*callanalyzer.Caller, 0)
	for _, mainPkg := range mains {
		targetsOfCurrPkgClient, targetOfCurrPkgServer, err := callanalyzer.AnalyzePackageCalls(mainPkg)
		if err != nil {
			fmt.Println(err.Error())
		}
		allTargetsClient = append(allTargetsClient, targetsOfCurrPkgClient...)
		allTargetsServer = append(allTargetsServer, targetOfCurrPkgServer...)
	}

	return allTargetsClient, allTargetsServer, nil
}
