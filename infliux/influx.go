package infliux

import (
	"context"
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
)

type VaultDetails struct {
	CredsPath   string `env:"INFLUX_VAULT_CREDS_PATH" envDefault:"secret/data/chewedfeed/influx"`
	DetailsPath string `env:"INFLUX_VAULT_DETAILS_PATH" envDefault:"secret/data/chewedfeed/details"`
}

type Details struct {
	Host     string `env:"INFLUX_HOSTNAME" envDefault:"http://db.chewed-k8s.net:8086"`
	User     string `env:"INFLUX_USERNAME"`
	Password string `env:"INFLUX_PASSWORD"`
	Bucket   string `env:"INFLUX_BUCKET"`
	Org      string `env:"INFLUX_ORG"`
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
	if s.VaultHelper != nil {
		return s.buildVault()
	}

	return s.buildGeneric()
}

func (s *System) buildGeneric() (*Details, error) {
	in := &Details{}
	if err := env.Parse(in); err != nil {
		return in, logs.Errorf("failed to parse influx env: %v", err)
	}

	s.Details = *in
	return in, nil
}

func (s *System) buildVault() (*Details, error) {
	in := &Details{}
	vh := *s.VaultHelper

	// Get Credentials
	if err := vh.GetSecrets(s.VaultDetails.DetailsPath); err != nil {
		return in, logs.Errorf("failed to get detail secrets from vault: %v", err)
	}
	if vh.Secrets() == nil {
		return in, logs.Error("no influx cred serets found")
	}

	username, err := vh.GetSecret("influx-username")
	if err != nil {
		return in, logs.Errorf("failed to get username: %v", err)
	}
	in.User = username

	password, err := vh.GetSecret("influx-password")
	if err != nil {
		return in, logs.Errorf("failed to get password: %v", err)
	}
	in.Password = password

	bucket, err := vh.GetSecret("influx-bucket")
	if err != nil {
		return in, logs.Errorf("failed to get bucket: %v", err)
	}
	in.Bucket = bucket

	org, err := vh.GetSecret("influx-org")
	if err != nil {
		return in, logs.Errorf("failed to get org: %v", err)
	}
	in.Org = org

	host, err := vh.GetSecret("influx-hostname")
	if err != nil {
		if err.Error() != fmt.Sprint("key not found") && err.Error() != fmt.Sprint("key not found: influx-hostname") {
			return in, logs.Errorf("failed to get host: %v", err)
		}
		host = "http://db.chewed-k8s.net:8086"
	}
	in.Host = host

	s.Details = *in
	return in, nil
}
