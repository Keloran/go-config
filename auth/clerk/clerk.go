package clerk

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
)

type Details struct {
	Key       string `env:"CLERK_SECRET_KEY" envDefault:""`
	PublicKey string `env:"NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY" envDefault:""`
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
		return nil, err
	}

	s.Details = *clerk
	return clerk, nil
}

func (s *System) buildVault() (*Details, error) {
	clerk := &Details{}
	vh := *s.VaultHelper

	if err := vh.GetSecrets(s.VaultDetails.DetailsPath); err != nil {
		return clerk, err
	}
	if vh.Secrets() == nil {
		return clerk, nil
	}

	if s.Details.Key == "" {
		secret, err := vh.GetSecret("clerk-key")
		if err != nil {
			return clerk, err
		}
		clerk.Key = secret
	} else {
		clerk.Key = s.Details.Key
	}

	if s.Details.PublicKey == "" {
		secret, err := vh.GetSecret("clerk-public-key")
		if err != nil {
			if err.Error() != fmt.Sprint("key: 'clerk-public-key' not found") {
				return clerk, err
			}
		}
		clerk.PublicKey = secret
	} else {
		clerk.PublicKey = s.Details.PublicKey
	}

	s.Details = *clerk
	return clerk, nil
}
