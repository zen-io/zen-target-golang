package golang

import (
	zen_targets "github.com/zen-io/zen-core/target"
)

var KnownTargets = zen_targets.TargetCreatorMap{
	"go_plugin": GolangPluginConfig{},
	"go_binary": GolangBinaryConfig{},
}
