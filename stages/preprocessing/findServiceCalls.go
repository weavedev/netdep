package preprocessing

import (
	"os"
	"path/filepath"
	"strings"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"
)

// LoadServiceCalls scans all the files of a given service directory and returns a list of
// Annotation from the comments in the format "//netdep: ..." that it discovers.
func LoadServiceCalls(servicePath string, serviceName string, internalCalls map[IntCall]string, clientTargets *[]*callanalyzer.CallTarget) error {
	files, err := os.ReadDir(servicePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".go" && !strings.HasSuffix(file.Name(), "_test.go") && !strings.HasSuffix(file.Name(), "pb.go") {
			// If the file is a .go file - parse it
			currClientTargets, err := parseMethods(filepath.Join(servicePath, file.Name()), internalCalls, serviceName)
			if err != nil {
				return err
			}
			*clientTargets = append(*clientTargets, *currClientTargets...)
		} else if file.IsDir() {
			// If the file is a directory - recursively look for .go files inside it
			err := LoadServiceCalls(filepath.Join(servicePath, file.Name()), serviceName, internalCalls, clientTargets)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
