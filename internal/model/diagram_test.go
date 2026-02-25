package model_test

import (
	"diagram-gen/internal/model"
	"testing"
)

func TestDiagramAddComponent(t *testing.T) {
	t.Parallel()
	d := &model.Diagram{
		Type:        model.DiagramTypeArchitecture,
		Components:  []model.Component{},
		Connections: []model.Connection{},
	}

	comp := model.Component{
		Type: model.ComponentTypeService,
		Name: "TestService",
	}
	d.AddComponent(comp)

	if len(d.Components) != 1 {
		t.Errorf("expected 1 component, got %d", len(d.Components))
	}
	if d.Components[0].Name != "TestService" {
		t.Errorf("expected component name 'TestService', got '%s'", d.Components[0].Name)
	}
}

func TestDiagramAddConnection(t *testing.T) {
	t.Parallel()
	d := &model.Diagram{
		Type:        model.DiagramTypeArchitecture,
		Components:  []model.Component{},
		Connections: []model.Connection{},
	}

	conn := model.Connection{
		Source: "ServiceA",
		Target: "ServiceB",
	}
	d.AddConnection(conn)

	if len(d.Connections) != 1 {
		t.Errorf("expected 1 connection, got %d", len(d.Connections))
	}
	if d.Connections[0].Source != "ServiceA" {
		t.Errorf("expected source 'ServiceA', got '%s'", d.Connections[0].Source)
	}
}

func TestDiagramGetComponentByName(t *testing.T) {
	t.Parallel()
	d := &model.Diagram{
		Type: model.DiagramTypeArchitecture,
		Components: []model.Component{
			{Type: model.ComponentTypeService, Name: "ServiceA"},
			{Type: model.ComponentTypeDatabase, Name: "DatabaseA"},
		},
	}

	comp := d.GetComponentByName("ServiceA")
	if comp == nil {
		t.Error("expected to find ServiceA")
		return
	}
	if comp.Type != model.ComponentTypeService {
		t.Errorf("expected type service, got %s", comp.Type)
	}

	notFound := d.GetComponentByName("Unknown")
	if notFound != nil {
		t.Error("expected nil for unknown component")
	}
}

func TestDiagramTypes(t *testing.T) {
	t.Parallel()
	if model.DiagramTypeArchitecture != "architecture" {
		t.Errorf("expected architecture, got %s", model.DiagramTypeArchitecture)
	}
	if model.DiagramTypeFlowchart != "flowchart" {
		t.Errorf("expected flowchart, got %s", model.DiagramTypeFlowchart)
	}
	if model.DiagramTypeNetwork != "network" {
		t.Errorf("expected network, got %s", model.DiagramTypeNetwork)
	}
}

func TestComponentTypes(t *testing.T) {
	t.Parallel()
	types := []model.ComponentType{
		model.ComponentTypeService,
		model.ComponentTypeDatabase,
		model.ComponentTypeQueue,
		model.ComponentTypeCache,
		model.ComponentTypeAPI,
		model.ComponentTypeUser,
		model.ComponentTypeExternal,
		model.ComponentTypeStorage,
		model.ComponentTypeGateway,
	}

	expected := []string{
		"service",
		"database",
		"queue",
		"cache",
		"api",
		"user",
		"external",
		"storage",
		"gateway",
	}

	for i, ct := range types {
		if string(ct) != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], ct)
		}
	}
}
