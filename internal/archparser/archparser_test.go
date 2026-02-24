package archparser_test

import (
	"testing"

	"diagram-gen/internal/archparser"
	"diagram-gen/internal/model"
)

func TestParseAnnotation(t *testing.T) {
	tests := []struct {
		name      string
		tag       string
		wantName  string
		wantType  model.ComponentType
		wantConns int
		wantErr   bool
	}{
		{
			name:      "basic annotation",
			tag:       `type=service,name=UserService,connectsTo=Database`,
			wantName:  "UserService",
			wantType:  model.ComponentTypeService,
			wantConns: 1,
		},
		{
			name:      "multiple connections",
			tag:       `type=api,name=APIGateway,connectsTo=AuthService;PaymentService`,
			wantName:  "APIGateway",
			wantType:  model.ComponentTypeAPI,
			wantConns: 2,
		},
		{
			name:      "with description",
			tag:       `type=database,name=PostgresDB,description=User data store`,
			wantName:  "PostgresDB",
			wantType:  model.ComponentTypeDatabase,
			wantConns: 0,
		},
		{
			name:    "missing name",
			tag:     `type=service`,
			wantErr: true,
		},
		{
			name:      "default type",
			tag:       `name=MyService`,
			wantName:  "MyService",
			wantType:  model.ComponentTypeService,
			wantConns: 0,
		},
		{
			name:      "with direction",
			tag:       `type=service,name=A,connectsTo=B,direction=bidirectional`,
			wantName:  "A",
			wantType:  model.ComponentTypeService,
			wantConns: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ann, err := archparser.ParseAnnotation(tt.tag)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if ann.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", ann.Name, tt.wantName)
			}
			if ann.ComponentType != tt.wantType {
				t.Errorf("type = %q, want %q", ann.ComponentType, tt.wantType)
			}
			if len(ann.ConnectsTo) != tt.wantConns {
				t.Errorf("connections = %d, want %d", len(ann.ConnectsTo), tt.wantConns)
			}
		})
	}
}

func TestAnnotationToComponent(t *testing.T) {
	ann, _ := archparser.ParseAnnotation(`type=database,name=Redis,description=Cache layer`)
	comp := ann.ToComponent()

	if comp.Name != "Redis" {
		t.Errorf("Name = %q, want %q", comp.Name, "Redis")
	}
	if comp.Type != model.ComponentTypeDatabase {
		t.Errorf("Type = %q, want %q", comp.Type, model.ComponentTypeDatabase)
	}
	if comp.Description != "Cache layer" {
		t.Errorf("Description = %q, want %q", comp.Description, "Cache layer")
	}
}

func TestAnnotationToConnections(t *testing.T) {
	ann, _ := archparser.ParseAnnotation(`type=service,name=A,connectsTo=B;C`)
	conns := ann.ToConnections()

	if len(conns) != 2 {
		t.Fatalf("len(connections) = %d, want 2", len(conns))
	}
	if conns[0].Source != "A" || conns[0].Target != "B" {
		t.Errorf("connection 0 = (%q, %q), want (A, B)", conns[0].Source, conns[0].Target)
	}
	if conns[1].Source != "A" || conns[1].Target != "C" {
		t.Errorf("connection 1 = (%q, %q), want (A, C)", conns[1].Source, conns[1].Target)
	}
}

func TestAnnotationToConnectionsEmpty(t *testing.T) {
	ann, _ := archparser.ParseAnnotation(`type=service,name=A`)
	conns := ann.ToConnections()

	if len(conns) != 0 {
		t.Errorf("len(connections) = %d, want 0", len(conns))
	}
}

func TestAnnotationToConnectionsWithDirection(t *testing.T) {
	ann, _ := archparser.ParseAnnotation(`type=service,name=A,connectsTo=B,direction=bidirectional`)
	conns := ann.ToConnections()

	if len(conns) != 1 {
		t.Fatalf("len(connections) = %d, want 1", len(conns))
	}
	if conns[0].Direction != model.ConnectionDirectionBidirectional {
		t.Errorf("Direction = %q, want %q", conns[0].Direction, model.ConnectionDirectionBidirectional)
	}
}

func TestAnnotationToConnectionsDefaultDirection(t *testing.T) {
	ann, _ := archparser.ParseAnnotation(`type=service,name=A,connectsTo=B`)
	conns := ann.ToConnections()

	if conns[0].Direction != model.ConnectionDirectionUnidirectional {
		t.Errorf("Direction = %q, want %q", conns[0].Direction, model.ConnectionDirectionUnidirectional)
	}
}

func TestParseAnnotationMultipleFields(t *testing.T) {
	ann, err := archparser.ParseAnnotation(`type=service,name=S,connectsTo=A;B,description=test desc,direction=bidirectional`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ann.Name != "S" {
		t.Errorf("name = %q, want %q", ann.Name, "S")
	}
	if len(ann.ConnectsTo) != 2 {
		t.Errorf("connectsTo = %v, want 2 elements", ann.ConnectsTo)
	}
	if ann.Description != "test desc" {
		t.Errorf("description = %q, want %q", ann.Description, "test desc")
	}
	if ann.Direction != model.ConnectionDirectionBidirectional {
		t.Errorf("direction = %q, want %q", ann.Direction, model.ConnectionDirectionBidirectional)
	}
}

func TestParseAnnotationAllComponentTypes(t *testing.T) {
	types := []string{"service", "database", "queue", "cache", "api", "user", "external", "storage", "gateway"}
	for _, compType := range types {
		ann, err := archparser.ParseAnnotation("type=" + compType + ",name=Test")
		if err != nil {
			t.Errorf("unexpected error for type %s: %v", compType, err)
		}
		if ann.ComponentType != model.ComponentType(compType) {
			t.Errorf("type = %q, want %q", ann.ComponentType, compType)
		}
	}
}

func TestParseAnnotationWithColonInValue(t *testing.T) {
	ann, err := archparser.ParseAnnotation(`type=service,name=A,description=has:colon`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ann.Description != "has:colon" {
		t.Errorf("description = %q, want %q", ann.Description, "has:colon")
	}
}

func TestParseAnnotationDiagramOnly(t *testing.T) {
	_, err := archparser.ParseAnnotation("diagram")
	if err == nil {
		t.Error("expected error for diagram-only tag")
	}
}

func TestParseAnnotationWithDiagramPrefix(t *testing.T) {
	ann, err := archparser.ParseAnnotation("diagram,name=Test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ann.Name != "Test" {
		t.Errorf("name = %q, want %q", ann.Name, "Test")
	}
}

func TestParseAnnotationInvalidKV(t *testing.T) {
	_, err := archparser.ParseAnnotation(`type`)
	if err == nil {
		t.Error("expected error for invalid key-value")
	}
}

func TestParseAnnotationEmpty(t *testing.T) {
	_, err := archparser.ParseAnnotation("")
	if err == nil {
		t.Error("expected error for empty tag")
	}
}

func TestParseAnnotationEmptyName(t *testing.T) {
	_, err := archparser.ParseAnnotation(`type=service,name=`)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestParseAnnotationSpacesInValues(t *testing.T) {
	ann, err := archparser.ParseAnnotation(`type = service , name = TestService `)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ann.Name != "TestService" {
		t.Errorf("name = %q, want %q", ann.Name, "TestService")
	}
}

func TestAnnotationToComponentWithDirection(t *testing.T) {
	ann, _ := archparser.ParseAnnotation(`type=service,name=A,direction=bidirectional`)
	comp := ann.ToComponent()
	if comp.Direction != model.ConnectionDirectionBidirectional {
		t.Errorf("Direction = %q, want %q", comp.Direction, model.ConnectionDirectionBidirectional)
	}
}

func TestAnnotationToConnectionsWhitespace(t *testing.T) {
	ann, _ := archparser.ParseAnnotation(`type=service,name=A,connectsTo= B ; C `)
	conns := ann.ToConnections()
	if len(conns) != 2 {
		t.Fatalf("len = %d, want 2", len(conns))
	}
	if conns[0].Target != "B" || conns[1].Target != "C" {
		t.Errorf("unexpected targets: %v", conns)
	}
}

func TestAnnotationToConnectionsEmptyTarget(t *testing.T) {
	ann, _ := archparser.ParseAnnotation(`type=service,name=A,connectsTo=;B;`)
	conns := ann.ToConnections()
	if len(conns) != 1 {
		t.Fatalf("len = %d, want 1 (empty strings filtered)", len(conns))
	}
}
