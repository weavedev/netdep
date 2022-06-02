package callanalyzer

import "golang.org/x/tools/go/ssa"

// Frame is a struct for keeping track of the traversal packages while looking for interesting functions
type Frame struct {
	// visited is a set of blocks that have been visited.
	visited []*ssa.Call
	// params maps a parameter inside a function to a argument value given in another frame
	params map[*ssa.Parameter]*ssa.Value
	// globals keeps a map the values associated with global variables
	globals map[*ssa.Global]*ssa.Value
	pkg     *ssa.Package
	// parent is necessary to recursively resolve variables (in different scopes)
	parent            *Frame
	targetsCollection *TargetsCollection
}

// hasVisited returns whether the block has already been visited.
func (f Frame) hasVisited(instruction *ssa.Call) bool {
	for _, instr := range f.visited {
		if instr == instruction {
			return true
		}
	}
	return false
}
