package vault

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv() // Clear all environment variables

		l, _, err := Build()
		assert.NoError(t, err)
		assert.Equal(t, "vault.vault", l.Host)
		assert.Equal(t, "", l.Port)
		assert.Equal(t, "root", l.Token)
		assert.Equal(t, "https://vault.vault", l.Address)
	})

	t.Run("with port", func(t *testing.T) {
		os.Clearenv()
		if err := os.Setenv("VAULT_PORT", "8080"); err != nil {
			assert.NoError(t, err)
		}

		v, _, err := Build()
		assert.NoError(t, err)
		assert.Equal(t, "vault.vault:8080", v.Address)
	})

	t.Run("with http prefix", func(t *testing.T) {
		os.Clearenv()
		if err := os.Setenv("VAULT_HOST", "http://localhost"); err != nil {
			assert.NoError(t, err)
		}

		v, _, err := Build()
		assert.NoError(t, err)
		assert.Equal(t, "http://localhost", v.Address)
	})
}
