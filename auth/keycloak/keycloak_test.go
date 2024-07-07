package keycloak

import (
	vaultHelper "github.com/keloran/vault-helper"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildVault(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
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
