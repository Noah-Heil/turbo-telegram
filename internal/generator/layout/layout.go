// Package layout provides layout algorithms for diagram positioning.
package layout

import (
	"diagram-gen/internal/model"
)

// Position represents coordinates in a diagram.
type Position struct {
	X float64
	Y float64
}

// Layout defines the interface for layout algorithms.
type Layout interface {
	Calculate(components []model.Component, connections []model.Connection) map[string]Position
	Name() string
}

// NewLayout creates a new layout instance based on the layout type.
func NewLayout(layoutType string) Layout {
	switch layoutType {
	case "layered":
		return &LayeredLayout{}
	case "isometric":
		return &IsometricLayout{}
	case "grid":
		return &GridLayout{}
	default:
		return &LayeredLayout{}
	}
}
