// Package stages defines different stages of analysis
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package stages

import (
	"encoding/json"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
	"sort"
)

/*
In the Building stages, the adjacency lists of each service are populated.
This is done by traversing the lists of endpoints/clients and looking for the other end of the connection.
The Building stages should handle
Refer to the Project plan, chapter 5.4 for more information.
*/

// Conn represents a connection tuple which consists of a service name
// and the amount of times a connection with / by that service was made.
type Conn struct {
	Service string `json:"service"`
	Amount  int    `json:"amount"`
}

type ServiceNode struct {
	ServiceName string
}

type Connection struct {
	Protocol  string             `json:"protocol"`
	Url       string             `json:"url"`
	Arguments []string           `json:"arguments"`
	Location  discovery.CallData `json:"location"`
}

type ServiceConnection struct {
	Service     ServiceNode  `json:"service"`
	Connection  []Connection `json:"connections"`
	Connections int          `json:"count"`
}

type ConnectionEdge struct {
	Connection Connection
	Source     *ServiceNode
	Target     *ServiceNode
}

// groupEdgesByServiceTargetAndSource creates a structure which you can use to query the edges using x[source][target] => array of edges
func groupEdgesByServiceTargetAndSource(edges []*ConnectionEdge) map[*ServiceNode]map[*ServiceNode][]*ConnectionEdge {
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
// Format of entries in the list is `"serviceName": [] Conn`
func ConstructAdjacencyList(nodes []*ServiceNode, edges []*ConnectionEdge) map[string][]ServiceConnection {
	adjacencyList := make(map[string][]ServiceConnection)
	groupedEdges := groupEdgesByServiceTargetAndSource(edges)

	for _, node := range nodes {
		adjacencyList[node.ServiceName] = make([]ServiceConnection, 0)

		// find the related edges in groupedEdges
		edgeSourceGroup, found := groupedEdges[node]

		if !found {
			continue
		}

		for targetServiceName, edgeGroup := range edgeSourceGroup {
			connectionList := make([]Connection, 0)

			for _, edge := range edgeGroup {
				connectionList = append(connectionList, edge.Connection)
			}

			// add the connection to the adjacencyList
			adjacencyList[node.ServiceName] = append(adjacencyList[node.ServiceName], ServiceConnection{
				Service:     *targetServiceName,
				Connection:  connectionList,
				Connections: len(connectionList),
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
func SerializeAdjacencyList(adjacencyList map[string][]ServiceConnection, pretty bool) (string, error) {
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
