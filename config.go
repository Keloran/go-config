package ConfigBuilder

import (
	"github.com/keloran/go-config/auth/clerk"
	"github.com/keloran/go-config/flags"
	"github.com/keloran/go-config/notify/resend"
	"net/http"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/keloran/go-config/auth/keycloak"
	"github.com/keloran/go-config/bugfixes"
	"github.com/keloran/go-config/database/mongo"
	"github.com/keloran/go-config/database/postgres"
	"github.com/keloran/go-config/influx"
	"github.com/keloran/go-config/local"
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
	Database postgres.System
	Keycloak keycloak.System
	Mongo    mongo.System
	Rabbit   rabbit.System
	Influx   influx.System
	Bugfixes bugfixes.System
	Clerk    clerk.System
	Resend   resend.System
	Flags    flags.System

	// Project level properties
	ProjectProperties map[string]interface{}
}

type BuildOption func(*Config) error

func NewConfig(vh vaultHelper.VaultHelper) *Config {
	return &Config{
		VaultHelper: &vh,
	}
}

func NewConfigNoVault() *Config {
	return &Config{}
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

// Database deprecated: use Postgres instead
func Database(cfg *Config) error {
	return Postgres(cfg)
}

func Postgres(cfg *Config) error {
	d := postgres.NewSystem()
	if cfg.VaultHelper != nil {
		vd := postgres.VaultDetails{}
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
	_, err := d.Build()
	if err != nil {
		return logs.Errorf("database failed to build: %v", err)
	}

	cfg.Database = *d

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

func Clerk(cfg *Config) error {
	c := clerk.NewSystem()
	if cfg.VaultHelper != nil {
		vd := c.VaultDetails
		if cfg.VaultPaths != (vault.Paths{}) {
			if cfg.VaultPaths.Clerk.Details != "" {
				vd.DetailsPath = cfg.VaultPaths.Clerk.Details
			}
		}

		c.Setup(vd, *cfg.VaultHelper)
	}
	_, err := c.Build()
	if err != nil {
		return logs.Errorf("failed to build clerk: %v", err)
	}
	cfg.Clerk = *c
	return nil
}

func Resend(cfg *Config) error {
	r := resend.NewSystem()
	if cfg.VaultHelper != nil {
		vd := r.VaultDetails
		if cfg.VaultPaths != (vault.Paths{}) {
			if cfg.VaultPaths.Resend.Details != "" {
				vd.DetailsPath = cfg.VaultPaths.Resend.Details
			}
		}

		r.Setup(vd, *cfg.VaultHelper)
	}

	_, err := r.Build()
	if err != nil {
		return logs.Errorf("failed to build resend: %v", err)
	}
	cfg.Resend = *r
	return nil
}

func Bugfixes(cfg *Config) error {
	b := bugfixes.NewSystem()
	if cfg.VaultHelper != nil {
		vd := b.VaultDetails
		if cfg.VaultPaths != (vault.Paths{}) {
			if cfg.VaultPaths.BugFixes.Details != "" {
				vd.DetailsPath = cfg.VaultPaths.BugFixes.Details
			}
		}

		b.Setup(vd, *cfg.VaultHelper)
	}

	_, err := b.Build()
	if err != nil {
		return logs.Errorf("failed to build bugfixes: %v", err)
	}

	bf := &logs.BugFixes{}
	bf.Setup(b.AgentKey, b.AgentSecret)
	b.Logger = bf

	cfg.Bugfixes = *b
	return nil
}

func Flags(cfg *Config) error {
	f, err := flags.Build()
	if err != nil {
		return logs.Errorf("failed to build flags: %v", err)
	}
	cfg.Flags = *f
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
