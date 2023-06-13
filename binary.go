package golang

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	zen_targets "github.com/zen-io/zen-core/target"
)

type GolangBinaryConfig struct {
	zen_targets.BaseFields `mapstructure:",squash"`
	Srcs                   []string `mapstructure:"srcs"`
	Out                    *string  `mapstructure:"out"`
	Path                   *string  `mapstructure:"path"`
	Flags                  string   `mapstructure:"flags"`
	Toolchain              *string  `mapstructure:"toolchain"`
}

func (gb GolangBinaryConfig) GetTargets(tcc *zen_targets.TargetConfigContext) ([]*zen_targets.Target, error) {
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

	return []*zen_targets.Target{
		zen_targets.NewTarget(
			gb.Name,
			zen_targets.WithSrcs(map[string][]string{"_srcs": gb.Srcs}),
			zen_targets.WithOuts([]string{*gb.Out}),
			zen_targets.WithVisibility(gb.Visibility),
			zen_targets.WithTools(tools),
			zen_targets.WithEnvVars(gb.Env),
			zen_targets.WithPassEnv(gb.PassEnv),
			zen_targets.WithTargetScript("build", &zen_targets.TargetScript{
				Deps: gb.Deps,
				Pre: func(target *zen_targets.Target, runCtx *zen_targets.RuntimeContext) error {
					if _, ok := target.Env["GOROOT"]; !ok {
						target.Env["GOROOT"] = target.Tools["golang"]
					}
					target.Env["PATH"] = fmt.Sprintf("%s/bin:%s", target.Tools["golang"], target.Env["PATH"])

					if _, ok := target.Env["GOBIN"]; !ok {
						target.Env["GOBIN"] = fmt.Sprintf("%s/bin/go", target.Tools["golang"])
					}

					target.Env["ZEN_DEBUG_CMD"] = fmt.Sprintf("%s/bin/go build -o %s %s", target.Tools["golang"], filepath.Join(target.Cwd, *gb.Out), gb.Flags)
					if runCtx.Shell {
						if gb.Path != nil {
							target.Cwd = filepath.Join(target.Cwd, *gb.Path)
						}
					}

					return nil
				},
				Run: func(target *zen_targets.Target, runCtx *zen_targets.RuntimeContext) error {
					if gb.Path != nil {
						target.Cwd = filepath.Join(target.Cwd, *gb.Path)
					}

					if err := runTidy(target); err != nil {
						return err
					}

					spl := strings.Split(target.Env["ZEN_DEBUG_CMD"], " ")
					cmd := exec.Command(spl[0], spl[1:]...)
					cmd.Dir = target.Cwd
					cmd.Env = target.GetEnvironmentVariablesList()
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						return fmt.Errorf("executing build: %w", err)
					}

					return nil
				},
			}),
		),
	}, nil
}
