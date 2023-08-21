package database

import (
	"fmt"
	vault_helper "github.com/keloran/vault-helper"
	"testing"
)

type MockVaultHelper struct {
	KVSecrets []vault_helper.KVSecret
}

func (m *MockVaultHelper) GetSecrets(path string) error {
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

func (m *MockVaultHelper) Secrets() []vault_helper.KVSecret {
	return m.KVSecrets
}

func TestBuild(t *testing.T) {
	mockVault := &MockVaultHelper{
		KVSecrets: []vault_helper.KVSecret{
			{Key: "password", Value: "testPassword"},
			{Key: "username", Value: "testUser"},
		},
	}

	vd := Setup("mockAddress", "mockToken")
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
