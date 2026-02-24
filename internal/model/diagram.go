package model

type DiagramType string

const (
	DiagramTypeArchitecture DiagramType = "architecture"
	DiagramTypeFlowchart    DiagramType = "flowchart"
	DiagramTypeNetwork      DiagramType = "network"
)

type ComponentType string

const (
	ComponentTypeService  ComponentType = "service"
	ComponentTypeDatabase ComponentType = "database"
	ComponentTypeQueue    ComponentType = "queue"
	ComponentTypeCache    ComponentType = "cache"
	ComponentTypeAPI      ComponentType = "api"
	ComponentTypeUser     ComponentType = "user"
	ComponentTypeExternal ComponentType = "external"
	ComponentTypeStorage  ComponentType = "storage"
	ComponentTypeGateway  ComponentType = "gateway"
	ComponentTypeUnknown  ComponentType = "unknown"
)

type ConnectionDirection string

const (
	ConnectionDirectionUnidirectional ConnectionDirection = "unidirectional"
	ConnectionDirectionBidirectional  ConnectionDirection = "bidirectional"
)

type Component struct {
	Type        ComponentType       `json:"type"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Direction   ConnectionDirection `json:"direction,omitempty"`
}

type Connection struct {
	Source    string              `json:"source"`
	Target    string              `json:"target"`
	Direction ConnectionDirection `json:"direction,omitempty"`
	Label     string              `json:"label,omitempty"`
}

type Diagram struct {
	Type        DiagramType  `json:"type"`
	Components  []Component  `json:"components"`
	Connections []Connection `json:"connections"`
}

func (d *Diagram) AddComponent(c Component) {
	d.Components = append(d.Components, c)
}

func (d *Diagram) AddConnection(c Connection) {
	d.Connections = append(d.Connections, c)
}

func (d *Diagram) GetComponentByName(name string) *Component {
	for i := range d.Components {
		if d.Components[i].Name == name {
			return &d.Components[i]
		}
	}
	return nil
}
