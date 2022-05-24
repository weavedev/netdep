package callanalyzer

import (
	"fmt"
	"golang.org/x/tools/go/ssa"
)

func PrintTraceToCall(trace []*ssa.Call, frame *Frame, config *AnalyserConfig) {
	traces := len(trace)
	minOffset := traces - config.maxTraceDepth

	for i, call := range trace {
		if i < minOffset {
			continue
		}

		_, file, position := getCallInformation(call.Pos(), frame.pkg)
		signature := call.Call.Signature().String()
		fmt.Printf("%d: %s:%s\t%s\n", traces-i, file, position, signature)
	}
}
