package model

import (
	"testing"
)

func TestDiagramAddComponent(t *testing.T) {
	d := &Diagram{
		Type:        DiagramTypeArchitecture,
		Components:  []Component{},
		Connections: []Connection{},
	}

	comp := Component{
		Type: ComponentTypeService,
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
	d := &Diagram{
		Type:        DiagramTypeArchitecture,
		Components:  []Component{},
		Connections: []Connection{},
	}

	conn := Connection{
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
	d := &Diagram{
		Type: DiagramTypeArchitecture,
		Components: []Component{
			{Type: ComponentTypeService, Name: "ServiceA"},
			{Type: ComponentTypeDatabase, Name: "DatabaseA"},
		},
	}

	comp := d.GetComponentByName("ServiceA")
	if comp == nil {
		t.Error("expected to find ServiceA")
	}
	if comp.Type != ComponentTypeService {
		t.Errorf("expected type service, got %s", comp.Type)
	}

	notFound := d.GetComponentByName("Unknown")
	if notFound != nil {
		t.Error("expected nil for unknown component")
	}
}

func TestDiagramTypes(t *testing.T) {
	if DiagramTypeArchitecture != "architecture" {
		t.Errorf("expected architecture, got %s", DiagramTypeArchitecture)
	}
	if DiagramTypeFlowchart != "flowchart" {
		t.Errorf("expected flowchart, got %s", DiagramTypeFlowchart)
	}
	if DiagramTypeNetwork != "network" {
		t.Errorf("expected network, got %s", DiagramTypeNetwork)
	}
}

func TestComponentTypes(t *testing.T) {
	types := []ComponentType{
		ComponentTypeService,
		ComponentTypeDatabase,
		ComponentTypeQueue,
		ComponentTypeCache,
		ComponentTypeAPI,
		ComponentTypeUser,
		ComponentTypeExternal,
		ComponentTypeStorage,
		ComponentTypeGateway,
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
