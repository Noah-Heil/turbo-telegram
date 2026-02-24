package layout

import (
	"diagram-gen/internal/model"
	"math"
)

// IsometricLayout arranges components in an isometric pattern.
type IsometricLayout struct{}

// Name returns the layout name.
func (l *IsometricLayout) Name() string {
	return "isometric"
}

// IsoProject converts 2D coordinates to isometric projection.
func IsoProject(x, y float64) (float64, float64) {
	return x - y, (x + y) / 2
}

// Calculate computes positions for components in an isometric layout.
func (l *IsometricLayout) Calculate(components []model.Component, connections []model.Connection) map[string]Position {
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

	layerGroups := make(map[int][]string)
	for _, comp := range components {
		layer := layers[comp.Name]
		layerGroups[layer] = append(layerGroups[layer], comp.Name)
	}

	baseX := 300.0
	baseY := 100.0
	spacingX := 200.0
	spacingY := 150.0

	for layer, comps := range layerGroups {
		for i, name := range comps {
			isoX, isoY := IsoProject(
				baseX+float64(i)*spacingX,
				baseY+float64(layer)*spacingY,
			)
			positions[name] = Position{
				X: isoX,
				Y: isoY - float64(layer)*spacingY*0.5,
			}
		}
	}

	if len(positions) == 0 {
		for i, comp := range components {
			isoX, isoY := IsoProject(
				baseX+float64(i%3)*spacingX,
				baseY+float64(i/3)*spacingY,
			)
			positions[comp.Name] = Position{
				X: isoX,
				Y: isoY,
			}
		}
	}

	_ = math.Max
	return positions
}
