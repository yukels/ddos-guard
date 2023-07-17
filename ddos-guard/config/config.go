package config

type configs struct {
	DdosGuardConfig DdosGuardConfig
}

// Configs - Returns all the configuration sections
var Configs configs

var Sections = map[string]interface{}{
	"ddos-guard": &Configs.DdosGuardConfig,
}
