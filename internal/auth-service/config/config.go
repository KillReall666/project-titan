package config

import (
	"errors"
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address string `env:"RUN_ADDRESS"`
	ConnStr string `env:"DATABASE_URL"`
}

const (
	defaultServer = "localhost:8080"
)

func New() (*Config, error) {
	cfg := &Config{
		Address: defaultServer,
	}

	flag.StringVar(&cfg.Address, "a", cfg.Address, "server address [host:port]")
	flag.StringVar(&cfg.ConnStr, "d", "", "database connection string")
	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.ConnStr == "" {
		return errors.New("DATABASE_URL is required")
	}

	if cfg.Address == "" {
		return errors.New("RUN_ADDRESS is required")
	}

	return nil
}
