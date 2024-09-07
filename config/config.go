package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App    `yaml:"app"`
		HTTP   `yaml:"http"`
		Log    `yaml:"log"`
		PG     `yaml:"postgres"`
		JWT    `yaml:"jwt"`
		Hasher `yaml:"hasher"`
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
		URL         string `env-required:"true" env:"PG_URL"`
	}

	JWT struct {
		SignKey  string        `env-required:"true" env:"JWT_SIGN_KEY"`
		TokenTTL time.Duration `env-required:"true" yaml:"token_ttl"`
	}

	Hasher struct {
		Salt string `env-required:"true" env:"PASSWORD_HASHER_SALT"`
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
