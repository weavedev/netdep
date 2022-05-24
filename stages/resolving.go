// Package stages defines different stages of analysis
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package stages

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

/*
In the Resolving stages, files containing Environment variables are traversed and stored in a map as follows:
Map<String(serviceName), Map<String(variable name), String(variable value)>>.
This map is the output of the resolving stages.
Refer to the Project plan, chapter 5.2 for more information.
*/

// ResolveEnvVars returns a map as described above, namely:
// map{ service: map{ var.name: var.value }}
func ResolveEnvVars(svcDir string) map[string]map[string]interface{} {
	m := make(map[string]map[string]interface{})

	// iterate through service directory
	items, err := ioutil.ReadDir(svcDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range items {
		// for every service within the directory, find .yaml files
		if item.IsDir() {
			env := make(map[string]interface{})
			yamlFiles := findYaml(svcDir+"/"+item.Name(), ".yaml")
			// for every .yaml file, create a map of env vars
			for _, file := range yamlFiles {
				envVars := envMap(file)
				for k, v := range envVars {
					env[k] = v
				}
			}
			// append env map to service name key
			if len(env) != 0 {
				m[item.Name()] = env
			}
		} else {
			continue
		}
	}
	return m
}

// given a directory, extract files with "ext" extension
func findYaml(root, ext string) []string {
	var yamlFiles []string
	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			log.Fatal(e)
		}
		if filepath.Ext(d.Name()) == ext {
			yamlFiles = append(yamlFiles, s)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return yamlFiles
}

// given a .yaml file, create a map of env vars(name, value)
func envMap(path string) map[string]interface{} {
	file, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		log.Fatal(err)
	}

	envVars := make(map[string]interface{})

	err2 := yaml.Unmarshal(file, &envVars)

	if err2 != nil {
		log.Fatal(err2)
	}
	return envVars
}

func MapEnvVarFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	var envVars map[string]string
	if err := json.NewDecoder(file).Decode(&envVars); err != nil {
		return nil, fmt.Errorf("the file cannot be parsed")
	}
	return envVars, nil
}
