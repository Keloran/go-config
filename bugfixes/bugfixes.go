package bugfixes

import (
	"context"
	"strings"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type Details struct {
	Server      string `env:"BUGFIXES_SERVER" envDefault:"https://api.bugfix.es/v1" vault:"bugfixes-server"`
	AgentKey    string `env:"BUGFIXES_AGENT_KEY" vault:"bugfixes-agentid"`
	AgentSecret string `env:"BUGFIXES_AGENT_SECRET" vault:"bugfixes-secret"`
}

type System struct {
	Context context.Context

	Details
	Logger *logs.BugFixes
}

func NewSystem() *System {
	return &System{
		Context: context.Background(),
	}
}

func (s *System) Build() (*Details, error) {
	bf := &Details{}
	if err := env.Parse(bf); err != nil {
		return bf, logs.Errorf("bugfixes: unable to parse env: %v", err)
	}

	s.Details = *bf
	if err := s.Validate(); err != nil {
		return bf, err
	}

	return bf, nil
}

func (s *System) Validate() error {
	if !strings.HasPrefix(s.Server, "http") {
		return logs.Error("bugfixes: unable to use server without protocol")
	}
	return nil
}
