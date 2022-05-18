package callanalyzer

import "golang.org/x/tools/go/ssa"

// Frame is a struct for keeping track of the traversal packages while looking for interesting functions
type Frame struct {
	// visited is a set of blocks that have been visited.
	visited map[*ssa.BasicBlock]bool
}

// hasVisited returns whether the block has already been visited.
func (f Frame) hasVisited(block *ssa.BasicBlock) bool {
	_, ok := f.visited[block]
	return ok
}
