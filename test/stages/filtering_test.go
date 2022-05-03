package stages

import (
	"go/ast"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"

	"github.com/stretchr/testify/assert"
)

/*
A test for the sample implementation of the resolution method
*/
func TestFiltering(t *testing.T) {
	var res = stages.ScanAndFilter("test")
	assert.Equal(t, map[string][]*ast.File(nil), res, "Expect nil as the output of the ScanAndFilter method")
}
