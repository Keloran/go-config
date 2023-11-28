package rabbit

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type Rabbit struct {
	Host     string `env:"RABBIT_HOSTNAME" envDefault:"" json:"host,omitempty"`
	Port     string `env:"RABBIT_PORT" envDefault:"" json:"port,omitempty"`
	Username string `env:"RABBIT_USERNAME" envDefault:"" json:"username,omitempty"`
	Password string `env:"RABBIT_PASSWORD" envDefault:"" json:"password,omitempty"`
	VHost    string `env:"RABBIT_VHOST" envDefault:"" json:"vhost,omitempty"`
}

func NewRabbit(host, port, username, password, vhost string) *Rabbit {
	return &Rabbit{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		VHost:    vhost,
	}
}

func Build() (*Rabbit, error) {
	r := &Rabbit{}
	if err := env.Parse(r); err != nil {
		return nil, logs.Errorf("rabbit: unable to parse rabbit: %v", err)
	}
	return r, nil
}
