// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

import (
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/output"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"

	"github.com/stretchr/testify/assert"
)

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
		Connection: output.Connection{
			Protocol:  "HTTP",
			Url:       "",
			Arguments: nil,
			Location:  discovery.CallData{},
		},
		Source: &node1,
		Target: &node2,
	}

	edge13 := output.ConnectionEdge{
		Connection: output.Connection{
			Protocol:  "HTTP",
			Url:       "",
			Arguments: nil,
			Location:  discovery.CallData{},
		},
		Source: &node1,
		Target: &node3,
	}

	nodes := []*output.ServiceNode{&node1, &node2, &node3}
	edges := []*output.ConnectionEdge{&edge12, &edge13}
	return nodes, edges
}

// Tests adjacency list construction based on data found in discovery stage
func TestConstructAdjacencyList(t *testing.T) {
	nodes, edges := createSmallTestGraph()
	res := output.ConstructAdjacencyList(nodes, edges)

	expected := map[string][]output.ServiceConnection{
		"Node1": {
			{
				Service:     *nodes[1],
				Connection:  []output.Connection{edges[0].Connection},
				Connections: 1,
			},
			{
				Service:     *nodes[2],
				Connection:  []output.Connection{edges[1].Connection},
				Connections: 1,
			},
		},
		"Node2": {},
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
	expected := "{\"Node1\":[{\"service\":{\"ServiceName\":\"Node2\"},\"connections\":[{\"protocol\":\"HTTP\",\"url\":\"\",\"arguments\":null,\"location\":{\"filepath\":\"\",\"line\":0}}],\"count\":1},{\"service\":{\"ServiceName\":\"Node3\"},\"connections\":[{\"protocol\":\"HTTP\",\"url\":\"\",\"arguments\":null,\"location\":{\"filepath\":\"\",\"line\":0}}],\"count\":1}],\"Node2\":[],\"Node3\":[]}"
	assert.Equal(t, expected, str)
}
