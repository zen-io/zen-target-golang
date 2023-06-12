package golang

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	ahoy_targets "gitlab.com/hidothealth/platform/ahoy/src/target"
)

type GolangBinaryConfig struct {
	ahoy_targets.BaseFields `mapstructure:",squash"`
	Srcs                    []string `mapstructure:"srcs"`
	Out                     *string  `mapstructure:"out"`
	Path                    *string  `mapstructure:"path"`
	Flags                   string   `mapstructure:"flags"`
	Toolchain               *string  `mapstructure:"toolchain"`
}

func (gb GolangBinaryConfig) GetTargets(tcc *ahoy_targets.TargetConfigContext) ([]*ahoy_targets.Target, error) {
	if gb.Out == nil {
		gb.Out = &gb.Name
	}

	tools := map[string]string{}
	if gb.Toolchain != nil {
		tools["golang"] = *gb.Toolchain
	} else if val, ok := tcc.KnownToolchains["golang"]; !ok {
		return nil, fmt.Errorf("golang toolchain is not configured")
	} else {
		tools["golang"] = val
	}

	return []*ahoy_targets.Target{
		ahoy_targets.NewTarget(
			gb.Name,
			ahoy_targets.WithSrcs(map[string][]string{"_srcs": gb.Srcs}),
			ahoy_targets.WithOuts([]string{*gb.Out}),
			ahoy_targets.WithVisibility(gb.Visibility),
			ahoy_targets.WithTools(tools),
			ahoy_targets.WithEnvVars(gb.Env),
			ahoy_targets.WithPassEnv(gb.PassEnv),
			ahoy_targets.WithTargetScript("build", &ahoy_targets.TargetScript{
				Deps: gb.Deps,
				Pre: func(target *ahoy_targets.Target, runCtx *ahoy_targets.RuntimeContext) error {
					if _, ok := target.Env["GOROOT"]; !ok {
						target.Env["GOROOT"] = target.Tools["golang"]
					}
					target.Env["PATH"] = fmt.Sprintf("%s/bin:%s", target.Tools["golang"], target.Env["PATH"])

					if _, ok := target.Env["GOBIN"]; !ok {
						target.Env["GOBIN"] = fmt.Sprintf("%s/bin/go", target.Tools["golang"])
					}

					return nil
				},
				Run: func(target *ahoy_targets.Target, runCtx *ahoy_targets.RuntimeContext) error {
					cwd := filepath.Join(target.Cwd, *gb.Path)
					if err := runTidy(cwd, target.GetEnvironmentVariablesList(), target); err != nil {
						return err
					}

					cmdArg := fmt.Sprintf("build -o %s %s", filepath.Join(target.Cwd, *gb.Out), gb.Flags)
					target.Debug(fmt.Sprintf("go %s", cmdArg))
					cmd := exec.Command("go", strings.Split(cmdArg, " ")...)
					cmd.Dir = cwd
					cmd.Env = target.GetEnvironmentVariablesList()
					cmd.Stdout = target
					cmd.Stderr = target
					if err := cmd.Run(); err != nil {
						return fmt.Errorf("executing build: %w", err)
					}

					return nil
				},
			}),
		),
	}, nil
}
