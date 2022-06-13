/*
Package nats contains NATS specific call analysis
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package nats

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type NatsCall struct {
	// The name of the package the method belongs to
	Communication string
	// The name of the call (i.e. name of function or some other target)
	MethodName string
	// Stream name of the message
	Subject string
	// The name of the service in which the call is made
	ServiceName string
	// The name of the file
	FileName string
	// Line of code
	PositionInFile string
}

func FindNATSCalls(serviceDir string) ([]*NatsCall, []*NatsCall, error) {
	consumers := make([]*NatsCall, 0)
	producers := make([]*NatsCall, 0)
	var err error
	config := defaultNatsConfig()

	files, err := os.ReadDir(serviceDir)
	if err != nil {
		return consumers, producers, err
	}

	for _, file := range files {
		if file.IsDir() {
			servicePath := filepath.Join(serviceDir, file.Name())
			fileErr := filepath.Walk(servicePath, func(path string, info fs.FileInfo, e error) error {
				if e != nil {
					fmt.Println(e)
				}

				if filepath.Ext(info.Name()) == ".go" && !strings.HasSuffix(info.Name(), "_test.go") {
					cons, prod := findDependencies(path, file.Name(), config)
					consumers = append(consumers, cons...)
					producers = append(producers, prod...)
				}

				return nil
			})

			if fileErr != nil {
				fmt.Println(fileErr)
			}
		}
	}

	return consumers, producers, nil
}

func findDependencies(servicePath string, serviceName string, config NatsAnalysisConfig) ([]*NatsCall, []*NatsCall) {
	producers := make([]*NatsCall, 0)
	consumers := make([]*NatsCall, 0)

	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, servicePath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(node ast.Node) bool {
		switch funcCall := node.(type) {
		case *ast.CallExpr:
			switch methodCall := funcCall.Fun.(type) {
			case *ast.SelectorExpr:
				methodName := methodCall.Sel.Name
				if strings.Contains(methodName, "NotifyMsg") {
					subject := findSubject(funcCall.Args)
					if subject != "" {
						producerCall := &NatsCall{
							MethodName:     methodName,
							Communication:  config.communication,
							Subject:        subject,
							ServiceName:    serviceName,
							FileName:       fs.Position(methodCall.X.Pos()).Filename,
							PositionInFile: strconv.Itoa(fs.Position(methodCall.Sel.Pos()).Line),
						}
						producers = append(producers, producerCall)
					}
				} else if strings.Contains(methodName, "Subscribe") {
					subject := findSubject(funcCall.Args)
					if subject != "" {
						consumerCall := &NatsCall{
							MethodName:     methodName,
							Communication:  config.communication,
							Subject:        subject,
							ServiceName:    serviceName,
							FileName:       fs.Position(methodCall.X.Pos()).Filename,
							PositionInFile: strconv.Itoa(fs.Position(methodCall.Sel.Pos()).Line),
						}

						consumers = append(consumers, consumerCall)
					}
				}
			default:
				return true
			}
		default:
			return true
		}
		return true
	})

	return consumers, producers
}

func findSubject(args []ast.Expr) string {
	for _, argument := range args {
		switch subjectArg := argument.(type) {
		case *ast.SelectorExpr:
			if strings.Contains(subjectArg.Sel.Name, "Subject") {
				return subjectArg.Sel.Name
			}
		default:
			continue
		}
	}

	return ""
}
