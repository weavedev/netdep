// Package discovery defines discovery of clients calls and endpoints
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package discovery

import (
	"fmt"
	"path/filepath"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/preprocessing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/helpers"

	"github.com/stretchr/testify/assert"
)

/*
TestEnvVarResolution tests getEnv substitution. After passing
the env vars to discovery, they can be properly substituted
*/
func TestEnvVarResolution(t *testing.T) {
	projDir := filepath.Join(helpers.RootDir, "test", "example")
	svcDir := filepath.Join(helpers.RootDir, "test", "example", "env_svc")
	services, _ := preprocessing.FindServices(svcDir)
	destinationURL := "127.0.0.1:8081"
	env := map[string]map[string]string{
		"env_variable": {
			"FOO": destinationURL,
		},
	}

	configWithEnv := callanalyzer.DefaultConfigForFindingHTTPCalls()
	configWithEnv.SetEnv(env)
	resC, resS := discoverAllServices(projDir, services, &configWithEnv)

	assert.Equal(t, 1, len(resC), "Expect 1 interesting call")
	assert.Equal(t, 0, len(resS), "Expect 0 interesting call")
	assert.Equal(t, destinationURL, resC[0].RequestLocation, fmt.Sprintf("Expect %s", destinationURL))
}
