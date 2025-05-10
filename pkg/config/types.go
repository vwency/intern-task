package config

type AppConfig struct {
	GRPC struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"grpc"`
	SubPub struct {
		MaxSubscribers  int    `mapstructure:"max_subscribers"`
		ShutdownTimeout string `mapstructure:"shutdown_timeout"`
	} `mapstructure:"subpub"`
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
}
