package config

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConfigPath = "configs/default-config.yaml"
	configFlagName    = "config"
)

// LoadConfig load config from yaml file, envs.
func LoadConfig() (Cfg, error) {
	pflag.String(configFlagName, defaultConfigPath, "config file path")
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return Cfg{}, err
	}

	viper.SetConfigFile(viper.GetString(configFlagName))

	err = viper.ReadInConfig()
	if err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return Cfg{}, fmt.Errorf("config file not found: %w", err)
		}

		return Cfg{}, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Cfg
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return Cfg{}, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}
