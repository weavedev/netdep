/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/

package callanalyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/preprocessing"
)

func TestReplaceTargetsAnnotations(t *testing.T) {
	target1 := &CallTarget{
		MethodName:      "a",
		RequestLocation: "b",
		IsResolved:      true,
		ServiceName:     "c",
		FileName:        "d",
		PositionInFile:  "5",
		TargetSvc:       "",
	}
	target2 := &CallTarget{
		MethodName:      "a1",
		RequestLocation: "",
		IsResolved:      false,
		ServiceName:     "c1",
		FileName:        "d1",
		PositionInFile:  "6",
		TargetSvc:       "",
	}

	targets := make([]*CallTarget, 0)
	targets = append(targets, target1, target2)

	annotations := make(map[string]map[preprocessing.Position]string)
	annotations["c1"] = make(map[preprocessing.Position]string)

	pos := preprocessing.Position{
		Filename: "d1",
		Line:     5,
	}

	annotations["c1"][pos] = "client url=http://localhost:50/get"

	config := DefaultConfigForFindingHTTPCalls(nil, annotations)

	expectedTarget := &CallTarget{
		MethodName:      "a1",
		RequestLocation: "http://localhost:50/get",
		IsResolved:      true,
		ServiceName:     "c1",
		FileName:        "d1",
		PositionInFile:  "6",
		TargetSvc:       "",
	}

	assert.NotEqual(t, expectedTarget, targets[1])
	ReplaceTargetsAnnotations(&targets, &config)
	assert.Equal(t, expectedTarget, targets[1])
}

func TestReplaceTargetsAnnotationsConfigNil(t *testing.T) {
	target1 := &CallTarget{
		MethodName:      "a",
		RequestLocation: "b",
		IsResolved:      true,
		ServiceName:     "c",
		FileName:        "d",
		PositionInFile:  "5",
		TargetSvc:       "",
	}
	targets := make([]*CallTarget, 0)
	targets = append(targets, target1)

	err := ReplaceTargetsAnnotations(&targets, nil)

	assert.Nil(t, err)
}

func TestReplaceTargetsAnnotationsAnnotationsNil(t *testing.T) {
	target1 := &CallTarget{
		MethodName:      "a",
		RequestLocation: "b",
		IsResolved:      true,
		ServiceName:     "c",
		FileName:        "d",
		PositionInFile:  "5",
		TargetSvc:       "",
	}
	targets := make([]*CallTarget, 0)
	targets = append(targets, target1)

	config := DefaultConfigForFindingHTTPCalls(nil, nil)

	err := ReplaceTargetsAnnotations(&targets, &config)

	assert.Nil(t, err)
}

func TestResolveAnnotationClientUrl(t *testing.T) {
	val := "client url=http://localhost:50/get"

	target := &CallTarget{
		MethodName:      "a",
		RequestLocation: "",
		IsResolved:      true,
		ServiceName:     "c",
		FileName:        "d",
		PositionInFile:  "5",
		TargetSvc:       "",
	}

	expectedTarget := &CallTarget{
		MethodName:      "a",
		RequestLocation: "http://localhost:50/get",
		IsResolved:      true,
		ServiceName:     "c",
		FileName:        "d",
		PositionInFile:  "5",
		TargetSvc:       "",
	}

	resolveAnnotation(val, target)

	assert.Equal(t, expectedTarget, target)
}

func TestResolveAnnotationClientTargetSvc(t *testing.T) {
	val := "client targetSvc=service2"

	target := &CallTarget{
		MethodName:      "a",
		RequestLocation: "",
		IsResolved:      true,
		ServiceName:     "c",
		FileName:        "d",
		PositionInFile:  "5",
		TargetSvc:       "",
	}

	expectedTarget := &CallTarget{
		MethodName:      "a",
		RequestLocation: "",
		IsResolved:      true,
		ServiceName:     "c",
		FileName:        "d",
		PositionInFile:  "5",
		TargetSvc:       "service2",
	}

	resolveAnnotation(val, target)

	assert.Equal(t, expectedTarget, target)
}

func TestResolveAnnotationClientBoth(t *testing.T) {
	val := "client url=http://localhost:50/get targetSvc=service2"

	target := &CallTarget{
		MethodName:      "a",
		RequestLocation: "",
		IsResolved:      true,
		ServiceName:     "c",
		FileName:        "d",
		PositionInFile:  "5",
		TargetSvc:       "",
	}

	expectedTarget := &CallTarget{
		MethodName:      "a",
		RequestLocation: "http://localhost:50/get",
		IsResolved:      true,
		ServiceName:     "c",
		FileName:        "d",
		PositionInFile:  "5",
		TargetSvc:       "service2",
	}

	resolveAnnotation(val, target)

	assert.Equal(t, expectedTarget, target)
}

func TestResolveAnnotationEndpointUrl(t *testing.T) {
	val := "endpoint url=http://localhost:50/get"

	target := &CallTarget{
		MethodName:      "a",
		RequestLocation: "",
		IsResolved:      true,
		ServiceName:     "c",
		FileName:        "d",
		PositionInFile:  "5",
		TargetSvc:       "",
	}

	expectedTarget := &CallTarget{
		MethodName:      "a",
		RequestLocation: "http://localhost:50/get",
		IsResolved:      true,
		ServiceName:     "c",
		FileName:        "d",
		PositionInFile:  "5",
		TargetSvc:       "",
	}

	resolveAnnotation(val, target)

	assert.Equal(t, expectedTarget, target)
}
