// Package discovery defines discovery of clients calls and endpoints
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package discovery

import (
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

// Discover finds client calls in the specified project directory
func Discover(projDir, svcDir string) ([]*callanalyzer.Target, error) {
	conf := callanalyzer.SSAConfig{
		Mode:    ssa.BuilderMode(0),
		SvcDir:  svcDir,
		ProjDir: projDir,
	}

	_, pkg, err := callanalyzer.CreateSSA(conf)
	if err != nil {
		return nil, err
	}

	mains := ssautil.MainPackages(pkg)

	allTargets := make([]*callanalyzer.Target, 0)
	for _, mainPkg := range mains {
		targetsOfCurrPkg, _ := callanalyzer.AnalyzePackageCalls(mainPkg)
		allTargets = append(allTargets, targetsOfCurrPkg...)
	}

	return allTargets, nil
}
