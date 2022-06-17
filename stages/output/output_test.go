// Package output
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package output

import (
	"strings"
	"testing"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"

	"github.com/stretchr/testify/assert"
)

// createSmallTestGraph creates a graph with three nodes, where node 1 had edges to node 2 and 3, and node 2 to node 3
func createSmallTestGraph() NodeGraph {
	node1 := ServiceNode{
		ServiceName:   "Node1",
		IsReferenced:  false,
		IsReferencing: true,
	}
	node2 := ServiceNode{
		ServiceName:   "Node2",
		IsReferenced:  true,
		IsReferencing: true,
	}
	node3 := ServiceNode{
		ServiceName:   "Node3",
		IsReferenced:  true,
		IsReferencing: false,
	}

	edge12 := ConnectionEdge{
		Call: NetworkCall{
			Protocol: "HTTP",
			URL:      "",
		},
		Source: &node1,
		Target: &node2,
	}

	edge13 := ConnectionEdge{
		Call: NetworkCall{
			Protocol: "HTTP",
			URL:      "",
		},
		Source: &node1,
		Target: &node3,
	}

	edge23 := ConnectionEdge{
		Call: NetworkCall{
			Protocol: "HTTP",
			URL:      "",
		},
		Source: &node2,
		Target: &node3,
	}

	nodes := []*ServiceNode{&node1, &node2, &node3}
	edges := []*ConnectionEdge{&edge12, &edge13, &edge23}

	return NodeGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

// Tests adjacency list construction based on data found in discovery stage
func TestConstructAdjacencyList(t *testing.T) {
	graph := createSmallTestGraph()
	res := ConstructAdjacencyList(graph)

	expected := AdjacencyList{
		"Node1": {
			{
				// node 2
				Service:       graph.Nodes[1].ServiceName,
				Calls:         []NetworkCall{graph.Edges[0].Call},
				NumberOfCalls: 1,
			},
			{
				// node 3
				Service:       graph.Nodes[2].ServiceName,
				Calls:         []NetworkCall{graph.Edges[1].Call},
				NumberOfCalls: 1,
			},
		},
		"Node2": {
			{
				// node 3
				Service:       graph.Nodes[2].ServiceName,
				Calls:         []NetworkCall{graph.Edges[2].Call},
				NumberOfCalls: 1,
			},
		},
		"Node3": {},
	}

	assert.Equal(t, expected, res)
}

// TestSerialiseOutputNull performs a sanity check for nil case
func TestSerialiseOutputNull(t *testing.T) {
	str, _ := SerializeAdjacencyList(nil, false)
	assert.Equal(t, "null", str)
}

// TestSerialiseOutput test a realistic output of a serialisation. More of a sanity check.
func TestSerialiseOutput(t *testing.T) {
	graph := createSmallTestGraph()
	list := ConstructAdjacencyList(graph)
	str, _ := SerializeAdjacencyList(list, false)
	expected := "{\"Node1\":[{\"service\":\"Node2\",\"calls\":[{\"protocol\":\"HTTP\",\"locations\":null}],\"count\":1},{\"service\":\"Node3\",\"calls\":[{\"protocol\":\"HTTP\",\"locations\":null}],\"count\":1}],\"Node2\":[{\"service\":\"Node3\",\"calls\":[{\"protocol\":\"HTTP\",\"locations\":null}],\"count\":1}],\"Node3\":[]}"
	assert.Equal(t, expected, str)
}

// TestPrintDiscoveredAnnotations test discovered annotation printing
func TestPrintDiscoveredAnnotations(t *testing.T) {
	annotations := make(map[string]map[callanalyzer.Position]string)
	annotations["a"] = make(map[callanalyzer.Position]string)
	annotations["b"] = make(map[callanalyzer.Position]string)

	pos1 := callanalyzer.Position{
		Filename: "d1",
		Line:     5,
	}
	pos2 := callanalyzer.Position{
		Filename: "d2",
		Line:     6,
	}
	pos3 := callanalyzer.Position{
		Filename: "a2",
		Line:     62,
	}

	annotations["a"][pos1] = "valuee1"
	annotations["a"][pos2] = "valuee2"
	annotations["b"][pos3] = "valuee3"

	str := PrintDiscoveredAnnotations(annotations)
	assert.True(t, strings.Contains(str, "Discovered annotations:"))
	assert.True(t, strings.Contains(str, "valuee1"))
	assert.True(t, strings.Contains(str, "valuee2"))
	assert.True(t, strings.Contains(str, "valuee3"))
}

// TestPrintDiscoveredAnnotationsEmpty test discovered annotation printing
func TestPrintDiscoveredAnnotationsEmpty(t *testing.T) {
	annotations := make(map[string]map[callanalyzer.Position]string)
	annotations["a"] = make(map[callanalyzer.Position]string)
	annotations["b"] = make(map[callanalyzer.Position]string)

	str := PrintDiscoveredAnnotations(annotations)
	assert.Equal(t, str, "[Discovered none]")
}

func TestConstructUnusedServicesLists(t *testing.T) {
	graph := createSmallTestGraph()
	node4 := ServiceNode{
		ServiceName:   "Node4",
		IsReferenced:  false,
		IsReferencing: false,
	}
	graph.Nodes = append(graph.Nodes, &node4)
	allServices := []string{"Node1", "Node2", "Node3", "Node5"}
	noReferenceToServices, noReferenceToAndFromServices := ConstructUnusedServicesLists(graph.Nodes, allServices)

	assert.Equal(t, []string{"Node1", "Node4", "Node5"}, noReferenceToServices)
	assert.Equal(t, []string{"Node4", "Node5"}, noReferenceToAndFromServices)
}

// TestPrintMethods runs untestable methods that only print so the coverage isn't affected
func TestPrintMethods(t *testing.T) {
	noReferenceToServices := []string{"svc1"}
	noReferenceToAndFromServices := []string{"svc2", "svc3"}
	PrintUnusedServices(noReferenceToServices, noReferenceToAndFromServices)

	trace := callanalyzer.CallTargetTrace{FileName: "file", PositionInFile: "1"}
	traces := []callanalyzer.CallTargetTrace{trace}
	callTarget := callanalyzer.CallTarget{
		PackageName: "pkg", MethodName: "method", RequestLocation: "url",
		IsResolved: true, ServiceName: "svc1", TargetSvc: "svc2", Trace: traces,
	}
	targets := make([]*callanalyzer.CallTarget, 0)
	targets = append(targets, &callTarget)
	PrintAnnotationSuggestions(targets)
}
