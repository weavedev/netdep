// Package discovery defines discovery of clients calls and endpoints
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package discovery

/*
In the Discovery stages, clients and endpoints are discovered and mapped to their parent service.
Refer to the Project plan, chapter 5.3 for more information.
*/

// CallData stores for each call the full path of the file in which it happens and the exact line in that file
type CallData struct {
	Filepath string `json:"filepath"`
	Line     int    `json:"line"`
}

// ServiceCalls stores for each service its name and the calls that it makes (strings of URLs / method names)
type ServiceCalls struct {
	Service string                `json:"service"`
	Calls   map[string][]CallData `json:"calls"`
}

// DiscoveredData is initialised and populated during the discovery stage.
// It stores a list of ServiceCalls for each service and a map of all handled endpoints / methods
// along with the name of the service that handles each one.
type DiscoveredData struct {
	ServCalls []ServiceCalls
	Handled   map[string]string
}

// FindCallersForEndpoint is a sample method for locating
// the callers of a specific endpoint, which is specified
// by the name of its parent service, its path in the target
// project, and its URI.
func FindCallersForEndpoint(parentService, endpointPath, endpointURI string) []interface{} {
	// This is a placeholder; the signature of this method might need to be changed.
	// Return empty slice for now.
	return nil
}
