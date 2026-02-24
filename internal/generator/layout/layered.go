package layout

import (
	"diagram-gen/internal/model"
)

// LayeredLayout arranges components in layers based on connections.
type LayeredLayout struct{}

// Name returns the layout name.
func (l *LayeredLayout) Name() string {
	return "layered"
}

// Calculate computes positions for components in a layered layout.
func (l *LayeredLayout) Calculate(components []model.Component, connections []model.Connection) map[string]Position {
	positions := make(map[string]Position)
	layers := make(map[string]int)

	for _, comp := range components {
		layers[comp.Name] = 0
	}

	maxIterations := len(components) * len(components)
	for iteration := 0; iteration < maxIterations; iteration++ {
		changed := false
		for _, conn := range connections {
			srcLayer, srcOk := layers[conn.Source]
			tgtLayer, tgtOk := layers[conn.Target]

			if srcOk && tgtOk && tgtLayer <= srcLayer {
				layers[conn.Target] = srcLayer + 1
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	layerWidths := make(map[int]int)
	for _, layer := range layers {
		layerWidths[layer]++
	}

	spacingX := 200.0
	spacingY := 120.0
	startX := 100.0
	startY := 100.0

	layerPositions := make(map[int]float64)
	for layer := range layerWidths {
		layerPositions[layer] = startX
	}

	for i, comp := range components {
		layer := layers[comp.Name]
		layerPositions[layer] += spacingX
		positions[comp.Name] = Position{
			X: layerPositions[layer],
			Y: startY + float64(layer)*spacingY,
		}
		_ = i
	}

	return positions
}
