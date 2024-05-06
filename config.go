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
	VaultHelper *vaultHelper.VaultHelper
  VaultPaths vault.VaultPaths

	Local    local.System
	Vault    vault.System
	Database database.Details
	Keycloak keycloak.System
	Mongo    mongo.System
	Rabbit   rabbit.System

	// Project level properties
	ProjectProperties map[string]interface{}
}

type BuildOption func(*Config) error

func NewConfig(vh *vaultHelper.VaultHelper) *Config {
	return &Config{
    VaultHelper: vh,
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
  d := database.NewSystem()
  if cfg.VaultHelper != nil {
    vd := database.VaultDetails{
      Address: cfg.Vault.Address,
      Token: cfg.Vault.Token,
    }
    d.Setup(vd, *cfg.VaultHelper)
  }
  db, err := d.Build()
  if err != nil {
    return logs.Errorf("database failed to build: %v", err)
  }


  cfg.Database = *db

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
  k := keycloak.NewSystem()
  if cfg.VaultHelper != nil {
    vd := keycloak.VaultDetails{
      Address: cfg.Vault.Address,
      Token: cfg.Vault.Token,
    }
    k.Setup(vd, *cfg.VaultHelper)
  }

  _, err := k.Build()
  if err != nil {
    return logs.Errorf("failed to build keycloak: %v", err)
  }
  cfg.Keycloak = *k
	return nil
}

func Rabbit(cfg *Config) error {
  r := rabbit.NewSystem(&http.Client{})
  if cfg.VaultHelper != nil {
    vd := rabbit.VaultDetails{
      Address: cfg.Vault.Address,
      Token: cfg.Vault.Token,
    }
    r.Setup(vd, *cfg.VaultHelper)
  }
  _, err := r.Build()
  if err != nil {
    return logs.Errorf("failed to build rabbit: %v", err)
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
