package keycloak

import (
	"fmt"
	vaulthelper "github.com/keloran/vault-helper"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockVaultHelper struct {
	KVSecrets []vaulthelper.KVSecret
	Lease     int
}

func (m *MockVaultHelper) GetSecrets(path string) error {
	if path == "" {
		return fmt.Errorf("path not found: %s", path)
	}

	return nil
}

func (m *MockVaultHelper) GetSecret(key string) (string, error) {
	for _, s := range m.Secrets() {
		for s.Key == key {
			return s.Value, nil
		}
	}

	return "", fmt.Errorf("key: '%s' not found", key)
}

func (m *MockVaultHelper) Secrets() []vaulthelper.KVSecret {
	return m.KVSecrets
}
func (m *MockVaultHelper) LeaseDuration() int {
	return m.Lease
}

func TestBuildVault(t *testing.T) {
	mockVault := &MockVaultHelper{
		KVSecrets: []vaulthelper.KVSecret{
			{Key: "keycloak-client", Value: "testClient"},
			{Key: "keycloak-secret", Value: "testSecret"},
			{Key: "keycloak-realm", Value: "testRealm"},
		},
	}

	vd := &VaultDetails{
		Address:     "mockAddress",
		Token:       "mockToken",
		DetailsPath: "tester",
	}
	d := NewSystem()
	d.Setup(*vd, mockVault)
	key, err := d.Build()

	assert.NoError(t, err)

	assert.Equal(t, "testClient", key.Client)
	assert.Equal(t, "testSecret", key.Secret)
	assert.Equal(t, "testRealm", key.Realm)
	assert.Equal(t, "https://keys.chewedfeed.com", key.Host)
}

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("KEYCLOAK_CLIENT", "testClient"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("KEYCLOAK_SECRET", "testSecret"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("KEYCLOAK_REALM", "testRealm"); err != nil {
		t.Fatal(err)
	}

	d := NewSystem()
	key, err := d.Build()
	assert.NoError(t, err)

	assert.Equal(t, "testClient", key.Client)
	assert.Equal(t, "testSecret", key.Secret)
	assert.Equal(t, "testRealm", key.Realm)
	assert.Equal(t, "https://keys.chewedfeed.com", key.Host)
}
