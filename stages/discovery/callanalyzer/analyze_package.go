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

type CallTargetTrace struct {
	Internal bool
	Pos      token.Pos
	// The name of the file in which the call is made
	FileName string
	// The line number in the file where the call is made
	PositionInFile string
}

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
	// Trace defines a stack trace for the call
	Trace []CallTargetTrace
}

// SubstitutionConfig holds interesting calls to substitute,
// as well as a map of the current service's environment
type SubstitutionConfig struct {
	substitutionCalls map[string]InterestingCall
	serviceEnv        map[string]string
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

// getPositionFromPos converts a token.Pos to a filename and line number
func getPositionFromPos(pos token.Pos, program *ssa.Program) (string, string) {
	if program == nil {
		return "", ""
	}

	file := program.Fset.File(pos)
	if file == nil {
		return "", ""
	}

	filePath := file.Name()

	base := 10
	// take the position of the call within the file and convert to string
	positionInFile := strconv.FormatInt(int64(program.Fset.Position(pos).Line), base)

	return filePath, positionInFile
}

// getFunctionQualifiers returns the function signature and the function package name
func getFunctionQualifiers(fn *ssa.Function) (string, string) {
	// .Pkg returns an obj of type *ssa.Package, whose .Pkg returns one of *type.Package
	// This is therefore not the grandparent package, but the *type.Package of the function
	calledFunctionPackage := ""
	if fn.Package() != nil && fn.Package().Pkg != nil {
		calledFunctionPackage = fn.Package().Pkg.Path() // e.g. net/http
	}

	return fn.RelString(nil), calledFunctionPackage
}

// getCallInformation creates a callTarget from a function and its trace
func getCallInformation(frame *Frame, fn *ssa.Function) *CallTarget {
	functionName, packageName := getFunctionQualifiers(fn)
	callTarget := defaultCallTarget(packageName, functionName)

	callTarget.ServiceName = frame.pkg.String()[strings.LastIndex(frame.pkg.String(), "/")+1:]

	// add trace
	for _, tracedCall := range frame.trace {
		filePath, position := getPositionFromPos(tracedCall.Pos(), frame.pkg.Prog)

		Internal := strings.Contains(filePath, string(os.PathSeparator)+callTarget.ServiceName+string(os.PathSeparator))

		newTrace := CallTargetTrace{
			// split package name and take the last item to get the service name
			FileName:       filePath[strings.LastIndex(filePath, string(os.PathSeparator)+callTarget.ServiceName+string(os.PathSeparator))+1:],
			PositionInFile: position,
			Internal:       Internal,
		}

		callTarget.Trace = append(callTarget.Trace, newTrace)
	}

	return callTarget
}

func analyzeCallToFunction(call *ssa.CallCommon, fn *ssa.Function, frame *Frame, config *AnalyserConfig) {
	wasInteresting := false

	// Qualified function name is: package + interface + function
	// TODO: handle parameter equivalence to other interface
	qualifiedFunctionNameOfTarget, functionPackage := getFunctionQualifiers(fn)

	// The following creates a copy of 'frame'.
	// This is the correct place for this because we are going to visit child blocks next.
	newFrame := *frame

	// copy trace and append current call
	copy(newFrame.trace, frame.trace)
	newFrame.trace = append(newFrame.trace, call)

	// offset when function was resolved to an invocation and the first parameter does not exist
	offset := len(fn.Params) - len(call.Args)

	// Keep track of given parameters for resolving
	for i, par := range fn.Params[offset:] {
		newFrame.params[par] = &call.Args[i]
	}

	// Keep a reference to the parent frame
	newFrame.parent = frame

	_, isInterestingClient := config.interestingCallsClient[qualifiedFunctionNameOfTarget]
	if isInterestingClient {
		// TODO: Resolve the arguments of the function call
		handleInterestingClientCall(call, fn, config, &newFrame)
		wasInteresting = true
	}

	_, isInterestingServer := config.interestingCallsServer[qualifiedFunctionNameOfTarget]
	if isInterestingServer {
		// TODO: Resolve the arguments of the function call
		handleInterestingServerCall(call, fn, config, &newFrame)
		wasInteresting = true
	}

	_, isIgnored := config.ignoreList[functionPackage]

	if isIgnored {
		// Do not recurse into the packageName if it is ignored
		return
	}

	// recurse into arguments if they are functions or calls themselves
	analyseCallArguments(call, frame, config)

	// at this point analyseCallArguments has been called so we can return
	if wasInteresting {
		return
	}

	frame.visited[call] = true

	// recurse into function blocks
	if fn.Blocks != nil {
		visitBlocks(fn.Blocks, &newFrame, config)
	}
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
func analyseCall(call *ssa.CallCommon, frame *Frame, config *AnalyserConfig) {
	if frame.hasVisited(call) || len(frame.trace) > config.maxTraversalDepth {
		return
	}

	fn := getFunctionFromCall(call, frame)

	if fn == nil {
		return
	}

	analyzeCallToFunction(call, fn, frame, config)
}

// analyseCallArguments goes over the call arguments and recurses into them
// given that they potentially contain another block of code. That is possible in two cases:
// 1. argument is a function. For example, a callback.
// 2. argument is another call. For example. http.Get(getEndpoint(smth))
func analyseCallArguments(call *ssa.CallCommon, fr *Frame, config *AnalyserConfig) {
	for _, argument := range call.Args {
		// visit function as argument
		if functionArg, ok := argument.(*ssa.Function); ok {
			visitBlocks(functionArg.Blocks, fr, config)
		}
	}
}

// handleInterestingServerCall collects the information about a supplied endpoint declaration
// and adds this information to the targetsServer data structure. If possible, also calls the function to resolve
// the parameters of the function call.
func handleInterestingServerCall(call *ssa.CallCommon, fn *ssa.Function, config *AnalyserConfig, frame *Frame) {
	qualifiedFunctionNameOfTarget, _ := getFunctionQualifiers(fn)
	interestingStuffServer := config.interestingCallsServer[qualifiedFunctionNameOfTarget]
	if interestingStuffServer.action != Output {
		return
	}
	// variables store the local variables of the call target
	var variables []string

	callTarget := getCallInformation(frame, fn)

	if call.Args != nil && len(interestingStuffServer.interestingArgs) > 0 {
		if qualifiedFunctionNameOfTarget == "(*github.com/gin-gonic/gin.Engine).Run" {
			variables, callTarget.IsResolved = resolveGinAddrSlice(call.Args[1])
			// TODO: parse the url
			callTarget.RequestLocation = strings.Join(variables, "")
		} else {
			// Since the environment can vary on a per-service basis,
			// a substConfig is created for the specific service
			substitutionConfig := getSubstConfig(config, callTarget.ServiceName)
			variables, callTarget.IsResolved = resolveParameters(call.Args, interestingStuffServer.interestingArgs, frame, substitutionConfig)
			// TODO: parse the url
			callTarget.RequestLocation = strings.Join(variables, "")
		}
	}

	if !callTarget.IsResolved && config.verbose {
		fmt.Println("Could not resolve variable(s) for call to " + qualifiedFunctionNameOfTarget)
		PrintTraceToCall(frame, config)
	}

	// Additional information about the call
	frame.targetsCollection.serverTargets = append(frame.targetsCollection.serverTargets, callTarget)
}

// getSubstConfig returns the substitution config (environment)
// for the specific service
func getSubstConfig(config *AnalyserConfig, service string) SubstitutionConfig {
	return SubstitutionConfig{
		config.substitutionCalls,
		config.environment[service],
	}
}

// defaultCallTarget returns a new callTarget with initialised packageName, functionName and IsResolved fields
func defaultCallTarget(packageName, functionName string) *CallTarget {
	return &CallTarget{
		packageName:     packageName,
		MethodName:      functionName,
		RequestLocation: "",
		IsResolved:      false,
		ServiceName:     "",
		Trace:           make([]CallTargetTrace, 0),
	}
}

// handleInterestingServerCall collects the information about a supplied http client call
// and adds this information to the targetClient data structure. If possible, also calls the function to resolve
// the parameters of the function call.
func handleInterestingClientCall(call *ssa.CallCommon, fn *ssa.Function, config *AnalyserConfig, frame *Frame) {
	qualifiedFunctionNameOfTarget, _ := getFunctionQualifiers(fn)
	interestingStuffClient := config.interestingCallsClient[qualifiedFunctionNameOfTarget]

	if interestingStuffClient.action != Output {
		return
	}

	// variables store the local variables of the call target
	var variables []string

	// callTarget holds all the details of the interesting call
	callTarget := getCallInformation(frame, fn)

	if call.Args != nil && len(interestingStuffClient.interestingArgs) > 0 {
		// Since the environment can vary on a per-service basis,
		// a substConfig is created for the specific service
		substitutionConfig := getSubstConfig(config, callTarget.ServiceName)
		variables, callTarget.IsResolved = resolveParameters(call.Args, interestingStuffClient.interestingArgs, frame, substitutionConfig)
		// TODO: parse the url
		callTarget.RequestLocation = strings.Join(variables, "")
	}

	if !callTarget.IsResolved && config.verbose {
		fmt.Println("Could not resolve variable(s) for call to " + qualifiedFunctionNameOfTarget)
		PrintTraceToCall(frame, config)
	}

	frame.targetsCollection.clientTargets = append(frame.targetsCollection.clientTargets, callTarget)
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
			analyseCall(&instruction.Call, fr, config)
		case *ssa.Go:
			analyseCall(&instruction.Call, fr, config)
		case *ssa.Defer:
			analyseCall(&instruction.Call, fr, config)
		case *ssa.Store:
			// for a store to a value
			if global, ok := instruction.Addr.(*ssa.Global); ok {
				// TODO: structure this in a way that doesn't corrupt the value
				// When recursing. Value might not correspond to actual value!

				if _, ok := fr.globals[global]; ok {
					// only save package globals!
					fr.globals[global] = &instruction.Val
				}
			}
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
	for _, block := range blocks {
		analyseInstructionsOfBlock(block, fr, config)
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
	if pkg == nil {
		return nil, nil, fmt.Errorf("no package given %v", pkg)
	}

	mainFunction := findFunctionInPackage(pkg, "main")
	initFunction := findFunctionInPackage(pkg, "init")

	// Find the main function
	if mainFunction == nil {
		return nil, nil, fmt.Errorf("no main function found in package %v", pkg)
	}

	baseFrame := Frame{
		trace: make([]*ssa.CallCommon, 0),
		// Reference to the final list of all _targets of the entire package
		pkg:     pkg,
		visited: make(map[*ssa.CallCommon]bool),
		params:  make(map[*ssa.Parameter]*ssa.Value),
		globals: make(map[*ssa.Global]*ssa.Value),
		// for the init function we should only pass once
		// as we don't expect to find a functional call in the setup
		singlePass: true,
		// targetsCollection is a pointer to the global target collection.
		targetsCollection: &TargetsCollection{
			make([]*CallTarget, 0),
			make([]*CallTarget, 0),
		},
	}

	// setup basic references to global variables
	for _, m := range pkg.Members {
		if globalPointer, ok := m.(*ssa.Global); ok {
			baseFrame.globals[globalPointer] = nil
		}
	}

	// Visit the init function for globals
	visitBlocks(initFunction.Blocks, &baseFrame, config)

	// rest visited
	baseFrame.visited = make(map[*ssa.CallCommon]bool)
	baseFrame.singlePass = false

	// Visit each of the block of the main function
	visitBlocks(mainFunction.Blocks, &baseFrame, config)

	// Here we can return the targets of the base frame: it is just a reference. All frames hold the same reference
	// to the targets collection.
	return baseFrame.targetsCollection.clientTargets, baseFrame.targetsCollection.serverTargets, nil
}
