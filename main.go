/*
Copyright © 2022 Team 1, Weave BV, TU Delft

*/
package main

import (
	"fmt"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/cmd"
	stages "lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
)

func main() {
	fmt.Println("haha")
	stages.FindCallersForEndpoint("", "", "")
	cmd.Execute()
}
