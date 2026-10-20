package clerk

import (
	"context"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type Details struct {
	Key       string `env:"CLERK_SECRET_KEY" envDefault:"" vault:"clerk-key"`
	PublicKey string `env:"NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY" envDefault:"" vault:"clerk-public-key"`
	DevUser   string `env:"CLERK_DEV_USER" envDefault:"" vault:"clerk-dev-user"`
}

type System struct {
	Context context.Context

	Details
}

func NewSystem() *System {
	return &System{
		Context: context.Background(),
	}
}

func (s *System) Build() (*Details, error) {
	clerk := &Details{}
	if err := env.Parse(clerk); err != nil {
		return nil, logs.Errorf("clerk: unable to parse env: %v", err)
	}

	s.Details = *clerk
	return clerk, nil
}
