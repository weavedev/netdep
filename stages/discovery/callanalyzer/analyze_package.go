/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package callanalyzer

import (
	"fmt"
	"go/token"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/fatih/color"

	"golang.org/x/tools/go/ssa"
)

// CallTargetTrace describes the physical location of a call in the stack trace
type CallTargetTrace struct {
	Internal       bool      // Internal defines an internal call
	Pos            token.Pos // Pos references the location of the call
	FileName       string    // FileName of the file in which the call is made
	PositionInFile string    // PositionInFile defines the line number in the file where the call is made
}

// CallTarget holds information about a certain call made by the analysed package.
// This used to be named "Caller" (which was slightly misleading, as it is in fact the target,
// thus rather a 'callee' than a 'caller'.
type CallTarget struct {
	PackageName     string            // PackageName is the name of the package the method belongs to
	MethodName      string            // MethodName is the name of the call (i.e. name of function or some other target)
	RequestLocation string            // RequestLocation is the URL of the entity
	IsResolved      bool              // IsResolved defines a flag describing whether the RequestLocation was resolved
	ServiceName     string            // ServiceName is the name of the service in which the call is made
	TargetSvc       string            // TargetSvc is the targeted service (in case the CallTarget is a client)
	Trace           []CallTargetTrace // Trace defines a stack trace for the call
}

// TraceAsStringArray maps the CallTargetTrace in Trace to a string
func (trace *CallTarget) TraceAsStringArray() []string {
	ret := make([]string, 0)
	for i := range trace.Trace {
		call := &trace.Trace[i]
		ret = append(ret, fmt.Sprintf("%s:%s", call.FileName, call.PositionInFile))
	}

	return ret
}

// SubstitutionConfig holds interesting calls to substitute,
// as well as a map of the current service's environment
type SubstitutionConfig struct {
	substitutionCalls map[string]InterestingCall // substitutionCalls defines a map of function signature to call
	serviceEnv        map[string]string          // serviceEnv holds the mapping of environment names to values
}

// TargetsCollection holds the output structures that are to be returned by the
// discovery stage
type TargetsCollection struct {
	clientTargets []*CallTarget // clientTargets is the collection of client calls (outgoing)
	serverTargets []*CallTarget // serverTargets is the collection of server calls (ingoing)
}

// getPositionFromPos converts a token.Pos to a filename and line number
func getPositionFromPos(pos token.Pos, program *ssa.Program) (string, string) {
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
		pos := tracedCall.Pos()
		filePath, position := getPositionFromPos(pos, frame.pkg.Prog)
		Internal := strings.Contains(filePath, string(os.PathSeparator)+callTarget.ServiceName+string(os.PathSeparator))

		newTrace := CallTargetTrace{
			// split package name and take the last item to get the service name
			FileName:       filePath[strings.LastIndex(filePath, fmt.Sprintf("%s%s%s", string(os.PathSeparator), callTarget.ServiceName, string(os.PathSeparator)))+1:],
			PositionInFile: position,
			Pos:            pos,
			Internal:       Internal,
		}

		callTarget.Trace = append(callTarget.Trace, newTrace)
	}

	return callTarget
}

// analyzeCallToFunction takes a call and its pointing function and analyses it (recursively)
// returns whether it found a path to an interesting call.
func analyzeCallToFunction(call *ssa.CallCommon, fn *ssa.Function, frame *Frame, config *AnalyserConfig) bool {
	wasInteresting := false

	// Qualified function name is: package + interface + function
	qualifiedFunctionNameOfTarget, functionPackage := getFunctionQualifiers(fn)

	if _, isIgnored := config.ignoreList[functionPackage]; isIgnored {
		// do not recurse on uninteresting packages
		return false
	}

	// The following creates a copy of 'frame'.
	// This is the correct place for this because we are going to visit child blocks next.
	newFrame := *frame

	// copy trace and append current call
	copy(newFrame.trace, frame.trace)
	newFrame.trace = append(newFrame.trace, call)
	newFrame.params = map[*ssa.Parameter]*ssa.Value{}

	// define offset when function was resolved to an invocation and the first parameter does not exist
	// this is the case for functions like `func (o obj) name (arg string) {}`
	offset := len(fn.Params) - len(call.Args)
	if offset < 0 {
		offset = 0
	}

	// Keep track of given parameters for resolving
	for i, par := range fn.Params[offset:] {
		newFrame.params[par] = &call.Args[i]
	}

	// Keep a reference to the parent frame
	newFrame.parent = frame

	_, isInterestingClient := config.interestingCallsClient[qualifiedFunctionNameOfTarget]
	if isInterestingClient {
		handleInterestingClientCall(call, fn, config, &newFrame)
		wasInteresting = true
	}

	_, isInterestingServer := config.interestingCallsServer[qualifiedFunctionNameOfTarget]
	if isInterestingServer {
		handleInterestingServerCall(call, fn, config, &newFrame)
		wasInteresting = true
	}

	// recurse into arguments if they are functions or calls themselves
	analyseCallArguments(call, frame, config)

	// do not recurse down on interesting calls
	if wasInteresting {
		// should be false even if interesting for singlePass
		frame.visited[call] = !frame.singlePass
		return true
	}

	// recurse into function blocks
	if fn.Blocks != nil {
		interesting := visitBlocks(fn.Blocks, &newFrame, config)

		// should be false even if interesting for singlePass
		frame.visited[call] = interesting && !frame.singlePass
	} else {
		frame.visited[call] = false
	}

	return false
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
func analyseCall(call *ssa.CallCommon, frame *Frame, config *AnalyserConfig) bool {
	if frame.hasVisited(call) || len(frame.trace) > config.maxTraversalDepth {
		return false
	}

	// fn := getFunctionFromCall(call, frame)
	fns, hasFn := frame.pointerMap[call]

	if !hasFn || len(fns) == 0 {
		return false
	}

	interesting := false
	for _, fn := range fns {
		wasInteresting := analyzeCallToFunction(call, fn, frame, config)
		if wasInteresting {
			interesting = true
		}
	}

	return interesting
}

// getHostFromAnnotation returns the resolved url using the annotated host name for a service
func getHostFromAnnotation(call *ssa.CallCommon, frame *Frame, config *AnalyserConfig, target *CallTarget) string {
	// absolute file path
	filePath := frame.pkg.Prog.Fset.File(call.Pos()).Name()
	// split path and form absolute path to service directory
	service := strings.Split(filePath, string(os.PathSeparator))
	service = service[:len(service)-1]
	serviceName := strings.Join(service, "/")

	annotations := config.annotations[serviceName]
	// look for annotated hostname
	for _, annotation := range annotations {
		if strings.HasPrefix(annotation, "host") {
			// if annotation is incorrectly formatted (eg. "hostsomething") the split will return ""
			// and the url will not be substituted from annotation
			host := strings.Join(strings.Split(annotation, "host ")[1:], "")
			resolvedURL, err := url.Parse(host)
			if err != nil {
				return target.RequestLocation
			}
			resolvedURL.Path = path.Join(resolvedURL.Path, target.RequestLocation)
			return resolvedURL.String()
		}
	}
	return target.RequestLocation
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

	callTarget.RequestLocation = getHostFromAnnotation(call, frame, config, callTarget)

	if !callTarget.IsResolved && config.Verbose {
		color.Yellow("Could not resolve variable(s) for call to " + qualifiedFunctionNameOfTarget)
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

// defaultCallTarget returns a new callTarget with initialised PackageName, functionName and IsResolved fields
func defaultCallTarget(packageName, functionName string) *CallTarget {
	return &CallTarget{
		PackageName:     packageName,
		MethodName:      functionName,
		RequestLocation: "",
		IsResolved:      false,
		ServiceName:     "",
		Trace:           []CallTargetTrace{},
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

	if !callTarget.IsResolved && config.Verbose {
		color.Yellow("Could not resolve variable(s) for call to " + qualifiedFunctionNameOfTarget)
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
func analyseInstructionsOfBlock(block *ssa.BasicBlock, fr *Frame, config *AnalyserConfig) bool {
	if block.Instrs == nil {
		return false
	}

	wasInteresting := false
	for _, instr := range block.Instrs {
		switch instruction := instr.(type) {
		case ssa.CallInstruction:
			interesting := analyseCall(instruction.Common(), fr, config)
			if interesting {
				wasInteresting = true
			}
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

	return wasInteresting
}

// visitBlocks visits each of the blocks in the specified 'blocks' list and analyses each of the block's instructions.
//
// Arguments:
// blocks is the array of blocks to analyse,
// fr keeps track of the traversal,
// config specifies the behaviour of the analyser,
// targets is a reference to the ultimate data structure that is to be completed and returned.
func visitBlocks(blocks []*ssa.BasicBlock, fr *Frame, config *AnalyserConfig) bool {
	wasInteresting := false

	for _, block := range blocks {
		interesting := analyseInstructionsOfBlock(block, fr, config)
		if interesting {
			wasInteresting = true
		}
	}

	return wasInteresting
}

func countVisited(visited map[*ssa.CallCommon]bool) (int, int) {
	interesting := 0
	count := 0
	for _, interstingVisit := range visited {
		count++
		if interstingVisit {
			interesting++
		}
	}
	return interesting, count
}

// AnalysePackageCalls takes a main package and finds all 'interesting' methods that are called
//
// Arguments:
// pkg is the package to analyse
// config specifies the behaviour of the analyser,
//
// Returns:
// List of pointers to callTargets, or an error if something went wrong.
func AnalysePackageCalls(pkg *ssa.Package, config *AnalyserConfig, pointerMap map[*ssa.CallCommon][]*ssa.Function) ([]*CallTarget, []*CallTarget, error) {
	if pkg == nil {
		return nil, nil, fmt.Errorf("no package given %v", pkg)
	}

	mainFunction := pkg.Func("main")
	initFunction := pkg.Func("init")

	// Find the main function
	if mainFunction == nil {
		return nil, nil, fmt.Errorf("no main function found in package %v", pkg)
	}

	baseFrame := Frame{
		trace: []*ssa.CallCommon{},
		// Reference to the final list of all _targets of the entire package
		pkg:        pkg,
		visited:    map[*ssa.CallCommon]bool{},
		params:     map[*ssa.Parameter]*ssa.Value{},
		globals:    map[*ssa.Global]*ssa.Value{},
		pointerMap: pointerMap,
		// for the init function we should only pass once
		// as we don't expect to find a functional call in the setup
		singlePass: true,
		// targetsCollection is a pointer to the global target collection.
		targetsCollection: &TargetsCollection{
			[]*CallTarget{},
			[]*CallTarget{},
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

	if config.Verbose {
		i, c := countVisited(baseFrame.visited)
		fmt.Printf("Init: Visited %d nodes, of which %d interesting\n", c, i)
	}

	// rest visited
	baseFrame.visited = map[*ssa.CallCommon]bool{}
	baseFrame.singlePass = false

	// Visit each of the block of the main function
	visitBlocks(mainFunction.Blocks, &baseFrame, config)

	if config.Verbose {
		i, c := countVisited(baseFrame.visited)
		fmt.Printf("Main: Visited %d nodes, of which %d interesting\n", c, i)
	}

	// Here we can return the targets of the base frame: it is just a reference. All frames hold the same reference
	// to the targets collection.
	return baseFrame.targetsCollection.clientTargets, baseFrame.targetsCollection.serverTargets, nil
}
