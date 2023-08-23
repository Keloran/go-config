package mongo

import (
	"context"
	"fmt"
	"os"
	"testing"

	vault_helper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockVaultHelper struct {
	KVSecrets []vault_helper.KVSecret
	Lease     int
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

func (m *MockVaultHelper) LeaseDuration() int {
	return m.Lease
}

type MockMongoClient struct {
	mock.Mock
}

func (m *MockMongoClient) Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*mongo.Client), args.Error(1)
}

func TestBuild(t *testing.T) {
	mockVault := &MockVaultHelper{
		KVSecrets: []vault_helper.KVSecret{
			{Key: "password", Value: "testPassword"},
			{Key: "username", Value: "testUser"},
		},
	}

	vd := Setup("testVaultAddress", "testVaultToken")
	m, err := Build(vd, mockVault)

	assert.NoError(t, err)
	assert.Equal(t, "testUser", m.Username)
	assert.Equal(t, "testPassword", m.Password)
}

func TestBuildCollections(t *testing.T) {
	os.Clearenv() // Clear all environment variables
	_ = os.Setenv("MONGO_COLLECTION_BOB", "bill")
	_ = os.Setenv("MONGO_COLLECTION_ALICE", "wonderland")
	_ = os.Setenv("MONGO_HOST", "localhost") // This should be ignored

	collections := BuildCollections()

	expected := map[string]string{
		"bob":   "bill",
		"alice": "wonderland",
	}

	assert.Equal(t, expected, collections, "Collections map did not match expected")
}

func TestBuildCollectionsNoMatch(t *testing.T) {
	os.Clearenv() // Clear all environment variables
	_ = os.Setenv("MONGO_HOST", "localhost")
	_ = os.Setenv("MONGO_USER", "user")

	collections := BuildCollections()

	assert.Empty(t, collections, "Expected Collections map to be empty")
}
