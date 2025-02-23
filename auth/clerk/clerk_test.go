package clerk

import (
	vaultHelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("CLERK_SECRET_KEY", "testKey"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY", "testPublicKey"); err != nil {
		t.Fatal(err)
	}

	c := NewSystem()
	ck, err := c.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", ck.Key)
	assert.Equal(t, "testPublicKey", ck.PublicKey)
}

func TestBuildVault(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "clerk_key", Value: "testKey"},
			{Key: "clerk_public_key", Value: "testPublicKey"},
		},
	}

	vd := &vaultHelper.VaultDetails{
		DetailsPath: "tester",
	}
	c := NewSystem()
	c.Setup(*vd, mockVault)
	ck, err := c.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", ck.Key)
	assert.Equal(t, "testPublicKey", ck.PublicKey)
}

func TestBuildVaultNoKey(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "clerk_key", Value: "testKey"},
			{Key: "clerk_public_key", Value: "testPublicKey"},
		},
	}

	vd := &vaultHelper.VaultDetails{
		DetailsPath: "tester",
	}
	c := NewSystem()
	c.Setup(*vd, mockVault)
	ck, err := c.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", ck.Key)
	assert.Equal(t, "testPublicKey", ck.PublicKey)
}

func TestBuildVaultNoPublicKey(t *testing.T) {
	os.Clearenv()
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "clerk_key", Value: "testKey"},
		},
	}

	vd := &vaultHelper.VaultDetails{
		DetailsPath: "tester",
	}
	c := NewSystem()
	c.Setup(*vd, mockVault)
	ck, err := c.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", ck.Key)
	assert.Equal(t, "", ck.PublicKey)
}
