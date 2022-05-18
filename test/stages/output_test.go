// Package output
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

import (
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/output"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"

	"github.com/stretchr/testify/assert"
)

// createSmallTestGraph creates a graph with three nodes, where node 1 had edges to node 2 and 3, and node 2 to node 3
func createSmallTestGraph() ([]*output.ServiceNode, []*output.ConnectionEdge) {
	node1 := output.ServiceNode{
		ServiceName: "Node1",
	}
	node2 := output.ServiceNode{
		ServiceName: "Node2",
	}
	node3 := output.ServiceNode{
		ServiceName: "Node3",
	}

	edge12 := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "",
			Arguments: nil,
			Location:  discovery.CallData{},
		},
		Source: &node1,
		Target: &node2,
	}

	edge13 := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "",
			Arguments: nil,
			Location:  discovery.CallData{},
		},
		Source: &node1,
		Target: &node3,
	}

	edge23 := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "",
			Arguments: nil,
			Location:  discovery.CallData{},
		},
		Source: &node2,
		Target: &node3,
	}

	nodes := []*output.ServiceNode{&node1, &node2, &node3}
	edges := []*output.ConnectionEdge{&edge12, &edge13, &edge23}
	return nodes, edges
}

// Tests adjacency list construction based on data found in discovery stage
func TestConstructAdjacencyList(t *testing.T) {
	nodes, edges := createSmallTestGraph()
	res := output.ConstructAdjacencyList(nodes, edges)

	expected := output.AdjacencyList{
		"Node1": {
			{
				// node 2
				Service:       *nodes[1],
				Calls:         []output.NetworkCall{edges[0].Call},
				NumberOfCalls: 1,
			},
			{
				// node 3
				Service:       *nodes[2],
				Calls:         []output.NetworkCall{edges[1].Call},
				NumberOfCalls: 1,
			},
		},
		"Node2": {
			{
				// node 3
				Service:       *nodes[2],
				Calls:         []output.NetworkCall{edges[2].Call},
				NumberOfCalls: 1,
			},
		},
		"Node3": {},
	}

	assert.Equal(t, expected, res)
}

// TestSerialiseOutputNull performs a sanity check for nil case
func TestSerialiseOutputNull(t *testing.T) {
	str, _ := output.SerializeAdjacencyList(nil, false)
	assert.Equal(t, "null", str)
}

// TestSerialiseOutput test a realistic output of a serialisation. More of a sanity check.
func TestSerialiseOutput(t *testing.T) {
	nodes, edges := createSmallTestGraph()
	list := output.ConstructAdjacencyList(nodes, edges)
	str, _ := output.SerializeAdjacencyList(list, false)
	expected := "{\"Node1\":[{\"service\":{\"serviceName\":\"Node2\"},\"calls\":[{\"protocol\":\"HTTP\",\"url\":\"\",\"arguments\":null,\"location\":{\"filepath\":\"\",\"line\":0}}],\"count\":1},{\"service\":{\"serviceName\":\"Node3\"},\"calls\":[{\"protocol\":\"HTTP\",\"url\":\"\",\"arguments\":null,\"location\":{\"filepath\":\"\",\"line\":0}}],\"count\":1}],\"Node2\":[{\"service\":{\"serviceName\":\"Node3\"},\"calls\":[{\"protocol\":\"HTTP\",\"url\":\"\",\"arguments\":null,\"location\":{\"filepath\":\"\",\"line\":0}}],\"count\":1}],\"Node3\":[]}"
	assert.Equal(t, expected, str)
}
