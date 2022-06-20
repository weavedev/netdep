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

// runPointerAnalysis perform Pointer Analysis on the given packages.
// Returns the pointer analysis result, which includes the call graph
func runPointerAnalysis(packages []*ssa.Package) (*pointer.Result, error) {
	var mains []*ssa.Package
	for _, pkg := range packages {
		if pkg == nil || pkg.Pkg == nil {
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

	return pointer.Analyze(ptConfig)
}

// FindCallPointer builds the points-to sets for calls
// It performs pointer analysis, and reformats the call graph
// The output is mapping from a call to its points-to set of functions
func FindCallPointer(packages []*ssa.Package) map[*ssa.CallCommon][]*ssa.Function {
	baseMap := map[*ssa.CallCommon][]*ssa.Function{}
	baseMapSet := map[*ssa.CallCommon]map[*ssa.Function]bool{}

	if len(packages) == 0 {
		return baseMap
	}

	pointerRes, err := runPointerAnalysis(packages)
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

			// get call
			call := edge.Site.Common()
			_, has := baseMap[call]

			// first time this call is found, create a new set
			if !has {
				baseMap[call] = []*ssa.Function{
					edge.Callee.Func,
				}
				baseMapSet[call] = map[*ssa.Function]bool{
					edge.Callee.Func: true,
				}
				continue
			}

			// otherwise, check for uniqueness
			_, isNotUnique := baseMapSet[call][edge.Callee.Func]

			// add to set
			if !isNotUnique {
				baseMap[call] = append(baseMap[call], edge.Callee.Func)
				baseMapSet[call][edge.Callee.Func] = true
			}
		}
	}

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

	if useConfig.Verbose {
		fmt.Println("Running pointer analysis...")
	}

	callPointer := FindCallPointer(packages)

	if useConfig.Verbose {
		fmt.Printf("Found mapping for %d calls\n", len(callPointer))
		fmt.Println("Running discover analysis...")
	}

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
