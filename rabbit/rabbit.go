package rabbit

import (
	"context"
	"fmt"
	"net/http"
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
	Address string
	Token   string

	CredPath    string `env:"RABBIT_VAULT_CREDS_PATH" envDefault:"secrets/data/chewedfeed/rabbitmq"`
	DetailsPath string `env:"RABBIT_VAULT_DETAILS_PATH" envDefault:"secrets/data/chewedfeed/details"`

	ExpireTime time.Time
}

type Details struct {
	Host           string `env:"RABBIT_HOSTNAME" envDefault:"" json:"host,omitempty"`
	ManagementHost string `env:"RABBIT_MANAGEMENT_HOSTNAME" envDefault:"" json:"management_host,omitempty"`
	Username       string `env:"RABBIT_USERNAME" envDefault:"" json:"username,omitempty"`
	Password       string `env:"RABBIT_PASSWORD" envDefault:"" json:"password,omitempty"`
	VHost          string `env:"RABBIT_VHOST" envDefault:"" json:"vhost,omitempty"`
	Queue          string `env:"RABBIT_QUEUE" envDefault:"" json:"queue,omitempty"`
}

type System struct {
	Context context.Context

	Details

	HTTPClient

	VaultDetails
	VaultHelper *vaulthelper.VaultHelper
}

func NewSystem(httpClient HTTPClient) *System {
	return &System{
		Context:    context.Background(),
		HTTPClient: httpClient,
	}
}

func (s *System) Setup(vd VaultDetails, vh vaulthelper.VaultHelper) {
	s.VaultDetails = vd
	s.VaultHelper = &vh
}

func (s *System) Build() (*Details, error) {
	gen, err := s.buildGeneric()
	if err != nil {
		return nil, err
	}

	if s.VaultHelper != nil {
		return s.buildVault()
	}

	return gen, nil
}

func (s *System) buildGeneric() (*Details, error) {
	rab := &Details{}

	if err := env.Parse(rab); err != nil {
		return nil, logs.Errorf("failed to parse env: %v", err)
	}
	s.Details = *rab

	return rab, nil
}

func (s *System) buildVault() (*Details, error) {
	rab := &Details{}
	vh := *s.VaultHelper

	if err := vh.GetSecrets(s.VaultDetails.DetailsPath); err != nil {
		return rab, logs.Errorf("failed to get rabbit secrets from vault: %v", err)
	}
	if vh.Secrets() == nil {
		return rab, logs.Error("no rabbit secrets found")
	}

	if rab.Username == "" {
		secret, err := vh.GetSecret("rabbit-username")
		if err != nil {
			return nil, logs.Errorf("failed to get username: %v", err)
		}
		rab.Username = secret
	}

	if rab.Password == "" {
		secret, err := vh.GetSecret("rabbit-password")
		if err != nil {
			return nil, logs.Errorf("failed to get password: %v", err)
		}
		rab.Password = secret
	}

	if rab.Host == "" {
		secret, err := vh.GetSecret("rabbit-hostname")
		if err != nil {
			return nil, logs.Errorf("failed to get hostname: %v", err)
		}
		rab.Host = secret
	}

	if rab.VHost == "" {
		secret, err := vh.GetSecret("rabbit-vhost")
		if err != nil {
			return nil, logs.Errorf("failed to get vhost: %v", err)
		}
		rab.VHost = secret
	}

	if rab.ManagementHost == "" {
		secret, err := vh.GetSecret("rabbit-management-hostname")
		if err != nil {
			return nil, logs.Errorf("failed to get management host: %v", err)
		}
		rab.ManagementHost = secret
	}

	if rab.Queue == "" {
		secret, err := vh.GetSecret("rabbit-queue")
		if err != nil {
			return nil, logs.Errorf("failed to get queue: %v", err)
		}
		rab.Queue = secret
	}
	s.Details = *rab

	return rab, nil
}

func (s *System) GetRabbitQueue() (interface{}, error) {
	if s.VaultHelper != nil && time.Now().Unix() > s.VaultDetails.ExpireTime.Unix() {
		_, err := s.Build()
		if err != nil {
			return nil, logs.Errorf("rabbit: unable to build rabbit: %v", err)
		}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/queues/%s/%s/get", s.Details.Host, s.Details.VHost, s.Details.Queue), nil)
	if err != nil {
		return nil, logs.Errorf("rabbit: unable to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(s.Details.Username, s.Details.Password)

	res, err := s.HTTPClient.Do(req)
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
