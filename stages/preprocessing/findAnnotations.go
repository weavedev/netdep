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

type Annotation struct {
	ServiceName string
	Position    token.Position
	Value       string
}

// LoadAnnotations scans all the files of a given service directory and returns a list of
// Annotation from the comments in the format "//netdep: ..." that it discovers.
func LoadAnnotations(servicePath string, serviceName string) ([]*Annotation, error) {
	files, err := os.ReadDir(servicePath)
	if err != nil {
		return nil, err
	}

	serviceAnnotations := make([]*Annotation, 0)

	for _, file := range files {
		if file.Name()[len(file.Name())-3:] == ".go" {
			// If the file is a .go file - parse it
			fileAnnotations, err := parseComments(path.Join(servicePath, file.Name()), serviceName)
			if err != nil {
				return nil, err
			}
			serviceAnnotations = append(serviceAnnotations, fileAnnotations...)
		} else if file.IsDir() {
			// If the file is a directory - recursively look for .go files inside it
			innerServiceAnnotations, err := LoadAnnotations(path.Join(servicePath, file.Name()), serviceName)
			if err != nil {
				return nil, err
			}
			serviceAnnotations = append(serviceAnnotations, innerServiceAnnotations...)
		}
	}

	return serviceAnnotations, nil
}

// parseComments parses the given file with a parser.ParseComments mode, filters out
// the comments which don't contain a substring "netdep:client" or "netdep:endpoint", generates an Annotation for
// every remaining comment and returns a list of them.
func parseComments(path string, serviceName string) ([]*Annotation, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	annotations := make([]*Annotation, 0)

	for _, commentGroup := range f.Comments {
		for _, comment := range commentGroup.List {
			if strings.HasPrefix(comment.Text, "//netdep:client") || strings.HasPrefix(comment.Text, "//netdep:endpoint") {
				ann := &Annotation{
					ServiceName: serviceName,
					Position:    fs.Position(comment.Slash),
					Value:       strings.Join(strings.Split(comment.Text, "netdep:")[1:], ""),
				}
				annotations = append(annotations, ann)
			}
		}
	}

	return annotations, nil
}
