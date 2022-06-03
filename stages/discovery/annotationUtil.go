package discovery

import (
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/preprocessing"
	"strconv"
	"strings"
)

func replaceTargetsAnnotations(callTargets *[]*callanalyzer.CallTarget, annotations map[string]map[preprocessing.Position]string) error {
	if annotations == nil {
		return nil
	}

	for i, callTarget := range *callTargets {
		if !callTarget.IsResolved {
			line, err := strconv.Atoi(callTarget.PositionInFile)
			if err != nil {
				return err
			}
			pos := preprocessing.Position{
				Filename: strings.ReplaceAll(callTarget.FileName, "\\", "/"),
				Line:     line - 1,
			}

			if ann, ex := annotations[callTarget.ServiceName][pos]; ex {
				newTarget := &callanalyzer.CallTarget{
					MethodName:      callTarget.MethodName,
					RequestLocation: "",
					IsResolved:      true,
					ServiceName:     callTarget.ServiceName,
					FileName:        callTarget.FileName,
					PositionInFile:  callTarget.PositionInFile,
					TargetSvc:       "",
				}
				resolveAnnotation(ann, newTarget)
				(*callTargets)[i] = newTarget
			}
		}
	}
	return nil
}

func resolveAnnotation(ann string, target *callanalyzer.CallTarget) {
	annType := strings.Split(ann, " ")[0]
	annData := strings.Split(ann, " ")[1:]

	switch annType {
	case "client":
		// client type can have url=... and targetSvc=...
		for _, param := range annData {
			if strings.Split(param, "=")[0] == "url" {
				target.RequestLocation = strings.Split(param, "=")[1]
			} else if strings.Split(param, "=")[0] == "targetSvc" {
				target.TargetSvc = strings.Split(param, "=")[1]
			}
		}
	case "server":
		// server type can have url=...
		for _, param := range annData {
			if strings.Split(param, "=")[0] == "url" {
				target.RequestLocation = strings.Split(param, "=")[1]
			}
		}
	}
}
