package callanalyzer

import "golang.org/x/tools/go/ssa"

type Frame struct {
	visited  []*ssa.BasicBlock
	Mappings map[string]ssa.Value
	OldFrame *Frame
}

func (f Frame) hasVisited(block *ssa.BasicBlock) bool {
	for _, b := range f.visited {
		if block == b {
			return true
		}
	}

	return false
}
