// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

/*
Copyright © 2022 TW Group 13C, Weave BV, TU Delft

In the Building stages, the adjacency lists of each service are populated.
This is done by traversing the lists of endpoints/clients and looking for the other end of the connection.
The Building stages should handle
Refer to the Project plan, chapter 5.4 for more information.
*/

// ConstructOutput
// is an example method for the building stage that could, for instance,
// json.Marshal/serialise a nice data structure into a string.
//
// TODO: Remove the following line when implementing this method
//goland:noinspection GoUnusedParameter
func ConstructOutput(dataStructure interface{}) string {
	return "{ type: \"error\", message: \"Not Implemented\" }"
}
