package stages

import "go/ast"

/*
Copyright © 2022 Team 1, Weave BV, TU Delft

In the Filtering stages, irrelevant files and directories are removed from the target project.
Refer to the Project plan, chapter 5.1 for more information.
*/

// ScanAndFilter
/**
Returns a map with:
- Key: service name
- Value: array of the services' ASTs per file
*/
func ScanAndFilter(svcDir string) map[string][]*ast.File {
	// TODO: perhaps, for each service, filter its contents?
	var servicesList = findAllServices(svcDir)
	for i := 0; i < len(servicesList); i++ {
		var _ = filter(servicesList[i], nil)
		//TODO: add to map the resulting AST array
	}
	filter("test", nil)

	return nil
}

// FindAllServices
/**
A method for finding all services, which takes the path of the svc directory as an argument
Returns a list of all services.
*/
// TODO: Remove the following line when implementing this method
//goland:noinspection GoUnusedParameter,GoUnusedFunction
func findAllServices(svcDir string) []string {
	// TODO extract a list of the paths of each service
	return nil
}

// Filter
/**
A placeholder method for filtering the directory of a specified service.
Return a list of ASTs (of each of the files).
*/
// TODO: Remove the following line when implementing this method
//goland:noinspection GoUnusedParameter,GoUnusedFunction
func filter(serviceLoc string, filterList []string) []*ast.File {
	// TODO: This is a placeholder; the signature of this method might need to be changed.
	// TODO: Loop over all subdirectories/files of this service, looking for relevant files
	// Return empty slice for now.
	return make([]*ast.File, 0)
}
