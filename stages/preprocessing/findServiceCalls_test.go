package preprocessing

import (
	"github.com/stretchr/testify/assert"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/helpers"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"
	"path/filepath"
	"testing"
)

func TestLoadServiceCalls(t *testing.T) {
	svcDir := filepath.Join(helpers.RootDir, "test", "sample", "http", "gin_handle")
	intCalls := make(map[IntCall]string)
	clientTargets := make([]*callanalyzer.CallTarget, 0)

	intCall := IntCall{
		Name:      "Default",
		NumParams: 0,
	}

	intCalls[intCall] = "serviceA"

	LoadServiceCalls(svcDir, "gin_handle", intCalls, &clientTargets)

	assert.Equal(t, 1, len(clientTargets))
	assert.Equal(t, "Default", clientTargets[0].MethodName)
}

func TestLoadServiceCallsInvalidPath(t *testing.T) {
	intCalls := make(map[IntCall]string)
	clientTargets := make([]*callanalyzer.CallTarget, 0)

	err := LoadServiceCalls("invalidPath", "serviceName", intCalls, &clientTargets)
	assert.NotNil(t, err)
}
