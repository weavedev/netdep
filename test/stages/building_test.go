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
			Calls:   map[string]int{"URL_2": 2},
		},
		{
			Service: "servB",
			Calls:   map[string]int{"URL_1": 1, "URL_3": 3},
		},
		{
			Service: "servC",
			Calls:   map[string]int{"URL_5": 1},
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
	m["servC"] = append(m["servC"], stages.Conn{Service: "Unknown Service", Amount: 2})
	res := stages.SerialiseOutput(m)
	expected := `{"servA":[{"service":"servB","amount":2}],"servB":[{"service":"servA","amount":1},{"service":"servC","amount":3}],"servC":[{"service":"Unknown Service","amount":2}]}`
	assert.Equal(t, expected, res)
}
