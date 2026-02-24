package generator

// ShapeType defines the shape style for diagram components.
type ShapeType string

const (
	// ShapeRectangle is a rectangle shape.
	ShapeRectangle ShapeType = "rectangle"
	// ShapeEllipse is an ellipse shape.
	ShapeEllipse ShapeType = "ellipse"
	// ShapeRounded is a rounded rectangle shape.
	ShapeRounded ShapeType = "rounded"
	// ShapeRhombus is a rhombus shape.
	ShapeRhombus ShapeType = "rhombus"
	// ShapeParallelogram is a parallelogram shape.
	ShapeParallelogram ShapeType = "parallelogram"
	// ShapeCylinder is a cylinder shape.
	ShapeCylinder ShapeType = "cylinder"
	// ShapeDocument is a document shape.
	ShapeDocument ShapeType = "document"
	// ShapeSwimlane is a swimlane shape.
	ShapeSwimlane ShapeType = "swimlane"
	// ShapeTriangle is a triangle shape.
	ShapeTriangle ShapeType = "triangle"
	// ShapeHexagon is a hexagon shape.
	ShapeHexagon ShapeType = "hexagon"
	// ShapeCloud is a cloud shape.
	ShapeCloud ShapeType = "cloud"
	// ShapeInternal is an internal shape.
	ShapeInternal ShapeType = "internal"
	// ShapeExternal is an external shape.
	ShapeExternal ShapeType = "external"
	// ShapeFolder is a folder shape.
	ShapeFolder ShapeType = "folder"

	// ShapeIsoCube is an isometric cube shape.
	ShapeIsoCube ShapeType = "mxgraph.isometric.cube"
	// ShapeIsoServer is an isometric server shape.
	ShapeIsoServer ShapeType = "mxgraph.isometric.server"
	// ShapeIsoDatabase is an isometric database shape.
	ShapeIsoDatabase ShapeType = "mxgraph.isometric.database"
	// ShapeIsoContainer is an isometric container shape.
	ShapeIsoContainer ShapeType = "mxgraph.isometric.container"
	// ShapeIsoCloud is an isometric cloud shape.
	ShapeIsoCloud ShapeType = "mxgraph.isometric.cloud"
	// ShapeIsoNetwork is an isometric network shape.
	ShapeIsoNetwork ShapeType = "mxgraph.isometric.network"
	// ShapeIsoCylinder is an isometric cylinder shape.
	ShapeIsoCylinder ShapeType = "mxgraph.isometric.cylinder"
	// ShapeImage is an image shape.
	ShapeImage ShapeType = "image"
)

// IsIsometric returns true if the shape is an isometric shape.
func (s ShapeType) IsIsometric() bool {
	switch s {
	case ShapeIsoCube, ShapeIsoServer, ShapeIsoDatabase,
		ShapeIsoContainer, ShapeIsoCloud, ShapeIsoNetwork, ShapeIsoCylinder:
		return true
	}
	return false
}

// IsBasic returns true if the shape is a basic shape.
func (s ShapeType) IsBasic() bool {
	switch s {
	case ShapeRectangle, ShapeEllipse, ShapeRounded, ShapeRhombus,
		ShapeParallelogram, ShapeCylinder, ShapeDocument, ShapeSwimlane,
		ShapeTriangle, ShapeHexagon, ShapeCloud, ShapeInternal,
		ShapeExternal, ShapeFolder:
		return true
	}
	return false
}

// GetShapeStyle returns the draw.io style string for a shape type.
func GetShapeStyle(shape ShapeType) string {
	switch shape {
	case ShapeRectangle:
		return "shape=rectangle;whiteSpace=wrap;html=1;"
	case ShapeEllipse:
		return "shape=ellipse;whiteSpace=wrap;html=1;"
	case ShapeRounded:
		return "shape=rounded;whiteSpace=wrap;html=1;rounded=1;"
	case ShapeRhombus:
		return "shape=rhombus;whiteSpace=wrap;html=1;"
	case ShapeParallelogram:
		return "shape=parallelogram;perimeter=parallelogramPerimeter;whiteSpace=wrap;html=1;"
	case ShapeCylinder:
		return "shape=cylinder;whiteSpace=wrap;html=1;boundedLbl=1;backgroundOutline=1;size=10;"
	case ShapeDocument:
		return "shape=document;whiteSpace=wrap;html=1;boundedLbl=1;"
	case ShapeSwimlane:
		return "shape=swimlane;horizontal=1;whiteSpace=wrap;html=1;"
	case ShapeTriangle:
		return "shape=triangle;whiteSpace=wrap;html=1;"
	case ShapeHexagon:
		return "shape=hexagon;perimeter=hexagonPerimeter;whiteSpace=wrap;html=1;"
	case ShapeCloud:
		return "shape=cloud;whiteSpace=wrap;html=1;"
	case ShapeInternal:
		return "shape=internal;whiteSpace=wrap;html=1;"
	case ShapeExternal:
		return "shape=external;whiteSpace=wrap;html=1;"
	case ShapeFolder:
		return "shape=folder;whiteSpace=wrap;html=1;"
	case ShapeIsoCube:
		return "shape=mxgraph.isometric.cube;"
	case ShapeIsoServer:
		return "shape=mxgraph.isometric.server;"
	case ShapeIsoDatabase:
		return "shape=mxgraph.isometric.database;"
	case ShapeIsoContainer:
		return "shape=mxgraph.isometric.container;"
	case ShapeIsoCloud:
		return "shape=mxgraph.isometric.cloud;"
	case ShapeIsoNetwork:
		return "shape=mxgraph.isometric.network;"
	case ShapeIsoCylinder:
		return "shape=mxgraph.isometric.cylinder;"
	default:
		return "shape=rectangle;whiteSpace=wrap;html=1;"
	}
}

// GetDefaultShapeForComponentType returns the default shape for a component type.
func GetDefaultShapeForComponentType(compType string) ShapeType {
	switch compType {
	case "service", "api", "gateway":
		return ShapeRounded
	case "database", "storage":
		return ShapeCylinder
	case "queue":
		return ShapeParallelogram
	case "cache":
		return ShapeRounded
	case "user":
		return ShapeEllipse
	case "external":
		return ShapeDocument
	case "iso:server":
		return ShapeIsoServer
	case "iso:database":
		return ShapeIsoDatabase
	case "iso:container":
		return ShapeIsoContainer
	case "iso:cloud":
		return ShapeIsoCloud
	case "iso:network":
		return ShapeIsoNetwork
	case "iso:cube":
		return ShapeIsoCube
	case "iso:cylinder":
		return ShapeIsoCylinder
	default:
		return ShapeRectangle
	}
}
