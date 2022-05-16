/*
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"testing"
)

func TestExecuteDepScanFull(t *testing.T) {
	depScanCmd := depScanCmd()
	depScanCmd.SetArgs([]string{
		"--project-directory", "../",
		"--service-directory", "./test/sample/http/aliased_call",
	})

	if err := depScanCmd.Execute(); err != nil {
		t.Error(err)
	}
}

func TestExecuteDepScanShorthand(t *testing.T) {
	depScanCmd := depScanCmd()
	depScanCmd.SetArgs([]string{
		"-p", "../",
		"-s", "./test/sample/http/aliased_call",
	})

	if err := depScanCmd.Execute(); err != nil {
		t.Error(err)
	}
}

func TestExecuteDepScanInvalidProjectDir(t *testing.T) {
	depScanCmd := depScanCmd()
	depScanCmd.SetArgs([]string{"--project-directory", "invalid"})

	if err := depScanCmd.Execute(); err != nil {
		expected := "invalid project directory specified: invalid"
		if expected != err.Error() {
			t.Errorf("Error actual = %v, and Expected = %v.", err, expected)
		}
	}
}

func TestExecuteDepScanInvalidServiceDir(t *testing.T) {
	depScanCmd := depScanCmd()
	depScanCmd.SetArgs([]string{"--service-directory", "invalid"})

	if err := depScanCmd.Execute(); err != nil {
		expected := "invalid service directory specified: invalid"
		if expected != err.Error() {
			t.Errorf("Error actual = %v, and Expected = %v.", err, expected)
		}
	}
}

func TestExecuteDepScanNoGoFiles(t *testing.T) {
	depScanCmd := depScanCmd()
	depScanCmd.SetArgs([]string{
		"-p", "../",
		"-s", "../test/example/",
	})

	if err := depScanCmd.Execute(); err == nil {
		t.Error("The error was not thrown, when testing for erroneous behaviour")
	}
}

func TestExecuteDepScanBasicCall(t *testing.T) {
	depScanCmd := depScanCmd()
	depScanCmd.SetArgs([]string{
		"-p", "../",
		"-s", "./test/sample/http/basic_call",
	})

	if err := depScanCmd.Execute(); err != nil {
		t.Error(err)
	}
}

func TestExecuteDepScanNoMainPackage(t *testing.T) {
	depScanCmd := depScanCmd()
	depScanCmd.SetArgs([]string{
		"-p", "../",
		"-s", "./test/example/pkg/http",
	})

	if err := depScanCmd.Execute(); err == nil {
		t.Error("The error was not thrown, when testing for erroneous behaviour")
	}
}
