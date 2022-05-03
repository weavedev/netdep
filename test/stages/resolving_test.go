package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

/*
A test for the sample implementation of the resolution method
*/
func TestResolving(t *testing.T) {
	var res = stages.ResolveEnvVars("test")

	var expected = map[string]map[string]interface{}{
		"testSampleService": {
			"VariableNameA": "1",
			"VariableNameB": "False",
		},
	}

	assert.Equal(t, expected, res, "Expected the resolution method to return the dummy map")
}
