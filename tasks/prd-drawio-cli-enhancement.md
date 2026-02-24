# PRD: Enhanced Draw.io CLI Diagram Generator

## 1. Introduction/Overview

The current CLI generates basic architecture diagrams with limited shape types and a simple grid layout. This feature aims to enhance the CLI to support the full breadth of the draw.io/mxGraph file format, including isometric shapes, advanced styling, swimlanes, custom layouts, and more. The goal is to enable infrastructure teams to generate professional-grade architecture diagrams directly from Go code annotations.

## 2. Goals

1. **Expand Shape Support**: Add isometric shapes (cube, server, database, container, cloud) and all standard mxGraph shapes from the draw.io format.
2. **Advanced Styling**: Support gradients, shadows, opacity, custom fonts, images, and all style parameters defined in fullDrawioOptions.md.
3. **Multiple Diagram Pages**: Enable multi-page diagram generation with separate pages for different views.
4. **Flexible Layout Algorithms**: Implement layered, orthogonal, and isometric layout engines.
5. **Swimlanes and Layers**: Support swimlane containers for grouping components and layer management.
6. **Compression Support**: Optionally output compressed+base64-encoded draw.io files.

## 3. User Stories

| ID | User Story |
|----|------------|
| US1 | As an infrastructure engineer, I want to generate isometric architecture diagrams so that my cloud diagrams look professional and visually appealing. |
| US2 | As a DevOps engineer, I want to specify custom styling (colors, fonts, gradients) in my annotations so the diagram matches my organization's branding. |
| US3 | As a platform engineer, I want to organize components into swimlanes (e.g., "AWS", "On-premise", "Kubernetes") so the diagram shows clear boundaries. |
| US4 | As a user, I want to choose between CLI flags, code annotations, or a config file to specify advanced options so I can use whichever method is most convenient. |
| US5 | As a user, I want to generate multi-page diagrams (e.g., "Network View", "Data Flow", "Components") so complex architectures are organized. |
| US6 | As a user, I want to embed custom SVG icons for cloud resources so the diagram uses official provider icons. |

## 4. Functional Requirements

### 4.1 Extended Shape Support

| ID | Requirement |
|----|-------------|
| FR1.1 | The generator must support all isometric shapes: `mxgraph.isometric.cube`, `mxgraph.isometric.server`, `mxgraph.isometric.database`, `mxgraph.isometric.container`, `mxgraph.isometric.cloud`, `mxgraph.isometric.network`, `mxgraph.isometric.cylinder`. |
| FR1.2 | The generator must support basic shapes: rectangle, ellipse, rhombus, parallelogram, cylinder, document, swimlane. |
| FR1.3 | The generator must support image embedding via base64-encoded SVGs. |
| FR1.4 | Users must be able to specify shape type via annotation: `@//diagram:type=iso:server` or via CLI `--shape iso:server`. |

### 4.2 Advanced Styling

| ID | Requirement |
|----|-------------|
| FR2.1 | The generator must support fill color, stroke color, stroke width, and opacity. |
| FR2.2 | The generator must support gradient fills with direction (north/south/east/west). |
| FR2.3 | The generator must support text styling: font family, font size, font color, bold/italic/underline. |
| FR2.4 | The generator must support dashed lines and custom dash patterns. |
| FR2.5 | The generator must support shadow and glass effects. |
| FR2.6 | Style must be specifiable via annotation: `@//diagram:fillColor=#dae8fc;strokeColor=#6c8ebf;gradientColor=#ffffff` or via config file. |

### 4.3 Multi-Page Support

| ID | Requirement |
|----|-------------|
| FR3.1 | The generator must support multiple `<diagram>` pages in a single `.drawio` file. |
| FR3.2 | Users must be able to group components into pages via annotation: `@//diagram:page=network-view`. |
| FR3.3 | The CLI must support `--page` flag to specify which page to generate. |

### 4.4 Layout Engines

| ID | Requirement |
|----|-------------|
| FR4.1 | The generator must support layered layout (default, hierarchical top-to-bottom). |
| FR4.2 | The generator must support isometric projection layout using 30° transformation. |
| FR4.3 | The generator must support grid layout (current default). |
| FR4.4 | Users must be able to select layout via CLI: `--layout layered`, `--layout isometric`, `--layout grid`. |

### 4.5 Swimlanes and Groups

| ID | Requirement |
|----|-------------|
| FR5.1 | The generator must support swimlane containers for grouping components. |
| FR5.2 | Users must be able to specify swimlane membership via annotation: `@//diagram:swimlane=AWS-Region-1`. |
| FR5.3 | Swimlanes must render with proper parent-child hierarchy in the mxGraph model. |

### 4.6 Edge Styles

| ID | Requirement |
|----|-------------|
| FR6.1 | The generator must support edge types: straight, orthogonal (right-angle), curved, elbow. |
| FR6.2 | The generator must support arrow styles: block, open, classic, diamond. |
| FR6.3 | The generator must support start and end arrowheads independently. |
| FR6.4 | Users must be able to specify edge style via annotation: `@//diagram:edgeStyle=elbowEdgeStyle;endArrow=block`. |

### 4.7 Compression

| ID | Requirement |
|----|-------------|
| FR7.1 | The generator must support optional compression (deflate + base64) for output. |
| FR7.2 | Users must be able to enable compression via CLI flag: `--compress`. |

### 4.8 Configuration Methods

| ID | Requirement |
|----|-------------|
| FR8.1 | CLI flags must support: `--isometric`, `--layout`, `--shape`, `--compress`, `--output`, `--config`. |
| Code annotations must support extended syntax: `@//diagram:key=value;key2=value2`. |
| FR8.2 | Config file (`.diagram-gen.yaml` or `.diagram-gen.json`) must support all options. |

### 4.9 Backward Compatibility (Breaking Changes OK)

| ID | Requirement |
|----|-------------|
| FR9.1 | Existing simple annotations (`@//diagram:name=Service`) must continue to work. |
| FR9.2 | The default output format may change to use improved defaults (e.g., better colors). |
| FR9.3 | Migration guide must be provided if default behavior changes. |

## 5. Non-Goals

- Visual editor/GUI integration within the CLI
- Export to other formats (PNG, SVG) - focus on draw.io XML only
- Real-time diagram preview
- Cloud provider API integration for automatic discovery
- Animation or interactive features in the output

## 6. Design Considerations

### 6.1 Annotation Syntax Extension

```
@//diagram:name=PaymentService;type=iso:server;page=payment;swimlane=AWS;fillColor=#dae8fc;edgeStyle=elbowEdgeStyle
```

### 6.2 CLI Flags

```
--isometric              Enable isometric mode (shortcut for --layout isometric)
--layout <layout>       Layout: grid, layered, isometric (default: layered)
--shape <shape>         Default shape for components
--compress              Compress output with deflate+base64
--config <path>         Path to config file
--output, -o            Output file path
```

### 6.3 Config File Format (.diagram-gen.yaml)

```yaml
diagram:
  layout: layered
  defaultShape: iso:server
  compress: false
  pages:
    - name: "Network View"
      components:
        - selector: "network"
    - name: "Services"
      components:
        - selector: "service"
styles:
  defaults:
    fillColor: "#dae8fc"
    strokeColor: "#6c8ebf"
    fontSize: 12
```

## 7. Technical Considerations

### 7.1 Architecture

```
internal/
  generator/
    drawio.go          # Main generator (refactor existing)
    shapes.go          # Shape type definitions and builders
    styles.go          # Style string builder
    layout/
      grid.go          # Grid layout engine
      layered.go       # Sugiyama layered layout
      isometric.go     # Isometric projection layout
    swimlane.go        # Swimlane container logic
    compress.go        # Compression utilities
```

### 7.2 Key Data Structures

```go
type ShapeType string

const (
    ShapeRectangle     ShapeType = "rectangle"
    ShapeEllipse       ShapeType = "ellipse"
    ShapeCylinder      ShapeType = "cylinder"
    ShapeIsoCube       ShapeType = "mxgraph.isometric.cube"
    ShapeIsoServer     ShapeType = "mxgraph.isometric.server"
    ShapeIsoDatabase   ShapeType = "mxgraph.isometric.database"
    // ... etc
)

type Style struct {
    Shape           string
    FillColor       string
    StrokeColor     string
    StrokeWidth     int
    Opacity         int
    GradientColor   string
    GradientDir     string
    FontSize        int
    FontFamily      string
    FontColor       string
    Dashed          bool
    Shadow          bool
    // ... etc
}
```

### 7.3 Dependencies

- Use standard library `compress/zlib` for compression
- No external dependencies required for basic draw.io generation

## 8. Success Metrics

- Support all shapes listed in fullDrawioOptions.md (100% coverage of documented features)
- Generate valid draw.io files that open in diagrams.net without errors
- Maintain backward compatibility for existing annotations
- Support three configuration methods (CLI, annotations, config file)

## 9. Open Questions

1. Should we support the `sketch=1` style option for hand-drawn look? (Not in initial phase)
2. How should we handle very large diagrams (100+ nodes)? Consider auto-layout with spacing hints.
3. Should we add a `validate` subcommand to check annotations without generating?
