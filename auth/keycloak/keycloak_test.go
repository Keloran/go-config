package keycloak

import (
	vaultHelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	testRealm  = "test-realm"
	testClient = "test-client"
	testSecret = "test-secret"
)

func TestBuildVault(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "keycloak-client", Value: testClient},
			{Key: "keycloak-secret", Value: testSecret},
			{Key: "keycloak-realm", Value: testRealm},
		},
	}

	vd := &VaultDetails{
		Address:     "mockAddress",
		Token:       "mockToken",
		DetailsPath: "tester",
	}
	d := NewSystem()
	d.Setup(*vd, mockVault)
	kc, err := d.Build()

	assert.NoError(t, err)

	assert.Equal(t, testClient, kc.Client)
	assert.Equal(t, testSecret, kc.Secret)
	assert.Equal(t, testRealm, kc.Realm)
	assert.Equal(t, "https://keys.chewedfeed.com", kc.Host)
}

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("KEYCLOAK_CLIENT", testClient); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("KEYCLOAK_SECRET", testSecret); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("KEYCLOAK_REALM", testRealm); err != nil {
		t.Fatal(err)
	}

	d := NewSystem()
	kc, err := d.Build()
	assert.NoError(t, err)

	assert.Equal(t, testClient, kc.Client)
	assert.Equal(t, testSecret, kc.Secret)
	assert.Equal(t, testRealm, kc.Realm)
	assert.Equal(t, "https://keys.chewedfeed.com", kc.Host)
}
