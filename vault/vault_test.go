package vault

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv() // Clear all environment variables

		l, err := Build()

		assert.NoError(t, err)
		assert.Equal(t, "localhost", l.Host)
		assert.Equal(t, "", l.Port)
		assert.Equal(t, "root", l.Token)
		assert.Equal(t, "https://localhost", l.Address)
	})

	t.Run("with port", func(t *testing.T) {
		os.Clearenv()
		_ = os.Setenv("VAULT_PORT", "8080")

		v, err := Build()

		assert.NoError(t, err)
		assert.Equal(t, "localhost:8080", v.Address)
	})

	t.Run("with http prefix", func(t *testing.T) {
		os.Clearenv()
		_ = os.Setenv("VAULT_HOST", "http://localhost")

		v, err := Build()

		assert.NoError(t, err)
		assert.Equal(t, "http://localhost", v.Address)
	})
}
