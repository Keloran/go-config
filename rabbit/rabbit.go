package rabbit

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Details struct {
	Host           string `env:"RABBIT_HOSTNAME" envDefault:"" json:"host,omitempty" vault:"rabbit-hostname"`
	ManagementHost string `env:"RABBIT_MANAGEMENT_HOSTNAME" envDefault:"" json:"management_host,omitempty" vault:"rabbit-management-hostname"`
	Username       string `env:"RABBIT_USERNAME" envDefault:"" json:"username,omitempty" vault:"rabbit-username"`
	Password       string `env:"RABBIT_PASSWORD" envDefault:"" json:"password,omitempty" vault:"rabbit-password"`
	VHost          string `env:"RABBIT_VHOST" envDefault:"" json:"vhost,omitempty" vault:"rabbit-vhost"`
	Queue          string `env:"RABBIT_QUEUE" envDefault:"" json:"queue,omitempty" vault:"rabbit-queue"`
}

type System struct {
	Context context.Context

	Details

	HTTPClient
}

func NewSystem(httpClient HTTPClient) *System {
	return &System{
		Context:    context.Background(),
		HTTPClient: httpClient,
	}
}

func (s *System) Build() (*Details, error) {
	return s.buildGeneric()
}

func (s *System) buildGeneric() (*Details, error) {
	rab := &Details{}

	if err := env.Parse(rab); err != nil {
		return nil, logs.Errorf("failed to parse env: %v", err)
	}
	s.Details = *rab

	return rab, nil
}

func (s *System) GetRabbitQueue() (interface{}, error) {
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
