package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/netDep/helpers"
)

func TestRunRootSmokeTest(t *testing.T) {
	dummyProjDir := filepath.Join(helpers.RootDir, "test", "example")
	dummySvcDir := filepath.Join(helpers.RootDir, "test", "example", "svc")

	osArgsBackup := os.Args
	os.Args = []string{"netDep.exe", "-p", dummyProjDir, "-s", dummySvcDir}
	err := runRoot()
	assert.Nil(t, err)
	os.Args = osArgsBackup
}
