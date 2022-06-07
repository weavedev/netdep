/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/

package callanalyzer

import (
	"strconv"
	"strings"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/preprocessing"
)

// ReplaceTargetsAnnotations replaces each unresolved callanalyzer.CallTarget with new a new target containing data
// obtained from the annotations (if they exist).
func ReplaceTargetsAnnotations(callTargets *[]*CallTarget, config *AnalyserConfig) error {
	if config == nil || config.annotations == nil {
		return nil
	}

	for _, callTarget := range *callTargets {
		if !callTarget.IsResolved {
			line, err := strconv.Atoi(callTarget.PositionInFile)
			if err != nil {
				return err
			}
			pos := preprocessing.Position{
				Filename: strings.ReplaceAll(callTarget.FileName, "\\", "/"),
				Line:     line - 1,
			}

			if ann, ex := config.annotations[callTarget.ServiceName][pos]; ex {
				resolveAnnotation(ann, callTarget)
			}
		}
	}
	return nil
}

/*
  resolveAnnotation populates the fields (RequestLocation or TargetSvc)
  of a CallTarget by extracting them from the annotation value string.

  Annotation format is currently:

  1) "//netdep:client url=... targetSvc=..."

  2) "//netdep:endpoint url=..."
*/
func resolveAnnotation(ann string, target *CallTarget) {
	annType := strings.Split(ann, " ")[0]
	annData := strings.Split(ann, " ")[1:]

	switch annType {
	case "client":
		// client type can have url=... and targetSvc=...
		for _, param := range annData {
			if strings.Split(param, "=")[0] == "url" {
				target.IsResolved = true
				target.RequestLocation = strings.Split(param, "=")[1]
			} else if strings.Split(param, "=")[0] == "targetSvc" {
				target.IsResolved = true
				target.TargetSvc = strings.Split(param, "=")[1]
			}
		}
	case "endpoint":
		// endpoint type can have url=...
		for _, param := range annData {
			if strings.Split(param, "=")[0] == "url" {
				target.IsResolved = true
				target.RequestLocation = strings.Split(param, "=")[1]
			}
		}
	}
}
