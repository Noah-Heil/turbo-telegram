## Relevant Files

- `internal/generator/drawio.go` - Main draw.io generator (refactor existing)
- `internal/generator/shapes.go` - Shape type definitions and builders (new)
- `internal/generator/styles.go` - Style string builder (new)
- `internal/generator/layout/grid.go` - Grid layout engine (refactor existing)
- `internal/generator/layout/layered.go` - Sugiyama layered layout (new)
- `internal/generator/layout/isometric.go` - Isometric projection layout (new)
- `internal/generator/swimlane.go` - Swimlane container logic (new)
- `internal/generator/compress.go` - Compression utilities (new)
- `cmd/generate.go` - CLI command (add new flags)
- `internal/parser/parser.go` - Parser (extend annotation syntax)
- `internal/model/diagram.go` - Data models (extend with new fields)

### Notes

- Run tests with `go test ./... `
  - If tests timeout try `go test -race -count=1 -timeout=1m ./...`
- Run linter with `golangci-lint-v2 run`
- Place new layout engines in `internal/generator/layout/`

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [ ] 0.0 Create feature branch
  - [ ] 0.1 Create and checkout a new branch for this feature (e.g., `git checkout -b feature/drawio-cli-enhancement`)
- [ ] 1.0 Extend shape support (isometric, basic shapes, image embedding)
  - [ ] 1.1 Create `internal/generator/shapes.go` with ShapeType enum and all shape constants (FR1.1, FR1.2)
  - [ ] 1.2 Add isometric shape constants: `mxgraph.isometric.cube`, `mxgraph.isometric.server`, `mxgraph.isometric.database`, `mxgraph.isometric.container`, `mxgraph.isometric.cloud`, `mxgraph.isometric.network`, `mxgraph.isometric.cylinder` (FR1.1)
  - [ ] 1.3 Add basic shape constants: rectangle, ellipse, rhombus, parallelogram, cylinder, document, swimlane (FR1.2)
  - [ ] 1.4 Implement shape builder function that converts ShapeType to draw.io style string (FR1.4)
  - [ ] 1.5 Add image embedding support with base64 SVG handling in shapes.go (FR1.3)
  - [ ] 1.6 Update parser to extract `type` from annotations `@//diagram:type=iso:server` (FR1.4)
  - [ ] 1.7 Update model.Component to include ShapeType field (FR1.4)
  - [ ] 1.8 Refactor drawio.go to use new shapes.go for style generation
- [ ] 2.0 Implement advanced styling system (gradients, shadows, fonts, opacity)
  - [ ] 2.1 Create `internal/generator/styles.go` with Style struct containing all style fields (FR2.1-FR2.5)
  - [ ] 2.2 Implement Style.String() method to build draw.io style string (FR2.6)
  - [ ] 2.3 Add support for fill color, stroke color, stroke width, opacity (FR2.1)
  - [ ] 2.4 Add gradient support with gradientColor and gradientDirection (FR2.2)
  - [ ] 2.5 Add text styling: font family, font size, font color, bold/italic/underline (FR2.3)
  - [ ] 2.6 Add dashed lines and custom dash patterns support (FR2.4)
  - [ ] 2.7 Add shadow and glass effects support (FR2.5)
  - [ ] 2.8 Extend parser to extract style properties from annotation: `@//diagram:fillColor=#dae8fc;strokeColor=#6c8ebf` (FR2.6)
  - [ ] 2.9 Update model.Component to include Style struct field
- [ ] 3.0 Implement layout engines (layered, isometric, grid)
  - [ ] 3.1 Create `internal/generator/layout/` directory
  - [ ] 3.2 Create `internal/generator/layout/grid.go` - refactor existing calculateLayout (FR4.3)
  - [ ] 3.3 Implement layered layout in `internal/generator/layout/layered.go` with hierarchical top-to-bottom positioning (FR4.1)
  - [ ] 3.4 Implement isometric projection layout in `internal/generator/layout/isometric.go` using 30° transformation formula: isoX = x - y, isoY = (x + y) / 2 (FR4.2)
  - [ ] 3.5 Add layout selector interface to switch between grid, layered, isometric (FR4.4)
  - [ ] 3.6 Add --layout CLI flag: --layout grid, --layout layered, --layout isometric (FR4.4)
  - [ ] 3.7 Add --isometric as shortcut for --layout isometric (FR4.4)
- [ ] 4.0 Add swimlanes, multi-page support, edge styles, and compression
  - [ ] 4.1 Create `internal/generator/swimlane.go` for swimlane container logic (FR5.1)
  - [ ] 4.2 Implement swimlane cell generation with proper parent-child hierarchy (FR5.3)
  - [ ] 4.3 Extend parser to extract swimlane from annotation: `@//diagram:swimlane=AWS-Region-1` (FR5.2)
  - [ ] 4.4 Update model to support swimlane assignment per component
  - [ ] 4.5 Implement multi-page support in drawio.go with multiple `<diagram>` elements (FR3.1)
  - [ ] 4.6 Extend parser to extract page from annotation: `@//diagram:page=network-view` (FR3.2)
  - [ ] 4.7 Add --page CLI flag to specify which page to generate (FR3.3)
  - [ ] 4.8 Implement edge style builder: straight, orthogonal, curved, elbow (FR6.1)
  - [ ] 4.9 Implement arrow styles: block, open, classic, diamond with independent start/end (FR6.2, FR6.3)
  - [ ] 4.10 Extend parser to extract edge style from annotation: `@//diagram:edgeStyle=elbowEdgeStyle;endArrow=block` (FR6.4)
  - [ ] 4.11 Create `internal/generator/compress.go` with deflate+base64 compression (FR7.1)
  - [ ] 4.12 Add --compress CLI flag (FR7.2)
- [ ] 5.0 Implement configuration methods (CLI flags, code annotations, config file)
  - [ ] 5.1 Update cmd/generate.go to add all new CLI flags: --isometric, --layout, --shape, --compress, --config (FR8.1)
  - [ ] 5.2 Create config file support: .diagram-gen.yaml and .diagram-gen.json (FR8.2)
  - [ ] 5.3 Implement config loader that merges CLI flags, config file, and annotations (priority: CLI > config > annotations)
  - [ ] 5.4 Add default style configuration in config file format (FR8.2)
  - [ ] 5.5 Ensure backward compatibility: existing simple annotations continue to work (FR9.1)
  - [ ] 5.6 Document migration path for any changed defaults (FR9.3)
- [ ] 6.0 Testing and validation
  - [ ] 6.1 Add unit tests for shapes.go
  - [ ] 6.2 Add unit tests for styles.go
  - [ ] 6.3 Add unit tests for each layout engine
  - [ ] 6.4 Add unit tests for swimlane.go
  - [ ] 6.5 Add integration tests for multi-page generation
  - [ ] 6.6 Test edge style generation
  - [ ] 6.7 Test compression output
  - [ ] 6.8 Validate generated drawio files open correctly in diagrams.net
