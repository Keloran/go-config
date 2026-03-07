package bugfixes

import (
	"context"
	"strings"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
)

type Details struct {
	Server      string `env:"BUGFIXES_SERVER" envDefault:"https://api.bugfix.es/v1"`
	AgentKey    string `env:"BUGFIXES_AGENT_KEY"`
	AgentSecret string `env:"BUGFIXES_AGENT_SECRET"`
}

type System struct {
	Context context.Context

	Details
	Logger *logs.BugFixes

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
	bf := &Details{}
	if err := env.Parse(bf); err != nil {
		return bf, logs.Errorf("bugfixes: unable to parse env: %v", err)
	}

	if !strings.HasPrefix(bf.Server, "http") {
		return bf, logs.Error("bugfixes: unable to use server without protocol")
	}

	s.Details = *bf
	return bf, nil
}

func isVaultKeyNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}

func (s *System) buildVault() (*Details, error) {
	bf := &Details{}
	vh := *s.VaultHelper

	if err := vh.GetSecrets(s.VaultDetails.DetailsPath); err != nil {
		return bf, logs.Errorf("bugfixes: unable to get detail secrets: %v", err)
	}

	if vh.Secrets() == nil {
		return bf, logs.Error("bugfixes: unable to find secrets")
	}

	if s.Details.AgentKey == "" {
		secret, err := vh.GetSecret("bugfixes-agentid")
		if err != nil {
			return bf, logs.Errorf("bugfixes: unable to get agent id: %v", err)
		}
		bf.AgentKey = secret
	} else {
		bf.AgentKey = s.Details.AgentKey
	}

	if s.Details.AgentSecret == "" {
		secret, err := vh.GetSecret("bugfixes-secret")
		if err != nil {
			return bf, logs.Errorf("bugfixes: unable to get secret: %v", err)
		}
		bf.AgentSecret = secret
	} else {
		bf.AgentSecret = s.Details.AgentSecret
	}

	if s.Details.Server == "" {
		secret, err := vh.GetSecret("bugfixes-server")
		if err != nil {
			if !isVaultKeyNotFound(err) {
				return bf, logs.Errorf("bugfixes: unable to get server: %v", err)
			}
			secret = "https://api.bugfix.es/v1"
		}
		bf.Server = secret
	} else {
		bf.Server = s.Details.Server
	}

	if !strings.HasPrefix(bf.Server, "http") {
		return bf, logs.Error("bugfixes: unable to use server without protocol")
	}

	s.Details = *bf

	return bf, nil
}
