package mongo

import (
	"fmt"
	"os"
	"testing"

	"github.com/bugfixes/go-bugfixes/logs"

	vaultHelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
)

type MockVaultHelper struct {
	KVSecrets []vaultHelper.KVSecret
	Lease     int
}

func (m *MockVaultHelper) GetSecrets(path string) error {
	if path == "" {
		return logs.Error("path not found")
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

func TestBuildCollections(t *testing.T) {
	os.Clearenv() // Clear all environment variables
	if err := os.Setenv("MONGO_COLLECTION_BOB", "bill"); err != nil {
		assert.NoError(t, err)
	}
	if err := os.Setenv("MONGO_COLLECTION_ALICE", "wonderland"); err != nil {
		assert.NoError(t, err)
	}
	if err := os.Setenv("MONGO_HOST", "localhost"); err != nil {
		assert.NoError(t, err)
	}

	collections := BuildCollections()

	expected := map[string]string{
		"bob":   "bill",
		"alice": "wonderland",
	}

	assert.Equal(t, expected, collections, "Collections map did not match expected")
}

func TestBuildCollectionsNoMatch(t *testing.T) {
	os.Clearenv() // Clear all environment variables
	if err := os.Setenv("MONGO_HOST", "localhost"); err != nil {
		assert.NoError(t, err)
	}
	if err := os.Setenv("MONGO_USER", "user"); err != nil {
		assert.NoError(t, err)
	}

	collections := BuildCollections()

	assert.Empty(t, collections, "Expected Collections map to be empty")
}

func TestBuildCollectionsNoHost(t *testing.T) {
	os.Clearenv() // Clear all environment variables
	if err := os.Setenv("MONGO_COLLECTION_BOB", "bill"); err != nil {
		assert.NoError(t, err)
	}
	if err := os.Setenv("MONGO_COLLECTION_ALICE", "wonderland"); err != nil {
		assert.NoError(t, err)
	}

	collections := BuildCollections()

	assert.NotEmpty(t, collections, "Expected Collections map to be empty")
}

func TestBuildCollectionsNoCollections(t *testing.T) {
	os.Clearenv() // Clear all environment variables
	if err := os.Setenv("MONGO_HOST", "localhost"); err != nil {
		assert.NoError(t, err)
	}

	collections := BuildCollections()

	assert.Empty(t, collections, "Expected Collections map to be empty")
}
