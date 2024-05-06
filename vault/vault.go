package vault

import (
	"fmt"
	"strings"
	"time"

	"github.com/bugfixes/go-bugfixes/logs"

	"github.com/caarlos0/env/v8"
)

type VaultPath struct {
  Details     string
  Credentials string
}

type VaultPaths struct {
  Database VaultPath
  Keycloak VaultPath
  Mongo    VaultPath
  Rabbit   VaultPath
  Influx   VaultPath
}

// System is the vault config
type System struct {
	Host       string `env:"VAULT_HOST" envDefault:"localhost"`
	Port       string `env:"VAULT_PORT" envDefault:""`
	Token      string `env:"VAULT_TOKEN" envDefault:"root"`
	Address    string
	ExpireTime time.Time
}

func NewVault(address, token string) *System {
	return &System{
		Address: address,
		Token:   token,
	}
}

func Build() (*System, error) {
	v := NewVault("", "")

	if err := env.Parse(v); err != nil {
		return v, logs.Errorf("vault: %v", err)
	}

	if strings.HasPrefix(v.Host, "http") {
		v.Address = v.Host
	}

	if v.Port != "" {
		v.Address = fmt.Sprintf("%s:%s", v.Host, v.Port)
	}

	if v.Address == "" {
		v.Address = fmt.Sprintf("https://%s", v.Host)
	}

	return v, nil
}
