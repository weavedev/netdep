/*
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"lab.weave.nl/internships/tud-2022/netDep/helpers"

	"github.com/stretchr/testify/assert"
)

func TestExecuteDepScanInvalidProjectDir(t *testing.T) {
	runDepScanCmd := RootCmd()
	runDepScanCmd.SetArgs([]string{"--project-directory", "invalid"})
	cwd, err := os.Getwd()
	assert.Nil(t, err)
	invalidPath := filepath.Join(cwd, "invalid")

	err = runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("invalid project directory specified: %s", invalidPath), err.Error())
}

func TestExecuteDepScanInvalidServiceDir(t *testing.T) {
	runDepScanCmd := RootCmd()
	runDepScanCmd.SetArgs([]string{"--service-directory", "invalid"})

	cwd, err := os.Getwd()
	assert.Nil(t, err)
	invalidPath := filepath.Join(cwd, "invalid")

	err = runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid service directory: "+invalidPath, err.Error())
}

func TestExecuteDepScanNoServicePackages(t *testing.T) {
	runDepScanCmd := RootCmd()
	svcDir := filepath.Join(helpers.RootDir, "test", "sample", "http", "aliased_call")

	runDepScanCmd.SetArgs([]string{
		"-p", helpers.RootDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "no service to analyse were found", err.Error())
}

func TestExecuteDepScanNoMainFunctionFound(t *testing.T) {
	runDepScanCmd := RootCmd()
	projDir := filepath.Join(helpers.RootDir, "test", "example") // root of the project
	svcDir := filepath.Join(helpers.RootDir, "test", "example", "pkg")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "no main function found in package package example/pkg/http", err.Error())
}

func TestExecuteDepScanNoGoFiles(t *testing.T) {
	runDepScanCmd := RootCmd()
	svcDir := filepath.Join(helpers.RootDir, "test", "empty")

	runDepScanCmd.SetArgs([]string{
		"-p", helpers.RootDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	// we found no go files, but that is an error for the builder
	assert.Equal(t, "no usable packages found", err.Error())
}

// TODO: Good weather tests currently stackoverflow,
// as there is no base case implemented. Uncomment these
// tests after merging with dev.
func TestExecuteDepScanExampleServicesVerbose(t *testing.T) {
	runDepScanCmd := RootCmd()

	projDir := filepath.Join(helpers.RootDir, "test", "example")
	svcDir := filepath.Join(helpers.RootDir, "test", "example", "svc")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
		"-v",
	})

	err := runDepScanCmd.Execute()
	assert.Nil(t, err)
}

func TestExecuteDepScanFull(t *testing.T) {
	runDepScanCmd := RootCmd()
	svcDir := filepath.Join(helpers.RootDir, "test", "sample", "http")
	runDepScanCmd.SetArgs([]string{
		"--project-directory", helpers.RootDir,
		"--service-directory", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.Nil(t, err)
}

func TestExecuteDepScanShortHand(t *testing.T) {
	runDepScanCmd := RootCmd()
	svcDir := filepath.Join(helpers.RootDir, "test", "sample", "http")

	runDepScanCmd.SetArgs([]string{
		"-p", helpers.RootDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.Nil(t, err)
}

func TestExecuteDepScanInvalidEnvVarFile(t *testing.T) {
	runDepScanCmd := RootCmd()
	svcDir := filepath.Join(helpers.RootDir, "test", "example", "svc")

	runDepScanCmd.SetArgs([]string{
		"-s", svcDir,
		"-e", "invalid",
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid environment variable file specified: invalid", err.Error())
}

func TestExecuteDepScanEnvFile(t *testing.T) {
	runDepScanCmd := RootCmd()

	projDir := filepath.Join(helpers.RootDir, "test", "sample") // root of the project
	svcDir := filepath.Join(helpers.RootDir, "test", "sample", "http")
	envVars := filepath.Join(helpers.RootDir, "test", "sample", "http", "env_variable", "env")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
		"-e", envVars,
	})

	err := runDepScanCmd.Execute()
	assert.Nil(t, err)
}

func TestExecuteDepScanEnvFileWrongFormat(t *testing.T) {
	runDepScanCmd := RootCmd()

	projDir := filepath.Join(helpers.RootDir, "test", "example") // root of the project
	svcDir := filepath.Join(helpers.RootDir, "test", "example", "svc")
	envVars := filepath.Join(helpers.RootDir, "test", "example", "svc", "node-basic-http", "values.yaml")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
		"-e", envVars,
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "the file cannot be parsed", err.Error())
}

func TestOutputToInvalidFile(t *testing.T) {
	err := printOutput("/../badPath/", "{\"key\": \"dummyJSON\"}", nil, nil)
	assert.NotNil(t, err)
}
