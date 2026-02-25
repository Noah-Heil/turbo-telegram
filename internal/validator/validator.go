// Package validator provides validation for diagram models.
package validator

import (
	"fmt"

	"diagram-gen/internal/model"
)

var validComponentTypes = map[model.ComponentType]bool{
	model.ComponentTypeService:  true,
	model.ComponentTypeDatabase: true,
	model.ComponentTypeQueue:    true,
	model.ComponentTypeCache:    true,
	model.ComponentTypeAPI:      true,
	model.ComponentTypeUser:     true,
	model.ComponentTypeExternal: true,
	model.ComponentTypeStorage:  true,
	model.ComponentTypeGateway:  true,
}

// ValidateDiagram validates a diagram model.
func ValidateDiagram(diagram *model.Diagram) error {
	if diagram == nil {
		return fmt.Errorf("diagram is nil")
	}

	if len(diagram.Components) == 0 {
		return fmt.Errorf("no components found in diagram")
	}

	componentNames := make(map[string]bool)
	for _, comp := range diagram.Components {
		if comp.Name == "" {
			return fmt.Errorf("component has empty name")
		}
		if _, exists := componentNames[comp.Name]; exists {
			return fmt.Errorf("duplicate component name: %s", comp.Name)
		}
		componentNames[comp.Name] = true

		if !validComponentTypes[comp.Type] {
			return fmt.Errorf("unknown component type: %s (valid types: service, database, queue, cache, api, user, external, storage, gateway)", comp.Type)
		}
	}

	for _, conn := range diagram.Connections {
		if conn.Source == "" || conn.Target == "" {
			return fmt.Errorf("connection has empty source or target")
		}
		if !componentNames[conn.Source] {
			return fmt.Errorf("connection references unknown source: %s", conn.Source)
		}
		if !componentNames[conn.Target] {
			return fmt.Errorf("connection references unknown target: %s", conn.Target)
		}
	}

	return nil
}

// ValidateComponentType checks if a component type is valid.
func ValidateComponentType(componentType model.ComponentType) bool {
	return validComponentTypes[componentType]
}
