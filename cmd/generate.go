package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"diagram-gen/internal/generator"
	"diagram-gen/internal/model"
	"diagram-gen/internal/parser"
	"diagram-gen/internal/validator"
)

var newGenerator = func() generator.Formatter {
	return generator.NewDrawIOGenerator()
}

func buildGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [file or directory]",
		Short: "Generate a diagram from Go source code",
		Long: `Parses Go source files and generates a draw.io diagram based on 
diagram struct tags. 

Example:
  diagram-gen generate ./internal/services/
  diagram-gen generate main.go -o diagram.drawio`,
		Args: cobra.ExactArgs(1),
		RunE: generateRunE,
	}

	cmd.Flags().StringP("output", "o", "diagram.drawio", "Output file path")
	cmd.Flags().StringP("type", "t", "architecture", "Diagram type (architecture, flowchart, network)")
	return cmd
}

func generateRunE(cmd *cobra.Command, args []string) error {
	inputPath := args[0]
	outputPath, _ := cmd.Flags().GetString("output")
	diagramType, _ := cmd.Flags().GetString("type")

	if outputPath == "" {
		outputPath = "diagram.drawio"
	}

	p := parser.New()
	diagram, err := p.Parse(inputPath)
	if err != nil {
		return fmt.Errorf("failed to parse input: %w", err)
	}

	if len(diagram.Components) == 0 {
		return fmt.Errorf("no diagram annotations found in %s", inputPath)
	}

	if err := validator.ValidateDiagram(diagram); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	diagram.Type = model.DiagramType(diagramType)

	gen := newGenerator()
	data, err := gen.Generate(diagram)
	if err != nil {
		return fmt.Errorf("failed to generate diagram: %w", err)
	}

	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	fmt.Printf("Generated %s diagram with %d components and %d connections\n",
		diagramType, len(diagram.Components), len(diagram.Connections))
	fmt.Printf("Output written to: %s\n", outputPath)

	return nil
}

var generateCmd = buildGenerateCmd()
