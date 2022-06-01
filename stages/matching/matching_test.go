package matching

import (
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/output"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"

	"github.com/stretchr/testify/assert"
)

func CreateSmallTestGraph() output.NodeGraph {
	node1 := output.ServiceNode{
		ServiceName: "Node1",
		IsUnknown:   false,
	}
	node2 := output.ServiceNode{
		ServiceName: "Node2",
		IsUnknown:   false,
	}
	node3 := output.ServiceNode{
		ServiceName: "Node3",
		IsUnknown:   false,
	}

	edge12 := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node2:80/URL_2",
			Arguments: nil,
			Location:  "./node1/path/to/some/file.go:24",
		},
		Source: &node1,
		Target: &node2,
	}

	edge13 := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node3:80/URL_3",
			Arguments: nil,
			Location:  "./node1/path/to/some/other/file.go:36",
		},
		Source: &node1,
		Target: &node3,
	}

	edge23a := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node3:80/URL_3",
			Arguments: nil,
			Location:  "./node1/path/to/some/file.go:245",
		},
		Source: &node2,
		Target: &node3,
	}

	edge23b := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node3:80/URL_3",
			Arguments: nil,
			Location:  "./node2/path/to/some/other/file.go:436",
		},
		Source: &node2,
		Target: &node3,
	}

	nodes := []*output.ServiceNode{&node1, &node2, &node3}
	edges := []*output.ConnectionEdge{&edge12, &edge13, &edge23a, &edge23b}

	return output.NodeGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

// test basic functionality of the matching functionality
func TestEmptyCaseCreateDependencyGraph(t *testing.T) {
	graph := CreateDependencyGraph(nil, nil)

	assert.Equal(t, make([]*output.ServiceNode, 0), graph.Nodes)
	assert.Equal(t, make([]*output.ConnectionEdge, 0), graph.Edges)
}

// test basic functionality of the matching functionality
func TestBasicCreateDependencyGraph(t *testing.T) {
	calls := []*callanalyzer.CallTarget{
		{
			RequestLocation: "http://Node2:80/URL_2",
			ServiceName:     "Node1",
			FileName:        "./node1/path/to/some/file.go",
			PositionInFile:  "24",
			IsResolved:      true,
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node1",
			FileName:        "./node1/path/to/some/other/file.go",
			PositionInFile:  "36",
			IsResolved:      true,
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node2",
			FileName:        "./node1/path/to/some/file.go",
			PositionInFile:  "245",
			IsResolved:      true,
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node2",
			FileName:        "./node2/path/to/some/other/file.go",
			PositionInFile:  "436",
			IsResolved:      true,
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
	expectedGraph := CreateSmallTestGraph()

	graph := CreateDependencyGraph(calls, endpoints)

	assert.Equal(t, len(expectedGraph.Nodes), len(graph.Nodes))
	for i := range expectedGraph.Nodes {
		assert.Equal(t, expectedGraph.Nodes[i], graph.Nodes[i])
	}

	assert.Equal(t, len(expectedGraph.Edges), len(graph.Edges))
	for i := range expectedGraph.Edges {
		assert.Equal(t, expectedGraph.Edges[i], graph.Edges[i])
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
			IsResolved:      true,
		},
		{
			RequestLocation: "http://Node2:80/URL_2",
			ServiceName:     "Node1",
			FileName:        "./node1/path/to/some/other/file.go",
			PositionInFile:  "436",
			IsResolved:      true,
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node1",
			FileName:        "./node1/path/to/some/other_file.go",
			PositionInFile:  "24",
			IsResolved:      true,
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

	unknownService := output.ServiceNode{
		ServiceName: "UnknownService",
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
		Target: &unknownService,
	}

	edge12b := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node2:80/URL_2",
			Arguments: nil,
			Location:  "./node1/path/to/some/other/file.go:436",
		},
		Source: &node1,
		Target: &unknownService,
	}

	edge13 := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "http://Node3:80/URL_3",
			Arguments: nil,
			Location:  "./node1/path/to/some/other_file.go:24",
		},
		Source: &node1,
		Target: &unknownService,
	}

	expectedNodes := []*output.ServiceNode{&node1, &unknownService}
	expectedEdges := []*output.ConnectionEdge{&edge12a, &edge12b, &edge13}

	graph := CreateDependencyGraph(calls, endpoints)

	assert.Equal(t, len(expectedNodes), len(graph.Nodes))
	for i := range expectedNodes {
		assert.Equal(t, expectedNodes[i], graph.Nodes[i])
	}

	assert.Equal(t, len(expectedEdges), len(graph.Edges))
	for i := range expectedEdges {
		assert.Equal(t, expectedEdges[i], graph.Edges[i])
	}
}
