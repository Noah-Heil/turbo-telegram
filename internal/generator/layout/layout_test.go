package layout_test

import (
	"testing"

	"diagram-gen/internal/generator/layout"
	"diagram-gen/internal/model"
)

func TestGridLayout(t *testing.T) {
	l := &layout.GridLayout{}

	if l.Name() != "grid" {
		t.Errorf("Name() = %q, want 'grid'", l.Name())
	}

	components := []model.Component{
		{Name: "A"}, {Name: "B"}, {Name: "C"},
	}

	pos := l.Calculate(components, nil)

	if len(pos) != 3 {
		t.Errorf("Calculate() returned %d positions, want 3", len(pos))
	}

	for name, p := range pos {
		if p.X <= 0 || p.Y <= 0 {
			t.Errorf("invalid position for %s: %v", name, p)
		}
	}
}

func TestGridLayoutFiveComponents(t *testing.T) {
	l := &layout.GridLayout{}
	components := []model.Component{{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}, {Name: "E"}}

	pos := l.Calculate(components, nil)
	if len(pos) != 5 {
		t.Errorf("Calculate() returned %d positions, want 5", len(pos))
	}
}

func TestGridLayoutSevenComponents(t *testing.T) {
	l := &layout.GridLayout{}
	components := []model.Component{{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}, {Name: "E"}, {Name: "F"}, {Name: "G"}}

	pos := l.Calculate(components, nil)
	if len(pos) != 7 {
		t.Errorf("Calculate() returned %d positions, want 7", len(pos))
	}
}

func TestGridLayoutEmpty(t *testing.T) {
	l := &layout.GridLayout{}
	pos := l.Calculate(nil, nil)

	if len(pos) != 0 {
		t.Errorf("Calculate() returned %d positions, want 0", len(pos))
	}
}

func TestLayeredLayout(t *testing.T) {
	l := &layout.LayeredLayout{}

	if l.Name() != "layered" {
		t.Errorf("Name() = %q, want 'layered'", l.Name())
	}

	components := []model.Component{
		{Name: "A"}, {Name: "B"}, {Name: "C"},
	}
	connections := []model.Connection{
		{Source: "A", Target: "B"},
		{Source: "B", Target: "C"},
	}

	pos := l.Calculate(components, connections)

	if len(pos) != 3 {
		t.Errorf("Calculate() returned %d positions, want 3", len(pos))
	}

	if pos["C"].Y <= pos["A"].Y {
		t.Error("C should be below A in layered layout")
	}
}

func TestLayeredLayoutEmpty(t *testing.T) {
	l := &layout.LayeredLayout{}
	pos := l.Calculate(nil, nil)

	if len(pos) != 0 {
		t.Errorf("Calculate() returned %d positions, want 0", len(pos))
	}
}

func TestLayeredLayoutNoConnections(t *testing.T) {
	l := &layout.LayeredLayout{}
	components := []model.Component{
		{Name: "A"}, {Name: "B"},
	}

	pos := l.Calculate(components, nil)

	if len(pos) != 2 {
		t.Errorf("Calculate() returned %d positions, want 2", len(pos))
	}
}

func TestIsometricLayout(t *testing.T) {
	l := &layout.IsometricLayout{}

	if l.Name() != "isometric" {
		t.Errorf("Name() = %q, want 'isometric'", l.Name())
	}

	components := []model.Component{
		{Name: "A"}, {Name: "B"}, {Name: "C"},
	}
	connections := []model.Connection{
		{Source: "A", Target: "B"},
		{Source: "B", Target: "C"},
	}

	pos := l.Calculate(components, connections)

	if len(pos) != 3 {
		t.Errorf("Calculate() returned %d positions, want 3", len(pos))
	}
}

func TestIsometricLayoutUnknownConnections(t *testing.T) {
	l := &layout.IsometricLayout{}
	components := []model.Component{{Name: "A"}, {Name: "B"}}
	connections := []model.Connection{{Source: "A", Target: "C"}, {Source: "X", Target: "B"}}

	pos := l.Calculate(components, connections)
	if len(pos) != 2 {
		t.Errorf("Calculate() returned %d positions, want 2", len(pos))
	}
}

func TestIsometricLayoutEmpty(t *testing.T) {
	l := &layout.IsometricLayout{}
	pos := l.Calculate(nil, nil)

	if len(pos) != 0 {
		t.Errorf("Calculate() returned %d positions, want 0", len(pos))
	}
}

func TestIsoProject(t *testing.T) {
	tests := []struct {
		x, y  float64
		wantX float64
		wantY float64
	}{
		{0, 0, 0, 0},
		{100, 0, 100, 50},
		{0, 100, -100, 50},
		{100, 100, 0, 100},
	}

	for _, tt := range tests {
		gotX, gotY := layout.IsoProject(tt.x, tt.y)
		if gotX != tt.wantX || gotY != tt.wantY {
			t.Errorf("IsoProject(%f, %f) = (%f, %f), want (%f, %f)",
				tt.x, tt.y, gotX, gotY, tt.wantX, tt.wantY)
		}
	}
}

func TestNewLayout(t *testing.T) {
	tests := []struct {
		layoutType string
		wantName   string
	}{
		{"grid", "grid"},
		{"layered", "layered"},
		{"isometric", "isometric"},
		{"unknown", "layered"},
		{"", "layered"},
	}

	for _, tt := range tests {
		l := layout.NewLayout(tt.layoutType)
		if l.Name() != tt.wantName {
			t.Errorf("NewLayout(%q).Name() = %q, want %q", tt.layoutType, l.Name(), tt.wantName)
		}
	}
}
