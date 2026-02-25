// Package model defines the data models for diagram generation.
package model

// DiagramType defines the type of diagram.
type DiagramType string

const (
	// DiagramTypeArchitecture is an architecture diagram.
	DiagramTypeArchitecture DiagramType = "architecture"
	// DiagramTypeFlowchart is a flowchart diagram.
	DiagramTypeFlowchart DiagramType = "flowchart"
	// DiagramTypeNetwork is a network diagram.
	DiagramTypeNetwork DiagramType = "network"
)

// ComponentType defines the type of component.
type ComponentType string

const (
	// ComponentTypeService is a service component.
	ComponentTypeService ComponentType = "service"
	// ComponentTypeDatabase is a database component.
	ComponentTypeDatabase ComponentType = "database"
	// ComponentTypeQueue is a queue component.
	ComponentTypeQueue ComponentType = "queue"
	// ComponentTypeCache is a cache component.
	ComponentTypeCache ComponentType = "cache"
	// ComponentTypeAPI is an API component.
	ComponentTypeAPI ComponentType = "api"
	// ComponentTypeUser is a user component.
	ComponentTypeUser ComponentType = "user"
	// ComponentTypeExternal is an external component.
	ComponentTypeExternal ComponentType = "external"
	// ComponentTypeStorage is a storage component.
	ComponentTypeStorage ComponentType = "storage"
	// ComponentTypeGateway is a gateway component.
	ComponentTypeGateway ComponentType = "gateway"
	// ComponentTypeUnknown is an unknown component type.
	ComponentTypeUnknown ComponentType = "unknown"
)

// ShapeType defines the shape style for components.
type ShapeType string

const (
	// ShapeTypeRectangle is a rectangle shape.
	ShapeTypeRectangle ShapeType = "rectangle"
	// ShapeTypeRounded is a rounded rectangle shape.
	ShapeTypeRounded ShapeType = "rounded"
	// ShapeTypeEllipse is an ellipse shape.
	ShapeTypeEllipse ShapeType = "ellipse"
	// ShapeTypeCylinder is a cylinder shape.
	ShapeTypeCylinder ShapeType = "cylinder"
	// ShapeTypeIsoServer is an isometric server shape.
	ShapeTypeIsoServer ShapeType = "iso:server"
	// ShapeTypeIsoDatabase is an isometric database shape.
	ShapeTypeIsoDatabase ShapeType = "iso:database"
	// ShapeTypeIsoCloud is an isometric cloud shape.
	ShapeTypeIsoCloud ShapeType = "iso:cloud"
	// ShapeTypeIsoCube is an isometric cube shape.
	ShapeTypeIsoCube ShapeType = "iso:cube"
	// ShapeTypeIsoContainer is an isometric container shape.
	ShapeTypeIsoContainer ShapeType = "iso:container"
)

// ConnectionDirection defines the direction of a connection.
type ConnectionDirection string

const (
	// ConnectionDirectionUnidirectional is unidirectional.
	ConnectionDirectionUnidirectional ConnectionDirection = "unidirectional"
	// ConnectionDirectionBidirectional is bidirectional.
	ConnectionDirectionBidirectional ConnectionDirection = "bidirectional"
)

// Component represents a node in the diagram.
type Component struct {
	Type        ComponentType       `json:"type"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Direction   ConnectionDirection `json:"direction,omitempty"`
	Shape       ShapeType           `json:"shape,omitempty"`
	Page        string              `json:"page,omitempty"`
	Swimlane    string              `json:"swimlane,omitempty"`
	Style       string              `json:"style,omitempty"`
	X           int                 `json:"x,omitempty"`
	Y           int                 `json:"y,omitempty"`
}

// Connection represents an edge between two components.
type Connection struct {
	Source     string              `json:"source"`
	Target     string              `json:"target"`
	Direction  ConnectionDirection `json:"direction,omitempty"`
	Label      string              `json:"label,omitempty"`
	Page       string              `json:"page,omitempty"`
	EdgeStyle  string              `json:"edgeStyle,omitempty"`
	StartArrow string              `json:"startArrow,omitempty"`
	EndArrow   string              `json:"endArrow,omitempty"`
}

// Diagram represents a complete diagram with components and connections.
type Diagram struct {
	Type        DiagramType  `json:"type"`
	Components  []Component  `json:"components"`
	Connections []Connection `json:"connections"`
	Pages       []Page       `json:"pages,omitempty"`
	Layout      string       `json:"layout,omitempty"`
	Compress    bool         `json:"compress,omitempty"`
}

// Page represents a page in a multi-page diagram.
type Page struct {
	Name        string       `json:"name"`
	Components  []Component  `json:"components,omitempty"`
	Connections []Connection `json:"connections,omitempty"`
}

// AddComponent adds a component to the diagram.
func (d *Diagram) AddComponent(c Component) {
	d.Components = append(d.Components, c)
}

// AddConnection adds a connection to the diagram.
func (d *Diagram) AddConnection(c Connection) {
	d.Connections = append(d.Connections, c)
}

// GetComponentByName returns a component by its name.
func (d *Diagram) GetComponentByName(name string) *Component {
	for i := range d.Components {
		if d.Components[i].Name == name {
			return &d.Components[i]
		}
	}
	return nil
}
