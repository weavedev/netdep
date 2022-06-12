package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
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
