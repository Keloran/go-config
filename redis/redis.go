package redis

import (
	"context"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type Details struct {
	Host     string `env:"REDIS_HOSTNAME" envDefault:"localhost:6379" json:"host,omitempty" vault:"redis-hostname"`
	Username string `env:"REDIS_USERNAME" json:"username,omitempty" vault:"redis-username"`
	Password string `env:"REDIS_PASSWORD" json:"password,omitempty" vault:"redis-password"`
	DB       int    `env:"REDIS_DB" envDefault:"0" json:"db,omitempty" vault:"redis-db"`
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
	r := &Details{}
	if err := env.Parse(r); err != nil {
		return nil, logs.Errorf("redis: unable to parse env: %v", err)
	}

	s.Details = *r
	return r, nil
}
