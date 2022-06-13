// Package discovery defines discovery of clients calls and endpoints
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package discovery

import (
	"fmt"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"
	"lab.weave.nl/internships/tud-2022/netDep/stages/output"
)

/*
In the Discovery stages, clients and endpoints are discovered and mapped to their parent service.
Refer to the Project plan, chapter 5.3 for more information.
*/

func FindCallPointer(packages []*ssa.Package) map[*ssa.CallCommon]*ssa.Function {
	baseMap := map[*ssa.CallCommon]*ssa.Function{}

	if packages == nil || len(packages) == 0 {
		return baseMap
	}

	var mains []*ssa.Package
	for _, pkg := range packages {
		if pkg == nil || pkg.Pkg == nil {
			if pkg != nil {
				fmt.Println("No package for " + pkg.String())
			}
			continue
		}

		if pkg.Pkg.Name() == "main" && pkg.Func("main") != nil {
			mains = append(mains, pkg)
		}
	}

	ptConfig := &pointer.Config{
		Mains:          mains,
		BuildCallGraph: true,
	}

	fmt.Println("Running pointer analysis...")

	pointerRes, err := pointer.Analyze(ptConfig)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	cg := pointerRes.CallGraph
	cg.DeleteSyntheticNodes()

	for _, node := range cg.Nodes {
		for _, edge := range node.Out {
			if edge.Site == nil {
				continue
			}
			baseMap[edge.Site.Common()] = edge.Callee.Func
		}
	}

	fmt.Printf("Found mapping for %d calls", len(baseMap))

	return baseMap
}

// DiscoverAll creates a combined list of all discovered calls in the given packages.
func DiscoverAll(packages []*ssa.Package, config *callanalyzer.AnalyserConfig) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, error) {
	allClientTargets := make([]*callanalyzer.CallTarget, 0)
	allServerTargets := make([]*callanalyzer.CallTarget, 0)

	useConfig := config

	if config == nil {
		newConfig := callanalyzer.DefaultConfigForFindingHTTPCalls()
		useConfig = &newConfig
	}

	callPointer := FindCallPointer(packages)

	for _, pkg := range packages {
		if pkg == nil {
			continue
		}

		clientCalls, serverCalls, err := callanalyzer.AnalysePackageCalls(pkg, useConfig, callPointer)
		if err != nil {
			return nil, nil, err
		}

		allClientTargets = append(allClientTargets, clientCalls...)
		allServerTargets = append(allServerTargets, serverCalls...)
	}

	err := callanalyzer.ReplaceTargetsAnnotations(&allClientTargets, config)
	if err != nil {
		return nil, nil, err
	}
	err = callanalyzer.ReplaceTargetsAnnotations(&allServerTargets, config)
	if err != nil {
		return nil, nil, err
	}

	// Filter the targets which are still unresolved and send them to the output stage
	// To print annotation location suggestions for the user
	unresolvedTargets := filterUnresolvedTargets(&allClientTargets, &allServerTargets)
	output.PrintAnnotationSuggestions(unresolvedTargets)

	return allClientTargets, allServerTargets, nil
}

// Discover finds client and server calls in the given packages
func Discover(pkg *ssa.Package, config *callanalyzer.AnalyserConfig) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, error) {
	// The current output data structure. TODO: add additional fields

	if config == nil {
		defaultConf := callanalyzer.DefaultConfigForFindingHTTPCalls()
		// Analyse each package with the default config
		return callanalyzer.AnalysePackageCalls(pkg, &defaultConf, map[*ssa.CallCommon]*ssa.Function{})
	} else {
		// Analyse each package
		return callanalyzer.AnalysePackageCalls(pkg, config, map[*ssa.CallCommon]*ssa.Function{})
	}
}

// filterUnresolvedTargets filters both client and server targets and returns a list of unresolved targets which is later
// passed on to the output stage to print annotation suggestions.
func filterUnresolvedTargets(clientTargets *[]*callanalyzer.CallTarget, serverTargets *[]*callanalyzer.CallTarget) []*callanalyzer.CallTarget {
	unresolvedTargets := make([]*callanalyzer.CallTarget, 0)

	for _, client := range *clientTargets {
		if !client.IsResolved {
			unresolvedTargets = append(unresolvedTargets, client)
		}
	}

	for _, server := range *serverTargets {
		if !server.IsResolved {
			unresolvedTargets = append(unresolvedTargets, server)
		}
	}

	return unresolvedTargets
}
