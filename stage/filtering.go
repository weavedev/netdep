package stage

import "go/ast"

/*
Copyright © 2022 Team 1, Weave BV, TU Delft

In the Filtering stage, irrelevant files and directories are removed from the target project.
Refer to the Project plan, chapter 5.1 for more information.
*/

/**
The filter method, which takes the path of the service directory as an argument
Returns a filtered AST.
*/
func filter(svcPath string) []*ast.File {
	//Return empty slice for now
	return make([]*ast.File, 0)
}
