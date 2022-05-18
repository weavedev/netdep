/*
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"testing"
)

func TestExecuteDepScanFull(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{
		"--project-directory", "../",
		"--service-directory", "./test/sample/http/aliased_call",
	})

	if err := runDepScanCmd.Execute(); err != nil {
		t.Error(err)
	}
}

func TestExecuteDepScanShorthand(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{
		"-p", "../",
		"-s", "./test/sample/http/aliased_call",
	})

	if err := runDepScanCmd.Execute(); err != nil {
		t.Error(err)
	}
}

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

func TestExecuteDepScanNoGoFiles(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{
		"-p", "../test/example",
		"-s", "/svc",
	})

	if err := runDepScanCmd.Execute(); err == nil {
		t.Error("The error was not thrown, when testing for erroneous behaviour")
	}
}

func TestExecuteDepScanBasicCall(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{
		"-p", "../",
		"-s", "./test/sample/http/basic_call",
	})

	if err := runDepScanCmd.Execute(); err != nil {
		t.Error(err)
	}
}

func TestExecuteDepScanNoMainPackage(t *testing.T) {
	runDepScanCmd := depScanCmd()
	runDepScanCmd.SetArgs([]string{
		"-p", "../",
		"-s", "./test/example/pkg/http",
	})

	if err := runDepScanCmd.Execute(); err == nil {
		t.Error("The error was not thrown, when testing for erroneous behaviour")
	}
}
