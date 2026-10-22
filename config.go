package config

import (
	"errors"
	"github.com/keloran/go-config/v2/auth/clerk"
	"github.com/keloran/go-config/v2/flags"
	"github.com/keloran/go-config/v2/notify/resend"
	"net/http"
	"os"

	bugfixesSDK "github.com/bugfixes/go-bugfixes"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/joho/godotenv"
	"github.com/keloran/go-config/v2/auth/keycloak"
	"github.com/keloran/go-config/v2/bugfixes"
	"github.com/keloran/go-config/v2/database/mongo"
	"github.com/keloran/go-config/v2/database/postgres"
	"github.com/keloran/go-config/v2/influx"
	"github.com/keloran/go-config/v2/local"
	"github.com/keloran/go-config/v2/onepassword"
	"github.com/keloran/go-config/v2/rabbit"
	"github.com/keloran/go-config/v2/redis"
	"github.com/keloran/go-config/v2/vault"

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
	Redis    redis.System
	Bugfixes bugfixes.System
	Clerk    clerk.System
	Resend   resend.System
	Flags    flags.System

	// Project level properties
	ProjectProperties ProjectProperties
	ProjectConfig     interface{}
}

type BuildOption func(*Config) error

// OnePassword loads unset environment variables from 1Password before later
// BuildOptions run. fields maps "item.field" (or an op:// reference) to an
// environment variable name. Without OP_SERVICE_ACCOUNT_TOKEN it is a no-op.
func OnePassword(vaultName string, fields map[string]string) BuildOption {
	return func(*Config) error {
		return onepassword.Load(vaultName, fields)
	}
}

// OnePasswordRequired is like OnePassword, but fails when
// OP_SERVICE_ACCOUNT_TOKEN is missing.
func OnePasswordRequired(vaultName string, fields map[string]string) BuildOption {
	return func(*Config) error {
		return onepassword.LoadRequired(vaultName, fields)
	}
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
		return logs.Errorf("config: unable to build local: %v", err)
	}

	cfg.Local = *l

	return nil
}

func Vault(cfg *Config) error {
	v, vh, err := vault.Build()
	if err != nil {
		return logs.Errorf("config: unable to build vault: %v", err)
	}

	cfg.Vault = *v
	cfg.VaultHelper = &vh

	return nil
}

// fillVault fills any `vault` tagged fields still unset by env, only when vault is configured
func fillVault(cfg *Config, name string, p vault.Path, details any) error {
	if cfg.VaultHelper == nil {
		return nil
	}
	if err := vault.Fill(*cfg.VaultHelper, p, details); err != nil {
		return logs.Errorf("%s: unable to fill from vault: %v", name, err)
	}
	return nil
}

func credsPath(p vault.Path, def string) vault.Path {
	if p.Credentials == "" {
		p.Credentials = def
	}
	return p
}

// Database deprecated: use Postgres instead
func Database(cfg *Config) error {
	return Postgres(cfg)
}

func Postgres(cfg *Config) error {
	d := postgres.NewSystem()
	if _, err := d.Build(); err != nil {
		return logs.Errorf("database: unable to build: %v", err)
	}
	if err := fillVault(cfg, "database", credsPath(cfg.VaultPaths.Database, "secret/data/chewedfeed/postgres"), &d.Details); err != nil {
		return err
	}
	cfg.Database = *d
	return nil
}

func Mongo(cfg *Config) error {
	m := mongo.NewSystem()
	if _, err := m.Build(); err != nil {
		return logs.Errorf("mongo: unable to build: %v", err)
	}
	if err := fillVault(cfg, "mongo", credsPath(cfg.VaultPaths.Mongo, "secret/data/chewedfeed/mongo"), &m.Details); err != nil {
		return err
	}
	cfg.Mongo = *m
	return nil
}

func Keycloak(cfg *Config) error {
	k := keycloak.NewSystem()
	if _, err := k.Build(); err != nil {
		return logs.Errorf("keycloak: unable to build: %v", err)
	}
	if err := fillVault(cfg, "keycloak", cfg.VaultPaths.Keycloak, &k.Details); err != nil {
		return err
	}
	cfg.Keycloak = *k
	return nil
}

func Rabbit(cfg *Config) error {
	r := rabbit.NewSystem(&http.Client{})
	if _, err := r.Build(); err != nil {
		return logs.Errorf("rabbit: unable to build: %v", err)
	}
	if err := fillVault(cfg, "rabbit", cfg.VaultPaths.Rabbit, &r.Details); err != nil {
		return err
	}
	cfg.Rabbit = *r
	return nil
}

func Influx(cfg *Config) error {
	i := influx.NewSystem()
	if _, err := i.Build(); err != nil {
		return logs.Errorf("influx: unable to build: %v", err)
	}
	if err := fillVault(cfg, "influx", cfg.VaultPaths.Influx, &i.Details); err != nil {
		return err
	}
	cfg.Influx = *i
	return nil
}

func Redis(cfg *Config) error {
	r := redis.NewSystem()
	if _, err := r.Build(); err != nil {
		return logs.Errorf("redis: unable to build: %v", err)
	}
	if err := fillVault(cfg, "redis", cfg.VaultPaths.Redis, &r.Details); err != nil {
		return err
	}
	cfg.Redis = *r
	return nil
}

func Clerk(cfg *Config) error {
	c := clerk.NewSystem()
	if _, err := c.Build(); err != nil {
		return logs.Errorf("clerk: unable to build: %v", err)
	}
	if err := fillVault(cfg, "clerk", cfg.VaultPaths.Clerk, &c.Details); err != nil {
		return err
	}
	cfg.Clerk = *c
	return nil
}

func Resend(cfg *Config) error {
	r := resend.NewSystem()
	if _, err := r.Build(); err != nil {
		return logs.Errorf("resend: unable to build: %v", err)
	}
	if err := fillVault(cfg, "resend", cfg.VaultPaths.Resend, &r.Details); err != nil {
		return err
	}
	cfg.Resend = *r
	return nil
}

func Bugfixes(cfg *Config) error {
	b := bugfixes.NewSystem()
	if _, err := b.Build(); err != nil {
		return logs.Errorf("bugfixes: unable to build: %v", err)
	}
	if err := fillVault(cfg, "bugfixes", cfg.VaultPaths.BugFixes, &b.Details); err != nil {
		return err
	}
	if err := b.Validate(); err != nil {
		return logs.Errorf("bugfixes: unable to build: %v", err)
	}

	logger := logs.Local()
	logger.Setup(b.AgentKey, b.AgentSecret)
	logger.SetConfig(bugfixesSDK.Config{Server: b.Server, Service: b.Service})
	b.Logger = logger

	cfg.Bugfixes = *b
	return nil
}

func Flags(cfg *Config) error {
	f, err := flags.Build()
	if err != nil {
		return logs.Errorf("flags: unable to build: %v", err)
	}
	cfg.Flags = *f
	return nil
}

func Build(opts ...BuildOption) (*Config, error) {
	cfg := &Config{}

	if err := cfg.Build(opts...); err != nil {
		return nil, logs.Errorf("config: unable to build: %v", err)
	}

	return cfg, nil
}

func BuildLocal(opts ...BuildOption) (*Config, error) {
	cfg := &Config{}

	if err := cfg.Build(opts...); err != nil {
		return nil, logs.Errorf("config: unable to build local: %v", err)
	}

	return cfg, nil
}

func BuildLocalVH(mockVault vaultHelper.VaultHelper, opts ...BuildOption) (*Config, error) {
	cfg := &Config{
		VaultHelper: &mockVault,
	}

	if err := cfg.Build(opts...); err != nil {
		return nil, logs.Errorf("config: unable to build with vault: %v", err)
	}

	return cfg, nil
}

func (c *Config) Build(opts ...BuildOption) error {
	if err := loadDotEnv(); err != nil {
		return logs.Errorf("config: unable to load .env: %v", err)
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return logs.Errorf("config: unable to apply option: %v", err)
		}
	}

	return nil
}

func loadDotEnv() error {
	if _, err := os.Stat(".env"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return logs.Errorf("failed to stat .env: %v", err)
	}

	if err := godotenv.Load(".env"); err != nil {
		return logs.Errorf("failed to load .env: %v", err)
	}

	return nil
}

type ProjectConfigurator interface {
	Build(*Config) error
}

func WithProjectConfigurator(pc ProjectConfigurator) BuildOption {
	return func(c *Config) error {
		if err := pc.Build(c); err != nil {
			return logs.Errorf("config: unable to apply project configurator: %v", err)
		}

		return nil
	}
}
