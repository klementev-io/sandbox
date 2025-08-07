package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/klementev-io/sandbox/internal/config"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("SB_LOGGER_LEVEL", "debug")

	want := config.Cfg{
		Service: config.Service{
			Name:    "sandbox",
			Version: "unknown",
		},
		APIServer: config.APIServer{
			Host:  "127.0.0.1",
			Port:  "8080",
			Pprof: true,
		},
		Logger: config.Logger{
			Level: "debug",
		},
	}

	got, err := config.LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}
