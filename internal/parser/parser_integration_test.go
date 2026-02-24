package parser

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	p := New()
	diagram, err := p.ParseFile("testdata/sample.go")
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(diagram.Components) == 0 {
		t.Error("expected components, got none")
	}

	if len(diagram.Connections) == 0 {
		t.Error("expected connections, got none")
	}

	t.Logf("Components: %d", len(diagram.Components))
	t.Logf("Connections: %d", len(diagram.Connections))

	for _, c := range diagram.Components {
		t.Logf("  Component: %s (%s)", c.Name, c.Type)
	}
	for _, c := range diagram.Connections {
		t.Logf("  Connection: %s -> %s", c.Source, c.Target)
	}
}
