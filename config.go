package config

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
	ProjectProperties ProjectProperties
	ProjectConfig     interface{}
}

type BuildOption func(*Config) error

type subsystemConfigurator[T any] struct {
	name       string
	system     *T
	setupVault func(*T, vault.Paths, vaultHelper.VaultHelper)
	build      func(*T) error
	assign     func(*T)
}

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
	return buildSubsystem(cfg, subsystemConfigurator[postgres.System]{
		name:   "database",
		system: d,
		setupVault: func(d *postgres.System, paths vault.Paths, vh vaultHelper.VaultHelper) {
			vd := postgres.VaultDetails{}
			if paths.Database.Details != "" {
				vd.DetailsPath = paths.Database.Details
			}
			if paths.Database.Credentials != "" {
				vd.CredPath = paths.Database.Credentials
			}
			d.Setup(vd, vh)
		},
		build: func(d *postgres.System) error {
			_, err := d.Build()
			return err
		},
		assign: func(d *postgres.System) {
			cfg.Database = *d
		},
	})
}

func Mongo(cfg *Config) error {
	m := mongo.NewSystem()
	return buildSubsystem(cfg, subsystemConfigurator[mongo.System]{
		name:   "mongo",
		system: m,
		setupVault: func(m *mongo.System, paths vault.Paths, vh vaultHelper.VaultHelper) {
			vd := mongo.VaultDetails{}
			if paths.Mongo.Details != "" {
				vd.DetailsPath = paths.Mongo.Details
			}
			if paths.Mongo.Credentials != "" {
				vd.CredPath = paths.Mongo.Credentials
			}
			m.Setup(vd, vh)
		},
		build: func(m *mongo.System) error {
			_, err := m.Build()
			return err
		},
		assign: func(m *mongo.System) {
			cfg.Mongo = *m
		},
	})
}

func Keycloak(cfg *Config) error {
	k := keycloak.NewSystem()
	return buildSubsystem(cfg, subsystemConfigurator[keycloak.System]{
		name:   "keycloak",
		system: k,
		setupVault: func(k *keycloak.System, paths vault.Paths, vh vaultHelper.VaultHelper) {
			vd := keycloak.VaultDetails{}
			if paths.Keycloak.Details != "" {
				vd.DetailsPath = paths.Keycloak.Details
			}
			k.Setup(vd, vh)
		},
		build: func(k *keycloak.System) error {
			_, err := k.Build()
			return err
		},
		assign: func(k *keycloak.System) {
			cfg.Keycloak = *k
		},
	})
}

func Rabbit(cfg *Config) error {
	r := rabbit.NewSystem(&http.Client{})
	return buildSubsystem(cfg, subsystemConfigurator[rabbit.System]{
		name:   "rabbit",
		system: r,
		setupVault: func(r *rabbit.System, paths vault.Paths, vh vaultHelper.VaultHelper) {
			vd := rabbit.VaultDetails{}
			if paths.Rabbit.Details != "" {
				vd.DetailsPath = paths.Rabbit.Details
			}
			if paths.Rabbit.Credentials != "" {
				vd.CredPath = paths.Rabbit.Credentials
			}
			r.Setup(vd, vh)
		},
		build: func(r *rabbit.System) error {
			_, err := r.Build()
			return err
		},
		assign: func(r *rabbit.System) {
			cfg.Rabbit = *r
		},
	})
}

func Influx(cfg *Config) error {
	i := influx.NewSystem()
	return buildSubsystem(cfg, subsystemConfigurator[influx.System]{
		name:   "influx",
		system: i,
		setupVault: func(i *influx.System, paths vault.Paths, vh vaultHelper.VaultHelper) {
			vd := influx.VaultDetails{}
			if paths.Influx.Details != "" {
				vd.DetailsPath = paths.Influx.Details
			}
			i.Setup(vd, vh)
		},
		build: func(i *influx.System) error {
			_, err := i.Build()
			return err
		},
		assign: func(i *influx.System) {
			cfg.Influx = *i
		},
	})
}

func Clerk(cfg *Config) error {
	c := clerk.NewSystem()
	return buildSubsystem(cfg, subsystemConfigurator[clerk.System]{
		name:   "clerk",
		system: c,
		setupVault: func(c *clerk.System, paths vault.Paths, vh vaultHelper.VaultHelper) {
			vd := c.VaultDetails
			if paths.Clerk.Details != "" {
				vd.DetailsPath = paths.Clerk.Details
			}
			c.Setup(vd, vh)
		},
		build: func(c *clerk.System) error {
			_, err := c.Build()
			return err
		},
		assign: func(c *clerk.System) {
			cfg.Clerk = *c
		},
	})
}

func Resend(cfg *Config) error {
	r := resend.NewSystem()
	return buildSubsystem(cfg, subsystemConfigurator[resend.System]{
		name:   "resend",
		system: r,
		setupVault: func(r *resend.System, paths vault.Paths, vh vaultHelper.VaultHelper) {
			vd := r.VaultDetails
			if paths.Resend.Details != "" {
				vd.DetailsPath = paths.Resend.Details
			}
			r.Setup(vd, vh)
		},
		build: func(r *resend.System) error {
			_, err := r.Build()
			return err
		},
		assign: func(r *resend.System) {
			cfg.Resend = *r
		},
	})
}

func Bugfixes(cfg *Config) error {
	b := bugfixes.NewSystem()
	return buildSubsystem(cfg, subsystemConfigurator[bugfixes.System]{
		name:   "bugfixes",
		system: b,
		setupVault: func(b *bugfixes.System, paths vault.Paths, vh vaultHelper.VaultHelper) {
			vd := b.VaultDetails
			if paths.BugFixes.Details != "" {
				vd.DetailsPath = paths.BugFixes.Details
			}
			b.Setup(vd, vh)
		},
		build: func(b *bugfixes.System) error {
			_, err := b.Build()
			if err != nil {
				return err
			}

			logger := logs.Local()
			logger.Setup(b.AgentKey, b.AgentSecret)
			b.Logger = logger

			return nil
		},
		assign: func(b *bugfixes.System) {
			cfg.Bugfixes = *b
		},
	})
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

func buildSubsystem[T any](cfg *Config, subsystem subsystemConfigurator[T]) error {
	if cfg.VaultHelper != nil && subsystem.setupVault != nil {
		subsystem.setupVault(subsystem.system, cfg.VaultPaths, *cfg.VaultHelper)
	}

	if err := subsystem.build(subsystem.system); err != nil {
		return logs.Errorf("failed to build %s: %v", subsystem.name, err)
	}

	subsystem.assign(subsystem.system)
	return nil
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
