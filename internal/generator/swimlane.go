package generator

import (
	"fmt"
	"strings"

	"diagram-gen/internal/model"
)

// Swimlane represents a swimlane container in the diagram.
type Swimlane struct {
	ID       string
	Name     string
	X        int
	Y        int
	Width    int
	Height   int
	Children []string
}

// BuildSwimlanes creates swimlane structures from components.
func BuildSwimlanes(components []model.Component, positions map[string]Position) []Swimlane {
	swimlaneMap := make(map[string]*Swimlane)

	for _, comp := range components {
		if comp.Swimlane == "" {
			continue
		}

		if _, exists := swimlaneMap[comp.Swimlane]; !exists {
			swimlaneMap[comp.Swimlane] = &Swimlane{
				ID:       "swimlane-" + comp.Swimlane,
				Name:     comp.Swimlane,
				Children: []string{},
			}
		}

		swimlaneMap[comp.Swimlane].Children = append(
			swimlaneMap[comp.Swimlane].Children,
			comp.Name,
		)
	}

	for name, sl := range swimlaneMap {
		if len(sl.Children) == 0 {
			delete(swimlaneMap, name)
			continue
		}

		minX := positions[sl.Children[0]].X
		maxX := minX
		minY := positions[sl.Children[0]].Y
		maxY := minY

		for _, child := range sl.Children {
			pos := positions[child]
			if pos.X < minX {
				minX = pos.X
			}
			if pos.X > maxX {
				maxX = pos.X
			}
			if pos.Y < minY {
				minY = pos.Y
			}
			if pos.Y > maxY {
				maxY = pos.Y
			}
		}

		sl.X = minX - 50
		sl.Y = minY - 80
		sl.Width = maxX - minX + 220
		sl.Height = maxY - minY + 180
	}

	result := make([]Swimlane, 0, len(swimlaneMap))
	for _, sl := range swimlaneMap {
		result = append(result, *sl)
	}

	return result
}

// GenerateSwimlaneXML generates XML for swimlane elements.
func GenerateSwimlaneXML(swimlanes []Swimlane, cellID *int) string {
	var sb strings.Builder

	for _, sl := range swimlanes {
		fmt.Fprintf(&sb, `        <mxCell id="%d" value="%s" style="shape=swimlane;horizontal=1;whiteSpace=wrap;html=1;fillColor=#f5f5f5;strokeColor=#666666;" vertex="1" parent="1">
          <mxGeometry x="%d" y="%d" width="%d" height="%d" as="geometry"/>
        </mxCell>
`, *cellID, sl.Name, sl.X, sl.Y, sl.Width, sl.Height)
		*cellID++
	}

	return sb.String()
}
