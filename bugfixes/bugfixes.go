package bugfixes

import (
	"context"
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
	"strings"
)

type Details struct {
	Server      string `env:"BUGIXES_SERVER" envDefault:"https://api.bugfix.es/v1"`
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
	if s.VaultHelper != nil {
		return s.buildVault()
	}

	return s.buildGeneric()
}

func (s *System) buildGeneric() (*Details, error) {
	bf := &Details{}
	if err := env.Parse(bf); err != nil {
		return bf, logs.Errorf("failed to parse bugfixes env: %v", err)
	}

	if !strings.HasPrefix(bf.Server, "http") {
		return bf, logs.Error("needs the protocol for server")
	}

	s.Details = *bf
	return bf, nil
}

func (s *System) buildVault() (*Details, error) {
	bf := &Details{}
	vh := *s.VaultHelper

	if err := vh.GetSecrets(s.VaultDetails.DetailsPath); err != nil {
		return bf, logs.Errorf("faiuled to get local bugfix details: %v", err)
	}

	if vh.Secrets() == nil {
		return bf, logs.Error("no bugfixes secrets found")
	}

	agent, err := vh.GetSecret("bugfixes-agentid")
	if err != nil {
		return bf, logs.Errorf("failed to get agentid: %v", err)
	}
	bf.AgentKey = agent

	secret, err := vh.GetSecret("bugfixes-secret")
	if err != nil {
		return bf, logs.Errorf("failed to get secret: %v", err)
	}
	bf.AgentSecret = secret

	server, err := vh.GetSecret("bugfixes-server")
	if err != nil {
		if err.Error() != fmt.Sprint("key: 'bugfixes-server' not found") {
			return bf, logs.Errorf("failed to get server: %v", err)
		}
		server = "https://api.bugfix.es/v1"
	}
	bf.Server = server
	if !strings.HasPrefix(server, "http") {
		return bf, logs.Error("needs the protocol for server")
	}

	s.Details = *bf

	return bf, nil
}
