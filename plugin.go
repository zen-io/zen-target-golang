package golang

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	zen_targets "github.com/zen-io/zen-core/target"
)

type GolangPluginConfig struct {
	zen_targets.BaseFields `mapstructure:",squash"`
	Srcs                   []string `mapstructure:"srcs"`
	Out                    string   `mapstructure:"out"`
	Flags                  string   `mapstructure:"flags"`
}

func (gp GolangPluginConfig) GetTargets(_ *zen_targets.TargetConfigContext) ([]*zen_targets.Target, error) {
	if gp.Out == "" {
		gp.Out = gp.Name
	}
	if !strings.HasSuffix(gp.Out, ".so") {
		gp.Out = fmt.Sprintf("%s.so", gp.Out)
	}

	opts := []zen_targets.TargetOption{
		zen_targets.WithSrcs(map[string][]string{"_srcs": gp.Srcs}),
		zen_targets.WithOuts([]string{gp.Out}),
		zen_targets.WithVisibility(gp.Visibility),
		zen_targets.WithEnvVars(gp.Env),
		zen_targets.WithTargetScript("build", &zen_targets.TargetScript{
			Deps: gp.Deps,
			Run: func(target *zen_targets.Target, runCtx *zen_targets.RuntimeContext) error {
				env_vars := target.GetEnvironmentVariablesList()

				splitCmd := strings.Split(fmt.Sprintf("build -buildmode=plugin -o %s %s", filepath.Join(target.Cwd, gp.Out), gp.Flags), " ")
				cmd := exec.Command("go", splitCmd...)
				cmd.Dir = target.Cwd
				cmd.Env = env_vars
				cmd.Stdout = target
				cmd.Stderr = target
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("executing build: %w", err)
				}

				return nil
			},
		}),
	}

	return []*zen_targets.Target{
		zen_targets.NewTarget(
			gp.Name,
			opts...,
		),
	}, nil
}
