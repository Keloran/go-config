package ConfigBuilder

import (
	"net/http"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/keloran/go-config/database"
	"github.com/keloran/go-config/influx"
	"github.com/keloran/go-config/keycloak"
	"github.com/keloran/go-config/local"
	"github.com/keloran/go-config/mongo"
	"github.com/keloran/go-config/rabbit"
	"github.com/keloran/go-config/vault"
	vaultHelper "github.com/keloran/vault-helper"
)

type Config struct {
	VaultHelper *vaultHelper.VaultHelper
	VaultPaths  vault.Paths
	VaultInject bool

	Local    local.System
	Vault    vault.System
	Database database.Details
	Keycloak keycloak.System
	Mongo    mongo.System
	Rabbit   rabbit.System
	Influx   influx.System

	// Project level properties
	ProjectProperties map[string]interface{}
}

type BuildOption func(*Config) error

func NewConfig(vh vaultHelper.VaultHelper) *Config {
	return &Config{
		VaultHelper: &vh,
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
	v, vh, err := vault.Build()
	if err != nil {
		return logs.Errorf("build vault: %v", err)
	}

	cfg.Vault = *v
	cfg.VaultHelper = &vh

	return nil
}

func Database(cfg *Config) error {
	d := database.NewSystem()
	if cfg.VaultHelper != nil {
		vd := database.VaultDetails{}
		if cfg.VaultPaths != (vault.Paths{}) {
			if cfg.VaultPaths.Database.Details != "" {
				vd.DetailsPath = cfg.VaultPaths.Database.Details
			}
			if cfg.VaultPaths.Database.Credentials != "" {
				vd.CredPath = cfg.VaultPaths.Database.Credentials
			}
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
	m := mongo.NewSystem()
	if cfg.VaultHelper != nil {
		vd := mongo.VaultDetails{}
		if cfg.VaultPaths != (vault.Paths{}) {
			if cfg.VaultPaths.Mongo.Details != "" {
				vd.DetailsPath = cfg.VaultPaths.Mongo.Details
			}
			if cfg.VaultPaths.Mongo.Credentials != "" {
				vd.CredPath = cfg.VaultPaths.Mongo.Credentials
			}
		}

		m.Setup(vd, *cfg.VaultHelper)
	}
	_, err := m.Build()
	if err != nil {
		return logs.Errorf("failed to build mongo: %v", err)
	}
	cfg.Mongo = *m

	return nil
}

func Keycloak(cfg *Config) error {
	k := keycloak.NewSystem()
	if cfg.VaultHelper != nil {
		vd := keycloak.VaultDetails{}
		if cfg.VaultPaths != (vault.Paths{}) {
			if cfg.VaultPaths.Keycloak.Details != "" {
				vd.DetailsPath = cfg.VaultPaths.Keycloak.Details
			}
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
		vd := rabbit.VaultDetails{}
		if cfg.VaultPaths != (vault.Paths{}) {
			if cfg.VaultPaths.Rabbit.Details != "" {
				vd.DetailsPath = cfg.VaultPaths.Rabbit.Details
			}
			if cfg.VaultPaths.Rabbit.Credentials != "" {
				vd.CredPath = cfg.VaultPaths.Rabbit.Credentials
			}
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

func Influx(cfg *Config) error {
	i := influx.NewSystem()
	if cfg.VaultHelper != nil {
		vd := influx.VaultDetails{}
		if cfg.VaultPaths != (vault.Paths{}) {
			if cfg.VaultPaths.Influx.Details != "" {
				vd.DetailsPath = cfg.VaultPaths.Influx.Details
			}
		}

		i.Setup(vd, *cfg.VaultHelper)
	}

	_, err := i.Build()
	if err != nil {
		return logs.Errorf("failed to build influx: %v", err)
	}

	cfg.Influx = *i

	return nil
}

func Build(opts ...BuildOption) (*Config, error) {
	cfg := &Config{}

	if err := cfg.Build(opts...); err != nil {
		return nil, logs.Errorf("build config: %v", err)
	}

	return cfg, nil
}

func BuildLocal(opts ...BuildOption) (*Config, error) {
	cfg := &Config{}

	if err := cfg.Build(opts...); err != nil {
		return nil, logs.Errorf("build config: %v", err)
	}

	return cfg, nil
}

func BuildLocalVH(mockVault vaultHelper.VaultHelper, opts ...BuildOption) (*Config, error) {
	cfg := &Config{
		VaultHelper: &mockVault,
	}

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
