package influx

import (
	vaultHelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("INFLUX_HOSTNAME", "testHost"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("INFLUX_TOKEN", "testToken"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("INFLUX_BUCKET", "testBucket"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("INFLUX_ORG", "testOrg"); err != nil {
		t.Fatal(err)
	}

	i := NewSystem()
	in, err := i.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testToken", in.Token)
	assert.Equal(t, "testBucket", in.Bucket)
	assert.Equal(t, "testOrg", in.Org)
	assert.Equal(t, "testHost", in.Host)
}

func TestBuildVault(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "influx-token", Value: "testToken"},
			{Key: "influx-bucket", Value: "testBucket"},
			{Key: "influx-hostname", Value: "testHost"},
			{Key: "influx-org", Value: "testOrg"},
		},
	}

	vd := &VaultDetails{
		DetailsPath: "tester",
	}
	i := NewSystem()
	i.Setup(*vd, mockVault)
	in, err := i.Build()
	assert.NoError(t, err)

	assert.Equal(t, "testToken", in.Token)
	assert.Equal(t, "testBucket", in.Bucket)
	assert.Equal(t, "testOrg", in.Org)
	assert.Equal(t, "testHost", in.Host)
}

func TestBuildVaultNoHost(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "influx-token", Value: "testToken"},
			{Key: "influx-bucket", Value: "testBucket"},
			{Key: "influx-org", Value: "testOrg"},
		},
	}

	vd := &VaultDetails{
		DetailsPath: "tester",
	}
	i := NewSystem()
	i.Setup(*vd, mockVault)
	in, err := i.Build()
	assert.NoError(t, err)

	assert.Equal(t, "testToken", in.Token)
	assert.Equal(t, "testBucket", in.Bucket)
	assert.Equal(t, "testOrg", in.Org)
	assert.Equal(t, "http://db.chewed-k8s.net:8086", in.Host)
}
