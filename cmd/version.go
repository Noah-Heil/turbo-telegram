package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of diagram-gen",
	Long:  `All software has versions. This is diagram-gen's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("diagram-gen version %s\n", version)
	},
}
