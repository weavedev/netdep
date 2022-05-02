package stage

import (
	"github.com/stretchr/testify/assert"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stage"
	"testing"
)

/**
A test for the sample implementation of the resolution method
*/
func TestResolving(t *testing.T) {
	var res = stage.ResolveEnvVars("test")

	var expected = map[string]map[string]interface{}{
		"testSampleService": {
			"VariableNameA": "1",
			"VariableNameB": "False",
		},
	}

	assert.Equal(t, expected, res, "Expected the resolution method to return the dummy map")
}
