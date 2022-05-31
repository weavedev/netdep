package callanalyzer

import (
	"fmt"

	"golang.org/x/tools/go/ssa"
)

func PrintTraceToCall(trace []*ssa.Call, frame *Frame, config *AnalyserConfig) {
	traces := len(trace)

	// displayStartingIndex is used to limit the number of traces to log
	// as defined by config.maxTraceDepth
	displayStartingIndex := traces - config.maxTraceDepth

	for i, call := range trace {
		// only show the last few items, as defined in the config
		if i < displayStartingIndex {
			continue
		}

		_, file, position := getCallInformation(call.Pos(), frame.pkg)
		signature := ""
		switch callee := call.Call.Value.(type) {
		case *ssa.Function:
			signature = callee.RelString(nil)
		default:
			signature = call.Call.Value.String()
		}
		fmt.Printf("%d: %s:%s\t%s\n", traces-i, file, position, signature)
	}
}
