// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

import (
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

func createSmallTestGraph() ([]*stages.ServiceNode, []*stages.ConnectionEdge) {

	node1 := stages.ServiceNode{
		ServiceName: "Node1",
	}
	node2 := stages.ServiceNode{
		ServiceName: "Node2",
	}
	node3 := stages.ServiceNode{
		ServiceName: "Node3",
	}

	edge12 := stages.ConnectionEdge{
		Connection: stages.Connection{
			Protocol:  "HTTP",
			Url:       "",
			Arguments: nil,
			Location:  discovery.CallData{},
		},
		Source: &node1,
		Target: &node2,
	}

	edge13 := stages.ConnectionEdge{
		Connection: stages.Connection{
			Protocol:  "HTTP",
			Url:       "",
			Arguments: nil,
			Location:  discovery.CallData{},
		},
		Source: &node1,
		Target: &node3,
	}

	nodes := []*stages.ServiceNode{&node1, &node2, &node3}
	edges := []*stages.ConnectionEdge{&edge12, &edge13}
	return nodes, edges
}

// Tests adjacency list construction based on data found in discovery stage
func TestConstructAdjacencyList(t *testing.T) {
	nodes, edges := createSmallTestGraph()
	res := stages.ConstructAdjacencyList(nodes, edges)

	expected := map[string][]stages.ServiceConnection{
		"Node1": {
			{
				Service:     *nodes[1],
				Connection:  []stages.Connection{edges[0].Connection},
				Connections: 1,
			},
			{
				Service:     *nodes[2],
				Connection:  []stages.Connection{edges[1].Connection},
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
	str, _ := stages.SerializeAdjacencyList(nil, false)
	assert.Equal(t, "null", str)
}

// TestSerialiseOutput test a realistic output of a serialisation. More of a sanity check.
func TestSerialiseOutput(t *testing.T) {
	nodes, edges := createSmallTestGraph()
	list := stages.ConstructAdjacencyList(nodes, edges)
	str, _ := stages.SerializeAdjacencyList(list, false)
	expected := "{\"Node1\":[{\"service\":{\"ServiceName\":\"Node2\"},\"connections\":[{\"protocol\":\"HTTP\",\"url\":\"\",\"arguments\":null,\"location\":{\"filepath\":\"\",\"line\":0}}],\"count\":1},{\"service\":{\"ServiceName\":\"Node3\"},\"connections\":[{\"protocol\":\"HTTP\",\"url\":\"\",\"arguments\":null,\"location\":{\"filepath\":\"\",\"line\":0}}],\"count\":1}],\"Node2\":[],\"Node3\":[]}"
	assert.Equal(t, expected, str)
}
