package stage

import (
	"github.com/stretchr/testify/assert"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
	"testing"
)

/**
A test for the sample implementation of the resolution method
*/
func TestBuilding(t *testing.T) {
	var res = stages.ConstructOutput(nil)
	assert.Equal(t, "{ type: \"error\", message: \"Not Implemented\" }", res, "Expect nil as the output of the sample ConstructOutput method")
}
