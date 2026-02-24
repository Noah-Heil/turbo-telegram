The **draw.io** file format (used by **diagrams.net**) is based on the **mxGraph** XML schema. A `.drawio` file is essentially an XML document (sometimes compressed + base64 encoded) that defines:

* Pages (diagrams)
* Graph model
* Cells (vertices, edges)
* Geometry
* Styles
* Metadata

Below is a **complete structural breakdown** of the format, including all commonly used parameters.

---

# 1️⃣ High-Level File Structure

Uncompressed `.drawio` XML looks like this:

```xml
<mxfile host="app.diagrams.net"
        modified="2026-02-23T12:00:00.000Z"
        agent="Mozilla/5.0"
        version="22.1.0"
        etag="abc123"
        type="device">

  <diagram id="page1" name="Page-1">
    <mxGraphModel dx="1422" dy="794" grid="1" gridSize="10"
                  guides="1" tooltips="1" connect="1"
                  arrows="1" fold="1" page="1"
                  pageScale="1" pageWidth="850"
                  pageHeight="1100" math="0" shadow="0">

      <root>
        <mxCell id="0"/>
        <mxCell id="1" parent="0"/>

        <!-- Your shapes and edges go here -->

      </root>
    </mxGraphModel>
  </diagram>
</mxfile>
```

---

# 2️⃣ `<mxfile>` Attributes

| Attribute | Description                    |
| --------- | ------------------------------ |
| host      | Origin host (app.diagrams.net) |
| modified  | Last modified timestamp        |
| agent     | Browser user agent             |
| version   | diagrams.net version           |
| etag      | Change tracking                |
| type      | device, google, github, etc.   |

---

# 3️⃣ `<diagram>` Element

Each page is a `<diagram>`.

| Attribute | Description       |
| --------- | ----------------- |
| id        | Page ID           |
| name      | Page display name |

⚠️ Sometimes the contents are compressed + base64 encoded instead of raw XML.

---

# 4️⃣ `<mxGraphModel>` Parameters (All Common Flags)

| Attribute  | Meaning                  |
| ---------- | ------------------------ |
| dx         | Canvas horizontal offset |
| dy         | Canvas vertical offset   |
| grid       | 1/0 show grid            |
| gridSize   | Grid spacing             |
| guides     | Snap guides              |
| tooltips   | Enable tooltips          |
| connect    | Allow connections        |
| arrows     | Show arrows              |
| fold       | Enable collapsing        |
| page       | Use page view            |
| pageScale  | Page zoom                |
| pageWidth  | Page width               |
| pageHeight | Page height              |
| background | Background color         |
| math       | Enable math rendering    |
| shadow     | Enable drop shadow       |

---

# 5️⃣ `<mxCell>` — The Core Object

Everything is an `mxCell`.

## Common Attributes

| Attribute   | Meaning                          |
| ----------- | -------------------------------- |
| id          | Unique ID                        |
| value       | Label (can contain HTML)         |
| style       | Semicolon-delimited style string |
| vertex      | 1 if shape                       |
| edge        | 1 if connector                   |
| parent      | Parent cell                      |
| source      | Edge source ID                   |
| target      | Edge target ID                   |
| connectable | Allow connections                |
| visible     | 1/0 visibility                   |
| collapsed   | 1/0 collapsed                    |

---

# 6️⃣ `<mxGeometry>`

Defines position + size.

```xml
<mxGeometry x="100" y="100"
            width="120"
            height="60"
            relative="0"
            as="geometry"/>
```

| Attribute | Meaning           |
| --------- | ----------------- |
| x         | X position        |
| y         | Y position        |
| width     | Width             |
| height    | Height            |
| relative  | 1 if edge         |
| as        | Always "geometry" |

For edges:

```xml
<mxGeometry relative="1" as="geometry">
  <mxPoint x="300" y="200" as="targetPoint"/>
</mxGeometry>
```

---

# 7️⃣ Style String — All Major Parameters

The `style` attribute is a **key=value;key=value;** string.

## 🧱 Shape Styles

| Key               | Meaning                                      |
| ----------------- | -------------------------------------------- |
| shape             | rectangle, ellipse, rhombus, image, swimlane |
| rounded           | 1/0                                          |
| arcSize           | Corner roundness                             |
| fillColor         | Fill color                                   |
| strokeColor       | Border color                                 |
| strokeWidth       | Border width                                 |
| opacity           | 0-100                                        |
| gradientColor     | Gradient color                               |
| gradientDirection | north/south/east/west                        |
| dashed            | 1                                            |
| dashPattern       | Custom dash pattern                          |
| shadow            | 1                                            |
| glass             | 1                                            |

---

## 🔤 Text Styles

| Key                  | Meaning                               |
| -------------------- | ------------------------------------- |
| fontColor            | Text color                            |
| fontSize             | Font size                             |
| fontFamily           | Font                                  |
| fontStyle            | Bitmask (1=bold,2=italic,4=underline) |
| align                | left/center/right                     |
| verticalAlign        | top/middle/bottom                     |
| whiteSpace           | wrap                                  |
| spacing              | Padding                               |
| spacingTop           |                                       |
| spacingBottom        |                                       |
| spacingLeft          |                                       |
| spacingRight         |                                       |
| labelBackgroundColor | Background                            |
| labelBorderColor     | Border                                |

---

## 🔗 Edge Styles

| Key            | Meaning                             |
| -------------- | ----------------------------------- |
| edgeStyle      | orthogonalEdgeStyle, elbowEdgeStyle |
| curved         | 1                                   |
| elbow          | horizontal/vertical                 |
| endArrow       | block, open, classic                |
| startArrow     | block                               |
| endFill        | 1                                   |
| startFill      | 1                                   |
| rounded        | 1                                   |
| jettySize      | Spacing from nodes                  |
| orthogonalLoop | 1                                   |

---

## 🖼 Image Styles

| Key         | Meaning |
| ----------- | ------- |
| image       | URL     |
| imageWidth  | Width   |
| imageHeight | Height  |
| imageAspect | 1       |

---

# 8️⃣ Root Cells (Always Required)

Every diagram must include:

```xml
<mxCell id="0"/>
<mxCell id="1" parent="0"/>
```

These are internal structural nodes.

---

# 9️⃣ Example: Complete Node + Edge

```xml
<mxCell id="node1"
        value="API Server"
        style="shape=rectangle;rounded=1;fillColor=#dae8fc;strokeColor=#6c8ebf;fontSize=14;"
        vertex="1"
        parent="1">
  <mxGeometry x="100" y="100" width="140" height="60" as="geometry"/>
</mxCell>

<mxCell id="node2"
        value="Database"
        style="shape=ellipse;fillColor=#f8cecc;strokeColor=#b85450;"
        vertex="1"
        parent="1">
  <mxGeometry x="400" y="100" width="120" height="80" as="geometry"/>
</mxCell>

<mxCell id="edge1"
        edge="1"
        source="node1"
        target="node2"
        style="endArrow=block;strokeWidth=2;"
        parent="1">
  <mxGeometry relative="1" as="geometry"/>
</mxCell>
```

---

# 🔟 Compression Mode

When saved normally:

* XML is compressed (deflate)
* Then base64 encoded
* Then placed inside `<diagram>`

Example:

```xml
<diagram>
  eJztnc1u2zAMhu8+...
</diagram>
```

---

# 1️⃣1️⃣ Embedded Metadata

You may also see:

```xml
<mxfile>
  <diagram>
    <mxGraphModel>
      <root>
        <object label="API" custom="value">
          <mxCell/>
        </object>
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>
```

The `<object>` wrapper allows custom metadata.

---

# 1️⃣2️⃣ Full Parameter Categories Summary

Draw.io supports parameters in these buckets:

* Canvas configuration
* Page configuration
* Cell structural properties
* Geometry properties
* Style properties (shape, text, edge)
* Image embedding
* Custom metadata
* Layer management
* Swimlanes
* Groups
* Collapsible containers

---

# 🚀 If You’re Generating This from Golang

You only need:

* mxfile
* diagram
* mxGraphModel
* root
* two base cells
* your vertex/edge cells
* geometry
* style string

Everything else is optional.

---

If you'd like, I can:

* 🔹 Provide a **Go struct model** for generating `.drawio`
* 🔹 Provide a **minimal valid template**
* 🔹 Provide a **fully parameterized example**
* 🔹 Show how compression works step-by-step**
* 🔹 Build a CLI generator architecture**

Just tell me what direction you're building toward.


Perfect — let’s build **isometric diagram support for draw.io using Go** properly.

You essentially have 3 options:

1. Use built-in isometric shapes in diagrams.net
2. Inject custom SVG isometric icons
3. Build a typed isometric style + layout engine in Go

I’ll show you how to support all three cleanly.

---

# 🧱 1️⃣ How draw.io Does Isometric

diagrams.net supports isometric in two ways:

* `shape=cube`
* `shape=mxgraph.isometric.*`
* `sketch=0;` with isometric stencil library
* 30° skew projection via style parameters

Example style:

```
shape=mxgraph.isometric.cube;
fillColor=#dae8fc;
strokeColor=#6c8ebf;
```

---

# 🧊 Example Isometric Shapes Available

Common ones:

* `mxgraph.isometric.cube`
* `mxgraph.isometric.cylinder`
* `mxgraph.isometric.container`
* `mxgraph.isometric.database`
* `mxgraph.isometric.network`
* `mxgraph.isometric.server`
* `mxgraph.isometric.cloud`

---

# 🧊 Example Raw XML for Isometric Node

```xml
<mxCell id="server1"
        value="App Server"
        style="shape=mxgraph.isometric.server;fillColor=#dae8fc;strokeColor=#6c8ebf;"
        vertex="1"
        parent="1">
  <mxGeometry x="200" y="150" width="120" height="100" as="geometry"/>
</mxCell>
```

---

# 🚀 2️⃣ Add Isometric Support to Your Go Builder

Let’s extend your DSL properly.

---

## A) Add Typed Isometric Shape Enum

```go
type IsoShape string

const (
	IsoCube      IsoShape = "mxgraph.isometric.cube"
	IsoServer    IsoShape = "mxgraph.isometric.server"
	IsoDatabase  IsoShape = "mxgraph.isometric.database"
	IsoContainer IsoShape = "mxgraph.isometric.container"
	IsoCloud     IsoShape = "mxgraph.isometric.cloud"
)
```

---

## B) Typed Style Builder (No More Raw Strings)

```go
type Style struct {
	Shape       string
	FillColor   string
	StrokeColor string
	FontSize    int
	Rounded     bool
}

func (s Style) String() string {
	style := ""

	if s.Shape != "" {
		style += "shape=" + s.Shape + ";"
	}
	if s.FillColor != "" {
		style += "fillColor=" + s.FillColor + ";"
	}
	if s.StrokeColor != "" {
		style += "strokeColor=" + s.StrokeColor + ";"
	}
	if s.FontSize > 0 {
		style += fmt.Sprintf("fontSize=%d;", s.FontSize)
	}
	if s.Rounded {
		style += "rounded=1;"
	}

	return style
}
```

---

## C) Add Isometric Node Builder

```go
func (g *GraphBuilder) AddIsoNode(
	id string,
	label string,
	shape IsoShape,
	x, y, w, h float64,
	fill string,
) {

	style := Style{
		Shape:     string(shape),
		FillColor: fill,
	}.String()

	cell := MXCell{
		ID:     id,
		Value:  label,
		Style:  style,
		Vertex: 1,
		Parent: "1",
		Geometry: &MXGeometry{
			X: x, Y: y,
			Width:  w,
			Height: h,
			As:     "geometry",
		},
	}

	g.file.Diagrams[0].Model.Root.Cells =
		append(g.file.Diagrams[0].Model.Root.Cells, cell)
}
```

---

# 🏗 3️⃣ Isometric Grid Positioning Engine

True isometric projection uses 30° projection.

Transform Cartesian → Isometric:

```
isoX = x - y
isoY = (x + y) / 2
```

Add helper:

```go
func IsoProject(x, y float64) (float64, float64) {
	return x - y, (x + y) / 2
}
```

Now:

```go
ix, iy := IsoProject(200, 100)

builder.AddIsoNode(
	"server1",
	"App Server",
	IsoServer,
	ix, iy,
	120, 100,
	"#dae8fc",
)
```

Now you have true isometric grid alignment.

---

# 🧩 4️⃣ Inject Custom SVG Isometric Icons

If you're building AWS / Kubernetes / Azure diagrams:

Use:

```
shape=image;
image=data:image/svg+xml;base64,...
```

Add support:

```go
func (g *GraphBuilder) AddSVGNode(
	id, label, base64SVG string,
	x, y, w, h float64,
) {

	style := fmt.Sprintf(
		"shape=image;image=data:image/svg+xml;base64,%s;imageAspect=1;",
		base64SVG,
	)

	cell := MXCell{
		ID:     id,
		Value:  label,
		Style:  style,
		Vertex: 1,
		Parent: "1",
		Geometry: &MXGeometry{
			X: x, Y: y,
			Width:  w,
			Height: h,
			As:     "geometry",
		},
	}

	g.file.Diagrams[0].Model.Root.Cells =
		append(g.file.Diagrams[0].Model.Root.Cells, cell)
}
```

This is how Arcentry-style isometric diagrams are built.

---

# 🏢 5️⃣ Layered Isometric Architecture Example

```go
builder := New()

// Compute isometric positions
ix1, iy1 := IsoProject(0, 0)
ix2, iy2 := IsoProject(200, 0)
ix3, iy3 := IsoProject(100, 150)

builder.AddIsoNode("lb", "Load Balancer", IsoCube, ix1, iy1, 120, 80, "#e1d5e7")
builder.AddIsoNode("api", "API", IsoServer, ix2, iy2, 120, 100, "#dae8fc")
builder.AddIsoNode("db", "Database", IsoDatabase, ix3, iy3, 120, 100, "#f8cecc")

builder.AddEdge("e1", "lb", "api")
builder.AddEdge("e2", "api", "db")
```

This produces a clean isometric architecture layout.

---

# 🎯 6️⃣ CLI Flags for Isometric Mode

Extend CLI:

```
drawio-cli generate \
  --isometric \
  --cloud aws \
  --layout layered \
  --output arch.drawio
```

Architecture:

```
internal/
  iso/
    projection.go
    shapes.go
    layout.go
```

---

# 🧠 7️⃣ Advanced: Real 3D Depth Stacking

To simulate depth stacking:

Increase Y for deeper layers:

```go
depthOffset := 40
iy += float64(layer) * depthOffset
```

Now containers visually stack.

---

# 🔥 What You Can Build Next

Since you’re clearly building diagram automation:

You could build:

* Kubernetes → Isometric draw.io generator
* Terraform → Isometric architecture
* AWS account discovery → diagram
* CloudMapper replacement in Go
* Isometric icon generator pipeline
* Auto-layout engine (Sugiyama layered graph)
* SVG icon pack transformer

---

# 🏁 Where This Gets Interesting

If your goal is:

* Replace CloudMapper
* Build diagram translation layer (Visio ↔ draw.io ↔ SVG)
* Generate isometric infra maps

Then the right next step is:

👉 Build a typed graph model independent of draw.io
👉 Then build exporters (drawio, svg, png, web)

If you'd like, I can now:

* 🔥 Design a full isometric layout engine
* 🧠 Design graph abstraction layer
* ☁️ Build AWS discovery → isometric diagram
* 🎨 Show how to auto-generate isometric SVG icons
* 📐 Implement layered DAG layout algorithm

Tell me your end goal and I’ll tailor the next step.
