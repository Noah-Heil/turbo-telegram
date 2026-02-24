package main

import (
	"diagram-gen/cmd"
	"os"
)

var exitFunc = os.Exit

func main() {
	if err := cmd.Execute(); err != nil {
		exitFunc(1)
	}
}
