// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

import (
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

// Tests adjacency list construction based on data found in discovery stage
func TestConstructAdjacencyList(t *testing.T) {
	servCalls := []discovery.ServiceCalls{
		{
			Service: "servA",
			Calls: map[string][]discovery.CallData{
				"URL_2": {
					{Filepath: "./path/to/some/file.go", Line: 24},
					{Filepath: "/path/to/some/otherfile.go", Line: 36},
				},
			},
		},
		{
			Service: "servB",
			Calls: map[string][]discovery.CallData{
				"URL_1": {
					{Filepath: "./path/to/some/cool/file.go", Line: 632},
				},
				"URL_3": {
					{Filepath: "./path1/to/some/file.go", Line: 245},
					{Filepath: "/path1/to/some/file.go", Line: 436},
					{Filepath: "/path1/to/some/file.go", Line: 623},
				},
			},
		},
		{
			Service: "servC",
			Calls: map[string][]discovery.CallData{
				"URL_5": {
					{Filepath: "./path5/to/some/file.go", Line: 215},
				},
			},
		},
	}
	data := &discovery.DiscoveredData{
		ServCalls: servCalls,
		Handled:   map[string]string{"URL_1": "servA", "URL_2": "servB", "URL_3": "servC"},
	}

	res := stages.ConstructAdjacencyList(data)
	expected := map[string][]stages.Conn{
		"servA": {
			{Service: "servB", Amount: 2},
		},
		"servB": {
			{Service: "servA", Amount: 1},
			{Service: "servC", Amount: 3},
		},
		"servC": {
			{Service: "Unknown Service", Amount: 1},
		},
	}

	assert.Equal(t, expected, res)
}

/*
	json package and its serialisation are already sufficiently tested by its developers.
	This unit test is essentially a sanity check to ensure the output is in the format we expect.
*/
func TestSerialiseOutput(t *testing.T) {
	m := make(map[string][]stages.Conn)
	m["servA"] = append(m["servA"], stages.Conn{Service: "servB", Amount: 2})
	m["servB"] = append(m["servB"], stages.Conn{Service: "servA", Amount: 1})
	m["servB"] = append(m["servB"], stages.Conn{Service: "servC", Amount: 3})
	m["servC"] = append(m["servC"], stages.Conn{Service: "Unknown Service", Amount: 1})

	servCalls := []discovery.ServiceCalls{
		{
			Service: "servA",
			Calls: map[string][]discovery.CallData{
				"URL_2": {
					{Filepath: "./path/to/some/file.go", Line: 24},
					{Filepath: "/path/to/some/otherfile.go", Line: 36},
				},
			},
		},
		{
			Service: "servB",
			Calls: map[string][]discovery.CallData{
				"URL_1": {
					{Filepath: "./path/to/some/cool/file.go", Line: 632},
				},
				"URL_3": {
					{Filepath: "./path1/to/some/file.go", Line: 245},
					{Filepath: "/path1/to/some/file.go", Line: 436},
					{Filepath: "/path1/to/some/file.go", Line: 623},
				},
			},
		},
		{
			Service: "servC",
			Calls: map[string][]discovery.CallData{
				"URL_5": {
					{Filepath: "./path5/to/some/file.go", Line: 215},
				},
			},
		},
	}

	adjList, callData := stages.SerialiseOutput(m, servCalls)

	expectedAdjList := "{\n\t\"servA\": [\n\t\t{\n\t\t\t\"service\": \"servB\",\n\t\t\t\"amount\": 2\n\t\t}\n\t],\n\t\"servB\": [\n\t\t{\n\t\t\t\"service\": \"servA\",\n\t\t\t\"amount\": 1\n\t\t},\n\t\t{\n\t\t\t\"service\": \"servC\",\n\t\t\t\"amount\": 3\n\t\t}\n\t],\n\t\"servC\": [\n\t\t{\n\t\t\t\"service\": \"Unknown Service\",\n\t\t\t\"amount\": 1\n\t\t}\n\t]\n}"
	expectedCallData := "[\n\t{\n\t\t\"service\": \"servA\",\n\t\t\"calls\": {\n\t\t\t\"URL_2\": [\n\t\t\t\t{\n\t\t\t\t\t\"filepath\": \"./path/to/some/file.go\",\n\t\t\t\t\t\"line\": 24\n\t\t\t\t},\n\t\t\t\t{\n\t\t\t\t\t\"filepath\": \"/path/to/some/otherfile.go\",\n\t\t\t\t\t\"line\": 36\n\t\t\t\t}\n\t\t\t]\n\t\t}\n\t},\n\t{\n\t\t\"service\": \"servB\",\n\t\t\"calls\": {\n\t\t\t\"URL_1\": [\n\t\t\t\t{\n\t\t\t\t\t\"filepath\": \"./path/to/some/cool/file.go\",\n\t\t\t\t\t\"line\": 632\n\t\t\t\t}\n\t\t\t],\n\t\t\t\"URL_3\": [\n\t\t\t\t{\n\t\t\t\t\t\"filepath\": \"./path1/to/some/file.go\",\n\t\t\t\t\t\"line\": 245\n\t\t\t\t},\n\t\t\t\t{\n\t\t\t\t\t\"filepath\": \"/path1/to/some/file.go\",\n\t\t\t\t\t\"line\": 436\n\t\t\t\t},\n\t\t\t\t{\n\t\t\t\t\t\"filepath\": \"/path1/to/some/file.go\",\n\t\t\t\t\t\"line\": 623\n\t\t\t\t}\n\t\t\t]\n\t\t}\n\t},\n\t{\n\t\t\"service\": \"servC\",\n\t\t\"calls\": {\n\t\t\t\"URL_5\": [\n\t\t\t\t{\n\t\t\t\t\t\"filepath\": \"./path5/to/some/file.go\",\n\t\t\t\t\t\"line\": 215\n\t\t\t\t}\n\t\t\t]\n\t\t}\n\t}\n]"
	assert.Equal(t, expectedAdjList, adjList)
	assert.Equal(t, expectedCallData, callData)
}
