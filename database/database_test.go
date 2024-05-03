package database

import (
	"fmt"
	"testing"

	vaultHelper "github.com/keloran/vault-helper"
)

type MockVaultHelper struct {
	KVSecrets []vaultHelper.KVSecret
	Lease     int
}

func (m *MockVaultHelper) GetSecrets(path string) error {
	if path == "" {
		return fmt.Errorf("path not found")
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

func (m *MockVaultHelper) Secrets() []vaultHelper.KVSecret {
	return m.KVSecrets
}

func (m *MockVaultHelper) LeaseDuration() int {
	return m.Lease
}

func TestBuild(t *testing.T) {
	mockVault := &MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "password", Value: "testPassword"},
			{Key: "username", Value: "testUser"},
		},
	}

	vd := Setup("mockAddress", "mockToken", false, nil)
	db, err := Build(vd, mockVault)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if db.Password != "testPassword" {
		t.Errorf("expected password to be 'testPassword', got %s", db.Password)
	}

	if db.User != "testUser" {
		t.Errorf("expected user to be 'testUser', got %s", db.User)
	}
}
