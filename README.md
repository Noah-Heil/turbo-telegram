# diagram-gen

A Go-based Cobra CLI tool that generates software diagrams from code annotations. Parse Go source files with struct tags to automatically generate draw.io compatible diagrams.

## Features

- Generate diagrams from Go struct tags
- Support for multiple component types (services, databases, queues, caches, etc.)
- Automatic layout calculation (grid, layered, isometric)
- Connection arrows between components
- draw.io XML output with optional compression
- Isometric shapes (cube, server, database, container, cloud)
- Advanced styling (gradients, shadows, fonts, opacity)
- Swimlane containers for grouping components
- Multiple diagram pages
- Edge styles (straight, orthogonal, curved, elbow)

## Installation

### From Source

```bash
git clone https://github.com/yourusername/diagram-gen.git
cd diagram-gen
go build -o bin/diagram-gen
```

### Using Go Install

```bash
go install github.com/yourusername/diagram-gen@latest
```

## Usage

```bash
# Generate diagram from a file
diagram-gen generate input.go -o diagram.drawio

# Generate diagram from a directory
diagram-gen generate ./internal/services/ -o architecture.drawio

# Specify diagram type
diagram-gen generate input.go -t architecture -o diagram.drawio

# Use isometric layout
diagram-gen generate input.go --layout isometric -o diagram.drawio

# Compress output
diagram-gen generate input.go --compress -o diagram.drawio
```

## CLI Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `diagram.drawio` | Output file path |
| `--type` | `-t` | `architecture` | Diagram type (architecture, flowchart, network) |
| `--layout` | | `layered` | Layout engine (grid, layered, isometric) |
| `--isometric` | | false | Shortcut for --layout isometric |
| `--shape` | | | Default shape for components |
| `--compress` | | false | Compress output with deflate+base64 |
| `--config` | | | Path to config file |
| `--page` | | | Generate specific page |

## Annotation Syntax

Add `diagram` struct tags to your Go code:

```go
package main

// API Gateway component
type APIGateway struct {
    Field string `diagram:"type=gateway,name=APIGateway,connectsTo=AuthService;UserService"`
}

// Authentication service
type AuthService struct {
    Field string `diagram:"type=service,name=AuthService,description=OAuth2 provider"`
}

// User service
type UserService struct {
    Field string `diagram:"type=service,name=UserService,connectsTo=UserDB;Cache"`
}

// User database
type UserDatabase struct {
    Field string `diagram:"type=database,name=UserDB"`
}

// Redis cache
type RedisCache struct {
    Field string `diagram:"type=cache,name=Cache"`
}
```

### Extended Tag Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Display name for the component |
| `type` | No | Component type (defaults to `service`) |
| `connectsTo` | No | Semicolon-separated list of target components |
| `description` | No | Optional description text |
| `direction` | No | Flow direction (`unidirectional` or `bidirectional`) |
| `page` | No | Page name for multi-page diagrams |
| `swimlane` | No | Swimlane container name |
| `shape` | No | Shape type (rectangle, ellipse, iso:server, iso:database, etc.) |
| `fillColor` | No | Fill color (hex format) |
| `strokeColor` | No | Stroke color (hex format) |
| `gradientColor` | No | Gradient end color |
| `gradientDirection` | No | Gradient direction (north, south, east, west) |
| `fontSize` | No | Font size |
| `fontFamily` | No | Font family |
| `edgeStyle` | No | Edge style (straightEdgeStyle, orthogonalEdgeEdgeStyle, curvedStyle, elbowEdgeStyle) |
| `endArrow` | No | End arrow style (block, open, classic, diamond) |

### Component Types

| Type | Shape | Color |
|------|-------|-------|
| `service` | Rounded rectangle | Blue |
| `api` | Rounded rectangle | Green |
| `gateway` | Rounded rectangle | Green |
| `database` | Cylinder | Orange |
| `queue` | Parallelogram | Yellow |
| `cache` | Dashed rectangle | Red |
| `user` | Ellipse | Purple |
| `external` | Document | Gray |
| `storage` | Cylinder | Yellow |

### Isometric Shapes

| Type | Shape |
|------|-------|
| `iso:cube` | Isometric cube |
| `iso:server` | Isometric server |
| `iso:database` | Isometric database |
| `iso:container` | Isometric container |
| `iso:cloud` | Isometric cloud |
| `iso:network` | Isometric network |
| `iso:cylinder` | Isometric cylinder |

Example:
```go
type CloudService struct {
    Field string `diagram:"type=service,shape=iso:cloud,name=CloudService"`
}
```

### Swimlanes

Group components into swimlanes:

```go
type ServiceA struct {
    Field string `diagram:"name=ServiceA,swimlane=AWS"`
}

type ServiceB struct {
    Field string `diagram:"name=ServiceB,swimlane=AWS"`
}

type DatabaseA struct {
    Field string `diagram:"name=DatabaseA,swimlane=On-Premise"`
}
```

### Multi-Page Diagrams

Organize components into pages:

```go
type NetworkComponent struct {
    Field string `diagram:"name=LoadBalancer,page=Network View"`
}

type ServiceComponent struct {
    Field string `diagram:"name=API,page=Services"`
}
```

### Styling

Apply custom styling:

```go
type MyService struct {
    Field string `diagram:"name=MyService,fillColor=#dae8fc,strokeColor=#6c8ebf,fontSize=14"`
}
```

## Example Output

Input file:
```go
package main

type UserService struct {
    Field string `diagram:"type=service,name=UserService,connectsTo=UserDB"`
}

type OrderService struct {
    Field string `diagram:"type=service,name=OrderService,connectsTo=OrderDB;PaymentGateway"`
}

type UserDB struct {
    Field string `diagram:"type=database,name=UserDB"`
}

type OrderDB struct {
    Field string `diagram:"type=database,name=OrderDB"`
}

type PaymentGateway struct {
    Field string `diagram:"type=external,name=PaymentGateway,description=Stripe"`
}
```

Generate diagram:
```bash
diagram-gen generate example.go -o diagram.drawio
```

The output will be a draw.io XML file with:
- UserService → UserDB arrow
- OrderService → OrderDB arrow
- OrderService → PaymentGateway arrow

## CI/CD Integration

### GitHub Actions

```yaml
name: Generate Architecture Diagram

on:
  push:
    paths:
      - '**.go'

jobs:
  diagram:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Build diagram-gen
        run: go build -o bin/diagram-gen
      
      - name: Generate diagram
        run: ./bin/diagram-gen generate ./internal/ -o docs/architecture.drawio
      
      - name: Commit and push
        run: |
          git config --local user.email "github-actions[bot]@users.noreply.github.com"
          git config --local user.name "github-actions[bot]"
          git add docs/architecture.drawio
          git commit -m "Update architecture diagram" || echo "No changes to commit"
          git push
```

## Development

### Run Tests

```bash
go test ./...
```

### Run Linter

```bash
golangci-lint-v2 run
```

### Build

```bash
go build -o bin/diagram-gen
```

## License

MIT
