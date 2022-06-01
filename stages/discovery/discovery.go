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

// Discover finds client calls in the specified project directory
func Discover(pkgsToAnalyse []*ssa.Package, config *callanalyzer.AnalyserConfig, envVariables map[string]map[string]string) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, error) {
	// The current output data structure. TODO: add additional fields
	allClientTargets := make([]*callanalyzer.CallTarget, 0)
	allServerTargets := make([]*callanalyzer.CallTarget, 0)

	for _, pkg := range pkgsToAnalyse {
		var (
			clientTargetsOfCurrPkg []*callanalyzer.CallTarget
			serverTargetsOfCurrPkg []*callanalyzer.CallTarget
			err                    error
		)

		if config == nil {
			defaultConf := callanalyzer.DefaultConfigForFindingHTTPCalls(envVariables)
			// Analyse each package with the default config
			clientTargetsOfCurrPkg, serverTargetsOfCurrPkg, err = callanalyzer.AnalysePackageCalls(pkg, &defaultConf)
		} else {
			// Analyse each package
			clientTargetsOfCurrPkg, serverTargetsOfCurrPkg, err = callanalyzer.AnalysePackageCalls(pkg, config)
		}
		if err != nil {
			return nil, nil, err
		}

		allClientTargets = append(allClientTargets, clientTargetsOfCurrPkg...)
		allServerTargets = append(allServerTargets, serverTargetsOfCurrPkg...)
	}
	return allClientTargets, allServerTargets, nil
}
