// Package helpers contains useful constants
// (and possibly, in the future, utility functions as well)
// for any test gofiles.
package helpers

import (
	"path/filepath"
	"runtime"
)

// RootDir is used in tests all over the project,
// in most cases to point the analyzer to a certain
// project directory, relative to the project RootDir
var RootDir = getRootDir() //nolint:gochecknoglobals

// getRootDir is implicitly called on init.
// Its value is stored in RootDir global variable.
func getRootDir() string {
	_, thisFilePath, _, _ := runtime.Caller(0)
	// The first filepath.Dir points to "helpers" directory;
	// The second nested filepath.Dir points to its parent, which is
	// the root of the project.
	return filepath.Dir(filepath.Dir(thisFilePath))
}
