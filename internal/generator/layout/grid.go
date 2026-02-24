package layout

import (
	"diagram-gen/internal/model"
)

// GridLayout arranges components in a grid pattern.
type GridLayout struct{}

// Name returns the layout name.
func (l *GridLayout) Name() string {
	return "grid"
}

// Calculate computes positions for components in a grid layout.
func (l *GridLayout) Calculate(components []model.Component, _ []model.Connection) map[string]Position {
	positions := make(map[string]Position)

	var cols int
	switch n := len(components); {
	case n <= 4:
		cols = n
	case n <= 6:
		cols = 3
	default:
		cols = 4
	}

	spacingX := 180.0
	spacingY := 120.0
	startX := 100.0
	startY := 100.0

	for i, comp := range components {
		row := i / cols
		col := i % cols
		positions[comp.Name] = Position{
			X: startX + float64(col)*spacingX,
			Y: startY + float64(row)*spacingY,
		}
	}

	return positions
}
