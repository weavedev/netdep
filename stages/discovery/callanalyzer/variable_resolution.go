/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/

package callanalyzer

import (
	"go/constant"
	"go/token"

	"golang.org/x/tools/go/ssa"
)

// resolveParameter resolves a parameter in a frame, recursively
func resolveParameter(par *ssa.Parameter, fr *Frame) (*ssa.Value, *Frame) {
	if fr != nil {
		parameterValue, hasParam := fr.params[par]
		if hasParam {
			recPar, isParam := (*parameterValue).(*ssa.Parameter)

			if isParam {
				return resolveParameter(recPar, fr.parent)
			} else {
				return parameterValue, fr
			}
		}
	}

	return nil, fr
}

// resolveVariable returns a string value of ssa.Value
// if the value can be resolved. It also returns a bool
// which indicates whether the variable was resolved.
func resolveVariable(value *ssa.Value, fr *Frame, config *AnalyserConfig) (string, bool) {
	switch val := (*value).(type) {
	case *ssa.Parameter:
		// (recursively) resolve a parameter to a value and return that value, if it is defined
		parValue, resFrame := resolveParameter(val, fr)

		if parValue != nil {
			return resolveVariable(parValue, resFrame, config)
		}

		return "unknown: the parameter was not resolved", false
	case *ssa.BinOp:
		switch val.Op { //nolint:exhaustive
		case token.ADD:
			left, isLeftResolved := resolveVariable(&val.X, fr, config)
			right, isRightResolved := resolveVariable(&val.Y, fr, config)
			if isRightResolved && isLeftResolved {
				return left + right, true
			}

			return left + right, false
		default:
			return "unknown: only ADD binary operation is supported", false
		}
	case *ssa.Const:
		switch val.Value.Kind() { //nolint:exhaustive
		case constant.String:
			return constant.StringVal(val.Value), true
		default:
			return "unknown: not a string constant", false
		}
	case *ssa.Call:
		// TODO: here shall the substitution happen
		if config.interestingCallsCommon["To be call value"].action == Substitute {
			return "unknown: interesting call that could be substituted (currently not implemented)", true
		}
		return "unknown: interesting call that is not supported", false

	default:
		return "unknown: unable to resolve", false
	}
}

// resolveParameters iterates over the parameters, resolving those where possible.
// It also keeps track of whether all variables could be resolved or not.
func resolveParameters(parameters []ssa.Value, positions []int, fr *Frame, config *AnalyserConfig) ([]string, bool) {
	stringParameters := make([]string, len(positions))
	wasResolved := true

	for i, idx := range positions {
		if idx < len(parameters) {
			variable, isResolved := resolveVariable(&parameters[idx], fr, config)
			if isResolved {
				stringParameters[i] = variable
			} else {
				wasResolved = false
			}
		}
	}

	return stringParameters, wasResolved
}

// resolveGinAddrSlice is a hardcoded solution to resolve the port address of a Run command from the "github.com/gin-gonic/gin" library.
// Returns list of strings that represent the slice, and bool value indicating whether the variable was resolved.
// TODO: implement a general way for resolving variables in slices
func resolveGinAddrSlice(value ssa.Value) ([]string, bool) {
	unresolvedType := []string{"unknown: var(" + value.Name() + ") = ??"}
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
						switch storeValType.Value.Kind() { //nolint:exhaustive
						case constant.String:
							return []string{constant.StringVal(storeValType.Value)}, true
						default:
							return unresolvedType, false
						}
					default:
						return unresolvedType, false
					}
				default:
					return unresolvedType, false
				}
			}
		default:
			return unresolvedType, false
		}
	case *ssa.Const:
		return []string{":8080"}, true
	}
	return unresolvedType, false
}
