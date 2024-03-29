package ConfigBuilder

import (
	"fmt"
	"os"
	"testing"

	vaulthelper "github.com/keloran/vault-helper"

	"github.com/stretchr/testify/assert"
)

type MockVaultHelper struct {
	KVSecrets []vaulthelper.KVSecret
	Lease     int
}

func (m *MockVaultHelper) GetSecrets(path string) error {
	if path == "" {
		return nil
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

func (m *MockVaultHelper) Secrets() []vaulthelper.KVSecret {
	return m.KVSecrets
}

func (m *MockVaultHelper) LeaseDuration() int {
	return m.Lease
}

func TestBuild(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv() // Clear all environment variables

		cfg, err := Build(Local)
		assert.NoError(t, err)
		assert.Equal(t, false, cfg.KeepLocal)
		assert.Equal(t, false, cfg.Development)
		assert.Equal(t, 80, cfg.HTTPPort)
	})

	t.Run("custom values", func(t *testing.T) {
		os.Clearenv()
		// Set custom environment variables for Local
		if err := os.Setenv("BUGFIXES_LOCAL_ONLY", "true"); err != nil {
			assert.NoError(t, err)
		}
		if err := os.Setenv("DEVELOPMENT", "true"); err != nil {
			assert.NoError(t, err)
		}
		if err := os.Setenv("HTTP_PORT", "8080"); err != nil {
			assert.NoError(t, err)
		}

		cfg, err := Build(Local)
		assert.NoError(t, err)
		assert.Equal(t, true, cfg.KeepLocal)
		assert.Equal(t, true, cfg.Development)
		assert.Equal(t, 8080, cfg.HTTPPort)
	})
}

func TestRabbit(t *testing.T) {
	t.Run("rabbit no set values", func(t *testing.T) {
		os.Clearenv()
		mockVault := &MockVaultHelper{
			KVSecrets: []vaulthelper.KVSecret{
				{Key: "password", Value: ""},
				{Key: "username", Value: ""},
				{Key: "vhost", Value: ""},
			},
		}

		cfg, err := BuildLocal(mockVault, Rabbit)
		assert.NoError(t, err)
		assert.Equal(t, "", cfg.Rabbit.Host)
	})
	t.Run("rabbit with values", func(t *testing.T) {
		os.Clearenv()
		mockVault := &MockVaultHelper{
			KVSecrets: []vaulthelper.KVSecret{
				{Key: "password", Value: "testPassword"},
				{Key: "username", Value: "testUser"},
				{Key: "vhost", Value: "testVhost"},
			},
		}

		cfg, err := BuildLocal(mockVault, Rabbit)
		assert.NoError(t, err)
		assert.Equal(t, "testUser", cfg.Rabbit.Username)
	})
}

func TestDatabase(t *testing.T) {
	mockVault := &MockVaultHelper{
		KVSecrets: []vaulthelper.KVSecret{
			{Key: "password", Value: "testPassword"},
			{Key: "username", Value: "testUser"},
		},
	}

	t.Run("database no set values", func(t *testing.T) {
		os.Clearenv()
		cfg, err := BuildLocal(mockVault, Database)
		assert.NoError(t, err)
		assert.Equal(t, "postgres.chewedfeed", cfg.Database.Host)
	})
	t.Run("database with values", func(t *testing.T) {
		os.Clearenv()
		cfg, err := BuildLocal(mockVault, Database)
		assert.NoError(t, err)
		assert.Equal(t, "testUser", cfg.Database.User)
	})
}

func TestKeycloak(t *testing.T) {
	t.Run("keycloak", func(t *testing.T) {
		os.Clearenv()
		cfg, err := Build(Keycloak)
		assert.NoError(t, err)
		assert.Equal(t, "", cfg.Keycloak.Host)
	})
}

func TestVault(t *testing.T) {
	t.Run("vault", func(t *testing.T) {
		os.Clearenv()
		cfg, err := Build(Vault)
		assert.NoError(t, err)
		assert.Equal(t, "localhost", cfg.Vault.Host)
	})
}

func TestMongo(t *testing.T) {
	mockVault := &MockVaultHelper{
		KVSecrets: []vaulthelper.KVSecret{
			{Key: "password", Value: "testPassword"},
			{Key: "username", Value: "testUser"},
		},
	}

	t.Run("mongo no set values", func(t *testing.T) {
		os.Clearenv()
		cfg, err := BuildLocal(mockVault, Mongo)
		assert.NoError(t, err)
		assert.Equal(t, "localhost", cfg.Mongo.Host)
	})
	t.Run("mongo with values", func(t *testing.T) {
		os.Clearenv()
		cfg, err := BuildLocal(mockVault, Mongo)
		assert.NoError(t, err)
		assert.Equal(t, "testUser", cfg.Mongo.Username)
	})
}

// Assuming ProjectConfigurator interface and Config structure are defined as shown previously

// MockProjectConfigurator is a mock implementation for testing purposes.
type MockProjectConfigurator struct{}

// Build simulates applying project-specific configurations.
func (mpc *MockProjectConfigurator) Build(opts ...BuildOption) error {
	_ = fmt.Sprintf("%v", opts)

	if err := os.Setenv("PROJECT_SPECIFIC_CONFIG", "true"); err != nil {
		return err
	}
	return nil
}

func TestProjectConfig(t *testing.T) {
	t.Run("project configuration", func(t *testing.T) {
		os.Clearenv()

		// Assuming this mock sets an environment variable as part of its configuration logic
		mockProjectConfigurator := MockProjectConfigurator{}
		mockLocalProject := WithProjectConfigurator(&mockProjectConfigurator)

		// Build configuration including the mock project configurator
		_, err := Build(mockLocalProject)
		assert.NoError(t, err)

		// Verify the project-specific configuration was recognized
		projectSpecificConfig, exists := os.LookupEnv("PROJECT_SPECIFIC_CONFIG")
		assert.True(t, exists)
		assert.Equal(t, "true", projectSpecificConfig)
	})
}
