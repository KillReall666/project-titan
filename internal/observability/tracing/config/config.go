package config

import (
	"github.com/caarlos0/env/v11"
	"time"
)

// Config Можно расширять: sampler ratio, headers, tls, env, etc.
type Config struct {
	ServiceName     string
	Endpoint        string // OTLP gRPC endpoint, например "localhost:4317"
	Insecure        bool   // для локалки true
	SampleAlways    bool   // для обучения true, для прода false
	ShutdownTimeout time.Duration
}

const (
	ServiceName  = "api-gateway"
	Endpoint     = "localhost:4317"
	Insecure     = true
	SampleAlways = true
)

func New() (*Config, error) {
	cfg := &Config{
		ServiceName:  ServiceName,
		Endpoint:     Endpoint,
		Insecure:     Insecure,
		SampleAlways: SampleAlways,
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil

}
