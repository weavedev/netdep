// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

// Tests adjacency list construction based on data found in discovery stage
func TestConstructAdjacencyList(t *testing.T) {
	servCalls := []stages.ServiceCalls{
		{
			Service: "servA",
			Calls: map[string][]stages.CallData{
				"URL_2": {
					{Filepath: "./path/to/some/file.go", Line: 24},
					{Filepath: "/path/to/some/otherfile.go", Line: 36},
				},
			},
		},
		{
			Service: "servB",
			Calls: map[string][]stages.CallData{
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
			Calls: map[string][]stages.CallData{
				"URL_5": {
					{Filepath: "./path5/to/some/file.go", Line: 215},
				},
			},
		},
	}
	data := &stages.DiscoveredData{
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

	servCalls := []stages.ServiceCalls{
		{
			Service: "servA",
			Calls: map[string][]stages.CallData{
				"URL_2": {
					{Filepath: "./path/to/some/file.go", Line: 24},
					{Filepath: "/path/to/some/otherfile.go", Line: 36},
				},
			},
		},
		{
			Service: "servB",
			Calls: map[string][]stages.CallData{
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
			Calls: map[string][]stages.CallData{
				"URL_5": {
					{Filepath: "./path5/to/some/file.go", Line: 215},
				},
			},
		},
	}

	adjList, callData := stages.SerialiseOutput(m, servCalls)

	expectedAdjList := `{"servA":[{"service":"servB","amount":2}],"servB":[{"service":"servA","amount":1},{"service":"servC","amount":3}],"servC":[{"service":"Unknown Service","amount":1}]}`
	expectedCallData := `[{"service":"servA","calls":{"URL_2":[{"filepath":"./path/to/some/file.go","line":24},{"filepath":"/path/to/some/otherfile.go","line":36}]}},{"service":"servB","calls":{"URL_1":[{"filepath":"./path/to/some/cool/file.go","line":632}],"URL_3":[{"filepath":"./path1/to/some/file.go","line":245},{"filepath":"/path1/to/some/file.go","line":436},{"filepath":"/path1/to/some/file.go","line":623}]}},{"service":"servC","calls":{"URL_5":[{"filepath":"./path5/to/some/file.go","line":215}]}}]`
	assert.Equal(t, expectedAdjList, adjList)
	assert.Equal(t, expectedCallData, callData)
}
