package golang

// import (
// 	"fmt"
// 	"github.com/zen-io/zen-core/target"
// 	"github.com/zen-io/zen/src/zen/own"

// 	"github.com/mitchellh/mapstructure"
// )

// type GolangToolchainConfig struct {
// 	Name       string            `mapstructure:"name"`
// 	Version    string            `mapstructure:"version"`
// 	Labels     []string          `mapstructure:"labels"`
// 	Env        map[string]string `mapstructure:"env"`
// 	PassEnv    []string          `mapstructure:"pass_env"`
// 	Deps       []string          `mapstructure:"deps"`
// 	Visibility []string          `mapstructure:"visibility"`
// }

// func (gtc GolangToolchainConfig) GetTargets(block interface{}, tcc *zen_targets.TargetConfigContext) ([]*zen_targets.Target, error) {
// 	mapstructure.Decode(block, &gtc)

// 	targets, err := own.RemoteFileConfig{
// 		Name:          fmt.Sprintf("%s_source", gtc.Name),
// 		Url:           fmt.Sprintf("https://go.dev/dl/go%s.{CONFIG.HOSTOS}-{CONFIG.HOSTARCH}.tar.gz", gtc.Version),
// 		Extract:       true,
// 		ExportedFiles: []string{"extract/go/"},
// 	}.ExportTargets(tcc)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return targets, nil
// }
