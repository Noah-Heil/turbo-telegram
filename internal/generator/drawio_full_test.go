package generator_test

import (
	"strings"
	"testing"

	"diagram-gen/internal/generator"
	"diagram-gen/internal/generator/layout"
	"diagram-gen/internal/model"
)

func TestGetEdgeStyle(t *testing.T) {
	tests := []struct {
		direction model.ConnectionDirection
		want      string
	}{
		{model.ConnectionDirectionUnidirectional, "endArrow"},
		{model.ConnectionDirectionBidirectional, "endArrow"},
		{"", "endArrow"},
	}

	for _, tt := range tests {
		t.Run(string(tt.direction), func(t *testing.T) {
			conn := model.Connection{Source: "A", Target: "B", Direction: tt.direction}
			g := &generator.DrawIOGenerator{}
			got := g.BuildEdgeStyle(conn)
			if !strings.Contains(got, tt.want) {
				t.Errorf("BuildEdgeStyle(%q) = %q, should contain %q", tt.direction, got, tt.want)
			}
		})
	}
}

func TestCalculateLayout(t *testing.T) {
	tests := []struct {
		name       string
		components []model.Component
	}{
		{
			name: "single component",
			components: []model.Component{
				{Type: model.ComponentTypeService, Name: "Service1"},
			},
		},
		{
			name: "four components",
			components: []model.Component{
				{Type: model.ComponentTypeService, Name: "S1"},
				{Type: model.ComponentTypeService, Name: "S2"},
				{Type: model.ComponentTypeService, Name: "S3"},
				{Type: model.ComponentTypeService, Name: "S4"},
			},
		},
		{
			name: "seven components",
			components: []model.Component{
				{Type: model.ComponentTypeService, Name: "S1"},
				{Type: model.ComponentTypeService, Name: "S2"},
				{Type: model.ComponentTypeService, Name: "S3"},
				{Type: model.ComponentTypeService, Name: "S4"},
				{Type: model.ComponentTypeService, Name: "S5"},
				{Type: model.ComponentTypeService, Name: "S6"},
				{Type: model.ComponentTypeService, Name: "S7"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layoutEngine := layout.NewLayout("grid")
			posMap := layoutEngine.Calculate(tt.components, nil)
			if len(posMap) != len(tt.components) {
				t.Errorf("layout length = %d, want %d", len(posMap), len(tt.components))
			}
			for _, pos := range posMap {
				if pos.X <= 0 || pos.Y <= 0 {
					t.Errorf("invalid position: %v", pos)
				}
			}
		})
	}
}

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<hello>", "&lt;hello&gt;"},
		{"hello&world", "hello&amp;world"},
		{`"quote"`, "&quot;quote&quot;"},
		{"'single'", "&apos;single&apos;"},
		{"plain", "plain"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := generator.EscapeXML(tt.input)
			if got != tt.expected {
				t.Errorf("EscapeXML(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGenerateWithManyComponents(t *testing.T) {
	gen := generator.NewDrawIOGenerator()

	components := make([]model.Component, 20)
	for i := range components {
		components[i] = model.Component{
			Type: model.ComponentTypeService,
			Name: "Service" + string(rune('A'+i)),
		}
	}

	diagram := &model.Diagram{
		Type:        model.DiagramTypeArchitecture,
		Components:  components,
		Connections: []model.Connection{},
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content := string(data)
	for i := range components {
		expected := "Service" + string(rune('A'+i))
		if !strings.Contains(content, expected) {
			t.Errorf("expected %s in output", expected)
		}
	}
}

func TestGenerateWithConnections(t *testing.T) {
	gen := generator.NewDrawIOGenerator()

	diagram := &model.Diagram{
		Type: model.DiagramTypeArchitecture,
		Components: []model.Component{
			{Type: model.ComponentTypeService, Name: "A"},
			{Type: model.ComponentTypeService, Name: "B"},
			{Type: model.ComponentTypeService, Name: "C"},
		},
		Connections: []model.Connection{
			{Source: "A", Target: "B"},
			{Source: "B", Target: "C"},
			{Source: "C", Target: "A"},
		},
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "source=") || !strings.Contains(content, "target=") {
		t.Error("expected source and target in edge elements")
	}
}

func TestGenerateEmptyDiagram(t *testing.T) {
	gen := generator.NewDrawIOGenerator()

	diagram := &model.Diagram{
		Type:        model.DiagramTypeArchitecture,
		Components:  []model.Component{},
		Connections: []model.Connection{},
	}

	_, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
}

func TestGenerateMultipleComponentTypes(t *testing.T) {
	gen := generator.NewDrawIOGenerator()

	diagram := &model.Diagram{
		Type: model.DiagramTypeArchitecture,
		Components: []model.Component{
			{Type: model.ComponentTypeService, Name: "Service"},
			{Type: model.ComponentTypeDatabase, Name: "Database"},
			{Type: model.ComponentTypeQueue, Name: "Queue"},
			{Type: model.ComponentTypeCache, Name: "Cache"},
			{Type: model.ComponentTypeUser, Name: "User"},
			{Type: model.ComponentTypeExternal, Name: "External"},
			{Type: model.ComponentTypeStorage, Name: "Storage"},
			{Type: model.ComponentTypeAPI, Name: "API"},
			{Type: model.ComponentTypeGateway, Name: "Gateway"},
		},
		Connections: []model.Connection{},
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content := string(data)
	for _, name := range []string{"Service", "Database", "Queue", "Cache", "User", "External", "Storage", "API", "Gateway"} {
		if !strings.Contains(content, name) {
			t.Errorf("expected %s in output", name)
		}
	}
}

func TestGenerateWithInvalidConnection(t *testing.T) {
	gen := generator.NewDrawIOGenerator()

	diagram := &model.Diagram{
		Type: model.DiagramTypeArchitecture,
		Components: []model.Component{
			{Type: model.ComponentTypeService, Name: "A"},
		},
		Connections: []model.Connection{
			{Source: "A", Target: "Unknown"},
			{Source: "Unknown2", Target: "A"},
		},
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "A") {
		t.Error("expected A in output")
	}
}

func TestCalculateLayoutEdgeCases(t *testing.T) {
	components := []model.Component{
		{Type: model.ComponentTypeService, Name: "S1"},
		{Type: model.ComponentTypeService, Name: "S2"},
		{Type: model.ComponentTypeService, Name: "S3"},
		{Type: model.ComponentTypeService, Name: "S4"},
		{Type: model.ComponentTypeService, Name: "S5"},
	}
	layoutEngine := layout.NewLayout("grid")
	posMap := layoutEngine.Calculate(components, nil)
	if len(posMap) != 5 {
		t.Errorf("layout length = %d, want 5", len(posMap))
	}
}
