package database

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	vaultHelper "github.com/keloran/vault-helper"
)

func TestBuildVault(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "password", Value: "testPassword"},
			{Key: "username", Value: "testUser"},
			{Key: "rds-port", Value: "1111"},
			{Key: "rds-db", Value: "testDB"},
			{Key: "rds-hostname", Value: "testHost"},
		},
	}

	vd := &VaultDetails{
		CredPath:    "tester",
		DetailsPath: "tester",
	}
	d := NewSystem()
	d.Setup(*vd, mockVault)
	db, err := d.Build()
	assert.NoError(t, err)

	assert.Equal(t, "testPassword", db.Password)
	assert.Equal(t, "testUser", db.User)
	assert.Equal(t, 1111, db.Port)
	assert.Equal(t, "testDB", db.DBName)
	assert.Equal(t, "testHost", db.Host)
}

func TestBuildVaultNoPort(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "password", Value: "testPassword"},
			{Key: "username", Value: "testUser"},
			{Key: "rds-db", Value: "testDB"},
			{Key: "rds-hostname", Value: "testHost"},
		},
	}

	vd := &VaultDetails{
		CredPath:    "tester",
		DetailsPath: "tester",
	}
	d := NewSystem()
	d.Setup(*vd, mockVault)
	db, err := d.Build()
	assert.NoError(t, err)

	assert.Equal(t, "testPassword", db.Password)
	assert.Equal(t, "testUser", db.User)
	assert.Equal(t, 5432, db.Port)
	assert.Equal(t, "testDB", db.DBName)
	assert.Equal(t, "testHost", db.Host)
}

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("RDS_HOSTNAME", "testHost"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("RDS_PORT", "1111"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("RDS_USERNAME", "testUser"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("RDS_PASSWORD", "testPassword"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("RDS_DB", "testDB"); err != nil {
		t.Fatal(err)
	}

	d := NewSystem()
	db, err := d.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testPassword", db.Password)
	assert.Equal(t, "testUser", db.User)
	assert.Equal(t, 1111, db.Port)
	assert.Equal(t, "testDB", db.DBName)
	assert.Equal(t, "testHost", db.Host)
}
