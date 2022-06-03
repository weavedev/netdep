// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package discovery

import (
	"os"
	"path"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/preprocessing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/helpers"

	"github.com/stretchr/testify/assert"
)

func discoverAllServices(projectDir string, services []string, config *callanalyzer.AnalyserConfig) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget) {
	resC := make([]*callanalyzer.CallTarget, 0)
	resS := make([]*callanalyzer.CallTarget, 0)

	// for each service
	for _, serviceDir := range services {
		// load packages
		packagesInService, err := preprocessing.LoadAndBuildPackages(projectDir, serviceDir)
		if err != nil {
			continue
		}

		// discover calls
		clientCalls, serviceCalls, err := DiscoverAll(packagesInService, config, nil)
		if err != nil {
			continue
		}

		// append and release
		resC = append(resC, clientCalls...)
		resS = append(resS, serviceCalls...)
	}

	return resC, resS
}

/*
A test for the sample implementation of the resolution method
*/
func TestDiscovery(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, "test", "sample", "http")
	services, _ := preprocessing.FindServices(svcDir)
	resC, _ := discoverAllServices(helpers.RootDir, services, nil)

	assert.Equal(t, 17, len(resC), "Expect 17 interesting call")
	assert.Equal(t, "net/http.Get", resC[0].MethodName, "Expect net/http.Get to be called")
}

func TestDiscoveryBasicCall(t *testing.T) {
	projDir := path.Join(helpers.RootDir, "test", "sample", "http", "basic_call")
	initial, _ := preprocessing.LoadAndBuildPackages(projDir, projDir)
	resC, _, _ := DiscoverAll(initial, nil, nil)

	assert.Equal(t, 1, len(resC), "Expect 1 interesting call")
	assert.Equal(t, "net/http.Get", resC[0].MethodName, "Expect net/http.Get to be called")
}

func TestDiscoveryBasicHandle(t *testing.T) {
	projDir := path.Join(helpers.RootDir, "test", "sample", "http", "basic_handle")
	initial, _ := preprocessing.LoadAndBuildPackages(projDir, projDir)
	_, resS, _ := DiscoverAll(initial, nil, nil)

	assert.Equal(t, 2, len(resS), "Expect 2 interesting calls")
	assert.Equal(t, "net/http.Handle", resS[0].MethodName, "Expect net/http.Handle to be called")
}

func TestDiscoveryBasicHandleFunc(t *testing.T) {
	projDir := path.Join(helpers.RootDir, "test", "sample", "http", "basic_handlefunc")
	initial, _ := preprocessing.LoadAndBuildPackages(projDir, projDir)
	_, resS, _ := DiscoverAll(initial, nil, nil)

	assert.Equal(t, 2, len(resS), "Expect 2 interesting calls")
	assert.Equal(t, "net/http.HandleFunc", resS[0].MethodName, "Expect net/http.HandleFunc to be called")
}

func TestDiscoveryGinHandle(t *testing.T) {
	projDir := path.Join(helpers.RootDir, path.Join("test/sample", path.Join("http", "gin_handle")))
	initial, _ := preprocessing.LoadAndBuildPackages(projDir, projDir)
	_, resS, _ := DiscoverAll(initial, nil, nil)

	assert.Equal(t, 2, len(resS), "Expect 2 interesting calls")
	assert.Equal(t, "(*github.com/gin-gonic/gin.RouterGroup).GET", resS[0].MethodName, "Expect (*github.com/gin-gonic/gin.RouterGroup).GET to be called")
}

func TestCallInfo(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, "test", "sample", "http")
	services, _ := preprocessing.FindServices(svcDir)
	res, _ := discoverAllServices(helpers.RootDir, services, nil)

	assert.Equal(t, "multiple_calls", res[5].ServiceName, "Expected service name multiple_calls.go")
	assert.Equal(t, "25", res[7].PositionInFile, "Expected line number 25")
	assert.Equal(t, "multiple_calls"+string(os.PathSeparator)+"multiple_calls.go", res[7].FileName, "Expected file name multiple_calls/multiple_calls.go")
}

func TestWrappedClientCall(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, "test", "sample", "http", "wrapped_client")
	initial, _ := preprocessing.LoadAndBuildPackages(helpers.RootDir, svcDir)
	res, _, _ := DiscoverAll(initial, nil, nil)

	assert.Equal(t, "wrapped_client", res[0].ServiceName, "Expected service name wrapped_client.go")
	// TODO: this should fail in the future (should be 28), but it now takes the last in the list.
	assert.Equal(t, "18", res[0].PositionInFile, "Expected line number 18")
	assert.Equal(t, true, res[0].IsResolved, "Expected call to be fully resolved")
	assert.Equal(t, "http://example.com/endpoint", res[0].RequestLocation, "Expected correct URL \"http://example.com/endpoint\"")
}

func TestGetEnvCall(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, "test", "sample", "http", "env_variable")

	destinationURL := "http://example.com/endpoint"
	env := map[string]map[string]string{
		"env_variable": {
			"FOO": destinationURL,
		},
	}

	config := callanalyzer.DefaultConfigForFindingHTTPCalls(env)
	initial, _ := preprocessing.LoadAndBuildPackages(helpers.RootDir, svcDir)
	res, _, _ := DiscoverAll(initial, &config, nil)

	assert.Equal(t, "env_variable", res[0].ServiceName, "Expected service name env_variable.go")
	assert.Equal(t, "11", res[0].PositionInFile, "Expected line number 11")
	assert.Equal(t, true, res[0].IsResolved, "Expected call to be fully resolved")
	assert.Equal(t, "http://example.com/endpoint", res[0].RequestLocation, "Expected correct URL \"http://example.com/endpoint\"")
}
