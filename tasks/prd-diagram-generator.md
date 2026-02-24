# PRD: Code-to-Diagram CLI Tool

## Introduction/Overview

This project creates a Go-based Cobra CLI tool that generates software diagrams from code annotations. Developers can annotate their Go code (or other languages) with struct tags, and the CLI will parse these annotations to produce visual diagrams. This follows an "infrastructure-as-code" approach where diagrams are version-controlled alongside the source code.

The tool addresses the problem of keeping documentation and diagrams in sync with code - a common pain point in software development where diagrams quickly become outdated.

## Goals

1. Enable developers to generate diagrams directly from code annotations using a simple CLI command
2. Support multiple diagram types: Architecture, Flowchart, and Network diagrams
3. Output to draw.io (XML) format as the initial target, with extensibility for Visio and isoflow formats
4. Provide a human-readable, declarative syntax for defining diagram elements
5. Integrate seamlessly into CI/CD pipelines for automated diagram generation

## User Stories

1. **As a** backend developer, **I want to** annotate my Go structs with tags that describe component relationships, **so that** I can generate architecture diagrams automatically without manual drawing.

2. **As a** DevOps engineer, **I want to** define my infrastructure in code with diagram annotations, **so that** the diagram stays in sync with the infrastructure code in my repository.

3. **As a** team lead, **I want** my team's diagrams to be version-controlled alongside code, **so that** diagram changes are tracked in git and reviewed via pull requests.

4. **As a** developer, **I want** a single CLI tool that can generate multiple output formats, **so that** I can use the format preferred by different stakeholders.

## Functional Requirements

1. The CLI must accept a Go source file or directory as input via command-line arguments
2. The CLI must parse struct tags in the format `diagram:"type=component,name=ServiceA,connectsTo=ServiceB"` to extract diagram metadata
3. The CLI must support the following struct tag keys:
   - `type`: Component type (service, database, queue, cache, api, user, etc.)
   - `name`: Display name for the component
   - `connectsTo`: Comma-separated list of components this element connects to
   - `description`: Optional description text
   - `direction`: Flow direction for connections (unidirectional, bidirectional)
4. The CLI must generate a draw.io compatible XML file as output
5. The CLI must support the following diagram types via a `--type` flag:
   - `architecture`: Layered system architecture (default)
   - `flowchart`: Process flow diagram
   - `network`: Network topology diagram
6. The CLI must support an output flag (`-o, --output`) to specify the output file path
7. The CLI must validate input and provide clear error messages for malformed annotations
8. The CLI must include a help command describing all available options
9. The CLI must include a `version` command to display the tool version
10. The generated diagram must include legend or key explaining the component shapes/colors

## Non-Goals (Out of Scope)

1. The tool will not parse code from languages other than Go in the initial version
2. The tool will not support real-time diagram editing or interactive preview
3. The tool will not integrate directly with draw.io/Visio APIs - it only generates exportable files
4. The tool will not maintain diagram state between runs - each run generates from scratch
5. The tool will not validate the logical correctness of the architecture (e.g., circular dependencies)

## Design Considerations

- Component shapes should follow standard conventions:
  - Services/APIs: Rounded rectangles
  - Databases: Cylinder shape
  - Queues: Rectangle with lines
  - Cache: Rectangle with lightning bolt
  - Users/External systems: Circle/actor shape
- Use consistent color coding for component types
- Include directional arrows for connections with labels where applicable
- Generate a simple, clean layout - manual repositioning in draw.io may be needed for complex diagrams

## Technical Considerations

- Use the Cobra framework for CLI structure
- Create a separate package for diagram generation logic to allow future format extensibility
- Define clear interfaces for output formatters (draw.io, visio, isoflow)
- Use Go's `go/ast` package for parsing source code and extracting struct tags
- Follow semantic versioning for releases
- Include unit tests for core parsing and generation logic

## Success Metrics

1. CLI successfully generates valid draw.io XML from annotated Go code within 2 weeks
2. Tool can parse and render at least 50 components without performance issues
3. Generated diagrams are importable and editable in draw.io without errors
4. CI/CD integration works in at least one major platform (GitHub Actions, GitLab CI)

## Open Questions

1. Should the tool support reading annotations from a separate YAML/JSON config file in addition to Go struct tags?
2. How should the tool handle layout positioning - auto-layout algorithm or grid-based placement?
3. Should there be a template system for customizing diagram styling?
4. Which output format should be prioritized after draw.io - Visio or isoflow?
