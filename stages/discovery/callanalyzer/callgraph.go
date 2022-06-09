package callanalyzer

import (
	"fmt"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"os"
	"strings"
)

func getCallTarget(trace []CalledFunction) *CallTarget {
	if len(trace) == 0 {
		return nil
	}
	rootCall := trace[0]
	if rootCall.call == nil || rootCall.function == nil {
		return nil
	}

	pkg := rootCall.function.Pkg
	MethodName, packageName := getFunctionQualifiers(rootCall.function)

	target := CallTarget{
		packageName:     packageName,
		MethodName:      MethodName,
		ServiceName:     pkg.String()[strings.LastIndex(pkg.String(), "/")+1:],
		RequestLocation: "",
		IsResolved:      false,
		Trace:           nil,
	}

	// add trace
	for _, tracedCall := range trace {
		if tracedCall.call == nil {
			continue
		}

		filePath, position := getPositionFromPos(tracedCall.call.Pos(), pkg.Prog)

		newTrace := CallTargetTrace{
			// split package name and take the last item to get the service name
			FileName:       filePath[strings.LastIndex(filePath, string(os.PathSeparator))+1:],
			PositionInFile: position,
		}

		target.Trace = append(target.Trace, newTrace)
	}

	return &target
}

type CalledFunction struct {
	call     *ssa.CallCommon
	function *ssa.Function
}

// GraphFrame is a struct for keeping track of the traversal packages while looking for interesting functions
type GraphFrame struct {
	// trace is a stack trace of previous calls.
	trace []CalledFunction
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
		trace:   make([]CalledFunction, 0),
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
	//fmt.Printf("Scanning %s, %d (%d)\n", node.String(), len(node.Out), len(frame.trace))
	if len(frame.trace) > frame.config.maxTraversalDepth {
		return false
	}

	foundInteresting := false

	for _, edge := range node.Out {
		outNode := edge.Callee

		// TODO: improve get call
		called := CalledFunction{
			call:     nil,
			function: outNode.Func,
		}

		if edge.Site != nil {
			called.call = edge.Site.Common()
		}

		isInterestingVisit, wasVisited := frame.visited[outNode.ID]
		if wasVisited && !isInterestingVisit {
			continue
		}

		frame.visited[outNode.ID] = true
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

		_, isInterestingClient := frame.config.interestingCallsClient[qualifiedFunctionNameOfTarget]
		_, isInterestingServer := frame.config.interestingCallsServer[qualifiedFunctionNameOfTarget]

		newFrame := *frame
		copy(newFrame.trace, frame.trace)
		newFrame.trace = append(newFrame.trace, called)

		if isInterestingClient || isInterestingServer {
			callTarget := getCallTarget(newFrame.trace)
			if isInterestingClient {
				frame.targetsCollection.clientTargets = append(frame.targetsCollection.clientTargets, callTarget)
			} else {
				frame.targetsCollection.serverTargets = append(frame.targetsCollection.serverTargets, callTarget)
			}

			frame.visited[outNode.ID] = true
			continue
		}

		found := uniquePaths(outNode, &newFrame)
		frame.visited[outNode.ID] = found

		if found {
			foundInteresting = true
		}
	}

	return foundInteresting
}
