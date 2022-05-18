/*
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"path"
	"runtime"
	"testing"
)

func TestExecuteDepScanInvalidProjectDir(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{"--project-directory", "invalid"})

	if err := runDepScanCmd.Execute(); err != nil {
		expected := "invalid project directory specified: invalid"
		if expected != err.Error() {
			t.Errorf("Error actual = %v, and Expected = %v.", err, expected)
		}
	}
}

func TestExecuteDepScanInvalidServiceDir(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{"--service-directory", "invalid"})

	if err := runDepScanCmd.Execute(); err != nil {
		expected := "invalid service directory specified: invalid"
		if expected != err.Error() {
			t.Errorf("Error actual = %v, and Expected = %v.", err, expected)
		}
	}
}

func TestExecuteDepScanNoServicePackages(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{
		"-p", "../",
		"-s", "./test/sample/http/aliased_call",
	})

	if err := runDepScanCmd.Execute(); err == nil {
		t.Error("The error was not thrown, when testing for erroneous behaviour")
	}
}

func TestExecuteDepScanNoMainFunctionFound(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{
		"-p", "../test/example",
		"-s", "./pkg",
	})

	if err := runDepScanCmd.Execute(); err == nil {
		t.Error("The error was not thrown, when testing for erroneous behaviour")
	}
}

func TestExecuteDepScanNoGoFiles(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{
		"-p", "../",
		"-s", "./test/empty",
	})

	if err := runDepScanCmd.Execute(); err == nil {
		t.Error("The error was not thrown, when testing for erroneous behaviour")
	}
}

// TODO: Good weather tests currently stackoverflow,
// as there is no base case implemented. Uncomment these
// tests after merging with dev.
func TestExecuteDepScanExampleServices(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{
		"-p", "../test/example",
		"-s", "/svc",
	})

	if err := runDepScanCmd.Execute(); err != nil {
		t.Error("The error was not thrown, when testing for erroneous behaviour")
	}
}

func TestExecuteDepScanFull(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{
		"--project-directory", "../",
		"--service-directory", "./test/sample/http",
	})

	if err := runDepScanCmd.Execute(); err != nil {
		t.Error(err)
	}
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

	if err := runDepScanCmd.Execute(); err != nil {
		t.Error(err)
	}
}
