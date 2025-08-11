package internal

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// LoadConfig load config from yaml file.
func LoadConfig[T any]() (*T, error) {
	const configFlagName = "c"

	pflag.String(configFlagName, "", "config file path")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, err
	}

	viper.SetConfigFile(viper.GetString(configFlagName))

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, fmt.Errorf("config file not found: %w", err)
		}

		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	cfg := new(T)
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}
