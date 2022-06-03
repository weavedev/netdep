package callanalyzer

import (
	"fmt"

	"golang.org/x/tools/go/ssa"
)

func PrintTraceToCall(frame *Frame, config *AnalyserConfig) {
	traces := len(frame.visited)

	// displayStartingIndex is used to limit the number of traces to log
	// as defined by config.maxTraceDepth
	displayStartingIndex := traces - config.maxTraceDepth

	// could use max(displayStartingIndex, 0), but this is more explicit.
	if displayStartingIndex < 0 {
		displayStartingIndex = 0
	}

	for i, call := range frame.trace[displayStartingIndex:] {
		file, position := getPositionFromPos(call.Pos(), frame.pkg.Prog)
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
