/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package callanalyzer

import (
	"fmt"
	"go/token"
	"os"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ssa"
)

// CallTarget holds information about a certain call made by the analysed package.
// This used to be named "Caller" (which was slightly misleading, as it is in fact the target,
// thus rather a 'callee' than a 'caller'.
type CallTarget struct {
	// The name of the package the method belongs to
	packageName string
	// The name of the call (i.e. name of function or some other target)
	MethodName string
	// The URL of the entity
	RequestLocation string
	// The name of the service in which the call is made
	ServiceName string
	// The name of the file in which the call is made
	FileName string
	// The line number in the file where the call is made
	PositionInFile string
}

// findFunctionInPackage finds the method by its name within the specified package.
// Important: it only looks for Exported functions
func findFunctionInPackage(pkg *ssa.Package, name string) *ssa.Function {
	// Find member
	member, hasSpecifiedMember := pkg.Members[name]
	if !hasSpecifiedMember {
		return nil
	}
	// Check that member is a Function
	foundFunction, ok := member.(*ssa.Function)
	if !ok {
		// Not a function
		return nil
	}
	return foundFunction
}

// getCallInformation returns the service, file and line number
// of a discovered call
//
// Arguments:
// pos is the position of the call
// frame is a structure for keeping track of the recursion and package
func getCallInformation(pos token.Pos, pkg *ssa.Package) (string, string, string) {
	// split package name and take the last item to get the service name
	service := pkg.String()[strings.LastIndex(pkg.String(), "/")+1:]

	// absolute file path
	filePath := pkg.Prog.Fset.File(pos).Name()
	// split absolute path to get the relative file path from the service directory
	parts := filePath[strings.LastIndex(filePath, string(os.PathSeparator)+service+string(os.PathSeparator))+1:]

	base := 10
	// take the position of the call within the file and convert to string
	position := strconv.FormatInt(int64(pkg.Prog.Fset.Position(pos).Line), base)

	return service, parts, position
}

// analyseCall recursively traverses the SSA, with call being the starting point,
// and using the environment specified in the frame
// Variables are only resolved if the call is 'interesting'
// Recursion is only continued if the call is not in the 'ignoreList'
//
// Arguments:
// call is the call to analyse,
// frame is a structure for keeping track of the recursion,
// config specifies how the analyser should behave, and
// targets is a reference to the ultimate data structure that is to be completed and returned.
func analyseCall(call *ssa.Call, frame *Frame, config *AnalyserConfig, targetsClient *[]*CallTarget, targetsServer *[]*CallTarget) {
	// The function call type can either be a *ssa.Function, an anonymous function type, or something else,
	// hence the switch. See https://pkg.go.dev/golang.org/x/tools/go/ssa#Call for all possibilities
	switch fnCallType := call.Call.Value.(type) {
	// TODO: handle other cases
	case *ssa.Function:
		// Qualified function name is: package + interface + function
		qualifiedFunctionNameOfTarget := fnCallType.RelString(nil)
		// .Pkg returns an obj of type *ssa.Package, whose .Pkg returns one of *type.Package
		// This is therefore not the grandparent package, but the *type.Package of the fnCall
		calledFunctionPackage := fnCallType.Pkg.Pkg.Path() // e.g. net/http

		if strings.Contains(qualifiedFunctionNameOfTarget, "Client") || strings.Contains(qualifiedFunctionNameOfTarget, ".Do") {
			fmt.Println("Client.Do: " + qualifiedFunctionNameOfTarget)
		}

		interestingStuffClient, isInterestingClient := config.interestingCallsClient[qualifiedFunctionNameOfTarget]
		if isInterestingClient {
			// TODO: Resolve the arguments of the function call
			handleInterestingClientCall(call, interestingStuffClient, calledFunctionPackage, qualifiedFunctionNameOfTarget, targetsClient, frame)
			return
		}

		interestingStuffServer, isInterestingServer := config.interestingCallsServer[qualifiedFunctionNameOfTarget]
		if isInterestingServer {
			// TODO: Resolve the arguments of the function call
			handleInterestingServerCall(call, interestingStuffServer, calledFunctionPackage, qualifiedFunctionNameOfTarget, targetsServer, frame)
			return
		}

		_, isIgnored := config.ignoreList[calledFunctionPackage]

		if isIgnored {
			// Do not recurse into the packageName if it is ignored
			return
		}
		// The following creates a copy of 'frame'.
		// This is the correct place for this because we are going to visit child blocks next.
		newFrame := *frame

		for _, arg := range call.Call.Args {
			if fn, ok := arg.(*ssa.Call); ok {
				analyseCall(fn, frame, config, targetsClient, targetsServer)
			}
			if fn, ok := arg.(*ssa.Function); ok {
				visitBlocks(fn.Blocks, &newFrame, config, targetsClient, targetsServer)
			}
		}

		if fnCallType.Blocks != nil {
			visitBlocks(fnCallType.Blocks, &newFrame, config, targetsClient, targetsServer)
		}
	default:
		// Unsupported call type
		return
	}
}

func handleInterestingServerCall(call *ssa.Call, interestingStuffServer InterestingCall, calledFunctionPackage string, qualifiedFunctionNameOfTarget string, targetsServer *[]*CallTarget, frame *Frame) {
	//nolint:nestif
	if interestingStuffServer.action == Output {
		requestLocation := ""
		if call.Call.Args != nil && len(interestingStuffServer.interestingArgs) > 0 {
			if qualifiedFunctionNameOfTarget == "(*github.com/gin-gonic/gin.Engine).Run" {
				requestLocation = strings.Join(resolveGinAddrSlice(call.Call.Args[1]), "")
			} else {
				requestLocation = strings.Join(resolveVariables(call.Call.Args, interestingStuffServer.interestingArgs), "")
			}
		}
		// Additional information about the call
		service, file, position := getCallInformation(call.Pos(), frame.pkg)

		callTarget := &CallTarget{
			packageName:     calledFunctionPackage,
			MethodName:      qualifiedFunctionNameOfTarget,
			RequestLocation: requestLocation,
			ServiceName:     service,
			FileName:        file,
			PositionInFile:  position,
		}

		// fmt.Println("Found call to function " + qualifiedFunctionNameOfTarget)

		*targetsServer = append(*targetsServer, callTarget)
		return
	} else if interestingStuffServer.action == Substitute {
		// TODO: implement substitution of env calls
	}
}

func handleInterestingClientCall(call *ssa.Call, interestingStuffClient InterestingCall, calledFunctionPackage string, qualifiedFunctionNameOfTarget string, targetsClient *[]*CallTarget, frame *Frame) {
	if interestingStuffClient.action == Output {
		requestLocation := ""
		if call.Call.Args != nil && len(interestingStuffClient.interestingArgs) > 0 {
			requestLocation = strings.Join(resolveVariables(call.Call.Args, interestingStuffClient.interestingArgs), "")
		}
		// Additional information about the call
		service, file, position := getCallInformation(call.Pos(), frame.pkg)

		callTarget := &CallTarget{
			packageName:     calledFunctionPackage,
			MethodName:      qualifiedFunctionNameOfTarget,
			RequestLocation: requestLocation,
			ServiceName:     service,
			FileName:        file,
			PositionInFile:  position,
		}

		// fmt.Println("Found call to function " + qualifiedFunctionNameOfTarget)

		*targetsClient = append(*targetsClient, callTarget)
		return
	} else if interestingStuffClient.action == Substitute {
		// TODO: implement substitution of env calls
	}
}

// analyseInstructionsOfBlock checks the type of each iteration in a block.
// If it finds a call, it analysed it to check if it is interesting.
//
// Arguments:
// blocks is the array of blocks to analyse,
// fr keeps track of the traversal,
// config specifies the behaviour of the analyser,
// targets is a reference to the ultimate data structure that is to be completed and returned.
func analyseInstructionsOfBlock(block *ssa.BasicBlock, fr *Frame, config *AnalyserConfig, targetsClient *[]*CallTarget, targetsServer *[]*CallTarget) {
	if block.Instrs == nil {
		return
	}

	for _, instr := range block.Instrs {
		switch instruction := instr.(type) {
		case *ssa.Call:
			analyseCall(instruction, fr, config, targetsClient, targetsServer)
		default:
			continue
		}
	}
}

// visitBlocks visits each of the blocks in the specified 'blocks' list and analyses each of the block's instructions.
//
// Arguments:
// blocks is the array of blocks to analyse,
// fr keeps track of the traversal,
// config specifies the behaviour of the analyser,
// targets is a reference to the ultimate data structure that is to be completed and returned.
func visitBlocks(blocks []*ssa.BasicBlock, fr *Frame, config *AnalyserConfig, targetsClient *[]*CallTarget, targetsServer *[]*CallTarget) {
	if len(fr.visited) > config.maxTraversalDepth {
		// fmt.Println("Traversal defaultMaxTraversalDepth is more than 16; terminate this recursion branch")
		return
	}

	for _, block := range blocks {
		if fr.hasVisited(block) || block == nil {
			continue
		}
		newFr := fr
		// Mark the block as visited
		newFr.visited[block] = true
		analyseInstructionsOfBlock(block, newFr, config, targetsClient, targetsServer)
	}
}

// AnalysePackageCalls takes a main package and finds all 'interesting' methods that are called
//
// Arguments:
// pkg is the package to analyse
// config specifies the behaviour of the analyser,
//
// Returns:
// List of pointers to callTargets, or an error if something went wrong.
func AnalysePackageCalls(pkg *ssa.Package, config *AnalyserConfig) ([]*CallTarget, []*CallTarget, error) {
	mainFunction := findFunctionInPackage(pkg, "main")

	// TODO: look for the init function will be useful if we want to know
	// the values of global file-scoped variables
	// initFunction := findFunctionInPackage(pkg, "init")

	// Find the main function
	if mainFunction == nil {
		return nil, nil, fmt.Errorf("no main function found in package %v", pkg)
	}

	baseFrame := Frame{
		visited: make(map[*ssa.BasicBlock]bool, 0),
		// Reference to the final list of all _targets of the entire package
		pkg: pkg,
	}

	targetsClient := make([]*CallTarget, 0)
	targetsServer := make([]*CallTarget, 0)

	// Visit each of the block of the main function
	visitBlocks(mainFunction.Blocks, &baseFrame, config, &targetsClient, &targetsServer)

	return targetsClient, targetsServer, nil
}
