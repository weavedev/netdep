package analyze

import "golang.org/x/tools/go/ssa"

type frame struct {
	visited  []*ssa.BasicBlock
	mappings map[string]ssa.Value
}

func (f frame) hasVisited(block *ssa.BasicBlock) bool {
	for _, b := range f.visited {
		if block == b {
			return true
		}
	}

	return false
}
