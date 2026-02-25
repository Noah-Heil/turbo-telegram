package layout

import (
	"testing"

	"diagram-gen/internal/model"
)

func TestFallbackIsometricPositions(t *testing.T) {
	components := []model.Component{{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}}
	positions := fallbackIsometricPositions(components, 300, 100, 200, 150)

	if len(positions) != len(components) {
		t.Fatalf("expected %d positions, got %d", len(components), len(positions))
	}

	if positions["A"].X == positions["B"].X {
		t.Error("expected different positions for A and B")
	}
}
