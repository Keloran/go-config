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

func TestConfig(t *testing.T) {
	t.Run("test config mock vault", func(t *testing.T) {
		mockVault := &MockVaultHelper{
			KVSecrets: []vaulthelper.KVSecret{
				{Key: "password", Value: "testPassword"},
				{Key: "username", Value: "testUser"},
				{Key: "rds-hostname", Value: "testHost"},
				{Key: "rds-db", Value: "testDB"},
			},
		}

		cfg := NewConfig(mockVault)
		err := cfg.Build(Local)
		assert.NoError(t, err)
	})
	t.Run("test config real vault", func(t *testing.T) {
		vh := vaulthelper.NewVault("tester", "tester")
		cfg := NewConfig(vh)
		err := cfg.Build(Local)
		assert.NoError(t, err)
	})
}

func TestBuild(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv() // Clear all environment variables

		cfg, err := Build(Local)
		fmt.Printf("%v", cfg)
		assert.NoError(t, err)
		assert.Equal(t, false, cfg.Local.KeepLocal)
		assert.Equal(t, false, cfg.Local.Development)
		assert.Equal(t, 80, cfg.Local.HTTPPort)
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
		assert.Equal(t, true, cfg.Local.KeepLocal)
		assert.Equal(t, true, cfg.Local.Development)
		assert.Equal(t, 8080, cfg.Local.HTTPPort)
	})
}

func TestRabbit(t *testing.T) {
	t.Run("rabbit no set values", func(t *testing.T) {
		os.Clearenv()

		cfg, err := BuildLocal(Rabbit)
		assert.NoError(t, err)
		assert.Equal(t, "", cfg.Rabbit.Host)
	})
	t.Run("rabbit with values", func(t *testing.T) {
		os.Clearenv()
		mockVault := &MockVaultHelper{
			KVSecrets: []vaulthelper.KVSecret{
				{Key: "rabbit-password", Value: "testPassword"},
				{Key: "rabbit-username", Value: "testUser"},
				{Key: "rabbit-vhost", Value: "testVhost"},
				{Key: "rabbit-hostname", Value: ""},
				{Key: "rabbit-management-hostname", Value: ""},
				{Key: "rabbit-queue", Value: ""},
			},
		}

		cfg, err := BuildLocalVH(mockVault, Rabbit)
		assert.NoError(t, err)
		assert.Equal(t, "testUser", cfg.Rabbit.Username)
	})
}

func TestDatabase(t *testing.T) {
	mockVault := &MockVaultHelper{
		KVSecrets: []vaulthelper.KVSecret{
			{Key: "password", Value: "testPassword"},
			{Key: "username", Value: "testUser"},
			{Key: "rds-hostname", Value: "testHost"},
			{Key: "rds-db", Value: "testDB"},
		},
	}

	t.Run("database no set values", func(t *testing.T) {
		os.Clearenv()
		cfg, err := BuildLocal(Database)
		assert.NoError(t, err)
		assert.Equal(t, "postgres.chewedfeed", cfg.Database.Host)
	})
	t.Run("database with values", func(t *testing.T) {
		os.Clearenv()
		cfg, err := BuildLocalVH(mockVault, Database)
		assert.NoError(t, err)
		assert.Equal(t, "testUser", cfg.Database.User)
	})
}

func TestKeycloak(t *testing.T) {
	t.Run("keycloak", func(t *testing.T) {
		os.Clearenv()
		cfg, err := Build(Keycloak)
		assert.NoError(t, err)
		assert.Equal(t, "https://keys.chewedfeed.com", cfg.Keycloak.Host)
	})
}

func TestVault(t *testing.T) {
	t.Run("vault", func(t *testing.T) {
		os.Clearenv()
		cfg, err := Build(Vault)
		assert.NoError(t, err)
		assert.Equal(t, "vault.vault", cfg.Vault.Host)
	})
}

func TestMongo(t *testing.T) {
	mockVault := &MockVaultHelper{
		KVSecrets: []vaulthelper.KVSecret{
			{Key: "password", Value: "testPassword"},
			{Key: "username", Value: "testUser"},
			{Key: "mongo-hostname", Value: "testHost"},
			{Key: "mongo-collections", Value: "tester:testerCollection"},
			{Key: "mongo-db", Value: "testDB"},
		},
	}

	t.Run("mongo no set values", func(t *testing.T) {
		os.Clearenv()
		cfg, err := BuildLocal(Mongo)
		assert.NoError(t, err)
		assert.Equal(t, "localhost", cfg.Mongo.Host)
	})
	t.Run("mongo with values", func(t *testing.T) {
		os.Clearenv()
		cfg, err := BuildLocalVH(mockVault, Mongo)
		assert.NoError(t, err)
		assert.Equal(t, "testUser", cfg.Mongo.Username)
	})
}

// Assuming ProjectConfigurator interface and Config structure are defined as shown previously

// MockProjectConfigurator is a mock implementation for testing purposes.
type MockProjectConfigurator struct{}

// Build simulates applying project-specific configurations.
func (mpc MockProjectConfigurator) Build(c *Config) error {
	if c.ProjectProperties == nil {
		c.ProjectProperties = make(map[string]interface{})
	}
	c.ProjectProperties["TestProperty"] = true
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
	})
}

func TestProjectBuild(t *testing.T) {
	t.Run("project build configuration", func(t *testing.T) {
		os.Clearenv()

		cfgNoProps, _ := Build(Local)
		assert.Equal(t, map[string]interface{}(map[string]interface{}(nil)), cfgNoProps.ProjectProperties)

		cfg, err := Build(Local, WithProjectConfigurator(MockProjectConfigurator{}))
		assert.NoError(t, err)

		assert.Equal(t, true, cfg.ProjectProperties["TestProperty"])
	})
}
