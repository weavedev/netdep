package callanalyzer

import (
	"go/types"

	"golang.org/x/tools/go/ssa"
)

// getInvokedFunctionFromCall takes a call and tries to determine which function is being invoked
// TODO resolve missing cases
func getInvokedFunctionFromCall(call *ssa.CallCommon, frame *Frame) *ssa.Function {
	pkg := call.Method.Pkg()
	name := call.Method.Name()
	program := frame.pkg.Prog

	callValue := &call.Value

	// resolve call to parameter
	if param, isParam := (*callValue).(*ssa.Parameter); isParam {
		parValue, _ := resolveParameter(param, frame)
		if parValue != nil {
			callValue = parValue
		}
	}

	// resolve type of value
	var methodType types.Type
	switch expr := (*callValue).(type) {
	case *ssa.MakeInterface:
		methodType = expr.X.Type()
	default:
		// methodType = expr.Type()
		methodType = call.Method.Type()
	}

	// lookup method of type
	methodSet := program.MethodSets.MethodSet(methodType)
	sel := methodSet.Lookup(pkg, name)

	if sel != nil {
		return program.MethodValue(sel)
	} else {
		return nil
	}
}

// getCallFunctionFromCall returns the function being call for static calls
func getCallFunctionFromCall(call *ssa.CallCommon, frame *Frame) *ssa.Function {
	// resolve parameter
	if param, isParam := call.Value.(*ssa.Parameter); isParam {
		parValue, _ := resolveParameter(param, frame)
		if paramFn, isFn := (*parValue).(*ssa.Function); isFn {
			// TODO: does this happen?
			return paramFn
		}
	}

	// helper function that handles a few cases
	return call.StaticCallee()
}

// getFunctionFromCall returns the function being called by either an invocation or a static call
func getFunctionFromCall(call *ssa.CallCommon, frame *Frame) *ssa.Function {
	if call.IsInvoke() {
		return getInvokedFunctionFromCall(call, frame)
	} else {
		return getCallFunctionFromCall(call, frame)
	}
}
