package archparser_test

import (
	"testing"

	"diagram-gen/internal/archparser"
	"diagram-gen/internal/model"
)

func TestParseAnnotationNewFields(t *testing.T) {
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

func TestParseAnnotationInvalid(t *testing.T) {
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
