package parser

import (
	"fmt"
	"strings"

	"diagram-gen/internal/model"
)

type Annotation struct {
	raw           string
	componentType model.ComponentType
	name          string
	connectsTo    []string
	description   string
	direction     model.ConnectionDirection
}

func ParseAnnotation(tag string) (*Annotation, error) {
	if tag == "" {
		return nil, fmt.Errorf("empty annotation tag")
	}

	tag = strings.Trim(tag, "`")

	ann := &Annotation{raw: tag}

	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "diagram" {
			continue
		}

		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid key-value pair: %s", part)
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		value = strings.Trim(value, `"`)

		switch key {
		case "type":
			ann.componentType = model.ComponentType(value)
		case "name":
			ann.name = value
		case "connectsTo":
			ann.connectsTo = strings.Split(value, ";")
			for i := range ann.connectsTo {
				ann.connectsTo[i] = strings.TrimSpace(ann.connectsTo[i])
			}
		case "description":
			ann.description = value
		case "direction":
			ann.direction = model.ConnectionDirection(value)
		}
	}

	if ann.name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if ann.componentType == "" {
		ann.componentType = model.ComponentTypeService
	}

	return ann, nil
}

func (a *Annotation) ToComponent() model.Component {
	return model.Component{
		Type:        a.componentType,
		Name:        a.name,
		Description: a.description,
		Direction:   a.direction,
	}
}

func (a *Annotation) ToConnections() []model.Connection {
	connections := make([]model.Connection, 0, len(a.connectsTo))
	direction := a.direction
	if direction == "" {
		direction = model.ConnectionDirectionUnidirectional
	}

	for _, target := range a.connectsTo {
		if target == "" {
			continue
		}
		connections = append(connections, model.Connection{
			Source:    a.name,
			Target:    target,
			Direction: direction,
		})
	}
	return connections
}
