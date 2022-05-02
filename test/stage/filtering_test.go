package stage

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
	"testing"
)

/**
A test for the sample implementation of the resolution method
*/
func TestFiltering(t *testing.T) {
	var res = stages.ScanAndFilter("test")
	assert.Equal(t, map[string][]*ast.File(nil), res, "Expect nil as the output of the ScanAndFilter method")
}
