package resend

import (
	"context"

	"github.com/caarlos0/env/v8"
)

type Details struct {
	Key string `env:"RESEND_KEY" envDefault:"" vault:"resend_key"`
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
	resend := &Details{}
	if err := env.Parse(resend); err != nil {
		return nil, err
	}

	s.Details = *resend
	return resend, nil
}
