package influx

import (
	"context"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type Details struct {
	Host   string `env:"INFLUX_HOSTNAME" envDefault:"http://db.chewed-k8s.net:8086" vault:"influx-hostname"`
	Token  string `env:"INFLUX_TOKEN" vault:"influx-token"`
	Bucket string `env:"INFLUX_BUCKET" vault:"influx-bucket"`
	Org    string `env:"INFLUX_ORG" vault:"influx-org"`
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
	in := &Details{}
	if err := env.Parse(in); err != nil {
		return in, logs.Errorf("influx: unable to parse env: %v", err)
	}

	s.Details = *in
	return in, nil
}
