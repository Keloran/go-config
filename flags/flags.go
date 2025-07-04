package flags

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type System struct {
	ProjectID     string `env:"FLAGS_PROJECT_ID" envDefault:""`
	AgentID       string `env:"FLAGS_AGENT_ID" envDefault:""`
	EnvironmentID string `env:"FLAGS_ENVIRONMENT_ID" envDefault:""`
}

func NewSystem() *System {
	return &System{}
}

func Build() (*System, error) {
	f := NewSystem()
	return f.buildLocal()
}

func (s *System) buildLocal() (*System, error) {
	if err := env.Parse(s); err != nil {
		return s, logs.Errorf("failed to parse flags config: %v", err)
	}

	return s, nil
}
