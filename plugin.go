package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	zen_targets "github.com/zen-io/zen-core/target"
)

type GolangPluginConfig struct {
	Name        string            `mapstructure:"name" desc:"Name for the target"`
	Description string            `mapstructure:"desc" desc:"Target description"`
	Labels      []string          `mapstructure:"labels" desc:"Labels to apply to the targets"` //
	Deps        []string          `mapstructure:"deps" desc:"Build dependencies"`
	PassEnv     []string          `mapstructure:"pass_env" desc:"List of environment variable names that will be passed from the OS environment, they are part of the target hash"`
	SecretEnv   []string          `mapstructure:"secret_env" desc:"List of environment variable names that will be passed from the OS environment, they are not used to calculate the target hash"`
	Env         map[string]string `mapstructure:"env" desc:"Key-Value map of static environment variables to be used"`
	Tools       map[string]string `mapstructure:"tools" desc:"Key-Value map of tools to include when executing this target. Values can be references"`
	Visibility  []string          `mapstructure:"visibility" desc:"List of visibility for this target"`
	Srcs        []string          `mapstructure:"srcs"`
	Out         string            `mapstructure:"out"`
	Path        *string           `mapstructure:"path"`
	Flags       string            `mapstructure:"flags"`
	Toolchain   *string           `mapstructure:"toolchain"`
}

func (gp GolangPluginConfig) GetTargets(tcc *zen_targets.TargetConfigContext) ([]*zen_targets.Target, error) {
	if gp.Out == "" {
		gp.Out = gp.Name
	}
	if !strings.HasSuffix(gp.Out, ".so") {
		gp.Out = fmt.Sprintf("%s.so", gp.Out)
	}

	tools := map[string]string{}
	if gp.Toolchain != nil {
		tools["golang"] = *gp.Toolchain
	} else if val, ok := tcc.KnownToolchains["golang"]; !ok {
		return nil, fmt.Errorf("golang toolchain is not configured")
	} else {
		tools["golang"] = val
	}

	opts := []zen_targets.TargetOption{
		zen_targets.WithSrcs(map[string][]string{"_srcs": gp.Srcs}),
		zen_targets.WithOuts([]string{gp.Out}),
		zen_targets.WithVisibility(gp.Visibility),
		zen_targets.WithTools(tools),
		zen_targets.WithEnvVars(gp.Env),
		zen_targets.WithTargetScript("build", &zen_targets.TargetScript{
			Deps: gp.Deps,

			Pre: func(target *zen_targets.Target, runCtx *zen_targets.RuntimeContext) error {
				if _, ok := target.Env["GOROOT"]; !ok {
					target.Env["GOROOT"] = target.Tools["golang"]
				}
				target.Env["PATH"] = fmt.Sprintf("%s/bin:%s", target.Tools["golang"], target.Env["PATH"])

				if _, ok := target.Env["GOBIN"]; !ok {
					target.Env["GOBIN"] = fmt.Sprintf("%s/bin/go", target.Tools["golang"])
				}

				target.Env["ZEN_DEBUG_CMD"] = fmt.Sprintf("%s/bin/go build -buildmode=plugin -o %s %s", target.Tools["golang"], filepath.Join(target.Cwd, gp.Out), gp.Flags)
				if gp.Path != nil {
					target.Cwd = filepath.Join(target.Cwd, *gp.Path)
				}

				return nil
			},
			Run: func(target *zen_targets.Target, runCtx *zen_targets.RuntimeContext) error {
				if err := target.Exec([]string{"go", "mod", "tidy"}, "tidy"); err != nil {
					return err
				}

				return target.Exec(strings.Split(target.Env["ZEN_DEBUG_CMD"], " "), "go build")
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
