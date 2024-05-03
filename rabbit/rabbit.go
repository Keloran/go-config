package rabbit

import (
	"context"
	"fmt"
	"net/http"
  "strconv"
  "time"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaulthelper "github.com/keloran/vault-helper"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type VaultHelper interface {
	GetSecrets(path string) error
	GetSecret(key string) (string, error)
	Secrets() []vaulthelper.KVSecret
	LeaseDuration() int
}

type VaultDetails struct {
	Address    string
	Token      string

  Path       string `env:"RABBIT_VAULT_PATH" envDefault:"secret/data/chewedfeed/rabbitmq"`
  CredPath   string `env:"RABBIT_VAULT_CREDS_PATH" envDefault:"secrets/data/chewedfeed/rabbitmq"`
  DetailsPath string `env:"RABBIT_VAULT_DETAILS_PATH" envDefault:"secrets/data/chewedfeed/details"`

	ExpireTime time.Time

  Exclusive bool
}

type System struct {
	Host           string `env:"RABBIT_HOSTNAME" envDefault:"" json:"host,omitempty"`
	ManagementHost string `env:"RABBIT_MANAGEMENT_HOSTNAME" envDefault:"" json:"management_host,omitempty"`
	Port           int    `env:"RABBIT_PORT" envDefault:"" json:"port,omitempty"`
	Username       string `env:"RABBIT_USERNAME" envDefault:"" json:"username,omitempty"`
	Password       string `env:"RABBIT_PASSWORD" envDefault:"" json:"password,omitempty"`
	VHost          string `env:"RABBIT_VHOST" envDefault:"" json:"vhost,omitempty"`
	Queue          string `env:"RABBIT_QUEUE" envDefault:"" json:"queue,omitempty"`

	VaultDetails
	HTTPClient
	VaultHelper
}

func NewRabbit(port int, host, username, password, vhost, management string, httpClient HTTPClient, vaultHelper VaultHelper) *System {
	return &System{
		Host:           host,
		Port:           port,
		Username:       username,
		Password:       password,
		VHost:          vhost,
		ManagementHost: management,

		HTTPClient:  httpClient,
		VaultHelper: vaultHelper,
	}
}

func Setup(vaultAddress, vaultToken string, exclusive bool) VaultDetails {
	return VaultDetails{
		Address: vaultAddress,
		Token:   vaultToken,

    Exclusive: exclusive,
	}
}

func vaultBuild(vd VaultDetails, vh vaulthelper.VaultHelper, client HTTPClient) (*System, error) {
  r:= NewRabbit(0, "", "", "", "", "", client, vh)

  // Creds
  if err := vh.GetSecrets(vd.CredPath); err != nil {
    return r, logs.Errorf("failed to get cred secrets from vault: %v", err)
  }
  if vh.Secrets() == nil {
    return r, logs.Error("no rabbit cred secrets found")
  }

  username, err := vh.GetSecret("username")
  if err != nil {
    return r, logs.Errorf("failed to get username: %v", err)
  }
  r.Username = username

  password, err := vh.GetSecret("password")
  if err != nil {
    return r, logs.Errorf("failed to get password: %v", err)
  }
  r.Password = password

  // Details
  if err := vh.GetSecrets(vd.DetailsPath); err != nil {
    return r, logs.Errorf("failed to get details secrets from vault: %v", err)
  }
  if vh.Secrets() == nil {
    return r, logs.Error("no rabbit detail secrets found")
  }

  host, err := vh.GetSecret("rabbit-hostname")
  if err != nil {
    return r, logs.Errorf("failed to get hostname: %v", err)
  }
  r.Host = host

  managementHost, err := vh.GetSecret("rabbit-management-host")
  if err != nil {
    return r, logs.Errorf("failed to get rabbit management host: %v", err)
  }
  r.ManagementHost = managementHost

  port, err := vh.GetSecret("rabbit-port")
  if err != nil {
    return r, logs.Errorf("failed to get rabbit port: %v", err)
  }
  if port != "" {
    iport, err := strconv.Atoi(port)
    if err != nil {
      return r, logs.Errorf("failed to parse rabbit port: %v", err)
    }
    r.Port = iport
  }

  vhost, err := vh.GetSecret("rabbit-vhost")
  if err != nil {
    return r, logs.Errorf("failed to get rabbit vhost: %v", err)
  }
  r.VHost = vhost

  queue, err := vh.GetSecret("rabbit-queue")
  if err != nil {
    return r, logs.Errorf("failed to get rabbit queue: %v", err)
  }
  r.Queue = queue

  return r, nil
}

func Build(vd VaultDetails, vh vaulthelper.VaultHelper, httpClient HTTPClient) (*System, error) {
	r := NewRabbit(0, "", "", "", "", "", httpClient, vh)
  r.VaultDetails = vd

  if vd.Exclusive {
    return vaultBuild(vd, vh, httpClient)
  }

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

func GetRabbitQueue(ctx context.Context, r System) (interface{}, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if time.Now().Unix() > r.VaultDetails.ExpireTime.Unix() {
		rb, err := Build(r.VaultDetails, vaulthelper.NewVault(r.VaultDetails.Address, r.VaultDetails.Token), r.HTTPClient)
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

	res, err := r.HTTPClient.Do(req)
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
