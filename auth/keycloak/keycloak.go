package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type Details struct {
	Client string `env:"KEYCLOAK_CLIENT" envDefault:"" json:"client,omitempty" vault:"keycloak-client"`
	Secret string `env:"KEYCLOAK_SECRET" envDefault:"" json:"secret,omitempty" vault:"keycloak-secret"`
	Realm  string `env:"KEYCLOAK_REALM" envDefault:"" json:"realm,omitempty" vault:"keycloak-realm"`
	Host   string `env:"KEYCLOAK_HOSTNAME" envDefault:"https://keys.chewedfeed.com" json:"host,omitempty" vault:"keycloak-host"`
}

type System struct {
	Context context.Context

	Details
}

func NewSystem() *System {
	return &System{
		Context: context.Background(),
	}
}

func (s *System) Build() (*Details, error) {
	return s.buildGeneric()
}

func (s *System) buildGeneric() (*Details, error) {
	key := &Details{}
	if err := env.Parse(key); err != nil {
		return nil, logs.Errorf("keycloak: unable to parse env: %v", err)
	}

	s.Details = *key
	return key, nil
}

func (s *System) GetClient(ctx context.Context) (*gocloak.GoCloak, *gocloak.JWT, error) {
	client := gocloak.NewClient(s.Host)
	token, err := client.LoginClient(ctx, s.Client, s.Secret, s.Realm)
	if err != nil {
		return nil, nil, logs.Errorf("keycloak: unable to login client: %v", err)
	}

	return client, token, nil
}
