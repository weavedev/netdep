// Package stages defines different stages of analysis
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

/*
In the Discovery stages, clients and endpoints are discovered and mapped to their parent service.
Refer to the Project plan, chapter 5.3 for more information.
*/

// ServiceCalls stores for each service its name and the calls that it makes (strings of URLs / method names)
type ServiceCalls struct {
	Service string
	Calls   map[string]int
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
