/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package callanalyzer

import (
	"fmt"
	"go/token"
	"go/types"
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
	// A flag describing whether the RequestLocation was resolved
	IsResolved bool
	// The name of the service in which the call is made
	ServiceName string
	// The name of the file in which the call is made
	FileName string
	// The line number in the file where the call is made
	PositionInFile string
}

// TargetsCollection holds the output structures that are to be returned by the
// discovery stage
type TargetsCollection struct {
	clientTargets []*CallTarget
	serverTargets []*CallTarget
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
func analyseCall(call *ssa.Call, frame *Frame, config *AnalyserConfig) {
	var fnCallType *ssa.Function

	if call.Call.IsInvoke() {
		// TODO: resolve a call to a method
		pkg := call.Call.Method.Pkg()
		name := call.Call.Method.Name()
		prog := frame.prog
		var mtype types.Type
		switch fnCallType := call.Call.Value.(type) {
		case *ssa.Parameter:
			parValue, _ := resolveParameter(fnCallType, frame)
			if parValue != nil {
				switch expr := (*parValue).(type) {
				case *ssa.MakeInterface:
					mtype = expr.X.Type()
				}
			}
			break
		default:
			break
		}

		if mtype == nil {
			mtype = call.Call.Method.Type()
		}

		mset := prog.MethodSets.MethodSet(mtype)
		sel := mset.Lookup(pkg, name)
		if sel != nil {
			fnCallType = prog.MethodValue(sel)
		}
	} else {
		fnCallType = call.Call.StaticCallee()
		if param, isParam := call.Call.Value.(*ssa.Parameter); isParam && fnCallType == nil {
			parValue, _ := resolveParameter(param, frame)
			if paramFn, isFn := (*parValue).(*ssa.Function); isFn {
				fnCallType = paramFn
			}
			return
		}
	}

	if fnCallType == nil {
		return
	}

	// Qualified function name is: package + interface + function
	// TODO: handle parameter equivalence to other interface
	qualifiedFunctionNameOfTarget := fnCallType.RelString(nil)
	// .Pkg returns an obj of type *ssa.Package, whose .Pkg returns one of *type.Package
	// This is therefore not the grandparent package, but the *type.Package of the fnCall
	calledFunctionPackage := fnCallType.Pkg.Pkg.Path() // e.g. net/http

	_, isInterestingClient := config.interestingCallsClient[qualifiedFunctionNameOfTarget]
	if isInterestingClient {
		// TODO: Resolve the arguments of the function call
		handleInterestingClientCall(call, config, calledFunctionPackage, qualifiedFunctionNameOfTarget, frame)
		return
	}

	_, isInterestingServer := config.interestingCallsServer[qualifiedFunctionNameOfTarget]
	if isInterestingServer {
		// TODO: Resolve the arguments of the function call
		handleInterestingServerCall(call, config, calledFunctionPackage, qualifiedFunctionNameOfTarget, frame)
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

	// Keep track of given parameters for resolving
	for i, par := range fnCallType.Params {
		newFrame.params[par] = &call.Call.Args[i]
	}

	// Keep a reference to the parent frame
	newFrame.parent = frame

	if fnCallType.Blocks != nil {
		visitBlocks(fnCallType.Blocks, &newFrame, config)
	}
}

// handleInterestingServerCall collects the information about a supplied endpoint declaration
// and adds this information to the targetsServer data structure. If possible, also calls the function to resolve
// the parameters of the function call.
func handleInterestingServerCall(call *ssa.Call, config *AnalyserConfig, packageName, qualifiedFunctionNameOfTarget string, frame *Frame) {
	interestingStuffServer := config.interestingCallsServer[qualifiedFunctionNameOfTarget]
	if interestingStuffServer.action == Output {
		requestLocation := ""
		var isResolved bool
		var variables []string

		if call.Call.Args != nil && len(interestingStuffServer.interestingArgs) > 0 {
			if qualifiedFunctionNameOfTarget == "(*github.com/gin-gonic/gin.Engine).Run" {
				variables, isResolved = resolveGinAddrSlice(call.Call.Args[1])
				// TODO: parse the url
				requestLocation = strings.Join(variables, "")
			} else {
				variables, isResolved = resolveParameters(call.Call.Args, interestingStuffServer.interestingArgs, frame, config)
				// TODO: parse the url
				requestLocation = strings.Join(variables, "")
			}
		}

		// Additional information about the call
		service, file, position := getCallInformation(call.Pos(), frame.pkg)

		callTarget := &CallTarget{
			packageName:     packageName,
			MethodName:      qualifiedFunctionNameOfTarget,
			RequestLocation: requestLocation,
			IsResolved:      isResolved,
			ServiceName:     service,
			FileName:        file,
			PositionInFile:  position,
		}

		// fmt.Println("Found call to function " + qualifiedFunctionNameOfTarget)

		frame.targetsCollection.serverTargets = append(frame.targetsCollection.serverTargets, callTarget)
		return
	}
}

// handleInterestingServerCall collects the information about a supplied http client call
// and adds this information to the targetClient data structure. If possible, also calls the function to resolve
// the parameters of the function call.
func handleInterestingClientCall(call *ssa.Call, config *AnalyserConfig, packageName, qualifiedFunctionNameOfTarget string, frame *Frame) {
	interestingStuffClient := config.interestingCallsClient[qualifiedFunctionNameOfTarget]

	if interestingStuffClient.action == Output {
		requestLocation := ""
		var isResolved bool
		var variables []string

		// Additional information about the call
		service, file, position := getCallInformation(call.Pos(), frame.pkg)

		if call.Call.Args != nil && len(interestingStuffClient.interestingArgs) > 0 {
			variables, isResolved = resolveParameters(call.Call.Args, interestingStuffClient.interestingArgs, frame, config)
			// TODO: parse the url
			requestLocation = strings.Join(variables, "")
		}

		callTarget := &CallTarget{
			packageName:     packageName,
			MethodName:      qualifiedFunctionNameOfTarget,
			RequestLocation: requestLocation,
			IsResolved:      isResolved,
			ServiceName:     service,
			FileName:        file,
			PositionInFile:  position,
		}

		// fmt.Println("Found call to function " + qualifiedFunctionNameOfTarget)

		frame.targetsCollection.clientTargets = append(frame.targetsCollection.clientTargets, callTarget)
		return
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
func analyseInstructionsOfBlock(block *ssa.BasicBlock, fr *Frame, config *AnalyserConfig) {
	if block.Instrs == nil {
		return
	}

	for _, instr := range block.Instrs {
		switch instruction := instr.(type) {
		case *ssa.Call:
			analyseCall(instruction, fr, config)
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
func visitBlocks(blocks []*ssa.BasicBlock, fr *Frame, config *AnalyserConfig) {
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
		analyseInstructionsOfBlock(block, newFr, config)
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
		pkg:    pkg,
		prog:   pkg.Prog,
		params: make(map[*ssa.Parameter]*ssa.Value),
		// targetsCollection is a pointer to the global target collection.
		targetsCollection: &TargetsCollection{
			make([]*CallTarget, 0),
			make([]*CallTarget, 0),
		},
	}

	// Visit each of the block of the main function
	visitBlocks(mainFunction.Blocks, &baseFrame, config)

	// Here we can return the targets of the base frame: it is just a reference. All frames hold the same reference
	// to the targets collection.
	return baseFrame.targetsCollection.clientTargets, baseFrame.targetsCollection.serverTargets, nil
}
