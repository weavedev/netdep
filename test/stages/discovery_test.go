package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

/*
A test for the sample implementation of the resolution method
*/
func TestDiscovery(t *testing.T) {
	var res = stages.FindCallersForEndpoint("testService", "testEndpoint", "https://example.com")
	assert.Equal(t, []interface{}(nil), res, "Expect nil as the output of the sample FindCallersForEndpoint method")
}
