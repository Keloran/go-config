package rabbit

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuild(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv()

		l, err := Build()
		assert.NoError(t, err)
		assert.Equal(t, "", l.Host)
		assert.Equal(t, "", l.Port)
		assert.Equal(t, "", l.Username)
		assert.Equal(t, "", l.Password)
		assert.Equal(t, "", l.VHost)
		assert.Equal(t, "", l.ManagementHost)
	})

	t.Run("with values", func(t *testing.T) {
		os.Clearenv()
		if err := os.Setenv("RABBIT_HOSTNAME", "http://localhost"); err != nil {
			assert.NoError(t, err)
		}
		r, err := Build()
		assert.NoError(t, err)
		assert.Equal(t, "http://localhost", r.Host)
	})
}
