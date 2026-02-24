package generator_test

import (
	"bytes"
	"strings"
	"testing"

	"diagram-gen/internal/generator"
	"diagram-gen/internal/generator/layout"
	"diagram-gen/internal/model"
)

func TestStyleString(t *testing.T) {
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
			got := tt.style.String()
			if got != tt.expected {
				t.Errorf("generator.Style.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestParseStyle(t *testing.T) {
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
			got := tt.shape.IsIsometric()
			if got != tt.expected {
				t.Errorf("IsIsometric() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestShapeTypeIsBasic(t *testing.T) {
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
			got := tt.shape.IsBasic()
			if got != tt.expected {
				t.Errorf("IsBasic() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSwimlaneBuild(t *testing.T) {
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
			got := generator.GetShapeStyle(tt.shape)
			if len(got) < len(tt.expected) || got[:len(tt.expected)] != tt.expected {
				t.Errorf("GetShapeStyle(%q) = %q, want prefix %q", tt.shape, got, tt.expected)
			}
		})
	}
}

func TestGetDefaultShapeForComponentType(t *testing.T) {
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
			got := generator.GetDefaultShapeForComponentType(tt.compType)
			if got != tt.expected {
				t.Errorf("GetDefaultShapeForComponentType(%q) = %q, want %q", tt.compType, got, tt.expected)
			}
		})
	}
}

func TestParseStyleOpacity(t *testing.T) {
	style := generator.ParseStyle("opacity=50")
	if style.Opacity != 50 {
		t.Errorf("Opacity = %d, want 50", style.Opacity)
	}
}

func TestParseStyleStrokeWidth(t *testing.T) {
	style := generator.ParseStyle("strokeWidth=3")
	if style.StrokeWidth != 3 {
		t.Errorf("StrokeWidth = %d, want 3", style.StrokeWidth)
	}
}

func TestParseStyleFontFamily(t *testing.T) {
	style := generator.ParseStyle("fontFamily=Arial")
	if style.FontFamily != "Arial" {
		t.Errorf("FontFamily = %q, want Arial", style.FontFamily)
	}
}

func TestParseStyleFontStyle(t *testing.T) {
	style := generator.ParseStyle("fontStyle=1")
	if style.FontStyle != 1 {
		t.Errorf("FontStyle = %d, want 1", style.FontStyle)
	}
}

func TestParseStyleImage(t *testing.T) {
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
	style := generator.ParseStyle("curved=1;elbow=horizontal")
	if !style.Curved {
		t.Error("Curved should be true")
	}
	if style.Elbow != "horizontal" {
		t.Errorf("Elbow = %q, want horizontal", style.Elbow)
	}
}

func TestParseStyleWhiteSpace(t *testing.T) {
	style := generator.ParseStyle("whiteSpace=wrap")
	if style.WhiteSpace != "wrap" {
		t.Errorf("WhiteSpace = %q, want wrap", style.WhiteSpace)
	}
}

func TestParseStyleAlign(t *testing.T) {
	style := generator.ParseStyle("align=center;verticalAlign=bottom")
	if style.Align != "center" {
		t.Errorf("Align = %q, want center", style.Align)
	}
	if style.VerticalAlign != "bottom" {
		t.Errorf("VerticalAlign = %q, want bottom", style.VerticalAlign)
	}
}

func TestMergeStylesOverrideAll(t *testing.T) {
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
	positions := map[string]generator.Position{
		"Service1": {X: 100, Y: 100},
	}

	swimlanes := generator.BuildSwimlanes(nil, positions)
	if len(swimlanes) != 0 {
		t.Errorf("expected 0 swimlanes for nil components, got %d", len(swimlanes))
	}
}

func TestSwimlaneNoSwimlane(t *testing.T) {
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
	gen := generator.NewDrawIOGenerator()
	if gen.Format() != "drawio" {
		t.Errorf("Format() = %q, want drawio", gen.Format())
	}
}

func TestDrawIOGeneratorWithOptions(t *testing.T) {
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
