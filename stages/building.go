// Package stages defines different stages of analysis
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package stages

import (
	"encoding/json"
	"fmt"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
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

// ConstructOutput constructs and returns the output of the tool as a string in JSON format.
func ConstructOutput(discoveredData *discovery.DiscoveredData) (string, string) {
	adjList := ConstructAdjacencyList(discoveredData)
	servCalls := discoveredData.ServCalls
	return SerialiseOutput(adjList, servCalls)
}

// ConstructAdjacencyList constructs an adjacency list of service dependencies.
// Format of entries in the list is `"serviceName": [] Conn`
func ConstructAdjacencyList(data *discovery.DiscoveredData) map[string][]Conn {
	m := make(map[string][]Conn)

	// Assuming DiscoveredData is something like currently defined in stages/discovery.go
	for _, servCall := range data.ServCalls {
		for k, callDataList := range servCall.Calls {
			if targetServ, ok := data.Handled[k]; ok {
				m[servCall.Service] = append(m[servCall.Service], Conn{targetServ, len(callDataList)})
			} else {
				m[servCall.Service] = append(m[servCall.Service], Conn{"Unknown Service", len(callDataList)})
			}
		}
	}

	return m
}

// SerialiseOutput serialises the given adjacency list and returns the output as a string in JSON format.
func SerialiseOutput(adjList map[string][]Conn, servCalls []discovery.ServiceCalls) (string, string) {
	outAdj, err := json.MarshalIndent(adjList, "", "\t")
	if err != nil {
		fmt.Println("JSON encode error")
		return "", ""
	}

	outCalls, err := json.MarshalIndent(servCalls, "", "\t")
	if err != nil {
		fmt.Println("JSON encode error")
		return string(outAdj), ""
	}

	return string(outAdj), string(outCalls)
}
