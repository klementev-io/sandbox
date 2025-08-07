package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envPrefix = "SB"
	flagName  = "config"
)

type Cfg struct {
	Service   Service   `mapstructure:"service"`
	APIServer APIServer `mapstructure:"api_server"`
	Logger    Logger    `mapstructure:"logger"`
}

type Service struct {
	Name    string
	Version string
}

type APIServer struct {
	Host  string `mapstructure:"host"`
	Port  string `mapstructure:"port"`
	Pprof bool   `mapstructure:"pprof"`
}

type Logger struct {
	Level string `mapstructure:"level"`
}

// LoadConfig load config from yaml file, envs.
func LoadConfig() (Cfg, error) {
	dir, err := os.Getwd()
	if err != nil {
		return Cfg{}, err
	}

	defaultPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "configs/default-config.yaml")

	pflag.String(flagName, defaultPath, "config file path")
	pflag.Parse()

	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return Cfg{}, err
	}

	viper.SetConfigFile(viper.GetString("config"))
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

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
