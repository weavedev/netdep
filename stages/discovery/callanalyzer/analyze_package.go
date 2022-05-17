/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package callanalyzer

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/ssa"
)

// locationIdx Stores Relevant Libraries
// their Relevant Methods and for each method
// a position of location in the Args of ssa.Call

//nolint
var (
	// ignoreList is a set of function names to not recurse into
	ignoreList = map[string]bool{
		"runtime":  true,
		"fmt":      true,
		"reflect":  true,
		"sync":     true,
		"internal": true,
		"syscall":  true,
		"unicode":  true,
		"time":     true,
	}
	locationIdx = map[string]map[string][]int{
		"net/http": {
			"Get":      []int{0, 1},
			"Post":     []int{0, 1},
			"Put":      []int{0, 1},
			"PostForm": []int{0, 1},
			"Head":     []int{0, 1},
			// "Do":                    []int{0},
			"NewRequest":            []int{1},
			"NewRequestWithContext": []int{2},
			// this is a bit different, as it uses http.Request
			// as an argument. This will be completed in the future.
		},
	}
	locationIdxServer = map[string]map[string][]int{
		"net/http": {
			"Handle":         []int{0},
			"HandleFunc":     []int{0},
			"ListenAndServe": []int{0},
		},
		"github.com/gin-gonic/gin": {
			"GET":     []int{1},
			"POST":    []int{1},
			"PUT":     []int{1},
			"DELETE":  []int{1},
			"PATCH":   []int{1},
			"HEAD":    []int{1},
			"OPTIONS": []int{1},
			"Run":     []int{1},
		},
	}
)

type Caller struct {
	requestLocation string
	library         string
	methodName      string
	// TODO: Add package name, filename, code line
}

func getMainFunction(pkg *ssa.Package) *ssa.Function {
	mainMember, hasMain := pkg.Members["main"]
	if !hasMain {
		return nil
	}
	mainFunction, ok := mainMember.(*ssa.Function)
	if !ok {
		return nil
	}

	return mainFunction
}

func discoverCall(call *ssa.Call) (*Caller, bool) {
	var caller *Caller
	var server bool

	//nolint
	switch call.Call.Value.(type) {
	case *ssa.Function:
		calledFunction, _ := call.Call.Value.(*ssa.Function)
		calledFunctionPackage := calledFunction.Pkg.Pkg.Path()

		relevantPackage, isRelevantPackage := locationIdx[calledFunctionPackage]
		relevantPackageServer, isRelevantPackageServer := locationIdxServer[calledFunctionPackage]

		if isRelevantPackage {
			indices, isRelevantFunction := relevantPackage[calledFunction.Name()]
			if call.Call.Args != nil && isRelevantFunction {
				arguments := resolveVariables(call.Call.Args, indices)
				caller = &Caller{
					requestLocation: strings.Join(arguments, ""),
					library:         calledFunctionPackage,
					methodName:      calledFunction.Name(),
				}
				return caller, server
			}
		}

		if isRelevantPackageServer {
			indices, isRelevantFunction := relevantPackageServer[calledFunction.Name()]
			if call.Call.Args != nil && isRelevantFunction {
				var arguments []string
				// Temporary hardcoded solution for resolving the port argument in Run command of the gin library
				if calledFunctionPackage == "github.com/gin-gonic/gin" && calledFunction.Name() == "Run" {
					arguments = resolveGinAddrSlice(call.Call.Args[1])
				} else {
					arguments = resolveVariables(call.Call.Args, indices)
				}

				caller = &Caller{
					requestLocation: strings.Join(arguments, "/"),
					library:         calledFunctionPackage,
					methodName:      calledFunction.Name(),
				}
				server = true
				return caller, server
			}
		}

		_, isIgnored := ignoreList[calledFunctionPackage]

		if calledFunction.Blocks != nil && !isIgnored {
			discoverBlocks(calledFunction.Blocks)
		}
		return caller, server
	default:
		return nil, server
	}
}

func discoverBlock(block *ssa.BasicBlock) ([]*Caller, []*Caller) {
	if block.Instrs == nil {
		return nil, nil
	}

	var clientCalls []*Caller
	var serverCalls []*Caller

	for _, instr := range block.Instrs {
		//nolint // can't rewrite switch with 1 case into if,
		// because .(type) is not allowed outside switch.
		switch instruction := instr.(type) {
		// Every complex Instruction is split into several instructions
		// so even if the call is part of variable assignment
		// or a loop it will be stored as a separate ssa.Call instruction
		case *ssa.Call:
			if call, server := discoverCall(instruction); server {
				serverCalls = append(serverCalls, call)
			} else {
				clientCalls = append(clientCalls, call)
			}
		}
	}

	return clientCalls, serverCalls
}

func discoverBlocks(blocks []*ssa.BasicBlock) ([]*Caller, []*Caller) {
	var clientCalls []*Caller
	var serverCalls []*Caller

	for _, block := range blocks {
		client, server := discoverBlock(block)
		clientCalls, serverCalls = append(clientCalls, client...), append(serverCalls, server...)
	}

	return clientCalls, serverCalls
}

func AnalyzePackageCalls(pkg *ssa.Package) ([]*Caller, []*Caller, error) {
	mainFunction := getMainFunction(pkg)

	if mainFunction == nil {
		return nil, nil, fmt.Errorf("no main function found")
	}

	clientCalls, serverCalls := discoverBlocks(mainFunction.Blocks)

	return clientCalls, serverCalls, nil
}
