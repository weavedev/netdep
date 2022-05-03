package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/tools/go/callgraph"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

/*
A test for the sample implementation of the resolution method
*/
func TestParsing(t *testing.T) {
	var res = stages.CreateCallGraph(nil)
	assert.Equal(t, callgraph.Graph{}, res, "Expect empty graph as the output of the CreateCallGraph method")
}
