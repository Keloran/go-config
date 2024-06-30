package keycloak

import (
	"context"
	"fmt"
	vaultHelper "github.com/keloran/vault-helper"

	"github.com/Nerzal/gocloak/v13"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type VaultDetails struct {
	Address string
	Token   string

	DetailsPath string `env:"KEYCLOAK_VAULT_DETAIL_PATH" envDefault:"secret/data/chewedfeed/details"`

	Exclusive bool
}

type Details struct {
	Client string `env:"KEYCLOAK_CLIENT" envDefault:"" json:"client,omitempty"`
	Secret string `env:"KEYCLOAK_SECRET" envDefault:"" json:"secret,omitempty"`
	Realm  string `env:"KEYCLOAK_REALM" envDefault:"" json:"realm,omitempty"`
	Host   string `env:"KEYCLOAK_HOSTNAME" envDefault:"https://keys.chewedfeed.com" json:"host,omitempty"`
}

type System struct {
	Context context.Context

	Details

	VaultDetails
	VaultHelper *vaultHelper.VaultHelper
}

func NewSystem() *System {
	return &System{
		Context: context.Background(),
	}
}

func (s *System) Setup(vd VaultDetails, vh vaultHelper.VaultHelper) {
	s.VaultDetails = vd
	s.VaultHelper = &vh
}

func (s *System) Build() (*Details, error) {
	if s.VaultHelper != nil {
		return s.buildVault()
	}

	return s.buildGeneric()
}

func (s *System) buildVault() (*Details, error) {
	key := &Details{}
	vh := *s.VaultHelper

	if err := vh.GetSecrets(s.VaultDetails.DetailsPath); err != nil {
		return key, logs.Errorf("failed to get detail secrets from vault: %v", err)
	}
	if vh.Secrets() == nil {
		return key, logs.Error("no keycloak secrets found")
	}

	clientId, err := vh.GetSecret("keycloak-client")
	if err != nil {
		return key, logs.Errorf("failed to get clientid: %v", err)
	}
	key.Client = clientId

	secret, err := vh.GetSecret("keycloak-secret")
	if err != nil {
		return key, logs.Errorf("failed to get secret: %v", err)
	}
	key.Secret = secret

	realm, err := vh.GetSecret("keycloak-realm")
	if err != nil {
		return key, logs.Errorf("failed to get realm: %v", err)
	}
	key.Realm = realm

	host, err := vh.GetSecret("keycloak-host")
	if err != nil {
		if err.Error() != fmt.Sprint("key: 'keycloak-host' not found") {
			return key, logs.Errorf("failed to get host: %v", err)
		}
		host = "https://keys.chewedfeed.com"
	}
	key.Host = host

	s.Details = *key
	return key, nil
}

func (s *System) buildGeneric() (*Details, error) {
	key := &Details{}
	if err := env.Parse(key); err != nil {
		return nil, logs.Errorf("failed to build generic: %v", err)
	}

	s.Details = *key
	return key, nil
}

func (s *System) GetClient(ctx context.Context) (*gocloak.GoCloak, *gocloak.JWT, error) {
	client := gocloak.NewClient(k.Host)
	token, err := client.LoginClient(ctx, k.Client, k.Secret, k.Realm)
	if err != nil {
		return nil, nil, logs.Errorf("failed to login client: %v", err)
	}

	return client, token, nil
}
