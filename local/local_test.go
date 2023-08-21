package local

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
		assert.Equal(t, false, l.KeepLocal)
		assert.Equal(t, false, l.Development)
		assert.Equal(t, 80, l.HTTPPort)
	})

	t.Run("custom values", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("BUGFIXES_LOCAL_ONLY", "true")
		os.Setenv("DEVELOPMENT", "true")
		os.Setenv("HTTP_PORT", "8080")

		l, err := Build()

		assert.NoError(t, err)
		assert.Equal(t, true, l.KeepLocal)
		assert.Equal(t, true, l.Development)
		assert.Equal(t, 8080, l.HTTPPort)
	})
}
