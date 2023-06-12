package golang

import (
	ahoy_targets "gitlab.com/hidothealth/platform/ahoy/src/target"
)

var KnownTargets = ahoy_targets.TargetCreatorMap{
	"go_plugin": GolangPluginConfig{},
	"go_binary": GolangBinaryConfig{},
}
