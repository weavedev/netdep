package stages

import (
	"fmt"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/output"
)

// CreateDependencyGraph creates the nodes and edges of a dependency graph, given the discovered calls and endpoints
func CreateDependencyGraph(calls []*callanalyzer.CallTarget, endpoints []*callanalyzer.CallTarget) ([]*output.ServiceNode, []*output.ConnectionEdge) {
	endpointMap := make(map[string]string)
	portMap := make(map[string]string)
	serviceMap := make(map[string]*output.ServiceNode)
	unknownServiceMap := make(map[string]*output.ServiceNode)

	nodes := make([]*output.ServiceNode, 0)
	edges := make([]*output.ConnectionEdge, 0)

	unknownCount := 0

	// create nodes
	for _, call := range calls {
		if _, ok := serviceMap[call.ServiceName]; !ok {
			serviceMap[call.ServiceName] = nil
		}
	}

	// find port definitions
	for _, call := range endpoints {
		if call.RequestLocation[0] == ':' {
			portMap[call.ServiceName] = call.RequestLocation
		}
	}

	for _, call := range endpoints {
		// register node
		if _, ok := serviceMap[call.ServiceName]; !ok {
			serviceMap[call.ServiceName] = nil
		}

		// set default port
		if _, ok := portMap[call.ServiceName]; !ok {
			portMap[call.ServiceName] = ":80"
		}

		port := portMap[call.ServiceName]

		if call.RequestLocation[0] == '/' {
			// register request
			endpointURL := fmt.Sprintf("http://%s%s%s", call.ServiceName, port, call.RequestLocation)
			endpointMap[endpointURL] = call.ServiceName
		}
	}

	// Populate nodes
	for serviceName, _ := range serviceMap {
		serviceNode := &output.ServiceNode{
			ServiceName: serviceName,
			IsUnknown:   false,
		}

		nodes = append(nodes, serviceNode)
		// save service name in a map for efficiency
		serviceMap[serviceName] = serviceNode
	}

	// Add edges
	for _, call := range calls {
		sourceNode := serviceMap[call.ServiceName]

		// TODO improve matching, compare URL
		// TODO handle dynamic urls like "/_var"
		targetServiceName, hasTarget := endpointMap[call.RequestLocation]
		var targetNode *output.ServiceNode

		// find target service
		if hasTarget {
			targetService, exists := serviceMap[targetServiceName]

			if exists {
				targetNode = targetService
			}
		}

		if targetNode == nil {
			// set unknown service
			// re-use if url is already used before
			oldUnknownService, exists := unknownServiceMap[call.RequestLocation]

			if exists {
				targetNode = oldUnknownService
			} else {
				// create new unknown node (new one, otherwise the graph becomes distorted)
				unknownCount += 1
				serviceNode := &output.ServiceNode{
					ServiceName: fmt.Sprintf("Unknown Service #%d", unknownCount),
					IsUnknown:   true,
				}
				nodes = append(nodes, serviceNode)
				targetNode = serviceNode
				unknownServiceMap[call.RequestLocation] = serviceNode
			}
		}

		connectionEdge := &output.ConnectionEdge{
			Call: output.NetworkCall{
				// TODO add more details
				Protocol:  "HTTP",
				URL:       call.RequestLocation,
				Arguments: nil,
				// TODO: handle stack trace?
				Location: fmt.Sprintf("%s:%s", call.FileName, call.PositionInFile),
			},
			Source: sourceNode,
			Target: targetNode,
		}

		edges = append(edges, connectionEdge)
	}

	return nodes, edges
}
