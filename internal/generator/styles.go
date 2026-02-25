package generator

import (
	"fmt"
	"strconv"
	"strings"
)

func parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

// Style defines draw.io style properties for diagram elements.
type Style struct {
	Shape         string
	FillColor     string
	StrokeColor   string
	StrokeWidth   int
	Opacity       int
	GradientColor string
	GradientDir   string
	FontSize      int
	FontFamily    string
	FontColor     string
	FontStyle     int
	Rounded       bool
	Dashed        bool
	DashPattern   string
	Shadow        bool
	Glass         bool
	WhiteSpace    string
	Align         string
	VerticalAlign string
	Image         string
	ImageWidth    int
	ImageHeight   int
	ImageAspect   bool
	EdgeStyle     string
	StartArrow    string
	EndArrow      string
	Curved        bool
	Elbow         string
	Orthogonal    bool
}

const (
	// FontStyleBold is bold font style.
	FontStyleBold = 1
	// FontStyleItalic is italic font style.
	FontStyleItalic = 2
	// FontStyleUnderline is underline font style.
	FontStyleUnderline = 4

	// GradientDirNorth is north gradient direction.
	GradientDirNorth = "north"
	// GradientDirSouth is south gradient direction.
	GradientDirSouth = "south"
	// GradientDirEast is east gradient direction.
	GradientDirEast = "east"
	// GradientDirWest is west gradient direction.
	GradientDirWest = "west"

	// AlignLeft is left text alignment.
	AlignLeft = "left"
	// AlignCenter is center text alignment.
	AlignCenter = "center"
	// AlignRight is right text alignment.
	AlignRight = "right"

	// VAlignTop is top vertical alignment.
	VAlignTop = "top"
	// VAlignMiddle is middle vertical alignment.
	VAlignMiddle = "middle"
	// VAlignBottom is bottom vertical alignment.
	VAlignBottom = "bottom"

	// WhiteSpaceWrap enables text wrapping.
	WhiteSpaceWrap = "wrap"

	// EdgeStyleOrthogonal is orthogonal edge style.
	EdgeStyleOrthogonal = "orthogonalEdgeStyle"
	// EdgeStyleElbow is elbow edge style.
	EdgeStyleElbow = "elbowEdgeStyle"
	// EdgeStyleCurved is curved edge style.
	EdgeStyleCurved = "curvedEdgeStyle"

	// ArrowBlock is a block arrow.
	ArrowBlock = "block"
	// ArrowOpen is an open arrow.
	ArrowOpen = "open"
	// ArrowClassic is a classic arrow.
	ArrowClassic = "classic"
	// ArrowDiamond is a diamond arrow.
	ArrowDiamond = "diamond"
	// ArrowNone is no arrow.
	ArrowNone = "none"
)

func (s Style) String() string {
	var parts []string

	if s.Shape != "" {
		parts = append(parts, "shape="+s.Shape)
	}
	if s.FillColor != "" {
		parts = append(parts, "fillColor="+s.FillColor)
	}
	if s.StrokeColor != "" {
		parts = append(parts, "strokeColor="+s.StrokeColor)
	}
	if s.StrokeWidth > 0 {
		parts = append(parts, fmt.Sprintf("strokeWidth=%d", s.StrokeWidth))
	}
	if s.Opacity > 0 && s.Opacity < 100 {
		parts = append(parts, fmt.Sprintf("opacity=%d", s.Opacity))
	}
	if s.GradientColor != "" {
		parts = append(parts, "gradientColor="+s.GradientColor)
		if s.GradientDir != "" {
			parts = append(parts, "gradientDirection="+s.GradientDir)
		}
	}
	if s.FontSize > 0 {
		parts = append(parts, fmt.Sprintf("fontSize=%d", s.FontSize))
	}
	if s.FontFamily != "" {
		parts = append(parts, "fontFamily="+s.FontFamily)
	}
	if s.FontColor != "" {
		parts = append(parts, "fontColor="+s.FontColor)
	}
	if s.FontStyle > 0 {
		parts = append(parts, fmt.Sprintf("fontStyle=%d", s.FontStyle))
	}
	if s.Rounded {
		parts = append(parts, "rounded=1")
	}
	if s.Dashed {
		parts = append(parts, "dashed=1")
	}
	if s.DashPattern != "" {
		parts = append(parts, "dashPattern="+s.DashPattern)
	}
	if s.Shadow {
		parts = append(parts, "shadow=1")
	}
	if s.Glass {
		parts = append(parts, "glass=1")
	}
	if s.WhiteSpace != "" {
		parts = append(parts, "whiteSpace="+s.WhiteSpace)
	}
	if s.Align != "" {
		parts = append(parts, "align="+s.Align)
	}
	if s.VerticalAlign != "" {
		parts = append(parts, "verticalAlign="+s.VerticalAlign)
	}
	if s.Image != "" {
		parts = append(parts, "image="+s.Image)
		if s.ImageWidth > 0 {
			parts = append(parts, fmt.Sprintf("imageWidth=%d", s.ImageWidth))
		}
		if s.ImageHeight > 0 {
			parts = append(parts, fmt.Sprintf("imageHeight=%d", s.ImageHeight))
		}
		if s.ImageAspect {
			parts = append(parts, "imageAspect=1")
		}
	}
	if s.EdgeStyle != "" {
		parts = append(parts, "edgeStyle="+s.EdgeStyle)
	}
	if s.StartArrow != "" {
		parts = append(parts, "startArrow="+s.StartArrow)
	}
	if s.EndArrow != "" {
		parts = append(parts, "endArrow="+s.EndArrow)
	}
	if s.Curved {
		parts = append(parts, "curved=1")
	}
	if s.Elbow != "" {
		parts = append(parts, "elbow="+s.Elbow)
	}
	if s.Orthogonal {
		parts = append(parts, "orthogonal=1")
	}

	parts = append(parts, "html=1")

	return strings.Join(parts, ";")
}

// ParseStyle parses a draw.io style string into a Style struct.
func ParseStyle(styleStr string) Style {
	style := Style{}
	if styleStr == "" {
		return style
	}

	parts := strings.Split(styleStr, ";")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "shape":
			style.Shape = value
		case "fillColor":
			style.FillColor = value
		case "strokeColor":
			style.StrokeColor = value
		case "strokeWidth":
			style.StrokeWidth = parseInt(value, 0)
		case "opacity":
			style.Opacity = parseInt(value, 0)
		case "gradientColor":
			style.GradientColor = value
		case "gradientDirection":
			style.GradientDir = value
		case "fontSize":
			style.FontSize = parseInt(value, 0)
		case "fontFamily":
			style.FontFamily = value
		case "fontColor":
			style.FontColor = value
		case "fontStyle":
			style.FontStyle = parseInt(value, 0)
		case "rounded":
			style.Rounded = value == "1"
		case "dashed":
			style.Dashed = value == "1"
		case "dashPattern":
			style.DashPattern = value
		case "shadow":
			style.Shadow = value == "1"
		case "glass":
			style.Glass = value == "1"
		case "whiteSpace":
			style.WhiteSpace = value
		case "align":
			style.Align = value
		case "verticalAlign":
			style.VerticalAlign = value
		case "image":
			style.Image = value
		case "imageWidth":
			style.ImageWidth = parseInt(value, 0)
		case "imageHeight":
			style.ImageHeight = parseInt(value, 0)
		case "imageAspect":
			style.ImageAspect = value == "1"
		case "edgeStyle":
			style.EdgeStyle = value
		case "startArrow":
			style.StartArrow = value
		case "endArrow":
			style.EndArrow = value
		case "curved":
			style.Curved = value == "1"
		case "elbow":
			style.Elbow = value
		case "orthogonal":
			style.Orthogonal = value == "1"
		}
	}

	return style
}

// MergeStyles merges two styles, with override taking precedence.
func MergeStyles(base, override Style) Style {
	if override.Shape != "" {
		base.Shape = override.Shape
	}
	if override.FillColor != "" {
		base.FillColor = override.FillColor
	}
	if override.StrokeColor != "" {
		base.StrokeColor = override.StrokeColor
	}
	if override.StrokeWidth > 0 {
		base.StrokeWidth = override.StrokeWidth
	}
	if override.Opacity > 0 {
		base.Opacity = override.Opacity
	}
	if override.GradientColor != "" {
		base.GradientColor = override.GradientColor
	}
	if override.GradientDir != "" {
		base.GradientDir = override.GradientDir
	}
	if override.FontSize > 0 {
		base.FontSize = override.FontSize
	}
	if override.FontFamily != "" {
		base.FontFamily = override.FontFamily
	}
	if override.FontColor != "" {
		base.FontColor = override.FontColor
	}
	if override.FontStyle > 0 {
		base.FontStyle = override.FontStyle
	}
	if override.Rounded {
		base.Rounded = true
	}
	if override.Dashed {
		base.Dashed = true
	}
	if override.DashPattern != "" {
		base.DashPattern = override.DashPattern
	}
	if override.Shadow {
		base.Shadow = true
	}
	if override.Glass {
		base.Glass = true
	}
	if override.WhiteSpace != "" {
		base.WhiteSpace = override.WhiteSpace
	}
	if override.Align != "" {
		base.Align = override.Align
	}
	if override.VerticalAlign != "" {
		base.VerticalAlign = override.VerticalAlign
	}

	return base
}
