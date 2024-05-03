package ConfigBuilder

import (
	"net/http"

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
  VaultExclusive bool
  VaultPaths vault.VaultPaths

	Local    local.System
	Vault    vault.System
	Database database.System
	Keycloak keycloak.System
	Mongo    mongo.System
	Rabbit   rabbit.System

	// Project level properties
	ProjectProperties map[string]interface{}
}

type BuildOption func(*Config) error

func NewConfig(vaultExclusive bool) *Config {
	return &Config{
		VaultExclusive: vaultExclusive,
	}
}

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

  if cfg.VaultExclusive {
    d, err := database.Build(database.Setup(cfg.Vault.Address, cfg.Vault.Token, true, &cfg.VaultPaths), vh)
    if err != nil {
      return logs.Errorf("build database: %v", err)
    }

    cfg.Database = *d
    return nil
  }

  d, err := database.Build(database.Setup(cfg.Vault.Address, cfg.Vault.Token, false, nil), vh)
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

	m, err := mongo.Build(mongo.Setup(cfg.Vault.Address, cfg.Vault.Token, cfg.VaultExclusive), vh)
	if err != nil {
		return logs.Errorf("build mongo: %v", err)
	}

	cfg.Mongo = *m

	return nil
}

func Keycloak(cfg *Config) error {
  vh := cfg.vaultHelper
  if vh == nil {
    vh = vaultHelper.NewVault(cfg.Vault.Address, cfg.Vault.Token)
  }

  k, err := keycloak.Build(keycloak.Setup(cfg.Vault.Address, cfg.Vault.Token, cfg.VaultExclusive), vh)
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

	r, err := rabbit.Build(rabbit.Setup(cfg.Vault.Address, cfg.Vault.Token, cfg.VaultExclusive), vh, &http.Client{})
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

type ProjectConfigurator interface {
	Build(*Config) error
}

func WithProjectConfigurator(pc ProjectConfigurator) BuildOption {
	return func(c *Config) error {
		if err := pc.Build(c); err != nil {
			return logs.Errorf("withProjectConfigurator: %v", err)
		}

		return nil
	}
}
