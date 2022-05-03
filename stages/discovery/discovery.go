package stages

import (
	"fmt"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callgraph"
)

func main() {

}

/*
Copyright © 2022 Team 1, Weave BV, TU Delft

In the Discovery stages, clients and endpoints are discovered and mapped to their parent service.
Refer to the Project plan, chapter 5.3 for more information.
*/

// FindCallersForEndpoint
/**
A sample method for locating the callers of a specific endpoint,
which is specified by the name of its parent service, its path in the target project,
and its URI.
*/
func FindCallersForEndpoint(parentService string, endpointPath string, endpointURI string) []interface{} {
	// This is a placeholder; the signature of this method might need to be changed.
	// Return empty slice for now.
	arr := make([]string, 1)
	arr[0] = "stages"
	fmt.Println(callgraph.DoCallGraph("", arr))
	return nil
}
