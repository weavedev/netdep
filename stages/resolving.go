// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

/*
In the Resolving stages, files containing Environment variables are traversed and stored in a map as follows:
Map<String(serviceName), Map<String(variable name), String(variable value)>>.
This map is the output of the resolving stages.
Refer to the Project plan, chapter 5.2 for more information.
*/

// ResolveEnvVars
// returns a map as described above, namely:
//
// map{ service: map{ var.name: var.value }}
func ResolveEnvVars(svcDir string) map[string]map[string]interface{} {
	//TODO: Implement the resolution of environment variables
	var testSvcName = svcDir + "SampleService"
	m := make(map[string]map[string]interface{})
	m[testSvcName] = make(map[string]interface{})
	m[testSvcName]["VariableNameA"] = "1"
	m[testSvcName]["VariableNameB"] = "False"
	return m
}
