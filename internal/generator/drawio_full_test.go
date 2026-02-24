package generator

import (
	"strings"
	"testing"

	"diagram-gen/internal/model"
)

func TestDrawIOGenerator(t *testing.T) {
	gen := NewDrawIOGenerator()

	diagram := &model.Diagram{
		Type: model.DiagramTypeArchitecture,
		Components: []model.Component{
			{Type: model.ComponentTypeService, Name: "UserService"},
			{Type: model.ComponentTypeDatabase, Name: "UserDB"},
			{Type: model.ComponentTypeAPI, Name: "API"},
		},
		Connections: []model.Connection{
			{Source: "UserService", Target: "UserDB"},
			{Source: "API", Target: "UserService"},
		},
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "<?xml") {
		t.Error("expected XML declaration")
	}
	if !strings.Contains(content, "UserService") {
		t.Error("expected UserService in output")
	}
	if !strings.Contains(content, "UserDB") {
		t.Error("expected UserDB in output")
	}
	if !strings.Contains(content, "mxCell") {
		t.Error("expected mxCell elements")
	}

	t.Logf("Generated %d bytes", len(data))
}

func TestDrawIOGeneratorFormat(t *testing.T) {
	gen := NewDrawIOGenerator()
	if gen.Format() != "drawio" {
		t.Errorf("Format() = %q, want %q", gen.Format(), "drawio")
	}
}

func TestGetShapeStyle(t *testing.T) {
	tests := []struct {
		compType model.ComponentType
		want     string
	}{
		{model.ComponentTypeService, "rounded=1"},
		{model.ComponentTypeDatabase, "shape=cylinder"},
		{model.ComponentTypeQueue, "shape=parallelogram"},
		{model.ComponentTypeCache, "dashed=1"},
		{model.ComponentTypeUser, "ellipse"},
		{model.ComponentTypeExternal, "shape=document"},
		{model.ComponentTypeStorage, "shape=cylinder"},
		{model.ComponentTypeAPI, "rounded=1"},
		{model.ComponentTypeGateway, "rounded=1"},
		{model.ComponentTypeUnknown, "rounded=1"},
	}

	for _, tt := range tests {
		t.Run(string(tt.compType), func(t *testing.T) {
			got := getShapeStyle(tt.compType)
			if !strings.Contains(got, tt.want) {
				t.Errorf("getShapeStyle(%q) = %q, should contain %q", tt.compType, got, tt.want)
			}
		})
	}
}

func TestGetEdgeStyle(t *testing.T) {
	tests := []struct {
		direction model.ConnectionDirection
		want      string
	}{
		{model.ConnectionDirectionUnidirectional, "endArrow=classic"},
		{model.ConnectionDirectionBidirectional, "endArrow=classic"},
		{"", "endArrow=classic"},
	}

	for _, tt := range tests {
		t.Run(string(tt.direction), func(t *testing.T) {
			got := getEdgeStyle(tt.direction)
			if !strings.Contains(got, tt.want) {
				t.Errorf("getEdgeStyle(%q) = %q, should contain %q", tt.direction, got, tt.want)
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
			layout := calculateLayout(tt.components)
			if len(layout) != len(tt.components) {
				t.Errorf("layout length = %d, want %d", len(layout), len(tt.components))
			}
			for _, pos := range layout {
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
			got := escapeXML(tt.input)
			if got != tt.expected {
				t.Errorf("escapeXML(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGenerateWithManyComponents(t *testing.T) {
	gen := NewDrawIOGenerator()

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
	gen := NewDrawIOGenerator()

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
	gen := NewDrawIOGenerator()

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
	gen := NewDrawIOGenerator()

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
	gen := NewDrawIOGenerator()

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
	layout := calculateLayout(components)
	if len(layout) != 5 {
		t.Errorf("layout length = %d, want 5", len(layout))
	}
}
