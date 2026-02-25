// Package cmd defines CLI commands.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"diagram-gen/internal/archparser"
	"diagram-gen/internal/generator"
	"diagram-gen/internal/model"
	"diagram-gen/internal/validator"
)

var (
	flagLayout    string
	flagIsometric bool
	flagCompress  bool
	flagShape     string
	flagConfig    string
	flagPage      string
)

func newGeneratorWithFlags() generator.Formatter {
	gen := generator.NewDrawIOGenerator()
	if flagIsometric {
		gen.LayoutType = "isometric"
	} else if flagLayout != "" {
		gen.LayoutType = flagLayout
	}
	gen.Compress = flagCompress
	return gen
}

var newGenerator = newGeneratorWithFlags

func buildGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [file or directory]",
		Short: "Generate a diagram from Go source code",
		Long: `Parses Go source files and generates a draw.io diagram based on 
diagram struct tags. 

Example:
  diagram-gen generate ./internal/services/
  diagram-gen generate main.go -o diagram.drawio
  diagram-gen generate main.go --layout isometric --compress`,
		Args: cobra.ExactArgs(1),
		RunE: generateRunE,
	}

	cmd.Flags().StringP("output", "o", "diagram.drawio", "Output file path")
	cmd.Flags().StringP("type", "t", "architecture", "Diagram type (architecture, flowchart, network)")
	cmd.Flags().StringVar(&flagLayout, "layout", "layered", "Layout type: grid, layered, isometric")
	cmd.Flags().BoolVar(&flagIsometric, "isometric", false, "Use isometric layout (shortcut for --layout isometric)")
	cmd.Flags().BoolVar(&flagCompress, "compress", false, "Compress output with deflate+base64")
	cmd.Flags().StringVar(&flagShape, "shape", "", "Default shape for components (e.g., iso:server, rounded, cylinder)")
	cmd.Flags().StringVar(&flagConfig, "config", "", "Path to config file (.diagram-gen.yaml or .diagram-gen.json)")
	cmd.Flags().StringVar(&flagPage, "page", "", "Page name to generate (for multi-page diagrams)")
	return cmd
}

func generateRunE(cmd *cobra.Command, args []string) error {
	inputPath := args[0]
	outputPath, _ := cmd.Flags().GetString("output")
	diagramType, _ := cmd.Flags().GetString("type")

	if outputPath == "" {
		outputPath = "diagram.drawio"
	}

	p := archparser.New()
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

	if flagIsometric {
		diagram.Layout = "isometric"
	} else if flagLayout != "" {
		diagram.Layout = flagLayout
	}

	diagram.Compress = flagCompress

	if flagPage != "" {
		filteredComps := []model.Component{}
		filteredConns := []model.Connection{}
		for _, comp := range diagram.Components {
			if comp.Page == "" || comp.Page == flagPage {
				filteredComps = append(filteredComps, comp)
			}
		}
		for _, conn := range diagram.Connections {
			if conn.Page == "" || conn.Page == flagPage {
				filteredConns = append(filteredConns, conn)
			}
		}
		diagram.Components = filteredComps
		diagram.Connections = filteredConns
	}

	gen := newGenerator()
	data, err := gen.Generate(diagram)
	if err != nil {
		return fmt.Errorf("failed to generate diagram: %w", err)
	}

	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	layoutType := diagram.Layout
	if layoutType == "" {
		layoutType = "layered"
	}

	fmt.Printf("Generated %s diagram (%s layout) with %d components and %d connections\n",
		diagramType, layoutType, len(diagram.Components), len(diagram.Connections))
	fmt.Printf("Output written to: %s\n", outputPath)

	if flagCompress {
		fmt.Println("Output compressed with deflate+base64")
	}

	return nil
}

var generateCmd = buildGenerateCmd()
