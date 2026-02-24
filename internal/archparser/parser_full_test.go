package archparser_test

import (
	"diagram-gen/internal/archparser"
	"os"
	"path/filepath"
	"testing"
)

func TestParseDirectory(t *testing.T) {
	p := archparser.New()

	dir := t.TempDir()

	file1 := filepath.Join(dir, "service.go")
	err := os.WriteFile(file1, []byte(`
package main

type UserService struct {
    Field string `+"`"+`diagram:"type=service,name=UserService,connectsTo=Database"`+"`"+`
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	file2 := filepath.Join(dir, "database.go")
	err = os.WriteFile(file2, []byte(`
package main

type Database struct {
    Field string `+"`"+`diagram:"type=database,name=Database"`+"`"+`
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	diagram, err := p.ParseDirectory(dir)
	if err != nil {
		t.Fatalf("ParseDirectory failed: %v", err)
	}

	if len(diagram.Components) != 2 {
		t.Errorf("expected 2 components, got %d", len(diagram.Components))
	}
}

func TestParse(t *testing.T) {
	p := archparser.New()

	diagram, err := p.Parse("testdata/sample.go")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(diagram.Components) == 0 {
		t.Error("expected components")
	}
}

func TestParseDirectoryEmpty(t *testing.T) {
	p := archparser.New()

	dir := t.TempDir()

	_, err := p.ParseDirectory(dir)
	if err != nil {
		t.Fatalf("ParseDirectory failed: %v", err)
	}
}

func TestParseNonExistent(t *testing.T) {
	p := archparser.New()

	_, err := p.Parse("nonexistent.go")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestParseStructTag(t *testing.T) {
	tests := []struct {
		name     string
		tagValue string
		key      string
		want     string
	}{
		{
			name:     "basic tag",
			tagValue: "`diagram:\"type=service,name=Test\"`",
			key:      "diagram",
			want:     "type=service,name=Test",
		},
		{
			name:     "no match",
			tagValue: "`json:\"name\"`",
			key:      "diagram",
			want:     "",
		},
		{
			name:     "only backticks",
			tagValue: "`diagram:\"test\"`",
			key:      "diagram",
			want:     "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := archparser.ParseStructTag(tt.tagValue, tt.key)
			if got != tt.want {
				t.Errorf("ParseStructTag() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseFileWithInvalidGo(t *testing.T) {
	p := archparser.New()

	tmpFile := "/tmp/invalid.go"
	err := os.WriteFile(tmpFile, []byte(`package main
invalid code here
`), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	_, err = p.ParseFile(tmpFile)
	if err == nil {
		t.Error("expected error for invalid Go code")
	}
}

func TestParseDirectoryWithInvalidFile(t *testing.T) {
	p := archparser.New()

	dir := t.TempDir()

	file1 := dir + "/invalid.go"
	err := os.WriteFile(file1, []byte(`invalid go code`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	diagram, err := p.ParseDirectory(dir)
	if err != nil {
		t.Fatalf("ParseDirectory should not fail for invalid file: %v", err)
	}

	if len(diagram.Components) != 0 {
		t.Errorf("expected 0 components, got %d", len(diagram.Components))
	}
}

func TestParseFileWithInvalidAnnotation(t *testing.T) {
	p := archparser.New()

	tmpFile := "/tmp/invalid_anno.go"
	err := os.WriteFile(tmpFile, []byte(`package main
type S struct { 
	F string `+"`"+`diagram:"invalid"`+"`"+` 
}
`), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	diagram, err := p.ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile should handle invalid annotation: %v", err)
	}

	if len(diagram.Components) != 0 {
		t.Errorf("expected 0 components (invalid annotation), got %d", len(diagram.Components))
	}
}

func TestParseFileWithNonStruct(t *testing.T) {
	p := archparser.New()

	tmpFile := "/tmp/interface.go"
	err := os.WriteFile(tmpFile, []byte(`package main
type MyInterface interface {
	Method() error
}
`), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	diagram, err := p.ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(diagram.Components) != 0 {
		t.Errorf("expected 0 components for interface, got %d", len(diagram.Components))
	}
}

func TestParseFileWithNoFields(t *testing.T) {
	p := archparser.New()

	tmpFile := "/tmp/empty.go"
	err := os.WriteFile(tmpFile, []byte(`package main
type Empty struct {}
`), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	diagram, err := p.ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(diagram.Components) != 0 {
		t.Errorf("expected 0 components for empty struct, got %d", len(diagram.Components))
	}
}

func TestParseDirectoryWithOnlySubdirs(t *testing.T) {
	p := archparser.New()

	dir := t.TempDir()
	subdir := dir + "/subdir"
	err := os.Mkdir(subdir, 0755)
	if err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	diagram, err := p.ParseDirectory(dir)
	if err != nil {
		t.Fatalf("ParseDirectory failed: %v", err)
	}

	if len(diagram.Components) != 0 {
		t.Errorf("expected 0 components for dir with no go files, got %d", len(diagram.Components))
	}
}

func TestParseDirectoryReadError(t *testing.T) {
	p := archparser.New()

	dir := "/nonexistent_dir_12345"
	_, err := p.ParseDirectory(dir)
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestParseWithFile(t *testing.T) {
	p := archparser.New()

	tmpFile := "/tmp/test_parse_file.go"
	err := os.WriteFile(tmpFile, []byte(`package main
type S struct { F string `+"`"+`diagram:"type=service,name=S"`+"`"+` }
`), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	diagram, err := p.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(diagram.Components) != 1 {
		t.Errorf("expected 1 component, got %d", len(diagram.Components))
	}
}

func TestParseWithDirectory(t *testing.T) {
	p := archparser.New()

	dir := t.TempDir()
	err := os.WriteFile(dir+"/a.go", []byte(`package main
type S struct { F string `+"`"+`diagram:"type=service,name=S"`+"`"+` }
`), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	diagram, err := p.Parse(dir)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(diagram.Components) != 1 {
		t.Errorf("expected 1 component, got %d", len(diagram.Components))
	}
}

func TestParseFileWithMultipleFields(t *testing.T) {
	p := archparser.New()

	tmpFile := "/tmp/multi_field.go"
	err := os.WriteFile(tmpFile, []byte(`package main
type UserService struct {
	Field1 string `+"`"+`diagram:"type=service,name=UserService"`+"`"+`
	Field2 string `+"`"+`diagram:"type=database,name=UserDB"`+"`"+`
}
`), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	diagram, err := p.ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(diagram.Components) != 2 {
		t.Errorf("expected 2 components, got %d", len(diagram.Components))
	}
}

func TestParseFileWithFieldWithoutTag(t *testing.T) {
	p := archparser.New()

	tmpFile := "/tmp/mixed.go"
	err := os.WriteFile(tmpFile, []byte(`package main
type S struct {
	Field1 string `+"`"+`diagram:"type=service,name=S"`+"`"+`
	Field2 string
}
`), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	diagram, err := p.ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(diagram.Components) != 1 {
		t.Errorf("expected 1 component (only tagged field), got %d", len(diagram.Components))
	}
}

func TestParseFileWithEmptyDiagramTag(t *testing.T) {
	p := archparser.New()

	tmpFile := "/tmp/empty_tag.go"
	backtick := string(rune(96))
	content := "package main\n" +
		"type S struct {\n" +
		"\tField1 string " + backtick + "diagram:\"\"" + backtick + "\n" +
		"}\n"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	diagram, err := p.ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	if len(diagram.Components) != 0 {
		t.Errorf("expected 0 components for empty tag, got %d", len(diagram.Components))
	}
}

func TestParseFileWithNonDiagramTag(t *testing.T) {
	p := archparser.New()

	tmpFile := "/tmp/non_diagram_tag.go"
	backtick := string(rune(96))
	content := "package main\n" +
		"type S struct {\n" +
		"\tField1 string " + backtick + "json:\"name\"" + backtick + "\n" +
		"}\n"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	diagram, err := p.ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	if len(diagram.Components) != 0 {
		t.Errorf("expected 0 components for non-diagram tag, got %d", len(diagram.Components))
	}
}
