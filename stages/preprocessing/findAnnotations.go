// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

type Position struct {
	Filename string
	Line     int
}

// LoadAnnotations scans all the files of a given service directory and returns a list of
// Annotation from the comments in the format "//netdep: ..." that it discovers.
func LoadAnnotations(servicePath string, serviceName string, annotations map[string]map[Position]string) error {
	files, err := os.ReadDir(servicePath)
	if err != nil {
		return err
	}

	annotations[serviceName] = make(map[Position]string)

	for _, file := range files {
		if file.Name()[len(file.Name())-3:] == ".go" {
			// If the file is a .go file - parse it
			err := parseComments(path.Join(servicePath, file.Name()), serviceName, annotations)
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			// If the file is a directory - recursively look for .go files inside it
			err := LoadAnnotations(path.Join(servicePath, file.Name()), serviceName, annotations)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// parseComments parses the given file with a parser.ParseComments mode, filters out
// the comments which don't contain a substring "netdep:client" or "netdep:endpoint", generates an Annotation for
// every remaining comment and returns a list of them.
func parseComments(path string, serviceName string, annotations map[string]map[Position]string) error {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, commentGroup := range f.Comments {
		for _, comment := range commentGroup.List {
			if strings.HasPrefix(comment.Text, "//netdep:") {
				tokenPos := fs.Position(comment.Slash)
				pos := Position{
					Filename: tokenPos.Filename,
					Line:     tokenPos.Line,
				}
				value := strings.Join(strings.Split(comment.Text, "netdep:")[1:], "")

				annotations[serviceName][pos] = value
			}
		}
	}

	return nil
}
