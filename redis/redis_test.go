package redis

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()
	r, err := NewSystem().Build()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:6379", r.Host)
	assert.Equal(t, 0, r.DB)

	t.Setenv("REDIS_HOSTNAME", "testHost:6379")
	t.Setenv("REDIS_PASSWORD", "testPassword")
	t.Setenv("REDIS_DB", "2")
	r, err = NewSystem().Build()
	assert.NoError(t, err)
	assert.Equal(t, "testHost:6379", r.Host)
	assert.Equal(t, "testPassword", r.Password)
	assert.Equal(t, 2, r.DB)
}
