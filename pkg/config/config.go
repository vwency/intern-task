package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Load(configPath string) (*AppConfig, error) {
	v := viper.New()

	// Базовый конфиг
	v.SetConfigName("config")
	v.AddConfigPath(configPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config: %w", err)
	}

	// Загрузка env-специфичного конфига
	env := DetectEnv()
	if env != "" {
		envConfigName := fmt.Sprintf("config.%s", env)
		v.SetConfigName(envConfigName)
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("failed to merge %s config: %w", env, err)
		}
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
