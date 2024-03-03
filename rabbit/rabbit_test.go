package rabbit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	vaulthelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct{}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if req.URL.Path == "/api/queues/testVhost/testQueue/get" {
		response := `[
			{"payload":"test message","payload_bytes":12,"redelivered":false}
		]`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(response)),
			Header:     make(http.Header),
		}, nil
	}
	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewBufferString("")),
		Header:     make(http.Header),
	}, nil
}

type MockVaultHelper struct {
	KVSecrets []vaulthelper.KVSecret
	Lease     int
}

func (m *MockVaultHelper) GetSecrets(path string) error {
	if path == "" {
		return nil
	}

	return nil // or simulate an error if needed
}
func (m *MockVaultHelper) GetSecret(key string) (string, error) {
	for _, s := range m.Secrets() {
		if s.Key == key {
			return s.Value, nil
		}
	}
	return "", fmt.Errorf("key not found")
}
func (m *MockVaultHelper) Secrets() []vaulthelper.KVSecret {
	return m.KVSecrets
}
func (m *MockVaultHelper) LeaseDuration() int {
	return m.Lease
}

func TestBuild(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv()

		mockVault := &MockVaultHelper{
			KVSecrets: []vaulthelper.KVSecret{
				{Key: "password", Value: ""},
				{Key: "username", Value: ""},
				{Key: "vhost", Value: ""},
			},
		}
		vd := Setup("mockAddress", "mockToken")

		l, err := Build(vd, mockVault, &MockHTTPClient{})
		assert.NoError(t, err)
		assert.Equal(t, "", l.Host)
		assert.Equal(t, 0, l.Port)
		assert.Equal(t, "", l.Username)
		assert.Equal(t, "", l.Password)
		assert.Equal(t, "", l.VHost)
		assert.Equal(t, "", l.ManagementHost)
	})

	t.Run("with values", func(t *testing.T) {
		mockVault := &MockVaultHelper{
			KVSecrets: []vaulthelper.KVSecret{
				{Key: "password", Value: "testPassword"},
				{Key: "username", Value: "testUser"},
				{Key: "vhost", Value: "testVhost"},
			},
		}

		vd := Setup("mockAddress", "mockToken")

		os.Clearenv()
		if err := os.Setenv("RABBIT_HOSTNAME", "http://localhost"); err != nil {
			assert.NoError(t, err)
		}
		r, err := Build(vd, mockVault, &MockHTTPClient{})
		assert.NoError(t, err)
		assert.Equal(t, "http://localhost", r.Host)
	})
}

func TestGetRabbitQueue(t *testing.T) {
	t.Run("successful queue retrieval", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{}
		mockVaultHelper := &MockVaultHelper{}

		// Setup Rabbit instance with mocks
		rabbit := Rabbit{
			Host:        "http://localhost",
			Username:    "testUser",
			Password:    "testPassword",
			VHost:       "testVhost",
			Queue:       "testQueue",
			HTTPClient:  mockHTTPClient,
			VaultHelper: mockVaultHelper,
			VaultDetails: VaultDetails{
				ExpireTime: time.Now().Add(time.Hour),
			},
		}

		// Test GetRabbitQueue
		result, err := GetRabbitQueue(context.Background(), rabbit)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("queue retrieval failure", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{}
		mockVaultHelper := &MockVaultHelper{}

		// Setup Rabbit instance with mocks and an invalid queue name to simulate failure
		rabbit := Rabbit{
			Host:        "http://localhost",
			Username:    "testUser",
			Password:    "testPassword",
			VHost:       "testVhost",
			Queue:       "invalidQueue",
			HTTPClient:  mockHTTPClient,
			VaultHelper: mockVaultHelper,
			VaultDetails: VaultDetails{
				ExpireTime: time.Now().Add(time.Hour),
			},
		}

		// Test GetRabbitQueue
		result, err := GetRabbitQueue(context.Background(), rabbit)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
