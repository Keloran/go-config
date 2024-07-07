package authentik

import (
	"context"
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
	authentik "goauthentik.io/api/v3"
)

type VaultDetails struct {
	Address string
	Token   string

	DetailsPath string `env:"AUTHENTIK_DETAILS_PATH" envDefault:"secret/data/chewedfeed/details"`
}

type Details struct {
	Host   string `env:"AUTHENTIK_HOST" envDefault:"https://auth.chewedfeed.com"`
	Client string `env:"AUTHENTIK_CLIENT_ID" envDefault:""`
	Secret string `env:"AUTHENTIK_CLIENT_SECRET" envDefault:""`
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

	clientId, err := vh.GetSecret("authentik-client")
	if err != nil {
		return key, logs.Errorf("failed to get clientid: %v", err)
	}
	key.Client = clientId

	secret, err := vh.GetSecret("authentik-secret")
	if err != nil {
		return key, logs.Errorf("failed to get secret: %v", err)
	}
	key.Secret = secret

	host, err := vh.GetSecret("authentik-host")
	if err != nil {
		if err.Error() != fmt.Sprint("key: 'keycloak-host' not found") {
			return key, logs.Errorf("failed to get host: %v", err)
		}
		host = "https://auth.chewedfeed.com"
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

func (s *System) GetClient(ctx context.Context) (*authentik.APIClient, error) {
	cfg := authentik.NewConfiguration()
	cfg.Host = s.Host

	client := authentik.NewAPIClient(cfg)
	return client, nil
}
