package stages

import (
	"fmt"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callgraph"
	"reflect"
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
	// We get the absolute URL of the project as a parameter.
	// Then we should locate all services and build a graph for each of them, separately.
	arr := make([]string, 1)
	arr[0] = "/Users/martynaskrupskis/Documents/code/main"
	var smth, _ = callgraph.DoCallGraph("", arr)
	fmt.Println(reflect.TypeOf(smth))
	return nil
}
