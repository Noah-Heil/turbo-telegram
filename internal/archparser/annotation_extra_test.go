package archparser_test

import (
	"testing"

	"diagram-gen/internal/archparser"
	"diagram-gen/internal/model"
)

func TestParseAnnotationNewFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		tag       string
		wantShape model.ShapeType
		wantPage  string
		wantSwim  string
	}{
		{
			name:      "with shape",
			tag:       `type=service,name=MyService,shape=iso:server`,
			wantShape: "iso:server",
		},
		{
			name:     "with page",
			tag:      `type=service,name=MyService,page=network`,
			wantPage: "network",
		},
		{
			name:     "with swimlane",
			tag:      `type=service,name=MyService,swimlane=AWS`,
			wantSwim: "AWS",
		},
		{
			name:     "with style fillColor",
			tag:      `type=service,name=MyService,fillColor=#dae8fc`,
			wantPage: "",
		},
		{
			name:     "with edgeStyle",
			tag:      `type=service,name=MyService,edgeStyle=elbowEdgeStyle`,
			wantPage: "",
		},
		{
			name:     "with startArrow",
			tag:      `type=service,name=MyService,startArrow=block`,
			wantPage: "",
		},
		{
			name:     "with endArrow",
			tag:      `type=service,name=MyService,endArrow=diamond`,
			wantPage: "",
		},
		{
			name:     "semicolon format",
			tag:      `type=service,name=S`,
			wantPage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ann, err := archparser.ParseAnnotation(tt.tag)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tt.wantShape != "" && ann.Shape != tt.wantShape {
				t.Errorf("shape = %q, want %q", ann.Shape, tt.wantShape)
			}
			if tt.wantPage != "" && ann.Page != tt.wantPage {
				t.Errorf("page = %q, want %q", ann.Page, tt.wantPage)
			}
			if tt.wantSwim != "" && ann.Swimlane != tt.wantSwim {
				t.Errorf("swimlane = %q, want %q", ann.Swimlane, tt.wantSwim)
			}
		})
	}
}

func TestParseAnnotationSemicolonInValue(t *testing.T) {
	t.Parallel()

	ann, err := archparser.ParseAnnotation(`name=Service,description=one;two,fillColor=#fff,strokeColor=#000`)
	if err != nil {
		t.Fatalf("ParseAnnotation error: %v", err)
	}
	if ann.Description != "one;two" {
		t.Fatalf("Description = %q, want one;two", ann.Description)
	}
	if ann.Style != "fillColor=#fff;strokeColor=#000" {
		t.Fatalf("Style = %q, want fillColor=#fff;strokeColor=#000", ann.Style)
	}
}

func TestParseAnnotationEdgeArrows(t *testing.T) {
	t.Parallel()

	ann, err := archparser.ParseAnnotation(`name=Service,edgeStyle=elbowEdgeStyle,startArrow=block,endArrow=classic`)
	if err != nil {
		t.Fatalf("ParseAnnotation error: %v", err)
	}
	if ann.EdgeStyle != "elbowEdgeStyle" {
		t.Fatalf("EdgeStyle = %q, want elbowEdgeStyle", ann.EdgeStyle)
	}
	if ann.StartArrow != "block" {
		t.Fatalf("StartArrow = %q, want block", ann.StartArrow)
	}
	if ann.EndArrow != "classic" {
		t.Fatalf("EndArrow = %q, want classic", ann.EndArrow)
	}
}

func TestParseAnnotationInvalid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		tag     string
		wantErr bool
	}{
		{
			name:    "empty tag",
			tag:     "",
			wantErr: true,
		},
		{
			name:    "missing name",
			tag:     "type=service",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := archparser.ParseAnnotation(tt.tag)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
