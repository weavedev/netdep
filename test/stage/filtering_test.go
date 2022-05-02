package stage

import (
	"github.com/stretchr/testify/assert"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stage"
	"testing"
)

/**
A test for the sample implementation of the resolution method
*/
func TestFiltering(t *testing.T) {
	var res = stage.ScanAndFilter("test")
	assert.Equal(t, make([]string, 0), res, "Expect nil as the output of the ScanAndFilter method")
}
