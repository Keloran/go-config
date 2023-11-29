package rabbit

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaulthelper "github.com/keloran/vault-helper"
)

type VaultDetails struct {
	Address string
	Path    string `env:"RABBIT_VAULT_PATH" envDefault:"secret/data/chewedfeed/rabbitmq"`
	Token   string
}

type Rabbit struct {
	Host           string `env:"RABBIT_HOSTNAME" envDefault:"" json:"host,omitempty"`
	ManagementHost string `env:"RABBIT_MANAGEMENT_HOSTNAME" envDefault:"" json:"management_host,omitempty"`
	Port           int    `env:"RABBIT_PORT" envDefault:"" json:"port,omitempty"`
	Username       string `env:"RABBIT_USERNAME" envDefault:"" json:"username,omitempty"`
	Password       string `env:"RABBIT_PASSWORD" envDefault:"" json:"password,omitempty"`
	VHost          string `env:"RABBIT_VHOST" envDefault:"" json:"vhost,omitempty"`
}

func NewRabbit(port int, host, username, password, vhost, management string) *Rabbit {
	return &Rabbit{
		Host:           host,
		Port:           port,
		Username:       username,
		Password:       password,
		VHost:          vhost,
		ManagementHost: management,
	}
}

func Setup(vaultAddress, vaultToken string) VaultDetails {
	return VaultDetails{
		Address: vaultAddress,
		Token:   vaultToken,
	}
}

func Build(vd VaultDetails, vh vaulthelper.VaultHelper) (*Rabbit, error) {
	r := NewRabbit(0, "", "", "", "", "")

	if err := env.Parse(r); err != nil {
		return nil, logs.Errorf("rabbit: unable to parse rabbit: %v", err)
	}

	// env rather than vault
	if r.Username != "" && r.Password != "" {
		return r, nil
	}

	if err := vh.GetSecrets(vd.Path); err != nil {
		return nil, logs.Errorf("rabbit: unable to get secrets: %v", err)
	}
	if vh.Secrets() == nil {
		return nil, logs.Error("rabbit: no secrets found")
	}
	if username, err := vh.GetSecret("username"); err == nil {
		r.Username = username
	} else {
		return nil, logs.Errorf("rabbit: unable to get username: %v", err)
	}
	if password, err := vh.GetSecret("password"); err == nil {
		r.Password = password
	} else {
		return nil, logs.Errorf("rabbit: unable to get password: %v", err)
	}
	if vhost, err := vh.GetSecret("vhost"); err == nil {
		r.VHost = vhost
	} else {
		return nil, logs.Errorf("rabbit: unable to get vhost: %v", err)
	}

	return r, nil
}
