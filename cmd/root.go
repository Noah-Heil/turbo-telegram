package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "diagram-gen",
	Short: "Generate software diagrams from code annotations",
	Long: `A CLI tool that parses Go source code annotations and generates 
draw.io compatible diagrams. Supports architecture, flowchart, and 
network diagram types.`,
}

// Execute runs the root command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(versionCmd)
}
