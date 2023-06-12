package golang

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	ahoy_targets "gitlab.com/hidothealth/platform/ahoy/src/target"
)

type GolangPluginConfig struct {
	ahoy_targets.BaseFields `mapstructure:",squash"`
	Srcs                    []string `mapstructure:"srcs"`
	Out                     string   `mapstructure:"out"`
	Flags                   string   `mapstructure:"flags"`
}

func (gp GolangPluginConfig) GetTargets(_ *ahoy_targets.TargetConfigContext) ([]*ahoy_targets.Target, error) {
	if gp.Out == "" {
		gp.Out = gp.Name
	}
	if !strings.HasSuffix(gp.Out, ".so") {
		gp.Out = fmt.Sprintf("%s.so", gp.Out)
	}

	opts := []ahoy_targets.TargetOption{
		ahoy_targets.WithSrcs(map[string][]string{"_srcs": gp.Srcs}),
		ahoy_targets.WithOuts([]string{gp.Out}),
		ahoy_targets.WithVisibility(gp.Visibility),
		ahoy_targets.WithEnvVars(gp.Env),
		ahoy_targets.WithTargetScript("build", &ahoy_targets.TargetScript{
			Deps: gp.Deps,
			Run: func(target *ahoy_targets.Target, runCtx *ahoy_targets.RuntimeContext) error {
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

	return []*ahoy_targets.Target{
		ahoy_targets.NewTarget(
			gp.Name,
			opts...,
		),
	}, nil
}
