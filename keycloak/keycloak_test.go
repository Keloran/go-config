package keycloak

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv()

		l, err := Build()
		assert.NoError(t, err)
		assert.Equal(t, "", l.Host)
		assert.Equal(t, "", l.Client)
		assert.Equal(t, "", l.Secret)
		assert.Equal(t, "", l.Realm)
	})

	t.Run("with values", func(t *testing.T) {
		os.Clearenv()

		if err := os.Setenv("KEYCLOAK_HOSTNAME", "http://localhost"); err != nil {
			assert.NoError(t, err)
		}
		l, err := Build()
		assert.NoError(t, err)
		assert.Equal(t, "http://localhost", l.Host)
	})
}
