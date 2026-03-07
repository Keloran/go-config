package influx

import (
	"context"
	"strings"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
)

type VaultDetails struct {
	CredsPath   string `env:"INFLUX_VAULT_CREDS_PATH" envDefault:"secret/data/chewedfeed/influx"`
	DetailsPath string `env:"INFLUX_VAULT_DETAILS_PATH" envDefault:"secret/data/chewedfeed/details"`
}

type Details struct {
	Host   string `env:"INFLUX_HOSTNAME" envDefault:"http://db.chewed-k8s.net:8086"`
	Token  string `env:"INFLUX_TOKEN"`
	Bucket string `env:"INFLUX_BUCKET"`
	Org    string `env:"INFLUX_ORG"`
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
	s.VaultHelper = &vh
	s.VaultDetails = vd
}

func (s *System) Build() (*Details, error) {
	gen, err := s.buildGeneric()
	if err != nil {
		return nil, err
	}

	if s.VaultHelper != nil {
		return s.buildVault()
	}

	return gen, nil
}

func (s *System) buildGeneric() (*Details, error) {
	in := &Details{}
	if err := env.Parse(in); err != nil {
		return in, logs.Errorf("influx: unable to parse env: %v", err)
	}

	s.Details = *in
	return in, nil
}

func isVaultKeyNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}

func (s *System) buildVault() (*Details, error) {
	in := &Details{}
	vh := *s.VaultHelper

	// Get Credentials
	if err := vh.GetSecrets(s.VaultDetails.DetailsPath); err != nil {
		return in, logs.Errorf("influx: unable to get detail secrets: %v", err)
	}
	if vh.Secrets() == nil {
		return in, logs.Error("influx: unable to find credential secrets")
	}

	if s.Details.Token == "" {
		secret, err := vh.GetSecret("influx-token")
		if err != nil {
			return in, logs.Errorf("influx: unable to get token: %v", err)
		}
		in.Token = secret
	}

	if s.Details.Bucket == "" {
		secret, err := vh.GetSecret("influx-bucket")
		if err != nil {
			return in, logs.Errorf("influx: unable to get bucket: %v", err)
		}
		in.Bucket = secret
	}

	if s.Details.Org == "" {
		secret, err := vh.GetSecret("influx-org")
		if err != nil {
			return in, logs.Errorf("influx: unable to get org: %v", err)
		}
		in.Org = secret
	}

	// get the host based on the token, since host has a default in env
	if s.Details.Token == "" {
		secret, err := vh.GetSecret("influx-hostname")
		if err != nil {
			if !isVaultKeyNotFound(err) {
				return in, logs.Errorf("influx: unable to get hostname: %v", err)
			}
			secret = "http://db.chewed-k8s.net:8086"
		}
		in.Host = secret
	}

	s.Details = *in
	return in, nil
}
