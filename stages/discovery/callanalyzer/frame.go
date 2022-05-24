package callanalyzer

import "golang.org/x/tools/go/ssa"

// Frame is a struct for keeping track of the traversal packages while looking for interesting functions
type Frame struct {
	// visited is a set of blocks that have been visited.
	visited map[*ssa.BasicBlock]bool
	// params maps a parameter inside a function to a argument value given in another frame
	params map[*ssa.Parameter]*ssa.Value
	pkg    *ssa.Package
	prog   *ssa.Program
	// parent is necessary to recursively resolve variables (in different scopes)
	parent            *Frame
	targetsCollection *TargetsCollection
}

// hasVisited returns whether the block has already been visited.
func (f Frame) hasVisited(block *ssa.BasicBlock) bool {
	_, ok := f.visited[block]
	return ok
}
