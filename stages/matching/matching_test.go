package matching

import (
	"testing"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/natsanalyzer"
	"lab.weave.nl/internships/tud-2022/netDep/structures"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"
	"lab.weave.nl/internships/tud-2022/netDep/stages/output"
)

func CreateSmallTestGraph() output.NodeGraph {
	node1 := output.ServiceNode{
		ServiceName:   "Node1",
		IsUnknown:     false,
		IsReferenced:  false,
		IsReferencing: true,
	}
	node2 := output.ServiceNode{
		ServiceName:   "Node2",
		IsUnknown:     false,
		IsReferenced:  true,
		IsReferencing: true,
	}
	node3 := output.ServiceNode{
		ServiceName:   "Node3",
		IsUnknown:     false,
		IsReferenced:  true,
		IsReferencing: false,
	}
	node4 := output.ServiceNode{
		ServiceName:   "Node4",
		IsUnknown:     false,
		IsReferenced:  true,
		IsReferencing: false,
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
	edge24 := output.ConnectionEdge{
		Call: output.NetworkCall{
			Protocol:   "servicecalls",
			MethodName: "SomeMethod",
			URL:        "",
			Arguments:  nil,
			Location:   "./node2/path/to/some/other/file1.go:42",
		},
		Source: &node2,
		Target: &node4,
	}

	nodes := []*output.ServiceNode{&node1, &node2, &node3, &node4}
	edges := []*output.ConnectionEdge{&edge12, &edge13, &edge23a, &edge23b, &edge24}

	return output.NodeGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

// test basic functionality of the matching functionality
func TestEmptyCaseCreateDependencyGraph(t *testing.T) {
	graph := CreateDependencyGraph(nil)

	assert.Equal(t, make([]*output.ServiceNode, 0), graph.Nodes)
	assert.Equal(t, make([]*output.ConnectionEdge, 0), graph.Edges)
}

// test basic functionality of the matching functionality
func TestBasicCreateDependencyGraph(t *testing.T) {
	calls := []*callanalyzer.CallTarget{
		{
			RequestLocation: "http://Node2:80/URL_2",
			ServiceName:     "Node1",
			Trace: []callanalyzer.CallTargetTrace{
				{
					FileName:       "./node1/path/to/some/file.go",
					PositionInFile: "24",
				},
			},
			IsResolved: true,
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node1",
			Trace: []callanalyzer.CallTargetTrace{
				{
					FileName:       "./node1/path/to/some/other/file.go",
					PositionInFile: "36",
				},
			},
			IsResolved: true,
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node2",
			Trace: []callanalyzer.CallTargetTrace{
				{
					FileName:       "./node1/path/to/some/file.go",
					PositionInFile: "245",
				},
			},
			IsResolved: true,
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node2",
			Trace: []callanalyzer.CallTargetTrace{
				{
					FileName:       "./node2/path/to/some/other/file.go",
					PositionInFile: "436",
				},
			},
			IsResolved: true,
		},
		{
			PackageName:     "servicecalls",
			MethodName:      "SomeMethod",
			RequestLocation: "SomeMethod",
			ServiceName:     "Node2",
			Trace: []callanalyzer.CallTargetTrace{
				{
					FileName:       "./node2/path/to/some/other/file1.go",
					PositionInFile: "42",
				},
			},
			IsResolved: true,
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
		{
			PackageName:     "servicecalls",
			MethodName:      "SomeMethod",
			RequestLocation: "SomeMethod",
			ServiceName:     "Node4",
		},
	}

	// reuse graph from output stage tests
	expectedGraph := CreateSmallTestGraph()

	dependencies := &structures.Dependencies{
		Calls:     calls,
		Endpoints: endpoints,
		Consumers: nil,
		Producers: nil,
	}

	graph := CreateDependencyGraph(dependencies)

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
			Trace: []callanalyzer.CallTargetTrace{
				{
					FileName:       "./node1/path/to/some/file.go",
					PositionInFile: "24",
				},
			},
			IsResolved: true,
		},
		{
			RequestLocation: "http://Node2:80/URL_2",
			ServiceName:     "Node1",
			Trace: []callanalyzer.CallTargetTrace{
				{
					FileName:       "./node1/path/to/some/other/file.go",
					PositionInFile: "436",
				},
			},
			IsResolved: true,
		},
		{
			RequestLocation: "http://Node3:80/URL_3",
			ServiceName:     "Node1",
			Trace: []callanalyzer.CallTargetTrace{
				{
					FileName:       "./node1/path/to/some/other_file.go",
					PositionInFile: "24",
				},
			},
			IsResolved: true,
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
		ServiceName:   "Node1",
		IsUnknown:     false,
		IsReferenced:  false,
		IsReferencing: true,
	}

	unknownService := output.ServiceNode{
		ServiceName:   "UnknownService",
		IsUnknown:     true,
		IsReferenced:  true,
		IsReferencing: true,
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

	dependencies := &structures.Dependencies{
		Calls:     calls,
		Endpoints: endpoints,
		Consumers: nil,
		Producers: nil,
	}

	graph := CreateDependencyGraph(dependencies)

	assert.Equal(t, len(expectedNodes), len(graph.Nodes))
	for i := range expectedNodes {
		assert.Equal(t, expectedNodes[i], graph.Nodes[i])
	}

	assert.Equal(t, len(expectedEdges), len(graph.Edges))
	for i := range expectedEdges {
		assert.Equal(t, expectedEdges[i], graph.Edges[i])
	}
}

func TestNatsExtension(t *testing.T) {
	call1 := &natsanalyzer.NatsCall{
		Communication:  "NATS",
		MethodName:     "Subscribe",
		Subject:        "HelloSubject",
		ServiceName:    "test",
		FileName:       "test.go",
		PositionInFile: "15",
	}

	call2 := &natsanalyzer.NatsCall{
		Communication:  "NATS",
		MethodName:     "Subscribe",
		Subject:        "ByeSubject",
		ServiceName:    "test",
		FileName:       "test.go",
		PositionInFile: "16",
	}

	call3 := &natsanalyzer.NatsCall{
		Communication:  "NATS",
		MethodName:     "ByeNotifyMsg",
		Subject:        "ByeSubject",
		ServiceName:    "test",
		FileName:       "testNotify.go",
		PositionInFile: "18",
	}

	call4 := &natsanalyzer.NatsCall{
		Communication:  "NATS",
		MethodName:     "HelloNotifyMsg",
		Subject:        "HelloSubject",
		ServiceName:    "test",
		FileName:       "testNotify.go",
		PositionInFile: "17",
	}

	call5 := &natsanalyzer.NatsCall{
		Communication:  "NATS",
		MethodName:     "AyoNotifyMsg",
		Subject:        "AyoSubject",
		ServiceName:    "ayo",
		FileName:       "ayo.go",
		PositionInFile: "1",
	}

	node1 := output.ServiceNode{
		ServiceName: "ayo",
		IsUnknown:   false,
	}
	node2 := output.ServiceNode{
		ServiceName: "test",
		IsUnknown:   false,
	}

	serviceMap := make(map[string]*output.ServiceNode)
	serviceMap["ayo"] = &node1
	serviceMap["test"] = &node2

	consumers := make([]*natsanalyzer.NatsCall, 2)
	consumers[0] = call1
	consumers[1] = call2

	producers := make([]*natsanalyzer.NatsCall, 3)
	producers[0] = call3
	producers[1] = call4
	producers[2] = call5

	hasUnknown := false
	nodes := make([]*output.ServiceNode, 0)

	edges := extendWithNats(consumers, producers, &hasUnknown, serviceMap, &nodes)
	assert.Equal(t, *edges[0].Source, node2)
	assert.Equal(t, *edges[0].Target, node2)
	assert.Equal(t, edges[0].Call.URL, "ByeSubject")
	assert.Equal(t, *edges[1].Source, node2)
	assert.Equal(t, *edges[1].Target, node2)
	assert.Equal(t, edges[1].Call.URL, "HelloSubject")
	assert.Equal(t, edges[2].Call.URL, "AyoSubject")
	assert.Equal(t, edges[2].Target.ServiceName, "UnknownService")
}
