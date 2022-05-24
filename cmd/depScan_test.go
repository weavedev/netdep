/*
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteDepScanInvalidProjectDir(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{"--project-directory", "invalid"})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid project directory specified: invalid", err.Error())
}

func TestExecuteDepScanInvalidServiceDir(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{"--service-directory", "invalid"})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid service directory specified: invalid", err.Error())
}

func TestExecuteDepScanNoServicePackages(t *testing.T) {
	runDepScanCmd := depScanCmd()

	// thisFilePath is ./cmd/depScan_test.go
	_, thisFilePath, _, _ := runtime.Caller(0)
	projDir := path.Dir(path.Dir(thisFilePath)) // root of the project
	svcDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "sample", "http", "aliased_call")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "no service to analyse were found", err.Error())
}

func TestExecuteDepScanNoMainFunctionFound(t *testing.T) {
	runDepScanCmd := depScanCmd()

	// thisFilePath is ./cmd/depScan_test.go
	_, thisFilePath, _, _ := runtime.Caller(0)
	projDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "example") // root of the project
	svcDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "example", "pkg")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "no main function found in package package example/pkg/http", err.Error())
}

func TestExecuteDepScanNoGoFiles(t *testing.T) {
	runDepScanCmd := depScanCmd()

	// thisFilePath is ./cmd/depScan_test.go
	_, thisFilePath, _, _ := runtime.Caller(0)
	projDir := path.Dir(path.Dir(thisFilePath)) // root of the project
	svcDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "empty")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "packages contain errors", err.Error())
}

// TODO: Good weather tests currently stackoverflow,
// as there is no base case implemented. Uncomment these
// tests after merging with dev.
func TestExecuteDepScanExampleServices(t *testing.T) {
	runDepScanCmd := depScanCmd()

	// thisFilePath is ./cmd/depScan_test.go
	_, thisFilePath, _, _ := runtime.Caller(0)
	projDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "example") // root of the project
	svcDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "example", "svc")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.Nil(t, err)
}

func TestExecuteDepScanFull(t *testing.T) {
	runDepScanCmd := depScanCmd()
	// thisFilePath is ./cmd/depScan_test.go
	_, thisFilePath, _, _ := runtime.Caller(0)
	projDir := path.Dir(path.Dir(thisFilePath)) // root of the project
	svcDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "sample", "http")
	runDepScanCmd.SetArgs([]string{
		"--project-directory", projDir,
		"--service-directory", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.Nil(t, err)
}

func TestExecuteDepScanShortHand(t *testing.T) {
	runDepScanCmd := depScanCmd()

	// thisFilePath is ./cmd/depScan_test.go
	_, thisFilePath, _, _ := runtime.Caller(0)
	projDir := path.Dir(path.Dir(thisFilePath)) // root of the project
	svcDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "sample", "http")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.Nil(t, err)
}

func TestExecuteDepScanInvalidEnvVarFile(t *testing.T) {
	runDepScanCmd := depScanCmd()
	_, thisFilePath, _, _ := runtime.Caller(0)
	svcDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "example", "svc")

	runDepScanCmd.SetArgs([]string{
		"-s", svcDir,
		"-e", "invalid",
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid environment variable file specified: invalid", err.Error())
}

func TestExecuteDepScanEnvFile(t *testing.T) {
	runDepScanCmd := depScanCmd()

	// thisFilePath is ./cmd/depScan_test.go
	_, thisFilePath, _, _ := runtime.Caller(0)
	projDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "example") // root of the project
	svcDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "example", "svc")
	envVars := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "example", "svc", "node-basic-http", "env")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
		"-e", envVars,
	})

	err := runDepScanCmd.Execute()
	assert.Nil(t, err)
}

func TestExecuteDepScanEnvFileWrongFormat(t *testing.T) {
	runDepScanCmd := depScanCmd()

	// thisFilePath is ./cmd/depScan_test.go
	_, thisFilePath, _, _ := runtime.Caller(0)
	projDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "example") // root of the project
	svcDir := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "example", "svc")
	envVars := path.Join(path.Dir(path.Dir(thisFilePath)), "test", "example", "svc", "node-basic-http", "values.yaml")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
		"-e", envVars,
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "the file cannot be parsed", err.Error())
}
