package cmd_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"diagram-gen/cmd"
	"diagram-gen/internal/generator"
	"diagram-gen/internal/model"
	"diagram-gen/internal/testutil"
)

func runCmd(t *testing.T, dir string, args ...string) error {
	t.Helper()
	testutil.LockCLI()
	defer testutil.UnlockCLI()

	oldArgs := os.Args
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	defer func() {
		os.Args = oldArgs
		_ = os.Chdir(oldWd)
	}()

	if dir != "" {
		if err := os.Chdir(dir); err != nil {
			t.Fatalf("Chdir failed: %v", err)
		}
	}

	os.Args = append([]string{"diagram-gen"}, args...)
	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("execute failed: %w", err)
	}
	return nil
}

func writeInputFile(t *testing.T, dir, name, contents string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	return path
}

func TestGenerateCommandSuccess(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	input := writeInputFile(t, dir, "input.go", "package main\n\n"+
		"type ServiceA struct {\n"+
		"\tField string `diagram:\"type=service,name=ServiceA,connectsTo=ServiceB,page=Page1\"`\n"+
		"}\n\n"+
		"type ServiceB struct {\n"+
		"\tField string `diagram:\"type=service,name=ServiceB\"`\n"+
		"}\n")
	output := filepath.Join(dir, "out.drawio")

	err := runCmd(t, "", "generate", input, "--output", output, "--layout", "grid", "--page", "Page1", "--compress")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if _, err := os.Stat(output); err != nil {
		t.Fatalf("expected output file: %v", err)
	}
}

func TestGenerateCommandOutputDefault(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	input := writeInputFile(t, dir, "input.go", "package main\n\n"+
		"type ServiceA struct {\n"+
		"\tField string `diagram:\"type=service,name=ServiceA\"`\n"+
		"}\n")

	testutil.LockCLI()
	defer testutil.UnlockCLI()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}

	err = cmd.RunGenerateForTest([]string{input, "--output=", "--layout="})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "diagram.drawio")); err != nil {
		t.Fatalf("expected default output file: %v", err)
	}
}

func TestGenerateCommandIsometric(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	input := writeInputFile(t, dir, "input.go", "package main\n\n"+
		"type ServiceA struct {\n"+
		"\tField string `diagram:\"type=service,name=ServiceA\"`\n"+
		"}\n")
	output := filepath.Join(dir, "iso.drawio")

	err := runCmd(t, "", "generate", input, "--output", output, "--isometric")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestRunGenerateForTestError(t *testing.T) {
	t.Parallel()

	testutil.LockCLI()
	defer testutil.UnlockCLI()

	err := cmd.RunGenerateForTest([]string{})
	if err == nil {
		t.Fatal("expected error from RunGenerateForTest")
	}
}

func TestVersionCommand(t *testing.T) {
	t.Parallel()

	err := runCmd(t, "", "version")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestGenerateCommandInvalidInput(t *testing.T) {
	t.Parallel()

	err := runCmd(t, "", "generate", "missing.go")
	if err == nil {
		t.Fatal("expected error for missing input")
	}
}

func TestGenerateCommandNoAnnotations(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	input := writeInputFile(t, dir, "input.go", `package main

type ServiceA struct {
	Field string `+"`json:\"name\"`"+`
}
`)

	err := runCmd(t, "", "generate", input)
	if err == nil {
		t.Fatal("expected error for missing annotations")
	}
}

func TestGenerateCommandInvalidType(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	input := writeInputFile(t, dir, "input.go", "package main\n\n"+
		"type ServiceA struct {\n"+
		"\tField string `diagram:\"type=invalid,name=ServiceA\"`\n"+
		"}\n")

	err := runCmd(t, "", "generate", input)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestGenerateCommandWriteError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	input := writeInputFile(t, dir, "input.go", "package main\n\n"+
		"type ServiceA struct {\n"+
		"\tField string `diagram:\"type=service,name=ServiceA\"`\n"+
		"}\n")

	err := runCmd(t, "", "generate", input, "--output", dir)
	if err == nil {
		t.Fatal("expected write error")
	}
}

func TestGenerateCommandGeneratorError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	input := writeInputFile(t, dir, "input.go", "package main\n\n"+
		"type ServiceA struct {\n"+
		"\tField string `diagram:\"type=service,name=ServiceA\"`\n"+
		"}\n")
	output := filepath.Join(dir, "out.drawio")

	testutil.LockCLI()
	defer testutil.UnlockCLI()

	cmd.SetGeneratorFactory(func() generator.Formatter {
		return errorGenerator{}
	})
	defer cmd.SetGeneratorFactory(nil)

	oldArgs := os.Args
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	defer func() {
		os.Args = oldArgs
		_ = os.Chdir(oldWd)
	}()

	os.Args = []string{"diagram-gen", "generate", input, "--output", output}
	err = cmd.Execute()
	if err == nil {
		t.Fatal("expected generator error")
	}
}

func TestExecuteError(t *testing.T) {
	t.Parallel()

	err := runCmd(t, "", "generate")
	if err == nil {
		t.Fatal("expected error from execute")
	}
}

type errorGenerator struct{}

func (e errorGenerator) Generate(_ *model.Diagram) ([]byte, error) {
	return nil, errors.New("generate failed")
}

func (e errorGenerator) Format() string {
	return "drawio"
}
