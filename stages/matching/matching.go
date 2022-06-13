// Package matching constructs a graph from the found calls in the discovery stage
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package matching

import (
	"fmt"
	"go/token"
	"sort"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/output"
)

// createEmptyNodes create a set of services, but populates them to nil
func createEmptyNodes(calls []*callanalyzer.CallTarget, endpoints []*callanalyzer.CallTarget) (map[string]*output.ServiceNode, []*output.ServiceNode) {
	nodes := make([]*output.ServiceNode, 0)
	serviceMap := make(map[string]*output.ServiceNode)

	combinedList := calls
	combinedList = append(combinedList, endpoints...)

	// create nodes
	for _, call := range combinedList {
		if _, ok := serviceMap[call.ServiceName]; !ok {
			serviceNode := &output.ServiceNode{
				ServiceName: call.ServiceName,
				IsUnknown:   false,
			}

			nodes = append(nodes, serviceNode)
			// save service name in a map for efficiency
			serviceMap[call.ServiceName] = serviceNode
		}
	}

	return serviceMap, nodes
}

// createBasicPortMap finds all port definition and sets them for the corresponding service
func createBasicPortMap(endpoints []*callanalyzer.CallTarget) map[string]string {
	portMap := make(map[string]string)

	// find port definitions
	for _, call := range endpoints {
		if len(call.RequestLocation) >= 1 && call.RequestLocation[0] == ':' {
			portMap[call.ServiceName] = call.RequestLocation
		}
	}

	return portMap
}

// createEndpointMap create a map of an endpoint to a service name
// TODO: this URL is very rudimentary currently
func createEndpointMap(endpoints []*callanalyzer.CallTarget) map[string]string {
	endpointMap := make(map[string]string)
	portMap := createBasicPortMap(endpoints)

	for _, call := range endpoints {
		// set default port
		if _, ok := portMap[call.ServiceName]; !ok {
			portMap[call.ServiceName] = ":80"
		}

		port := portMap[call.ServiceName]

		if call.RequestLocation == "" || call.RequestLocation[0] == '/' {
			// register request
			endpointURL := fmt.Sprintf("http://%s%s%s", call.ServiceName, port, call.RequestLocation)
			endpointMap[endpointURL] = call.ServiceName
		}
	}

	return endpointMap
}

func populateTraceMap(calls []*callanalyzer.CallTarget, traceMap map[token.Pos]int) {
	for _, target := range calls {
		for _, call := range target.Trace {
			if _, ok := traceMap[call.Pos]; !ok {
				traceMap[call.Pos] = 1
			} else {
				traceMap[call.Pos]++
			}
		}
	}
}

func getRelevantCallLocation(calls []callanalyzer.CallTargetTrace, traceMap map[token.Pos]int) []string {
	ret := make([]string, 0)
	//min := -1
	for _, call := range calls {
		//if !call.Internal {
		//	break
		//}
		//n, _ := traceMap[call.Pos]
		//
		//if min == -1 || n < min {
		//	ret = make([]string, 0)
		//	min = n
		//}
		//
		//if n > min {
		//	continue
		//}

		ret = append(ret, fmt.Sprintf("%s:%s", call.FileName, call.PositionInFile))
	}

	return ret
}

// CreateDependencyGraph creates the nodes and edges of a dependency graph, given the discovered calls and endpoints
func CreateDependencyGraph(calls []*callanalyzer.CallTarget, endpoints []*callanalyzer.CallTarget) output.NodeGraph {
	UnknownService := &output.ServiceNode{
		ServiceName: "UnknownService",
		IsUnknown:   true,
	}

	traceMap := map[token.Pos]int{}
	edges := make([]*output.ConnectionEdge, 0)
	serviceMap, nodes := createEmptyNodes(calls, endpoints)
	endpointMap := createEndpointMap(endpoints)
	hasUnknown := false

	populateTraceMap(calls, traceMap)
	populateTraceMap(endpoints, traceMap)

	// Add edges (eg. matching)
	// This order is guaranteed because calls is an array
	for _, call := range calls {
		sourceNode := serviceMap[call.ServiceName]
		targetServiceName, isResolved := findTargetNodeName(call, endpointMap)

		var targetNode *output.ServiceNode

		if target, ok := serviceMap[targetServiceName]; ok && isResolved {
			targetNode = target
			// Set target to UnknownService if not found
			// There are 3 possibilities for this scenario:
			// 1. endpoint definition of call.RequestLocation wasn't resolved correctly.
			// 2. call.RequestLocation references external API, which is not contained
			// in the endpointMap.
			// 3. The call.RequestLocation itself was not resolved correctly.
			// In the future this distinction could be made.
		} else {
			targetNode = UnknownService
		}

		// If at least one unknown target has been found,
		// add it to a list of nodes.
		if targetNode == UnknownService && !hasUnknown {
			nodes = append(nodes, UnknownService)
			hasUnknown = true
		}

		callLocations := getRelevantCallLocation(call.Trace, traceMap)

		connectionEdge := &output.ConnectionEdge{
			Call: output.NetworkCall{
				// TODO add more details
				Protocol: "HTTP",
				URL:      call.RequestLocation,
				// TODO: handle stack trace?
				Locations: callLocations,
			},
			Source: sourceNode,
			Target: targetNode,
		}

		edges = append(edges, connectionEdge)
	}

	// ensure alphabetical order for nodes (to prevent flaky tests)
	sort.Slice(nodes, func(i, j int) bool {
		x := nodes[i]
		y := nodes[j]

		return x.ServiceName < y.ServiceName
	})

	return output.NodeGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

// findTargetNodeName returns a name of the target service
// If call was unresolved in discovery stage it by default returns false
// If call was resolved, but the endpointMap does not contain the URL,
// then empty string and false is returned.
// Otherwise, a name of the target service is returned.
func findTargetNodeName(call *callanalyzer.CallTarget, endpointMap map[string]string) (string, bool) {
	if !call.IsResolved {
		return "", false
	}

	// TODO improve matching, compare URL
	// TODO handle dynamic urls like "/_var"
	targetServiceName, hasTarget := endpointMap[call.RequestLocation]

	return targetServiceName, hasTarget
}
