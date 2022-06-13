// Package output
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package output

import (
	"strings"
	"testing"

	"lab.weave.nl/internships/tud-2022/netDep/stages/preprocessing"

	"github.com/stretchr/testify/assert"
)

// createSmallTestGraph creates a graph with three nodes, where node 1 had edges to node 2 and 3, and node 2 to node 3
func createSmallTestGraph() NodeGraph {
	node1 := ServiceNode{
		ServiceName: "Node1",
	}
	node2 := ServiceNode{
		ServiceName: "Node2",
	}
	node3 := ServiceNode{
		ServiceName: "Node3",
	}

	edge12 := ConnectionEdge{
		Call: NetworkCall{
			Protocol:  "HTTP",
			URL:       "",
			Arguments: nil,
			Location:  "",
		},
		Source: &node1,
		Target: &node2,
	}

	edge13 := ConnectionEdge{
		Call: NetworkCall{
			Protocol:  "HTTP",
			URL:       "",
			Arguments: nil,
			Location:  "",
		},
		Source: &node1,
		Target: &node3,
	}

	edge23 := ConnectionEdge{
		Call: NetworkCall{
			Protocol:  "HTTP",
			URL:       "",
			Arguments: nil,
			Location:  "",
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
	expected := "{\"Node1\":[{\"service\":\"Node2\",\"calls\":[{\"protocol\":\"HTTP\",\"location\":\"\"}],\"count\":1},{\"service\":\"Node3\",\"calls\":[{\"protocol\":\"HTTP\",\"location\":\"\"}],\"count\":1}],\"Node2\":[{\"service\":\"Node3\",\"calls\":[{\"protocol\":\"HTTP\",\"location\":\"\"}],\"count\":1}],\"Node3\":[]}"
	assert.Equal(t, expected, str)
}

// TestPrintDiscoveredAnnotations test discovered annotation printing
func TestPrintDiscoveredAnnotations(t *testing.T) {
	annotations := make(map[string]map[preprocessing.Position]string)
	annotations["a"] = make(map[preprocessing.Position]string)
	annotations["b"] = make(map[preprocessing.Position]string)

	pos1 := preprocessing.Position{
		Filename: "d1",
		Line:     5,
	}
	pos2 := preprocessing.Position{
		Filename: "d2",
		Line:     6,
	}
	pos3 := preprocessing.Position{
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
	annotations := make(map[string]map[preprocessing.Position]string)
	annotations["a"] = make(map[preprocessing.Position]string)
	annotations["b"] = make(map[preprocessing.Position]string)

	str := PrintDiscoveredAnnotations(annotations)
	assert.Equal(t, str, "[Discovered none]")
}
