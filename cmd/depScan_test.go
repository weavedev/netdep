package cmd

import (
	"testing"
)

func Test_ExecuteDepScan_Full(t *testing.T) {

	depScanCmd := depScanCmd()
	depScanCmd.SetArgs([]string{
		"--project-directory", "./",
		"--service-directory", "./",
	})

	if err := depScanCmd.Execute(); err != nil {
		t.Error(err)
	}
}

func Test_ExecuteDepScan_Shorthand(t *testing.T) {

	depScanCmd := depScanCmd()
	depScanCmd.SetArgs([]string{
		"-p", "./",
		"-s", "./",
	})

	if err := depScanCmd.Execute(); err != nil {
		t.Error(err)
	}
}

func Test_ExecuteDepScan_InvalidProjectDir(t *testing.T) {

	depScanCmd := depScanCmd()
	depScanCmd.SetArgs([]string{"--project-directory", "invalid"})

	if err := depScanCmd.Execute(); err != nil {
		expected := "invalid project directory specified: invalid"
		if expected != err.Error() {
			t.Errorf("Error actual = %v, and Expected = %v.", err, expected)
		}
	}
}

func Test_ExecuteDepScan_InvalidServiceDir(t *testing.T) {

	depScanCmd := depScanCmd()
	depScanCmd.SetArgs([]string{"--service-directory", "invalid"})

	if err := depScanCmd.Execute(); err != nil {
		expected := "invalid service directory specified: invalid"
		if expected != err.Error() {
			t.Errorf("Error actual = %v, and Expected = %v.", err, expected)
		}
	}
}
