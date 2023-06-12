package golang

import (
	"fmt"
	"io"
	"os/exec"
)

func runTidy(cwd string, env []string, logger io.Writer) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = cwd
	cmd.Env = env
	cmd.Stdout = logger
	cmd.Stderr = logger
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("executing tidy: %w", err)
	}

	return nil
}
