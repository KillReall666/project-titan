package config

import (
	"errors"
	"flag"
	"github.com/caarlos0/env/v11"
	"time"
)

type Config struct {
	Address string `env:"RUN_ADDRESS"`
	ConnStr string `env:"DATABASE_URL"`
	JWT     JWT
}

type JWT struct {
	Secret     string        `env:"SECRET_KEY"`
	AccessTLL  time.Duration `env:"ACCESS_TOKEN_TTL"`
	RefreshTLL time.Duration `env:"REFRESH_TOKEN_TTL"`
	Issuer     string        `env:"ISSUER_AUTH_SERVICE"`
}

const (
	defaultServer  = "localhost:8080"
	JWTSecretKey   = "secret-for-signing-key"
	accessTokenTTL = time.Minute * 60
)

func New() (*Config, error) {
	cfg := &Config{
		Address: defaultServer,
		JWT: JWT{
			Secret:    JWTSecretKey,
			AccessTLL: accessTokenTTL,
		},
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
