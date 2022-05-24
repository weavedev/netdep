// Package output
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateSmallTestGraph() NodeGraph {
	node1 := ServiceNode{
		ServiceName: "Node1",
		IsUnknown:   false,
	}
	node2 := ServiceNode{
		ServiceName: "Node2",
		IsUnknown:   false,
	}
	node3 := ServiceNode{
		ServiceName: "Node3",
		IsUnknown:   false,
	}

	edge12 := ConnectionEdge{
		Call: NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node2:80/URL_2",
			Arguments: nil,
			Location:  "./node1/path/to/some/file.go:24",
		},
		Source: &node1,
		Target: &node2,
	}

	edge13 := ConnectionEdge{
		Call: NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node3:80/URL_3",
			Arguments: nil,
			Location:  "./node1/path/to/some/other/file.go:36",
		},
		Source: &node1,
		Target: &node3,
	}

	edge23a := ConnectionEdge{
		Call: NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node3:80/URL_3",
			Arguments: nil,
			Location:  "./node1/path/to/some/file.go:245",
		},
		Source: &node2,
		Target: &node3,
	}

	edge23b := ConnectionEdge{
		Call: NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node3:80/URL_3",
			Arguments: nil,
			Location:  "./node2/path/to/some/other/file.go:436",
		},
		Source: &node2,
		Target: &node3,
	}

	nodes := []*ServiceNode{&node1, &node2, &node3}
	edges := []*ConnectionEdge{&edge12, &edge13, &edge23a, &edge23b}

	return NodeGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

// Tests adjacency list construction based on data found in discovery stage
func TestConstructAdjacencyList(t *testing.T) {
	graph := CreateSmallTestGraph()
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
				Calls:         []NetworkCall{graph.Edges[2].Call, graph.Edges[3].Call},
				NumberOfCalls: 2,
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
	graph := CreateSmallTestGraph()
	list := ConstructAdjacencyList(graph)
	str, _ := SerializeAdjacencyList(list, false)
	expected := "{\"Node1\":[{\"service\":\"Node2\",\"calls\":[{\"protocol\":\"HTTP\",\"url\":\"http://Node2:80/URL_2\",\"arguments\":null,\"location\":\"./node1/path/to/some/file.go:24\"}],\"count\":1},{\"service\":\"Node3\",\"calls\":[{\"protocol\":\"HTTP\",\"url\":\"http://Node3:80/URL_3\",\"arguments\":null,\"location\":\"./node1/path/to/some/other/file.go:36\"}],\"count\":1}],\"Node2\":[{\"service\":\"Node3\",\"calls\":[{\"protocol\":\"HTTP\",\"url\":\"http://Node3:80/URL_3\",\"arguments\":null,\"location\":\"./node1/path/to/some/file.go:245\"},{\"protocol\":\"HTTP\",\"url\":\"http://Node3:80/URL_3\",\"arguments\":null,\"location\":\"./node2/path/to/some/other/file.go:436\"}],\"count\":2}],\"Node3\":[]}"

	assert.Equal(t, expected, str)
}
