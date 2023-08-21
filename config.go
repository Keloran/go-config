package ConfigBuilder

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	"github.com/keloran/go-config/local"
)

type Config struct {
	local.Local
	//Vault
	//Services
}

type BuildOption func(*Config) error

func BuildLocal(cfg *Config) error {
	l, err := local.Build()
	if err != nil {
		return logs.Errorf("build local: %v", err)
	}

	cfg.Local = *l

	return nil
}

func BuildVault(cfg *Config) error {
	return nil
}

func BuildServices(cfg *Config) error {
	return nil
}

func Build(opts ...BuildOption) (*Config, error) {
	cfg := &Config{}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, logs.Errorf("build configOptions: %v", err)
		}
	}

	if err := env.Parse(cfg); err != nil {
		return nil, logs.Errorf("parse config: %v", err)
	}

	return cfg, nil
}
