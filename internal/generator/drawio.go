package generator

import (
	"fmt"
	"strings"

	"diagram-gen/internal/model"
)

type DrawIOGenerator struct{}

func NewDrawIOGenerator() *DrawIOGenerator {
	return &DrawIOGenerator{}
}

func (g *DrawIOGenerator) Format() string {
	return "drawio"
}

func (g *DrawIOGenerator) Generate(diagram *model.Diagram) ([]byte, error) {
	var sb strings.Builder

	components := diagram.Components
	connections := diagram.Connections

	layout := calculateLayout(components)

	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<mxfile host="app.diagrams.net">
  <diagram name="Architecture Diagram">
    <mxGraphModel dx="1200" dy="800" grid="1" gridSize="10" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="1" pageScale="1" pageWidth="1200" pageHeight="900" math="0" shadow="0">
      <root>
        <mxCell id="0" />
        <mxCell id="1" parent="0" />
`)

	cellID := 2

	for i, comp := range components {
		pos := layout[i]
		shapeStyle := getShapeStyle(comp.Type)

		sb.WriteString(fmt.Sprintf(`        <mxCell id="%d" value="%s" style="%s" vertex="1" parent="1">
          <mxGeometry x="%d" y="%d" width="120" height="60" as="geometry" />
        </mxCell>
`, cellID, escapeXML(comp.Name), shapeStyle, pos.X, pos.Y))
		cellID++
	}

	compIDMap := make(map[string]int)
	for i := range components {
		compIDMap[components[i].Name] = i + 2
	}

	for _, conn := range connections {
		sourceID, ok1 := compIDMap[conn.Source]
		targetID, ok2 := compIDMap[conn.Target]
		if !ok1 || !ok2 {
			continue
		}

		edgeStyle := getEdgeStyle(conn.Direction)
		sb.WriteString(fmt.Sprintf(`        <mxCell id="%d" style="%s" edge="1" parent="1" source="%d" target="%d">
          <mxGeometry as="geometry" />
        </mxCell>
`, cellID, edgeStyle, sourceID, targetID))
		cellID++
	}

	sb.WriteString(`      </root>
    </mxGraphModel>
  </diagram>
</mxfile>`)

	return []byte(sb.String()), nil
}

type Position struct {
	X int
	Y int
}

func calculateLayout(components []model.Component) []Position {
	layout := make([]Position, len(components))

	cols := 3
	if len(components) <= 4 {
		cols = len(components)
	} else if len(components) <= 6 {
		cols = 3
	} else {
		cols = 4
	}

	spacingX := 180
	spacingY := 120
	startX := 100
	startY := 100

	for i := range components {
		row := i / cols
		col := i % cols
		layout[i] = Position{
			X: startX + col*spacingX,
			Y: startY + row*spacingY,
		}
	}

	return layout
}

func getShapeStyle(compType model.ComponentType) string {
	switch compType {
	case model.ComponentTypeService:
		return "rounded=1;whiteSpace=wrap;html=1;fillColor=#dae8fc;strokeColor=#6c8ebf;fontSize=12;fontStyle=1"
	case model.ComponentTypeAPI:
		return "rounded=1;whiteSpace=wrap;html=1;fillColor=#d5e8d4;strokeColor=#82b366;fontSize=12;fontStyle=1"
	case model.ComponentTypeDatabase:
		return "shape=cylinder;whiteSpace=wrap;html=1;boundedLbl=1;backgroundOutline=1;size=10;fillColor=#ffe6cc;strokeColor=#d79b00;fontSize=12"
	case model.ComponentTypeQueue:
		return "shape=parallelogram;perimeter=parallelogramPerimeter;whiteSpace=wrap;html=1;fillColor=#fff2cc;strokeColor=#d6b656;fontSize=12"
	case model.ComponentTypeCache:
		return "rounded=1;whiteSpace=wrap;html=1;fillColor=#f8cecc;strokeColor=#b85450;fontSize=12;dashed=1"
	case model.ComponentTypeUser:
		return "ellipse;whiteSpace=wrap;html=1;fillColor=#e1d5e7;strokeColor=#9673a6;fontSize=12"
	case model.ComponentTypeExternal:
		return "shape=document;whiteSpace=wrap;html=1;boundedLbl=1;fillColor=#f5f5f5;strokeColor=#666666;fontSize=12"
	case model.ComponentTypeStorage:
		return "shape=cylinder;whiteSpace=wrap;html=1;boundedLbl=1;backgroundOutline=1;size=10;fillColor=#fff2cc;strokeColor=#d6b656;fontSize=12"
	case model.ComponentTypeGateway:
		return "rounded=1;whiteSpace=wrap;html=1;fillColor=#d5e8d4;strokeColor=#82b366;fontSize=12;fontStyle=1"
	default:
		return "rounded=1;whiteSpace=wrap;html=1;fillColor=#ffffff;strokeColor=#000000;fontSize=12"
	}
}

func getEdgeStyle(direction model.ConnectionDirection) string {
	if direction == model.ConnectionDirectionBidirectional {
		return "endArrow=classic;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0"
	}
	return "endArrow=classic;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0"
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
