package database

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	vaultHelper "github.com/keloran/vault-helper"
)

type MockVaultHelper struct {
	KVSecrets []vaultHelper.KVSecret
	Lease     int
}

func (m *MockVaultHelper) GetSecrets(path string) error {
	if path == "" {
		return fmt.Errorf("path not found: %s", path)
	}

	return nil // or simulate an error if needed
}

func (m *MockVaultHelper) GetSecret(key string) (string, error) {
	for _, s := range m.Secrets() {
		if s.Key == key {
			return s.Value, nil
		}
	}
	return "", fmt.Errorf("key not found: %s", key)
}

func (m *MockVaultHelper) Secrets() []vaultHelper.KVSecret {
	return m.KVSecrets
}

func (m *MockVaultHelper) LeaseDuration() int {
	return m.Lease
}

func TestBuildVault(t *testing.T) {
	mockVault := &MockVaultHelper{
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
	mockVault := &MockVaultHelper{
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
