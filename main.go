package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("Project directory is required")
	}

	args := []string{"build", ".", "errors"}
	cmd := exec.Command("go", args...)
	cmd.Dir = proj

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("'go build' failed: %s", err)
	}

	_, err := fmt.Fprintln(out, "Go Build: SUCCESS")

	return err

}