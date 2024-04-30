package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type System struct {
	Client string `env:"KEYCLOAK_CLIENT" envDefault:"" json:"client,omitempty"`
	Secret string `env:"KEYCLOAK_SECRET" envDefault:"" json:"secret,omitempty"`
	Realm  string `env:"KEYCLOAK_REALM" envDefault:"" json:"realm,omitempty"`
	Host   string `env:"KEYCLOAK_HOSTNAME" envDefault:"" json:"host,omitempty"`
}

func NewKeycloak(client, secret, realm, host string) *System {
	return &System{
		Client: client,
		Secret: secret,
		Realm:  realm,
		Host:   host,
	}
}

func Build() (*System, error) {
	k := NewKeycloak("", "", "", "")
	if err := env.Parse(k); err != nil {
		return nil, logs.Errorf("keycloak: unable to parse keycloak: %v", err)
	}
	return k, nil
}

func (k *System) GetClient(ctx context.Context) (*gocloak.GoCloak, *gocloak.JWT, error) {
	client := gocloak.NewClient(k.Host)
	token, err := client.LoginClient(ctx, k.Client, k.Secret, k.Realm)
	if err != nil {
		return nil, nil, logs.Errorf("failed to login client: %v", err)
	}

	return client, token, nil
}
