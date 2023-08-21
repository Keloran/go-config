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
    assert.Equal(t, 3000, l.GRPCPort)
	})

	t.Run("custom values", func(t *testing.T) {
		os.Clearenv()
		_ = os.Setenv("BUGFIXES_LOCAL_ONLY", "true")
		_ = os.Setenv("DEVELOPMENT", "true")
		_ = os.Setenv("HTTP_PORT", "8080")
    _ = os.Setenv("GRPC_PORT", "9999")

		l, err := Build()

		assert.NoError(t, err)
		assert.Equal(t, true, l.KeepLocal)
		assert.Equal(t, true, l.Development)
		assert.Equal(t, 8080, l.HTTPPort)
    assert.Equal(t, 9999, l.GRPCPort)
	})
}
