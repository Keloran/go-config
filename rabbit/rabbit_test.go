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
	})

	t.Run("with values", func(t *testing.T) {
		os.Clearenv()
		_ = os.Setenv("RABBIT_HOSTNAME", "http://localhost")
		r, err := Build()
		assert.NoError(t, err)
		assert.Equal(t, "http://localhost", r.Host)
	})
}