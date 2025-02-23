package resend

import (
	"context"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
)

type Details struct {
	Key string `env:"RESEND_KEY" envDefault:""`
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
	resend := &Details{}
	vh := *s.VaultHelper

	if err := vh.GetSecrets(s.VaultDetails.DetailsPath); err != nil {
		return resend, err
	}
	if vh.Secrets() == nil {
		return resend, nil
	}

	if resend.Key == "" {
		secret, err := vh.GetSecret("resend_key")
		if err != nil {
			return resend, err
		}
		resend.Key = secret
	}

	s.Details = *resend
	return resend, nil
}
