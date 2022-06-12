package main

import (
	"os"
	"path/filepath"
	"testing"

	"lab.weave.nl/internships/tud-2022/netDep/helpers"
)

func TestRunRootSmokeTest(_ *testing.T) {
	dummyProjDir := filepath.Join(helpers.RootDir, "test", "example")
	dummySvcDir := filepath.Join(helpers.RootDir, "test", "example", "svc")

	osArgsBackup := os.Args
	os.Args = []string{"netDep.exe", "-p", dummyProjDir, "-s", dummySvcDir}
	runRoot()
	os.Args = osArgsBackup
}
