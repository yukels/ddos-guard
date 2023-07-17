package config

import (
	"path"
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
)

const (
	defaultConfigDir = "/etc/ddos-guard"
)

type configActionFunc func(ctx context.Context, config interface{}, sectionName string) error

// ConfigDir base configuration folder
var ConfigDir string

// validate interface for config files
type validate interface {
	Validate(ctx context.Context, configDir string) error
}

// InitFromDir initializes configuration from given configuration directory, validates and decrypts it (if needed)
func InitFromDir(ctx context.Context, configDir string, sections map[string]interface{}) {
	ConfigDir = configDir
	for sectionName, section := range sections {
		sectionViper := viper.New()
		sectionViper.AddConfigPath(configDir)
		sectionViper.SetConfigName(sectionName)
		sectionViper.SetConfigType("json")

		err := sectionViper.ReadInConfig()
		if err != nil {
			log.Log(ctx).WithError(errors.Errorf("Fatal error reading %s config file: %s", sectionName, err)).Panic("Config file read error")
		}

		sectionViper.Unmarshal(section)
		if err = validateConfig(ctx, section, sectionName); err != nil {
			log.Log(ctx).WithError(errors.Errorf("Error validating %s config file; %s", sectionName, err)).Panic("Config file validation error")
			return
		}
	}
}

// InitFromDefaultDir initializes configuration from default directory.
func InitFromDefaultDir(ctx context.Context, sections map[string]interface{}) {
	InitFromDir(ctx, defaultConfigDir, sections)
}

// GetAbsPath get absolute path for provider entity based on the configuration folder
func GetAbsPath(ctx context.Context, dir string) string {
	if !path.IsAbs(dir) {
		dir = path.Clean(path.Join(ConfigDir, dir))
	}
	return dir
}

func multiConfig(ctx context.Context, config interface{}, sectionName string, action configActionFunc) error {
	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Ptr {
		return errors.Errorf("The provided type MUST be Ptr, section [%s]", sectionName)
	}

	if rv.Elem().Kind() == reflect.Map {
		for _, k := range rv.Elem().MapKeys() {
			if err := action(ctx, rv.Elem().MapIndex(k).Interface(), sectionName); err != nil {
				return err
			}
		}
		return nil
	}

	return action(ctx, config, sectionName)
}

func validateConfig(ctx context.Context, config interface{}, sectionName string) error {
	return multiConfig(ctx, config, sectionName, validateSingleConfig)
}

// validateSingleConfig checks if a config struct can apply validate - applies it if true
func validateSingleConfig(ctx context.Context, config interface{}, sectionName string) error {
	confVal, okVal := config.(validate)
	if okVal {
		if err := confVal.Validate(ctx, ConfigDir); err != nil {
			return err
		}
	}
	return nil
}
