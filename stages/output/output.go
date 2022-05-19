// Package output defines the different ways of output in the tool
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package output

import (
	"encoding/json"
	"sort"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
)

/*
In the Building stages, the adjacency lists of each service are populated.
This is done by traversing the lists of endpoints/clients and looking for the other end of the connection.
The Building stages should handle
Refer to the Project plan, chapter 5.4 for more information.
*/

// NetworkCall represents a call that can be made in the network
type NetworkCall struct {
	Protocol  string             `json:"protocol"`
	URL       string             `json:"url"`
	Arguments []string           `json:"arguments"`
	Location  discovery.CallData `json:"location"`
}

// ServiceNode represents a node in the output graph, which is a Service
type ServiceNode struct {
	ServiceName string `json:"serviceName"`
}

// ConnectionEdge represents a directed edge in the output graph
type ConnectionEdge struct {
	Call   NetworkCall
	Source *ServiceNode
	Target *ServiceNode
}

// ServiceCallList holds the NetworkCall's related to a Service, used in the AdjacencyList
type ServiceCallList struct {
	Service       ServiceNode   `json:"service"`
	Calls         []NetworkCall `json:"calls"`
	NumberOfCalls int           `json:"count"`
}

type (
	AdjacencyList  map[string][]ServiceCallList
	GroupedEdgeMap map[*ServiceNode]map[*ServiceNode][]*ConnectionEdge
)

// groupEdgesByServiceTargetAndSource creates a structure which you can use to query the edges using x[source][target] => array of edges
func groupEdgesByServiceTargetAndSource(edges []*ConnectionEdge) GroupedEdgeMap {
	outputMap := make(map[*ServiceNode]map[*ServiceNode][]*ConnectionEdge)

	for _, edge := range edges {
		sourceMap, hasSourceMap := outputMap[edge.Source]

		// create the target structure
		if !hasSourceMap {
			outputMap[edge.Source] = make(map[*ServiceNode][]*ConnectionEdge)
			sourceMap = outputMap[edge.Source]
		}

		targetList, hasTargetList := sourceMap[edge.Target]

		// create the connection list
		if !hasTargetList {
			sourceMap[edge.Target] = make([]*ConnectionEdge, 0)
		}

		// add the edge to the right group
		sourceMap[edge.Target] = append(targetList, edge)
	}

	return outputMap
}

// ConstructAdjacencyList constructs an adjacency list of service dependencies.
// In its current representation this is a map to a list of adjacent nodes.
func ConstructAdjacencyList(nodes []*ServiceNode, edges []*ConnectionEdge) AdjacencyList {
	adjacencyList := make(map[string][]ServiceCallList)
	groupedEdges := groupEdgesByServiceTargetAndSource(edges)

	for _, node := range nodes {
		adjacencyList[node.ServiceName] = make([]ServiceCallList, 0)

		// find the related edges in groupedEdges
		edgeSourceGroup, found := groupedEdges[node]

		if !found {
			continue
		}

		for targetServiceName, edgeGroup := range edgeSourceGroup {
			callList := make([]NetworkCall, 0)

			for _, edge := range edgeGroup {
				callList = append(callList, edge.Call)
			}

			// add the connection to the adjacencyList
			adjacencyList[node.ServiceName] = append(adjacencyList[node.ServiceName], ServiceCallList{
				Service:       *targetServiceName,
				Calls:         callList,
				NumberOfCalls: len(callList),
			})
		}

		// sort the list based on service name, to make sure the order is always the same (for testing)
		sort.Slice(adjacencyList[node.ServiceName], func(i, j int) bool {
			x := adjacencyList[node.ServiceName][i]
			y := adjacencyList[node.ServiceName][j]

			return x.Service.ServiceName < y.Service.ServiceName
		})
	}

	return adjacencyList
}

// SerializeAdjacencyList serialises a given adjacencyList in JSON format
func SerializeAdjacencyList(adjacencyList AdjacencyList, pretty bool) (string, error) {
	var output []byte
	var err error

	if pretty {
		output, err = json.MarshalIndent(adjacencyList, "", "\t")
	} else {
		output, err = json.Marshal(adjacencyList)
	}

	if err != nil {
		return "null", err
	}

	return string(output), err
}