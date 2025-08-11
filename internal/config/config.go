package config

type Cfg struct {
	Service     string      `mapstructure:"service"`
	APIServer   APIServer   `mapstructure:"api_server"`
	PprofServer PprofServer `mapstructure:"pprof_server"`
	Log         Log         `mapstructure:"log"`
}

type APIServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type PprofServer struct {
	Host   string `mapstructure:"host"`
	Port   string `mapstructure:"port"`
	Enable bool   `mapstructure:"enable"`
}

type Log struct {
	Level string `mapstructure:"level"`
}
