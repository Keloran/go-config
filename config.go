package ConfigBuilder

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/keloran/go-config/database"
	"github.com/keloran/go-config/keycloak"
	"github.com/keloran/go-config/local"
	"github.com/keloran/go-config/mongo"
	"github.com/keloran/go-config/rabbit"
	"github.com/keloran/go-config/vault"
	vaultHelper "github.com/keloran/vault-helper"
)

type Config struct {
	vaultHelper vaultHelper.VaultHelper

	local.Local
	vault.Vault
	database.Database
	keycloak.Keycloak
	mongo.Mongo
	rabbit.Rabbit
	ProjectConfigurator
}

type BuildOption func(*Config) error

func Local(cfg *Config) error {
	l, err := local.Build()
	if err != nil {
		return logs.Errorf("build local: %v", err)
	}

	cfg.Local = *l

	return nil
}

func Vault(cfg *Config) error {
	v, err := vault.Build()
	if err != nil {
		return logs.Errorf("build vault: %v", err)
	}

	cfg.Vault = *v

	return nil
}

func Database(cfg *Config) error {
	vh := cfg.vaultHelper
	if vh == nil {
		vh = vaultHelper.NewVault(cfg.Vault.Address, cfg.Vault.Token)
	}

	d, err := database.Build(database.Setup(cfg.Vault.Address, cfg.Vault.Token), vh)
	if err != nil {
		return logs.Errorf("build database: %v", err)
	}

	cfg.Database = *d

	return nil
}

func Mongo(cfg *Config) error {
	vh := cfg.vaultHelper
	if vh == nil {
		vh = vaultHelper.NewVault(cfg.Vault.Address, cfg.Vault.Token)
	}

	m, err := mongo.Build(mongo.Setup(cfg.Vault.Address, cfg.Vault.Token), vh)
	if err != nil {
		return logs.Errorf("build mongo: %v", err)
	}

	cfg.Mongo = *m

	return nil
}

func Keycloak(cfg *Config) error {
	k, err := keycloak.Build()
	if err != nil {
		return logs.Errorf("build keycloak: %v", err)
	}

	cfg.Keycloak = *k

	return nil
}

func Rabbit(cfg *Config) error {
	vh := cfg.vaultHelper
	if vh == nil {
		vh = vaultHelper.NewVault(cfg.Vault.Address, cfg.Vault.Token)
	}

	r, err := rabbit.Build(rabbit.Setup(cfg.Vault.Address, cfg.Vault.Token), vh)
	if err != nil {
		return logs.Errorf("build rabbit: %v", err)
	}
	cfg.Rabbit = *r

	return nil
}

func Build(opts ...BuildOption) (*Config, error) {
	cfg := &Config{}
	if err := cfg.Build(opts...); err != nil {
		return nil, logs.Errorf("build config: %v", err)
	}

	return cfg, nil
}

func BuildLocal(mockVault vaultHelper.VaultHelper, opts ...BuildOption) (*Config, error) {
	cfg := &Config{}

	cfg.vaultHelper = mockVault
	if err := cfg.Build(opts...); err != nil {
		return nil, logs.Errorf("build config: %v", err)
	}

	return cfg, nil
}

func (c *Config) Build(opts ...BuildOption) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return logs.Errorf("build configOptions: %v", err)
		}
	}

	return nil
}

// Project Configurators
type ProjectConfigurator interface {
	Build(opts ...BuildOption) error
}

func WithProjectConfigurator(pc ProjectConfigurator) BuildOption {
	return func(c *Config) error {
		if err := pc.Build(); err != nil {
			return logs.Errorf("withProjectConfigurator: %v", err)
		}

		return nil
	}
}
