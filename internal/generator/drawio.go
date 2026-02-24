package generator

import (
	"fmt"
	"strings"

	"diagram-gen/internal/generator/layout"
	"diagram-gen/internal/model"
)

// DrawIOGenerator generates draw.io compatible diagrams.
type DrawIOGenerator struct {
	LayoutType string
	Compress   bool
}

// NewDrawIOGenerator creates a new DrawIOGenerator with default settings.
func NewDrawIOGenerator() *DrawIOGenerator {
	return &DrawIOGenerator{
		LayoutType: "layered",
		Compress:   false,
	}
}

// Format returns the output format name.
func (g *DrawIOGenerator) Format() string {
	return "drawio"
}

// Generate creates draw.io XML from a diagram model.
func (g *DrawIOGenerator) Generate(diagram *model.Diagram) ([]byte, error) {
	layoutType := diagram.Layout
	if layoutType == "" {
		layoutType = g.LayoutType
	}

	layoutEngine := layout.NewLayout(layoutType)
	positions := layoutEngine.Calculate(diagram.Components, diagram.Connections)

	intPositions := make(map[string]Position)
	for name, pos := range positions {
		intPositions[name] = Position{
			X: int(pos.X),
			Y: int(pos.Y),
		}
	}

	swimlanes := BuildSwimlanes(diagram.Components, intPositions)

	pages := g.BuildPages(diagram)

	var sb strings.Builder

	if g.Compress || diagram.Compress {
		sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<mxfile host="app.diagrams.net">
`)
		for _, page := range pages {
			pageXML := g.GeneratePageXML(page, swimlanes, intPositions)
			compressed, err := CompressXML([]byte(pageXML))
			if err != nil {
				compressed = []byte(pageXML)
			}
			fmt.Fprintf(&sb, `  <diagram name="%s">
    %s
  </diagram>
`, page.Name, compressed)
		}
		sb.WriteString(`</mxfile>`)
	} else {
		sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<mxfile host="app.diagrams.net">
`)
		for _, page := range pages {
			pageXML := g.GeneratePageXML(page, swimlanes, intPositions)
			sb.WriteString(pageXML)
		}
		sb.WriteString(`</mxfile>`)
	}

	return []byte(sb.String()), nil
}

func (g *DrawIOGenerator) BuildPages(diagram *model.Diagram) []model.Page {
	if len(diagram.Pages) > 0 {
		return diagram.Pages
	}

	pageMap := make(map[string]*model.Page)
	pageMap["default"] = &model.Page{Name: "Architecture Diagram"}

	for _, comp := range diagram.Components {
		pageName := comp.Page
		if pageName == "" {
			pageName = "default"
		}

		if _, exists := pageMap[pageName]; !exists {
			pageMap[pageName] = &model.Page{Name: pageName}
		}
		pageMap[pageName].Components = append(pageMap[pageName].Components, comp)
	}

	for _, conn := range diagram.Connections {
		pageName := conn.Page
		if pageName == "" {
			pageName = "default"
		}

		if _, exists := pageMap[pageName]; !exists {
			pageMap[pageName] = &model.Page{Name: pageName}
		}
		pageMap[pageName].Connections = append(pageMap[pageName].Connections, conn)
	}

	pages := make([]model.Page, 0, len(pageMap))
	for _, page := range pageMap {
		pages = append(pages, *page)
	}

	return pages
}

func (g *DrawIOGenerator) GeneratePageXML(page model.Page, swimlanes []Swimlane, positions map[string]Position) string {
	var sb strings.Builder

	components := page.Components
	connections := page.Connections

	fmt.Fprintf(&sb, `  <diagram name="%s">
    <mxGraphModel dx="1200" dy="800" grid="1" gridSize="10" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="1" pageScale="1" pageWidth="1200" pageHeight="900" math="0" shadow="0">
      <root>
        <mxCell id="0" />
        <mxCell id="1" parent="0" />
`, page.Name)

	cellID := 2

	if len(swimlanes) > 0 {
		sb.WriteString(GenerateSwimlaneXML(swimlanes, &cellID))
	}

	for _, comp := range components {
		pos := positions[comp.Name]
		if pos.X == 0 && pos.Y == 0 {
			pos = Position{X: 100 + (cellID-2)*50, Y: 100 + (cellID-2)*30}
		}

		shapeStyle := g.BuildComponentStyle(comp)

		width := 120
		height := 60
		if comp.Shape == "iso:server" || comp.Shape == "iso:database" {
			height = 80
		}

		fmt.Fprintf(&sb, `        <mxCell id="%d" value="%s" style="%s" vertex="1" parent="1">
          <mxGeometry x="%d" y="%d" width="%d" height="%d" as="geometry" />
        </mxCell>
`, cellID, EscapeXML(comp.Name), shapeStyle, pos.X, pos.Y, width, height)
		cellID++
	}

	compIDMap := make(map[string]int)
	for i := range components {
		compIDMap[components[i].Name] = i + 2
		if len(swimlanes) > 0 {
			compIDMap[components[i].Name] = i + 2 + len(swimlanes)
		}
	}

	for _, conn := range connections {
		sourceID, ok1 := compIDMap[conn.Source]
		targetID, ok2 := compIDMap[conn.Target]
		if !ok1 || !ok2 {
			continue
		}

		edgeStyle := g.BuildEdgeStyle(conn)
		fmt.Fprintf(&sb, `        <mxCell id="%d" style="%s" edge="1" parent="1" source="%d" target="%d">
          <mxGeometry as="geometry" />
        </mxCell>
`, cellID, edgeStyle, sourceID, targetID)
		cellID++
	}

	sb.WriteString(`      </root>
    </mxGraphModel>
  </diagram>
`)

	return sb.String()
}

func (g *DrawIOGenerator) BuildComponentStyle(comp model.Component) string {
	var style Style

	if comp.Shape != "" {
		shapeType := GetDefaultShapeForComponentType(string(comp.Shape))
		style.Shape = string(shapeType)
	} else {
		shapeType := GetDefaultShapeForComponentType(string(comp.Type))
		style.Shape = string(shapeType)
	}

	switch comp.Type {
	case model.ComponentTypeService, model.ComponentTypeAPI, model.ComponentTypeGateway:
		style.FillColor = "#dae8fc"
		style.StrokeColor = "#6c8ebf"
		style.FontStyle = FontStyleBold
	case model.ComponentTypeDatabase, model.ComponentTypeStorage:
		style.FillColor = "#ffe6cc"
		style.StrokeColor = "#d79b00"
	case model.ComponentTypeQueue:
		style.FillColor = "#fff2cc"
		style.StrokeColor = "#d6b656"
	case model.ComponentTypeCache:
		style.FillColor = "#f8cecc"
		style.StrokeColor = "#b85450"
		style.Dashed = true
	case model.ComponentTypeUser:
		style.FillColor = "#e1d5e7"
		style.StrokeColor = "#9673a6"
	case model.ComponentTypeExternal:
		style.FillColor = "#f5f5f5"
		style.StrokeColor = "#666666"
	default:
		style.FillColor = "#ffffff"
		style.StrokeColor = "#000000"
	}

	if comp.Shape == "iso:server" || comp.Shape == "iso:database" || comp.Shape == "iso:container" || comp.Shape == "iso:cloud" {
		style.FillColor = "#dae8fc"
		style.StrokeColor = "#6c8ebf"
	}

	if comp.Style != "" {
		override := ParseStyle(comp.Style)
		style = MergeStyles(style, override)
	}

	style.FontSize = 12
	style.WhiteSpace = WhiteSpaceWrap

	return style.String()
}

func (g *DrawIOGenerator) BuildEdgeStyle(conn model.Connection) string {
	var style Style

	if conn.EdgeStyle != "" {
		style.EdgeStyle = conn.EdgeStyle
	}

	if conn.StartArrow != "" {
		style.StartArrow = conn.StartArrow
	} else if conn.Direction == model.ConnectionDirectionBidirectional {
		style.StartArrow = ArrowClassic
	}

	if conn.EndArrow != "" {
		style.EndArrow = conn.EndArrow
	} else {
		style.EndArrow = ArrowClassic
	}

	return style.String()
}

// Position represents coordinates in the diagram.
type Position struct {
	X int
	Y int
}

func EscapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
