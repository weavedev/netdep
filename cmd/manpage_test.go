package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenManPage(t *testing.T) {
	err := genManPageToDir(RootCmd(), os.TempDir())
	if err != nil {
		panic(err)
	}
	targetFile := filepath.Join(os.TempDir(), "netDep.1")
	assert.FileExists(t, targetFile)
	_ = os.Remove(targetFile)
}

func TestGenManPageCmd(t *testing.T) {
	cmd := GenManpageCmd(RootCmd())
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	targetFile := filepath.Join(wd, "netDep.1")
	assert.FileExists(t, targetFile)
	_ = os.Remove(targetFile)
}
