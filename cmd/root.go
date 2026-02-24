package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "diagram-gen",
	Short: "Generate software diagrams from code annotations",
	Long: `A CLI tool that parses Go source code annotations and generates 
draw.io compatible diagrams. Supports architecture, flowchart, and 
network diagram types.`,
}

var exitFunc = os.Exit

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(versionCmd)
}
