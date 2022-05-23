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
	"strings"
)

// resolveVariable Resolves a supplied variable value, only in the cases that are supported by the tool:
//
// - string concatenation (see BinOp),
//
// - string literal
//
// - call to os.GetEnv
//
// - other InterestingCalls with the action Substitute
func resolveVariable(value ssa.Value, config *AnalyserConfig) string {
	switch val := value.(type) {
	case *ssa.Parameter:
		return "[[Unknown]]"
	case *ssa.BinOp:
		switch val.Op {
		case token.ADD:
			return resolveVariable(val.X, config) + resolveVariable(val.Y, config)
		}
		return "[[OP]]"
	case *ssa.Const:
		switch val.Value.Kind() {
		case constant.String:
			return constant.StringVal(val.Value)
		}
		return "[[CONST]]"
	case *ssa.Call:
		// TODO: here shall the substitution happen
		if config.interestingCallsCommon["To be call value"].action == Substitute {
			return "[[Interesting Call]]"
		}
		return "[[CALL]]"
	}
	return "var(" + value.Name() + ") = ??"
}

// resolveParameters resolves each of the parameters supplied to a function.
func resolveParameters(parameters []ssa.Value, positions []int, config *AnalyserConfig) []string {
	stringParameters := make([]string, len(positions))
	for i, idx := range positions {
		if idx < len(parameters) {
			variable := resolveVariable(parameters[idx], config)
			if !strings.HasPrefix(variable, "not a constant") {
				stringParameters[i] = variable
			}
		}
	}

	return stringParameters
}

// resolveGinAddrSlice is a hardcoded solution to resolve the port
// address of a Run command from the "github.com/gin-gonic/gin" library.
// TODO: implement a general way for resolving variables in slices
func resolveGinAddrSlice(value ssa.Value) []string {
	switch val := value.(type) {
	case *ssa.Slice:
		switch val1Type := val.X.(type) {
		case *ssa.Alloc:
			block := val1Type.Block()
			for i := range block.Instrs {
				// Iterate through instruction of the block in reverse order
				// In the case of the Gin library, the last Store instruction
				// in the block contains the address value we're looking for
				switch instruction := block.Instrs[len(block.Instrs)-1-i].(type) {
				case *ssa.Store:
					switch storeValType := instruction.Val.(type) {
					case *ssa.Const:
						switch storeValType.Value.Kind() {
						case constant.String:
							return []string{constant.StringVal(storeValType.Value)}
						}
					}
				}
			}
		}
	case *ssa.Const:
		return []string{":8080"}
	}
	return []string{"var(" + value.Name() + ") = ??"}
}
