package main

import (
	"os"
	"testing"

	"diagram-gen/cmd"
)

func TestMainExit(t *testing.T) {
	t.Parallel()
	oldArgs := os.Args
	oldExit := exitFunc
	defer func() {
		os.Args = oldArgs
		exitFunc = oldExit
	}()

	var exitCode int
	exitFunc = func(code int) { exitCode = code }

	os.Args = []string{"diagram-gen", "generate"}
	main()
	if exitCode != 1 {
		t.Errorf("exitCode = %d, want 1", exitCode)
	}
}

func TestMainSuccess(t *testing.T) {
	t.Parallel()
	oldArgs := os.Args
	oldExit := exitFunc
	defer func() {
		os.Args = oldArgs
		exitFunc = oldExit
	}()

	exitFunc = func(_ int) {}

	os.Args = []string{"diagram-gen", "--help"}
	main()
}

func TestCmdExecute(t *testing.T) {
	t.Parallel()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"diagram-gen", "--help"}
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}
