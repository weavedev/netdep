package stage

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/callgraph"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stage"
	"testing"
)

/**
A test for the sample implementation of the resolution method
*/
func TestDiscovery(t *testing.T) {
	var res = stage.FindCallersForEndpoint("testService", "testEndpoint", "https://example.com")
	assert.Equal(t, callgraph.Graph{}, res, "Expect empty graph as the output of the sample CreateCallGraph method")
}
