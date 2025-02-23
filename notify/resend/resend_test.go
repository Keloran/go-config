package resend

import (
	vaultHelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("RESEND_KEY", "testKey"); err != nil {
		t.Fatal(err)
	}

	r := NewSystem()
	rd, err := r.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", rd.Key)
}

func TestBuildVault(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "resend_key", Value: "testKey"},
		},
	}

	vd := &vaultHelper.VaultDetails{
		DetailsPath: "tester",
	}
	r := NewSystem()
	r.Setup(*vd, mockVault)
	rd, err := r.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", rd.Key)
}

func TestBuildVaultNoKey(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "resend_key", Value: "testKey"},
		},
	}

	vd := &vaultHelper.VaultDetails{
		DetailsPath: "tester",
	}
	r := NewSystem()
	r.Setup(*vd, mockVault)
	rd, err := r.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", rd.Key)
}
