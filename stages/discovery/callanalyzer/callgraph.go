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

// resolveValue Resolves a supplied ssa.Value, only in the cases that are supported by the tool:
// - string concatenation (see BinOp),
// - string literal
// - call to os.GetEnv
// - other InterestingCalls with the action Substitute.
// It also returns a bool which indicates whether the variable was resolved.
func resolveArgValue(value *ssa.Value, trace []*callgraph.Edge) (string, bool) {
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
			return resolveArgValue(&previousCall.Args[paramIndex], trace[:traceLength-1])
		}

		return "unknown: the parameter was not resolved", false
	case *ssa.Global:
		//if globalValue, ok := fr.globals[val]; ok {
		//	return resolveValue(globalValue, fr, substConf)
		//}

		return "unknown: the global was not resolved", false

	case *ssa.UnOp:
		return resolveArgValue(&val.X, trace)

	case *ssa.BinOp:
		switch val.Op { //nolint:exhaustive
		case token.ADD:
			left, isLeftResolved := resolveArgValue(&val.X, trace)
			right, isRightResolved := resolveArgValue(&val.Y, trace)
			if isRightResolved && isLeftResolved {
				return left + right, true
			}

			return left + right, false
		default:
			return "unknown: only ADD binary operation is supported", false
		}
	case *ssa.Const:
		switch val.Value.Kind() { //nolint:exhaustive
		case constant.String:
			return constant.StringVal(val.Value), true
		default:
			return "unknown: not a string constant", false
		}
	//case *ssa.Call:
	//	return handleSubstitutableCall(val, substConf)
	default:
		return "unknown: the parameter was not resolved", false
	}
}

func resolveParameterOfCall(trace []*callgraph.Edge, positions []int) ([]string, bool) {
	stringParameters := make([]string, len(positions))
	wasResolved := true

	callSite := trace[len(trace)-1].Site

	if callSite == nil {
		return stringParameters, wasResolved
	}

	parameters := callSite.Common().Args

	for i, idx := range positions {
		if idx < len(parameters) {
			variable, isResolved := resolveArgValue(&parameters[idx], trace)
			if isResolved {
				stringParameters[i] = variable
			} else {
				wasResolved = false
			}
		}
	}

	return stringParameters, wasResolved
}

func getCallTarget(trace []*callgraph.Edge, interest InterestingCall) *CallTarget {
	if len(trace) == 0 {
		return nil
	}
	rootCall := trace[0]
	finalCall := trace[len(trace)-1]

	pkg := rootCall.Caller.Func.Pkg
	program := pkg.Prog
	MethodName, packageName := getFunctionQualifiers(finalCall.Callee.Func)

	variables, IsResolved := resolveParameterOfCall(trace, interest.interestingArgs)

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

		filePath, position := getPositionFromPos(tracedCall.Callee.Func.Pos(), program)

		newTrace := CallTargetTrace{
			// split package name and take the last item to get the service name
			FileName:       filePath[strings.LastIndex(filePath, string(os.PathSeparator))+1:],
			PositionInFile: position,
		}

		target.Trace = append(target.Trace, newTrace)
	}

	return &target
}

// GraphFrame is a struct for keeping track of the traversal packages while looking for interesting functions
type GraphFrame struct {
	// trace is a stack trace of previous calls.
	trace []*callgraph.Edge
	// visited is shared between frames and keeps track of which nodes have been visited
	// to prevent repetitive visits.
	visited           map[int]bool
	program           *ssa.Program
	parent            *Frame
	config            *AnalyserConfig
	targetsCollection *TargetsCollection
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

	// query

	fmt.Println("Running pointer analysis...")
	ptares, err := pointer.Analyze(ptConfig)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	cg := ptares.CallGraph
	cg.DeleteSyntheticNodes()

	fmt.Printf("Finding nodes (%d)...\n", len(cg.Nodes))

	baseFrame := GraphFrame{
		trace:   make([]*callgraph.Edge, 0),
		visited: make(map[int]bool),
		program: pkgs[0].Prog,
		parent:  nil,
		config:  config,
		targetsCollection: &TargetsCollection{
			make([]*CallTarget, 0),
			make([]*CallTarget, 0),
		},
	}

	for _, edge := range cg.Root.Out {
		if edge.Callee.Func.Name() == "main" {
			baseFrame.visited[edge.Callee.ID] = true
			uniquePaths(edge.Callee, &baseFrame)
		}
	}

	count := 0
	interestingCount := 0
	for _, v := range baseFrame.visited {
		count++
		if v {
			interestingCount++
		}
	}

	fmt.Printf("%d, %d found (%d/%d)\n", len(baseFrame.targetsCollection.clientTargets), len(baseFrame.targetsCollection.serverTargets), interestingCount, count)

	return baseFrame.targetsCollection.clientTargets, baseFrame.targetsCollection.serverTargets, nil
}

func uniquePaths(node *callgraph.Node, frame *GraphFrame) bool {
	//fmt.Printf("Node scanning %s, %d (%d)\n", node.String(), len(node.Out), len(frame.trace))
	if len(frame.trace) > frame.config.maxTraversalDepth {
		return false
	}

	foundInteresting := false
	frame.visited[node.ID] = false

	for _, edge := range node.Out {
		//fmt.Printf("Edge scanning %s\n", outNode.String())
		outNode := edge.Callee

		// TODO: improve get call
		shouldSkip := false
		for _, tracedEdge := range frame.trace {
			if tracedEdge.Callee.ID == outNode.ID || tracedEdge.Caller.ID == outNode.ID {
				shouldSkip = true
				break
			}
		}

		if shouldSkip {
			continue
		}

		isInterestingVisit, wasVisited := frame.visited[outNode.ID]
		if wasVisited && !isInterestingVisit {
			continue
		}

		frame.visited[outNode.ID] = false
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

		newFrame := *frame
		copy(newFrame.trace, frame.trace)
		newFrame.trace = append(newFrame.trace, edge)

		if isInterestingClient || isInterestingServer {
			if isInterestingClient {
				callTarget := getCallTarget(newFrame.trace, interestingClient)
				frame.targetsCollection.clientTargets = append(frame.targetsCollection.clientTargets, callTarget)
			} else {
				callTarget := getCallTarget(newFrame.trace, interestingServer)
				frame.targetsCollection.serverTargets = append(frame.targetsCollection.serverTargets, callTarget)
			}

			foundInteresting = true
			frame.visited[outNode.ID] = true
			continue
		}

		found := uniquePaths(outNode, &newFrame)
		frame.visited[outNode.ID] = found

		if found {
			foundInteresting = true
		}
	}

	frame.visited[node.ID] = foundInteresting
	return foundInteresting
}
