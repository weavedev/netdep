package callanalyzer

import "golang.org/x/tools/go/ssa"

// Frame is a struct for keeping track of the traversal packages while looking for interesting functions
type Frame struct {
	trace             []*ssa.CallCommon             // trace is a stack trace of previous calls.
	visited           map[*ssa.CallCommon]bool      // visited is shared between frames and keeps track of which nodes have been visited
	params            map[*ssa.Parameter]*ssa.Value // params maps a parameter inside a function to a argument value given in another frame
	globals           map[*ssa.Global]*ssa.Value    // globals keeps a map the values associated with global variables
	pkg               *ssa.Package                  // pkg references the service package
	parent            *Frame                        // parent is necessary to recursively resolve variables (in different scopes)
	targetsCollection *TargetsCollection            // targetsCollection is a reference to the collection of found calls
	singlePass        bool                          // singlePass defines if we should check visited or trace for performance
	pointerMap        map[*ssa.CallCommon]*ssa.Function
}

// hasVisited returns whether the block has already been trace.
func (f Frame) hasVisited(call *ssa.CallCommon) bool {
	interesting, visited := f.visited[call]
	if visited && !interesting {
		return true
	}

	for _, callee := range f.trace {
		if callee == call {
			return true
		}
	}

	return false
}
