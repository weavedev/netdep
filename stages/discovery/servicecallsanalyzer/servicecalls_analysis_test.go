package servicecallsanalyzer

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/netDep/helpers"
	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"
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

func TestFindServiceCalls(t *testing.T) {
	serviceCallsDir := filepath.Join(helpers.RootDir, "test", "sample", "servicecalls")

	internalCalls, _, _ := ParseServiceCallsPackage(serviceCallsDir)

	posOne := IntCall{
		Name:      "FirstMethod",
		NumParams: 3,
	}
	posTwo := IntCall{
		Name:      "SecondMethod",
		NumParams: 3,
	}
	posThree := IntCall{
		Name:      "ThirdMethod",
		NumParams: 3,
	}

	assert.Equal(t, 3, len(internalCalls))

	assert.Equal(t, "test", internalCalls[posOne])
	assert.Equal(t, "test", internalCalls[posTwo])
	assert.Equal(t, "test", internalCalls[posThree])
}

func TestFindServiceCallsEmptyDir(t *testing.T) {
	serviceCallsDir := ""
	internalCalls, serverTargets, _ := ParseServiceCallsPackage(serviceCallsDir)

	assert.Equal(t, 0, len(internalCalls))
	assert.Equal(t, 0, len(*serverTargets))
}

func TestFindServiceCallsInvalidDir(t *testing.T) {
	serviceCallsDir := "invalidDir"
	internalCalls, serverTargets, err := ParseServiceCallsPackage(serviceCallsDir)

	assert.Equal(t, 0, len(internalCalls))
	assert.Equal(t, 0, len(*serverTargets))
	assert.NotNil(t, err)
}
