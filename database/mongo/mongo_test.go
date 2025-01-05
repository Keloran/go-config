package mongo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
