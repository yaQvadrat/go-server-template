package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)


type (
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Log     `yaml:"log"`
		PG      `yaml:"postgres"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"`
		Version string `env-required:"true" yaml:"version"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level"`
	}

	PG struct {
		MaxPoolSize int    `env-required:"true" yaml:"max_pool_size"`
		URL         string `env-required:"true"                      env:"PG_URL"`
	}
)


func New(configPath string) (*Config, error) {
	cfg := &Config{}

	// Reading config from .yaml file
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("config - New - cleanenv.ReadConfig: %w", err)
	}

	// Reading (and override if it possible) config from .env
	if err := cleanenv.UpdateEnv(cfg); err != nil {
		return nil, fmt.Errorf("config - New - cleanenv.UpdateEnv: %w", err)
	}

	return cfg, nil
}
