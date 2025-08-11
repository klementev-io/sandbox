package config

type Cfg struct {
	Service   Service   `mapstructure:"service"`
	APIServer APIServer `mapstructure:"api_server"`
	Log       Log       `mapstructure:"log"`
}

type Service struct {
	Name string
}

type APIServer struct {
	Host  string `mapstructure:"host"`
	Port  string `mapstructure:"port"`
	Pprof bool   `mapstructure:"pprof"`
}

type Log struct {
	Level string `mapstructure:"level"`
}
