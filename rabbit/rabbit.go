package rabbit

import (
	"context"
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaulthelper "github.com/keloran/vault-helper"
	"net/http"
	"time"
)

type VaultDetails struct {
	Address    string
	Path       string `env:"RABBIT_VAULT_PATH" envDefault:"secret/data/chewedfeed/rabbitmq"`
	Token      string
	ExpireTime time.Time
}

type Rabbit struct {
	Host           string `env:"RABBIT_HOSTNAME" envDefault:"" json:"host,omitempty"`
	ManagementHost string `env:"RABBIT_MANAGEMENT_HOSTNAME" envDefault:"" json:"management_host,omitempty"`
	Port           int    `env:"RABBIT_PORT" envDefault:"" json:"port,omitempty"`
	Username       string `env:"RABBIT_USERNAME" envDefault:"" json:"username,omitempty"`
	Password       string `env:"RABBIT_PASSWORD" envDefault:"" json:"password,omitempty"`
	VHost          string `env:"RABBIT_VHOST" envDefault:"" json:"vhost,omitempty"`
	Queue          string `env:"RABBIT_QUEUE" envDefault:"" json:"queue,omitempty"`

	VaultDetails
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

	r.VaultDetails.ExpireTime = time.Now().Add(time.Duration(vh.LeaseDuration()) * time.Second)

	return r, nil
}

func GetRabbitQueue(ctx context.Context, r Rabbit) (interface{}, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if time.Now().Unix() > r.VaultDetails.ExpireTime.Unix() {
		rb, err := Build(r.VaultDetails, vaulthelper.NewVault(r.VaultDetails.Address, r.VaultDetails.Token))
		if err != nil {
			return nil, logs.Errorf("rabbit: unable to build rabbit: %v", err)
		}
		r = *rb
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/queues/%s/%s/get", r.Host, r.VHost, r.Queue), nil)
	if err != nil {
		return nil, logs.Errorf("rabbit: unable to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(r.Username, r.Password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, logs.Errorf("rabbit: unable to get queue: %v", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			_ = logs.Errorf("rabbit: unable to close response: %v", err)
		}
	}()

	if res.StatusCode != 200 {
		return nil, logs.Errorf("rabbit: unable to get queue: %v", res.Status)
	}

	return res, nil
}
