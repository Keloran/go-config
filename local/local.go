package local

import (
	"os"
	"strings"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

// System is the local config
type System struct {
	KeepLocal   bool `env:"BUGFIXES_LOCAL_ONLY" envDefault:"false"`
	Development bool `env:"DEVELOPMENT" envDefault:"false"`
	HTTPPort    int  `env:"HTTP_PORT" envDefault:"80"`
	GRPCPort    int  `env:"GRPC_PORT" envDefault:"3000"`
	EnvMap      map[string]string
}

func NewSystem(local, dev bool, http, grpc int) *System {
	return &System{
		KeepLocal:   local,
		Development: dev,
		HTTPPort:    http,
		GRPCPort:    grpc,
	}
}

func Build() (*System, error) {
	l := NewSystem(false, false, 80, 3000)
	if err := env.Parse(l); err != nil {
		return l, logs.Errorf("failed to parse local config: %v", err)
	}

	l.getAllEnvironment()

	return l, nil
}

func (l *System) getAllEnvironment() {
	envVars := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			envVars[pair[0]] = pair[1]
		}
	}
	l.EnvMap = envVars
}

func (l *System) GetValue(key string) string {
	if value, ok := l.EnvMap[key]; ok {
		return value
	}
	return ""
}
