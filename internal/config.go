package internal

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"path"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServer HTTPServerConfig `mapstructure:"http_server"`
	Logger     LoggerConfig     `mapstructure:"logger"`
}

type HTTPServerConfig struct {
	Host      string `mapstructure:"host"`
	Port      string `mapstructure:"port"`
	PprofPort string `mapstructure:"pprof_port"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func LoadConfig() (Config, error) {
	const (
		envPrefix      = "SX"
		configFlagName = "c"
	)

	cfgFilePath := flag.String(configFlagName, "./configs/config.yaml", "config file path")
	flag.Parse()

	viper.SetConfigName(path.Base(*cfgFilePath))
	viper.SetConfigType(path.Ext(*cfgFilePath)[1:])
	viper.AddConfigPath(path.Dir(*cfgFilePath))

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	bindEnvs()

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			slog.Default().WarnContext(context.Background(), "config file not found, using environment variables")
		} else {
			return Config{}, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}

func bindEnvs() {
	_ = viper.BindEnv("http_server.host")
	_ = viper.BindEnv("http_server.port")
	_ = viper.BindEnv("http_server.pprof_port")
	_ = viper.BindEnv("logger.level")
	_ = viper.BindEnv("logger.format")
}
