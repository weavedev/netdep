package stages

import (
	"fmt"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/output"
)

func CreateDependencyGraph(discoveredData *discovery.DiscoveredData) ([]*output.ServiceNode, []*output.ConnectionEdge) {
	nodes := make([]*output.ServiceNode, 0)
	edges := make([]*output.ConnectionEdge, 0)

	serviceMap := make(map[string]*output.ServiceNode)
	unknownCount := 0

	// create service nodes
	for _, serviceData := range discoveredData.ServCalls {
		serviceNode := &output.ServiceNode{
			ServiceName: serviceData.Service,
		}

		nodes = append(nodes, serviceNode)

		// save service name in a map for efficiency
		serviceMap[serviceData.Service] = serviceNode
	}

	// for each service node
	for _, serviceData := range discoveredData.ServCalls {
		sourceNode := serviceMap[serviceData.Service]

		// for each service this service calls
		for foundURL, serviceRemoteCalls := range serviceData.Calls {
			var targetNode *output.ServiceNode

			// determine target node
			if targetServiceName, ok := discoveredData.Handled[foundURL]; ok {
				targetNode = serviceMap[targetServiceName]
			} else {
				// create new unknown node (new one, otherwise the graph becomes distorted)
				unknownCount += 1
				serviceNode := &output.ServiceNode{
					ServiceName: fmt.Sprintf("Unknown Service #%d", unknownCount),
				}
				nodes = append(nodes, serviceNode)
				targetNode = serviceNode
			}

			// create all edges between these two nodes
			for _, remoteCall := range serviceRemoteCalls {
				connectionEdge := &output.ConnectionEdge{
					Call: output.NetworkCall{
						// TODO add more details
						Protocol:  "HTTP",
						URL:       foundURL,
						Arguments: nil,
						Location:  remoteCall,
					},
					Source: sourceNode,
					Target: targetNode,
				}

				edges = append(edges, connectionEdge)
			}
		}
	}

	return nodes, edges
}
