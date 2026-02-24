package validator

import (
	"testing"

	"diagram-gen/internal/model"
)

func TestValidateDiagram(t *testing.T) {
	tests := []struct {
		name    string
		diagram *model.Diagram
		wantErr bool
	}{
		{
			name: "valid diagram",
			diagram: &model.Diagram{
				Components: []model.Component{
					{Type: model.ComponentTypeService, Name: "ServiceA"},
					{Type: model.ComponentTypeDatabase, Name: "DatabaseA"},
				},
				Connections: []model.Connection{
					{Source: "ServiceA", Target: "DatabaseA"},
				},
			},
			wantErr: false,
		},
		{
			name:    "nil diagram",
			diagram: nil,
			wantErr: true,
		},
		{
			name: "empty components",
			diagram: &model.Diagram{
				Components: []model.Component{},
			},
			wantErr: true,
		},
		{
			name: "empty component name",
			diagram: &model.Diagram{
				Components: []model.Component{
					{Type: model.ComponentTypeService, Name: ""},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicate component names",
			diagram: &model.Diagram{
				Components: []model.Component{
					{Type: model.ComponentTypeService, Name: "ServiceA"},
					{Type: model.ComponentTypeDatabase, Name: "ServiceA"},
				},
			},
			wantErr: true,
		},
		{
			name: "unknown component type",
			diagram: &model.Diagram{
				Components: []model.Component{
					{Type: model.ComponentType("unknown"), Name: "ServiceA"},
				},
			},
			wantErr: true,
		},
		{
			name: "connection to unknown source",
			diagram: &model.Diagram{
				Components: []model.Component{
					{Type: model.ComponentTypeService, Name: "ServiceA"},
				},
				Connections: []model.Connection{
					{Source: "UnknownService", Target: "ServiceA"},
				},
			},
			wantErr: true,
		},
		{
			name: "connection to unknown target",
			diagram: &model.Diagram{
				Components: []model.Component{
					{Type: model.ComponentTypeService, Name: "ServiceA"},
				},
				Connections: []model.Connection{
					{Source: "ServiceA", Target: "UnknownDB"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty connection source",
			diagram: &model.Diagram{
				Components: []model.Component{
					{Type: model.ComponentTypeService, Name: "ServiceA"},
				},
				Connections: []model.Connection{
					{Source: "", Target: "ServiceA"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDiagram(tt.diagram)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDiagram() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateComponentType(t *testing.T) {
	tests := []struct {
		componentType model.ComponentType
		want          bool
	}{
		{model.ComponentTypeService, true},
		{model.ComponentTypeDatabase, true},
		{model.ComponentTypeQueue, true},
		{model.ComponentTypeCache, true},
		{model.ComponentTypeAPI, true},
		{model.ComponentTypeUser, true},
		{model.ComponentTypeExternal, true},
		{model.ComponentTypeStorage, true},
		{model.ComponentTypeGateway, true},
		{model.ComponentType("unknown"), false},
		{model.ComponentType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.componentType), func(t *testing.T) {
			if got := ValidateComponentType(tt.componentType); got != tt.want {
				t.Errorf("ValidateComponentType(%q) = %v, want %v", tt.componentType, got, tt.want)
			}
		})
	}
}
