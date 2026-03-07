package clerk

import (
	"context"
	"strings"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
)

type Details struct {
	Key       string `env:"CLERK_SECRET_KEY" envDefault:""`
	PublicKey string `env:"NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY" envDefault:""`
	DevUser   string `env:"CLERK_DEV_USER" envDefault:""`
}

type System struct {
	Context context.Context

	Details

	VaultDetails vaultHelper.VaultDetails
	VaultHelper  *vaultHelper.VaultHelper
}

func NewSystem() *System {
	return &System{
		Context: context.Background(),
	}
}

func (s *System) Setup(vd vaultHelper.VaultDetails, vh vaultHelper.VaultHelper) {
	s.VaultDetails = vd
	s.VaultHelper = &vh
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
	clerk := &Details{}
	if err := env.Parse(clerk); err != nil {
		return nil, logs.Errorf("clerk: unable to parse env: %v", err)
	}

	s.Details = *clerk
	return clerk, nil
}

func isVaultKeyNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}

func (s *System) buildVault() (*Details, error) {
	clerk := &Details{}
	vh := *s.VaultHelper

	if err := vh.GetSecrets(s.VaultDetails.DetailsPath); err != nil {
		return clerk, logs.Errorf("clerk: unable to get detail secrets: %v", err)
	}
	if vh.Secrets() == nil {
		return clerk, nil
	}

	if s.Details.Key == "" {
		secret, err := vh.GetSecret("clerk-key")
		if err != nil {
			return clerk, logs.Errorf("clerk: unable to get key: %v", err)
		}
		clerk.Key = secret
	} else {
		clerk.Key = s.Details.Key
	}

	if s.Details.PublicKey == "" {
		secret, err := vh.GetSecret("clerk-public-key")
		if err != nil {
			if !isVaultKeyNotFound(err) {
				return clerk, logs.Errorf("clerk: unable to get public key: %v", err)
			}
		}
		clerk.PublicKey = secret
	} else {
		clerk.PublicKey = s.Details.PublicKey
	}

	if s.Details.DevUser == "" {
		secret, err := vh.GetSecret("clerk-dev-user")
		if err != nil {
			if !isVaultKeyNotFound(err) {
				return clerk, logs.Errorf("clerk: unable to get dev user: %v", err)
			}
		}
		clerk.DevUser = secret
	} else {
		clerk.DevUser = s.Details.DevUser
	}

	s.Details = *clerk
	return clerk, nil
}
