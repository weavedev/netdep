package discovery

import (
	"github.com/stretchr/testify/assert"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
	"path"
	"runtime"
	"testing"
)

/*
A test for the sample implementation of the resolution method
*/
func TestEnvVarResolution(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Dir(path.Dir(thisFileParent))
	svcDir := path.Join(path.Dir(path.Dir(thisFileParent)), "test", "sample", "http", "env_variable")

	initial, _ := stages.LoadServices(projDir, svcDir)
	resC, _, _ := Discover(initial)
	assert.Equal(t, 13, len(resC), "Expect 12 interesting call")
	assert.Equal(t, "net/http.Get", resC[0].MethodName, "Expect net/http.Get to be called")
}
