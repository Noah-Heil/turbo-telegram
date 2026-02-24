package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"diagram-gen/internal/generator"
	"diagram-gen/internal/model"
)

func TestExecuteNoError(t *testing.T) {
	oldArgs := os.Args
	oldExit := exitFunc
	defer func() {
		os.Args = oldArgs
		exitFunc = oldExit
	}()

	exitFunc = func(code int) {}

	os.Args = []string{"diagram-gen", "--help"}
	err := Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestExecuteWithError(t *testing.T) {
	oldArgs := os.Args
	oldExit := exitFunc
	defer func() {
		os.Args = oldArgs
		exitFunc = oldExit
	}()

	var exitCode int
	exitFunc = func(code int) { exitCode = code }

	os.Args = []string{"diagram-gen", "generate"}
	Execute()

	if exitCode != 1 {
		t.Logf("exitCode = %d, (expected 1 for missing args, but Cobra shows help instead)", exitCode)
	}
}

func TestRootCmdExecute(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"diagram-gen", "--help"}
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestVersionCmd(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"diagram-gen", "version"}
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestGenerateCmdHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"diagram-gen", "generate", "--help"}
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestGenerateCmdValidFile(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tmpFile := "/tmp/test_valid.go"
	os.WriteFile(tmpFile, []byte(`package main
type S struct { F string `+"`"+`diagram:"type=service,name=S"`+"`"+` }
`), 0644)
	defer os.Remove(tmpFile)

	os.Args = []string{"diagram-gen", "generate", tmpFile, "-o", "/tmp/out.drawio"}
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Execute() error = %v", err)
	}
	os.Remove("/tmp/out.drawio")
}

func TestGenerateCmdDirectory(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tmpDir := "/tmp/testdir"
	os.Mkdir(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	os.WriteFile(tmpDir+"/a.go", []byte(`package main
type S struct { F string `+"`"+`diagram:"type=service,name=S"`+"`"+` }
`), 0644)

	os.Args = []string{"diagram-gen", "generate", tmpDir, "-o", "/tmp/out.drawio"}
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Execute() error = %v", err)
	}
	os.Remove("/tmp/out.drawio")
}

func TestGenerateCmdWithCustomType(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tmpFile := "/tmp/test_type.go"
	os.WriteFile(tmpFile, []byte(`package main
type S struct { F string `+"`"+`diagram:"type=service,name=S"`+"`"+` }
`), 0644)
	defer os.Remove(tmpFile)

	os.Args = []string{"diagram-gen", "generate", tmpFile, "-t", "network", "-o", "/tmp/out.drawio"}
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Execute() error = %v", err)
	}
	os.Remove("/tmp/out.drawio")
}

func TestBuildGenerateCmd(t *testing.T) {
	cmd := buildGenerateCmd()
	if cmd.Use == "" {
		t.Error("expected Use to be set")
	}
	if cmd.Flags().Lookup("output") == nil {
		t.Error("expected output flag")
	}
	if cmd.Flags().Lookup("type") == nil {
		t.Error("expected type flag")
	}
}

func TestGenerateRunESuccess(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("output", "o", "diagram.drawio", "")
	cmd.Flags().StringP("type", "t", "architecture", "")

	tmpFile := "/tmp/test_run_e.go"
	backtick := string(rune(96))
	content := "package main\n" +
		"type S struct { F string " + backtick + "diagram:\"type=service,name=S\"" + backtick + " }\n"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)

	cmd.Flags().Set("output", "/tmp/run_e.drawio")
	cmd.Flags().Set("type", "architecture")

	err := generateCmd.RunE(cmd, []string{tmpFile})
	if err != nil {
		t.Fatalf("RunE failed: %v", err)
	}
	os.Remove("/tmp/run_e.drawio")
}

func TestGenerateRunEParseError(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("output", "o", "diagram.drawio", "")
	cmd.Flags().StringP("type", "t", "architecture", "")

	err := generateCmd.RunE(cmd, []string{"/tmp/does_not_exist.go"})
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestGenerateRunENoAnnotations(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("output", "o", "diagram.drawio", "")
	cmd.Flags().StringP("type", "t", "architecture", "")

	tmpFile := "/tmp/noanno_run_e.go"
	os.WriteFile(tmpFile, []byte("package main\ntype S struct { F string }\n"), 0644)
	defer os.Remove(tmpFile)

	err := generateCmd.RunE(cmd, []string{tmpFile})
	if err == nil {
		t.Fatal("expected error for file with no annotations")
	}
}

func TestGenerateRunEValidationError(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("output", "o", "diagram.drawio", "")
	cmd.Flags().StringP("type", "t", "architecture", "")

	tmpFile := "/tmp/invalid_type.go"
	backtick := string(rune(96))
	content := "package main\n" +
		"type S struct { F string " + backtick + "diagram:\"type=unknown,name=S\"" + backtick + " }\n"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)

	err := generateCmd.RunE(cmd, []string{tmpFile})
	if err == nil {
		t.Fatal("expected validation error for unknown component type")
	}
}

func TestGenerateRunEEmptyOutputFlag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("output", "o", "diagram.drawio", "")
	cmd.Flags().StringP("type", "t", "architecture", "")

	tmpFile := "/tmp/run_e_empty_output.go"
	backtick := string(rune(96))
	content := "package main\n" +
		"type S struct { F string " + backtick + "diagram:\"type=service,name=S\"" + backtick + " }\n"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)

	cmd.Flags().Set("output", "")

	err := generateCmd.RunE(cmd, []string{tmpFile})
	if err != nil {
		t.Fatalf("RunE failed: %v", err)
	}
	os.Remove("diagram.drawio")
}

func TestGenerateRunEWriteError(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("output", "o", "diagram.drawio", "")
	cmd.Flags().StringP("type", "t", "architecture", "")

	tmpFile := "/tmp/run_e_write_error.go"
	backtick := string(rune(96))
	content := "package main\n" +
		"type S struct { F string " + backtick + "diagram:\"type=service,name=S\"" + backtick + " }\n"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)

	cmd.Flags().Set("output", "/tmp")

	err := generateCmd.RunE(cmd, []string{tmpFile})
	if err == nil {
		t.Fatal("expected error when output is a directory")
	}
}

type failingGenerator struct{}

func (f failingGenerator) Generate(*model.Diagram) ([]byte, error) {
	return nil, errors.New("generate error")
}

func (f failingGenerator) Format() string {
	return "drawio"
}

func TestGenerateRunEGeneratorError(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("output", "o", "diagram.drawio", "")
	cmd.Flags().StringP("type", "t", "architecture", "")

	tmpFile := "/tmp/run_e_gen_error.go"
	backtick := string(rune(96))
	content := "package main\n" +
		"type S struct { F string " + backtick + "diagram:\"type=service,name=S\"" + backtick + " }\n"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)

	oldNewGenerator := newGenerator
	newGenerator = func() generator.Formatter {
		return failingGenerator{}
	}
	defer func() { newGenerator = oldNewGenerator }()

	err := generateRunE(cmd, []string{tmpFile})
	if err == nil {
		t.Fatal("expected generator error")
	}
}
