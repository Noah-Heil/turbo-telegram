// Package archparser provides parsing of Go source files for diagram annotations.
package archparser

import (
	"fmt"
	"strings"

	"diagram-gen/internal/model"
)

// Annotation represents a parsed diagram annotation.
type Annotation struct {
	Raw           string
	ComponentType model.ComponentType
	Name          string
	ConnectsTo    []string
	Description   string
	Direction     model.ConnectionDirection
	Shape         model.ShapeType
	Page          string
	Swimlane      string
	Style         string
	EdgeStyle     string
	StartArrow    string
	EndArrow      string
}

// ParseAnnotation parses a diagram annotation string.
func ParseAnnotation(tag string) (*Annotation, error) {
	if tag == "" {
		return nil, fmt.Errorf("empty annotation tag")
	}

	tag = strings.Trim(tag, "`")

	ann := &Annotation{Raw: tag}

	parts := splitKeyValuePairs(tag)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "diagram" {
			continue
		}

		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		value = strings.Trim(value, `"`)

		switch key {
		case "type":
			ann.ComponentType = model.ComponentType(value)
		case "name":
			ann.Name = value
		case "connectsTo":
			ann.ConnectsTo = strings.Split(value, ";")
			for i := range ann.ConnectsTo {
				ann.ConnectsTo[i] = strings.TrimSpace(ann.ConnectsTo[i])
			}
		case "description":
			ann.Description = value
		case "direction":
			ann.Direction = model.ConnectionDirection(value)
		case "shape":
			ann.Shape = model.ShapeType(value)
		case "page":
			ann.Page = value
		case "swimlane":
			ann.Swimlane = value
		case "fillColor", "strokeColor", "fontColor", "gradientColor",
			"fontSize", "strokeWidth", "opacity", "rounded",
			"dashed", "shadow", "glass":
			if ann.Style != "" {
				ann.Style += ";"
			}
			ann.Style += key + "=" + value
		case "edgeStyle":
			ann.EdgeStyle = value
		case "startArrow":
			ann.StartArrow = value
		case "endArrow":
			ann.EndArrow = value
		}
	}

	if ann.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if ann.ComponentType == "" {
		ann.ComponentType = model.ComponentTypeService
	}

	return ann, nil
}

func splitKeyValuePairs(s string) []string {
	isInValue := false

	insideValue := make(map[int]bool)
	for i := 0; i < len(s); i++ {
		if isInValue && s[i] == ';' {
			insideValue[i] = true
		}

		if s[i] == '=' {
			isInValue = true
		} else if isInValue && (s[i] == ',' || s[i] == ';') {
			isInValue = false
		}
	}

	var result []string
	var current strings.Builder

	for i := 0; i < len(s); i++ {
		c := s[i]

		if c == ';' && insideValue[i] {
			current.WriteByte(c)
			continue
		}

		if c == ',' || c == ';' {
			part := strings.TrimSpace(current.String())
			if part != "" {
				result = append(result, part)
			}
			current.Reset()
			continue
		}

		current.WriteByte(c)
	}

	part := strings.TrimSpace(current.String())
	if part != "" {
		result = append(result, part)
	}

	return result
}

// ToComponent converts the annotation to a Component model.
func (a *Annotation) ToComponent() model.Component {
	return model.Component{
		Type:        a.ComponentType,
		Name:        a.Name,
		Description: a.Description,
		Direction:   a.Direction,
		Shape:       a.Shape,
		Page:        a.Page,
		Swimlane:    a.Swimlane,
		Style:       a.Style,
	}
}

// ToConnections converts the annotation to Connection models.
func (a *Annotation) ToConnections() []model.Connection {
	connections := make([]model.Connection, 0, len(a.ConnectsTo))
	direction := a.Direction
	if direction == "" {
		direction = model.ConnectionDirectionUnidirectional
	}

	for _, target := range a.ConnectsTo {
		if target == "" {
			continue
		}
		connections = append(connections, model.Connection{
			Source:     a.Name,
			Target:     target,
			Direction:  direction,
			EdgeStyle:  a.EdgeStyle,
			StartArrow: a.StartArrow,
			EndArrow:   a.EndArrow,
			Page:       a.Page,
		})
	}
	return connections
}
