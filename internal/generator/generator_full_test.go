package generator_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"diagram-gen/internal/generator"
	"diagram-gen/internal/generator/layout"
	"diagram-gen/internal/model"
)

func TestStyleString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		style    generator.Style
		expected string
	}{
		{
			name: "basic shape",
			style: generator.Style{
				Shape:     "rectangle",
				FillColor: "#dae8fc",
			},
			expected: "shape=rectangle;fillColor=#dae8fc;html=1",
		},
		{
			name: "with font",
			style: generator.Style{
				FontSize:  14,
				FontColor: "#000000",
			},
			expected: "fontSize=14;fontColor=#000000;html=1",
		},
		{
			name: "with gradient",
			style: generator.Style{
				FillColor:     "#dae8fc",
				GradientColor: "#ffffff",
				GradientDir:   "north",
			},
			expected: "fillColor=#dae8fc;gradientColor=#ffffff;gradientDirection=north;html=1",
		},
		{
			name: "rounded and shadow",
			style: generator.Style{
				Rounded: true,
				Shadow:  true,
			},
			expected: "rounded=1;shadow=1;html=1",
		},
		{
			name: "dashed",
			style: generator.Style{
				Dashed:      true,
				DashPattern: "5 5",
			},
			expected: "dashed=1;dashPattern=5 5;html=1",
		},
		{
			name: "image",
			style: generator.Style{
				Image:       "data:image/svg+xml;base64,abc",
				ImageWidth:  100,
				ImageHeight: 100,
				ImageAspect: true,
			},
			expected: "image=data:image/svg+xml;base64,abc;imageWidth=100;imageHeight=100;imageAspect=1;html=1",
		},
		{
			name: "edge style",
			style: generator.Style{
				EdgeStyle:  "elbowEdgeStyle",
				StartArrow: "block",
				EndArrow:   "classic",
			},
			expected: "edgeStyle=elbowEdgeStyle;startArrow=block;endArrow=classic;html=1",
		},
		{
			name:     "empty",
			style:    generator.Style{},
			expected: "html=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.style.String()
			if got != tt.expected {
				t.Errorf("generator.Style.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestParseStyle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		styleStr      string
		wantFillColor string
		wantFontSize  int
		wantRounded   bool
	}{
		{
			name:          "basic fill color",
			styleStr:      "fillColor=#dae8fc",
			wantFillColor: "#dae8fc",
		},
		{
			name:         "font size",
			styleStr:     "fontSize=14",
			wantFontSize: 14,
		},
		{
			name:        "rounded",
			styleStr:    "rounded=1",
			wantRounded: true,
		},
		{
			name:        "not rounded",
			styleStr:    "rounded=0",
			wantRounded: false,
		},
		{
			name:          "multiple values",
			styleStr:      "fillColor=#dae8fc;strokeColor=#6c8ebf;fontSize=14",
			wantFillColor: "#dae8fc",
			wantFontSize:  14,
		},
		{
			name:     "empty",
			styleStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := generator.ParseStyle(tt.styleStr)
			if got.FillColor != tt.wantFillColor {
				t.Errorf("FillColor = %q, want %q", got.FillColor, tt.wantFillColor)
			}
			if got.FontSize != tt.wantFontSize {
				t.Errorf("FontSize = %d, want %d", got.FontSize, tt.wantFontSize)
			}
			if got.Rounded != tt.wantRounded {
				t.Errorf("Rounded = %v, want %v", got.Rounded, tt.wantRounded)
			}
		})
	}
}

func TestMergeStyles(t *testing.T) {
	t.Parallel()
	base := generator.Style{
		FillColor:   "#ffffff",
		StrokeColor: "#000000",
		FontSize:    12,
	}

	override := generator.Style{
		FillColor: "#dae8fc",
		FontSize:  14,
		Shadow:    true,
	}

	got := generator.MergeStyles(base, override)

	if got.FillColor != "#dae8fc" {
		t.Errorf("FillColor = %q, want #dae8fc", got.FillColor)
	}
	if got.StrokeColor != "#000000" {
		t.Errorf("StrokeColor = %q, want #000000", got.StrokeColor)
	}
	if got.FontSize != 14 {
		t.Errorf("FontSize = %d, want 14", got.FontSize)
	}
	if !got.Shadow {
		t.Error("Shadow should be true")
	}
}

func TestShapeTypeIsIsometric(t *testing.T) {
	t.Parallel()
	tests := []struct {
		shape    generator.ShapeType
		expected bool
	}{
		{generator.ShapeIsoCube, true},
		{generator.ShapeIsoServer, true},
		{generator.ShapeIsoDatabase, true},
		{generator.ShapeIsoContainer, true},
		{generator.ShapeIsoCloud, true},
		{generator.ShapeIsoNetwork, true},
		{generator.ShapeIsoCylinder, true},
		{generator.ShapeRectangle, false},
		{generator.ShapeEllipse, false},
		{generator.ShapeCylinder, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.shape), func(t *testing.T) {
			t.Parallel()
			got := tt.shape.IsIsometric()
			if got != tt.expected {
				t.Errorf("IsIsometric() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestShapeTypeIsBasic(t *testing.T) {
	t.Parallel()
	tests := []struct {
		shape    generator.ShapeType
		expected bool
	}{
		{generator.ShapeRectangle, true},
		{generator.ShapeEllipse, true},
		{generator.ShapeRounded, true},
		{generator.ShapeRhombus, true},
		{generator.ShapeParallelogram, true},
		{generator.ShapeCylinder, true},
		{generator.ShapeDocument, true},
		{generator.ShapeSwimlane, true},
		{generator.ShapeTriangle, true},
		{generator.ShapeHexagon, true},
		{generator.ShapeCloud, true},
		{generator.ShapeIsoCube, false},
		{generator.ShapeIsoServer, false},
		{generator.ShapeImage, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.shape), func(t *testing.T) {
			t.Parallel()
			got := tt.shape.IsBasic()
			if got != tt.expected {
				t.Errorf("IsBasic() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSwimlaneBuild(t *testing.T) {
	t.Parallel()
	components := []model.Component{
		{Name: "Service1", Swimlane: "AWS"},
		{Name: "Service2", Swimlane: "AWS"},
		{Name: "Service3", Swimlane: "Azure"},
		{Name: "Service4", Swimlane: ""},
	}

	positions := map[string]generator.Position{
		"Service1": {X: 100, Y: 100},
		"Service2": {X: 200, Y: 100},
		"Service3": {X: 100, Y: 200},
		"Service4": {X: 300, Y: 300},
	}

	swimlanes := generator.BuildSwimlanes(components, positions)

	if len(swimlanes) != 2 {
		t.Errorf("expected 2 swimlanes, got %d", len(swimlanes))
	}

	for _, sl := range swimlanes {
		if sl.Name == "AWS" {
			if len(sl.Children) != 2 {
				t.Errorf("AWS swimlane should have 2 children, got %d", len(sl.Children))
			}
		}
	}
}

func TestSwimlaneGenerateXML(t *testing.T) {
	t.Parallel()
	swimlanes := []generator.Swimlane{
		{Name: "AWS", X: 50, Y: 50, Width: 300, Height: 200},
	}
	cellID := 100

	got := generator.GenerateSwimlaneXML(swimlanes, &cellID)

	if cellID != 101 {
		t.Errorf("cellID should be 101, got %d", cellID)
	}

	expected := `mxCell id="100"`
	if !contains(got, expected) {
		t.Errorf("expected XML to contain %q, got: %s", expected, got)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestCompressXML(t *testing.T) {
	t.Parallel()
	xml := []byte(`<mxfile><diagram>test</diagram></mxfile>`)

	compressed, err := generator.CompressXML(xml)
	if err != nil {
		t.Fatalf("CompressXML failed: %v", err)
	}

	if len(compressed) == 0 {
		t.Error("compressed data should not be empty")
	}
}

func TestCompressAndEncode(t *testing.T) {
	t.Parallel()
	xml := []byte(`<mxfile><diagram>test</diagram></mxfile>`)

	encoded, err := generator.CompressAndEncode(xml)
	if err != nil {
		t.Fatalf("CompressAndEncode failed: %v", err)
	}

	if encoded == "" {
		t.Error("encoded string should not be empty")
	}
}

func TestGetShapeStyle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		shape    generator.ShapeType
		expected string
	}{
		{generator.ShapeRectangle, "shape=rectangle;"},
		{generator.ShapeEllipse, "shape=ellipse;"},
		{generator.ShapeRounded, "shape=rounded;"},
		{generator.ShapeRhombus, "shape=rhombus;"},
		{generator.ShapeParallelogram, "shape=parallelogram;"},
		{generator.ShapeCylinder, "shape=cylinder;"},
		{generator.ShapeDocument, "shape=document;"},
		{generator.ShapeSwimlane, "shape=swimlane;"},
		{generator.ShapeTriangle, "shape=triangle;"},
		{generator.ShapeHexagon, "shape=hexagon;"},
		{generator.ShapeCloud, "shape=cloud;"},
		{generator.ShapeInternal, "shape=internal;"},
		{generator.ShapeExternal, "shape=external;"},
		{generator.ShapeFolder, "shape=folder;"},
		{generator.ShapeIsoCube, "shape=mxgraph.isometric.cube;"},
		{generator.ShapeIsoServer, "shape=mxgraph.isometric.server;"},
		{generator.ShapeIsoDatabase, "shape=mxgraph.isometric.database;"},
		{generator.ShapeIsoContainer, "shape=mxgraph.isometric.container;"},
		{generator.ShapeIsoCloud, "shape=mxgraph.isometric.cloud;"},
		{generator.ShapeIsoNetwork, "shape=mxgraph.isometric.network;"},
		{generator.ShapeIsoCylinder, "shape=mxgraph.isometric.cylinder;"},
		{generator.ShapeImage, "shape=rectangle;"},
		{"unknown", "shape=rectangle;"},
	}

	for _, tt := range tests {
		t.Run(string(tt.shape), func(t *testing.T) {
			t.Parallel()
			got := generator.GetShapeStyle(tt.shape)
			if len(got) < len(tt.expected) || got[:len(tt.expected)] != tt.expected {
				t.Errorf("GetShapeStyle(%q) = %q, want prefix %q", tt.shape, got, tt.expected)
			}
		})
	}
}

func TestGetDefaultShapeForComponentType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		compType string
		expected generator.ShapeType
	}{
		{"service", generator.ShapeRounded},
		{"api", generator.ShapeRounded},
		{"gateway", generator.ShapeRounded},
		{"database", generator.ShapeCylinder},
		{"storage", generator.ShapeCylinder},
		{"queue", generator.ShapeParallelogram},
		{"cache", generator.ShapeRounded},
		{"user", generator.ShapeEllipse},
		{"external", generator.ShapeDocument},
		{"iso:server", generator.ShapeIsoServer},
		{"iso:database", generator.ShapeIsoDatabase},
		{"iso:container", generator.ShapeIsoContainer},
		{"iso:cloud", generator.ShapeIsoCloud},
		{"iso:network", generator.ShapeIsoNetwork},
		{"iso:cube", generator.ShapeIsoCube},
		{"iso:cylinder", generator.ShapeIsoCylinder},
		{"unknown_type", generator.ShapeRectangle},
	}

	for _, tt := range tests {
		t.Run(tt.compType, func(t *testing.T) {
			t.Parallel()
			got := generator.GetDefaultShapeForComponentType(tt.compType)
			if got != tt.expected {
				t.Errorf("GetDefaultShapeForComponentType(%q) = %q, want %q", tt.compType, got, tt.expected)
			}
		})
	}
}

func TestParseStyleOpacity(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("opacity=50")
	if style.Opacity != 50 {
		t.Errorf("Opacity = %d, want 50", style.Opacity)
	}
}

func TestParseStyleStrokeWidth(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("strokeWidth=3")
	if style.StrokeWidth != 3 {
		t.Errorf("StrokeWidth = %d, want 3", style.StrokeWidth)
	}
}

func TestParseStyleFontFamily(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("fontFamily=Arial")
	if style.FontFamily != "Arial" {
		t.Errorf("FontFamily = %q, want Arial", style.FontFamily)
	}
}

func TestParseStyleFontStyle(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("fontStyle=1")
	if style.FontStyle != 1 {
		t.Errorf("FontStyle = %d, want 1", style.FontStyle)
	}
}

func TestParseStyleImage(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("image=https://example.com/image.png")
	if style.Image != "https://example.com/image.png" {
		t.Errorf("Image = %q, want https://example.com/image.png", style.Image)
	}

	style2 := generator.ParseStyle("imageWidth=50;imageHeight=60;imageAspect=1")
	if style2.ImageWidth != 50 {
		t.Errorf("ImageWidth = %d, want 50", style2.ImageWidth)
	}
	if style2.ImageHeight != 60 {
		t.Errorf("ImageHeight = %d, want 60", style2.ImageHeight)
	}
	if !style2.ImageAspect {
		t.Error("ImageAspect should be true")
	}
}

func TestParseStyleEdgeAndArrow(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("edgeStyle=orthogonalEdgeStyle;startArrow=block;endArrow=diamond")
	if style.EdgeStyle != "orthogonalEdgeStyle" {
		t.Errorf("EdgeStyle = %q, want orthogonalEdgeStyle", style.EdgeStyle)
	}
	if style.StartArrow != "block" {
		t.Errorf("StartArrow = %q, want block", style.StartArrow)
	}
	if style.EndArrow != "diamond" {
		t.Errorf("EndArrow = %q, want diamond", style.EndArrow)
	}
}

func TestParseStyleCurvedAndElbow(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("curved=1;elbow=horizontal")
	if !style.Curved {
		t.Error("Curved should be true")
	}
	if style.Elbow != "horizontal" {
		t.Errorf("Elbow = %q, want horizontal", style.Elbow)
	}
}

func TestParseStyleWhiteSpace(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("whiteSpace=wrap")
	if style.WhiteSpace != "wrap" {
		t.Errorf("WhiteSpace = %q, want wrap", style.WhiteSpace)
	}
}

func TestParseStyleAlign(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("align=center;verticalAlign=bottom")
	if style.Align != "center" {
		t.Errorf("Align = %q, want center", style.Align)
	}
	if style.VerticalAlign != "bottom" {
		t.Errorf("VerticalAlign = %q, want bottom", style.VerticalAlign)
	}
}

func TestMergeStylesOverrideAll(t *testing.T) {
	t.Parallel()
	base := generator.Style{
		FillColor:   "#ffffff",
		StrokeColor: "#000000",
		FontSize:    12,
		Rounded:     false,
	}

	override := generator.Style{
		FillColor:   "#dae8fc",
		StrokeColor: "#6c8ebf",
		FontSize:    14,
		Rounded:     true,
	}

	got := generator.MergeStyles(base, override)

	if got.FillColor != "#dae8fc" {
		t.Errorf("FillColor = %q, want #dae8fc", got.FillColor)
	}
	if got.StrokeColor != "#6c8ebf" {
		t.Errorf("StrokeColor = %q, want #6c8ebf", got.StrokeColor)
	}
	if got.FontSize != 14 {
		t.Errorf("FontSize = %d, want 14", got.FontSize)
	}
	if !got.Rounded {
		t.Error("Rounded should be true")
	}
}

func TestSwimlaneEmptyComponents(t *testing.T) {
	t.Parallel()
	positions := map[string]generator.Position{
		"Service1": {X: 100, Y: 100},
	}

	swimlanes := generator.BuildSwimlanes(nil, positions)
	if len(swimlanes) != 0 {
		t.Errorf("expected 0 swimlanes for nil components, got %d", len(swimlanes))
	}
}

func TestSwimlaneNoSwimlane(t *testing.T) {
	t.Parallel()
	components := []model.Component{
		{Name: "Service1", Swimlane: ""},
		{Name: "Service2", Swimlane: ""},
	}

	positions := map[string]generator.Position{
		"Service1": {X: 100, Y: 100},
		"Service2": {X: 200, Y: 100},
	}

	swimlanes := generator.BuildSwimlanes(components, positions)
	if len(swimlanes) != 0 {
		t.Errorf("expected 0 swimlanes when none have swimlane set, got %d", len(swimlanes))
	}
}

func TestDrawIOGenerator(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	if gen.Format() != "drawio" {
		t.Errorf("Format() = %q, want drawio", gen.Format())
	}
}

func TestDrawIOGeneratorWithOptions(t *testing.T) {
	t.Parallel()
	gen := &generator.DrawIOGenerator{
		LayoutType: "isometric",
		Compress:   true,
	}

	if gen.LayoutType != "isometric" {
		t.Errorf("LayoutType = %q, want isometric", gen.LayoutType)
	}
	if !gen.Compress {
		t.Error("Compress should be true")
	}
}

func TestDrawIOGenerate(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "Service1", Type: model.ComponentTypeService},
		},
		Connections: []model.Connection{
			{Source: "Service1", Target: "Database"},
		},
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestDrawIOGenerateWithLayout(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "A"}, {Name: "B"}, {Name: "C"},
		},
		Connections: []model.Connection{
			{Source: "A", Target: "B"},
			{Source: "B", Target: "C"},
		},
		Layout: "layered",
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestDrawIOGenerateWithCompress(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "Service1"},
		},
		Compress: true,
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestDrawIOBuildPages(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "Service1", Page: "page1"},
			{Name: "Service2", Page: "page2"},
		},
		Connections: []model.Connection{
			{Source: "Service1", Target: "Service2", Page: "page1"},
		},
	}

	pages := gen.BuildPages(diagram)
	if len(pages) < 2 {
		t.Errorf("expected at least 2 pages, got %d", len(pages))
	}
}

func TestDrawIOBuildPagesDefault(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "Service1"},
			{Name: "Service2"},
		},
		Connections: []model.Connection{
			{Source: "Service1", Target: "Service2"},
		},
	}

	pages := gen.BuildPages(diagram)
	if len(pages) != 1 {
		t.Errorf("expected 1 page, got %d", len(pages))
	}
}

func TestBuildComponentStyleWithShape(t *testing.T) {
	t.Parallel()
	gen := &generator.DrawIOGenerator{}
	comp := model.Component{
		Name:  "Test",
		Type:  model.ComponentTypeService,
		Shape: "iso:server",
		Style: "fillColor=#0000ff",
	}

	style := gen.BuildComponentStyle(comp)
	if style == "" {
		t.Error("style should not be empty")
	}
}

func TestBuildEdgeStyleWithOptions(t *testing.T) {
	t.Parallel()
	gen := &generator.DrawIOGenerator{}
	conn := model.Connection{
		Source:     "A",
		Target:     "B",
		Direction:  model.ConnectionDirectionBidirectional,
		EdgeStyle:  "elbowEdgeStyle",
		StartArrow: "block",
		EndArrow:   "diamond",
	}

	style := gen.BuildEdgeStyle(conn)
	if style == "" {
		t.Error("style should not be empty")
	}
}

func TestBuildSwimlanesMultiple(t *testing.T) {
	t.Parallel()
	components := []model.Component{
		{Name: "S1", Swimlane: "AWS"},
		{Name: "S2", Swimlane: "AWS"},
		{Name: "S3", Swimlane: "AWS"},
		{Name: "DB1", Swimlane: "DB"},
	}

	positions := map[string]generator.Position{
		"S1":  {X: 100, Y: 100},
		"S2":  {X: 200, Y: 100},
		"S3":  {X: 300, Y: 100},
		"DB1": {X: 400, Y: 100},
	}

	swimlanes := generator.BuildSwimlanes(components, positions)
	if len(swimlanes) != 2 {
		t.Errorf("expected 2 swimlanes, got %d", len(swimlanes))
	}
}

func TestMergeStylesEdgeCases(t *testing.T) {
	t.Parallel()
	base := generator.Style{
		WhiteSpace:    "wrap",
		Align:         "center",
		VerticalAlign: "middle",
	}

	override := generator.Style{
		WhiteSpace:    "nowrap",
		Align:         "left",
		GradientColor: "#ffffff",
	}

	got := generator.MergeStyles(base, override)

	if got.WhiteSpace != "nowrap" {
		t.Errorf("WhiteSpace = %q, want nowrap", got.WhiteSpace)
	}
	if got.Align != "left" {
		t.Errorf("Align = %q, want left", got.Align)
	}
	if got.VerticalAlign != "middle" {
		t.Errorf("VerticalAlign = %q, want middle", got.VerticalAlign)
	}
	if got.GradientColor != "#ffffff" {
		t.Errorf("GradientColor = %q, want #ffffff", got.GradientColor)
	}
}

func TestGridLayoutWithManyComponents(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("grid")
	components := make([]model.Component, 10)
	for i := range components {
		components[i] = model.Component{Name: string(rune('A' + i))}
	}

	pos := l.Calculate(components, nil)
	if len(pos) != 10 {
		t.Errorf("expected 10 positions, got %d", len(pos))
	}
}

func TestLayeredLayoutWithNoConnections(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("layered")
	components := []model.Component{
		{Name: "A"}, {Name: "B"}, {Name: "C"},
	}

	pos := l.Calculate(components, nil)
	if len(pos) != 3 {
		t.Errorf("expected 3 positions, got %d", len(pos))
	}
}

func TestIsometricLayoutWithNoConnections(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("isometric")
	components := []model.Component{
		{Name: "A"}, {Name: "B"}, {Name: "C"},
	}

	pos := l.Calculate(components, nil)
	if len(pos) != 3 {
		t.Errorf("expected 3 positions, got %d", len(pos))
	}
}

func TestCompressXMLWriter(t *testing.T) {
	t.Parallel()
	xml := []byte(`<mxfile><diagram>test</diagram></mxfile>`)

	var buf bytes.Buffer
	err := generator.CompressXMLWriter(xml, &buf)
	if err != nil {
		t.Fatalf("CompressXMLWriter failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("compressed data should not be empty")
	}
}

func TestStyleStringFull(t *testing.T) {
	t.Parallel()
	style := generator.Style{
		Shape:         "rectangle",
		FillColor:     "#dae8fc",
		StrokeColor:   "#6c8ebf",
		StrokeWidth:   2,
		Opacity:       80,
		GradientColor: "#ffffff",
		GradientDir:   "north",
		FontSize:      14,
		FontFamily:    "Arial",
		FontColor:     "#000000",
		FontStyle:     1,
		Rounded:       true,
		Dashed:        true,
		DashPattern:   "5 5",
		Shadow:        true,
		Glass:         true,
		WhiteSpace:    "wrap",
		Align:         "center",
		VerticalAlign: "middle",
	}

	got := style.String()
	if !strings.Contains(got, "shape=rectangle") {
		t.Errorf("missing shape: %s", got)
	}
	if !strings.Contains(got, "fillColor=#dae8fc") {
		t.Errorf("missing fillColor: %s", got)
	}
}

func TestBuildSwimlanesEdgeCases(t *testing.T) {
	t.Parallel()
	components := []model.Component{
		{Name: "S1", Swimlane: "AWS"},
		{Name: "S2", Swimlane: "AWS"},
	}

	positions := map[string]generator.Position{
		"S1": {X: 150, Y: 150},
		"S2": {X: 250, Y: 250},
	}

	swimlanes := generator.BuildSwimlanes(components, positions)
	if len(swimlanes) != 1 {
		t.Fatalf("expected 1 swimlane, got %d", len(swimlanes))
	}

	sl := swimlanes[0]
	if sl.X == 0 {
		t.Error("swimlane X should be set")
	}
	if sl.Width == 0 {
		t.Error("swimlane Width should be set")
	}
}

func TestDrawIOGenerateFull(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "A", Type: model.ComponentTypeService},
			{Name: "B", Type: model.ComponentTypeDatabase},
			{Name: "C", Type: model.ComponentTypeAPI},
		},
		Connections: []model.Connection{
			{Source: "A", Target: "B"},
			{Source: "B", Target: "C", EdgeStyle: "elbowEdgeStyle"},
		},
		Layout: "layered",
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestDrawIOGenerateWithSwimlanes(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "A", Swimlane: "AWS"},
			{Name: "B", Swimlane: "AWS"},
		},
		Connections: []model.Connection{
			{Source: "A", Target: "B"},
		},
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestDrawIOBuildPagesMultiple(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "A", Page: "Page1"},
			{Name: "B", Page: "Page1"},
			{Name: "C", Page: "Page2"},
			{Name: "D", Page: "Page2"},
			{Name: "E"},
		},
		Connections: []model.Connection{
			{Source: "A", Target: "B", Page: "Page1"},
			{Source: "C", Target: "D", Page: "Page2"},
			{Source: "E", Target: "A"},
		},
	}

	pages := gen.BuildPages(diagram)
	if len(pages) < 2 {
		t.Errorf("expected at least 2 pages, got %d", len(pages))
	}
}

func TestDrawIOGeneratePageXML(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	page := model.Page{
		Name: "Test",
		Components: []model.Component{
			{Name: "S1", Type: model.ComponentTypeService},
		},
		Connections: []model.Connection{
			{Source: "S1", Target: "S2"},
		},
	}

	data := gen.GeneratePageXML(page, nil, map[string]generator.Position{"S1": {X: 100, Y: 100}})
	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestDrawIOGeneratePageXMLWithSwimlanes(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	swimlanes := []generator.Swimlane{{Name: "AWS", X: 50, Y: 50, Width: 300, Height: 200}}
	page := model.Page{
		Name: "Test",
		Components: []model.Component{
			{Name: "S1", Type: model.ComponentTypeService, Swimlane: "AWS"},
		},
	}

	data := gen.GeneratePageXML(page, swimlanes, map[string]generator.Position{"S1": {X: 100, Y: 100}})
	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestGridLayoutVariousSizes(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("grid")

	tests := []int{1, 2, 3, 4, 5, 7, 10, 15}
	for _, n := range tests {
		components := make([]model.Component, n)
		for i := range components {
			components[i] = model.Component{Name: string(rune('A' + i))}
		}
		pos := l.Calculate(components, nil)
		if len(pos) != n {
			t.Errorf("for n=%d: expected %d positions, got %d", n, n, len(pos))
		}
	}
}

func TestCompressXMLError(t *testing.T) {
	t.Parallel()
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	compressed, err := generator.CompressXML(largeData)
	if err != nil {
		t.Fatalf("CompressXML failed: %v", err)
	}

	if len(compressed) == 0 {
		t.Error("compressed data should not be empty")
	}
}

type errorWriter struct{}

func (e *errorWriter) Write(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("write error")
}

type errorWriteCloser struct {
	errorWriter
}

func (e *errorWriteCloser) Close() error {
	return fmt.Errorf("close error")
}

func TestCompressXMLWriteError(t *testing.T) {
	t.Parallel()
	xml := []byte(`<mxfile><diagram>test</diagram></mxfile>`)
	writer := &errorWriteCloser{}
	err := generator.CompressXMLWriter(xml, writer)
	if err == nil {
		t.Error("expected error from write, got nil")
	}
}

func TestCompressAndEncodeLargeData(t *testing.T) {
	t.Parallel()
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	encoded, err := generator.CompressAndEncode(largeData)
	if err != nil {
		t.Fatalf("CompressAndEncode failed: %v", err)
	}

	if encoded == "" {
		t.Error("encoded string should not be empty")
	}
}

func TestParseIntInvalid(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("opacity=invalid")
	if style.Opacity != 0 {
		t.Errorf("Opacity = %d, want 0 for invalid value", style.Opacity)
	}

	style2 := generator.ParseStyle("strokeWidth=abc")
	if style2.StrokeWidth != 0 {
		t.Errorf("StrokeWidth = %d, want 0 for invalid value", style2.StrokeWidth)
	}

	style3 := generator.ParseStyle("fontSize=")
	if style3.FontSize != 0 {
		t.Errorf("FontSize = %d, want 0 for empty value", style3.FontSize)
	}
}

func TestStyleStringOpacity100(t *testing.T) {
	t.Parallel()
	style := generator.Style{
		Opacity: 100,
		Shape:   "rectangle",
	}
	got := style.String()
	if strings.Contains(got, "opacity=100") {
		t.Errorf("opacity=100 should not be included: %s", got)
	}
}

func TestParseStyleNoEquals(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("invalidpart")
	if style.Shape != "" {
		t.Errorf("Shape = %q, want empty", style.Shape)
	}

	style2 := generator.ParseStyle("fillColor=val1;strokeColor=val2")
	if style2.FillColor != "val1" {
		t.Errorf("FillColor = %q, want val1", style2.FillColor)
	}
	if style2.StrokeColor != "val2" {
		t.Errorf("StrokeColor = %q, want val2", style2.StrokeColor)
	}
}

func TestMergeStylesEmpty(t *testing.T) {
	t.Parallel()
	base := generator.Style{
		FillColor: "#ffffff",
	}
	override := generator.Style{}

	got := generator.MergeStyles(base, override)
	if got.FillColor != "#ffffff" {
		t.Errorf("FillColor = %q, want #ffffff", got.FillColor)
	}
}

func TestMergeStylesFont(t *testing.T) {
	t.Parallel()
	base := generator.Style{
		FontFamily: "Arial",
		FontColor:  "#000000",
	}
	override := generator.Style{
		FontFamily: "Helvetica",
		FontSize:   14,
	}

	got := generator.MergeStyles(base, override)
	if got.FontFamily != "Helvetica" {
		t.Errorf("FontFamily = %q, want Helvetica", got.FontFamily)
	}
	if got.FontSize != 14 {
		t.Errorf("FontSize = %d, want 14", got.FontSize)
	}
	if got.FontColor != "#000000" {
		t.Errorf("FontColor = %q, want #000000", got.FontColor)
	}
}

func TestMergeStylesEdgeArrow(t *testing.T) {
	t.Parallel()
	base := generator.Style{
		Shadow:        true,
		Glass:         false,
		WhiteSpace:    "wrap",
		Align:         "center",
		VerticalAlign: "top",
	}
	override := generator.Style{
		Shadow:        true,
		Glass:         true,
		WhiteSpace:    "nowrap",
		Align:         "left",
		VerticalAlign: "bottom",
	}

	got := generator.MergeStyles(base, override)
	if !got.Shadow {
		t.Error("Shadow should be true")
	}
	if !got.Glass {
		t.Error("Glass should be true")
	}
	if got.WhiteSpace != "nowrap" {
		t.Errorf("WhiteSpace = %q, want nowrap", got.WhiteSpace)
	}
	if got.Align != "left" {
		t.Errorf("Align = %q, want left", got.Align)
	}
	if got.VerticalAlign != "bottom" {
		t.Errorf("VerticalAlign = %q, want bottom", got.VerticalAlign)
	}
}

func TestSwimlaneBuildMissingPositions(t *testing.T) {
	t.Parallel()
	components := []model.Component{
		{Name: "S1", Swimlane: "AWS"},
		{Name: "S2", Swimlane: "AWS"},
	}

	positions := map[string]generator.Position{}

	swimlanes := generator.BuildSwimlanes(components, positions)
	if len(swimlanes) != 1 {
		t.Fatalf("expected 1 swimlane, got %d", len(swimlanes))
	}
}

func TestDrawIOBuildPagesNoComponentPage(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "Service1"},
		},
		Connections: []model.Connection{
			{Source: "Service1", Target: "Service2"},
		},
	}

	pages := gen.BuildPages(diagram)
	if len(pages) != 1 {
		t.Errorf("expected 1 page, got %d", len(pages))
	}
}

func TestDrawIOGeneratePageXMLUnknownConnection(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	page := model.Page{
		Name: "Test",
		Components: []model.Component{
			{Name: "S1", Type: model.ComponentTypeService},
		},
		Connections: []model.Connection{
			{Source: "S1", Target: "Unknown"},
		},
	}

	data := gen.GeneratePageXML(page, nil, map[string]generator.Position{"S1": {X: 100, Y: 100}})
	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestDrawIOGenerateCompressError(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "A"},
		},
		Compress: true,
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestDrawIOBuildPagesNoComponents(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{},
		Connections: []model.Connection{
			{Source: "A", Target: "B"},
		},
	}

	pages := gen.BuildPages(diagram)
	if len(pages) != 1 {
		t.Errorf("expected 1 page, got %d", len(pages))
	}
}

func TestDrawIOBuildPagesConnectionsOnly(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{},
		Connections: []model.Connection{
			{Source: "A", Target: "B", Page: "Page1"},
			{Source: "C", Target: "D", Page: "Page2"},
		},
	}

	pages := gen.BuildPages(diagram)
	if len(pages) != 3 {
		t.Errorf("expected 3 pages (Page1, Page2, default), got %d", len(pages))
	}
}

func TestDrawIOBuildPagesMixedPages(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "A", Page: "Page1"},
			{Name: "B"},
		},
		Connections: []model.Connection{
			{Source: "A", Target: "B", Page: "Page1"},
			{Source: "C", Target: "D"},
		},
	}

	pages := gen.BuildPages(diagram)
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}
}

func TestGeneratePageXMLWithZeroPosition(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	page := model.Page{
		Name: "Test",
		Components: []model.Component{
			{Name: "S1", Type: model.ComponentTypeService},
			{Name: "S2", Type: model.ComponentTypeDatabase},
		},
		Connections: []model.Connection{},
	}

	positions := map[string]generator.Position{
		"S1": {X: 0, Y: 0},
		"S2": {X: 0, Y: 0},
	}

	data := gen.GeneratePageXML(page, nil, positions)
	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestGeneratePageXMLWithConnections(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	page := model.Page{
		Name: "Test",
		Components: []model.Component{
			{Name: "S1", Type: model.ComponentTypeService},
			{Name: "S2", Type: model.ComponentTypeDatabase},
			{Name: "S3", Type: model.ComponentTypeAPI},
		},
		Connections: []model.Connection{
			{Source: "S1", Target: "S2"},
			{Source: "S2", Target: "S3"},
		},
	}

	positions := map[string]generator.Position{
		"S1": {X: 100, Y: 100},
		"S2": {X: 200, Y: 100},
		"S3": {X: 300, Y: 100},
	}

	data := gen.GeneratePageXML(page, nil, positions)
	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestGridLayoutEdgeCases(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("grid")

	components := []model.Component{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
		{Name: "E"},
		{Name: "F"},
		{Name: "G"},
		{Name: "H"},
		{Name: "I"},
		{Name: "J"},
		{Name: "K"},
		{Name: "L"},
		{Name: "M"},
		{Name: "N"},
		{Name: "O"},
		{Name: "P"},
		{Name: "Q"},
		{Name: "R"},
		{Name: "S"},
		{Name: "T"},
		{Name: "U"},
		{Name: "V"},
		{Name: "W"},
		{Name: "X"},
		{Name: "Y"},
		{Name: "Z"},
	}

	pos := l.Calculate(components, nil)
	if len(pos) != len(components) {
		t.Errorf("expected %d positions, got %d", len(components), len(pos))
	}
}

func TestIsometricLayoutWithConnections(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("isometric")
	components := []model.Component{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
		{Name: "E"},
		{Name: "F"},
		{Name: "G"},
		{Name: "H"},
	}
	connections := []model.Connection{
		{Source: "A", Target: "B"},
		{Source: "B", Target: "C"},
		{Source: "C", Target: "D"},
		{Source: "D", Target: "E"},
		{Source: "E", Target: "F"},
		{Source: "F", Target: "G"},
		{Source: "G", Target: "H"},
	}

	pos := l.Calculate(components, connections)
	if len(pos) != len(components) {
		t.Errorf("expected %d positions, got %d", len(components), len(pos))
	}
}

func TestIsometricLayoutManyComponents(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("isometric")
	components := make([]model.Component, 20)
	for i := range components {
		components[i] = model.Component{Name: string(rune('A' + i))}
	}
	connections := make([]model.Connection, 0)
	for i := 0; i < len(components)-1; i++ {
		connections = append(connections, model.Connection{
			Source: components[i].Name,
			Target: components[i+1].Name,
		})
	}

	pos := l.Calculate(components, connections)
	if len(pos) != len(components) {
		t.Errorf("expected %d positions, got %d", len(components), len(pos))
	}
}

func TestStyleStringEmpty(t *testing.T) {
	t.Parallel()
	style := generator.Style{}
	got := style.String()
	if got != "html=1" {
		t.Errorf("expected only html=1, got: %s", got)
	}
}

func TestParseStyleMultiple(t *testing.T) {
	t.Parallel()
	style := generator.ParseStyle("shape=rectangle;fillColor=#ffffff;strokeColor=#000000;strokeWidth=2;opacity=50;gradientColor=#cccccc;gradientDirection=north;fontSize=14;fontFamily=Arial;fontColor=#000000;fontStyle=1;rounded=1;dashed=1;dashPattern=5 5;shadow=1;glass=1;whiteSpace=wrap;align=center;verticalAlign=middle;image=test.png;imageWidth=100;imageHeight=200;imageAspect=1;edgeStyle=orthogonalEdgeStyle;startArrow=block;endArrow=classic;curved=1;elbow=horizontal;orthogonal=1")

	if style.Shape != "rectangle" {
		t.Errorf("Shape = %q, want rectangle", style.Shape)
	}
	if style.FillColor != "#ffffff" {
		t.Errorf("FillColor = %q, want #ffffff", style.FillColor)
	}
	if style.StrokeColor != "#000000" {
		t.Errorf("StrokeColor = %q, want #000000", style.StrokeColor)
	}
	if style.StrokeWidth != 2 {
		t.Errorf("StrokeWidth = %d, want 2", style.StrokeWidth)
	}
	if style.Opacity != 50 {
		t.Errorf("Opacity = %d, want 50", style.Opacity)
	}
	if style.GradientColor != "#cccccc" {
		t.Errorf("GradientColor = %q, want #cccccc", style.GradientColor)
	}
	if style.GradientDir != "north" {
		t.Errorf("GradientDir = %q, want north", style.GradientDir)
	}
	if style.FontSize != 14 {
		t.Errorf("FontSize = %d, want 14", style.FontSize)
	}
	if style.FontFamily != "Arial" {
		t.Errorf("FontFamily = %q, want Arial", style.FontFamily)
	}
	if style.FontColor != "#000000" {
		t.Errorf("FontColor = %q, want #000000", style.FontColor)
	}
	if style.FontStyle != 1 {
		t.Errorf("FontStyle = %d, want 1", style.FontStyle)
	}
	if !style.Rounded {
		t.Error("Rounded should be true")
	}
	if !style.Dashed {
		t.Error("Dashed should be true")
	}
	if style.DashPattern != "5 5" {
		t.Errorf("DashPattern = %q, want 5 5", style.DashPattern)
	}
	if !style.Shadow {
		t.Error("Shadow should be true")
	}
	if !style.Glass {
		t.Error("Glass should be true")
	}
	if style.WhiteSpace != "wrap" {
		t.Errorf("WhiteSpace = %q, want wrap", style.WhiteSpace)
	}
	if style.Align != "center" {
		t.Errorf("Align = %q, want center", style.Align)
	}
	if style.VerticalAlign != "middle" {
		t.Errorf("VerticalAlign = %q, want middle", style.VerticalAlign)
	}
	if style.Image != "test.png" {
		t.Errorf("Image = %q, want test.png", style.Image)
	}
	if style.ImageWidth != 100 {
		t.Errorf("ImageWidth = %d, want 100", style.ImageWidth)
	}
	if style.ImageHeight != 200 {
		t.Errorf("ImageHeight = %d, want 200", style.ImageHeight)
	}
	if !style.ImageAspect {
		t.Error("ImageAspect should be true")
	}
	if style.EdgeStyle != "orthogonalEdgeStyle" {
		t.Errorf("EdgeStyle = %q, want orthogonalEdgeStyle", style.EdgeStyle)
	}
	if style.StartArrow != "block" {
		t.Errorf("StartArrow = %q, want block", style.StartArrow)
	}
	if style.EndArrow != "classic" {
		t.Errorf("EndArrow = %q, want classic", style.EndArrow)
	}
	if !style.Curved {
		t.Error("Curved should be true")
	}
	if style.Elbow != "horizontal" {
		t.Errorf("Elbow = %q, want horizontal", style.Elbow)
	}
	if !style.Orthogonal {
		t.Error("Orthogonal should be true")
	}
}

func TestMergeStylesAllFields(t *testing.T) {
	t.Parallel()
	base := generator.Style{
		Shape:         "rectangle",
		FillColor:     "#ffffff",
		StrokeColor:   "#000000",
		StrokeWidth:   1,
		Opacity:       50,
		GradientColor: "#cccccc",
		GradientDir:   "north",
		FontSize:      12,
		FontFamily:    "Arial",
		FontColor:     "#000000",
		FontStyle:     0,
		Rounded:       false,
		Dashed:        false,
		DashPattern:   "",
		Shadow:        false,
		Glass:         false,
		WhiteSpace:    "wrap",
		Align:         "center",
		VerticalAlign: "middle",
	}

	override := generator.Style{
		Shape:         "ellipse",
		FillColor:     "#000000",
		StrokeColor:   "#ffffff",
		StrokeWidth:   2,
		Opacity:       80,
		GradientColor: "#eeeeee",
		GradientDir:   "south",
		FontSize:      14,
		FontFamily:    "Helvetica",
		FontColor:     "#ffffff",
		FontStyle:     1,
		Rounded:       true,
		Dashed:        true,
		DashPattern:   "3 3",
		Shadow:        true,
		Glass:         true,
		WhiteSpace:    "nowrap",
		Align:         "left",
		VerticalAlign: "top",
	}

	got := generator.MergeStyles(base, override)

	if got.Shape != "ellipse" {
		t.Errorf("Shape = %q, want ellipse", got.Shape)
	}
	if got.FillColor != "#000000" {
		t.Errorf("FillColor = %q, want #000000", got.FillColor)
	}
	if got.StrokeColor != "#ffffff" {
		t.Errorf("StrokeColor = %q, want #ffffff", got.StrokeColor)
	}
	if got.StrokeWidth != 2 {
		t.Errorf("StrokeWidth = %d, want 2", got.StrokeWidth)
	}
	if got.Opacity != 80 {
		t.Errorf("Opacity = %d, want 80", got.Opacity)
	}
	if got.GradientColor != "#eeeeee" {
		t.Errorf("GradientColor = %q, want #eeeeee", got.GradientColor)
	}
	if got.GradientDir != "south" {
		t.Errorf("GradientDir = %q, want south", got.GradientDir)
	}
	if got.FontSize != 14 {
		t.Errorf("FontSize = %d, want 14", got.FontSize)
	}
	if got.FontFamily != "Helvetica" {
		t.Errorf("FontFamily = %q, want Helvetica", got.FontFamily)
	}
	if got.FontColor != "#ffffff" {
		t.Errorf("FontColor = %q, want #ffffff", got.FontColor)
	}
	if got.FontStyle != 1 {
		t.Errorf("FontStyle = %d, want 1", got.FontStyle)
	}
	if !got.Rounded {
		t.Error("Rounded should be true")
	}
	if !got.Dashed {
		t.Error("Dashed should be true")
	}
	if got.DashPattern != "3 3" {
		t.Errorf("DashPattern = %q, want 3 3", got.DashPattern)
	}
	if !got.Shadow {
		t.Error("Shadow should be true")
	}
	if !got.Glass {
		t.Error("Glass should be true")
	}
	if got.WhiteSpace != "nowrap" {
		t.Errorf("WhiteSpace = %q, want nowrap", got.WhiteSpace)
	}
	if got.Align != "left" {
		t.Errorf("Align = %q, want left", got.Align)
	}
	if got.VerticalAlign != "top" {
		t.Errorf("VerticalAlign = %q, want top", got.VerticalAlign)
	}
}

func TestBuildSwimlanesAllPositions(t *testing.T) {
	t.Parallel()
	components := []model.Component{
		{Name: "S1", Swimlane: "AWS"},
		{Name: "S2", Swimlane: "AWS"},
		{Name: "S3", Swimlane: "AWS"},
	}

	positions := map[string]generator.Position{
		"S1": {X: 100, Y: 100},
		"S2": {X: 200, Y: 100},
		"S3": {X: 300, Y: 100},
	}

	swimlanes := generator.BuildSwimlanes(components, positions)
	if len(swimlanes) != 1 {
		t.Fatalf("expected 1 swimlane, got %d", len(swimlanes))
	}

	if len(swimlanes[0].Children) != 3 {
		t.Errorf("expected 3 children, got %d", len(swimlanes[0].Children))
	}
}

func TestDrawIOGenerateWithCompressFallback(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGeneratorForTest()
	diagram := &model.Diagram{
		Components: []model.Component{
			{Name: "A"},
		},
		Connections: []model.Connection{},
	}

	data, err := gen.Generate(diagram)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}

	content := string(data)
	if !strings.Contains(content, "<mxfile") {
		t.Error("expected mxfile in output")
	}
}

func TestBuildPagesWithEmptyDiagram(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Components:  []model.Component{},
		Connections: []model.Connection{},
	}

	pages := gen.BuildPages(diagram)
	if len(pages) != 1 {
		t.Errorf("expected 1 page, got %d", len(pages))
	}
}

func TestBuildPagesWithPagesInModel(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	diagram := &model.Diagram{
		Pages: []model.Page{
			{Name: "Page1", Components: []model.Component{{Name: "A"}}},
			{Name: "Page2", Components: []model.Component{{Name: "B"}}},
		},
	}

	pages := gen.BuildPages(diagram)
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}
}

func TestGeneratePageXMLNoSwimlane(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	page := model.Page{
		Name: "Test",
		Components: []model.Component{
			{Name: "S1", Type: model.ComponentTypeService},
		},
		Connections: []model.Connection{},
	}

	data := gen.GeneratePageXML(page, nil, map[string]generator.Position{"S1": {X: 100, Y: 100}})
	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestGeneratePageXMLWithSwimlaneParent(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	swimlanes := []generator.Swimlane{
		{Name: "AWS", X: 50, Y: 50, Width: 300, Height: 200},
	}
	page := model.Page{
		Name: "Test",
		Components: []model.Component{
			{Name: "S1", Type: model.ComponentTypeService, Swimlane: "AWS"},
		},
		Connections: []model.Connection{},
	}

	data := gen.GeneratePageXML(page, swimlanes, map[string]generator.Position{"S1": {X: 100, Y: 100}})
	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestGridLayoutWith5Components(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("grid")
	components := []model.Component{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
		{Name: "E"},
	}

	pos := l.Calculate(components, nil)
	if len(pos) != 5 {
		t.Errorf("expected 5 positions, got %d", len(pos))
	}
}

func TestGridLayoutWith6Components(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("grid")
	components := []model.Component{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
		{Name: "E"},
		{Name: "F"},
	}

	pos := l.Calculate(components, nil)
	if len(pos) != 6 {
		t.Errorf("expected 6 positions, got %d", len(pos))
	}
}

func TestGridLayoutWith7Components(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("grid")
	components := []model.Component{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
		{Name: "E"},
		{Name: "F"},
		{Name: "G"},
	}

	pos := l.Calculate(components, nil)
	if len(pos) != 7 {
		t.Errorf("expected 7 positions, got %d", len(pos))
	}
}

func TestGridLayoutWith8Components(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("grid")
	components := []model.Component{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
		{Name: "E"},
		{Name: "F"},
		{Name: "G"},
		{Name: "H"},
	}

	pos := l.Calculate(components, nil)
	if len(pos) != 8 {
		t.Errorf("expected 8 positions, got %d", len(pos))
	}
}

func TestGridLayoutWith1Component(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("grid")
	components := []model.Component{
		{Name: "A"},
	}

	pos := l.Calculate(components, nil)
	if len(pos) != 1 {
		t.Errorf("expected 1 position, got %d", len(pos))
	}
}

func TestGridLayoutWith2Components(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("grid")
	components := []model.Component{
		{Name: "A"},
		{Name: "B"},
	}

	pos := l.Calculate(components, nil)
	if len(pos) != 2 {
		t.Errorf("expected 2 positions, got %d", len(pos))
	}
}

func TestGridLayoutWith3Components(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("grid")
	components := []model.Component{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
	}

	pos := l.Calculate(components, nil)
	if len(pos) != 3 {
		t.Errorf("expected 3 positions, got %d", len(pos))
	}
}

func TestGridLayoutWith4Components(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("grid")
	components := []model.Component{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
	}

	pos := l.Calculate(components, nil)
	if len(pos) != 4 {
		t.Errorf("expected 4 positions, got %d", len(pos))
	}
}

func TestIsometricLayoutWithNoConnectionsAndNoLayers(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("isometric")
	components := []model.Component{}

	pos := l.Calculate(components, nil)
	if len(pos) != 0 {
		t.Errorf("expected 0 positions, got %d", len(pos))
	}
}

func TestIsometricLayoutEmptyComponents(t *testing.T) {
	t.Parallel()
	l := layout.NewLayout("isometric")
	components := []model.Component{}

	pos := l.Calculate(components, nil)
	if len(pos) != 0 {
		t.Errorf("expected 0 positions, got %d", len(pos))
	}
}

func TestStyleStringGradientDirNoGradientColor(t *testing.T) {
	t.Parallel()
	style := generator.Style{
		GradientDir: "north",
	}
	got := style.String()
	if strings.Contains(got, "gradientDirection") {
		t.Errorf("gradientDirection should not be included without gradientColor: %s", got)
	}
}

func TestStyleStringImageWithoutDimensions(t *testing.T) {
	t.Parallel()
	style := generator.Style{
		Image: "data:image/svg+xml;base64,abc",
	}
	got := style.String()
	if !strings.Contains(got, "image=data:image/svg+xml;base64,abc") {
		t.Errorf("missing image: %s", got)
	}
	if strings.Contains(got, "imageWidth=") || strings.Contains(got, "imageHeight=") {
		t.Errorf("unexpected image dimensions: %s", got)
	}
}

func TestStyleStringImageAspectOnly(t *testing.T) {
	t.Parallel()
	style := generator.Style{
		Image:       "data:image/svg+xml;base64,abc",
		ImageAspect: true,
	}
	got := style.String()
	if !strings.Contains(got, "imageAspect=1") {
		t.Errorf("missing imageAspect: %s", got)
	}
}

func TestStyleStringEdgeOptions(t *testing.T) {
	t.Parallel()
	style := generator.Style{
		Curved:     true,
		Elbow:      "horizontal",
		Orthogonal: true,
	}
	got := style.String()
	if !strings.Contains(got, "curved=1") {
		t.Errorf("missing curved: %s", got)
	}
	if !strings.Contains(got, "elbow=horizontal") {
		t.Errorf("missing elbow: %s", got)
	}
	if !strings.Contains(got, "orthogonal=1") {
		t.Errorf("missing orthogonal: %s", got)
	}
}

func TestStyleStringGradientBoth(t *testing.T) {
	t.Parallel()
	style := generator.Style{
		GradientColor: "#cccccc",
		GradientDir:   "north",
	}
	got := style.String()
	if !strings.Contains(got, "gradientColor=#cccccc") {
		t.Errorf("missing gradientColor: %s", got)
	}
	if !strings.Contains(got, "gradientDirection=north") {
		t.Errorf("missing gradientDirection: %s", got)
	}
}

type failingWriter struct{}

func (f *failingWriter) Write(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("write error")
}

type failingWriteCloser struct {
	failingWriter
}

func (f *failingWriteCloser) Close() error {
	return fmt.Errorf("close error")
}

func TestCompressXMLWithFailingWriter(t *testing.T) {
	t.Parallel()
	xml := []byte(`<mxfile><diagram>test</diagram></mxfile>`)
	var buf failingWriter
	err := generator.CompressXMLWriter(xml, &buf)
	if err == nil {
		t.Error("expected error from write, got nil")
	}
}

func TestCompressXMLWithFailingCloser(t *testing.T) {
	t.Parallel()
	xml := []byte(`<mxfile><diagram>test</diagram></mxfile>`)
	var buf failingWriteCloser
	err := generator.CompressXMLWriter(xml, &buf)
	if err == nil {
		t.Error("expected error from close, got nil")
	}
}

func TestStyleStringImageWidth(t *testing.T) {
	t.Parallel()
	style := generator.Style{
		Image:      "test.png",
		ImageWidth: 100,
	}
	got := style.String()
	if !strings.Contains(got, "imageWidth=100") {
		t.Errorf("missing imageWidth: %s", got)
	}
}

func TestStyleStringImageHeight(t *testing.T) {
	t.Parallel()
	style := generator.Style{
		Image:       "test.png",
		ImageHeight: 100,
	}
	got := style.String()
	if !strings.Contains(got, "imageHeight=100") {
		t.Errorf("missing imageHeight: %s", got)
	}
}

func TestGeneratePageXMLWithIsoShape(t *testing.T) {
	t.Parallel()
	gen := generator.NewDrawIOGenerator()
	page := model.Page{
		Name: "Test",
		Components: []model.Component{
			{Name: "Server", Shape: "iso:server"},
			{Name: "DB", Shape: "iso:database"},
		},
		Connections: []model.Connection{},
	}

	positions := map[string]generator.Position{
		"Server": {X: 100, Y: 100},
		"DB":     {X: 200, Y: 100},
	}

	data := gen.GeneratePageXML(page, nil, positions)
	if len(data) == 0 {
		t.Error("generated data should not be empty")
	}
}

func TestBuildSwimlanesPositionsEdge(t *testing.T) {
	t.Parallel()
	components := []model.Component{
		{Name: "S1", Swimlane: "AWS"},
		{Name: "S2", Swimlane: "AWS"},
		{Name: "S3", Swimlane: "AWS"},
		{Name: "S4", Swimlane: "AWS"},
		{Name: "S5", Swimlane: "AWS"},
	}

	positions := map[string]generator.Position{
		"S1": {X: 100, Y: 100},
		"S2": {X: 200, Y: 200},
		"S3": {X: 300, Y: 150},
		"S4": {X: 150, Y: 300},
		"S5": {X: 250, Y: 250},
	}

	swimlanes := generator.BuildSwimlanes(components, positions)
	if len(swimlanes) != 1 {
		t.Fatalf("expected 1 swimlane, got %d", len(swimlanes))
	}
}

func TestCompressXMLEmpty(t *testing.T) {
	t.Parallel()
	compressed, err := generator.CompressXML([]byte{})
	if err != nil {
		t.Fatalf("CompressXML failed: %v", err)
	}
	if len(compressed) == 0 {
		t.Error("compressed data should not be empty")
	}
}

func TestCompressXMLWithInvalidLevel(t *testing.T) {
	t.Parallel()
	xml := []byte(`<mxfile><diagram>test</diagram></mxfile>`)
	_, err := generator.CompressXMLWithLevel(xml, 100)
	if err == nil {
		t.Error("expected error for invalid compression level")
	}
}

func TestCompressAndEncodeWithInvalidLevel(t *testing.T) {
	t.Parallel()
	xml := []byte(`<mxfile><diagram>test</diagram></mxfile>`)
	_, err := generator.CompressAndEncodeWithLevel(xml, 100)
	if err == nil {
		t.Error("expected error for invalid compression level")
	}
}

func TestCompressXMLWithValidLevels(t *testing.T) {
	t.Parallel()
	xml := []byte(`<mxfile><diagram>test</diagram></mxfile>`)

	levels := []int{-1, 1, 9, 0, 2}
	levelNames := []string{"DefaultCompression", "BestSpeedCompression", "BestCompression", "NoCompression", "HuffmanOnly"}
	for i, level := range levels {
		compressed, err := generator.CompressXMLWithLevel(xml, level)
		if err != nil {
			t.Errorf("unexpected error for level %d (%s): %v", level, levelNames[i], err)
		}
		if len(compressed) == 0 {
			t.Errorf("compressed data empty for level %d (%s)", level, levelNames[i])
		}
	}
}

func TestCompressAndEncodeWithValidLevels(t *testing.T) {
	t.Parallel()
	xml := []byte(`<mxfile><diagram>test</diagram></mxfile>`)

	levels := []int{-1, 1, 9, 0, 2}
	levelNames := []string{"DefaultCompression", "BestSpeedCompression", "BestCompression", "NoCompression", "HuffmanOnly"}
	for i, level := range levels {
		encoded, err := generator.CompressAndEncodeWithLevel(xml, level)
		if err != nil {
			t.Errorf("unexpected error for level %d (%s): %v", level, levelNames[i], err)
		}
		if encoded == "" {
			t.Errorf("encoded data empty for level %d (%s)", level, levelNames[i])
		}
	}
}
