package stages

import (
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/output"

	"github.com/stretchr/testify/assert"
)

// test basic functionality of the matching functionality
func TestEmptyCaseCreateDependencyGraph(t *testing.T) {
	nodes, edges := stages.CreateDependencyGraph(nil, nil)

	assert.Equal(t, make([]*output.ServiceNode, 0), nodes)
	assert.Equal(t, make([]*output.ConnectionEdge, 0), edges)
}

// test basic functionality of the matching functionality
func TestBasicCreateDependencyGraph(t *testing.T) {
	calls := []*callanalyzer.CallTarget{
		{
			RequestLocation: "http://Node2:80/URL_2",
			ServiceName:     "Node1",
			FileName:        "./node1/path/to/some/file.go",
			PositionInFile:  "24",
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node1",
			FileName:        "./node1/path/to/some/other/file.go",
			PositionInFile:  "36",
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node2",
			FileName:        "./node1/path/to/some/file.go",
			PositionInFile:  "245",
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node2",
			FileName:        "./node2/path/to/some/other/file.go",
			PositionInFile:  "436",
		},
	}

	endpoints := []*callanalyzer.CallTarget{
		{
			RequestLocation: "/URL_1",
			ServiceName:     "Node1",
		},
		{
			RequestLocation: "/URL_2",
			ServiceName:     "Node2",
		},
		{
			RequestLocation: "/URL_3",
			ServiceName:     "Node3",
		},
	}

	// reuse graph from output stage tests
	expectedNodes, expectedEdges := CreateSmallTestGraph()

	nodes, edges := stages.CreateDependencyGraph(calls, endpoints)

	assert.Equal(t, len(expectedNodes), len(nodes))
	for i := range expectedNodes {
		assert.Equal(t, expectedNodes[i], nodes[i])
	}

	assert.Equal(t, len(expectedEdges), len(edges))
	for i := range expectedEdges {
		assert.Equal(t, expectedEdges[i], edges[i])
	}
}

// test functionality with unknown service
func TestWithUnknownService(t *testing.T) {
	// input data
	calls := []*callanalyzer.CallTarget{
		{
			RequestLocation: "http://Node2:80/URL_2",
			ServiceName:     "Node1",
			FileName:        "./node1/path/to/some/file.go",
			PositionInFile:  "24",
		},
		{
			RequestLocation: "http://Node2:80/URL_2",
			ServiceName:     "Node1",
			FileName:        "./node1/path/to/some/other/file.go",
			PositionInFile:  "436",
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node1",
			FileName:        "./node1/path/to/some/other_file.go",
			PositionInFile:  "24",
		},
	}

	endpoints := []*callanalyzer.CallTarget{
		{
			RequestLocation: "/URL_1",
			ServiceName:     "Node1",
		},
	}

	// output data
	node1 := output.ServiceNode{
		ServiceName: "Node1",
		IsUnknown:   false,
	}

	node2 := output.ServiceNode{
		ServiceName: "Unknown Service #1",
		IsUnknown:   true,
	}

	node3 := output.ServiceNode{
		ServiceName: "Unknown Service #2",
		IsUnknown:   true,
	}

	edge12a := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node2:80/URL_2",
			Arguments: nil,
			Location:  "./node1/path/to/some/file.go:24",
		},
		Source: &node1,
		Target: &node2,
	}

	edge12b := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node2:80/URL_2",
			Arguments: nil,
			Location:  "./node1/path/to/some/other/file.go:436",
		},
		Source: &node1,
		Target: &node2,
	}

	edge13 := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node3:80/URL_3",
			Arguments: nil,
			Location:  "./node1/path/to/some/other_file.go:24",
		},
		Source: &node1,
		Target: &node3,
	}

	expectedNodes := []*output.ServiceNode{&node1, &node2, &node3}
	expectedEdges := []*output.ConnectionEdge{&edge12a, &edge12b, &edge13}

	nodes, edges := stages.CreateDependencyGraph(calls, endpoints)

	assert.Equal(t, len(expectedNodes), len(nodes))
	for i := range expectedNodes {
		assert.Equal(t, expectedNodes[i], nodes[i])
	}

	assert.Equal(t, len(expectedEdges), len(edges))
	for i := range expectedEdges {
		assert.Equal(t, expectedEdges[i], edges[i])
	}
}
