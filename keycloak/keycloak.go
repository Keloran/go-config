package keycloak

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type Keycloak struct {
	Client string `env:"KEYCLOAK_CLIENT" envDefault:"" json:"client,omitempty"`
	Secret string `env:"KEYCLOAK_SECRET" envDefault:"" json:"secret,omitempty"`
	Realm  string `env:"KEYCLOAK_REALM" envDefault:"" json:"realm,omitempty"`
	Host   string `env:"KEYCLOAK_HOSTNAME" envDefault:"" json:"host,omitempty"`
}

func NewKeycloak(client, secret, realm, host string) *Keycloak {
	return &Keycloak{
		Client: client,
		Secret: secret,
		Realm:  realm,
		Host:   host,
	}
}

func Build() (*Keycloak, error) {
	k := NewKeycloak("", "", "", "")
	if err := env.Parse(k); err != nil {
		return nil, logs.Errorf("keycloak: unable to parse keycloak: %v", err)
	}
	return k, nil
}
