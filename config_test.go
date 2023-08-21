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
		// Check default values for Local
		assert.Equal(t, false, cfg.KeepLocal)
		assert.Equal(t, false, cfg.Development)
		assert.Equal(t, 80, cfg.HTTPPort)
	})

	t.Run("custom values", func(t *testing.T) {
		os.Clearenv()
		// Set custom environment variables for Local
		os.Setenv("BUGFIXES_LOCAL_ONLY", "true")
		os.Setenv("DEVELOPMENT", "true")
		os.Setenv("HTTP_PORT", "8080")

		cfg, err := Build(Local)

		assert.NoError(t, err)
		// Check custom values for Local
		assert.Equal(t, true, cfg.KeepLocal)
		assert.Equal(t, true, cfg.Development)
		assert.Equal(t, 8080, cfg.HTTPPort)
	})
}
