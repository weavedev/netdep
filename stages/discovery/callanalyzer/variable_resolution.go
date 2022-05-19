/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/

//nolint:gocritic,exhaustive
package callanalyzer

import (
	"go/constant"
	"go/token"
	"golang.org/x/tools/go/ssa"
)

func resolveVariable(value ssa.Value) (string, bool) {
	switch val := value.(type) {
	case *ssa.Parameter:
		return "unknown: the parameter was not resolved", false
	case *ssa.BinOp:
		switch val.Op {
		case token.ADD:
			left, isLeftResolved := resolveVariable(val.X)
			right, isRightResolved := resolveVariable(val.Y)
			if isRightResolved && isLeftResolved {
				return left + right, true
			}

			return left + right, false
		}
		return "unknown: only string concatenation can be resolved", false
	case *ssa.Const:
		switch val.Value.Kind() {
		case constant.String:
			return constant.StringVal(val.Value), true
		}
		return "unknown: constant is not a string", false
	}

	return "unknown: var(" + value.Name() + ") = ??", false
}

func resolveVariables(parameters []ssa.Value, positions []int) ([]string, bool) {
	stringParameters := make([]string, len(positions))
	var isResolved = true

	for i, idx := range positions {
		if idx < len(parameters) {
			variable, isResolved := resolveVariable(parameters[idx])
			if isResolved {
				stringParameters[i] = variable
			} else {
				isResolved = false
			}
		}
	}

	return stringParameters, isResolved
}

// resolveGinAddrSlice is a hardcoded solution to resolve the port address of a Run command from the "github.com/gin-gonic/gin" library
// TODO: implement a general way for resolving variables in slices
func resolveGinAddrSlice(value ssa.Value) ([]string, bool) {
	switch val := value.(type) {
	case *ssa.Slice:
		switch val1Type := val.X.(type) {
		case *ssa.Alloc:
			block := val1Type.Block()
			for i := range block.Instrs {
				// Iterate through instruction of the block in reverse order
				// In the case of the Gin library, the last Store instruction in the block contains the address value we're looking for
				switch instruction := block.Instrs[len(block.Instrs)-1-i].(type) {
				case *ssa.Store:
					switch storeValType := instruction.Val.(type) {
					case *ssa.Const:
						switch storeValType.Value.Kind() {
						case constant.String:
							return []string{constant.StringVal(storeValType.Value)}, true
						}
					}
				}
			}
		}
	case *ssa.Const:
		return []string{":8080"}, true
	}
	return []string{"var(" + value.Name() + ") = ??"}, false
}
