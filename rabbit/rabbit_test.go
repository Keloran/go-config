package rabbit

import (
	"fmt"
	vaulthelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

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

		l, err := Build(vd, mockVault)
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
		r, err := Build(vd, mockVault)
		assert.NoError(t, err)
		assert.Equal(t, "http://localhost", r.Host)
	})
}
