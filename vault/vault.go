package vault

import (
	"fmt"
	vaultHelper "github.com/keloran/vault-helper"
	"strings"
	"time"

	"github.com/bugfixes/go-bugfixes/logs"

	"github.com/caarlos0/env/v8"
)

type Path struct {
	Details     string
  Local string
	Credentials string
}

type Paths struct {
	Database Path
	Keycloak Path
	Mongo    Path
	Rabbit   Path
	Influx   Path
  BugFixes Path
}

// System is the vault config
type System struct {
	Host       string `env:"VAULT_HOST" envDefault:"vault.vault"`
	Port       string `env:"VAULT_PORT" envDefault:""`
	Token      string `env:"VAULT_TOKEN" envDefault:"root"`
	Address    string
	ExpireTime time.Time
}

func NewSystem(address, token string) *System {
	return &System{
		Address: address,
		Token:   token,
	}
}

func Build() (*System, vaultHelper.VaultHelper, error) {
	v := NewSystem("", "")

	if err := env.Parse(v); err != nil {
		return v, nil, logs.Errorf("vault: %v", err)
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

	vh := vaultHelper.NewVault(v.Address, v.Token)

	return v, vh, nil
}
