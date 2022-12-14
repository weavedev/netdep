// Package output defines the different ways of output in the tool
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"
)

/*
In the Building stages, the adjacency lists of each service are populated.
This is done by traversing the lists of endpoints/clients and looking for the other end of the connection.
The Building stages should handle
Refer to the Project plan, chapter 5.4 for more information.
*/

// NetworkCall represents a call that can be made in the network
type NetworkCall struct {
	Protocol   string   `json:"protocol"`
	URL        string   `json:"url,omitempty"`
	MethodName string   `json:"methodName,omitempty"`
	Arguments  []string `json:"arguments,omitempty"`
	Locations  []string `json:"locations"`
}

// ServiceNode represents a node in the output graph, which is a Service
type ServiceNode struct {
	ServiceName   string `json:"serviceName"`
	IsUnknown     bool   `json:"isUnknown"`
	IsReferenced  bool   `json:"isReferenced"`
	IsReferencing bool   `json:"isReferencing"`
	// Hostname    []string `json:"hostname"`
	// Endpoints   []string `json:"endpoints"`
}

// ConnectionEdge represents a directed edge in the output graph
type ConnectionEdge struct {
	Call   NetworkCall
	Source *ServiceNode
	Target *ServiceNode
}

// ServiceCallList holds the NetworkCall's related to a Service, used in the AdjacencyList
type ServiceCallList struct {
	Service       string        `json:"service"`
	Calls         []NetworkCall `json:"calls"`
	NumberOfCalls int           `json:"count"`
}

type NodeGraph struct {
	Nodes []*ServiceNode
	Edges []*ConnectionEdge
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
func ConstructAdjacencyList(graph NodeGraph) AdjacencyList {
	adjacencyList := make(map[string][]ServiceCallList)
	groupedEdges := groupEdgesByServiceTargetAndSource(graph.Edges)

	for _, node := range graph.Nodes {
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
				Service:       targetServiceName.ServiceName,
				Calls:         callList,
				NumberOfCalls: len(callList),
			})
		}

		// sort the list based on service name, to make sure the order is always the same (for testing)
		sort.Slice(adjacencyList[node.ServiceName], func(i, j int) bool {
			x := adjacencyList[node.ServiceName][i]
			y := adjacencyList[node.ServiceName][j]

			return x.Service < y.Service
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

// PrintAnnotationSuggestions prints suggestions to add annotations for the list of callanalyzer.CallTarget it's provided.
// Intended to be used for unresolved targets.
func PrintAnnotationSuggestions(targets []*callanalyzer.CallTarget) {
	for _, target := range targets {
		color.HiCyan("%s:%s couldn't be resolved. ", target.Trace[0].FileName, target.Trace[0].PositionInFile)
		color.HiCyan("Add an annotation above it in the format \"//netdep:client ...\" or \"//netdep:endpoint ...\"")
	}
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

// ConstructUnusedServicesLists constructs 2 lists containing unreferenced services and unreferenced
// services that don't make any calls to other services, respectively
func ConstructUnusedServicesLists(services []*ServiceNode, allServices []string) ([]string, []string) {
	var noReferenceToServices []string
	var noReferenceToAndFromServices []string
	// list of services in which dependencies have been discovered
	var servicesInGraph []string

	for _, service := range services {
		if !service.IsReferenced && !service.IsReferencing {
			noReferenceToAndFromServices = append(noReferenceToAndFromServices, service.ServiceName)
		}
		if !service.IsReferenced {
			noReferenceToServices = append(noReferenceToServices, service.ServiceName)
		}
		servicesInGraph = append(servicesInGraph, service.ServiceName)
	}

	var allServiceNames []string
	for _, service := range allServices {
		// get only the service names from the absolute paths stored in allServices
		allServiceNames = append(allServiceNames, service[strings.LastIndex(service, string(os.PathSeparator))+1:])
	}
	for _, service := range allServiceNames {
		// add services in which no dependencies have been found
		if !contains(servicesInGraph, service) {
			noReferenceToAndFromServices = append(noReferenceToAndFromServices, service)
			noReferenceToServices = append(noReferenceToServices, service)
		}
	}
	return noReferenceToServices, noReferenceToAndFromServices
}

// PrintUnusedServices prints all the unused services
func PrintUnusedServices(noReferenceToServices []string, noReferenceToAndFromServices []string) {
	color.HiCyan("Unreferenced services: ")
	for _, service := range noReferenceToServices {
		color.HiWhite("\t%s\n", service)
	}
	color.HiCyan("Unreferenced services that don't make any calls: ")
	for _, service := range noReferenceToAndFromServices {
		color.HiWhite("\t%s\n", service)
	}
}

// PrintDiscoveredAnnotations prints all the discovered annotations if the tool was run with the verbose flag.
func PrintDiscoveredAnnotations(annotations map[string]map[callanalyzer.Position]string) string {
	type Annotation struct {
		ServiceName string
		Position    string
		Value       string
	}

	annotationList := make([]*Annotation, 0)

	for serName, serMap := range annotations {
		for pos, val := range serMap {
			ann := &Annotation{
				ServiceName: serName,
				Position:    fmt.Sprintf("%s:%s", pos.Filename, strconv.Itoa(pos.Line)),
				Value:       val,
			}
			annotationList = append(annotationList, ann)
		}
	}

	discoveredAnnotations := ""

	if len(annotationList) != 0 {
		discoveredAnnotations += "Discovered annotations:\n\n"

		for _, ann := range annotationList {
			discoveredAnnotations += "Service name: " + ann.ServiceName + "\n"
			discoveredAnnotations += "Position: " + ann.Position + "\n"
			discoveredAnnotations += "Value: " + ann.Value + "\n\n"
		}
	}

	color.Magenta(discoveredAnnotations)
	return discoveredAnnotations
}
