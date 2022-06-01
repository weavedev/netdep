package discovery

import (
	"fmt"
	"path"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/helpers"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

/*
TestEnvVarResolution tests getEnv substitution. After passing
the env vars to discovery, they can be properly substituted
*/
func TestEnvVarResolution(t *testing.T) {
	projDir := path.Join(helpers.RootDir, "test", "example")
	svcDir := path.Join(helpers.RootDir, "test", "example", "env_svc")
	initial, _ := stages.LoadServices(projDir, svcDir)
	destinationURL := "127.0.0.1:8081"
	env := map[string]map[string]string{
		"env_variable": {
			"FOO": destinationURL,
		},
	}
	// configWithEnv := callanalyzer.DefaultConfigForFindingHTTPCalls(env)
	resC, _, _ := Discover(initial, nil, env)
	assert.Equal(t, 1, len(resC), "Expect 1 interesting call")
	assert.Equal(t, destinationURL, resC[0].RequestLocation, fmt.Sprintf("Expect %s", destinationURL))
}
