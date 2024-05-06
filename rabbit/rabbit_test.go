package rabbit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	vaulthelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct{}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if req.URL.Path == "testHost/api/queues/testVHost/testQueue/get" {
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

func TestVaultBuild(t *testing.T) {
	mockVault := &MockVaultHelper{
		KVSecrets: []vaulthelper.KVSecret{
			{Key: "rabbit-hostname", Value: "testHost"},
			{Key: "rabbit-management-hostname", Value: "testManagementHost"},
			{Key: "rabbit-username", Value: "testUsername"},
			{Key: "rabbit-password", Value: "testPassword"},
			{Key: "rabbit-vhost", Value: "testVHost"},
			{Key: "rabbit-queue", Value: "testQueue"},
		},
	}

	vd := &VaultDetails{
		Address:     "mockAddress",
		Token:       "mockToken",
		DetailsPath: "tester",
	}

	mockHTTPClient := &MockHTTPClient{}

	d := NewSystem(mockHTTPClient)
	d.Setup(*vd, mockVault)
	rab, err := d.Build()
	assert.NoError(t, err)

	assert.Equal(t, "testUsername", rab.Username)
	assert.Equal(t, "testPassword", rab.Password)
	assert.Equal(t, "testHost", rab.Host)
	assert.Equal(t, "testManagementHost", rab.ManagementHost)
	assert.Equal(t, "testVHost", rab.VHost)
	assert.Equal(t, "testQueue", rab.Queue)
}

func TestGenericBuild(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("RABBIT_HOSTNAME", "testHost"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("RABBIT_USERNAME", "testUsername"); err != nil {
		t.Fatal(err)
	}

	mockHTTPClient := &MockHTTPClient{}
	d := NewSystem(mockHTTPClient)
	rab, err := d.Build()
	assert.NoError(t, err)

	assert.Equal(t, "testHost", rab.Host)
	assert.Equal(t, "testUsername", rab.Username)
}

func TestGetRabbitQueue(t *testing.T) {
	mockVault := &MockVaultHelper{
		KVSecrets: []vaulthelper.KVSecret{
			{Key: "rabbit-hostname", Value: "testHost"},
			{Key: "rabbit-management-hostname", Value: "testManagementHost"},
			{Key: "rabbit-username", Value: "testUsername"},
			{Key: "rabbit-password", Value: "testPassword"},
			{Key: "rabbit-vhost", Value: "testVHost"},
			{Key: "rabbit-queue", Value: "testQueue"},
		},
	}

	vd := &VaultDetails{
		Address:     "mockAddress",
		Token:       "mockToken",
		DetailsPath: "tester",
	}

	mockHTTPClient := &MockHTTPClient{}

	d := NewSystem(mockHTTPClient)
	d.Setup(*vd, mockVault)
	_, err := d.Build()
	assert.NoError(t, err)

	result, err := d.GetRabbitQueue()
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
