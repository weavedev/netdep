package callanalyzer

import "golang.org/x/tools/go/ssa"

// Frame is a struct for keeping track of the traversal packages while looking for interesting functions
type Frame struct {
	// trace is a stack trace of previous calls.
	trace []*ssa.Call
	// visited is shared between frames and keeps track of which nodes have been visited
	// to prevent repetitive visits.
	visited map[*ssa.Call]bool
	// params maps a parameter inside a function to a argument value given in another frame
	params map[*ssa.Parameter]*ssa.Value
	// globals keeps a map the values associated with global variables
	globals map[*ssa.Global]*ssa.Value
	pkg     *ssa.Package
	// parent is necessary to recursively resolve variables (in different scopes)
	parent            *Frame
	targetsCollection *TargetsCollection
	// singlePass defines if we should check visited or trace for performance
	singlePass bool
}

// hasVisited returns whether the block has already been trace.
func (f Frame) hasVisited(call *ssa.Call) bool {
	if f.singlePass {
		_, visited := f.visited[call]
		return visited
	} else {
		for _, callee := range f.trace {
			if callee == call {
				return true
			}
		}

		return false
	}
}
