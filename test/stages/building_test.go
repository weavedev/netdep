package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

/**
A test for the sample implementation of the resolution method
*/
func TestBuilding(t *testing.T) {
	var res = stages.ConstructOutput(nil)
	assert.Equal(t, "{ type: \"error\", message: \"Not Implemented\" }", res, "Expect nil as the output of the sample ConstructOutput method")
}
