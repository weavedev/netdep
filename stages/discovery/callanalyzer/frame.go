package callanalyzer

import "golang.org/x/tools/go/ssa"

type Frame struct {
	visited map[*ssa.BasicBlock]bool
}

func (f Frame) hasVisited(block *ssa.BasicBlock) bool {
	_, ok := f.visited[block]
	return ok
}
