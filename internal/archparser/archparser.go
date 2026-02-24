package archparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"diagram-gen/internal/model"
)

// Parser parses Go source files for diagram annotations.
type Parser struct {
	fset *token.FileSet
}

// New creates a new Parser.
func New() *Parser {
	return &Parser{
		fset: token.NewFileSet(),
	}
}

// ParseFile parses a single Go file for diagram annotations.
func (p *Parser) ParseFile(path string) (*model.Diagram, error) {
	f, err := parser.ParseFile(p.fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", path, err)
	}

	diagram := &model.Diagram{
		Type:        model.DiagramTypeArchitecture,
		Components:  []model.Component{},
		Connections: []model.Connection{},
	}

	ast.Inspect(f, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range structType.Fields.List {
			if field.Tag == nil {
				continue
			}

			tag := ParseStructTag(field.Tag.Value, "diagram")
			if tag == "" {
				continue
			}

			ann, err := ParseAnnotation(tag)
			if err != nil {
				continue
			}

			component := ann.ToComponent()
			diagram.AddComponent(component)

			connections := ann.ToConnections()
			for _, conn := range connections {
				diagram.AddConnection(conn)
			}
		}

		return true
	})

	return diagram, nil
}

// ParseDirectory parses all Go files in a directory.
func (p *Parser) ParseDirectory(dirPath string) (*model.Diagram, error) {
	diagram := &model.Diagram{
		Type:        model.DiagramTypeArchitecture,
		Components:  []model.Component{},
		Connections: []model.Connection{},
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		fileDiagram, err := p.ParseFile(filePath)
		if err != nil {
			continue
		}

		diagram.Components = append(diagram.Components, fileDiagram.Components...)
		diagram.Connections = append(diagram.Connections, fileDiagram.Connections...)
	}

	return diagram, nil
}

// Parse parses a file or directory for diagram annotations.
func (p *Parser) Parse(inputPath string) (*model.Diagram, error) {
	info, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("cannot access input path: %w", err)
	}

	if info.IsDir() {
		return p.ParseDirectory(inputPath)
	}

	return p.ParseFile(inputPath)
}

func ParseStructTag(tagValue, key string) string {
	tagValue = tagValue[1 : len(tagValue)-1]
	tagValue = strings.ReplaceAll(tagValue, `\"`, `"`)

	tag := strings.TrimSpace(tagValue)
	if strings.HasPrefix(tag, key+":") {
		tag = strings.TrimPrefix(tag, key+":")
		tag = strings.Trim(tag, "\"")
		return tag
	}
	return ""
}
