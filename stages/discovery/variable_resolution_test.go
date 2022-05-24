package discovery

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"
	"path"
	"runtime"
	"testing"
)

/*
TestEnvVarResolution tests getEnv substitution. For now with no asserts,
as we don't have a nice way to pass env values to the discovery stage.
*/
func TestEnvVarResolution(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Dir(path.Dir(thisFileParent))
	svcDir := path.Join(path.Dir(path.Dir(thisFileParent)), "test", "sample", "http", "env_variable")

	initial, _ := stages.LoadServices(projDir, svcDir)
	destinationURL := "127.0.0.1:8081"
	env := map[string]map[string]string{
		"env_service": {
			"FOO": destinationURL,
		},
	}
	configWithEnv := callanalyzer.DefaultConfigForFindingHTTPCalls(env)
	resC, _, _ := Discover(initial, &configWithEnv)
	assert.Equal(t, 1, len(resC), "Expect 1 interesting call")
	assert.Equal(t, destinationURL, resC[0].RequestLocation, fmt.Sprintf("Expect %s", destinationURL))
}
