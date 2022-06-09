package callanalyzer

import (
	"fmt"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"os"
	"strings"
)

func getCallTarget(trace []*ssa.CallCommon, pkg *ssa.Package) *CallTarget {
	rootCall := trace[len(trace)-1]

	if rootCall == nil {
		return nil
	}

	rootFn := rootCall.StaticCallee()
	MethodName, packageName := getFunctionQualifiers(rootFn)

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
		if tracedCall == nil {
			continue
		}

		filePath, position := getPositionFromPos(tracedCall.Pos(), pkg.Prog)

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
	trace []*ssa.CallCommon
	// visited is shared between frames and keeps track of which nodes have been visited
	// to prevent repetitive visits.
	visited           map[int]bool
	program           *ssa.Program
	parent            *Frame
	config            *AnalyserConfig
	targetsCollection *TargetsCollection
}

func (f GraphFrame) hasVisited(call *ssa.CallCommon, node *callgraph.Node) bool {
	if _, ok := f.visited[node.ID]; ok {
		return true
	}

	for _, callee := range f.trace {
		if callee == call {
			return true
		}
	}

	return false
}

func AnalyzeUsingCallGraph(pkgs []*ssa.Package, config *AnalyserConfig) ([]*CallTarget, []*CallTarget, error) {
	mains := ssautil.MainPackages(pkgs)

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
		trace:   make([]*ssa.CallCommon, 0),
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

	fmt.Printf("%d, %d found\n", len(baseFrame.targetsCollection.clientTargets), len(baseFrame.targetsCollection.serverTargets))
	return baseFrame.targetsCollection.clientTargets, baseFrame.targetsCollection.serverTargets, nil
}

func uniquePaths(node *callgraph.Node, frame *GraphFrame) {
	//fmt.Printf("Scanning %s, %d (%d)\n", node.String(), len(node.Out), len(frame.trace))
	for _, edge := range node.Out {
		outNode := edge.Callee

		// TODO: improve get call
		var call *ssa.CallCommon = nil
		if edge.Site != nil {
			call = edge.Site.Common()
		}

		if frame.hasVisited(call, outNode) {
			continue
		}

		// check for interesting call
		outFunc := outNode.Func

		qualifiedFunctionNameOfTarget, functionPackage := getFunctionQualifiers(outFunc)

		// check ignored
		_, isIgnored := frame.config.ignoreList[functionPackage]
		if isIgnored {
			continue
		}

		_, isInterestingClient := frame.config.interestingCallsClient[qualifiedFunctionNameOfTarget]
		_, isInterestingServer := frame.config.interestingCallsServer[qualifiedFunctionNameOfTarget]

		if isInterestingClient || isInterestingServer {
			frame.trace = append(frame.trace, call)
			callTarget := getCallTarget(frame.trace, outFunc.Pkg)
			if isInterestingClient {
				frame.targetsCollection.clientTargets = append(frame.targetsCollection.clientTargets, callTarget)
			} else {
				frame.targetsCollection.serverTargets = append(frame.targetsCollection.serverTargets, callTarget)
			}

			continue
		}

		newFrame := *frame

		//newVisited := make(map[int]bool)
		//for k, _ := range frame.visited {
		//	newVisited[k] = true
		//}

		newFrame.visited[outNode.ID] = true
		//newFrame.visited = newVisited
		copy(newFrame.trace, frame.trace)
		newFrame.trace = append(newFrame.trace, call)

		uniquePaths(outNode, &newFrame)
	}
}
