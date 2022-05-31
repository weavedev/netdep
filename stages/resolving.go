// Package stages defines different stages of analysis
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package stages

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

/*
In the Resolving stages, the file supplied by the user containing Environment variables is traversed and stored in a map as follows:
Map<String(serviceName), Map<String(variable name), String(variable value)>>.
This map is the output of the resolving stages.
Refer to the Project plan, chapter 5.2 for more information.
*/

// MapEnvVars returns a map as described above, namely:
// map{ service: map{ var.name: var.value }}

func MapEnvVars(path string) (map[string]map[string]string, error) {
	file, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("the file cannot be parsed")
	}

	envVars := make(map[string]map[string]string)

	err2 := yaml.Unmarshal(file, &envVars)

	if err2 != nil {
		return nil, fmt.Errorf("the file cannot be parsed")
	}
	return envVars, nil
}
