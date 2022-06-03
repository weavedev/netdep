package callanalyzer

import "golang.org/x/tools/go/ssa"

// Frame is a struct for keeping track of the traversal packages while looking for interesting functions
type Frame struct {
	// trace is a list of previous calls.
	trace []*ssa.Call
	// visited is shared between frames and keeps track of which nodes have been visited
	// to prevent repetitive visits.
	visited map[*ssa.Call]bool
	pkg     *ssa.Package

	// params maps a parameter inside a function to a argument value given in another frame
	params map[*ssa.Parameter]*ssa.Value

	// parent is necessary to recursively resolve variables (in different scopes)
	parent            *Frame
	targetsCollection *TargetsCollection
}

// hasVisited returns whether the block has already been trace.
func (f Frame) hasVisited(call *ssa.Call) bool {
	_, visited := f.visited[call]

	for _, called := range f.trace {
		if called == call {
			return true
		}
	}

	return visited
}
