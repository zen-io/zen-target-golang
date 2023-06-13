package golang

import (
	"fmt"
	"os/exec"

	zen_targets "github.com/zen-io/zen-core/target"
)

func runTidy(target *zen_targets.Target) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = target.Cwd
	cmd.Env = target.GetEnvironmentVariablesList()
	cmd.Stdout = target
	cmd.Stderr = target
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("executing tidy: %w", err)
	}

	return nil
}
