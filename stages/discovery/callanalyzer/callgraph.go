package callanalyzer

import (
	"fmt"
	"go/constant"
	"go/token"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"os"
	"strings"
)

// DiscoverFrame is a struct for keeping track of the traversal packages while looking for interesting functions
type DiscoverFrame struct {
	// trace is a stack trace of previous calls.
	trace []*callgraph.Edge
	// visited is shared between frames and keeps track of which nodes have been visited
	// to prevent repetitive visits.
	visited       map[int]int
	globals       map[*ssa.Global]*ssa.Value
	program       *ssa.Program
	parent        *Frame
	config        *AnalyserConfig
	result        *pointer.Result
	clientTargets []*CallTarget
	serverTargets []*CallTarget
	singlePass    bool
}

const maxRevisits = 1

// resolveValue Resolves a supplied ssa.Value, only in the cases that are supported by the tool:
// - string concatenation (see BinOp),
// - string literal
// - call to os.GetEnv
// - other InterestingCalls with the action Substitute.
// It also returns a bool which indicates whether the variable was resolved.
func resolveArgValue(value *ssa.Value, trace []*callgraph.Edge, frame *DiscoverFrame) (string, bool) {
	if value == nil {
		return "unknown: the give value is null", false
	}

	traceLength := len(trace)
	edge := trace[traceLength-1]
	calledFunc := edge.Callee.Func
	//call := edge.Site.Common()

	switch val := (*value).(type) {
	case *ssa.Parameter:
		// (recursively) resolve a parameter to a value and return that value, if it is defined
		paramIndex := -1
		for index, param := range calledFunc.Params {
			if param.Name() == val.Name() {
				paramIndex = index
			}
		}
		if paramIndex > -1 && traceLength > 1 {
			previousCall := trace[traceLength-2].Site.Common()
			return resolveArgValue(&previousCall.Args[paramIndex], trace[:traceLength-1], frame)
		}

		return "unknown: the parameter was not resolved", false
	case *ssa.Global:
		// TODO
		if p, ok := frame.result.Queries[val]; ok {
			return p.String(), true
		}

		return "unknown: the global was not resolved", false

	case *ssa.UnOp:
		return resolveArgValue(&val.X, trace, frame)

	case *ssa.BinOp:
		switch val.Op { //nolint:exhaustive
		case token.ADD:
			left, isLeftResolved := resolveArgValue(&val.X, trace, frame)
			right, isRightResolved := resolveArgValue(&val.Y, trace, frame)
			if isRightResolved && isLeftResolved {
				return left + right, true
			}

			return left + right, false
		default:
			return "unknown: only ADD binary operation is supported", false
		}
	case *ssa.Const:
		if val.Value != nil {
			switch val.Value.Kind() { //nolint:exhaustive
			case constant.String:
				return constant.StringVal(val.Value), true
			default:
				return "unknown: not a string constant", false
			}
		}
		return "unknown: not a string constant", false
	//case *ssa.Call:
	// TODO
	//	return handleSubstitutableCall(val, substConf)
	default:
		return "unknown: the parameter was not resolved", false
	}
}

func resolveParameterOfCall(trace []*callgraph.Edge, positions []int, frame *DiscoverFrame) ([]string, bool) {
	stringParameters := make([]string, len(positions))
	wasResolved := true

	callSite := trace[len(trace)-1].Site

	if callSite == nil {
		return stringParameters, wasResolved
	}

	arguments := callSite.Common().Args

	for i, idx := range positions {
		if idx < len(arguments) {
			variable, isResolved := resolveArgValue(&arguments[idx], trace, frame)
			if isResolved {
				stringParameters[i] = variable
			} else {
				wasResolved = false
			}
		}
	}

	return stringParameters, wasResolved
}

func getCallTarget(trace []*callgraph.Edge, interest InterestingCall, frame *DiscoverFrame) *CallTarget {
	if len(trace) == 0 {
		return nil
	}
	rootCall := trace[0]
	targetCall := trace[len(trace)-1]

	pkg := rootCall.Caller.Func.Pkg
	program := pkg.Prog
	MethodName, packageName := getFunctionQualifiers(targetCall.Callee.Func)

	variables, IsResolved := resolveParameterOfCall(trace, interest.interestingArgs, frame)

	target := CallTarget{
		packageName:     packageName,
		MethodName:      MethodName,
		ServiceName:     pkg.String()[strings.LastIndex(pkg.String(), "/")+1:],
		RequestLocation: strings.Join(variables, ""),
		IsResolved:      IsResolved,
		Trace:           nil,
	}

	// add trace
	for _, tracedCall := range trace {
		if tracedCall.Site == nil {
			continue
		}

		call := tracedCall.Site.Common()
		pos := call.Pos()

		filePath, position := getPositionFromPos(pos, program)

		Internal := strings.Contains(filePath, string(os.PathSeparator)+target.ServiceName+string(os.PathSeparator))

		newTrace := CallTargetTrace{
			// split package name and take the last item to get the service name
			FileName:       filePath[strings.LastIndex(filePath, string(os.PathSeparator)+target.ServiceName+string(os.PathSeparator))+1:],
			PositionInFile: position,
			Pos:            pos,
			Internal:       Internal,
		}

		target.Trace = append(target.Trace, newTrace)
	}

	return &target
}

func AnalyzeUsingCallGraph(pkgs []*ssa.Package, config *AnalyserConfig) ([]*CallTarget, []*CallTarget, error) {
	if pkgs == nil || len(pkgs) == 0 {
		return nil, nil, fmt.Errorf("no packages given")
	}

	var mains []*ssa.Package
	for _, pkg := range pkgs {
		if pkg == nil || pkg.Pkg == nil {
			if pkg != nil {
				fmt.Println("No package for " + pkg.String())
			}
			continue
		}

		if pkg.Pkg.Name() == "main" && pkg.Func("main") != nil {
			mains = append(mains, pkg)
		}
	}

	ptConfig := &pointer.Config{
		Mains:          mains,
		BuildCallGraph: true,
	}

	//for _, pkg := range pkgs {
	//	if pkg == nil {
	//		continue
	//	}
	//	for _, mem := range pkg.Members {
	//		if globalPointer, ok := mem.(*ssa.Global); ok {
	//			ptConfig.AddQuery(globalPointer)
	//		}
	//	}
	//}
	// query

	if config.verbose {
		fmt.Println("Running pointer analysis...")
	}

	pointerRes, err := pointer.Analyze(ptConfig)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	cg := pointerRes.CallGraph
	cg.DeleteSyntheticNodes()

	if config.verbose {
		fmt.Printf("Finding nodes (%d)...\n", len(cg.Nodes))
	}

	baseFrame := DiscoverFrame{
		trace: []*callgraph.Edge{},
		visited: map[int]int{
			cg.Root.ID: -1,
		},
		program:       pkgs[0].Prog,
		parent:        nil,
		config:        config,
		clientTargets: []*CallTarget{},
		serverTargets: []*CallTarget{},
		result:        pointerRes,
		//singlePass:    true,
	}

	for _, edge := range cg.Root.Out {
		if edge.Callee.Func.Name() == "main" {
			baseFrame.visited[edge.Callee.ID] = -1
			uniquePaths(edge.Callee, &baseFrame)
		}
	}

	count := 0
	interestingCount := 0
	for _, v := range baseFrame.visited {
		count++
		if v > -1 {
			interestingCount++
		}
	}

	if config.verbose {
		fmt.Printf("%d, %d found (%d/%d)\n", len(baseFrame.clientTargets), len(baseFrame.serverTargets), interestingCount, count)
	}

	return baseFrame.clientTargets, baseFrame.serverTargets, nil
}

func uniquePaths(node *callgraph.Node, frame *DiscoverFrame) bool {
	if frame.config.verbose {
		fmt.Printf("Node scanning %s, %d (%d)\n", node.String(), len(node.Out), len(frame.trace))
	}

	if len(frame.trace) > frame.config.maxTraversalDepth {
		frame.visited[node.ID] = -1
		return false
	}

	foundInteresting := false
	_, nodeVisited := frame.visited[node.ID]
	if !nodeVisited {
		frame.visited[node.ID] = -1
	}

	if len(node.Out) == 0 {
		return false
	}

	for _, edge := range node.Out {
		//fmt.Printf("Edge scanning %s\n", outNode.String())
		outNode := edge.Callee

		// TODO: improve get call
		shouldSkip := false
		for _, tracedEdge := range frame.trace {
			if tracedEdge == edge || tracedEdge.Callee.ID == outNode.ID || tracedEdge.Caller.ID == outNode.ID {
				shouldSkip = true
				break
			}
		}

		visited, wasVisited := frame.visited[outNode.ID]
		if wasVisited && (visited == -1 || visited > maxRevisits) {
			shouldSkip = true
		}

		if shouldSkip {
			continue
		}

		if !wasVisited {
			frame.visited[outNode.ID] = -1
		}

		// check for interesting call
		outFunc := outNode.Func

		qualifiedFunctionNameOfTarget, functionPackage := getFunctionQualifiers(outFunc)

		// check ignored
		_, isIgnored := frame.config.ignoreList[functionPackage]
		if isIgnored {
			continue
		}

		rootPackage := strings.Split(functionPackage, "/")[0]
		_, isIgnored = frame.config.ignoreList[rootPackage]
		if isIgnored {
			continue
		}

		interestingClient, isInterestingClient := frame.config.interestingCallsClient[qualifiedFunctionNameOfTarget]
		interestingServer, isInterestingServer := frame.config.interestingCallsServer[qualifiedFunctionNameOfTarget]

		frame.trace = append(frame.trace, edge)

		if isInterestingClient || isInterestingServer {
			if isInterestingClient {
				callTarget := getCallTarget(frame.trace, interestingClient, frame)
				frame.clientTargets = append(frame.clientTargets, callTarget)
			} else {
				callTarget := getCallTarget(frame.trace, interestingServer, frame)
				frame.serverTargets = append(frame.serverTargets, callTarget)
			}

			if frame.config.verbose {
				fmt.Printf("Found new trace! %d, %d\n", len(frame.clientTargets), len(frame.serverTargets))
			}

			foundInteresting = true
			frame.visited[outNode.ID] = 1

			//	pop trace
			frame.trace = frame.trace[:len(frame.trace)-1]
			continue
		}

		found := uniquePaths(outNode, frame)
		//	pop trace
		frame.trace = frame.trace[:len(frame.trace)-1]

		if found {
			foundInteresting = true
			frame.visited[outNode.ID]++
		}
	}

	if foundInteresting && !nodeVisited {
		frame.visited[node.ID] = 0
	}
	return foundInteresting
}
