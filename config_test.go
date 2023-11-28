package ConfigBuilder

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	t.Run("rabbit", func(t *testing.T) {
		os.Clearenv()
		cfg, err := Build(Rabbit)
		assert.NoError(t, err)
		assert.Equal(t, "", cfg.Rabbit.Host)
	})
}

//func TestDatabase(t *testing.T) {
//	t.Run("database", func(t *testing.T) {
//		os.Clearenv()
//		cfg, err := Build(Database)
//		assert.NoError(t, err)
//		assert.Equal(t, "", cfg.Database.Host)
//	})
//}

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

//func TestMongo(t *testing.T) {
//	t.Run("mongo", func(t *testing.T) {
//		os.Clearenv()
//		cfg, err := Build(Mongo)
//		assert.NoError(t, err)
//		assert.Equal(t, "", cfg.Mongo.Host)
//	})
//}
