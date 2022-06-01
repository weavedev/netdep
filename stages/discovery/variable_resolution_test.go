package discovery

import (
	"fmt"
	"path"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/helpers"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"
)

/*
TestEnvVarResolution tests getEnv substitution. For now with no asserts,
as we don't have a nice way to pass env values to the discovery stage.
*/
func TestEnvVarResolution(t *testing.T) {
	projDir := path.Join(helpers.RootDir, "test", "example")
	svcDir := path.Join(helpers.RootDir, "test", "example", "env_svc")
	services, _ := stages.LoadServices(svcDir)
	destinationURL := "127.0.0.1:8081"
	env := map[string]map[string]string{
		"env_variable": {
			"FOO": destinationURL,
		},
	}
	configWithEnv := callanalyzer.DefaultConfigForFindingHTTPCalls(env)
	resC, _, _ := discoverAllServices(projDir, services, &configWithEnv)
	assert.Equal(t, 1, len(resC), "Expect 1 interesting call")
	assert.Equal(t, destinationURL, resC[0].RequestLocation, fmt.Sprintf("Expect %s", destinationURL))
}
