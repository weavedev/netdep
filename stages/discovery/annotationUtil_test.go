// Package discovery defines discovery of clients calls and endpoints
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package discovery

import (
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/preprocessing"
)

func TestReplaceTargetsAnnotations(t *testing.T) {
	target1 := &callanalyzer.CallTarget{
		MethodName:      "a",
		RequestLocation: "b",
		IsResolved:      true,
		ServiceName:     "c",
		TargetSvc:       "",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d",
				PositionInFile: "5",
			},
		},
	}
	target2 := &callanalyzer.CallTarget{
		MethodName:      "a1",
		RequestLocation: "",
		IsResolved:      false,
		ServiceName:     "c1",
		TargetSvc:       "",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d1",
				PositionInFile: "6",
			},
		},
	}

	targets := make([]*callanalyzer.CallTarget, 0)
	targets = append(targets, target1, target2)

	annotations := make(map[string]map[preprocessing.Position]string)
	annotations["c1"] = make(map[preprocessing.Position]string)

	pos := preprocessing.Position{
		Filename: "d1",
		Line:     5,
	}

	annotations["c1"][pos] = "client url=http://localhost:50/get"

	config := callanalyzer.DefaultConfigForFindingHTTPCalls()
	config.SetAnnotations(annotations)

	expectedTarget := &callanalyzer.CallTarget{
		MethodName:      "a1",
		RequestLocation: "http://localhost:50/get",
		IsResolved:      true,
		ServiceName:     "c1",
		TargetSvc:       "",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d1",
				PositionInFile: "6",
			},
		},
	}

	assert.NotEqual(t, expectedTarget, targets[1])
	callanalyzer.ReplaceTargetsAnnotations(&targets, &config)
	assert.Equal(t, expectedTarget, targets[1])
}

func TestReplaceTargetsAnnotationsConfigNil(t *testing.T) {
	target1 := &callanalyzer.CallTarget{
		MethodName:      "a",
		RequestLocation: "b",
		IsResolved:      true,
		ServiceName:     "c",
		TargetSvc:       "",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d",
				PositionInFile: "5",
			},
		},
	}
	targets := make([]*callanalyzer.CallTarget, 0)
	targets = append(targets, target1)

	err := callanalyzer.ReplaceTargetsAnnotations(&targets, nil)

	assert.Nil(t, err)
}

func TestReplaceTargetsAnnotationsAnnotationsNil(t *testing.T) {
	target1 := &callanalyzer.CallTarget{
		MethodName:      "a",
		RequestLocation: "b",
		IsResolved:      true,
		ServiceName:     "c",
		TargetSvc:       "",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d",
				PositionInFile: "5",
			},
		},
	}
	targets := make([]*callanalyzer.CallTarget, 0)
	targets = append(targets, target1)

	config := callanalyzer.DefaultConfigForFindingHTTPCalls()

	err := callanalyzer.ReplaceTargetsAnnotations(&targets, &config)

	assert.Nil(t, err)
}

func TestResolveAnnotationClientUrl(t *testing.T) {
	val := "client url=http://localhost:50/get"

	target := &callanalyzer.CallTarget{
		MethodName:      "a",
		RequestLocation: "",
		IsResolved:      true,
		ServiceName:     "c",
		TargetSvc:       "",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d",
				PositionInFile: "5",
			},
		},
	}

	expectedTarget := &callanalyzer.CallTarget{
		MethodName:      "a",
		RequestLocation: "http://localhost:50/get",
		IsResolved:      true,
		ServiceName:     "c",
		TargetSvc:       "",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d",
				PositionInFile: "5",
			},
		},
	}

	callanalyzer.ResolveAnnotation(val, target)

	assert.Equal(t, expectedTarget, target)
}

func TestResolveAnnotationClientTargetSvc(t *testing.T) {
	val := "client targetSvc=service2"

	target := &callanalyzer.CallTarget{
		MethodName:      "a",
		RequestLocation: "",
		IsResolved:      true,
		ServiceName:     "c",
		TargetSvc:       "",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d",
				PositionInFile: "5",
			},
		},
	}

	expectedTarget := &callanalyzer.CallTarget{
		MethodName:      "a",
		RequestLocation: "",
		IsResolved:      true,
		ServiceName:     "c",
		TargetSvc:       "service2",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d",
				PositionInFile: "5",
			},
		},
	}

	callanalyzer.ResolveAnnotation(val, target)

	assert.Equal(t, expectedTarget, target)
}

func TestResolveAnnotationClientBoth(t *testing.T) {
	val := "client url=http://localhost:50/get targetSvc=service2"

	target := &callanalyzer.CallTarget{
		MethodName:      "a",
		RequestLocation: "",
		IsResolved:      true,
		ServiceName:     "c",
		TargetSvc:       "",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d",
				PositionInFile: "5",
			},
		},
	}

	expectedTarget := &callanalyzer.CallTarget{
		MethodName:      "a",
		RequestLocation: "http://localhost:50/get",
		IsResolved:      true,
		ServiceName:     "c",
		TargetSvc:       "service2",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d",
				PositionInFile: "5",
			},
		},
	}

	callanalyzer.ResolveAnnotation(val, target)

	assert.Equal(t, expectedTarget, target)
}

func TestResolveAnnotationEndpointUrl(t *testing.T) {
	val := "endpoint url=http://localhost:50/get"

	target := &callanalyzer.CallTarget{
		MethodName:      "a",
		RequestLocation: "",
		IsResolved:      true,
		ServiceName:     "c",
		TargetSvc:       "",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d",
				PositionInFile: "5",
			},
		},
	}

	expectedTarget := &callanalyzer.CallTarget{
		MethodName:      "a",
		RequestLocation: "http://localhost:50/get",
		IsResolved:      true,
		ServiceName:     "c",
		TargetSvc:       "",
		Trace: []callanalyzer.CallTargetTrace{
			{
				FileName:       "d",
				PositionInFile: "5",
			},
		},
	}

	callanalyzer.ResolveAnnotation(val, target)

	assert.Equal(t, expectedTarget, target)
}
