package stages

import (
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/output"

	"github.com/stretchr/testify/assert"
)

// test basic functionality of the matching functionality
func TestEmptyCaseCreateDependencyGraph(t *testing.T) {
	discoveredData := &discovery.DiscoveredData{
		ServCalls: nil,
		Handled:   nil,
	}

	nodes, edges := stages.CreateDependencyGraph(discoveredData)

	assert.Equal(t, make([]*output.ServiceNode, 0), nodes)
	assert.Equal(t, make([]*output.ConnectionEdge, 0), edges)
}

// test basic functionality of the matching functionality
func TestBasicCreateDependencyGraph(t *testing.T) {
	discoveredData := &discovery.DiscoveredData{
		ServCalls: []discovery.ServiceCalls{
			{
				Service: "Node1",
				Calls: map[string][]discovery.CallData{
					"URL_2": {
						{Filepath: "./node1/path/to/some/file.go", Line: 24},
					},
					"URL_3": {
						{Filepath: "./node1/path/to/some/other/file.go", Line: 36},
					},
				},
			},
			{
				Service: "Node2",
				Calls: map[string][]discovery.CallData{
					"URL_3": {
						{Filepath: "./node1/path/to/some/file.go", Line: 245},
						{Filepath: "./node2/path/to/some/other/file.go", Line: 436},
					},
				},
			},
			{
				Service: "Node3",
				Calls:   map[string][]discovery.CallData{},
			},
		},
		Handled: map[string]string{"URL_1": "Node1", "URL_2": "Node2", "URL_3": "Node3"},
	}

	// reuse graph from output stage tests
	expectedNodes, expectedEdges := CreateSmallTestGraph()

	nodes, edges := stages.CreateDependencyGraph(discoveredData)

	assert.Equal(t, expectedNodes, nodes)
	assert.Equal(t, expectedEdges, edges)
}

// test functionality with unknown service
func TestWithUnknownService(t *testing.T) {
	// input data
	discoveredData := &discovery.DiscoveredData{
		ServCalls: []discovery.ServiceCalls{
			{
				Service: "Node1",
				Calls: map[string][]discovery.CallData{
					"URL_1": {
						{Filepath: "./node1/path/to/some/file.go", Line: 24},
						{Filepath: "./node1/path/to/some/other/file.go", Line: 436},
					},
					"URL_2": {
						{Filepath: "./node1/path/to/some/other_file.go", Line: 24},
					},
				},
			},
		},
		Handled: map[string]string{},
	}

	// output data
	node1 := output.ServiceNode{
		ServiceName: "Node1",
	}

	node2 := output.ServiceNode{
		ServiceName: "Unknown Service #1",
	}

	node3 := output.ServiceNode{
		ServiceName: "Unknown Service #2",
	}

	edge12a := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "URL_1",
			Arguments: nil,
			Location:  discovery.CallData{Filepath: "./node1/path/to/some/file.go", Line: 24},
		},
		Source: &node1,
		Target: &node2,
	}

	edge12b := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "URL_1",
			Arguments: nil,
			Location:  discovery.CallData{Filepath: "./node1/path/to/some/other/file.go", Line: 436},
		},
		Source: &node1,
		Target: &node2,
	}

	edge13 := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:  "HTTP",
			URL:       "URL_2",
			Arguments: nil,
			Location:  discovery.CallData{Filepath: "./node1/path/to/some/other_file.go", Line: 24},
		},
		Source: &node1,
		Target: &node3,
	}

	expectedNodes := []*output.ServiceNode{&node1, &node2, &node3}
	expectedEdges := []*output.ConnectionEdge{&edge12a, &edge12b, &edge13}

	nodes, edges := stages.CreateDependencyGraph(discoveredData)

	assert.Equal(t, expectedNodes, nodes)
	assert.Equal(t, expectedEdges, edges)
}
