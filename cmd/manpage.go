package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
	"path/filepath"
)

// GenManpageCmd returns a cobra command that
// generates a manpage for the specified targetCommand
func GenManpageCmd(targetCommand *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genManPage",
		Short: "Generate a manpage entry for netDep",
		Long: `Outputs a manpage file in the current directory. 
The filename is netDep.1, where 1 stands for manPages that relate to "executable shell commands".`,

		RunE: func(cmd *cobra.Command, args []string) error {

			header := &doc.GenManHeader{
				Title:   "netDep",
				Section: "1",
			}
			currentWd, err := os.Getwd()
			if err != nil {
				return err
			}
			err = doc.GenManTree(targetCommand, header, currentWd)
			if err != nil {
				return err
			}
			fmt.Printf("successfully generated manpage entry: %s\n", filepath.Join(currentWd, "netDep.1"))
			return nil
		},
	}
	return cmd
}
