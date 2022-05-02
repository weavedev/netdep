package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var funcs = map[string]bool{"NewRequest": true, "NewClient": true, "Post": true, "Get": true, "NewAuthClient": true,
	"FetchClient": true, "FetchSchema": true, "ValidateQueryModel": true, "Do": true}

func printServices(path string, funcs map[string]bool) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Println("SERVICE: " + file.Name())
			servicePath := path + file.Name()

			filepath.Walk(servicePath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Fatalf(err.Error())
				}
				if info.Name()[len(info.Name())-3:] == ".go" && !strings.Contains(info.Name(), "_test") {
					fmt.Println("-File Name:", info.Name())
					fileToAST(path, funcs)
				}
				return nil
			})
		}
		fmt.Println("------------------")
	}
}

func fileToAST(path string, funcs map[string]bool) {

	//Create a FileSet to work with
	fset := token.NewFileSet()
	//Parse the file and create an AST
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ast.Inspect(file, func(n ast.Node) bool {
		// Find Function Call Statements
		funcCall, ok := n.(*ast.CallExpr)
		if ok {
			i, ok := funcCall.Fun.(*ast.SelectorExpr)
			if ok {
				_, ok := funcs[i.Sel.Name]
				if ok {
					//fmt.Println(i.Sel.Name)
					fmt.Print("--Func: " + i.Sel.Name)
					fmt.Println(funcCall.Args)
				}
			}
		}
		return true
	})
}
