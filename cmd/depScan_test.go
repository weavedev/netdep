/*
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"path"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/helpers"

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
	svcDir := path.Join(helpers.RootDir, "test", "sample", "http", "aliased_call")

	runDepScanCmd.SetArgs([]string{
		"-p", helpers.RootDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "no service to analyse were found", err.Error())
}

func TestExecuteDepScanNoMainFunctionFound(t *testing.T) {
	runDepScanCmd := depScanCmd()
	projDir := path.Join(helpers.RootDir, "test", "example") // root of the project
	svcDir := path.Join(helpers.RootDir, "test", "example", "pkg")

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
	svcDir := path.Join(helpers.RootDir, "test", "empty")

	runDepScanCmd.SetArgs([]string{
		"-p", helpers.RootDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.NotNil(t, err)
	// we found no go files, but that is an error for the builder
	assert.Equal(t, "packages contain errors", err.Error())
}

// TODO: Good weather tests currently stackoverflow,
// as there is no base case implemented. Uncomment these
// tests after merging with dev.
func TestExecuteDepScanExampleServices(t *testing.T) {
	runDepScanCmd := depScanCmd()

	projDir := path.Join(helpers.RootDir, "test", "example")
	svcDir := path.Join(helpers.RootDir, "test", "example", "svc")

	runDepScanCmd.SetArgs([]string{
		"-p", projDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.Nil(t, err)
}

func TestExecuteDepScanFull(t *testing.T) {
	runDepScanCmd := depScanCmd()
	svcDir := path.Join(helpers.RootDir, "test", "sample", "http")
	runDepScanCmd.SetArgs([]string{
		"--project-directory", helpers.RootDir,
		"--service-directory", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.Nil(t, err)
}

func TestExecuteDepScanShortHand(t *testing.T) {
	runDepScanCmd := depScanCmd()
	svcDir := path.Join(helpers.RootDir, "test", "sample", "http")

	runDepScanCmd.SetArgs([]string{
		"-p", helpers.RootDir,
		"-s", svcDir,
	})

	err := runDepScanCmd.Execute()
	assert.Nil(t, err)
}
