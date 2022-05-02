package stage

/*
Copyright © 2022 Team 1, Weave BV, TU Delft

In the Resolving stage, files containing Environment variables are traversed and stored in a map as follows:
Map<String(serviceName), Map<String(variable name), String(variable value)>>.
This map is the output of the resolving stage.
Refer to the Project plan, chapter 5.2 for more information.
*/

func ResolveEnvVars(svcDir string) map[string]map[string]interface{} {
	//TODO: Implement the resolution of environment variables
	var testSvcName = svcDir + "SampleService"
	m := make(map[string]map[string]interface{})
	m[testSvcName] = make(map[string]interface{})
	m[testSvcName]["VariableNameA"] = "1"
	m[testSvcName]["VariableNameB"] = "False"
	return m
}
