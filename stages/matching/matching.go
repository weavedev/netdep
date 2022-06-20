// Package matching constructs a graph from the found calls in the discovery stage
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package matching

import (
	"fmt"
	"sort"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/natsanalyzer"
	"lab.weave.nl/internships/tud-2022/netDep/structures"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"
	"lab.weave.nl/internships/tud-2022/netDep/stages/output"
)

// createEmptyNodes create a set of services, but populates them to nil
func createEmptyNodes(dependencies *structures.Dependencies) (map[string]*output.ServiceNode, []*output.ServiceNode) {
	nodes := make([]*output.ServiceNode, 0)
	serviceMap := make(map[string]*output.ServiceNode)

	combinedList := dependencies.Calls
	combinedList = append(combinedList, dependencies.Endpoints...)

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

	combinedNats := dependencies.Consumers
	combinedNats = append(combinedNats, dependencies.Producers...)

	// extend nodes with NATS only services
	for _, call := range combinedNats {
		if _, ok := serviceMap[call.ServiceName]; !ok {
			serviceNode := &output.ServiceNode{
				ServiceName: call.ServiceName,
				IsUnknown:   false,
			}

			nodes = append(nodes, serviceNode)
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
		} else if call.PackageName == "servicecalls" {
			endpointURL := call.RequestLocation
			endpointMap[endpointURL] = call.ServiceName
		}
	}

	return endpointMap
}

// CreateDependencyGraph creates the nodes and edges of a dependency graph, given the discovered calls and endpoints
func CreateDependencyGraph(dependencies *structures.Dependencies) output.NodeGraph {
	if dependencies == nil {
		return output.NodeGraph{Nodes: make([]*output.ServiceNode, 0), Edges: make([]*output.ConnectionEdge, 0)}
	}

	UnknownService := &output.ServiceNode{
		ServiceName:   "UnknownService",
		IsUnknown:     true,
		IsReferenced:  true,
		IsReferencing: true,
	}

	edges := make([]*output.ConnectionEdge, 0)
	serviceMap, nodes := createEmptyNodes(dependencies)
	endpointMap := createEndpointMap(dependencies.Endpoints)
	hasUnknown := false
	edges = append(edges, extendWithNats(dependencies.Consumers, dependencies.Producers, &hasUnknown, serviceMap, &nodes)...)

	// Add edges (eg. matching). This order is guaranteed because calls is an array
	for _, call := range dependencies.Calls {
		sourceNode := serviceMap[call.ServiceName]
		sourceNode.IsReferencing = true
		targetServiceName, isResolved := findTargetNodeName(call, endpointMap)

		var targetNode *output.ServiceNode

		if target, ok := serviceMap[targetServiceName]; ok && isResolved {
			targetNode = target
			targetNode.IsReferenced = true
			// Set target to UnknownService if not found. There are 3 possibilities for this scenario:
			// 1. endpoint definition of call.RequestLocation wasn't resolved correctly.
			// 2. call.RequestLocation references external API, which is not contained in the endpointMap.
			// 3. The call.RequestLocation itself was not resolved correctly. In the future this distinction could be made.
		} else {
			targetNode = UnknownService
		}
		// In case of servicecalls scanning some services use the methods in module or proto definitions which
		// Make the tool think that it's a self reference. This is a clear case of false positives.
		if targetNode.ServiceName == sourceNode.ServiceName {
			continue
		}
		// If at least one unknown target has been found,
		// add it to a list of nodes.
		if targetNode == UnknownService && !hasUnknown {
			nodes = append(nodes, UnknownService)
			hasUnknown = true
		}

		// Default values
		protocol := "HTTP"
		url := call.RequestLocation
		methodName := ""

		// If the call was discovered via servicecalls package scanning
		// Edit the values with the servicecalls specific data
		if call.PackageName == "servicecalls" {
			protocol = call.PackageName
			url = ""
			methodName = call.RequestLocation
		}

		connectionEdge := &output.ConnectionEdge{
			Call: output.NetworkCall{
				Protocol:   protocol,
				URL:        url,
				Arguments:  nil,
				MethodName: methodName,
				Locations:  call.TraceAsStringArray(),
			},
			Source: sourceNode,
			Target: targetNode,
		}

		edges = append(edges, connectionEdge)
	}
	// ensure alphabetical order for nodes (to prevent flaky tests)
	sortNodes(&nodes)
	return output.NodeGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

// sortNodes sorts the nodes in alphabetical order of their service names.
func sortNodes(nodes *[]*output.ServiceNode) {
	sort.Slice(*nodes, func(i, j int) bool {
		x := (*nodes)[i]
		y := (*nodes)[j]

		return x.ServiceName < y.ServiceName
	})
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

// extendWithNats extends the Connection Edges data structure
// with discovered NATS edges. This was required, because NATS
// can have one-to-many dependencies, where as something like HTTP
// is one-to-one.
func extendWithNats(consumers []*natsanalyzer.NatsCall, producers []*natsanalyzer.NatsCall, hasUnknown *bool, services map[string]*output.ServiceNode, nodes *[]*output.ServiceNode) []*output.ConnectionEdge {
	edges := make([]*output.ConnectionEdge, 0)

	if consumers == nil || producers == nil {
		return edges
	}

	UnknownService := &output.ServiceNode{
		ServiceName: "UnknownService",
		IsUnknown:   true,
	}

	// for each producer we  find all  the consumer
	// if it has no consumer, we mark it as an edge
	// with unknown target. There is a small chance
	// that producer's messages are never consumed,
	// but the greater chance is that its consumer
	// was not discovered.
	for _, producer := range producers {
		hasConsumer := false
		for _, consumer := range consumers {
			if producer.Subject == consumer.Subject {
				hasConsumer = true
				edge := &output.ConnectionEdge{
					Call: output.NetworkCall{
						// TODO add more details
						Protocol:  producer.Communication,
						URL:       producer.Subject,
						Arguments: nil,
						// TODO: handle stack trace?
						Locations: []string{fmt.Sprintf("%s:%s", producer.FileName, producer.PositionInFile)},
					},
					// Always hits, because services was populated using consumers and producers
					Source: services[producer.ServiceName],
					Target: services[consumer.ServiceName],
				}

				edges = append(edges, edge)
			}
		}

		if !hasConsumer {
			edge := &output.ConnectionEdge{
				Call: output.NetworkCall{
					// TODO add more details
					Protocol:  producer.Communication,
					URL:       producer.Subject,
					Arguments: nil,
					// TODO: handle stack trace?
					Locations: []string{fmt.Sprintf("%s:%s", producer.FileName, producer.PositionInFile)},
				},
				// Always hits, because services was populated using consumers and producers
				Source: services[producer.ServiceName],
				Target: UnknownService,
			}

			if !*hasUnknown {
				*nodes = append(*nodes, UnknownService)
			}

			edges = append(edges, edge)
			*hasUnknown = true
		}
	}

	return edges
}
