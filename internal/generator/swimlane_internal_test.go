package generator

import (
	"testing"

	"diagram-gen/internal/model"
)

func TestPruneEmptySwimlanes(t *testing.T) {
	t.Parallel()
	swimlaneMap := map[string]*Swimlane{
		"empty":  {Name: "empty", Children: []string{}},
		"filled": {Name: "filled", Children: []string{"A"}},
	}

	pruneEmptySwimlanes(swimlaneMap)

	if _, ok := swimlaneMap["empty"]; ok {
		t.Fatal("expected empty swimlane to be pruned")
	}
	if _, ok := swimlaneMap["filled"]; !ok {
		t.Fatal("expected filled swimlane to remain")
	}
}

func TestBuildSwimlanesBounds(t *testing.T) {
	t.Parallel()
	components := []model.Component{
		{Name: "A", Swimlane: "Lane"},
		{Name: "B", Swimlane: "Lane"},
		{Name: "C", Swimlane: "Lane"},
		{Name: "D", Swimlane: ""},
	}
	positions := map[string]Position{
		"A": {X: 200, Y: 200},
		"B": {X: 100, Y: 300},
		"C": {X: 300, Y: 100},
		"D": {X: 400, Y: 400},
	}

	swimlanes := BuildSwimlanes(components, positions)
	if len(swimlanes) != 1 {
		t.Fatalf("expected 1 swimlane, got %d", len(swimlanes))
	}
	if swimlanes[0].Width == 0 || swimlanes[0].Height == 0 {
		t.Fatal("expected swimlane bounds to be set")
	}
}
