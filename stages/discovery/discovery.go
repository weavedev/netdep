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

// DiscoverAll creates a combined list of all discovered calls in the given packages.
func DiscoverAll(packages []*ssa.Package, config *callanalyzer.AnalyserConfig) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, error) {
	allClientTargets := make([]*callanalyzer.CallTarget, 0)
	allServerTargets := make([]*callanalyzer.CallTarget, 0)

	for _, pkg := range packages {
		clientCalls, serverCalls, err := Discover(pkg, config)
		if err != nil {
			return nil, nil, err
		}

		allClientTargets = append(allClientTargets, clientCalls...)
		allServerTargets = append(allServerTargets, serverCalls...)
	}

	return allClientTargets, allServerTargets, nil
}

// Discover finds client and server calls in the given packages
func Discover(pkg *ssa.Package, config *callanalyzer.AnalyserConfig) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, error) {
	// The current output data structure. TODO: add additional fields

	if config == nil {
		defaultConf := callanalyzer.DefaultConfigForFindingHTTPCalls()
		// Analyse each package with the default config
		return callanalyzer.AnalysePackageCalls(pkg, &defaultConf)
	} else {
		// Analyse each package
		return callanalyzer.AnalysePackageCalls(pkg, config)
	}
}
