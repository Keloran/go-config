package local

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

// Local is the local config
type Local struct {
	KeepLocal   bool `env:"BUGFIXES_LOCAL_ONLY" envDefault:"false"`
	Development bool `env:"DEVELOPMENT" envDefault:"false"`
	HTTPPort    int  `env:"HTTP_PORT" envDefault:"80"`
	GRPCPort    int  `env:"GRPC_PORT" envDefault:"3000"`
}

func Build() (*Local, error) {
	l := &Local{}
	if err := env.Parse(l); err != nil {
		return l, logs.Errorf("failed to parse local config: %v", err)
	}
	return l, nil
}
