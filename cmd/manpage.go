package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// genManPageToDir is a private function that generates the manPage
// and writes it to a file in the specified directory
func genManPageToDir(targetCommand *cobra.Command, outputDir string) error {
	header := &doc.GenManHeader{
		Title:   "netDep",
		Section: "1",
	}

	err := doc.GenManTree(targetCommand, header, outputDir)
	if err != nil {
		return err
	}
	fmt.Printf("successfully generated manpage entry: %s\n", filepath.Join(outputDir, "netDep.1"))
	return nil
}

// GenManpageCmd returns a cobra command that
// generates a manpage for the specified targetCommand
// in the current working directory
func GenManpageCmd(targetCommand *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genManPage",
		Short: "Generate a manpage entry for netDep",
		Long: `Outputs a manpage file in the current directory. 
The filename is netDep.1, where 1 stands for manPages that relate to "executable shell commands".`,

		RunE: func(cmd *cobra.Command, args []string) error {
			currentWd, err := os.Getwd()
			if err != nil {
				return err
			}
			return genManPageToDir(targetCommand, currentWd)
		},
	}
	return cmd
}
