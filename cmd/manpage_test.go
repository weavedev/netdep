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
