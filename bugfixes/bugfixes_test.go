package bugfixes

import (
	"errors"
	vaultHelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("BUGFIXES_AGENT_KEY", "testKey"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("BUGFIXES_AGENT_SECRET", "testSecret"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("BUGFIXES_SERVER", "http://bob.bob"); err != nil {
		t.Fatal(err)
	}

	b := NewSystem()
	bf, err := b.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", bf.AgentKey)
	assert.Equal(t, "testSecret", bf.AgentSecret)
	assert.Equal(t, "http://bob.bob", bf.Server)
}

func TestBuildVault(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "bugfixes-agentid", Value: "testKey"},
			{Key: "bugfixes-secret", Value: "testSecret"},
			{Key: "bugfixes-server", Value: "http://bob.bob"},
		},
	}

	vd := &vaultHelper.VaultDetails{
		DetailsPath: "tester",
	}
	b := NewSystem()
	b.Setup(*vd, mockVault)
	bf, err := b.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", bf.AgentKey)
	assert.Equal(t, "testSecret", bf.AgentSecret)
	assert.Equal(t, "http://bob.bob", bf.Server)
}

func TestBuildVaultNoHost(t *testing.T) {
	os.Clearenv()
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "bugfixes-agentid", Value: "testKey"},
			{Key: "bugfixes-secret", Value: "testSecret"},
		},
	}

	vd := &vaultHelper.VaultDetails{
		DetailsPath: "tester",
	}
	b := NewSystem()
	b.Setup(*vd, mockVault)
	bf, err := b.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", bf.AgentKey)
	assert.Equal(t, "testSecret", bf.AgentSecret)
	assert.Equal(t, "https://api.bugfix.es/v1", bf.Server)
}

func TestBuildVaultInvalidHost(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "bugfixes-agentid", Value: "testKey"},
			{Key: "bugfixes-secret", Value: "testSecret"},
			{Key: "bugfixes-server", Value: "bob.bob"},
		},
	}

	vd := &vaultHelper.VaultDetails{
		DetailsPath: "tester",
	}
	b := NewSystem()
	b.Setup(*vd, mockVault)
	_, err := b.Build()
	assert.Error(t, errors.New("needs the protocol for server"), err)
}
